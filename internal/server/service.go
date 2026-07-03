package server

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// --- Service management for `vibecast service` subcommand ---
//
// Supported platforms:
//   - Linux: systemd
//   - macOS: launchd
//   - Windows: not supported (user must register as a Windows Service manually)
//
// The service runs vibecast with the same flags the user specified at
// registration time, so the binary path, addr, storage, and db are all
// captured into the service unit file.

// serviceConfig holds the resolved binary path and flags to pass to the
// service unit.
type serviceConfig struct {
	exePath    string // absolute path to the vibecast binary
	addr       string
	storageDir string
	dbPath     string
}

// resolveServiceConfig builds a serviceConfig from the CLI flags.
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

// --- systemd (Linux) ---

const systemdUnitTemplate = `[Unit]
Description=Vibecast Static Site Hosting
After=network.target

[Service]
Type=simple
ExecStart=%s --addr %s --storage %s --db %s
Restart=on-failure
RestartSec=5
WorkingDirectory=%s

[Install]
WantedBy=multi-user.target
`

const systemdUnitName = "vibecast.service"
const systemdUnitPath = "/etc/systemd/system/" + systemdUnitName

func systemdInstall(cfg *serviceConfig) error {
	// Resolve working directory (parent of the binary or /opt/vibecast).
	workDir := filepath.Dir(cfg.exePath)

	unit := fmt.Sprintf(systemdUnitTemplate,
		cfg.exePath, cfg.addr, cfg.storageDir, cfg.dbPath, workDir)

	if err := writeFileAsRoot(systemdUnitPath, unit); err != nil {
		return fmt.Errorf("write unit file: %w", err)
	}

	// systemctl daemon-reload + enable + start
	for _, cmd := range [][]string{
		{"systemctl", "daemon-reload"},
		{"systemctl", "enable", "vibecast"},
		{"systemctl", "start", "vibecast"},
	} {
		if out, err := exec.Command(cmd[0], cmd[1:]...).CombinedOutput(); err != nil {
			return fmt.Errorf("%s: %w (%s)", strings.Join(cmd, " "), err, strings.TrimSpace(string(out)))
		}
	}
	return nil
}

func systemdStatus() error {
	cmd := exec.Command("systemctl", "status", "vibecast", "--no-pager")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run() // exit code conveys status
	return nil
}

func systemdStop() error {
	if out, err := exec.Command("systemctl", "stop", "vibecast").CombinedOutput(); err != nil {
		return fmt.Errorf("stop: %w (%s)", err, strings.TrimSpace(string(out)))
	}
	return nil
}

func systemdRestart() error {
	if out, err := exec.Command("systemctl", "restart", "vibecast").CombinedOutput(); err != nil {
		return fmt.Errorf("restart: %w (%s)", err, strings.TrimSpace(string(out)))
	}
	return nil
}

