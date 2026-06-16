# Project Commonwealth — Patch Notes

<!-- Style: professional, player-facing. No decorative emoji, sparing bold, no hype or filler.
     Each release: "## X.Y.Z — Theme", a short lead paragraph, then "### Section"s.
     Bullets where they make sense: independent changes (fixes, additions, tweaks) get one
     bullet each — lead with the change, then the why/impact in the same bullet. Keep prose
     paragraphs for narrative that doesn't enumerate (a new system's introduction, a rename).
     The "## X.Y.Z" section ships verbatim as the GitHub release body (build.yml), so it must
     read complete on its own. Shipped sections are the published record — don't restyle them. -->

## 0.7.1 — Ore rebalance

A worldgen follow-up to 0.7.0's ore retune, acting on the first playtest reports: silver was
everywhere, every vein was roughly the same size, and there was no find worth mounting an
expedition over. Ore generation only affects newly generated chunks.

### Silver is scarce again

- Silver now comes from silver country. The stray silver bands in the snowy lapis, tin, and nickel
  veins are gone (replaced by iron, copper, and gold), and silver no longer appears among the
  scattered small ores. What remains: the dedicated hill vein — now about one hill deposit in ten —
  a thin band in lead country, and traces in the rare mithril and palladium veins.

### Vein sizes now vary

- Common-metal veins (coal, iron, copper, tin, zinc, bauxite, lignite) are smaller — routine finds
  that keep a region supplied without carpeting it. Mid-tier and rare veins are unchanged.
- New mega veins: iron in mountains, copper in badlands, coal in taiga and forest. A mega vein is an
  enormous deposit at roughly one per three hundred chunks of its region — a find worth building a
  mining outpost and a trade route around.

### Regional mix retune

- Veins are lighter: ore density inside ordinary veins is halved, so a deposit fills less of its
  footprint. Where veins spawn and how often is unchanged, and the rare mega veins stay as dense as
  before.
- Oceans are no longer salt-only. Warm and lukewarm oceans and coasts carry very rare copper, zinc,
  lead, gold, and silver deposits (5% of finds each); deep, cold, and frozen oceans carry nickel,
  copper, and palladium (5% each) plus a trace of mithril (2.5%). Salt remains the dominant ocean
  find.
- Forests gain gold (5%) and copper (15%) alongside coal; hills gain gold (5%); mountains gain
  diamond (2%) and redstone (8%).
- Jungles trade some jade and diamond for lapis: lapis is now 20% of jungle finds, with jade and
  diamond each down about ten points.
- Swamps gain iron (10% of finds) alongside lignite.

### World generation

- Archwood forests are much rarer. Ars Nouveau's default biome weight handed a single biome
  roughly a sixth of the overworld; it is now at the minimum non-zero weight, so the magic
  forest is a find rather than a backdrop. Applies to newly generated chunks.

### Mods

- Prospected ore veins now draw on Xaero's maps, with a toggle button on the fullscreen map —
  it turns out this feature had never actually worked: the ore generation mod shipped the map
  renderers but never registered them. The pack now runs its own maintained build of the mod,
  which wires the integration up properly and keeps the map mods current.
- Ore veins now follow the surface landscape. Vein eligibility used to be decided by the
  underground cave biome at depth, which silently blanked or hijacked cells under every region;
  it is now decided by the biome you see on the surface. Palladium, which previously leaned on
  cave biomes, is re-homed to mountains with richer deposits. Applies to newly generated chunks.
- Diet - AppleSeed Edition updates to 2.1.0 and installs from Modrinth again. The dedicated-server
  startup fix the pack had carried as a custom build is in the official release, which also adds a
  filter to the automatic food-group calculator.
- Villager farmers replant crops in single-player worlds again. Smarter Farmers was marked
  server-only, so single-player installs never received it; it now installs everywhere.

### Interface

- Mining tools show their mining speed on the tooltip again, and the line reflects any Efficiency
  enchantment on the tool (for example, "Mining speed: 34 (base 8, Efficiency V)"). This restores the
  feature lost when the Mining Speed Tooltips mod was removed in 0.6.3.
- The recipe viewer and creative search now list only the pack's Create: Armorer guns, ammunition, and
  attachments. The stock TaCZ guns and attachments, which were already uncraftable, no longer show up
  without a recipe.
- Inventory sorting is now Inventory Essentials, replacing Inventory Profiles Next and its two
  companion mods. Sorting and bulk item movement work the same or better, with one mod instead of
  three and none of the old overlay quirks; saved inventory profiles and locked slots, which saw
  little use, are gone with it.

### Fixes

