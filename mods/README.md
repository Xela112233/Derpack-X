# mods/

Each `.pw.toml` file is a packwiz manifest for one mod. **Don't edit them by hand for normal stuff** — use the workflows (Add mod, Remove mod, Update mods). The workflows compute hashes and update the index correctly.

But sometimes the workflows aren't enough — you need a beta version, an older version, or a specific platform. That's when you edit these files directly.

## What a manifest looks like

```toml
name = "Some Mod"
filename = "somemod-1.21.1-1.2.3.jar"
side = "both"
pin = true                              # optional, see below

[download]
url = "https://cdn.modrinth.com/data/PROJECTID/versions/VERSIONID/somemod-1.21.1-1.2.3.jar"
hash-format = "sha512"
hash = "abc123..."

[update]
[update.modrinth]
mod-id = "PROJECTID"
version = "VERSIONID"
```

CurseForge mods use `[update.curseforge]` with `project-id` and `file-id` (numeric) instead of `mod-id` / `version`.

## The pin field — read this before editing manually

`pin = true` tells packwiz "leave this manifest alone during `update --all`". This is what you want when you've manually picked a non-latest version (older, beta, specific build) and don't want it auto-bumped.

**Critical:** `pin` MUST go at the top level of the file, NOT inside `[update]`. TOML treats everything after a `[section]` header as part of that section, so this:

```toml
[update]
pin = true                # WRONG — packwiz ignores this
```

won't work. This does:

```toml
name = "Some Mod"
side = "both"
pin = true                # right — top level, before [download]
[download]
...
```

When `pin = true` is set, `packwiz update <slug>` will skip this mod, AND `packwiz update --all` will skip it. So the workflows can't fix anything for a pinned mod — you have to edit the manifest by hand AND compute the hash yourself.

## Common edits

### Downgrade to an older Modrinth version

1. Go to the Modrinth page → Versions tab → find the version you want
2. **Right-click the green Download button** → Copy Link. URL looks like `https://cdn.modrinth.com/data/PROJECTID/versions/VERSIONID/filename.jar`
3. The 8-character random string between `versions/` and the filename is the version ID
4. Edit four things in the manifest:
   - `filename` (the jar name from the URL)
   - `url` under `[download]` (the full URL)
   - `version` under `[update.modrinth]` (the VERSIONID — must match the URL)
   - `hash` — see "Computing the hash" below
5. Add `pin = true` at the top level so future `update --all` doesn't undo this
6. Commit

### Switch a mod from Modrinth to CurseForge (or vice versa)

The `[update.modrinth]` and `[update.curseforge]` sections have different fields, so editing in place is brittle. Easier:

1. **Remove mod** workflow with the current slug
2. **Add mod** workflow specifying the other source

For betas only on CF, after Add mod runs, edit the manifest to pin to the beta's file-id (see below) and recompute the hash.

### Pin a CurseForge beta or specific file

CF manifests look like:

```toml
[update.curseforge]
file-id = 5247113
project-id = 392898
```

To switch to a specific (e.g. beta) file:

1. CF page → Files tab → All Files → click the file you want
2. URL ends in `/files/<NUMBER>` — that number is the file-id
3. Edit `file-id` in the manifest
4. Update `filename` and `url` to match the new file (download URL is on the same CF page)
5. Recompute the hash (see below)
6. Add `pin = true` at the top level

### Change client-only / server-only / both

The `side` field controls which artifacts the mod gets bundled into:

- `side = "both"` — bundled into client AND server packs
- `side = "client"` — only in client artifacts (Prism, mrpack)
- `side = "server"` — only in server artifact

Most mods are `both`. Stuff like Xaero's Minimap, IPN, EMI/JEI should be `client`. Server-side performance mods or admin tools would be `server`.

Just change the field and commit. Build workflow respects this on the next run.

## Computing the hash

When you edit a manifest by hand to point at a different file, the existing `hash` value is wrong (it was computed from the previous file). Three ways to fix this, in order of convenience:

### Option 1: Compute hash workflow (recommended)

Run the **Compute hash** workflow with the mod's slug as input. The workflow downloads the file from the URL in the manifest, computes the SHA-512 hash, writes it back into the manifest, and commits. Works for pinned mods (unlike Update mods).

Use this whenever you edit a manifest's URL by hand and need the hash refreshed.

### Option 2: Compute it yourself

If you don't want to wait for a workflow run, the hash is a SHA-512 of the actual jar contents.

PowerShell on Windows:

```powershell
$url = "<the new download URL>"
Invoke-WebRequest $url -OutFile $env:TEMP\mod.jar
(Get-FileHash -Algorithm SHA512 $env:TEMP\mod.jar).Hash.ToLower()
```

Bash on Linux/WSL:

```bash
curl -sL "<URL>" | sha512sum | awk '{print $1}'
```

Paste the resulting lowercase hex string into the `hash =` field (in quotes), commit.

### Option 3: Let Update mods compute it (un-pinned mods only)

1. Set `hash = ""` (empty string)
2. Make sure `pin = true` is NOT set
3. Commit
4. Run **Update mods** workflow with this mod's slug

packwiz will download the file, compute the hash, and write it back. Then if you want it pinned, edit the file again to add `pin = true`.

This works but takes more steps than Option 1. Mostly useful if you're un-pinning anyway.

## After editing a manifest

If the URL changed, run **Compute hash** with the mod's slug to refresh the hash and commit. This works for both pinned and un-pinned mods.

If you only changed `side` or other non-URL fields, just commit — no hash recomputation needed. The build workflow runs `packwiz refresh` itself, which updates the index pointer to your manifest.

## Don't edit

- `index.toml` in the repo root — it's auto-managed by `packwiz refresh` (the build workflow runs this)
- The hash field, unless you computed it yourself or set it to `""` to be regenerated by Update mods
- Filenames that don't match what the URL actually serves (will 404 at install time)

## Reference

- packwiz manifest format: https://packwiz.infra.link/reference/pack-format/mod-toml/
- packwiz pin command: https://packwiz.infra.link/reference/commands/packwiz/pin/
- packwiz update command: https://packwiz.infra.link/reference/commands/packwiz/update/
