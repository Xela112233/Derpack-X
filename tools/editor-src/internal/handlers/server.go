package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/derpack/derpack-edit/internal/config"
	"github.com/derpack/derpack-edit/internal/packwiz"
)

// Server holds the dependencies used by all handlers.
type Server struct {
	RepoRoot string
	Runner   *packwiz.Runner
	Config   *config.Config
}

// errorResp is the standard error envelope returned to the UI.
type errorResp struct {
	Error string `json:"error"`
}

// writeJSON sends a JSON response with the given status code.
func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("writeJSON: %v", err)
	}
}

// writeError sends a 4xx or 5xx error response.
func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, errorResp{Error: msg})
}

// requirePost rejects non-POST requests with 405.
func requirePost(w http.ResponseWriter, r *http.Request) bool {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST")
		writeError(w, http.StatusMethodNotAllowed, "POST required")
		return false
	}
	return true
}

// decodeJSON reads the request body as JSON into v. Returns false (and writes
// an error response) on failure.
func decodeJSON(w http.ResponseWriter, r *http.Request, v interface{}) bool {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return false
	}
	return true
}
