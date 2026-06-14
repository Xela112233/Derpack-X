# `pcmc-conquest` — fork vs. addon, decided per feature

> **Verdict: build it as a MineColonies *addon* (public API) + a small set of *targeted mixins*. No fork.**
> This is the recommendation the handoff [`CONQUEST-FORK-VS-ADDON-HANDOFF.md`](CONQUEST-FORK-VS-ADDON-HANDOFF.md)
> asked for. It extends the addon-vs-fork bar already set in
> [`GOVERNANCE-REALMS-SCOPE.md`](GOVERNANCE-REALMS-SCOPE.md) §5 to the three features that scope adds, and
> updates [`CONQUEST-MOD-SCOPE.md`](CONQUEST-MOD-SCOPE.md) §5 to match.
>
> **Status: design investigation, nothing built.** The mod lives in its own repo (per
> [`CUSTOM-MODS.md`](CUSTOM-MODS.md)); this doc is the architecture decision, not code.

---

## 0. The bar (adopted verbatim from realms-scope §5)

An **addon suffices** when MineColonies' **public API + targeted mixins** reach the behavior. A **fork is
justified only** for colony/guard behavior the public API provably can't express **and** that a targeted
mixin can't reach cleanly — and it must be budgeted as **replacing a core pack dependency**, not an
addon-sized task: a fork ships modified MineColonies to **every** player, forfeits the SYSTEMS.md §3
production-route role and upstream **snapshot** updates, can't coexist with unmodified colonies, and must be
**continuously rebased** on the snapshots this pack runs. A mixin, by contrast, is **snapshot-scoped and
reversible**. So the escalation ladder is **public API → access-widener → targeted mixin → (last resort)
fork** — you stop at the first rung that works.

