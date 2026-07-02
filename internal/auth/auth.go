package auth

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"time"

	"static-host/internal/db"

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

// GetSessionToken extracts the Bearer token from the Authorization header.
// Supports both "Authorization: Bearer <token>" and "Authorization: <token>".
func GetSessionToken(r *http.Request) string {
	auth := r.Header.Get("Authorization")
	if auth == "" {
		return ""
	}
	auth = strings.TrimSpace(auth)
	// Strip "Bearer " prefix if present
	if strings.HasPrefix(strings.ToLower(auth), "bearer ") {
		return strings.TrimSpace(auth[7:])
	}
	return auth
}

// GetSiteToken extracts the site access token from the Authorization header
// or from the ?token= query parameter (for browser navigation).
func GetSiteToken(r *http.Request) string {
	// Try Authorization header first
	if token := GetSessionToken(r); token != "" {
		return token
	}
	// Fall back to query parameter for browser navigation
	return r.URL.Query().Get("token")
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
