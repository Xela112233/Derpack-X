package curseforge

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	baseURL   = "https://api.curseforge.com"
	userAgent = "derpack-edit/1.0 (modpack editor)"

	// Minecraft's CurseForge game ID, used in /v1/mods/search.
	minecraftGameID = 432
	// "Mods" class ID for Minecraft, used to narrow search to a unique result.
	// Per CF docs: "coupled with classId will result in a unique result".
	minecraftModsClassID = 6
)

// File is a single CurseForge file (their term for what Modrinth calls a
// version), simplified to the fields our UI needs.
type File struct {
	ID           int64    `json:"id"`           // file ID — used for [update.curseforge].file-id
	ModID        int64    `json:"modId"`        // parent project's mod ID
	DisplayName  string   `json:"displayName"`  // human title, e.g. "Create Aeronautics Compat 0.2.1"
	FileName     string   `json:"fileName"`     // jar name
	ReleaseType  int      `json:"releaseType"`  // 1=release, 2=beta, 3=alpha
	FileDate     string   `json:"fileDate"`     // ISO 8601
	FileLength   int64    `json:"fileLength"`   // bytes
	DownloadURL  string   `json:"downloadUrl"`  // direct CDN URL — empty when author opted out
	GameVersions []string `json:"gameVersions"` // includes loader entries like "NeoForge" alongside MC versions
}

// Mod represents a CurseForge project. We only need a handful of fields:
// the numeric ID for [update.curseforge].project-id, the human name for the
// manifest, and the AllowModDistribution flag to detect a project-level
// opt-out before we try to fetch any files.
type Mod struct {
	ID                   int64  `json:"id"`
	Name                 string `json:"name"`
	Slug                 string `json:"slug"`
	AllowModDistribution bool   `json:"allowModDistribution"`
}

// listFilesResp is the API envelope around the file list endpoints.
type listFilesResp struct {
	Data []File `json:"data"`
}

// searchModsResp wraps GET /v1/mods/search.
type searchModsResp struct {
	Data []Mod `json:"data"`
}

// getFilesReq is the body for POST /v1/mods/files.
type getFilesReq struct {
	FileIDs []int64 `json:"fileIds"`
}

// Client wraps the CurseForge REST API. Requires an API key from
// https://console.curseforge.com/.
type Client struct {
	APIKey string
	http   *http.Client
}

// New returns a client. An empty APIKey makes API calls fail with a clear
// error instead of an obscure 401.
func New(apiKey string) *Client {
	return &Client{
		APIKey: strings.TrimSpace(apiKey),
		http:   &http.Client{Timeout: 20 * time.Second},
	}
}

// authedRequest builds an HTTP request with the auth + accept headers the
// CurseForge API requires.
func (c *Client) authedRequest(method, fullURL string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, fullURL, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("x-api-key", c.APIKey)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	return req, nil
}

// checkStatus turns common error responses into descriptive Go errors so the
// UI can surface "check your API key" rather than a raw HTTP status.
func checkStatus(resp *http.Response, callURL string) error {
	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return fmt.Errorf("CurseForge API rejected request (HTTP %d). Check that your API key is valid", resp.StatusCode)
	}
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("curseforge %s: %s: %s", callURL, resp.Status, string(body))
	}
	return nil
}

