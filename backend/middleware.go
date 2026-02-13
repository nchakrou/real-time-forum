package backend

import (
	"database/sql"
	"net/http"
)

func GetUserIDFromRequest(DB *sql.DB, r *http.Request) (int, error) {
	c, err := r.Cookie("session_token")
	if err != nil {
		return 0, err
	}
	token := c.Value

	var userID int

	err = DB.QueryRow(
		"SELECT user_id FROM sessions WHERE token = ? AND expires_at > datetime('now')",
		token,
	).Scan(&userID)
	if err != nil {
		return 0, err
	}

	return userID, nil
}

func AuthRequired(DB *sql.DB, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if _, err := GetUserIDFromRequest(DB, r); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	}
}

func NotAuthRequired(DB *sql.DB, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if _, err := GetUserIDFromRequest(DB, r); err == nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		next.ServeHTTP(w, r)
	}
}
