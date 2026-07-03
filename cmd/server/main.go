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
	// Custom usage so config flags show --addr style, action flags show -version style.
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: vibecast [options]\n\nOptions:\n")
		fmt.Fprintf(os.Stderr, "  --addr <addr>\n    	listen address (default \":8080\", env VIBECAST_ADDR)\n")
		fmt.Fprintf(os.Stderr, "  --storage <dir>\n    	site files storage directory (default \"./data/sites\", env VIBECAST_STORAGE)\n")
		fmt.Fprintf(os.Stderr, "  --db <path>\n    	SQLite database path (default \"./data/vibecast.db\", env VIBECAST_DB)\n")
		fmt.Fprintf(os.Stderr, "  -v, -version\n    	print version and exit\n")
		fmt.Fprintf(os.Stderr, "  -update\n    	check for updates and self-update\n")
		fmt.Fprintf(os.Stderr, "  -h, -help\n    	show this help message\n")
	}

	addr := flag.String("addr", getEnv("VIBECAST_ADDR", ":8080"), "listen address")
	storageDir := flag.String("storage", getEnv("VIBECAST_STORAGE", "./data/sites"), "site files storage directory")
	dbPath := flag.String("db", getEnv("VIBECAST_DB", "./data/vibecast.db"), "SQLite database path")
	versionFlag := flag.Bool("version", false, "print version and exit")
	vShort := flag.Bool("v", false, "shorthand for -version")
	updateFlag := flag.Bool("update", false, "check for updates and self-update")
	flag.Parse()

	if *versionFlag || *vShort {
		fmt.Printf("vibecast v%s\n", version)
		return
	}

	if *updateFlag {
		if err := server.RunUpdateCLI(version); err != nil {
			log.Fatalf("Update failed: %v", err)
		}
		return
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
