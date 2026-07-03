package server

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// --- Service registration (Hermes-style) ---
//
// `vibecast setup` writes a service unit file and enables it. After that,
// the user manages the service with standard system commands:
//
//   Linux:   systemctl --user status/stop/restart vibecast
//   macOS:   launchctl list/stop/start com.vibecast
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
	return &serviceConfig{
		exePath:    exe,
		addr:       addr,
		storageDir: storageDir,
		dbPath:     dbPath,
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

	workDir := filepath.Dir(cfg.exePath)
	unit := fmt.Sprintf(systemdUnitTemplate,
		cfg.exePath, cfg.addr, cfg.storageDir, cfg.dbPath, workDir)

	unitPath := systemdUnitPath()
	if err := os.WriteFile(unitPath, []byte(unit), 0644); err != nil {
		return fmt.Errorf("write unit file: %w", err)
	}

	// Enable linger so the user service survives logout.
	enableLinger()

	// daemon-reload + enable --now
	for _, args := range [][]string{
		{"--user", "daemon-reload"},
		{"--user", "enable", "--now", "vibecast"},
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
    <key>StandardOutPath</key>
    <string>/tmp/vibecast.log</string>
    <key>StandardErrorPath</key>
    <string>/tmp/vibecast-error.log</string>
</dict>
</plist>
`

const launchdLabel = "com.vibecast"
const launchdPlistPath = "/Library/LaunchDaemons/com.vibecast.plist"

func launchdSetup(cfg *serviceConfig) error {
	workDir := filepath.Dir(cfg.exePath)
	plist := fmt.Sprintf(launchdPlistTemplate,
		cfg.exePath, cfg.addr, cfg.storageDir, cfg.dbPath, workDir)

	if err := writeFileAsRoot(launchdPlistPath, plist); err != nil {
		return fmt.Errorf("write plist: %w", err)
	}
	if out, err := exec.Command("launchctl", "load", launchdPlistPath).CombinedOutput(); err != nil {
		return fmt.Errorf("launchctl load: %w (%s)", err, strings.TrimSpace(string(out)))
	}
	return nil
}

func launchdTeardown() error {
	_ = exec.Command("launchctl", "unload", launchdPlistPath).Run()
	if err := removeFileAsRoot(launchdPlistPath); err != nil {
		return fmt.Errorf("remove plist: %w", err)
	}
	return nil
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
