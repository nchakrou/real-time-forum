package backend

import (
	"database/sql"
	"errors"
	"net/http"
	"strings"
	"time"
)

type User struct {
	ID       int    `json:"id"`
	Nickname string `json:"nickname"`
}

var (
	ErrNoSession      = errors.New("no session found")
	ErrInvalidSession = errors.New("invalid or expired session")
)

func GetUserIDFromRequest(db *sql.DB, r *http.Request) (User, error) {
	var user User

	cookie, err := r.Cookie("session_token")
	if err != nil {
		return user, ErrNoSession
	}

	token := strings.TrimSpace(cookie.Value)
	if token == "" {
		return user, ErrNoSession
	}

	query := `
		SELECT s.user_id, u.nickname
		FROM sessions s
		JOIN users u ON s.user_id = u.id
		WHERE s.token = ? AND s.expires_at > ?
	`

	err = db.QueryRow(query, token, time.Now()).Scan(&user.ID, &user.Nickname)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, ErrInvalidSession
		}
		return user, err
	}

	return user, nil
}

func Middleware(DB *sql.DB, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if _, err := GetUserIDFromRequest(DB, r); err != nil {
			WriteJSONError(w, http.StatusUnauthorized, "login required")
			return
		}
		next.ServeHTTP(w, r)
	}
}
