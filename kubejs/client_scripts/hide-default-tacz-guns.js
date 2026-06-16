// Hide the default (non-armorer) TaCZ guns / ammo / attachments from the recipe viewer
// and creative search (#307).
//
// The pack's gun roster is the Create: Armorer set, crafted through Create (createimmersivetacz
// recipes), and the stock TaCZ content is already UNCRAFTABLE — `remove-default-tacz-guns.js`
// strips the whole `tacz:gun_smith_table_crafting` type and the bench blocks. But the default
// guns/ammo/attachments still SHOW in JEI and creative, which reads as broken content ("there's
// an AK in JEI with no recipe"). This hides them.
//
// Why it can't be a plain id blacklist (config/jei/blacklist.cfg): every gun is the SAME item id,
// `tacz:modern_kinetic_gun`, differing only by a `minecraft:custom_data` GunId — the pack's
// guns are `GunId: create_armorer:*` riding on that base item (likewise `tacz:attachment` +
// AttachmentId). TaCZ registers a JEI subtype per GunId/AttachmentId, so the viewer lists each
// variant separately, but they share one item id — an id-level hide would take the armorer guns
// with them. So we hide by VARIANT, keeping only the `create_armorer:*` ones.
//
// API (KubeJS 1.21.1): `JEIEvents.hideItems` is superseded by `RecipeViewerEvents.removeEntries`
// (one call covering JEI/REI/EMI). We pass a custom Ingredient that tests each viewer stack:
//   - `tacz:ammo`                : ALL variants are default-pack — the pack's real ammo is its own
//                                  `createimmersivetacz:*` items (different ids), so hide every one.
//   - `tacz:modern_kinetic_gun`  : hide unless the variant's component data names `create_armorer`.
//   - `tacz:attachment`          : same create_armorer keep-test.
// Substring test (not an exact component match) on purpose: the keeper stacks may carry extra
// custom_data keys (e.g. GunFireMode) that an `Ingredient.of(id).except(id[custom_data={…}])`
// exact-match would miss and wrongly hide. It also means we never have to enumerate TaCZ's
// bundled default-pack GunId list — anything that isn't ours is hidden by exclusion.
//
// The shared carriers themselves (`tacz:modern_kinetic_gun` / `tacz:ammo` with no/empty variant)
// stay craftable and functional — this is presentation only, client-side, no recipe/gameplay change.
//
// VERIFY IN PLAYTEST (static CI can't check this — the removeEntries iteration over TaCZ's
// per-GunId subtypes and the component-string read are runtime behaviour):
//   - JEI/creative show ONLY create_armorer:* guns (14: rifle_assult_crane, shotgun_db_stone,
//     cannon_40mm_salamander, gl_revolver_devastator, mg_platemag_flywheel, special_melee_wrench,
//     pistol_auto_stress, shotgun_pump_bearing, pistol_revolver_torque, rifle_assult_roller,
//     sniper_semi_m1, smg_auto_crank, sniper_semi_clockwork, special_melee_atomic) — no stock AK/M4/etc.
//   - Only create_armorer:* attachments remain; no default tacz attachments.
//   - No `tacz:ammo` entries (the createimmersivetacz:* shells/casings stay).
//   - Armorer guns still craft (Mechanical Crafter) and fire; client log clean on load + hover.
//   - If the keepers ALSO vanished: the component-string read found no 'create_armorer' marker —
//     log one keeper's `stack.componentString` and widen `variantData()` to the field it uses.
//   - If nothing hid: removeEntries isn't iterating subtypes here — fall back to the per-GunId
//     `event.remove('tacz:modern_kinetic_gun[minecraft:custom_data={GunId:"<id>"}]')` form with the
//     in-game-enumerated default GunId list, or TaCZ's own default-pack toggle if one exists.

const CARRIERS = ['tacz:modern_kinetic_gun', 'tacz:attachment'] // variant-gated: keep create_armorer
const KEEP_MARKER = 'create_armorer'

// Best-effort read of a stack's component/NBT text, across the shapes KubeJS may expose it as.
function variantData(stack) {
  let s = ''
  try { if (stack.componentString) s += '' + stack.componentString } catch (err) { /* ignore */ }
  try { if (stack.nbt) s += ' ' + stack.nbt } catch (err) { /* ignore */ }
  return s
}

RecipeViewerEvents.removeEntries('item', event => {
  event.remove(Ingredient.custom(stack => {
    try {
      const id = stack.id
      if (id === 'tacz:ammo') return true               // all default ammo (real ammo is createimmersivetacz:*)
      if (CARRIERS.indexOf(id) >= 0) {
        return variantData(stack).indexOf(KEEP_MARKER) < 0 // hide every non-armorer variant
      }
      return false
    } catch (err) {
      return false // never hide on an unexpected read error
    }
  }))
})
