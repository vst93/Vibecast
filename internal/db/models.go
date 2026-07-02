package db

import (
	"database/sql"
	"time"
)

type User struct {
	ID        int64
	Email     string
	Password  string // bcrypt hash
	IsAdmin   bool
	CreatedAt time.Time
}

type Site struct {
	ID            int64
	UserID        int64
	Slug          string
	Name          string
	Password      string // bcrypt hash, empty = no protection
	PasswordPlain string // plaintext password for display
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// --- Users ---

func CreateUser(db *sql.DB, email, hashedPassword string, isAdmin bool) (*User, error) {
	adminVal := 0
	if isAdmin {
		adminVal = 1
	}
	res, err := db.Exec(
		`INSERT INTO users (email, password, is_admin) VALUES (?, ?, ?)`,
		email, hashedPassword, adminVal,
	)
	if err != nil {
		return nil, err
	}
	id, _ := res.LastInsertId()
	return &User{ID: id, Email: email, Password: hashedPassword, IsAdmin: isAdmin}, nil
}

func GetUserByEmail(db *sql.DB, email string) (*User, error) {
	u := &User{}
	var adminVal int
	err := db.QueryRow(
		`SELECT id, email, password, is_admin, created_at FROM users WHERE email = ?`,
		email,
	).Scan(&u.ID, &u.Email, &u.Password, &adminVal, &u.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	u.IsAdmin = adminVal == 1
	return u, err
}

func GetUserByID(db *sql.DB, id int64) (*User, error) {
	u := &User{}
	var adminVal int
	err := db.QueryRow(
		`SELECT id, email, password, is_admin, created_at FROM users WHERE id = ?`,
		id,
	).Scan(&u.ID, &u.Email, &u.Password, &adminVal, &u.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	u.IsAdmin = adminVal == 1
	return u, err
}

func CountUsers(db *sql.DB) (int64, error) {
	var count int64
	err := db.QueryRow(`SELECT COUNT(*) FROM users`).Scan(&count)
	return count, err
}

func ListUsers(db *sql.DB) ([]*User, error) {
	rows, err := db.Query(
		`SELECT id, email, password, is_admin, created_at FROM users ORDER BY created_at ASC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []*User
	for rows.Next() {
		u := &User{}
		var adminVal int
		if err := rows.Scan(&u.ID, &u.Email, &u.Password, &adminVal, &u.CreatedAt); err != nil {
			return nil, err
		}
		u.IsAdmin = adminVal == 1
		users = append(users, u)
	}
	return users, rows.Err()
}

func UpdateUserAdmin(db *sql.DB, id int64, isAdmin bool) error {
	adminVal := 0
	if isAdmin {
		adminVal = 1
	}
	_, err := db.Exec(`UPDATE users SET is_admin = ? WHERE id = ?`, adminVal, id)
	return err
}

func DeleteUser(db *sql.DB, id int64) error {
	_, err := db.Exec(`DELETE FROM users WHERE id = ?`, id)
	return err
}

// --- Sessions ---

func CreateSession(db *sql.DB, userID int64, token string, expiresAt time.Time) error {
	_, err := db.Exec(
		`INSERT INTO sessions (user_id, token, expires_at) VALUES (?, ?, ?)`,
		userID, token, expiresAt,
	)
	return err
}

func GetSession(db *sql.DB, token string) (*User, error) {
	u := &User{}
	var adminVal int
	err := db.QueryRow(
		`SELECT u.id, u.email, u.password, u.is_admin, u.created_at
		 FROM sessions s JOIN users u ON s.user_id = u.id
		 WHERE s.token = ? AND s.expires_at > datetime('now')`,
		token,
	).Scan(&u.ID, &u.Email, &u.Password, &adminVal, &u.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	u.IsAdmin = adminVal == 1
	return u, err
}

func DeleteSession(db *sql.DB, token string) error {
	_, err := db.Exec(`DELETE FROM sessions WHERE token = ?`, token)
	return err
}

// --- Sites ---

func CreateSite(db *sql.DB, userID int64, slug, name, hashedPassword, plainPassword string) (*Site, error) {
	res, err := db.Exec(
		`INSERT INTO sites (user_id, slug, name, password, password_plain) VALUES (?, ?, ?, ?, ?)`,
		userID, slug, name, hashedPassword, plainPassword,
	)
	if err != nil {
		return nil, err
	}
	id, _ := res.LastInsertId()
	return &Site{ID: id, UserID: userID, Slug: slug, Name: name, Password: hashedPassword, PasswordPlain: plainPassword}, nil
}

func GetSiteBySlug(db *sql.DB, slug string) (*Site, error) {
	s := &Site{}
	err := db.QueryRow(
		`SELECT id, user_id, slug, name, password, password_plain, created_at, updated_at FROM sites WHERE slug = ?`,
		slug,
	).Scan(&s.ID, &s.UserID, &s.Slug, &s.Name, &s.Password, &s.PasswordPlain, &s.CreatedAt, &s.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return s, err
}

func GetSiteByID(db *sql.DB, id int64) (*Site, error) {
	s := &Site{}
	err := db.QueryRow(
		`SELECT id, user_id, slug, name, password, password_plain, created_at, updated_at FROM sites WHERE id = ?`,
		id,
	).Scan(&s.ID, &s.UserID, &s.Slug, &s.Name, &s.Password, &s.PasswordPlain, &s.CreatedAt, &s.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return s, err
}

func ListSitesByUser(db *sql.DB, userID int64) ([]*Site, error) {
	rows, err := db.Query(
		`SELECT id, user_id, slug, name, password, password_plain, created_at, updated_at FROM sites WHERE user_id = ? ORDER BY created_at DESC`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var sites []*Site
	for rows.Next() {
		s := &Site{}
		if err := rows.Scan(&s.ID, &s.UserID, &s.Slug, &s.Name, &s.Password, &s.PasswordPlain, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, err
		}
		sites = append(sites, s)
	}
	return sites, rows.Err()
}

func ListAllSites(db *sql.DB) ([]*Site, error) {
	rows, err := db.Query(
		`SELECT id, user_id, slug, name, password, password_plain, created_at, updated_at FROM sites ORDER BY created_at DESC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var sites []*Site
	for rows.Next() {
		s := &Site{}
		if err := rows.Scan(&s.ID, &s.UserID, &s.Slug, &s.Name, &s.Password, &s.PasswordPlain, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, err
		}
		sites = append(sites, s)
	}
	return sites, rows.Err()
}

func UpdateSite(db *sql.DB, id int64, name, hashedPassword, plainPassword string) error {
	_, err := db.Exec(
		`UPDATE sites SET name = ?, password = ?, password_plain = ?, updated_at = datetime('now') WHERE id = ?`,
		name, hashedPassword, plainPassword, id,
	)
	return err
}

func DeleteSite(db *sql.DB, id int64) error {
	_, err := db.Exec(`DELETE FROM sites WHERE id = ?`, id)
	return err
}

func CountSites(db *sql.DB) (int64, error) {
	var count int64
	err := db.QueryRow(`SELECT COUNT(*) FROM sites`).Scan(&count)
	return count, err
}

// --- Site Sessions ---

func CreateSiteSession(db *sql.DB, siteID int64, token string, expiresAt time.Time) error {
	_, err := db.Exec(
		`INSERT INTO site_sessions (site_id, token, expires_at) VALUES (?, ?, ?)`,
		siteID, token, expiresAt,
	)
	return err
}

func GetSiteSession(db *sql.DB, token string) (*Site, error) {
	s := &Site{}
	err := db.QueryRow(
		`SELECT s.id, s.user_id, s.slug, s.name, s.password, s.password_plain, s.created_at, s.updated_at
		 FROM site_sessions ss JOIN sites s ON ss.site_id = s.id
		 WHERE ss.token = ? AND ss.expires_at > datetime('now')`,
		token,
	).Scan(&s.ID, &s.UserID, &s.Slug, &s.Name, &s.Password, &s.PasswordPlain, &s.CreatedAt, &s.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return s, err
}

// --- Settings ---

func GetSetting(db *sql.DB, key string) (string, error) {
	var val string
	err := db.QueryRow(`SELECT value FROM settings WHERE key = ?`, key).Scan(&val)
	if err == sql.ErrNoRows {
		return "", nil
	}
	return val, err
}

func GetSettingBool(db *sql.DB, key string, defaultVal bool) bool {
	val, err := GetSetting(db, key)
	if err != nil || val == "" {
		return defaultVal
	}
	return val == "1" || val == "true"
}

func SetSetting(db *sql.DB, key, value string) error {
	_, err := db.Exec(
		`INSERT INTO settings (key, value) VALUES (?, ?)
		 ON CONFLICT(key) DO UPDATE SET value = excluded.value`,
		key, value,
	)
	return err
}

func GetAllSettings(db *sql.DB) (map[string]string, error) {
	rows, err := db.Query(`SELECT key, value FROM settings`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	m := make(map[string]string)
	for rows.Next() {
		var k, v string
		if err := rows.Scan(&k, &v); err != nil {
			return nil, err
		}
		m[k] = v
	}
	return m, rows.Err()
}

// GetStats returns dashboard statistics.
type Stats struct {
	Users  int64 `json:"users"`
	Sites  int64 `json:"sites"`
	Admins int64 `json:"admins"`
}

func GetStats(db *sql.DB) (*Stats, error) {
	s := &Stats{}
	if err := db.QueryRow(`SELECT COUNT(*) FROM users`).Scan(&s.Users); err != nil {
		return nil, err
	}
	if err := db.QueryRow(`SELECT COUNT(*) FROM sites`).Scan(&s.Sites); err != nil {
		return nil, err
	}
	if err := db.QueryRow(`SELECT COUNT(*) FROM users WHERE is_admin = 1`).Scan(&s.Admins); err != nil {
		return nil, err
	}
	return s, nil
}

// GetSettings returns a structured settings object.
type Settings struct {
	OpenRegistration bool `json:"openRegistration"`
	AllowPublicAccess bool `json:"allowPublicAccess"`
}

func GetSettings(db *sql.DB) (*Settings, error) {
	s := &Settings{}
	val, err := GetSetting(db, "open_registration")
	if err != nil {
		return nil, err
	}
	s.OpenRegistration = val == "1" || val == "true"
	val2, err := GetSetting(db, "allow_public_access")
	if err != nil {
		return nil, err
	}
	s.AllowPublicAccess = val2 == "1" || val2 == "true"
	return s, nil
}
