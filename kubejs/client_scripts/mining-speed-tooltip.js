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
// API note — third rewrite; why this one uses the NATIVE NeoForge event (#302).
// History, all verified against kubejs-neoforge-2101.7.2-build.368 bytecode:
//   1. `ItemEvents.tooltip` — KubeJS 6 API, removed in the 1.21 rewrite; fails the whole
//      script at load ("Unknown event"). Confirmed in-game.
//   2. `ItemEvents.modifyTooltips` (static per-item lines) — loaded clean, rendered
//      nothing in the maintainer playtest. Never root-caused; one suspect is the Text
//      wrapper rewrite (`Text.gray` now takes a MutableComponent, not a String), whose
//      failure the per-item try/catch would have swallowed silently.
//   3. `ItemEvents.dynamicTooltips(itemId, cb)` — INERT BY DESIGN, the actual #302 bug:
//      the event's required STRING target is not the hovered item's id. It is an
//      arbitrary dynamic-action id, and KubeJSClientEventHandler.onItemTooltip only
//      posts the event for `DynamicTextAction` entries carried by a *registered*
//      ItemTooltipData (from `modifyTooltips`'s `text.dynamic(id)` or a server datapack
//      `kubejs/text/tooltip` file). With no such registration the event is never posted,
//      so per-item-id subscriptions can never fire — no line, no error.
// Rather than re-enter that plumbing through the one path that already failed in-game
// unexplained (attempt 2), this rewrite steps off it entirely:
// `NativeEvents.onEvent(ItemTooltipEvent, ...)` — KubeJS's native-event bridge, which
// registers a plain game-bus listener (reload-tracked per ScriptType, so F3-style script
// reloads don't double-register). ItemTooltipEvent is the exact event KubeJS's own
// tooltip handler consumes; we append straight onto `event.getToolTip()`, the list the
// screen renders, with zero KubeJS tooltip machinery in between. Components are built
// with vanilla API only (Component.literal + ChatFormatting) to keep Rhino wrapper
// overload resolution out of the failure surface.
//
// Mechanics (unchanged from the previous attempt, both verified against Mojang's
// official 1.21.1 mappings):
//   - Efficiency in 1.21 is data-driven (data/minecraft/enchantment/efficiency.json):
//     `levels_squared, added: 1` onto `player.mining_efficiency`, i.e. bonus = lvl² + 1,
//     applied by vanilla only while the tool is the correct one. The line shows
//     base + bonus with the base alongside, e.g. "Mining speed: 34 (base 8, Efficiency V)".
//   - Haste/conduit boosts are player state, not item info — deliberately not shown.
//   - Swords and shears also carry tool components (cobweb/leaf speeds) but aren't
//     mining tools; the DiggerItem filter keeps the line off them on purpose.
//   - A digger's component keeps defaultMiningSpeed at 1.0 and carries the real tier
//     speed in its mineable-tag rule, so the fastest rule speed is the base shown.
//   - Defensive: the whole per-hover handler is try/catch-wrapped — a tooltip must
//     never crash the client; worst case is a missing line.
//
// VERIFY IN PLAYTEST (static CI can't check this): iron pickaxe shows "Mining speed: 8"
// (netherite 9); iron pickaxe + Efficiency V shows "Mining speed: 34 (base 8,
// Efficiency V)"; a sword/shears shows no line; a modded digger (paxel/hammer) shows
// one; the line appears in both inventory and JEI tooltips; no kubejs errors in the
// client log on load or on hover. Stamina line (zagwar, PR #320): iron tools show
// "Stamina drain: x1 (~10 stone per shank)", netherite x0.5 (~20), wood/gold x2 (~5),
// modded-tier diggers x1.
//
// The stamina line mirrors server_scripts/mining-hunger.js (the mining-costs-hunger
// pressure mechanic): drain = EXHAUSTION_PER_WEIGHT x block weight x tool-tier
// multiplier, and the per-TOOL part is the tier multiplier, anchored player-facing as
// "standard blocks per shank" (one drumstick = 2 points = 8 exhaustion; 8 / (0.8 x mult)
// = 10/mult on a weight-1.0 block — zagwar's preferred framing). Client and server
// scripts can't share code, so the tier table and id rules below are a deliberate copy
// of toolMultiplier() there — KEEP THEM IN SYNC (a note in mining-hunger.js points back
// here).

const $ItemTooltipEvent = Java.loadClass('net.neoforged.neoforge.event.entity.player.ItemTooltipEvent')
const $DataComponents = Java.loadClass('net.minecraft.core.component.DataComponents')
const $DiggerItem = Java.loadClass('net.minecraft.world.item.DiggerItem')
const $Component = Java.loadClass('net.minecraft.network.chat.Component')
const $ChatFormatting = Java.loadClass('net.minecraft.ChatFormatting')

const ROMAN = ['I', 'II', 'III', 'IV', 'V', 'VI', 'VII', 'VIII', 'IX', 'X']

// Mirror of mining-hunger.js toolMultiplier() — same substring rules, same order, same
// values (wood and gold share the wood tier; an unrecognised tool material is neutral).
function staminaMultiplier(itemId) {
  const p = itemId.indexOf(':') === -1 ? itemId : itemId.substring(itemId.indexOf(':') + 1)
  if (p.indexOf('netherite') !== -1) return 0.5
  if (p.indexOf('diamond') !== -1) return 0.75
  if (p.indexOf('iron') !== -1) return 1.0
  if (p.indexOf('stone') !== -1) return 1.5
  if (p.indexOf('wood') !== -1 || p.indexOf('gold') !== -1) return 2.0
  return 1.0
}

function efficiencyLevel(stack) {
  try {
    const enchants = stack.get($DataComponents.ENCHANTMENTS) // ItemEnchantments (EMPTY default)
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

NativeEvents.onEvent($ItemTooltipEvent, event => {
  try {
    const stack = event.getItemStack()
    if (!(stack.getItem() instanceof $DiggerItem)) return
    const tool = stack.get($DataComponents.TOOL)
    if (!tool) return
    let base = tool.defaultMiningSpeed()
    tool.rules().forEach(rule => {
      const s = rule.speed()
      if (s.isPresent() && s.get() > base) base = s.get()
    })
    if (base <= 1) return
    base = Math.round(base * 10) / 10

    let text
    const lvl = efficiencyLevel(stack)
    if (lvl > 0) {
      const total = Math.round((base + lvl * lvl + 1) * 10) / 10
      const name = lvl <= 10 ? ROMAN[lvl - 1] : `${lvl}`
      text = `Mining speed: ${total} (base ${base}, Efficiency ${name})`
    } else {
      text = `Mining speed: ${base}`
    }
    event.getToolTip().add($Component.literal(text).withStyle($ChatFormatting.GRAY))

    const mult = staminaMultiplier(stack.kjs$getId())
    // One shank (drumstick) = 2 hunger points = 8 exhaustion; a weight-1.0 block costs
    // 0.8 x mult, so blocks per shank = 10 / mult. (Framing per zagwar, PR #320.)
    const perShank = Math.round(10 / mult)
    event.getToolTip().add($Component.literal(`Stamina drain: x${mult} (~${perShank} stone per shank)`).withStyle($ChatFormatting.GRAY))
  } catch (err) {
    // swallow — a tooltip must never crash the client
  }
})
