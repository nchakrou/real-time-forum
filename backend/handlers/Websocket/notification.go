package Websocket

import (
	"database/sql"
	"log"
)

type Notification struct {
	Type      string `json:"type"`
	From      string `json:"from"`
	Content   string `json:"content"`
	PostID    int    `json:"post_id,omitempty"`
	CreatedAt string `json:"timestamp"`
}

func (hub *Hub) SendNotification(db *sql.DB, fromID int, toID int, fromUsername string, content string, postID int) {

	hub.mu.Lock()
	defer hub.mu.Unlock()

	_, err := db.Exec(`
    INSERT INTO notifications (sender_id, receiver_id, type, content, post_id)
    VALUES (?, ?, ?, ?, ?)`,
		fromID,
		toID,
		"like",
		content,
		postID,
	)
	if err != nil {
		log.Println("Error notification1:", err)

	}

	for _, conn := range hub.Clients[toID] {
		resp := Notification{
			Type:      "notification",
			From:      fromUsername,
			Content:   content,
			PostID:    postID,
			CreatedAt: "",
		}
		if err := conn.WriteJSON(resp); err != nil {
			log.Println("Error sending notification:", err)
			conn.Close()
			continue
		}

	}

}
