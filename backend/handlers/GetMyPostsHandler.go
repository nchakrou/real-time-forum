package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
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

		offsetstr := r.URL.Query().Get("offset")
		offset, err := strconv.Atoi(offsetstr)
		if err != nil {
			backend.WriteJSONError(w, http.StatusBadRequest, "invalid offset")
			log.Println("Error converting offset to integer:", err)
			return
		}

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
			LIMIT 10 OFFSET ?
		`, userId, userId, offset)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var posts = struct {
			Posts []Post
			IsEnd bool
		}{
			Posts: []Post{},
			IsEnd: false,
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

		if len(posts.Posts) < 10 {
			posts.IsEnd = true
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(posts)
	}
}
