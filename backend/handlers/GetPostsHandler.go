package handlers

import (
	"database/sql"
	"encoding/json"
	"forum/backend"
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
	UserValue  int      `json:"userValue"`
}

func GetPostsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			backend.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		user, err := backend.GetUserIDFromRequest(db, r)
		userID := 0
		if err != nil {
			backend.WriteJSONError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}
		userID = user.ID

		category := r.URL.Query().Get("category")

		var rows *sql.Rows

		if category == "" {
			rows, err = db.Query(`
				SELECT 
					p.id, p.title, p.content, p.likes, p.dislikes, p.comments,
					GROUP_CONCAT(c.name) AS categories,
					COALESCE(l.value, 0) as userValue
				FROM posts p
				JOIN post_categories pc ON pc.post_id = p.id
				JOIN categories c ON c.id = pc.category_id
				LEFT JOIN likes l ON l.post_id = p.id AND l.user_id = ?
				GROUP BY p.id 
				ORDER BY p.created_at DESC
			`, userID)

		} else {
			rows, err = db.Query(`
				SELECT 
					p.id, p.title, p.content, p.likes, p.dislikes, p.comments,
					(SELECT GROUP_CONCAT(c2.name)
					 FROM post_categories pc2
					 JOIN categories c2 ON c2.id = pc2.category_id
					 WHERE pc2.post_id = p.id) AS categories,
					COALESCE(l.value, 0) as userValue
				FROM posts p
				LEFT JOIN likes l ON l.post_id = p.id AND l.user_id = ?
				WHERE p.id IN (
					SELECT pc3.post_id
					FROM post_categories pc3
					JOIN categories c3 ON c3.id = pc3.category_id
					WHERE c3.name = ?
				)
				ORDER BY p.created_at DESC
			`, userID, category)
		}

		if err != nil {
			backend.WriteJSONError(w, http.StatusInternalServerError, "something went wrong")
			log.Println("Error querying posts:", err)
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
				backend.WriteJSONError(w, http.StatusInternalServerError, "scan error")
				log.Println("Error scanning post:", err)
				return
			}
			post.Categories = strings.Split(category, ",")
			posts.Posts = append(posts.Posts, post)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(posts)
	}
}
