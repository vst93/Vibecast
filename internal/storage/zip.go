package storage

import (
	"archive/zip"
	"errors"
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

var (
	ErrSiteSizeLimit = errors.New("site size limit exceeded")
	ErrSiteEmpty     = errors.New("site has no files")
)

// FileEntry describes one path inside a site's storage directory.
type FileEntry struct {
	Name string `json:"name"`
	Size int64  `json:"size"`
	Dir  bool   `json:"dir"`
}

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

// excludedDirs are directories whose contents are junk / metadata
// and should be stripped during zip extraction.
var excludedDirs = map[string]bool{
	"__MACOSX":     true, // macOS resource-fork metadata
	".git":         true, // git repository metadata
	".svn":         true, // subversion metadata
	".hg":          true, // mercurial metadata
	"node_modules": true, // npm dependencies (usually huge, not needed for static deploy)
}

// excludedFiles are individual filenames to strip.
var excludedFiles = map[string]bool{
	".DS_Store":      true, // macOS desktop store
	"Thumbs.db":      true, // Windows thumbnail cache
	"desktop.ini":    true, // Windows folder config
	".gitattributes": true,
	".gitignore":     true,
}

// isExcludedPath returns true if any path segment is an excluded dir
// or the final filename is an excluded file.
func isExcludedPath(name string) bool {
	if name == "" {
		return false
	}
	segs := strings.Split(name, "/")
	for _, seg := range segs {
		if excludedDirs[seg] {
			return true
		}
	}
	return excludedFiles[segs[len(segs)-1]]
}

// IsBlockedExtension returns true if the file extension is in the blocklist.
func IsBlockedExtension(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	return blockedExtensions[ext]
}

// SaveSingleFile saves one uploaded file while preserving all other site files.
// A file with the same basename is replaced after the final site size is checked.
// maxFileSize is also the maximum aggregate site size (0 = default 100 MB).
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

	if err := os.MkdirAll(destDir, 0755); err != nil {
		return 0, fmt.Errorf("mkdir site: %w", err)
	}
	target := filepath.Join(destDir, filename)
	var replacedSize int64
	if info, err := os.Lstat(target); err == nil {
		if !info.Mode().IsRegular() {
			return 0, fmt.Errorf("upload path conflicts with non-file: %s", filename)
		}
		replacedSize = info.Size()
	} else if !os.IsNotExist(err) {
		return 0, fmt.Errorf("stat existing file: %w", err)
	}

	currentSize, err := DirectorySize(destDir)
	if err != nil {
		return 0, err
	}

	out, err := os.CreateTemp(destDir, ".upload-")
	if err != nil {
		return 0, fmt.Errorf("create temp file: %w", err)
	}
	tmpName := out.Name()
	defer os.Remove(tmpName)

	n, err := io.Copy(out, io.LimitReader(src, maxFileSize+1))
	if err != nil {
		_ = out.Close()
		return 0, fmt.Errorf("write file: %w", err)
	}
	if err := out.Close(); err != nil {
		return 0, fmt.Errorf("close file: %w", err)
	}

	if n > maxFileSize {
		return 0, fmt.Errorf("file too large: %d bytes > %d limit", n, maxFileSize)
	}
	projectedSize := currentSize - replacedSize + n
	if projectedSize > maxFileSize {
		return 0, fmt.Errorf("%w: %d bytes > %d limit", ErrSiteSizeLimit, projectedSize, maxFileSize)
	}
	if err := replaceFile(tmpName, target); err != nil {
		return 0, err
	}

	return n, nil
}

// ExtractZipResult holds extraction results including skipped dangerous files.
type ExtractZipResult struct {
	WebRoot    string   // actual web root directory
	Skipped    []string // list of skipped file names (dangerous types)
	TotalFiles int      // number of files extracted
	TotalSize  int64    // total bytes extracted
	SiteSize   int64    // final aggregate bytes in the merged site
}

