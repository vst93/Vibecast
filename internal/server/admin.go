package server

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"vibecast/internal/db"
)

// paginationParams extracts page, perPage, offset, and search from query params.
func paginationParams(r *http.Request) (page, perPage, offset int, search string) {
	page, _ = strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	perPage, _ = strconv.Atoi(r.URL.Query().Get("perPage"))
	if perPage < 1 || perPage > 100 {
		perPage = 10
	}
	offset = (page - 1) * perPage
	search = strings.TrimSpace(r.URL.Query().Get("q"))
	return
}

// normalizeSiteURL ensures the URL has a scheme (defaults to https).
func normalizeSiteURL(u string) string {
	u = strings.TrimSpace(u)
	if u == "" {
		return ""
	}
	if !strings.Contains(u, "://") {
		u = "https://" + u
	}
	return strings.TrimRight(u, "/")
}

// publicSettings returns non-sensitive settings for unauthenticated clients.
func (s *Server) publicSettings(w http.ResponseWriter, r *http.Request) {
	settings, err := db.GetSettings(s.database)
	if err != nil {
		writeJSON(w, 200, jsonResp{Data: map[string]interface{}{"openRegistration": true, "maxUploadSize": 50}})
		return
	}
	maxMB := db.GetSettingInt(s.database, "max_upload_size", 50)
	if maxMB < 1 {
		maxMB = 50
	}
	writeJSON(w, 200, jsonResp{Data: map[string]interface{}{
		"openRegistration": settings.OpenRegistration,
		"maxUploadSize":    maxMB,
		"maxSitesPerUser":  db.GetSettingInt(s.database, "max_sites_per_user", 30),
		"siteBaseUrl":      s.getSiteBaseURL(),
	}})
}

// adminStats returns dashboard statistics.
func (s *Server) adminStats(w http.ResponseWriter, r *http.Request, user *db.User) {
	if r.Method != http.MethodGet {
		writeJSON(w, 405, jsonResp{Error: tMsg(r, "method_not_allowed")})
		return
	}
	stats, err := db.GetStats(s.database)
	if err != nil {
		writeJSON(w, 500, jsonResp{Error: tMsg(r, "get_stats_failed")})
		return
	}
	writeJSON(w, 200, jsonResp{Data: stats})
}

