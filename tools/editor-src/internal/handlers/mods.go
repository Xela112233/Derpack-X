package handlers

import (
	"fmt"
	"net/http"
	"path"
	"strings"

	"github.com/derpack/derpack-edit/internal/curseforge"
	"github.com/derpack/derpack-edit/internal/hashutil"
	"github.com/derpack/derpack-edit/internal/packwiz"
)

type addReq struct {
	Source string `json:"source"` // "mr", "cf", or "url" (self-hosted)
	Slug   string `json:"slug"`
	Side   string `json:"side"` // "both", "client", "server", or ""

	// Self-hosted fields. Only used when Source == "url".
	Name     string `json:"name,omitempty"`
	URL      string `json:"url,omitempty"`
	Filename string `json:"filename,omitempty"`
}

type opResp struct {
	OK     bool   `json:"ok"`
	Slug   string `json:"slug,omitempty"`
	Output string `json:"output,omitempty"`
	Error  string `json:"error,omitempty"`
}

func (s *Server) HandleAddMod(w http.ResponseWriter, r *http.Request) {
	if !requirePost(w, r) {
		return
	}
	var req addReq
	if !decodeJSON(w, r, &req) {
		return
	}
	req.Slug = strings.TrimSpace(req.Slug)
	if req.Slug == "" {
		writeError(w, http.StatusBadRequest, "slug is required")
		return
	}

	var (
		out string
		err error
	)
	switch req.Source {
	case "mr":
		out, err = s.Runner.AddModrinth(req.Slug, req.Side)
	case "cf":
		// CurseForge adds go through our own resolution path (not
		// `packwiz cf add`) so we can pass the user's API key, detect
		// opt-outs before writing the manifest, and apply the "manifest
		// only written if fetch + hash succeed" contract.
		out, err = s.addCurseForgeMod(req.Slug)
	case "url":
		out, err = s.addSelfHostedMod(req.Slug, req.Name, req.URL, req.Filename, req.Side)
	default:
		writeError(w, http.StatusBadRequest, "source must be 'mr', 'cf', or 'url'")
		return
	}

	if err != nil {
		writeJSON(w, http.StatusOK, opResp{
			OK:     false,
			Slug:   req.Slug,
			Output: out,
			Error:  err.Error(),
		})
		return
	}

	// For mr/cf, apply the user's side choice (packwiz may pick wrong side
	// from upstream metadata). For self-hosted, side is already set during
	// manifest creation.
	if req.Side != "" && req.Source != "url" {
		if err := s.applySide(req.Slug, req.Side); err != nil {
			// Mod was added but side wasn't applied. Surface as a warning.
			writeJSON(w, http.StatusOK, opResp{
				OK:     true,
				Slug:   req.Slug,
				Output: out + "\nNOTE: side could not be set: " + err.Error(),
			})
			return
		}
	}

	writeJSON(w, http.StatusOK, opResp{OK: true, Slug: req.Slug, Output: out})
}

