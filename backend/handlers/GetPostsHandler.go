package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

type Post struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

func GetPostsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		rows, err := db.Query(`
		SELECT title, content FROM posts
		ORDER BY created_at DESC
		LIMIT 10
		`)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println("Error querying posts:", err)
			return
		}
		defer rows.Close()

		var posts []Post
		for rows.Next() {
			var post Post
			if err := rows.Scan(&post.Title, &post.Content); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Println("Error scanning post:", err)
				return
			}
			posts = append(posts, post)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(posts)
	}
}
