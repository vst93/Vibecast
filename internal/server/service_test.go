package server

import (
	"encoding/xml"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestResolveUserPath(t *testing.T) {
	home := filepath.Join(t.TempDir(), "vibecast user")
	abs := filepath.Join(t.TempDir(), "vibecast", "data.db")

	tests := []struct {
		name string
		path string
		want string
	}{
		{name: "dot relative", path: "./data/sites", want: filepath.Join(home, "data", "sites")},
		{name: "home shorthand", path: "~/data/vibecast.db", want: filepath.Join(home, "data", "vibecast.db")},
		{name: "home itself", path: "~", want: home},
		{name: "absolute unchanged", path: abs, want: abs},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := resolveUserPath(home, tt.path)
			if err != nil {
				t.Fatalf("resolveUserPath() error = %v", err)
			}
			if got != tt.want {
				t.Fatalf("resolveUserPath() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestResolveUserPathRejectsEscape(t *testing.T) {
	home := filepath.Join(t.TempDir(), "vibecast")
	if _, err := resolveUserPath(home, filepath.Join("..", "outside.db")); err == nil {
		t.Fatal("resolveUserPath() accepted a relative path outside user home")
	}
}

func TestRenderSystemdUnitUsesQuotedAbsolutePaths(t *testing.T) {
	cfg := &serviceConfig{
		exePath:    "/home/vibe cast/bin/vibe%cast",
		addr:       ":8080",
		storageDir: "/home/vibe cast/data/sites",
		dbPath:     "/home/vibe cast/data/vibecast.db",
		homeDir:    "/home/vibe cast",
	}

	unit := renderSystemdUnit(cfg)
	want := `ExecStart="/home/vibe cast/bin/vibe%%cast" --addr ":8080" --storage "/home/vibe cast/data/sites" --db "/home/vibe cast/data/vibecast.db"`
	if !strings.Contains(unit, want) {
		t.Fatalf("systemd unit does not contain safely quoted ExecStart:\n%s", unit)
	}
	if !strings.Contains(unit, `WorkingDirectory="/home/vibe cast"`) {
		t.Fatalf("systemd unit does not use the user home as working directory:\n%s", unit)
	}
}

func TestRenderLaunchdPlistEscapesValuesAndSetsUser(t *testing.T) {
	cfg := &serviceConfig{
		exePath:    "/Users/a&b/bin/vibecast",
		addr:       "<:8080>",
		storageDir: "/Users/a&b/data/sites",
		dbPath:     "/Users/a&b/data/vibecast.db",
		homeDir:    "/Users/a&b",
		userName:   "a&b",
	}

	plist := renderLaunchdPlist(cfg)
	for _, want := range []string{
		"<key>UserName</key>\n    <string>a&amp;b</string>",
		"<string>/Users/a&amp;b/data/vibecast.db</string>",
		"<string>&lt;:8080&gt;</string>",
		"<key>WorkingDirectory</key>\n    <string>/Users/a&amp;b</string>",
		"<key>HOME</key>\n        <string>/Users/a&amp;b</string>",
		"<string>/Users/a&amp;b/Library/Logs/vibecast.log</string>",
	} {
		if !strings.Contains(plist, want) {
			t.Fatalf("launchd plist missing %q:\n%s", want, plist)
		}
	}

	decoder := xml.NewDecoder(strings.NewReader(plist))
	for {
		if _, err := decoder.Token(); err != nil {
			if err == io.EOF {
				break
			}
			t.Fatalf("launchd plist is not valid XML: %v", err)
		}
	}
}

func TestRestartExecutablePathAndArgs(t *testing.T) {
	dir := t.TempDir()
	exe := filepath.Join(dir, "vibecast")
	if err := os.WriteFile(exe, []byte("binary"), 0755); err != nil {
		t.Fatal(err)
	}

	gotExe, err := restartExecutablePath(exe)
	if err != nil {
		t.Fatalf("restartExecutablePath() error = %v", err)
	}
	if gotExe != exe {
		t.Fatalf("restartExecutablePath() = %q, want %q", gotExe, exe)
	}

	cfg := &Config{
		Addr:       ":18099",
		StorageDir: filepath.Join(dir, "data", "sites"),
		DBPath:     filepath.Join(dir, "data", "vibecast.db"),
	}
	want := []string{exe, "--addr", cfg.Addr, "--storage", cfg.StorageDir, "--db", cfg.DBPath}
	if got := restartArgs(exe, cfg); !reflect.DeepEqual(got, want) {
		t.Fatalf("restartArgs() = %#v, want %#v", got, want)
	}
}

func TestRestartState(t *testing.T) {
	s := &Server{restartCh: make(chan struct{})}
	if s.IsRestarting() {
		t.Fatal("new server unexpectedly reports a pending restart")
	}
	s.beginRestart()
	s.beginRestart()
	if !s.IsRestarting() {
		t.Fatal("server did not retain pending restart state")
	}
}
