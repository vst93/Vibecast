package server

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

// --- Service registration (Hermes-style) ---
//
// `vibecast setup` writes a service unit file and enables it. After that,
// the user manages the service with standard system commands:
//
//   Linux:   systemctl --user status/stop/restart vibecast
//   macOS:   sudo launchctl list/stop/start com.vibecast
//
// Only `setup` (install) and `uninstall` are handled here — no wrapper
// commands for status/stop/restart. The user already knows systemctl.
//
// Windows is not supported; the user is told to use nssm or Task Scheduler.

// serviceConfig holds the resolved binary path and flags for the service unit.
type serviceConfig struct {
	exePath    string
	addr       string
	storageDir string
	dbPath     string
	homeDir    string
	userName   string
}

// ResolveUserPaths makes relative data paths independent of the process
// working directory. Relative paths are rooted at the current user's home so
// a service restart from /usr/local/bin cannot silently open a new database.
// Explicit absolute paths remain unchanged for backwards compatibility.
func ResolveUserPaths(storageDir, dbPath string) (string, string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", "", fmt.Errorf("cannot determine user home: %w", err)
	}
	storageDir, err = resolveUserPath(home, storageDir)
	if err != nil {
		return "", "", fmt.Errorf("resolve storage path: %w", err)
	}
	dbPath, err = resolveUserPath(home, dbPath)
	if err != nil {
		return "", "", fmt.Errorf("resolve database path: %w", err)
	}
	return storageDir, dbPath, nil
}

