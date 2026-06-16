#!/usr/bin/env python3
"""Validate the GTMOGS ore-vein layer (kubejs/data/pcmc/gtmogs/ore_vein/).

Catches the worldgen mistakes that are valid JSON and therefore sail through
pr-checks green but crash or silently misbehave at world creation:

  - a band keyed `target` instead of the required plural `targets` (crashes worldgen)
  - missing bands / missing top-level keys
  - placed block states that exist in no installed mod (checked against
    tools/mod-data digests; minecraft: ids are trusted). Caveat: the digests are
    ASSET-derived (blockstates/item models), so an id passing this check can still
    be an unregistered datagen ghost (the samurai-silver v0.7.1 crash) — corroborate
    new ids against the mod's lang/loot/recipes before tagging them anywhere.
  - a vein whose `biomes` tag file does not exist
  - small-ore features whose configured/placed halves or biome-modifier entry
    are out of sync (a placed_feature pointing at a missing configured_feature
    crashes world creation)

Also prints the per-region weight shares (including megas) so a tuning change
can be sanity-checked against the intended rarity math without booting the game.

Run from the repo root:  python3 scripts/validate-ore-veins.py
Exit code 0 = clean, 1 = problems found.
"""
import glob
import json
import os
import re
import sys

VEIN_DIR = "kubejs/data/pcmc/gtmogs/ore_vein"
TAG_DIR = "kubejs/data/pcmc/tags/worldgen/biome"
SMALL_CONF = "kubejs/data/pcmc/worldgen/configured_feature"
SMALL_PLACED = "kubejs/data/pcmc/worldgen/placed_feature"
SMALL_MODIFIER = "kubejs/data/pcmc/neoforge/biome_modifier/small_ores.json"
BANDS = ("primary", "secondary", "between", "sporadic")

problems = []


def known_ids():
    ids = set()
    for f in glob.glob("tools/mod-data/by-mod/*.txt"):
        for m in re.finditer(r"\b([a-z0-9_]+:[a-z0-9_/]+)\b", open(f, errors="ignore").read()):
            ids.add(m.group(1))
    return ids


BAND_FRACTION = {"primary": 3 / 9, "secondary": 3 / 9, "between": 2 / 9, "sporadic": 1 / 9}


def material(state):
    """minecraft:deepslate_iron_ore -> iron; create_new_age:magnetite_block -> magnetite; tfmg:bauxite -> bauxite."""
    name = state.split(":", 1)[1]
    for prefix in ("deepslate_",):
        if name.startswith(prefix):
            name = name[len(prefix):]
    for suffix in ("_ore", "_block"):
        if name.endswith(suffix):
            name = name[: -len(suffix)]
    return name


def check_veins(ids):
    weights = {}
    compositions = {}
    for f in sorted(glob.glob(f"{VEIN_DIR}/*.json")):
        name = os.path.basename(f)[:-5]
        d = json.load(open(f))
        for key in ("cluster_size", "density", "weight", "biomes", "height_range", "generator"):
            if key not in d:
                problems.append(f"{name}: missing top-level `{key}`")
        gen = d.get("generator", {})
        comp = {}
        for band in BANDS:
            b = gen.get(band)
            if b is None:
                problems.append(f"{name}: missing band `{band}`")
                continue
            if "targets" not in b:
                problems.append(f"{name}/{band}: no `targets` (plural!) — a singular `target` "
                                f"parses as JSON but CRASHES world creation. keys={list(b)}")
                continue
            if "layers" not in b:
                problems.append(f"{name}/{band}: missing `layers`")
            for t in b["targets"]:
                if "target" not in t or "state" not in t:
                    problems.append(f"{name}/{band}: target entry missing inner `target`/`state`")
                    continue
                state = t["state"].get("Name", "?")
                if state.split(":")[0] != "minecraft" and state not in ids:
                    problems.append(f"{name}/{band}: block `{state}` not found in tools/mod-data digests")
            first = b["targets"][0].get("state", {}).get("Name") if b["targets"] else None
            if first:
                mat = material(first)
                comp[mat] = comp.get(mat, 0) + BAND_FRACTION[band]
        biomes = d.get("biomes", "")
        if isinstance(biomes, str) and biomes.startswith("#pcmc:"):
            tag = biomes.split(":", 1)[1]
            if not os.path.exists(f"{TAG_DIR}/{tag}.json"):
                problems.append(f"{name}: biome tag file {TAG_DIR}/{tag}.json does not exist")
        else:
            problems.append(f"{name}: `biomes` should be exactly one `#pcmc:vein_*` tag (got {biomes!r}) "
                            f"— edit the tag file to move a vein, never this field")
        weights[name] = d.get("weight", 0)
        compositions[name] = comp
    return weights, compositions