// ExtractZip extracts a zip reader contents into destDir.
// It strips common top-level directory if all entries share one (e.g. "site/" prefix).
// Dangerous file types are skipped. Zip bombs are rejected.
// Existing files not present in the archive are preserved. Matching paths are
// overwritten. maxSingleFileSize is also the aggregate site limit.
func ExtractZip(r io.ReaderAt, size int64, destDir string, maxSingleFileSize int64) (*ExtractZipResult, error) {
	if maxSingleFileSize <= 0 {
		maxSingleFileSize = defaultMaxSingleFileSize
	}
	tmpDir, err := os.MkdirTemp(filepath.Dir(destDir), filepath.Base(destDir)+".tmp-")
	if err != nil {
		return nil, fmt.Errorf("mkdir tmp: %w", err)
	}
	committed := false
	defer func() {
		if !committed {
			_ = os.RemoveAll(tmpDir)
		}
	}()
	if err := cloneDirectory(destDir, tmpDir); err != nil {
		return nil, fmt.Errorf("clone existing site: %w", err)
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
	for _, f := range zr.File {
		if f.Name == "" {
			continue
		}
		if f.Name[0] == '/' || strings.Contains(f.Name, "..") {
			return nil, fmt.Errorf("unsafe path in zip: %s", f.Name)
		}
		// Skip dotfiles and excluded junk dirs (e.g. __MACOSX) for prefix detection
		firstSeg := f.Name
		if idx := strings.Index(f.Name, "/"); idx > 0 {
			firstSeg = f.Name[:idx]
		}
		if strings.HasPrefix(firstSeg, ".") || excludedDirs[firstSeg] {
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

		// Skip junk / metadata files and directories (__MACOSX, .git, .DS_Store, etc.)
		if isExcludedPath(name) {
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
			if info, err := os.Lstat(target); err == nil && !info.IsDir() {
				return nil, fmt.Errorf("zip path conflicts with existing file: %s", name)
			} else if err != nil && !os.IsNotExist(err) {
				return nil, fmt.Errorf("stat %s: %w", target, err)
			}
			if err := os.MkdirAll(target, 0755); err != nil {
				return nil, fmt.Errorf("mkdir %s: %w", target, err)
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
			return nil, fmt.Errorf("mkdir parent: %w", err)
		}
		if info, err := os.Lstat(target); err == nil {
			if info.IsDir() {
				return nil, fmt.Errorf("zip path conflicts with existing directory: %s", name)
			}
			if err := os.Remove(target); err != nil {
				return nil, fmt.Errorf("replace existing file %s: %w", name, err)
			}
		} else if !os.IsNotExist(err) {
			return nil, fmt.Errorf("stat existing file %s: %w", name, err)
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

	finalSize, err := DirectorySize(tmpDir)
	if err != nil {
		return nil, err
	}
	if finalSize > maxSingleFileSize {
		return nil, fmt.Errorf("%w: %d bytes > %d limit", ErrSiteSizeLimit, finalSize, maxSingleFileSize)
	}
	if err := swapDirectory(tmpDir, destDir); err != nil {
		return nil, err
	}
	committed = true

	result.WebRoot = destDir
	result.TotalSize = totalSize
	result.SiteSize = finalSize
	return result, nil
}

// ListFiles recursively lists regular files and directories below root.
func ListFiles(root string) ([]FileEntry, int64, error) {
	entries := make([]FileEntry, 0)
	var totalSize int64
	if _, err := os.Stat(root); os.IsNotExist(err) {
		return entries, 0, nil
	} else if err != nil {
		return nil, 0, fmt.Errorf("stat site directory: %w", err)
	}
	err := filepath.WalkDir(root, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if path == root {
			return nil
		}
		if entry.Type()&os.ModeSymlink != 0 {
			return nil
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		item := FileEntry{Name: filepath.ToSlash(rel), Dir: entry.IsDir()}
		if !entry.IsDir() {
			info, err := entry.Info()
			if err != nil {
				return err
			}
			if !info.Mode().IsRegular() {
				return nil
			}
			item.Size = info.Size()
			totalSize += info.Size()
		}
		entries = append(entries, item)
		return nil
	})
	if err != nil {
		return nil, 0, fmt.Errorf("list site files: %w", err)
	}
	return entries, totalSize, nil
}

// DirectorySize returns the sum of regular file sizes below root.
func DirectorySize(root string) (int64, error) {
	_, total, err := ListFiles(root)
	return total, err
}

// ClearDirectory removes all site content and recreates an empty directory.
func ClearDirectory(destDir string) error {
	if err := os.RemoveAll(destDir); err != nil {
		return fmt.Errorf("clear site directory: %w", err)
	}
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("recreate site directory: %w", err)
	}
	return nil
}

// ZipDirectory writes all regular files below root to a ZIP stream.
func ZipDirectory(dst io.Writer, root string) (int, int64, error) {
	entries, totalSize, err := ListFiles(root)
	if err != nil {
		return 0, 0, err
	}
	fileCount := 0
	for _, entry := range entries {
		if !entry.Dir {
			fileCount++
		}
	}
	if fileCount == 0 {
		return 0, 0, ErrSiteEmpty
	}

	zw := zip.NewWriter(dst)
	for _, entry := range entries {
		path := filepath.Join(root, filepath.FromSlash(entry.Name))
		info, err := os.Stat(path)
		if err != nil {
			_ = zw.Close()
			return 0, 0, fmt.Errorf("stat zip entry: %w", err)
		}
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			_ = zw.Close()
			return 0, 0, fmt.Errorf("create zip header: %w", err)
		}
		header.Name = entry.Name
		if entry.Dir {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}
		writer, err := zw.CreateHeader(header)
		if err != nil {
			_ = zw.Close()
			return 0, 0, fmt.Errorf("create zip entry: %w", err)
		}
		if entry.Dir {
			continue
		}
		file, err := os.Open(path)
		if err != nil {
			_ = zw.Close()
			return 0, 0, fmt.Errorf("open zip entry: %w", err)
		}
		_, copyErr := io.Copy(writer, file)
		closeErr := file.Close()
		if copyErr != nil {
			_ = zw.Close()
			return 0, 0, fmt.Errorf("write zip entry: %w", copyErr)
		}
		if closeErr != nil {
			_ = zw.Close()
			return 0, 0, fmt.Errorf("close zip entry: %w", closeErr)
		}
	}
	if err := zw.Close(); err != nil {
		return 0, 0, fmt.Errorf("finalize zip: %w", err)
	}
	return fileCount, totalSize, nil
}

func replaceFile(tmpName, target string) error {
	backup, err := reservePath(filepath.Dir(target), ".backup-")
	if err != nil {
		return err
	}
	hadTarget := false
	if _, err := os.Lstat(target); err == nil {
		hadTarget = true
		if err := os.Rename(target, backup); err != nil {
			return fmt.Errorf("backup existing file: %w", err)
		}
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("stat existing file: %w", err)
	}
	if err := os.Rename(tmpName, target); err != nil {
		if hadTarget {
			_ = os.Rename(backup, target)
		}
		return fmt.Errorf("commit uploaded file: %w", err)
	}
	if hadTarget {
		_ = os.Remove(backup)
	}
	return nil
}

func cloneDirectory(src, dst string) error {
	if _, err := os.Stat(src); os.IsNotExist(err) {
		return nil
	} else if err != nil {
		return err
	}
	return filepath.WalkDir(src, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if path == src {
			return nil
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		target := filepath.Join(dst, rel)
		if entry.IsDir() {
			return os.MkdirAll(target, 0755)
		}
		info, err := entry.Info()
		if err != nil {
			return err
		}
		if !info.Mode().IsRegular() {
			return nil
		}
		if err := os.Link(path, target); err == nil {
			return nil
		}
		return copyFile(path, target, info.Mode().Perm())
	})
}

func copyFile(src, dst string, mode os.FileMode) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_EXCL, mode)
	if err != nil {
		return err
	}
	_, copyErr := io.Copy(out, in)
	closeErr := out.Close()
	if copyErr != nil {
		return copyErr
	}
	return closeErr
}

func swapDirectory(tmpDir, destDir string) error {
	if err := os.MkdirAll(filepath.Dir(destDir), 0755); err != nil {
		return fmt.Errorf("mkdir storage parent: %w", err)
	}
	backup, err := reservePath(filepath.Dir(destDir), filepath.Base(destDir)+".backup-")
	if err != nil {
		return err
	}
	hadDest := false
	if _, err := os.Stat(destDir); err == nil {
		hadDest = true
		if err := os.Rename(destDir, backup); err != nil {
			return fmt.Errorf("backup site directory: %w", err)
		}
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("stat site directory: %w", err)
	}
	if err := os.Rename(tmpDir, destDir); err != nil {
		if hadDest {
			_ = os.Rename(backup, destDir)
		}
		return fmt.Errorf("commit site directory: %w", err)
	}
	if hadDest {
		_ = os.RemoveAll(backup)
	}
	return nil
}

func reservePath(dir, pattern string) (string, error) {
	file, err := os.CreateTemp(dir, pattern)
	if err != nil {
		return "", fmt.Errorf("reserve backup path: %w", err)
	}
	name := file.Name()
	if err := file.Close(); err != nil {
		_ = os.Remove(name)
		return "", err
	}
	if err := os.Remove(name); err != nil {
		return "", err
	}
	return name, nil
}

// DeleteSiteDir removes the site storage directory.
func DeleteSiteDir(baseDir, slug string) error {
	return os.RemoveAll(filepath.Join(baseDir, slug))
}
