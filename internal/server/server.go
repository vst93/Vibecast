package server

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"strings"

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

	// Sites API (auth required)
	mux.HandleFunc("/api/sites", auth.RequireAuth(s.database, s.handleSites))
	mux.HandleFunc("/api/sites/", auth.RequireAuth(s.database, s.handleSite))

	// Admin UI
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
func (s *Server) generateUniqueSlug(base string) (string, error) {
	slug := slugify(base)
	suffix := 1
	for {
		existing, err := db.GetSiteBySlug(s.database, slug)
		if err != nil {
			return "", err
		}
		if existing == nil {
			return slug, nil
		}
		suffix++
		slug = fmt.Sprintf("%s-%d", slugify(base), suffix)
	}
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
