# qbelt

Quick two-pulley belt length calculation.

**Converted from:** `QBELT.C` (M. W. Klotz), `MWKC/WorkshopUtilities/belt.zip`
**Go source:** `MWKGo/qbelt/qbelt.go`

## Purpose

A faster, purely interactive alternative to [belt](belt.md) for the
common case of just two pulleys: given both diameters and the
separation between their centers, computes each pulley's wrap angle
and wrap length, the straight belt span between them, and the total
belt length. No data file is needed.

## Inputs

| Prompt | Default |
|---|---|
| Diameter of smaller pulley | 2 |
| Diameter of larger pulley | 6 |
| Separation between pulley centers | 6 |

## Output

For each pulley (labeled smaller/larger, swapped automatically if
entered the other way around): its diameter, wrap angle (degrees),
and wrap length. Then the belt span between the pulleys and the total
belt length.

## Method

`twoPulleyWrap` is the standard external-tangent two-pulley belt
formula: half the radius difference's arcsine gives the angle each
pulley's wrap deviates from a half-circle (the smaller pulley wraps
slightly less than half, the larger slightly more — the two always
sum to a full circle), and the straight span is the same on both
sides by symmetry. `pulley` and `pcd` (the archive's other two
companion programs) share this identical formula, each solving for a
different unknown; see `docs/calculators/pulley.md` and
`docs/calculators/pcd.md`.

## Worked Example

No worked numeric example was included with the original program.
As independently verifiable checks, this conversion's tests confirm
the two wrap angles always sum to exactly 2π regardless of pulley
size; that two equal-diameter pulleys each wrap exactly half the
circle; that the reported total is exactly the sum of its own parts
(spans and wraps); and, using `BELT.DAT`'s own "conical pulley"
example pulley pair (diameters 1.4 and 0.603 at a 2.5 separation),
that the resulting belt length comes out close to the value that
example's own doc comment implies (8.21) — the same fixed length
[pulley](pulley.md) and [pcd](pcd.md) both search for by default.
