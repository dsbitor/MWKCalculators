# ddh

Differential dividing head calculations.

**Converted from:** `DDH.C` (M. W. Klotz), `MWKC/WorkshopUtilities/ddh.zip`
**Go source:** `MWKGo/ddh/ddh.go`

## Purpose

A differential dividing head (DDH) is a dividing head augmented with
gearing between the spindle and hole plate, so the hole plate itself
turns a little as the spindle turns — a fraction of the spindle rate
set by the gear ratio used. That differential motion lets a DDH reach
divisions a plain dividing head's hole plates alone can't. This
computes, for a required number of divisions: first whether a plain
(non-differential) solution exists at all — the same calculation
`divhead` itself performs — and, if not, a differential solution: a
hole plate, an indexing increment near it, and a change gear train
between spindle and hole plate that makes the combination work.

## Data setup

A DDH's worm-gear ratio, rapid indexing plate, available hole-circle
plates, and available change gears are all specific to one owner's
equipment, so — like `diffthrd`, `divhead`, `gearatio`, and `change`
in earlier conversion groups — this program reads them from the
user's own database (`userdata.db`) rather than the shared
`reference.db`; see `ai/plans/c-to-go-conversion-plan.md`,
"Data-file strategy for Tier 2". These tables ship empty. Import all
three parts before first use:

```
ddh -import-settings my-ddh-settings.csv
ddh -import-holes my-ddh-holes.csv
ddh -import-gears my-ddh-gears.csv
```

The settings CSV needs `id` (always `1`), `worm_gear_ratio`, and
`rapid_index_holes` (0 or negative if there is no rapid indexing
plate). The holes CSV needs one column, `holes` (every hole circle
across every plate the DDH has, pooled together — `DDH.DAT`'s own
three plates, A, B, and C, are combined the same way). The gears CSV
needs one column, `teeth`. Worked examples built from the original
archive's own shipped `DDH.DAT` (a 40:1 worm gear ratio, a 24-hole
rapid indexing plate, 18 pooled hole circles across three plates, and
a 12-gear set) ship at `MWKGo/ddh/testdata/`.

## Inputs

| Prompt | Default |
|---|---|
| Number of workpiece divisions | 67 |
| Maximum gear pairs to examine (only if a differential search is needed) | 1 |
| Allowable gear ratio matching accuracy (%) (only if a differential search is needed) | 0.01 |

## Output

The worm gear ratio, whether a rapid indexing plate is available,
turns required, and then either: a rapid indexing plate step count; a
plain full-turns-plus-hole-plate solution (matching `divhead`'s own
output exactly); or, if neither exists, every differential solution
found — grouped by hole plate and indexing increment, each showing
whether the hole plate and crank turn the same way or opposite ways,
the gear ratio required, and every gear chain achieving it within
tolerance.

## Method

The non-differential stage is identical to `divhead`'s own
calculation; see `docs/calculators/divhead.md`, "Method" for the full
explanation of `turnsRequired` and the hole-plate search.

When no plain solution exists, the differential search tries, for
each available hole plate, five candidate indexing increments centered
on `floor(ratio × plate / divisions)` (`DDH.C`'s own search window of
±2 around this estimate) — a non-positive increment is skipped as not
physically meaningful. Each candidate implies a required change-gear
ratio; the gear-chain search finding trains for that ratio (no gear
position reused within one chain, ratios applied by multiplying or
dividing toward the target, every match within tolerance reported) is
identical to `gearatio`'s own search — see
`docs/calculators/gearatio.md`, "Method". `coding-style.md` Rule 2
replaces the original's interactive keypress abort with an explicit
search-evaluation cap shared across every hole plate and indexing
increment tried in one run.

## Worked Example

No worked numeric example with expected output is available. As
independently verifiable checks, this conversion's tests confirm
`turnsRequired` against `divhead`'s own worked example (40:1 ratio, 14
divisions gives 2 full turns plus 6/7); that `DDH.DAT`'s own default
of 67 divisions has no plain hole-plate solution (none of its 18 hole
circles is a multiple of 67); that a differential search for that same
67-division case does find a solution, and that every returned gear
chain's stages actually multiply out to its recorded achieved ratio;
and that `rapidIndexSolution` correctly accepts an exact multiple and
rejects both a missing plate and a non-multiple.
