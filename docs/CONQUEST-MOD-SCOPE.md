# Colony Factions & Warfare — `pcmc-conquest` scope (working name)

> **Status: CONCEPT / early scope — nothing built, nothing accepted.** Captures the vision worked through
> in the session of **2026-06-13**: a MineColonies-based **faction + army** mod that **replaces Valarian
> Conquest** (`valarian-conquest`), reuses VC's *ideas* under a strict **clean-room** rule (no VC code —
> VC is **All Rights Reserved**), and **integrates with `pcmc-realms`** (the governance "laws/government"
> mod) at the **MineColonies guard system**. The working name `pcmc-conquest` is **provisional** — rename
> freely; it only nods to the VC lineage it replaces.
>
> This is a **sibling** to the governance trio (`pcmc-territory` → `pcmc-realms` → `pcmc-mint`), **not** a
> fourth part of it. Governance is *political* (tiers, laws, taxes, minting); this mod is *military +
> faction-identity* (joinable NPC powers, faction kit, armies). They meet at the MineColonies guard. Per
> [`CUSTOM-MODS.md`](CUSTOM-MODS.md) it would live in its **own repo** under `theasshats` and reach the
> pack via the mod-mirror packwiz pattern — it is **not** developed in this repo, and the web sandbox
> **cannot compile or run NeoForge**, so all work here is design-only (built/playtested on the box).
>
> **The architecture question — fork MineColonies vs. build an addon — is handed off in
> [`CONQUEST-FORK-VS-ADDON-HANDOFF.md`](CONQUEST-FORK-VS-ADDON-HANDOFF.md).** This doc is the feature
> scope that handoff reasons against.

---

## 1. The vision (maintainer, 2026-06-13)

> "Get rid of VC and do a MineColonies fork [or addon] and use VC systems for inspiration but none of
> their code. I'd like to make **NPC MineColonies factions that players could join**, **faction armor and
> equipment**, implement **more army mechanics** and build to the mod, and **integrate the upcoming
> laws/government mod into the guard system**."

The "upcoming laws/government mod" is **our own** `pcmc-realms` (the governance Part 2, already scoped in
[`GOVERNANCE-REALMS-SCOPE.md`](GOVERNANCE-REALMS-SCOPE.md)). So this mod does not invent a law engine — it
**consumes** `pcmc-realms` laws and gives them **military teeth** through the guard system.

## 2. Why a custom mod (the path here)

The session walked the options for tying Valarian Conquest into MineColonies and concluded:

- **Engine-level integration of VC and MineColonies is impossible** without Java: they are separate entity
  systems with separate AI; KubeJS/datapacks cannot bridge entity brains. A bridge needs a compiled mod.
- **VC is All Rights Reserved** — its internals can't be lifted or redistributed. A deep bridge often needs
  exactly that, so VC-as-base is off the table on licensing alone.
- **So: drop VC, and build the faction/army layer fresh on MineColonies**, reusing VC's *ideas* clean-room.
  MineColonies is the right base because it already ships the colony, citizen, and **guard** substrate this
  vision extends.

## 3. The license picture (the decisive finding)

Read from the jar metadata in the offline digest (`tools/mod-data/by-mod/`):

| Mod | License | Consequence |
| --- | --- | --- |
| **MineColonies / Structurize / minecolonies-compatibility** | **GPL-3.0** | **Fork *and* addon are both legal.** A fork/derivative must itself be GPL-3.0 with public source (copyleft) — fine for a `theasshats` mod (`pcmc-killfeed` is already GPL-3.0). |
| **Valarian Conquest** | **All Rights Reserved** | No code may be reused. Clean-room only — *ideas* are free, *code* is not. This mod must touch **zero** VC source. |

The GPL-3.0 finding is what makes this vision viable where the VC bridge was not: the base mod *permits*
derivatives. The remaining question is **engineering, not law** — fork vs addon (§5, and the handoff).

## 4. Feature scope (each with a feasibility lean)

These are the four pillars from the vision, mapped to where they likely land. **Leans are provisional —
the handoff's API spikes confirm them.**

| Feature | What it is | Likely shape | Notes / risk |
| --- | --- | --- | --- |
| **NPC factions players can join** | Pre-existing, NPC-run colony "powers" a player can pledge to / be recruited by — a persistent allegiance + reputation layer, VC's faction system reimagined on colonies. | **Addon if** MineColonies exposes programmatic colony creation + non-player ownership + an allegiance hook; **fork-or-cut if** NPC-owned colonies aren't reachable via API. **The hardest feature.** | Distinct from `pcmc-realms`, which is *player-founded political entities*. NPC-run joinable factions are new. |
| **Faction armor & equipment** | Faction-themed weapons/armor/shields/banners (VC had ~319 items), plus guards/citizens wearing their faction's kit by allegiance. | **Content is easy** (items/models/recipes — even partly data-driven). The **auto-equip-by-allegiance** hook on guards/citizens needs an API check. | Strong second-system weave candidate: gate faction kit through Create/MineColonies (mirrors VC's retired M-05 armorsmith→Create gate). |
| **More army mechanics** | The **unbounded** pillar — formations, ranks, offensive units, marching/attacking on command, sieges, inter-faction war. | **Mixed.** Defensive guard aggro is native (addon). **Offensive** control (march/attack/leave borders) is exactly the "reshape guard target-acquisition / pursue past border" that `GOVERNANCE-REALMS-SCOPE.md` §5 flagged as **fork territory**. | **Must be scoped to a concrete bounded list before estimating** — "army mechanics" balloons without limit. The MineColonies **raid system** (barbarian/pirate waves) is a possible reuse lever for army AI — investigate. |
| **Laws/government → guard system** | `pcmc-realms` laws enforced by guards: a soft-law violation makes the offender attackable by the colony's guards. | **Addon — already designed.** `pcmc-realms` §5 specifies the mechanism: flip the offender's MineColonies rank to a Fight-Guards rank via `IPermissions#setPlayerRank` for a "wanted" window; guards aggro natively. | **This pillar is the most solved.** This mod is a **consumer** of the `pcmc-realms` wanted-signal — reuse it, don't rebuild it. |

