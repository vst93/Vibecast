package auth

import (
	"database/sql"
	"net/http"

	"vibecast/internal/db"
)

// RequireAdmin is middleware that requires a valid admin session.
func RequireAdmin(database *sql.DB, next func(http.ResponseWriter, *http.Request, *db.User)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := CurrentUser(r, database)
		if user == nil {
			http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
			return
		}
		if !user.IsAdmin {
			http.Error(w, `{"error":"forbidden: admin only"}`, http.StatusForbidden)
			return
		}
		next(w, r, user)
	}
}
