# Derpack X: Create

A Create-focused Minecraft modpack for 1.21.1 / NeoForge, built around [Create Aeronautics](https://modrinth.com/mod/create-aeronautics).

---

## For players: how to install

Pick **one** of the three downloads on the latest [release](../../releases) page:

| File | What it is | When to use |
|------|-----------|-------------|
| `derpack-x-create-prism-bundled-X.Y.Z.zip` | Prism instance with all mod jars baked in | **Recommended.** Drop it into Prism, play. Big file. |
| `derpack-x-create-prism-installer-X.Y.Z.zip` | Prism instance that downloads mods on first launch | Smaller download, auto-updates on every launch. Needs internet on first run. |
| `derpack-x-create-X.Y.Z.mrpack` | Modrinth-format pack file | If you use the Modrinth app, or want to import into Prism manually. |

**Steps:** download → open Prism Launcher → Add Instance → Import from zip → pick the file → launch.

Memory and Java args are pre-configured. Java 21 is required (Prism will offer to download it if missing).

See [`docs/PRISM-SETUP.md`](docs/PRISM-SETUP.md) for the full walkthrough.

---

## For collaborators (you, me, contributors)

This repo doesn't store mod jars — it stores a list of *what mods are in the pack* (with URLs). When we cut a release, GitHub Actions downloads the jars and builds the zips users actually install.

The list of mods to read first:

- **[`docs/MODLIST.md`](docs/MODLIST.md)** — human-readable list of every mod, what it does, and why it's in the pack. Start here.
- **[`docs/EDITING.md`](docs/EDITING.md)** — how to add a mod, drop in a config, or change something. Written for people who've never touched packwiz.

### How to make changes

Most edits happen through the GitHub UI — no local install needed:

- **Add a mod** → Actions tab → "Add mod" workflow → Run
- **Remove a mod** → Actions tab → "Remove mod" workflow → Run
- **Update mods** → Actions tab → "Update mods" workflow → Run (blank slug = update all)
- **Drop in a config / shaderpack / KubeJS script** → upload via the web UI, then trigger "Update mods" to refresh the index

Full walkthrough in [`docs/EDITING.md`](docs/EDITING.md).

### Where things live

```
.
├── pack.toml              # Pack name, version, Minecraft + loader version
├── mods/                  # One small text file per mod (auto-managed, don't edit by hand)
├── config/                # Mod configs that ALWAYS get applied to the user's instance
├── defaultconfigs/        # Mod configs only applied on FRESH install (user can override)
├── kubejs/                # KubeJS scripts for recipe tweaks etc.
├── resourcepacks/         # Bundled resource packs (.zip files)
├── shaderpacks/           # Bundled shaderpacks (.zip files)
├── docs/                  # Human-facing documentation
├── scripts/               # Build helpers (you almost never run these — CI does)
└── .github/workflows/     # CI: builds the release zips when you publish a release
```

Each folder has its own `README.md` explaining what goes inside it.

### Cutting a release

1. Bump `version` in `pack.toml` (e.g. `0.1.0` → `0.1.1`)
2. Commit and push to `main`
3. On GitHub: Releases → Draft a new release → tag `v0.1.1` → Publish
4. Wait ~5 minutes for the build, then check the release page for attached zips

---

## Known issues

- **Iris/Oculus shaders glitch on Aeronautics ships.** Disable shaders if airships look broken. The mod team is aware.
- **Some non-Aeronautics mods misbehave on physics ships.** That's why `aeronauticscompat` is in the pack — it patches Etched, Cobblemon, sentry turrets, sleeping, chains, and others.

## License

Pack configuration and scripts: see [LICENSE](LICENSE) (MIT). Bundled mods retain their original licenses; this repo does not redistribute mod jars — the build pipeline downloads them at release time.
