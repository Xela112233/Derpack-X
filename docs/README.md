# `docs/` — index

The durable documentation for Project Commonwealth. This folder is **reference**, not scratch: every file here is
meant to stay useful. Ephemeral session/handoff prompts don't belong in `docs/` (see the convention at
the bottom). `docs/` is fully `.packwizignore`'d, so nothing here touches the packwiz index or CI.

Start at `../CLAUDE.md` (the north star + guardrails) and `DESIGN.md` (the why). Everything else hangs
off those two.

## Canonical reference — "this is how it is"

| Doc | Purpose |
|---|---|
| [`DESIGN.md`](DESIGN.md) | **The why** behind the architecture — the systems-loop rationale and all the architectural decisions. If a decision feels weird, look here first. |
| [`SYSTEMS.md`](SYSTEMS.md) | **The canonical structural model** — the systems as one causal loop (scarcity → pressure → production → economy; *Eco*-inspired). How the parts interlock: production (Create/magic), the economy (trade + logistics/aeronautics), pressure (survival), the recipe/gating spine, the ecological cost. The old flat "five systems" framing is retired here. |
| [`CURATION.md`](CURATION.md) | **The curation rubric** — the charter every keep/cut pass scores mods against (the 8 gates). |
| [`RECIPES.md`](RECIPES.md) | The "Made through Create" recipe convention **and the live triage ledger** — the most-referenced doc; the source of truth for recipe changes. |
| [`CONNECTIVITY.md`](CONNECTIVITY.md) | The "one web, not many clusters" connectivity metric — turns the north-star web goal into a number you can re-check (`tools/connectivity/`). |
| [`ORE-GENERATION.md`](ORE-GENERATION.md) | The GTMOGS mix-vein ore-gen model (GTNH-style regional veins) and its config. |
| [`GUN-PACKS.md`](GUN-PACKS.md) | The gun system — TaCZ wired into Create via two add-ons. |
| [`CREATE-SPINE.md`](CREATE-SPINE.md) | **The ratified Create-spine design** — the tier ladder, the cost model + tier tokens (Part 2b-ii), boss gates, the aeronautics control thread, and the ratification record. |
| [`SERVER.md`](SERVER.md) | Client *and* server from one `pack.toml` — delivery paths and the `side`-metadata rules. |
| [`SPAWN-GATING.md`](SPAWN-GATING.md) | The spawn-control / mob-curation work from the #83 pass — what each file does and what needs in-game verification. |

## Setup & how-to

| Doc | Purpose |
|---|---|
| [`PRISM-SETUP.md`](PRISM-SETUP.md) | Player-facing install walkthrough (Prism Launcher). |
| [`EDITING.md`](EDITING.md) | How to add mods / change configs / ship updates — the editor flow plus the packwiz CLI fallback. |
| [`CI-CHECKS.md`](CI-CHECKS.md) | What the merge-gating `pr-checks.yml` validates (and how to make the checks required). |

## Planning & process — living docs

