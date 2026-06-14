# Governance Part 2 — `pcmc-realms` scope (issue #260)

> **Status: SCOPE, not accepted, nothing built.** This is the focused design scope for **Part 2 —
> Government** of the three-part governance plan: the mod **`theasshats/pcmc-realms`** (`pcmc_realms`).
> It stands on **Part 1** (`pcmc-territory`, the claims/resolution substrate — already spiked, see its
> PR #1) and is consumed by **Part 3** (`pcmc-mint`, the treasury/minting layer — later). It refines
> and supersedes the law-engine half (§4) of [`GOVERNANCE-MOD-SPEC.md`](GOVERNANCE-MOD-SPEC.md), which
> remains the whole-project (3-mod) reference. Per [`docs/CUSTOM-MODS.md`](CUSTOM-MODS.md) this mod lives
> in its **own repo** and reaches the pack via the mod-mirror packwiz pattern — it is **not** developed
> in this repo. This doc is the design reference that will seed that repo.
>
> **What this scope adds over the spec:** (1) the **hard-law / soft-law enforcement taxonomy** the
> maintainers asked for — the central reframe; (2) the **exact Part 1 API contract** Part 2 compiles
> against, read from the shipped Part 1 source; (3) the **MineColonies-guard soft-enforcement
> mechanism**; (4) the **territory-mod feature recommendations** for Part 2, folded in.

---

## 1. Where Part 2 sits

```
  pcmc-territory (Part 1)          pcmc-realms (Part 2)              pcmc-mint (Part 3)
  ─────────────────────           ────────────────────             ──────────────────
  "who governs this chunk?"  ──►   tiers + hierarchy + LAWS    ──►  treasury + minting
  entities, members, claims        (this doc)                       (registers TAX/STIPEND
  reverse-index + resolver         + a pluggable LawType registry    effects into Part 2)
```

Part 2 **owns**: the tier ladder and promotion; federation composition and the recursive sub-region
tree (parent/child links); member **roles for governance acts** (who may legislate/enforce); and the
**law engine** — the typed-law registry, the leaf→root precedence resolver, and the three enforcement
modes (§3). It introduces **no new claims and no new entity identity** — an entity's UUID, members,
colony IDs and claim keys all live in Part 1; Part 2 hangs *political* state off that same UUID in its
**own** `SavedData` (spec §8: "persistence partitions per mod"). Money movement (tax skim, stipend
deposit, fine collection) is **Part 3's** job — Part 2 defines the *law* and the *trigger*, Part 3
supplies the *coin* (§4, §6).

**Lower mods never reference upper ones.** `pcmc_territory` knows nothing of tiers, laws, or money;
`pcmc_realms` knows nothing of treasuries. The only coupling is the published APIs in §2 and the
`LawType` registry Part 2 publishes for Part 3.

---

## 2. The Part 1 contract Part 2 compiles against (read from shipped source)

Part 1's `:mod` exposes one public package — `com.theasshats.pcmcterritory.api` — and its javadoc is
explicit: *"Part 2 consumes only this package: it must not reach into `data`/`integration`/`core`
directly."* The surface, verbatim from the Part 1 PR branch:

### `TerritoryApi` (static facade)

| Method | Returns | Use in Part 2 |
| --- | --- | --- |
| `resolve(ServerLevel, ChunkPos)` | `List<UUID>` — governing chain, **leaf first** | The entry point for every position-local law check. |
| `getEntity(ServerLevel, UUID)` | `Optional<EntitySnapshot>` | Read an entity's members/colonies/claims by id. |
| `getEntityByName(ServerLevel, String)` | `Optional<EntitySnapshot>` | Command/UX lookups. |
| `toTerritoryChunk(ServerLevel, ChunkPos)` | `TerritoryChunk` | Helper. |

> **Critical clarification on `resolve` — Part 2 owns the ancestry walk.** In Part 1 the returned chain
> is length **0** (wilderness) or **1** (the single leaf entity governing that chunk; colony borders take
> strict precedence and *shield* their chunks, so there is at most one leaf). Part 1 has no concept of
> parents, so it can never return more than the leaf. **Part 2 builds the full `leaf → … → root` chain
> itself** by taking `resolve()[0]` (the leaf id) and walking its **own** `parentId` map. The Part 1
> javadoc phrase "Part 2 extends the chain" describes the *effective* chain the consumer sees — **not** a
> call back into territory. Do not try to push hierarchy down into `pcmc-territory`; that would violate
> the "lower mod never references upper" rule.

### `EntitySnapshot` (immutable record)

`record EntitySnapshot(UUID id, String name, Map<UUID,Role> members, Set<Integer> colonyIds, Set<ClaimKey> claimKeys)`

Point-in-time view; its own javadoc notes *"Part 2 entities will extend this record with additional
fields like tier, parent, and children."* We do that **by association** (a parallel `Map<UUID,RealmGov>`
in Part 2's `SavedData`), **not** by editing the record — the record stays Part 1's, stable.

### `TerritoryEvents` (on `NeoForge.EVENT_BUS`)

`Created`, `Removed`, `Changed` — each carries an `EntitySnapshot`. Part 2 **listens** instead of
polling. Load-bearing for Part 2:

- **`Removed`** → cascade cleanup: drop the entity's `RealmGov`, detach it from any parent's `childIds`,
  and decide the fate of a federation that loses a member (re-evaluate tier / dissolve if below the
  composition floor). This must be handled or Part 2's tree leaks dangling ids.
- **`Changed`** → membership/colony/claim edits; invalidate any cached jurisdiction footprint.

### `Role` (the governance permission ladder)

`enum Role { CITIZEN, OFFICER, LEADER }` with `atLeast(Role)` (ordinal-ordered). Part 2's command
permissions gate on this (§7). **API-hygiene note for the territory mod:** `Role` lives in the `core`
package but is surfaced through `api.EntitySnapshot.members()`, so Part 2 must import a `core` type to
read it — a small leak across the "api only" boundary. **Recommendation back to Part 1: promote `Role`
into the `api` package** (or re-export it) so the contract is self-contained. Track as a territory issue.

### Versions / coordinates Part 2 builds against

Part 1 declares soft-deps **MineColonies `1.1.1327`** and **OPAC `0.26.2`** (snapshot builds — minor API
drift between bumps is possible; pin the dev coordinates and re-verify on each MineColonies bump). Part 2
hard-depends on `pcmc_territory`; soft-deps MineColonies (for guard enforcement, §5) and, in Part 3, Numismatics.

---

## 3. The law model — **hard vs soft** (the central reframe)

The spec's §4 modelled *every* law as a NeoForge event that **cancels the action** (PVP → cancel the
damage, build → cancel the break). The maintainers' direction reframes this: **what the mod may
automate, and what it must leave to people (and guards), are different things.** A law carries an
**enforcement mode**, and the mode — not the law type — decides whether the mod *acts* or merely
*signals*.

