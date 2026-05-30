package handlers

import (
	"net/http"
	"path"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/derpack/derpack-edit/internal/curseforge"
	"github.com/derpack/derpack-edit/internal/hashutil"
	"github.com/derpack/derpack-edit/internal/modrinth"
	"github.com/derpack/derpack-edit/internal/packwiz"
)

// cfOptOutMessage is the user-facing explanation when a CurseForge author
// has disabled third-party distribution. Used in several spots that all
// surface the same condition, so the wording lives in one place.
const cfOptOutMessage = "This mod's author has disabled third-party distribution on CurseForge. " +
	"packwiz cannot auto-fetch the jar. To use this mod, you'll need to self-host the file " +
	"(see docs for the self-hosted workflow)."

// cfPageURLPattern matches the CurseForge page URLs users typically copy
// from their browser. Captures the file ID.
//
// Matches:
//
//	https://www.curseforge.com/minecraft/mc-mods/<slug>/files/<fileId>
//	https://www.curseforge.com/minecraft/mc-mods/<slug>/download/<fileId>
//	https://legacy.curseforge.com/minecraft/mc-mods/<slug>/files/<fileId>
//	https://www.curseforge.com/projects/<projectId>/files/<fileId>   (older form)
//
// Does NOT match direct CDN URLs (media.forgecdn.net, edge.forgecdn.net) —
// those are already direct downloads and don't need resolution.
var cfPageURLPattern = regexp.MustCompile(
	`curseforge\.com/(?:minecraft/mc-mods/[^/]+|projects/\d+)/(?:files|download)/(\d+)`)

// extractCFFileID returns the file ID from a CurseForge page URL, or 0 if
// the URL doesn't look like one we recognize.
func extractCFFileID(url string) int64 {
	m := cfPageURLPattern.FindStringSubmatch(url)
	if m == nil {
		return 0
	}
	id, err := strconv.ParseInt(m[1], 10, 64)
	if err != nil {
		return 0
	}
	return id
}

// ----- GET /api/mods/{slug}/versions ------------------------------------

type versionEntry struct {
	ID            string   `json:"id"`             // version ID (Modrinth) or file ID (CF)
	VersionNumber string   `json:"version_number"` // e.g. "2.29.33"
	Name          string   `json:"name"`
	Type          string   `json:"type"` // "release", "beta", "alpha"
	GameVersions  []string `json:"game_versions"`
	Loaders       []string `json:"loaders"`
	DatePublished string   `json:"date_published"`
	Downloads     int      `json:"downloads"`
	URL           string   `json:"url"`
	Filename      string   `json:"filename"`
	Size          int64    `json:"size"`
}

// HandleListVersions resolves the slug to a manifest, looks up the source
// (Modrinth/CF), then queries that source's API for available versions.
func (s *Server) HandleListVersions(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	if slug == "" {
		writeError(w, http.StatusBadRequest, "missing slug in path")
		return
	}

	mod, err := loadModBySlug(s.RepoRoot, slug)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	switch mod.Source() {
	case "mr":
		s.listModrinthVersions(w, mod)
	case "cf":
		s.listCurseForgeVersions(w, mod)
	default:
		writeError(w, http.StatusBadRequest,
			"mod has no recognized update source (must be Modrinth or CurseForge)")
	}
}

func (s *Server) listModrinthVersions(w http.ResponseWriter, mod *packwiz.Mod) {
	if mod.Update.Modrinth == nil {
		writeError(w, http.StatusBadRequest, "mod is not a Modrinth mod")
		return
	}

	// Use the loaders/game versions from pack.toml to filter results.
	pack, err := packwiz.LoadPack(s.RepoRoot)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	loaders, mcVersions := loadersAndVersions(pack)

	versions, err := modrinth.ListVersions(mod.Update.Modrinth.ModID, loaders, mcVersions)
	if err != nil {
		writeError(w, http.StatusBadGateway, err.Error())
		return
	}

	out := make([]versionEntry, 0, len(versions))
	for _, v := range versions {
		f := v.PrimaryFile()
		if f == nil {
			continue
		}
		out = append(out, versionEntry{
			ID:            v.ID,
			VersionNumber: v.VersionNumber,
			Name:          v.Name,
			Type:          v.VersionType,
			GameVersions:  v.GameVersions,
			Loaders:       v.Loaders,
			DatePublished: v.DatePublished,
			Downloads:     v.Downloads,
			URL:           f.URL,
			Filename:      f.Filename,
			Size:          f.Size,
		})
	}
	// API returns newest first already; preserve that order. Keep sort by
	// date_published descending as a safety net in case the API changes.
	sort.SliceStable(out, func(i, j int) bool {
		return out[i].DatePublished > out[j].DatePublished
	})

	writeJSON(w, http.StatusOK, out)
}

