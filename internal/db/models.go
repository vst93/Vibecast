package db

import (
	"database/sql"
	"strconv"
	"strings"
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
	OrgOpen       bool   // open to organization members
	OrgPinned     bool   // pinned to org list (visible to org members)
	OwnerEmail    string // only populated in admin views (JOIN with users)
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type Organization struct {
	ID         int64
	OwnerID    int64
	InviteCode string
	Name       string
	CreatedAt  time.Time
}

type OrgMember struct {
	ID      int64
	OrgID   int64
	UserID  int64
	Email   string // populated via JOIN
	IsOwner bool   // true if this member is the org owner
	JoinedAt time.Time
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

// ListUsersPaged returns a page of users with optional search on email.
func ListUsersPaged(db *sql.DB, search string, limit, offset int) ([]*User, error) {
	search = "%" + strings.ToLower(search) + "%"
	rows, err := db.Query(
		`SELECT id, email, password, is_admin, created_at
		 FROM users
		 WHERE LOWER(email) LIKE ?
		 ORDER BY created_at ASC
		 LIMIT ? OFFSET ?`,
		search, limit, offset,
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

// CountUsersWithSearch returns total user count with optional search.
func CountUsersWithSearch(db *sql.DB, search string) (int64, error) {
	search = "%" + strings.ToLower(search) + "%"
	var count int64
	err := db.QueryRow(
		`SELECT COUNT(*) FROM users WHERE LOWER(email) LIKE ?`,
		search,
	).Scan(&count)
	return count, err
}

func UpdateUserAdmin(db *sql.DB, id int64, isAdmin bool) error {
	adminVal := 0
	if isAdmin {
		adminVal = 1
	}
	_, err := db.Exec(`UPDATE users SET is_admin = ? WHERE id = ?`, adminVal, id)
	return err
}

func UpdateUserPassword(db *sql.DB, id int64, hashedPassword string) error {
	_, err := db.Exec(`UPDATE users SET password = ? WHERE id = ?`, hashedPassword, id)
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

// CleanupExpiredSessions deletes all expired user sessions and site sessions.
// Returns the total number of rows deleted.
func CleanupExpiredSessions(db *sql.DB) (int64, error) {
	res1, err := db.Exec(`DELETE FROM sessions WHERE expires_at <= datetime('now')`)
	if err != nil {
		return 0, err
	}
	n1, _ := res1.RowsAffected()
	res2, err := db.Exec(`DELETE FROM site_sessions WHERE expires_at <= datetime('now')`)
	if err != nil {
		return n1, nil
	}
	n2, _ := res2.RowsAffected()
	return n1 + n2, nil
}

// --- Sites ---

func CreateSite(db *sql.DB, userID int64, slug, name, hashedPassword, plainPassword string, orgOpen, orgPinned bool) (*Site, error) {
	orgVal := 0
	if orgOpen {
		orgVal = 1
	}
	pinnedVal := 0
	if orgPinned {
		pinnedVal = 1
	}
	res, err := db.Exec(
		`INSERT INTO sites (user_id, slug, name, password, password_plain, org_open, org_pinned) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		userID, slug, name, hashedPassword, plainPassword, orgVal, pinnedVal,
	)
	if err != nil {
		return nil, err
	}
	id, _ := res.LastInsertId()
	return &Site{ID: id, UserID: userID, Slug: slug, Name: name, Password: hashedPassword, PasswordPlain: plainPassword, OrgOpen: orgOpen, OrgPinned: orgPinned}, nil
}

func GetSiteBySlug(db *sql.DB, slug string) (*Site, error) {
	s := &Site{}
	var orgOpen, orgPinned int
	err := db.QueryRow(
		`SELECT id, user_id, slug, name, password, password_plain, org_open, org_pinned, created_at, updated_at FROM sites WHERE slug = ?`,
		slug,
	).Scan(&s.ID, &s.UserID, &s.Slug, &s.Name, &s.Password, &s.PasswordPlain, &orgOpen, &orgPinned, &s.CreatedAt, &s.UpdatedAt)
	s.OrgOpen = orgOpen == 1
	s.OrgPinned = orgPinned == 1
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return s, err
}

func GetSiteByID(db *sql.DB, id int64) (*Site, error) {
	s := &Site{}
	var orgOpen, orgPinned int
	err := db.QueryRow(
		`SELECT id, user_id, slug, name, password, password_plain, org_open, org_pinned, created_at, updated_at FROM sites WHERE id = ?`,
		id,
	).Scan(&s.ID, &s.UserID, &s.Slug, &s.Name, &s.Password, &s.PasswordPlain, &orgOpen, &orgPinned, &s.CreatedAt, &s.UpdatedAt)
	s.OrgOpen = orgOpen == 1
	s.OrgPinned = orgPinned == 1
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return s, err
}

func ListSitesByUser(db *sql.DB, userID int64) ([]*Site, error) {
	rows, err := db.Query(
		`SELECT id, user_id, slug, name, password, password_plain, org_open, org_pinned, created_at, updated_at FROM sites WHERE user_id = ? ORDER BY created_at DESC`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var sites []*Site
	for rows.Next() {
		s := &Site{}
		var orgOpen, orgPinned int
		if err := rows.Scan(&s.ID, &s.UserID, &s.Slug, &s.Name, &s.Password, &s.PasswordPlain, &orgOpen, &orgPinned, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, err
		}
		s.OrgOpen = orgOpen == 1
		s.OrgPinned = orgPinned == 1
		sites = append(sites, s)
	}
	return sites, rows.Err()
}

// ListSitesByUserPaged returns a page of sites for a user with optional search.
// search matches against name or slug (case-insensitive LIKE).
func ListSitesByUserPaged(db *sql.DB, userID int64, search string, limit, offset int) ([]*Site, error) {
	search = "%" + strings.ToLower(search) + "%"
	rows, err := db.Query(
		`SELECT id, user_id, slug, name, password, password_plain, org_open, org_pinned, created_at, updated_at
		 FROM sites
		 WHERE user_id = ? AND (LOWER(name) LIKE ? OR LOWER(slug) LIKE ?)
		 ORDER BY created_at DESC
		 LIMIT ? OFFSET ?`,
		userID, search, search, limit, offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var sites []*Site
	for rows.Next() {
		s := &Site{}
		var orgOpen, orgPinned int
		if err := rows.Scan(&s.ID, &s.UserID, &s.Slug, &s.Name, &s.Password, &s.PasswordPlain, &orgOpen, &orgPinned, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, err
		}
		s.OrgOpen = orgOpen == 1
		s.OrgPinned = orgPinned == 1
		sites = append(sites, s)
	}
	return sites, rows.Err()
}

// CountSitesByUser returns the total number of sites for a user (with optional search).
func CountSitesByUser(db *sql.DB, userID int64, search string) (int64, error) {
	search = "%" + strings.ToLower(search) + "%"
	var count int64
	err := db.QueryRow(
		`SELECT COUNT(*) FROM sites WHERE user_id = ? AND (LOWER(name) LIKE ? OR LOWER(slug) LIKE ?)`,
		userID, search, search,
	).Scan(&count)
	return count, err
}

func ListAllSites(db *sql.DB) ([]*Site, error) {
	rows, err := db.Query(
		`SELECT id, user_id, slug, name, password, password_plain, org_open, org_pinned, created_at, updated_at FROM sites ORDER BY created_at DESC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var sites []*Site
	for rows.Next() {
		s := &Site{}
		var orgOpen, orgPinned int
		if err := rows.Scan(&s.ID, &s.UserID, &s.Slug, &s.Name, &s.Password, &s.PasswordPlain, &orgOpen, &orgPinned, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, err
		}
		s.OrgOpen = orgOpen == 1
		s.OrgPinned = orgPinned == 1
		sites = append(sites, s)
	}
	return sites, rows.Err()
}

// ListAllSitesWithOwnerPaged returns a page of all sites with owner email and optional search.
// search matches against site name, slug, or owner email.
func ListAllSitesWithOwnerPaged(db *sql.DB, search string, limit, offset int) ([]*Site, error) {
	search = "%" + strings.ToLower(search) + "%"
	rows, err := db.Query(
		`SELECT s.id, s.user_id, s.slug, s.name, s.password, s.password_plain, s.org_open, s.org_pinned, u.email, s.created_at, s.updated_at
		 FROM sites s JOIN users u ON s.user_id = u.id
		 WHERE LOWER(s.name) LIKE ? OR LOWER(s.slug) LIKE ? OR LOWER(u.email) LIKE ?
		 ORDER BY s.created_at DESC
		 LIMIT ? OFFSET ?`,
		search, search, search, limit, offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var sites []*Site
	for rows.Next() {
		s := &Site{}
		var orgOpen, orgPinned int
		if err := rows.Scan(&s.ID, &s.UserID, &s.Slug, &s.Name, &s.Password, &s.PasswordPlain, &orgOpen, &orgPinned, &s.OwnerEmail, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, err
		}
		s.OrgOpen = orgOpen == 1
		s.OrgPinned = orgPinned == 1
		sites = append(sites, s)
	}
	return sites, rows.Err()
}

// CountAllSitesWithOwner returns total count of sites with optional search.
func CountAllSitesWithOwner(db *sql.DB, search string) (int64, error) {
	search = "%" + strings.ToLower(search) + "%"
	var count int64
	err := db.QueryRow(
		`SELECT COUNT(*)
		 FROM sites s JOIN users u ON s.user_id = u.id
		 WHERE LOWER(s.name) LIKE ? OR LOWER(s.slug) LIKE ? OR LOWER(u.email) LIKE ?`,
		search, search, search,
	).Scan(&count)
	return count, err
}

func UpdateSite(db *sql.DB, id int64, name, hashedPassword, plainPassword string, orgOpen, orgPinned bool) error {
	orgVal := 0
	if orgOpen {
		orgVal = 1
	}
	pinnedVal := 0
	if orgPinned {
		pinnedVal = 1
	}
	_, err := db.Exec(
		`UPDATE sites SET name = ?, password = ?, password_plain = ?, org_open = ?, org_pinned = ?, updated_at = datetime('now') WHERE id = ?`,
		name, hashedPassword, plainPassword, orgVal, pinnedVal, id,
	)
	return err
}

func DeleteSite(db *sql.DB, id int64) error {
	_, err := db.Exec(`DELETE FROM sites WHERE id = ?`, id)
	return err
}

// GetAllSlugs returns all site slugs from the database.
func GetAllSlugs(db *sql.DB) ([]string, error) {
	rows, err := db.Query(`SELECT slug FROM sites`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var slugs []string
	for rows.Next() {
		var s string
		if err := rows.Scan(&s); err != nil {
			return nil, err
		}
		slugs = append(slugs, s)
	}
	return slugs, rows.Err()
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
	var orgOpen, orgPinned int
	err := db.QueryRow(
		`SELECT s.id, s.user_id, s.slug, s.name, s.password, s.password_plain, s.org_open, s.org_pinned, s.created_at, s.updated_at
		 FROM site_sessions ss JOIN sites s ON ss.site_id = s.id
		 WHERE ss.token = ? AND ss.expires_at > datetime('now')`,
		token,
	).Scan(&s.ID, &s.UserID, &s.Slug, &s.Name, &s.Password, &s.PasswordPlain, &orgOpen, &orgPinned, &s.CreatedAt, &s.UpdatedAt)
	s.OrgOpen = orgOpen == 1
	s.OrgPinned = orgPinned == 1
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

// GetSettingInt returns an integer setting with a default value.
func GetSettingInt(db *sql.DB, key string, defaultVal int) int {
	val, err := GetSetting(db, key)
	if err != nil || val == "" {
		return defaultVal
	}
	n, err := strconv.Atoi(val)
	if err != nil {
		return defaultVal
	}
	return n
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
	OpenRegistration  bool   `json:"openRegistration"`
	AllowPublicAccess bool   `json:"allowPublicAccess"`
	DomainRestriction bool   `json:"domainRestriction"`
	AllowedDomains    string `json:"allowedDomains"`
	MaxUploadSize     int    `json:"maxUploadSize"`
	MaxSitesPerUser   int    `json:"maxSitesPerUser"`
	SiteAccessDomains string `json:"siteAccessDomains"`
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
	val3, err := GetSetting(db, "domain_restriction")
	if err != nil {
		return nil, err
	}
	s.DomainRestriction = val3 == "1" || val3 == "true"
	val4, err := GetSetting(db, "allowed_domains")
	if err != nil {
		return nil, err
	}
	s.AllowedDomains = val4
	val5, err := GetSetting(db, "max_upload_size")
	if err != nil {
		return nil, err
	}
	s.MaxUploadSize = 50 // default 50 MB
	if val5 != "" {
		if n, err := strconv.Atoi(val5); err == nil && n > 0 {
			s.MaxUploadSize = n
		}
	}
	val6, err := GetSetting(db, "max_sites_per_user")
	if err != nil {
		return nil, err
	}
	s.MaxSitesPerUser = 30 // default 30
	if val6 != "" {
		if n, err := strconv.Atoi(val6); err == nil && n >= 0 {
			s.MaxSitesPerUser = n
		}
	}
	val7, err := GetSetting(db, "site_access_domains")
	if err != nil {
		return nil, err
	}
	s.SiteAccessDomains = val7
	return s, nil
}

// --- Site Visits ---

// RecordVisit records a single site visit with daily and monthly granularity.
func RecordVisit(db *sql.DB, siteID int64, visitDate, visitMonth string) error {
	_, err := db.Exec(
		`INSERT INTO site_visits (site_id, visit_date, visit_month) VALUES (?, ?, ?)`,
		siteID, visitDate, visitMonth,
	)
	return err
}

// VisitStats holds visit counts for a site.
type VisitStats struct {
	Today int64 `json:"today"`
	Month int64 `json:"month"`
	Total int64 `json:"total"`
}

// GetVisitStats returns today/monthly/total visit counts for a site.
func GetVisitStats(db *sql.DB, siteID int64) (*VisitStats, error) {
	v := &VisitStats{}
	today := time.Now().Format("2006-01-02")
	month := time.Now().Format("2006-01")

	if err := db.QueryRow(
		`SELECT COUNT(*) FROM site_visits WHERE site_id = ? AND visit_date = ?`,
		siteID, today,
	).Scan(&v.Today); err != nil {
		return nil, err
	}
	if err := db.QueryRow(
		`SELECT COUNT(*) FROM site_visits WHERE site_id = ? AND visit_month = ?`,
		siteID, month,
	).Scan(&v.Month); err != nil {
		return nil, err
	}
	if err := db.QueryRow(
		`SELECT COUNT(*) FROM site_visits WHERE site_id = ?`,
		siteID,
	).Scan(&v.Total); err != nil {
		return nil, err
	}
	return v, nil
}

// GetBatchVisitStats returns visit stats for multiple site IDs in one query set.
// Returns map[siteID]*VisitStats.
func GetBatchVisitStats(db *sql.DB, siteIDs []int64) (map[int64]*VisitStats, error) {
	result := make(map[int64]*VisitStats)
	if len(siteIDs) == 0 {
		return result, nil
	}

	today := time.Now().Format("2006-01-02")
	month := time.Now().Format("2006-01")

	// Build placeholders
	placeholders := make([]string, len(siteIDs))
	args := make([]interface{}, len(siteIDs))
	for i, id := range siteIDs {
		placeholders[i] = "?"
		args[i] = id
	}
	phStr := strings.Join(placeholders, ",")

	// Total counts
	totalArgs := make([]interface{}, len(siteIDs))
	copy(totalArgs, args)
	rows, err := db.Query(
		`SELECT site_id, COUNT(*) FROM site_visits WHERE site_id IN (`+phStr+`) GROUP BY site_id`,
		totalArgs...,
	)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var id int64
		var cnt int64
		if err := rows.Scan(&id, &cnt); err != nil {
			rows.Close()
			return nil, err
		}
		if result[id] == nil {
			result[id] = &VisitStats{}
		}
		result[id].Total = cnt
	}
	rows.Close()

	// Today counts
	todayArgs := append([]interface{}{today}, args...)
	rows2, err := db.Query(
		`SELECT site_id, COUNT(*) FROM site_visits WHERE visit_date = ? AND site_id IN (`+phStr+`) GROUP BY site_id`,
		todayArgs...,
	)
	if err != nil {
		return nil, err
	}
	for rows2.Next() {
		var id int64
		var cnt int64
		if err := rows2.Scan(&id, &cnt); err != nil {
			rows2.Close()
			return nil, err
		}
		if result[id] == nil {
			result[id] = &VisitStats{}
		}
		result[id].Today = cnt
	}
	rows2.Close()

	// Month counts
	monthArgs := append([]interface{}{month}, args...)
	rows3, err := db.Query(
		`SELECT site_id, COUNT(*) FROM site_visits WHERE visit_month = ? AND site_id IN (`+phStr+`) GROUP BY site_id`,
		monthArgs...,
	)
	if err != nil {
		return nil, err
	}
	for rows3.Next() {
		var id int64
		var cnt int64
		if err := rows3.Scan(&id, &cnt); err != nil {
			rows3.Close()
			return nil, err
		}
		if result[id] == nil {
			result[id] = &VisitStats{}
		}
		result[id].Month = cnt
	}
	rows3.Close()

	// Ensure all requested IDs have an entry
	for _, id := range siteIDs {
		if result[id] == nil {
			result[id] = &VisitStats{}
		}
	}

	return result, nil
}

// --- Organizations ---

func CreateOrganization(db *sql.DB, ownerID int64, name, inviteCode string) (*Organization, error) {
	res, err := db.Exec(
		`INSERT INTO organizations (owner_id, invite_code, name) VALUES (?, ?, ?)`,
		ownerID, inviteCode, name,
	)
	if err != nil {
		return nil, err
	}
	id, _ := res.LastInsertId()
	// Owner is also a member
	_, _ = db.Exec(`INSERT OR IGNORE INTO org_members (org_id, user_id) VALUES (?, ?)`, id, ownerID)
	return &Organization{ID: id, OwnerID: ownerID, InviteCode: inviteCode, Name: name}, nil
}

func GetOrganizationByID(db *sql.DB, id int64) (*Organization, error) {
	o := &Organization{}
	err := db.QueryRow(
		`SELECT id, owner_id, invite_code, name, created_at FROM organizations WHERE id = ?`,
		id,
	).Scan(&o.ID, &o.OwnerID, &o.InviteCode, &o.Name, &o.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return o, err
}

func GetOrganizationByInviteCode(db *sql.DB, code string) (*Organization, error) {
	o := &Organization{}
	err := db.QueryRow(
		`SELECT id, owner_id, invite_code, name, created_at FROM organizations WHERE invite_code = ?`,
		code,
	).Scan(&o.ID, &o.OwnerID, &o.InviteCode, &o.Name, &o.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return o, err
}

func GetOrganizationByOwner(db *sql.DB, ownerID int64) (*Organization, error) {
	o := &Organization{}
	err := db.QueryRow(
		`SELECT id, owner_id, invite_code, name, created_at FROM organizations WHERE owner_id = ?`,
		ownerID,
	).Scan(&o.ID, &o.OwnerID, &o.InviteCode, &o.Name, &o.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return o, err
}

func DeleteOrganization(db *sql.DB, id int64) error {
	_, err := db.Exec(`DELETE FROM organizations WHERE id = ?`, id)
	return err
}

// GetUserOrganization returns the organization the user belongs to (as member or owner), or nil.
func GetUserOrganization(db *sql.DB, userID int64) (*Organization, error) {
	o := &Organization{}
	err := db.QueryRow(
		`SELECT o.id, o.owner_id, o.invite_code, o.name, o.created_at
		 FROM organizations o
		 LEFT JOIN org_members m ON m.org_id = o.id
		 WHERE o.owner_id = ? OR m.user_id = ?
		 LIMIT 1`,
		userID, userID,
	).Scan(&o.ID, &o.OwnerID, &o.InviteCode, &o.Name, &o.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return o, err
}

// IsOrgMember checks if userID is a member of orgID (owner counts as member).
func IsOrgMember(db *sql.DB, orgID, userID int64) (bool, error) {
	// Check if owner
	var oOrgID int64
	err := db.QueryRow(`SELECT id FROM organizations WHERE id = ? AND owner_id = ?`, orgID, userID).Scan(&oOrgID)
	if err == nil {
		return true, nil
	}
	if err != sql.ErrNoRows {
		return false, err
	}
	// Check if member
	var id int64
	err = db.QueryRow(`SELECT id FROM org_members WHERE org_id = ? AND user_id = ?`, orgID, userID).Scan(&id)
	if err == sql.ErrNoRows {
		return false, nil
	}
	return err == nil, err
}

func JoinOrganization(db *sql.DB, orgID, userID int64) error {
	_, err := db.Exec(`INSERT OR IGNORE INTO org_members (org_id, user_id) VALUES (?, ?)`, orgID, userID)
	return err
}

func LeaveOrganization(db *sql.DB, orgID, userID int64) error {
	_, err := db.Exec(`DELETE FROM org_members WHERE org_id = ? AND user_id = ?`, orgID, userID)
	return err
}

// RemoveOrgMember removes a member from an organization (used by org owner).
func RemoveOrgMember(db *sql.DB, orgID, userID int64) error {
	_, err := db.Exec(`DELETE FROM org_members WHERE org_id = ? AND user_id = ?`, orgID, userID)
	return err
}

// CountOrgMembers returns the number of non-owner members in an org.
func CountOrgMembers(db *sql.DB, orgID int64) (int64, error) {
	var count int64
	err := db.QueryRow(`SELECT COUNT(*) FROM org_members WHERE org_id = ?`, orgID).Scan(&count)
	return count, err
}

// ListOrgMembersPaged returns a page of org members with optional search.
// Owner is always first in the list.
func ListOrgMembersPaged(db *sql.DB, orgID int64, search string, limit, offset int) ([]*OrgMember, error) {
	search = "%" + strings.ToLower(search) + "%"
	rows, err := db.Query(
		`SELECT m.id, m.org_id, m.user_id, u.email, (o.owner_id = m.user_id) AS is_owner, m.joined_at
		 FROM org_members m
		 JOIN users u ON m.user_id = u.id
		 JOIN organizations o ON m.org_id = o.id
		 WHERE m.org_id = ? AND LOWER(u.email) LIKE ?
		 ORDER BY is_owner DESC, m.joined_at ASC
		 LIMIT ? OFFSET ?`,
		orgID, search, limit, offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var members []*OrgMember
	for rows.Next() {
		m := &OrgMember{}
		var isOwner int
		if err := rows.Scan(&m.ID, &m.OrgID, &m.UserID, &m.Email, &isOwner, &m.JoinedAt); err != nil {
			return nil, err
		}
		m.IsOwner = isOwner == 1
		members = append(members, m)
	}
	return members, rows.Err()
}

// CountOrgMembersWithSearch returns total member count with optional search.
func CountOrgMembersWithSearch(db *sql.DB, orgID int64, search string) (int64, error) {
	search = "%" + strings.ToLower(search) + "%"
	var count int64
	err := db.QueryRow(
		`SELECT COUNT(*)
		 FROM org_members m
		 JOIN users u ON m.user_id = u.id
		 WHERE m.org_id = ? AND LOWER(u.email) LIKE ?`,
		orgID, search,
	).Scan(&count)
	return count, err
}

// ListPinnedOrgSitesPaged returns sites pinned to the org that belong to org members.
// Only returns name and slug - no password/config info.
func ListPinnedOrgSitesPaged(db *sql.DB, orgID int64, search string, limit, offset int) ([]*Site, error) {
	search = "%" + strings.ToLower(search) + "%"
	rows, err := db.Query(
		`SELECT s.id, s.slug, s.name
		 FROM sites s
		 JOIN org_members m ON s.user_id = m.user_id AND m.org_id = ?
		 WHERE s.org_pinned = 1 AND LOWER(s.name) LIKE ?
		 ORDER BY s.updated_at DESC
		 LIMIT ? OFFSET ?`,
		orgID, search, limit, offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var sites []*Site
	for rows.Next() {
		s := &Site{}
		if err := rows.Scan(&s.ID, &s.Slug, &s.Name); err != nil {
			return nil, err
		}
		sites = append(sites, s)
	}
	return sites, rows.Err()
}

// CountPinnedOrgSites returns total count of pinned sites in an org.
func CountPinnedOrgSites(db *sql.DB, orgID int64, search string) (int64, error) {
	search = "%" + strings.ToLower(search) + "%"
	var count int64
	err := db.QueryRow(
		`SELECT COUNT(*)
		 FROM sites s
		 JOIN org_members m ON s.user_id = m.user_id AND m.org_id = ?
		 WHERE s.org_pinned = 1 AND LOWER(s.name) LIKE ?`,
		orgID, search,
	).Scan(&count)
	return count, err
}
