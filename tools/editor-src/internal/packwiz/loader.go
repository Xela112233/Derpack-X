package packwiz

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/BurntSushi/toml"
)

// LoadPack reads pack.toml from the given repo root.
func LoadPack(repoRoot string) (*PackToml, error) {
	path := filepath.Join(repoRoot, "pack.toml")
	var p PackToml
	if _, err := toml.DecodeFile(path, &p); err != nil {
		return nil, fmt.Errorf("read %s: %w", path, err)
	}
	return &p, nil
}

// LoadMods scans <repoRoot>/mods/*.pw.toml and returns each parsed manifest.
// The Mod.Slug field is populated from the filename (sans .pw.toml).
// Returns an empty slice if mods/ doesn't exist.
func LoadMods(repoRoot string) ([]*Mod, error) {
	dir := filepath.Join(repoRoot, "mods")
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return []*Mod{}, nil
		}
		return nil, fmt.Errorf("read %s: %w", dir, err)
	}

	mods := make([]*Mod, 0, len(entries))
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if !strings.HasSuffix(name, ".pw.toml") {
			continue
		}
		slug := strings.TrimSuffix(name, ".pw.toml")
		path := filepath.Join(dir, name)

		m, err := loadMod(path)
		if err != nil {
			// Skip broken manifests but log; the UI can show them as errored.
			// For v1 we just skip silently and a future endpoint can list them.
			fmt.Fprintf(os.Stderr, "warning: failed to parse %s: %v\n", path, err)
			continue
		}
		m.Slug = slug
		mods = append(mods, m)
	}

	sort.Slice(mods, func(i, j int) bool {
		return mods[i].Slug < mods[j].Slug
	})
	return mods, nil
}

func loadMod(path string) (*Mod, error) {
	var m Mod
	if _, err := toml.DecodeFile(path, &m); err != nil {
		return nil, err
	}
	return &m, nil
}

// SaveMod writes a Mod back to its manifest file. The file path is
// <repoRoot>/mods/<slug>.pw.toml.
func SaveMod(repoRoot string, m *Mod) error {
	if m.Slug == "" {
		return fmt.Errorf("cannot save mod with empty slug")
	}
	path := filepath.Join(repoRoot, "mods", m.Slug+".pw.toml")
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create %s: %w", path, err)
	}
	defer f.Close()

	enc := toml.NewEncoder(f)
	enc.Indent = ""
	if err := enc.Encode(m); err != nil {
		return fmt.Errorf("encode %s: %w", path, err)
	}
	return nil
}

// ManifestPath returns the path to a mod's manifest within the repo.
func ManifestPath(repoRoot, slug string) string {
	return filepath.Join(repoRoot, "mods", slug+".pw.toml")
}

// ManifestExists checks if mods/<slug>.pw.toml exists.
func ManifestExists(repoRoot, slug string) bool {
	_, err := os.Stat(ManifestPath(repoRoot, slug))
	return err == nil
}
