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
- **Consented combat (duels / arenas).** Mutual opt-in (a `/duel` accept) or a jurisdiction-designated
  arena zone exempts both parties. A later nicety, not MVP.

**Generalize the check, not just the cases.** Model it as a pluggable **`Justification`** step the
SOCIAL evaluator runs for *any* law — so CONTRABAND can later exempt in-transit goods, TRESPASS a player
fleeing combat, etc., without re-plumbing each law. Self-defense is simply the first registered
justification. **Risk is low:** this is vanilla combat bookkeeping (a small per-player map + existing
revenge timestamps), not a new external dependency like the guard integration it gates.

### Mechanism (leading hypothesis — verify in the spike)

MineColonies guards already attack players natively via the colony **permission system**: each colony
ranks players, and a rank with the **"Fight Guards" permission** (the **Hostile** rank by default) is
attacked by that colony's guards on sight. (Confirmed from the MineColonies wiki/issue tracker:
"Fight Guards" causes guards to treat the player as an enemy; it is intended for hostile-ranked players.)

So Part 2's `SOCIAL` enforcement is: **on a violation in a colony-governed jurisdiction, set the
offender's MineColonies rank in that colony to a Fight-Guards rank for a time-boxed "wanted" window.**
The guards do the rest natively — no custom combat AI, and the "no guards nearby → get away with it"
behaviour falls out for free (a Hostile rank with no guards in range simply does nothing). When the
wanted timer expires, restore the prior rank.

### Spike (do this first, before any enforcement code)

1. **Is rank assignment API-reachable?** Verify `IColony#getPermissions()` exposes programmatic
   rank/player assignment (something like `setPlayerRank` / `addPlayer(UUID, Rank)` — names to confirm
   by decompiling the pinned `1.1.1327`, *do not guess*). Compile a real call in a dev env, exactly as
   Part 1 spiked its lookups.
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

### Edge cases to design for

- **OPAC-only jurisdictions have no guards.** A realm governing land via bound OPAC claims (no colony)
  has nothing to aggro — `SOCIAL` enforcement there is *purely* player-driven (wanted status + bounty,
  no NPC teeth). That's acceptable and on-theme; document it so it isn't read as a bug.
- **Aeronautics ships over a border** (per spec §9 / `aeronauticscompat`): a violation resolves by the
  block position under the actor; a ship crossing a no-PvP border mid-fight is governed by whatever chunk
  it's over at the hit. Flag for playtest whether that feels right.
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
   tangible* (land budget) and weaves governance into the claim economy. Sits in 2b/later.
3. **Hard-limiting claims (open).** Enforcing those allowances (preventing over-allowance claims)
   depends on OPAC exposing **either** a *settable per-player claim limit* **or** a *cancellable
   claim-creation event* — neither verified yet. This is a **shared spike with the territory mod**
   (its open question too). If neither exists, allowances degrade to *soft* (a `SOCIAL` "illegal
   over-claim" rather than a hard cap) — which fits the taxonomy fine.

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

---

## 9. Commands (extends Part 1's `/realm` tree)

Part 1 already ships `/realm found | info | whogoverns | debug bindclaim`. Part 2 adds (all
`Role`-gated via `atLeast`):

```
/realm promote                         (LEADER; attempt next municipality tier — §8)
/realm federate <name>                 (LEADER; compose a federation from consenting members)
/realm carve <name> <tier>             (LEADER; recursive sub-region, KINGDOM+)
/realm member <add|remove|role> ...    (LEADER; manage members/roles)
/realm law set <type> <value>          (LEADER; issue/replace a law)
/realm law unset <type>                (LEADER; clear a law → falls back to parent/vanilla)
/realm law list                        (CITIZEN; what's in force here, with the deciding tier)
/realm fine <player> <amount> [reason] (OFFICER; OFFICER-settled penalty — auto-withdraw via Part 3)
/realm wanted [player]                 (CITIZEN; view active wanted status in this jurisdiction)
/realm mint <amount>                   (LEADER; Part 3 — federation+, against reserves)
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

---

_Refs: [`GOVERNANCE-MOD-SPEC.md`](GOVERNANCE-MOD-SPEC.md) (the whole-project 3-mod spec this refines —
§4 law engine, §3 tiers, §8 parts); [`GOVERNANCE.md`](GOVERNANCE.md) (scoping + path comparison);
[`CUSTOM-MODS.md`](CUSTOM-MODS.md) (own-repo + mod-mirror pattern); `docs/SYSTEMS.md` §3a
(event-driven/position-local perf doctrine); issue #260; `theasshats/pcmc-territory` PR #1 (the Part 1
source this contract is read from) and its issue #3 (resolver thread-safety)._
