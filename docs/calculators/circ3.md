# circ3

Radius of a circle passing through three points, given only
the distances between them.

**Converted from:** `CIRC3.C` (M. W. Klotz, 7/99),
`MWKC/Math/circ3.zip`
**Go source:** `MWKGo/circ3/circ3.go`

## Purpose

Finds the radius of the unique circle that passes through
three points, when the three points can't be measured
directly by coordinates but the distances between each pair of
points can be measured (for example, three points marked on a
broken gear fragment, where the original diameter needs to be
recovered).

Number the three points 1 to 3 in any order, measure the
distance between each pair, and enter the three distances. The
program reports the radius of the circle passing through all
three, provided they aren't collinear.

To find the physical center once the radius is known: set a
compass to the reported radius and draw a circle around each
of the three points. The point where the three circles
intersect is the center.

## Related, Not Yet Converted

The same zip file also contains `CIRC3C.C`, a variant that
takes the three points as Cartesian coordinates instead of
pairwise distances, and reports both the center coordinates and
the radius. It has not been converted yet; when it is, it
belongs in its own `MWKGo/circ3c/` directory per the project's
one-file-one-program convention, not folded into this one.

## Inputs

| Prompt | Default |
|---|---|
| length 1-2 | 2 |
| length 1-3 | 1.1 |
| length 2-3 | 1.1 |

If the three distances don't satisfy the triangle inequality,
the program reports that there is no solution and prompts
again.

## Output

Radius and diameter of the circle passing through all three
points.

## Method

Places point 1 at the origin and point 2 at distance d12 along
the x-axis, finds point 3 as the intersection of a circle of
radius d13 around point 1 and a circle of radius d23 around
point 2, then finds the circumcenter as the intersection of the
perpendicular bisectors of two of the triangle's sides.

## Worked Example

The original `.TXT` file gives no worked numeric example. As an
independently verifiable check: distances of 3, 4, and 5 (a
right triangle) give a circumradius of exactly 2.5, half the
hypotenuse, regardless of which distance is entered as d12,
d13, or d23.
