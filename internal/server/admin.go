package server

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"static-host/internal/db"
)

// adminStats returns dashboard statistics.
func (s *Server) adminStats(w http.ResponseWriter, r *http.Request, user *db.User) {
	if r.Method != http.MethodGet {
		writeJSON(w, 405, jsonResp{Error: "method not allowed"})
		return
	}
	stats, err := db.GetStats(s.database)
	if err != nil {
		writeJSON(w, 500, jsonResp{Error: "failed to get stats"})
		return
	}
	writeJSON(w, 200, jsonResp{Data: stats})
}

// adminListUsers lists all users (GET) or toggles admin status (PUT /api/admin/users/{id}).
func (s *Server) adminListUsers(w http.ResponseWriter, r *http.Request, user *db.User) {
	if r.Method == http.MethodGet {
		users, err := db.ListUsers(s.database)
		if err != nil {
			writeJSON(w, 500, jsonResp{Error: "failed to list users"})
			return
		}
		var list []map[string]interface{}
		for _, u := range users {
			list = append(list, map[string]interface{}{
				"id":        u.ID,
				"email":     u.Email,
				"isAdmin":   u.IsAdmin,
				"createdAt": u.CreatedAt,
			})
		}
		if list == nil {
			list = []map[string]interface{}{}
		}
		writeJSON(w, 200, jsonResp{Data: list})
		return
	}
	// POST fallback — list users
	if r.Method == http.MethodPost {
		users, err := db.ListUsers(s.database)
		if err != nil {
			writeJSON(w, 500, jsonResp{Error: "failed to list users"})
			return
		}
		var list []map[string]interface{}
		for _, u := range users {
			list = append(list, map[string]interface{}{
				"id":        u.ID,
				"email":     u.Email,
				"isAdmin":   u.IsAdmin,
				"createdAt": u.CreatedAt,
			})
		}
		if list == nil {
			list = []map[string]interface{}{}
		}
		writeJSON(w, 200, jsonResp{Data: list})
		return
	}
	writeJSON(w, 405, jsonResp{Error: "method not allowed"})
}

// adminUserAction handles /api/admin/users/{id} — PUT to toggle admin, DELETE to remove.
func (s *Server) adminUserAction(w http.ResponseWriter, r *http.Request, user *db.User) {
	pathParts := strings.SplitN(strings.TrimPrefix(r.URL.Path, "/api/admin/users/"), "/", 2)
	targetID, err := strconv.ParseInt(pathParts[0], 10, 64)
	if err != nil {
		writeJSON(w, 400, jsonResp{Error: "invalid user ID"})
		return
	}

	target, err := db.GetUserByID(s.database, targetID)
	if err != nil || target == nil {
		writeJSON(w, 404, jsonResp{Error: "user not found"})
		return
	}

	switch r.Method {
	case http.MethodPut:
		// Toggle admin
		if target.ID == user.ID {
			writeJSON(w, 400, jsonResp{Error: "cannot modify your own admin status"})
			return
		}
		newVal := !target.IsAdmin
		if err := db.UpdateUserAdmin(s.database, target.ID, newVal); err != nil {
			writeJSON(w, 500, jsonResp{Error: "failed to update"})
			return
		}
		writeJSON(w, 200, jsonResp{Message: "updated", Data: map[string]interface{}{"isAdmin": newVal}})

	case http.MethodDelete:
		// Delete user
		if target.ID == user.ID {
			writeJSON(w, 400, jsonResp{Error: "cannot delete yourself"})
			return
		}
		// Delete user's sites from disk
		sites, _ := db.ListSitesByUser(s.database, target.ID)
		for _, site := range sites {
			_ = os.RemoveAll(filepath.Join(s.config.StorageDir, site.Slug))
		}
		if err := db.DeleteUser(s.database, target.ID); err != nil {
			writeJSON(w, 500, jsonResp{Error: "failed to delete user"})
			return
		}
		writeJSON(w, 200, jsonResp{Message: "deleted"})

	default:
		writeJSON(w, 405, jsonResp{Error: "method not allowed"})
	}
}

// adminListAllSites lists all sites (GET).
func (s *Server) adminListAllSites(w http.ResponseWriter, r *http.Request, user *db.User) {
	if r.Method != http.MethodGet {
		writeJSON(w, 405, jsonResp{Error: "method not allowed"})
		return
	}
	sites, err := db.ListAllSites(s.database)
	if err != nil {
		writeJSON(w, 500, jsonResp{Error: "failed to list sites"})
		return
	}
	var list []map[string]interface{}
	for _, site := range sites {
		list = append(list, siteToJSON(site))
	}
	if list == nil {
		list = []map[string]interface{}{}
	}
	writeJSON(w, 200, jsonResp{Data: list})
}

// adminSiteAction handles /api/admin/sites/{id} — DELETE to remove any site.
func (s *Server) adminSiteAction(w http.ResponseWriter, r *http.Request, user *db.User) {
	pathParts := strings.SplitN(strings.TrimPrefix(r.URL.Path, "/api/admin/sites/"), "/", 2)
	siteID, err := strconv.ParseInt(pathParts[0], 10, 64)
	if err != nil {
		writeJSON(w, 400, jsonResp{Error: "invalid site ID"})
		return
	}

	site, err := db.GetSiteByID(s.database, siteID)
	if err != nil || site == nil {
		writeJSON(w, 404, jsonResp{Error: "site not found"})
		return
	}

	switch r.Method {
	case http.MethodDelete:
		if err := db.DeleteSite(s.database, site.ID); err != nil {
			writeJSON(w, 500, jsonResp{Error: "failed to delete site"})
			return
		}
		_ = os.RemoveAll(filepath.Join(s.config.StorageDir, site.Slug))
		writeJSON(w, 200, jsonResp{Message: "deleted"})

	default:
		writeJSON(w, 405, jsonResp{Error: "method not allowed"})
	}
}

// adminHandleSettings handles both GET and PUT for settings.
func (s *Server) adminHandleSettings(w http.ResponseWriter, r *http.Request, user *db.User) {
	switch r.Method {
	case http.MethodGet:
		settings, err := db.GetSettings(s.database)
		if err != nil {
			writeJSON(w, 500, jsonResp{Error: "failed to get settings"})
			return
		}
		writeJSON(w, 200, jsonResp{Data: settings})
	case http.MethodPut:
		var body struct {
			OpenRegistration bool `json:"openRegistration"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			writeJSON(w, 400, jsonResp{Error: "invalid JSON"})
			return
		}
		if err := db.SetSetting(s.database, "open_registration", strconv.FormatBool(body.OpenRegistration)); err != nil {
			writeJSON(w, 500, jsonResp{Error: "failed to update settings"})
			return
		}
		writeJSON(w, 200, jsonResp{Message: "settings updated"})
	default:
		writeJSON(w, 405, jsonResp{Error: "method not allowed"})
	}
}

// handleAdminPage serves the admin dashboard SPA.
func (s *Server) handleAdminPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(adminPageHTML))
}