- Fixed a client crash when an AzureLib-animated mob (for example the enenra) came into view. The
  performance mod's parallel resource reloader was racing AzureLib's animation parser, corrupting it
  during world load and then crashing the game the next time one of those mobs was drawn. The
  parallel reloader is now turned off; world loads are marginally slower in exchange.

## 0.7.0 — The Create spine

The tech route gets its backbone. Create and its addons now climb four tiers — manual, water and wind,
steam, then electricity — and each step is earned through deeper production chains and scarcer
materials, never timers or research screens. Alongside the spine: a curation pass that removed twelve
mods, retuned regional ore generation, and rebuilt weapons, ammunition, and capture tools as Create
production lines.

### The Create spine

Progression now runs on a few telltale components. Steam-tier work is built on the basic burner and
steel; the electric tier's mark is the capacitor, which itself takes energised iron — proof of a powered
base; the top tier asks for boss metals, and its generator magnets take silver, the occult metal found
only in dedicated silver country. If a recipe asks for one of these, that is the game telling you what
tier you are stepping into.

- Andesite alloy, the root of nearly every Create machine, is dearer: a full iron ingot per alloy by
  hand, three nuggets in the Mechanical Mixer. The mixer is the intended route; the hand craft is the
  bootstrap.
- Heat is no longer free. Blaze burners give no passive heat: the new basic burner provides low heat for
  everyday cooking and the first boilers, while blaze burners and blaze cakes carry the hotter recipes.
  The ladder is literal — a blaze burner is built on a basic burner, so each heat rung contains the one
  below it.
- From the steam tier up, key machines call for boss materials: a Mowzie's trophy for the advanced
  energiser, Ancient Remnant metal for the airship's gyroscopic bearing, Maledictus cursium for the
  netherite jetpack.
- Create: New Age's reactor is the end-game generator. Its rods are sheathed in ignitium won from Ignis,
  and its casings call for TFMG steel and circuit boards — the full industrial chain and the end-game
  hunt both stand between a base and reactor power.
- Jetpacks consolidated into one line: Create Jetpack. The base jetpack is an electric-age Mechanical
  Crafter build around a capacitor; the netherite upgrade requires cursium, and the smithing-table
  shortcut no longer works. Create Stuff & Additions' four parallel jetpacks are no longer craftable
  and no longer appear in JEI. Create: Storage's backpack flight upgrade remains as the top of the
  flight ladder: it is built from a netherite jetpack and capacitors, so backpack flight comes after
  the flagship, not around it.
- Power gear is built from the machine it works like, each on the Mechanical Crafter: the portable drill
  carries a real Ore Excavation drill head, the flamethrower a TFMG steel fluid tank, the grappling
  whisk an aeronautics plunger launcher, the fan component a Crafts & Additions electric motor. The
  brass exoskeleton runs on capacitors; the andesite and copper exoskeletons remain earlier entry points.
- Create: Ore Excavation is an end-game operation: the drilling machine and extractor require TFMG
  circuit boards and Crafts & Additions electric motors, and drill heads are steel-bodied Crafter work
  that wears out in use. Finding veins stays early — the vein finder senses through an electron tube, the vein atlas binds
  with a Create clipboard.
- The treadmill joins the hand crank and water wheel as recognized early power — modest, and capped by
  whatever walks on it.
- TFMG's voltage system is removed for now: its generators, cables, transformers, electric machines,
  meters, and electric lights are no longer craftable and no longer appear in the recipe viewer. The
  pack runs two electricity systems whose
  conversion rates are balanced against each other — Crafts & Additions and Create: New Age — and
  TFMG's grid was a third with rates of its own. None of TFMG's production machinery is affected; the
  blast furnace, distillation, centrifuge, and winding machine all run on rotation, and gas lamps stay.
  The voltage system returns, properly integrated, with the planned electricity overhaul after 1.0.
- The hot-air balloon's adjustable burner is built around the basic burner — its only accepted fire
  source now, a raw coal block no longer qualifies — so the heat ladder carries from the kitchen to
  the boiler to the sky.
- Aeronautics flight machinery is Mechanical Crafter work: the adjustable burner, steam vent,
  propeller bearings, and smart propeller all assemble on the Crafter rather than the crafting
  table. Simple parts — propellers and goggles — stay table-craftable. Envelopes are machine work
  too: a deployer presses a stick onto wool to make them, a basin of dye recolours them, and washing
  returns them to white; the crafting-table recipes are gone.
- Create: Aeroworks is added as the top tier of flight: its digital controls are what make large or
  complex flying machines practical to pilot, and they sit at the electric tier alongside the rest of
  the end-game machinery.

### Ore generation

- Veins are larger and spread further apart than 0.6.3 — bigger bodies on a wider grid, so a deposit is
  a real find and then a real supply. Prospecting (right-click an ore) and the map overlay carry the
  finding.
