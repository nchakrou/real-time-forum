package Websocket

import (
	"database/sql"
	"forum/backend"
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
		if userid == senderID {
			continue
		}
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
func (h *Hub) Join(db *sql.DB, user backend.User, conn *websocket.Conn) {
	var res response
	res.Type = "join"
	res.From = user.Nickname
	h.mu.Lock()
	defer h.mu.Unlock()
	for userid, c := range h.Clients {
		if userid == user.ID {
			continue
		}
		for _, con := range c {

			con.WriteJSON(res)
		}
	}

}
func (h *Hub) GetChatUsers(db *sql.DB, senderID int, conn *websocket.Conn) {
	query := `
		SELECT 
			u.id, 
			u.nickname, 
			m.content, 
			m.created_at
		FROM users u
		LEFT JOIN messages m ON m.id = (
			SELECT id 
			FROM messages 
			WHERE (sender_id = ? AND receiver_id = u.id) 
			   OR (sender_id = u.id AND receiver_id = ?)
			ORDER BY created_at DESC 
			LIMIT 1
		)
		WHERE u.id != ?
		ORDER BY 
			CASE WHEN m.created_at IS NULL THEN 1 ELSE 0 END, 
			m.created_at DESC,
			u.nickname ASC
	`
	rows, err := db.Query(query, senderID, senderID, senderID)
	if err != nil {
		log.Println("Error querying chat users:", err)
		return
	}
	defer rows.Close()

	var chats []Chat

	h.mu.Lock()
	defer h.mu.Unlock()
	for rows.Next() {
		var userID int
		var nickname string
		var content sql.NullString
		var createdAt sql.NullString

		if err := rows.Scan(&userID, &nickname, &content, &createdAt); err != nil {
			log.Println("Error scanning user:", err)
			continue
		}

		// Check if user is online
		isOnline := false
		if conns, exists := h.Clients[userID]; exists && len(conns) > 0 {
			isOnline = true
		}

		var lastMsg *LastMessage
		if content.Valid && createdAt.Valid {
			lastMsg = &LastMessage{
				Content:   content.String,
				CreatedAt: createdAt.String,
			}
		}

		chats = append(chats, Chat{
			Target:      nickname,
			IsOnline:    isOnline,
			LastMessage: lastMsg,
		})
	}

	var res response
	res.Type = "chat_users"
	res.Chat = chats

	if err := conn.WriteJSON(res); err != nil {
		log.Println("Error writing final chat users JSON:", err)
	}
}
