package modrinth

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	baseURL   = "https://api.modrinth.com/v2"
	userAgent = "derpack-edit/1.0 (modpack editor)"
)

// Version is a single version of a Modrinth project, simplified to the fields
// our UI needs.
type Version struct {
	ID            string   `json:"id"`             // version ID, e.g. "abc123XY"
	VersionNumber string   `json:"version_number"` // e.g. "2.29.33"
	Name          string   `json:"name"`           // human title
	VersionType   string   `json:"version_type"`   // "release", "beta", "alpha"
	GameVersions  []string `json:"game_versions"`  // e.g. ["1.21", "1.21.1"]
	Loaders       []string `json:"loaders"`        // e.g. ["neoforge"]
	DatePublished string   `json:"date_published"` // ISO 8601
	Downloads     int      `json:"downloads"`
	Files         []File   `json:"files"`
}

type File struct {
	URL      string            `json:"url"`
	Filename string            `json:"filename"`
	Size     int64             `json:"size"`
	Hashes   map[string]string `json:"hashes"`
	Primary  bool              `json:"primary"`
}

// PrimaryFile returns the primary jar file for this version, falling back to
// the first one if no file is marked primary.
func (v *Version) PrimaryFile() *File {
	for i := range v.Files {
		if v.Files[i].Primary {
			return &v.Files[i]
		}
	}
	if len(v.Files) > 0 {
		return &v.Files[0]
	}
	return nil
}

var httpClient = &http.Client{
	Timeout: 20 * time.Second,
}

// ListVersions fetches all versions of a Modrinth project, optionally filtered
// by loader and game version.
func ListVersions(modID string, loaders []string, gameVersions []string) ([]Version, error) {
	u := fmt.Sprintf("%s/project/%s/version", baseURL, url.PathEscape(modID))

	q := url.Values{}
	if len(loaders) > 0 {
		q.Set("loaders", jsonArray(loaders))
	}
	if len(gameVersions) > 0 {
		q.Set("game_versions", jsonArray(gameVersions))
	}
	if len(q) > 0 {
		u += "?" + q.Encode()
	}

	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch %s: %w", u, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("modrinth %s: %s: %s", u, resp.Status, string(body))
	}

	var versions []Version
	if err := json.NewDecoder(resp.Body).Decode(&versions); err != nil {
		return nil, fmt.Errorf("decode versions: %w", err)
	}
	return versions, nil
}

// jsonArray formats a string slice as a JSON array literal, which is what
// Modrinth's API expects for these query params.
func jsonArray(xs []string) string {
	quoted := make([]string, len(xs))
	for i, x := range xs {
		// Use json.Marshal for proper escaping
		b, _ := json.Marshal(x)
		quoted[i] = string(b)
	}
	return "[" + strings.Join(quoted, ",") + "]"
}
