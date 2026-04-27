# Editing this pack

Quick reference for adding mods, dropping in configs, tweaking behavior. Two ways to do most things:

- **Through GitHub** (no install required) — click buttons in the Actions tab. Works for adding/removing/updating mods.
- **Locally with `packwiz` installed** — needed for dropping in config files, KubeJS scripts, or anything that requires inspecting actual game output. See "Local setup" near the bottom.

Default to the GitHub flow when you can. It leaves a clean paper trail of every change in your Actions history.

---

## I want to ADD a mod

1. Go to the repo on GitHub → **Actions** tab.
2. Click **Add mod** in the sidebar.
3. Click **Run workflow** (top right).
4. Fill in the form:
   - **Source**: Modrinth or CurseForge — pick whichever has the mod
   - **Slug**: the last part of the mod's URL
     - Modrinth: `modrinth.com/mod/jei` → slug is `jei`
     - CurseForge: `curseforge.com/minecraft/mc-mods/create-aeronautics-compatability` → slug is `create-aeronautics-compatability`
   - **Side**: usually `both`. Use `client` for things like minimaps and JEI; `server` is rare.
5. Click **Run workflow**.
6. Wait ~30 seconds. The workflow runs `packwiz add`, commits, and pushes to `main`.
7. Refresh the repo — you'll see a new commit "Add `<slug>` from `<source>`" by `github-actions[bot]`.

If the workflow fails: open the run, scroll to the failed step, and the log will say why (most common: wrong slug, mod not on the platform you picked, no version available for 1.21.1 NeoForge).

After adding mods, also update [`MODLIST.md`](MODLIST.md) so the human-readable list stays accurate. You can edit it directly in the GitHub UI (pencil icon).

---

## I want to REMOVE a mod

Actions → **Remove mod** → Run workflow → enter the slug → run.

---

## I want to UPDATE mods

Actions → **Update mods** → Run workflow:
- Leave the slug blank to update **everything** to latest versions
- Or enter a specific slug to update just that one

The workflow checks Modrinth/CurseForge for newer versions and bumps the URL+hash automatically.

---

## I want to drop in a CONFIG file

You have two folders, and the difference matters:

### `config/` — always overwritten on install

Use this for configs you want to **enforce** for every user. Every time someone installs the pack (or the auto-updater runs), files here will overwrite whatever they had.

**Example:** You spent two hours tuning JEI's bookmark sorting. Drop your `jei-client.ini` in `config/jei-client.ini`. Every user gets it.

### `defaultconfigs/` — only seeded on first install

Use this for configs that you want to **suggest** but not enforce. The user gets your version on first install, and any changes they make stick around — future updates won't clobber them.

**Example:** You set `xaerominimap.txt` defaults to a friendly zoom level, but if a user prefers something else, you don't want to keep stomping on their preference. Drop it in `defaultconfigs/`.

### Which do I use?

Rule of thumb: **`config/` for things that affect gameplay or balance** (Create speed multipliers, mob spawn rates, KubeJS scripts), **`defaultconfigs/` for things that are personal preference** (minimap zoom, JEI search filter style).

When in doubt, use `defaultconfigs/`. It's friendlier to users.

### Where do I get the config files in the first place?

Run the pack once, configure the mod in-game, then go grab the file from `.minecraft/config/<modname>.toml` (or whatever the file is). Drop it in the repo using one of:

- **GitHub web UI**: navigate to `config/` or `defaultconfigs/` → "Add file" → "Upload files" → drag the file in → commit.
- **Local checkout** with `packwiz` (see Local setup below).

After uploading via the web UI, trigger the **Update mods** workflow (with no slug) — it does a `packwiz refresh` which is needed to register the new file in `index.toml`.

---

## I want to add a RESOURCEPACK or SHADERPACK

Same flow as configs:

1. In the GitHub UI, navigate to `resourcepacks/` or `shaderpacks/` → "Add file" → "Upload files" → drag the `.zip` in → commit.
2. Trigger **Update mods** with no slug to refresh the index.

---

## I want to add a KUBEJS script (custom recipe, etc.)

KubeJS lets you write JavaScript files to add/remove/change recipes. Drop your `.js` file in:

