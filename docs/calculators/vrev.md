# vrev

Calculations for various sized holes arranged on a bolt circle.

**Converted from:** `VREV.C` (M. W. Klotz), `MWKC/WorkshopUtilities/vrev.zip`
**Go source:** `MWKGo/vrev/vrev.go`

## Purpose

Lays out holes of *varying* diameter evenly around a bolt circle
given a desired edge-to-edge spacing — the tap-guide and
multi-diameter punch-holder problem, where a fixed hole-to-hole gap
(rather than a fixed angular spacing) is wanted despite the holes
being different sizes. Reports the minimum bolt circle diameter that
fits every hole with at least that spacing, then, for a chosen
(possibly larger) bolt circle diameter, the exact position and angle
of each hole.

## Data setup

The hole list is specific to one job (the particular set of holes
being laid out this time), not universal reference data or reusable
equipment configuration, so — like `calibrat` in the same conversion
group — this program reads its input fresh from a file named on the
command line each run, in the same `STARTOFDATA`/`ENDOFDATA` text
format the original program used; see
`ai/plans/c-to-go-conversion-plan.md`, "Data-file strategy for
Tier 2".

```
vrev -holes my-punch-holder.dat
```

The file needs one hole diameter per line, in the order the holes
will be drilled. A worked example built from `VREV.TXT`'s own pin
punch holder example ships at `MWKGo/vrev/testdata/example.dat`.

## Inputs

| Prompt | Default |
|---|---|
| Desired hole-to-hole spacing | 0.1 |
| Desired bolt circle diameter | the computed minimum |

## Output

The minimum and maximum hole sizes; the minimum bolt circle diameter;
the actual hole-to-hole spacing at the chosen bolt circle diameter
(larger than requested if a bolt circle bigger than the minimum was
chosen); each hole's index, diameter, (x,y) position, and angle; and
the theoretical minimum and recommended stock diameters to cut the
part from.

## Method

For a bolt circle of radius R, each hole of diameter d subtends an
angle of `2*asin(0.5*d/R)`, and each hole-to-hole gap subtends
`spacing/R` (in radians, converted to degrees). The minimum bolt
circle radius is the one where the sum of every hole's angle plus
every gap's angle totals exactly 360°, found by the same zoom-in grid
search technique used by `cseg` and `flywheel` in Tier 1 (repeatedly
scanning ten steps across a window and zooming in on the best point,
bounded to 200 refinement passes per `coding-style.md` Rule 2).
Arcsine arguments are clamped to [-1,1] via `internal/mwktrig`
(shared with several Tier 1 conversions), matching the original's own
`ASND` macro.

Once a bolt circle diameter is chosen, the actual resulting spacing
is recomputed from the total angle the holes themselves occupy (not
including gaps) subtracted from 360°, divided evenly among the gaps;
holes are then placed around the circle starting from angle zero,
walking forward by each hole's own half-width, the gap, and the next
hole's half-width.

## Worked Example

`VREV.TXT`'s own pin-punch-holder example (eight holes: two each of
0.265625, 0.328125, 0.40625 in, plus one 0.484375 in) is this
conversion's worked example, reproduced in its tests to within the
four decimal places the original documentation shows: a minimum bolt
circle diameter of 1.1915, an actual spacing of 0.1236 at a chosen
diameter of 1.25, the full eight-hole position table, and theoretical
(1.7344) and recommended (1.9816) stock diameters.
