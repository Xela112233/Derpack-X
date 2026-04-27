#!/usr/bin/env bash
# Builds a .mrpack file from the packwiz manifest.
# Output: dist/ntnh-aeronautics-<version>.mrpack

set -euo pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "${REPO_ROOT}"

if ! command -v packwiz >/dev/null 2>&1; then
    echo "ERROR: packwiz not in PATH. Run scripts/setup.sh first." >&2
    exit 1
fi

# Pull version out of pack.toml
VERSION="$(grep -E '^version\s*=' pack.toml | head -1 | sed -E 's/.*"([^"]+)".*/\1/')"
NAME="$(grep -E '^name\s*=' pack.toml | head -1 | sed -E 's/.*"([^"]+)".*/\1/' | tr '[:upper:]' '[:lower:]' | sed 's/[^a-z0-9]/-/g; s/--*/-/g; s/^-//; s/-$//')"

mkdir -p dist
OUT="dist/${NAME}-${VERSION}.mrpack"

echo "==> Refreshing packwiz index"
packwiz refresh

echo "==> Exporting Modrinth modpack"
packwiz mr export -o "${OUT}"

echo
echo "Built: ${OUT}"
ls -lh "${OUT}"
