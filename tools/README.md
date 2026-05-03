# Derpack Editor

A local web app for managing the modpack without using the command line or GitHub Actions workflows.

## Quick start

1. Make sure you've cloned the repo with GitHub Desktop and have the latest `main` (or your working branch).
2. Open the repo folder in File Explorer.
3. Double-click `tools\derpack-edit.exe`.
4. A console window appears, then your browser opens to `http://localhost:8765`.
5. Edit mods. Use GitHub Desktop to commit and push when done.

## What it does

- Lists every mod in the pack
- Add a mod by slug (Modrinth or CurseForge)
- Remove a mod
- Pin / unpin a mod (pinned mods don't get auto-updated)
- Set a specific version of a mod (pick from Modrinth versions or paste a URL)
- Update one mod or all mods at once (skips pinned)
- Recompute a mod's hash (after manual manifest edits)
- Build the modpack and copy it directly into your local Prism instance for testing

## What it doesn't do

- **Doesn't handle git.** Use GitHub Desktop for clone, branch, commit, push, and PRs. This tool only edits files in your working directory — git records those changes, and GitHub Desktop helps you push them.
- **Doesn't run CI builds for distribution.** Releases still happen via the existing build workflow when you publish a release on GitHub.
- **Doesn't replace the GitHub Actions workflows.** They still work — the editor is a friendlier alternative for local editing.

## Build & Launch in Prism

When you click the green ▶ Build & Launch in Prism button:

1. The editor builds the modpack locally (downloads jars, copies configs, builds Prism instance scaffold)
2. Copies everything into the Prism Launcher instance you've configured under Settings (⚙)
3. Tells you to open Prism and click Play

**Requirements:**
- Java 21 must be on your system PATH. If you have Prism installed and have launched a modpack at least once, you might already have a Java install at `C:\Program Files\Eclipse Adoptium\jdk-21...\bin`. Add that to your PATH if needed.
- A configured Prism instance — open Settings (⚙) and either pick from auto-detected instances or paste a path.

**What gets synced:**
- `mods/` — replaced entirely from the build
- `config/`, `defaultconfigs/`, `kubejs/`, `resourcepacks/`, `shaderpacks/` — replaced from the repo

**What doesn't get touched:**
- Worlds and saves
- Anything outside the dirs listed above

## Local config

The editor stores your Prism instance path in `.editor/config.json`. That file is gitignored — your friend won't see your path, you won't see theirs. Each user has their own.

## Running it

The .exe is committed to the repo at `tools/derpack-edit.exe`. You don't need to build anything yourself; CI rebuilds it whenever the source under `tools/editor-src/` changes.

To update the editor itself (rare): edit files under `tools/editor-src/`, commit, push to main. The build-editor workflow rebuilds and commits a new `tools/derpack-edit.exe`.

## Troubleshooting

**"could not find pack.toml"**  
You launched the editor from outside the repo. Move the `.exe` into the repo's `tools/` folder, or open a terminal in the repo and run it from there.

**"could not bind to any port"**  
Something else is using ports 8765-8774. Close other apps and try again, or restart your machine.

**Browser doesn't open automatically**  
The console window will print the URL — copy it into your browser manually.

**Build fails with "java not found"**  
Install Java 21 (Adoptium Temurin is the easy choice: https://adoptium.net/) or add Prism's bundled Java to your PATH.

**Build & Launch hangs**  
First build downloads ~500MB of mod jars and can take 5+ minutes. Watch the console window for progress messages. If it's truly hung (no output for 10+ minutes), close the console (X) to kill it, then re-launch.

**My changes aren't showing up in the table**  
Click the ↻ Refresh button. The editor reads from disk fresh, but the UI caches the list between operations.
