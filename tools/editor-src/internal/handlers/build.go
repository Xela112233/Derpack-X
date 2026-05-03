package handlers

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/derpack/derpack-edit/internal/builder"
)

// HandleBuild builds a Prism instance into <repo>/.editor/build-staging/.
// The instance directory is staged but NOT copied to the user's Prism folder
// — that's a separate launch step. This is useful for testing/inspecting the
// staged output.
func (s *Server) HandleBuild(w http.ResponseWriter, r *http.Request) {
	if !requirePost(w, r) {
		return
	}

	b := &builder.Builder{
		RepoRoot:   s.RepoRoot,
		PackwizBin: s.Runner.Binary,
	}

	// Use a long timeout — bundling 200+ mods over slow connections takes time.
	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Minute)
	defer cancel()

	res, err := b.Build(ctx, nil)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"ok":          true,
		"staging_dir": res.StagingDir,
		"mod_count":   res.ModCount,
		"logs":        res.Logs,
	})
}

// HandleLaunchPrism builds + syncs the result into the configured Prism
// instance folder. After this completes, the user opens Prism manually and
// clicks Play (we don't try to launch Prism itself).
func (s *Server) HandleLaunchPrism(w http.ResponseWriter, r *http.Request) {
	if !requirePost(w, r) {
		return
	}
	if s.Config.PrismInstancePath == "" {
		writeError(w, http.StatusBadRequest,
			"Prism instance path is not configured. Set it in the config first.")
		return
	}
	dst := s.Config.PrismInstancePath
	if _, err := os.Stat(dst); err != nil {
		writeError(w, http.StatusBadRequest,
			"configured Prism path does not exist: "+dst)
		return
	}

	b := &builder.Builder{
		RepoRoot:   s.RepoRoot,
		PackwizBin: s.Runner.Binary,
	}

	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Minute)
	defer cancel()

	res, err := b.Build(ctx, nil)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "build failed: "+err.Error())
		return
	}

	// Sync .minecraft/mods/ from staging into the Prism instance.
	srcMods := filepath.Join(res.StagingDir, ".minecraft", "mods")
	dstMc := resolveMinecraftDir(dst)
	dstMods := filepath.Join(dstMc, "mods")

	if err := replaceModsDir(srcMods, dstMods); err != nil {
		writeError(w, http.StatusInternalServerError,
			fmt.Sprintf("synced build to staging but failed to copy to Prism: %v", err))
		return
	}

	// Also sync configs that affect gameplay or visuals.
	// Matches the dirs that the builder stages, so the Prism instance ends up
	// reflecting the full repo state.
	for _, d := range []string{"config", "defaultconfigs", "kubejs", "resourcepacks", "shaderpacks"} {
		src := filepath.Join(res.StagingDir, ".minecraft", d)
		if _, err := os.Stat(src); err != nil {
			continue
		}
		dst := filepath.Join(dstMc, d)
		// Replace the dir to ensure removed configs go away, but we don't nuke
		// the whole .minecraft — only these specific dirs.
		if err := os.RemoveAll(dst); err != nil {
			writeError(w, http.StatusInternalServerError,
				fmt.Sprintf("clear %s in instance: %v", d, err))
			return
		}
		if err := copyDirShallow(src, dst); err != nil {
			writeError(w, http.StatusInternalServerError,
				fmt.Sprintf("copy %s to instance: %v", d, err))
			return
		}
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"ok":           true,
		"mod_count":    res.ModCount,
		"prism_path":   dst,
		"hint":         "Open Prism Launcher and click Play on the configured instance.",
	})
}

// resolveMinecraftDir returns the .minecraft subdir of a Prism instance,
// falling back to the path itself if there's no .minecraft. Prism instances
// always have .minecraft, but if the user pointed at .minecraft directly we
// honor that.
func resolveMinecraftDir(instancePath string) string {
	mc := filepath.Join(instancePath, ".minecraft")
	if _, err := os.Stat(mc); err == nil {
		return mc
	}
	return instancePath
}

// replaceModsDir wipes dst, then copies all jars from src into it. Uses
// copyDirShallow because mods/ has no subdirectories.
func replaceModsDir(src, dst string) error {
	if err := os.RemoveAll(dst); err != nil {
		return err
	}
	if err := os.MkdirAll(dst, 0o755); err != nil {
		return err
	}
	return copyDirShallow(src, dst)
}

// copyDirShallow copies src to dst recursively, preserving structure.
// (Same logic as builder.copyDir but local to the handler package to avoid
// exposing it from builder for now.)
func copyDirShallow(src, dst string) error {
	if err := os.MkdirAll(dst, 0o755); err != nil {
		return err
	}
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}
	for _, e := range entries {
		s := filepath.Join(src, e.Name())
		d := filepath.Join(dst, e.Name())
		if e.IsDir() {
			if err := copyDirShallow(s, d); err != nil {
				return err
			}
			continue
		}
		if err := copyFile(s, d); err != nil {
			return err
		}
	}
	return nil
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	return err
}
