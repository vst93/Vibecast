package server

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync/atomic"
	"syscall"
	"time"

	"vibecast/internal/db"
)

// GitHub mirror hosts for China mainland — tried in order for binary downloads.
// Each entry is a prefix prepended to the full GitHub URL.
// Empty string = direct GitHub (used as last resort).
// NOTE: mirrors only proxy release downloads, NOT the GitHub API.
var githubMirrors = []string{
	"https://ghfast.top/",
	"https://gh-proxy.com/",
	"",
}

// updateInProgress is an atomic flag (0 = idle, 1 = running) preventing
// concurrent update operations from the admin API.
var updateInProgress int32

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

// fetchLatestRelease queries GitHub API for the latest release.
// Strategy: try direct api.github.com first, then fall back to mirrors.
func fetchLatestRelease() (*releaseInfo, error) {
	ghAPIURL := "https://api.github.com/repos/" + vibecastRepo + "/releases/latest"

	// Build candidate URLs: direct first, then mirrors as fallback.
	// (Mirrors return 403 for API calls, but direct may fail in CN networks.)
	var apiURLs []string
	apiURLs = append(apiURLs, ghAPIURL) // direct first
	for _, mirror := range githubMirrors {
		if mirror == "" {
			continue // already added direct above
		}
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
		break // success
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
	// Expected name: vibecast-{version}-{os}-{arch}[.exe]
	prefix := fmt.Sprintf("vibecast-%s-%s-%s", rel.TagName, goos, goarch)
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

// downloadAsset downloads an asset via mirror proxies (mirror-first, direct fallback)
// and returns the local file path. If totalSize > 0 and progress != nil, reports
// download progress via the callback.
func downloadAsset(assetURL string, totalSize int64, progress func(downloaded, total int64)) (string, error) {
	// Build candidate download URLs: mirrors first, then direct.
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

		var reader io.Reader = resp.Body
		total := totalSize
		if total <= 0 {
			total = resp.ContentLength
		}
		if progress != nil && total > 0 {
			reader = &progressReader{r: resp.Body, total: total, fn: progress}
		}
		_, err = io.Copy(tmpFile, reader)
		resp.Body.Close()
		if err != nil {
			lastErr = err
			continue
		}
		tmpFile.Close()
		return tmpPath, nil
	}

	tmpFile.Close()
	os.Remove(tmpPath)
	return "", fmt.Errorf("failed to download asset: %w", lastErr)
}

// downloadAndVerifyAsset downloads the binary, then verifies it against
// SHA256SUMS if available. Returns the path to the verified temp file.
func downloadAndVerifyAsset(assetURL, assetName string) (string, error) {
	tmpPath, err := downloadAsset(assetURL, 0, nil)
	if err != nil {
		return "", err
	}

	// Try to fetch SHA256SUMS from the same release download base.
	// assetURL looks like: https://github.com/vst93/Vibecast/releases/download/VERSION/assetName
	// SHA256SUMS is at:    .../releases/download/VERSION/SHA256SUMS
	sumsURL := assetURL[:strings.LastIndex(assetURL, "/")+1] + "SHA256SUMS"

	sumsOK, verifyErr := verifySHA256(tmpPath, assetName, sumsURL)
	if !sumsOK {
		if verifyErr == errNoChecksum {
			// Old release without SHA256SUMS — skip but warn.
			fmt.Fprintf(os.Stderr, "WARNING: %s\n", tStatic("updateNoChecksum"))
		} else {
			// Checksum mismatch — delete and fail.
			os.Remove(tmpPath)
			return "", fmt.Errorf("%s: %w", tStatic("updateVerifyFailed"), verifyErr)
		}
	}

	return tmpPath, nil
}

var errNoChecksum = fmt.Errorf("no checksum file available")

// verifySHA256 downloads SHA256SUMS, extracts the hash for assetName, computes
// the SHA256 of localFile, and compares. Returns (true, nil) on match,
// (true, err) on mismatch, (false, errNoChecksum) if sums file unavailable.
func verifySHA256(localFile, assetName, sumsURL string) (bool, error) {
	// Try mirrors then direct for the sums file.
	var urls []string
	for _, mirror := range githubMirrors {
		urls = append(urls, mirror+sumsURL)
	}

	var sumsBody []byte
	var got bool
	for _, url := range urls {
		resp, err := httpClient.Get(url)
		if err != nil {
			continue
		}
		if resp.StatusCode != 200 {
			resp.Body.Close()
			continue
		}
		sumsBody, err = io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			continue
		}
		got = true
		break
	}

	if !got {
		return false, errNoChecksum
	}

	// Extract expected hash for this asset.
	expectedHash := extractSHA256FromSums(string(sumsBody), assetName)
	if expectedHash == "" {
		return false, errNoChecksum
	}

	// Compute actual hash.
	f, err := os.Open(localFile)
	if err != nil {
		return true, fmt.Errorf("open for hashing: %w", err)
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return true, fmt.Errorf("hashing: %w", err)
	}
	actualHash := hex.EncodeToString(h.Sum(nil))

	if !strings.EqualFold(actualHash, expectedHash) {
		return true, fmt.Errorf("hash mismatch: expected %s, got %s", expectedHash, actualHash)
	}
	return true, nil
}

