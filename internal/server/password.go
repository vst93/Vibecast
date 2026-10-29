package server

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"vibecast/internal/auth"
	"vibecast/internal/db"
)

// passwordPageHandler shows the password gate page and processes password submissions.
// GET  /p/{slug}  → show password form
// POST /p/{slug}  → validate password, return token (JSON) or redirect with cookie
//
// A ?next= parameter carries the sub-path the visitor originally requested, so
// deep links into a protected site resume at the right page after unlocking.
func (s *Server) passwordPageHandler(w http.ResponseWriter, r *http.Request) {
	// Domain isolation: if site_access_domains is configured, only allow
	// password gate to be accessed from those domains.
	if !s.isHostAllowedForSites(r) {
		writeVisitorError(w, http.StatusForbidden, "This address does not serve site content",
			"Uploaded sites are served from a separate domain. Open the link you were given, or ask the owner for it.")
		return
	}

	slug := strings.TrimPrefix(r.URL.Path, "/p/")
	slug = strings.SplitN(slug, "/", 2)[0]
	if slug == "" {
		http.NotFound(w, r)
		return
	}

	site, err := db.GetSiteBySlug(s.database, slug)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}
	if site == nil {
		http.NotFound(w, r)
		return
	}

	next := sanitizeSiteReturnTarget(r.URL.Query().Get("next"))
	// Redirects climb back to the mounted root relatively, so they keep working
	// behind a reverse-proxy prefix the server cannot see.
	contentURL := relativeSiteContentURL(r.URL.Path, slug, next)

	if site.Password == "" {
		redirectSeeOther(w, contentURL)
		return
	}

	// Owners and authorized organization members do not need the public site
	// password when they are already logged in.
	if currentUser := auth.CurrentUser(r, s.database); currentUser != nil &&
		(currentUser.ID == site.UserID || s.sameOrgAsSiteOwner(site, currentUser.ID)) {
		redirectSeeOther(w, contentURL)
		return
	}

	switch r.Method {
	case http.MethodGet, http.MethodHead:
		// A valid site cookie should not show the form again — send the
		// visitor straight through, e.g. after a bookmark or a back button.
		if s.hasValidSiteSession(r, site) {
			redirectSeeOther(w, contentURL)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Cache-Control", "no-store")
		_, _ = w.Write([]byte(passwordPageHTML(slug, site.Name, next, "")))
	case http.MethodPost:
		s.handlePasswordSubmit(w, r, site, slug, next)

	default:
		w.Header().Set("Allow", "GET, HEAD, POST")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handlePasswordSubmit validates a submitted password and issues a site session.
func (s *Server) handlePasswordSubmit(w http.ResponseWriter, r *http.Request, site *db.Site, slug, next string) {
	isJSON := strings.Contains(r.Header.Get("Content-Type"), "application/json")

	rateLimited := func() {
		errMsg := "Too many attempts. Please try again later."
		if isJSON {
			writeJSON(w, http.StatusTooManyRequests, jsonResp{Error: errMsg})
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Cache-Control", "no-store")
		w.WriteHeader(http.StatusTooManyRequests)
		_, _ = w.Write([]byte(passwordPageHTML(slug, site.Name, next, errMsg)))
	}
	if !authRateLimiter.allow(clientIP(r)) {
		rateLimited()
		return
	}

	var password string
	if isJSON {
		var body struct {
			Password string `json:"password"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			writeJSON(w, http.StatusBadRequest, jsonResp{Error: tMsg(r, "invalid_json")})
			return
		}
		password = body.Password
	} else {
		// ParseForm covers both application/x-www-form-urlencoded and
		// multipart bodies. A parse failure used to be ignored, which made
		// the submission look unresponsive; report it instead.
		if err := r.ParseForm(); err != nil {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.Header().Set("Cache-Control", "no-store")
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(passwordPageHTML(slug, site.Name, next, "Could not read the submitted form. Please try again.")))
			return
		}
		password = r.FormValue("password")
	}

	if !auth.CheckPassword(site.Password, password) {
		// The password page is public-facing and always English.
		errMsg := "Incorrect password"
		if isJSON {
			writeJSON(w, http.StatusUnauthorized, jsonResp{Error: errMsg})
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Cache-Control", "no-store")
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(passwordPageHTML(slug, site.Name, next, errMsg)))
		return
	}

	// Create site session.
	token := auth.GenerateToken()
	expires := time.Now().Add(7 * 24 * time.Hour)
	if err := db.CreateSiteSession(s.database, site.ID, token, expires); err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	if isJSON {
		// Return token for API clients.
		writeJSON(w, http.StatusOK, jsonResp{
			Message: "authenticated",
			Data: map[string]interface{}{
				"token": token,
				"url":   "s/" + slug + "/" + next,
			},
		})
		return
	}

	// One cookie per site, named after the slug. Naming keeps sites from
	// overwriting each other without relying on Path scoping, which cannot work
	// behind a proxy prefix: a cookie scoped to "/s/{slug}/" is never sent back
	// for "/v/s/{slug}/", so a correct password would look like it was ignored.
	secure := r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https"
	http.SetCookie(w, &http.Cookie{
		Name:     auth.SiteCookieName(slug),
		Value:    token,
		Path:     "/",
		MaxAge:   7 * 24 * 3600,
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
	})
	// Drop a legacy root-scoped cookie from older versions so it cannot shadow
	// the per-site cookie.
	http.SetCookie(w, &http.Cookie{
		Name:     "site_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
	})
	redirectSeeOther(w, relativeSiteContentURL(r.URL.Path, slug, next))
}

// hasValidSiteSession reports whether the request carries an unexpired site
// session token for this exact site.
func (s *Server) hasValidSiteSession(r *http.Request, site *db.Site) bool {
	token := auth.GetSiteToken(r, site.Slug)
	if token == "" {
		return false
	}
	session, err := db.GetSiteSession(s.database, token)
	if err != nil || session == nil {
		return false
	}
	return session.ID == site.ID
}

// redirectSeeOther sends a 303 See Other to an absolute path.
func redirectSeeOther(w http.ResponseWriter, location string) {
	w.Header().Set("Location", location)
	w.WriteHeader(http.StatusSeeOther)
}
