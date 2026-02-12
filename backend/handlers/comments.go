package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"forum/backend"
)

func HandleAddComment(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(map[string]string{
				"status":  "error",
				"message": "Method not allowed",
			})
			return
		}

		userID, err := backend.GetUserIDFromRequest(db, r)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{
				"status":  "error",
				"message": "Unauthorized",
			})
			return
		}

		postID := r.FormValue("post_id")
		content := r.FormValue("comment")
		if postID == "" || content == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"status":  "error",
				"message": "Missing post_id or comment",
			})
			return
		}

		var dummy int
		err = db.QueryRow("SELECT id FROM posts WHERE id = ?", postID).Scan(&dummy)
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"status":  "error",
				"message": "Post not found",
			})
			return
		}

		_, err = db.Exec("INSERT INTO comments (post_id, user_id, content) VALUES (?, ?, ?)", postID, userID, content)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{
				"status":  "error",
				"message": "Database error",
			})
			return
		}

		json.NewEncoder(w).Encode(map[string]string{
			"status": "ok",
		})
	}
}
