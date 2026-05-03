package handlers

import (
	"net/http"
	"strings"

	"github.com/derpack/derpack-edit/internal/packwiz"
)

type pinReq struct {
	Slug   string `json:"slug"`
	Pinned bool   `json:"pinned"`
}

func (s *Server) HandlePin(w http.ResponseWriter, r *http.Request) {
	if !requirePost(w, r) {
		return
	}
	var req pinReq
	if !decodeJSON(w, r, &req) {
		return
	}
	req.Slug = strings.TrimSpace(req.Slug)
	if req.Slug == "" {
		writeError(w, http.StatusBadRequest, "slug is required")
		return
	}
	if !packwiz.ManifestExists(s.RepoRoot, req.Slug) {
		writeError(w, http.StatusNotFound, "no manifest for slug: "+req.Slug)
		return
	}

	var (
		out string
		err error
	)
	if req.Pinned {
		out, err = s.Runner.Pin(req.Slug)
	} else {
		out, err = s.Runner.Unpin(req.Slug)
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
	writeJSON(w, http.StatusOK, opResp{OK: true, Slug: req.Slug, Output: out})
}
