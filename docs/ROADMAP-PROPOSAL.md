# Roadmap & scope proposal — the 1.0 line, the 2.0 horizon, and magic

> **Status: PROPOSAL — for Xela112233 + zagwar to discuss. Not accepted.** Nothing here is canon until
> folded into `docs/ROADMAP.md` / `docs/SYSTEMS.md`. Drafted 2026-06-14, against the verified state below.
> It proposes three moves and lays out the alternatives so they can be argued with, not rubber-stamped.

## 0. The three moves (TL;DR)

1. **Draw a 1.0 line.** 1.0 = a complete co-op pack built on *buildable* content (config / recipe /
   curation / shipped-mods work that fits the odd-even cadence). The custom-mod megaprojects move to a
   named **2.0**.
2. **Demote magic** from a co-equal production pillar to a **woven optional vertical** — it isn't carrying
   production weight and forcing it to is expensive.
3. **Name the 2.0 arc:** the **"Commonwealth" governance/civilization update** (`pcmc-territory → realms →
   mint` + `conquest`), which also absorbs the **deep economy** (dynamic pricing + player-minted currency)
   and sits alongside the already-deferred **2.0 electricity overhaul** (#282).

The throughline: **stop cramming 2.0-sized custom-mod work into the pre-1.0 march.** It's the least-certain,
longest-lead work, and keeping it on the 1.0 critical path is how 1.0 either slips forever or ships things
half-built.

---

## 1. Verified state (June 2026)

- **Shipped:** `v0.7.0 — The Create spine` (`pack.toml` = 0.7.0; PATCHNOTES current). MC 1.21.1 / NeoForge
  21.1.233.
- **In flight:** `v0.7.1` (Xaero's, ore-vein validation, IPN→Inventory Essentials), `v0.8.0` (Stabilization
  I — Distant Horizons, SSRD, the **weave project**).
- **Live strategic branches:** `weaving-plan` (the Phase-2 weave-ledger — 3,000+ ratified weave decisions,
  the "taste weaving" engine), `outreach-plan` (`OUTREACH.md` + `ASSETS.md` — the public launch),
  `governance-plan` (zagwar's 4-mod governance/conquest program).
- **Issue health:** **143 open issues, every one milestoned — zero untriaged.** The triage invariant holds.
- **Cadence is working.** The odd/even machine (feature pillar → stabilization) has shipped cleanly through
  0.7.0. This proposal does not touch the cadence; it re-draws *what falls inside 1.0 vs after it*.

**The problem this addresses.** The road to 1.0 is well-defined, but it's scoped to include work that is
genuinely 2.0-sized and genuinely uncertain:
- **v0.13.0 — Economy & logistics is already the largest milestone at 30 open issues**, bundling the
  *shippable* economy (Numismatics, Trading Floor, Bountiful, aeronautics) with the **unsolved** dynamic
  pricing + player-minted currency (#221, #136, #150) — the part that, per the maintainer, "may need a fork
  or a whole new mod."
- **Governance (#260)** is a year-plus custom-mod program (`docs/GOVERNANCE*.md`); it currently sits in
  Backlog but is informally leaned into v0.13.0. It cannot fit one release.
- The **2.0 electricity overhaul (#282)** is *already* correctly deferred post-1.0 — the precedent this
  proposal extends to the other megaprojects.

---

## 2. The 1.0 line — what "done" means

**Principle:** 1.0 ships on content/config/curation/shipped-mods work that fits the cadence. Anything that
needs multi-month custom-mod development is 2.0. Below, each milestone with its **real open count**, owner,
and what changes.

| Milestone | Open | Goal (unchanged unless noted) | Owner | Scope change in this proposal |
|---|---:|---|---|---|
| `v0.8.0 — Stabilization I` | 17 | Profile/balance 0.7; land the **weave project** (clusters C-1…C-8) | shared | — (in progress) |
| `v0.9.0 — Survival` | 20 | Temperature × diet × seasons pressure, tuned | **Xela** | — (low mod-risk, long soak) |
| `v0.10.0 — Stabilization II` | 3 | Profile/balance 0.9 | shared | — |
| **`v0.11.0 — Colony`** *(was "Magic & MineColonies")* | 11 | **MineColonies as the load-bearing non-Create route** (cheap-basics calibration + colony lock list); **magic = a light balance/curation pass only** | **zagwar** (colony), Xela (magic-light) | **Magic demoted (§4):** drop deep Create-gating + locked-exclusives (#146), shelve Arcana (#118/#80); keep the shipped #75 web as-is |
| `v0.12.0 — Stabilization III` | 4 | Profile/balance 0.11; **#309 split-assessment** | shared | #309 largely **resolved** by the §3 deferral (see below) |
| `v0.13.0 — Economy & logistics` | 30 | **The *shippable* economy + the logistics/aeronautics arm** | shared | **Scoped down:** keep Numismatics/Trading Floor/Bountiful/loot-wiring + the full aeronautics/transport ladder; **defer the deep economy** (#221 dynamic pricing, #136 coin-tier minting, #150 player-minted faucet/sink, #240 player-driven value) **to 2.0** |
| `v0.14.0 — Stabilization IV` | 2 | Profile/balance 0.13 | shared | — |
| `v0.15.0 — Polish & site` | 16 | Wiki, onboarding, QoL, JEI, shaders/render decisions, **the rename (#212)**, the **weave review**, and the **public launch** (`outreach-plan`) | shared | — (this is where outreach + weave-review land) |
| `v1.0.0 — Release [NFR]` | 13 | Feature-frozen perf/RAM tuning, ore-gen finalized, CI required, ship the public build | **Xela** (release/box) | — |

**The two 1.0 scope cuts that matter:**

- **v0.11.0 → the Colony pillar.** MineColonies is load-bearing (cheap-basics route, locked exclusives, and
  the substrate the whole 2.0 governance program stands on). Magic is not. So this milestone re-centers on
  the colony route; magic drops to a light pass (§4). Of v0.11.0's 11 issues, ~8 are magic — most lighten or
  defer, so this milestone *shrinks*.
- **v0.13.0 → the shippable economy only.** Everything that uses mods already in the pack stays
  (coins, stalls, bounties, the airship/transport ladder, loot/mob-drop input wiring). The **unsolved**
  pricing/minting/currency work (#221/#136/#150/#240) — the part with no proven tooling — defers to 2.0,
  where it belongs *with* the governance minting layer (they're the same problem). This also de-risks the
  **#309 split question**: with the uncertain half removed, v0.13.0 is logistics + shippable-trade, tractable
  as one pillar — the split likely becomes unnecessary.

**Governance (#260) is not on the 1.0 path at all** — it stays in Backlog and becomes the 2.0 headline (§3).

---

## 3. The 2.0+ horizon — the megaproject arc

These are the genuinely big, genuinely differentiating bets. They share three traits: each is **custom-mod
or deep-systems work**, each is **multi-month and can't be built/verified in the web sandbox**, and each is
the **least schedule-predictable** work in the project. That is exactly why they belong *after* a shipped
1.0, not stuffed inside it.

- **2.0 — "Commonwealth" (the governance / civilization update).** The `pcmc-territory → pcmc-realms →
  pcmc-mint` trio + `pcmc-conquest`, scoped in `docs/GOVERNANCE*.md` (zagwar's #260). This is the pack's true
  differentiator — the thing that makes it "a cooperative political-economy survival pack," not "another
  Create pack." It's a year-plus effort; give it a real home as **the** 2.0, with the staged build order the
  governance docs already define (territory → realms-2a → hierarchy → mint, conquest alongside).
- **The deep economy** — dynamic pricing + player-minted currency (#221/#136/#150/#240). This is governance's
  **economic capstone** (minting *is* a realm power), so it folds into the 2.0 Commonwealth update rather
  than living a separate life. This is also where the "currency may need its own mod/fork" work lands.
- **2.x — the electricity overhaul (#282).** Create: Power Grid as the power spine, Create: Nuclear's
  return. Already deferred here by the v0.7.0 power review (`docs/POWER-MODS-REVIEW.md`); listed for
  completeness.
- **Optional — a magic re-anchor** (§4), *only if* post-1.0 play shows the loop actually needs a fourth
  producer (unlikely, given Create + MineColonies + bosses + governance).

**This line is already encoded in the player-facing ladder.** `wiki/economy-progression.md` (just drafted)
marks rungs 1–5 *live* (1.0) and rungs 6–9 — realms, federation, factions/armies, mint — *planned* (2.0).
The roadmap and the player map agree on where the line falls.

**Post-1.0 cadence:** keep the odd/even rhythm for live content (no world resets; freshness from curated
updates), exactly as `ROADMAP.md` already states — `2.1` feature → `2.2` perf → … The Commonwealth update is
large enough that it is itself a sequence of odd/even pairs, not a single release.

---

## 4. Magic — the solution

**Diagnosis (from the research pass).** Magic is *not* weakly connected — its connectivity islands are small
and the docs already accept them. Magic is weakly **load-bearing**: it has **no locked exclusives**, so
nobody in the trade loop actually *needs* a magic specialist. `SYSTEMS.md` frames magic as a co-equal
production route to Create, but it isn't carrying that weight — and the "non-Create producer with exclusives"
role is already filled, *better*, by **MineColonies** (cheap basics + exclusives + it's the governance
substrate) and **bosses** (high-tier drops). Magic is the redundant fourth leg. That redundancy is the "less
of a place" feeling.

**The current plan (v0.11.0) is the expensive option.** Making magic load-bearing means #146 (deep
Create-gating + defining locked exclusives) **plus** building `pcmc-arcana` (a NeoForge capability-bridge mod,
currently zero lines written) — and deep Create-gating **inverts** magic's self-contained apparatus theme.
That's a multi-mod, theme-bending investment to prop up a route the loop doesn't need.

**Three options:**

- **(A) Demote to a woven optional vertical — recommended.** Retire the "Create **or** Magic (production)"
  framing in `SYSTEMS.md`; the producers are **Create + MineColonies + bosses**, with magic a *flavor / side
  vertical*. For 1.0, keep magic exactly as it ships **today** (the #75 recipe-web, lightly Create-woven) —
  it's a perfectly good optional vertical. **Shelve `pcmc-arcana`.** Re-title v0.11.0 to the **Colony
  pillar**; magic shrinks to a light balance/curation pass. *Optional:* consolidate the magic surface — the
  research flags **Iron's Spellbooks** as the most droppable of the three cores (Ars is the hub, Occultism
  has the Create bridge via Occult Engineering) — to cut maintenance.
- **(B) Cut magic entirely.** Tightest identity (purely industrial-political-civilization), lowest
  maintenance — but loses player variety, discards the shipped #75 weave, and removes a whole fantasy
  register some players want. More aggressive than needed.
- **(C) Re-anchor hard now** (the current v0.11.0 plan). Most expensive, inverts the theme, props up a
  redundant route. Not recommended under the throughput lens.

**Recommendation: (A).** It matches reality, frees Xela's bandwidth and shrinks v0.11.0, removes the
awkward "magic is co-equal to Create" fiction, and preserves the option to re-anchor in 2.0 if play ever
demands it. **This is Xela's pillar**, so the call is his — which is presumably why the seam is visible from
the inside.

**Doc impact if (A) is accepted:** `SYSTEMS.md` §3 (drop magic from the production producers; note it as a
woven vertical), `ROADMAP.md` v0.11.0 (retitle + re-scope), `CUSTOM-MODS.md` (mark `pcmc-arcana` shelved),
`DESIGN.md` (the "Create or magic" north-star line).

---

## 5. Dates — a cadence-based timeline (calibrate against your real velocity)

No invented calendar dates. This is a **model**; correct the per-release durations with your actual
wall-clock and the windows move accordingly.

**Assumptions (state your real numbers):** two part-time maintainers + LLM tooling; an **odd feature pillar
≈ 4–8 weeks** implementation + soak; an **even stabilization ≈ 2–4 weeks**. Survival (v0.9.0) runs long — a
**full seasons cycle is weeks of soak** by design. Economy (v0.13.0) is the heaviest content pillar.

| Phase | Milestones remaining | Rough span |
|---|---|---|
| Now → mid-1.0 | v0.8 (in progress) · v0.9 Survival · v0.10 | ~3–5 months |
| mid-1.0 → late | v0.11 Colony · v0.12 · v0.13 Economy · v0.14 | ~3–5 months |
| 1.0 run-in | v0.15 Polish + launch · v1.0 perf freeze | ~1–2 months |

- **1.0 target window: roughly H1–H2 2027** (optimistic ~early 2027 at the fast end of the ranges;
  conservative ~late 2027). The two swing factors are the **Survival soak** and the **Economy (v0.13)** size.
- **2.0 "Commonwealth": 2028.** It's a year-plus custom-mod program *on top of* 1.0. Its build is gated on
  box-side spikes (the governance/conquest docs' `[needs box]` lists), so it's the least predictable work —
  the core reason it must not sit on the 1.0 critical path.

**The key schedule insight:** the 1.0 pillars are config/content and reasonably estimable; the 2.0
megaprojects are custom-mod and inherently not. Mixing them (the status quo) makes the *whole* schedule as
unpredictable as its least-predictable piece. Separating them gives 1.0 an estimable, defensible date.

---

## 6. Ownership

**Reality check:** **127 of 143 open issues are unassigned** (Xela 11, zagwar 9, a few co-assigned). So
ownership today is effectively pillar-level by intent, not issue-level by assignment. Proposal: **assign a
pillar lead per milestone now** (one click each), so the milestone bars read as someone's responsibility.

| Track | Lead | Scope |
|---|---|---|
| **Survival (v0.9), ore-gen, playtest** | **Xela** | Runs the live server/box; owns the survival interlock, ore-gen tuning, and the per-pillar playtests |
| **Magic (now light)** | **Xela** | The v0.11 magic balance/curation pass under option (A) |
| **Colony / MineColonies (v0.11)** | **zagwar** | The load-bearing colony route; the colony lock list |
| **Recipes / weave** | **zagwar** (#17) + shared | The recipe tracker; the `weaving-plan` project rides the odd thunderdomes |
| **Economy & logistics (v0.13)** | **shared** | Biggest pillar; consider a logistics-lead / economy-lead split if #309 still bites |
| **Release / CI / perf (v1.0)** | **Xela** | Box-side perf gates (#205/#48), CI-required (#79), ore bootstrap (#81) |
| **2.0 Commonwealth (governance/conquest)** | **zagwar** | His #260 program; the multi-month custom-mod arc |
| **Outreach / launch** | **shared** | `outreach-plan`; lands with v0.15 |
| **Stabilizations + curation thunderdomes** | **shared** | Both maintainers, every even milestone |

---

## 7. Loose ends the documentation review surfaced (a docs-pass, not blockers)

- **`ROADMAP.md` is stale at the top:** "Current release: v0.6.3" (we're at 0.7.0), and the v0.7.0 milestone
  block still reads as upcoming though it shipped. Refresh the header + move v0.7.0 to "shipped."
- **#17 (`needed-for-release`) sits in `Backlog / living`** — the one place a release-blocker is parked
  outside a release milestone. Either re-milestone it or drop the label.
- **Minor stale sections:** `docs/archive/` notes and `BOOT-LOG-BASELINE.md` carry closed-issue (#119/#120/
  #121) sections that can be trimmed. Archive docs (RELEASE-CADENCE, MODLIST-AUDIT) are correctly archived;
  no action.
- These are captured here so they're not lost; they're a cleanup pass, not part of this proposal's decisions.

---

## 8. The decisions for you two

1. **Accept the 1.0 line?** 1.0 = buildable content; defer custom-mod megaprojects to 2.0.
2. **Magic — A, B, or C?** Recommendation: **A** (demote to optional vertical, shelve Arcana, re-title
   v0.11 to Colony). *(Xela's call — his pillar.)*
3. **Name 2.0 = "Commonwealth"** (governance/conquest + deep economy)? And accept the **electricity overhaul
   (#282)** as a separate 2.x.
4. **Defer the deep economy** (#221/#136/#150/#240) out of v0.13.0 into 2.0 — and does that retire the #309
   split?
5. **Consolidate the magic mods** (drop Iron's Spellbooks), or keep all three as-is under option A?
6. **Assign pillar leads** on the milestones now (close the 127-unassigned gap)?

Once these are settled, the follow-up is a docs-pass folding the accepted moves into `ROADMAP.md`,
`SYSTEMS.md`, `CUSTOM-MODS.md`, and `DESIGN.md` — and re-titling/​re-scoping the v0.11.0 milestone on GitHub.

---

_Refs: `docs/ROADMAP.md` (the canonical map this extends), `docs/SYSTEMS.md` (the model magic-demotion
touches), `docs/GOVERNANCE*.md` (the 2.0 program, on `governance-plan`), `docs/POWER-MODS-REVIEW.md` (#282
2.0 electricity), `docs/CUSTOM-MODS.md` (`pcmc-arcana` status), `wiki/economy-progression.md` (the
player-facing 1.0/2.0 ladder), `docs/CURATION.md` (the keep/cut rubric for the magic-consolidation call).
Issue data: 143 open, all milestoned, as of 2026-06-14._
