package handlers

import (
	"database/sql"
	"encoding/json"
	"forum/backend"
	"log"
	"net/http"
	"strings"
)

func GetMyPostsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		userId, err := backend.GetUserIDFromRequest(db, r)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		rows, err := db.Query(`
		SELECT 
		p.id, p.title, p.content, p.likes, p.dislikes, p.comments,
		 GROUP_CONCAT(c.name) AS categories
	FROM posts p
	 JOIN post_categories pc ON pc.post_id = p.id
	 JOIN categories c ON c.id = pc.category_id
	WHERE p.user_id = ?
	GROUP BY p.id
	ORDER BY p.created_at DESC
		`, userId)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println("Error querying posts:", err)
			return
		}
		defer rows.Close()
		var posts []Post
		for rows.Next() {
			var post Post
			var category string
			if err := rows.Scan(&post.Id, &post.Title, &post.Content, &post.Likes, &post.Dislikes, &post.Comments, &category); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Println("Error scanning post:", err)
				return
			}

			post.Categories = strings.Split(category, ",")
			posts = append(posts, post)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(posts)
	}
}
