package Websocket

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"forum/backend"

	"github.com/gorilla/websocket"
)

type Hub struct {
	Clients      map[int][]*websocket.Conn
	mu           sync.Mutex
	UserSessions map[int]string
	TypingByConn map[*websocket.Conn]map[int]bool
	TypingCount  map[int]map[int]int
}
type request struct {
	Type      string `json:"type"`
	Target    string `json:"target,omitempty"`
	Message   string `json:"message,omitempty"`
	CreatedAt int64  `json:"timestamp,omitempty"`
	LastID    int    `json:"lastID,omitempty"`
}
type response struct {
	Type      string   `json:"type"`
	Users     []string `json:"users,omitempty"`
	From      string   `json:"from,omitempty"`
	Message   string   `json:"message,omitempty"`
	CreatedAt string   `json:"timestamp,omitempty"`
	Chat      []Chat   `json:"chat,omitempty"`
}
type LastMessage struct {
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
}

type Chat struct {
	Target      string       `json:"target"`
	CreatedAt   string       `json:"created_at,omitempty"`
	IsOnline    bool         `json:"is_online"`
	LastMessage *LastMessage `json:"last_message,omitempty"`
}

func WsHandler(db *sql.DB, hub *Hub) http.HandlerFunc {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	hub.Clients = make(map[int][]*websocket.Conn)
	hub.UserSessions = make(map[int]string)
	if hub.TypingByConn == nil {
		hub.TypingByConn = make(map[*websocket.Conn]map[int]bool)
	}
	if hub.TypingCount == nil {
		hub.TypingCount = make(map[int]map[int]int)
	}
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := backend.GetUserIDFromRequest(db, r)
		if err != nil {
			log.Println("Error getting user ID from request:", err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		userid := user.ID
		c, err := r.Cookie("session_token")

		hub.mu.Lock()
		Session := hub.UserSessions[userid]
		clients := hub.Clients[userid]

		shouldClose := Session != "" && Session != c.Value

		if shouldClose {
			delete(hub.Clients, userid)
		}

		hub.mu.Unlock()

		if shouldClose {
			for _, conn := range clients {
				conn.Close()
			}
		}
		hub.mu.Lock()
		hub.UserSessions[userid] = c.Value
		hub.mu.Unlock()
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("Error upgrading to WebSocket:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer func() {
			conn.Close()
			hub.mu.Lock()
			isEmpty := len(hub.Clients[userid]) == 0
			hub.mu.Unlock()

			if isEmpty {
				hub.ClearTypingForConn(db, conn, userid, user.Nickname)
				hub.Leave(db, user)
			}
		}()
		hub.mu.Lock()
		hub.Clients[userid] = append(hub.Clients[userid], conn)
		hub.mu.Unlock()
		hub.Join(db, user, conn)
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("Error reading message:", err)
				break
			}
			var req request
			err = json.Unmarshal(message, &req)
			if err != nil {
				log.Println("Error unmarshaling message:", err)
				break
			}
			switch req.Type {
			case "online_users":
				hub.OnlineUsers(db, userid, conn)
			case "message":
				hub.SendPrivateMessage(conn, db, userid, req.Target, req.Message, user.Nickname)
			case "getChat":
				fmt.Println("getChat", req.Target)
				hub.GetMessages(w, db, userid, req.Target, conn, user.Nickname, req.LastID)
			case "get_notifications":
				hub.GetNotifications(db, userid, conn)
			case "get_chat_users":
				hub.GetChatUsers(db, userid, conn)
			case "typing":
				hub.HandleTypingEvent(db, conn, userid, req.Target, user.Nickname, true)
			case "stop_typing":
				hub.HandleTypingEvent(db, conn, userid, req.Target, user.Nickname, false)
			case "logout":
				hub.mu.Lock()
				conns := append([]*websocket.Conn(nil), hub.Clients[userid]...)
				hub.mu.Unlock()

				for _, c := range conns {
					c.Close()
				}
			}
		}
	}
}

func (h *Hub) Leave(db *sql.DB, user backend.User) {
	h.mu.Lock()

	var allConns []*websocket.Conn
	for _, conns := range h.Clients {
		allConns = append(allConns, conns...)
	}

	h.mu.Unlock()

	for _, conn := range allConns {
		go func(c *websocket.Conn) {
			err := c.WriteJSON(response{
				Type: "user_offline",
				From: user.Nickname,
			})
			if err != nil {
				log.Println("Error sending message:", err)
			}
		}(conn)
	}
}
