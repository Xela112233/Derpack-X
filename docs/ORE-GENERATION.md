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

Prefer **vanilla biome tags** — Terralith biomes inherit them, so one tag covers many biomes:

`#minecraft:is_mountain` · `is_hill` · `is_forest` · `is_taiga` · `is_jungle` · `is_savanna`
· `is_badlands` · `is_beach` · `is_ocean` · `is_river` · `is_nether` · `is_overworld`

Use a **specific biome ID** only for special cases (e.g. diamonds *only* in `terralith:amethyst_rainforest`).
You can also list several: `"biomes": ["terralith:alpine_highlands", "terralith:rocky_mountains"]`.

---

## Ores to place

**Vanilla (overworld):** coal, iron, copper, gold, redstone, lapis, diamond, emerald.

**Modded (overworld) — installed in this pack** *(block IDs to be VERIFIED from jars once network allows; best-guess shown):*
- **Zinc** — Create — `create:zinc_ore`, `create:deepslate_zinc_ore`
- **Uranium** — Create: New Age — `create_new_age:uranium_ore` (+ deepslate?) *(verify)*
- **Silver** — Occultism — `occultism:silver_ore`, `occultism:silver_ore_deepslate` *(verify)*

**Nether / dimensional (optional, later pass):** `nether_gold_ore`, `nether_quartz_ore`,
`ancient_debris`; Occultism `iesnium_ore` (Other dimension).

---

## Distribution table — COLLABORATOR FILLS THIS IN

Pre-filled with a starting proposal. Edit freely; add rows; assign several ores to the same biome.

| Ore | Biome(s) / tag | Size | Rarity (1/N) | Y-band | Notes |
|---|---|---|---|---|---|
| Coal | `#is_taiga`, `#is_forest` | 40 | 5 | 0…136 | common-ish |
| Iron ✅ | `#is_mountain` | 48 | 8 | -24…64 | example built |
| Copper ✅ | `#is_badlands` | 48 | 7 | -16…112 | example built |
| Gold | `#is_savanna`, `#is_badlands` | 40 | 9 | -48…40 | |
| Redstone | desert biomes | 40 | 8 | -60…-8 | |
| Lapis | snowy/icy biomes | 28 | 11 | -32…40 | |
| Diamond | jungle + deepslate | 24 | 14 | -60…-16 | endgame rare |
| Emerald | mountain peaks | 16 | 9 | 64…256 | |
| Zinc (Create) | `#is_overworld` plains/savanna | 40 | 7 | 0…80 | |
| Uranium (New Age) | ? | 24 | 14 | -60…-20 | strong — keep rare |
| Silver (Occultism) | ? | 32 | 10 | -32…48 | |

**Starter-trickle floor (thinned vanilla, everywhere):** coal/iron/copper at `chance` 10–12,
others 12+. Tune so a fresh player can find *some* iron within ~10 min, but real supply = veins.

---

## Terralith biome list (choose from these)

> From **Terralith 2.5.x** knowledge — **cross-check against the jar** (`data/terralith/worldgen/biome/`)
> when jar access is available. Vanilla biomes also still generate. Tip: most of these carry the
> vanilla tags above, so tag-targeting usually suffices.

**Hot / arid:** `ancient_sands` · `arid_highlands` · `ashen_savanna` · `brushland` · `bryce_canyon`
· `desert_canyon` · `desert_oasis` · `desert_spires` · `fractured_savanna` · `hot_shrubland`
· `lush_desert` · `red_oasis` · `sandstone_valley` · `savanna_badlands` · `warped_mesa`
· `white_mesa` · `yellowstone`

**Temperate / forest:** `blooming_valley` · `blooming_plateau` · `cloud_forest` · `forested_highlands`
· `haze_mountain` · `highlands` · `hot_shrubland` · `lavender_forest` · `lavender_valley`
· `lush_valley` · `moonlight_grove` · `moonlight_valley` · `orchid_swamp` · `rocky_jungle`
· `sakura_grove` · `sakura_valley` · `shrubland` · `skylands` · `temperate_highlands`
· `tropical_jungle` · `tropical_rainforest` · `valley_clearing`

**Mountainous:** `alpine_grove` · `alpine_highlands` · `emerald_peaks` · `granite_cliffs`
· `jungle_mountains` · `mountain_steppe` · `painted_mountains` · `rocky_mountains`
· `scarlet_mountains` · `stony_spires` · `windswept_spires` · `yosemite_cliffs` · `yosemite_lowlands`

**Cold / snowy:** `cold_shrubland` · `frozen_cliffs` · `glacial_chasm` · `ice_marsh` · `rocky_shrubland`
· `shield` · `shield_clearing` · `siberian_grove` · `siberian_taiga` · `snowy_badlands`
· `snowy_cliffs` · `snowy_maple_forest` · `snowy_shrubland` · `steppe` · `wintry_forest` · `wintry_lowlands`

**Volcanic / unique:** `amethyst_canyon` · `amethyst_rainforest` · `basalt_cliffs` · `caldera`
· `mirage_isles` · `volcanic_crater` · `volcanic_peaks` · `white_cliffs`

**Coastal / aquatic:** `alpha_islands` · `alpha_islands_winter` · `gravel_beach` · `warm_river`

**Cave biomes:** `andesite_caves` · `crystal_caves` · `deep_caves` · `desert_caves` · `diorite_caves`
· `fractured_caves` · `frostfire_caves` · `granite_caves` · `ice_caves` · `infested_caves`
· `mantle_caves` · `thermal_caves` · `tuff_caves` · `underground_jungle` (+ vanilla `lush_caves`, `dripstone_caves`, `deep_dark`)

---

## Applying & testing

1. Files live in `kubejs/data/` → loaded as a global datapack automatically.
2. **Worldgen only affects NEW chunks** — test in a fresh world or unexplored area; `/reload`
   won't regenerate existing terrain.
3. Playtest checklist: veins appear in the right biomes at the expected size/rarity; vanilla ore
   is scarce but a starter amount still exists; no `worldgen`/`biome_modifier` errors in the log.

> **Open verification items (blocked until jar access / playtest):** exact modded ore block IDs
> (zinc/uranium/silver); whether KubeJS loads `neoforge/biome_modifier` from its data folder
> (expected yes); final Terralith biome IDs.
