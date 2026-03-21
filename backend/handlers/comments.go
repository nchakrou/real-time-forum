package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"forum/backend"
)

// insert the comment into the db
func HandleAddComment(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.Method != http.MethodPost {
			backend.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
			log.Println("Method not allowed")
			return
		}

		user, err := backend.GetUserIDFromRequest(db, r)
		userID := user.ID
		if err != nil {
			backend.WriteJSONError(w, http.StatusUnauthorized, "You must be logged in to leave a comment.")
			return
		}

		postIDstr := r.FormValue("post_id")
		content := r.FormValue("comment")
		if postIDstr == "" || content == "" {
			backend.WriteJSONError(w, http.StatusBadRequest, "Please write something before posting your comment.")
			return
		}

		postID, err := strconv.Atoi(postIDstr)
		if err != nil {
			backend.WriteJSONError(w, http.StatusBadRequest, "We couldn't identify which post you're commenting on.")
			log.Println("Error converting post_id to int:", err)
			return
		}

		var dummy int
		err = db.QueryRow("SELECT id FROM posts WHERE id = ?", postID).Scan(&dummy)
		if err == sql.ErrNoRows {
			backend.WriteJSONError(w, http.StatusBadRequest, "The post you're trying to comment on no longer exists.")
			return
		}

		_, err = db.Exec("INSERT INTO comments (post_id, user_id, content) VALUES (?, ?, ?)", postID, userID, content)
		if err != nil {
			backend.WriteJSONError(w, http.StatusInternalServerError, "something went wrong. Please try again later.")
			log.Println("Error inserting comment into database:", err)
			return
		}

		_, err = db.Exec("UPDATE posts SET comments = comments + 1 WHERE id = ?", postID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println("Error updating comment count:", err)
			return
		}
		var commentsCount int
		err = db.QueryRow(
			"SELECT COUNT(*) FROM comments WHERE post_id = ?",
			postID,
		).Scan(&commentsCount)
		
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
			_, _ = db.Exec(
			"UPDATE posts SET comments = ? WHERE id = ?",
			commentsCount, postID,
		)	
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":   "ok",
			"comments": commentsCount,
		})
	}
}
