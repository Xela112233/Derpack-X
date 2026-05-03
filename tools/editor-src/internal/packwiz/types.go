package packwiz

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
// CurseForge will be populated for a given mod.
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

// Source returns "mr", "cf", or "" depending on which update block is present.
func (m *Mod) Source() string {
	if m.Update.Modrinth != nil {
		return "mr"
	}
	if m.Update.CurseForge != nil {
		return "cf"
	}
	return ""
}
