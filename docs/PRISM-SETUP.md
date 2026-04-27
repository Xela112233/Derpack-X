# Prism Launcher setup

## First-time install (you and friends)

1. **Prism Launcher** → Add Instance → Custom
   - Name: `Derpack X: Create`
   - Minecraft version: `1.21.1`
   - Mod loader: `NeoForge 21.1.105` (or whatever version is pinned in `pack.toml`)

2. **Allocate memory** → Edit Instance → Settings → Java
   - Min: `8192M`
   - Max: `12288M`
   - Java args: paste Aikar's flags (same set as the NTNH server tuning):
     ```
     -XX:+UseG1GC -XX:+ParallelRefProcEnabled -XX:MaxGCPauseMillis=200
     -XX:+UnlockExperimentalVMOptions -XX:+DisableExplicitGC -XX:+AlwaysPreTouch
     -XX:G1NewSizePercent=30 -XX:G1MaxNewSizePercent=40 -XX:G1HeapRegionSize=8M
     -XX:G1ReservePercent=20 -XX:G1HeapWastePercent=5 -XX:G1MixedGCCountTarget=4
     -XX:InitiatingHeapOccupancyPercent=15 -XX:G1MixedGCLiveThresholdPercent=90
     -XX:G1RSetUpdatingPauseTimePercent=5 -XX:SurvivorRatio=32
     -XX:+PerfDisableSharedMem -XX:MaxTenuringThreshold=1
     ```

3. **Install mods** (two options):

   **Option A — packwiz-installer (recommended once you self-host the manifest):**
   - Edit Instance → Settings → Custom commands → Pre-launch:
     ```
     "$INST_JAVA" -jar "$INST_MC_DIR/packwiz-installer-bootstrap.jar" -g \
       https://your-host.example.com/derpack-x-create/pack.toml
     ```
   - Drop `packwiz-installer-bootstrap.jar` (from `.tools/`) into the instance's `.minecraft/`.
   - Launch — it'll sync mods on every start.

   **Option B — manual zip drop:**
   - Download the latest server pack zip from the GitHub release.
   - Extract `mods/`, `config/`, `defaultconfigs/`, `kubejs/` into the instance's `.minecraft/`.

4. **Launch**.

## Sharing with friends

Easiest path: export the configured Prism instance (Right-click instance → Export) and send them the resulting zip + a one-line note about Java args. They drop it into Prism's instances folder and it's ready.

## Updating

- packwiz-installer flow: just relaunch, it syncs.
- Manual flow: download the new server zip, overwrite `mods/` and `config/`.
