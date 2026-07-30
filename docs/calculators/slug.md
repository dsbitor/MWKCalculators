# slug

Chain-drilling layout for opening a large hole in plate stock.

**Converted from:** `SLUG.C` (M. W. Klotz, 10/00),
`MWKC/WorkshopUtilities/slug.zip`
**Go source:** `MWKGo/slug/slug.go`
**Original documentation:** `SLUG.TXT`, inside `MWKC/WorkshopUtilities/slug.zip` (not included in this conversion)

## Purpose

Drilling is the most efficient mechanical way to remove metal
without commercial cutting equipment. Chain drilling a circle of
holes around the perimeter of a large hole, then freeing the
resulting "slug", is a practical way to open that hole. Given
the hole's final diameter, a radial allowance left for finish
machining, the drill diameter, and an approximate web thickness
wanted between holes, this program computes how many holes to
drill, on what diameter, and the actual web thickness that
results.

This program shares its layout mathematics with `plate`
(`PLATE.C`, in the same zip file), which solves the mirror-image
problem of cutting a plate free rather than opening a hole; both
are implemented here on top of a shared `internal/chaindrill`
package.

## Inputs

| Prompt | Default |
|---|---|
| Diameter of final hole | 3 in |
| Radial allowance for finish machining | 0.05 in |
| Drill diameter | 0.25 in |
| Approximate web thickness | 0.05 in |

## Output

Number of holes, diameter of the drilling circle, the resulting
web thickness, the angle between adjacent holes, and the
chordal distance between them.

## Method

See `internal/chaindrill`. `slug` shrinks the drilling circle
inward from the final diameter (`outward = false`), since the
as-drilled hole must be smaller than the finished one: finish
machining enlarges it out to size.

## Worked Example

`SLUG.TXT` gives two full worked examples at the documented
default inputs, varying only the drill diameter:

| Drill diameter | Holes | Drilling circle | Web thickness | Angle | Chord |
|---|---|---|---|---|---|
| 0.25 in | 27 | 2.650 in | 0.058 in | 13.333 deg | 0.308 in |
| 0.375 in | 18 | 2.525 in | 0.064 in | 20.000 deg | 0.438 in |

This conversion reproduces both exactly.
