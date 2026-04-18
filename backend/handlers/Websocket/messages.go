package Websocket

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type Msg struct {
	ID        int    `json:"id"`
	Type      string `json:"type"`
	From      string `json:"from"`
	Message   string `json:"message,omitempty"`
	Content   string `json:"content,omitempty"`
	CreatedAt string `json:"CreatedAt"`
}

type resChat struct {
	Type     string `json:"type"`
	Messages []Msg  `json:"Messages"`
}

type TypingResponse struct {
	Type     string `json:"type"`
	From     string `json:"from"`
	IsTyping bool   `json:"is_typing"`
}

func (hub *Hub) SendPrivateMessage(db *sql.DB, fromID int, toUsername string, message string, fromUsername string) {

	toID, err := TargetID(db, toUsername)
	if err != nil {
		log.Println("Error getting target ID:", err)
		return
	}

	tx, err := db.Begin()
	if err != nil {
		log.Println("Error starting transaction:", err)
		return
	}
	defer tx.Rollback()

	_, err = tx.Exec(
		"INSERT INTO messages (sender_id, receiver_id, content) VALUES (?, ?, ?)",
		fromID, toID, message,
	)
	if err != nil {
		log.Println("Error saving message:", err)
		return
	}

	if fromID != toID {
		_, err = tx.Exec(`
		INSERT INTO notifications (sender_id, receiver_id, type, content)
		VALUES (?, ?, ?, ?)`,
			fromID, toID, "message", message,
		)
	}
	if err != nil {
		log.Println("Error saving notification:", err)
		return
	}

	if err := tx.Commit(); err != nil {
		log.Println("Error committing transaction:", err)
		return
	}

	now := time.Now().Format(time.RFC3339)

	hub.mu.Lock()
	conns := make([]*websocket.Conn, len(hub.Clients[toID]))
	copy(conns, hub.Clients[toID])
	hub.mu.Unlock()

	for i, conn := range conns {
		resp := Msg{
			Type:      "private_message",
			From:      fromUsername,
			Message:   message,
			CreatedAt: now,
		}

		if err := conn.WriteJSON(resp); err != nil {
			log.Println("Error sending message:", err)
			conn.Close()

			hub.mu.Lock()
			clients := hub.Clients[toID]
			if i < len(clients) {
				hub.Clients[toID] = append(clients[:i], clients[i+1:]...)
			}
			hub.mu.Unlock()
			continue
		}
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

func (hub *Hub) GetMessages(w http.ResponseWriter, db *sql.DB, fromID int, toUsername string, conn *websocket.Conn, fromUsername string, lastID int) {

	toID, err := TargetID(db, toUsername)
	if err != nil {
		log.Println("Error getting target ID:", err)
		conn.WriteJSON(struct {
			Type    string `json:"type"`
			Message string `json:"message"`
			Code    int    `json:"code"`
		}{
			Type:    "error",
			Code:    404,
			Message: "user not found",
		})

		return
	}

	_, err = db.Exec(`
UPDATE notifications
SET is_read = 1
WHERE sender_id = ? AND receiver_id = ? AND type = 'message'
`, toID, fromID)

	if err != nil {
		log.Println("Error updating notifications:", err)
		conn.WriteJSON(struct {
			Type    string `json:"type"`
			Message string `json:"message"`
			Code    int    `json:"code"`
		}{
			Type:    "error",
			Code:    500,
			Message: "internal server error",
		})
		return
	}

	query := `
	select m.id, u.nickname, m.content, m.created_at 
	from messages m 
	join users u on m.sender_id = u.id 
	where ((m.sender_id = ? AND m.receiver_id = ?) 
	OR (m.sender_id = ? AND m.receiver_id = ?))`

	args := []interface{}{fromID, toID, toID, fromID}

	if lastID > 0 {
		query += " AND m.id < ?"
		args = append(args, lastID)
	}

	query += " ORDER BY m.id DESC LIMIT 10"

	rows, err := db.Query(query, args...)

	if err != nil {
		log.Println("Error getting messages:", err)
		conn.WriteJSON(struct {
			Type    string `json:"type"`
			Message string `json:"message"`
			Code    int    `json:"code"`
		}{
			Type:    "error",
			Code:    500,
			Message: "internal server error",
		})
		return
	}
	defer rows.Close()

	var Msgs []Msg

	for rows.Next() {
		var msg Msg
		if err := rows.Scan(&msg.ID, &msg.From, &msg.Message, &msg.CreatedAt); err != nil {
			log.Println("Error scanning message:", err)
			continue
		}

		msg.Type = "message"

		Msgs = append(Msgs, msg)
	}

	if err := conn.WriteJSON(resChat{
		Type:     "chat_history",
		Messages: Msgs,
	}); err != nil {
		log.Println("Error sending messages:", err)
		return
	}
}

func (hub *Hub) SendTypingStatus(db *sql.DB, fromID int, toUsername string, fromUsername string, isTyping bool) {
	toID, err := TargetID(db, toUsername)
	if err != nil {
		log.Println("Error getting target ID:", err)
		return
	}
	if fromID == toID {
		return
	}

	resp := TypingResponse{
		Type:     "typing",
		From:     fromUsername,
		IsTyping: isTyping,
	}

	hub.mu.Lock()
	conns := make([]*websocket.Conn, len(hub.Clients[toID]))
	copy(conns, hub.Clients[toID])
	hub.mu.Unlock()

	for _, conn := range conns {
		if err := conn.WriteJSON(resp); err != nil {
			log.Println("Error sending typing status:", err)
		}
	}
}
