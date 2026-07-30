# divhead

Dividing head calculations.

**Converted from:** `DIVHEAD.C` (M. W. Klotz), `MWKC/WorkshopUtilities/divhead.zip`
**Go source:** `MWKGo/divhead/divhead.go`

## Purpose

A dividing head lets a milling machine cut equally spaced features
(gear teeth, bolt holes, flutes) around a workpiece, by cranking a
worm-driven spindle a controlled fraction of a turn between cuts.
Given how many equal divisions are wanted, this calculator works out
how many full crank turns are needed and, when that isn't a whole
number, either a rapid indexing plate shortcut (if your head has one
and it applies) or which index-plate hole circle and hole count
covers the remaining fraction of a turn.

## Data setup

A dividing head's worm-gear ratio, rapid indexing plate, and set of
available hole-circle plates are all specific to one owner's
equipment — like `diffthrd` earlier in this conversion group, this
program reads from the user's own database (`userdata.db`) rather
than the shared `reference.db`; see
`ai/plans/c-to-go-conversion-plan.md`, "Data-file strategy for
Tier 2". Both tables ship empty. Import your own dividing head's
settings and available hole-circle plates before first use:

```
divhead -import-settings my-divhead-settings.csv
divhead -import-holes my-divhead-holes.csv
```

The settings CSV needs three columns, `id,worm_gear_ratio,
rapid_index_holes` (one row, `id` always `1`; use `-1` for
`rapid_index_holes` if your head has no rapid indexing plate — the
original convention `DIVHEAD.C` itself uses). The holes CSV needs one
column, `holes`, one row per hole circle. Worked examples built from
the original `DIVHEAD.DAT`'s own default configuration (a 40:1 head
with a 24-hole rapid indexing plate and three standard Brown & Sharpe
style plates) ship at `MWKGo/divhead/testdata/example-settings.csv`
and `example-holes.csv`.

## Inputs

| Prompt | Default |
|---|---|
| Number of workpiece divisions | 14 |

## Output

The worm gear ratio, rapid indexing plate (if any), and the crank
turns required; then either a rapid-index shortcut, an exact
whole-turn count, or every hole-circle plate that can make up the
remaining fraction of a turn (a plate qualifies only if its hole
count is an exact multiple of the reduced remainder's denominator).

## Method

Turns required = worm gear ratio / divisions, split into whole turns
and a remainder fraction reduced to lowest terms via the greatest
common divisor. If the rapid indexing plate's hole count is an exact
multiple of the number of divisions, stepping `rapid holes /
divisions` holes on it is reported as the (only) solution — a
dedicated shortcut that skips the general hole-circle search
entirely, exactly as `DIVHEAD.C` does. Otherwise, if the division is
exact (no remainder), the whole-turn count alone is the answer. Failing
both of those, every available hole-circle plate whose hole count is
an exact multiple of the reduced remainder's denominator is reported,
each with how many holes to step on it — there can be more than one
valid plate, and all of them are listed as alternatives ("and" before
the first, "or" before the rest, matching the original's own
phrasing).

## Worked Example

`DIVHEAD.DAT`'s own default configuration (40:1 ratio, 24-hole rapid
index plate, plates of 15-20, 21-33, and 37-49 holes) is this
conversion's worked example, reproduced exactly in its tests: for 14
divisions, the rapid index plate does not apply (24 is not a multiple
of 14), so the search falls through to the reduced remainder 6/7,
which only two of the eighteen shipped plates (21 and 49 holes) are
multiples of — 18 holes on the 21-hole plate, or 42 holes on the
49-hole plate, both correct alternatives for the same 6/7-turn
remainder after 2 full crank turns.
