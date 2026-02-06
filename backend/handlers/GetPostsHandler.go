package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Post struct {
	Id         int      `json:"id"`
	Title      string   `json:"title"`
	Content    string   `json:"content"`
	Likes      int      `json:"likes"`
	Dislikes   int      `json:"dislikes"`
	Comments   int      `json:"comments"`
	Categories []string `json:"categories"`
}

func GetPostsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		rows, err := db.Query(`
		SELECT 
		p.id, p.title, p.content, p.likes, p.dislikes, p.comments,
		 c.name
	FROM posts p
	 JOIN post_categories pc ON pc.post_id = p.id
	 JOIN categories c ON c.id = pc.category_id
	 GROUP BY p.id 
	ORDER BY p.created_at DESC
	LIMIT 10
		`)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println("Error querying posts:", err)
			return
		}
		defer rows.Close()

		var posts = make(map[int]Post)
		for rows.Next() {
			var post Post
			var category string
			if err := rows.Scan(&post.Id, &post.Title, &post.Content, &post.Likes, &post.Dislikes, &post.Comments, &category); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Println("Error scanning post:", err)
				return
			}

			if existingPost, exists := posts[post.Id]; exists {
				existingPost.Categories = append(existingPost.Categories, category)
				posts[post.Id] = existingPost
			} else {
				post.Categories = []string{category}
				posts[post.Id] = post
			}
		}
		fmt.Println(posts)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(posts)

	}
}