- `kubejs/server_scripts/` — for recipes, loot, server-side stuff
- `kubejs/client_scripts/` — for tooltips, JEI hide/show, item renaming
- `kubejs/startup_scripts/` — for things that run at game launch (rare)

If you're just changing recipes, put it in `server_scripts/`.

You can create the file directly in the GitHub web UI (navigate to the folder → "Add file" → "Create new file"). Then trigger **Update mods** to refresh the index.

Example recipe removal:

```js
ServerEvents.recipes(event => {
  event.remove({ output: 'minecraft:stick' })
})
```

See [`../kubejs/README.md`](../kubejs/README.md) for more examples.

---

## I want to publish a release

1. Bump the `version` in `pack.toml` (edit it directly in GitHub's web UI — pencil icon — and commit).
2. Releases tab → "Draft a new release"
3. "Choose a tag" → type `v0.1.1` → "Create new tag on publish"
4. Title: `v0.1.1` (or something descriptive)
5. Description: write what changed, or just hit "Generate release notes"
6. Click **Publish release**

GitHub Actions will then:
- Download every mod jar
- Build the .mrpack, both Prism zips, and the server zip
- Attach all four to the release page

This takes ~5 minutes. If something fails, check the **Actions** tab — the failed step will be red and you can read the error log.

---

## Local setup (only if you need it)

The GitHub workflows handle adding/removing/updating mods, and the web UI handles file uploads. But a few tasks still benefit from a local checkout:

- Bulk-editing many config files at once
- Writing or testing KubeJS scripts (so you can see syntax errors before pushing)
- Running a local test build before tagging a release

If you want to do this:

1. Install [packwiz](https://packwiz.infra.link/installation/): `go install github.com/packwiz/packwiz@latest` (with [Go](https://go.dev/dl/) installed), or `brew install packwiz` on Mac.
2. Clone the repo: `git clone <repo-url>`.
3. From the repo root you can now run any packwiz command. Common ones:
   - `packwiz refresh` — re-hash files after editing configs
   - `packwiz serve` — local HTTP server for testing
   - `./scripts/build-prism-bundled.sh` — build a test zip locally

If you'd rather not install anything: every task above can be done through the GitHub web UI plus the workflows. It's just slightly slower for bulk edits.

---

## Things you should NOT edit by hand

- `mods/*.pw.toml` — these are auto-managed. Use the **Add mod** / **Update mods** workflows.
- `index.toml` — packwiz manages this. The workflows refresh it automatically.

---

## Common gotchas

- **Workflow says "no version found"** — the mod doesn't have a 1.21.1 NeoForge build, or you picked the wrong source (some mods are CF-only or Modrinth-only).
- **Wrong version got picked by Add mod** — the workflow uses `--yes` which auto-accepts. If a mod has multiple candidates and the workflow picked badly, remove it and add it from a local checkout (no `--yes`) so you can pick.
- **"hash mismatch" on the build workflow** — usually means a file in `config/`, `defaultconfigs/`, `kubejs/` etc. was edited without a refresh. Trigger the **Update mods** workflow with no slug — it does a refresh as part of its run.
- **Config changes don't apply for existing users** — if you put it in `defaultconfigs/`, it only applies on first install. If they already installed, they keep their old version. Use `config/` if you want to force-update.

---

## Quick FAQ

**Q: Where do I put a jar I downloaded?**
A: You don't. We don't ship jars in the repo — only URLs. Use the **Add mod** workflow to register a mod from Modrinth or CurseForge. The build pipeline downloads jars at release time.

**Q: What if a mod is *only* on some random website, not Modrinth or CurseForge?**
A: Rare case, currently needs a local packwiz install: `packwiz url add <n> <url>`. Tell me if this comes up and I can add a workflow for it too.

**Q: I broke something. How do I undo?**
A: Every workflow run is a commit. Repo → Commits → find the bad one → "Revert" button. GitHub will open a PR to undo it.

**Q: I want to test the pack without making a release.**
A: Trigger the **Build modpack** workflow manually (Actions tab → Build modpack → Run workflow). It produces all four artifacts in the run's "Artifacts" section, downloadable for 90 days. Inspect the bundled zip to see if mods are there before cutting a real release.
