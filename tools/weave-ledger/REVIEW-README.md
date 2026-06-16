# Weave-review generator — how to update

`docs/WEAVE-REVIEW.md` integrates the mod-review sweep with the weave-ledger candidate corpus. It is
**generated and re-runnable** — the corpus keeps growing (targeting ~50 convergence passes), so refresh
the doc instead of hand-editing it.

## Update after a new pass

1. Land the new pass on the **`weaving-plan`** branch so its `CANDIDATES.md` is regenerated (the
   corpus is **not** vendored on this branch — we reference it).
2. Run:
   ```
   git fetch origin weaving-plan        # get the latest corpus
   python3 tools/weave-ledger/gen-weave-review.py   # reads it via `git show origin/weaving-plan:...`
   ```
   (or `--candidates <path>` to read a local checkout, `--ref <branch>` for a different ref.)
3. Read the run summary. It prints `unmapped` and `missing-DX` counts and lists them in a **Needs
   triage** banner at the top of the doc:
   - **unmapped** → add the modid to `review-map.json` as `{"slug": "...", "milestone": N}`.
   - **missing DX** → add the modid to `review-dx.json` with your take.
4. Re-run until both are 0, commit the regenerated doc + the data files.

New mods that appear at higher passes can't slip through silently — they surface in the triage banner and
in an "Unmapped — needs triage" section rather than being dropped.

## The files (edit these, never the generated doc)

| File | What it is |
|---|---|
| `gen-weave-review.py` | the generator (`--threshold N` min passes to show, default 2; `--per-mod K` rows/mod) |
| `review-map.json` | `modid → {slug, milestone}` — the schedule mapping (which milestone reviews the mod) |
| `review-dx.json` | `modid → DX note` — the per-mod take; **persists across passes** |
| `review-preamble.md` | the static intro (thunderdome framing, synthesis); prepended verbatim |
| `CANDIDATES.md` *(corpus, input)* | **not vendored** — read from the `weaving-plan` branch via `git show` |

`modid` keys are the ledger's mod ids (as they appear in `CANDIDATES.md`), which include version suffixes
for some mods (e.g. `trek-b0.6.1.1`) — match them exactly.

## Threshold

Default shows candidates seen in **≥2 passes**. As passes accumulate, consider raising it
(`--threshold 5`) to focus on the Gate-2 core (≥5 passes), which is the high-confidence authoring set.