// listCurseForgeVersions handles the CF source. Requires a CF API key
// configured in user settings; returns a clear error if missing.
//
// If every file from the API has an empty downloadUrl (the author opted out
// of third-party distribution at the project level), returns a structured
// {opted_out: true, hint: ...} response so the UI can render a real
// explanation instead of a misleading "no matching versions" message.
func (s *Server) listCurseForgeVersions(w http.ResponseWriter, mod *packwiz.Mod) {
	if mod.Update.CurseForge == nil {
		writeError(w, http.StatusBadRequest, "mod is not a CurseForge mod")
		return
	}
	if s.Config.CurseForgeAPIKey == "" {
		// 200 with a structured error so the UI can render a "configure key"
		// prompt instead of a generic failure. The status field is "needs_key".
		writeJSON(w, http.StatusOK, map[string]any{
			"needs_key": true,
			"hint":      "Add a CurseForge API key in Settings to enable version picking. Get a free key at https://console.curseforge.com/.",
		})
		return
	}

	pack, err := packwiz.LoadPack(s.RepoRoot)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
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

	cf := curseforge.New(s.Config.CurseForgeAPIKey)
	files, err := cf.ListFiles(int64(mod.Update.CurseForge.ProjectID), mcVersion, loaderName)
	if err != nil {
		writeError(w, http.StatusBadGateway, err.Error())
		return
	}

	out := make([]versionEntry, 0, len(files))
	optedOutCount := 0
	for _, f := range files {
		if f.DownloadURL == "" {
			// Author opted out of third-party distribution for this file.
			// Track it so we can surface a clear message if NOTHING resolves.
			optedOutCount++
			continue
		}
		// CF doesn't have a separate "version_number" field — DisplayName is
		// usually a clean version string ("1.21.1-NF-0.2.1") that's fine to
		// surface as both name and version.
		out = append(out, versionEntry{
			ID:            strconv.FormatInt(f.ID, 10),
			VersionNumber: f.DisplayName,
			Name:          f.DisplayName,
			Type:          curseforge.ReleaseTypeName(f.ReleaseType),
			GameVersions:  f.GameVersions,
			Loaders:       []string{loaderName},
			DatePublished: f.FileDate,
			URL:           f.DownloadURL,
			Filename:      f.FileName,
			Size:          f.FileLength,
		})
	}

	// If we found files for this MC/loader combo but every single one is
	// opted-out, the user needs to know that — otherwise the picker says
	// "no matching versions" and they assume it's a filter problem.
	if len(out) == 0 && optedOutCount > 0 {
		writeJSON(w, http.StatusOK, map[string]any{
			"opted_out": true,
			"hint":      cfOptOutMessage,
		})
		return
	}

	sort.SliceStable(out, func(i, j int) bool {
		return out[i].DatePublished > out[j].DatePublished
	})

	writeJSON(w, http.StatusOK, out)
}

// loadersAndVersions converts the [versions] table from pack.toml into
// loader names and MC versions.
//
//	[versions]
//	minecraft = "1.21.1"
//	neoforge = "21.1.180"
//
// → loaders=["neoforge"], mcVersions=["1.21.1"].
func loadersAndVersions(pack *packwiz.PackToml) (loaders, mcVersions []string) {
	knownLoaders := map[string]bool{
		"forge": true, "neoforge": true, "fabric": true, "quilt": true,
	}
	for k, v := range pack.Versions {
		if k == "minecraft" {
			mcVersions = append(mcVersions, v)
		} else if knownLoaders[k] {
			loaders = append(loaders, k)
		}
	}
	return loaders, mcVersions
}

// ----- POST /api/mods/set-version ---------------------------------------

type setVersionReq struct {
	Slug string `json:"slug"`

	// Path 1: from Modrinth/CurseForge picker.
	VersionID string `json:"version_id,omitempty"`

	// Path 2: paste URL (works for any source).
	URL      string `json:"url,omitempty"`
	Filename string `json:"filename,omitempty"` // optional, derived from URL if absent
}

