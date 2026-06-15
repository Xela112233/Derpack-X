// Mining costs hunger — a Survival-pillar pressure mechanic.
//
// Breaking a solid block makes the player hungrier. The cost scales with a per-block
// "weight" so clearing dirt is cheap, hauling out a vein of ore or chopping into solid
// metal is dear, and near-massless things (torches, crops, flowers, redstone wire) are
// free. This raises baseline food demand, which is exactly the pressure that pulls on the
// diet/food economy (see docs/SYSTEMS.md, the Survival pillar). Design + tuning notes:
// docs/MINING-HUNGER.md.
//
// VERIFY IN PLAYTEST (this can't be checked statically — see the PR's Playtest block):
//   - block-break actually drains hunger at roughly the rates below, and
//   - it does NOT fire for Create drills / TNT / other no-player breaks, nor in creative.

// ---- Tunables --------------------------------------------------------------

// Exhaustion added per "weight 1.0" block. Vanilla drains 1 food point per 4.0 exhaustion,
// and a full hunger bar is 20 points (= 80 exhaustion), so 0.8 here means exactly one full
// hunger bar per 100 standard blocks — the design target. (Saturation absorbs the first
// hits, as it does for any exhaustion, so the bar starts visibly dropping once it's gone.)
const EXHAUSTION_PER_WEIGHT = 0.8

// Tool-tier multiplier on the drain — better tools tire you less. Values per zagwar: a
// bare hand is hardest (2.5×), each tier of tool easier, down to netherite (0.5×). This is
// why the "100 blocks per bar" anchor below is calibrated to an IRON tool (1.0×): wood/none
// compress it, diamond/netherite stretch it. Final drain = rate × block weight × this.
// ⚠️ The tooltip "Stamina drain" line (client_scripts/mining-speed-tooltip.js) mirrors
// these values and toolMultiplier()'s id rules — KEEP THEM IN SYNC when tuning.
const T_NONE      = 2.5  //  bare hand / non-mining item (sword, block, food, …)
const T_WOOD      = 2.0  //  wooden (and golden — same harvest tier as wood)
const T_STONE     = 1.5
const T_IRON      = 1.0  //  the anchor: iron tool on a Normal block = 100/bar
const T_DIAMOND   = 0.75
const T_NETHERITE = 0.5
const T_UNKNOWN   = 1.0  //  a tool we recognise as a tool but whose tier we can't read (modded) — neutral

// Item-id suffixes that mark something as a mining tool. Swords/weapons are deliberately
// excluded — holding one counts as "no tool" for mining. "_paxel" covers the common modded multitool.
const TOOL_SUFFIXES = ['_pickaxe', '_axe', '_shovel', '_hoe', '_paxel']

// Weight tiers (multipliers on the rate above). These deliberately MIRROR the Create:
// Aeronautics physics-value mass ranking — light envelopes/wood/wool sit well below
// stone/dirt, metal and ore sit well above (oak planks are ~0.5 kpg in Aeronautics) — so
// "heavier to fly is heavier to mine." They are NOT read from Aeronautics at runtime (it
// exposes no stable API for resolved block mass); they track its mass tiers by category.
// Tune freely — this list is the balancing surface.
const W_FREE   = 0.0  //   no cost — instant-break / near-massless: torches, flowers, redstone wire, leaves
const W_CROP   = 0.25 //  ~400/bar (iron hoe) — harvesting crops; a light farming cost (per zagwar)
const W_LIGHT  = 0.5  //  200/bar — wood, wool, glass, sand, ice, foliage solids
const W_NORMAL = 1.0  //  100/bar — stone, dirt, concrete, bricks (the anchor)
const W_HEAVY  = 1.6  //  ~62/bar — ores, obsidian, deepslate, dense machines
const W_VHEAVY = 2.2  //  ~45/bar — solid metal / storage blocks, netherite, ancient debris

// ---- Classification --------------------------------------------------------

// Matched against the block id's path (the part after "namespace:"). First tier that
// matches wins, in the order checked below. Substring tokens are chosen to be collision-
// free against the exact ids in the lighter tiers; ambiguous cases use exact matches.

