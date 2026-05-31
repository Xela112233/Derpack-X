// Derpack X — "Made through Create" recipe overhaul.
//
// Principle: route crafting through Create PARTS + METHODS, woven through every branch.
// This is NOT a grind — coherence over tedium. Full convention + parts palette: docs/RECIPES.md
//
// Per-branch conversions live in sibling files (10-metals.js, 20-tools.js, 30-magic.js, ...),
// added as we read each mod's recipes (from jars on hand, or the tools/mod-data/recipes/ digest).
// This file is shared notes + a load check. It changes no recipes on its own.

ServerEvents.recipes(event => {
  console.info('[derpack] recipe-overhaul scripts loaded')

  // ─────────────────────────────────────────────────────────────────────────────
  // PATTERN TEMPLATES — copy into a per-branch file, then uncomment + fill real ids.
  // ─────────────────────────────────────────────────────────────────────────────
  //
  // (A) Gate a craftable behind a Create part (one gate per tier — don't stack):
  //   event.remove({ output: 'somemod:gizmo' })
  //   event.shaped('somemod:gizmo', ['SBS', 'BPB', 'SBS'], {
  //     S: 'create:brass_sheet', B: '#c:ingots/brass', P: 'create:precision_mechanism'
  //   })
  //
  // (B) Build it via a Create METHOD instead of the crafting table:
  //   event.remove({ output: 'somemod:metal_plate' })
  //   event.recipes.create.pressing('somemod:metal_plate', '#c:ingots/somemetal')
  //   event.recipes.create.mixing('2x somemod:bronze_ingot', ['#c:ingots/copper', '#c:ingots/tin'])
  //
  // (C) Ore → metal: add a Create crushing REWARD; leave the vanilla 1:1 smelt intact (not a tax):
  //   event.recipes.create.crushing(
  //     ['2x somemod:crushed_raw_metal', Item.of('create:experience_nugget').withChance(0.2)],
  //     'somemod:metal_ore'
  //   )
  //
  // Reminder: target #c: tags where possible (almost-unified collapses dupes); always pair a
  // remove() with a replacement so nothing becomes uncraftable.
})
