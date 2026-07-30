# polygon

Properties of regular polygons.

**Converted from:** `POLYGON.C`, `MWKC/Math/polygon.zip`
**Go source:** `MWKGo/polygon/polygon.go`

## Purpose

Given the number of sides, and any one of four possible size
specifications (side length; distance across flats, for an even side
count; distance from a flat to the opposite vertex, for an odd side
count; or the diameter of the circumscribed or inscribed circle),
this program computes every other property of the regular polygon:
angles, both circle diameters, side length, perimeter, and area.
Units don't matter; output uses whatever unit was used for the one
size specification given.

## Inputs

| Prompt | Default |
|---|---|
| Number of polygon sides | 6 |
| Length of side | (skip if unknown) |
| Size across flats (even sides) | (skip if unknown) |
| Size flat-to-opposite-vertex (odd sides) | (skip if unknown) |
| Diameter of circumscribed circle | (skip if unknown) |
| Diameter of inscribed circle | (skip if unknown) |

At least one of the five size prompts must be answered.

## Output

Central angle, vertex angle, circumscribed and inscribed circle
diameters, distance across flats or flat-to-vertex (whichever applies
to the side count's parity), side length, perimeter, and area.

## Method

All derived from the circumradius (half the circumscribed circle
diameter):

```
centralAngle = 360 / sides
apothem = circumradius * cos(centralAngle/2)
sideLength = 2 * circumradius * sin(centralAngle/2)
area = 0.5 * sideLength * apothem * sides
```

Each of the four alternate size specifications is converted to a
circumradius by inverting this same relationship before the above is
applied.

## Worked Example

No worked numeric example was included with the original program. As
independently verifiable checks: a regular hexagon's side length
equals its circumradius exactly, and a square's distance across
flats equals its side length exactly; both well known geometric
properties, confirmed in this conversion's tests, along with a
round-trip check that each alternate-input formula exactly inverts
the corresponding output field.
