# Governance mod — Path B technical spec (issue #260)

> **Status: DRAFT spec, not accepted.** This is the technical design for **Path B** from
> `docs/GOVERNANCE.md` — a purpose-built NeoForge mod that adds the hierarchical-governance layer
> *on top of* the claim systems the pack already runs, owning no claims of its own. It assumes the
> governance direction is pursued; the decision to commit is still open. Nothing here is built.
>
> Per `docs/CUSTOM-MODS.md`, a custom mod lives in **its own repo** under `theasshats` and reaches the
> pack via the mod-mirror packwiz pattern — it is *not* developed in this repo. The implementation is
> split into **three parts — claims → government → economy (§8)** — shipping as **three separate mods
> (decided)**: **`theasshats/pcmc-territory`**, **`theasshats/pcmc-realms`**, **`theasshats/pcmc-mint`**,
> dependency chain economy → government → claims. This doc is the design reference that seeds those
> repos.

---

## 1. Design principle — a layer, not a claim system

The single most important constraint (from #260): **do not introduce a third claim system.** The pack
already runs two — **MineColonies** colony borders and **Open Parties and Claims (OPAC)** chunk claims.
The governance mod **reads** those and hangs political state off them. An "entity" owns *references*
(colony IDs, OPAC claim keys), never raw claims. This is what keeps it from "fighting MineColonies."

```
   MineColonies colony ──┐
                         ├──► governs ──► ENTITY (municipality | federation)
   OPAC claim region  ───┘                  │
                                            ├─ tier, parent, children (recursive tree)
                                            ├─ members / capital
                                            ├─ laws  (typed rules, tier-priority resolver)
                                            └─ treasury + mint policy (Numismatics bridge)
```

---

## 2. Data model

All persisted as a server-level `SavedData` (`RealmsSavedData`, attached to the overworld), serialized to
NBT. No per-chunk data of our own — chunk→entity resolution goes through the claim mods' lookups (§6).

### Entity

| Field | Type | Notes |
| --- | --- | --- |
| `id` | UUID | Stable identity. |
| `name` | String | Player-set. |
| `tier` | enum `Tier` | See §3. |
| `kind` | enum `{MUNICIPALITY, FEDERATION}` | Municipality = solo-founded leaf; federation = composed. |
| `parentId` | UUID? | Null for a top-level entity. The recursive link. |
| `childIds` | Set<UUID> | Sub-regions / member municipalities. |
| `colonyIds` | Set<int> | MineColonies colonies this entity governs (municipalities mainly). |
| `claimKeys` | Set<ClaimKey> | OPAC claim references, if used directly. |
| `capitalId` | UUID? | For federations: the member that is the seat. |
| `members` | Map<UUID playerId, Role> | Role ∈ `{LEADER, OFFICER, CITIZEN}`. |
| `laws` | List<Law> | §4. |
| `treasury` | Treasury | §5. |
| `createdTick` | long | For promotion cooldowns / audit. |

### Tier

`Tier` is an ordered enum carrying its **rank** (for law precedence) and its **kind**:

```
MUNICIPALITY kind:  SETTLEMENT(0) → VILLAGE(1) → TOWN(2) → CITY(3)
FEDERATION   kind:  COUNTY(4)     → KINGDOM(5) → EMPIRE(6)
```

Rank is the integer in parens — used directly by the law resolver (§4) and to validate parent/child
nesting (a child's rank must be < its parent's).

---

## 3. Tier ladder & promotion

**Inputs are free from MineColonies** — no new metrics to invent or tick:

- **population** = sum of `IColony.getCitizenManager().getCitizens().size()` across `colonyIds`.
- **footprint** = claimed-chunk count (MineColonies claim radius + any OPAC claims).
- **development** = highest building level in the colony (Town Hall level is the cleanest single proxy).

Promotion is **player-initiated and validated**, not automatic (avoids surprise tier-ups and per-tick
scanning):

- A municipality leader runs `/realm promote`; the mod checks the next tier's thresholds (configurable in
  `realms-server.toml`) and either promotes or reports the unmet criterion.
- Federations are **created by composition**, not promoted into: `/realm federate <name>` requires ≥ N
  consenting member municipalities (config; default 2) of at least tier TOWN. The new federation's tier
  is the lowest legal federation tier (COUNTY) and climbs as more/larger members join.
- **Recursion:** a federation at rank ≥ KINGDOM may `/realm carve <name>` a sub-region (a child federation
  or municipality) one rank below itself, reassigning some of its members/colonies to the child.

Thresholds live in config so they're tunable per-server without a code change. All checks fire **on the
command**, never on a tick.

---

## 4. Law engine

A **law is a typed rule**, not free text — so it can be enforced and resolved deterministically.

