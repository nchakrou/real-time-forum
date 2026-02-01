package backend

import (
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func InitDB(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	db.SetMaxOpenConns(1)
	db.SetConnMaxLifetime(time.Minute * 10)

	_, err = db.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		db.Close()
		return nil, err
	}

	if err := Schema(db); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

func Schema(db *sql.DB) error {
	schema := `
CREATE TABLE IF NOT EXISTS users (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	nickname TEXT NOT NULL UNIQUE,
	first_name TEXT NOT NULL,
	last_name TEXT NOT NULL,
	username TEXT NOT NULL UNIQUE,
	age INTEGER NOT NULL,
	email TEXT NOT NULL UNIQUE,
	password TEXT NOT NULL,
	gender TEXT NOT NULL,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

	CREATE TABLE IF NOT EXISTS sessions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		token TEXT NOT NULL UNIQUE,
		expires_at DATETIME NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	);

	CREATE INDEX IF NOT EXISTS idx_sessions_token ON sessions(token);
	CREATE UNIQUE INDEX IF NOT EXISTS idx_sessions_user ON sessions(user_id);
	`

	_, err := db.Exec(schema)
	return err
}
