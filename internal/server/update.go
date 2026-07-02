package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"vibecast/internal/db"
)

// GitHub mirror hosts for China mainland — tried in order.
// Each entry is a prefix prepended to the full GitHub URL.
// Empty string = direct GitHub (used as last resort).
var githubMirrors = []string{
	"https://ghfast.top/",
	"https://gh-proxy.com/",
	"",
}

// releaseInfo represents a GitHub release.
type releaseInfo struct {
	TagName string  `json:"tag_name"`
	Name    string  `json:"name"`
	Body    string  `json:"body"`
	Assets  []asset `json:"assets"`
	HTMLURL string  `json:"html_url"`
}

type asset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
	Size               int64  `json:"size"`
}

const vibecastRepo = "vst93/Vibecast"

// fetchLatestRelease queries GitHub API (via mirrors) for the latest release.
func fetchLatestRelease() (*releaseInfo, error) {
	// The GitHub API URL — mirrors proxy it by prepending their base
	ghAPIURL := "https://api.github.com/repos/" + vibecastRepo + "/releases/latest"

	// Build candidate URLs: mirrors first, then direct
	var apiURLs []string
	for _, mirror := range githubMirrors {
		apiURLs = append(apiURLs, mirror+ghAPIURL)
	}

	var body []byte
	var lastErr error
	for _, url := range apiURLs {
		resp, err := httpClient.Get(url)
		if err != nil {
			lastErr = err
			continue
		}
		if resp.StatusCode != 200 {
			resp.Body.Close()
			lastErr = fmt.Errorf("HTTP %d from %s", resp.StatusCode, url)
			continue
		}
		body, err = io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			lastErr = err
			continue
		}
		// Success
		break
	}
	if body == nil {
		return nil, fmt.Errorf("failed to fetch release info: %w", lastErr)
	}

	var rel releaseInfo
	if err := json.Unmarshal(body, &rel); err != nil {
		return nil, fmt.Errorf("failed to parse release info: %w", err)
	}
	return &rel, nil
}

// findAsset finds the matching binary asset for the current OS/arch.
func findAsset(rel *releaseInfo) *asset {
	goos := runtime.GOOS
	goarch := runtime.GOARCH
	prefix := fmt.Sprintf("vibecast-%s-%s", rel.TagName, goos)
	if goos == "windows" {
		prefix += ".exe"
	}
	for i, a := range rel.Assets {
		if a.Name == prefix || strings.HasPrefix(a.Name, prefix) {
			return &rel.Assets[i]
		}
	}
	// Fallback: match by OS/arch in name (without version)
	for i, a := range rel.Assets {
		name := strings.ToLower(a.Name)
		if strings.Contains(name, goos) && strings.Contains(name, goarch) {
			return &rel.Assets[i]
		}
	}
	return nil
}

// downloadAsset downloads an asset via mirror proxies and returns the local file path.
func downloadAsset(assetURL string) (string, error) {
	// Build candidate download URLs: mirrors first, then direct
	var urls []string
	for _, mirror := range githubMirrors {
		urls = append(urls, mirror+assetURL)
	}

	tmpFile, err := os.CreateTemp("", "vibecast-update-*")
	if err != nil {
		return "", fmt.Errorf("create temp file: %w", err)
	}
	tmpPath := tmpFile.Name()

	var lastErr error
	for _, url := range urls {
		// Reset file
		tmpFile.Truncate(0)
		tmpFile.Seek(0, 0)

		resp, err := httpClient.Get(url)
		if err != nil {
			lastErr = err
			continue
		}
		if resp.StatusCode != 200 {
			resp.Body.Close()
			lastErr = fmt.Errorf("HTTP %d from %s", resp.StatusCode, url)
			continue
		}
		_, err = io.Copy(tmpFile, resp.Body)
		resp.Body.Close()
		if err != nil {
			lastErr = err
			continue
		}
		tmpFile.Close()
		// Success
		return tmpPath, nil
	}

	tmpFile.Close()
	os.Remove(tmpPath)
	return "", fmt.Errorf("failed to download asset: %w", lastErr)
}