- Vein markers outline in the ore's own material colour, replacing 0.6.3's temporary hot-pink debugging
  outline.
- The uranium vein no longer generates (Create: Nuclear left the pack); thorium veins now run pure
  thorium, and lead country was rebuilt as an industrial corridor — lead with zinc, silver, and copper
  bands.
- Bauxite is rarer in jungles, which read as jade and deep-diamond country again; salt veins are
  deliberately small on the coasts.
- Every block the regional veins can place is collected under the `pcmc:vein_ores` block tag, so map
  tools, filters, and commands can address vein ore in one expression.

### Curation

- The Let's Do family was trimmed: Vinery, Meadow, Bakery, and Farm & Charm removed (their food and
  cooking overlapped existing mods); Beachparty stays for its unmatched beach furniture. Wine-making
  returns through Create: Winery — grapes and apples press into pomace and must, the Wine Cellar matures
  nine wines, a cider, and two juices, all through Create machines.
- Create Ultimate Factory removed — renewable diamonds and netherite work against the scarcity economy.
- Create Metalwork removed — its molten casting duplicated TFMG, and most of its content targeted absent
  mods.
- Create Compressed removed — storage compression already served by Create's own options.
- Create: Nuclear removed for this cycle — the pack ships one end-game reactor (New Age's); the nuclear
  option is planned to return with the power-grid overhaul after 1.0.
- Traveler's Backpack removed — Create: Storage's backpack line builds through Create and keeps
  backpacks on the tech spine.
- Brazil Legends and Create: Northstar - Redux removed (the latter pending a post-beta re-review); Too
  Many Paintings removed in favour of Immersive Paintings, with custom image uploads disabled
  server-side.
- Alex's Mobs and Naturalist duplicated five animals; the Alex's Mobs versions are kept and the
  Naturalist copies no longer spawn, with their fur and catfish items unified.
- Create Aeronautics: Claims is added — a claim block that extends land-claim protection to airships,
  working with the pack's Open Parties and Claims.

### Weapons and capture tools

- Firearms are Create projects: guns assemble on the Mechanical Crafter from manufactured parts — steel
  barrels, deployer-built firing mechanisms, precision triggers, fluid-filled primers — and the
  bench-style gun table stays out of the pack.
- The self-loading weapons (SMG, both shotguns, the rifles, the snipers) cycle on a Create Big Cannons
  recoil spring; the three heavies (machine gun, grenade launcher, 40mm field cannon) are built around a
  steel autocannon barrel from the same mod. Small arms grow out of the artillery industry.
- Ammunition is a four-stage production line: casings cut leaner from brass (half the old yields), then
  each round takes a primer, a gunpowder fill, a deployed lead projectile, and a crimping press. The
  40mm shell is finished with a Big Cannons impact fuze, and cannon powder charges are filled from
  the armory's gunpowder fluid at a spout — the hand-packed gunpowder route is gone, so the two
  calibers share one propellant line.
- The telephoto sniper scope gains an electron-tube rangefinder.
- Heavy autocannon parts can now actually be cast: Advanced Technologies' casting recipes pointed at
  a fluid tag nothing registered, so the whole line was unpourable. The tag is now filled, and steel
  casts draw the same molten steel as the rest of the pack.
- Mob lassos are priced for what they do: the golden lasso is built around Create electronics, the
  hostile lasso takes a monstrous horn from the Ancient Remnant, and the diamond lasso merges the
  golden, aqua, and hostile lassos with diamonds and one boss trophy of your choice. Bosses themselves
  cannot be lassoed, by any tier.

### Recipes and survival

- Steel is consolidated under TFMG: ingot, nugget, block, and plate all resolve to the TFMG item, so
  steel reads as one material across every recipe. The same goes for the molten form — melting steel
  down in the Big Cannons foundry now yields TFMG's molten steel, which pours cannons and casts back
  to ingots; Big Cannons' duplicate fluid is gone, and everything runs on one scale: 144 mb to the
  ingot, 16 to the nugget, with the foundry's melting and casting updated to match. One known gap:
  cannon casts measure their volume in the mod's old units, so cast cannons currently cost fewer
  ingots than intended — a recalibration is planned. The steel plate
  also has a single recipe now — TFMG's press line — with Ironworks' one-step press shortcut removed.
  Bronze remains with Create: Ironworks.
- End-game electricity reaches outside the tech route: the layered magnet at the heart of New Age's
  generators takes a silver core, pulling a region most tech bases don't hold — mine silver country or
  trade with whoever does.
- Your machines warm you under Cold Sweat: blaze burners by their heat level, the basic burner gently,
  kitchen stoves and ovens, New Age's electric heater, and TFMG's blast stove and coke oven — a working
  steelworks is a warm place to winter.
- Immersive Paintings' craftables print through Create: a painting is dye deployed onto paper, graffiti
  is dye, gunpowder, and a bucket mixed in a basin, and the glow variants take glow ink at a deployer.

### Game Discs

The Game Discs handheld console and its parts build through Create. Pong is craftable to get you
started; the other six titles are found in structure chests while exploring, the rarer ones turning up
less often.

## 0.6.3 — Ore generation, prospecting, and mining hunger

Fixes to regional ore generation and the prospecting guide, plus a new survival mechanic: mining now costs hunger.

### Ore generation

- Aluminium (bauxite), lignite, and fireclay no longer generate in every overworld biome. They rode along in a background stone feature, which sidestepped the regional vein system; they now come from dedicated regional veins — bauxite in jungle and savanna, and lignite, with fireclay as its secondary band, in swamps and lowland plains.
- The background stone feature still generates as terrain decoration but no longer carries any ore. Lead's stray secondary source there was removed for the same reason; lead remains available from its regional vein.

### Prospecting

- The in-game prospecting guide was corrected: veins are located by right-clicking an ore block, not by breaking it — breaking does nothing. The guide now covers the right-click trigger, the map overlay toggles, marking veins as depleted, waypoints, sharing prospecting data with teammates, and the JEI ore-vein diagram.
- Vein markers on the map now draw with a hot pink outline so they stand out against terrain. Previously the outline was transparent and effectively invisible.

### Mining hunger

Breaking blocks now drains hunger, so a big dig needs provisioning — or machinery.

- The cost scales with the block: roughly 100 stone or dirt per hunger bar with an iron tool; planks, wool, and glass cost about half that; ores and metal blocks noticeably more.
- The tool matters: each tier from wood up to netherite tires you less, and digging with no tool at all drains fastest.
- Harvesting crops costs a small amount. Decorative blocks such as torches, flowers, and saplings cost nothing, and creative mode is exempt.
- Automation pays nothing: Create drills and other machine breakers drain nobody's hunger, so machines remain the answer to industrial-scale digging.

### Fixes

- Dying no longer disconnects you from the server. A library bundled inside another mod (Forgified Fabric API) crashed while opening the death screen, and the resulting packet error kicked the player to the server list on every death. The pack now ships the library explicitly at a newer build that contains the fix.
- TxniLib and the two small interface mods that pulled it in — Immersive Armor HUD and Mining Speed Tooltips — have been removed. The library carried its own hidden copy of the Fabric API build involved in the death-screen bug; the pack now ships that library once, explicitly, at the fixed version.
- Create Big Cannons moved back from 5.11.5 to 5.11.3. The Advanced Technologies addon is built against 5.11.3 and calls a block-damage method whose signature changed in the newer build, so firing a heavy autocannon armour-piercing round at a block crashed the server.
- Firing a cannon no longer crashes the client while rendering its muzzle smoke. The smoke particle uses a vertex format the renderer's optimized particle path did not expect; that path is now disabled for billboard particles, which fall back to the standard renderer with no visible difference.

### Other

- This build bundles an X-ray resource pack (disabled by default) as a debugging aid for verifying the ore generation changes during the pre-alpha. It will be removed in a later release.


## 0.6.2 — Project Commonwealth

The pack has been renamed to Project Commonwealth (motto: Magna Communitas). The new name runs through the installer, the launcher instance, the documentation, and the pack's internal data. Mods, recipes, balance, and world generation are unchanged from 0.6.1.

Existing worlds and installs continue to work, and ore generation is unchanged.

## 0.6.1 — Dedicated server fix

Diet – AppleSeed Edition, which provides the diet and nutrition system, crashed dedicated servers on startup. It now uses a patched build that resolves the crash. Single-player and client installs are unaffected.

## 0.6.0 — Foundation

This release concentrates on performance and maintenance: a new renderer, a full pass of mod updates,
spawn tuning, and a set of fixes and removals.

### Performance

Sodium has been added. It replaces the default terrain renderer, which carried out its work on the
CPU rather than the graphics card and was the primary cause of low frame rates at an established
base. In testing, frame rates approximately doubled. Sodium could not be used previously because it
left all world textures invisible alongside Create: Aeronautics; a compatibility layer now reconciles
the two, so Sodium runs together with Aeronautics and its airships. Sodium is client-side only and
has no effect on the server. Shaders remain unsupported with Aeronautics.

Natural spawn rates have been reduced for several of the more demanding creature mods — Born in
Chaos, Mowzie's Mobs, and Grimoire of Gaia. These mobs were appearing often enough to take up a
noticeable share of server tick time. They still spawn, only less frequently. This is an initial
adjustment and will be revised with further play.

Chunky, a chunk pre-generation tool, has been added. Running `/chunky radius <n>` followed by
`/chunky start` on the server generates the surrounding area in advance, so chunks are saved before
players reach them and live world generation does not stall a session.

Memory and garbage-collection settings were tuned. The client instance ships with G1GC at 8 GB minimum /
12 GB maximum and a 50 ms collection-pause target, so the brief periodic GC hitches are shorter and less
noticeable on high-refresh displays. The setup guide adds a system-RAM-to-heap table and an optional
Generational ZGC configuration for machines with 24 GB or more (`docs/PRISM-SETUP.md`). These ship in the
installer, so a fresh install or reinstall picks them up.

### Updates

Approximately a month of mod updates has been applied across the pack. NeoForge has been raised to
21.1.233 to meet the requirements of the updated mods. Update your launcher's loader version to
21.1.233 before playing.

Configured has been added. It provides in-game configuration screens for mods that support it,
reachable from the mod list, so per-mod settings can be changed without editing config files by hand.

### Fixes

The locked-slot indicators from Inventory Profiles Next no longer appear as stray arrows on the
hotbar. Slot locking continues to function normally.

Cactus is no longer offered as a building material in the cases where it produced broken decoration
blocks.

Iron's Spellbooks spellbooks now fit the dedicated spellbook accessory slot.

### Recovered content

A number of recipes and block drops that several mods intended to ship were being discarded at load
because of malformed files, with no visible sign beyond a missing recipe or an empty drop. Each was
checked against the mod's own files and restored.

Recipes: the Farm and Charm mincer now processes the ores and logs from Create: TFMG, Occultism,
Iron's Spells, Let's Do Vinery, and Create that it was meant to; Create: Compressed's crushed-pile
washing, dough, and polished rose quartz work again; and the cherry-sapling cutting recipe works.
Crushed raw silver smelts, blasts, and washes into Occultism silver. Create: Compressed's compressed
crushed uranium pile, which previously had no use here, now washes into Create: Nuclear uranium powder.

Block drops: the Create treadmill, the komainu, kawauso, and tanuki statues, the Stoneholm cleric
chest, and the Trail and Tales pottery cooking pot and lantern fruits now drop correctly. Smokeleaf
Industries houses can generate in plains again.

### Spawns and content

Rotten Creatures spawning has been tuned per mob so each undead appears where it fits. The Undead
Miner is limited to mineshafts and mining dungeons; Frostbitten and the Glacial Hunter to cold biomes
and snowy forts; Swampy to swamps; the Mummy to deserts, pyramids, and crypts; Burned to ruined
portals and Nether structures; and Dead Beard to ocean ruins, with a smaller chance aboard the When
Dungeons Arise pirate ships. Spawns draw on the relevant vanilla structures together with those from
YUNG's Better Dungeons and Dungeons and Taverns.

### Food and diet

The two Spice of Life diet mods have been replaced by Diet – AppleSeed Edition, a single nutrition
system that tracks five food groups and grants small benefits for eating a varied diet. It reads each
food's groups from its recipe automatically, so it covers the whole pack without per-item setup. The
separate Touhou Little Maid diet addon is unrelated and stays.

### Removed

Bosses of Mass Destruction has been removed, along with its library dependency. Its four
structure-bosses overlapped L_Ender's Cataclysm, which remains as the pack's endgame boss mod.

Umapyoi and five unused library mods were also removed. Nothing that depended on any of these remains
in the pack.

## 0.5.4 — Guns crafted through Create

Guns, ammunition, and attachments are now built through their Create recipes, not assembled at the
gun smith table. The table no longer offers them — nor the stock TaCZ weapons — so the path is the
Create machinery the recipes were always meant to run on. Use JEI to find a gun's recipe. The gun
smith table blocks are no longer craftable, since they no longer do anything. Guns and parts you
already have are unaffected.

## 0.5.3 — Gunsmith delivery fix

Builds made with the pack editor were shipping without the Create: Armorer gun pack, so its
guns, ammo, and attachments were missing and the gun smith table had nothing to craft. The editor
now bundles the gun pack the same way the build scripts do, so a fresh build has the full Armorer
lineup. No recipes or content changed; this only restores files that should have been included.

The stock TaCZ gun smith table can no longer be crafted, so the Create: Armorer table is the one
survival path to building guns. The stock block still exists in the creative tab.

## 0.5.2 — Magic Web

The magic mods are now tied together. Ars Nouveau sits at the center: its Source acts as a common
currency that other schools' essences convert into and out of, the occult mods funnel through one shared
reagent, and a number of mob drops and metals that had no purpose now feed the magic systems. Everything
here is an addition — no existing recipe changed, so these are extra paths rather than replacements.
Conversion rates are a first pass and will be tuned with play.

### Ars Nouveau at the center
- Source gems and Iron's Spellbooks arcane essence convert into each other at the Imbuement Chamber and
  the Enchanting Apparatus, with a Galosphere shard as the catalyst so a round trip costs a crystal
  instead of being free.
- Source also produces Occultism's basic essence, and Occultism's higher otherworld essence converts
  back into the arcane line.

### The occult branch
- Ars Source and Born in Chaos spirits can be combined with water in a heated mixing basin to make
  Occult Engineering's spirit solution — the same reagent it already runs on.
- Born in Chaos spiritual dust and ethereal spirits transmute into Occultism essences on the Spirit Fire.

### Drops with a purpose
- Drops that previously did nothing now feed the web: Grimoire of Gaia's soulfire, shiny pearls, and its
  organ drops; Mowzie's Mobs' ice crystals, glowing jelly, and foliaath seeds. Reward gear is left alone.
- Modular Golems' upgrade base gains an alternative route through Ars Source.

### Magic metals through Create
- Iron's Spellbooks mithril ore and Occultism's iesnium ore can now be run through the crushing wheels
  for their raw material, alongside their existing smelting routes.

### Other
- A few metal building blocks can be recycled in the crushing wheels: Northstar tungsten sheetmetal and
  grates back into tungsten, and Quark raw-copper bricks back into raw copper.

### Known rough edges
- Conversion rates and Source costs are first-pass; expect adjustments. The deeper code-level weave
  (a block that moves Source and mana directly, cross-school spell power) is still in progress and not
  part of this release.

## 0.5.1 — Ore Veins

Ore generation is reworked to a GregTech-style model. Instead of ore scattered evenly through the
ground, ores are concentrated in large regional veins of related materials, vanilla overworld ore
generation is turned off, and a layer of scattered small ores plus a prospecting system help you find
the veins. Worldgen only affects newly generated chunks, so use a fresh world or unexplored terrain to
see it.

### Ore veins
- Ore now comes from large mix veins. A single vein holds several related ores in layers: a headline
  ore through the upper part, a second ore through the lower part, and smaller amounts of two more
  mixed in. An iron vein in the mountains, for example, also carries magnetite, some nickel, and a
  little gold.
- Veins are regional. The biome decides what generates: iron and magnetite in mountains, copper in
  badlands, redstone and lithium in deserts, coal in forests and taiga, zinc and tin in plains,
  diamond and jade in jungles, lapis in snowy biomes, emerald in peaks, salt along coasts, and so on.
  Twenty-one veins in all, covering the vanilla ores and the pack's modded metals.
- Vanilla overworld ore generation is disabled. Veins and small ores are the only sources, so where
  you settle now determines what you can mine, and what you will need to trade for.

### Small ores
- Single ore blocks of the common metals are scattered across the world. They are enough to get a
  first set of tools, and they double as indicators: a cluster of one type hints at a matching vein
  below.

### Prospecting
- Break an ore block to locate nearby veins of that material; a message confirms how many were found.
  With a supported map mod the located veins can be shown on the map with waypoints, and
  `/gtmogs share_prospection_data <player>` shares your finds with the rest of the group.

### Mods added
- GregTech Modern Ore Generation Standalone (the vein engine) and its libraries LDLib and
  Configuration.

### Known rough edges (tuning continues)
- Balance is not final. Ongoing work covers vein rarity, salt density along coasts, and ore depth.
- Cave pockets follow the local cave biome rather than the surface biome above them, so veins are
  currently absent in those pockets; the scattered small ores cover caving in the meantime. Making
  cave ore follow the surface region is the priority for the next ore-gen pass.

## 0.5.0 — Made Through Create

Crafting across the pack now routes through Create. Most tech, gear, machines, and a fair amount of
decoration are built from Create parts (sheets, casings, mechanisms) and methods (pressing, mixing,
sequenced assembly) rather than plain crafting tables, so the mods connect through the Create progression
instead of standing alone. The recipe viewer also changed from EMI to JEI, and a mob and spawn curation
pass landed. Recipe changes take effect on new crafts immediately; no world reset is required.

### Recipe viewer: EMI replaced by JEI
EMI does not display Create's processing recipes on 1.21.1 (pressing, mixing, and similar showed no recipe,
including for base Create items), which made it unusable for this pack. JEI replaces it. JER (loot, worldgen,
villager trades) and JEED (potion-effect descriptions) cover the extras EMI provided.

