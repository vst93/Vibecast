package storage

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const (
	maxDecompressedSize = 500 * 1024 * 1024 // 500 MB total uncompressed
	maxFileCount        = 10000             // max files in a single zip
)

// maxSingleFileSize is the default per-file limit (100 MB).
// ExtractZip and SaveSingleFile accept a custom limit via parameter.
const defaultMaxSingleFileSize = 100 * 1024 * 1024

// blockedExtensions are file types that are dangerous to serve or execute.
// They are either server-side scripts, native executables, or system-level files.
var blockedExtensions = map[string]bool{
	// Server-side scripts
	".php": true, ".php3": true, ".php4": true, ".php5": true, ".phtml": true,
	".cgi": true, ".pl": true, ".py": true, ".rb": true, ".sh": true, ".bash": true,
	".asp": true, ".aspx": true, ".jsp": true, ".node": true,
	// Native executables / binaries
	".exe": true, ".bat": true, ".cmd": true, ".com": true, ".scr": true,
	".msi": true, ".dll": true, ".so": true, ".dylib": true, ".bin": true,
	".jar": true, ".app": true, ".run": true, ".out": true,
	".apk": true, ".deb": true, ".rpm": true, ".dmg": true, ".pkg": true,
	".iso": true, ".img": true,
	// Web server configs (irrelevant for Go, but could confuse)
	".htaccess": true, ".htpasswd": true,
	// System / shell
	".ps1": true, ".psm1": true, ".vbs": true, ".wsf": true,
	// Misc potentially dangerous
	".reg": true, ".lnk": true, ".desktop": true,
}

// IsBlockedExtension returns true if the file extension is in the blocklist.
func IsBlockedExtension(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	return blockedExtensions[ext]
}

// SaveSingleFile saves a single uploaded file to the site directory.
// It replaces the entire site content (same as ZIP deploy — atomic replace).
// maxFileSize is the per-file size limit in bytes (0 = use default 100 MB).
// Returns the file size.
func SaveSingleFile(src io.Reader, filename string, destDir string, maxFileSize int64) (int64, error) {
	if maxFileSize <= 0 {
		maxFileSize = defaultMaxSingleFileSize
	}
	if IsBlockedExtension(filename) {
		return 0, fmt.Errorf("file type not allowed: %s", filepath.Ext(filename))
	}

	// Sanitize filename — keep it simple, no path traversal
	filename = filepath.Base(filename)
	if filename == "" || filename == "." || filename == ".." {
		return 0, fmt.Errorf("invalid filename")
	}

	tmpDir := destDir + ".tmp"
	_ = os.RemoveAll(tmpDir)
	if err := os.MkdirAll(tmpDir, 0755); err != nil {
		return 0, fmt.Errorf("mkdir tmp: %w", err)
	}

	target := filepath.Join(tmpDir, filename)
	out, err := os.OpenFile(target, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return 0, fmt.Errorf("create file: %w", err)
	}

	n, err := io.Copy(out, src)
	if err != nil {
		out.Close()
		return 0, fmt.Errorf("write file: %w", err)
	}
	out.Close()

	if n > maxFileSize {
		_ = os.RemoveAll(tmpDir)
		return 0, fmt.Errorf("file too large: %d bytes > %d limit", n, maxFileSize)
	}

	// Atomic swap
	_ = os.RemoveAll(destDir)
	if err := os.Rename(tmpDir, destDir); err != nil {
		return 0, fmt.Errorf("rename %s -> %s: %w", tmpDir, destDir, err)
	}

	return n, nil
}

// ExtractZipResult holds extraction results including skipped dangerous files.
type ExtractZipResult struct {
	WebRoot     string   // actual web root directory
	Skipped     []string // list of skipped file names (dangerous types)
	TotalFiles  int      // number of files extracted
	TotalSize   int64    // total bytes extracted
}

