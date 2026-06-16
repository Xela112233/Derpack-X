// Mining-speed tooltip — in-pack replacement for the removed Mining Speed Tooltips mod.
//
// That mod was cut in 0.6.3 along with txnilib, whose jar secretly bundled the broken
// Forgified Fabric API behind the death-screen disconnect (#286, #297), and the one
// replacement candidate surveyed (item-tooltips) doesn't show mining speed. This script
// is the whole feature: every mining tool (anything extending DiggerItem — pickaxes,
// axes, shovels, hoes, and modded paxels/hammers that subclass it) gets a tooltip line
// with its mining speed, read from the `minecraft:tool` data component where 1.21.1
// keeps it — and the line reflects the hovered stack's Efficiency enchantment.
//
// Why this rev (#302): the previous rev called `ItemEvents.dynamicTooltips(id, …)` ALONE
// and rendered no line. Per the KubeJS 1.21 docs, `dynamicTooltips` is only the *content*
// half — the handler never fires unless the item is first MARKED as carrying a dynamic
// section via `ItemEvents.modifyTooltips` + `tooltip.dynamic(<key>)`. The earlier
// "requiredTarget = the stack's own id" reading of the bytecode was wrong: the target is
// an arbitrary dynamic-section key you attach in modifyTooltips, not the item id. So this
// rev pairs the two events on one shared key, and passes a single Text to `event.add`
// (the dynamicTooltips signature is `add(Component)`, not `add(Component[])` — the old
// array arg was the second reason nothing showed).
//
//   - modifyTooltips: mark every digger id with the `pcmc:mining_speed` dynamic key (once,
//     at load) so the section exists on those items.
//   - dynamicTooltips('pcmc:mining_speed', …): per-hover, stack-aware — look the hovered
//     item's base speed up in the load-time table and render base + Efficiency bonus.
//   - Efficiency in 1.21 is data-driven (data/minecraft/enchantment/efficiency.json):
//     `levels_squared, added: 1` onto `player.mining_efficiency`, i.e. bonus = lvl² + 1,
//     applied by vanilla only while the tool is the correct one. The line shows
//     base + bonus with the base alongside, e.g. "Mining speed: 34 (base 8, Efficiency V)".
//   - Haste/conduit boosts are player state, not item info — deliberately not shown.
//
// Notes:
//   - Swords and shears also carry tool components (cobweb/leaf speeds) but aren't
//     mining tools; the DiggerItem filter keeps the line off them on purpose.
//   - A digger's component keeps defaultMiningSpeed at 1.0 and carries the real tier
//     speed in its mineable-tag rule, so the fastest rule speed is the base shown.
//   - Defensive: per-item and per-hover try/catch — one odd modded item loses only its
//     own line, and an enchant-read failure falls back to the base line.
//
// VERIFY IN PLAYTEST (static CI can't check this): iron pickaxe shows "Mining speed: 8"
// (netherite 9); iron pickaxe + Efficiency V shows "Mining speed: 34 (base 8,
// Efficiency V)"; a sword/shears shows no line; a modded tool shows one; no kubejs
// errors in the client log on load or on hover.

const $BuiltInRegistries = Java.loadClass('net.minecraft.core.registries.BuiltInRegistries')
const $DataComponents = Java.loadClass('net.minecraft.core.component.DataComponents')
const $DiggerItem = Java.loadClass('net.minecraft.world.item.DiggerItem')

const DYNAMIC_KEY = 'pcmc:mining_speed'
const ROMAN = ['I', 'II', 'III', 'IV', 'V', 'VI', 'VII', 'VIII', 'IX', 'X']

// Load-time pass: id -> base mining speed, for every DiggerItem whose fastest rule beats 1.
const BASE = {}
$BuiltInRegistries.ITEM.forEach(item => {
  try {
    if (!(item instanceof $DiggerItem)) return
    const tool = item.components().get($DataComponents.TOOL)
    if (!tool) return
    let base = tool.defaultMiningSpeed()
    tool.rules().forEach(rule => {
      const s = rule.speed()
      if (s.isPresent() && s.get() > base) base = s.get()
    })
    if (base <= 1) return
    const id = $BuiltInRegistries.ITEM.getKey(item).toString()
    BASE[id] = Math.round(base * 10) / 10
  } catch (err) {
    // swallow — worst case this one item is missing its line, never a load failure
  }
})

function efficiencyLevel(stack) {
  try {
    const enchants = stack.enchantments // kjs$getEnchantments -> ItemEnchantments
    if (!enchants || enchants.isEmpty()) return 0
    let level = 0
    enchants.entrySet().forEach(entry => {
      if (entry.getKey().getRegisteredName() === 'minecraft:efficiency') level = entry.getIntValue()
    })
    return level
  } catch (err) {
    return 0 // unreadable enchantments -> just show the base line
  }
}

// Mark every digger id as carrying the dynamic section (the half the old rev was missing).
ItemEvents.modifyTooltips(event => {
  Object.keys(BASE).forEach(id => {
    try {
      event.modify(id, tooltip => tooltip.dynamic(DYNAMIC_KEY))
    } catch (err) {
      // swallow — a single bad id must not abort the whole registration pass
    }
  })
})

// Supply the content per hover, reading the hovered stack's base + Efficiency.
ItemEvents.dynamicTooltips(DYNAMIC_KEY, event => {
  try {
    const base = BASE[event.item.id]
    if (base === undefined) return
    const lvl = efficiencyLevel(event.item)
    if (lvl > 0) {
      const total = Math.round((base + lvl * lvl + 1) * 10) / 10
      const name = lvl <= 10 ? ROMAN[lvl - 1] : `${lvl}`
      event.add(Text.gray(`Mining speed: ${total} (base ${base}, Efficiency ${name})`))
    } else {
      event.add(Text.gray(`Mining speed: ${base}`))
    }
  } catch (err) {
    // swallow — a tooltip must never crash the client
  }
})