### Recipe overhaul (#62)
- Tech, gear, and machines route through Create, gated by complexity: everyday gear takes a single Create
  part, while endgame machines such as the MFFS force-field gear are assembled on the Mechanical Crafters
  through sequenced assembly. Affected mods include Immersive Armors, Samurai Dynasty, Modular Golems,
  Northstar, MFFS, SecurityCraft, Drones, Traveler's Backpack, the MineColonies quarries, and gliders.
- Steel and bronze are produced through Create at the source. The iron-to-steel blasting and MFFS smelt
  shortcuts were removed; steel now comes from Create Ironworks, TFMG, or Big Cannons.
- Decoration and kitchenware are included: Macaw's metal doors and windows, Supplementaries metal decoration,
  and Farmer's Delight cooking gear use Create-pressed sheets at comparable cost.
- Orphan worldgen materials (Quark stones, Create Nuclear autunite, TFMG asphalt) gained a Create use rather
  than dead-ending.
- The joke "bible" item is now a Totem-of-Undying-gated holy book, and Meadow cheese converts to and from
  Create cheese.

### Material unification
Duplicate metals from different mods are collapsed to one item each via almost-unified: steel and bronze to
Create Ironworks, lead and cast iron to TFMG, and pasta to Farmer's Delight. Note that Galosphere's "silver"
is its palladium under a legacy id, and is intentionally left out of silver unification.

