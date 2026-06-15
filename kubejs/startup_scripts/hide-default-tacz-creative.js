// Hide TaCZ's default-pack guns/ammo/attachments from the creative menu (#307).
//
// The pack's gun roster is the Create: Armorer set (`create_armorer:*` gun ids riding on
// the generic tacz items, Create-crafted via createimmersivetacz). The stock TaCZ content
// is already uncraftable — remove-default-tacz-guns.js strips the whole
// gun_smith_table_crafting type — but the default pack's variants still showed in TaCZ's
// creative tabs, which reads as broken content ("there's an AK with no recipe").
//
// Why hide instead of disabling the default gun pack: the default pack can't be dropped
// wholesale — two armorer shotguns load `tacz:12g` (their gun data's "ammo" field), and
// the pack's Create path crafts it (createimmersivetacz twelve_gauge_shell_fill), so the
// default pack's ammo indexes must stay loaded. Presentation-only cleanup instead.
//
// Mechanics (verified against tacz-neoforge-1.21.1-1.1.8-hotfix-r1 +
// kubejs-neoforge-2101.7.2-build.368 bytecode):
//   - TaCZ registers exactly 15 creative tabs: 7 gun types + 6 attachment types + ammo +
//     other (names extracted from ModCreativeTabs's DeferredRegister calls, listed below).
//   - Gun/ammo/attachment variants are ItemStacks of the generic tacz items carrying
//     their pack id in custom_data; TaCZ's own API statics (IGun.getIGunOrNull etc.)
//     read it back, so the predicate filters by id namespace — `create_armorer:*`
//     variants survive untouched.
//   - StartupEvents.modifyCreativeTab posts per BuildCreativeModeTabContentsEvent and
//     its remove() walks both parent and search entries, so the search tab (which JEI
//     also feeds from) is cleaned too.
//   - The four gun-table workbenches are dead ends (their recipes are removed and every
//     bench tab is empty), so they're hidden from the `other` tab as well.
//
// Exception kept visible: `tacz:12g` (see above — real, craftable pack content).
//
// The JEI side of #307 lives in client_scripts/hide-default-tacz-jei.js (same predicate;
// startup and client scripts can't share code).
//
// VERIFY IN PLAYTEST: TaCZ creative tabs show only create_armorer variants (+ tacz:12g in
// the ammo tab); the four workbenches are gone from the `other` tab; targets/statue/ammo
// boxes still show; armorer guns still craft and fire.

const $IGun = Java.loadClass('com.tacz.guns.api.item.IGun')
const $IAmmo = Java.loadClass('com.tacz.guns.api.item.IAmmo')
const $IAttachment = Java.loadClass('com.tacz.guns.api.item.IAttachment')

// Default-pack variant = a TaCZ gun/ammo/attachment stack whose pack id sits in the
// `tacz` namespace (the default pack), bar the one armorer-shared ammo. An id-less
// stack of those items is a broken placeholder — hide it too.
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

// The dead gun-table workbenches (recipes removed in remove-default-tacz-guns.js).
const DEAD_TABLES = ['tacz:gun_smith_table', 'tacz:workbench_a', 'tacz:workbench_b', 'tacz:workbench_c']

const TACZ_TABS = [
  'tacz:pistol', 'tacz:sniper', 'tacz:rifle', 'tacz:shotgun', 'tacz:smg', 'tacz:rpg', 'tacz:mg',
  'tacz:scope', 'tacz:muzzle', 'tacz:stock', 'tacz:grip', 'tacz:extended_mag', 'tacz:laser',
  'tacz:ammo', 'tacz:other'
]

TACZ_TABS.forEach(tab => {
  StartupEvents.modifyCreativeTab(tab, event => {
    event.remove(stack => isDefaultTaczVariant(stack) || DEAD_TABLES.includes(stack.kjs$getId()))
  })
})
