package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"forum/backend"
)

// select comments with limit & offset for "Load more"
func HandleGetComments(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.Method != http.MethodGet {
			backend.WriteJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}

		postIDstr := r.URL.Query().Get("post_id")
		offsetStr := r.URL.Query().Get("offset")
		limitStr := r.URL.Query().Get("limit")

		if postIDstr == "" {
			backend.WriteJSONError(w, http.StatusBadRequest, "something went wrong. Please try again.")
			return
		}

		postID, err := strconv.Atoi(postIDstr)
		if err != nil {
			backend.WriteJSONError(w, http.StatusBadRequest, "invalid post_id")
			return
		}

		// Default values
		offset := 0
		limit := 10

		if offsetStr != "" {
			offset, _ = strconv.Atoi(offsetStr)
		}
		if limitStr != "" {
			limit, _ = strconv.Atoi(limitStr)
		}

		rows, err := db.Query(`
			SELECT users.nickname, comments.content, comments.created_at
			FROM comments
			JOIN users ON users.id = comments.user_id
			WHERE comments.post_id = ?
			ORDER BY comments.created_at DESC
			LIMIT ? OFFSET ?
		`, postID, limit, offset)

		if err != nil {
			log.Println("DB query error:", err)
			backend.WriteJSONError(w, http.StatusInternalServerError, "something went wrong. Please try again later.")
			return
		}
		defer rows.Close()

		var comments []map[string]string
		for rows.Next() {
			var username, content, createdAt string
			if err := rows.Scan(&username, &content, &createdAt); err != nil {
				log.Println("scan err:", err)
				continue
			}
			comments = append(comments, map[string]string{
				"username":  username,
				"content":   content,
				"createdAt": createdAt,
			})
		}

		json.NewEncoder(w).Encode(comments)
	}
}
