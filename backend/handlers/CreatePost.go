package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"forum/backend"
	"io"
	"log"
	"net/http"
)

type post struct {
	Title      string `json:"title"`
	Content    string `json:"content"`
	Categories []int  `json:"categories"`
}

func CreatePostHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("ok")

		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			log.Println("Method not allowed")
			return
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println("Error reading request body:", err)
			return
		}
		var postData post
		err = json.Unmarshal(body, &postData)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Println("Error unmarshaling JSON:", err)
			return
		}
		userid := backend.GetUserIDFromRequest(db, r)
		if userid == 0 {
			w.WriteHeader(http.StatusUnauthorized)
			log.Println("Unauthorized: no valid session")
			return
		}
		res, err := db.Exec("INSERT INTO posts (user_id, title, content) VALUES (?, ?, ?)", userid, postData.Title, postData.Content)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
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
				w.WriteHeader(http.StatusBadRequest)
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