### The three enforcement modes

| Mode | Who acts | What the mod does | Maintainer's examples |
| --- | --- | --- | --- |
| **`AUTOMATIC`** (hard) | the mod, no human | Applies a mechanical effect directly and reliably. Reserved for **money/state the government legitimately controls** — never for blocking a player's *physical* action. | **Tax**: a cut comes out of a transaction automatically. **Stipend**: a weekly citizen payment is deposited automatically. |
| **`OFFICER`** (hybrid) | a human officer triggers; the mod settles | The *judgment* is a person's; the *settlement* is automatic. | **Fine**: a designated officer enacts the fine (the mod does not auto-detect the crime), **but the withdrawal of funds is automatic** once enacted. |
| **`SOCIAL`** (soft) | the player base / MineColonies guards | The mod **never blocks the action**. It records the violation and raises a **consequence signal** (guard aggression, "wanted" status, optional bounty). **If no enforcer is present, nothing happens** — the offender gets away with it. A matching act counts only if it **fails an exception check** first (e.g. self-defense; §5). | **PvP ban**: making PvP illegal does **not** stop me killing another player — instead MineColonies guards attack me; if none are near, I get away with it. |

### Why this split (the design rationale)

- **Hard-blocking behavioral actions feels bad and is often unfair or unenforceable.** Cancelling a
  PvP hit, a block break, or item use silently overrides player agency and creates "why didn't that
  work?" confusion. Worse, it can't capture intent or context. So **physical/behavioral laws are
  `SOCIAL`**: the world enforces them, and enforcement is *emergent* — present where there are guards or
  willing players, absent where there aren't. That **absence is a feature**: a lawless frontier is a
  real place, and "is anyone actually enforcing this?" becomes a meaningful question.
- **Money the government controls is fair to automate.** A polity collecting tax on a sale, paying its
  citizens a stipend, or collecting a fine an officer has *already judged* — these are legitimate,
  predictable state functions. Automating them is reliable and removes book-keeping toil. So **fiscal
  effects are `AUTOMATIC`**.
- **The fine is deliberately hybrid** because it bridges the two: the *crime* (a soft-law violation) is
  judged by a person (no false-positive auto-fining), but the *penalty* (coin) is exactly the kind of
  state money-movement that should be automatic and inescapable. Soft trigger, hard settlement.

This is the rule of thumb for slotting any future law: **does it move government money, or does it
restrain a player's action?** Money → `AUTOMATIC` (or `OFFICER` if it needs a human verdict first);
action → `SOCIAL`.

---

## 4. Law catalog

Each law is a typed rule resolved leaf→root by tier (spec §4 resolver: highest tier rank wins; a
subordinate fills gaps its parent leaves unset; documented same-rank tie-break). What changes here is
the **enforcement mode** column and the new fiscal/penalty types.

| Law / act | Mode | Mechanism (Part 2 unless noted) | Lands in | MVP? |
| --- | --- | --- | --- | --- |
| **PVP** (`ALLOW`/`DENY`) | `SOCIAL` | On player→player damage in a `DENY` jurisdiction: **do not cancel**; mark the *unprovoked* attacker "wanted" (self-defense exempt, §5) and trigger guard aggression. No guards → no effect. | 2a | ✅ |
| **TRESPASS / non-member entry** | `SOCIAL` | Non-member in a restricted jurisdiction → guard aggression / wanted. (Hard *block* of build/use is already the claim mods' job — don't duplicate it; see note below.) | 2a | optional |
| **CONTRABAND / ITEM_BAN** (tag/item set) | `SOCIAL` | Possessing/using a banned item in-region flags a violation → guard aggression. (A `customs`-style border confiscation is a later `AUTOMATIC` variant.) | 2b | later |
| **CURFEW** (tick window) | `SOCIAL` | Being out in-region during the window → wanted/guard aggro. TPS-sensitive; defer. | later | later |
| **SECESSION** (`ALLOW`/`DENY`) | `SOCIAL` | Federation-scope. Leaving is **never hard-blocked** (at-will, §8.1); on `DENY`, secession lets the abandoned faction's top officer declare the seceder (+ its new host) hostile. | 2b | later |
| **TAX** (rate, bps) | `AUTOMATIC` | On a taxable transaction in-region, skim `rate` into the entity treasury. Part 2 defines the type + trigger point; **Part 3 moves the coin** (Numismatics). Trigger hook is the shared risk (§9). | 3 (type stubbed in 2) | later |
| **STIPEND** (amount, period) | `AUTOMATIC` | On a recurring schedule (§6), pay `amount` from the treasury to each citizen. Part 2 owns the **scheduler + member iteration**; **Part 3 moves the coin**. Skips when underfunded. | 3 (type stubbed in 2) | later |
| **FINE** (amount, target, reason) | `OFFICER` | An OFFICER+ runs `/realm fine <player> <amount> [reason]`; the mod **auto-withdraws** from the target's account (Part 3) into the treasury. Not a standing rule — an *enforcement action*, usually issued *for* a soft-law violation. | 3 (command in 2) | later |
| **TIER-GATED CLAIM ALLOWANCE** | `AUTOMATIC` | Government grants citizens an OPAC claim pool that scales with tier; leaders hold a larger grantable pool. From the territory feature recs (§7). Depends on an OPAC capability spike. | 2b/later | later |

