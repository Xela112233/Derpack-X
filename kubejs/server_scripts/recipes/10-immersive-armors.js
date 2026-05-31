// Immersive Armors — "made through Create" (TEMPLATE: first worked branch).
//
// Principle: armor is plate -> METAL/MECHANICAL sets rebuilt from Create pressed SHEETS
// (+ a mechanism core on marquee sets). ORGANIC sets (wooden, bone, slime, prismarine, robe,
// wither) have no coherent Create routing -> LEFT UNTOUCHED (coherence over a forced grind).
//
// CONVENTION: each swap carries the ORIGINAL ingredients as an `orig:` comment (from the
// mod-data recipe digest) so review can weigh the author's intent and catch accidental
// easier/harder shifts. (Digest gives ingredients, not exact shape/counts.)

ServerEvents.recipes(event => {
  const swap = (id, pattern, key) => { event.remove({ output: id }); event.shaped(id, pattern, key) }

  // ── HEAVY — iron plate. NOTE: orig used iron_block (dense), so sheets may be a touch cheaper. ──
  // orig heavy_helmet:     iron_block + iron_ingot
  swap('immersive_armors:heavy_helmet',     ['SSS', 'S S'],        { S: 'create:iron_sheet' })
  // orig heavy_chestplate: iron_block + iron_ingot
  swap('immersive_armors:heavy_chestplate', ['S S', 'SSS', 'SSS'], { S: 'create:iron_sheet' })
  // orig heavy_leggings:   iron_block + iron_ingot
  swap('immersive_armors:heavy_leggings',   ['SSS', 'S S', 'S S'], { S: 'create:iron_sheet' })
  // orig heavy_boots:      iron_block + iron_ingot
  swap('immersive_armors:heavy_boots',      ['S S', 'S S'],        { S: 'create:iron_sheet' })

  // ── STEAMPUNK — brass + gears + mechanism core (orig was gold + piston/clock/compass gadgetry) ──
  // orig steampunk_helmet:     gold_ingot + redstone_torch + white_stained_glass
  swap('immersive_armors:steampunk_helmet',     ['BBB', 'B B'],        { B: 'create:brass_sheet' })
  // orig steampunk_chestplate: gold_ingot + piston + clock + compass
  swap('immersive_armors:steampunk_chestplate', ['B B', 'BPB', 'BBB'], { B: 'create:brass_sheet', P: 'create:precision_mechanism' })
  // orig steampunk_leggings:   gold_ingot + piston
  swap('immersive_armors:steampunk_leggings',   ['BBB', 'BGB', 'B B'], { B: 'create:brass_sheet', G: 'create:cogwheel' })
  // orig steampunk_boots:      gold_ingot + piston
  swap('immersive_armors:steampunk_boots',      ['B B', 'BGB'],        { B: 'create:brass_sheet', G: 'create:cogwheel' })

  // ── WARRIOR — leather reinforced with iron plate (orig was iron_ingot + the leather piece) ──
  // orig warrior_helmet:     iron_ingot + leather_helmet + bone
  swap('immersive_armors:warrior_helmet',     ['SSS', 'SLS'],        { S: 'create:iron_sheet', L: 'minecraft:leather_helmet' })
  // orig warrior_chestplate: iron_ingot + leather + leather_chestplate
  swap('immersive_armors:warrior_chestplate', ['S S', 'SLS', 'SSS'], { S: 'create:iron_sheet', L: 'minecraft:leather_chestplate' })
  // orig warrior_leggings:   iron_ingot + leather_leggings
  swap('immersive_armors:warrior_leggings',   ['SSS', 'SLS', 'S S'], { S: 'create:iron_sheet', L: 'minecraft:leather_leggings' })
  // orig warrior_boots:      iron_ingot + leather_boots
  swap('immersive_armors:warrior_boots',      ['S S', 'SLS'],        { S: 'create:iron_sheet', L: 'minecraft:leather_boots' })

  // ── DIVINE — endgame upgrade of robe. KEEP the original gates (emerald on helmet, TOTEM on
  //    chest) so it isn't made easier; add golden plate + a mechanism core for the Create routing. ──
  // orig divine_helmet:     robe_helmet + gold_ingot + emerald            (emerald gate KEPT)
  swap('immersive_armors:divine_helmet',     ['SES', 'SRS'],        { S: 'create:golden_sheet', E: 'minecraft:emerald', R: 'immersive_armors:robe_helmet' })
  // orig divine_chestplate: robe_chestplate + gold_ingot + totem_of_undying   (TOTEM gate KEPT)
  swap('immersive_armors:divine_chestplate', ['STS', 'SPS', 'SRS'], { S: 'create:golden_sheet', T: 'minecraft:totem_of_undying', P: 'create:precision_mechanism', R: 'immersive_armors:robe_chestplate' })
  // orig divine_leggings:   robe_leggings + gold_ingot
  swap('immersive_armors:divine_leggings',   ['SSS', 'SRS', 'S S'], { S: 'create:golden_sheet', R: 'immersive_armors:robe_leggings' })
  // orig divine_boots:      robe_boots + gold_ingot
  swap('immersive_armors:divine_boots',      ['S S', 'SRS'],        { S: 'create:golden_sheet', R: 'immersive_armors:robe_boots' })

  // ── LEFT AS-IS (no coherent Create routing; converting = forced grind):
  //    wooden, bone, slime, prismarine, robe, wither. Robe stays the soft entry divine builds on.
})