> **Reframed in Part 2 scope.** The enforcement taxonomy below (every law "cancels the action") is
> superseded by the **hard-law / soft-law** model in **[`GOVERNANCE-REALMS-SCOPE.md`](GOVERNANCE-REALMS-SCOPE.md)**
> — the authoritative Part 2 (`pcmc-realms`) design. A law now carries an **enforcement mode**:
> `AUTOMATIC` (mod applies a fiscal effect — tax/stipend), `OFFICER` (a human enacts, the mod settles —
> fines), or `SOCIAL` (the mod never blocks; guards/players enforce — PvP). Read that doc for §4 details;
> the sketch here is kept for the data-model shape only.

```java
record Law(UUID id, LawType type, EnforcementMode mode, ScopeRef scope, Value value, boolean enabled) {}

enum EnforcementMode { AUTOMATIC, OFFICER, SOCIAL }  // see GOVERNANCE-REALMS-SCOPE.md §3

enum LawType {
  PVP,             // SOCIAL — value: ALLOW | DENY (guard aggro on DENY, not a cancel)
  ITEM_BAN,        // SOCIAL — value: item/tag set (contraband → guard aggro)
  TAX,             // AUTOMATIC — value: rate (bps); Part 3 moves the coin
  STIPEND,         // AUTOMATIC — value: amount + period; scheduled payout, Part 3 moves the coin
  FINE,            // OFFICER — an enforcement *action* (not a standing rule); auto-withdraw via Part 3
  CURFEW/optional, // SOCIAL — value: tick window (defer — TPS-sensitive)
  // BLOCK_INTERACT dropped: claim protection is OPAC/MineColonies' job, not a governance law.
}
```

### Enforcement

Each enforceable `LawType` maps to a **NeoForge event handler** that resolves the governing entity for
the action's position (§6), walks that entity's ancestry, and applies the winning law (§ resolver below):

| LawType | Mode | Event / trigger | Effect |
| --- | --- | --- | --- |
| PVP | SOCIAL | `LivingIncomingDamageEvent` (player→player) | **Do not cancel.** Mark the *unprovoked* attacker "wanted" (self-defense exempt) + trigger MineColonies guard aggro (scope §5); no guards → no effect. |
| ITEM_BAN | SOCIAL | `PlayerInteractEvent` / region-entry inventory check | Flag a violation → guard aggro; not a hard use-block. |
| TAX | AUTOMATIC | hook on Trading Floor sale + (optional) colony output | Skim `rate` into the entity treasury (§5; Part 3 moves the coin). |
| STIPEND | AUTOMATIC | coarse scheduler (scope §6) | Pay `amount` per citizen from the treasury (Part 3). |
| FINE | OFFICER | `/realm fine` command | Auto-withdraw `amount` from the target into the treasury (Part 3). |

*(BLOCK_INTERACT removed — claim protection is OPAC/MineColonies' job; a hard build/break law would
duplicate it. Griefing-as-a-crime is the SOCIAL path instead. See scope §4.)*

All handlers are **event-driven and position-local** — no global world scan, matching the perf doctrine
in `SYSTEMS.md` §3a. The chunk→entity lookup must be **O(1)-ish and cached** (§6), since these events are
hot paths (block break, combat).

### Precedence resolver (the core rule #260 asks for)

For a given action at a position:

1. Resolve the **leaf entity** governing that position (§6).
2. Walk **leaf → root** via `parentId`, collecting any law of the action's `LawType`.
3. **Highest tier rank wins.** An EMPIRE(6) PVP=DENY overrides a CITY(3) PVP=ALLOW. *(A higher tier sets
   the floor; a subordinate cannot override its parent — only fill gaps the parent leaves unset.)*
4. **Same-rank tie** (can only happen for two unrelated entities over the same position — a claim overlap):
   resolve by the more-specific (smaller-footprint) entity, else the older `createdTick`. Document the
   rule; overlaps should be rare because claims come from the underlying mods.
5. If no entity in the chain sets that `LawType`, **vanilla/pack default applies** (no law = no change).

This gives the cascading "empire > kingdom > county > city > settlement" semantics with a defined
conflict rule, and makes sub-region laws **subordinate to** (never override) their parent — the natural
reading of a hierarchy.

---

## 5. Treasury & minting (the hard part — Numismatics bridge)

This is the **least certain** part and is deferred to **Part 3 (§8)**, behind laws/hierarchy. Two
shapes were considered; the desk check has since settled it:

- **Option 1 — ledger + entity notes (recommended).** Each entity has a `Treasury` = a Numismatics
  **bank account** (Numismatics already has accounts/cards) plus the authority to issue **entity notes** —
  a custom item denominated in that entity's currency, **backed by reserves** held in the treasury account.
  A note is redeemable at the issuing entity's bank for its reserve value. Inflation control = **you can
  only mint notes against reserves** (a hard backing ratio in config), so an entity can't print free money.
  Rival currencies "compete" by their redemption trust; exchange is manual (player-to-player) or via a
  Trading Floor stall. This rides Numismatics' existing money rather than forking the coin system.
