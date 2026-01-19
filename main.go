package main

import (
	"encoding/json"
	"fmt"
	"io"
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
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "frontend/index.html")
	})

	http.HandleFunc("/src/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.URL.Path)
		http.ServeFile(w, r, "frontend/"+r.URL.Path)
	})
	http.HandleFunc("/api/login", func(w http.ResponseWriter, r *http.Request) {
		var p login
		jsonData, _ := io.ReadAll(r.Body)
		json.Unmarshal(jsonData, &p)
		username := p.Username
		password := p.Password
		fmt.Println(username, password)
		if pass, ok := users[username]; ok && pass == password {
			cookie := http.Cookie{
				Name:  "session",
				Value: username,
				Path:  "/",
			}
			http.SetCookie(w, &cookie)
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("success"))
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("unauthorized"))
		}
	})
	http.HandleFunc("/api/islogged", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("ww")
		cookie, err := r.Cookie("session")
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("unauthorized"))
		}
		fmt.Println(err, cookie)

	})
	http.HandleFunc("/api/logout", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("fdgdf")
		http.SetCookie(w, &http.Cookie{Name: "session", Value: "", MaxAge: -1, Path: "/"})
	})
	fmt.Println("Server started at http://localhost:8081")
	http.ListenAndServe(":8081", nil)

}
