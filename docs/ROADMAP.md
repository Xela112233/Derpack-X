# Project Commonwealth — Roadmap (multi-phase)

> Living roadmap for the build-out. **Goal:** a cooperative PvPvE Create pack for a small crew (~10)
> where every mod earns its place by anchoring to one of the systems — **Create** or **magic**
> (production), **survival** (pressure), or the **economy** (distribution & trade, incl.
> logistics/aeronautics) — or it's cut. They form one causal loop (**scarcity → pressure → production →
> economy**, inspired by *Eco*): scarcity + locked exclusives drive an emergent but **load-bearing** trade
> economy, and complex tech unlocks via MineColonies or boss drops. (The old flat "five systems" framing
> is retired — see `docs/SYSTEMS.md`, **the canonical structural model**.) Also: `docs/DESIGN.md` (the
> why) and `docs/RECIPES.md` (the recipe convention + triage ledger, merged with #62).
>
> **Standing caveat:** the dev sandbox is **headless** — everything is authored + structurally
> validated (JSON valid, JS syntax OK); **an in-game playtest is required before any of it merges**
> (JEI shows the new recipes/veins, the originals are gone, no errors in the log).

**Status:** ✅ done & on `main` · 🟦 in review (open PR) · 🚧 in progress · ⏳ queued · 💤 not started

**Current release:** **v0.6.3** (Foundation line). MC **1.21.1** / NeoForge **21.1.233**.

## Where the work lives

`main` now carries **Phases 0–6**: tooling + ground-truth digest, ore-gen, guns, design docs,
`pr-checks.yml`, the **recipe overhaul (#62, v0.4.6)**, the **curation/spawn-gating pass (#88, v0.5.0)**,
the **GTMOGS ore-gen rework (#82, v0.5.1)**, and the **magic web (#75, v0.5.2)**.

**No stacked feature PRs remain.** The custom mods (Project Commonwealth Arcana #80, Kill Feed #196) moved to their
own repos (`theasshats/pcmc-arcana`, `theasshats/pcmc-killfeed`) — those prototype PRs are
superseded; see `docs/CUSTOM-MODS.md`. Everything else is now issue work, not stacked PRs.

| Release | Shipped |
|---|---|
| v0.4.6 | Recipe overhaul (#62) — "made through Create" + almost-unified dedup + EMI→JEI |
| v0.5.0 | Curation / spawn-gating (#88) — mob/boss cuts + In Control! gating + kobold swap |
| v0.5.1 | Ore veins (#82) — GTMOGS GregTech-style mix veins replace the per-ore model |
| v0.5.2 | Magic web (#75) — KubeJS bridges weaving Ars / Iron's / Occultism into the Create spine |

## Phases

### Phase 0 — Tooling & ground truth ✅ (on `main`)
packwiz pack + editor + build CI. **Ground-truth digest** at `tools/mod-data/` (facts + recipe + loot
indexes) via **Actions ▸ "Ground-truth digest"** — re-run after modlist changes.

### Phase 1 — Eco foundation: ore-gen + scarcity ✅ on `main` — **model settled: GTMOGS**
The **#56 per-ore model has been retired**; **PR #82 (GTMOGS GregTech-style mix veins) shipped as the
v1 model in v0.5.1**, strictly regional. The #56-vs-#82 decision (**#93**) is **closed**. Remaining
ore-gen work is all release-tuning:
- **#81** — early-game ore bootstrap, now that the vanilla starter trickle is gone (`needed-for-release`).
- **#318** — the consolidated ore-gen tracker (absorbed #116 and #58 in v0.7.1): tuning surface, rarity math, the kingdom/monopoly design lens, surface-biome cave-ore want, release exit checklist. #81/#261/#290/#296/#298/#300 ride as sub-issues.
- **#115** — player-facing prospecting wiki page (kicked off by the GTMOGS move).

### Phase 2 — Guns / combat ✅ on `main` (#55, #59)
TaCZ + Create: Immersive TaCZ + Create: Armorer gunpack + Create Big Cannons. Merged and playtested.
- **Later:** gun→Create recipe reconciliation + sequenced-assembly depth (#112); gun loot (#18).

### Phase 3 — Recipe overhaul: "made through Create" ✅ on `main` (#62, v0.4.6)
Routed crafting through Create parts + methods across the pack (`docs/RECIPES.md` = convention + full
triage ledger). Metals foundation + almost-unified dedup; per-mod gating; the decoration/equipment
weave; worldgen-island bridges; EMI→JEI + JER/JEED; the `tools/connectivity/` connectivity toolkit.
- **Dedup follow-ups:** #103 (functional-duplicate parts), #101 (limestones), #102 (carbon/graphene), #113 (connectivity islands).
- **Boot-log correctness:** #119 (dropped recipes), #120 (mangled loot tables), #121 (remaining ERROR noise).

### Phase 3b — Magic web ✅ on `main` (#75, v0.5.2)
KubeJS bridges weaving Ars Nouveau / Iron's Spells / Occultism into the Create spine.
- **Follow-ups:** #122 (post-merge balance, playtest), #123 (flagship depth, conceded from #75).
- **Code track:** #80 (Project Commonwealth Arcana) — the bridges KubeJS can't express; #118 tracks its P3 deviation.

### Phase 4 — Loot pass ⏳ (#18)
`lootjs` + the loot index: pull structure/dungeon loot back from free end-game gear; seed Numismatics
coins. Feeds the economy pillar (#90). Not started.

### Phase 5 — Economy (Numismatics) 💤 #90
Coin values, vendors + pricing, Trading Floor, **Bountiful** bounties, and **wiring the inputs** (mob
drops / structure loot → sellables + bounties). The loop's highest-leverage,
thinnest-in-integration system. **#129** (mob drops in the connectivity tool) is the enabler.

### Phase 6 — Curation & de-dup ✅ #88 on `main` (v0.5.0) — follow-ups open
The mob/boss + duplicate-mechanic pass landed (arphex/creeper-overhaul/mutant-monsters cuts, In Control!
spawn-gating, kobold swap). Remaining curation: **#94** (un-anchored mods + dup clusters), **#106/#108** (spawn-gating allowlist +
behavioral review), **#21** (living removal tracker), **#87** (Create Stuff'N Additions balance).
_Settled: #60 umapyoi cut, #100 orphan libs swept, #107 Ender Moon kept (flavor → v0.9.0)._ MineColonies is settled as a woven on-ramp (#92).

### Phase 7 — Stability & QoL 💤 (needs a live server)
Playtest pre-gen #98 (force-multiplier) · RAM/GC #48 · flight interactions #43 · claims #25 ·
purple-arrows bug #3 · shaders caveat #1 · render distance #2. (TPS audit #42 resolved in v0.3.5/0.3.6.)

### Phase 8 — CI / merge-gating ✅ (`pr-checks.yml`)
PR-triggered checks gate merges (packwiz index freshness, manifest lint incl. the `pin` gotcha, KubeJS
JS/JSON + config TOML parse, Go build/vet for the editor + site). Full rundown in `docs/CI-CHECKS.md`.
- **Open (#79):** flip the checks to *required* on `main` via a branch-protection ruleset — #76 has
  merged, so this is unblocked (a one-time admin toggle).
- **CI hygiene:** #127 (consolidate the 3 overlapping index-refresh workflows).

## Needed for release

The set that gates the next milestone release (tagged `needed-for-release`). Everything else is
post-release polish or ongoing curation.

- [x] **Ore-gen model decided** — **#93 closed**; GTMOGS (#82) shipped as the v1 model in v0.5.1.
- [ ] **Early-game ore bootstrap** — #81 (vanilla trickle removed; define the replacement). Phase 1.
- [ ] **Ore-gen tuned for release** — #318 (the consolidated tracker; absorbed #116/#58 in v0.7.1). Phase 1.
- [x] **Create recipe spine landed** — #62 (v0.4.6). Follow-ups: #101/#102/#103/#113.
- [x] **Pack renamed** — #78 (renamed to **Project Commonwealth**).
- [ ] **CI required on `main`** — #79 (turn on the `pr-checks.yml` ruleset; Phase 8).
- [x] **Design/goal docs merged** — PR #61; the systems-loop direction is canon (`docs/SYSTEMS.md`).
- [x] **Release artifacts sane** — `.mrpack` dropped (#73); Prism installer is the single ~5 MB artifact.

**Explicit non-goals** (keep off the release path): a server-pack artifact, NeoForge auto-bump,
bundled-jar distribution — all deliberately out, see `docs/DESIGN.md`.

## Milestones — the road to 1.0 (odd/even cadence)

Open issues are organized into GitHub **milestones** on an **odd/even release cadence** (this section is the canonical cadence map; the original proposal + extended rationale is archived at `docs/archive/RELEASE-CADENCE.md`): the minor version's parity sets the job. **ODD = feature / integration** — add the pillar's mods, weave them in structurally, then curate by merit (**the thunderdome** — see below). **EVEN = stabilization** — profile, balance, and bug-fix what the prior odd added; *no new mods or features*. One **pillar per odd version**, each followed by its own **even stabilization pass** so feature debt never stacks; the version count is open (extend with odd/even pairs as the pillars need). **`1.0.0` is even → a feature-frozen performance patch.** **Every issue maps to a milestone, and so does every *new* issue** — the pillar's odd version, the stabilization even that follows it, or `📋 Backlog / living`. (Milestones are created/edited on GitHub; this list is the map. See `docs/TRIAGE.md` for the process.)

**Two kinds of curation — both apply the #157 rubric, on different gates.** The **odd thunderdome** is curation *by competition*: the version's new adds **and** the existing related set are reviewed against each other on merit (system-anchor, no-redundant-mechanic, woven-not-floating); winners stay, losers are cut or reworked, then the version freezes. The **even stabilization** is *performance pruning*: is a mod's TPS/RAM cost justified (spark — #147/#160)? Cut or optimize the expensive, fix what broke, balance. **Web % is a compass, not a gate** — use the connectivity tool (#129) to *find* off-web clusters worth a human look, never as a release target (target it and it gets gamed — Goodhart). The deep *taste*-weaving question now has working rules: the **seam ladder** (tier T_n ≥ n+1 cross-mod seams) and the **functional-part rule** (gate on the part the item works like), both in `docs/CREATE-SPINE.md` Parts 1 + 2b-ii.

**Why this cadence (rationale, folded in from the former RELEASE-CADENCE doc).** One **themed pillar per odd version** rather than mixing pillars or chasing a per-release web-% target, because: themed milestones **motivate** ("ship the Economy update" feels like an accomplishment in a way "raise cohesion to 78%" doesn't); each pillar release is a natural **ecosystem-survey** moment to pull in what's new and good; and a per-release web-% target is **gameable** (Goodhart — the compass note above). Weaving has **two depths** — *structural plumbing* (tags unified, no hard dups, recipes valid, nothing orphaned; largely automatable) and *taste weaving* (the cross-system connections that are actually good and themed — the hard part), the latter governed by the seam-ladder + functional-part rules in `docs/CREATE-SPINE.md`. **Post-1.0:** keep the rhythm for live content updates (no world resets; freshness from curated updates) — `1.1` odd (feature/thunderdome) → `1.2` even (perf/balance) → …

### The cycle — order of operations (one odd/even pair)

Every pillar runs the same loop (formalized after v0.7.0 ran it informally). Steps 5–6 overlap by
design: the cycle is a pipeline, not a queue.

1. **Kickoff** (the odd milestone opens). Triage the milestone against the pillar's goal line — close
   the stale, re-home strays, file the gaps — until the milestone states its exit criterion. **Arm the
   instruments:** regenerate the recipe dump if recipes will be authored, record the connectivity
   baseline, switch the pillar's diagnostic configs ON (e.g. the gtmogs debug flags).
2. **Implementation.** Work the milestone on the version branch, folded into the integration PR; every
   change lands with its playtest items on that PR (CLAUDE.md law), and **patchnotes/docs are written
   with the change, never after**. Close-on-stage as work lands (close + `playtest-blocked`), so the
   milestone bar reads implementation progress. Ends in a **freeze**: milestone fully implemented, CI
   green, PATCHNOTES convention-styled, docs/wiki matching the data. After the freeze, only playtest
   fixes land.
3. **Smoke playtest** (the MVP gate). One human pass over the PR checklist on a fresh world + dedicated
   server: it boots, a world creates, every change exists and works. Failures: fix → re-freeze → re-run
   the failed items. Passes drop the `playtest-blocked` labels. Exit: checklist ticked, labels gone, PR
   un-drafted.
4. **Ship the feature patch.** Verify the version bump, merge; `release.yml` tags and publishes; deploy
   to the box.
5. **Soak** (the live playtest — which *is* the even milestone's kickoff). Real play on the shipped
   release, aimed at *finding work*: balance misses, bugs, perf, pacing feel. Diagnostics stay ON
   (debug logs, spark). Everything found files straight onto the even milestone as it's found; the
   measured passes live here (the cost-dial sweep, the pacing targets from `CREATE-SPINE.md` Part 1).
6. **Stabilization implementation** (parallel with the soak). Work the even milestone while play
   continues: balance, bugs, perf, the review-tail checklists — **no new mods or features** (cadence
   law). The *next* odd's design and triage may begin in the background — specs, issue filing, dumps —
   **paper only, no feature code**.
7. **Extended smoke, then ship the stab patch.** A longer verification pass — balance needs feel-time,
   not just existence checks — plus the standard smoke items, with **diagnostics switched OFF** for
   ship. Merge, release, deploy.
8. **Repeat.** The next odd's kickoff opens on a stabilized base.

**`v0.6.0 — Foundation` · even / stabilization** — *Goal: clean boot log, deterministic CI, a fast playtest loop, a perf baseline, modlist hygiene — decks cleared for the pillar work.* The initial cleanup, not a feature dump (already even-aligned).
- Quick wins / CI: #127 consolidate index-refresh workflows · #131 auto-regen the ground-truth digest · #154 config/ vs defaultconfigs audit. _Done in #166: #111 pin · #100 orphan-lib sweep · #126 docs cleanup + index. #105 prune branches → recurring at v0.7.0._
- Boot-log correctness: #119 dropped recipes · #120 mangled loot tables · #121 remaining ERROR noise.
- Perf baseline (server's up): #147 spark TPS routine + target · #151 Create contraption guardrails.
- Modlist hygiene (chartered by the #157 rubric — initial keep/cut cleanup, not a feature thunderdome): #157 define the rubric (`docs/CURATION.md`) · #83 umbrella · #106 structure allowlist · #108 spawn-gating review. _Done in #166: #60 umapyoi cut. Moved out: #107 Ender Moon (kept → flavor at v0.11.0), #110 textures (→ v1.0.0 living tracker)._
- Design charter: #218 ratify the structural model (`docs/SYSTEMS.md` — the Eco-style production/trade loop) + reconcile the docs to it; #157 curation rubric (`docs/CURATION.md`). _Done in #166: both landed + ratified; SoL → Diet (AppleSeed) food swap._
- Server's up: #98 pre-gen test world · #3 purple-arrows repro · #38 inventory-sort verify + close.

**`v0.7.0 — Create spine` · odd / feature** — *Goal: the recipe graph reads as one connected web — no functional-duplicate parts, materials unified, Create addons cohere; **the recipe/gating spine is built** — cost model applied (~3× basics, GregTech-style step-depth for tiers) and the KubeJS locking + progression-gating mechanism stood up with the Create-spine lock list. Everything builds on this, so it goes first.*
#219 recipe cost model (3× basics + step-depth) · #220 KubeJS locking + gating mechanism + first lock list · #103 duplicate parts (steel-plate family) · #101 limestones · #102 carbon/graphene · #113 connectivity islands · #132 Create-addon cohesion · #145 kinetic power tier ladder · #112 guns through sequenced assembly · #87 Create Stuff'N Additions · #17 recipes tracker (zagwar) · #105 prune merged branches (recurring hygiene sweep).
**Closes with the thunderdome** — the Create-related set reviewed by feature merit. _(Per `docs/SYSTEMS.md`, each feature pillar refines its content against the model and sets its locks as it's built.)_

**`v0.8.0 — Stabilization I` · even** — *Goal: profile, balance, and bug-fix everything 0.7 added; no new mods.* Spark the Create-spine feature dump (#147 routine), run the mod-conflict catalog (#160), perf-prune by the #157 perf gate, clear any boot-log regressions. Plus #129 — mob drops into the connectivity tool + rename (dev tooling pulled forward from the economy pillar; the weave/bounty enabler).

**`v0.9.0 — Survival` · odd / feature** — *Goal: the survival interlock (temperature × diet × seasons-driven scarcity) works in play — the world pushing back is real and tuned.* Pressure goes first by design: it generates the demand the later pillars serve (production exclusives at v0.11.0, distribution at v0.13.0), it's config-heavy rather than mod-risky (Cold Sweat / Serene Seasons / Diet shipped in v0.6.0), and it needs the longest tuning soak — a full seasons cycle is weeks of real server time.
Survival: #91 interlock (umbrella over #124/#99) · #124 altitude-driven cold (lands *before* the airship build-out, so ships contend with cold from day one) · #99 midnight-thoughts anchor · #152 author Cold Sweat world config · #155 death-penalty + gravestone tuning (settles PvPvE stakes before airship-combat balance at v0.13.0) · #238 perishability as continuous demand · #263 First Aid eval · #225 crop growth × season/temperature · #197 Cold Sweat heat-source range · #217 absolute darkness · #270 Terralith companion curation · #301 mining-hunger review · #25 open-parties-and-claims (the stakes infrastructure #155 tunes against) · #251 mod-review sweep (Survival core). _(Diet — AppleSeed replaced the SoL forks in v0.6.0; tune its groups/buffs here.)_ Flavor: #107 Ender Moon weave (kept; earn its place via a system tie).
**Closes with the thunderdome.**

**`v0.10.0 — Stabilization II` · even** — *Goal: profile, balance, and bug-fix everything 0.9 added; no new mods.*

**`v0.11.0 — Magic & MineColonies` · odd / feature** — *Goal: the two non-Create production routes completed, one owner each (`docs/SYSTEMS.md` §3 — they're the mirrored alternate routes, both with locked exclusives). **Magic (Xela112233):** the magic web balanced and bridged into one progression, **with its own locked exclusives** (a magic specialist is a needed trade partner — reachable also via the MineColonies route, mirroring Create). **MineColonies (zagwar):** the colony route stood up — the **colony lock list** (outputs only a colony can produce, per the incremental #220 pattern) and the **cheap-basics calibration** against the ~3× Create cost. Completes the supply side, so the economy pillar has specialists to distribute between.*
Magic (the #80 Arcana satellites): #122 post-merge balance · #123 flagship depth · #118 Soul-Reaping intent · #146 gate Ars/Iron's/Occultism behind Create + define magic's locked exclusives (per #220) · #235 End Remastered eyes ↔ boss drops. Colonies: #92 progression on-ramp (complex Create tech behind colony/boss progress) · #259 boss respawning (keeps the boss half of the high-tier fork reachable for late joiners). Sweep: #249 (Magic + MineColonies).
**Closes with the thunderdome.**

**`v0.12.0 — Stabilization III` · even** — *Goal: profile, balance, and bug-fix everything 0.11 added; no new mods.* Carries #309 — the scheduled assessment of splitting v0.13.0 into logistics + economy pillar versions.

**`v0.13.0 — Economy & logistics` · odd / feature** — *Goal: the **load-bearing** economy is playable end-to-end — coins valued, player-driven pricing, bounties pay, mob/structure inputs wired, the dynamic-pricing / player-minted-currency question resolved — **and its logistics/aeronautics movement arm is woven** (airship payoff + transport ladder), balanced against the scarcity lens (logistics must relieve regional scarcity without erasing it — `docs/SYSTEMS.md` §5).* The economy *is* the distribution stage and aeronautics is its physical arm, so they build together. It lands **last** by design: this pillar has come to encompass whole-pack balance (pricing, faucet/sink, export specialization), which needs the things being priced — magic/colony exclusives, transport costs — and the pressure that makes anyone buy to exist first. **Hedge the tooling risk early:** run #221 as a research spike during v0.9–v0.11 so this build starts with the pricing/currency tooling question answered (#129 already moved up to v0.8.0 on the same logic).
Economy: #221 economy tooling (dynamic pricing + player-minted currency) · #240 player-driven value model (retire vendor-sell + NPC coin faucets) · #138 shop/transaction mechanism · #136 coin tiers + price sheet · #137 Bountiful bounty pools · #150 faucet/sink integrity audit · #139 regional export specialization · #178 Companions coins vs Numismatics · #167 horse-betting loop · #209 gravel→iron rates + Ex Nihilo question · #298 finite excavation veins · #300 ore-gen review against the loop · #90 Numismatics + Trading Floor + Bountiful + wire inputs · #18 loot pass. Logistics/aeronautics: #280 Aeronautics keeps review · #43 flight-system interactions · #144 airship assembly/docking infra · #143 ship power/fuel economy · #142 audit Create machinery on ships · #141 transport/logistics ladder · #140 airship combat balance (Big Cannons) · #304 trains + logistics clusters audit/weave · #277 map-gating (Xaero's mid/late, Antique Atlas early — the travel-meaningfulness lens) · #125 Touhou-maids weave. Sweep: #247 (Economy/logistics/aeronautics).
**Closes with the thunderdome.** _(Split assessment scheduled — #309, on v0.12.0: if this is still ~20+ open there, split logistics (v0.13.0) from economy-value (v0.15.0) via the cadence's insert-a-pair rule, polish sliding to v0.17.0. Logistics goes first: transport cost is a pricing input.)_

**`v0.14.0 — Stabilization IV` · even** — *Goal: profile, balance, and bug-fix everything 0.13 added; no new mods.*

**`v0.15.0 — Polish & site + weave review` · odd / feature — LAST CALL** — *Goal: wiki, onboarding, QoL, claims, and the open decisions settled — plus the comprehensive weave review (all systems checked against each other — the loop seams) before the freeze.*
#115 player wiki · #148 in-game onboarding quest/guide · #153 JEI recipe-viewer cleanup pass · #77 in-site issue submission · #134 map server to new Cloudflare address · #70 WEB_DIR bind-mount · #13 Create: Harmonics eval · #1 shaders decision · #2 far-field render distance.
**Last call for major additions/integrations**; closes with the final thunderdome + the full weave review (the loop seams *between* the systems), then feature freeze.

**`v1.0.0 — Release` · even / perf patch [NFR]** — *Goal: feature-frozen — final performance + RAM tuning, ore-gen finalized, renamed, CI required → ship the public build.*
#149 v1.0.0 readiness checklist (meta-tracker) · #79 require CI on `main` · #81 early-game ore bootstrap · #318 ore-gen consolidated tracker (absorbed #116/#58) · #78 rename the pack (zagwar) · #48 GC/RAM (spark; server's up) · #110 custom-item textures (living tracker, needed-for-release).

**`📋 Backlog / living`** — *the holding pen so nothing is untriaged; pull into a version when it's next up.* #9 mod-ideas wishlist · #21 mod-removal suggestions · #40 quarterly port watch. Recurring passes that ride the cadence: #160 (mod-conflict catalog — runs on the even stabilization passes) · #161 (mod-update / changelog — feeds the next odd thunderdome).

**Recently closed/merged:** #75 (magic web, v0.5.2) · #82 (GTMOGS ore-gen, v0.5.1) · #88 (curation, v0.5.0) ·
#93 (ore-gen model decided → GTMOGS) · #62 (recipe overhaul, v0.4.6) · #51 (meadow cheese, in #62) ·
#117/#128/#130 (patch notes + site reframe) · #42 (TPS) · #73 (`.mrpack` drop) · #45 (meme cull) ·
#55/#59 (guns) · #56 (per-ore gen, since retired for GTMOGS) · #57/#63/#65 (digest/worldgen) ·
#61 (design docs) · #76/#85/#89 (CI + auto-resolver + workflow docs) · #86 (→ #88) · #133/#135/#156/#158 (docs sync + milestone map + triage playbook + curation-thunderdome scheduling).

## Cross-cutting notes
- **Verification gates everything** — headless can't validate KubeJS/worldgen; playtest per phase.
- **Network:** the sandbox reaches GitHub only (not Modrinth/CF); jar ground truth comes via `tools/mod-data/`.
- **Duplicate *items*** (lead/silver/tin/steel ×N) are auto-collapsed by `almost-unified` + `polymorph` —
  but **never unify `galosphere:*silver*`** (it's palladium, not silver); pin almost-unified (#111) to protect the config.
