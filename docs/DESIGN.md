# Derpack X — Design

This document captures the **why** behind the architectural decisions in Derpack X. The README and per-folder docs cover the **what**: how to install, how to add a mod, where things live. This is for when someone is staring at something and wondering why it's structured that way.

If a decision feels weird, look here first. If it isn't here, it should be — propose an addition.

---

## What is this thing

Derpack X is a Minecraft 1.21.1 / NeoForge modpack distributed via [packwiz](https://packwiz.infra.link/) manifests. The repo doesn't contain mod jars; it contains URLs and hashes. CI builds two release artifacts (a `.mrpack` and a Prism installer zip), both small. End-users' launchers fetch jars from each mod's authoritative source at install time. A custom Go-based local editor handles routine pack maintenance through a web UI; GitHub Actions workflows are a fallback.

---

## The design spine — breadth, with a Create economy running through it

Derpack X is, on purpose, a **kitchen-sink survival pack**: a large, varied content set — Farmer's Delight and a dozen food addons, Yung's structures, MineColonies, the Ars / Iron's Spells / Occultism magic stack, ~16 mob mods, Macaw's decoration, Create Aeronautics vehicles, and more. Breadth is a feature, not an accident — the friend group wants a world with a lot to do.

Breadth's failure mode is **incoherence**: a heap of mods that never talk to each other, each its own self-contained progression. We avoid that with a single **spine** running through everything — a **Create-driven scarcity economy**. Ores are scarce and regional (large veins in specific biomes; vanilla thinned to a trickle), and the core loop is **mine → process through Create → manufacture goods → sell for currency (Numismatics)**.

The governing rule: **the spine should penetrate every branch as far as is feasible.** Wherever a mod offers progression, power, or resources, route it *through* Create and the economy rather than around them:

- **Magic and adventure gear** is gated behind Create-made components, so those branches plug into the factory instead of sidestepping it.
- **"Free resource" shortcuts** — kitchen-sink factory blocks, colonist auto-mining, cheap backpacks, auto-cooking — are gated or tuned so they don't quietly undo scarcity.
- **Structure / dungeon loot** is pulled back from free end-game gear, and may seed currency instead.
- **MineColonies** stays as a village/social layer, but its auto-production is gated so it *feeds* the Create loop rather than replacing it.

So the lens for every content decision is one question: **does this plug into the Create spine, or bypass it?** We don't cut content to get focus — we integrate it. Bypasses are the thing to fix, not the mods themselves. The detailed, per-branch integration (recipe changes tracked in #17, loot changes in #18) happens phase by phase as we reach each cluster; this section is the *why* behind all of it.

---

## 1. Why packwiz

Why this isn't just a folder of jars. A naive modpack distribution model would zip the mod jars and host the zip somewhere. That model has several problems we wanted to avoid:

1. **Licensing.** Most mod licenses don't permit redistribution. Bundling jars puts the maintainer in legal gray territory at best, and burns trust with mod authors at worst.
2. **File size.** A zip with ~200 mod jars is 500+ MB. Hosting and downloading that for every release is wasteful when the actual *change* between versions is usually a handful of files.
3. **Updates.** A zip is a snapshot. To stay current, the user has to re-download and re-import. That's enough friction that most users just don't update.
4. **Collaboration.** Multiple maintainers pushing big binary blobs into git is painful. Diffs are unreadable, merge conflicts are a nightmare, and the repo's history bloats fast.
5. **Version tracking.** Plain-text manifests have meaningful diffs. You can see at a glance which mod was bumped, which version was pinned, when something was added or removed — both in PR review and in `git log`.

packwiz addresses all of these. The repo holds tiny `.pw.toml` manifests — one per mod — describing where to fetch the jar (Modrinth or CurseForge) and what its hash should be. Releases ship those manifests as a small `.mrpack` (~100 KB) or a Prism installer zip (~5 MB). The user's launcher fetches the actual jars from each mod's authoritative source at install time, with hash verification.

The trade-off: every install needs internet on first launch, and ongoing operation depends on Modrinth/CurseForge staying reachable. Acceptable — anyone playing modded MC is already online.

---

## 2. Why version-named branches (provisionally)

For now, branches are named after the next release version. If `main` is at `0.3.2`, work-in-progress goes on a branch called `0.3.3` for patches, `0.4.0` for bigger work. The branch name and the eventual release tag are the same string.

**This convention is provisional.** The team is small (one maintainer + one regular collaborator), the workflow is still settling, and the "right" branching strategy depends on patterns we haven't seen yet. Expect this to change.

The convention works *for now* because:

- We can usually answer "what's the next version?" at any moment, so the branch name is obvious.
- There's typically only one branch in flight at a time.
- Branch name = release tag means no translation needed when cutting a release.

Alternatives we'll probably revisit when the workflow matures:

- **Personal branches** (`zagwar/whatever`, `xela/cool-stuff`). Likely the right move once parallel work becomes common — each contributor gets their own lane and merge ordering becomes about review, not version sequencing. Probably the next step if the version-named branch starts feeling like a collision point.
- **Feature branches** (`add-create-cafe`, `fix-quark-crash`). The traditional one-feature-per-branch model. Worth considering if we ever want to track work by *what* it is, not *when* it ships.
- **Continuous delivery** — drop the concept of versioned releases entirely; the installer-style instance always pulls from `main`. Genuinely interesting since the installer-style distribution already auto-syncs. Would let us ship the moment a PR merges, no version dance. Big workflow shift, not on the near-term radar but worth keeping in mind.

No version-naming dogma. Pick the convention that matches how the team actually works at any given moment, and rewrite this section when you do.

---

## 3. Why the editor exists

Why a custom editor instead of just using packwiz CLI or the GitHub Actions workflows. Both alternatives existed before the editor and worked fine for a single technical maintainer. The editor came out of a specific collaboration friction.

The original workflow was: maintainer runs `packwiz` commands locally; collaborator triggers GitHub Actions workflows in the browser ("Add mod" with a slug field, "Update mods," etc.). That worked, sort of. The collaborator was actually editing locally with raw jars and never touching the repo's workflows at all — the friction of "click button, fill form, wait 30 seconds, refresh page" was high enough that he just downloaded jars manually and told the maintainer what to add. The repo and the actually-played pack drifted apart constantly.

The editor exists to close that gap. It's a single .exe that runs from the repo, opens a local web UI, and exposes packwiz operations as buttons. The collaborator can edit the pack the same way the maintainer does, branch and PR through GitHub Desktop, never touch git or YAML or a terminal. The maintainer benefits too — many edits are faster through the UI than through CLI flags or workflow forms.

Three principles drive the design:

1. **The editor doesn't do git.** GitHub Desktop already handles clone/branch/commit/push/PR with a polished UI. Reimplementing any of that would be wasted effort and a worse experience. The editor only edits files in the working directory; git records those changes; GitHub Desktop helps push them.
2. **The editor doesn't replace anything.** The packwiz CLI still works. The GitHub Actions workflows still work. The editor is *an additional path*, not a forced one. If it breaks, the underlying mechanisms are unchanged.
3. **The editor is a thin wrapper.** It shells out to packwiz for operations packwiz already supports (add, remove, pin, update, refresh). It only reimplements operations that needed new logic — Modrinth API calls for the version picker, hash recomputation for pinned mods, Prism instance bundling for local testing.

Net effect: the collaborator's friction dropped from "I'll just edit my own folder and tell you" to "I added these mods on branch 0.4.0, PR open." That's the win that justified building the thing.

---

## 4. Why we don't redistribute jars

The build only ships manifests, not jars. Releases attach two artifacts: a `.mrpack` (manifests bundled in Modrinth's format) and a Prism installer zip (manifests + a launcher script). Both are small (~5 MB or less). The user's launcher fetches actual mod jars from each mod's authoritative source — Modrinth or CurseForge — at install time.

Earlier in the project we also shipped a "bundled" Prism zip with all jars baked in. It worked, but we dropped it for two reasons:

1. **Licensing.** Most mod licenses don't permit re-hosting jars. Some explicitly forbid it. By distributing a zip with someone else's jar in it, the maintainer is the one violating that license, even though they're not making money or claiming authorship. It's the kind of thing that quietly damages the relationship between modpack maintainers and mod authors, and at scale it makes mod authors more restrictive about future licensing.
2. **GitHub release asset cap.** GitHub releases have a 2 GiB per-asset limit. The bundled zip was approaching that ceiling as the pack grew. The server-pack variant blew past it. Rather than build infrastructure to host large artifacts elsewhere (S3, R2, a self-hosted box), it was cleaner to drop the bundled artifact entirely.

The trade-off: every install needs internet on first launch. Acceptable — anyone playing modded MC is online anyway, and Modrinth/CurseForge are reliable enough that "fetch on install" is a non-event.

The editor's "Build & Launch in Prism" feature still does jar bundling, but it's local-only — never published. That's how the maintainer tests pack changes before pushing. The result lives on their machine and is never uploaded anywhere.

---

## 5. Why no server pack right now

There is no server pack artifact, and there isn't going to be one. When we first built the release pipeline, a server pack was the fourth release artifact. It got dropped because:

1. **The 2 GiB GitHub release asset cap.** With ~200 mods, the bundled server zip pushed past 2 GiB. GitHub returned `size must be less than 2147483648` and refused the upload. Workable with splits, but ugly.
2. **No actual users.** No one is hosting a server. Building an artifact nobody downloads is wasted CI time.
3. **A server doesn't need a special build anyway.** packwiz-installer-bootstrap works server-side just like it does client-side. A future server admin can run it directly against the repo's `pack.toml` URL — same mechanism the client installer zip uses. The CI doesn't need to produce anything special; the server pulls what it needs.

The build step has been removed from `build.yml`. If we ever want to host a server, the procedure will be:

1. SSH to the server box
2. Install Java 21 and grab `packwiz-installer-bootstrap.jar`
3. Run it against `https://raw.githubusercontent.com/<repo>/main/pack.toml` with `-s server`
4. Configure the launch script and start it up

No CI involvement. The server stays in sync the same way clients do — every restart pulls latest from `main`, exactly like the installer-style client behavior.

---

## 6. We had a build cache, then removed it

Earlier iterations of the build pipeline had a content-keyed cache in `build.yml`. It hashed the build inputs (manifests, configs, scripts), looked up that hash in GitHub Actions cache, and reused artifacts on hits. This was meaningful when the build pipeline included a bundled-jars zip — fetching ~500 MB of jars takes minutes, and skipping it when nothing material changed was a real win.

Once we stopped distributing bundled jars (see "Why we don't redistribute jars"), the build dropped to ~30 seconds end to end. Caching saved at most 25 seconds while adding two cache steps, a fingerprint computation, and a "version field stripped from `pack.toml`" hack to keep test/release builds compatible. The complexity outweighed the savings.

The cache logic was removed. If pack distribution ever shifts back to a model where builds are slow (server-side jar fetching for some artifact, generated content from KubeJS, etc.), reintroduce caching. The version-field-stripping pattern is the non-obvious bit worth remembering: hash everything *except* the version, so a workflow_dispatch test build at `0.3.3-dev` can be reused for the `0.3.3` release tag.

---

## 7. Why batch-add forces side=both

The Wishlist (batch-add) feature in the editor always sets newly-added mods to `side = "both"`, even if Modrinth says otherwise. The single Add Mod path lets the user pick a side; batch-add doesn't ask, it just forces both.

The reason: Modrinth's `client_side` / `server_side` metadata is unreliable in practice. Many mods have it set wrong. The mod author copied a template, didn't think hard about the field, and now `client_side: required, server_side: unsupported` ends up in the metadata even though the mod absolutely requires the server too. packwiz reads this faithfully and writes `side = "client"` to the manifest. The mod then quietly fails to register on dedicated servers, and the maintainer spends an evening figuring out why.

When you're adding mods one at a time, you're paying attention. You see the side selector, you pick "both" if that's right. When you're batch-adding 20 mods from a wishlist, you're not paying that level of attention to each one. The most common case for batch-adding is "I want all these content mods" and the answer for those is *almost always* both. So we default to that.

The escape hatch: any mod's side can be flipped from the table afterward via the side dropdown on its row. If a mod really is client-only (a minimap, a tooltip mod, a settings menu mod), one click puts it back. The flow is "default to both, fix the few exceptions" rather than "trust upstream metadata, fix the many exceptions."

A more conservative design would inspect each mod's metadata, default to the upstream value, and require explicit override. We chose the opinionated default because the failure mode of "wrong side, mod silently doesn't load" is much worse than "right side, but I had to fix three of them."

---

## 8. Why the editor doesn't do git

The editor only edits files in the working directory. It doesn't touch git, doesn't commit, doesn't push, doesn't open PRs. All of that is GitHub Desktop's job.

The temptation to add git operations was real. The editor already knows when an operation modifies the manifests; adding a "Commit and Push" button would be a few lines of code. But every time we considered it, the same answer came up: GitHub Desktop already does this better than we could.

What GitHub Desktop gives you for free:

- A polished diff viewer with syntax highlighting
- Visual conflict resolution
- Branch management (create, switch, delete, fetch, pull, push)
- PR creation that integrates with GitHub's review UI
- Authentication via stored OAuth tokens, no `git config` setup
- History browser with revert
- Stashing for in-progress work

Reimplementing any of that in our editor would mean shipping a worse version of features the user already has. And duplicating those features creates two states to keep in sync — what the editor thinks is committed vs. what GitHub Desktop thinks. Diverging UIs lead to bugs where the user trusts one and the other is the source of truth.

The contract is sharper for being narrow: the editor edits files; git records changes; GitHub Desktop helps push them. Each tool has one job. If a contributor breaks something with the editor, GitHub Desktop's diff view shows it in the same UI they'd use for any other change, and recovery is the same as recovering any other bad commit.

One concession: the editor *displays* git information (current branch, files changed since last commit) when those would help the user decide what to do. It just never *modifies* git state.

---

## 9. Operations baked into the binary

The editor's operations — Add, Remove, Pin, Set Version, Build & Launch, etc. — are hardcoded in the binary. A v2 alternative would let the repo declare custom operations via something like `.editor/operations.toml`, dynamically rendered as buttons.

The honest reason v1 is baked-in: it was simpler to ship. We had a working set of ~10 operations that covered everything we needed; deferring "make this configurable" let us ship the whole editor in a few sessions instead of designing a config schema first.

A note from session 1 of editor design called this out as a "v2 idea" and it has stayed there. The trigger to actually build v2 would be wanting pack-specific automation that doesn't belong in the binary — something like a custom validator, a per-pack tweaks profile, or pack-specific deployment hooks. None of that has come up yet.

If v2 happens, the migration would be: define an `operations.toml` schema (operation name, button label, what command to invoke, what arguments to collect from the user), implement a loader that reads it on editor startup, render dynamic buttons. The existing baked-in operations would become the defaults when no `operations.toml` exists.

Until then, baking operations into the binary is a "good enough" choice rather than the right one. CI rebuilds the binary on every push to `tools/editor-src/`, so adding a button is a code change but not a heavy one.

---

## When updating this doc

If you're making a decision that future-you might forget the rationale for, add it here. If you're undoing a decision that's documented here, update or remove the section — don't leave stale "why we did X" notes after we've stopped doing X.

If a decision is clearly mundane ("we used Go because we know Go") it doesn't need a section. If it's contentious, novel, or a trap we already fell into once, it does.