// compareVersions compares two version strings (semantic versioning).
// Returns: 1 if a > b, -1 if a < b, 0 if equal.
func compareVersions(a, b string) int {
	// Strip non-alphanumeric prefixes (v, etc.)
	a = strings.TrimLeft(a, "vV ")
	b = strings.TrimLeft(b, "vV ")

	partsA := strings.Split(a, ".")
	partsB := strings.Split(b, ".")

	maxLen := len(partsA)
	if len(partsB) > maxLen {
		maxLen = len(partsB)
	}

	for i := 0; i < maxLen; i++ {
		var na, nb int
		if i < len(partsA) {
			fmt.Sscanf(partsA[i], "%d", &na)
		}
		if i < len(partsB) {
			fmt.Sscanf(partsB[i], "%d", &nb)
		}
		if na > nb {
			return 1
		}
		if na < nb {
			return -1
		}
	}
	return 0
}

// selfReplace replaces the current binary with the new one.
func selfReplace(newBinaryPath string) error {
	currentExe, err := os.Executable()
	if err != nil {
		return fmt.Errorf("get executable path: %w", err)
	}
	currentExe, _ = filepath.EvalSymlinks(currentExe)

	// Make new binary executable
	if err := os.Chmod(newBinaryPath, 0755); err != nil {
		return fmt.Errorf("chmod new binary: %w", err)
	}

	// On Linux, we can't overwrite a running binary directly.
	// Strategy: backup current → copy new into place → clean up backup.
	backupPath := currentExe + ".old"

	// Remove stale backup if exists
	_ = os.Remove(backupPath)

	// Try rename (same filesystem), fall back to copy
	if err := os.Rename(currentExe, backupPath); err == nil {
		// Rename worked (same FS) — now copy new binary into place
		if copyErr := copyFile(newBinaryPath, currentExe); copyErr != nil {
			// Restore backup on failure
			os.Rename(backupPath, currentExe)
			return fmt.Errorf("install new binary: %w", copyErr)
		}
		os.Chmod(currentExe, 0755)
		_ = os.Remove(backupPath)
		return nil
	}

	// Cross-device or rename failed — overwrite in place via copy
	// On Linux, writing to a running binary is fine if we use O_TRUNC
	if err := copyFile(newBinaryPath, currentExe); err != nil {
		return fmt.Errorf("install new binary (copy): %w", err)
	}
	os.Chmod(currentExe, 0755)
	return nil
}

// copyFile copies a file from src to dst, preserving permissions.
func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	info, err := srcFile.Stat()
	if err != nil {
		return err
	}

	dstFile, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, info.Mode())
	if err != nil {
		return err
	}
	defer dstFile.Close()

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return err
	}
	return nil
}

// --- HTTP Handlers ---

// adminCheckUpdate handles GET /api/admin/update/check
func (s *Server) adminCheckUpdate(w http.ResponseWriter, r *http.Request, user *db.User) {
	if r.Method != http.MethodGet {
		writeJSON(w, 405, jsonResp{Error: tMsg(r, "method_not_allowed")})
		return
	}

	currentVersion := s.version
	if currentVersion == "" {
		currentVersion = "dev"
	}

	rel, err := fetchLatestRelease()
	if err != nil {
		writeJSON(w, 200, jsonResp{
			Data: map[string]interface{}{
				"currentVersion": currentVersion,
				"updateAvailable": false,
				"error":           err.Error(),
			},
		})
		return
	}

	latestVersion := strings.TrimPrefix(rel.TagName, "v")
	updateAvailable := false
	if currentVersion != "dev" {
		updateAvailable = compareVersions(latestVersion, currentVersion) > 0
	}

	// Find matching asset
	asset := findAsset(rel)
	var assetInfo map[string]interface{}
	if asset != nil {
		assetInfo = map[string]interface{}{
			"name": asset.Name,
			"size": asset.Size,
		}
	}

	writeJSON(w, 200, jsonResp{
		Data: map[string]interface{}{
			"currentVersion":  currentVersion,
			"latestVersion":   latestVersion,
			"updateAvailable": updateAvailable,
			"releaseName":     rel.Name,
			"releaseNotes":    rel.Body,
			"releaseURL":      rel.HTMLURL,
			"asset":           assetInfo,
		},
	})
}

