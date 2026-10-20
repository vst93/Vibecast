package auth

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"time"

	"vibecast/internal/db"

	"golang.org/x/crypto/bcrypt"
)

const (
	SessionDuration = 7 * 24 * time.Hour
)

// GenerateToken generates a random 32-byte hex token.
func GenerateToken() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		panic(fmt.Sprintf("crypto/rand failed: %v", err))
	}
	return hex.EncodeToString(b)
}

// HashPassword bcrypts a plaintext password.
func HashPassword(plain string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	return string(b), err
}

// CheckPassword compares a plaintext password against a bcrypt hash.
func CheckPassword(hash, plain string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain)) == nil
}

// getTokenFromHeader extracts the Bearer token from the Authorization header only.
// Internal helper — does NOT check cookies.
func getTokenFromHeader(r *http.Request) string {
	auth := r.Header.Get("Authorization")
	if auth == "" {
		return ""
	}
	auth = strings.TrimSpace(auth)
	if strings.HasPrefix(strings.ToLower(auth), "bearer ") {
		return strings.TrimSpace(auth[7:])
	}
	return auth
}

// GetSessionToken extracts the user session token from the Authorization header, or
// falls back to the vibecast_session cookie (set on login for browser navigation).
func GetSessionToken(r *http.Request) string {
	if token := getTokenFromHeader(r); token != "" {
		return token
	}
	// Fall back to cookie (browser navigation to /s/ pages)
	if cookie, err := r.Cookie("vibecast_session"); err == nil && cookie.Value != "" {
		return cookie.Value
	}
	return ""
}

// SiteCookieName returns the cookie name that holds a site's access token.
//
// Cookies are scoped per site by NAME rather than by Path on purpose. Path has
// to be an absolute prefix, which the server cannot know when a reverse proxy
// strips the mount point (e.g. serving /v/s/{slug}/ as /s/{slug}/): a cookie
// scoped to "/s/{slug}/" is never sent back for "/v/s/{slug}/", so a correct
// password would bounce the visitor straight back to the gate. A per-site name
// with Path="/" works behind any mount, and still keeps one site from
// overwriting another's token.
func SiteCookieName(slug string) string {
	return "site_token_" + slug
}

// GetSiteToken extracts the site access token from the Authorization header or
// from the per-site cookie set by the password gate.
// NOTE: does NOT use vibecast_session cookie — that's a user login token, not a site token.
func GetSiteToken(r *http.Request, slug string) string {
	// Try Authorization header first (API clients)
	if token := getTokenFromHeader(r); token != "" {
		return token
	}
	// Per-site cookie (browser navigation — set by password gate)
	if c, err := r.Cookie(SiteCookieName(slug)); err == nil && c.Value != "" {
		return c.Value
	}
	// Legacy cookie from earlier versions. It is not site-scoped, so callers
	// must verify the resolved session belongs to the requested site before
	// granting access (both staticHandler and hasValidSiteSession do).
	if c, err := r.Cookie("site_token"); err == nil && c.Value != "" {
		return c.Value
	}
	return ""
}

// CurrentUser resolves the user from the Authorization header token. Returns nil if not authenticated.
func CurrentUser(r *http.Request, database *sql.DB) *db.User {
	token := GetSessionToken(r)
	if token == "" {
		return nil
	}
	user, err := db.GetSession(database, token)
	if err != nil || user == nil {
		return nil
	}
	return user
}

// RequireAuth is middleware that requires a valid user session via Bearer token.
func RequireAuth(database *sql.DB, next func(http.ResponseWriter, *http.Request, *db.User)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := CurrentUser(r, database)
		if user == nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, `{"error":"unauthorized"}`)
			return
		}
		next(w, r, user)
	}
}

// isTLSRequest returns true if the request was made over HTTPS (directly or
// behind a reverse proxy with X-Forwarded-Proto).
func isTLSRequest(r *http.Request) bool {
	if r.TLS != nil {
		return true
	}
	// Check X-Forwarded-Proto for reverse proxy setups
	if proto := r.Header.Get("X-Forwarded-Proto"); proto == "https" {
		return true
	}
	return false
}

// SetSessionCookie sets the vibecast_session cookie so browser navigation to /s/
// pages carries auth (for org_open bypass). Path "/" covers the entire site.
// When the request is over TLS, the Secure flag is set automatically.
func SetSessionCookie(w http.ResponseWriter, r *http.Request, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "vibecast_session",
		Value:    token,
		Path:     "/",
		MaxAge:   int(SessionDuration.Seconds()),
		HttpOnly: true,
		Secure:   isTLSRequest(r),
		SameSite: http.SameSiteLaxMode,
	})
}

// ClearSessionCookie expires the vibecast_session cookie.
func ClearSessionCookie(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "vibecast_session",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   isTLSRequest(r),
		SameSite: http.SameSiteLaxMode,
	})
}