// addCurseForgeMod is the replacement for `packwiz cf add`. The bare packwiz
// invocation has no way to pass our stored CF API key, and when packwiz
// can't authenticate it falls back to writing manifests with empty download
// URLs that crash packwiz-installer at launch time. This path:
//
//  1. Looks up the project by slug via the CF API
//  2. Checks the project-level AllowModDistribution flag before doing
//     anything else (so opt-out gets a clear up-front error)
//  3. Picks the newest release matching the pack's MC version + loader
//  4. Computes the SHA-512 hash by downloading the jar
//  5. Writes the manifest and refreshes the index
//
// Manifest only hits disk if every step succeeds.
func (s *Server) addCurseForgeMod(slug string) (string, error) {
	if s.Config.CurseForgeAPIKey == "" {
		return "", fmt.Errorf("CurseForge API key required to add CF mods. Configure one in Settings (free key from https://console.curseforge.com/)")
	}

	cf := curseforge.New(s.Config.CurseForgeAPIKey)

	// 1. Slug → project ID + opt-out flag.
	proj, err := cf.GetModBySlug(slug)
	if err != nil {
		return "", fmt.Errorf("look up CurseForge mod %q: %w", slug, err)
	}
	if proj == nil {
		return "", fmt.Errorf("CurseForge mod %q not found (check the slug)", slug)
	}
	if !proj.AllowModDistribution {
		return "", fmt.Errorf("%s: %s", proj.Name, cfOptOutMessage)
	}

	// 2. Look at the pack to determine which MC version + loader to match.
	pack, err := packwiz.LoadPack(s.RepoRoot)
	if err != nil {
		return "", fmt.Errorf("load pack.toml: %w", err)
	}
	loaders, mcVersions := loadersAndVersions(pack)
	mcVersion := ""
	if len(mcVersions) > 0 {
		mcVersion = mcVersions[0]
	}
	loaderName := ""
	if len(loaders) > 0 {
		loaderName = loaders[0]
	}

	// 3. List files filtered to MC + loader, then pick the newest release.
	files, err := cf.ListFiles(proj.ID, mcVersion, loaderName)
	if err != nil {
		return "", fmt.Errorf("list CurseForge files for %q: %w", slug, err)
	}
	if len(files) == 0 {
		return "", fmt.Errorf("no CurseForge files for %s match Minecraft %s + %s. Check the mod page for available versions",
			proj.Name, mcVersion, loaderName)
	}

	// CF returns files newest-first. Prefer a release-channel file; fall
	// back to beta/alpha if no release is available for this combo. Skip
	// files with empty downloadUrl (per-file opt-out).
	var picked *curseforge.File
	for i := range files {
		if files[i].DownloadURL == "" {
			continue
		}
		if files[i].ReleaseType == 1 {
			picked = &files[i]
			break
		}
	}
	if picked == nil {
		for i := range files {
			if files[i].DownloadURL != "" {
				picked = &files[i]
				break
			}
		}
	}
	if picked == nil {
		// All files for this combo were opted-out. Treat as project-level
		// opt-out for the user — same fix needed either way.
		return "", fmt.Errorf("%s: %s", proj.Name, cfOptOutMessage)
	}

	// 4. Hash the file BEFORE writing the manifest. If this fails (network,
	// 403, etc.), nothing changes on disk.
	hash, size, err := hashutil.SHA512OfURL(picked.DownloadURL)
	if err != nil {
		return "", fmt.Errorf("download %s: %w", picked.FileName, err)
	}

	// 5. Build the manifest and write it.
	mod := &packwiz.Mod{
		Name:     proj.Name,
		Filename: picked.FileName,
		Slug:     slug,
		// Side is intentionally left unset here. HandleAddMod's applySide
		// call writes the user's choice afterward; leaving this empty
		// avoids a redundant write of the wrong value.
		Download: packwiz.Download{
			URL:        picked.DownloadURL,
			HashFormat: "sha512",
			Hash:       hash,
		},
		Update: packwiz.Update{
			CurseForge: &packwiz.CurseForgeUpdate{
				ProjectID: int(proj.ID),
				FileID:    int(picked.ID),
			},
		},
	}
	if err := packwiz.SaveMod(s.RepoRoot, mod); err != nil {
		return "", fmt.Errorf("save manifest: %w", err)
	}
	if _, err := s.Runner.Refresh(); err != nil {
		return "", fmt.Errorf("manifest saved but `packwiz refresh` failed: %w", err)
	}

	return fmt.Sprintf("Added %s — CF file %d, %s", proj.Name, picked.ID, humanBytes(size)), nil
}

// addSelfHostedMod adds a mod backed by a direct URL (no Modrinth / CF
// upstream). The user supplies a slug, name, URL, side, and optionally a
// filename. The manifest has no [update.*] block, so `packwiz update` and
// `packwiz update --all` correctly skip it (packwiz prints "A supported
// update system for ... cannot be found" and moves on).
//
// As with addCurseForgeMod, the manifest only hits disk if the hash succeeds.
func (s *Server) addSelfHostedMod(slug, name, urlStr, filename, side string) (string, error) {
	slug = strings.TrimSpace(slug)
	name = strings.TrimSpace(name)
	urlStr = strings.TrimSpace(urlStr)
	filename = strings.TrimSpace(filename)

	if name == "" {
		return "", fmt.Errorf("name is required for self-hosted mods")
	}
	if urlStr == "" {
		return "", fmt.Errorf("URL is required for self-hosted mods")
	}
	if filename == "" {
		filename = path.Base(urlStr)
	}
	if filename == "" || filename == "." || filename == "/" {
		return "", fmt.Errorf("could not derive a filename from URL; provide one explicitly")
	}

	if packwiz.ManifestExists(s.RepoRoot, slug) {
		return "", fmt.Errorf("a mod with slug %q already exists; remove it first or use a different slug", slug)
	}

	hash, size, err := hashutil.SHA512OfURL(urlStr)
	if err != nil {
		return "", fmt.Errorf("download %s: %w", urlStr, err)
	}

	mod := &packwiz.Mod{
		Name:     name,
		Filename: filename,
		Slug:     slug,
		Download: packwiz.Download{
			URL:        urlStr,
			HashFormat: "sha512",
			Hash:       hash,
		},
		// No Update block — packwiz will treat this as having no update system.
	}
	// "both" is packwiz's default and is represented as an empty side field.
	if side != "" && side != "both" {
		mod.Side = side
	}

	if err := packwiz.SaveMod(s.RepoRoot, mod); err != nil {
		return "", fmt.Errorf("save manifest: %w", err)
	}
	if _, err := s.Runner.Refresh(); err != nil {
		return "", fmt.Errorf("manifest saved but `packwiz refresh` failed: %w", err)
	}

	return fmt.Sprintf("Added %s (self-hosted, %s)", name, humanBytes(size)), nil
}

