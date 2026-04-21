package handlers

import (
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"

	"forum/backend"
)

type post struct {
	Title      string `json:"title"`
	Content    string `json:"content"`
	Categories []int  `json:"categories"`
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
		if strings.TrimSpace(postData.Title) == "" || strings.TrimSpace(postData.Content) == "" {
			backend.WriteJSONError(w, http.StatusBadRequest, "Please fill in all the fields.")
			log.Println("Empty title or content")
			return
		}
		if len(postData.Title) > 100 || len(postData.Content) > 1000 {
			backend.WriteJSONError(w, http.StatusBadRequest, "Title must be less than 100 characters and content must be less than 1000 characters.")
			log.Println("Title or content too long")
			return
		}
		if len(postData.Categories) == 0 {
			backend.WriteJSONError(w, http.StatusBadRequest, "Please select at least one category for your post.")
			log.Println("No categories selected")
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
