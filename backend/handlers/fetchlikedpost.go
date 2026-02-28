package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"forum/backend"
)

type Pot struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	UserID    int       `json:"user_id"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
	Likes     int       `json:"likes"`
	Dislikes  int       `json:"dislikes"`
}

func HandleLikedPosts(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// 🔹 Check Method
		if r.Method != http.MethodGet {
			log.Println("❌ Method not allowed:", r.Method)
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// 🔹 Check Auth
		user, err := backend.GetUserIDFromRequest(db, r)
		if err != nil {
			log.Println("❌ Unauthorized error:", err)
			http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
			return
		}

		log.Println("✅ User ID:", user.ID)

		query := `
		SELECT 
			p.id,
			p.title,
			p.content,
			p.user_id,
			u.nickname,
			p.created_at,
			p.likes,
			p.dislikes
		FROM posts p
		INNER JOIN likes l ON l.post_id = p.id
		LEFT JOIN users u ON u.id = p.user_id
		WHERE l.user_id = ? AND l.value = 1
		ORDER BY p.created_at DESC;
		`

		// 🔹 Execute Query
		rows, err := db.Query(query, user.ID)
		if err != nil {
			log.Println("❌ Query error:", err)
			http.Error(w, "Database query error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var posts []Pot

		// 🔹 Loop rows
		for rows.Next() {
			var p Pot

			err := rows.Scan(
				&p.ID,
				&p.Title,
				&p.Content,
				&p.UserID,
				&p.Username,
				&p.CreatedAt,
				&p.Likes,
				&p.Dislikes,
			)

			if err != nil {
				log.Println("❌ Scan error:", err)
				http.Error(w, "Scan error: "+err.Error(), http.StatusInternalServerError)
				return
			}

			posts = append(posts, p)
		}

		// 🔹 Check iteration error
		if err = rows.Err(); err != nil {
			log.Println("❌ Rows iteration error:", err)
			http.Error(w, "Rows error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// 🔹 Debug number of posts
		log.Println("✅ Total liked posts:", len(posts))

		// 🔹 Set Header
		w.Header().Set("Content-Type", "application/json")

		// 🔹 Encode Response
		if err := json.NewEncoder(w).Encode(posts); err != nil {
			log.Println("❌ JSON Encode error:", err)
			http.Error(w, "JSON encode error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		log.Println("✅ Response sent successfully")
	}
}