| Doc | Purpose |
|---|---|
| [`ROADMAP.md`](ROADMAP.md) | **The plan** — multi-phase build-out and the milestone map + the odd/even cadence (the road to 1.0). The sequencing source of truth; can lag `main`. |
| [`TRIAGE.md`](TRIAGE.md) | Issue-triage playbook (re-label / milestone / assign / close / sequence). |
| [`CREATE-SPINE-IMPL.md`](CREATE-SPINE-IMPL.md) | The spine **build map** — file-by-file implementation record, the gate table, staged-work status. Pairs with `CREATE-SPINE.md` (design). |
| [`POWER-MODS-REVIEW.md`](POWER-MODS-REVIEW.md) | The power-mod landscape review — decided for v0.7.0; kept as the basis for the 2.0 Power Grid overhaul (#282). |
| [`MOD-REVIEW.md`](MOD-REVIEW.md) / [`WEAVE-REVIEW.md`](WEAVE-REVIEW.md) | Generated review ledgers (`tools/weave-ledger/`) — the pre-1.0 sweep and weave-status snapshots. Regenerate, don't hand-edit. |
| [`PERFORMANCE.md`](PERFORMANCE.md) | Performance ledger — measured findings (spark) and the perf rules of thumb. |
| [`CUSTOM-MODS.md`](CUSTOM-MODS.md) | The custom-mods program (pcmc-arcana etc.) — what we build ourselves and why. |
| [`GOVERNANCE.md`](GOVERNANCE.md) | Player-governance scoping (#260) — municipalities→empires, cascading laws, entity minting. **SCOPING, not accepted** — survey + 3 implementation paths. (Governance reopened in `SYSTEMS.md`.) |
| [`GOVERNANCE-MOD-SPEC.md`](GOVERNANCE-MOD-SPEC.md) | Path B technical spec for #260 — the custom NeoForge governance mod (data model, law-precedence resolver, Numismatics minting bridge, claim-mod integration, MVP phasing). **DRAFT, not accepted.** |
| [`GOVERNANCE-REALMS-SCOPE.md`](GOVERNANCE-REALMS-SCOPE.md) | Governance Part 2 — `pcmc-realms` scope: the hard/soft law taxonomy, the Part-1 API contract, and the **MineColonies-guard soft-enforcement mechanism + addon-vs-fork decision aid** (§5). **SCOPE, not accepted.** |
| [`CONQUEST-MOD-SCOPE.md`](CONQUEST-MOD-SCOPE.md) | Colony Factions & Warfare mod (`pcmc-conquest`, working name) — the VC-replacement faction/army mod scope: NPC joinable factions, faction equipment, army mechanics, and the `pcmc-realms` law→guard seam. **CONCEPT, nothing built.** |
| [`CONQUEST-FORK-VS-ADDON-HANDOFF.md`](CONQUEST-FORK-VS-ADDON-HANDOFF.md) | Handoff brief — decide (per feature) whether `pcmc-conquest` is a MineColonies **addon** or a **fork**, against the `GOVERNANCE-REALMS-SCOPE.md` §5 bar. Investigation tasks + spikes. |

## Player-facing copy

| Doc | Purpose |
|---|---|
| [`../wiki/`](../wiki/) | **The player wiki source** — progression, prospecting, weapons/gear guides (staging for the site; see `wiki/README.md`). |
| [`PATCHNOTES.md`](PATCHNOTES.md) | Hand-maintained changelog. **Kept and maintained** (CLAUDE.md treats it as canonical player-facing copy) — `build.yml` cuts GitHub Releases on tag, but the prose release notes live here. Professional plain prose only. |

## Archive — `archive/`

Historical handoff docs for **merged** PRs. Kept for the conventions and file maps, not for current
priorities. See [`archive/`](archive/):

- `archive/RECIPE-OVERHAUL-HANDOFF.md` — the #62 "Made through Create" overhaul (merged ~v0.4.6). Live convention now lives in `RECIPES.md`.
- `archive/MAGIC-WEB-HANDOFF.md` — the #75 magic-web recipe weave (merged ~v0.5.2). Live design in `RECIPES.md` ("Magic web").
- `archive/MAGIC-WEB-PLAYTEST.md` — the #75 in-game playtest checklist.
- `archive/MODLIST-AUDIT.md` — the modlist-vs-goals snapshot. **Superseded** by `SYSTEMS.md` (model) + `CURATION.md` (rubric) + `ROADMAP.md` (plan); kept for its §A/§B/§C analysis. Archived v0.6.0.
- `archive/RELEASE-CADENCE.md` — the odd/even cadence proposal + extended rationale. **Accepted**; the canonical cadence now lives in `ROADMAP.md` (rationale folded in). Archived v0.6.0.
- `archive/ROTTEN-CREATURES-STRUCTURES.md` — the #106 spawn-allowlist review sheet. Durable rules live in `SPAWN-GATING.md`. Archived v0.6.0.
- `archive/V0.7.0-PHASE-PLAN.md` + `archive/V0.7.0-SPINE-HANDOFF.md` — the v0.7.0 spine phase plan and the org-rename resume doc. **Executed in full**; the living records are `CREATE-SPINE.md` (design) + `CREATE-SPINE-IMPL.md` (build map). Archived v0.7.0.
- `archive/ARCANA-BRIDGE-MOD.md` — scoping/design for the **Project Commonwealth Arcana** code-level magic-weave helper mod (PR #80). The mod moved to its own repo (`theasshats/pcmc-arcana`); the living reference is `CUSTOM-MODS.md`. Archived v0.6.0.

### Still-live docs that were *candidates* (kept, with a note)

- **`BOOT-LOG-BASELINE.md`** — **kept** as a living reference (the accepted boot-log-noise catalog future passes check against). Trim the #119/#120/#121-specific sections now those have closed.

## Conventions

**The canonical split — don't duplicate the systems preamble.** Several docs open with a version of
the systems-loop framing; they are **not** interchangeable:
- `SYSTEMS.md` is **canonical for the structure** (the *how it interlocks* — the loop). The model itself.
- `DESIGN.md` is **canonical for the rationale** (the *why*) — cite it; don't re-argue it elsewhere.
- `ROADMAP.md` is the **plan** (the *when / in what order*) and the **cadence map** (odd/even).
- `CURATION.md` is the **keep/cut charter** (the gates a mod must pass).
- (`archive/MODLIST-AUDIT.md` is the old **snapshot audit** — historical, superseded by the above.)

**`EDITING.md` ↔ `../mods/README.md`.** `EDITING.md` is the **task walkthrough** (add a mod, change a
config, ship an update, with the editor + packwiz CLI fallback); `mods/README.md` is the **manifest reference**
(the `.pw.toml` format, the `pin`-placement gotcha, hash workflows). Link across for depth; keep the
walkthrough out of the reference and vice-versa.

**No ephemeral handoff / pickup prompts in `docs/`.** "Paste this to start the next session" prompts and
per-PR handoff briefs rot the moment work moves. Keep them in the **PR description** (or a gitignored
scratch dir) while live. When a PR with a handoff doc **merges**, either fold its durable conventions into
the canonical doc it points at (e.g. `RECIPES.md`) or move the handoff to `archive/` as a historical
record — don't leave it loose in `docs/` claiming to be current.
