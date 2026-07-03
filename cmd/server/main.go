package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"vibecast/internal/server"
)

var version = "dev"

func main() {
	// Custom usage: Options (config flags with --) and Commands (subcommands).
	flag.CommandLine.SetOutput(os.Stderr)
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: vibecast [options] [command]\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		fmt.Fprintf(os.Stderr, "  --addr <addr>\n    \tlisten address (default \":8080\", env VIBECAST_ADDR)\n")
		fmt.Fprintf(os.Stderr, "  --storage <dir>\n    \tsite files storage directory (default \"./data/sites\", env VIBECAST_STORAGE)\n")
		fmt.Fprintf(os.Stderr, "  --db <path>\n    \tSQLite database path (default \"./data/vibecast.db\", env VIBECAST_DB)\n")
		fmt.Fprintf(os.Stderr, "\nCommands:\n")
		fmt.Fprintf(os.Stderr, "  version, v   print version and exit\n")
		fmt.Fprintf(os.Stderr, "  update       check for updates and self-update\n")
		fmt.Fprintf(os.Stderr, "  help, h      show this help message\n")
	}

	addr := flag.String("addr", getEnv("VIBECAST_ADDR", ":8080"), "listen address")
	storageDir := flag.String("storage", getEnv("VIBECAST_STORAGE", "./data/sites"), "site files storage directory")
	dbPath := flag.String("db", getEnv("VIBECAST_DB", "./data/vibecast.db"), "SQLite database path")
	flag.Parse()

	// Subcommands (bare words, no dash prefix)
	args := flag.Args()
	if len(args) > 0 {
		switch args[0] {
		case "version", "v":
			fmt.Printf("vibecast v%s\n", version)
			return
		case "update":
			if err := server.RunUpdateCLI(version); err != nil {
				log.Fatalf("Update failed: %v", err)
			}
			return
		case "help", "h":
			flag.Usage()
			return
		default:
			fmt.Fprintf(os.Stderr, "unknown command: %s\n\n", args[0])
			flag.Usage()
			os.Exit(2)
		}
	}

	cfg := &server.Config{
		Addr:       *addr,
		StorageDir: *storageDir,
		DBPath:     *dbPath,
		Version:    version,
	}

	srv, err := server.New(cfg)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}
	defer srv.Close()

	// Graceful shutdown
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh
		fmt.Println("\nShutting down...")
		srv.Close()
		os.Exit(0)
	}()

	fmt.Printf("Vibecast v%s\n", version)
	fmt.Printf("Build with vibe. Cast instantly.\n")
	fmt.Printf("────────────────────────────\n")
	fmt.Printf("Listening:  http://localhost%s\n", *addr)
	fmt.Printf("Storage:   %s\n", *storageDir)
	fmt.Printf("Database:  %s\n", *dbPath)
	fmt.Printf("────────────────────────────\n")
	fmt.Printf("Dashboard: http://localhost%s/dashboard\n", *addr)

	// Use net.Listen + http.Server so we can gracefully shut down / restart.
	ln, err := net.Listen("tcp", *addr)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	hs := &http.Server{Handler: srv.Router()}
	srv.SetHTTPServer(hs)
	if err := hs.Serve(ln); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server failed: %v", err)
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
