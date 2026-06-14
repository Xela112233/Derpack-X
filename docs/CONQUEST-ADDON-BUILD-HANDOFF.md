# Handoff — building `pcmc-conquest` as a MineColonies add-on

> **For:** the session (or maintainer) that **builds** the colony-factions-and-warfare mod, on a machine
> with a **NeoForge dev/build env** — i.e. the box or the mod's own repo, **not** the web sandbox (which
> can't compile or run NeoForge, and can't even pull the MineColonies jar — see §6).
>
> **Architecture is already decided:** a **MineColonies add-on** on the public API plus **two targeted
> mixins** — **not a fork**. That decision and its evidence are
> [`CONQUEST-FORK-VS-ADDON.md`](CONQUEST-FORK-VS-ADDON.md); the feature scope and the public-API reuse
> surface are [`CONQUEST-MOD-SCOPE.md`](CONQUEST-MOD-SCOPE.md) (§5a is the API building-block table). This
> doc is the **build plan** — read those two first, then work this.

---

## 0. Read first (in order)

1. [`CONQUEST-FORK-VS-ADDON.md`](CONQUEST-FORK-VS-ADDON.md) — the per-feature verdict, the source evidence,
   and the **§7 box spike checklist** (this handoff's first milestone is running that checklist).
2. [`CONQUEST-MOD-SCOPE.md`](CONQUEST-MOD-SCOPE.md) — the four features, §5a the API reuse surface, §6 the
   `pcmc-realms` seam, and the clean-room VC rule.
3. [`CUSTOM-MODS.md`](CUSTOM-MODS.md) — the own-repo + mod-mirror packwiz delivery pattern this mod ships
   through (the `pcmc-killfeed`/`pcmc-arcana` model).
4. [`GOVERNANCE-REALMS-SCOPE.md`](GOVERNANCE-REALMS-SCOPE.md) §5 — the wanted-signal + guard rank-flip this
   mod **consumes** (it does not rebuild the law engine).

## 1. What's settled (do not re-litigate)

- **Add-on, not fork.** Public `com.minecolonies.api.*` + two targeted mixins into core
  `AbstractBuildingGuards`. The fork question is closed unless the §3 mixins prove snapshot-unstable
  (`CONQUEST-FORK-VS-ADDON.md` §6 is the only reopen condition).
- **GPL-3.0.** MineColonies is GPL-3.0, so the add-on is GPL-3.0 with public source (like `pcmc-killfeed`).
- **Clean-room VC.** VC is All Rights Reserved. The inventory
  (`tools/mod-data/by-mod/valarian_conquest-4.2.1.1-neoforge-1.21.1.txt`) is a **feature checklist**, never
  a code source. VC built parallel summon-entity soldiers; we deliberately **don't** — we drive
  MineColonies' own colonies/citizens/guards/raiders.
- **Own repo + mod-mirror.** New `theasshats` repo, `build.yml` attaches the jar to a `v*` release, the pack
  references it by a packwiz manifest with the release sha1 (`CUSTOM-MODS.md`). Not built in the pack repo.
- **Pinned target:** MineColonies `1.1.1327-1.21.1-snapshot`, NeoForge `21.1.x` per the pack's `pack.toml`.

## 2. Build order (milestones)

Each milestone ends by ticking its part of the `CONQUEST-FORK-VS-ADDON.md` §7 spike table — *API compiles →
works at runtime* — before the next starts. Sequenced so the riskiest/most-foundational lands first.

**M0 — Repo + dev env.** Create the `theasshats` mod repo (clone the `pcmc-killfeed` GPL-3.0 scaffolding:
license, `build.yml`, mod metadata). Add MineColonies + Structurize as `compileOnly`/dev dependencies from
the ldtteam maven (`https://ldtteam.jfrog.io/...`). Confirm a no-op mod boots a dev client **with
MineColonies present**.

**M1 — Spike pass (de-risk before features).** Run the §7 checklist as throwaway test code: FakePlayer
`createColony`, `addPlayer`/`setPlayerRank` join, `setCustomTexture`, `IGuardBuilding` rally/follow, and
`IRaiderManager.raiderEvent`. Land the two mixins (§3) and confirm they load against the **pinned jar**.
**This milestone answers every [needs box] item** and is the real go/no-go for the design.

**M2 — Faction state + NPC colonies (Pillar 1).** Own `SavedData` for faction identity/standing/membership,
keyed by colony id. Found NPC faction colonies (FakePlayer/`setOwnerAbandoned`, over a pre-placed Town Hall
if M1 shows it's required). Join/leave/promote via the permission API. This is the spine the rest hangs on.

**M3 — Faction equipment (Pillar 2).** Content first (clean-room armor/weapons/shields/banners — a strong
**second-system weave**: gate the kit through Create/MineColonies production per the systems loop). Then
auto-equip by allegiance: livery via `setCustomTexture`, kit via the gear/request pipeline.

**M4 — Army mechanics (Pillar 3).** Ranks/promotion (guard ranks + hut levels); offensive march via the
rally/follow/patrol primitives (+ mixin #1 for arbitrary range); formations as a squad controller over those
primitives; sieges/inter-faction war via `IRaiderManager.raiderEvent` aimed at the enemy colony (register
faction-themed raider subtypes). Plus mixin #2 for forced guard mode.

**M5 — Laws → guards seam (Pillar 4).** Subscribe to the `pcmc-realms` **wanted-signal** (event/tag/
capability — `GOVERNANCE-REALMS-SCOPE.md` §5) and flip the offender's rank to `getRankHostile()` for the
wanted window. Build this **with** `pcmc-realms` at the seam; don't duplicate the wanted table.

**M6 — Pack integration + playtest.** Add the mod-mirror manifest to the pack on a version branch, drop
Valarian Conquest in the same branch (manifest removal + its worldgen/Create-weave cleanup, tracked
separately per scope §7), and carry the per-feature `## Playtest` checklist (CLAUDE.md's playtest rule —
green CI is not a load check).

## 3. The two mixins (the only non-API pieces — keep them surgical)

Both target core `com.minecolonies.core.colony.buildings.AbstractBuildingGuards`. Prefer an **access
widener/transformer** if it suffices before a mixin; keep each to the single method it needs; treat them as
**snapshot-scoped** (re-verify on every MineColonies bump — their stability is the fork-reopen tripwire).

1. **Rally-range leash.** Core `getRallyLocation()` re-applies a 500-block + `TELESCOPE`-research cap every
   tick and nulls an out-of-range rally. Mixin to lift/relax the cap for faction armies marching beyond a
   colony.
2. **Guard-mode setter.** The GUARD/PATROL/FOLLOW mode is keyed by core `GUARD_TASK` behind a read-only
   `api.ISetting` (no public setter). Mixin/AT to flip mode programmatically.

If MineColonies upstream later adds a public guard-mode setter and a rally-range config, **both mixins
disappear** — file that as an upstream feature request (§5) before writing them.

## 4. Open design questions to resolve while building

- **NPC factions vs. `pcmc-realms` entities** (scope §8): a separate entity kind (own `SavedData` by colony
  id — the leaning answer) or pre-seeded realms the governance layer already understands? Decide at the seam.
- **Pre-placed Town Hall** for `createColony` (spike #2): if required, faction colonies need a worldgen/
  structure or admin placement step — an addon-side recast, design it into M2.
- **Custom raider horde types** (spike #9): does `RaidSettings` admit a registered faction horde, or only the
  built-in barbarian/pirate/etc.? Gates faction livery on attackers in M4.
- **Milestone home in the pack roadmap:** colony/faction/army core leans **v0.11.0 (Magic & MineColonies)**;
  the law→guard seam leans **v0.13.0 (Economy & governance, #260)**. Likely a `Backlog / living` project
  until split — maintainer call.

## 5. Cross-repo actions (cannot be done from a `project-commonwealth`-scoped session)

Capture-only here; file by hand or from a correctly-scoped session:

- **Create the `theasshats/pcmc-conquest` repo** and its first issues: the §7 spike checklist (M1), the two
  mixins as tracked snapshot-scoped tasks, the per-pillar milestones above.
- **Add CDN egress** (`*.forgecdn.net`, `ldtteam.jfrog.io`, `api.modrinth.com`/`cdn.modrinth.com`) to any
  web session expected to pull/decompile the MineColonies jar for byte-exact API checks — the
  investigation session could reach **only** `raw.githubusercontent.com`.
- **Upstream feature request to MineColonies** (GPL-3.0): a public guard-mode setter + optional rally-range
  config — would delete both §3 mixins. Worth asking before writing them.
- **Pack-side, when the jar exists:** add the `mods/pcmc-conquest.pw.toml` mod-mirror manifest, remove
  `valarian-conquest`, `./tools/packwiz refresh`, and add the playtest checklist to the version PR.

## 6. Why this can't be finished in the web sandbox

The web sandbox **can't compile or run NeoForge**, and its egress allowlist excludes `forgecdn`,
`ldtteam.jfrog.io`, and Modrinth, so it **can't even pull the MineColonies jar** — only
`raw.githubusercontent.com` is reachable. That's why the architecture investigation verified against the
GitHub **source** (`version/1.21`@`aeb1ad8`, the build branch for `1.1.1327`) and flagged every runtime
claim and the byte-exact-vs-jar check as **[needs box]**. The build itself (M0–M6) is a box/mod-repo job.

---

_Refs: [`CONQUEST-FORK-VS-ADDON.md`](CONQUEST-FORK-VS-ADDON.md) (the decided architecture + §7 spike list);
[`CONQUEST-MOD-SCOPE.md`](CONQUEST-MOD-SCOPE.md) (feature scope + §5a API reuse surface);
[`CONQUEST-FORK-VS-ADDON-HANDOFF.md`](CONQUEST-FORK-VS-ADDON-HANDOFF.md) (the prior architecture-decision
handoff — now answered); [`CUSTOM-MODS.md`](CUSTOM-MODS.md) (own-repo + mod-mirror);
[`GOVERNANCE-REALMS-SCOPE.md`](GOVERNANCE-REALMS-SCOPE.md) §5 (the wanted-signal this mod consumes); issue
#260 (governance). MineColonies pinned `1.1.1327-1.21.1-snapshot`._
</content>
</invoke>
