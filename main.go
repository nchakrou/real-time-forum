package main

import (
	"fmt"
	"forum/backend"
	"forum/backend/handlers"
	"forum/backend/handlers/Websocket"
	"forum/backend/handlers/auth"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
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
	hub := &Websocket.Hub{
		Clients: make(map[int][]*websocket.Conn),
	}
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

	http.HandleFunc("/api/login", auth.LoginHandler(db))
	http.HandleFunc("/api/register", auth.RegisterHandler(db))
	http.HandleFunc("/api/islogged", auth.IsLogged(db))
	http.HandleFunc("/api/logout", func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{Name: "session_token", Value: "", MaxAge: -1, Path: "/"})
	})
	http.HandleFunc("/api/like", handlers.HandleLike(db, "post"))
	http.HandleFunc("/api/like-comment", handlers.HandleLike(db, "comment"))
	http.HandleFunc("/api/add-comment", handlers.HandleAddComment(db))
	http.HandleFunc("/api/comments", handlers.HandleGetComments(db))

	http.HandleFunc("/api/posts", handlers.GetPostsHandler(db))
	http.HandleFunc("/api/createpost", handlers.CreatePostHandler(db))
	http.HandleFunc("/api/myposts", handlers.GetMyPostsHandler(db))
	http.HandleFunc("/ws", Websocket.WsHandler(db, hub))

	fmt.Println("Server started at http://localhost:8081")
	log.Fatal(http.ListenAndServe(":8081", nil))

}
