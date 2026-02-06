package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strings"
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

		category := r.URL.Query().Get("category")
		var rows *sql.Rows
		var err error
		if category == "" {
			rows, err = db.Query(`
		SELECT 
		p.id, p.title, p.content, p.likes, p.dislikes, p.comments,
		GROUP_CONCAT(c.name) AS categories
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
		} else {
			rows, err = db.Query(`
			SELECT 
    p.id, p.title, p.content, p.likes, p.dislikes, p.comments,
    (SELECT GROUP_CONCAT(c2.name)
     FROM post_categories pc2
     JOIN categories c2 ON c2.id = pc2.category_id
     WHERE pc2.post_id = p.id
    ) AS categories
FROM posts p
WHERE p.id IN (
   SELECT pc3.post_id
   FROM post_categories pc3
   JOIN categories c3 ON c3.id = pc3.category_id
   WHERE c3.name = ?
)
ORDER BY p.created_at DESC
LIMIT 10

			`, category)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Println("Error querying posts with category:", err)
				return
			}
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
		if err := json.NewEncoder(w).Encode(posts); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println("Error encoding posts to JSON:", err)
			return
		}

	}
}
