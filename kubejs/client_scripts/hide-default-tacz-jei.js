// Hide TaCZ's default-pack guns/ammo/attachments from JEI's ingredient list (#307).
//
// Companion to startup_scripts/hide-default-tacz-creative.js — same predicate, same
// rationale (the default pack must stay *loaded* for tacz:12g, but its uncraftable
// variants shouldn't be browsable). JEI builds its item list from the creative search
// entries, so the creative-side removal already covers most of this; this hide is the
// belt-and-suspenders that keeps JEI clean even if a tab's search entries slip through
// (and it's the lever the issue names). Startup and client scripts run in separate
// script contexts, so the predicate is duplicated, not shared.
//
// API (verified against kubejs-neoforge-2101.7.2-build.368 bytecode):
// RecipeViewerEvents.removeEntries('item', ...) is posted by KubeJS's JEI plugin when
// the JEI runtime comes up; event.remove() accepts a JS function wrapped into an
// ItemPredicate and tests every listed ItemStack (per-GunId variants are distinct
// entries thanks to TaCZ's subtype interpreters).
//
// VERIFY IN PLAYTEST: JEI's item list shows only create_armorer guns/ammo/attachments
// (+ tacz:12g); no bare AK/Glock/etc. entries; the four TaCZ workbenches are gone;
// armorer items' Create recipes still look up fine from JEI.

const $IGun = Java.loadClass('com.tacz.guns.api.item.IGun')
const $IAmmo = Java.loadClass('com.tacz.guns.api.item.IAmmo')
const $IAttachment = Java.loadClass('com.tacz.guns.api.item.IAttachment')

function isDefaultTaczVariant(stack) {
  try {
    const iGun = $IGun.getIGunOrNull(stack)
    if (iGun) {
      const id = iGun.getGunId(stack)
      return !id || id.getNamespace() === 'tacz'
    }
    const iAmmo = $IAmmo.getIAmmoOrNull(stack)
    if (iAmmo) {
      const id = iAmmo.getAmmoId(stack)
      return !id || (id.getNamespace() === 'tacz' && id.toString() !== 'tacz:12g')
    }
    const iAttachment = $IAttachment.getIAttachmentOrNull(stack)
    if (iAttachment) {
      const id = iAttachment.getAttachmentId(stack)
      return !id || id.getNamespace() === 'tacz'
    }
  } catch (err) {
    return false // unreadable stack -> leave it alone
  }
  return false
}

const DEAD_TABLES = ['tacz:gun_smith_table', 'tacz:workbench_a', 'tacz:workbench_b', 'tacz:workbench_c']

RecipeViewerEvents.removeEntries('item', event => {
  event.remove(stack => isDefaultTaczVariant(stack) || DEAD_TABLES.includes(stack.kjs$getId()))
})
