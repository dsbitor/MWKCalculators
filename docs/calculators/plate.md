# plate

Chain-drilling layout for cutting a circular plate free from flat
stock.

**Converted from:** `PLATE.C` (M. W. Klotz, 6/03),
`MWKC/WorkshopUtilities/slug.zip`
**Go source:** `MWKGo/plate/plate.go`

## Purpose

Chain drilling a circle of holes and cutting out the resulting
scalloped disc is a practical way to cut a circular plate from
sheet stock without commercial cutting equipment. Given the
plate's final diameter, a radial allowance left for finish
machining, the drill diameter, and an approximate web thickness
wanted between holes, this program computes how many holes to
drill, on what diameter, and the actual web thickness that
results (since the hole count must be a whole number, the
requested web thickness is only approximate).

This program shares its layout mathematics with `slug`
(`SLUG.C`, in the same zip file), which solves the mirror-image
problem of opening a large hole rather than cutting a plate
free; both are implemented here on top of a shared
`internal/chaindrill` package. `SLUG.TXT`, written for the
`slug` side of the problem, is this program's only surviving
documentation and is carried forward on the `slug` page instead
of duplicated here.

## Inputs

| Prompt | Default |
|---|---|
| Diameter of final plate | 3 in |
| Radial allowance for finish machining | 0.05 in |
| Drill diameter | 0.25 in |
| Approximate web thickness | 0.05 in |

## Output

Number of holes, diameter of the drilling circle, the resulting
web thickness, the angle between adjacent holes, and the
chordal distance between them (useful for laying out the holes
with dividers rather than coordinates).

## Method

See `internal/chaindrill`. `plate` grows the drilling circle
outward past the final diameter (`outward = true`), since the
as-drilled piece must be larger than the finished plate: finish
machining removes the scalloped edge down to size.

## Worked Example

No worked example specific to `plate` was included with the
original program (`SLUG.TXT` documents `slug`'s worked examples
only). As an independently verifiable check: growing the
drilling circle outward and then shrinking the same amount back
inward returns to the original final diameter, confirmed in
`internal/chaindrill`'s own tests.