**Don't re-implement claim protection.** OPAC and MineColonies *already* hard-protect their claims
(only permitted players build/break inside). A `BLOCK_INTERACT = DENY` "hard cancel" law would duplicate
that. So Part 2 does **not** ship a hard build/break law; griefing-as-a-*crime* (consequences for a
permitted-but-unwelcome act, or acts in unprotected land) is the `SOCIAL` TRESPASS/PVP path instead.
This keeps Part 2 a *governance* layer, not a third protection system — the same discipline that keeps
it from being a third claim system.

**MVP (slice 2a) law set:** PVP (`SOCIAL`) is the one enforceable law that proves the whole model end to
end (resolver + enforcement-mode dispatch + guard integration) without needing Part 3's money. TAX/
STIPEND/FINE **types and commands are stubbed** in 2a (registered into the `LawType` registry, no-op or
"requires pcmc-mint" until Part 3 wires the coin) so the registry contract is exercised early.

---

## 5. Soft-law enforcement via MineColonies guards (the new technical lynchpin)

This is the mechanism that makes `SOCIAL` laws real, and the **highest-risk new dependency**. The
maintainers' target behavior: *"if I make a law that makes PvP illegal that shouldn't mean I can't kill
another player; instead MineColonies guards should attack me, or if none are near I should get away with
it."*

### A violation must clear an **exception check** first (self-defense, and others)

A soft law that flags *every* matching act punishes the wrong player — most obviously, **defending
yourself**: if someone attacks you, retaliating must not make *you* wanted. So the SOCIAL pipeline has
three steps, not two: **detect a candidate act → run a justification check → only an *unjustified* act
raises the wanted signal** (and the guard aggro below). A justified act is a full no-op — no wanted
status, no guards, no bounty.

**Self-defense (the primary built-in, MVP).** The rule is **the first striker is the violator.** When A
attacks B unprovoked in a no-PvP jurisdiction, A is flagged; when B hits back, B is *not* flagged,
because A is B's recent aggressor. Mechanism: a short-lived per-pair **aggression record** — on every
player→player hit note `(attacker, target, tick)` (vanilla already tracks revenge via
`LivingEntity#getLastHurtByMob` + timestamp; a small `Map` keyed by attacker suffices). On a candidate
hit, if the *target* damaged the *attacker* within a config **combat window** (default ~15s), the hit is
retaliation → **suppress**. The original aggressor stays flagged from their own first strike; the
defender never is.

Further justifications, in rough priority (keep MVP to self-defense; the rest are cheap once the check
exists):

- **The target is already wanted/an outlaw → open season.** Attacking a player currently wanted in this
  jurisdiction is *sanctioned*, not a crime — the attacker isn't flagged. Citizens become effective
  deputies, and it pairs naturally with bounties (hunting an outlaw is lawful). Strong MVP candidate —
  same wanted-table lookup, no new data.
