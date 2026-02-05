package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Post struct {
	Title    string `json:"title"`
	Content  string `json:"content"`
	Likes    int    `json:"likes"`
	Dislikes int    `json:"dislikes"`
	Comments int    `json:"comments"`
}

func GetPostsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		rows, err := db.Query(`
		SELECT title, content, likes,dislikes,comments FROM posts
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
			if err := rows.Scan(&post.Title, &post.Content, &post.Likes, &post.Dislikes, &post.Comments); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Println("Error scanning post:", err)
				return
			}
			posts = append(posts, post)
		}
		fmt.Println(posts)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(posts)

	}
}
