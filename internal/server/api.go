package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"static-host/internal/auth"
	"static-host/internal/db"
	"static-host/internal/storage"
)

const maxUploadSize = 100 << 20 // 100 MB

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
		writeJSON(w, 405, jsonResp{Error: "method not allowed"})
		return
	}
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, 400, jsonResp{Error: "invalid JSON"})
		return
	}
	body.Email = strings.TrimSpace(strings.ToLower(body.Email))
	if body.Email == "" || len(body.Password) < 6 {
		writeJSON(w, 400, jsonResp{Error: "email required and password must be at least 6 characters"})
		return
	}
	if len(body.Password) > 72 {
		writeJSON(w, 400, jsonResp{Error: "password too long (max 72 chars)"})
		return
	}

	// Check if registration is open (unless this is the first user — first user becomes admin)
	userCount, _ := db.CountUsers(s.database)
	isFirstUser := userCount == 0
	if !isFirstUser {
		if !db.GetSettingBool(s.database, "open_registration", true) {
			writeJSON(w, 403, jsonResp{Error: "registration is closed"})
			return
		}
	}

	hashed, err := auth.HashPassword(body.Password)
	if err != nil {
		writeJSON(w, 500, jsonResp{Error: "failed to hash password"})
		return
	}
	user, err := db.CreateUser(s.database, body.Email, hashed, isFirstUser)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE") {
			writeJSON(w, 409, jsonResp{Error: "email already registered"})
			return
		}
		writeJSON(w, 500, jsonResp{Error: "failed to create user"})
		return
	}

	// Auto-login
	token := auth.GenerateToken()
	expires := time.Now().Add(auth.SessionDuration)
	if err := db.CreateSession(s.database, user.ID, token, expires); err != nil {
		writeJSON(w, 500, jsonResp{Error: "failed to create session"})
		return
	}
	auth.SetSessionCookie(w, token)

	writeJSON(w, 201, jsonResp{
		Message: "registered",
		Data: map[string]interface{}{
			"id":      user.ID,
			"email":   user.Email,
			"isAdmin": user.IsAdmin,
		},
	})
}

func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, 405, jsonResp{Error: "method not allowed"})
		return
	}
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, 400, jsonResp{Error: "invalid JSON"})
		return
	}
	body.Email = strings.TrimSpace(strings.ToLower(body.Email))

	user, err := db.GetUserByEmail(s.database, body.Email)
	if err != nil {
		writeJSON(w, 500, jsonResp{Error: "internal error"})
		return
	}
	if user == nil || !auth.CheckPassword(user.Password, body.Password) {
		writeJSON(w, 401, jsonResp{Error: "invalid email or password"})
		return
	}

	token := auth.GenerateToken()
	expires := time.Now().Add(auth.SessionDuration)
	if err := db.CreateSession(s.database, user.ID, token, expires); err != nil {
		writeJSON(w, 500, jsonResp{Error: "failed to create session"})
		return
	}
	auth.SetSessionCookie(w, token)

	writeJSON(w, 200, jsonResp{
		Message: "logged in",
		Data: map[string]interface{}{
			"id":      user.ID,
			"email":   user.Email,
			"isAdmin": user.IsAdmin,
		},
	})
}

func (s *Server) handleLogout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, 405, jsonResp{Error: "method not allowed"})
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
		writeJSON(w, 401, jsonResp{Error: "unauthorized"})
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

// --- Sites API ---

func (s *Server) handleSites(w http.ResponseWriter, r *http.Request, user *db.User) {
	switch r.Method {
	case http.MethodGet:
		s.listSites(w, r, user)
	case http.MethodPost:
		s.createSite(w, r, user)
	default:
		writeJSON(w, 405, jsonResp{Error: "method not allowed"})
	}
}

func (s *Server) handleSite(w http.ResponseWriter, r *http.Request, user *db.User) {
	// /api/sites/{id} or /api/sites/{id}/deploy
	pathParts := strings.SplitN(strings.TrimPrefix(r.URL.Path, "/api/sites/"), "/", 2)
	siteIDStr := pathParts[0]
	if siteIDStr == "" {
		writeJSON(w, 400, jsonResp{Error: "site ID required"})
		return
	}

	var siteID int64
	fmt.Sscanf(siteIDStr, "%d", &siteID)

	site, err := db.GetSiteByID(s.database, siteID)
	if err != nil || site == nil {
		writeJSON(w, 404, jsonResp{Error: "site not found"})
		return
	}
	if site.UserID != user.ID {
		writeJSON(w, 403, jsonResp{Error: "forbidden"})
		return
	}

	// Check for /deploy sub-action
	if len(pathParts) > 1 && pathParts[1] == "deploy" {
		if r.Method != http.MethodPost {
			writeJSON(w, 405, jsonResp{Error: "method not allowed"})
			return
		}
		s.deploySite(w, r, user, site)
		return
	}

	switch r.Method {
	case http.MethodGet:
		writeJSON(w, 200, jsonResp{Data: siteToJSON(site)})
	case http.MethodPut:
		s.updateSite(w, r, user, site)
	case http.MethodDelete:
		s.deleteSite(w, r, user, site)
	default:
		writeJSON(w, 405, jsonResp{Error: "method not allowed"})
	}
}

