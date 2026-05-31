// Immersive Armors — "made through Create" (TEMPLATE: first worked branch).
//
// Principle in action: armor is plate, so the METAL/MECHANICAL sets are rebuilt from Create
// pressed SHEETS (+ a mechanism core on the marquee sets). The ORGANIC/thematic sets
// (wooden, bone, slime, prismarine, robe, wither) have no coherent Create routing, so they are
// LEFT UNTOUCHED on purpose — coherence over a forced grind. Costs stay ~comparable to vanilla;
// we're changing the *flavor* of the craft, not adding tedium.
//
// Recipes verified against tools/mod-data/recipes/immersive_armors-*.txt.

ServerEvents.recipes(event => {
  const swap = (id, pattern, key) => { event.remove({ output: id }); event.shaped(id, pattern, key) }

  // ── HEAVY — iron plate armor: pressed iron sheets (was iron blocks/ingots) ──
  swap('immersive_armors:heavy_helmet',     ['SSS', 'S S'],        { S: 'create:iron_sheet' })
  swap('immersive_armors:heavy_chestplate', ['S S', 'SSS', 'SSS'], { S: 'create:iron_sheet' })
  swap('immersive_armors:heavy_leggings',   ['SSS', 'S S', 'S S'], { S: 'create:iron_sheet' })
  swap('immersive_armors:heavy_boots',      ['S S', 'S S'],        { S: 'create:iron_sheet' })

  // ── STEAMPUNK — the natural Create fit: brass plate + gears + a precision-mechanism core ──
  swap('immersive_armors:steampunk_helmet',     ['BBB', 'B B'],        { B: 'create:brass_sheet' })
  swap('immersive_armors:steampunk_chestplate', ['B B', 'BPB', 'BBB'], { B: 'create:brass_sheet', P: 'create:precision_mechanism' })
  swap('immersive_armors:steampunk_leggings',   ['BBB', 'BGB', 'B B'], { B: 'create:brass_sheet', G: 'create:cogwheel' })
  swap('immersive_armors:steampunk_boots',      ['B B', 'BGB'],        { B: 'create:brass_sheet', G: 'create:cogwheel' })

  // ── WARRIOR — leather reinforced with iron plating (light touch) ──
  swap('immersive_armors:warrior_helmet',     ['SSS', 'SLS'],        { S: 'create:iron_sheet', L: 'minecraft:leather_helmet' })
  swap('immersive_armors:warrior_chestplate', ['S S', 'SLS', 'SSS'], { S: 'create:iron_sheet', L: 'minecraft:leather_chestplate' })
  swap('immersive_armors:warrior_leggings',   ['SSS', 'SLS', 'S S'], { S: 'create:iron_sheet', L: 'minecraft:leather_leggings' })
  swap('immersive_armors:warrior_boots',      ['S S', 'SLS'],        { S: 'create:iron_sheet', L: 'minecraft:leather_boots' })

  // ── DIVINE — endgame: gilded robe + a mechanism core (keeps the robe base it upgrades from) ──
  swap('immersive_armors:divine_helmet',     ['SSS', 'SRS'],        { S: 'create:golden_sheet', R: 'immersive_armors:robe_helmet' })
  swap('immersive_armors:divine_chestplate', ['S S', 'SPS', 'SRS'], { S: 'create:golden_sheet', P: 'create:precision_mechanism', R: 'immersive_armors:robe_chestplate' })
  swap('immersive_armors:divine_leggings',   ['SSS', 'SRS', 'S S'], { S: 'create:golden_sheet', R: 'immersive_armors:robe_leggings' })
  swap('immersive_armors:divine_boots',      ['S S', 'SRS'],        { S: 'create:golden_sheet', R: 'immersive_armors:robe_boots' })

  // ── LEFT AS-IS (no coherent Create routing — converting would be a forced grind):
  //    wooden, bone, slime, prismarine, robe, wither. Robe stays the soft entry that divine builds on.
})