The maintainer's steer (2026-06-14): **"tie into MineColonies as much as possible rather than reinventing
the wheel."** That points the same direction — reuse colonies, citizens, guards, the raid system, and the
rally mechanism rather than building parallel entities (which is what Valarian Conquest did). It also argues
*against* a fork independently of capability: a fork still needs all those MineColonies systems (it modifies
them, it doesn't remove them), so it buys no extra reuse — only extra maintenance.

## 1. How this was verified (and what still needs the box)

The decisive half — *"does this public API exist?"* — is answerable statically from MineColonies' GPL-3.0
source, and it was. **The jar itself could not be pulled in the web sandbox** (`forgecdn`,
`ldtteam.jfrog.io`, and Modrinth are all outside the egress allowlist — only `raw.githubusercontent.com` is
reachable), so the findings are read from the **GitHub `version/1.21` source at commit `aeb1ad8`** — the
branch the pinned `1.1.1327-1.21.1-snapshot` jar is built from — not from a decompile of the exact jar.
These guard/colony/raid API surfaces have been stable across many MineColonies versions, so the shape is
very likely byte-identical; **byte-exact confirmation against `1.1.1327` is a [needs box] item** (decompile
the jar, or add the CDN hosts to egress and re-pull). Every *runtime* claim ("does it actually behave this
way in-game") is flagged **[needs box]** and collected in the spike checklist (§7).

## 2. Per-feature verdict

| Feature | Verdict | Why |
| --- | --- | --- |
| **NPC joinable factions** | **addon** | Colony creation, NPC ownership, and the join/allegiance hook are all public API. |
| **Faction equipment by allegiance** | **addon** | Items are content; livery + worn-gear assignment are reachable via public API + the native gear pipeline. |
| **Army mechanics** (ranks, offensive march, formations, sieges/war) | **addon + targeted mixin** | The drive primitives (rally/patrol/follow, raid trigger) are public; only the arbitrary-range leash and the guard-mode setter need a small mixin into one or two core methods. |
| **Laws → guards** | **addon** (already solved) | Consumes the `pcmc-realms` wanted-signal via `setPlayerRank`; no new fork question. |

**Global call: addon core + a handful of targeted mixins. No fork.** A fork is reconsidered only under a
narrow, currently-unmet condition (§6).

---

## 3. NPC joinable factions — **addon**

The handoff called this "likely the hardest." It is the most API-dependent, but the public surface covers
it. The two sub-questions — *programmatic colony creation* and *non-player ownership* — both resolve yes,
and the *join/allegiance* hook falls out of the permission system.

**Colony creation (public).**
`com.minecolonies.api.colony.IColonyManager`:
```java
IColony createColony(@NotNull Level w, BlockPos pos, @NotNull Player player,
                     @NotNull String colonyName, @NotNull String pack);
```
Server-side colony creation is a public call. The catch is the `@NotNull Player` (documented "the player
that creates the colony - owner"): there is no `createColony(...UUID...)` overload, so an **NPC owner needs a
NeoForge `FakePlayer`** (a synthetic-UUID server player — the standard pattern automation mods use) passed
here, or owner reassignment immediately after (below).

**Non-player ownership (public).**
`com.minecolonies.api.colony.permissions.IPermissions` exposes the full owner surface:
- `boolean setOwner(Player player)` — set/replace owner (FakePlayer works).
- `void setOwnerAbandoned()` — mark the colony **ownerless** (an NPC-run colony with no human owner — the
  cleanest expression of "an NPC power").
- `UUID getOwner()`, `Map.Entry<UUID, ColonyPlayer> getOwnerEntry()`, `Rank getRankOwner()`.

**Join / allegiance hook (public).** Pledging a *player* to a faction is just adding their UUID to that
colony's permission roster at a chosen rank:
- `boolean addPlayer(@NotNull UUID id, String name, Rank rank)` / `addPlayer(@NotNull GameProfile, Rank)`.
- `boolean setPlayerRank(UUID id, Rank rank, Level world)` (also the §5 wanted-flip).
- Rank accessors: `getRankFriend()`, `getRankOfficer()`, `getRankHostile()`, `getRankNeutral()`.

So "join the Barathian faction" = `addPlayer(playerUUID, name, getRankFriend())` in the Barathian NPC
colony; promotion/desertion = `setPlayerRank(...)`; declaring you an enemy = flip to `getRankHostile()`
(which doubles as the guard-aggro lever). **No fork — this is the public permission API.**

**Faction identity / reputation / join-state** itself is **your own `SavedData`**, hung off the colony id
(exactly as `pcmc-realms` hangs `RealmGov` off the territory entity, realms-scope §2). MineColonies stores
*who is in the colony at what rank*; the *faction flavor, standing, and reputation curve* are yours.

**[needs box] confirmations:**
- A `FakePlayer`-owned or `setOwnerAbandoned()` colony **persists across save/load and ticks** (citizens
  spawn, guards work, no owner-null crash).
- Whether `createColony` requires a **physical Town Hall block** already placed at `pos` (it likely expects
  the hut to exist) — i.e. faction colonies probably need a **pre-placed Town Hall** (structure/worldgen or
  an admin/datapack placement), then `createColony` over it. This is a *recast detail, not a fork trigger*:
  pre-placing the hut is worldgen/structure work, all addon-side.
- Whether two abandoned/NPC colonies can mutually set each other `Hostile` for faction-vs-faction state.

---

## 4. Faction equipment by allegiance — **addon**

**Content is trivial and pure addon:** faction armor sets, weapons, shields, banners (Valarian Conquest
shipped ~319 such items as a *feature checklist* — reimagined clean-room, never its code). Items, models,
and recipes are datapack/registration; this is the easy half and a strong **second-system weave**
candidate — gate the faction kit through Create/MineColonies production (the retired VC armorsmith→Create
gate, reimagined), so it earns its place per the systems loop.

**Auto-equip + livery by allegiance (public enough):**
- **Appearance:** `ICitizenData.setCustomTexture(UUID texture)` is public — per-citizen faction livery.
- **Worn gear:** `ICitizenData.getEntity()` returns `Optional<AbstractEntityCitizen>` (a `LivingEntity`),
  so vanilla `setItemSlot(EquipmentSlot, ItemStack)` reaches the worn slots; and guards natively equip from
  the building's **gear/request pipeline**, so the clean route is to feed each faction's building its kit
  and let citizens equip it the normal way (the maintainer's "reuse, don't reinvent" path).

**[needs box]:** does the citizen AI / request system **overwrite** forced gear on its next tick (if so,
prefer supplying the kit through the request pipeline over force-setting slots, or a small mixin to pin the
slot); does a custom citizen texture register and render cleanly. Neither is a fork trigger.

---

## 5. Army mechanics — **addon + targeted mixin** (all four sub-mechanics confirmed in scope)

The maintainer confirmed the bounded list is **all of**: unit ranks & promotion, offensive march/attack-move,
formations / grouped movement, and sieges & inter-faction war. Each maps below. The headline: the **drive
primitives are public**; only **two** behaviors need a targeted mixin, and neither approaches a fork.

**5a. Unit ranks & promotion — addon.** Maps onto MineColonies guard **ranks** (`IPermissions` rank API)
and **guard-building levels** (the combat academy / barracks / archery / guard-tower hut ladder). Your
unit-tier and promotion *state* is your own; the in-world effect (a higher-rank guard) is set through the
colony API. No core changes.

**5b. Offensive march / attack-move — addon within range, addon+mixin for free range.**
The public guard-control interface `com.minecolonies.api.colony.buildings.IGuardBuilding` is richer than
§5 assumed. It exposes the load-bearing primitives directly:
- `void setRallyLocation(ILocation)` / `ILocation getRallyLocation()` — move guards to a rally point and
  fight there (the guard AI reads `getRallyLocation()` every tick).
- `void setGuardPos(BlockPos)` — set the guard stance position.
- `void addPatrolTarget(BlockPos)` / `resetPatrolTargets()` / `getNextPatrolTarget()` — waypoint patrols.
- `void setPlayerToFollow(Player)` — make a building's guards follow a (fake)player anywhere in-world (the
  guard-help-scroll behavior, but as API).

So **"march my faction's guards to a point and fight there" is achievable through the public API, no fork** —
this overturns §5's worst-case read ("no public way to command guards offensively"). The §5 caution was
about *target acquisition* (`AbstractEntityAIGuard.checkForTarget`, core, no inject hook), but the **rally /
follow / patrol** path reaches the same outcome (guards go where you point and engage what's there) without
touching target selection.

**Two behaviors still need a targeted mixin** (verified in core, not exposed on the api interface):
1. **Arbitrary-range march.** Core `AbstractBuildingGuards.getRallyLocation()` re-applies a **leash every
   tick**: a rally point outside the owning colony is allowed **only** with the `TELESCOPE` research **and**
   within **500 blocks** of colony center; otherwise it nulls the rally. `setRallyLocation()` past that is
   silently dropped. Marching an army across the map = a small mixin/AT that neuters that one check.
2. **Forcing guard mode** (GUARD / PATROL / FOLLOW). The mode is a settings value keyed by core
   `AbstractBuildingGuards.GUARD_TASK`; the `api.ISetting` interface is **read-only** (`getValue()`, no
   public setter), so flipping mode programmatically needs the core key + a cast — i.e. an AT/mixin.

Both are **one or two methods in one core class** — the textbook targeted-mixin case the handoff §4 says to
prefer over a fork (reversible, snapshot-scoped). `scepterguard` and `scroll_guard_help` are
**player-item-only** wrappers over these same `IGuardBuilding` primitives; useful as a behavior model, not
as an API.

**5c. Formations / grouped movement — addon.** A controller issuing per-unit rally/patrol/guard-pos targets
(5b primitives) over a squad, or over your own/raider entities, is **your** coordination code — it doesn't
modify MineColonies internals. Tighter formation-keeping (spacing, facing) is custom movement logic on top,
still addon. No fork.

**5d. Sieges & inter-faction war — addon, via the raid substrate.** This is the heaviest sub-mechanic and
the one most tempted toward a fork — but MineColonies' **own raid system is drivable through the public
API**, which is exactly the "reuse, don't reinvent" answer:
`com.minecolonies.api.colony.managers.interfaces.IRaiderManager`:
```java
RaidSpawnResult raiderEvent(@NotNull RaidSettings raidSettings); // "Trigger a specific type of raid on a colony."
void setRaidNextNight(RaidSettings raidSettings);                // schedule the next raid
int  calculateRaiderAmount(int raidLevel);  boolean isRaided();  // sizing + state
```
So **a faction's offensive force = a `raiderEvent(RaidSettings)` aimed at the target colony**, using
MineColonies' own mobile attacking raiders (`AbstractEntityMinecoloniesRaider`) that path to and assault the
colony — **no custom unit AI, no guard-AI fork.** Faction-themed attackers = register **your own**
`AbstractEntityMinecoloniesRaider` subtypes (addon content / livery), if `RaidSettings` admits a custom
horde type.

**[needs box] for the army pillar:**
- Confirm the **500-block leash actually nulls** an out-of-range `setRallyLocation` at runtime (it's the one
  thing forcing mixin #1).
- Confirm `raiderEvent(RaidSettings)` **spawns and paths** a raid at a chosen colony on command, and inspect
  `RaidSettings` for whether a **custom/registered horde type** (faction livery) is accepted, or only the
  built-in barbarian/pirate/norsemen/mummy/drowned hordes.
- Confirm the rally-leash + guard-mode mixins **load and hold** against the `1.1.1327` jar (and gauge how
  often they'd re-break across the snapshot cadence — that's the fork-reconsideration trigger, §6).
- Confirm faction-vs-faction `Hostile`-rank cross-aggro between two NPC colonies' guards.

---

## 6. When a fork *would* be reconsidered (and why it isn't now)

Per §5's bar, a fork must clear **"replacing a core pack dependency,"** not "adding a feature." Nothing in
the four pillars clears it: every capability is reached by public API or a **targeted, reversible mixin into
one or two core guard methods**. A fork would still depend on all the same MineColonies systems (it modifies
them, doesn't remove them) while **forfeiting** upstream snapshot updates, the SYSTEMS.md §3 production-route
role, and coexistence with unmodified colonies — and a fork still creates **no guards where there are none**
(OPAC-only land, ship-realms), so it doesn't even solve §5's guard-gap edge cases.

The **only** scenario that reopens the fork question: the **rally-leash / guard-mode mixins prove unstable
across the MineColonies snapshot cadence** (they re-break on most bumps and the rebase cost rivals a fork's)
**and** arbitrary-range offensive march is a hard must-have **and** no access-widener / narrow coremod holds
the seam. Even then, escalate **mixin → access-widener → narrow coremod** before a full fork. Treat that as a
*later, separate* decision triggered by observed mixin instability, not by the feature list.

## 7. Spike checklist (run on the box — statics pre-filled)

Mirrors the Part 1/Part 2 spike pattern: *API exists? → compiles? → works at runtime? → addon or fork?* The
"API exists / compiles" column is filled from the source read above; the runtime column is the box's job.

| # | Spike | API exists (static) | Compiles | Runtime [needs box] |
| --- | --- | --- | --- | --- |
| 1 | `IColonyManager.createColony(Level,BlockPos,Player,String,String)` with a **FakePlayer** owner | ✅ public; signature confirmed | ☐ | ☐ FakePlayer/abandoned-owner colony persists across save/load, ticks, no owner-null crash |
| 2 | Does `createColony` need a **pre-placed Town Hall** at `pos`? | ❓ likely yes (recast, not fork) | n/a | ☐ create over a placed hut vs. bare pos |
| 3 | Join hook: `IPermissions.addPlayer(UUID,name,getRankFriend())` then `setPlayerRank(...)` | ✅ public | ☐ | ☐ player gains/loses colony membership + permissions live |
| 4 | Two NPC colonies set each other `getRankHostile()` → mutual guard aggro | ✅ public rank API | ☐ | ☐ cross-colony guard aggro actually fires |
| 5 | Citizen livery `ICitizenData.setCustomTexture(UUID)` + gear via `getEntity()`/request pipeline | ✅ public | ☐ | ☐ texture renders; forced gear isn't overwritten by AI/requests |
| 6 | Offensive march: `IGuardBuilding.setRallyLocation/setGuardPos/addPatrolTarget/setPlayerToFollow` | ✅ public | ☐ | ☐ guards move to rally/patrol and fight there |
| 7 | Arbitrary-range march: **mixin/AT** neutering core `AbstractBuildingGuards.getRallyLocation()` 500-block+TELESCOPE leash | ✅ leash confirmed in core | ☐ | ☐ mixin lets rally exceed 500 blocks; ☐ confirm leash nulls it *without* the mixin |
| 8 | Force guard mode: **mixin/AT** on core `GUARD_TASK` (read-only `api.ISetting`) | ✅ confirmed core-only | ☐ | ☐ mode flips GUARD/PATROL/FOLLOW programmatically |
| 9 | Sieges/war: `IRaiderManager.raiderEvent(RaidSettings)` / `setRaidNextNight` aimed at a colony | ✅ public; "trigger a raid on a colony" | ☐ | ☐ raid spawns+paths at target; ☐ does `RaidSettings` accept a **custom horde type** for faction livery |
| 10 | Mixin stability: do #7/#8 load against the **exact `1.1.1327` jar** and survive a snapshot bump | ❓ source matches branch; jar byte-exactness unverified | ☐ | ☐ load against jar; gauge rebreak frequency (fork-reconsider trigger) |
| 11 | Byte-exact API check vs. the **pinned jar** (decompile `1.1.1327`, or add `forgecdn`/`ldtteam`/Modrinth to egress and re-pull) | ❓ verified vs. `version/1.21`@`aeb1ad8`, not the jar | n/a | ☐ |

## 8. Prior art (read, don't reinvent)

- **MineColonies: War 'N Taxes** — an *addon* adding taxes + war to colonies: standing evidence the public
  surface carries war-class features without a fork. (Exact API surface it uses is a worthwhile confirm; the
  agent pass couldn't fully enumerate it — a [needs box]/follow-up read.)
- **minecolonies-compatibility** (already in the pack, GPL-3.0) — adds jobs/research incl. a TACZ "Gunner"
  job: prior art that the colony↔combat-mod seam is addon-reachable (jobs + research are public extension
  points).
- **Native offensive surfaces in MineColonies itself** — the rally-guards banner, guard scepter, and
  guard-help scroll already drive guards offensively via `IGuardBuilding`; the mod's own items are the
  behavior model to copy, not fork.

## 9. Notes for a correctly-scoped session (cross-repo — cannot act from here)

This session is scoped to `theasshats/project-commonwealth`; it **cannot** write to the `pcmc-conquest` mod
repo or to `ldtteam/minecolonies`. Capture, to be filed by hand or a correctly-scoped session:
- **Add CDN egress** (`*.forgecdn.net`, `ldtteam.jfrog.io`, `api.modrinth.com`/`cdn.modrinth.com`) for any
  future session that must decompile the **pinned jar** for byte-exact API checks — this session could only
  reach `raw.githubusercontent.com`.
- **Mod-repo issues** to open when `pcmc-conquest` exists: the §7 spike checklist as the first milestone;
  the two targeted mixins (rally-leash, guard-mode) as their own tracked, snapshot-scoped tasks.
- **Possible upstream ask** (MineColonies, GPL-3.0): a public guard-mode setter and an optional rally-range
  config would remove mixins #7/#8 entirely — worth a feature request before writing the mixins.

---

_Refs: [`CONQUEST-FORK-VS-ADDON-HANDOFF.md`](CONQUEST-FORK-VS-ADDON-HANDOFF.md) (the task brief);
[`CONQUEST-MOD-SCOPE.md`](CONQUEST-MOD-SCOPE.md) §5 (architecture stance, updated to match this verdict);
[`GOVERNANCE-REALMS-SCOPE.md`](GOVERNANCE-REALMS-SCOPE.md) §5 (the guard mechanism + the addon-vs-fork bar
this extends); [`CUSTOM-MODS.md`](CUSTOM-MODS.md) (own-repo + mod-mirror delivery); `CLAUDE.md` ("where
you're running" — sandbox/egress limits). MineColonies source verified at `ldtteam/minecolonies`
`version/1.21`@`aeb1ad8` (the build branch for the pinned `1.1.1327-1.21.1-snapshot`); byte-exact-vs-jar is
[needs box]. VC content inventory (`tools/mod-data/by-mod/valarian_conquest-4.2.1.1-neoforge-1.21.1.txt`) used
as a clean-room feature checklist only — never its code (VC is All Rights Reserved)._
</content>
</invoke>
