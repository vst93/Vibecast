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
		fmt.Fprintf(os.Stderr, "%s\n\n", server.TCLIMsg("cli_usage"))
		fmt.Fprintf(os.Stderr, "%s\n", server.TCLIMsg("cli_options"))
		fmt.Fprintf(os.Stderr, "  --addr <addr>\n    	%s (default \":8080\", env VIBECAST_ADDR)\n", server.TCLIMsg("cli_addr"))
		fmt.Fprintf(os.Stderr, "  --storage <dir>\n    	%s (default \"./data/sites\", env VIBECAST_STORAGE)\n", server.TCLIMsg("cli_storage"))
		fmt.Fprintf(os.Stderr, "  --db <path>\n    	%s (default \"./data/vibecast.db\", env VIBECAST_DB)\n", server.TCLIMsg("cli_db"))
		fmt.Fprintf(os.Stderr, "\n%s\n", server.TCLIMsg("cli_commands"))
		fmt.Fprintf(os.Stderr, "  version, v   %s\n", server.TCLIMsg("cli_version_cmd"))
		fmt.Fprintf(os.Stderr, "  update       %s\n", server.TCLIMsg("cli_update_cmd"))
		fmt.Fprintf(os.Stderr, "  service      %s\n", server.TCLIMsg("cli_service_cmd"))
		fmt.Fprintf(os.Stderr, "  help, h      %s\n", server.TCLIMsg("cli_help_cmd"))
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
		case "service":
			// `vibecast service <action>` — install/status/stop/restart/uninstall
			serviceAction := ""
			if len(args) > 1 {
				serviceAction = args[1]
			}
			if serviceAction == "" {
				fmt.Fprintf(os.Stderr, "%s\n\n", server.TCLIMsg("svc_usage"))
				fmt.Fprintf(os.Stderr, "  vibecast service install    # %s\n", server.TCLIMsg("svc_install_desc"))
				fmt.Fprintf(os.Stderr, "  vibecast service status     # %s\n", server.TCLIMsg("svc_status_desc"))
				fmt.Fprintf(os.Stderr, "  vibecast service stop       # %s\n", server.TCLIMsg("svc_stop_desc"))
				fmt.Fprintf(os.Stderr, "  vibecast service restart    # %s\n", server.TCLIMsg("svc_restart_desc"))
				fmt.Fprintf(os.Stderr, "  vibecast service uninstall  # %s\n", server.TCLIMsg("svc_uninstall_desc"))
				return
			}
			if err := server.RunServiceCLI(serviceAction, *addr, *storageDir, *dbPath); err != nil {
				log.Fatalf("Service: %v", err)
			}
			return
		case "help", "h":
			flag.Usage()
			return
		default:
			fmt.Fprintf(os.Stderr, "%s: %s\n\n", server.TCLIMsg("cli_unknown_cmd"), args[0])
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
	fmt.Printf("%s  http://localhost%s\n", server.TCLIMsg("cli_listening"), *addr)
	fmt.Printf("%s   %s\n", server.TCLIMsg("cli_storage_label"), *storageDir)
	fmt.Printf("%s  %s\n", server.TCLIMsg("cli_database"), *dbPath)
	fmt.Printf("────────────────────────────\n")
	fmt.Printf("%s http://localhost%s/dashboard\n", server.TCLIMsg("cli_dashboard"), *addr)

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
