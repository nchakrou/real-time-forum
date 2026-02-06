package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func RegisterHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req struct {
			FirstName string `json:"firstName"`
			LastName  string `json:"lastName"`
			Nickname  string `json:"nickname"`
			Age       string `json:"age"`
			Email     string `json:"email"`
			Password  string `json:"password"`
			Gender    string `json:"gender"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Println("Error decoding JSON:", err)
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		if req.FirstName == "" || req.LastName == "" || req.Nickname == "" ||
			req.Age == "" || req.Email == "" || req.Password == "" || req.Gender == "" {
			log.Println("Missing fields in registration request")
			http.Error(w, "Missing fields", http.StatusBadRequest)
			return
		}

		ageInt, err := strconv.Atoi(req.Age)
		if err != nil || ageInt <= 0 {
			http.Error(w, "Invalid age", http.StatusBadRequest)
			return
		}

		hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}

		res, err := db.Exec(
			`INSERT INTO users(nickname, first_name, last_name, age, email, password, gender)
			 VALUES (?, ?, ?, ?, ?, ?, ?)`,
			req.Nickname,
			req.FirstName,
			req.LastName,
			ageInt,
			req.Email,
			hashed,
			req.Gender,
		)
		if err != nil {
			log.Println("Error inserting user:", err)
			http.Error(w, "User already exists or DB error", http.StatusConflict)
			return
		}

		// auto login //
		userID, err := res.LastInsertId()
		if err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}

		token := generateRandomToken()
		exp := time.Now().Add(24 * time.Hour)

		_, err = db.Exec("INSERT INTO sessions(token, user_id, expires_at) VALUES (?, ?, ?)",
			token, userID, exp)
		if err != nil {
			http.Error(w, "Session creation failed", http.StatusInternalServerError)
			return
		}
		http.SetCookie(w, &http.Cookie{
			Name:     "session_token",
			Value:    token,
			Expires:  exp,
			Path:     "/",
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
		})
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{
			"status": "ok",
		})
	}
}
