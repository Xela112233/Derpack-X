package main

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"
	"time"

	"github.com/derpack/derpack-edit/internal/config"
	"github.com/derpack/derpack-edit/internal/handlers"
	"github.com/derpack/derpack-edit/internal/packwiz"
)

//go:embed ui/*
var uiAssets embed.FS

const (
	preferredPort = 8765
	maxPortTries  = 10
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("derpack-edit: %v", err)
	}
}

func run() error {
	// Resolve repo root. If we're run as `tools/derpack-edit.exe` from the
	// repo root, the working dir IS the repo root. If invoked elsewhere, we
	// fall back to assuming the binary lives in <repo>/tools/.
	repoRoot, err := resolveRepoRoot()
	if err != nil {
		return err
	}
	log.Printf("repo root: %s", repoRoot)

	// Extract the embedded packwiz binary if needed.
	binPath, err := packwiz.EnsureBinary(repoRoot)
	if err != nil {
		return fmt.Errorf("setup packwiz: %w", err)
	}
	log.Printf("packwiz: %s", binPath)

	// Load saved config (Prism path, etc.).
	cfg, err := config.Load(repoRoot)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	srv := &handlers.Server{
		RepoRoot: repoRoot,
		Runner:   &packwiz.Runner{Binary: binPath, WorkingDir: repoRoot},
		Config:   cfg,
	}

	// Find a port to bind to.
	listener, port, err := openListener()
	if err != nil {
		return fmt.Errorf("listen: %w", err)
	}

	mux := http.NewServeMux()

	// API endpoints.
	mux.HandleFunc("/api/pack", srv.HandleGetPack)
	mux.HandleFunc("/api/mods", srv.HandleListMods)
	mux.HandleFunc("/api/mods/add", srv.HandleAddMod)
	mux.HandleFunc("/api/mods/remove", srv.HandleRemoveMod)
	mux.HandleFunc("/api/mods/pin", srv.HandlePin)
	mux.HandleFunc("/api/mods/set-version", srv.HandleSetVersion)
	mux.HandleFunc("/api/mods/compute-hash", srv.HandleComputeHash)
	mux.HandleFunc("/api/mods/update", srv.HandleUpdate)
	mux.HandleFunc("GET /api/config", srv.HandleGetConfig)
	mux.HandleFunc("POST /api/config", srv.HandleSetConfig)
	mux.HandleFunc("/api/build", srv.HandleBuild)
	mux.HandleFunc("/api/launch-prism", srv.HandleLaunchPrism)
	// Pattern routes (Go 1.22+) for variable path segments.
	mux.HandleFunc("GET /api/mods/{slug}/versions", srv.HandleListVersions)

	// Static UI served from the embedded assets.
	uiSubFS, err := fs.Sub(uiAssets, "ui")
	if err != nil {
		return fmt.Errorf("ui assets: %w", err)
	}
	mux.Handle("/", http.FileServer(http.FS(uiSubFS)))

	httpSrv := &http.Server{
		Handler:      mux,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 60 * time.Second,
	}

	// Run server in a goroutine so we can handle shutdown.
	errCh := make(chan error, 1)
	go func() {
		errCh <- httpSrv.Serve(listener)
	}()

	url := fmt.Sprintf("http://localhost:%d", port)
	fmt.Printf("\n  Derpack Editor running at %s\n", url)
	fmt.Printf("  Press Ctrl+C to stop.\n\n")

	// Try to open the browser. Best-effort; don't fail if it doesn't work.
	if err := openBrowser(url); err != nil {
		log.Printf("could not open browser automatically: %v", err)
		fmt.Printf("  Open the URL above in your browser manually.\n\n")
	}

	// Wait for either Ctrl+C or a server error.
	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-stopCh:
		fmt.Printf("\nReceived %s, shutting down...\n", sig)
	case err := <-errCh:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("server: %w", err)
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return httpSrv.Shutdown(ctx)
}

// resolveRepoRoot finds the directory containing pack.toml, starting from CWD
// and walking up. Falls back to the binary's parent's parent (assuming
// .../tools/derpack-edit.exe) if no pack.toml is found from CWD.
func resolveRepoRoot() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	if root := findPackToml(cwd); root != "" {
		return root, nil
	}

	// Try the binary's directory.
	exe, err := os.Executable()
	if err == nil {
		exeDir := filepath.Dir(exe)
		if root := findPackToml(exeDir); root != "" {
			return root, nil
		}
	}

	return "", fmt.Errorf("could not find pack.toml in current directory or any parent. Run derpack-edit from the modpack repo root, or place it in <repo>/tools/")
}

func findPackToml(start string) string {
	dir := start
	for i := 0; i < 10; i++ { // walk up at most 10 levels
		if _, err := os.Stat(filepath.Join(dir, "pack.toml")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return ""
		}
		dir = parent
	}
	return ""
}

func openListener() (net.Listener, int, error) {
	for i := 0; i < maxPortTries; i++ {
		port := preferredPort + i
		ln, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
		if err == nil {
			return ln, port, nil
		}
	}
	return nil, 0, fmt.Errorf("could not bind to any port in %d-%d", preferredPort, preferredPort+maxPortTries-1)
}

// openBrowser tries to open the given URL in the user's default browser.
// Best-effort; returns an error if the launch couldn't be initiated.
func openBrowser(url string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", "", url)
	case "darwin":
		cmd = exec.Command("open", url)
	default:
		cmd = exec.Command("xdg-open", url)
	}
	return cmd.Start()
}
