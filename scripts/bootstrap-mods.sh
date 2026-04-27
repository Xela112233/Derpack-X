#!/usr/bin/env bash
# Bootstraps the initial mod list. Run once after `packwiz init`.
# Idempotent: packwiz add will skip mods that are already present.

set -euo pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "${REPO_ROOT}"

if ! command -v packwiz >/dev/null 2>&1; then
    echo "ERROR: packwiz not in PATH." >&2
    exit 1
fi

# Helper - tolerate a slug not existing on a given platform
add_mr() { packwiz mr add "$1" --yes 2>&1 | sed "s/^/  [mr $1] /" || echo "  [mr $1] FAILED, skip"; }
add_cf() { packwiz cf add "$1" --yes 2>&1 | sed "s/^/  [cf $1] /" || echo "  [cf $1] FAILED, skip"; }

echo "==> Required: Create + Sable + Aeronautics + Compat"
add_mr create
add_mr sable
add_mr create-aeronautics
add_cf create-aeronautics-compatability

echo "==> Performance"
add_mr embeddium
add_mr embeddiumplus
add_mr ferrite-core
add_mr modernfix
add_mr saturn
add_mr oculus

echo "==> QoL"
add_mr jei
add_mr jade
add_mr mouse-tweaks
add_mr xaeros-minimap
add_mr xaeros-world-map
add_mr inventory-profiles-next

echo "==> Create addons"
add_mr create-steam-n-rails
add_mr create-new-age
add_mr create-enchantment-industry
add_cf create-diesel-generators

echo "==> Adventure / content"
add_mr alexs-mobs
add_mr farmers-delight
add_mr sophisticated-backpacks
add_mr irons-spells-n-spellbooks
add_mr l_endsmobs    # Cataclysm slug varies; adjust manually
add_mr supplementaries

echo
echo "==> Refreshing index"
packwiz refresh

echo
echo "Done. Review mods/*.pw.toml and commit."
echo "Some slugs may have failed if they don't exist on Modrinth — check the log,"
echo "edit this script, or add them manually with 'packwiz cf add <slug>'."
