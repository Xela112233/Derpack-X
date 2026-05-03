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
- Add a mod (Modrinth or CurseForge) by slug
- Remove a mod
- Pin / unpin a mod (pinned mods don't get auto-updated)
- More features coming in future versions: set specific version, compute hash, build & launch in Prism

## What it doesn't do

- **Doesn't handle git.** Use GitHub Desktop for clone, branch, commit, push, PR. This tool only edits files in your working directory.
- **Doesn't run CI builds.** Pushing to a branch still triggers the build workflow on GitHub.
- **Doesn't replace the GitHub workflows.** They still work; the editor is just a friendlier alternative.

## Running it

The .exe is committed to the repo at `tools/derpack-edit.exe`. You don't need to build anything yourself.

If you want to update the editor itself (rare), edit files under `tools/editor-src/` and push — a workflow will rebuild and commit a new `tools/derpack-edit.exe`.

## Troubleshooting

**"could not find pack.toml"**  
You launched the editor from outside the repo. Move the `.exe` into the repo's `tools/` folder, or open a terminal in the repo and run it from there.

**"could not bind to any port"**  
Something else is using ports 8765–8774. Close other apps and try again, or restart your machine.

**Browser doesn't open automatically**  
The console window will print the URL — copy it into your browser manually.

**My changes aren't showing up**  
Click the ↻ Refresh button in the top right. The editor reads from disk on each request, but the UI caches the list until you refresh.
