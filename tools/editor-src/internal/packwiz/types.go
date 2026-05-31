package packwiz

import "fmt"

// PackToml is the top-level pack.toml manifest at the repo root.
type PackToml struct {
	Name       string                 `toml:"name"`
	Author     string                 `toml:"author,omitempty"`
	Version    string                 `toml:"version,omitempty"`
	PackFormat string                 `toml:"pack-format"`
	Index      IndexRef               `toml:"index"`
	Versions   map[string]string      `toml:"versions"`
	Export     map[string]interface{} `toml:"export,omitempty"`
}

type IndexRef struct {
	File       string `toml:"file"`
	HashFormat string `toml:"hash-format"`
	Hash       string `toml:"hash,omitempty"`
}

// Mod is a single mod's .pw.toml manifest. Mirrors packwiz's own struct so
// round-tripping preserves fields we don't explicitly touch.
type Mod struct {
	Name     string   `toml:"name"`
	Filename string   `toml:"filename"`
	Side     string   `toml:"side,omitempty"`
	Pin      bool     `toml:"pin,omitempty"`
	Download Download `toml:"download"`
	Update   Update   `toml:"update,omitempty"`

	// Slug is the manifest filename without .pw.toml — not in the file itself.
	// Populated by the loader.
	Slug string `toml:"-"`
}

type Download struct {
	URL        string `toml:"url"`
	HashFormat string `toml:"hash-format"`
	Hash       string `toml:"hash"`
}

// Update holds source-specific update metadata. At most one of Modrinth or
// CurseForge will be populated for a given mod. A mod with neither is a
// self-hosted entry (packwiz skips it during `update --all`).
type Update struct {
	Modrinth   *ModrinthUpdate   `toml:"modrinth,omitempty"`
	CurseForge *CurseForgeUpdate `toml:"curseforge,omitempty"`
}

type ModrinthUpdate struct {
	ModID   string `toml:"mod-id"`
	Version string `toml:"version"`
}

type CurseForgeUpdate struct {
	FileID    int `toml:"file-id"`
	ProjectID int `toml:"project-id"`
}

// Source returns "mr", "cf", or "" depending on which update block is
// present. "" indicates a self-hosted mod with no upstream.
func (m *Mod) Source() string {
	if m.Update.Modrinth != nil {
		return "mr"
	}
	if m.Update.CurseForge != nil {
		return "cf"
	}
	return ""
}

// PageURL returns the human-facing project page for the mod on its source
// (Modrinth or CurseForge), or "" if the mod has no recognized update source
// (e.g. self-hosted URL-only manifests).
//
// Modrinth: the manifest filename matches the project's canonical slug in
// almost all cases. Modrinth keeps redirects for renamed projects, so this
// stays valid even after a slug change. If a slug ever does 404, switching
// to m.Update.Modrinth.ModID would be a one-line fix.
//
// CurseForge: the manifest only stores the numeric project ID. CurseForge's
// "/projects/<id>" URL redirects to the canonical slug page, so we let
// CurseForge resolve it rather than calling the API to look up the slug.
func (m *Mod) PageURL() string {
	switch {
	case m.Update.Modrinth != nil && m.Slug != "":
		return "https://modrinth.com/mod/" + m.Slug
	case m.Update.CurseForge != nil && m.Update.CurseForge.ProjectID != 0:
		return fmt.Sprintf("https://www.curseforge.com/projects/%d", m.Update.CurseForge.ProjectID)
	}
	return ""
}
