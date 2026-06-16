# Finding Ores: Prospecting

Project Commonwealth uses GregTech-style ore generation through the GTMOGS mod (Gregtech: Modern Ore
Generation Standalone). Ores are concentrated in large veins of related materials instead of
being scattered evenly, and most vanilla overworld ore generation is turned off
(`removeVanillaOreGen` and `removeVanillaLargeOreVeins` both default on). Veins, plus a thin
layer of small ores, are the only sources. Prospecting is how you locate the veins and put
them on your map.

> The details below are taken from the shipped GTMOGS jar (`gtmogs-1.21.1-1.0.6`). Where a number
> or message is quoted, it is the mod's actual default or string, not an estimate.

## How prospecting works

You prospect by **right-clicking** a vein ore block. Breaking it does nothing — there is no
block-break trigger and no prospector tool item in this version; the only trigger is a
right-click (use) on a block that belongs to a registered vein.

When you right-click such a block, the server scans for veins within `oreBlockProspectRange`
(default **24** blocks) and adds any it finds to your prospecting records. A chat message
confirms the result:

- `Prospected %d new veins!` when several are found.
- `Prospected %s!` naming the vein when a single new one is found.

If you get no message, either you broke the block instead of right-clicking it, or that block
type is not part of any vein. An empty hand is the simplest way to right-click without doing
anything else.

## Small ores are indicators

Single ore blocks are scattered across the overworld: coal, iron, copper, tin, zinc, gold,
nickel, redstone, lapis, and the occasional diamond. They are enough for a first set of
tools, but their real job is to point at veins. When you find one, right-click it and check
your map. (Silver is not among them — it only comes from veins, in silver country.)

## Step by step

1. Find an ore block — a scattered small ore, or any vein ore you have already exposed.
2. Right-click the ore block (empty hand is simplest). Breaking it will not prospect.
3. Read the chat confirmation (`Prospected N new veins!`). This is the reliable signal that
   prospecting happened.
4. Open your map and enable the ore-vein layer (below) to see where the located veins are.
5. Travel to a vein and dig down. Veins are several blocks tall and hold a few related ores in
   layers, so one shaft usually returns more than the headline ore.
6. When a vein is mined out, mark it depleted so it stops cluttering your map.

## The map overlay

Prospected veins are drawn by whichever supported map mod you have. GTMOGS ships integrations
for Xaero's (this pack's map mod), FTB Chunks, and JourneyMap; all three default on in the
config and need a game restart to toggle. The integration only **displays** veins you have
already prospected — it does not prospect for you, and an empty overlay almost always means you
have not right-clicked an ore yet, not that anything is broken.

Look for these toggles in your map mod's interface (the exact placement depends on the map mod;
in JourneyMap they appear under a "Prospection layers" option group):

- "Show GT Ore Veins" — the main overlay.
- "Show Depleted Veins" / "Hide Depleted Veins" — controls whether mined-out veins still draw.
- "Show Bedrock Fluid Veins" — a separate layer (not used by this pack's overworld setup).

## Managing veins on the map

Interacting with a vein marker gives two actions:

- "Mark as Depleted" — flags a mined-out vein; it then shows as `Depleted` and can be hidden
  with the toggle above.
- "Toggle Waypoint" — creates or removes a waypoint for the vein. Creating one confirms with
  `Created waypoint named %s!`.

## Sharing with your group

Run `/gtmogs share_prospection_data <player>` to send the veins you have located to a teammate.
They receive a notification (`<you> is sharing prospecting data with you!`). On a shared world
this is part of how a group divides work — a player who has mapped a region can hand that
knowledge on, or trade it.

## Looking veins up in JEI

GTMOGS adds an "Ore Vein Diagram" category to JEI. Looking up an ore there shows the vein it
belongs to along with its Weight, Chance, Dimensions, and Spawn Range, so you can confirm which
biome and depth to search before you set out.

## Where to look

Veins are regional — the biome at the vein's depth decides what generates. A rough guide:

- Mountains and hills: iron, magnetite, lead, silver, nickel, palladium (deep).
- Badlands and savanna: copper, gold, thorium (deep), bauxite (savanna).
- Desert: redstone, lithium.
- Forest and taiga: coal.
- Plains: zinc, tin, lignite, fireclay.
- Swamps: lignite, fireclay.
- Jungle: diamond (deep), jade, bauxite.
- Snowy biomes: lapis.
- Mountain peaks: emerald.
- Coasts and oceans: salt.

A few regions also hide **mega veins** — iron in mountains, copper in badlands, coal in taiga and
forest. They are far larger than ordinary veins and far rarer (on the order of one per few hundred
chunks of their region). Prospecting reveals them the same way, and a single find can supply a
settlement or anchor a trade route.

Depth matters as much as biome. Marquee ores such as diamond, thorium, and mithril sit deep;
emerald sits high in peaks. The JEI diagram above gives the exact spawn range per vein.

## Configuration knobs

These live in the GTMOGS config (`config/gtmogs.yaml`). The pack ships its own values where they
differ from the jar defaults:

- `oreBlockProspectRange` = 24 — how far a right-click scans for veins.
- `oreVeinGridSize` = 5 and `oreVeinRandomOffset` = 20 (pack values; the jar defaults are 3 and 12)
  — veins sit on a wide grid with a random jitter, so once you find one the next sits at a roughly
  predictable interval.
- `xaerosMapIntegration` / `ftbChunksIntegration` / `journeyMapIntegration` = all on.
- `oreIconSize` = 32, and the pack sets `borderColor` = `material`, so vein markers outline in the
  ore's own colour.
- `removeVanillaOreGen` / `removeVanillaLargeOreVeins` = on (vanilla ore is disabled, as above).

## Tips

- Veins generate on a regular grid, so once you find one others sit at predictable intervals.
- A single vertical shaft through the center of a vein covers every depth it could be at, so
  there is no need to strip-mine.
- Caves cut through veins too; right-click ore you find while caving the same way.

## The prospecting kit

Create: Ore Excavation adds three tools that work alongside right-click prospecting, and all three
are early-game crafts:

- **Vein finder** — built around an electron tube, it points you toward buried deposits from the
  surface.
- **Vein atlas** — a clipboard-bound survey book that keeps a record of what you have found.
- **Sample drill** — a small brass machine that bores a test column so you can check what is below
  without committing to a shaft.

## Two vein systems

The pack has two different kinds of "vein," and they are found and used differently:

- **Ore veins** (the GregTech-style system above) are the pack's ordinary ore supply: regional,
  hand-mined, revealed by right-clicking ore.
- **Drilling veins** are separate, deeper deposits that only Create: Ore Excavation's machines can
  tap. The vein finder, atlas, and sample drill locate them. You cannot dig these out by hand.

## End-game: industrial vein drilling

Tapping a drilling vein is an end-game operation. The drilling machine and its fluid-pumping
counterpart, the extractor, require etched circuit boards and electric motors from the top of the
industrial chain, and they only run on serious rotational power. Drill heads are steel-bodied,
assembled on the Mechanical Crafter, and wear out with use; the heavier veins also consume fluid
(the netherite vein drinks lava). Expect to spend drills, power, and logistics to run one — finding
veins is free, but emptying the planet is not.

## Admin / debug

Operators can force a vein for testing with `/gtmogs place_vein <vein> <pos>`, which confirms
with `Placed vein %s at position %s` (or a failure message). This places a vein directly and is
a fast way to verify a vein definition without searching the world.
