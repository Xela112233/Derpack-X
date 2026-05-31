# Modlist & Path-Forward Audit

_356 manifests, audited against the stated goals (CLAUDE.md / DESIGN.md / Phase-1 plan).
Slug + knowledge based — items marked ⚠ need confirming from the jar/EMI; the ground-truth
CI digest will nail those. This is a planning artifact, not a final decision list._

## Stated goals (the yardstick)
- **Create-driven survival economy "in the spirit of Eco":** scarce, regional, large-vein ores →
  goods manufactured through Create → sold for currency (**Numismatics**).
- **Heavy / GregTech-like recipe conversion** — route progression through Create machines.
- **Built around Create Aeronautics** (airships / planes / vehicles).
- **Guns = TaCZ + Create Big Cannons.**
- Friend-group, pre-alpha. Joke mods allowed *if intentional*.

## On-track check — ✅ mostly
| Goal | State |
|---|---|
| Scarce regional veins | ✅ done — 20 ores (PR #56, playtest #58) |
| Guns = TaCZ + Big Cannons | ✅ TaCZ + Immersive TaCZ + Armorer (#55); Big Cannons + Adv-Tech present |
| Currency economy | 🟡 Numismatics + Trading Floor present, Blockchain gated; vendors/coin-values pending |
| GregTech recipe conversion | ⏳ next phase (#17) — scope in §C |
| Create Aeronautics core | ✅ present (+ aeronautics-compatability) |
| Joke-mod audit | 🟡 this doc + #45 / #60 |

**Two drift risks the modlist reveals:** (1) the pack is **much broader than "Create economy"** —
heavy food / structure / magic / colony content; (2) several **parallel progression+economy
systems** (MineColonies, magic mods) that compete with "Create is how you make things." Neither is
fatal, but both need a deliberate keep/gate call (§D).

## A. Fit flags

**A1 — Pure memes (remove candidates; mostly already decided in #45):**
`mmmmmmmmmmmm`, `staaaaaaaaaaaack` (mega-stacks → **directly breaks the scarcity economy**),
`brazil-legends`, `create-more-pipe-bombs-in-packages` (grief), `umapyoi` (→ #60).
**Keep — harmless flavor:** `just-blahaj`, music discs (`game-discs`, `net-music`, `immersivemusicmod`,
`steves-lava-chicken…`), `emoji-type`.
**Rework not remove:** `the-bible` (recipe rework, #17).

**A2 — "Meme" but functional → keep:**
`touhou-little-maid` (maids = automated labor), `smokeleaf-industries` (a cash-crop *processing*
chain — actually **fits** the grow→process→sell economy; treat as an economy crop).

**A3 — Scope-creep content (beyond the stated TaCZ+Cannons / Create theme — keep or gate?):**
- Melee/adventure weapons: `epic-samurais`, `swashbucklers`, `valarian-conquest` ⚠ — weapon trees
  beyond the stated gun/cannon lineup.
- Magic stack: `ars-nouveau` + `ars-n-spells`, `irons-spells-n-spellbooks`, `occultism` +
  `occult-engineering`, `necronomicon`, `modular-golems` — full parallel power systems (big
  recipe-gate candidates, §C).
- `mffs` (force fields), `security-craft` ⚠ — base defense; minor fit for a PvE friend group.
- `sky-whale-ship` ⚠, `companion` / `companions-mod` / `buddy-drones` / `followers-teleport-too` ⚠ —
  companion/novelty; review value.

## B. Duplicate / overlapping **mechanics** (pick one)
> Duplicate **items** (lead×2, silver×2, tin, etc.) are largely auto-handled by **`almost-unified`**
> + **`polymorph`**. Below are overlapping **mods/mechanics**, which unifiers do *not* fix.

| Cluster | Mods | Recommendation |
|---|---|---|
| Inventory sort | `inventory-profiles-next` · `mouse-tweaks` · `invtweaks-emu-for-ipn` | IPN owns sort; Mouse Tweaks = drag-fill only (#38; note: `clientsort` isn't actually installed) |
| Accessory API | `curios` + `accessories` + `accessories-compat-layer` | dual slot APIs — keep both only if mods require each; confirm the compat-layer resolves ⚠ |
| Tree growth | `dynamictrees`(+`-terralith`,+`plus`) **vs** `oh-the-trees-youll-grow` | two growth systems → likely conflict; pick one ⚠ |
| Moon events | `enhanced-celestials` · `ender-moon` · `horde-moon` | overlapping; pick one (enhanced-celestials broadest) ⚠ |
| Food diversity | `spice-of-life-maids-dream` + `spice-of-life-onion` | two Spice-of-Life forks → redundant; keep one |
| Portable storage | `travelersbackpack` + sophisticated backpacks (`sophisticated-core`) | overlap (also a recipe-gate target, §C) |
| AFK / sleep | `afk-sleep` + `detect-afk` | overlap; one suffices ⚠ |

## C. Recipe-change targets (GregTech-conversion scope → #17)
Where "just craft it cheaply" undercuts the Create economy — gate behind Create processing, raise
costs, or nerf the loot source:
- **OP factory shortcuts:** `create-ultimate-factory`, `create-factory` ⚠ — kitchen-sink blocks that
  can trivialize automation; review their progression/recipes.
- **Magic items:** gate key ars/iron's-spells/occultism/golem items behind Create-made components
  (brass, precision mechanisms…) so magic plugs **into** the Create spine, not around it.
- **Easy storage/travel:** sophisticated backpacks, `travelersbackpack`, `create-storage` — push
  upgrades behind Create materials.
- **Auto-cooking:** `cooking-for-blockheads` (free food auto-craft) — gate or accept (food ≠ economy core).
- **Structure/dungeon loot** (Yung's ×5, `dungeons-and-taverns`, `when-dungeons-arise-seven-seas`,
  `towns-and-towers`, `endrem`, `valarian-conquest`) — hand out gear/valuables that bypass
  mine-and-sell. **Loot pass (#18):** strip free metal/diamond gear; optionally seed Numismatics
  coins as rare loot.

## D. The two big design calls (decide before the recipe overhaul)
1. **MineColonies vs "Create is the factory."** `minecolonies` (+ `structurize`, `domum-ornamentum`,
   `stylecolonies`, `ct-overhaul-village`, `multi-piston`, `towntalk`, `blockui`) is a whole **rival
   production+labor economy** — colonists auto-mine/farm/craft, which **short-circuits the
   scarcity-and-Create loop**. Decide: keep as a *village/social* layer with production
   **disabled/gated**, or cut. This is the single biggest fit question.
2. **Breadth vs focus.** Food (Farmer's Delight + ~10 Delight/"Let's Do" addons), decoration (5×
   Macaw's + Supplementaries + Handcrafted + …), structures (Yung's ×5 + more), mobs (~16, #42) —
   all large. Nothing breaks the theme, but this is a **kitchen-sink survival pack with a
   Create-economy spine**, not a focused Create pack. Fine if intended — just naming it so recipe
   scope is set realistically.

## E. Path forward (sequenced)
1. **Land current work** — playtest #58/#59 → merge #55/#56; merge #57 + run the digest.
2. **Curation pass** (#45 / #21 / #60) — remove agreed memes; resolve the §B duplicate-mechanic
   clusters (quick config/manifest edits).
3. **Decide the two big calls** (§D) — MineColonies role + breadth — *these set the recipe scope*.
4. **GregTech recipe overhaul** (#17) — metals/magic/storage through Create; needs the recipe digest.
5. **Loot pass** (#18) — structure/dungeon loot to fit scarcity; optionally seed coins.
6. **Economy finish** — Numismatics vendors/values + Trading Floor.
7. **Stability** (#42/#43/#48) and **QoL** (#38, …).

> Confidence: the Create / food / decoration / structure / colony / magic / library clusters are
> high-confidence. ⚠ items (and exact "what does X add") get confirmed by the ground-truth digest +
> an in-game EMI/JEI pass.
