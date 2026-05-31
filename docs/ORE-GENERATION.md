# Ore Generation — Design & Configuration Guide

> GregTech-style ores: **rare, large veins, biome-specific** — so regions specialize and
> trade matters. We **thin** vanilla ore generation but keep a small **starter trickle** so
> new players near spawn aren't hard-stuck if they can't find a vein yet.
>
> **Collaborator:** fill in / edit the **Distribution table** below (multiple ores per biome
> is fine). Pick biomes from the **Terralith biome list** at the bottom, or — easier — target
> a **biome tag** that covers many biomes at once.

---

## How it works (mechanism)

KubeJS loads `kubejs/data/` as a **global datapack** — no extra mod, applies to every world.
All files below live under `kubejs/data/`. Each ore is up to **4 small JSON files**:

| File | Purpose | Key dial |
|---|---|---|
| `derpack/worldgen/configured_feature/vein_<ore>.json` | the ore body | `size` = blocks per vein |
| `derpack/worldgen/placed_feature/vein_<ore>.json` | how often / how deep | `rarity_filter.chance`, `height_range` |
| `derpack/neoforge/biome_modifier/vein_<ore>.json` | **which biomes** get it | `biomes` (tag or list) |
| `minecraft/worldgen/placed_feature/ore_<ore>_*.json` | **thin** vanilla default | `rarity_filter.chance` (higher = rarer) |

A complete worked example is already in the repo: **iron** (mountains) and **copper** (badlands).
Copy that pattern for each new ore.

## Tuning knobs

- **`size`** — vein volume. Vanilla ≈ 8–9; our veins use **24–48** for the "big deposit" feel.
- **`rarity_filter.chance: N`** — one placement attempt per **1/N chunks**. Higher N = rarer.
  Veins use 5–14; thinned vanilla uses 9–12.
- **`height_range`** — `trapezoid` (concentrated mid-band) or `uniform` (even). Match real ore depth.
- **`discard_chance_on_air_exposure`** — 0.0 = veins show in caves/cliffs (good for findability).
- **Multiple ores in one biome:** just add several `biome_modifier` files all targeting the same
  biome/tag (e.g. mountains get both iron **and** emerald). No conflict.
- **Starter trickle:** thin vanilla to `chance` ~10–12 (rare) but **do not remove it** — keep a
  thin coal/iron/copper baseline everywhere so newbies can bootstrap; veins are the bulk source.

## Biome targeting: tags (easy) vs IDs (precise)

Prefer **vanilla biome tags** — Terralith routes its biomes into them via `#terralith:reference/*`
tags (verified), so one vanilla tag covers all the matching Terralith biomes too:

`#minecraft:is_mountain` · `is_hill` · `is_forest` · `is_taiga` · `is_jungle` · `is_savanna`
· `is_badlands` · `is_beach` · `is_deep_ocean` · `is_river` · `is_overworld`

> Note (verified): there is **no** `is_desert`, `is_ocean` (shallow), or `is_snowy` tag — for those,
> list specific biomes (e.g. `minecraft:desert`, `terralith:desert_oasis`, snowy biomes by name).

Use a **specific biome ID** only for special cases (e.g. diamonds *only* in `terralith:amethyst_rainforest`).
You can also list several: `"biomes": ["terralith:alpine_highlands", "terralith:rocky_mountains"]`.

---

## Ores to place

**Vanilla (overworld):** coal, iron, copper, gold, redstone, lapis, diamond, emerald.

**Modded (overworld) — installed in this pack** *(block IDs VERIFIED from the jars):*
- **Zinc** — Create — `create:zinc_ore`, `create:deepslate_zinc_ore`
- **Thorium** — Create: New Age — `create_new_age:thorium_ore` (no deepslate variant). *New Age adds
  **thorium, not uranium**; it also adds `create_new_age:magnetite_block` — vein it only if we want it.*
- **Silver** — Occultism — `occultism:silver_ore`, `occultism:silver_ore_deepslate`

**Nether / dimensional (optional, later pass):** `nether_gold_ore`, `nether_quartz_ore`,
`ancient_debris`; Occultism `occultism:iesnium_ore` (nether).

### Disabling each mod's default ore gen (so our veins are the only source)

