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
	Binary     string // absolute path to packwiz binary
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
//
// Note: packwiz doesn't currently expose --side as a flag for `mr add`.
// The side has to be set after the fact by editing the manifest if needed —
// see handlers.applySide.
func (r *Runner) AddModrinth(slug, side string) (string, error) {
	_ = side
	return r.Run("mr", "add", slug, "--yes")
}

// CurseForge mods are NOT added via packwiz CLI. The bare `packwiz cf add`
// invocation has no path to pass our stored CF Core API key, and when
// packwiz can't authenticate it falls back to writing manifests with empty
// download URLs that crash packwiz-installer at launch time.
//
// CF adds go through handlers.addCurseForgeMod instead, which uses the
// CF API directly with the user's key, detects project-level opt-outs
// before writing anything, and enforces the "manifest only written if
// fetch + hash succeed" contract.

// Remove deletes a mod's manifest and re-indexes.
func (r *Runner) Remove(slug string) (string, error) {
	return r.Run("remove", slug, "--yes")
}

// Update bumps a single mod to its latest version. Skips pinned mods.
func (r *Runner) Update(slug string) (string, error) {
	return r.Run("update", slug, "--yes")
}

// UpdateAll bumps all mods to their latest versions, skipping pinned ones.
// Self-hosted mods (no [update.*] block) are skipped by packwiz with a
// "supported update system cannot be found" message.
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
