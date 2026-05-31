#!/usr/bin/env bash
# Extract licensing-SAFE facts from mod jars into a text digest for offline reference.
#
# Captures ONLY functional identifiers + tiny structural metadata + DERIVED indexes:
#   - mod id/version/deps (META-INF/neoforge.mods.toml headers)
#   - block ids, item ids
#   - neoforge biome_modifier contents (small structural json)
#   - c: common tag membership
#   - recipe INDEX  : per recipe,     "id | type | referenced ids"  (derived, NOT verbatim)
#   - loot INDEX    : per loot table,  "id | type | referenced ids"  (derived, NOT verbatim)
# Never copies verbatim recipes/loot/lang/models/textures -> safe to commit.
#
# Single Python pass (one zip open per jar, in-process) — fast enough for the full pack.
# Usage: scripts/extract-mod-data.sh <jars-dir> [out-dir]
set -euo pipefail
python3 - "$@" <<'PY'
import sys, os, re, json, zipfile, glob

jars_dir = sys.argv[1]
out = sys.argv[2] if len(sys.argv) > 2 else "tools/mod-data"
bymod, recdir, lootdir = (os.path.join(out, d) for d in ("by-mod", "recipes", "loot"))
for d in (bymod, recdir, lootdir):
    os.makedirs(d, exist_ok=True)

RL        = re.compile(r'^[a-z0-9_.-]+:[a-z0-9_./-]+$')
BLOCK     = re.compile(r'^assets/([^/]+)/blockstates/(.+)\.json$')
ITEM      = re.compile(r'^assets/([^/]+)/models/item/(.+)\.json$')
CTAG      = re.compile(r'^data/c/tags/(.+)\.json$')
RECIPE    = re.compile(r'^data/([^/]+)/recipes?/(.+)\.json$')
LOOT      = re.compile(r'^data/([^/]+)/loot_tables?/(.+)\.json$')
TOMLKEY   = re.compile(r'^\s*(modId|version|displayName|license|versionRange)\s*=', re.I)

def refs(obj):
    found = set()
    def walk(x):
        if isinstance(x, str):
            if RL.match(x): found.add(x)
        elif isinstance(x, dict):
            for v in x.values(): walk(v)
        elif isinstance(x, list):
            for v in x: walk(v)
    walk(obj)
    return " ".join(sorted(found))

ores_idx = open(os.path.join(out, "INDEX-ores.txt"), "w", encoding="utf-8")
bmod_idx = open(os.path.join(out, "INDEX-biome-modifiers.txt"), "w", encoding="utf-8")
count = 0

for jar in sorted(glob.glob(os.path.join(jars_dir, "*.jar"))):
    name = os.path.basename(jar)[:-4]
    try:
        zf = zipfile.ZipFile(jar)
    except Exception:
        continue
    blocks, items, bmods, ctags, rec_lines, loot_lines, toml = [], [], [], [], [], [], []
    try:
        toml = [l.strip() for l in zf.read("META-INF/neoforge.mods.toml").decode("utf-8", "replace").splitlines() if TOMLKEY.match(l)]
    except Exception:
        pass
    for n in zf.namelist():
        m = BLOCK.match(n)
        if m: blocks.append(f"{m.group(1)}:{m.group(2)}"); continue
        m = ITEM.match(n)
        if m: items.append(f"{m.group(1)}:{m.group(2)}"); continue
        if "/neoforge/biome_modifier/" in n and n.endswith(".json"):
            try: bmods.append((n.split("/biome_modifier/")[-1], json.dumps(json.loads(zf.read(n)))))
            except Exception: pass
            continue
        m = CTAG.match(n)
        if m:
            try: ctags.append((m.group(1), json.dumps(json.loads(zf.read(n)))))
            except Exception: pass
            continue
        m = RECIPE.match(n)
        if m:
            try: obj = json.loads(zf.read(n))
            except Exception: continue
            t = obj.get("type", "?") if isinstance(obj, dict) else "?"
            rec_lines.append(f"{m.group(1)}:{m.group(2)} | {t} | {refs(obj)}")
            continue
        m = LOOT.match(n)
        if m:
            try: obj = json.loads(zf.read(n))
            except Exception: continue
            t = obj.get("type", "?") if isinstance(obj, dict) else "?"
            loot_lines.append(f"{m.group(1)}:{m.group(2)} | {t} | {refs(obj)}")
    zf.close()
    with open(os.path.join(bymod, name + ".txt"), "w", encoding="utf-8") as f:
        f.write(f"# {name}\n\n## mod metadata (id / version / deps)\n" + "\n".join(toml))
        f.write("\n\n## blocks\n" + "\n".join(sorted(set(blocks))))
        f.write("\n\n## items\n" + "\n".join(sorted(set(items))))
        f.write("\n\n## neoforge biome_modifiers\n" + "\n".join(f"{a} :: {b}" for a, b in sorted(bmods)))
        f.write("\n\n## c: common tag membership\n" + "\n".join(f"{a} :: {b}" for a, b in sorted(ctags)) + "\n")
    if rec_lines:
        open(os.path.join(recdir, name + ".txt"), "w", encoding="utf-8").write("\n".join(sorted(rec_lines)) + "\n")
    if loot_lines:
        open(os.path.join(lootdir, name + ".txt"), "w", encoding="utf-8").write("\n".join(sorted(loot_lines)) + "\n")
    for b in sorted(set(blocks)):
        if "_ore" in b: ores_idx.write(f"{name}: {b}\n")
    for a, b in sorted(bmods):
        bmod_idx.write(f"{name}\t{a}\t{b}\n")
    count += 1

ores_idx.close(); bmod_idx.close()
p = os.path.join(out, "INDEX-ores.txt")
open(p, "w", encoding="utf-8").write("\n".join(sorted(open(p, encoding="utf-8").read().splitlines())) + "\n")
open(os.path.join(out, "README.txt"), "w", encoding="utf-8").write(
    "Mod-data digest — IDs / tags / biome-modifiers + recipe & loot INDEXES (derived, not verbatim).\n"
    f"Generated by scripts/extract-mod-data.sh over {count} jars.\n\n"
    "  by-mod/<jar>.txt            blocks, items, biome_modifiers, c: tags, deps\n"
    "  recipes/<jar>.txt           per-recipe 'id | type | referenced ids'\n"
    "  loot/<jar>.txt              per-loot-table 'id | type | referenced ids'\n"
    "  INDEX-ores.txt              every *_ore blockstate across all mods\n"
    "  INDEX-biome-modifiers.txt   every neoforge biome_modifier across all mods\n\n"
    "Derived indexes only (no verbatim recipe/loot/lang/assets). Pair with in-game EMI for code-injected ores.\n")
print(f"Extracted {count} mods -> {out}")
PY
