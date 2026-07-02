package storage

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// ExtractZip extracts a zip reader contents into destDir.
// It strips common top-level directory if all entries share one (e.g. "site/" prefix).
// Returns the actual web root (destDir or a subdirectory of it).
func ExtractZip(r io.ReaderAt, size int64, destDir string) (string, error) {
	tmpDir := destDir + ".tmp"
	_ = os.RemoveAll(tmpDir)
	if err := os.MkdirAll(tmpDir, 0755); err != nil {
		return "", fmt.Errorf("mkdir tmp: %w", err)
	}

	zr, err := zip.NewReader(r, size)
	if err != nil {
		return "", fmt.Errorf("zip reader: %w", err)
	}

	// Detect common prefix
	prefix := ""
	for i, f := range zr.File {
		// Only consider actual directory entries or files for prefix detection
		if f.Name[0] == '/' || strings.Contains(f.Name, "..") {
			return "", fmt.Errorf("unsafe path in zip: %s", f.Name)
		}
		parts := strings.SplitN(f.Name, "/", 2)
		if len(parts) < 2 {
			// File at root level — no common prefix to strip
			prefix = ""
			break
		}
		if i == 0 {
			prefix = parts[0] + "/"
		} else if !strings.HasPrefix(f.Name, prefix) {
			prefix = ""
			break
		}
	}

	for _, f := range zr.File {
		// Strip the common prefix
		name := f.Name
		if prefix != "" && strings.HasPrefix(name, prefix) {
			name = strings.TrimPrefix(name, prefix)
		}
		if name == "" {
			continue
		}

		target := filepath.Join(tmpDir, name)

		// Prevent path traversal
		if !strings.HasPrefix(filepath.Clean(target), filepath.Clean(tmpDir)+string(os.PathSeparator)) {
			return "", fmt.Errorf("path traversal detected: %s", f.Name)
		}

		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(target, 0755); err != nil {
				return "", fmt.Errorf("mkdir %s: %w", target, err)
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
			return "", fmt.Errorf("mkdir parent: %w", err)
		}

		rc, err := f.Open()
		if err != nil {
			return "", fmt.Errorf("open zip entry %s: %w", f.Name, err)
		}

		out, err := os.OpenFile(target, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			rc.Close()
			return "", fmt.Errorf("create file %s: %w", target, err)
		}

		if _, err := io.Copy(out, rc); err != nil {
			rc.Close()
			out.Close()
			return "", fmt.Errorf("write file %s: %w", target, err)
		}
		rc.Close()
		out.Close()
	}

	// Atomic-ish swap: remove old dir, rename tmp
	_ = os.RemoveAll(destDir)
	if err := os.Rename(tmpDir, destDir); err != nil {
		// Fallback: try copy
		return "", fmt.Errorf("rename %s -> %s: %w", tmpDir, destDir, err)
	}

	return destDir, nil
}

// DeleteSite removes the site storage directory.
func DeleteSiteDir(baseDir, slug string) error {
	return os.RemoveAll(filepath.Join(baseDir, slug))
}