// ExtractZip extracts a zip reader contents into destDir.
// It strips common top-level directory if all entries share one (e.g. "site/" prefix).
// Dangerous file types are skipped. Zip bombs are rejected.
// maxSingleFileSize is the per-file limit in bytes (0 = use default 100 MB).
func ExtractZip(r io.ReaderAt, size int64, destDir string, maxSingleFileSize int64) (*ExtractZipResult, error) {
	if maxSingleFileSize <= 0 {
		maxSingleFileSize = defaultMaxSingleFileSize
	}
	tmpDir := destDir + ".tmp"
	_ = os.RemoveAll(tmpDir)
	if err := os.MkdirAll(tmpDir, 0755); err != nil {
		return nil, fmt.Errorf("mkdir tmp: %w", err)
	}

	zr, err := zip.NewReader(r, size)
	if err != nil {
		return nil, fmt.Errorf("zip reader: %w", err)
	}

	// Check file count
	if len(zr.File) > maxFileCount {
		return nil, fmt.Errorf("zip contains too many files (%d > %d limit)", len(zr.File), maxFileCount)
	}

	// Detect common prefix
	// Skip dotfiles (e.g. .DS_Store) when detecting prefix — they shouldn't
	// prevent prefix stripping for the actual site files.
	prefix := ""
	for i, f := range zr.File {
		if f.Name[0] == '/' || strings.Contains(f.Name, "..") {
			return nil, fmt.Errorf("unsafe path in zip: %s", f.Name)
		}
		// Skip dotfiles for prefix detection
		firstSeg := f.Name
		if idx := strings.Index(f.Name, "/"); idx > 0 {
			firstSeg = f.Name[:idx]
		}
		if strings.HasPrefix(firstSeg, ".") {
			continue
		}
		parts := strings.SplitN(f.Name, "/", 2)
		if len(parts) < 2 {
			prefix = ""
			break
		}
		if prefix == "" {
			prefix = parts[0] + "/"
		} else if !strings.HasPrefix(f.Name, prefix) {
			prefix = ""
			break
		}
		_ = i
	}

	result := &ExtractZipResult{}
	var totalSize int64

	for _, f := range zr.File {
		// Strip the common prefix
		name := f.Name
		if prefix != "" && strings.HasPrefix(name, prefix) {
			name = strings.TrimPrefix(name, prefix)
		}
		if name == "" {
			continue
		}

		// Check for dangerous file extensions
		ext := strings.ToLower(filepath.Ext(name))
		if blockedExtensions[ext] {
			result.Skipped = append(result.Skipped, name)
			continue
		}

		// Check single file size (uncompressed)
		if f.UncompressedSize64 > uint64(maxSingleFileSize) {
			return nil, fmt.Errorf("file too large: %s (%d bytes > %d limit)", name, f.UncompressedSize64, maxSingleFileSize)
		}

		// Check total decompressed size
		totalSize += int64(f.UncompressedSize64)
		if totalSize > maxDecompressedSize {
			return nil, fmt.Errorf("zip bomb detected: total uncompressed size exceeds %d MB limit", maxDecompressedSize/(1024*1024))
		}

		target := filepath.Join(tmpDir, name)

		// Prevent path traversal
		if !strings.HasPrefix(filepath.Clean(target), filepath.Clean(tmpDir)+string(os.PathSeparator)) {
			return nil, fmt.Errorf("path traversal detected: %s", f.Name)
		}

		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(target, 0755); err != nil {
				return nil, fmt.Errorf("mkdir %s: %w", target, err)
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
			return nil, fmt.Errorf("mkdir parent: %w", err)
		}

		rc, err := f.Open()
		if err != nil {
			return nil, fmt.Errorf("open zip entry %s: %w", f.Name, err)
		}

		out, err := os.OpenFile(target, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			rc.Close()
			return nil, fmt.Errorf("create file %s: %w", target, err)
		}

		if _, err := io.Copy(out, rc); err != nil {
			rc.Close()
			out.Close()
			return nil, fmt.Errorf("write file %s: %w", target, err)
		}
		rc.Close()
		out.Close()

		result.TotalFiles++
	}

	// Atomic-ish swap: remove old dir, rename tmp
	_ = os.RemoveAll(destDir)
	if err := os.Rename(tmpDir, destDir); err != nil {
		return nil, fmt.Errorf("rename %s -> %s: %w", tmpDir, destDir, err)
	}

	result.WebRoot = destDir
	result.TotalSize = totalSize
	return result, nil
}

// DeleteSiteDir removes the site storage directory.
func DeleteSiteDir(baseDir, slug string) error {
	return os.RemoveAll(filepath.Join(baseDir, slug))
}
