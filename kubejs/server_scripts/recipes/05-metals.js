// Metals — source the alloy metals through Create (FOUNDATION pass; runs before the per-mod files).
//
// Bronze & brass are already Create-only here (Ironworks / Metalwork create:mixing of copper+tin /
// copper+zinc). STEEL had clean Create paths (TFMG casting, Metalwork molten->compacting, Create Big
// Cannons compacting, Ironworks) AND two trivial NON-Create shortcuts. Closing those routes ALL steel
// through Create — and because steel is shared via #c:ingots/steel (almost-unified), this
// auto-propagates to every steel-using mod (Samurai, MFFS, …) with no per-mod edits.
//
// Not a grind: steel just becomes a mid-tier "needs a bit of Create running" metal (one heating/
// casting step), then it's used freely. Vanilla iron/copper/gold smelting is untouched.
//
// Verified against tools/mod-data/recipes/ — the Create steel sources below remain after removal.

ServerEvents.recipes(event => {
  // Remove non-Create steel shortcuts BY RECIPE ID (never by output — that would nuke the Create paths).
  // orig: samurai_dynasty:steel_ingot_from_blasting_iron_ingot = blast(minecraft:iron_ingot) -> steel
  event.remove({ id: 'samurai_dynasty:steel_ingot_from_blasting_iron_ingot' })
  // orig: mffs:steel_compound = shaped(c:ingots/iron + coal) -> steel_compound
  // orig: mffs:steel_ingot    = smelt(mffs:steel_compound)   -> steel
  event.remove({ id: 'mffs:steel_compound' })
  event.remove({ id: 'mffs:steel_ingot' })

  // KEPT — the Create steel chain (steel stays craftable through these):
  //   tfmg:casting/steel · createmetalwork …compacting(molten_steel) · createbigcannons:compacting/forge_steel_ingot · create_ironworks
  // ALREADY Create-only: bronze & brass (Ironworks/Metalwork create:mixing).
  // FLAG (left for review, lower priority): occultengineering smelt/blast of brass_dust -> create:brass_ingot
  //   is a non-Create brass alt-path — revisit if brass should be Create-only too.
})
