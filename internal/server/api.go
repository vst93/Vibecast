package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"vibecast/internal/auth"
	"vibecast/internal/db"
	"vibecast/internal/storage"
)

const defaultMaxUploadSize = 50 << 20 // 50 MB default, overridable via admin settings

// getMaxUploadSize reads the configured max upload size from DB settings.
func (s *Server) getMaxUploadSize() int64 {
	mb := db.GetSettingInt(s.database, "max_upload_size", 50)
	if mb < 1 {
		mb = 50
	}
	return int64(mb) << 20
}

type jsonResp struct {
	Error   string      `json:"error,omitempty"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

// --- Auth handlers ---

func (s *Server) handleRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, 405, jsonResp{Error: tMsg(r, "method_not_allowed")})
		return
	}
	var body struct {
		Email      string `json:"email"`
		Password   string `json:"password"`
		Confirm    string `json:"confirm"`
		CaptchaID  string `json:"captchaId"`
		CaptchaCode string `json:"captchaCode"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, 400, jsonResp{Error: tMsg(r, "invalid_json")})
		return
	}
	body.Email = strings.TrimSpace(strings.ToLower(body.Email))
	if body.Email == "" || len(body.Password) < 6 {
		writeJSON(w, 400, jsonResp{Error: tMsg(r, "email_password_invalid")})
		return
	}
	if len(body.Password) > 72 {
		writeJSON(w, 400, jsonResp{Error: tMsg(r, "password_too_long")})
		return
	}
	if body.Password != body.Confirm {
		writeJSON(w, 400, jsonResp{Error: tMsg(r, "password_mismatch")})
		return
	}
	// Verify captcha
	if !verifyCaptcha(body.CaptchaID, body.CaptchaCode) {
		writeJSON(w, 400, jsonResp{Error: tMsg(r, "captcha_incorrect")})
		return
	}

	// Check if registration is open (unless this is the first user — first user becomes admin)
	userCount, _ := db.CountUsers(s.database)
	isFirstUser := userCount == 0
	if !isFirstUser {
		if !db.GetSettingBool(s.database, "open_registration", true) {
			writeJSON(w, 403, jsonResp{Error: tMsg(r, "registration_closed")})
			return
		}
		// Check domain restriction
		if db.GetSettingBool(s.database, "domain_restriction", false) {
			allowedDomains, _ := db.GetSetting(s.database, "allowed_domains")
			if !isEmailDomainAllowed(body.Email, allowedDomains) {
				writeJSON(w, 403, jsonResp{Error: tMsg(r, "domain_not_allowed")})
				return
			}
		}
	}

	hashed, err := auth.HashPassword(body.Password)
	if err != nil {
		writeJSON(w, 500, jsonResp{Error: tMsg(r, "hash_failed")})
		return
	}
	user, err := db.CreateUser(s.database, body.Email, hashed, isFirstUser)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE") {
			writeJSON(w, 409, jsonResp{Error: tMsg(r, "email_taken")})
			return
		}
		writeJSON(w, 500, jsonResp{Error: tMsg(r, "create_user_failed")})
		return
	}

	// Auto-login: create session token
	token := auth.GenerateToken()
	expires := time.Now().Add(auth.SessionDuration)
	if err := db.CreateSession(s.database, user.ID, token, expires); err != nil {
		writeJSON(w, 500, jsonResp{Error: tMsg(r, "create_session_failed")})
		return
	}
	auth.SetSessionCookie(w, token)

	writeJSON(w, 201, jsonResp{
		Message: "registered",
		Data: map[string]interface{}{
			"id":      user.ID,
			"email":   user.Email,
			"isAdmin": user.IsAdmin,
			"token":   token,
		},
	})
}

