package server

import (
	"fmt"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"static-host/internal/db"
)

// staticHandler serves files from the site's storage directory.
// Handles path safety, MIME detection, index.html fallback, and password protection.
func (s *Server) staticHandler(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.SplitN(strings.TrimPrefix(r.URL.Path, "/s/"), "/", 2)
	slug := pathParts[0]
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

	// Password protection check
	if site.Password != "" {
		cookieName := "vibeshare_site_" + slug
		cookie, err := r.Cookie(cookieName)
		if err != nil || cookie.Value == "" {
			http.Redirect(w, r, "/p/"+slug, http.StatusSeeOther)
			return
		}
		ss, err := db.GetSiteSession(s.database, cookie.Value)
		if err != nil || ss == nil || ss.ID != site.ID {
			http.Redirect(w, r, "/p/"+slug, http.StatusSeeOther)
			return
		}
	}

	// Determine sub-path
	subPath := ""
	if len(pathParts) > 1 {
		subPath = pathParts[1]
	}

	siteDir := filepath.Join(s.config.StorageDir, slug)
	s.serveStaticFile(w, r, siteDir, subPath)
}

// serveStaticFile reads a file from disk and serves it with proper MIME and headers.
func (s *Server) serveStaticFile(w http.ResponseWriter, r *http.Request, siteDir, subPath string) {
	// Clean and safe-join the path
	cleanPath := filepath.Clean("/" + subPath)

	fullPath := filepath.Join(siteDir, cleanPath)

	// Ensure the resolved path is within siteDir
	absSiteDir, _ := filepath.Abs(siteDir)
	absFullPath, _ := filepath.Abs(fullPath)
	if !strings.HasPrefix(absFullPath, absSiteDir) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	stat, err := os.Stat(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			// SPA fallback: try index.html at root
			indexPath := filepath.Join(siteDir, "index.html")
			if idxStat, err := os.Stat(indexPath); err == nil && !idxStat.IsDir() {
				s.writeFile(w, r, indexPath, idxStat)
				return
			}
			http.NotFound(w, r)
			return
		}
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	if stat.IsDir() {
		indexPath := filepath.Join(fullPath, "index.html")
		if idxStat, err := os.Stat(indexPath); err == nil && !idxStat.IsDir() {
			s.writeFile(w, r, indexPath, idxStat)
			return
		}
		http.NotFound(w, r)
		return
	}

	s.writeFile(w, r, fullPath, stat)
}

// writeFile sets headers and writes the file content to the response.
func (s *Server) writeFile(w http.ResponseWriter, r *http.Request, path string, stat os.FileInfo) {
	ext := filepath.Ext(path)
	ct := mime.TypeByExtension(ext)
	if ct == "" {
		ct = "application/octet-stream"
	}
	w.Header().Set("Content-Type", ct)
	w.Header().Set("Content-Length", fmt.Sprintf("%d", stat.Size()))
	w.Header().Set("Last-Modified", stat.ModTime().UTC().Format(http.TimeFormat))
	w.Header().Set("Cache-Control", "public, max-age=3600")

	// Conditional requests
	if t, err := time.Parse(http.TimeFormat, r.Header.Get("If-Modified-Since")); err == nil {
		if !stat.ModTime().Truncate(time.Second).After(t.Truncate(time.Second)) {
			w.WriteHeader(http.StatusNotModified)
			return
		}
	}

	http.ServeFile(w, r, path)
}
