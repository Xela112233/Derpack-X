# Issue Triage Playbook

Quick process for an issue-triage session — re-label, **milestone**, re-assign, close, and sequence.
Written from the post-0.5.2 triage; **update it as the project moves** (issue numbers below will drift).

## 0 · Get ground truth first — the ROADMAP can lag `main`
The live tracker is truth; `docs/ROADMAP.md` can lag `main`. Always:
1. `git log --oneline -40 | grep -iE "merge pull request|bump version"` — what actually shipped + the current version.
2. `pack.toml` → current `version`.
3. List **open PRs** (only those are unmerged work) and **open issues** (`state OPEN`, `perPage 100`).
4. Reconcile: an "open" issue may already be resolved by a merged PR — flag those as **close candidates**.

## 1 · Parse the issue list (too big for context)
`mcp__github__list_issues` returns 100 KB+. Offload to a subagent: slice the saved tool-result file in
~80 KB spans (`open(path).read()[A:B]` via python) and return a table of number / title / labels /
milestone / body-summary / cross-refs. Note: the list API **omits `assignees`** — read those per-issue
(`issue_read`) only if needed.

## 2 · Re-label — work within the fixed vocabulary (don't invent labels)
16 labels in use: `add-mod` · `bug` · `content` · `discussion` · `enhancement` · `gameplay` ·
`low-priority` · `needed-for-release` · `performance` · `playtest-blocked` · `qol` · `remove-mod` ·
`site` · `tooling-ci` · `update-mod` · `worldgen`. (Confirm against the live set; it can grow.)
- Every issue gets ≥1 **area** label and, where apt, a **type**.
- `needed-for-release` = the `v1.0.0` gate set **only** (§6). Keep it tight.
- `playtest-blocked` = can't finish without an in-game check. Also stays on a **closed** issue whose fix
  has landed but isn't verified in-game yet (see §5): `is:closed label:playtest-blocked` in a milestone =
  the "fixed but unverified" count. Drop the label once playtest passes.
- `issue_write` with `labels` **replaces** the whole set — include the labels you're keeping.

## 2a · Effort estimates (`size:*` labels)
Each issue also carries a **`size:*`** effort label — the one sanctioned addition to the 16-label
vocabulary above. The estimate is **solo hands-on time**, decision assumed already settled, and it
**excludes** boot/settle/playtest waits (flag those separately with `playtest-blocked`).

| Size | Band | Typical |
|---|---|---|
| `size:XS` | ≤30 min | one config/manifest edit, a pin, a doc line, a quick decision |
| `size:S` | 30–60 min | a small recipe/KubeJS/config change, a single-mod keep/cut, a settings toggle |
| `size:M` | 1–3 hr | multi-file content, an investigation + change, a tooling script |
| `size:L` | 3–6 hr / half-day | a pillar build-out chunk, multi-mod balance, worldgen + playtest, a new-mod skeleton, a hard bug |
| `size:XL` | multi-day — **decompose** | pillar umbrellas / recurring frameworks; split into sub-issues |

`size:XL` umbrellas **double-count** their sub-issues — don't sum them into a project total. The
mechanical "size every open issue" pass is the `/estimate-issues` command.

## 3 · Milestone every issue (the road to 1.0)
The backlog is organized into GitHub **milestones**, not loose buckets, and **every issue — including
every new one — must have one** (an unmilestoned issue is untriaged). Pick the milestone whose **goal**
it serves; ongoing trackers and not-yet-scheduled work go to **`📋 Backlog / living`**. The map — each
milestone's goal / exit-criterion and the current ordered set (Foundation → the pillar versions →
**curation passes** → `v1.0.0 Release [NFR]` → `📋 Backlog / living`) — lives in `docs/ROADMAP.md` →
"Milestones — the road to 1.0"; **treat that as the source of truth** (milestones get added and reordered
on GitHub, so don't hardcode the set here). Set the milestone via `issue_write` (`milestone` = the
milestone number). Front-load small/foundational work; the maintainer prefers **not**
leading with release blockers.

