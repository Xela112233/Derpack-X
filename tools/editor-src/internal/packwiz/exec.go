package packwiz

import (
	"bytes"
	"fmt"
	"os/exec"
)

// Runner runs packwiz as a subprocess from a given working directory.
// The Binary field should be the absolute path to the packwiz executable
// (typically extracted to .editor/packwiz-cache/ on first run).
type Runner struct {
	Binary string // absolute path to packwiz binary
	WorkingDir string // repo root, where pack.toml lives
}

// Run invokes packwiz with the given args from r.WorkingDir.
// Returns combined stdout+stderr and any exec error.
func (r *Runner) Run(args ...string) (string, error) {
	cmd := exec.Command(r.Binary, args...)
	cmd.Dir = r.WorkingDir

	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf

	if err := cmd.Run(); err != nil {
		return buf.String(), fmt.Errorf("packwiz %v: %w (output: %s)", args, err, buf.String())
	}
	return buf.String(), nil
}

// Refresh runs `packwiz refresh` to update the index after manifest edits.
func (r *Runner) Refresh() (string, error) {
	return r.Run("refresh")
}

// AddModrinth adds a Modrinth mod by slug or project ID.
// Side may be "both", "client", or "server" (or empty for default).
func (r *Runner) AddModrinth(slug, side string) (string, error) {
	args := []string{"mr", "add", slug, "--yes"}
	// packwiz doesn't currently expose --side as a flag for `mr add`.
	// The side has to be set after the fact by editing the manifest if needed.
	_ = side
	return r.Run(args...)
}

// AddCurseForge adds a CurseForge mod by slug.
func (r *Runner) AddCurseForge(slug string) (string, error) {
	return r.Run("cf", "add", slug, "--yes")
}

// Remove deletes a mod's manifest and re-indexes.
func (r *Runner) Remove(slug string) (string, error) {
	return r.Run("remove", slug, "--yes")
}

// Update bumps a single mod to its latest version. Skips pinned mods.
func (r *Runner) Update(slug string) (string, error) {
	return r.Run("update", slug, "--yes")
}

// UpdateAll bumps all mods to their latest versions, skipping pinned ones.
func (r *Runner) UpdateAll() (string, error) {
	return r.Run("update", "--all", "--yes")
}

// Pin marks a mod as pinned (won't be touched by `update --all`).
func (r *Runner) Pin(slug string) (string, error) {
	return r.Run("pin", slug)
}

// Unpin clears the pinned flag.
func (r *Runner) Unpin(slug string) (string, error) {
	return r.Run("unpin", slug)
}