func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, 405, jsonResp{Error: tMsg(r, "method_not_allowed")})
		return
	}
	var body struct {
		Email       string `json:"email"`
		Password    string `json:"password"`
		CaptchaID   string `json:"captchaId"`
		CaptchaCode string `json:"captchaCode"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, 400, jsonResp{Error: tMsg(r, "invalid_json")})
		return
	}
	body.Email = strings.TrimSpace(strings.ToLower(body.Email))
	// Verify captcha
	if !verifyCaptcha(body.CaptchaID, body.CaptchaCode) {
		writeJSON(w, 400, jsonResp{Error: tMsg(r, "captcha_incorrect")})
		return
	}

	user, err := db.GetUserByEmail(s.database, body.Email)
	if err != nil {
		writeJSON(w, 500, jsonResp{Error: tMsg(r, "internal_error")})
		return
	}
	if user == nil || !auth.CheckPassword(user.Password, body.Password) {
		writeJSON(w, 401, jsonResp{Error: tMsg(r, "invalid_credentials")})
		return
	}

	token := auth.GenerateToken()
	expires := time.Now().Add(auth.SessionDuration)
	if err := db.CreateSession(s.database, user.ID, token, expires); err != nil {
		writeJSON(w, 500, jsonResp{Error: tMsg(r, "create_session_failed")})
		return
	}
	auth.SetSessionCookie(w, token)

	writeJSON(w, 200, jsonResp{
		Message: "logged in",
		Data: map[string]interface{}{
			"id":      user.ID,
			"email":   user.Email,
			"isAdmin": user.IsAdmin,
			"token":   token,
		},
	})
}

func (s *Server) handleLogout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, 405, jsonResp{Error: tMsg(r, "method_not_allowed")})
		return
	}
	token := auth.GetSessionToken(r)
	if token != "" {
		_ = db.DeleteSession(s.database, token)
	}
	auth.ClearSessionCookie(w)
	writeJSON(w, 200, jsonResp{Message: "logged out"})
}

func (s *Server) handleMe(w http.ResponseWriter, r *http.Request) {
	user := auth.CurrentUser(r, s.database)
	if user == nil {
		writeJSON(w, 401, jsonResp{Error: tMsg(r, "unauthorized")})
		return
	}
	writeJSON(w, 200, jsonResp{
		Data: map[string]interface{}{
			"id":      user.ID,
			"email":   user.Email,
			"isAdmin": user.IsAdmin,
		},
	})
}

// handleCaptcha generates and returns an arithmetic captcha question.
func (s *Server) handleCaptcha(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, 405, jsonResp{Error: tMsg(r, "method_not_allowed")})
		return
	}
	id, svg := generateCaptcha()
	writeJSON(w, 200, jsonResp{
		Data: map[string]interface{}{
			"id":    id,
			"image": svg,
		},
	})
}

// handleChangePassword handles PUT /api/auth/change-password.
// Requires current password + new password. Also validates captcha.
func (s *Server) handleChangePassword(w http.ResponseWriter, r *http.Request, user *db.User) {
	if r.Method != http.MethodPut {
		writeJSON(w, 405, jsonResp{Error: tMsg(r, "method_not_allowed")})
		return
	}
	var body struct {
		OldPassword string `json:"oldPassword"`
		NewPassword string `json:"newPassword"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, 400, jsonResp{Error: tMsg(r, "invalid_json")})
		return
	}
	if len(body.NewPassword) < 6 {
		writeJSON(w, 400, jsonResp{Error: tMsg(r, "new_password_too_short")})
		return
	}
	if len(body.NewPassword) > 72 {
		writeJSON(w, 400, jsonResp{Error: tMsg(r, "password_too_long")})
		return
	}
	if !auth.CheckPassword(user.Password, body.OldPassword) {
		writeJSON(w, 403, jsonResp{Error: tMsg(r, "current_password_wrong")})
		return
	}
	hashed, err := auth.HashPassword(body.NewPassword)
	if err != nil {
		writeJSON(w, 500, jsonResp{Error: tMsg(r, "hash_failed")})
		return
	}
	if err := db.UpdateUserPassword(s.database, user.ID, hashed); err != nil {
		writeJSON(w, 500, jsonResp{Error: tMsg(r, "update_password_failed")})
		return
	}
	writeJSON(w, 200, jsonResp{Message: "password changed"})
}

// --- Sites API ---

func (s *Server) handleSites(w http.ResponseWriter, r *http.Request, user *db.User) {
	switch r.Method {
	case http.MethodGet:
		s.listSites(w, r, user)
	case http.MethodPost:
		s.createSite(w, r, user)
	default:
		writeJSON(w, 405, jsonResp{Error: tMsg(r, "method_not_allowed")})
	}
}

