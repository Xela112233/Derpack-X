// Package builder assembles a local Prism Launcher instance from the current
// repo state. Equivalent to scripts/build-prism-bundled.sh but in Go, so it
// works on Windows without bash.
//
// The build does NOT zip — it produces a directory structure ready to be
// copied into the user's existing Prism instance folder (or imported as-is).
//
// High-level flow:
//   1. Read pack metadata from pack.toml.
//   2. Write Prism scaffold files (mmc-pack.json, instance.cfg) into a staging
//      folder.
//   3. Copy config/, defaultconfigs/, kubejs/, resourcepacks/, shaderpacks/
//      into <staging>/.minecraft/.
//   4. Start `packwiz serve` on a free port.
//   5. Run packwiz-installer-bootstrap (Java) against that server to fetch
//      every mod jar into <staging>/.minecraft/mods/.
//   6. Stop the server.
//   7. Clean up bootstrap leftovers.
package builder

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/derpack/derpack-edit/internal/packwiz"
)

// Result describes a completed build.
type Result struct {
	StagingDir string // absolute path to the assembled instance dir
	ModCount   int
	Logs       []string // human-readable progress log
}

// Builder runs builds. Reuse one across requests.
type Builder struct {
	RepoRoot     string
	PackwizBin   string // absolute path to packwiz executable
	BootstrapJar string // absolute path to packwiz-installer-bootstrap.jar (downloaded on first call if absent)
}

// Build assembles a Prism instance into <RepoRoot>/.editor/build-staging/.
// Returns the path on success.
func (b *Builder) Build(ctx context.Context, log func(string)) (*Result, error) {
	logf := func(format string, args ...any) {
		msg := fmt.Sprintf(format, args...)
		if log != nil {
			log(msg)
		}
	}
	res := &Result{Logs: []string{}}
	collect := func(format string, args ...any) {
		msg := fmt.Sprintf(format, args...)
		res.Logs = append(res.Logs, msg)
		logf(msg)
	}

	// 1. Read pack metadata.
	pack, err := packwiz.LoadPack(b.RepoRoot)
	if err != nil {
		return nil, fmt.Errorf("read pack.toml: %w", err)
	}
	mcVersion := pack.Versions["minecraft"]
	neoforgeVersion := pack.Versions["neoforge"]
	if mcVersion == "" || neoforgeVersion == "" {
		return nil, fmt.Errorf("pack.toml missing minecraft or neoforge version")
	}

	staging := filepath.Join(b.RepoRoot, ".editor", "build-staging")
	if err := os.RemoveAll(staging); err != nil {
		return nil, fmt.Errorf("clean staging: %w", err)
	}
	mcDir := filepath.Join(staging, ".minecraft")
	if err := os.MkdirAll(mcDir, 0o755); err != nil {
		return nil, fmt.Errorf("mkdir staging: %w", err)
	}
	collect("Staging dir: %s", staging)

	// 2. Scaffold Prism instance files.
	if err := writeMmcPack(staging, mcVersion, neoforgeVersion); err != nil {
		return nil, err
	}
	if err := writeInstanceCfg(staging, pack.Name, pack.Version); err != nil {
		return nil, err
	}
	collect("Wrote mmc-pack.json and instance.cfg")

	// 3. Copy config dirs.
	for _, d := range []string{"config", "defaultconfigs", "kubejs", "resourcepacks", "shaderpacks"} {
		src := filepath.Join(b.RepoRoot, d)
		if _, err := os.Stat(src); err != nil {
			continue // skip missing dirs
		}
		dst := filepath.Join(mcDir, d)
		if err := copyDir(src, dst); err != nil {
			return nil, fmt.Errorf("copy %s: %w", d, err)
		}
		collect("Copied %s/", d)
	}

	// 4. Ensure bootstrap jar present.
	if err := b.ensureBootstrapJar(); err != nil {
		return nil, fmt.Errorf("get installer bootstrap: %w", err)
	}

	// 5. Run packwiz serve + bootstrap to fetch jars.
	port, stop, err := startPackwizServe(ctx, b.PackwizBin, b.RepoRoot)
	if err != nil {
		return nil, fmt.Errorf("start packwiz serve: %w", err)
	}
	defer stop()
	collect("Started packwiz serve on :%d", port)

	if err := waitForServer(ctx, port); err != nil {
		return nil, fmt.Errorf("packwiz serve never came up: %w", err)
	}

	collect("Fetching mod jars (this may take a few minutes)...")
	if err := runBootstrap(ctx, b.BootstrapJar, mcDir, port); err != nil {
		return nil, fmt.Errorf("install mods: %w", err)
	}

	// 6. Clean up bootstrap leftovers.
	for _, f := range []string{"packwiz.json", "packwiz-installer-bootstrap.jar"} {
		_ = os.Remove(filepath.Join(mcDir, f))
	}

	// 7. Sanity check.
	jars, _ := filepath.Glob(filepath.Join(mcDir, "mods", "*.jar"))
	res.ModCount = len(jars)
	if res.ModCount == 0 {
		return nil, fmt.Errorf("no mod jars after install — something went wrong fetching")
	}
	collect("Bundled %d mod jars", res.ModCount)

	res.StagingDir = staging
	return res, nil
}

// ----- helpers ----------------------------------------------------------

func writeMmcPack(staging, mcVer, neoforgeVer string) error {
	contents := fmt.Sprintf(`{
    "components": [
        {
            "important": true,
            "uid": "net.minecraft",
            "version": "%s"
        },
        {
            "uid": "net.neoforged",
            "version": "%s"
        }
    ],
    "formatVersion": 1
}
`, mcVer, neoforgeVer)
	return os.WriteFile(filepath.Join(staging, "mmc-pack.json"), []byte(contents), 0o644)
}