def check_small_ores():
    conf = {os.path.basename(p)[:-5] for p in glob.glob(f"{SMALL_CONF}/small_*.json")}
    placed = {os.path.basename(p)[:-5] for p in glob.glob(f"{SMALL_PLACED}/small_*.json")}
    listed = {f.split(":", 1)[1] for f in json.load(open(SMALL_MODIFIER))["features"]}
    for missing in placed - conf:
        problems.append(f"small ores: placed_feature {missing} has no configured_feature (CRASHES world creation)")
    for missing in conf - placed:
        problems.append(f"small ores: configured_feature {missing} has no placed_feature (dead file)")
    for missing in listed - placed:
        problems.append(f"small ores: biome modifier lists {missing} but no placed_feature exists (CRASHES)")
    for orphan in placed - listed:
        problems.append(f"small ores: {orphan} exists but is not in the small_ores biome modifier (never generates)")


def region_shares(weights, compositions):
    # region -> member veins, mirrors the vein_* tag files; update when regions move.
    # Singletons get their own rows: every anchor in that region is that vein.
    regions = {
        "hills": ["tin", "lead", "magnetite", "silver", "gold_hills"],
        "mountains": ["iron", "lead", "magnetite", "nickel", "emerald", "palladium", "iron_mega", "diamond_mountains", "redstone_mountains"],
        "badlands": ["copper", "gold", "lithium", "thorium", "copper_mega"],
        "taiga": ["coal", "coal_mega"],
        "forest": ["coal", "coal_mega", "gold_forest", "copper_forest"],
        "savanna": ["zinc", "gold", "nickel", "bauxite"],
        "plains": ["tin", "zinc", "lignite"],
        "desert": ["redstone", "lithium"],
        "jungle": ["bauxite", "jade", "diamond", "lapis_jungle"],
        "snowy [singleton]": ["lapis"],
        "warm oceans/coasts": ["salt", "copper_warm_ocean", "zinc_warm_ocean", "lead_warm_ocean", "gold_warm_ocean", "silver_warm_ocean"],
        "deep/cold oceans": ["salt", "nickel_deep_ocean", "copper_deep_ocean", "palladium_deep_ocean", "mithril_deep_ocean"],
        "swamps": ["lignite", "iron_swamp"],
        "Terralith specials [singleton]": ["mithril"],
    }
    print("\nper-region shares (anchor every grid^2 chunks; grid 5 = 25).")
    print("ore output model: bands P/S/B/Sp = 3/3/2/1 of 9 layers (between/sporadic are engine-")
    print("reduced -> upper bounds; vein-size differences NOT folded in - megas are ~10-25x volume):")
    for region, members in regions.items():
        total = sum(weights.get(m, 0) for m in members)
        if not total:
            continue
        parts = ", ".join(f"{m} {weights.get(m, 0) / total * 100:.0f}%" for m in members)
        print(f"  {region} (sum {total}): {parts}")
        ores = {}
        for m in members:
            share = weights.get(m, 0) / total
            for ore, frac in compositions.get(m, {}).items():
                ores[ore] = ores.get(ore, 0) + share * frac
        ranked = sorted(ores.items(), key=lambda kv: -kv[1])
        print("      ore output: " + " · ".join(f"{o} {v * 100:.0f}%" for o, v in ranked))


def main():
    ids = known_ids()
    weights, compositions = check_veins(ids)
    check_small_ores()
    print(f"{len(weights)} veins checked against {len(ids)} known ids:",
          "ALL CLEAN" if not problems else f"{len(problems)} PROBLEM(S)")
    for p in problems:
        print("  !!", p)
    region_shares(weights, compositions)
    sys.exit(1 if problems else 0)


if __name__ == "__main__":
    main()
