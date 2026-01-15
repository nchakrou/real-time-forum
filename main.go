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
	http.HandleFunc("/check", func(w http.ResponseWriter, r *http.Request) {
		var p login
		jsonData, _ := io.ReadAll(r.Body)
		json.Unmarshal(jsonData, &p)
		username := p.Username
		password := p.Password
		fmt.Println(username, password)
		if pass, ok := users[username]; ok && pass == password {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("success"))
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("unauthorized"))
		}
	})
	fmt.Println("Server started at http://localhost:8081")
	http.ListenAndServe(":8081", nil)
}
