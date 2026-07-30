# change

Lathe change gear calculations.

**Converted from:** `CHANGE.C` (M. W. Klotz), `MWKC/WorkshopUtilities/change.zip`
**Go source:** `MWKGo/change/change.go`
**Original documentation:** `CHANGE.TXT`, inside `MWKC/WorkshopUtilities/change.zip` (not included in this conversion)

## Purpose

Finds which of a lathe's own change gears, chained together (up to a
chosen chain length), drive the leadscrew at the ratio needed to cut a
required thread pitch — the classic lathe threading calculation, and
the direct ancestor of `gearatio` (an earlier conversion group's
stand-alone gear-ratio finder, adapted from this program's own
search).

## Data setup

A lathe's leadscrew's effective pitch(es) and its own set of available
change gears are specific to one owner's equipment, so — like
`diffthrd`, `divhead`, and `gearatio` in earlier conversion groups —
this program reads them from the user's own database (`userdata.db`)
rather than the shared `reference.db`; see
`ai/plans/c-to-go-conversion-plan.md`, "Data-file strategy for
Tier 2". These tables ship empty. Import all three parts before first
use:

```
change -import-settings my-change-settings.csv
change -import-pitches my-change-pitches.csv
change -import-gears my-change-gears.csv
```

The settings CSV needs `id` (always `1`) and `max_pairs` (the largest
gear chain length to search, 1 to 9 — `CHANGE.C`'s own fixed bound).
The pitches CSV needs one column, `pitch_tpi`: most lathes have one
leadscrew pitch, but a lathe with both a quick-change gearbox and a
change-gear banjo can offer several "effective leadscrew pitches", and
a different one may let a required ratio be reached with fewer gear
pairs. The gears CSV needs one column, `teeth`. Worked examples built
from the original archive's own shipped `CHANGE.DAT` (a single 8 tpi
leadscrew pitch, a 17-gear set from 20 to 100 teeth, and a suggested
starting chain length of 2) ship at `MWKGo/change/testdata/`.

## Inputs

| Prompt | Default |
|---|---|
| Type of thread to cut - Imperial or Metric | Imperial |
| Thread pitch to cut (tpi, if Imperial) | 21 |
| Thread pitch to cut (mm, if Metric) | 1.0 |
| Allowable thread pitch accuracy (%) | 0.01 |

## Output

For each configured leadscrew pitch: the required gear ratio, then
every gear chain found within tolerance — the gears used (`teeth:teeth`
per stage, joined by `-` for chains longer than one stage), the
chain's actual combined ratio, its error from the target, and `**`
marking the lowest-error chain found for that leadscrew pitch. If more
than one leadscrew pitch is configured, each solution line is also
prefixed with which pitch it belongs to. If none is found within the
examined chain lengths, a "NO SOLUTION FOUND" message.

## Method

The gear-chain search — no gear position reused within one chain,
each stage's ratio applied by multiplying or dividing the running
total toward the target and swapping the displayed pair order to
match, every match within tolerance reported rather than only the
first — is identical to `gearatio`'s own search; see
`docs/calculators/gearatio.md`, "Method" for the full explanation.
`change` adds an outer loop over each configured leadscrew pitch (each
one implying its own required ratio, `leadscrew pitch / desired
pitch`), and reports the single lowest-error solution found for each
pitch with `**`, matching `CHANGE.C`'s own per-pitch "best" marker
exactly. `coding-style.md` Rule 2 replaces the original's interactive
keypress abort with an explicit search-evaluation cap shared across
every leadscrew pitch tried in one run.

`CHANGE.TXT` explains that this program deliberately eliminates the
duplicate output its own predecessor, `CHANGEX.C` (bundled in the same
archive), produced for gear trains that are geometrically different
arrangements of the same gears but mathematically identical ratios
(`A:B-C:D` and `C:D-A:B` both give `R = AC/BD`) — `CHANGEX` is
deferred to a future clean-up group alongside `combi`, `ratio`, and
`belt`, since `CHANGE` supersedes it for every case that doesn't
specifically need to see those redundant variants.

## Worked Example

No worked numeric example with expected output was included with the
original program. As independently verifiable checks, this
conversion's tests confirm an exact single-pair ratio present in the
shipped gear set (20:40 = 0.5) is found; that, using `CHANGE.DAT`'s
own 8 tpi leadscrew pitch and suggested starting chain length of two,
a common 20 tpi thread (a real, practical scenario) is solvable within
0.5% tolerance; that solutions for two different leadscrew pitches are
each tagged with, and their recorded ratio and error match, their own
pitch; and that an unreachable target correctly reports no solution.
