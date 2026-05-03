#!/usr/bin/env bash
# Build script for derpack-edit.
#
# Downloads the packwiz binary for Windows, embeds it into the Go binary,
# and cross-compiles a Windows .exe.
#
# Usage:
#   bash build.sh                 # default: builds for windows/amd64
#   GOOS=linux GOARCH=amd64 bash build.sh   # builds for current platform (dev)
#
# Output: bin/derpack-edit.exe (or bin/derpack-edit on non-Windows)

set -euo pipefail

cd "$(dirname "$0")"

OUT_DIR="${OUT_DIR:-bin}"
TARGET_GOOS="${GOOS:-windows}"
TARGET_GOARCH="${GOARCH:-amd64}"

# Latest packwiz binary URL — may need updating if packwiz reorganizes their
# release format. Pin a specific release for reproducibility.
PACKWIZ_VERSION="${PACKWIZ_VERSION:-latest}"
PACKWIZ_URL_BASE="https://nightly.link/packwiz/packwiz/workflows/go/main"

# For now pin to the prebuilt nightly (packwiz's distribution model uses
# Actions artifacts rather than tagged releases for binaries).
case "${TARGET_GOOS}" in
  windows) PACKWIZ_ARTIFACT="packwiz-windows-amd64.zip"; PACKWIZ_BIN="packwiz.exe" ;;
  linux)   PACKWIZ_ARTIFACT="packwiz-linux-amd64.zip"; PACKWIZ_BIN="packwiz" ;;
  darwin)  PACKWIZ_ARTIFACT="packwiz-darwin-arm64.zip"; PACKWIZ_BIN="packwiz" ;;
  *) echo "unsupported GOOS: $TARGET_GOOS"; exit 1 ;;
esac

PACKWIZ_URL="${PACKWIZ_URL_BASE}/${PACKWIZ_ARTIFACT}"

mkdir -p "$OUT_DIR"
TMP_DIR="$(mktemp -d)"
trap "rm -rf $TMP_DIR" EXIT

echo "==> Downloading packwiz from $PACKWIZ_URL"
if ! curl -fsSL "$PACKWIZ_URL" -o "$TMP_DIR/packwiz.zip"; then
  echo "ERROR: could not download packwiz binary."
  echo "       The URL pattern in this script may be outdated."
  echo "       Update PACKWIZ_URL_BASE in build.sh to point at a working source."
  exit 1
fi

echo "==> Extracting packwiz binary"
unzip -q -o "$TMP_DIR/packwiz.zip" -d "$TMP_DIR/extracted"

# Find the binary inside the extracted folder (packwiz zips have it at root).
PACKWIZ_PATH=""
for candidate in "$TMP_DIR/extracted/$PACKWIZ_BIN" "$TMP_DIR/extracted/packwiz" "$TMP_DIR/extracted/packwiz.exe"; do
  if [[ -f "$candidate" ]]; then
    PACKWIZ_PATH="$candidate"
    break
  fi
done

if [[ -z "$PACKWIZ_PATH" ]]; then
  echo "ERROR: could not find packwiz binary in zip. Contents:"
  ls -la "$TMP_DIR/extracted"
  exit 1
fi

echo "==> Embedding packwiz at internal/packwiz/assets/packwiz.bin"
cp "$PACKWIZ_PATH" internal/packwiz/assets/packwiz.bin

echo "==> Compiling derpack-edit ($TARGET_GOOS/$TARGET_GOARCH)"
ext=""
[[ "$TARGET_GOOS" == "windows" ]] && ext=".exe"

GOOS="$TARGET_GOOS" GOARCH="$TARGET_GOARCH" \
  go build -o "$OUT_DIR/derpack-edit$ext" \
    -ldflags="-s -w" \
    .

echo ""
echo "Built: $OUT_DIR/derpack-edit$ext"
ls -lh "$OUT_DIR/derpack-edit$ext"
