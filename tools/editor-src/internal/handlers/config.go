package handlers

import (
	"net/http"
	"os"
	"path/filepath"
	"runtime"

	"github.com/derpack/derpack-edit/internal/config"
)

// configResp is what GET /api/config returns. Includes default Prism path
// suggestions for the OS so the UI can offer auto-detect.
type configResp struct {
	PrismInstancePath  string   `json:"prism_instance_path"`
	PrismDefaultRoot   string   `json:"prism_default_root,omitempty"` // %APPDATA%\PrismLauncher\instances on Windows
	DetectedInstances  []string `json:"detected_instances,omitempty"` // names of folders found there
}

func (s *Server) HandleGetConfig(w http.ResponseWriter, r *http.Request) {
	resp := configResp{
		PrismInstancePath: s.Config.PrismInstancePath,
	}

	if root := defaultPrismInstancesRoot(); root != "" {
		resp.PrismDefaultRoot = root
		if entries, err := os.ReadDir(root); err == nil {
			for _, e := range entries {
				if e.IsDir() {
					// Ignore Prism's internal folders.
					if e.Name() == "_LAUNCHER_TEMP" || e.Name() == "_MMC_TEMP" {
						continue
					}
					resp.DetectedInstances = append(resp.DetectedInstances, e.Name())
				}
			}
		}
	}

	writeJSON(w, http.StatusOK, resp)
}

type setConfigReq struct {
	PrismInstancePath string `json:"prism_instance_path"`
}

func (s *Server) HandleSetConfig(w http.ResponseWriter, r *http.Request) {
	if !requirePost(w, r) {
		return
	}
	var req setConfigReq
	if !decodeJSON(w, r, &req) {
		return
	}

	// Validate that the path looks like a Prism instance. We accept either:
	// - the instance dir itself (containing instance.cfg)
	// - the .minecraft subdir (since some users might point at that directly)
	if req.PrismInstancePath != "" {
		st, err := os.Stat(req.PrismInstancePath)
		if err != nil {
			writeError(w, http.StatusBadRequest,
				"path does not exist: "+req.PrismInstancePath)
			return
		}
		if !st.IsDir() {
			writeError(w, http.StatusBadRequest,
				"path is not a directory: "+req.PrismInstancePath)
			return
		}
		// Check for telltale files. Either instance.cfg (instance root) or
		// mods/+config/ (already inside .minecraft).
		instanceCfg := filepath.Join(req.PrismInstancePath, "instance.cfg")
		minecraftSub := filepath.Join(req.PrismInstancePath, ".minecraft")
		if _, err1 := os.Stat(instanceCfg); err1 != nil {
			if _, err2 := os.Stat(minecraftSub); err2 != nil {
				writeError(w, http.StatusBadRequest,
					"path doesn't look like a Prism instance: missing instance.cfg and .minecraft/. Make sure you point at the instance folder, not the launcher folder.")
				return
			}
		}
	}

	s.Config.PrismInstancePath = req.PrismInstancePath
	if err := config.Save(s.RepoRoot, s.Config); err != nil {
		writeError(w, http.StatusInternalServerError, "save config: "+err.Error())
		return
	}
	writeJSON(w, http.StatusOK, opResp{OK: true})
}

// defaultPrismInstancesRoot returns the standard location of Prism Launcher
// instances on the current OS, or "" if unknown.
func defaultPrismInstancesRoot() string {
	switch runtime.GOOS {
	case "windows":
		appdata := os.Getenv("APPDATA")
		if appdata == "" {
			return ""
		}
		return filepath.Join(appdata, "PrismLauncher", "instances")
	case "darwin":
		home, err := os.UserHomeDir()
		if err != nil {
			return ""
		}
		return filepath.Join(home, "Library", "Application Support", "PrismLauncher", "instances")
	case "linux":
		// Prism on Linux uses XDG_DATA_HOME or ~/.local/share/PrismLauncher
		if dataHome := os.Getenv("XDG_DATA_HOME"); dataHome != "" {
			return filepath.Join(dataHome, "PrismLauncher", "instances")
		}
		home, err := os.UserHomeDir()
		if err != nil {
			return ""
		}
		return filepath.Join(home, ".local", "share", "PrismLauncher", "instances")
	}
	return ""
}
