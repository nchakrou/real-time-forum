package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"regexp"
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
		emailrg := regexp.MustCompile(`^[A-Za-z0-9._%+\-]+@[A-Za-z0-9.\-]+\.[A-Za-z]{2,}$`)

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
			writeJSONError(w, http.StatusBadRequest, "json", "invalid JSON")
			return
		}
		if !emailrg.MatchString(req.Email) {
			writeJSONError(w, http.StatusBadRequest, "email", "invalid email")
			return
		}
		if req.FirstName == "" {
			writeJSONError(w, http.StatusBadRequest, "firstname", "first name is required")
				return
		}
		if req.LastName == "" {
			writeJSONError(w, http.StatusBadRequest, "lastname", "last name is required")
			return
		}
		if req.Nickname == "" {
			writeJSONError(w, http.StatusBadRequest, "nickname","nickname is required")
			return
		}
		if req.Age == "" {
			writeJSONError(w, http.StatusBadRequest, "age","age is required")
			return
		}
		if req.Email == "" {
			writeJSONError(w, http.StatusBadRequest, "email","email is required")
			return
		}
		if req.Password == "" {
			writeJSONError(w, http.StatusBadRequest, "password","password is required")
			return
		}
		if req.Gender == "" {
			writeJSONError(w, http.StatusBadRequest, "gender","gender is required")
			return
		}
		if len(req.Password) < 8 {
			writeJSONError(w, http.StatusBadRequest, "password","password must be > 8")
			return
		}
		ageInt, err := strconv.Atoi(req.Age)
		if err != nil || ageInt <= 0 {
			writeJSONError(w, http.StatusBadRequest,"age","invalid age")
			return
		}

		hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		res, err := db.Exec(
			`INSERT INTO users(nickname, first_name, last_name, age, email, password, gender)
			 VALUES (?, ?, ?, ?, ?, ?, ?)`,
			req.Nickname, req.FirstName, req.LastName, ageInt, req.Email, hashed, req.Gender,
		)
		if err != nil {
			log.Println("Error inserting user:", err)
			writeJSONError(w, http.StatusConflict, "db", "user already exists or db error")
			return
		}

		// auto login //
		userID, err := res.LastInsertId()
		if err != nil {
			writeJSONError(w, http.StatusBadRequest, "db", "could not get user id")
			return
		}

		token := generateRandomToken()
		exp := time.Now().Add(24 * time.Hour)

		_, err = db.Exec("INSERT INTO sessions(token, user_id, expires_at) VALUES (?, ?, ?)",
			token, userID, exp)
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, "db", "session creation failed")
			return
		}
		http.SetCookie(w, &http.Cookie{
			Name:     "session_token",
			Value:    token,
			Expires:  exp,
			Path:     "/",
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteLaxMode,
		})
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	}
}

func writeJSONError(w http.ResponseWriter, status int, field, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"error":  message, "field":  field, "status": "error",
	})
}
