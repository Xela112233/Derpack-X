# Roadmap & scope proposal — the beta arc (Beta 1.0 → 3.0 = Release 1.0), and magic

> **Status: PROPOSAL — for Xela112233 + zagwar to discuss. Not accepted.** Nothing here is canon until
> folded into `docs/ROADMAP.md` / `docs/SYSTEMS.md`. Drafted 2026-06-14; revised with maintainer feedback
> (magic skill-gate, compressed timeline, parallel 2.0, ownership). It proposes the moves and the
> alternatives so they can be argued with, not rubber-stamped.

> **Versioning — this reframes the whole arc.** Everything below is the **beta program**: the three majors
> are **Beta 1.0 → Beta 2.0 → Beta 3.0**. **Beta 3.0 *is* the official Release 1.0** — when the living world
> closes the Eco loop and the design vision is complete, the pack exits beta and that milestone ships as
> **Release 1.0**, the first non-beta public version. So the **early-August launch is the public *beta***
> (Beta 1.0, the foundation); the true **Release 1.0** lands at the end of the arc (~2027+), not in August.
> Mind the two distinct "1.0"s: **Beta 1.0** = the August foundation; **Release 1.0** = Beta 3.0, the
> design-complete graduation. Below, "1.0 / 2.0 / 3.0" mean the betas unless they say "Release." *(Open
> mechanical detail for you two: do the betas carry 1.x/2.x/3.x numbers, beta-tagged, or does the pack stay
> 0.x until Release 1.0? — a labeling call.)*

## 0. The moves (TL;DR)

1. **Draw a 1.0 line.** 1.0 = a complete co-op pack on *buildable* content (config / recipe / curation /
   shipped-mods work). Target: **early August 2026.** The custom-mod megaprojects move to **2.0**.
2. **Magic = a skill-gated specialization** (the middle ground between demote and deep-re-anchor): a few
   players invest in a **bare-bones skill gate** to become magic specialists whose output brings value to
   their group/realm. Keeps magic *and* gives it a real Eco-style place. **Evaluate a full skill/
   specialization system for 2.0** — which doubles as the loop's missing "hard specialization" lever.
3. **2.0 = parallel background tracks**, not a linear sequel. The **"Commonwealth" governance/conquest**
   program, the **deep economy**, the **electricity overhaul**, the **full skill system**, and tech
   expansions — **extending the tech tree to space** and **machine miniaturization** — run as *simultaneous*
   milestones in the background **while 1.0 is live and played**. Target: **end of 2026 likely, end of
   spring 2027 at the latest.**

**The throughline (revised around agent velocity):** agent-driven authoring is no longer the schedule
driver — a recent mod backport estimated at 1–2 weeks landed in under two hours. So the binding constraint
is **box-side verification** (playtest, perf, time-based soak) that agents *can't* do, of which there is one
human and one server. Custom-mod work carries the heaviest verification burden (the governance docs'
`[needs box]` lists), so it must run as **background-parallel tracks that don't compete with 1.0's
verification throughput** — that's *why* it's 2.0, not because it's slow to write.

*(Beyond 2.0, a **Beta 3.0 horizon — "the living world" (which ships as the official Release 1.0)** —
anchored on the ecological-cost mechanic — is sketched after §3: a direction marker, not a decision here.)*

---

## 1. Verified state (June 2026)