## 4 · Re-assign — two people, clear lanes
- **Xela112233** (you, admin) own implementation **and everything that runs on the live server** — you
  run the box, so **server tests / playtests / RAM tuning are yours** (e.g. #25 claims, #48 GC/RAM,
  #98 pre-gen, #43 flight, the #58 ore-gen *playtest*). zagwar typically hands server tests to you.
- **zagwar** (write; co-maintainer) owns: his living trackers (#17 recipes, #18 loot, #21 removals) and
  playtest **verdicts / design calls** (#60 umapyoi keep/cut, #58 rarity/regionality *review*). Joint
  **decisions** (#78 rename, #1 shaders) → assign **both**.
- Assign the **active milestone's** work to its owner; leave the long tail **unassigned** as backlog —
  don't self-assign all 50+ (noise; self-actions don't notify anyway).

## 5 · Close — on stage, so the milestone bar ticks up
- **Close an issue when its fix lands on the version branch** (e.g. committed to `v0.6.0`), not when the
  branch later merges to `main`. Closing on stage is what makes the milestone % climb *as work happens*,
  which is the at-a-glance "how much of this version is built?" gauge — close-on-merge would leave it
  flat-then-spike and re-surface settled work in triage.
- **Tag the closed issue `playtest-blocked` until it's verified in-game** (§2). The close means "fix
  implemented"; the label means "not yet confirmed." Drop the label once playtest passes. A failed
  playtest just **reopens** the issue — normal.
- **"Version done" ≠ milestone at 100%.** Milestone 100% = all work *implemented*. The hard ship signal is
  the **`vX.Y.Z` release tag** (`release.yml`, §6). Truly shipped = milestone 100% **and** no
  `playtest-blocked` left **and** the tag cut.
- A close still gets a one-line comment pointing to live follow-ups where apt (e.g. #116/#58 → #318),
  and the PR's `Closes #N` keyword is fine to leave in (a no-op once the issue is already closed).
- **Living trackers stay open by design:** #9 (mod ideas), #17 (recipes), #18 (loot), #21 (removals).
- Confirm with the maintainer before closing anything **ambiguous** — a partial fix, a multi-issue PR, or
  a design-call close. Routine "fix landed on the branch" closes don't need a prompt.

## 6 · Release gate — `v1.0.0` / `needed-for-release`
The `needed-for-release`-tagged set is the `v1.0.0 — Release` milestone, deliberately last: rename (#78),
CI-required (#79), ore-gen tuned + playtested (#318, with #81 as a sub-issue), RAM tuned (#48). Tag flow: bump `version`
in `pack.toml` → merge to `main` → tag `vX.Y.Z`.

## 7 · The systems-loop yardstick (why an issue matters)
Every mod/feature anchors to ≥1 system — **Create · magic** (production) · **survival** (pressure) ·
**economy** (distribution & trade, incl. logistics/aeronautics) — the systems being one causal loop
(`docs/SYSTEMS.md`; the old "five systems" framing is retired), which is also why the pillar milestones
are shaped the way they are. Target: **2**. A one-system mod is a candidate for a *second* weave; a
no-system mod is a **cut** candidate (curation, `docs/CURATION.md`). See `CLAUDE.md` + `docs/DESIGN.md` +
`docs/SYSTEMS.md`.

## 8 · After triage — update the docs
Refresh `docs/ROADMAP.md` (milestones + needed-for-release + recently-merged) so the next session
starts from truth.

## Workflow guardrails (from `CLAUDE.md`)
- **Branch:** rename the harness's random name to `claude/<topic>` before pushing; never commit to `main`.
- **PRs:** don't auto-open; fold related work into the branch you're on. Gameplay/content PRs stay
  **draft + carry a `## Playtest` checklist**; docs-only PRs don't. After opening, `subscribe_pr_activity`.
- **Gotchas:** `galosphere:silver_ingot` is palladium (never unify as silver); don't re-add EMI; Sodium
  runs only via the pinned Sable-compat bridge — never bare, never Embeddium (would break Aeronautics;
  see CLAUDE.md / #98); itemphysic + FoamFix are client-only; Aeronautics breaks under shaders.
