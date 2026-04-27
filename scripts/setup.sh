#!/usr/bin/env bash
# Sets up the dev environment for working on this modpack.
# Idempotent — safe to run repeatedly.

set -euo pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
TOOLS_DIR="${REPO_ROOT}/.tools"
mkdir -p "${TOOLS_DIR}"

echo "==> Repo root: ${REPO_ROOT}"

# --- packwiz ---
if ! command -v packwiz >/dev/null 2>&1; then
    echo "==> packwiz not found in PATH"
    echo "    Install it with one of:"
    echo "      go install github.com/packwiz/packwiz@latest"
    echo "      paru -S packwiz-bin                              # Arch / Fedora COPR"
    echo "      brew install packwiz                             # macOS"
    echo "    or download a binary from https://github.com/packwiz/packwiz/releases"
    exit 1
fi
echo "==> packwiz: $(packwiz --version 2>/dev/null || echo present)"

# --- packwiz-installer-bootstrap (for Prism / clients) ---
BOOTSTRAP_JAR="${TOOLS_DIR}/packwiz-installer-bootstrap.jar"
if [[ ! -f "${BOOTSTRAP_JAR}" ]]; then
    echo "==> Downloading packwiz-installer-bootstrap"
    curl -fsSL -o "${BOOTSTRAP_JAR}" \
        "https://github.com/packwiz/packwiz-installer-bootstrap/releases/latest/download/packwiz-installer-bootstrap.jar"
fi
echo "==> bootstrap jar: ${BOOTSTRAP_JAR}"

echo
echo "Done. Next steps:"
echo "  packwiz init                # if pack.toml doesn't exist yet"
echo "  packwiz mr add create        # add a Modrinth mod"
echo "  packwiz cf add <slug>        # add a CurseForge mod"
echo "  packwiz refresh              # update index hashes"
echo "  ./scripts/build-mrpack.sh    # produce a .mrpack"
