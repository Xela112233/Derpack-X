# derpack-edit — source

Local web app for managing the modpack. Single Go binary, embedded packwiz, embedded UI.

## Repository layout

```
editor-src/
├── go.mod, go.sum
├── main.go                 # entry point — server, browser launch, route registration
├── build.sh                # local build script (rarely needed; CI handles this)
├── ui/
│   ├── index.html
│   ├── app.js
│   └── style.css
└── internal/
    ├── packwiz/            # packwiz manifest types, loader, exec wrapper, embedded binary
    ├── modrinth/           # Modrinth API client (version listing)
    ├── hashutil/           # SHA-512 of remote URLs
    ├── builder/            # assembles a Prism instance locally
    ├── handlers/           # HTTP handlers for each API endpoint
    └── config/             # user config persistence (.editor/config.json)
```

## Local development

You need Go 1.22+ installed.

First-time setup: place a packwiz Windows binary at `internal/packwiz/assets/packwiz.bin`. Either download from <https://nightly.link/packwiz/packwiz/workflows/go/main> (rename the `.exe` to `.bin`) or build from source:

```sh
git clone https://github.com/packwiz/packwiz /tmp/packwiz
cd /tmp/packwiz
go build -o /path/to/derpack-edit/tools/editor-src/internal/packwiz/assets/packwiz.bin .
```

Then iterate normally:

```sh
cd tools/editor-src
go mod tidy   # first time only
go run .      # launches the server and opens browser
```

## Build for release

CI handles this on every push to `main` that touches `tools/editor-src/**`. Workflow file: `.github/workflows/build-editor.yml`.

The workflow:

1. Checks out packwiz source
2. Cross-compiles packwiz for `windows/amd64`
3. Embeds it into the editor binary
4. Cross-compiles the editor for `windows/amd64`
5. Commits the resulting `tools/derpack-edit.exe` if it changed (with `[skip ci]` to avoid loops)

To force a rebuild manually, use the workflow's "Run workflow" button (workflow_dispatch trigger).

## API endpoints

All routes return JSON.

| Method | Path | Body | Returns |
|---|---|---|---|
| GET | `/api/pack` | — | `{name, version, author, pack_format}` |
| GET | `/api/mods` | — | array of mod entries |
| POST | `/api/mods/add` | `{source, slug, side}` | op result |
| POST | `/api/mods/remove` | `{slug}` | op result |
| POST | `/api/mods/pin` | `{slug, pinned}` | op result |
| POST | `/api/mods/update` | `{slug}` (empty for all) | op result |
| POST | `/api/mods/set-version` | `{slug, version_id}` or `{slug, url, filename?}` | op result |
| POST | `/api/mods/compute-hash` | `{slug}` | op result |
| GET | `/api/mods/{slug}/versions` | — | array of version entries (Modrinth only) |
| GET | `/api/config` | — | `{prism_instance_path, prism_default_root, detected_instances}` |
| POST | `/api/config` | `{prism_instance_path}` | op result |
| POST | `/api/build` | — | `{ok, staging_dir, mod_count, logs}` |
| POST | `/api/launch-prism` | — | `{ok, mod_count, prism_path, hint}` |

Op result shape:

```json
{ "ok": true|false, "slug"?: "...", "output"?: "...", "error"?: "..." }
```

Errors that don't make it to a structured op result come back as HTTP errors:

```json
{ "error": "..." }
```

## Architecture notes

**Subprocess vs reimplementation.** The editor uses packwiz as a subprocess for all the operations that have direct CLI equivalents (`add`, `remove`, `pin`, `unpin`, `update`, `refresh`). It reimplements:

- Version listing — Modrinth API, since `packwiz mr list-versions` doesn't exist
- Hash computation — direct HTTP fetch + SHA-512, faster than `packwiz refresh` re-validating everything
- Build flow — reimplemented in Go to avoid bash dependency on Windows

**Embedded packwiz.** `//go:embed assets/packwiz.bin` bundles the packwiz binary into the editor `.exe`. On first run, it's extracted to `.editor/packwiz-cache/packwiz.exe` and invoked via `os/exec`.

**Embedded UI.** `//go:embed ui/*` bundles the HTML/CSS/JS into the binary. Served by Go's `http.FileServer`.

**Config storage.** `.editor/config.json`. Gitignored.

**Build staging.** `.editor/build-staging/`. Gitignored. Cleaned on every build.

## Adding a new endpoint

1. Add a handler method on `Server` in a file under `internal/handlers/`
2. Register the route in `main.go` under "API endpoints"
3. Add the corresponding fetch call in `ui/app.js`
4. Wire up the UI element (button, modal, etc.) in `index.html`

## Testing

There aren't unit tests yet. Manual testing is the protocol:

1. `go run .`
2. Walk through each operation in the browser
3. Verify the corresponding files in `mods/` change correctly
4. Verify the status pane logs the operation

If anything breaks, the Go terminal will surface the error. The browser DevTools console catches frontend issues.

## Future work

The biggest deferred item is **operations being repo-driven** instead of baked into the binary. Idea: a `.editor/operations.toml` file in the repo declares custom buttons and what they do. Lets the pack maintainer add pack-specific automations without rebuilding the binary. Not in v1 to keep scope tight.