func systemdUninstall() error {
	// Stop + disable first
	for _, cmd := range [][]string{
		{"systemctl", "stop", "vibecast"},
		{"systemctl", "disable", "vibecast"},
	} {
		_ = exec.Command(cmd[0], cmd[1:]...).Run() // ignore errors if not running
	}
	// Remove the unit file
	if err := removeFileAsRoot(systemdUnitPath); err != nil {
		return fmt.Errorf("remove unit file: %w", err)
	}
	// daemon-reload
	_ = exec.Command("systemctl", "daemon-reload").Run()
	return nil
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

func launchdInstall(cfg *serviceConfig) error {
	workDir := filepath.Dir(cfg.exePath)

	plist := fmt.Sprintf(launchdPlistTemplate,
		cfg.exePath, cfg.addr, cfg.storageDir, cfg.dbPath, workDir)

	if err := writeFileAsRoot(launchdPlistPath, plist); err != nil {
		return fmt.Errorf("write plist: %w", err)
	}

	// load + start
	if out, err := exec.Command("launchctl", "load", launchdPlistPath).CombinedOutput(); err != nil {
		return fmt.Errorf("launchctl load: %w (%s)", err, strings.TrimSpace(string(out)))
	}
	return nil
}

func launchdStatus() error {
	// launchctl list prints all loaded services; grep for our label
	out, err := exec.Command("launchctl", "list").CombinedOutput()
	if err != nil {
		return fmt.Errorf("launchctl list: %w", err)
	}
	for _, line := range strings.Split(string(out), "\n") {
		if strings.Contains(line, launchdLabel) {
			fmt.Println(line)
			// Also show recent log
			if logData, err := os.ReadFile("/tmp/vibecast.log"); err == nil && len(logData) > 0 {
				lines := strings.Split(string(logData), "\n")
				start := len(lines) - 15
				if start < 0 {
					start = 0
				}
				fmt.Println("\n--- Recent log ---")
				for _, l := range lines[start:] {
					if l != "" {
						fmt.Println(l)
					}
				}
			}
			return nil
		}
	}
	fmt.Println("Vibecast service is not loaded.")
	return nil
}

func launchdStop() error {
	if out, err := exec.Command("launchctl", "unload", launchdPlistPath).CombinedOutput(); err != nil {
		return fmt.Errorf("launchctl unload: %w (%s)", err, strings.TrimSpace(string(out)))
	}
	return nil
}

func launchdRestart() error {
	// unload then load
	_ = exec.Command("launchctl", "unload", launchdPlistPath).Run()
	if out, err := exec.Command("launchctl", "load", launchdPlistPath).CombinedOutput(); err != nil {
		return fmt.Errorf("launchctl load: %w (%s)", err, strings.TrimSpace(string(out)))
	}
	return nil
}

func launchdUninstall() error {
	// unload first
	_ = exec.Command("launchctl", "unload", launchdPlistPath).Run()
	// remove plist
	if err := removeFileAsRoot(launchdPlistPath); err != nil {
		return fmt.Errorf("remove plist: %w", err)
	}
	return nil
}

// --- helpers for root-owned file operations ---

// writeFileAsRoot writes content to a path, using sudo if the current user
// is not root. On Linux/macOS, service unit files live in /etc or /Library
// which require root.
func writeFileAsRoot(path, content string) error {
	if os.Geteuid() == 0 {
		return os.WriteFile(path, []byte(content), 0644)
	}
	// Use tee via sudo to write the file content
	cmd := exec.Command("sudo", "tee", path)
	cmd.Stdin = strings.NewReader(content)
	cmd.Stdout = os.Stdout // suppress
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	// Set permissions
	_ = exec.Command("sudo", "chmod", "644", path).Run()
	return nil
}

// removeFileAsRoot removes a file, using sudo if not root.
func removeFileAsRoot(path string) error {
	if os.Geteuid() == 0 {
		return os.Remove(path)
	}
	if out, err := exec.Command("sudo", "rm", "-f", path).CombinedOutput(); err != nil {
		return fmt.Errorf("%s (%s)", err, strings.TrimSpace(string(out)))
	}
	return nil
}

// --- RunServiceCLI is called from main.go when `vibecast service` is invoked ---

// RunServiceCLI handles the `vibecast service <action>` subcommand.
// Actions: install, status, stop, restart, uninstall
// The install action captures the current --addr, --storage, --db flags.
func RunServiceCLI(action, addr, storageDir, dbPath string) error {
	if runtime.GOOS == "windows" {
		fmt.Println("⚠ " + TCLIMsg("svc_windows_unsupported"))
		fmt.Println()
		fmt.Println(TCLIMsg("svc_windows_hint"))
		return nil
	}

	switch action {
	case "install":
		cfg, err := resolveServiceConfig(addr, storageDir, dbPath)
		if err != nil {
			return err
		}
		fmt.Printf("%s: %s\n", TCLIMsg("svc_installing"), cfg.exePath)
		fmt.Printf("  --addr %s --storage %s --db %s\n", cfg.addr, cfg.storageDir, cfg.dbPath)
		fmt.Println()

		switch runtime.GOOS {
		case "linux":
			if err := systemdInstall(cfg); err != nil {
				return fmt.Errorf("%s: %w", TCLIMsg("svc_install_failed"), err)
			}
		case "darwin":
			if err := launchdInstall(cfg); err != nil {
				return fmt.Errorf("%s: %w", TCLIMsg("svc_install_failed"), err)
			}
		default:
			return fmt.Errorf("%s", TCLIMsg("svc_unsupported"))
		}

		fmt.Printf("✓ %s\n", TCLIMsg("svc_installed"))
		fmt.Printf("  vibecast service status   # %s\n", TCLIMsg("svc_status_cmd"))
		fmt.Printf("  vibecast service stop     # %s\n", TCLIMsg("svc_stop_cmd"))
		fmt.Printf("  vibecast service restart  # %s\n", TCLIMsg("svc_restart_cmd"))
		fmt.Printf("  vibecast service uninstall # %s\n", TCLIMsg("svc_uninstall_cmd"))
		return nil

	case "status":
		switch runtime.GOOS {
		case "linux":
			return systemdStatus()
		case "darwin":
			return launchdStatus()
		default:
			return fmt.Errorf("%s", TCLIMsg("svc_unsupported"))
		}

	case "stop":
		switch runtime.GOOS {
		case "linux":
			if err := systemdStop(); err != nil {
				return err
			}
		case "darwin":
			if err := launchdStop(); err != nil {
				return err
			}
		default:
			return fmt.Errorf("%s", TCLIMsg("svc_unsupported"))
		}
		fmt.Printf("✓ %s\n", TCLIMsg("svc_stopped"))
		return nil

	case "restart":
		switch runtime.GOOS {
		case "linux":
			if err := systemdRestart(); err != nil {
				return err
			}
		case "darwin":
			if err := launchdRestart(); err != nil {
				return err
			}
		default:
			return fmt.Errorf("%s", TCLIMsg("svc_unsupported"))
		}
		fmt.Printf("✓ %s\n", TCLIMsg("svc_restarted"))
		return nil

	case "uninstall":
		switch runtime.GOOS {
		case "linux":
			if err := systemdUninstall(); err != nil {
				return err
			}
		case "darwin":
			if err := launchdUninstall(); err != nil {
				return err
			}
		default:
			return fmt.Errorf("%s", TCLIMsg("svc_unsupported"))
		}
		fmt.Printf("✓ %s\n", TCLIMsg("svc_uninstalled"))
		return nil

	default:
		return fmt.Errorf("%s: %s", TCLIMsg("svc_unknown_action"), action)
	}
}