func writeInstanceCfg(staging, packName, packVersion string) error {
	contents := fmt.Sprintf(`[General]
ConfigVersion=1.2
InstanceType=OneSix
JavaArchitecture=64
JavaVersion=21
JvmArgs=-XX:+UseG1GC -XX:+ParallelRefProcEnabled -XX:MaxGCPauseMillis=200 -XX:+UnlockExperimentalVMOptions -XX:+DisableExplicitGC -XX:+AlwaysPreTouch -XX:G1NewSizePercent=30 -XX:G1MaxNewSizePercent=40 -XX:G1HeapRegionSize=8M -XX:G1ReservePercent=20 -XX:G1HeapWastePercent=5 -XX:G1MixedGCCountTarget=4 -XX:InitiatingHeapOccupancyPercent=15 -XX:G1MixedGCLiveThresholdPercent=90 -XX:G1RSetUpdatingPauseTimePercent=5 -XX:SurvivorRatio=32 -XX:+PerfDisableSharedMem -XX:MaxTenuringThreshold=1
MaxMemAlloc=12288
MinMemAlloc=8192
OverrideJavaArgs=true
OverrideMemory=true
PermGen=256
iconKey=default
name=%s %s
notes=Built locally from %s by derpack-edit.
`, packName, packVersion, packName)
	return os.WriteFile(filepath.Join(staging, "instance.cfg"), []byte(contents), 0o644)
}

// copyDir recursively copies srcDir into dstDir.
func copyDir(srcDir, dstDir string) error {
	if err := os.MkdirAll(dstDir, 0o755); err != nil {
		return err
	}
	return filepath.Walk(srcDir, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(srcDir, p)
		dst := filepath.Join(dstDir, rel)
		if info.IsDir() {
			return os.MkdirAll(dst, 0o755)
		}
		return copyFile(p, dst, info.Mode())
	})
}

func copyFile(src, dst string, mode os.FileMode) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}
	out, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, mode)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}

// ensureBootstrapJar downloads packwiz-installer-bootstrap.jar to .editor/
// if it isn't already there. Idempotent.
func (b *Builder) ensureBootstrapJar() error {
	if b.BootstrapJar == "" {
		b.BootstrapJar = filepath.Join(b.RepoRoot, ".editor", "packwiz-installer-bootstrap.jar")
	}
	if _, err := os.Stat(b.BootstrapJar); err == nil {
		return nil
	}

	const url = "https://github.com/packwiz/packwiz-installer-bootstrap/releases/latest/download/packwiz-installer-bootstrap.jar"
	if err := os.MkdirAll(filepath.Dir(b.BootstrapJar), 0o755); err != nil {
		return err
	}

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("download bootstrap: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("download bootstrap: %s", resp.Status)
	}

	out, err := os.Create(b.BootstrapJar)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, resp.Body)
	return err
}

// startPackwizServe starts `packwiz serve` on a free port and returns the port
// plus a stop function.
func startPackwizServe(ctx context.Context, packwizBin, repoRoot string) (port int, stop func(), err error) {
	port, err = freePort()
	if err != nil {
		return 0, nil, err
	}

	cmd := exec.CommandContext(ctx, packwizBin, "serve", "--port", fmt.Sprint(port))
	cmd.Dir = repoRoot
	// Capture output for debug; not currently surfaced.
	cmd.Stdout = io.Discard
	cmd.Stderr = io.Discard

	if err := cmd.Start(); err != nil {
		return 0, nil, err
	}

	stop = func() {
		if cmd.Process != nil {
			_ = cmd.Process.Kill()
		}
		_ = cmd.Wait()
	}
	return port, stop, nil
}

func freePort() (int, error) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return 0, err
	}
	defer ln.Close()
	return ln.Addr().(*net.TCPAddr).Port, nil
}

func waitForServer(ctx context.Context, port int) error {
	url := fmt.Sprintf("http://localhost:%d/pack.toml", port)
	deadline := time.Now().Add(15 * time.Second)
	for time.Now().Before(deadline) {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		resp, err := http.Get(url)
		if err == nil {
			resp.Body.Close()
			if resp.StatusCode == 200 {
				return nil
			}
		}
		time.Sleep(500 * time.Millisecond)
	}
	return fmt.Errorf("timeout waiting for packwiz serve on :%d", port)
}

// runBootstrap invokes `java -jar packwiz-installer-bootstrap.jar -g -s client <pack-url>`
// from the .minecraft directory, fetching all mods.
func runBootstrap(ctx context.Context, jarPath, mcDir string, port int) error {
	packURL := fmt.Sprintf("http://localhost:%d/pack.toml", port)

	javaBin := "java"
	cmd := exec.CommandContext(ctx, javaBin, "-jar", jarPath, "-g", "-s", "client", packURL)
	cmd.Dir = mcDir

	out, err := cmd.CombinedOutput()
	if err != nil {
		// Most common failure: java not in PATH. Surface a clear message.
		if strings.Contains(err.Error(), "executable file not found") {
			return fmt.Errorf("java not found in PATH. Install Java 21+ and ensure it's on PATH, or open Prism once (Prism bundles Java) and configure your PATH to include it. Underlying: %w", err)
		}
		return fmt.Errorf("bootstrap failed: %w (output: %s)", err, string(out))
	}
	return nil
}