// Metal/dense names — any "<name>_block" is treated as solid metal storage (very heavy).
// Catches modded metals generically (e.g. tin_block, steel_block) without listing each.
const METALS = [
  'iron', 'gold', 'copper', 'netherite', 'diamond', 'emerald', 'steel', 'bronze', 'tin',
  'lead', 'silver', 'zinc', 'brass', 'nickel', 'platinum', 'tungsten', 'titanium',
  'aluminum', 'aluminium', 'osmium', 'uranium', 'cobalt', 'electrum', 'invar',
  'constantan', 'palladium'
]
const VHEAVY_EXACT = ['ancient_debris']

// Heavy: ores and dense stone/machines. Ores are matched on a word boundary (below) so
// "*_ore" lands here while incidental "ore" substrings (e.g. "reactor_core") do not.
const HEAVY_CONTAINS = ['obsidian', 'deepslate', 'blackstone', 'basalt']
const isOre = (p) => p === 'ore' || p.endsWith('_ore') || p.indexOf('_ore_') !== -1
const HEAVY_EXACT = ['anvil', 'chipped_anvil', 'damaged_anvil', 'lodestone', 'respawn_anchor', 'enchanting_table']

// Free (no hunger): near-massless blocks. Suffix/exact lists are kept separate from the
// LIGHT substrings so e.g. exact "snow" (the layer) is free while "snow_block" stays light.
const FREE_SUFFIX = ['_torch', '_sapling', '_propagule', '_button', '_pressure_plate', '_carpet', '_banner', '_sign', '_candle', '_tulip', '_head', '_skull']
const FREE_CONTAINS = ['leaves', '_fungus']
const FREE_EXACT = [
  'torch', 'candle', 'redstone_wire', 'repeater', 'comparator', 'lever', 'tripwire', 'tripwire_hook',
  'string', 'cobweb', 'rail', 'powered_rail', 'detector_rail', 'activator_rail', 'flower_pot', 'end_rod',
  'lightning_rod', 'chain', 'scaffolding', 'ladder', 'lantern', 'soul_lantern', 'snow', 'powder_snow',
  'sugar_cane', 'cactus', 'bamboo', 'vine', 'weeping_vines', 'twisting_vines', 'kelp', 'kelp_plant',
  'seagrass', 'tall_seagrass', 'sea_pickle', 'lily_pad', 'dead_bush', 'short_grass', 'tall_grass',
  'fern', 'large_fern', 'melon_stem',
  'pumpkin_stem', 'attached_melon_stem', 'attached_pumpkin_stem', 'brown_mushroom', 'red_mushroom',
  'crimson_roots', 'warped_roots', 'nether_sprouts', 'glow_lichen', 'sculk_vein', 'hanging_roots',
  'spore_blossom', 'small_dripleaf', 'big_dripleaf', 'azalea', 'flowering_azalea', 'pink_petals',
  'torchflower', 'frogspawn', 'chorus_flower', 'chorus_plant', 'poppy', 'dandelion', 'blue_orchid',
  'allium', 'azure_bluet', 'oxeye_daisy', 'cornflower', 'lily_of_the_valley', 'wither_rose', 'sunflower',
  'lilac', 'rose_bush', 'peony', 'wildflowers', 'closed_eyeblossom', 'open_eyeblossom', 'leaf_litter',
  'firefly_bush'
]

// Crops: harvestable farm plants. A light cost so big manual harvests aren't free, while
// routine farming stays cheap. "_crop" catches modded crops (e.g. farmersdelight tomato_crop).
// Maturity isn't checked — breaking any crop block costs the same light amount.
const CROP_EXACT = ['wheat', 'carrots', 'potatoes', 'beetroots', 'nether_wart', 'cocoa', 'sweet_berry_bush', 'cave_vines', 'cave_vines_plant', 'torchflower_crop', 'pitcher_crop', 'pitcher_plant']
const CROP_CONTAINS = ['_crop']

