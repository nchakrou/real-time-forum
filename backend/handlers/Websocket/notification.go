package Websocket

import (
	"database/sql"
	"log"

	"github.com/gorilla/websocket"
)

type Notification struct {
	Type      string `json:"type"`
	From      string `json:"from"`
	Message   string `json:"message"`
	CreatedAt string `json:"CreatedAt"`
}

func (hub *Hub) GetNotifications(db *sql.DB, userID int, conn *websocket.Conn) {

	rows, err := db.Query(`
	SELECT u.nickname, n.content, n.created_at
	FROM notifications n
	JOIN users u ON n.sender_id = u.id
	WHERE n.receiver_id = ? AND n.is_read = 0
	ORDER BY n.created_at DESC
	`, userID)

	if err != nil {
		log.Println("Error getting notifications:", err)
		return
	}
	defer rows.Close()

	var notifs []Notification

	for rows.Next() {

		var n Notification

		if err := rows.Scan(&n.From, &n.Message, &n.CreatedAt); err != nil {
			log.Println("Error scanning notification:", err)
			continue
		}

		n.Type = "notification"

		notifs = append(notifs, n)
	}

	if err := conn.WriteJSON(map[string]interface{}{
		"type": "notifications_history",
		"data": notifs,
	}); err != nil {
		log.Println("Error sending notifications:", err)
	}
}
