package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"forum/backend"
)

type LikeRequest struct {
	Value int `json:"value"` // 1 = like, -1 = dislike
}

type LikeResponse struct {
	Likes     int `json:"likes"`
	Dislikes  int `json:"dislikes"`
	UserValue int `json:"userValue"`
}

func HandleLike(db *sql.DB, target string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("like")
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		userID, err := backend.GetUserIDFromRequest(db, r)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		var req LikeRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		if req.Value != 1 && req.Value != -1 {
			http.Error(w, "Invalid value", http.StatusBadRequest)
			return
		}

		idStr := r.URL.Query().Get("id")
		id, err := strconv.Atoi(idStr)
		if err != nil || id <= 0 {
			http.Error(w, "Invalid id", http.StatusBadRequest)
			return
		}

		table, column := "likes", "post_id"
		if target == "comment" {
			column = "comment_id"
		}

		tx, err := db.Begin()
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}
		defer tx.Rollback()

		var existing int
		err = tx.QueryRow(
			"SELECT value FROM "+table+" WHERE user_id=? AND "+column+"=?",
			userID, id,
		).Scan(&existing)

		userValue := 0

		if err == sql.ErrNoRows {

			_, err = tx.Exec(
				"INSERT INTO "+table+" (user_id,"+column+",value) VALUES (?,?,?)",
				userID, id, req.Value,
			)
			if err != nil {
				http.Error(w, "Database error", http.StatusInternalServerError)
				return
			}
			userValue = req.Value

		} else if err == nil {

			if existing == req.Value {

				_, err = tx.Exec(
					"DELETE FROM "+table+" WHERE user_id=? AND "+column+"=?",
					userID, id,
				)
				if err != nil {
					http.Error(w, "Database error", http.StatusInternalServerError)
					return
				}
				userValue = 0

			} else {

				_, err = tx.Exec(
					"UPDATE "+table+" SET value=? WHERE user_id=? AND "+column+"=?",
					req.Value, userID, id,
				)
				if err != nil {
					http.Error(w, "Database error", http.StatusInternalServerError)
					return
				}
				userValue = req.Value
			}

		} else {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		var likes, dislikes int

		err = tx.QueryRow(
			"SELECT COUNT(*) FROM "+table+" WHERE "+column+"=? AND value=1",
			id,
		).Scan(&likes)
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}
		err = tx.QueryRow(
			"SELECT COUNT(*) FROM "+table+" WHERE "+column+"=? AND value=-1",
			id,
		).Scan(&dislikes)
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}
		if target == "post" {
			_, err = tx.Exec(
				"UPDATE posts SET likes=?, dislikes=? WHERE id=?",
				likes, dislikes, id,
			)
			if err != nil {
				http.Error(w, "Database error", http.StatusInternalServerError)
				return
			}
		}
		if err = tx.Commit(); err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}
		resp := LikeResponse{
			Likes:     likes,
			Dislikes:  dislikes,
			UserValue: userValue,
		}
		fmt.Println(resp)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}
