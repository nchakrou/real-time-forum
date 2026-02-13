package Websocket

import (
	"database/sql"
	"log"

	"github.com/gorilla/websocket"
)

func (h *Hub) OnlineUsers(db *sql.DB, senderID int, conn *websocket.Conn) {
	query := `SELECT nickname FROM users WHERE id =?`
	var res response
	res.Type = "online_users"
	h.mu.Lock()
	defer h.mu.Unlock()

	for userid := range h.Clients {
		var username string
		if err := db.QueryRow(query, userid).Scan(&username); err != nil {
			log.Println("Error getting username:", err)
			continue
		}
		res.Users = append(res.Users, username)
	}
	if err := conn.WriteJSON(res); err != nil {
		log.Println("Error writing JSON:", err)
	}
}
