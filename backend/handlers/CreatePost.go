package handlers

import (
	"database/sql"
	"encoding/json"
	"forum/backend"
	"io"
	"log"
	"net/http"
)

type post struct {
	Title      string `json:"title"`
	Content    string `json:"content"`
	Categories []int  `json:"categories"`
	IsEnd      bool   `json:"is_end"`
}

func CreatePostHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPost {
			backend.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
			log.Println("Method not allowed")
			return
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			backend.WriteJSONError(w, http.StatusBadRequest, "something went wrong. Please try again.")
			log.Println("Error reading request body:", err)
			return
		}
		var postData post
		err = json.Unmarshal(body, &postData)
		if err != nil {
			backend.WriteJSONError(w, http.StatusBadRequest, "Something is wrong with your post's format. Please check and try again.")
			log.Println("Error unmarshaling JSON:", err)
			return
		}
		user, err := backend.GetUserIDFromRequest(db, r)
		userid := user.ID
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			log.Println("Error getting user ID from request:", err)
			return
		}
		res, err := db.Exec("INSERT INTO posts (user_id, title, content) VALUES (?, ?, ?)", userid, postData.Title, postData.Content)
		if err != nil {
			backend.WriteJSONError(w, http.StatusInternalServerError, "something went wrong. Please try again later.")
			log.Println("Error inserting post into database:", err)
			return
		}

		postID, err := res.LastInsertId()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println("Error getting last insert ID:", err)
			return
		}
		for _, categoryID := range postData.Categories {
			if categoryID <= 0 || categoryID > 7 {
				backend.WriteJSONError(w, http.StatusBadRequest, "Please select a valid category for your post.")
				log.Println("Invalid category ID:", categoryID)
				return
			}
			_, err = db.Exec("INSERT INTO post_categories (post_id, category_id) VALUES (?, ?)", postID, categoryID)

			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Println("Error inserting post category into database:", err)
				return
			}
		}
		w.WriteHeader(http.StatusCreated)

	}
}
