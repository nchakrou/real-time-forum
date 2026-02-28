package handlers

import (
	"database/sql"
	"encoding/json"
	"forum/backend"
	"net/http"
	"time"
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

		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		user, err := backend.GetUserIDFromRequest(db, r)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

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

		rows, err := db.Query(query, user.ID)
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var posts []Pot

		for rows.Next() {
			var p Pot
			if err := rows.Scan(
				&p.ID,
				&p.Title,
				&p.Content,
				&p.UserID,
				&p.Username,
				&p.CreatedAt,
				&p.Likes,
				&p.Dislikes,
			); err != nil {
				http.Error(w, "Scan error", http.StatusInternalServerError)
				return
			}
			posts = append(posts, p)
		}

		if err = rows.Err(); err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(posts)
	}
}
