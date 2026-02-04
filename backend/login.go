package backend

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func LoginHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var creds struct {
			Username string `json:"username"` // nickname أو email
			Password string `json:"password"`
		}

		if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		if creds.Username == "" || creds.Password == "" {
			http.Error(w, "Nickname/email and password required", http.StatusBadRequest)
			return
		}

		var userID int64
		var passwordHash string

		err := db.QueryRow(
			"SELECT id, password FROM users WHERE nickname = ? OR email = ?",
			creds.Username, creds.Username,
		).Scan(&userID, &passwordHash)

		if err == sql.ErrNoRows {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}
		if err != nil {
			log.Println("DB error:", err)
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}

		if bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(creds.Password)) != nil {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		token := generateRandomToken()
		exp := time.Now().Add(24 * time.Hour)

		db.Exec("DELETE FROM sessions WHERE user_id = ?", userID)

		_, err = db.Exec(
			"INSERT INTO sessions(token, user_id, expires_at) VALUES (?, ?, ?)",
			token, userID, exp,
		)
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
		json.NewEncoder(w).Encode(map[string]string{
			"status": "ok",
		})
	}
}

func generateRandomToken() string {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(b)
}
