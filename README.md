# Derpack X: Create

A Create-focused modpack for Minecraft 1.21.1 / NeoForge built around [Create Aeronautics](https://modrinth.com/mod/create-aeronautics) — physics-driven airships, planes, off-road vehicles, and the rest of the Create ecosystem.

> **Status:** Pre-alpha. Pinned to Create Aeronautics 1.1.3 (released April 2026). Expect breakage on Aeronautics updates.

## Quick start (Prism Launcher)

1. Install [Prism Launcher](https://prismlauncher.org/).
2. Add a new instance: **Minecraft 1.21.1 → NeoForge 21.1.105**.
3. Allocate **8GB minimum** (12GB recommended for the kitchen-sink mod list).
4. Close Prism.
5. Drop the contents of the latest [release](../../releases) into the instance's `.minecraft` folder, overwriting.
6. Launch.

> Mod jars are not in this repo. They live in the GitHub Releases attached zip (private — message me for access if you're not on the friend list).

## Structure

```
.
├── pack.toml              # packwiz manifest
├── mods/                  # one .pw.toml per mod (URL + hash, no jars)
├── config/                # per-mod configs we ship
├── defaultconfigs/        # configs only applied to fresh instances
├── kubejs/                # KubeJS scripts (recipes, tweaks)
├── resourcepacks/         # bundled resource packs
├── shaderpacks/           # bundled shaders (note: Aeronautics has shader issues)
├── scripts/               # build + sync helpers
├── docs/
│   └── MODLIST.md         # human-readable mod list
└── .github/workflows/     # CI: build .mrpack on tag
```

## Development workflow

```bash
# one-time setup
./scripts/setup.sh

# add a mod
packwiz mr add create-aeronautics      # Modrinth
packwiz cf add aeronautics-compat      # CurseForge

# refresh hashes after a mod update
packwiz refresh

# build a .mrpack locally
./scripts/build-mrpack.sh
```

Tag a release (`v0.1.0`) and the GitHub Action will publish a `.mrpack` automatically.

## Known issues

- **Iris/Oculus shaders have visual bugs with Aeronautics ships** (per the mod page). Disable shaders if airships look broken.
- **Mods that don't use `Sable Companion` may misbehave on physics ships.** That's why `aeronauticscompat` is in the pack — it patches Etched, Cobblemon, sentry turrets, sleeping, chains, etc.
- **`itemphysic` and `FoamFix`** — leaving notes here from prior pack experience: client-only, do not put on the server.

## License

Pack configuration and scripts: see [LICENSE](LICENSE). Bundled mods retain their original licenses; this repo does not redistribute mod jars.
