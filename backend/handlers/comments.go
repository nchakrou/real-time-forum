package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"forum/backend"
)

func HandleAddComment(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusBadRequest)
			return
		}

		userID, err := backend.GetUserIDFromRequest(db, r)
		if err != nil {
			http.Error(w, "Method not allowed", http.StatusBadRequest)
			return
		}

		postID := r.FormValue("post_id")
		content := r.FormValue("comment")
		if postID == "" || content == "" {
			http.Error(w, "Champs manquants", http.StatusBadRequest)
			return
		}
		var dummy int 
		err = db.QueryRow("SELECT id FROM posts WHERE id = ?", postID).Scan(&dummy)

		if err == sql.ErrNoRows {
			http.Error(w, "Method not allowed", http.StatusBadRequest)
			return
		}
		_, err = db.Exec("INSERT INTO comments (post_id, user_id, comment) VALUES (?, ?, ?)", postID, userID, content)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Erreur base de données", http.StatusInternalServerError)
			return
		}

		rows, err := db.Query(`
			SELECT u.username, c.comment, c.created_at
			FROM comments c
			JOIN users u ON u.id = c.user_id
			WHERE c.post_id = ?
			ORDER BY c.created_at DESC`, postID)
		if err != nil {
			http.Error(w, "Erreur base de données", http.StatusInternalServerError)
			return
		}
		defer rows.Close()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status": "ok",
		})
	}
}
