#!/usr/bin/env bash
# Builds the "jars bundled" Prism instance zip.
#
# Approach: spin up `packwiz serve` (its built-in local HTTP server),
# point packwiz-installer-bootstrap at it, let it fetch every jar
# (including CurseForge ones, which the bootstrap resolves via the CF API).
# Then bundle the resulting mods/ folder into a Prism instance zip.
#
# This is the documented, supported path — same one users would run on
# their own machine, just driven by us at build time.
#
# Result: dist/<slug>-prism-bundled-<version>.zip

set -euo pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "${REPO_ROOT}"

PACK_NAME="$(grep -E '^name\s*=' pack.toml | head -1 | sed -E 's/.*"([^"]+)".*/\1/')"
PACK_VERSION="$(grep -E '^version\s*=' pack.toml | head -1 | sed -E 's/.*"([^"]+)".*/\1/')"
SLUG="$(echo "${PACK_NAME}" | tr '[:upper:]' '[:lower:]' | sed 's/[^a-z0-9]/-/g; s/--*/-/g; s/^-//; s/-$//')"

INSTANCE_NAME="${SLUG}-prism-bundled-${PACK_VERSION}"
STAGING="dist/staging/${INSTANCE_NAME}"
OUT="dist/${INSTANCE_NAME}.zip"

rm -rf "${STAGING}"
mkdir -p dist

# 1. Build the Prism instance skeleton
bash scripts/build-prism-skeleton.sh "${STAGING}" "bundled"

# 2. Get the bootstrap jar
BOOTSTRAP_JAR="${REPO_ROOT}/.tools/packwiz-installer-bootstrap.jar"
if [[ ! -f "${BOOTSTRAP_JAR}" ]]; then
    mkdir -p "${REPO_ROOT}/.tools"
    echo "==> Downloading packwiz-installer-bootstrap"
    curl -fsSL -o "${BOOTSTRAP_JAR}" \
        "https://github.com/packwiz/packwiz-installer-bootstrap/releases/latest/download/packwiz-installer-bootstrap.jar"
fi

# 3. Start `packwiz serve` in the background, fetch jars, kill it
PORT=8765
echo "==> Starting packwiz serve on :${PORT}"
packwiz serve --port "${PORT}" >/tmp/packwiz-serve.log 2>&1 &
SERVE_PID=$!
trap "kill ${SERVE_PID} 2>/dev/null || true" EXIT

# Wait for server to be reachable (max ~10s)
for i in $(seq 1 20); do
    if curl -fs "http://localhost:${PORT}/pack.toml" >/dev/null 2>&1; then
        break
    fi
    sleep 0.5
done
if ! curl -fs "http://localhost:${PORT}/pack.toml" >/dev/null 2>&1; then
    echo "ERROR: packwiz serve never came up. Log:" >&2
    cat /tmp/packwiz-serve.log >&2
    exit 1
fi

echo "==> Fetching jars via packwiz-installer-bootstrap"
(
    cd "${STAGING}/.minecraft"
    java -jar "${BOOTSTRAP_JAR}" -g -s client "http://localhost:${PORT}/pack.toml"
)

# Kill the server
kill ${SERVE_PID} 2>/dev/null || true
trap - EXIT

# bootstrap leaves these behind; not needed in the shipped zip
rm -f "${STAGING}/.minecraft/packwiz.json" "${STAGING}/.minecraft/packwiz-installer-bootstrap.jar"

# Sanity check
mod_count="$(find "${STAGING}/.minecraft/mods" -name '*.jar' 2>/dev/null | wc -l)"
if [[ "${mod_count}" -eq 0 ]]; then
    echo "ERROR: no mod jars in ${STAGING}/.minecraft/mods after install" >&2
    exit 1
fi
echo "==> ${mod_count} mod jars bundled"

# 4. README for friends
cat > "${STAGING}/README.txt" <<EOF
${PACK_NAME} ${PACK_VERSION} — Prism Instance (bundled jars)

Install:
  1. Open Prism Launcher
  2. Add Instance -> Import from zip
  3. Pick this zip
  4. Launch

Memory: 8-12 GB (configured)
Java:   21 (Prism will prompt to download if missing)
EOF

# 5. Zip
echo "==> Packaging zip"
(cd dist/staging && zip -qr "../${INSTANCE_NAME}.zip" "${INSTANCE_NAME}")

ls -lh "${OUT}"
echo "Built: ${OUT}"
