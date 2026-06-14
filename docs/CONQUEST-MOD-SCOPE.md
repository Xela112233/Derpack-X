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

## 5. Architecture stance — **decided: addon + targeted mixins, no fork**

The fork-vs-addon investigation is **done** — see [`CONQUEST-FORK-VS-ADDON.md`](CONQUEST-FORK-VS-ADDON.md)
for the per-feature verdict, the source evidence, and the box spike checklist. Outcome:

| Feature | Verdict |
| --- | --- |
| NPC joinable factions | **addon** — `IColonyManager.createColony` + a FakePlayer/`setOwnerAbandoned` owner; the join hook is `IPermissions.addPlayer(UUID, rank)` / `setPlayerRank` (all public API). |
| Faction equipment by allegiance | **addon** — items are content; livery via `ICitizenData.setCustomTexture`, worn gear via `getEntity()` + the native gear/request pipeline. |
| Army mechanics (ranks, offensive march, formations, sieges/war) | **addon + targeted mixin** — guard rally/patrol/follow (`IGuardBuilding`) and the raid trigger (`IRaiderManager.raiderEvent`) are public; only the 500-block rally **leash** and the guard-**mode** setter need a small mixin into one core class. Sieges/war ride the **raid system**, not custom unit AI. |
| Laws → guards | **addon** — already solved; consumes the `pcmc-realms` wanted-signal via `setPlayerRank`. |

**Global call: addon-first, confirmed.** Every capability is reachable through MineColonies' public API or a
**targeted, reversible mixin** into one or two core guard methods — none clears §5's fork bar ("replacing a
core pack dependency"). An addon keeps the real MineColonies (its updates and its `SYSTEMS.md` §3
production-route role), coexists, and matches the maintainer's "reuse, don't reinvent" steer (reuse
colonies/citizens/guards/raids/rally rather than VC-style parallel entities). A **fork** is reconsidered
**only** if the two guard mixins prove unstable across the snapshot cadence *and* arbitrary-range offensive
march is a hard must-have *and* no access-widener/coremod holds — a later, separate, evidence-triggered
call, not part of this scope. Findings are read from MineColonies' GPL-3.0 source (`version/1.21`@`aeb1ad8`);
**byte-exact confirmation against the pinned `1.1.1327` jar and every runtime behavior are [needs box]** (the
web sandbox can't pull the jar — `forgecdn`/`ldtteam`/Modrinth are outside its egress allowlist).

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

The fork/addon determiners are now **resolved** by [`CONQUEST-FORK-VS-ADDON.md`](CONQUEST-FORK-VS-ADDON.md)
(statics confirmed against the 1.21 source; runtime checks are its §7 box spike list). Remaining open items
are recast details and scheduling, not architecture.

- [x] **Programmatic colony creation + NPC/non-player ownership** → **yes, public API.**
      `IColonyManager.createColony(...)` + a FakePlayer owner or `IPermissions.setOwnerAbandoned()`. Likely
      needs a **pre-placed Town Hall** at the position (a worldgen/structure recast, not a fork). Addon.
- [x] **Offensive guard behavior** → **yes.** Rally/patrol/follow primitives are public on `IGuardBuilding`;
      only **arbitrary-range march** (core 500-block leash) and **forcing guard mode** need a targeted mixin.
      §5's "no public target-injection hook" stands for `AbstractEntityAIGuard`, but the rally path reaches
      the same outcome without it.
- [x] **Faction equipment by allegiance** → **yes.** `ICitizenData.setCustomTexture(UUID)` for livery;
      `getEntity()` + the native gear/request pipeline for worn kit. Addon.
- [x] **Raid system as army substrate** → **yes.** `IRaiderManager.raiderEvent(RaidSettings)` triggers a raid
      at a colony — the siege/inter-faction-war substrate, reusing MineColonies' own raiders.
- [x] **Bounded "army mechanics" list** → **confirmed (maintainer, 2026-06-14): all four** — unit ranks &
      promotion, offensive march/attack-move, formations/grouped movement, sieges & inter-faction war. All
      land addon or addon+mixin (see the army section of the verdict doc).
- [ ] How do **NPC joinable factions** relate to `pcmc-realms`' player-founded entities — a separate entity
      kind (own `SavedData` keyed by colony id, the leaning answer), or pre-seeded realms the governance
      layer already understands? Design at the seam with `pcmc-realms`.
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
