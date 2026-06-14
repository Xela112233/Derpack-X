# Handoff — fork vs. addon for `pcmc-conquest` (the colony-warfare mod)

> **For:** a future instance investigating the architecture of the colony-factions-and-warfare mod
> (working name `pcmc-conquest`), ideally one that can reach a **NeoForge dev/build env** or will write
> the **spikes for the box** (the web sandbox cannot compile or run NeoForge).
>
> **Task:** decide — and justify, **per feature** — whether `pcmc-conquest` should be a **MineColonies
> addon** (a separate GPL mod on MineColonies' public API + targeted mixins) or a **fork** of MineColonies
> (a modified MineColonies shipped in place of upstream). Produce a recommendation against the bar already
> established in [`GOVERNANCE-REALMS-SCOPE.md`](GOVERNANCE-REALMS-SCOPE.md) §5.
>
> **Status of the surrounding work:** the mod's feature scope is
> [`CONQUEST-MOD-SCOPE.md`](CONQUEST-MOD-SCOPE.md) (CONCEPT, nothing built). The guard-enforcement
> mechanism and a first addon-vs-fork analysis already exist in the governance work — **do not re-derive
> them; build on them.**

---

## 0. Read first (in order)

1. **[`CONQUEST-MOD-SCOPE.md`](CONQUEST-MOD-SCOPE.md)** — the four features you're deciding architecture
   for, the GPL-3.0 license finding, and the open questions.
2. **[`GOVERNANCE-REALMS-SCOPE.md`](GOVERNANCE-REALMS-SCOPE.md) §5** — *the prior art.* It already (a)
   specifies the MineColonies **guard soft-enforcement mechanism** (rank-flip via
   `IPermissions#setPlayerRank` → native guard aggro), and (b) contains a full **"addon vs. fork (decision
   aid)"** that concluded **addon-first** for guard enforcement, with an explicit bar for when a fork is
   justified. **Your job extends that analysis to the three *new* features this mod adds.**
3. **[`GOVERNANCE.md`](GOVERNANCE.md)** — the 3-mod governance plan and how `pcmc-realms` (the
   "laws/government mod") relates; the path survey (Path A/B/C) and why custom mods win.
4. **[`CUSTOM-MODS.md`](CUSTOM-MODS.md)** — the own-repo + mod-mirror delivery pattern any verdict must fit.
5. **`CLAUDE.md`** — the "where you're running" section: sandbox can't launch the game; runtime claims need
   the box.
6. **`tools/mod-data/by-mod/valarian_conquest-4.2.1.1-neoforge-1.21.1.txt`** — the VC content inventory
   this mod reimagines **clean-room** (ideas only; VC is ARR — never read or reuse its code).

## 1. What's already settled — do not re-litigate

- **License:** MineColonies / Structurize / minecolonies-compatibility are **GPL-3.0** → **fork and addon
  are both legal.** A derivative is GPL-3.0 with public source (copyleft), which fits `theasshats`
  (`pcmc-killfeed` is GPL-3.0). **VC is All Rights Reserved** → clean-room, zero VC code.
- **Project shape:** Tier-2 custom mod in its **own repo**, mod-mirror packwiz delivery, **built/playtested
  on the box** — not in the web sandbox.
- **The guard-enforcement mechanism is solved** (§5): rank-flip → native guard aggro, addon-sufficient via
  public API. The "laws → guards" feature **reuses `pcmc-realms`' wanted-signal** — it is *not* part of
  this fork/addon question.
- **§5's bar for a fork (adopt it verbatim):** an **addon suffices** when MineColonies' **public API +
  targeted mixins** reach the behavior. A **fork is justified only** for guard/colony **behavior the public
  API provably can't express** — §5's examples: reshaping guard **target acquisition** (no public hook in
  `AbstractEntityAIGuard#checkForTarget`/`decide`), guards **pursuing past the colony border**, graduated/
  non-lethal response. And: **budget a fork as *replacing a core pack dependency*, not an addon-sized
  task** — it ships modified MineColonies to **every** player and must **continuously rebase** on the
  MineColonies snapshots this pack runs.

## 2. What's NEW here vs. §5 (the crux you must resolve)

§5 only had to make *one* player attackable by *existing* guards. This mod adds three things §5 never
needed — each is a fresh fork-vs-addon question:

| New capability | Why it might exceed the addon bar | Verdict hinges on |
| --- | --- | --- |
| **NPC-run joinable factions** | Needs colonies that exist and are owned **without a player** founding them, plus a **join/allegiance** hook. MineColonies colonies are player-anchored (Town Hall). | Does the API expose **programmatic colony creation + non-player ownership**? If not → fork, recast (pre-placed structure colonies / a parallel faction registry that *uses* colonies for guards), or cut. **Likely the hardest.** |
| **Offensive army mechanics** | Native guards are **defensive** (aggro within colony, leash to border). Marching/attacking on command and leaving the border is precisely §5's "reshape target acquisition / pursue past border" = **fork territory**. | Is there **any** public way to command guards offensively or inject targets/waypoints? Re-confirm §5's "no public target hook." Check the **raid system** as an alternative AI substrate (it already spawns mobile attacking NPCs). |
| **Faction equipment by allegiance** | Content (items) is trivial; **auto-equipping guards/citizens by faction** (slots + visuals) is the API question. | Does the API let you **set/lock a citizen's equipment + appearance** programmatically? Likely addon (citizens already equip gear), but confirm the by-allegiance hook. |

