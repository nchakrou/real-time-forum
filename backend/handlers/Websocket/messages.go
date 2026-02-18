package Websocket

import (
	"database/sql"
	"log"
)

func (hub *Hub) SendPrivateMessage(db *sql.DB, fromID int, toUsername string, message string, fromUsername string) {
	toID, err := TargetID(db, toUsername)
	if err != nil {
		return
	}
	hub.mu.Lock()
	defer hub.mu.Unlock()

	for _, conn := range hub.Clients[toID] {
		resp := response{
			Type:    "message",
			From:    fromUsername,
			Message: message,
		}
		if err := conn.WriteJSON(resp); err != nil {
			log.Println("Error sending message to client:", err)
			conn.Close()
			continue
		}
	}
	_, err = db.Exec("INSERT INTO messages (sender_id, receiver_id, content) VALUES (?, ?, ?)", fromID, toID, message)
	if err != nil {
		log.Println("Error saving message to database:", err)
		return
	}
}

func TargetID(db *sql.DB, username string) (int, error) {
	var id int
	err := db.QueryRow("SELECT id FROM users WHERE nickname = ?", username).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}
