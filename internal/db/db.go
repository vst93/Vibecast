package db

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

func Open(path string) (*sql.DB, error) {
	dsn := fmt.Sprintf("%s?_pragma=journal_mode(WAL)&_pragma=busy_timeout(5000)&_pragma=foreign_keys(ON)", path)
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}
	db.SetMaxOpenConns(1) // SQLite single-writer
	if err := migrate(db); err != nil {
		db.Close()
		return nil, fmt.Errorf("migrate: %w", err)
	}
	return db, nil
}

func migrate(db *sql.DB) error {
	schema := `
	CREATE TABLE IF NOT EXISTS users (
		id          INTEGER PRIMARY KEY AUTOINCREMENT,
		email       TEXT NOT NULL UNIQUE,
		password    TEXT NOT NULL,
		created_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS sessions (
		id          INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id     INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		token       TEXT NOT NULL UNIQUE,
		expires_at  DATETIME NOT NULL,
		created_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
	);
	CREATE INDEX IF NOT EXISTS idx_sessions_token ON sessions(token);

	CREATE TABLE IF NOT EXISTS sites (
		id           INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id      INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		slug         TEXT NOT NULL UNIQUE,
		name         TEXT NOT NULL DEFAULT '',
		password     TEXT NOT NULL DEFAULT '',  -- bcrypt hash, empty = no password
		created_at   DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at   DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
	);
	CREATE INDEX IF NOT EXISTS idx_sites_user ON sites(user_id);
	CREATE INDEX IF NOT EXISTS idx_sites_slug ON sites(slug);

	CREATE TABLE IF NOT EXISTS site_sessions (
		id          INTEGER PRIMARY KEY AUTOINCREMENT,
		site_id     INTEGER NOT NULL REFERENCES sites(id) ON DELETE CASCADE,
		token       TEXT NOT NULL UNIQUE,
		expires_at  DATETIME NOT NULL,
		created_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
	);
	CREATE INDEX IF NOT EXISTS idx_site_sessions_token ON site_sessions(token);
	`
	_, err := db.Exec(schema)
	return err
}