Vanilla **and** every ore mod inject their ore into `#minecraft:is_overworld` at worldgen step
`underground_ores`. To make an ore vein-based + biome-specific, **override the mod's biome modifier**
(replace its `features` with `minecraft:empty`) and add our own vein. Verified override targets
(create the same path under `kubejs/data/…` to shadow the mod's file):

| Ore | File to override |
|---|---|
| Zinc | `create/neoforge/biome_modifier/zinc_ore.json` |
| Thorium | `create_new_age/neoforge/biome_modifier/thorium_ore.json` |
| Magnetite | `create_new_age/neoforge/biome_modifier/magnetite_block.json` |
| Silver | `occultism/neoforge/biome_modifier/add_ore_silver.json` + `…/add_ore_silver_deepslate.json` |

Leave the nether ones (`occultism` iesnium, Create `striated_ores_nether`) alone unless we do a nether pass.
Create also has `striated_ores_overworld` (its stone-stripe ores) — decide later whether to keep.

---

## Distribution table — COLLABORATOR FILLS THIS IN

Pre-filled with a starting proposal. Edit freely; add rows; assign several ores to the same biome.

| Ore | Biome(s) / tag | Size | Rarity (1/N) | Y-band | Notes |
|---|---|---|---|---|---|
| Coal | `#is_taiga`, `#is_forest` | 40 | 5 | 0…136 | common-ish |
| Iron ✅ | `#is_mountain` | 48 | 8 | -24…64 | example built |
| Copper ✅ | `#is_badlands` | 48 | 7 | -16…112 | example built |
| Gold | `#is_savanna`, `#is_badlands` | 40 | 9 | -48…40 | |
| Redstone | desert biomes (by name) | 40 | 8 | -60…-8 | no `is_desert` tag |
| Lapis | snowy biomes (by name) | 28 | 11 | -32…40 | no `is_snowy` tag |
| Diamond | jungle + deepslate | 24 | 14 | -60…-16 | endgame rare |
| Emerald | mountain peaks | 16 | 9 | 64…256 | |
| Zinc (Create) | `#is_overworld` plains/savanna | 40 | 7 | 0…80 | |
| Thorium (New Age) | ? (deep, rare) | 24 | 14 | -60…-20 | radioactive fuel — keep rare |
| Silver (Occultism) | ? | 32 | 10 | -32…48 | |

**Starter-trickle floor (thinned vanilla, everywhere):** coal/iron/copper at `chance` 10–12,
others 12+. Tune so a fresh player can find *some* iron within ~10 min, but real supply = veins.

---

## Terralith biome list — VERIFIED (95 biomes, Terralith 2.5.8)

> Exact IDs from the jar (`data/terralith/worldgen/biome/`). Vanilla biomes also still generate.
> Easiest path is to target a **vanilla tag** (above); use these IDs only for specific placements.

**Surface (`terralith:…`):** `alpha_islands` · `alpha_islands_winter` · `alpine_grove` · `alpine_highlands`
· `amethyst_canyon` · `amethyst_rainforest` · `ancient_sands` · `arid_highlands` · `ashen_savanna`
· `basalt_cliffs` · `birch_taiga` · `blooming_plateau` · `blooming_valley` · `brushland` · `bryce_canyon`
· `caldera` · `cloud_forest` · `cold_shrubland` · `deep_warm_ocean` · `desert_canyon` · `desert_oasis`
· `desert_spires` · `emerald_peaks` · `forested_highlands` · `fractured_savanna` · `frozen_cliffs`
· `glacial_chasm` · `granite_cliffs` · `gravel_beach` · `gravel_desert` · `haze_mountain` · `highlands`
· `hot_shrubland` · `ice_marsh` · `jungle_mountains` · `lavender_forest` · `lavender_valley` · `lush_desert`
· `lush_valley` · `mirage_isles` · `moonlight_grove` · `moonlight_valley` · `orchid_swamp` · `painted_mountains`
· `red_oasis` · `rocky_jungle` · `rocky_mountains` · `rocky_shrubland` · `sakura_grove` · `sakura_valley`
· `sandstone_valley` · `savanna_badlands` · `savanna_slopes` · `scarlet_mountains` · `shield` · `shield_clearing`
· `siberian_grove` · `siberian_taiga` · `skylands_autumn` · `skylands_spring` · `skylands_summer`
· `skylands_winter` · `snowy_badlands` · `snowy_cherry_grove` · `snowy_maple_forest` · `snowy_shield`
· `steppe` · `stony_spires` · `temperate_highlands` · `tropical_jungle` · `valley_clearing` · `volcanic_crater`
· `volcanic_peaks` · `warm_river` · `warped_mesa` · `white_cliffs` · `white_mesa` · `windswept_spires`
· `wintry_forest` · `wintry_lowlands` · `yellowstone` · `yosemite_cliffs` · `yosemite_lowlands`

**Cave (`terralith:cave/…` — the `cave/` is part of the ID):** `cave/andesite_caves` · `cave/deep_caves`
· `cave/diorite_caves` · `cave/frostfire_caves` · `cave/fungal_caves` · `cave/granite_caves`
· `cave/infested_caves` · `cave/mantle_caves` · `cave/thermal_caves` · `cave/tuff_caves`
· `cave/underground_jungle`  *(+ vanilla `lush_caves`, `dripstone_caves`, `deep_dark`)*

---

## Applying & testing

1. Files live in `kubejs/data/` → loaded as a global datapack automatically.
2. **Worldgen only affects NEW chunks** — test in a fresh world or unexplored area; `/reload`
   won't regenerate existing terrain.
3. Playtest checklist: veins appear in the right biomes at the expected size/rarity; vanilla ore
   is scarce but a starter amount still exists; no `worldgen`/`biome_modifier` errors in the log.

> **Verified from jars:** modded ore IDs (zinc/thorium/silver) and the per-mod biome-modifier override
> paths; the full 95-biome Terralith list; Terralith routes biomes into vanilla tags. **Still to confirm
> in-game:** that KubeJS loads `neoforge/biome_modifier` overrides from `kubejs/data/`.
