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
			backend.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		var creds struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
			backend.WriteJSONError(w, http.StatusBadRequest, "something went wrong. Please check your entry.")
			return
		}
		if creds.Username == "" || creds.Password == "" {
			backend.WriteJSONError(w, http.StatusBadRequest, "Please enter both your nickname/email and password.")
			return
		}
		var userID int64
		var passwordHash string

		err := db.QueryRow(
			"SELECT id, password FROM users WHERE nickname = ? OR email = ?",
			creds.Username, creds.Username,
		).Scan(&userID, &passwordHash)
		if err == sql.ErrNoRows {
			backend.WriteJSONError(w, http.StatusUnauthorized, "invalid credentials")
			return
		}
		if err != nil {
			log.Println("DB error during login:", err)
			backend.WriteJSONError(w, http.StatusInternalServerError, "something went wrong. Please try again.")
			return
		}
		if bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(creds.Password)) != nil {
			backend.WriteJSONError(w, http.StatusUnauthorized, "invalid credentials")
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
			backend.WriteJSONError(w, http.StatusInternalServerError, "something went wrong. Please try again.")
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "session_token",
			Value:    token,
			Expires:  exp,
			Path:     "/",
			HttpOnly: true,
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

func Logout(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := backend.GetUserIDFromRequest(db, r)
		if err != nil || userID.ID == 0 {
			w.WriteHeader(http.StatusUnauthorized)
			log.Println("Unauthorized logout attempt:", err)
			return
		}

		_, err = db.Exec("DELETE FROM sessions WHERE user_id = ?", userID.ID)
		if err != nil {
			log.Println("Error deleting session during logout:", err)
			backend.WriteJSONError(w, http.StatusInternalServerError, "something went wrong. Please try again.")
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "session_token",
			Value:    "",
			MaxAge:   -1,
			Path:     "/",
			HttpOnly: true,
		})

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status": "ok",
		})
	}
}
