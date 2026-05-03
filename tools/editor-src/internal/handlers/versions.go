package handlers

import (
	"net/http"
	"path"
	"sort"
	"strconv"
	"strings"

	"github.com/derpack/derpack-edit/internal/hashutil"
	"github.com/derpack/derpack-edit/internal/modrinth"
	"github.com/derpack/derpack-edit/internal/packwiz"
)

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
		writeError(w, http.StatusNotImplemented,
			"CurseForge version listing is not supported in this version. Use the 'Paste URL' tab to set a specific CF file.")
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

// loadersAndVersions extracts loader names and MC versions from pack.toml.
// The "versions" map looks like:
//   [versions]
//   minecraft = "1.21.1"
//   neoforge = "21.1.180"
//
// We turn that into loaders=["neoforge"], mcVersions=["1.21.1"].
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

	// Path 1: from Modrinth picker.
	VersionID string `json:"version_id,omitempty"`

	// Path 2: paste URL (works for any source).
	URL      string `json:"url,omitempty"`
	Filename string `json:"filename,omitempty"` // optional, derived from URL if absent
}

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

	// Resolve URL + filename + (optionally) version ID.
	var (
		url, filename, versionID string
	)
	switch {
	case req.VersionID != "":
		// Modrinth picker path. We already have the mod ID; look up the version.
		if mod.Update.Modrinth == nil {
			writeError(w, http.StatusBadRequest, "version_id only works for Modrinth mods")
			return
		}
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

	case req.URL != "":
		url = strings.TrimSpace(req.URL)
		if req.Filename != "" {
			filename = req.Filename
		} else {
			filename = path.Base(url)
		}
		// We can't infer the version ID from a pasted URL, but if the URL is a
		// Modrinth CDN URL (https://cdn.modrinth.com/data/<modID>/versions/<versionID>/<file>),
		// we can extract it.
		if vid := extractModrinthVersionID(url); vid != "" {
			versionID = vid
		}

	default:
		writeError(w, http.StatusBadRequest, "must provide either version_id or url")
		return
	}

	// Mutate the manifest.
	mod.Filename = filename
	mod.Download.URL = url
	mod.Download.HashFormat = "sha512"
	mod.Download.Hash = "" // placeholder, computed below
	if mod.Update.Modrinth != nil && versionID != "" {
		mod.Update.Modrinth.Version = versionID
	}

	if err := packwiz.SaveMod(s.RepoRoot, mod); err != nil {
		writeError(w, http.StatusInternalServerError, "save manifest: "+err.Error())
		return
	}

	// Compute hash and save again.
	hash, size, err := hashutil.SHA512OfURL(url)
	if err != nil {
		writeJSON(w, http.StatusOK, opResp{
			OK:    false,
			Slug:  req.Slug,
			Error: "saved manifest but failed to compute hash: " + err.Error(),
		})
		return
	}
	mod.Download.Hash = hash
	if err := packwiz.SaveMod(s.RepoRoot, mod); err != nil {
		writeError(w, http.StatusInternalServerError, "save manifest with hash: "+err.Error())
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
