package backend

import (
	"database/sql"
	"net/http"
)

type User struct {
	ID       int    `json:"id"`
	Nickname string `json:"nickname"`
}

func GetUserIDFromRequest(DB *sql.DB, r *http.Request) (User, error) {
	c, err := r.Cookie("session_token")
	var user User
	if err != nil {
		return user, err
	}
	token := c.Value

	err = DB.QueryRow(
		`SELECT s.user_id, u.nickname
FROM sessions s
JOIN users u ON s.user_id = u.id
WHERE s.token = ? AND s.expires_at > datetime('now');
`,
		token,
	).Scan(&user.ID, &user.Nickname)
	if err != nil {
		return user, err
	}

	return user, nil
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