// HandleSetVersion changes a mod to a specific version, either via picker
// (version_id) or by pasting a URL. The contract is: the on-disk manifest
// is only mutated if both the source resolution AND the hash computation
// succeed. A failure at any point leaves the existing manifest untouched.
func (s *Server) HandleSetVersion(w http.ResponseWriter, r *http.Request) {
	if !requirePost(w, r) {
		return
	}
	var req setVersionReq
	if !decodeJSON(w, r, &req) {
		return
	}
	req.Slug = strings.TrimSpace(req.Slug)
	if req.Slug == "" {
		writeError(w, http.StatusBadRequest, "slug is required")
		return
	}

	mod, err := loadModBySlug(s.RepoRoot, req.Slug)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	// Resolve URL + filename + (optionally) version ID / file ID. Nothing
	// hits disk in this section.
	var (
		url, filename, versionID string
		cfFileID                 int64
	)
	switch {
	case req.VersionID != "":
		// Picker path. The version_id is a Modrinth version ID for MR mods,
		// or a CurseForge numeric file ID for CF mods.
		switch mod.Source() {
		case "mr":
			versions, err := modrinth.ListVersions(mod.Update.Modrinth.ModID, nil, nil)
			if err != nil {
				writeError(w, http.StatusBadGateway, err.Error())
				return
			}
			var picked *modrinth.Version
			for i := range versions {
				if versions[i].ID == req.VersionID {
					picked = &versions[i]
					break
				}
			}
			if picked == nil {
				writeError(w, http.StatusNotFound, "version_id not found for this mod")
				return
			}
			f := picked.PrimaryFile()
			if f == nil {
				writeError(w, http.StatusInternalServerError, "version has no files")
				return
			}
			url = f.URL
			filename = f.Filename
			versionID = picked.ID

		case "cf":
			if s.Config.CurseForgeAPIKey == "" {
				writeError(w, http.StatusBadRequest,
					"CurseForge API key required to set version by ID. Add one in Settings, or use the Paste URL tab.")
				return
			}
			fileID, parseErr := strconv.ParseInt(req.VersionID, 10, 64)
			if parseErr != nil {
				writeError(w, http.StatusBadRequest, "CurseForge version_id must be a numeric file ID")
				return
			}
			// One CF API call instead of listing + filtering. Returns the
			// File directly by ID.
			cf := curseforge.New(s.Config.CurseForgeAPIKey)
			picked, err := cf.GetFile(fileID)
			if err != nil {
				writeError(w, http.StatusBadGateway, err.Error())
				return
			}
			if picked.DownloadURL == "" {
				writeError(w, http.StatusBadRequest, cfOptOutMessage)
				return
			}
			url = picked.DownloadURL
			filename = picked.FileName
			cfFileID = picked.ID

		default:
			writeError(w, http.StatusBadRequest, "version_id only works for Modrinth or CurseForge mods")
			return
		}

	case req.URL != "":
		pasted := strings.TrimSpace(req.URL)

		// If this looks like a CurseForge page URL, resolve via the CF API
		// to get the actual download URL. Direct GETs to CF page URLs return
		// 403 — that's what was producing the "saved manifest but failed to
		// compute hash" error before this fix.
		if fileID := extractCFFileID(pasted); fileID != 0 {
			if s.Config.CurseForgeAPIKey == "" {
				writeError(w, http.StatusBadRequest,
					"This looks like a CurseForge page URL. Add a CurseForge API key in Settings to resolve it, or paste the direct media.forgecdn.net URL instead.")
				return
			}
			cf := curseforge.New(s.Config.CurseForgeAPIKey)
			picked, ferr := cf.GetFile(fileID)
			if ferr != nil {
				writeError(w, http.StatusBadGateway, ferr.Error())
				return
			}
			if picked.DownloadURL == "" {
				writeError(w, http.StatusBadRequest, cfOptOutMessage)
				return
			}
			url = picked.DownloadURL
			filename = picked.FileName
			cfFileID = picked.ID
			// Honor a user-supplied filename override.
			if req.Filename != "" {
				filename = req.Filename
			}
		} else {
			// Not a CF page URL — proceed with the pasted URL as-is.
			url = pasted
			if req.Filename != "" {
				filename = req.Filename
			} else {
				filename = path.Base(url)
			}
			// If the URL is a Modrinth CDN URL, we can extract the version ID
			// for the [update.modrinth] block.
			if vid := extractModrinthVersionID(url); vid != "" {
				versionID = vid
			}
		}

	default:
		writeError(w, http.StatusBadRequest, "must provide either version_id or url")
		return
	}

	// Stage manifest mutations in memory; nothing is written to disk yet.
	mod.Filename = filename
	mod.Download.URL = url
	mod.Download.HashFormat = "sha512"
	if mod.Update.Modrinth != nil && versionID != "" {
		mod.Update.Modrinth.Version = versionID
	}
	if mod.Update.CurseForge != nil && cfFileID != 0 {
		mod.Update.CurseForge.FileID = int(cfFileID)
	}

	// Compute the hash BEFORE writing the manifest. If the download fails,
	// nothing changes on disk — the prior manifest is preserved exactly.
	hash, size, err := hashutil.SHA512OfURL(url)
	if err != nil {
		writeJSON(w, http.StatusOK, opResp{
			OK:    false,
			Slug:  req.Slug,
			Error: "failed to download from URL: " + err.Error(),
		})
		return
	}
	mod.Download.Hash = hash

	// Now save — manifest hits disk only if the hash succeeded.
	if err := packwiz.SaveMod(s.RepoRoot, mod); err != nil {
		writeError(w, http.StatusInternalServerError, "save manifest: "+err.Error())
		return
	}

	if _, err := s.Runner.Refresh(); err != nil {
		writeJSON(w, http.StatusOK, opResp{
			OK:     false,
			Slug:   req.Slug,
			Output: msgWithSize(filename, size),
			Error:  "set version, but `packwiz refresh` failed: " + err.Error(),
		})
		return
	}

	writeJSON(w, http.StatusOK, opResp{
		OK:     true,
		Slug:   req.Slug,
		Output: msgWithSize(filename, size),
	})
}

