package server

import (
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"static-host/internal/auth"
	"static-host/internal/db"
)

// extraMimeTypes extends Go's built-in mime map with modern frontend file types.
var extraMimeTypes = map[string]string{
	// JavaScript
	".mjs":   "text/javascript; charset=utf-8",
	".js":    "text/javascript; charset=utf-8",
	".map":   "application/json; charset=utf-8", // source maps
	// Fonts
	".woff2": "font/woff2",
	".woff":  "font/woff",
	".ttf":   "font/ttf",
	".otf":   "font/otf",
	".eot":   "application/vnd.ms-fontobject",
	// Images
	".webp":  "image/webp",
	".avif":  "image/avif",
	".svg":   "image/svg+xml",
	".ico":   "image/x-icon",
	// Documents
	".webmanifest": "application/manifest+json",
	".json":  "application/json; charset=utf-8",
	".xml":   "application/xml; charset=utf-8",
	".csv":   "text/csv; charset=utf-8",
	".txt":   "text/plain; charset=utf-8",
	".html":  "text/html; charset=utf-8",
	".htm":   "text/html; charset=utf-8",
	".css":   "text/css; charset=utf-8",
	// Video / Audio
	".mp4":   "video/mp4",
	".webm":  "video/webm",
	".mp3":   "audio/mpeg",
	".ogg":   "audio/ogg",
	".wav":   "audio/wav",
	".flac":  "audio/flac",
	// Other
	".wasm":  "application/wasm",
	".pdf":   "application/pdf",
}

// getContentType returns the MIME type for a file extension.
// Falls back to Go's built-in mime map, then to application/octet-stream.
func getContentType(ext string) string {
	ext = strings.ToLower(ext)
	if ct, ok := extraMimeTypes[ext]; ok {
		return ct
	}
	// Try Go's built-in map
	if ct := mime.TypeByExtension(ext); ct != "" {
		return ct
	}
	return "application/octet-stream"
}

// staticHandler serves files from the site's storage directory.
// Handles path safety, MIME detection, index.html fallback, SPA fallback, and password protection.
func (s *Server) staticHandler(w http.ResponseWriter, r *http.Request) {
	// Check if public access is allowed
	if !db.GetSettingBool(s.database, "allow_public_access", true) {
		http.Error(w, "Public access is disabled", http.StatusForbidden)
		return
	}

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

	// Password protection check — token from Authorization header or ?token= query param
	if site.Password != "" {
		token := auth.GetSiteToken(r)
		if token == "" {
			http.Redirect(w, r, "/p/"+slug, http.StatusSeeOther)
			return
		}
		ss, err := db.GetSiteSession(s.database, token)
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

	// Ensure the resolved path is within siteDir (prevent path traversal)
	absSiteDir, _ := filepath.Abs(siteDir)
	absFullPath, _ := filepath.Abs(fullPath)
	if !strings.HasPrefix(absFullPath, absSiteDir) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// Also reject dotfiles (e.g. .env, .git, .htaccess) for security
	relPath, _ := filepath.Rel(absSiteDir, absFullPath)
	for _, seg := range strings.Split(relPath, string(os.PathSeparator)) {
		if strings.HasPrefix(seg, ".") && seg != "." && seg != ".." {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
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
// Uses http.ServeContent to handle conditional requests, range requests, and ETag.
func (s *Server) writeFile(w http.ResponseWriter, r *http.Request, path string, stat os.FileInfo) {
	ext := filepath.Ext(path)
	ct := getContentType(ext)

	// Security headers
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("X-Frame-Options", "SAMEORIGIN")
	w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
	w.Header().Set("Content-Type", ct)

	// Cache-Control: static assets can be cached aggressively, HTML should not
	if ext == ".html" || ext == ".htm" {
		w.Header().Set("Cache-Control", "no-cache, must-revalidate")
	} else {
		w.Header().Set("Cache-Control", "public, max-age=86400") // 24h
	}

	// Open the file
	f, err := os.Open(path)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}
	defer f.Close()

	// Use ServeContent — handles Last-Modified, If-Modified-Since, Range, ETag
	modTime := stat.ModTime()
	http.ServeContent(w, r, filepath.Base(path), modTime, f)
}
