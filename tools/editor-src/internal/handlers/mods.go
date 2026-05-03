package handlers

import (
	"net/http"
	"strings"

	"github.com/derpack/derpack-edit/internal/packwiz"
)

type addReq struct {
	Source string `json:"source"` // "mr" or "cf"
	Slug   string `json:"slug"`
	Side   string `json:"side"` // "both", "client", "server", or ""
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
		out, err = s.Runner.AddCurseForge(req.Slug)
	default:
		writeError(w, http.StatusBadRequest, "source must be 'mr' or 'cf'")
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

	// Apply the user's side choice. packwiz sets the side based on the mod's
	// Modrinth metadata (client_side / server_side), but that's frequently
	// wrong — most Create addons need "both" but get tagged "client" because
	// the author misconfigured their Modrinth project. The user's pick wins.
	if req.Side != "" {
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