func msgWithSize(filename string, size int64) string {
	return "set " + filename + " (" + humanBytes(size) + ")"
}

func humanBytes(b int64) string {
	const unit = 1024
	if b < unit {
		return strconv.FormatInt(b, 10) + " B"
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return strconv.FormatFloat(float64(b)/float64(div), 'f', 1, 64) + " " + string("KMGT"[exp]) + "B"
}

// extractModrinthVersionID pulls "abc123XY" from a Modrinth CDN URL like
// https://cdn.modrinth.com/data/PROJECTID/versions/VERSIONID/filename.jar.
func extractModrinthVersionID(url string) string {
	const marker = "/versions/"
	idx := strings.Index(url, marker)
	if idx < 0 {
		return ""
	}
	rest := url[idx+len(marker):]
	end := strings.Index(rest, "/")
	if end < 0 {
		return ""
	}
	return rest[:end]
}

// ----- POST /api/mods/compute-hash --------------------------------------

type computeHashReq struct {
	Slug string `json:"slug"`
}

func (s *Server) HandleComputeHash(w http.ResponseWriter, r *http.Request) {
	if !requirePost(w, r) {
		return
	}
	var req computeHashReq
	if !decodeJSON(w, r, &req) {
		return
	}
	req.Slug = strings.TrimSpace(req.Slug)
	if req.Slug == "" {
		writeError(w, http.StatusBadRequest, "slug is required")
		return
	}

	mod, err := loadModBySlug(s.RepoRoot, req.Slug)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	if mod.Download.URL == "" {
		writeError(w, http.StatusBadRequest, "mod has no download URL")
		return
	}

	hash, size, err := hashutil.SHA512OfURL(mod.Download.URL)
	if err != nil {
		writeJSON(w, http.StatusOK, opResp{
			OK:    false,
			Slug:  req.Slug,
			Error: err.Error(),
		})
		return
	}

	mod.Download.Hash = hash
	mod.Download.HashFormat = "sha512"
	if err := packwiz.SaveMod(s.RepoRoot, mod); err != nil {
		writeError(w, http.StatusInternalServerError, "save manifest: "+err.Error())
		return
	}
	if _, err := s.Runner.Refresh(); err != nil {
		writeJSON(w, http.StatusOK, opResp{
			OK:     false,
			Slug:   req.Slug,
			Output: msgWithSize(mod.Filename, size),
			Error:  "computed hash, but `packwiz refresh` failed: " + err.Error(),
		})
		return
	}

	writeJSON(w, http.StatusOK, opResp{
		OK:     true,
		Slug:   req.Slug,
		Output: msgWithSize(mod.Filename, size),
	})
}

// ----- POST /api/mods/update --------------------------------------------

type updateReq struct {
	Slug string `json:"slug"` // empty => update all
}

func (s *Server) HandleUpdate(w http.ResponseWriter, r *http.Request) {
	if !requirePost(w, r) {
		return
	}
	var req updateReq
	if !decodeJSON(w, r, &req) {
		return
	}
	req.Slug = strings.TrimSpace(req.Slug)

	var (
		out string
		err error
	)
	if req.Slug == "" {
		out, err = s.Runner.UpdateAll()
	} else {
		if !packwiz.ManifestExists(s.RepoRoot, req.Slug) {
			writeError(w, http.StatusNotFound, "no manifest for slug: "+req.Slug)
			return
		}
		out, err = s.Runner.Update(req.Slug)
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
	writeJSON(w, http.StatusOK, opResp{
		OK:     true,
		Slug:   req.Slug,
		Output: out,
	})
}

// ----- helpers ----------------------------------------------------------

// loadModBySlug finds a single mod by its slug. Returns an error if not found.
func loadModBySlug(repoRoot, slug string) (*packwiz.Mod, error) {
	if !packwiz.ManifestExists(repoRoot, slug) {
		return nil, &slugNotFoundError{slug: slug}
	}
	mods, err := packwiz.LoadMods(repoRoot)
	if err != nil {
		return nil, err
	}
	for _, m := range mods {
		if m.Slug == slug {
			return m, nil
		}
	}
	return nil, &slugNotFoundError{slug: slug}
}

type slugNotFoundError struct{ slug string }

func (e *slugNotFoundError) Error() string { return "mod not found: " + e.slug }