### Curation and spawns (#88)
- Removed: Arphex, Creeper Overhaul, Mutant Monsters, and the orphaned Sophisticated Core library.
- Spawn control through In Control!: Grimoire of Gaia kobolds no longer spawn naturally (the Kobolds mod
  covers that role and inherits the drops), Rotten Creatures spawn only inside structures, and Mutants and
  Zombies spawn only during blood moons.
- Added: Create Aeronautics Dyeable Components, Create Stuff 'N Additions with its Sable and Aeronautics
  compat, furniture/Create/Sable compatibility, and Sable Collision Damage.

### Infrastructure
- A recipe-connectivity tool (tools/recipe-graph/) measures how connected the pack's recipes are, with an
  interactive map regenerated by CI.
- Added the five-systems design docs and roadmap, the recipe conventions, the website releases page, and CI
  checks on every pull request (packwiz index freshness, manifest lint, automatic conflict resolution).

### Open after 0.5.0
Spawn-gating observation (#108) and the modded-structure allowlist (#106), the magic-mod weave (#75 and #80),
the ore-generation model decision (#93), and functional-duplicate parts such as heavy_plate versus steel
sheet (#103).

## 0.4.6 — World Generation

> **Headline:** this release makes custom ore generation actually work. **0.4.4 and 0.4.5
> crash the moment you click _Create New World_** (a malformed worldgen file aborts the whole
> registry load) — 0.4.6 fixes that and ships the full GregTech-style ore layer. Worldgen only
> affects **newly generated chunks**: use a fresh world or explore unloaded terrain to see it.

### Critical fixes (0.4.5 could not create worlds)
- **Fixed world-creation crash** from biome modifiers that listed multiple biome tags in a JSON
  array (`"biomes": ["#tag", "#tag"]`). That form is invalid for a HolderSet and failed registry
  load with `Failed to load registries`. All 21 veins now resolve their biomes through proper
  custom biome **tags** (`#pcmc:vein_<ore>`).
- **Fixed second world-creation crash** (`Unbound values … [minecraft:ore_copper, minecraft:ore_diamond]`):
  the copper/diamond vanilla-thinning files pointed at configured features that don't exist in 1.21
  (the real names are `ore_copper_small` / `ore_diamond_small`).

### Ore generation — the system
GregTech-style economy: **rare, large, biome-specific veins** so regions specialize and trade
matters, with vanilla ore **thinned to a starter trickle** so new players aren't hard-stuck.

- **21 modded + vanilla ores veined**, each as a configured feature (the ore body, size 16–48 vs
  vanilla's ~8), a placed feature (depth + rarity), and a biome modifier keyed to a biome tag:
  coal, iron, copper, gold, redstone, lapis, diamond, emerald, zinc, thorium, magnetite, silver,
  uranium, lead, lithium, nickel, jade, tin, salt, mithril, **palladium**.
- **Vanilla ore thinned** (rarity 12–16, iron's small vein at 6) — a deliberate scarce baseline
  near spawn; veins are the real supply.
- **Mod default ore gen disabled** so our veins are the only overworld source. Shadowed with
  `neoforge:none`: Create zinc; Create: New Age thorium + magnetite; Occultism silver (+deepslate);
  Create: Nuclear uranium + lead; TFMG lead/lithium/nickel; Create: Ironworks tin (+ raw-tin
  deposits); irons_spellbooks mithril; samurai_dynasty jade; expandeddelight salt.

### Underground findability (3D biomes)
MC 1.21 biomes are 3D — underground is mostly **cave biomes** (vanilla dripstone/lush/deep_dark
plus Terralith/Galosphere/Northstar/Quark, pooled in `#c:is_underground`), which our surface-keyed
veins didn't cover, so caving turned up almost nothing. **Every vein tag now also includes the cave
biomes**, so veins are findable wherever you mine. _Trade-off under review_ (issue #58): cave pockets
are currently region-agnostic.

### New / changed ores this release
- **Palladium (Galosphere)** — added a rare, regional vein (mountains + dripstone + deep_dark,
  `chance 14`, deliberately **not** in the global cave coverage — it's the prestige exception). Also
  **stopped Galosphere's default palladium flood**: 1.5.x renamed *Silver → Palladium* but kept the
  old `ore_silver_*` feature IDs, so it was dumping small palladium blobs into every overworld biome —
  now disabled (the large crystal-canyons deposit is left as a regional bonus).
- **Born in Chaos infected diamond** — disabled its all-biome overworld injection so diamonds come
  only from the regional diamond vein + thinned vanilla. No replacement vein (by design).

### Still open (needs zagwar — see issue #58)
- Rarity/size pass — current numbers are a first cut.
- Regionality model — keep "findable everywhere," go strict-regional, or hybrid.

### Other
- Editor: removed the server `WriteTimeout` so long operations (build, batch-add, version checks)
  no longer fail with `i/o timeout`. _(Shipped in 0.4.5; noted here for completeness.)_
- `docs/ORE-GENERATION.md` expanded: HolderSet array-vs-tag gotcha, 3D/cave-biome notes, the
  shadow table, and per-ore decisions.

---

## 0.4.5 — TaCZ × Create guns + tooling

### Gun integration (Timeless & Classics × Create)
- **Stock TaCZ guns removed** — the pack is **Armorer-only**: vanilla TaCZ default guns are stripped
  (via KubeJS recipe removal, after the config route proved unreliable) so firearms come exclusively
  through the Create-integrated path.
- **All Immersive TaCZ recipes re-authored to the Create 6.0.10 schema** and routed through the
  Create **Armorer / gun smith table**, requiring Create + Immersive components (with an EMI-visible
  fallback). Recipe difficulty bumped **+2–3 notches** to fit progression.
- Added `create-immersive-tacz-integration` and the `create-armorer` pack (shipped as a native TaCZ
  pack rather than `kubejs/data`).

### Tooling / infrastructure
- **Editor:** removed the server `WriteTimeout`, so long operations (build, batch-add, version
  checks) no longer die with `write tcp …: i/o timeout`.
- **CI ground-truth digest:** new workflow + extractor that indexes every mod's real block/ore IDs,
  biome modifiers, tags, recipes, and loot tables into `tools/mod-data/`. This is what let later
  worldgen work verify IDs from the jars instead of guessing (e.g. catching the Galosphere
  Silver→Palladium rename).

---

## 0.4.4 — Eco-foundation: ore generation + Blockchain gating (PR #56)

> ⚠️ **Shipped broken** — this release introduced the ore-gen layer but a malformed worldgen file
> crashes world creation. **Do not use 0.4.4/0.4.5 for new worlds; use 0.4.6.** Documented here
> because the original notes were bare-bones.

### Ore generation — first cut
- Introduced the whole **GregTech-style vein system**: rare, large, biome-specific veins for **16
  modded + vanilla ores** (configured feature + placed feature + biome modifier each), with vanilla
  ore **thinned to a starter trickle**.
- **Disabled mods' default `#is_overworld` ore gen** via `neoforge:none` shadows so the veins are the
  intended source — Create zinc, Create: New Age thorium + magnetite, Occultism silver (+deepslate),
  Create: Nuclear uranium + lead, TFMG lead/lithium/nickel, Create: Ironworks tin, samurai_dynasty
  jade, irons_spellbooks mithril, expandeddelight salt. Let's Do Meadow alpine ores left regional.
- Modded ore IDs and biome-modifier override paths verified against the mod jars.
- _(Known issue, fixed in 0.4.6: the biome modifiers used an invalid tag-in-array form and crashed
  world creation; copper/diamond thinning referenced nonexistent configured features.)_

### Create: Blockchain gating
- **Currency Miner re-tuned into a slow sink, not a printer:** "Cost per Coin" gated to ~**10×**
  energy with a steep global ramp (`config/createblockchain-common.toml`), so minting currency is an
  expensive late-game activity rather than an early exploit.

### Project
- Added `CLAUDE.md` (repo/agent guidance) and `docs/ORE-GENERATION.md` (design + per-ore tables).