// applySide rewrites mods/<slug>.pw.toml to set the Side field, then refreshes.
// Passing side = "both" clears the field entirely (packwiz default = both).
func (s *Server) applySide(slug, side string) error {
	mods, err := packwiz.LoadMods(s.RepoRoot)
	if err != nil {
		return err
	}
	for _, m := range mods {
		if m.Slug == slug {
			if side == "both" {
				m.Side = ""
			} else {
				m.Side = side
			}
			if err := packwiz.SaveMod(s.RepoRoot, m); err != nil {
				return err
			}
			_, err := s.Runner.Refresh()
			return err
		}
	}
	return nil
}

// HandleSetSide changes a mod's side field on an existing manifest.
// Useful for fixing mods that packwiz tagged wrong from Modrinth metadata.
func (s *Server) HandleSetSide(w http.ResponseWriter, r *http.Request) {
	if !requirePost(w, r) {
		return
	}
	var req struct {
		Slug string `json:"slug"`
		Side string `json:"side"` // "both", "client", "server"
	}
	if !decodeJSON(w, r, &req) {
		return
	}
	req.Slug = strings.TrimSpace(req.Slug)
	if req.Slug == "" {
		writeError(w, http.StatusBadRequest, "slug is required")
		return
	}
	if req.Side != "both" && req.Side != "client" && req.Side != "server" {
		writeError(w, http.StatusBadRequest, "side must be 'both', 'client', or 'server'")
		return
	}
	if !packwiz.ManifestExists(s.RepoRoot, req.Slug) {
		writeError(w, http.StatusNotFound, "no manifest for slug: "+req.Slug)
		return
	}
	if err := s.applySide(req.Slug, req.Side); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, opResp{OK: true, Slug: req.Slug, Output: "side set to " + req.Side})
}

type removeReq struct {
	Slug string `json:"slug"`
}

func (s *Server) HandleRemoveMod(w http.ResponseWriter, r *http.Request) {
	if !requirePost(w, r) {
		return
	}
	var req removeReq
	if !decodeJSON(w, r, &req) {
		return
	}
	req.Slug = strings.TrimSpace(req.Slug)
	if req.Slug == "" {
		writeError(w, http.StatusBadRequest, "slug is required")
		return
	}

	if !packwiz.ManifestExists(s.RepoRoot, req.Slug) {
		writeJSON(w, http.StatusOK, opResp{
			OK:    false,
			Slug:  req.Slug,
			Error: "no manifest exists for slug: " + req.Slug,
		})
		return
	}

	out, err := s.Runner.Remove(req.Slug)
	if err != nil {
		writeJSON(w, http.StatusOK, opResp{
			OK:     false,
			Slug:   req.Slug,
			Output: out,
			Error:  err.Error(),
		})
		return
	}
	writeJSON(w, http.StatusOK, opResp{OK: true, Slug: req.Slug, Output: out})
}

// HandleRefreshIndex runs `packwiz refresh` to rebuild index.toml. Useful
// after hand-edits to manifests, and required after the self-hosted add
// flow (which is handled internally by addSelfHostedMod, but external
// users may also want a manual button).
func (s *Server) HandleRefreshIndex(w http.ResponseWriter, r *http.Request) {
	if !requirePost(w, r) {
		return
	}
	out, err := s.Runner.Refresh()
	if err != nil {
		writeJSON(w, http.StatusOK, opResp{
			OK:     false,
			Output: out,
			Error:  err.Error(),
		})
		return
	}
	writeJSON(w, http.StatusOK, opResp{OK: true, Output: out})
}
