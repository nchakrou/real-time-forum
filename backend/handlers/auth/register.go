package auth

import (
	"database/sql"
	"encoding/json"
	"forum/backend"
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
			backend.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
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
			backend.WriteJSONError(w, http.StatusBadRequest, "something went wrong. Please try again.")
			return
		}
		if !emailrg.MatchString(req.Email) {
			backend.WriteJSONError(w, http.StatusBadRequest, "invalid email")
			return
		}
		if req.FirstName == "" || req.LastName == "" || req.Nickname == "" || req.Age == "" || req.Email == "" || req.Password == "" || req.Gender == "" {
			backend.WriteJSONError(w, http.StatusBadRequest, "All fields are required to create your account.")
			return
		}
		if len(req.Password) < 8 {
			backend.WriteJSONError(w, http.StatusBadRequest, "Your password is too short—please use at least 8 characters.")
			return
		}
		ageInt, err := strconv.Atoi(req.Age)
		if err != nil || ageInt < 18 || ageInt > 80 {
			backend.WriteJSONError(w, http.StatusBadRequest, "Please enter a valid age from 18 to 80.")
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
			backend.WriteJSONError(w, http.StatusConflict, "This nickname or email is already in use. Try another one!")
			return
		}

		// auto login //
		userID, err := res.LastInsertId()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		token := generateRandomToken()
		exp := time.Now().Add(24 * time.Hour)

		_, err = db.Exec("INSERT INTO sessions(token, user_id, expires_at) VALUES (?, ?, ?)",
			token, userID, exp)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
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
