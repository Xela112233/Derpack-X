package packwiz

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
)

// embeddedBinary is the packwiz executable, embedded at build time.
// At dev time this file is empty; the real binary is downloaded by the
// build script and placed at this path before `go build`.
//
// On Windows the build script will fetch packwiz_windows_amd64.exe; on
// other platforms (only used for dev iteration here) it falls back to
// the Linux binary.
//
//go:embed assets/packwiz.bin
var embeddedBinary []byte

// EnsureBinary writes the embedded packwiz binary to .editor/packwiz-cache/
// if it isn't already there, and returns the absolute path. Idempotent.
func EnsureBinary(repoRoot string) (string, error) {
	if len(embeddedBinary) == 0 {
		return "", fmt.Errorf("packwiz binary was not embedded at build time (assets/packwiz.bin is empty)")
	}

	cacheDir := filepath.Join(repoRoot, ".editor", "packwiz-cache")
	if err := os.MkdirAll(cacheDir, 0o755); err != nil {
		return "", fmt.Errorf("mkdir %s: %w", cacheDir, err)
	}

	binName := "packwiz"
	// We assume Windows targets since that's our v1 platform; the .exe
	// suffix is harmless on Unix and Windows requires it.
	binName += ".exe"

	binPath := filepath.Join(cacheDir, binName)

	// If it already exists with the correct size, skip writing. Avoids a
	// pointless write on every launch and respects file locks if packwiz
	// happens to be running.
	if st, err := os.Stat(binPath); err == nil && st.Size() == int64(len(embeddedBinary)) {
		return binPath, nil
	}

	if err := os.WriteFile(binPath, embeddedBinary, 0o755); err != nil {
		return "", fmt.Errorf("write %s: %w", binPath, err)
	}
	return binPath, nil
}
