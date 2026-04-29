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
	To        string `json:"to,omitempty"`
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

func (hub *Hub) SendPrivateMessage(myconn *websocket.Conn, db *sql.DB, fromID int, toUsername string, message string, fromUsername string) {
	toID, err := TargetID(db, toUsername)
	if err != nil {
		myconn.WriteJSON(struct {
			Type    string `json:"type"`
			Message string `json:"message"`
			Code    int    `json:"code"`
		}{
			Type:    "error",
			Code:    404,
			Message: "user not found",
		})
		log.Println("Error getting target ID:", err)
		return
	}

	tx, err := db.Begin()
	if err != nil {
		myconn.WriteJSON(struct {
			Type    string `json:"type"`
			Message string `json:"message"`
			Code    int    `json:"code"`
		}{
			Type:    "error",
			Code:    500,
			Message: "internal server error",
		})
		log.Println("Error starting transaction:", err)
		return
	}
	defer tx.Rollback()
	if fromID == toID {
		myconn.WriteJSON(struct {
			Type    string `json:"type"`
			Message string `json:"message"`
			Code    int    `json:"code"`
		}{
			Type:    "error",
			Code:    400,
			Message: "You cannot send a message to yourself.",
		})
		log.Println("User cannot send message to themselves")
		return
	}

	if len(message) == 0 || len(message) > 1000 {
		myconn.WriteJSON(struct {
			Type    string `json:"type"`
			Message string `json:"message"`
			Code    int    `json:"code"`
		}{
			Type:    "error",
			Code:    400,
			Message: "Message must be between 1 and 1000 characters.",
		})
		log.Println("Invalid message length:", len(message))
		return
	}

	_, err = tx.Exec(
		"INSERT INTO messages (sender_id, receiver_id, content) VALUES (?, ?, ?)",
		fromID, toID, message,
	)
	if err != nil {
		myconn.WriteJSON(struct {
			Type    string `json:"type"`
			Message string `json:"message"`
			Code    int    `json:"code"`
		}{
			Type:    "error",
			Code:    500,
			Message: "internal server error",
		})
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
		myconn.WriteJSON(struct {
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

	if err := tx.Commit(); err != nil {
		log.Println("Error committing transaction:", err)
		myconn.WriteJSON(struct {
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

	now := time.Now().Format(time.RFC3339)

	hub.mu.Lock()
	toconns := append([]*websocket.Conn(nil), hub.Clients[toID]...)
	fromconns := append([]*websocket.Conn(nil), hub.Clients[fromID]...)
	hub.mu.Unlock()
	for i, conn := range toconns {
		resp := Msg{
			Type:      "private_message",
			From:      fromUsername,
			Message:   message,
			CreatedAt: now,
		}
		if err := conn.WriteJSON(resp); err != nil {

			log.Println("WS write error:", err)
			hub.mu.Lock()
			clients := hub.Clients[toID]
			if i < len(clients) {
				hub.Clients[toID] = append(clients[:i], clients[i+1:]...)
			}
			hub.mu.Unlock()
		}

	}
	for i, conn := range fromconns {
		if conn == myconn {
			continue
		}
		resp := Msg{
			Type:      "private_message",
			From:      fromUsername,
			To:        toUsername,
			Message:   message,
			CreatedAt: now,
		}
		if err := conn.WriteJSON(resp); err != nil {
			log.Println("WS write error:", err)

			hub.mu.Lock()
			clients := hub.Clients[fromID]
			if i < len(clients) {
				hub.Clients[fromID] = append(clients[:i], clients[i+1:]...)
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

func (hub *Hub) SendTypingStatus(db *sql.DB, fromID int, toUsername, fromUsername string, isTyping bool) {
	toID, err := TargetID(db, toUsername)
	if err != nil || fromID == toID {
		return
	}
	hub.sendTypingToUserID(toID, fromUsername, isTyping)
}

func (hub *Hub) sendTypingToUserID(toID int, fromUsername string, isTyping bool) {
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

func (hub *Hub) HandleTypingEvent(db *sql.DB, conn *websocket.Conn, fromID int, toUsername, fromUsername string, isTyping bool) {
	toID, err := TargetID(db, toUsername)
	if err != nil || fromID == toID {
		return
	}

	shouldBroadcast := false
	broadcastValue := false

	hub.mu.Lock()

	if hub.TypingByConn[conn] == nil {
		hub.TypingByConn[conn] = make(map[int]bool)
	}
	if hub.TypingCount[fromID] == nil {
		hub.TypingCount[fromID] = make(map[int]int)
	}

	connTargets := hub.TypingByConn[conn]
	fromCounts := hub.TypingCount[fromID]

	if isTyping {
		if !connTargets[toID] {
			connTargets[toID] = true
			fromCounts[toID]++
			if fromCounts[toID] == 1 {
				shouldBroadcast = true
				broadcastValue = true
			}
		}
	} else {
		if connTargets[toID] {
			delete(connTargets, toID)

			if fromCounts[toID] > 0 {
				fromCounts[toID]--
			}
			if fromCounts[toID] <= 0 {
				delete(fromCounts, toID)
				shouldBroadcast = true
				broadcastValue = false
			}

			if len(connTargets) == 0 {
				delete(hub.TypingByConn, conn)
			}
			if len(fromCounts) == 0 {
				delete(hub.TypingCount, fromID)
			}
		}
	}

	hub.mu.Unlock()

	if shouldBroadcast {
		hub.sendTypingToUserID(toID, fromUsername, broadcastValue)
	}
}

func (hub *Hub) ClearTypingForConn(db *sql.DB, conn *websocket.Conn, fromID int, fromUsername string) {
	hub.mu.Lock()
	targetsSet := hub.TypingByConn[conn]
	targetIDs := make([]int, 0, len(targetsSet))
	for toID := range targetsSet {
		targetIDs = append(targetIDs, toID)
	}
	hub.mu.Unlock()

	for _, toID := range targetIDs {
		shouldBroadcast := false

		hub.mu.Lock()
		connTargets := hub.TypingByConn[conn]
		fromCounts := hub.TypingCount[fromID]

		if connTargets != nil && connTargets[toID] {
			delete(connTargets, toID)

			if fromCounts != nil {
				if fromCounts[toID] > 0 {
					fromCounts[toID]--
				}
				if fromCounts[toID] <= 0 {
					delete(fromCounts, toID)
					shouldBroadcast = true
				}
				if len(fromCounts) == 0 {
					delete(hub.TypingCount, fromID)
				}
			}
			if len(connTargets) == 0 {
				delete(hub.TypingByConn, conn)
			}
		}
		hub.mu.Unlock()

		if shouldBroadcast {
			hub.sendTypingToUserID(toID, fromUsername, false)
		}
	}
}