- **Defense of another.** Striking someone who is *currently attacking a fellow member/citizen* extends
  self-defense to third parties. A little more state (who's attacking whom) → 2b+.
- **Consented combat (duels / arenas).** Mutual opt-in or a jurisdiction-designated arena zone exempts
  both parties. The pack already ships **`opacpvp` (OPAC PvP Support, v1.0.1)** — a per-player
  mutual-consent PvP toggle for *party* members (`/opacpvp`; both must enable, and it *hard-blocks* member
  friendly-fire by default). A ready-made consent surface to read for this exception, and it works on the
  party-consent axis — orthogonal to the territorial PvP law, so the two compose rather than conflict
  (member friendly-fire is opacpvp's call; outsiders in-jurisdiction are the governance law's). Later, not MVP.

**Generalize the check, not just the cases.** Model it as a pluggable **`Justification`** step the
SOCIAL evaluator runs for *any* law — so CONTRABAND can later exempt in-transit goods, TRESPASS a player
fleeing combat, etc., without re-plumbing each law. Self-defense is simply the first registered
justification. **Risk is low:** this is vanilla combat bookkeeping (a small per-player map + existing
revenge timestamps), not a new external dependency like the guard integration it gates.

### Mechanism (leading hypothesis — verify in the spike)

MineColonies guards already attack players natively via the colony **permission system**: each colony
ranks players, and a rank with the **"Fight Guards" permission** (the **Hostile** rank by default) is
attacked by that colony's guards on sight. Confirmed on two points: the **rank setter is public API** —
`IPermissions#setPlayerRank(UUID, Rank, Level)` (plus `add(UUID, name, Rank)`) in
`com.minecolonies.api.colony.permissions`, the surface addon developers use — and guard targeting
**natively attacks pvp-hostile players** (`AbstractEntityAIGuard` target acquisition; upstream specifically
fixed "guards not targeting pvp-hostile players"). Guards carry a ~30s combat timer (`COMBAT_TIME = 30*20t`)
and a per-guard attack list. *(Exact signatures still to be confirmed against the pinned snapshot in the
spike — the pack runs MineColonies snapshots.)*

So Part 2's `SOCIAL` enforcement is: **on a violation in a colony-governed jurisdiction, set the
offender's MineColonies rank in that colony to a Fight-Guards rank for a time-boxed "wanted" window.**
The guards do the rest natively — no custom combat AI, and the "no guards nearby → get away with it"
behaviour falls out for free (a Hostile rank with no guards in range simply does nothing). When the
wanted timer expires, restore the prior rank.

### Spike (do this first, before any enforcement code)

1. **Confirm the rank API on the pinned snapshot.** `IPermissions#setPlayerRank(UUID, Rank, Level)` is the
   target (public `api.colony.permissions`); verify its exact signature on `1.1.1327` and compile a real
   call in a dev env, exactly as Part 1 spiked its lookups. (Known to exist; snapshots can shift it.)
2. **Does setting it actually make guards aggro a player in-game?** Static compile ≠ runtime — this needs
   the box (a guard, a Hostile-ranked test player). Add it to the Part 2 playtest checklist.
3. **Restore semantics.** Confirm we can read the prior rank and put it back; make wanted-state
   all-or-nothing (don't strand a player as Hostile if the server stops mid-window).

### Fallbacks if the rank API isn't usable

- **Custom target-goal injection** into nearby guard entities (add the offender to a guard's target
  list directly) — heavier, more brittle across MineColonies bumps.
- **Decouple via a "wanted" signal**: Part 2 maintains the wanted table and *emits* it (an event /
  scoreboard tag / capability); guard aggression, bounties (economy pillar), and `pcmc-killfeed` overlays
  subscribe. This is the most robust shape regardless of the guard API and is probably worth building
  **anyway** as the canonical soft-law output, with the rank-flip as one consumer.

### Delivering the guard integration: MineColonies addon vs. fork (decision aid)

The maintainers are weighing a **MineColonies addon** (a separate mod calling MineColonies' public API)
against **forking MineColonies** (shipping a modified MineColonies in place of upstream). The framework for
the build instance to decide:

**What soft-law guard enforcement actually needs:** (1) make a specific player attackable by a specific
colony's guards on demand; (2) scope it to that jurisdiction and to *nearby* guards, with "none near →
nothing happens"; (3) time-box it and restore.

**The rank/permission path covers all three through the *public API* — so an addon is sufficient for the
core mechanism:**
- `IPermissions#setPlayerRank(...)` flips the offender to a Fight-Guards (Hostile) rank → guards attack
  natively (req 1).
- Rank is **per-colony**, so aggression auto-scopes to the jurisdiction; guard **vision range** gives
  locality and "no guards → get away with it" for free (req 2).
- Read the prior rank, restore it when the window ends (req 3).
- **Prior art:** *MineColonies: War 'N Taxes* is an existing **addon** adding taxes + war to colonies —
  evidence the public surface carries this class of feature without a fork.

**What would justify a fork — only guard *behaviour the binary rank can't express*:** graduated response
(warn → fine → attack), non-lethal/subdue or **jailing**, guards **pursuing past the colony border**, or
targeting keyed to offense severity / a bounty rather than the on/off Hostile rank. Target acquisition
lives in `AbstractEntityAIGuard` (`checkForTarget`/`decide`) with **no public hook to inject a target**, so
reshaping *how* guards choose targets means editing MineColonies source. (A fork does **not** create guards
where there are none — OPAC-only land and ship-realms still have no guards aboard.)

**A fork's cost is high and pack-wide:** it **replaces upstream MineColonies for every player** (blast
radius = the whole pack, vs. an additive soft-dep that breaks nothing if removed); the pack runs
MineColonies **snapshots** (`1.1.1327`), so a fork must **continuously rebase** on upstream or freeze and
forfeit updates — a standing tax this repo's snapshot-drift notes already call out; and it's two custom
mods to maintain instead of one.

**Recommendation: addon-first.** Build the integration as MineColonies-API calls inside `pcmc-realms` (or a
thin companion), gated behind the §5 spike. Treat a fork as a **later, separate** decision triggered only by
a concrete need for custom guard *behaviour* the rank system provably can't deliver — and budget it as
*replacing a core pack dependency*, not an addon-sized task. Decide against that bar.

### Edge cases to design for

- **OPAC-only jurisdictions have no guards.** A realm governing land via bound OPAC claims (no colony)
  has nothing to aggro — `SOCIAL` enforcement there is *purely* player-driven (wanted status + bounty,
  no NPC teeth). That's acceptable and on-theme; document it so it isn't read as a bug.
- **Airships are their own jurisdiction (and extraterritorial)** — maintainer decision: a claimed Create
  Aeronautics ship is a city/outpost-tier entity that stays its own faction's soil even inside another
  faction's claims (an embassy). Its own design point — see §7.1. For enforcement: a ship has **no
  MineColonies guards aboard**, so `SOCIAL` enforcement on a ship is player/bounty-driven (like the
  OPAC-only case above); hard laws (tax/stipend/fine) and the wanted-signal still apply aboard.
- **Wanted decay & escalation.** Single window length in config for MVP; later, repeat offences could
  escalate (longer window, auto-bounty). Keep MVP simple.

---

## 6. Recurring hard laws need a scheduler (stipend)

Most `AUTOMATIC`/`OFFICER` effects are **event-driven** (tax skims on a transaction, a fine on a command)
— cheap and position-local, matching the `SYSTEMS.md` §3a perf doctrine. **STIPEND is the exception: it's
*time*-driven.** Part 2 needs a **coarse scheduler**:

- Fire on a long period (config; default ~ one in-game week, i.e. every *N* ticks), **never per-tick**.
- On fire, for each entity with a STIPEND law: iterate its citizens (`EntitySnapshot.members`), pay
  `amount` each from the treasury (Part 3), **stop/skip when the treasury can't cover** the run (no debt,
  no negative balances) and log it.
- Stagger across entities if many exist, so a payout tick doesn't spike. Persist the "last paid" tick in
  Part 2's `SavedData` so a restart doesn't double-pay or skip.

The scheduler is Part 2; the per-citizen credit is Part 3. This is the only always-resident recurring
cost Part 2 adds, and it's coarse by construction.

---

## 7. Territory-mod feature recommendations, folded in

The Part 1 README/PR explicitly hands these to Part 2 (the maintainer asked that they be considered):

1. **Combined city/faction borders.** A realm's jurisdiction is the **union** of its `colonyIds`
   footprint **and** its bound OPAC `claimKeys` — colony borders take strict precedence and *shield*
   their chunks from outside claims (Part 1 already implements the shielding). Part 2's law resolution
   therefore covers both colony land and bound-claim land under one realm. **No new work to enable** —
   it's how `resolve()` already behaves; Part 2 just treats the leaf uniformly.
2. **Tier-gated claim allowances** (a *governance* feature, in the catalog as `AUTOMATIC`). Citizens
   start with a small OPAC claim pool; advancing a realm's tier (city growth) raises the per-citizen
   allowance; leaders hold a larger pool to **grant**. This makes the tier ladder *mean something
   tangible* (land budget) and weaves governance into the claim economy. Sits in 2b/later. **Formalized by
   the charter system (§7.2)** — a charter is the grant vehicle; a bare allowance is its simplest form.
3. **Hard-limiting claims (open).** Enforcing those allowances (preventing over-allowance claims)
   depends on OPAC exposing **either** a *settable per-player claim limit* **or** a *cancellable
   claim-creation event* — neither verified yet. This is a **shared spike with the territory mod**
   (its open question too). If neither exists, allowances degrade to *soft* (a `SOCIAL` "illegal
   over-claim" rather than a hard cap) — which fits the taxonomy fine.

### 7.1 Airships are their own jurisdiction (city/outpost-tier) — and extraterritorial (embassies)

**Maintainer decisions:** (a) a claimed Create Aeronautics airship is **its own city/outpost** for law — a
first-class jurisdiction, not terrain it flies over; and (b) a ship inside another faction's claims **stays
its own (primary) faction's soil — an embassy** (extraterritoriality).

**Substrate.** The pack ships **`aeroclaims` (Create Aeronautics: Claims, v0.9.0, Modrinth `CwZ8q37q`)** —
claim a Create Aeronautics **sub-level** (the ship) via a claim block + GUI. Separate from OPAC chunk claims
but rides OPAC's **parties/allies** and shared **claim budget**, transferable both ways
(`/aeroclaims transfer to opac|aero`). So a ship already carries an owner/party — the hook a realm binds to.

**Model.** A ship-realm is a **municipality-kind entity bound to an aeroclaims sub-level** rather than a
colony — a *mobile outpost*. It holds laws, promotes, and can join a federation (a flagship as a member of
its nation) like any municipality; its footprint is the ship, not chunks. **Founding a ship-realm requires
an `AIRSHIP` charter from a high-tier entity (§7.2)** — anyone may aeroclaim a ship for protection, but a
charter is what turns it into a recognized outpost.

**Embassy / extraterritoriality — the precedence rule.** For an entity **aboard a claimed sub-level the
ship-realm governs and overrides whatever territory the ship is over or inside** — its primary faction's
laws apply on deck even within another faction's colony/claim, and the leaf→root walk follows the *ship's*
hierarchy, not the host's. Resolution checks **"on a claimed sub-level?" first**; only off-ship does it fall
to chunk resolution. The general principle: **a realm's own claims are sovereign enclaves (embassies),
governed by their owner, not the host they sit inside.** *(Open: Part 1 currently has colony borders
strict-shield their chunks, so a chunk-level enclave inside a colony is hidden today. The embassy rule is
unambiguous for **sub-levels** (ship over ground); whether it also overturns colony-shielding for **chunk**
enclaves is a precedence call to settle with Part 1 — §12.)*

**This is territory's job (issues to file, §13).** `TerritoryApi.resolve(level, chunkPos)` is **chunk-based**
and won't see a sub-level claim — so both the **resolution source** (map entity → sub-level → owner → realm)
and the **embassy precedence** (sub-level wins over the chunk beneath) are **`pcmc-territory` changes** (the
resolver owns "who governs this position"). Captured as territory issues in §13 (this session can't write to
that repo). Part 2 doesn't block on it — until it lands, ship governance simply can't resolve.

**Enforcement aboard.** Hard laws (tax/stipend/fine) and the wanted-signal work on a ship-realm normally.
**Soft enforcement has no native teeth aboard** — no MineColonies guards ride a ship — so on-ship `SOCIAL`
enforcement is player/bounty-driven, like an OPAC-only jurisdiction (§5). (Ship-mounted turrets/cannons as
"guards" is a speculative later weave, not scoped.)

### 7.2 Charters — high-tier entities license special claims and enterprises

**Maintainer decision:** special claims/enterprises — **airship claims (§7.1), mines, farms,** and the like
— are **chartered**: a player or lower entity may establish one only under a **charter** granted by a
**high-tier political entity**. Chartering is a *power of scale*, like minting.

- **Who may issue.** Config-gated; default **CITY-tier municipalities and all federations (COUNTY+)**.
  Settlements/villages/towns can't charter — chartering (with minting) is what high tiers are *for*,
  reinforcing the climb. (Minting stays federation-only; chartering is deliberately a touch broader to
  include high cities, per the maintainer.)
- **What a charter is.** A typed grant record `{issuer, holder (player|entity), type, terms, obligations,
  revocable}`.
  - **Type** — an **extensible registry** (mirrors the `LawType` pattern): `AIRSHIP` (a sub-level outpost →
    a ship-realm, §7.1), `MINE`, `FARM`, … (`TRADING_POST`/`FACTORY` later). MVP of the feature: the
    framework + `AIRSHIP` + one or two of `MINE`/`FARM`.
  - **Terms** — allowance/count, optional area or duration, revocability.
  - **Obligations** — a chartered enterprise typically owes the issuer **tax/tribute** (`AUTOMATIC` fiscal
    hook, Part 3): the **value-loop weave** — charters turn high-tier political power into a revenue stream,
    funnelling mine/farm/airship output up into the issuer's treasury.
- **Subsumes "tier-gated claim allowances" (§7 item 2).** A charter is the formal grant vehicle; the bare
  allowance grant is just its simplest form.
- **Layer, not a claim system (enforcement).** The claim mods still own *protection* (anyone may aeroclaim
  a ship or fence a field). The charter gates **political/economic recognition + special rights** —
  founding a **ship-realm**, or registering a **recognized mine/farm enterprise** (allowance + tax terms) —
  a Part 2 registration act, so it's **hard at the founding layer** (the mod refuses an unchartered
  founding). *Hard-gating the raw claim itself* depends on the **OPAC/aeroclaims claim-limit API** (the same
  open spike as §7 item 3); without it, running an **unchartered** enterprise is a **`SOCIAL` "unlicensed"
  violation** (wanted/fineable) — which fits the taxonomy.
- **Scope:** **2b/later** — charters need the tier ladder + federations to exist first.

**Design anchor.** A strong second-weave: **Economy** (governance issues value-bearing licenses; tax flows
up) × **Survival/production** (mines, farms) × **logistics/aeronautics** (airships) — it makes high tiers
economically meaningful and gives the production pillars a political dimension.

---

## 8. Tiers, promotion, federation (from spec §3, unchanged — recap)

Carried from `GOVERNANCE-MOD-SPEC.md` §3; no change, recapped so this doc stands alone:

- **Tiers** (rank in parens): municipality `SETTLEMENT(0) → VILLAGE(1) → TOWN(2) → CITY(3)`; federation
  `COUNTY(4) → KINGDOM(5) → EMPIRE(6)`. Stored in Part 2's `RealmGov`, keyed by the territory entity UUID.
- **Promotion is player-initiated and validated on the command** (no per-tick scanning). Inputs are free
  MineColonies metrics read via `EntitySnapshot.colonyIds` → MineColonies API: population (citizen count),
  footprint (claimed chunks), development (Town Hall level). Thresholds in `realms-server.toml`.
- **Federations are composed, not promoted** — `/realm federate` needs ≥ N (config, default 2) consenting
  members of ≥ TOWN. **Recursive `carve`** lets a KINGDOM+ spawn a child one rank below.
- **An entity may bind an `aeroclaims` sub-level (a mobile outpost) instead of a colony** (§7.1) —
  ship-realms sit at the municipality kind and promote/federate like any other; their footprint metric is
  the ship, not claimed chunks.

### 8.1 Secession is at-will — forced confederation & inter-entity hostility

**Principle: confederations are voluntary and at-will.** A member may **leave its parent at any time**
(`/realm secede`). A federation is not a prison — and, mirroring the rule that a law never *cancels* an
action (§3), **a federation cannot hard-block secession even by passing a law against it.** The *only*
thing that truly holds an entity in is a **server-admin lock** (Exception 1). Everything a faction itself
can do makes leaving *costly*, not impossible.

**Exception 1 — admin-forced confederation (the one hard lock; real-world time).** OP admins can force an
entity into a confederation and **lock it from leaving for an IRL (wall-clock) duration** — e.g. to enforce
peace terms when a war's loser won't honor them. Persist a `lockedUntil` **epoch timestamp** (real time,
not game ticks — it elapses across restarts and while the server is offline); `/realm secede` is refused
until it passes, then at-will resumes. `/realm admin confederate <entity> <parent> <duration>`;
`/realm admin unlock <entity>` releases early. Default duration in config.

**Exception 2 — a "no-secession" law (SOCIAL: leaving is allowed but is an act of war).** A federation may
pass `SECESSION = DENY`. Per the taxonomy this **does not block** the leave (at-will holds); instead, on
secession the abandoned faction's **top officer gets an OFFICER-mode decision** to **declare the seceder
hostile** — and, if the seceder **joins another faction, that host faction too**. Resolved via a pending
decision (`/realm hostility decide …`; clickable prompt / GUI later); the officer may also **ignore** it
(leaving was lawful by default).

- **Whose officer decides — OPEN (maintainer-flagged).** "The highest-ranking officer in the whole faction
  *or* the sub-region, depending on what type of entity leaves" — to settle during the build. The fork: a
  **direct member** leaving → the **root federation's** top officer; a **carved sub-region** leaving → its
  **immediate parent's** top officer (the entity actually abandoned), perhaps escalating to the root.
  *(Leaning: the immediate parent's top officer — it's the entity directly harmed — but recorded as open.)*

**Inter-entity hostility ("war") = the entity-scale analog of player "wanted" (§5).** Declaring entity B
hostile means **B's members are auto-wanted within A's territory** (guards aggro them via §5) and **PvP
between A and B is lawful** (the "outlaw → open season" exception, §5, applied faction-wide). Stored as
directional hostile pairs on the entity in Part 2. The secession decision is its **first trigger**; a
general `/realm war declare / sue-for-peace` surface is a natural later extension (open question, §12).

**Scope:** 2b — it needs federations + the hierarchy, and it **reuses §5's wanted machinery** (no new
enforcement path). At-will `/realm secede` itself (minus the war consequences) can land as soon as
federations exist.

