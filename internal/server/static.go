package server

import (
	"fmt"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"vibecast/internal/auth"
	"vibecast/internal/db"
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
	// Office documents
	".doc":   "application/msword",
	".docx":  "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
	".xls":   "application/vnd.ms-excel",
	".xlsx":  "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
	".ppt":   "application/vnd.ms-powerpoint",
	".pptx":  "application/vnd.openxmlformats-officedocument.presentationml.presentation",
	".odt":   "application/vnd.oasis.opendocument.text",
	".ods":   "application/vnd.oasis.opendocument.spreadsheet",
	".odp":   "application/vnd.oasis.opendocument.presentation",
	".rtf":   "application/rtf",
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
	// Check if public access is allowed (password-protected sites are still accessible)
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

	// If public access is disabled, only password-protected sites can be accessed
	if !db.GetSettingBool(s.database, "allow_public_access", true) && site.Password == "" {
		http.Error(w, "Public access is disabled", http.StatusForbidden)
		return
	}

	// Password protection check — token from Authorization header or ?token= query param
	if site.Password != "" {
		token := auth.GetSiteToken(r)
		if token == "" || func() bool {
			ss, err := db.GetSiteSession(s.database, token)
			return err != nil || ss == nil || ss.ID != site.ID
		}() {
			http.Redirect(w, r, prefURL(r, "/p/"+slug), http.StatusSeeOther)
			return
		}
	}

	// Record visit (async, non-blocking — don't slow down the request)
	go func() {
		now := time.Now()
		_ = db.RecordVisit(s.database, site.ID, now.Format("2006-01-02"), now.Format("2006-01"))
	}()

	// Determine sub-path
	subPath := ""
	if len(pathParts) > 1 {
		subPath = pathParts[1]
	}

	siteDir := filepath.Join(s.config.StorageDir, slug)
	s.serveStaticFile(w, r, siteDir, slug, subPath)
}

// serveStaticFile reads a file from disk and serves it with proper MIME and headers.
func (s *Server) serveStaticFile(w http.ResponseWriter, r *http.Request, siteDir, slug, subPath string) {
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
			// Try directory listing on the site root
			if subPath == "" {
				s.serveDirListing(w, r, siteDir, slug, "")
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
		// No index.html — serve directory listing (nginx-style)
		s.serveDirListing(w, r, fullPath, slug, subPath)
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

// serveDirListing renders an nginx-style directory listing.
func (s *Server) serveDirListing(w http.ResponseWriter, r *http.Request, dirPath, slug, subPath string) {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	// Build breadcrumb
	baseURL := prefURL(r, "/s/" + slug + "/")
	parts := strings.Split(strings.TrimSuffix(subPath, "/"), "/")
	var crumbs []string
	crumbs = append(crumbs, `<a href="`+baseURL+`">/</a>`)
	acc := ""
	for i, p := range parts {
		if p == "" {
			continue
		}
		acc += p + "/"
		sep := ""
		if i > 0 || subPath != "" {
			sep = "/"
		}
		crumbs = append(crumbs, `<a href="`+baseURL+acc+`">`+p+`</a>`+sep)
	}

	var b strings.Builder
	b.WriteString(`<!DOCTYPE html><html><head><meta charset="UTF-8"><meta name="viewport" content="width=device-width,initial-scale=1">`)
	b.WriteString(`<style>body{font-family:system-ui,sans-serif;background:`)
	b.WriteString(`var(--ink,#0c1117);color:var(--text,#e6edf3);margin:0;padding:2rem}
a{color:#39d353;text-decoration:none}a:hover{text-decoration:underline}
h1{font-size:1.1rem;font-weight:600;margin-bottom:1rem}
table{border-collapse:collapse;width:100%;max-width:800px}
th,td{text-align:left;padding:6px 12px;border-bottom:1px solid #30363d;font-size:.85rem}
th{color:#7d8590;font-size:.75rem;text-transform:uppercase;font-weight:600}
.dir{font-weight:600}.size{color:#7d8590;text-align:right;font-family:monospace}</style>`)
	b.WriteString(`</head><body><h1>`)
	b.WriteString(strings.Join(crumbs, " / "))
	b.WriteString(`</h1><table><thead><tr><th>Name</th><th>Size</th></tr></thead><tbody>`)

	// ".." link for subdirectories
	if subPath != "" {
		parent := baseURL
		if parts := strings.Split(strings.TrimSuffix(subPath, "/"), "/"); len(parts) > 1 {
			parent = baseURL + strings.Join(parts[:len(parts)-1], "/") + "/"
		}
		b.WriteString(`<tr><td class="dir"><a href="` + parent + `">../</a></td><td>-</td></tr>`)
	}

	for _, e := range entries {
		name := e.Name()
		if strings.HasPrefix(name, ".") {
			continue
		}
		info, _ := e.Info()
		size := "-"
		if !e.IsDir() && info != nil {
			size = formatSize(info.Size())
		}
		displayName := name
		href := baseURL
		if subPath != "" {
			href += subPath
			if !strings.HasSuffix(href, "/") {
				href += "/"
			}
		}
		href += name
		cls := ""
		if e.IsDir() {
			displayName += "/"
			href += "/"
			cls = ` class="dir"`
		}
		b.WriteString(`<tr><td` + cls + `><a href="` + href + `">` + displayName + `</a></td><td class="size">` + size + `</td></tr>`)
	}
	b.WriteString(`</tbody></table></body></html>`)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache, must-revalidate")
	w.WriteHeader(200)
	w.Write([]byte(b.String()))
}

// formatSize converts a byte count to a human-readable string.
func formatSize(n int64) string {
	if n < 1024 {
		return fmt.Sprintf("%d B", n)
	}
	if n < 1024*1024 {
		return fmt.Sprintf("%.1f KB", float64(n)/1024)
	}
	if n < 1024*1024*1024 {
		return fmt.Sprintf("%.1f MB", float64(n)/(1024*1024))
	}
	return fmt.Sprintf("%.1f GB", float64(n)/(1024*1024*1024))
}