- **Shipped:** `v0.7.0 — The Create spine` (`pack.toml` = 0.7.0). MC 1.21.1 / NeoForge 21.1.233.
- **In flight:** `v0.7.1` (Xaero's, ore-vein validation), `v0.8.0` (Stabilization I — the **weave project**).
- **Live strategic branches:** `weaving-plan` (3,000+ ratified weave decisions), `outreach-plan`
  (`OUTREACH.md` + `ASSETS.md` — the public launch), `governance-plan` (zagwar's 4-mod program).
- **Issue health:** **143 open, every one milestoned — zero untriaged.**
- **Cadence works.** Odd feature → even stabilization has shipped cleanly through 0.7.0. This proposal
  doesn't touch the cadence; it re-draws *what's inside 1.0 vs after it* and compresses the clock.

**The problem:** the road to 1.0 is well-defined but scoped to include 2.0-sized custom-mod work — most
visibly **v0.13.0 (30 open issues)**, which bundles the shippable economy with the *unsolved* dynamic
pricing + player-minted currency (#221/#136/#150), and **governance (#260)**, a multi-mod program leaning
informally into v0.13. The **2.0 electricity overhaul (#282)** is *already* deferred — the precedent this
extends.

---

## 2. The Beta 1.0 line — milestones, owners, targets

**Principle:** 1.0 ships on content/config/curation/shipped-mods work. Anything needing multi-month
box-verified custom-mod development is 2.0. Targets assume agent-authored content + serialized box
playtests (see §5); they're tight by design to hit **early August**.

| Milestone | Open | Owner | Target | Key change in this proposal |
|---|---:|---|---|---|
| `v0.8.0 — Stabilization I` | 17 | shared (weave: Xela) | late June | — (in progress; lands the weave project) |
| `v0.9.0 — Survival` | 20 | **zagwar + Xela** | late June | — (soak overlaps into live 1.0; see §5 risk) |
| `v0.10.0 — Stabilization II` | 3 | shared | early July | — |
| **`v0.11.0 — Colony, territory & the magic skill-gate`** *(was "Magic & MineColonies")* | 11 | **zagwar** (colony + territory) · **Xela** (magic skill-gate) | early–mid July | MineColonies as the load-bearing route **+ ship `pcmc-territory`** (unified colony+OPAC borders — the spiked claims substrate) **+ magic re-cast as a skill-gated specialization (§4)** |
| `v0.12.0 — Stabilization III` | 4 | shared | mid July | #309 split largely **resolved** by the §3 deferral |
| `v0.13.0 — Economy & logistics` | 30 | **zagwar** (economy) · **Xela** (logistics) | mid–late July | **Scoped down:** keep coins/stalls/bounties/loot-wiring + the full aeronautics/transport ladder; **defer dynamic pricing + minting (#221/#136/#150/#240) to 2.0** |
| `v0.14.0 — Stabilization IV` | 2 | shared | late July | — |
| `v0.15.0 — Polish & site + launch` | 16 | shared | late July | Wiki, onboarding, JEI, the rename (#212), the **weave review**, and the **public launch** (`outreach-plan`) |
| `v1.0.0 — Release [NFR]` | 13 | **Xela** (box) | **early August** | Feature-frozen perf/RAM, ore-gen final, CI required, ship |

**The two 1.0 scope cuts that matter:**

- **v0.11.0 → Colony + a magic *skill-gate*.** MineColonies is load-bearing (cheap-basics route, locked
  exclusives, and the substrate the 2.0 governance program stands on); zagwar owns it. Magic is re-cast (§4)
  as a **skill-gated specialization** — Xela's pillar, now a defined small build (a bare-bones skill gate +
  config), not the expensive deep-re-anchor. The milestone re-centers on the colony route with magic as a
  contained, well-shaped add.
- **v0.13.0 → the shippable economy only.** Everything using mods already in the pack stays. The **unsolved**
  pricing/minting work (#221/#136/#150/#240) defers to 2.0, where it belongs *with* governance's minting
  layer (same problem). Ownership splits cleanly — **economy = zagwar, logistics = Xela** — which also
  de-risks the **#309 split** question: with the uncertain half gone, the pillar is tractable as one.

**Territory pulled into Beta 1.0 — the one custom-mod exception (v0.11).** `pcmc-territory` (the unified
colony + OPAC "who-governs-this-chunk" borders layer) is **already spiked (PR #1)**, **read-only and
low-runtime-risk** (it owns no claims — it just resolves the two systems already in the pack), and **useful
standalone** (a co-op group gets unified territory/borders without needing laws). Shipping it in 1.0 delivers
value now *and* de-risks 2.0: Realms/Conquest/Mint drop onto a foundation that already shipped. The
"custom-mod ⇒ 2.0" rule still holds for the heavy, `[needs box]`-laden tracks; Territory is the justified
exception.

**The rest of governance (#260) is not on the 1.0 path** — only its claims substrate (above) is; Realms,
Mint, and Conquest are 2.0 tracks (§3).

---

## 3. The Beta 2.0 horizon — parallel background tracks

2.0 is **not** a linear sequel. It's a set of **simultaneous milestones worked in the background while 1.0
is live and being played** — enabled by agent-parallel authoring, gated only by box-verification
throughput. Each is its own milestone/branch; they start *now* (authoring can begin before 1.0 ships) and
land as 2.x content drops between **end of 2026 (likely)** and **end of spring 2027 (latest)**.

| Track (proposed milestone) | What it is | Owner | Depends on |
|---|---|---|---|
| **Territory** (`pcmc-territory`) — *pulled to Beta 1.0 (v0.11)* | The claims substrate ("who governs this chunk"); **ships in 1.0** — the 2.0 tracks build on it | zagwar | — (foundation, in Beta 1.0) |
| **2.0-B — Realms** (`pcmc-realms`) | Tiers, hierarchy, the hard/soft law engine, guard enforcement | zagwar | Territory (shipped in Beta 1.0) |
| **2.0-C — Mint + deep economy** (`pcmc-mint`) | Treasury, taxes, charters, minting **+ dynamic pricing / player-minted currency** (#221/#136/#150/#240) | zagwar | Realms; may be its own mod/fork |
| **2.0-D — Conquest** (`pcmc-conquest`) | Joinable NPC factions, faction kit, armies; replaces Valarian Conquest | zagwar | Realms wanted-signal (seam only — parallel) |
| **2.0-E — Electricity overhaul** (#282 Create: Power Grid) | The power-spine reimagining; Nuclear's return | Xela | independent (gated by #281 perf) |
| **2.0-F — Skill / specialization system** | Generalize the 1.0 magic skill-gate into a full Eco-style specialization system across routes; **magic 2.0 rides it** | Xela | the 1.0 magic gate (pilot) |
| **2.0-G — Space tech tier** (Ad Astra — lead candidate) | Extend the tech tree off-world: an orbital/planetary endgame — the new top scarcity frontier and the pack's clear endgame goal | Xela | Electricity (E) + logistics/aeronautics |
| **2.0-H — Machine miniaturization** (microchip/controller — mod TBD) | A chip/controller tier that shrinks machine footprint; late-game compaction | Xela | Electricity (E) — the PCB/circuit tier |

**Why this shape works:** the governance work (B Realms → C Mint) builds on the **Territory layer already
shipped in Beta 1.0**; conquest (D) runs parallel off the shared wanted-signal seam; electricity (E) and
the skill system (F) are fully independent; and the tech tracks (G space, H miniaturization) build on E's
power and circuit tiers. Agents author all seven 2.0 tracks (B–H) concurrently; the **box-verification
queue** (one human, one server) is the real serializer — which is exactly why these stay off 1.0's
verification path.

**Strategic notes:**
- **The deep economy folds into Mint (2.0-C).** Dynamic pricing + minting *is* governance's economic
  capstone — same problem, same home. This is where the "currency may need a fork or whole new mod" work
  lives.
- **The full skill system (2.0-F) is the loop's missing "hard specialization" lever.** `SYSTEMS.md` flags
  that our specialization is *soft* — a determined player can go self-sufficient and erode trade, with
  locked exclusives the only anti-erosion device. A real skill/specialization system (Eco's hard skill-cap
  analog) is the mechanism that *hardens* it across the board. So 2.0-F isn't just "magic 2.0" — it's the
  fix for an open loop-design problem, with magic as its 1.0 pilot.
- **Space (2.0-G) closes the loop at the top — it's not a bolt-on.** `SYSTEMS.md` §5 says growth renews
  scarcity *one tier up* — farther, scarcer, more dangerous frontiers. Space is literally that: off-world
  resources re-open scarcity → production → trade at a new tier, and it hands the pack the clear endgame
  *goal* the video essay praised in Create: Above and Beyond. Three anchors — Create (top production tier),
  scarcity/economy (frontier resources → new trade), logistics/aeronautics (the movement arm goes orbital).
- **Miniaturization (2.0-H) is a perf lever, not just QoL.** Smaller machines = fewer block entities =
  direct relief for the **#151 chunk-save serialization spikes** the perf docs flag. So it earns its place
  twice — a Create tech reward *and* a late-game perf win — built on the electricity overhaul's circuit/PCB
  tier (chips are PCB-tier outputs).

**This line matches the player ladder.** `wiki/economy-progression.md` already marks rungs 1–5 *live* (1.0)
and rungs 6–9 — realms, federation, factions, mint — *planned* (2.0). The roadmap and the player map agree.

---

## The Beta 3.0 horizon — the living world, shipping as Release 1.0 (a sketch)

> Further out and more speculative than 2.0 — a *sketch* to capture direction, not a plan. 3.0 follows 2.0
> (so ~2027+), and like 2.0 it's a set of parallel background tracks, not a linear release.

**The arc.** 1.0 is the industrial foundation; 2.0 is civilization + reaching space; **3.0 is the world
becoming a living system that responds to — and is damaged by — the economy the players have built.** It
closes the one edge of the Eco loop still open: **production → world**. After 3.0 the causal loop
(scarcity → pressure → production → economy, and back to scarcity *through the world itself*) is finally
self-driving, not just designed. **This is the out-of-beta line:** with the loop complete and the design
vision "done," Beta 3.0 ships as the official **Release 1.0** — the pack's first non-beta public version,
and the point freshness shifts from build-out to the live-service content cadence.

### What fits 3.0 (already seeded; lands here by its own logic)

- **Ecological cost — the flagship (`SYSTEMS.md` §3a).** Production damages the shared world: heavy
  extraction degrades the local environment and worsens scarcity; industrial output feeds pressure
  (pollution → local temperature into Cold Sweat's range, fouled water/crops). It's **adopted in principle
  but parked as unsolved** — two hard blockers (an infinite world defeats naive enforcement; it's
  TPS-heavy). **Why it's 3.0, not 2.0:** the fix is to scope damage to a player's *claimed/settled
  footprint*, which **needs 2.0's claims/realms to exist first**, and it wants 2.0's miniaturization/perf
  headroom. Likely custom KubeJS/datapack, event-driven and locality-bounded. The single most
  loop-completing feature left.
- **Dynamic weather as survival pressure (`SYSTEMS.md` "Open/unsettled").** Storms, severe weather,
  exposure — considered, parked as a TPS cost. Pairs naturally with ecology (industry → climate) and fits
  the living-world theme; affordable once the 2.0 perf work (miniaturization, electricity) lands.
- **Full maturity of the player economy** — if 2.0-C's dynamic pricing only partly lands, a market that
  genuinely *reacts* to ecological scarcity shifts matures here (ecology makes prices move for real reasons).

### What we could add to 3.0 (candidates to evaluate)

- **A computational / programmable layer** — CC: Tweaked, or Create logic/displays: programmable automation,
  the "intelligence" tier after industry and space. Anchors to Create production as a top automation tier.
- **Deeper frontiers** — beyond the first space tier (2.0-G): more planets/dimensions, deep-space resources —
  the loop re-opened one more tier up. (If 2.0-G ships one destination, 3.0 ships the system.)
- **Diplomacy & treaties** — formal inter-realm diplomacy, trade agreements, alliances atop 2.0's
  governance: the social endgame for a maturing server. Anchors to governance + economy.
- **Endgame challenge content** — raids, mega-bosses, dungeon dimensions for fully-geared groups. Anchors to
  danger-pressure + the boss-drop production route.
- **Live-service cadence** — formalize the "freshness from curated content, no resets" principle into a
  seasonal/event rhythm for the public server (rotating limited-time content), keeping a public pack alive
  long-term.

**Same shape as 2.0:** parallel background tracks, agent-authored, gated by box verification, landing
rolling. The ecological mechanic is the one that *defines* 3.0's identity; the rest are optional weaves
around it. Scope it properly once 2.0 is landing — this is a direction marker, not a commitment.

---

## 4. Magic — the solution (the skill-gated middle ground)

**Diagnosis.** Magic isn't weakly *connected* (its islands are small and accepted); it's weakly
**load-bearing** — no locked exclusives, so nobody *needs* a magic specialist, and the non-Create producer
role is already filled better by MineColonies + bosses. That redundancy is the "less of a place" feeling.

**The move (option D — recommended): make magic a *specialization*, not a *production route*.** A few
players invest in a **bare-bones skill gate** to become magic specialists; their magic output (buffs,
reagents, enchants, services) brings value to their **group, and later their realm**. This is the Eco core —
*specialization drives interdependence* — and it gives magic a real place **without** the deep Create-gating
that option C needed (which inverted magic's self-contained theme). Magic stops pretending to be a parallel
Create and becomes *the thing the realm's mage does*.

- **1.0 — bare-bones skill gate.** Lock the meaningful magic outputs behind a single **Arcana/Magic skill** a
  player levels by doing magic. The vehicle: **Project MMO (PMMO)** is the natural fit — it natively gates
  use/craft behind a "magic" skill, deploys minimally (one skill, the rest off), and **is exactly what 2.0
  scales up** (1.0-bare-bones and 2.0-full become the same mod, not two systems). The lighter alternative is
  a **KubeJS progression gate** (leverages the #220 locking mechanism, no new mod) if a full skill mod feels
  heavy for 1.0. *Pick the mechanism — that's a decision (§8).* No new mod is fine for the *concept*; the
  cost is one mod-add + config + a playtest, in the v0.11 odd/feature window.
- **In 1.0 there are no formal realms yet**, so "value to their realm" reads as **value to their co-op group/
  colony** — a specialist others trade with. When formal realms arrive in 2.0, magic specialists slot in as
  realm roles (the court mage), and the gate generalizes into the **full specialization system (2.0-F)**.

**Options also considered:**
- **(A) Demote to optional vertical** — cheapest, but leaves magic with no real pull; the skill-gate is
  barely more work and gives it a *reason to exist*. Superseded by D.
- **(B) Cut entirely** — tightest identity, but discards the shipped #75 weave and a whole register some
  players want. More than needed.
- **(C) Deep re-anchor now** (the current v0.11 plan: gate all magic behind Create + build `pcmc-arcana`) —
  most expensive, inverts the theme, props up a redundant *production route*. D gets magic's "real place"
  for a fraction of the cost by making it a *specialization* instead.

**Doc impact if D lands:** `SYSTEMS.md` §3 (magic re-cast from co-equal producer to a skill-gated
specialization), `ROADMAP.md` v0.11 (retitle + the skill-gate scope), `CUSTOM-MODS.md` (`pcmc-arcana`
deferred to the 2.0 skill track, not 1.0), `DESIGN.md` (the "Create or magic" line). *Optional* consolidation
of the magic mods (the research flags Iron's Spellbooks as the most droppable core) stays available as a
curation call.

---

## 5. Dates — built on agent velocity (the binding constraint is the box, not the keyboard)

The reframe that makes the targets real: **agents collapse authoring time** (the 1–2-week-backport-in-2-hours
data point), so the schedule is no longer gated by *writing* content. It's gated by what agents **can't** do
— in-game playtest, perf measurement, and **time-based soak** — all of which run on **one human + one
server**. Plan accordingly: **parallelize authoring, serialize verification.**

### 1.0 — target early August 2026 (~7 weeks)

The §2 targets pull to early August: v0.8 (late June) → v1.0 (early August). Agents author the pillars
**concurrently** on branches; the human serializes box playtests as each lands. The critical path is the
**playtest queue**, not the build.

- **#1 risk: the Survival seasons soak.** A full seasons cycle is weeks of wall-clock that can't be
  agent-accelerated, and early-August is too tight to fit one. **Accept it explicitly:** ship v0.9 on an
  existence-and-function check and let seasons *tuning* continue **into live 1.0** as a post-launch patch,
  rather than holding the release for a full cycle. Everything else fits the queue.
- Secondary risks: the v0.13 breadth (mitigated by the deep-economy deferral) and the v1.0 perf gate
  (#205/#48) — both box-verification items; protect the queue.

### 2.0 — a rolling timeline (end of 2026 likely, spring 2027 latest)

The eight tracks (§3) don't ship as one big-bang release — they **land on a rolling basis as each clears box
verification**, in the background while 1.0 is played. Authoring starts now (Territory is already spiked);
the bulk of *landing* is post-launch, when the box frees from 1.0.

| Window | Tracks landing | Note |
|---|---|---|
| **Now → 1.0** (Jun–early Aug) | Author-ahead only: **E** Electricity, **F** skill-system design + the 1.0 magic-gate pilot (Territory itself ships *in* Beta 1.0, v0.11) | Overlaps 1.0; verification queued behind it |
| **Launch → Sep 2026** | **Foundations:** **B** Realms MVP (guard-enforced law, on the Beta-1.0 Territory layer), **E** Electricity overhaul | 1.0 live → the box frees for 2.0 |
| **Oct–Nov 2026** | **Parallel build-out:** **B** Realms hierarchy, **D** Conquest, **F** Skill system, **G** Space, **H** Miniaturization | Independent tracks land as ready (rolling 2.x) |
| **Nov–Dec 2026** | **The hard last mile:** **C** Mint + deep economy (dynamic pricing / player-minted currency) | The least-certain track, deliberately last |
| **Jan → spring 2027** | Integration, balance, buffer for any verification-backed slippage | Absolute-latest landing for the custom-mod tracks |

The binding constraint throughout is the **box-verification queue** (one human, one server), now shared with
live-1.0 support; the custom-mod tracks each carry `[needs box]` spikes. Agents keep authoring ahead of it,
so **box throughput is the real 2.0 schedule.** End-2026 if the queue keeps pace; spring-2027 is the buffer.

**Calibrate against your real cadence.** These windows assume the playtest/soak queue keeps moving; if box
time is scarcer than assumed, the dates stretch — but the *authoring* won't be the reason.

---

## 6. Ownership

127 of 143 open issues are unassigned — ownership is pillar-level by intent. Proposal: **assign these leads
on the milestones now.**

| Track | Lead | Notes |
|---|---|---|
| **Survival (v0.9), ore-gen, playtest** | **zagwar + Xela** | Both; Xela runs the box/soak |
| **Colony / MineColonies (v0.11)** | **zagwar** | The load-bearing non-Create route |
| **Magic skill-gate (v0.11)** | **Xela** | The bare-bones gate (§4) |
| **Economy (v0.13)** | **zagwar** | Coins/shops/bounties/faucet-sink |
| **Logistics / aeronautics (v0.13)** | **Xela** | Airships/transport ladder |
| **Weave project** | **shared, mostly Xela** | `weaving-plan`, rides the thunderdomes |
| **Release / CI / perf (v1.0)** | **Xela** | Box-side gates (#205/#48/#79/#81) |
| **Territory (Beta 1.0, v0.11) + 2.0 governance & conquest (B–D)** | **zagwar** | The claims substrate ships in 1.0; Realms/Mint/Conquest in 2.0 |
| **2.0 tech tracks — electricity (E), skill system (F), space (G), miniaturization (H)** | **Xela** | #282 + the magic-gate generalization + the new top-tier tech |
| **Outreach / launch** | **shared** | `outreach-plan`; lands with v0.15 |

---

## 7. Loose ends the review surfaced (a docs-pass, not blockers)

- **`ROADMAP.md` is stale at the top:** "Current release: v0.6.3" (we're at 0.7.0) and the v0.7.0 block reads
  as upcoming though it shipped. Refresh on the docs-pass.
- **#17 (`needed-for-release`) sits in `Backlog / living`** — the one release-blocker parked outside a
  release milestone. Re-milestone or drop the label.
- Minor stale closed-issue (#119/#120/#121) sections in `BOOT-LOG-BASELINE.md` can be trimmed. Archive docs
  are correctly archived.

---

## 8. The decisions for you two

1. **Accept the 1.0 line + the end-of-August target?** (Buildable content in; custom-mod megaprojects to 2.0.)
2. **Magic — option D (skill-gated specialization)?** And **which mechanism for the 1.0 bare-bones gate —
   PMMO (recommended, scales to 2.0) or a KubeJS progression gate (no new mod)?**
3. **2.0 as parallel background tracks** (B–H in §3, on the Territory layer pulled into Beta 1.0), started
   now, landing end-2026 / spring-2027?
4. **Defer the deep economy** (#221/#136/#150/#240) out of v0.13 into 2.0-C — and does that retire the #309
   split?
5. **Adopt the full skill system (2.0-F) as the loop's "hard specialization" lever**, or keep specialization
   soft and skill-gate magic only?
6. **Consolidate the magic mods** (drop Iron's Spellbooks) under D, or keep all three?
7. **Confirm the two new tech-expansion tracks** — 2.0-G space (Ad Astra the lead candidate) and 2.0-H
   machine miniaturization — and **source their mods** (an existing mod, or a small custom one for the
   chip/controller if nothing fits)?
8. **Assign the §6 pillar leads** now (close the 127-unassigned gap)?
9. **3.0 direction (far-out):** does "the living world" — the ecological-cost mechanic (§3a) as the 3.0
   flagship — resonate, or should 3.0 anchor elsewhere (deep-space, a computational tier, diplomacy)?

Once settled, the follow-up is a docs-pass folding the accepted moves into `ROADMAP.md`, `SYSTEMS.md`,
`CUSTOM-MODS.md`, `DESIGN.md`, re-titling the v0.11 milestone, and creating the 2.0 track milestones on
GitHub.

---

_Refs: `docs/ROADMAP.md`, `docs/SYSTEMS.md` (the model magic re-cast touches; the soft→hard specialization
note), `docs/GOVERNANCE*.md` (the 2.0 A–D tracks, on `governance-plan`), `docs/POWER-MODS-REVIEW.md` (2.0-E,
#282), `docs/CUSTOM-MODS.md` (`pcmc-arcana` → 2.0-F), `wiki/economy-progression.md` (the player-facing 1.0/2.0
ladder), `docs/CURATION.md` (the magic-consolidation call). Issue data: 143 open, all milestoned, 2026-06-14._