## 5. Architecture stance (pending the handoff)

**Lean: addon-first**, mirroring the conclusion `GOVERNANCE-REALMS-SCOPE.md` §5 already reached for the
guard-enforcement mechanism. An addon keeps the real MineColonies (its updates, and its load-bearing
"cheap-basics + locked-exclusives" production-route role in `SYSTEMS.md`), coexists, and is additive. A
**fork replaces MineColonies pack-wide** (blast radius = the whole pack) and must **continuously rebase**
on the MineColonies **snapshots** this pack runs — a standing tax. So a fork is justified **only** by a
concrete must-have that the public API provably can't express (per §5's bar) — most plausibly the
**offensive army mechanics**. The handoff makes the call **per feature**, not globally.

## 6. Relationship to `pcmc-realms` (the guard/law seam)

```
  pcmc-realms (governance)                 pcmc-conquest (factions + army)
  ─────────────────────────                ──────────────────────────────
  laws, tiers, federations, the            joinable NPC factions, faction
  "wanted" / soft-law signal (§5)   ──►     kit, army mechanics; CONSUMES
  rank-flip guard enforcement              the wanted-signal to drive guard
                                           aggression + faction warfare
```

`pcmc-realms` §5 already recommends building a **decoupled "wanted signal"** (event/tag/capability) as the
canonical soft-law output, with the guard rank-flip as one consumer. `pcmc-conquest` would be **another
consumer** of that same signal — faction hostility, bounty-hunting, and killfeed overlays all subscribe to
one wanted table. **Design these two together at the seam**; don't fork the wanted concept.

## 7. Decisions already made this session

- **Drop Valarian Conquest** (a manifest removal; its medieval-combat flavor is absorbed here). Removing VC
  also drops its worldgen structures and its retired Create weave — minor cleanup, track separately.
- **Clean-room VC** — ideas only, never code.
- **MineColonies is the base** (GPL-3.0 permits it).
- **Tier-2 custom-mod project** — own `theasshats` repo, mod-mirror packwiz delivery, built on the box.
- **Addon-first lean**, fork reserved for a proven API gap (handoff decides per feature).
- **`pcmc-realms` is the "laws/government" mod** — this mod consumes it, doesn't duplicate it.

## 8. Open questions (the fork-vs-addon determiners → the handoff)

These are what the handoff must resolve; listed here so the scope is honest about what's unsettled.

- [ ] Does MineColonies' public API allow **programmatic colony creation + NPC/non-player ownership** (for
      joinable NPC factions)? If not, is the feature fork-only, or recast (e.g. pre-placed structure colonies)?
- [ ] Can guards be given **offensive** behavior (march/attack on command, pursue past the colony border)
      via API/mixin, or only native defensive aggro? (§5: `AbstractEntityAIGuard` has **no public target-
      injection hook** — re-confirm on the pinned snapshot.)
- [ ] Can citizens/guards be made to **wear faction-specific equipment by allegiance** (slots + visuals)
      through the API?
- [ ] Can the MineColonies **raid system** be reused/extended as the army-AI substrate, rather than building
      unit AI from scratch?
- [ ] **Bounded list for "army mechanics"** — which specific mechanics are in scope? (Each maps to
      addon/mixin/fork differently; the pillar is meaningless to estimate until enumerated.)
- [ ] How do **NPC joinable factions** relate to `pcmc-realms`' player-founded entities — a separate entity
      kind, or pre-seeded realms the governance layer already understands?
- [ ] Milestone home: the colony/faction/army core leans **v0.11.0 (Magic & MineColonies, zagwar's
      pillar)**; the law-into-guard seam leans **v0.13.0 (Economy & governance, #260)**. Likely a
      `Backlog / living` project until split. Maintainer call.

---

_Refs: [`GOVERNANCE-REALMS-SCOPE.md`](GOVERNANCE-REALMS-SCOPE.md) (§5 — the guard mechanism + addon-vs-fork
decision aid this mod reuses); [`GOVERNANCE.md`](GOVERNANCE.md) (the 3-mod governance plan + path survey);
[`GOVERNANCE-MOD-SPEC.md`](GOVERNANCE-MOD-SPEC.md) (whole-project spec — confirms factions/army/equipment
are **not** in governance scope); [`CUSTOM-MODS.md`](CUSTOM-MODS.md) (own-repo + mod-mirror pattern);
[`CONQUEST-FORK-VS-ADDON-HANDOFF.md`](CONQUEST-FORK-VS-ADDON-HANDOFF.md) (the architecture handoff);
`tools/mod-data/by-mod/valarian_conquest-4.2.1.1-neoforge-1.21.1.txt` (the VC content inventory this mod
reimagines clean-room); `SYSTEMS.md` §3 (the production-route role a fork would inherit); issue #260
(governance)._
