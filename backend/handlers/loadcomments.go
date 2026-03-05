package handlers

import (
	"database/sql"
	"encoding/json"
	"forum/backend"
	"log"
	"net/http"
	"strconv"
)

// select all comments that the post has and show it
func HandleGetComments(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.Method != http.MethodGet {
			backend.WriteJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}

		postIDstr := r.URL.Query().Get("post_id")
		if postIDstr == "" {
			backend.WriteJSONError(w, http.StatusBadRequest, "something went wrong. Please try again.")
			return
		}

		postID, err := strconv.Atoi(postIDstr)
		if err != nil {
			backend.WriteJSONError(w, http.StatusBadRequest, "something went wrong. Please try again.")
			return
		}

		rows, err := db.Query(`
			SELECT users.nickname, comments.content
			FROM comments
			JOIN users ON users.id = comments.user_id
			WHERE comments.post_id = ?
			ORDER BY comments.id ASC
		`, postID)
		if err != nil {
			log.Println("DB query error:", err)
			backend.WriteJSONError(w, http.StatusInternalServerError, "something went wrong. Please try again later.")
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