// Light: planks, wool, glass, loose ground, foliage solids. Raw timber (logs, wood, stems,
// hyphae) is NOT light — it falls through to Normal; only processed planks sit here.
const LIGHT_CONTAINS = ['_planks', 'wool', 'glass', 'ice']
const LIGHT_EXACT = [
  'sand', 'red_sand', 'soul_sand', 'soul_soil', 'gravel', 'suspicious_sand', 'suspicious_gravel',
  'mud', 'packed_mud', 'hay_block', 'sponge', 'wet_sponge', 'glowstone', 'shroomlight',
  'ochre_froglight', 'verdant_froglight', 'pearlescent_froglight', 'slime_block', 'honey_block',
  'honeycomb_block', 'dried_kelp_block', 'snow_block', 'bamboo_block', 'stripped_bamboo_block',
  'bamboo_mosaic', 'nether_wart_block', 'warped_wart_block', 'brown_mushroom_block', 'red_mushroom_block',
  'mushroom_stem', 'moss_block', 'pale_moss_block', 'dripstone_block', 'pointed_dripstone', 'bone_block',
  'melon', 'pumpkin', 'carved_pumpkin', 'jack_o_lantern', 'sculk', 'sculk_catalyst', 'sculk_sensor',
  'sculk_shrieker'
]

const includesAny = (s, arr) => arr.some(x => s.includes(x))
const equalsAny = (s, arr) => arr.indexOf(s) !== -1
const endsWithAny = (s, arr) => arr.some(x => s.endsWith(x))

function hungerWeight(id) {
  const p = id.indexOf(':') === -1 ? id : id.substring(id.indexOf(':') + 1)

  // Very heavy: solid metal / dense storage blocks.
  if (equalsAny(p, VHEAVY_EXACT)) return W_VHEAVY
  if (p.endsWith('_block') && includesAny(p, METALS)) return W_VHEAVY

  // Heavy: ores and dense stone/machines.
  if (isOre(p) || includesAny(p, HEAVY_CONTAINS) || equalsAny(p, HEAVY_EXACT)) return W_HEAVY

  // Crops: light farming cost (checked before Free, since these would otherwise read as free).
  if (equalsAny(p, CROP_EXACT) || includesAny(p, CROP_CONTAINS)) return W_CROP

  // Free: near-massless foliage / decoration / mechanisms.
  if (equalsAny(p, FREE_EXACT) || endsWithAny(p, FREE_SUFFIX) || includesAny(p, FREE_CONTAINS)) return W_FREE

  // Light: wood, wool, glass, loose ground.
  if (includesAny(p, LIGHT_CONTAINS) || equalsAny(p, LIGHT_EXACT)) return W_LIGHT

  // Everything else (stone, dirt, concrete, bricks, terracotta, …).
  return W_NORMAL
}

// Multiplier from the tool held in the main hand. Recognised vanilla tiers map to zagwar's
// values; a recognised-but-unknown-tier tool (modded material) is neutral; anything that
// isn't a mining tool (empty hand, sword, a block, food) is the hardest "no tool" case.
function toolMultiplier(itemId) {
  const p = itemId.indexOf(':') === -1 ? itemId : itemId.substring(itemId.indexOf(':') + 1)
  if (!p || p === 'air') return T_NONE
  if (!endsWithAny(p, TOOL_SUFFIXES)) return T_NONE
  if (p.indexOf('netherite') !== -1) return T_NETHERITE
  if (p.indexOf('diamond') !== -1) return T_DIAMOND
  if (p.indexOf('iron') !== -1) return T_IRON
  if (p.indexOf('stone') !== -1) return T_STONE
  if (p.indexOf('wood') !== -1 || p.indexOf('gold') !== -1) return T_WOOD
  return T_UNKNOWN
}

// ---- Event -----------------------------------------------------------------

BlockEvents.broken(event => {
  const player = event.player
  // No player → Create drills, TNT, fire, pistons, etc. These pay no hunger.
  if (!player) return

  try {
    if (player.isCreative() || player.isSpectator()) return
  } catch (e) {
    // If those accessors aren't present on this build, fall through —
    // causeFoodExhaustion is already a no-op for creative players in vanilla.
  }

  let weight
  try {
    weight = hungerWeight(String(event.block.id))
  } catch (e) {
    weight = W_NORMAL // classification must never break block-breaking
  }
  if (weight <= 0) return

  let toolMult = T_IRON // neutral default if the held item can't be read
  try {
    toolMult = toolMultiplier(String(event.player.mainHandItem.id))
  } catch (e) { /* keep neutral default */ }

  const amount = EXHAUSTION_PER_WEIGHT * weight * toolMult
  try {
    player.causeFoodExhaustion(amount)
  } catch (e) {
    try { player.getFoodData().addExhaustion(amount) } catch (e2) { /* give up silently */ }
  }
})
