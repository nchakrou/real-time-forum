package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"

	"forum/backend"
)

func GetMyPostsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			backend.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		user, err := backend.GetUserIDFromRequest(db, r)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		userId := user.ID

		rows, err := db.Query(`
			SELECT 
				p.id, p.title, p.content, p.likes, p.dislikes, p.comments,
				GROUP_CONCAT(c.name) AS categories,
				COALESCE(l.value, 0) as userValue
			FROM posts p
			JOIN post_categories pc ON pc.post_id = p.id
			JOIN categories c ON c.id = pc.category_id
			LEFT JOIN likes l ON l.post_id = p.id AND l.user_id = ?
			WHERE p.user_id = ?
			GROUP BY p.id
			ORDER BY p.created_at DESC 
		`, userId, userId)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var posts = struct {
			Posts []Post
		}{
			Posts: []Post{},
		}

		for rows.Next() {
			var post Post
			var category string

			err := rows.Scan(&post.Id, &post.Title, &post.Content, &post.Likes, &post.Dislikes, &post.Comments, &category, &post.UserValue)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			post.Categories = strings.Split(category, ",")
			posts.Posts = append(posts.Posts, post)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(posts)
	}
}