func (s *Server) handleSite(w http.ResponseWriter, r *http.Request, user *db.User) {
	// /api/sites/{id} or /api/sites/{id}/deploy
	pathParts := strings.SplitN(strings.TrimPrefix(r.URL.Path, "/api/sites/"), "/", 2)
	siteIDStr := pathParts[0]
	if siteIDStr == "" {
		writeJSON(w, 400, jsonResp{Error: tMsg(r, "site_id_required")})
		return
	}

	var siteID int64
	fmt.Sscanf(siteIDStr, "%d", &siteID)

	site, err := db.GetSiteByID(s.database, siteID)
	if err != nil || site == nil {
		writeJSON(w, 404, jsonResp{Error: tMsg(r, "site_not_found")})
		return
	}
	if site.UserID != user.ID {
		writeJSON(w, 403, jsonResp{Error: tMsg(r, "forbidden")})
		return
	}

	// Check for /deploy sub-action
	if len(pathParts) > 1 && pathParts[1] == "deploy" {
		if r.Method != http.MethodPost {
			writeJSON(w, 405, jsonResp{Error: tMsg(r, "method_not_allowed")})
			return
		}
		s.deploySite(w, r, user, site)
		return
	}

	// Check for /files sub-action (file tree listing)
	if len(pathParts) > 1 && pathParts[1] == "files" {
		if r.Method != http.MethodGet {
			writeJSON(w, 405, jsonResp{Error: tMsg(r, "method_not_allowed")})
			return
		}
		s.siteFileTree(w, r, site)
		return
	}

	switch r.Method {
	case http.MethodGet:
		writeJSON(w, 200, jsonResp{Data: s.siteToJSON(site)})
	case http.MethodPut:
		s.updateSite(w, r, user, site)
	case http.MethodDelete:
		s.deleteSite(w, r, user, site)
	default:
		writeJSON(w, 405, jsonResp{Error: tMsg(r, "method_not_allowed")})
	}
}

func (s *Server) listSites(w http.ResponseWriter, r *http.Request, user *db.User) {
	// Pagination
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	perPage, _ := strconv.Atoi(r.URL.Query().Get("perPage"))
	if perPage < 1 || perPage > 100 {
		perPage = 10
	}
	search := strings.TrimSpace(r.URL.Query().Get("q"))

	offset := (page - 1) * perPage
	sites, err := db.ListSitesByUserPaged(s.database, user.ID, search, perPage, offset)
	if err != nil {
		writeJSON(w, 500, jsonResp{Error: tMsg(r, "list_sites_failed")})
		return
	}
	total, err := db.CountSitesByUser(s.database, user.ID, search)
	if err != nil {
		writeJSON(w, 500, jsonResp{Error: tMsg(r, "count_sites_failed")})
		return
	}

	// Batch fetch visit stats for all sites on this page
	siteIDs := make([]int64, len(sites))
	for i, site := range sites {
		siteIDs[i] = site.ID
	}
	visitStats, _ := db.GetBatchVisitStats(s.database, siteIDs)

	var list []map[string]interface{}
	for _, site := range sites {
		vs := visitStats[site.ID]
		list = append(list, s.siteToJSONWithVisits(site, vs))
	}
	if list == nil {
		list = []map[string]interface{}{}
	}
	writeJSON(w, 200, jsonResp{Data: map[string]interface{}{
		"items":   list,
		"total":   total,
		"page":    page,
		"perPage": perPage,
	}})
}