// extractSHA256FromSums parses a SHA256SUMS file and returns the hash for the
// given asset name. Format: "<hash>  <filename>" per line.
func extractSHA256FromSums(sums, assetName string) string {
	for _, line := range strings.Split(sums, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) >= 2 && fields[1] == assetName {
			return fields[0]
		}
	}
	return ""
}

// progressReader wraps an io.Reader and reports progress via callback.
type progressReader struct {
	r        io.Reader
	total    int64
	read     int64
	fn       func(downloaded, total int64)
	lastRep  time.Time
}

func (pr *progressReader) Read(p []byte) (int, error) {
	n, err := pr.r.Read(p)
	pr.read += int64(n)
	// Throttle progress callback to every 100ms.
	now := time.Now()
	if n > 0 && now.Sub(pr.lastRep) >= 100*time.Millisecond {
		pr.fn(pr.read, pr.total)
		pr.lastRep = now
	}
	if err == io.EOF && pr.fn != nil {
		pr.fn(pr.read, pr.total) // final
	}
	return n, err
}

// compareVersions compares two version strings.
// Supports both semantic versioning (1.2.3) and date-based versioning (YYYYMMDD-N).
// Returns: 1 if a > b, -1 if a < b, 0 if equal.
func compareVersions(a, b string) int {
	a = strings.TrimLeft(a, "vV ")
	b = strings.TrimLeft(b, "vV ")

	// Date-based format: YYYYMMDD-N (e.g. "20260703-6")
	// Split on '-' and compare date first, then suffix number.
	if strings.Contains(a, "-") || strings.Contains(b, "-") {
		partsA := strings.SplitN(a, "-", 2)
		partsB := strings.SplitN(b, "-", 2)
		var dateA, dateB int
		fmt.Sscanf(partsA[0], "%d", &dateA)
		fmt.Sscanf(partsB[0], "%d", &dateB)
		if dateA != dateB {
			if dateA > dateB {
				return 1
			}
			return -1
		}
		// Same date — compare suffix number (default 0 if no suffix).
		var subA, subB int
		if len(partsA) > 1 {
			fmt.Sscanf(partsA[1], "%d", &subA)
		}
		if len(partsB) > 1 {
			fmt.Sscanf(partsB[1], "%d", &subB)
		}
		if subA > subB {
			return 1
		}
		if subA < subB {
			return -1
		}
		return 0
	}

	// Semantic versioning: split on '.'
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
// On Windows, the running binary is locked, so we rename the old binary
// aside and copy the new one in. On Unix, we can overwrite in place.
func selfReplace(newBinaryPath string) error {
	currentExe, err := os.Executable()
	if err != nil {
		return fmt.Errorf("get executable path: %w", err)
	}
	// Resolve symlinks, but keep the original path if resolution fails
	// (e.g. on some macOS or container setups where /proc/self/exe is weird).
	resolved, err := filepath.EvalSymlinks(currentExe)
	if err == nil && resolved != "" {
		currentExe = resolved
	}
	if currentExe == "" {
		return fmt.Errorf("could not determine executable path")
	}

	if err := os.Chmod(newBinaryPath, 0755); err != nil {
		return fmt.Errorf("chmod new binary: %w", err)
	}

	// Windows: cannot overwrite a running .exe — must rename it aside first.
	if runtime.GOOS == "windows" {
		backupPath := currentExe + ".old"
		_ = os.Remove(backupPath)
		if err := os.Rename(currentExe, backupPath); err != nil {
			return fmt.Errorf("cannot replace running binary (close vibecast first): %w", err)
		}
		if copyErr := copyFile(newBinaryPath, currentExe); copyErr != nil {
			os.Rename(backupPath, currentExe) // restore
			return fmt.Errorf("install new binary: %w", copyErr)
		}
		os.Chmod(currentExe, 0755)
		_ = os.Remove(backupPath)
		return nil
	}

	// Unix: try rename (same filesystem → fast + atomic), fallback to copy.
	backupPath := currentExe + ".old"
	_ = os.Remove(backupPath)

	if err := os.Rename(currentExe, backupPath); err == nil {
		// Rename succeeded — copy new binary into place.
		if copyErr := copyFile(newBinaryPath, currentExe); copyErr != nil {
			os.Rename(backupPath, currentExe) // restore old binary
			return fmt.Errorf("install new binary: %w", copyErr)
		}
		os.Chmod(currentExe, 0755)
		_ = os.Remove(backupPath)
		return nil
	}

	// Rename failed (cross-device or permission) — overwrite in place via copy.
	// On Linux, writing to a running binary is fine with O_TRUNC.
	if err := copyFile(newBinaryPath, currentExe); err != nil {
		return fmt.Errorf("install new binary (copy): %w", err)
	}
	os.Chmod(currentExe, 0755)
	return nil
}

// copyFile copies a file from src to dst, preserving permissions.
func copyFile(src, dst string) error {
	if src == "" {
		return fmt.Errorf("source path is empty")
	}
	if dst == "" {
		return fmt.Errorf("destination path is empty")
	}
	srcFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("open source: %w", err)
	}
	defer srcFile.Close()

	info, err := srcFile.Stat()
	if err != nil {
		return fmt.Errorf("stat source: %w", err)
	}

	dstFile, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, info.Mode())
	if err != nil {
		// Common on Linux when the binary is installed in a system directory
		// (e.g. /usr/local/bin) and the process doesn't have write permission.
		if os.IsPermission(err) {
			return fmt.Errorf("permission denied writing to %s — try running as root or the same user that owns the binary", dst)
		}
		return fmt.Errorf("open destination: %w", err)
	}
	defer dstFile.Close()

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return fmt.Errorf("copy data: %w", err)
	}
	// Sync to ensure data is flushed to disk before we return.
	if err := dstFile.Sync(); err != nil {
		// Sync failure is not fatal on all platforms (e.g. Windows).
		_ = err
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
				"currentVersion":  currentVersion,
				"updateAvailable": false,
				"error":           err.Error(),
			},
		})
		return
	}

	latestVersion := strings.TrimPrefix(rel.TagName, "v")
	updateAvailable := false
	if currentVersion == "dev" {
		// Dev builds can always update to a release.
		updateAvailable = true
	} else {
		updateAvailable = compareVersions(latestVersion, currentVersion) > 0
	}

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

	// Concurrency guard: only one update at a time.
	if !atomic.CompareAndSwapInt32(&updateInProgress, 0, 1) {
		writeJSON(w, 409, jsonResp{Error: tMsg(r, "updateInProgress")})
		return
	}
	defer atomic.StoreInt32(&updateInProgress, 0)

	currentVersion := s.version
	if currentVersion == "" {
		currentVersion = "dev"
	}

	rel, err := fetchLatestRelease()
	if err != nil {
		writeJSON(w, 500, jsonResp{Error: tMsg(r, "update_fetch_failed") + ": " + err.Error()})
		return
	}

	asset := findAsset(rel)
	if asset == nil {
		writeJSON(w, 404, jsonResp{Error: tMsg(r, "update_asset_not_found")})
		return
	}

	// Download + verify SHA256.
	tmpPath, err := downloadAndVerifyAsset(asset.BrowserDownloadURL, asset.Name)
	if err != nil {
		writeJSON(w, 500, jsonResp{Error: tMsg(r, "update_download_failed") + ": " + err.Error()})
		return
	}
	defer os.Remove(tmpPath)

	info, err := os.Stat(tmpPath)
	if err != nil || info.Size() == 0 {
		writeJSON(w, 500, jsonResp{Error: tMsg(r, "update_download_failed")})
		return
	}

	if err := selfReplace(tmpPath); err != nil {
		errStr := err.Error()
		// Provide platform-specific guidance for common failures.
		if strings.Contains(errStr, "permission denied") {
			writeJSON(w, 403, jsonResp{Error: tMsg(r, "update_permission_denied")})
			return
		}
		if runtime.GOOS == "windows" && (strings.Contains(errStr, "being used by another process") || strings.Contains(errStr, "Access is denied") || strings.Contains(errStr, "cannot replace running binary")) {
			writeJSON(w, 409, jsonResp{Error: tMsg(r, "update_windows_locked")})
			return
		}
		writeJSON(w, 500, jsonResp{Error: tMsg(r, "update_install_failed") + ": " + errStr})
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

// adminRestartUpdate handles POST /api/admin/update/restart
// Restarts the server process by gracefully shutting down then exec'ing the
// new binary in-place (Unix only). On Windows, returns a message to restart manually.
func (s *Server) adminRestartUpdate(w http.ResponseWriter, r *http.Request, user *db.User) {
	if r.Method != http.MethodPost {
		writeJSON(w, 405, jsonResp{Error: tMsg(r, "method_not_allowed")})
		return
	}

	if runtime.GOOS == "windows" {
		writeJSON(w, 200, jsonResp{
			Message: "manual_restart_required",
			Data:    map[string]interface{}{"platform": "windows"},
		})
		return
	}

	// Respond to client BEFORE shutting down, so the HTTP response is sent.
	writeJSON(w, 200, jsonResp{Message: "restarting"})

	// Shut down in a goroutine so the response is flushed first.
	go func() {
		time.Sleep(500 * time.Millisecond) // let response flush

		// Gracefully close the database.
		s.database.Close()

		// Close the HTTP server listener so the port is released.
		if s.httpServer != nil {
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()
			_ = s.httpServer.Shutdown(ctx)
		}

		// Replace the current process with the new binary.
		exe, err := os.Executable()
		if err != nil {
			fmt.Fprintf(os.Stderr, "restart: cannot find executable: %v\n", err)
			return
		}
		if err := syscall.Exec(exe, os.Args, os.Environ()); err != nil {
			fmt.Fprintf(os.Stderr, "restart: exec failed: %v\n", err)
		}
	}()
}

// getLocalIPs returns all non-loopback IPv4 addresses of the server.
func getLocalIPs() []string {
	var ips []string
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ips
	}
	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				ips = append(ips, ipNet.IP.String())
			}
		}
	}
	return ips
}