The **faction identity / reputation / join state** itself is almost certainly **addon** — hang it off the
entity in your **own `SavedData`**, exactly as `pcmc-realms` hangs `RealmGov` off the territory entity
(realms-scope §2). No fork needed for *state*; the fork question is only about *behavior*.

## 3. Spikes (the investigation — most need the box)

Mirror how Part 1/Part 2 spiked their APIs: **compile a real call against the pinned MineColonies snapshot
(`1.1.1327`), then confirm runtime behavior in-game.** For each, record: *API exists? → compiles? → works
at runtime? → addon or fork?*

1. **Colony lifecycle API.** Find the public surface for **creating** a colony and **setting its owner**
   (look in `com.minecolonies.api.colony` — `IColonyManager`, colony creation/claim entry points). Can a
   colony be created server-side, unowned/NPC-owned, and persist? → gates **NPC joinable factions**.
2. **Guard control API.** Re-confirm there is **no** public hook to inject a guard target or order an
   advance (§5: `AbstractEntityAIGuard` target acquisition is private). Check for any **patrol / guard-
   mode / rallypoint** API. → gates **offensive army mechanics** (the most likely fork trigger).
3. **Raid system as army substrate.** MineColonies ships a **raid** system (barbarians/pirates/etc.:
   `com.minecolonies.api.colony.colonyEvents` / the raid manager). Can it be **driven** to spawn a
   faction's attacking force at a target colony? If yes, **army mechanics may be an addon** that reuses raid
   AI instead of writing unit AI or forking guard AI. **High-leverage — check early.**
4. **Citizen equipment API.** Can you set/lock a citizen's worn equipment and visuals by job/allegiance? →
   gates **faction equipment auto-equip**.
5. **Prior art (read, don't reinvent):** **MineColonies: War 'N Taxes** (an *addon* adding taxes + war —
   evidence the public surface carries war-class features without a fork), **minecolonies-compatibility**
   (the pack already ships it; adds jobs/research, including a TACZ "Gunner" job — colony↔combat-mod glue
   prior art), and any faction/diplomacy MineColonies addons on the loader. What did they do via API vs.
   mixin?

## 4. Decision framework (the deliverable's spine)

- **Per-feature verdict**, not one global call: each of the four features →
  **`addon` / `addon+mixin` / `fork-required` / `cut-or-defer`**, with the spike evidence behind it.
- **Then a global recommendation.** Most likely outcome (state plainly if the spikes bear it out): **addon
  core + targeted mixins**, with a fork reserved **only** if **offensive army mechanics *and* NPC-colony
  ownership both** prove API-impossible **and** are must-haves. If only army-offense needs more than the
  API gives, prefer **mixins into the specific guard-AI methods** over a full fork (a mixin is reversible
  and snapshot-scoped; a fork is pack-wide and forever-rebased).
- **Apply §5's fork bar honestly:** "decide against that bar" — a fork must clear *replacing a core pack
  dependency*, not *adding a feature*.
- **Note the coexistence consequence** in the recommendation: addon keeps real MineColonies + its
  production-route role (`SYSTEMS.md` §3) + its snapshot updates; a fork forfeits all three and blocks
  running both. A fork also does **not** create guards where there are none (OPAC-only land, ship-realms).

## 5. Constraints & non-goals

- **Clean-room VC** — reimagine VC's *ideas* (factions, faction kit, army feel) with **zero** VC code. VC
  is ARR. The digest inventory is a *feature checklist*, not a code source.
- **Sandbox can't compile or playtest** — every "does it work at runtime" claim needs the box. Write the
  spikes as a checklist a maintainer (or a box-connected session) runs; don't assert runtime behavior from
  a static read.
- **GPL-3.0 copyleft** — whichever path, the result is GPL-3.0 with public source (a fork *must* publish
  modified MineColonies source; an addon publishes its own).
- **Cross-repo:** a session scoped to `theasshats/project-commonwealth` **cannot** write issues/code to the
  mod's own repo or to MineColonies (the realms-scope §13 hit this). Capture any MineColonies-side findings
  or mod-repo tasks as **notes in your deliverable** to be filed by a correctly-scoped session or by hand.

## 6. Deliverable

1. A **recommendation** — the per-feature fork/addon table (§4) + a global call, with spike evidence.
2. The **spike results** (API exists/compiles/runtime, per §3), or, if no build env, the spikes written up
   as a runnable checklist for the box.
3. An update to [`CONQUEST-MOD-SCOPE.md`](CONQUEST-MOD-SCOPE.md) §5 (Architecture stance) reflecting the
   verdict, and any new open questions surfaced.
4. Fold all of it into the governance branch (or its successor) — **do not** open a new PR unprompted
   (CLAUDE.md branch rules).

---

_Refs: [`CONQUEST-MOD-SCOPE.md`](CONQUEST-MOD-SCOPE.md); [`GOVERNANCE-REALMS-SCOPE.md`](GOVERNANCE-REALMS-SCOPE.md)
§2 (Part-1 API contract pattern), §5 (guard mechanism + addon-vs-fork decision aid — **the model for this
handoff**); [`GOVERNANCE.md`](GOVERNANCE.md); [`CUSTOM-MODS.md`](CUSTOM-MODS.md); `CLAUDE.md` ("where you're
running" — sandbox limits); MineColonies snapshot `1.1.1327` (the pinned coordinates to spike against);
MineColonies *War 'N Taxes* (addon war prior art)._