// adminApplyUpdate handles POST /api/admin/update/apply
func (s *Server) adminApplyUpdate(w http.ResponseWriter, r *http.Request, user *db.User) {
	if r.Method != http.MethodPost {
		writeJSON(w, 405, jsonResp{Error: tMsg(r, "method_not_allowed")})
		return
	}

	currentVersion := s.version
	if currentVersion == "" {
		currentVersion = "dev"
	}

	// Fetch latest release
	rel, err := fetchLatestRelease()
	if err != nil {
		writeJSON(w, 500, jsonResp{Error: tMsg(r, "update_fetch_failed") + ": " + err.Error()})
		return
	}

	// Find matching asset
	asset := findAsset(rel)
	if asset == nil {
		writeJSON(w, 404, jsonResp{Error: tMsg(r, "update_asset_not_found")})
		return
	}

	// Download
	tmpPath, err := downloadAsset(asset.BrowserDownloadURL)
	if err != nil {
		writeJSON(w, 500, jsonResp{Error: tMsg(r, "update_download_failed") + ": " + err.Error()})
		return
	}
	defer os.Remove(tmpPath)

	// Verify the downloaded file is not empty
	info, err := os.Stat(tmpPath)
	if err != nil || info.Size() == 0 {
		writeJSON(w, 500, jsonResp{Error: tMsg(r, "update_download_failed")})
		return
	}

	// Replace the binary
	if err := selfReplace(tmpPath); err != nil {
		writeJSON(w, 500, jsonResp{Error: tMsg(r, "update_install_failed") + ": " + err.Error()})
		return
	}

	writeJSON(w, 200, jsonResp{
		Message: "updated",
		Data: map[string]interface{}{
			"previousVersion": currentVersion,
			"newVersion":      strings.TrimPrefix(rel.TagName, "v"),
		},
	})
}

// --- CLI: `vibecast update` ---

// RunUpdateCLI is called from main.go when `vibecast update` is invoked.
func RunUpdateCLI(currentVersion string) error {
	fmt.Printf("Vibecast v%s\n", currentVersion)
	fmt.Printf("────────────────────────────\n")
	fmt.Printf("Checking for updates...\n")

	rel, err := fetchLatestRelease()
	if err != nil {
		return fmt.Errorf("failed to check for updates: %w", err)
	}

	latestVersion := strings.TrimPrefix(rel.TagName, "v")
	fmt.Printf("Latest release: v%s\n", latestVersion)

	if currentVersion != "dev" && compareVersions(latestVersion, currentVersion) <= 0 {
		fmt.Printf("✓ You are already running the latest version.\n")
		return nil
	}

	if currentVersion == "dev" {
		fmt.Printf("Current version: dev (always allows update)\n")
	}

	fmt.Printf("Update available! v%s → v%s\n", currentVersion, latestVersion)
	if rel.Name != "" {
		fmt.Printf("Release: %s\n", rel.Name)
	}
	if rel.Body != "" {
		fmt.Printf("\n%s\n", rel.Body)
	}
	fmt.Printf("────────────────────────────\n")

	// Find matching asset
	asset := findAsset(rel)
	if asset == nil {
		return fmt.Errorf("no matching binary found for %s/%s", runtime.GOOS, runtime.GOARCH)
	}
	fmt.Printf("Downloading %s (%s)...\n", asset.Name, formatSize(asset.Size))

	// Download
	tmpPath, err := downloadAsset(asset.BrowserDownloadURL)
	if err != nil {
		return fmt.Errorf("download failed: %w", err)
	}
	defer os.Remove(tmpPath)

	// Verify
	info, err := os.Stat(tmpPath)
	if err != nil || info.Size() == 0 {
		return fmt.Errorf("downloaded file is empty or invalid")
	}
	fmt.Printf("✓ Downloaded (%s)\n", formatSize(info.Size()))

	// Replace
	fmt.Printf("Installing...\n")
	if err := selfReplace(tmpPath); err != nil {
		return fmt.Errorf("installation failed: %w", err)
	}

	fmt.Printf("✓ Updated to v%s\n", latestVersion)
	fmt.Printf("Please restart vibecast to apply the update.\n")
	return nil
}

// httpClient is a shared HTTP client with a reasonable timeout.
var httpClient = &http.Client{
	Timeout: 5 * time.Minute, // 5 minutes — allows time for large binary downloads via slow mirrors
}