---

## 9. Commands (extends Part 1's `/realm` tree)

Part 1 already ships `/realm found | info | whogoverns | debug bindclaim`. Part 2 adds (all
`Role`-gated via `atLeast`):

```
/realm promote                         (LEADER; attempt next municipality tier — §8)
/realm federate <name>                 (LEADER; compose a federation from consenting members)
/realm carve <name> <tier>             (LEADER; recursive sub-region, KINGDOM+)
/realm secede                          (LEADER; leave your parent — at-will unless admin-locked, §8.1)
/realm member <add|remove|role> ...    (LEADER; manage members/roles)
/realm law set <type> <value>          (LEADER; issue/replace a law)
/realm law unset <type>                (LEADER; clear a law → falls back to parent/vanilla)
/realm law list                        (CITIZEN; what's in force here, with the deciding tier)
/realm fine <player> <amount> [reason] (OFFICER; OFFICER-settled penalty — auto-withdraw via Part 3)
/realm wanted [player]                 (CITIZEN; view active wanted status in this jurisdiction)
/realm hostility decide <entity> <hostile|ignore>  (OFFICER; resolve a secession war decision — §8.1)
/realm hostility list                  (CITIZEN; entities at war with this realm)
/realm charter grant <type> <holder>   (LEADER, CITY+/federation; license a mine/farm/airship — §7.2)
/realm charter list                    (CITIZEN; charters issued by / held in this realm)
/realm charter revoke <id>             (LEADER, issuer; revoke a charter)
/realm mint <amount>                   (LEADER; Part 3 — federation+, against reserves)

/realm admin confederate <entity> <parent> <duration>  (OP; force into a confederation, IRL leave-lock — §8.1)
/realm admin unlock <entity>           (OP; release an admin leave-lock early)
```

