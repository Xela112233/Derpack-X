#!/usr/bin/env bash
# Builds a server-ready zip: configs + KubeJS + a fetch-mods.sh that
# pulls every jar from the URLs in mods/*.pw.toml.
#
# Usage on a Linux server (like ishimura):
#   unzip ntnh-aeronautics-server-<version>.zip
#   cd ntnh-aeronautics-server
#   ./fetch-mods.sh
#   ./run.sh

set -euo pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "${REPO_ROOT}"

VERSION="$(grep -E '^version\s*=' pack.toml | head -1 | sed -E 's/.*"([^"]+)".*/\1/')"
NAME="derpack-x-create-server-${VERSION}"
WORK="dist/${NAME}"

rm -rf "${WORK}"
mkdir -p "${WORK}"

# Copy what a server actually needs
for d in config defaultconfigs kubejs; do
    if [[ -d "${d}" ]]; then
        cp -r "${d}" "${WORK}/"
    fi
done

# Generate fetch-mods.sh from packwiz manifests
cat > "${WORK}/fetch-mods.sh" <<'EOF'
#!/usr/bin/env bash
# Auto-generated. Downloads server-side mod jars from the packwiz manifest URLs.
set -euo pipefail
mkdir -p mods
EOF

# Walk mods/*.pw.toml and emit curl lines, skipping client-only mods
if [[ -d mods ]]; then
    for f in mods/*.pw.toml; do
        [[ -e "$f" ]] || continue
        side="$(grep -E '^side\s*=' "$f" | sed -E 's/.*"([^"]+)".*/\1/' || echo both)"
        if [[ "${side}" == "client" ]]; then
            continue
        fi
        url="$(grep -E '^url\s*=' "$f" | sed -E 's/.*"([^"]+)".*/\1/')"
        if [[ -n "${url}" ]]; then
            echo "echo \"  fetching $(basename "$f" .pw.toml)\"" >> "${WORK}/fetch-mods.sh"
            echo "curl -fsSL -o \"mods/$(basename "$f" .pw.toml).jar\" \"${url}\"" >> "${WORK}/fetch-mods.sh"
        fi
    done
fi

chmod +x "${WORK}/fetch-mods.sh"

# Minimal launcher
cat > "${WORK}/run.sh" <<'EOF'
#!/usr/bin/env bash
# Aikar's flags, 12G heap (matches your existing NTNH tuning)
set -e
java -Xms12G -Xmx12G \
    -XX:+UseG1GC -XX:+ParallelRefProcEnabled -XX:MaxGCPauseMillis=200 \
    -XX:+UnlockExperimentalVMOptions -XX:+DisableExplicitGC \
    -XX:+AlwaysPreTouch -XX:G1NewSizePercent=30 -XX:G1MaxNewSizePercent=40 \
    -XX:G1HeapRegionSize=8M -XX:G1ReservePercent=20 -XX:G1HeapWastePercent=5 \
    -XX:G1MixedGCCountTarget=4 -XX:InitiatingHeapOccupancyPercent=15 \
    -XX:G1MixedGCLiveThresholdPercent=90 -XX:G1RSetUpdatingPauseTimePercent=5 \
    -XX:SurvivorRatio=32 -XX:+PerfDisableSharedMem -XX:MaxTenuringThreshold=1 \
    -Dusing.aikars.flags=https://mcflags.emc.gs -Daikars.new.flags=true \
    @user_jvm_args.txt @libraries/net/neoforged/neoforge/*/unix_args.txt \
    "$@"
EOF
chmod +x "${WORK}/run.sh"

cat > "${WORK}/README.txt" <<EOF
Derpack X: Create — Server Pack ${VERSION}

1. Install NeoForge 21.1.x for Minecraft 1.21.1 in this directory:
     https://neoforged.net/

2. Run:
     ./fetch-mods.sh

3. Accept the EULA in eula.txt

4. Start the server:
     ./run.sh nogui

Notes:
- Tune memory in run.sh if 12G is wrong for your box.
- Configs in config/ overwrite defaults; defaultconfigs/ apply only on first run.
- Aeronautics is in active development — back up worlds before updating.
EOF

(cd dist && zip -qr "${NAME}.zip" "${NAME}")
echo "Built: dist/${NAME}.zip"
ls -lh "dist/${NAME}.zip"
