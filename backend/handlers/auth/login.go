package auth

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"forum/backend"

	"golang.org/x/crypto/bcrypt"
)

func LoginHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			writeJSONError(w, http.StatusMethodNotAllowed, "method", "method not allowed")
			return
		}

		var creds struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
			writeJSONError(w, http.StatusBadRequest, "json", "invalid JSON")
			return
		}
		if creds.Username == "" || creds.Password == "" {
			writeJSONError(w, http.StatusBadRequest, "credentials", "username and password required")
			return
		}
		var userID int64
		var passwordHash string

		err := db.QueryRow(
			"SELECT id, password FROM users WHERE nickname = ? OR email = ?",
			creds.Username, creds.Username,
		).Scan(&userID, &passwordHash)
		if err == sql.ErrNoRows {
			writeJSONError(w, http.StatusUnauthorized, "credentials", "invalid credentials")
			return
		}
		if err != nil {
			log.Println("DB error during login:", err)
			writeJSONError(w, http.StatusInternalServerError, "db", "something went wrong")
			return
		}
		if bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(creds.Password)) != nil {
			writeJSONError(w, http.StatusUnauthorized, "credentials", "invalid credentials")
			return
		}

		token := generateRandomToken()
		exp := time.Now().Add(24 * time.Hour)

		_, _ = db.Exec("DELETE FROM sessions WHERE user_id = ?", userID)

		_, err = db.Exec(
			"INSERT INTO sessions(token, user_id, expires_at) VALUES (?, ?, ?)",
			token, userID, exp,
		)
		if err != nil {
			log.Println("Session insert error:", err)
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
		json.NewEncoder(w).Encode(map[string]string{
			"status": "ok",
		})
	}
}

func IsLogged(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := backend.GetUserIDFromRequest(db, r)
		if err != nil || userID.ID == 0 {
			w.WriteHeader(http.StatusUnauthorized)
			log.Println("Unauthorized:", err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"nickname": userID.Nickname,
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