GUI (entity/law/treasury screens, `/realm map` border render) is **post-MVP polish** (spec §"Cross-cutting
polish") — commands are the MVP so the system is playable before any screen work.

---

## 10. Slices, gates, and the public surface for Part 3

Aligned with spec §8 Part 2, made concrete:

| Slice | Delivers | Gate (in-game) |
| --- | --- | --- |
| **2a — single-entity laws (MVP)** | Tier field + `/realm promote` validation; the **law engine + leaf→root resolver**; the **three enforcement modes** with **PVP (`SOCIAL`) live** end-to-end (guard integration, §5); TAX/STIPEND/FINE **types + commands stubbed** into the `LawType` registry (no coin yet). No federation. | A single colony bans PvP; a guard attacks an offender inside it; an offender escapes where no guard stands. Independently shippable. |
| **2b — hierarchy** | Federation composition (consent, capital), recursive `carve`, parent/child links, governance roles; the full leaf→root cascade + same-rank tie-break; tier-gated claim allowances (pending the OPAC spike, §7). | Empire > kingdom > city precedence verified: a higher tier's PvP ban overrides a subordinate's `ALLOW`. |

**Public surface Part 2 publishes for Part 3:** a **pluggable `LawType` registry** plus an **effect
hook** so Part 3 registers TAX/STIPEND/FINE *effects* (the money movement) without Part 2 knowing
anything about coins — Part 2 owns the *rule, scope, precedence, schedule and trigger*; Part 3 owns the
*credit/debit*. Version it deliberately (breaking the registry API = major bump; Part 3 pins a range),
exactly as Part 1↔Part 2.

---

## 11. Spikes & risks (ranked)

1. **MineColonies guard targeting (the `SOCIAL` lynchpin, §5).** Highest risk because the whole soft-law
   model leans on it. Spike the rank API *and* runtime guard aggro before building enforcement. Have the
   "wanted-signal" fallback ready.
2. **The TAX transaction hook (shared with Part 3).** A desk check found **no public sale event** on
   Create: Trading Floor — assume a mixin until proven otherwise. Part 2 can define the TAX *type* now;
   the *trigger* is a Part 3 spike. Don't block 2a on it.
3. **OPAC claim-limit / cancellable-claim capability (§7).** Gates *hard* tier-gated allowances. If
   absent, degrade to a `SOCIAL` over-claim law. Shared spike with the territory mod.
4. **Thread-safety.** Part 1's `TerritoryResolver` is **server-thread-only** (it dropped
   `ConcurrentHashMap`; see territory issue #3). Part 2's law handlers run **on the server thread**
   (combat/interaction/command events) so they're in-bounds — but **never call `TerritoryApi.resolve`
   off-thread** (no async, no render-thread queries). State this as a Part 2 invariant.
5. **MineColonies snapshot drift.** The pack runs MineColonies snapshots; the rank API and metric
   lookups can shift between bumps. Pin dev coordinates; re-verify §5 and §8 on each MineColonies bump.

---

## 12. Open questions for the maintainers

- [ ] **Wanted-status canonical form:** build the decoupled "wanted signal" (event/tag) as the primary
      soft-law output (with the guard rank-flip as one consumer), or wire the rank-flip directly for MVP
      and generalize later? *(Recommendation: build the signal — it's robust to the guard API and feeds
      bounties + killfeed too.)*
- [ ] **Bounties on soft-law violation:** auto-post a Numismatics bounty when a player goes "wanted"
      (ties the soft law into the economy value-loop), or keep bounties a separate manual system for now?
- [ ] **TRESPASS / CONTRABAND / CURFEW in MVP scope, or PVP-only for 2a?** *(Recommendation: PVP-only for
      2a; it proves the whole pipeline. Add the rest in 2b once the mode dispatch is solid.)*
- [ ] **Which exceptions ship in MVP (§5), and the combat-window length?** Self-defense is non-negotiable;
      is "attacking an already-wanted outlaw is open season" also MVP (recommended — same lookup), or do
      defense-of-another / consented duels wait for 2b? *(Recommendation: self-defense + outlaw-open-season
      in 2a; default combat window ~15s, config.)*
- [ ] **Stipend funding when the treasury is dry:** skip the run (recommended), pay partial (first-come),
      or accrue arrears? *(Recommendation: skip + log; no debt.)*
- [ ] **Fine target scope:** can an officer fine *any* player who offended in-jurisdiction, or only
      members? Can a player be fined while not present/online (auto-withdraw on next login)?
- [ ] **Should `Role` be promoted into Part 1's `api` package** (§2) so Part 2's contract is
      self-contained? *(Recommendation: yes — file it on the territory mod.)*
- [ ] **Embassy vs. colony-shielding precedence (§7.1):** sub-level (ship) claims clearly override the
      ground they sit on; should a realm's **chunk** enclave inside another faction's colony also be an
      embassy, overturning Part 1's current colony-shielding for that case? *(Recommendation: ships-only
      embassy for MVP; revisit chunk enclaves with Part 1.)*
