# Serving the pack — clients and servers

Project Commonwealth is distributed as packwiz manifests and installs cleanly as **both** a
client and a dedicated server. This is the reference for the server side and for
keeping the two in sync.

## Two delivery paths, one `pack.toml`

- **Client** — Prism. CI builds a Prism installer zip on release (`build.yml`);
  `scripts/build-prism-*.sh` build it locally. Clients get every mod (`side`
  `both` + `client`).
- **Server** — packwiz `side` filtering. The live box runs the
  `itzg/minecraft-server` image with `TYPE=NEOFORGE`; its packwiz-installer pulls
  the **server-bound** mods (`side` `both` + `server`, never `client`) from
  `pack.toml`. `scripts/build-server.sh` bakes the same set into a standalone
  server zip.

Both read the same `pack.toml`, so the pack is the single source of truth and the
two sides can't drift on content — only on `side` tags (below).

## The deploy config lives in its own repo

The box's run config — `docker-compose.yml` + `auto-update.sh` — is **not** pack
content. It lives in its own repo, **`theasshats/pcmc-server`** (split out
like the site); see that repo's README for the deploy runbook.

`auto-update.sh` there resolves one immutable commit from a channel (`release` =
production / `branch` = playtest, e.g. `v0.7.0`), reads `pack.toml` at that
commit, and pins **both** the NeoForge version and a commit-pinned `PACKWIZ_URL`
into `.env` together — so the loader and the mod set always come from the same
commit and can't desync on a restart.

The only repo-side server task is the NeoForge `[versions]` bump in `pack.toml`
(a human decision; the deploy script only mirrors it onto the box).

## Side metadata — the rule that keeps both installs healthy

A client-only mod that the server pulls can crash it or, worse, break
client/server registry parity so players can't connect. So:

- **New mods default to `side = "both"`** (Modrinth's per-version side data is
  often wrong). Narrow to `client`/`server` only **later, with a reason.**
- **`client`** = pure rendering / UI / input with no server half and no
  registered content (safe for the server to lack): the Sodium bridge
  (`sodium`, `sable-new-sodium-compat`), minimaps, `controlling`, `mouse-tweaks`,
  texture/animation/culling mods, etc. — 22 mods today.
- **`server`** = no client half — **0 today, and the bar is higher than "no client half":**
  `side = "server"` also means **single-player never gets the mod** (SP runs everything in the
  client install), so a behavior mod narrowed to `server` silently vanishes from SP worlds.
  `smarter-farmers` was exactly this (narrowed reason-less in PR #88; farmers didn't replant in
  SP until #312 caught it — flipped back to `both`). Narrow to `server` only for a mod that is
  *meaningless* in single-player, and write the reason in the manifest.
- Everything else is `both`. Mods that *look* client but **must** be `both`
  because they have a server half: `appleskin`, `jade`, `ping-wheel`,
  `simple-voice-chat`, `no-chat-reports`, `playeranimator`, and any mod that
  registers items/blocks (music discs, content). Do not narrow these.

The `side` tags themselves are correct — no client-only mod is mis-routed to the
server. A separate, mod-level failure mode still exists: a `both`-tagged mod can
carry its own server-side crash (e.g. loading a client class from common code),
which only a real server boot catches — not CI, not the `side` tags. Live
example: **appleseed 2.0.3** did exactly this (a client `KeyMapping` constructed
from its common mod constructor) and blocked server boot; the pack carried the
fix as a pinned fork build (`theasshats/Diet---AppleSeed-Edition-fix`, 2.0.3.1)
until upstream shipped it — **appleseed 2.1.0** ("fix the initialization bug on
dedicated server") includes it, and since v0.7.1 the manifest points back at
Modrinth, unpinned. The lesson stands: only a real server boot catches this class.

### Optional: slimming the server (playtest-gated)

A handful of pure-cosmetic client mods currently ride `both` and could be
narrowed to `client` to slim the server (cosmetic-only, no registered content):
`fogoverrides`, `particle-rain`, `blood-n-particles-mod`,
`better-animations-collection`, `immersive-armor-hud`, `build-guide`,
`durability-tooltip`, `better-ping-display`,
`immersivemusicmod`, `sound-physics-remastered`. This is an optimization, not a
fix — narrow only after a playtest confirms clients still connect (a wrong narrow
on a content-registering mod breaks registry parity). Leaving them `both` is
harmless.
