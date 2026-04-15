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
	Clients map[int][]*websocket.Conn
	mu      sync.Mutex
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
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := backend.GetUserIDFromRequest(db, r)
		userid := user.ID
		if err != nil {
			log.Println("Error getting user ID from request:", err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("Error upgrading to WebSocket:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer func() {
			conn.Close()
			hub.mu.Lock()
			conns := hub.Clients[userid]
			for i, c := range conns {
				if c == conn {
					hub.Clients[userid] = append(conns[:i], conns[i+1:]...)
					break
				}
			}
			hub.mu.Unlock()
			if len(hub.Clients[userid]) == 0 {
				delete(hub.Clients, userid)
				hub.Leave(db, user)
			}
			log.Printf("User %d disconnected", userid)
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
			fmt.Println(req)
			switch req.Type {
			case "online_users":
				hub.OnlineUsers(db, userid, conn)
			case "message":
				hub.SendPrivateMessage(db, userid, req.Target, req.Message, user.Nickname)
			case "getChat":
				hub.GetMessages(db, userid, req.Target, conn, user.Nickname, req.LastID)
			case "get_notifications":
				hub.GetNotifications(db, userid, conn)
			case "get_chat_users":
				hub.GetChatUsers(db, userid, conn)
			case "typing":
				hub.SendTypingStatus(db, userid, req.Target, user.Nickname, true)
			case "stop_typing":
				hub.SendTypingStatus(db, userid, req.Target, user.Nickname, false)
			}
		}
	}
}

func (h *Hub) Leave(db *sql.DB, user backend.User) {
	h.mu.Lock()
	for _, conns := range h.Clients {
		for _, conn := range conns {
			log.Printf("Sending offline message to user %d", user.ID)
			err := conn.WriteJSON(response{
				Type: "user_offline",
				From: user.Nickname,
			})
			if err != nil {
				log.Println("Error sending message to user:", err)
			}
		}
	}
	h.mu.Unlock()
}
