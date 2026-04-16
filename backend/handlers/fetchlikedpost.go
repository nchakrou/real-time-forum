package handlers

import (
	"database/sql"
	"encoding/json"
	"forum/backend"
	"net/http"
	"strings"
	"time"
)

type LikedPost struct {
	ID         int       `json:"id"`
	Title      string    `json:"title"`
	Content    string    `json:"content"`
	UserID     int       `json:"user_id"`
	Username   string    `json:"username"`
	CreatedAt  time.Time `json:"created_at"`
	Likes      int       `json:"likes"`
	Dislikes   int       `json:"dislikes"`
	Comments   int       `json:"comments"`
	UserValue  int       `json:"userValue"`
	Categories []string  `json:"categories"`
}

func HandleLikedPosts(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			backend.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		user, err := backend.GetUserIDFromRequest(db, r)
		if err != nil {
			backend.WriteJSONError(w, http.StatusUnauthorized, "login required")
			return
		}

		query := `
            SELECT p.id, p.title, p.content, p.user_id, u.nickname, p.created_at, p.likes, p.dislikes, p.comments, l.value as userValue,
                   GROUP_CONCAT(c.name) AS categories
            FROM posts p
            INNER JOIN likes l ON l.post_id = p.id
            LEFT JOIN users u ON u.id = p.user_id
            JOIN post_categories pc ON pc.post_id = p.id
            JOIN categories c ON c.id = pc.category_id
            WHERE l.user_id = ? AND l.value = 1
            GROUP BY p.id
            ORDER BY p.created_at DESC`

		rows, err := db.Query(query, user.ID)
		if err != nil {
			backend.WriteJSONError(w, http.StatusInternalServerError, "something went wrong")
			return
		}
		defer rows.Close()

		posts := make([]LikedPost, 0)

		for rows.Next() {
			var p LikedPost
			var categories string
			if err := rows.Scan(
				&p.ID, &p.Title, &p.Content, &p.UserID,
				&p.Username, &p.CreatedAt, &p.Likes, &p.Dislikes, &p.Comments, &p.UserValue, &categories,
			); err != nil {
				backend.WriteJSONError(w, http.StatusInternalServerError, "something went wrong")
				return
			}
			p.Categories = strings.Split(categories, ",")
			posts = append(posts, p)
		}

		if err = rows.Err(); err != nil {
			backend.WriteJSONError(w, http.StatusInternalServerError, "something went wrong")
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(posts)
	}
}
