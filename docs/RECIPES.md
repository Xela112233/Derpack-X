# Recipe Overhaul — "Made through Create"

> **Principle:** re-author recipes across the pack so crafting flows *through* Create — its
> **parts** and **methods** — weaving Create through every branch. **This is NOT a grind.** The bar
> is *coherence* ("of course you'd build this with Create parts"), never extra steps for their own sake.

## Guardrails (keeps it "Create-ified," not GregTech)
- One Create-part "gate" per tier is plenty — don't stack requirements.
- Keep yields fair; don't stretch a 1-step craft into five unless the item truly warrants it.
- `sequenced_assembly` / multi-stage processing is reserved for **marquee** items, not everyday gear.
- Ore→metal uses Create **crushing** (the ore-doubling players already like) — a *reward*, not a tax;
  the vanilla 1:1 smelt path stays.
- Never break a chain: if you `remove` a recipe, always add the Create replacement.

## Create "parts palette" (the recipe vocabulary)
Ingredients used to signal/gate tiers (verified in Create 6):
- **Tier-0 connective:** `create:andesite_alloy` — the universal "made with Create" marker
- **Shafts / cogs:** `create:shaft`, `create:cogwheel`, `create:large_cogwheel`
- **Sheets (pressed):** `create:iron_sheet`, `create:copper_sheet`, `create:brass_sheet`, `create:golden_sheet`
- **Casings:** `create:andesite_casing`, `create:copper_casing`, `create:brass_casing`
- **Tier-up gate:** `create:precision_mechanism`, `create:electron_tube`, `create:brass_ingot`
- **Cross-mod metals via tags:** `#c:ingots/<metal>`, `#c:plates/<metal>`, `#c:nuggets/<metal>`
  (`almost-unified` collapses duplicates, so tags are safe)
- **Higher metal tiers — already Create-made by addons, *use* don't rebuild:** bronze
  (`create_ironworks:bronze_ingot`/`_sheet`), steel (`#c:ingots/steel` — TFMG + Ironworks) and
  TFMG's steel parts (`tfmg:steel_mechanism`, `tfmg:steel_casing`, `tfmg:steel_cogwheel`), plus
  Create: Metalwork's molten-metal casting (`create:compacting` from molten buckets).

## Create "methods palette"
- `crushing` — ore → crushed (+byproduct): the ore-processing path
- `pressing` — ingot → sheet
- `mixing` — alloys / compounds (heated when the recipe needs it)
- `mechanical_crafting` — large multi-block crafts (marquee items)
- `deploying`, `filling` / `emptying`, `cutting`, `milling`, `compacting` — situational
- *reserved, sparing:* `sequenced_assembly` — special items only

## KubeJS mechanics
Files live in `kubejs/server_scripts/recipes/`, numbered for readability. Core pattern:
```js
ServerEvents.recipes(event => {
  event.remove({ output: 'somemod:fancy_gizmo' })           // kill the bypass recipe
  event.shaped('somemod:fancy_gizmo', ['SBS','BPB','SBS'], {  // re-add the Create way
    S: 'create:brass_sheet', B: '#c:ingots/brass', P: 'create:precision_mechanism'
  })
})
```
- Target by **tags** where possible so a change spans mods.
- `remove({ output })` kills all recipes for an item without needing its exact id.
- Verify in-game: JEI/EMI shows the new recipe; no recipe errors in the log.

## Scope & sequence (per-branch)
Each branch's key craftables get a Create-part gate or Create-method build:
1. **Metals & base materials — already Create-ified; leverage, don't rebuild.** Create (sheets,
   andesite_alloy), **Create: Metalwork** (molten-metal melting + `compacting`), **TFMG** (steel
   tier + sheets), and **Create: Ironworks** (bronze/steel + sheets) already route metal processing
   through Create. We *use* these parts as the gating vocabulary below — only add a missing
   pressing/crushing recipe if a specific pack metal lacks one. **The real work is steps 2–5.**
2. **Tools / armor / machines** — gate tier-ups behind `precision_mechanism` / sheets.
3. **Magic** (Ars, Iron's Spells, Occultism) — foci/altars want Create parts.
4. **Guns / vehicles** — TaCZ already routed via Immersive TaCZ; extend to Big Cannons / Aeronautics.
5. **Storage / QoL** — backpack & upgrade tiers behind Create parts.

Pure bypass-recipe removals happen alongside each branch.

## Ground truth
- Mods with jars on hand → recipes read directly for precise `remove` + replace.
- Everything else → the **`tools/mod-data/recipes/`** digest (CI) lists each recipe as
  `id | type | referenced-ids`; grep an item to find what to remove. Run **Actions ▸ Ground-truth
  digest** to populate it.
