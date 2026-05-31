// Samurai Dynasty — "made through Create" (restraint case).
//
// Samurai already runs on STEEL (c:ingots/steel), a Create-made metal in this pack (TFMG /
// Ironworks / Metalwork). So it's largely Create-routed already — converting every steel recipe
// would be the over-gating we avoid. ONE visible Create touch: the flagship "steel samurai" plate
// set is built from Create-pressed steel PLATES (#c:plates/steel -> create_ironworks:steel_sheet
// via create:pressing), keeping the jade ornament on the helmet. Weapons, the light/ninja/basic
// sets, stone, cloth and netherite upgrades are coherent as-is and left untouched.
//
// STEEL SOURCING (samurai's iron->steel blasting, mffs's smelt) is a pack-wide question -> deferred
// to a dedicated metals/steel pass, not touched here.
//
// Verified against tools/mod-data/recipes/samurai_dynasty-*.txt.

ServerEvents.recipes(event => {
  const swap = (id, pattern, key) => { event.remove({ output: id }); event.shaped(id, pattern, key) }
  const P = '#c:plates/steel' // Create-pressed steel plate

  // orig steel_samurai_helmet:     c:gems/jade + c:ingots/steel + iron_ingot   (jade ornament KEPT)
  swap('samurai_dynasty:steel_samurai_helmet',     ['SJS', 'S S'],        { S: P, J: '#c:gems/jade' })
  // orig steel_samurai_chestplate: c:ingots/steel + iron_ingot
  swap('samurai_dynasty:steel_samurai_chestplate', ['S S', 'SSS', 'SSS'], { S: P })
  // orig steel_samurai_leggings:   c:ingots/steel + iron_ingot
  swap('samurai_dynasty:steel_samurai_leggings',   ['SSS', 'S S', 'S S'], { S: P })
  // orig steel_samurai_boots:      c:ingots/steel + iron_ingot
  swap('samurai_dynasty:steel_samurai_boots',      ['S S', 'S S'],        { S: P })

  // LEFT AS-IS (already steel-routed / coherent): all weapons (katana, kunai, naginata, …),
  // light_samurai_* + steel_* + steel_ninja_* sets, jade/spirit-stone blocks, cloth/kimono,
  // and the netherite smithing upgrades.
})
