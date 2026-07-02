package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"vibecast/internal/server"
)

var version = "dev"

func main() {
	addr := flag.String("addr", getEnv("VIBECAST_ADDR", ":8080"), "listen address")
	storageDir := flag.String("storage", getEnv("VIBECAST_STORAGE", "./data/sites"), "site files storage directory")
	dbPath := flag.String("db", getEnv("VIBECAST_DB", "./data/vibecast.db"), "SQLite database path")
	versionFlag := flag.Bool("version", false, "print version and exit")
	flag.Parse()

	if *versionFlag {
		fmt.Printf("vibecast v%s\n", version)
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

	if err := http.ListenAndServe(*addr, srv.Router()); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