// ListFiles returns all files for a CurseForge mod, optionally filtered to
// the given Minecraft game version. The loader is filtered client-side
// because CF's gameVersions field mixes loader names and MC versions in one
// list.
func (c *Client) ListFiles(modID int64, mcVersion, loaderName string) ([]File, error) {
	if c.APIKey == "" {
		return nil, fmt.Errorf("no CurseForge API key configured (add one in Settings)")
	}
	if modID == 0 {
		return nil, fmt.Errorf("modID is zero; this mod's manifest has no project-id")
	}

	apiURL := fmt.Sprintf("%s/v1/mods/%d/files?pageSize=50", baseURL, modID)
	if mcVersion != "" {
		// CF supports a gameVersion query param that does an exact match on
		// MC version. Loader is filtered client-side below.
		apiURL = fmt.Sprintf("%s&gameVersion=%s", apiURL, url.QueryEscape(mcVersion))
	}

	req, err := c.authedRequest(http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch %s: %w", apiURL, err)
	}
	defer resp.Body.Close()

	if err := checkStatus(resp, apiURL); err != nil {
		return nil, err
	}

	var envelope listFilesResp
	if err := json.NewDecoder(resp.Body).Decode(&envelope); err != nil {
		return nil, fmt.Errorf("decode CF response: %w", err)
	}

	if loaderName == "" {
		return envelope.Data, nil
	}

	// Filter to files that list this loader. CF's gameVersions field is a
	// flat string array mixing MC versions and loader names ("1.21.1",
	// "NeoForge", "Server", etc.) — case-insensitive match on the loader.
	want := strings.ToLower(loaderName)
	out := envelope.Data[:0]
	for _, f := range envelope.Data {
		for _, gv := range f.GameVersions {
			if strings.ToLower(gv) == want {
				out = append(out, f)
				break
			}
		}
	}
	return out, nil
}

// GetFile fetches a single CurseForge file by ID. Uses POST /v1/mods/files
// because it's the only endpoint that doesn't also require the parent mod ID
// — useful when we have a file ID from a pasted page URL but not the mod ID.
//
// Returns an error if the file ID isn't found.
func (c *Client) GetFile(fileID int64) (*File, error) {
	if c.APIKey == "" {
		return nil, fmt.Errorf("no CurseForge API key configured (add one in Settings)")
	}
	if fileID == 0 {
		return nil, fmt.Errorf("fileID is zero")
	}

	body, err := json.Marshal(getFilesReq{FileIDs: []int64{fileID}})
	if err != nil {
		return nil, err
	}
	apiURL := baseURL + "/v1/mods/files"
	req, err := c.authedRequest(http.MethodPost, apiURL, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch CF file %d: %w", fileID, err)
	}
	defer resp.Body.Close()

	if err := checkStatus(resp, apiURL); err != nil {
		return nil, err
	}

	var envelope listFilesResp
	if err := json.NewDecoder(resp.Body).Decode(&envelope); err != nil {
		return nil, fmt.Errorf("decode CF response: %w", err)
	}
	if len(envelope.Data) == 0 {
		return nil, fmt.Errorf("CurseForge file %d not found", fileID)
	}
	return &envelope.Data[0], nil
}

// GetModBySlug looks up a CurseForge Minecraft mod by its slug. Per the CF
// docs, search with gameId + classId + slug returns a unique result (or
// none). Returns nil with no error if no mod matches.
func (c *Client) GetModBySlug(slug string) (*Mod, error) {
	if c.APIKey == "" {
		return nil, fmt.Errorf("no CurseForge API key configured (add one in Settings)")
	}
	slug = strings.TrimSpace(slug)
	if slug == "" {
		return nil, fmt.Errorf("slug is empty")
	}

	apiURL := fmt.Sprintf("%s/v1/mods/search?gameId=%d&classId=%d&slug=%s",
		baseURL, minecraftGameID, minecraftModsClassID, url.QueryEscape(slug))

	req, err := c.authedRequest(http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("search CF for %q: %w", slug, err)
	}
	defer resp.Body.Close()

	if err := checkStatus(resp, apiURL); err != nil {
		return nil, err
	}

	var envelope searchModsResp
	if err := json.NewDecoder(resp.Body).Decode(&envelope); err != nil {
		return nil, fmt.Errorf("decode CF search response: %w", err)
	}

	// Prefer an exact slug match. CF's search is normally already exact when
	// classId is supplied, but be defensive in case the API loosens.
	for i := range envelope.Data {
		if strings.EqualFold(envelope.Data[i].Slug, slug) {
			return &envelope.Data[i], nil
		}
	}
	if len(envelope.Data) > 0 {
		return &envelope.Data[0], nil
	}
	return nil, nil
}

// ReleaseTypeName turns the integer into a human label.
func ReleaseTypeName(rt int) string {
	switch rt {
	case 1:
		return "release"
	case 2:
		return "beta"
	case 3:
		return "alpha"
	default:
		return "unknown"
	}
}
