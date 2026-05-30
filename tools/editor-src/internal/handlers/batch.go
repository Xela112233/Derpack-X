package handlers

import (
	"net/http"
	"strings"
)

type batchAddReq struct {
	Lines []string `json:"lines"`
}

// batchAddResult is per-line result, returned in the same order as input.
type batchAddResult struct {
	Line   string `json:"line"`
	Slug   string `json:"slug,omitempty"`   // resolved slug (if line parsed)
	Source string `json:"source,omitempty"` // "mr" or "cf"
	OK     bool   `json:"ok"`
	Error  string `json:"error,omitempty"`
	Output string `json:"output,omitempty"`
}

type batchAddResp struct {
	Results []batchAddResult `json:"results"`
}

// HandleBatchAdd processes a list of wishlist-format lines like:
//
//	mr:supplementaries
//	cf:create-mob-spawners
//	# this is a comment
//	supplementaries (defaults to mr)
//
// Each line is independent — failures don't stop processing.
func (s *Server) HandleBatchAdd(w http.ResponseWriter, r *http.Request) {
	if !requirePost(w, r) {
		return
	}
	var req batchAddReq
	if !decodeJSON(w, r, &req) {
		return
	}

	results := make([]batchAddResult, 0, len(req.Lines))
	for _, raw := range req.Lines {
		line := strings.TrimSpace(raw)
		// Skip empties and comments.
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		res := batchAddResult{Line: raw}
		source, slug := parseWishlistLine(line)
		if slug == "" {
			res.Error = "could not parse line; expected 'mr:slug', 'cf:slug', or just 'slug'"
			results = append(results, res)
			continue
		}
		res.Source = source
		res.Slug = slug

		var (
			out string
			err error
		)
		switch source {
		case "mr":
			out, err = s.Runner.AddModrinth(slug, "")
		case "cf":
			// CF goes through our own resolution path so we can pass the
			// API key, detect opt-outs cleanly, and enforce the no-half-
			// state contract.
			out, err = s.addCurseForgeMod(slug)
		}

		if err != nil {
			res.Error = err.Error()
			res.Output = out
		} else {
			// Force "both" — packwiz often picks up wrong side metadata
			// from Modrinth for Create addons and similar. The user can
			// flip a specific mod with the per-row Side button if needed.
			// CF additions already write with no Side field (matching
			// "both"), but applying again is harmless and keeps the path
			// uniform.
			if sideErr := s.applySide(slug, "both"); sideErr != nil {
				// Mod added but side wasn't normalized. Annotate but don't fail.
				res.OK = true
				res.Output = out + "\nNOTE: side normalization failed: " + sideErr.Error()
			} else {
				res.OK = true
				res.Output = out
			}
		}
		results = append(results, res)
	}

	writeJSON(w, http.StatusOK, batchAddResp{Results: results})
}

// parseWishlistLine accepts forms:
//
//	mr:slug
//	cf:slug
//	slug         (defaults to "mr")
//
// Returns ("", "") if the line couldn't be parsed.
func parseWishlistLine(line string) (source, slug string) {
	if i := strings.Index(line, ":"); i > 0 {
		prefix := strings.ToLower(strings.TrimSpace(line[:i]))
		rest := strings.TrimSpace(line[i+1:])
		if rest == "" {
			return "", ""
		}
		switch prefix {
		case "mr", "modrinth":
			return "mr", rest
		case "cf", "curseforge":
			return "cf", rest
		default:
			return "", ""
		}
	}
	// No prefix — default to Modrinth, treat the whole line as the slug.
	if isPlausibleSlug(line) {
		return "mr", line
	}
	return "", ""
}

// isPlausibleSlug rejects lines with whitespace or wildly non-slug chars,
// since those are probably typos rather than real slugs.
func isPlausibleSlug(s string) bool {
	if s == "" {
		return false
	}
	for _, c := range s {
		if c == '-' || c == '_' || c == '.' {
			continue
		}
		if (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') {
			continue
		}
		return false
	}
	return true
}
