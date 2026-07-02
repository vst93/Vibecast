package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"static-host/internal/auth"
	"static-host/internal/db"
)

// passwordPageHandler shows the password gate page and processes password submissions.
// GET /p/{slug}  → show password form
// POST /p/{slug} → validate password, set cookie, redirect to site
func (s *Server) passwordPageHandler(w http.ResponseWriter, r *http.Request) {
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
	if site.Password == "" {
		// Not protected — redirect to site
		http.Redirect(w, r, "/s/"+slug+"/", http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodGet {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, passwordPageHTML(slug, site.Name))
		return
	}

	if r.Method == http.MethodPost {
		var body struct {
			Password string `json:"password"`
		}
		contentType := r.Header.Get("Content-Type")
		if strings.Contains(contentType, "application/json") {
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				writeJSON(w, 400, jsonResp{Error: "invalid JSON"})
				return
			}
		} else {
			// Form submission
			r.ParseForm()
			body.Password = r.FormValue("password")
		}

		if !auth.CheckPassword(site.Password, body.Password) {
			if strings.Contains(contentType, "application/json") {
				writeJSON(w, 401, jsonResp{Error: "incorrect password"})
			} else {
				w.Header().Set("Content-Type", "text/html; charset=utf-8")
				w.WriteHeader(http.StatusUnauthorized)
				fmt.Fprint(w, passwordPageHTMLWithErr(slug, site.Name, "密码错误，请重试"))
			}
			return
		}

		// Create site session
		token := auth.GenerateToken()
		expires := time.Now().Add(7 * 24 * time.Hour)
		if err := db.CreateSiteSession(s.database, site.ID, token, expires); err != nil {
			http.Error(w, "Internal error", http.StatusInternalServerError)
			return
		}

		// Set cookie
		cookieName := "vibeshare_site_" + slug
		http.SetCookie(w, &http.Cookie{
			Name:     cookieName,
			Value:    token,
			Path:     "/",
			MaxAge:   int((7 * 24 * time.Hour).Seconds()),
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
		})

		if strings.Contains(contentType, "application/json") {
			writeJSON(w, 200, jsonResp{Message: "authenticated"})
		} else {
			http.Redirect(w, r, "/s/"+slug+"/", http.StatusSeeOther)
		}
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}