- [ ] **On-ship soft enforcement (§7.1):** no guards ride a ship — is player/bounty enforcement enough for
      ship jurisdictions, or is a "turret/cannon-as-guard" weave worth a later spike?
- [ ] **Guard integration delivery (§5):** addon vs. fork — recommendation is addon-first; the build
      instance makes the final call against the §5 decision aid once the spike confirms `setPlayerRank`.
- [ ] **Charter issuer threshold (§7.2):** CITY-tier municipalities + all federations (recommended), or
      restrict to federations like minting, or open it to TOWN+?
- [ ] **Chartered-enterprise types in the charter MVP (§7.2):** `AIRSHIP` first only, or `AIRSHIP` +
      `MINE`/`FARM` together? And do charters carry **tax obligations** from the start (Part 3) or ship
      rights-only first? *(Recommendation: rights-only framework in 2b; obligations land with Part 3.)*
- [ ] **Secession war decision — whose officer (§8.1)?** Maintainer-flagged: direct-member leave → root
      federation's top officer; carved-sub-region leave → immediate parent's top officer (leaning), escalate
      to root? Settle during the build.
- [ ] **What "hostile/war" does beyond wanted-propagation (§8.1):** just guards-aggro + lawful PvP
      (recommended MVP), or also sever trade / void embassies (§7.1) / other effects? And is a general
      `/realm war declare ↔ sue-for-peace` surface in 2b scope, or only the secession trigger for now?
