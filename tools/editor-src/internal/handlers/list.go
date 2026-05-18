package handlers

import (
	"net/http"

	"github.com/derpack/derpack-edit/internal/packwiz"
)

// packResp is what GET /api/pack returns.
type packResp struct {
	Name       string `json:"name"`
	Version    string `json:"version"`
	Author     string `json:"author"`
	PackFormat string `json:"pack_format"`
}

// HandleGetPack returns metadata from pack.toml.
func (s *Server) HandleGetPack(w http.ResponseWriter, r *http.Request) {
	p, err := packwiz.LoadPack(s.RepoRoot)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, packResp{
		Name:       p.Name,
		Version:    p.Version,
		Author:     p.Author,
		PackFormat: p.PackFormat,
	})
}

// modResp is a single mod's data as the UI sees it.
type modResp struct {
	Slug    string `json:"slug"`
	Name    string `json:"name"`
	Version string `json:"version"` // for Modrinth: version ID; for CF: file ID; for URL: filename
	Side    string `json:"side"`
	Pinned  bool   `json:"pinned"`
	Source  string `json:"source"` // "mr", "cf", or ""
	URL     string `json:"url"`
	PageURL string `json:"page_url"` // Modrinth/CurseForge project page; "" if no source
}

// HandleListMods returns all mods from mods/*.pw.toml.
func (s *Server) HandleListMods(w http.ResponseWriter, r *http.Request) {
	mods, err := packwiz.LoadMods(s.RepoRoot)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	out := make([]modResp, 0, len(mods))
	for _, m := range mods {
		entry := modResp{
			Slug:    m.Slug,
			Name:    m.Name,
			Side:    m.Side,
			Pinned:  m.Pin,
			Source:  m.Source(),
			URL:     m.Download.URL,
			PageURL: m.PageURL(),
		}
		// "Version" surfaces source-specific identifiers for display.
		if m.Update.Modrinth != nil {
			entry.Version = m.Update.Modrinth.Version
		} else if m.Update.CurseForge != nil {
			entry.Version = ""
		}
		out = append(out, entry)
	}
	writeJSON(w, http.StatusOK, out)
}