func (s *Server) listSites(w http.ResponseWriter, r *http.Request, user *db.User) {
	sites, err := db.ListSitesByUser(s.database, user.ID)
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

func (s *Server) createSite(w http.ResponseWriter, r *http.Request, user *db.User) {
	var body struct {
		Name     string `json:"name"`
		Slug     string `json:"slug"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, 400, jsonResp{Error: "invalid JSON"})
		return
	}
	body.Name = strings.TrimSpace(body.Name)
	if body.Name == "" {
		writeJSON(w, 400, jsonResp{Error: "name required"})
		return
	}

	slug := strings.TrimSpace(body.Slug)
	if slug == "" {
		var err error
		slug, err = s.generateUniqueSlug(body.Name)
		if err != nil {
			writeJSON(w, 500, jsonResp{Error: "failed to generate slug"})
			return
		}
	} else {
		slug = slugify(slug)
		if !isValidSlug(slug) {
			writeJSON(w, 400, jsonResp{Error: "invalid slug (2-63 chars, a-z0-9- only)"})
			return
		}
		existing, err := db.GetSiteBySlug(s.database, slug)
		if err != nil {
			writeJSON(w, 500, jsonResp{Error: "failed to check slug"})
			return
		}
		if existing != nil {
			writeJSON(w, 409, jsonResp{Error: "slug already taken"})
			return
		}
	}

	hashedPwd := ""
	if body.Password != "" {
		if len(body.Password) < 4 {
			writeJSON(w, 400, jsonResp{Error: "site password must be at least 4 characters"})
			return
		}
		h, err := auth.HashPassword(body.Password)
		if err != nil {
			writeJSON(w, 500, jsonResp{Error: "failed to hash password"})
			return
		}
		hashedPwd = h
	}

	site, err := db.CreateSite(s.database, user.ID, slug, body.Name, hashedPwd, body.Password)
	if err != nil {
		writeJSON(w, 500, jsonResp{Error: "failed to create site"})
		return
	}

	writeJSON(w, 201, jsonResp{
		Message: "site created",
		Data:    siteToJSON(site),
	})
}

func (s *Server) updateSite(w http.ResponseWriter, r *http.Request, user *db.User, site *db.Site) {
	var body struct {
		Name     string `json:"name"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, 400, jsonResp{Error: "invalid JSON"})
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
			writeJSON(w, 400, jsonResp{Error: "site password must be at least 4 characters"})
			return
		}
		h, err := auth.HashPassword(body.Password)
		if err != nil {
			writeJSON(w, 500, jsonResp{Error: "failed to hash password"})
			return
		}
		hashedPwd = h
		plainPwd = body.Password
	}
	if err := db.UpdateSite(s.database, site.ID, name, hashedPwd, plainPwd); err != nil {
		writeJSON(w, 500, jsonResp{Error: "failed to update site"})
		return
	}
	writeJSON(w, 200, jsonResp{Message: "updated"})
}

func (s *Server) deleteSite(w http.ResponseWriter, r *http.Request, user *db.User, site *db.Site) {
	if err := db.DeleteSite(s.database, site.ID); err != nil {
		writeJSON(w, 500, jsonResp{Error: "failed to delete site"})
		return
	}
	_ = storage.DeleteSiteDir(s.config.StorageDir, site.Slug)
	writeJSON(w, 200, jsonResp{Message: "deleted"})
}

func (s *Server) deploySite(w http.ResponseWriter, r *http.Request, user *db.User, site *db.Site) {
	// Limit upload size
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)

	file, header, err := r.FormFile("file")
	if err != nil {
		writeJSON(w, 400, jsonResp{Error: "file upload required (field name: 'file')"})
		return
	}
	defer file.Close()

	if !strings.HasSuffix(strings.ToLower(header.Filename), ".zip") {
		writeJSON(w, 400, jsonResp{Error: "only .zip files are accepted"})
		return
	}

	// Read file into memory (we need size for zip.NewReader)
	data, err := io.ReadAll(file)
	if err != nil {
		writeJSON(w, 400, jsonResp{Error: "failed to read upload"})
		return
	}

	siteDir := fmt.Sprintf("%s/%s", s.config.StorageDir, site.Slug)
	if _, err := storage.ExtractZip(bytesReader(data), int64(len(data)), siteDir); err != nil {
		writeJSON(w, 500, jsonResp{Error: fmt.Sprintf("failed to extract zip: %v", err)})
		return
	}

	writeJSON(w, 200, jsonResp{
		Message: "deployed",
		Data: map[string]interface{}{
			"slug":  site.Slug,
			"url":   fmt.Sprintf("/s/%s/", site.Slug),
			"files": header.Filename,
		},
	})
}

func siteToJSON(site *db.Site) map[string]interface{} {
	protected := site.Password != ""
	return map[string]interface{}{
		"id":            site.ID,
		"slug":          site.Slug,
		"name":          site.Name,
		"protected":     protected,
		"password":      site.PasswordPlain,
		"storagePath":   site.Slug,
		"url":           fmt.Sprintf("/s/%s/", site.Slug),
		"createdAt":     site.CreatedAt,
		"updatedAt":     site.UpdatedAt,
	}
}

// bytesReader is a helper to avoid importing bytes in the import block above.
func bytesReader(data []byte) io.ReaderAt {
	return bytesReaderType{data}
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
