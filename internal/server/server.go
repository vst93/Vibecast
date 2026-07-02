package server

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"strings"

	crand "crypto/rand"

	"static-host/internal/auth"
	"static-host/internal/db"
)

// Config holds server configuration.
type Config struct {
	Addr       string // listen address, e.g. ":8080"
	StorageDir string // path to site files storage
	DBPath     string // path to SQLite database
}

// Server holds application state.
type Server struct {
	config   *Config
	database *sql.DB
}

// New creates a new Server instance.
func New(cfg *Config) (*Server, error) {
	if err := os.MkdirAll(cfg.StorageDir, 0755); err != nil {
		return nil, fmt.Errorf("create storage dir: %w", err)
	}
	database, err := db.Open(cfg.DBPath)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}
	return &Server{config: cfg, database: database}, nil
}

// Close closes the database connection.
func (s *Server) Close() error {
	return s.database.Close()
}

// Router returns the main HTTP handler with all routes registered.
func (s *Server) Router() http.Handler {
	mux := http.NewServeMux()

	// Static site serving: /s/{slug}/...
	mux.HandleFunc("/s/", s.staticHandler)

	// Password gate page: /p/{slug}
	mux.HandleFunc("/p/", s.passwordPageHandler)

	// API routes
	mux.HandleFunc("/api/auth/register", s.handleRegister)
	mux.HandleFunc("/api/auth/login", s.handleLogin)
	mux.HandleFunc("/api/auth/logout", s.handleLogout)
	mux.HandleFunc("/api/auth/me", s.handleMe)
	mux.HandleFunc("/api/auth/captcha", s.handleCaptcha)
	mux.HandleFunc("/api/auth/change-password", auth.RequireAuth(s.database, s.handleChangePassword))

	// Sites API (auth required)
	mux.HandleFunc("/api/sites", auth.RequireAuth(s.database, s.handleSites))
	mux.HandleFunc("/api/sites/", auth.RequireAuth(s.database, s.handleSite))

	// Admin API (admin only)
	mux.HandleFunc("/api/admin/stats", auth.RequireAdmin(s.database, s.adminStats))
	mux.HandleFunc("/api/admin/users", auth.RequireAdmin(s.database, s.adminListUsers))
	mux.HandleFunc("/api/admin/users/", auth.RequireAdmin(s.database, s.adminUserAction))
	mux.HandleFunc("/api/admin/sites", auth.RequireAdmin(s.database, s.adminListAllSites))
	mux.HandleFunc("/api/admin/sites/", auth.RequireAdmin(s.database, s.adminSiteAction))
	mux.HandleFunc("/api/admin/settings", auth.RequireAdmin(s.database, s.adminHandleSettings))

	// Admin UI
	mux.HandleFunc("/admin", s.handleAdminPage)
	mux.HandleFunc("/admin/", s.handleAdminPage)

	// Dashboard UI
	mux.HandleFunc("/dashboard", s.handleDashboard)
	mux.HandleFunc("/dashboard/", s.handleDashboard)

	// Landing page
	mux.HandleFunc("/", s.handleIndex)

	return s.recoverMiddleware(s.logMiddleware(mux))
}

func (s *Server) recoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rv := recover(); rv != nil {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (s *Server) logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}

// handleIndex serves the landing page.
func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, landingPageHTML)
}

// handleDashboard serves the admin dashboard SPA.
func (s *Server) handleDashboard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, dashboardHTML)
}

// slugify converts a string to a URL-safe slug.
func slugify(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	s = strings.Map(func(r rune) rune {
		if r >= 'a' && r <= 'z' || r >= '0' && r <= '9' {
			return r
		}
		if r == '-' || r == '_' {
			return r
		}
		return '-'
	}, s)
	for strings.Contains(s, "--") {
		s = strings.ReplaceAll(s, "--", "-")
	}
	s = strings.Trim(s, "-")
	if s == "" {
		s = "site"
	}
	return s
}

// generateUniqueSlug ensures the slug is unique.
// Uses a random suffix instead of sequential numbers to avoid leaking site counts.
func (s *Server) generateUniqueSlug(base string) (string, error) {
	slug := slugify(base)
	existing, err := db.GetSiteBySlug(s.database, slug)
	if err != nil {
		return "", err
	}
	if existing == nil {
		return slug, nil
	}
	// Append random suffix — not sequential, so counts can't be guessed
	for i := 0; i < 20; i++ {
		suffix := randomSuffix(4)
		candidate := fmt.Sprintf("%s-%s", slugify(base), suffix)
		ex, err := db.GetSiteBySlug(s.database, candidate)
		if err != nil {
			return "", err
		}
		if ex == nil {
			return candidate, nil
		}
	}
	return "", fmt.Errorf("could not generate unique slug")
}

// randomSuffix generates a short random alphanumeric string.
func randomSuffix(n int) string {
	const chars = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = chars[randInt(len(chars))]
	}
	return string(b)
}

// randInt returns a non-negative random int < n.
func randInt(n int) int {
	return int(mrandUint32()) % n
}

// mrandUint32 returns a pseudo-random uint32 using crypto/rand.
func mrandUint32() uint32 {
	var b [4]byte
	_, _ = crand.Read(b[:])
	return uint32(b[0])<<24 | uint32(b[1])<<16 | uint32(b[2])<<8 | uint32(b[3])
}

// isValidSlug checks if a slug matches the allowed pattern.
func isValidSlug(slug string) bool {
	if len(slug) < 2 || len(slug) > 63 {
		return false
	}
	for _, r := range slug {
		if !(r >= 'a' && r <= 'z' || r >= '0' && r <= '9' || r == '-') {
			return false
		}
	}
	return true
}
