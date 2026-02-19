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
}
type response struct {
	Type      string   `json:"type"`
	Users     []string `json:"users,omitempty"`
	From      string   `json:"from,omitempty"`
	Message   string   `json:"message,omitempty"`
	CreatedAt string   `json:"timestamp,omitempty"`
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
		defer conn.Close()
		hub.mu.Lock()
		hub.Clients[userid] = append(hub.Clients[userid], conn)
		hub.mu.Unlock()
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
				hub.GetMessages(db, userid, req.Target, conn, user.Nickname)
			}
		}
	}
}