// adminListUsers lists all users with pagination and search.
func (s *Server) adminListUsers(w http.ResponseWriter, r *http.Request, user *db.User) {
	if r.Method != http.MethodGet {
		writeJSON(w, 405, jsonResp{Error: tMsg(r, "method_not_allowed")})
		return
	}
	page, limit, offset, search := paginationParams(r)

	users, err := db.ListUsersPaged(s.database, search, limit, offset)
	if err != nil {
		writeJSON(w, 500, jsonResp{Error: tMsg(r, "list_users_failed")})
		return
	}
	total, err := db.CountUsersWithSearch(s.database, search)
	if err != nil {
		writeJSON(w, 500, jsonResp{Error: tMsg(r, "count_users_failed")})
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
	writeJSON(w, 200, jsonResp{Data: map[string]interface{}{
		"items": list,
		"total": total,
		"page":  page,
		"limit": limit,
	}})
}

// adminUserAction handles /api/admin/users/{id} — PUT to toggle admin, DELETE to remove.
func (s *Server) adminUserAction(w http.ResponseWriter, r *http.Request, user *db.User) {
	pathParts := strings.SplitN(strings.TrimPrefix(r.URL.Path, "/api/admin/users/"), "/", 2)
	targetID, err := strconv.ParseInt(pathParts[0], 10, 64)
	if err != nil {
		writeJSON(w, 400, jsonResp{Error: tMsg(r, "invalid_user_id")})
		return
	}

	target, err := db.GetUserByID(s.database, targetID)
	if err != nil || target == nil {
		writeJSON(w, 404, jsonResp{Error: tMsg(r, "user_not_found")})
		return
	}

	switch r.Method {
	case http.MethodPut:
		if target.ID == user.ID {
			writeJSON(w, 400, jsonResp{Error: tMsg(r, "cannot_modify_self_admin")})
			return
		}
		newVal := !target.IsAdmin
		if err := db.UpdateUserAdmin(s.database, target.ID, newVal); err != nil {
			writeJSON(w, 500, jsonResp{Error: tMsg(r, "update_failed")})
			return
		}
		writeJSON(w, 200, jsonResp{Message: "updated", Data: map[string]interface{}{"isAdmin": newVal}})

	case http.MethodDelete:
		if target.ID == user.ID {
			writeJSON(w, 400, jsonResp{Error: tMsg(r, "cannot_delete_self")})
			return
		}
		sites, _ := db.ListSitesByUser(s.database, target.ID)
		for _, site := range sites {
			_ = os.RemoveAll(filepath.Join(s.config.StorageDir, site.Slug))
		}
		if err := db.DeleteUser(s.database, target.ID); err != nil {
			writeJSON(w, 500, jsonResp{Error: tMsg(r, "delete_user_failed")})
			return
		}
		writeJSON(w, 200, jsonResp{Message: "deleted"})

	default:
		writeJSON(w, 405, jsonResp{Error: tMsg(r, "method_not_allowed")})
	}
}

// adminListAllSites lists all sites with pagination, search, and owner email.
func (s *Server) adminListAllSites(w http.ResponseWriter, r *http.Request, user *db.User) {
	if r.Method != http.MethodGet {
		writeJSON(w, 405, jsonResp{Error: tMsg(r, "method_not_allowed")})
		return
	}
	page, limit, offset, search := paginationParams(r)

	sites, err := db.ListAllSitesWithOwnerPaged(s.database, search, limit, offset)
	if err != nil {
		writeJSON(w, 500, jsonResp{Error: tMsg(r, "list_sites_failed")})
		return
	}
	total, err := db.CountAllSitesWithOwner(s.database, search)
	if err != nil {
		writeJSON(w, 500, jsonResp{Error: tMsg(r, "count_sites_failed")})
		return
	}
	var list []map[string]interface{}
	// Batch fetch visit stats
	siteIDs := make([]int64, len(sites))
	for i, site := range sites {
		siteIDs[i] = site.ID
	}
	visitStats, _ := db.GetBatchVisitStats(s.database, siteIDs)
	for _, site := range sites {
		vs := visitStats[site.ID]
		list = append(list, s.siteToJSONWithVisits(site, vs))
	}
	if list == nil {
		list = []map[string]interface{}{}
	}
	writeJSON(w, 200, jsonResp{Data: map[string]interface{}{
		"items": list,
		"total": total,
		"page":  page,
		"limit": limit,
	}})
}

// adminSiteAction handles /api/admin/sites/{id} — DELETE to remove any site.
func (s *Server) adminSiteAction(w http.ResponseWriter, r *http.Request, user *db.User) {
	pathParts := strings.SplitN(strings.TrimPrefix(r.URL.Path, "/api/admin/sites/"), "/", 2)
	siteID, err := strconv.ParseInt(pathParts[0], 10, 64)
	if err != nil {
		writeJSON(w, 400, jsonResp{Error: tMsg(r, "invalid_site_id")})
		return
	}

	site, err := db.GetSiteByID(s.database, siteID)
	if err != nil || site == nil {
		writeJSON(w, 404, jsonResp{Error: tMsg(r, "site_not_found")})
		return
	}

	// Check for /files sub-action (file tree listing for any site)
	if len(pathParts) > 1 && pathParts[1] == "files" {
		if r.Method != http.MethodGet {
			writeJSON(w, 405, jsonResp{Error: tMsg(r, "method_not_allowed")})
			return
		}
		s.siteFileTree(w, r, site)
		return
	}

	// Check for /password sub-action (admin can view any site's password)
	if len(pathParts) > 1 && pathParts[1] == "password" {
		if r.Method != http.MethodGet {
			writeJSON(w, 405, jsonResp{Error: tMsg(r, "method_not_allowed")})
			return
		}
		s.sitePassword(w, r, site)
		return
	}

	switch r.Method {
	case http.MethodDelete:
		if err := db.DeleteSite(s.database, site.ID); err != nil {
			writeJSON(w, 500, jsonResp{Error: tMsg(r, "delete_site_failed")})
			return
		}
		_ = os.RemoveAll(filepath.Join(s.config.StorageDir, site.Slug))
		writeJSON(w, 200, jsonResp{Message: "deleted"})

	default:
		writeJSON(w, 405, jsonResp{Error: tMsg(r, "method_not_allowed")})
	}
}

// adminHandleSettings handles both GET and PUT for settings.
func (s *Server) adminHandleSettings(w http.ResponseWriter, r *http.Request, user *db.User) {
	switch r.Method {
	case http.MethodGet:
		settings, err := db.GetSettings(s.database)
		if err != nil {
			writeJSON(w, 500, jsonResp{Error: tMsg(r, "get_settings_failed")})
			return
		}
		writeJSON(w, 200, jsonResp{Data: settings})
	case http.MethodPut:
		var body struct {
			OpenRegistration  bool   `json:"openRegistration"`
			AllowPublicAccess bool   `json:"allowPublicAccess"`
			DomainRestriction bool   `json:"domainRestriction"`
			AllowedDomains    string `json:"allowedDomains"`
			MaxUploadSize     int    `json:"maxUploadSize"`
			MaxSitesPerUser   int    `json:"maxSitesPerUser"`
			SiteBaseURL       string `json:"siteBaseUrl"`
			}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			writeJSON(w, 400, jsonResp{Error: tMsg(r, "invalid_json")})
			return
		}
		if err := db.SetSetting(s.database, "open_registration", strconv.FormatBool(body.OpenRegistration)); err != nil {
			writeJSON(w, 500, jsonResp{Error: tMsg(r, "update_settings_failed")})
			return
		}
		if err := db.SetSetting(s.database, "allow_public_access", strconv.FormatBool(body.AllowPublicAccess)); err != nil {
			writeJSON(w, 500, jsonResp{Error: tMsg(r, "update_settings_failed")})
			return
		}
		if err := db.SetSetting(s.database, "domain_restriction", strconv.FormatBool(body.DomainRestriction)); err != nil {
			writeJSON(w, 500, jsonResp{Error: tMsg(r, "update_settings_failed")})
			return
		}
		if err := db.SetSetting(s.database, "allowed_domains", body.AllowedDomains); err != nil {
			writeJSON(w, 500, jsonResp{Error: tMsg(r, "update_settings_failed")})
			return
		}
		if body.MaxUploadSize > 0 {
			if err := db.SetSetting(s.database, "max_upload_size", strconv.Itoa(body.MaxUploadSize)); err != nil {
				writeJSON(w, 500, jsonResp{Error: tMsg(r, "update_settings_failed")})
				return
			}
		}
		if body.MaxSitesPerUser >= 0 {
			if err := db.SetSetting(s.database, "max_sites_per_user", strconv.Itoa(body.MaxSitesPerUser)); err != nil {
				writeJSON(w, 500, jsonResp{Error: tMsg(r, "update_settings_failed")})
				return
			}
		}
		if err := db.SetSetting(s.database, "site_base_url", normalizeSiteURL(body.SiteBaseURL)); err != nil {
			writeJSON(w, 500, jsonResp{Error: tMsg(r, "update_settings_failed")})
			return
		}
		writeJSON(w, 200, jsonResp{Message: "settings updated"})
	default:
		writeJSON(w, 405, jsonResp{Error: tMsg(r, "method_not_allowed")})
	}
}

// handleAdminPage serves the admin dashboard SPA.
func (s *Server) handleAdminPage(w http.ResponseWriter, r *http.Request) {
	// Domain isolation: block admin panel from site content domains
	if s.isHostBlockedForAdmin(r) {
		http.Error(w, "Forbidden: admin panel is not accessible from this domain", http.StatusForbidden)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(adminPageHTML))
}
