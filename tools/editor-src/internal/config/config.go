package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Config holds user-specific app settings, persisted to .editor/config.json.
type Config struct {
	// PrismInstancePath is the directory of the Prism Launcher instance to
	// copy bundled mods into when "Launch in Prism" is clicked.
	// Example: C:\Users\Alex\AppData\Roaming\PrismLauncher\instances\Derpack
	PrismInstancePath string `json:"prism_instance_path,omitempty"`
}

// Load reads .editor/config.json from the repo root, returning a zero-value
// Config if the file doesn't exist.
func Load(repoRoot string) (*Config, error) {
	path := filepath.Join(repoRoot, ".editor", "config.json")
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &Config{}, nil
		}
		return nil, fmt.Errorf("read %s: %w", path, err)
	}

	var c Config
	if err := json.Unmarshal(data, &c); err != nil {
		return nil, fmt.Errorf("parse %s: %w", path, err)
	}
	return &c, nil
}

// Save writes the config to .editor/config.json.
func Save(repoRoot string, c *Config) error {
	dir := filepath.Join(repoRoot, ".editor")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("mkdir %s: %w", dir, err)
	}

	path := filepath.Join(dir, "config.json")
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal config: %w", err)
	}

	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("write %s: %w", path, err)
	}
	return nil
}
