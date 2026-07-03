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
	db.SetMaxOpenConns(1)
	if err := migrate(db); err != nil {
		db.Close()
		return nil, fmt.Errorf("migrate: %w", err)
	}
	return db, nil
}

func migrate(db *sql.DB) error {
	// Step 1: create tables fresh
	schema := `
	CREATE TABLE IF NOT EXISTS users (
		id          INTEGER PRIMARY KEY AUTOINCREMENT,
		email       TEXT NOT NULL UNIQUE,
		password    TEXT NOT NULL,
		is_admin    INTEGER NOT NULL DEFAULT 0,
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
		id            INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id       INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		slug          TEXT NOT NULL UNIQUE,
		name          TEXT NOT NULL DEFAULT '',
		password      TEXT NOT NULL DEFAULT '',
		password_plain TEXT NOT NULL DEFAULT '',
		created_at    DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at    DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
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

	CREATE TABLE IF NOT EXISTS settings (
		key   TEXT PRIMARY KEY,
		value TEXT NOT NULL
	);

	CREATE TABLE IF NOT EXISTS site_visits (
		id          INTEGER PRIMARY KEY AUTOINCREMENT,
		site_id     INTEGER NOT NULL REFERENCES sites(id) ON DELETE CASCADE,
		visit_date  TEXT NOT NULL,  -- YYYY-MM-DD
		visit_month TEXT NOT NULL,  -- YYYY-MM
		created_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
	);
	CREATE INDEX IF NOT EXISTS idx_visits_site ON site_visits(site_id);
	CREATE INDEX IF NOT EXISTS idx_visits_date ON site_visits(visit_date);
	CREATE INDEX IF NOT EXISTS idx_visits_month ON site_visits(visit_month);
	`
	if _, err := db.Exec(schema); err != nil {
		return err
	}

	// Step 2: add columns if they doesn't exist (for existing DBs)
	_, _ = db.Exec(`ALTER TABLE users ADD COLUMN is_admin INTEGER NOT NULL DEFAULT 0`)
	_, _ = db.Exec(`ALTER TABLE sites ADD COLUMN password_plain TEXT NOT NULL DEFAULT ''`)

	// Step 3: seed default settings
	_, _ = db.Exec(`INSERT OR IGNORE INTO settings (key, value) VALUES ('open_registration', '1')`)
	_, _ = db.Exec(`INSERT OR IGNORE INTO settings (key, value) VALUES ('allow_public_access', '1')`)
	_, _ = db.Exec(`INSERT OR IGNORE INTO settings (key, value) VALUES ('domain_restriction', '0')`)
	_, _ = db.Exec(`INSERT OR IGNORE INTO settings (key, value) VALUES ('allowed_domains', '')`)
	_, _ = db.Exec(`INSERT OR IGNORE INTO settings (key, value) VALUES ('max_sites_per_user', '30')`)

	return nil
}
