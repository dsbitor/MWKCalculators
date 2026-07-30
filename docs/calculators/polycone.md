# polycone

Geometry of a "polycone" (pyramid-like) solid.

**Converted from:** `POLYCONE.C`, `MWKC/WorkshopUtilities/polycone.zip`
**Go source:** `MWKGo/polycone/polycone.go`

## Purpose

A "polycone" is a solid whose base is a regular polygon and whose
sides are triangular facets all meeting at a single apex above the
base's center (a tetrahedron is a three-sided polycone; a pyramid is
a four-sided one; stained-glass lampshades are often frustums of
polycones). Given the number of base sides, the length of one base
side, and the polycone's perpendicular height, this program computes
the base's dimensions, each facet's dimensions and angles, and the
solid's total surface area and volume.

## Inputs

| Prompt | Default |
|---|---|
| Number of polygon sides | 4 |
| Length of polygon side | 6 |
| Cone height | 10 |

## Output

Base circumscribed/inscribed circle diameter and area; facet height,
edge length, base angle, tip angle, face angle, edge angle, the
exterior and interior angles between adjacent facets, and facet area;
total surface area and volume.

## Method

```
halfChord = 0.5*sideLength
halfAngle = 180/sides
radius = halfChord / sin(halfAngle)          (circumscribed circle)
apothem = radius * cos(halfAngle)
baseArea = sides * apothem * halfChord

facetHeight = hypot(height, apothem)
edgeLength = hypot(height, radius)
baseAngle = acos(halfChord / edgeLength)
tipAngle = 2*(90 - baseAngle)
faceAngle = acos(apothem / facetHeight)
edgeAngle = acos(radius / edgeLength)
facetArea = 0.5 * facetHeight * halfChord

totalArea = baseArea + sides*facetArea
totalVolume = baseArea * height / 3
```

The exterior angle between two adjacent facets is found from the
angle between two vectors: one perpendicular to a facet's base edge
within that facet's plane, and the same vector rotated by the base's
full central angle to lie in the adjacent facet.

## Worked Example

No worked numeric example was included with the original program. As
independently verifiable checks: facet height and edge length are
each the hypotenuse of the polycone's own height and a base dimension
(the Pythagorean theorem, independent of this code's own separate
formula for each), the tip angle is always exactly twice the
complement of the base angle by construction, the face angle's
tangent must equal height divided by apothem (rise over run), and
total area is exactly the base area plus every facet's area; all
confirmed in this conversion's tests.