// adminSystemInfo handles GET /api/admin/system-info
func (s *Server) adminSystemInfo(w http.ResponseWriter, r *http.Request, user *db.User) {
	if r.Method != http.MethodGet {
		writeJSON(w, 405, jsonResp{Error: tMsg(r, "method_not_allowed")})
		return
	}
	v := s.version
	if v == "" {
		v = "dev"
	}
	writeJSON(w, 200, jsonResp{Data: map[string]interface{}{
		"version":     v,
		"storagePath": s.config.StorageDir,
		"dbPath":      s.config.DBPath,
		"listenAddr":  s.config.Addr,
		"goVersion":   runtime.Version(),
		"os":          runtime.GOOS,
		"arch":        runtime.GOARCH,
		"startTime":   startTime,
		"localIPs":    getLocalIPs(),
	}})
}

// startTime is set at package init.
var startTime = time.Now().Format(time.RFC3339)

// --- CLI: `vibecast update` ---

// RunUpdateCLI is called from main.go when `vibecast update` is invoked.
func RunUpdateCLI(currentVersion string) error {
	fmt.Printf("Vibecast v%s\n", currentVersion)
	fmt.Printf("────────────────────────────\n")
	fmt.Printf("%s\n", TCLIMsg("cli_checking"))

	rel, err := fetchLatestRelease()
	if err != nil {
		return fmt.Errorf("%s: %w", TCLIMsg("cli_fetch_failed"), err)
	}

	latestVersion := strings.TrimPrefix(rel.TagName, "v")
	fmt.Printf("%s v%s\n", TCLIMsg("cli_latest_rel"), latestVersion)

	if currentVersion != "dev" && compareVersions(latestVersion, currentVersion) <= 0 {
		fmt.Printf("✓ %s\n", TCLIMsg("cli_up_to_date"))
		return nil
	}

	if currentVersion == "dev" {
		fmt.Printf("%s\n", TCLIMsg("cli_dev_version"))
	}

	fmt.Printf("%s v%s → v%s\n", TCLIMsg("cli_update_avail"), currentVersion, latestVersion)
	if rel.Name != "" {
		fmt.Printf("%s %s\n", TCLIMsg("cli_release"), rel.Name)
	}
	if rel.Body != "" {
		fmt.Printf("\n%s\n", rel.Body)
	}
	fmt.Printf("────────────────────────────\n")

	asset := findAsset(rel)
	if asset == nil {
		return fmt.Errorf("%s %s/%s", TCLIMsg("cli_no_binary"), runtime.GOOS, runtime.GOARCH)
	}
	fmt.Printf("%s %s (%s)...\n", TCLIMsg("cli_downloading"), asset.Name, formatSize(asset.Size))

	// Download with progress.
	tmpPath, err := downloadAsset(asset.BrowserDownloadURL, asset.Size, func(downloaded, total int64) {
		pct := float64(downloaded) / float64(total) * 100
		fmt.Printf("\r%s... %.0f%% [%s/%s]   ", TCLIMsg("cli_downloading"), pct, formatSize(downloaded), formatSize(total))
	})
	if err != nil {
		return fmt.Errorf("%s: %w", TCLIMsg("cli_dl_failed"), err)
	}
	defer os.Remove(tmpPath)
	fmt.Printf("\r✓ %s (%s)                    \n", TCLIMsg("cli_downloaded"), formatSize(asset.Size))

	// Verify SHA256.
	sumsURL := asset.BrowserDownloadURL[:strings.LastIndex(asset.BrowserDownloadURL, "/")+1] + "SHA256SUMS"
	sumsOK, verifyErr := verifySHA256(tmpPath, asset.Name, sumsURL)
	if !sumsOK {
		if verifyErr == errNoChecksum {
			fmt.Printf("⚠ %s\n", TCLIMsg("updateNoChecksum"))
		} else {
			return fmt.Errorf("%s: %w", TCLIMsg("updateVerifyFailed"), verifyErr)
		}
	} else {
		fmt.Printf("✓ %s\n", TCLIMsg("cli_checksum_ok"))
	}

	info, err := os.Stat(tmpPath)
	if err != nil || info.Size() == 0 {
		return fmt.Errorf("%s", TCLIMsg("cli_empty_file"))
	}

	fmt.Printf("%s\n", TCLIMsg("cli_installing"))
	if err := selfReplace(tmpPath); err != nil {
		errStr := err.Error()
		if strings.Contains(errStr, "permission denied") {
			if runtime.GOOS == "windows" {
				return fmt.Errorf("%s", TCLIMsg("update_windows_locked"))
			}
			return fmt.Errorf("%s", TCLIMsg("update_permission_denied"))
		}
		if runtime.GOOS == "windows" && (strings.Contains(errStr, "being used by another process") || strings.Contains(errStr, "Access is denied") || strings.Contains(errStr, "cannot replace running binary")) {
			return fmt.Errorf("%s", TCLIMsg("update_windows_locked"))
		}
		return fmt.Errorf("%s: %w", TCLIMsg("cli_install_failed"), err)
	}

	fmt.Printf("✓ %s v%s\n", TCLIMsg("cli_updated"), latestVersion)
	fmt.Printf("%s\n", TCLIMsg("cli_restart_hint"))
	return nil
}

// httpClient is a shared HTTP client with a reasonable timeout.
var httpClient = &http.Client{
	Timeout: 5 * time.Minute,
}