- **Option 2 — new coin denominations: ruled out.** Verified against upstream source
  (`Layers-of-Railways/CreateNumismatics`, LGPL): `Coin` is a **fixed Java enum** (SPUR/BEVEL/SPROCKET/…),
  not a registry — there is no runtime denomination registration without mixins into the coin system.
  Option 1 is the shape. (LGPL also means linking a separate mod against Numismatics is license-clean;
  whether the bank-account credit/debit surface is reachable as API or needs a small mixin is the Part 3
  spike's first question.)

**Minting authority by tier:** config-gated; default **only FEDERATION-kind entities (COUNTY+) may mint**,
so a lone settlement can't issue currency — minting is a privilege of scale, reinforcing the federation
incentive. TAX revenue (§4) flows into the treasury account; that's the **value-loop weave** #260 wants —
taxes/tribute fund the mint, minted notes circulate via Trading Floor, governance becomes a real trade
driver, not flavor.

---

## 6. Claim integration & chunk→entity resolution

The performance-critical primitive: **given a (level, chunkPos), which entity governs it?**

- **MineColonies:** `IColonyManager.getInstance().getColonyByDimensionAndId(...)` / position lookup gives a
  colony; we map `colonyId → entityId` via a reverse index kept in `RealmsSavedData`.
- **OPAC:** resolve via OPAC's published API (the `xaero.pac.api` claims manager) to a claim owner/party,
  mapped `claimKey → entityId` the same way.
- **Cache:** a `Long2ObjectMap<UUID>` (packed chunk key → leaf entity id), invalidated on claim change /
  entity edit. This keeps the hot-path handlers (§4) cheap. The ancestry walk (leaf→root) is over the
  in-memory entity tree, not storage.

Federation governance over a chunk is **derived**, not stored per-chunk: a chunk resolves to its leaf
(municipality) entity, and the ancestry walk reaches the federations above it. So claims only ever attach
at the municipality level; higher tiers govern transitively.

---

## 7. Commands & UI (MVP = commands; GUI later)

`/realm` tree (permission-checked against member `Role`):

```
/realm found <name>                 (bind current colony as a SETTLEMENT)
/realm info [name]                  (tier, members, laws, treasury)
/realm promote                      (attempt next municipality tier)
/realm federate <name>              (compose a federation from consenting members)
/realm carve <name> <tier>          (recursive sub-region)
/realm law set <type> <value>       (issue/replace a law)
/realm law list
/realm member <add|remove|role> ...
/realm mint <amount>                (federation+, against reserves)
/realm map                          (optional: render governed borders — reuse OPAC/MineColonies render)
```

A GUI (entity screen, law list, treasury) is a **v2** nicety; commands are the MVP so the system is
playable before any screen work.

---

## 8. Scope — three parts: claims → government → economy

The full vision is multi-month; it splits into **three sequential parts with hard boundaries** —
**claims** (territory resolution), **government** (hierarchy + laws), **economy** (treasury + minting).
Each part is independently playtestable on the box, each starts with its own spike, and each ships as
**its own mod** — see "Packaging" below.

### Part 1 — Claims (the territory layer)

The read-only substrate everything else stands on: **given a position, which entity governs it?**

| Delivers | Detail |
| --- | --- |
| Spike | Verify the §6 assumptions: MineColonies position/colony lookup, OPAC claims-manager lookup, both reachable from a NeoForge 1.21.1 dev env. |
| Entity registry | `RealmsSavedData` + a minimal `Entity` (id, name, members, `colonyIds`, `claimKeys`) — **no tiers, laws, or treasury yet**. |
| Resolution | The colonyId→entity / claimKey→entity reverse indexes and the chunk→entity cache + invalidation (§6). |
| Commands | `/realm found`, `/realm info`, plus a `/realm whogoverns` debug query for the chunk under the player. |
| Degradation | Soft-dep no-ops when MineColonies/OPAC absent (§9). |

**Gate:** founding binds a real colony; `whogoverns` answers correctly with no measurable hot-path cost;
the mod owns zero claims.

**Public surface for the parts above it:** `resolve(level, chunkPos) → entity chain (leaf→root)` plus
entity CRUD events — kept in a clean `api` package from day one. This is the contract Part 2 consumes.

### Part 2 — Government (the political layer)

> **Full Part 2 scope: [`GOVERNANCE-REALMS-SCOPE.md`](GOVERNANCE-REALMS-SCOPE.md)** — the hard/soft
> enforcement model, the law catalog, the exact Part 1 API contract (`TerritoryApi`/`EntitySnapshot`/
> `TerritoryEvents`/`Role`) read from the shipped territory source, the MineColonies-guard mechanism,
> the territory-mod feature recommendations, commands, slices, and the spike list. The summary below
> is the project-level view.

The hierarchy and the law engine, standing on Part 1's resolution. Two slices:

| Slice | Delivers | Gate |
| --- | --- | --- |
| **2a — single-entity laws (MVP)** | Tier field + `/realm promote` validation (§3); the **law engine + resolver** (§4) with PVP and BLOCK_INTERACT enforcement. **No federation yet.** | A single municipality sets & enforces PvP/build laws on its colony — independently shippable. |
| **2b — hierarchy** | Federation composition (consent, capital), recursive `carve`, parent/child links, roles; the full leaf→root cascade + tie-breaks. | Empire>kingdom>city precedence verified in-game. |

**Public surface for Part 3:** a **pluggable `LawType` registry**, so the economy part registers TAX
without the law engine knowing anything about money.

### Part 3 — Economy (the value layer)

The Numismatics bridge (§5) — the least certain part, deliberately last.

| Delivers | Detail |
| --- | --- |
| Spike | The Numismatics account credit/debit surface (API or small mixin — LGPL source, either workable), and **whether Trading Floor exposes any sale hook** — a desk check found no public event, so assume a mixin until proven otherwise. |
| Treasury + TAX | Entity treasury = Numismatics bank account; TAX `LawType` registered into Part 2's engine; revenue flows per §4/§5. |
| Minting | Option-1 reserve-backed entity notes, hard backing ratio in config, federation-tier gate. |
| Weave | Notes circulate via Trading Floor stalls; exchange is player-to-player or stall-mediated. |

**Gate:** taxes fund a treasury; a note is minted against reserves and redeemed at the issuing bank.

### Cross-cutting polish (last)

GUI (entity/law/treasury screens), `/realm map` border render, ITEM_BAN, config surface — after the
three parts, blocking none of them.

### Packaging — three separate mods (decided)

The parts ship as **three separate mods**, each in its own repo per `docs/CUSTOM-MODS.md`:

| Part | Repo (mod id) | Depends on |
| --- | --- | --- |
| 1 — Claims | `theasshats/pcmc-territory` (`pcmc_territory`) | MineColonies + OPAC (soft, §9) |
| 2 — Government | `theasshats/pcmc-realms` (`pcmc_realms`) | `pcmc_territory` (hard) |
| 3 — Economy | `theasshats/pcmc-mint` (`pcmc_mint`) | `pcmc_realms` + Numismatics (hard) |

This buys independent versioning and release cadence (a treasury fix doesn't re-ship the law engine),
a smaller blast radius per release, and Part 1 shipping alone as a useful "who governs here" utility.
The cost — three repos/CI/manifests and a cross-mod API kept stable — is accepted; the release
pipeline is the same one `pcmc-killfeed` already runs.

Consequences of the split:

- The boundary contracts are **inter-mod APIs**, not internal seams: Part 1's `resolve`/entity-events
  surface and Part 2's `LawType` registry are published packages the mods above compile against —
  version them deliberately (breaking the API = major bump, the dependents pin a range).
- **Persistence partitions per mod.** Each mod owns its own `SavedData` keyed by entity UUID
  (territory: the core entity record; realms: tier/hierarchy/laws; mint: treasury) — no shared NBT
  blob, so an upper mod being absent never corrupts a lower one's state. §2's unified table is the
  *logical* model across all three.
- **Lower mods never reference upper ones**: `pcmc_territory` knows nothing of laws or money.

---

## 9. Perf & compat notes

- **Event-driven, position-local only** — no global tick scan (matches `SYSTEMS.md` §3a doctrine and the
  pack's tight TPS budget). The only always-resident cost is the chunk→entity cache (small) and the
  in-memory entity tree.
- **Server-authoritative.** All law checks server-side; the mod is `side = "both"` but enforcement never
  trusts the client. (Map render is the only client-side piece.)
- **Soft-dep on the claim mods.** If MineColonies or OPAC is absent, the corresponding resolver no-ops
  rather than crashing — so the mod degrades instead of hard-failing. (In-pack both are always present.)
- **`aeronauticscompat`-style courtesy:** laws resolve by block position, so a moving Aeronautics ship is
  governed by whatever chunk it's over — acceptable; flag for playtest whether that feels right over
  contested borders.

---

_Refs: `docs/GOVERNANCE.md` (the scoping doc + path comparison this realizes); `docs/CUSTOM-MODS.md` (the
own-repo + mod-mirror pattern this mod follows); `docs/SYSTEMS.md` §3a (the event-driven/locality-bounded
perf doctrine); issue #260._