- [ ] **Admin leave-lock default duration (§8.1)** and the list of *other* "few exceptions" to at-will, if
      any beyond the admin lock.

---

## 13. Cross-repo: `pcmc-territory` issues to file

> Ship-as-jurisdiction (§7.1) needs Part 1 (territory) resolver work. **This session's GitHub scope is
> limited to `theasshats/project-commonwealth` and cannot write to `theasshats/pcmc-territory`** (an
> `issue_write` there was denied; the repo-scoping tools to add it aren't available in this session). So
> these are captured here to file on that repo — by a session scoped to it, or by hand.

**Issue 1 — Resolve `aeroclaims` sub-level claims (claimed airships resolve as ungoverned).**
`resolve(level, chunkPos)` is chunk-based (OPAC chunk claims + MineColonies colonies); an `aeroclaims`
claim is bound to a **sub-level**, not a chunk, so a player on a claimed ship resolves to *ungoverned*.
Add a resolution source mapping *entity → sub-level it's on → owning party/claim → entity*. Spike whether
`aeroclaims` exposes the sub-level's owning party/owner via API. Not blocking Part 2 MVP (resolve by the
ground chunk until then), but the chunk-only assumption should be a documented limitation.

**Issue 2 — Precedence: a claimed sub-level is extraterritorial over the territory beneath it (embassy).**
When a claimed ship sits inside/over another faction's claim, occupants **aboard** are governed by the
**ship-realm**, not the host. Define the resolver precedence "on a claimed sub-level wins over the chunk
claim/colony beneath," and decide whether the same embassy rule extends to **chunk** enclaves (a realm's
OPAC claim inside another's colony), which would revisit Part 1's current colony-shielding. May fold into
Issue 1 as its precedence half.

---

_Refs: [`GOVERNANCE-MOD-SPEC.md`](GOVERNANCE-MOD-SPEC.md) (the whole-project 3-mod spec this refines —
§4 law engine, §3 tiers, §8 parts); [`GOVERNANCE.md`](GOVERNANCE.md) (scoping + path comparison);
[`CUSTOM-MODS.md`](CUSTOM-MODS.md) (own-repo + mod-mirror pattern); `docs/SYSTEMS.md` §3a
(event-driven/position-local perf doctrine); issue #260; `theasshats/pcmc-territory` PR #1 (the Part 1
source this contract is read from) and its issue #3 (resolver thread-safety); `aeroclaims` (Modrinth
`CwZ8q37q`), `opacpvp` (OPAC PvP Support), and MineColonies *War 'N Taxes* (addon prior art)._
