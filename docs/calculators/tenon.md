# tenon

Depth of cut for a regular polygonal tenon.

**Converted from:** `TENON.C` (M. W. Klotz, 6/00),
`MWKC/WorkshopUtilities/tenon.zip`
**Go source:** `MWKGo/tenon/tenon.go`

## Purpose

Cutting a regular polygonal tenon (most often square or
hexagonal) on the end of cylindrical stock, using a rotary
indexer on the milling machine, needs the depth of cut below the
stock's surface for each flat. An even-sided tenon is easily
specified by its across-flats dimension, but an odd-sided one
(the author's example: a pentagonal tenon to serve as an "Allen
key" for tamper-proof bolts) is better specified by the diameter
of the circle that just circumscribes it. This program computes
the depth of cut either way, along with the angle to rotate
between cuts.

## Inputs

| Prompt | Default |
|---|---|
| Diameter of stock | 0.5 |
| Number of sides on tenon | 5 |
| Distance across flats (even sides only) | 0.25 |
| Diameter of circle circumscribing tenon (odd sides only) | 0.25 |

Only one of the last two prompts appears, depending on whether
the number of sides is even or odd.

## Output

Stock diameter, number of flats, angle between adjacent flats,
diameter of the circumscribing circle, and depth of cut.

## Method

```
anglePerSide = 360 / sides
halfAngle = anglePerSide / 2

# even sides: derive the circumscribed diameter from across flats
circumscribedDiameter = acrossFlats / cos(halfAngle)

# odd sides: circumscribed diameter is given directly

depthOfCut = 0.5*stockDiameter - 0.5*circumscribedDiameter*cos(halfAngle)
```

## Worked Example

No worked example was included with the original program. As an
independently verifiable check: the even-sides formula matches
standard across-flats-to-across-corners identities for regular
polygons (a square's across-corners is across-flats times
`sqrt(2)`; a hexagon's is across-flats times `2/sqrt(3)`), both
well known and independent of this code.
