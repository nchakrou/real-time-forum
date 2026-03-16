package Websocket

import (
	"database/sql"
	"fmt"
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
func (h *Hub) GetChatUsers(db *sql.DB, senderID int, conn *websocket.Conn) {
	query := `SELECT 
    u.nickname,
    m.created_at
FROM messages m
JOIN (
    SELECT MAX(id) AS last_id
    FROM messages
    WHERE (sender_id = ? OR receiver_id = ?) AND sender_id != receiver_id
    GROUP BY
        CASE
            WHEN sender_id = ? THEN receiver_id
            ELSE sender_id
        END
) t
ON m.id = t.last_id
JOIN users u
ON u.id = CASE
            WHEN m.sender_id = ? THEN m.receiver_id
            ELSE m.sender_id
           END
ORDER BY m.created_at DESC;`

	var res response
	res.Type = "chat_users"
	h.mu.Lock()
	defer h.mu.Unlock()
	rows, err := db.Query(query, senderID, senderID, senderID, senderID)
	var chat []Chat
	if err != nil {
		log.Println("Error getting chat users:", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var temp Chat
		if err := rows.Scan(&temp.Target, &temp.CreatedAt); err != nil {
			log.Println("Error scanning chat users:", err)
			continue
		}
		chat = append(chat, temp)
	}
	res.Chat = chat
	fmt.Println(res)
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
