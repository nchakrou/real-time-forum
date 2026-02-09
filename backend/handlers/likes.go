package handlers

import (
	"database/sql"
	"encoding/json"
	"forum/backend"
	"log"
	"net/http"
	"strconv"
)

type LikeRequest struct {
	Value int `json:"value"` // 1 = like -1 = dislike
}

type LikeResponse struct {
	Likes     int `json:"likes"`
	Dislikes  int `json:"dislikes"`
	UserValue int `json:"userValue"`
}

func HandleLike(db *sql.DB, target string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		userID, err := backend.GetUserIDFromRequest(db, r)
		if err != nil {
			log.Println("err1", err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		log.Println("err2", err)

		var req LikeRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Println("err3", err)
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		idStr := r.URL.Query().Get("id")
		id, _ := strconv.Atoi(idStr)
		if id == 0 {
			log.Println("id", idStr)
			http.Error(w, "Missing id", http.StatusBadRequest)
			return
		}

		table, column := "likes", "post_id"
		if target == "comment" {
			column = "comment_id"
		}

		var existing int
		err = db.QueryRow("SELECT kind FROM "+table+" WHERE user_id=? AND "+column+"=?", userID, id).Scan(&existing)
		userValue := 0

		if err == sql.ErrNoRows {
			db.Exec("INSERT INTO "+table+" (user_id,"+column+",kind) VALUES (?,?,?)", userID, id, req.Value)
			userValue = req.Value
		} else if err == nil {
			if existing == req.Value {
				db.Exec("DELETE FROM "+table+" WHERE user_id=? AND "+column+"=?", userID, id)
				userValue = 0
			} else {
				db.Exec("UPDATE "+table+" SET kind=? WHERE user_id=? AND "+column+"=?", req.Value, userID, id)
				userValue = req.Value
			}
		}

		var likes, dislikes int
		db.QueryRow("SELECT COUNT(*) FROM "+table+" WHERE "+column+"=? AND kind=1", id).Scan(&likes)
		db.QueryRow("SELECT COUNT(*) FROM "+table+" WHERE "+column+"=? AND kind=-1", id).Scan(&dislikes)

		resp := LikeResponse{Likes: likes, Dislikes: dislikes, UserValue: userValue}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}