func resolveUserPath(home, path string) (string, error) {
	if path == "" {
		return path, nil
	}
	if path == "~" {
		return filepath.Clean(home), nil
	}
	if strings.HasPrefix(path, "~/") || strings.HasPrefix(path, `~\`) {
		path = path[2:]
	}
	if filepath.IsAbs(path) {
		return filepath.Clean(path), nil
	}
	resolved := filepath.Clean(filepath.Join(home, path))
	rel, err := filepath.Rel(home, resolved)
	if err != nil {
		return "", err
	}
	if rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("relative path %q escapes user home", path)
	}
	return resolved, nil
}

func resolveServiceConfig(addr, storageDir, dbPath string) (*serviceConfig, error) {
	exe, err := os.Executable()
	if err != nil {
		return nil, fmt.Errorf("cannot determine executable path: %w", err)
	}
	exe, err = filepath.Abs(exe)
	if err != nil {
		return nil, fmt.Errorf("cannot resolve executable path: %w", err)
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("cannot determine user home: %w", err)
	}
	currentUser, err := user.Current()
	if err != nil {
		return nil, fmt.Errorf("cannot determine current user: %w", err)
	}
	storageDir, dbPath, err = ResolveUserPaths(storageDir, dbPath)
	if err != nil {
		return nil, err
	}
	return &serviceConfig{
		exePath:    exe,
		addr:       addr,
		storageDir: storageDir,
		dbPath:     dbPath,
		homeDir:    home,
		userName:   currentUser.Username,
	}, nil
}

// --- systemd user service (Linux) ---

const systemdUnitTemplate = `[Unit]
Description=Vibecast Static Site Hosting
After=network-online.target
Wants=network-online.target

[Service]
Type=simple
ExecStart=%s --addr %s --storage %s --db %s
WorkingDirectory=%s
Restart=on-failure
RestartSec=5
StandardOutput=journal
StandardError=journal

[Install]
WantedBy=default.target
`

const systemdUnitName = "vibecast.service"

func systemdQuote(value string) string {
	// Percent has specifier semantics in systemd, including inside quotes.
	return strconv.Quote(strings.ReplaceAll(value, "%", "%%"))
}

func renderSystemdUnit(cfg *serviceConfig) string {
	return fmt.Sprintf(systemdUnitTemplate,
		systemdQuote(cfg.exePath), systemdQuote(cfg.addr),
		systemdQuote(cfg.storageDir), systemdQuote(cfg.dbPath),
		systemdQuote(cfg.homeDir))
}

func systemdUserDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "systemd", "user")
}

func systemdUnitPath() string {
	return filepath.Join(systemdUserDir(), systemdUnitName)
}

func systemdSetup(cfg *serviceConfig) error {
	dir := systemdUserDir()
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("create systemd user dir: %w", err)
	}

	unit := renderSystemdUnit(cfg)

	unitPath := systemdUnitPath()
	if err := os.WriteFile(unitPath, []byte(unit), 0644); err != nil {
		return fmt.Errorf("write unit file: %w", err)
	}

	// Enable linger so the user service survives logout.
	enableLinger()

	// Reload and restart even when setup is rerun, so an existing service picks
	// up newly resolved absolute paths immediately.
	for _, args := range [][]string{
		{"--user", "daemon-reload"},
		{"--user", "enable", "vibecast"},
		{"--user", "restart", "vibecast"},
	} {
		if out, err := exec.Command("systemctl", args...).CombinedOutput(); err != nil {
			return fmt.Errorf("systemctl %s: %w (%s)",
				strings.Join(args, " "), err, strings.TrimSpace(string(out)))
		}
	}
	return nil
}

func systemdTeardown() error {
	// stop + disable (ignore errors if not running/installed)
	for _, args := range [][]string{
		{"--user", "stop", "vibecast"},
		{"--user", "disable", "vibecast"},
	} {
		_ = exec.Command("systemctl", args...).Run()
	}
	// remove unit file
	unitPath := systemdUnitPath()
	if err := os.Remove(unitPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("remove unit file: %w", err)
	}
	// daemon-reload
	_ = exec.Command("systemctl", "--user", "daemon-reload").Run()
	return nil
}

// enableLinger enables systemd linger for the current user so user services
// survive logout. Uses loginctl; ignores failure (not all systems have it).
func enableLinger() {
	user := os.Getenv("USER")
	if user == "" {
		user = os.Getenv("LOGNAME")
	}
	if user == "" {
		return
	}
	// loginctl enable-linger may need sudo on some systems; try without first.
	if err := exec.Command("loginctl", "enable-linger", user).Run(); err != nil {
		// Try with sudo
		_ = exec.Command("sudo", "loginctl", "enable-linger", user).Run()
	}
}

// --- launchd (macOS) ---

const launchdPlistTemplate = `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>com.vibecast</string>
    <key>UserName</key>
    <string>%s</string>
    <key>ProgramArguments</key>
    <array>
        <string>%s</string>
        <string>--addr</string>
        <string>%s</string>
        <string>--storage</string>
        <string>%s</string>
        <string>--db</string>
        <string>%s</string>
    </array>
    <key>RunAtLoad</key>
    <true/>
    <key>KeepAlive</key>
    <true/>
    <key>WorkingDirectory</key>
    <string>%s</string>
    <key>EnvironmentVariables</key>
    <dict>
        <key>HOME</key>
        <string>%s</string>
    </dict>
    <key>StandardOutPath</key>
    <string>%s</string>
    <key>StandardErrorPath</key>
    <string>%s</string>
</dict>
</plist>
`

const launchdLabel = "com.vibecast"
const launchdPlistPath = "/Library/LaunchDaemons/com.vibecast.plist"

func xmlEscape(value string) string {
	var buf bytes.Buffer
	_ = xml.EscapeText(&buf, []byte(value))
	return buf.String()
}

func renderLaunchdPlist(cfg *serviceConfig) string {
	logDir := filepath.Join(cfg.homeDir, "Library", "Logs")
	return fmt.Sprintf(launchdPlistTemplate,
		xmlEscape(cfg.userName), xmlEscape(cfg.exePath), xmlEscape(cfg.addr), xmlEscape(cfg.storageDir),
		xmlEscape(cfg.dbPath), xmlEscape(cfg.homeDir), xmlEscape(cfg.homeDir),
		xmlEscape(filepath.Join(logDir, "vibecast.log")),
		xmlEscape(filepath.Join(logDir, "vibecast-error.log")))
}

func launchdSetup(cfg *serviceConfig) error {
	if err := os.MkdirAll(filepath.Join(cfg.homeDir, "Library", "Logs"), 0755); err != nil {
		return fmt.Errorf("create launchd log directory: %w", err)
	}
	// setup is intentionally idempotent: unload the old definition before
	// replacing it so rerunning setup applies changed paths and user identity.
	_, _ = runAsRoot("launchctl", "unload", launchdPlistPath)

	plist := renderLaunchdPlist(cfg)

	if err := writeFileAsRoot(launchdPlistPath, plist); err != nil {
		return fmt.Errorf("write plist: %w", err)
	}
	if out, err := runAsRoot("launchctl", "load", launchdPlistPath); err != nil {
		return fmt.Errorf("launchctl load: %w (%s)", err, strings.TrimSpace(string(out)))
	}
	return nil
}

func launchdTeardown() error {
	_, _ = runAsRoot("launchctl", "unload", launchdPlistPath)
	if err := removeFileAsRoot(launchdPlistPath); err != nil {
		return fmt.Errorf("remove plist: %w", err)
	}
	return nil
}

func runAsRoot(name string, args ...string) ([]byte, error) {
	if os.Geteuid() == 0 {
		return exec.Command(name, args...).CombinedOutput()
	}
	sudoArgs := append([]string{name}, args...)
	return exec.Command("sudo", sudoArgs...).CombinedOutput()
}

// --- root-owned file helpers (macOS only; Linux uses user-level dir) ---

func writeFileAsRoot(path, content string) error {
	if os.Geteuid() == 0 {
		return os.WriteFile(path, []byte(content), 0644)
	}
	cmd := exec.Command("sudo", "tee", path)
	cmd.Stdin = strings.NewReader(content)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	_ = exec.Command("sudo", "chmod", "644", path).Run()
	return nil
}

func removeFileAsRoot(path string) error {
	if os.Geteuid() == 0 {
		return os.Remove(path)
	}
	if out, err := exec.Command("sudo", "rm", "-f", path).CombinedOutput(); err != nil {
		return fmt.Errorf("%s (%s)", err, strings.TrimSpace(string(out)))
	}
	return nil
}

// --- RunSetupCLI handles `vibecast setup` (install service) ---
// --- RunUninstallCLI handles `vibecast uninstall` (remove service) ---

func RunSetupCLI(addr, storageDir, dbPath string) error {
	if runtime.GOOS == "windows" {
		fmt.Println("⚠ " + TCLIMsg("svc_windows_unsupported"))
		fmt.Println()
		fmt.Println(TCLIMsg("svc_windows_hint"))
		return nil
	}

	cfg, err := resolveServiceConfig(addr, storageDir, dbPath)
	if err != nil {
		return err
	}

	fmt.Printf("%s: %s\n", TCLIMsg("svc_installing"), cfg.exePath)
	fmt.Printf("  --addr %s --storage %s --db %s\n", cfg.addr, cfg.storageDir, cfg.dbPath)
	fmt.Println()

	switch runtime.GOOS {
	case "linux":
		if err := systemdSetup(cfg); err != nil {
			return fmt.Errorf("%s: %w", TCLIMsg("svc_install_failed"), err)
		}
	case "darwin":
		if err := launchdSetup(cfg); err != nil {
			return fmt.Errorf("%s: %w", TCLIMsg("svc_install_failed"), err)
		}
	default:
		return fmt.Errorf("%s", TCLIMsg("svc_unsupported"))
	}

	fmt.Printf("✓ %s\n", TCLIMsg("svc_installed"))
	fmt.Println()
	// Tell the user how to manage the service with standard commands.
	if runtime.GOOS == "linux" {
		fmt.Println(TCLIMsg("svc_manage_hint_linux"))
	} else if runtime.GOOS == "darwin" {
		fmt.Println(TCLIMsg("svc_manage_hint_macos"))
	}
	fmt.Println()
	fmt.Printf("  vibecast uninstall  # %s\n", TCLIMsg("svc_uninstall_cmd"))
	return nil
}

func RunUninstallCLI() error {
	if runtime.GOOS == "windows" {
		fmt.Println("⚠ " + TCLIMsg("svc_windows_unsupported"))
		return nil
	}

	switch runtime.GOOS {
	case "linux":
		if err := systemdTeardown(); err != nil {
			return fmt.Errorf("%s: %w", TCLIMsg("svc_uninstall_failed"), err)
		}
	case "darwin":
		if err := launchdTeardown(); err != nil {
			return fmt.Errorf("%s: %w", TCLIMsg("svc_uninstall_failed"), err)
		}
	default:
		return fmt.Errorf("%s", TCLIMsg("svc_unsupported"))
	}

	fmt.Printf("✓ %s\n", TCLIMsg("svc_uninstalled"))
	return nil
}
