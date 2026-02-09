package main

import (
	"fmt"
	"forum/backend"
	"forum/backend/handlers"
	"log"
	"net/http"
)

var users = map[string]string{
	"1": "123",
	"2": "1234",
}

type login struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func main() {

	db, err := backend.InitDB("forum.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "frontend/index.html")
	})

	http.HandleFunc("/src/", func(w http.ResponseWriter, r *http.Request) {

		http.ServeFile(w, r, "frontend/"+r.URL.Path)
	})

	http.HandleFunc("/api/login", handlers.LoginHandler(db))
	http.HandleFunc("/api/register", handlers.RegisterHandler(db))
	http.HandleFunc("/api/islogged", handlers.IsLogged(db))
	http.HandleFunc("/api/logout", func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{Name: "session_token", Value: "", MaxAge: -1, Path: "/"})
	})
	http.HandleFunc("/api/like", handlers.HandleLike(db, "post"))
	http.HandleFunc("/api/like-comment", handlers.HandleLike(db, "comment"))

	http.HandleFunc("/api/posts", handlers.GetPostsHandler(db))
	http.HandleFunc("/api/createpost", handlers.CreatePostHandler(db))
	http.HandleFunc("/api/myposts", handlers.GetMyPostsHandler(db))

	fmt.Println("Server started at http://localhost:8081")
	log.Fatal(http.ListenAndServe(":8081", nil))

}
