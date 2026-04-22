package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"forum/backend"
	"forum/backend/handlers"
	"forum/backend/handlers/Websocket"
	"forum/backend/handlers/auth"

	"github.com/gorilla/websocket"
)


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

	http.HandleFunc("/frontend/", func(w http.ResponseWriter, r *http.Request) {
		infos, err := os.Stat(r.URL.Path[1:])
		if err != nil {
			http.Error(w,"internal server error", http.StatusInternalServerError)
			return
		}
		if infos.IsDir() {
		http.Error(w,"internal server error", http.StatusInternalServerError)
			return
		}
		if os.IsNotExist(err) {
			http.Error(w,"internal server error", http.StatusInternalServerError)
			return
		}

		http.ServeFile(w, r, r.URL.Path[1:])
	})

	http.HandleFunc("/api/login", auth.LoginHandler(db))
	http.HandleFunc("/api/register", auth.RegisterHandler(db))
	http.HandleFunc("/api/islogged", auth.IsLogged(db))
	http.HandleFunc("/api/logout", auth.Logout(db))

	http.HandleFunc("/api/liked-posts", backend.Middleware(db, handlers.HandleLikedPosts(db)))
	http.HandleFunc("/api/posts", backend.Middleware(db, handlers.GetPostsHandler(db)))
	http.HandleFunc("/api/createpost", backend.Middleware(db, handlers.CreatePostHandler(db)))
	http.HandleFunc("/api/myposts", backend.Middleware(db, handlers.GetMyPostsHandler(db)))

	http.HandleFunc("/api/like", backend.Middleware(db, handlers.HandleLike(db, "post")))
	http.HandleFunc("/api/like-comment", backend.Middleware(db, handlers.HandleLike(db, "comment")))
	http.HandleFunc("/api/add-comment", backend.Middleware(db, handlers.HandleAddComment(db)))
	http.HandleFunc("/api/comments", backend.Middleware(db, handlers.HandleGetComments(db)))

	http.HandleFunc("/ws", backend.Middleware(db, Websocket.WsHandler(db, hub)))

	fmt.Println("Server started at http://localhost:8081")
	log.Fatal(http.ListenAndServe(":8081", nil))

}
