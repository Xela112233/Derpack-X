# Mining stamina

Breaking blocks makes you hungry. Every solid block you mine adds a little exhaustion — the same
hidden meter that sprinting and jumping feed — so a long dig eats into your food bar just like a
long sprint would. Clearing dirt is cheap, hauling out a vein of ore is noticeably dearer, and
chopping through solid metal is hard labor. Light decoration — torches, flowers, crops' free
cousins — costs nothing.

The practical upshot: pack food for ore runs, and treat a big excavation as something your kitchen
has to fund. That is deliberate — the pack's survival systems are meant to pull on the food economy
(see the design docs if you care about the why; this page is the what).

## How much it costs

The cost is a product of two things: **what you break** and **what you break it with**.

Every mining tool shows its half of the equation in its tooltip, right under the mining-speed line:

> Stamina drain: x1 (~10 stone per shank)

"Per shank" means one drumstick on the hunger bar. The number is calibrated to standard blocks —
stone, dirt, concrete — so the weight classes below scale it up or down from there.

### What you break

| Class | Examples | Stone-equivalent cost |
|---|---|---|
| Free | Torches, flowers, foliage, redstone wire, rails, snow layers | Nothing |
| Crops | Wheat, carrots, modded crops | A quarter block — light farming cost |
| Light | Planks, wool, glass, sand, gravel, ice | Half a block |
| Standard | Stone, dirt, concrete, bricks, logs | The anchor: ~10 per shank with an iron tool |
| Heavy | Ores, obsidian, deepslate, basalt | ~1.6 blocks each |
| Very heavy | Metal storage blocks, ancient debris | ~2.2 blocks each |

The ranking deliberately mirrors what blocks weigh on a Create: Aeronautics airship — heavier to
fly is heavier to mine.

### What you break it with

Better tools tire you less. A bare hand (or a sword — it is not a mining tool) is the worst way to
dig; each tier of real tool cuts the drain:

| Tool | Drain | Stone per shank |
|---|---|---|
| Bare hand, sword, anything not a tool | x2.5 | ~4 |
| Wooden or golden tool | x2 | ~5 |
| Stone tool | x1.5 | ~7 |
| Iron tool | x1 | ~10 |
| Diamond tool | x0.75 | ~13 |
| Netherite tool | x0.5 | ~20 |

Modded tools whose material the pack does not recognize count as iron. Efficiency speeds up your
mining but does not change the drain — you get tired per block, not per second.

## What does not cost stamina

- **Machines.** Create drills, TNT, and anything else that breaks blocks without a player swing
  pay nothing. Automating the quarry is the intended escape hatch.
- **Creative and spectator mode.**
- **Saturation absorbs the first hits.** As with any exhaustion, the hidden saturation buffer
  drains before the visible bar moves — after a good meal you can dig for a while before the
  drumsticks start dropping.

## Reading the numbers

A full hunger bar is ten shanks, so an iron pickaxe pays one full bar per ~100 stone — roughly a
7×7 room, two blocks high. The same room in deepslate costs about a bar and a half; dug out by
hand, nearly four. A netherite pick halves all of it. If a dig feels mis-priced, say so — the
weights are a tuning surface and playtest reports move them.
