package server

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	crand "crypto/rand"

	"vibecast/internal/auth"
	"vibecast/internal/db"
)

// Config holds server configuration.
type Config struct {
	Addr       string // listen address, e.g. ":8080"
	StorageDir string // path to site files storage
	DBPath     string // path to SQLite database
	Version    string // build version (injected via ldflags)
}

// Server holds application state.
type Server struct {
	config     *Config
	database   *sql.DB
	version    string
	httpServer *http.Server // set by main.go for graceful shutdown / restart
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
	s := &Server{config: cfg, database: database, version: cfg.Version}

	// Start a background goroutine to clean up expired sessions every hour.
	go s.sessionCleanupLoop()

	return s, nil
}

// sessionCleanupLoop periodically deletes expired session rows.
func (s *Server) sessionCleanupLoop() {
	// Run once at startup
	_, _ = db.CleanupExpiredSessions(s.database)

	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()
	for range ticker.C {
		_, _ = db.CleanupExpiredSessions(s.database)
	}
}

// Close closes the database connection.
func (s *Server) Close() error {
	return s.database.Close()
}

// SetHTTPServer stores a reference to the running http.Server, used for
// graceful shutdown during restart.
func (s *Server) SetHTTPServer(hs *http.Server) {
	s.httpServer = hs
}

// Router returns the main HTTP handler with all routes registered.
func (s *Server) Router() http.Handler {
	mux := http.NewServeMux()

	// Static site serving: /s/{slug}/...
	mux.HandleFunc("/s/", s.staticHandler)

	// Password gate page: /p/{slug}
	mux.HandleFunc("/p/", s.passwordPageHandler)

	// API routes
	mux.HandleFunc("/api/auth/register", rateLimitMiddleware(s.handleRegister))
	mux.HandleFunc("/api/auth/login", rateLimitMiddleware(s.handleLogin))
	mux.HandleFunc("/api/auth/logout", s.handleLogout)
	mux.HandleFunc("/api/auth/me", s.handleMe)
	mux.HandleFunc("/api/auth/captcha", s.handleCaptcha)
	mux.HandleFunc("/api/settings", s.publicSettings)
	mux.HandleFunc("/api/version", s.handleVersion)
	mux.HandleFunc("/api/auth/change-password", auth.RequireAuth(s.database, s.handleChangePassword))

	// Sites API (auth required)
	mux.HandleFunc("/api/sites", auth.RequireAuth(s.database, s.handleSites))
	mux.HandleFunc("/api/sites/", auth.RequireAuth(s.database, s.handleSite))

	// Organization API (auth required)
	mux.HandleFunc("/api/org", auth.RequireAuth(s.database, s.handleOrg))
	mux.HandleFunc("/api/org/", auth.RequireAuth(s.database, s.handleOrgAction))

	// Admin API (admin only)
	mux.HandleFunc("/api/admin/stats", auth.RequireAdmin(s.database, s.adminStats))
	mux.HandleFunc("/api/admin/users", auth.RequireAdmin(s.database, s.adminListUsers))
	mux.HandleFunc("/api/admin/users/", auth.RequireAdmin(s.database, s.adminUserAction))
	mux.HandleFunc("/api/admin/sites", auth.RequireAdmin(s.database, s.adminListAllSites))
	mux.HandleFunc("/api/admin/sites/", auth.RequireAdmin(s.database, s.adminSiteAction))
	mux.HandleFunc("/api/admin/settings", auth.RequireAdmin(s.database, s.adminHandleSettings))
	mux.HandleFunc("/api/admin/cleanup", auth.RequireAdmin(s.database, s.adminCleanup))
	mux.HandleFunc("/api/admin/update/check", auth.RequireAdmin(s.database, s.adminCheckUpdate))
	mux.HandleFunc("/api/admin/update/apply", auth.RequireAdmin(s.database, s.adminApplyUpdate))
	mux.HandleFunc("/api/admin/update/status", auth.RequireAdmin(s.database, s.adminUpdateStatus))
	mux.HandleFunc("/api/admin/update/restart", auth.RequireAdmin(s.database, s.adminRestartUpdate))
	mux.HandleFunc("/api/admin/system-info", auth.RequireAdmin(s.database, s.adminSystemInfo))

	// Admin UI
	mux.HandleFunc("/admin", s.handleAdminPage)
	mux.HandleFunc("/admin/", s.handleAdminPage)

	// Dashboard UI
	mux.HandleFunc("/dashboard", s.handleDashboard)
	mux.HandleFunc("/dashboard/", s.handleDashboard)

	// Landing page
	mux.HandleFunc("/", s.handleIndex)

	return s.recoverMiddleware(s.logMiddleware(s.bodyLimitMiddleware(s.adminDomainMiddleware(mux))))
}

// adminDomainMiddleware blocks API and admin/dashboard routes from site content
// adminDomainMiddleware is a no-op placeholder kept for backward compatibility.
// Domain isolation is now handled only at /s/ and /p/ routes via isHostAllowedForSites.
func (s *Server) adminDomainMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}

// maxJSONBodySize limits JSON API request bodies to prevent memory exhaustion.
// Deploy endpoints already use their own MaxBytesReader for large uploads.
const maxJSONBodySize = 1 << 20 // 1 MB

// bodyLimitMiddleware wraps the handler with a request body size limit.
// Only applies to POST/PUT/PATCH with JSON content type — multipart uploads
// (deploy) are exempt since they use their own MaxBytesReader.
func (s *Server) bodyLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ct := r.Header.Get("Content-Type")
		// Only limit JSON/text requests, not multipart form uploads
		if strings.Contains(ct, "application/json") || strings.Contains(ct, "text/") {
			r.Body = http.MaxBytesReader(w, r.Body, maxJSONBodySize)
		}
		next.ServeHTTP(w, r)
	})
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

// handleVersion returns the build version.
func (s *Server) handleVersion(w http.ResponseWriter, r *http.Request) {
	v := s.version
	if v == "" {
		v = "dev"
	}
	writeJSON(w, 200, jsonResp{Data: map[string]string{"version": v}})
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

// generateUniqueSlug generates a random 12-character slug.
// Ignores the site name entirely — slugs are unguessable random strings.
func (s *Server) generateUniqueSlug(_ string) (string, error) {
	for i := 0; i < 30; i++ {
		slug := randomSuffix(12)
		ex, err := db.GetSiteBySlug(s.database, slug)
		if err != nil {
			return "", err
		}
		if ex == nil {
			return slug, nil
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
