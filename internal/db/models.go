package db

import (
	"database/sql"
	"time"
)

type User struct {
	ID        int64
	Email     string
	Password  string // bcrypt hash
	CreatedAt time.Time
}

type Site struct {
	ID        int64
	UserID    int64
	Slug      string
	Name      string
	Password  string // bcrypt hash, empty = no protection
	CreatedAt time.Time
	UpdatedAt time.Time
}

// CreateUser inserts a new user. Password should already be bcrypt-hashed.
func CreateUser(db *sql.DB, email, hashedPassword string) (*User, error) {
	res, err := db.Exec(
		`INSERT INTO users (email, password) VALUES (?, ?)`,
		email, hashedPassword,
	)
	if err != nil {
		return nil, err
	}
	id, _ := res.LastInsertId()
	return &User{ID: id, Email: email, Password: hashedPassword}, nil
}

func GetUserByEmail(db *sql.DB, email string) (*User, error) {
	u := &User{}
	err := db.QueryRow(
		`SELECT id, email, password, created_at FROM users WHERE email = ?`,
		email,
	).Scan(&u.ID, &u.Email, &u.Password, &u.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return u, err
}

func GetUserByID(db *sql.DB, id int64) (*User, error) {
	u := &User{}
	err := db.QueryRow(
		`SELECT id, email, password, created_at FROM users WHERE id = ?`,
		id,
	).Scan(&u.ID, &u.Email, &u.Password, &u.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return u, err
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
	err := db.QueryRow(
		`SELECT u.id, u.email, u.password, u.created_at
		 FROM sessions s JOIN users u ON s.user_id = u.id
		 WHERE s.token = ? AND s.expires_at > datetime('now')`,
		token,
	).Scan(&u.ID, &u.Email, &u.Password, &u.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return u, err
}

func DeleteSession(db *sql.DB, token string) error {
	_, err := db.Exec(`DELETE FROM sessions WHERE token = ?`, token)
	return err
}

// --- Sites ---

func CreateSite(db *sql.DB, userID int64, slug, name, hashedPassword string) (*Site, error) {
	res, err := db.Exec(
		`INSERT INTO sites (user_id, slug, name, password) VALUES (?, ?, ?, ?)`,
		userID, slug, name, hashedPassword,
	)
	if err != nil {
		return nil, err
	}
	id, _ := res.LastInsertId()
	return &Site{ID: id, UserID: userID, Slug: slug, Name: name, Password: hashedPassword}, nil
}

func GetSiteBySlug(db *sql.DB, slug string) (*Site, error) {
	s := &Site{}
	err := db.QueryRow(
		`SELECT id, user_id, slug, name, password, created_at, updated_at FROM sites WHERE slug = ?`,
		slug,
	).Scan(&s.ID, &s.UserID, &s.Slug, &s.Name, &s.Password, &s.CreatedAt, &s.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return s, err
}

func GetSiteByID(db *sql.DB, id int64) (*Site, error) {
	s := &Site{}
	err := db.QueryRow(
		`SELECT id, user_id, slug, name, password, created_at, updated_at FROM sites WHERE id = ?`,
		id,
	).Scan(&s.ID, &s.UserID, &s.Slug, &s.Name, &s.Password, &s.CreatedAt, &s.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return s, err
}

func ListSitesByUser(db *sql.DB, userID int64) ([]*Site, error) {
	rows, err := db.Query(
		`SELECT id, user_id, slug, name, password, created_at, updated_at FROM sites WHERE user_id = ? ORDER BY created_at DESC`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sites []*Site
	for rows.Next() {
		s := &Site{}
		if err := rows.Scan(&s.ID, &s.UserID, &s.Slug, &s.Name, &s.Password, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, err
		}
		sites = append(sites, s)
	}
	return sites, rows.Err()
}

func UpdateSite(db *sql.DB, id int64, name, hashedPassword string) error {
	_, err := db.Exec(
		`UPDATE sites SET name = ?, password = ?, updated_at = datetime('now') WHERE id = ?`,
		name, hashedPassword, id,
	)
	return err
}

func DeleteSite(db *sql.DB, id int64) error {
	_, err := db.Exec(`DELETE FROM sites WHERE id = ?`, id)
	return err
}

// --- Site Sessions (for password-protected sites) ---

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
		`SELECT s.id, s.user_id, s.slug, s.name, s.password, s.created_at, s.updated_at
		 FROM site_sessions ss JOIN sites s ON ss.site_id = s.id
		 WHERE ss.token = ? AND ss.expires_at > datetime('now')`,
		token,
	).Scan(&s.ID, &s.UserID, &s.Slug, &s.Name, &s.Password, &s.CreatedAt, &s.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return s, err
}