func (s *Server) createSite(w http.ResponseWriter, r *http.Request, user *db.User) {
	var body struct {
		Name     string `json:"name"`
		Password string `json:"password"`
		OrgOpen  bool   `json:"orgOpen"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, 400, jsonResp{Error: tMsg(r, "invalid_json")})
		return
	}
	body.Name = strings.TrimSpace(body.Name)
	if body.Name == "" {
		writeJSON(w, 400, jsonResp{Error: tMsg(r, "name_required")})
		return
	}

	// If orgOpen is requested, user must be in an org
	if body.OrgOpen {
		org, err := db.GetUserOrganization(s.database, user.ID)
		if err != nil {
			writeJSON(w, 500, jsonResp{Error: tMsg(r, "internal_error")})
			return
		}
		if org == nil {
			writeJSON(w, 403, jsonResp{Error: tMsg(r, "org_open_requires_org")})
			return
		}
	}

	// Slug is always auto-generated — 12 random chars, no user input
	slug, err := s.generateUniqueSlug("")
	if err != nil {
		writeJSON(w, 500, jsonResp{Error: tMsg(r, "generate_slug_failed")})
		return
	}

	hashedPwd := ""
	if body.Password != "" {
		if len(body.Password) < 4 {
			writeJSON(w, 400, jsonResp{Error: tMsg(r, "site_password_too_short")})
			return
		}
		h, err := auth.HashPassword(body.Password)
		if err != nil {
			writeJSON(w, 500, jsonResp{Error: tMsg(r, "hash_failed")})
			return
		}
		hashedPwd = h
	}

	// If public access is disabled, sites must have password protection
	if body.Password == "" {
		if !db.GetSettingBool(s.database, "allow_public_access", true) {
			writeJSON(w, 403, jsonResp{Error: tMsg(r, "public_access_disabled")})
			return
		}
	}

	// Check site limit per user
	maxSites := db.GetSettingInt(s.database, "max_sites_per_user", 30)
	if maxSites > 0 {
		count, _ := db.CountSitesByUser(s.database, user.ID, "")
		if count >= int64(maxSites) {
			writeJSON(w, 403, jsonResp{Error: tMsg(r, "site_limit_reached")})
			return
		}
	}

	site, err := db.CreateSite(s.database, user.ID, slug, body.Name, hashedPwd, body.Password, body.OrgOpen)
	if err != nil {
		writeJSON(w, 500, jsonResp{Error: tMsg(r, "create_site_failed")})
		return
	}

	writeJSON(w, 201, jsonResp{
		Message: "site created",
		Data:    s.siteToJSON(site),
	})
}

func (s *Server) updateSite(w http.ResponseWriter, r *http.Request, user *db.User, site *db.Site) {
	var body struct {
		Name     string `json:"name"`
		Password string `json:"password"`
		OrgOpen  *bool  `json:"orgOpen"` // pointer to distinguish unset from false
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, 400, jsonResp{Error: tMsg(r, "invalid_json")})
		return
	}
	name := strings.TrimSpace(body.Name)
	if name == "" {
		name = site.Name
	}
	hashedPwd := site.Password // keep existing by default
	plainPwd := site.PasswordPlain
	if body.Password != "" {
		if len(body.Password) < 4 {
			writeJSON(w, 400, jsonResp{Error: tMsg(r, "site_password_too_short")})
			return
		}
		h, err := auth.HashPassword(body.Password)
		if err != nil {
			writeJSON(w, 500, jsonResp{Error: tMsg(r, "hash_failed")})
			return
		}
		hashedPwd = h
		plainPwd = body.Password
	} else if site.Password == "" {
		// Trying to keep a site public when public access is disabled
		if !db.GetSettingBool(s.database, "allow_public_access", true) {
			writeJSON(w, 403, jsonResp{Error: tMsg(r, "public_access_disabled")})
			return
		}
	}

	// Determine orgOpen value
	orgOpen := site.OrgOpen
	if body.OrgOpen != nil {
		orgOpen = *body.OrgOpen
		// If enabling orgOpen, user must be in an org
		if orgOpen {
			org, err := db.GetUserOrganization(s.database, user.ID)
			if err != nil {
				writeJSON(w, 500, jsonResp{Error: tMsg(r, "internal_error")})
				return
			}
			if org == nil {
				writeJSON(w, 403, jsonResp{Error: tMsg(r, "org_open_requires_org")})
				return
			}
		}
	}

	if err := db.UpdateSite(s.database, site.ID, name, hashedPwd, plainPwd, orgOpen); err != nil {
		writeJSON(w, 500, jsonResp{Error: tMsg(r, "update_site_failed")})
		return
	}
	writeJSON(w, 200, jsonResp{Message: "updated"})
}

func (s *Server) deleteSite(w http.ResponseWriter, r *http.Request, user *db.User, site *db.Site) {
	if err := db.DeleteSite(s.database, site.ID); err != nil {
		writeJSON(w, 500, jsonResp{Error: tMsg(r, "delete_site_failed")})
		return
	}
	_ = storage.DeleteSiteDir(s.config.StorageDir, site.Slug)
	writeJSON(w, 200, jsonResp{Message: "deleted"})
}

func (s *Server) deploySite(w http.ResponseWriter, r *http.Request, user *db.User, site *db.Site) {
	// Limit upload size (configurable via admin settings)
	maxSize := s.getMaxUploadSize()
	r.Body = http.MaxBytesReader(w, r.Body, maxSize)

	file, header, err := r.FormFile("file")
	if err != nil {
		if err.Error() == "http: request body too large" || strings.Contains(err.Error(), "too large") {
			writeJSON(w, 413, jsonResp{Error: tMsg(r, "file_too_large")})
			return
		}
		writeJSON(w, 400, jsonResp{Error: tMsg(r, "file_required")})
		return
	}
	defer file.Close()

	filename := strings.ToLower(header.Filename)
	isZip := strings.HasSuffix(filename, ".zip")
	siteDir := fmt.Sprintf("%s/%s", s.config.StorageDir, site.Slug)

	// Check single file size before processing
	if header.Size > maxSize {
		writeJSON(w, 413, jsonResp{Error: tMsg(r, "file_too_large")})
		return
	}

	if isZip {
		// ZIP deploy — extract and replace entire site
		data, err := io.ReadAll(file)
		if err != nil {
			if err.Error() == "http: request body too large" || strings.Contains(err.Error(), "too large") {
				writeJSON(w, 413, jsonResp{Error: tMsg(r, "file_too_large")})
				return
			}
			writeJSON(w, 400, jsonResp{Error: tMsg(r, "read_upload_failed")})
			return
		}

		result, err := storage.ExtractZip(bytesReader(data), int64(len(data)), siteDir, maxSize)
		if err != nil {
			errMsg := err.Error()
			if strings.Contains(errMsg, "file too large") && strings.Contains(errMsg, "limit") {
				// Extract the filename from the error message
				writeJSON(w, 413, jsonResp{Error: tMsg(r, "zip_file_too_large") + ": " + errMsg})
				return
			}
			if strings.Contains(errMsg, "zip bomb") {
				writeJSON(w, 413, jsonResp{Error: tMsg(r, "zip_bomb")})
				return
			}
			writeJSON(w, 500, jsonResp{Error: tMsg(r, "extract_zip_failed") + ": " + errMsg})
			return
		}

		respData := map[string]interface{}{
			"slug":     site.Slug,
			"url":      fmt.Sprintf("s/%s/", site.Slug),
			"files":    header.Filename,
			"fileCount": result.TotalFiles,
			"totalSize": result.TotalSize,
		}
		if len(result.Skipped) > 0 {
			respData["skipped"] = result.Skipped
		}

		writeJSON(w, 200, jsonResp{
			Message: "deployed",
			Data:    respData,
		})
		return
	}

	// Single file deploy — save file, replace entire site content
	if storage.IsBlockedExtension(filename) {
		writeJSON(w, 400, jsonResp{Error: tMsg(r, "file_type_blocked")})
		return
	}

	fileSize, err := storage.SaveSingleFile(file, header.Filename, siteDir, maxSize)
	if err != nil {
		writeJSON(w, 500, jsonResp{Error: tMsg(r, "save_file_failed") + ": " + err.Error()})
		return
	}

	writeJSON(w, 200, jsonResp{
		Message: "deployed",
		Data: map[string]interface{}{
			"slug":      site.Slug,
			"url":        fmt.Sprintf("s/%s/", site.Slug),
			"files":      header.Filename,
			"fileCount":  1,
			"totalSize":  fileSize,
		},
	})
}

func (s *Server) siteToJSON(site *db.Site) map[string]interface{} {
	protected := site.Password != ""
	publicAccessDisabled := !db.GetSettingBool(s.database, "allow_public_access", true)
	return map[string]interface{}{
		"id":                   site.ID,
		"slug":                 site.Slug,
		"name":                 site.Name,
		"protected":            protected,
		"password":             site.PasswordPlain,
		"storagePath":          site.Slug,
		"url":                  fmt.Sprintf("s/%s/", site.Slug),
		"createdAt":            site.CreatedAt,
		"updatedAt":            site.UpdatedAt,
		"publicAccessDisabled": publicAccessDisabled,
		"ownerEmail":           site.OwnerEmail,
		"orgOpen":              site.OrgOpen,
		"visits":               map[string]int64{"today": 0, "month": 0, "total": 0},
	}
}

// siteToJSONWithVisits builds the site JSON with visit stats included.
func (s *Server) siteToJSONWithVisits(site *db.Site, visits *db.VisitStats) map[string]interface{} {
	m := s.siteToJSON(site)
	if visits != nil {
		m["visits"] = map[string]int64{
			"today": visits.Today,
			"month": visits.Month,
			"total": visits.Total,
		}
	}
	return m
}

// bytesReader is a helper to avoid importing bytes in the import block above.
func bytesReader(data []byte) io.ReaderAt {
	return bytesReaderType{data}
}

// siteFileTree returns a flat list of files in the site's storage directory.
func (s *Server) siteFileTree(w http.ResponseWriter, r *http.Request, site *db.Site) {
	siteDir := filepath.Join(s.config.StorageDir, site.Slug)
	type fileEntry struct {
		Name string `json:"name"`
		Size int64  `json:"size"`
		Dir  bool   `json:"dir"`
	}
	var files []fileEntry
	entries, err := os.ReadDir(siteDir)
	if err != nil {
		writeJSON(w, 200, jsonResp{Data: []fileEntry{}})
		return
	}
	for _, e := range entries {
		if strings.HasPrefix(e.Name(), ".") {
			continue
		}
		info, _ := e.Info()
		size := int64(0)
		if info != nil {
			size = info.Size()
		}
		files = append(files, fileEntry{Name: e.Name(), Size: size, Dir: e.IsDir()})
	}
	writeJSON(w, 200, jsonResp{Data: files})
}

type bytesReaderType struct {
	data []byte
}

func (b bytesReaderType) ReadAt(p []byte, off int64) (int, error) {
	if off >= int64(len(b.data)) {
		return 0, io.EOF
	}
	n := copy(p, b.data[off:])
	return n, nil
}

// isEmailDomainAllowed checks if the email's domain is in the allowed list.
// allowedDomains is a comma/newline separated list of domains.
func isEmailDomainAllowed(email, allowedDomains string) bool {
	if allowedDomains == "" {
		return false // restriction is on but no domains configured → block all
	}
	at := strings.LastIndex(email, "@")
	if at < 0 {
		return false
	}
	domain := strings.ToLower(email[at+1:])
	for _, d := range strings.Split(allowedDomains, "\n") {
		d = strings.TrimSpace(strings.ToLower(d))
		if d == "" {
			continue
		}
		if d == domain {
			return true
		}
	}
	return false
}

// adminCleanup handles orphaned directory cleanup.
// GET  → scans storage dir, returns list of orphaned directories (not in DB)
// POST → deletes the orphaned directories (requires confirm:true in body)
func (s *Server) adminCleanup(w http.ResponseWriter, r *http.Request, user *db.User) {
	// Get all slugs from DB
	dbSlugs, err := db.GetAllSlugs(s.database)
	if err != nil {
		writeJSON(w, 500, jsonResp{Error: "failed to get site list"})
		return
	}
	slugSet := make(map[string]bool)
	for _, sl := range dbSlugs {
		slugSet[sl] = true
	}

	// Scan storage directory
	entries, err := os.ReadDir(s.config.StorageDir)
	if err != nil {
		writeJSON(w, 500, jsonResp{Error: "failed to read storage directory"})
		return
	}

	var orphans []string
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		name := e.Name()
		if strings.HasPrefix(name, ".") {
			continue
		}
		if !slugSet[name] {
			orphans = append(orphans, name)
		}
	}

	if r.Method == http.MethodGet {
		writeJSON(w, 200, jsonResp{
			Data: map[string]interface{}{
				"orphans": orphans,
				"count":   len(orphans),
			},
		})
		return
	}

	if r.Method == http.MethodPost {
		var body struct {
			Confirm bool `json:"confirm"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			writeJSON(w, 400, jsonResp{Error: tMsg(r, "invalid_json")})
			return
		}
		if !body.Confirm {
			writeJSON(w, 400, jsonResp{Error: "confirmation required"})
			return
		}

		deleted := 0
		for _, slug := range orphans {
			dir := filepath.Join(s.config.StorageDir, slug)
			if err := os.RemoveAll(dir); err != nil {
				continue
			}
			deleted++
		}

		writeJSON(w, 200, jsonResp{
			Message: "cleanup complete",
			Data: map[string]interface{}{
				"deleted": deleted,
				"total":   len(orphans),
			},
		})
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}
