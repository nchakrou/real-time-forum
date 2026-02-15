package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)
// select all comments that the post has and show it
func HandleGetComments(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(map[string]string{
				"status":  "error",
				"message": "Method not allowed",
			})
			return
		}

		postID := r.URL.Query().Get("post_id")
		if postID == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"status":  "error",
				"message": "Missing post_id",
			})
			return
		}

		pid, err := strconv.Atoi(postID)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"status":  "error",
				"message": "Invalid post_id",
			})
			return
		}

		rows, err := db.Query(`
			SELECT users.nickname, comments.content
			FROM comments
			JOIN users ON users.id = comments.user_id
			WHERE comments.post_id = ?
			ORDER BY comments.id ASC
		`, pid)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println("DB query error:", err)
			json.NewEncoder(w).Encode(map[string]string{
				"status":  "error",
				"message": "Database error",
			})
			return
		}
		defer rows.Close()

		var comments []map[string]string
		for rows.Next() {
			var username, content string
			if err := rows.Scan(&username, &content); err != nil {
				log.Println("scan err:", err)
				continue
			}
			comments = append(comments, map[string]string{
				"username": username,
				"content":  content,
			})
		}

		json.NewEncoder(w).Encode(comments)
	}
}
