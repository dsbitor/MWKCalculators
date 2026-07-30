# collet

Bore diameter for a slotted collet holding polygonal stock.

**Converted from:** `COLLET.C` (M. W. Klotz, 5/04),
`MWKC/WorkshopUtilities/collet.zip`
**Go source:** `MWKGo/collet/collet.go`

## Purpose

John Way's technique (Machinist's Workshop, June/July 2004,
Vol. 17, No. 3) for holding square or hexagonal stock without a
dedicated broach: bore a plain cylindrical hole in the collet,
then saw one slot per side of the stock. The corners of the
stock locate in the slots while the cylindrical bore clamps down
on the stock's flat faces. Given the number of sides, the
stock's across-flats dimension, and the slot width, this program
computes the required bore diameter. Valid for stock with an
even number of sides; the author suggests 4, 6, or 8 as the
practical range.

## Inputs

| Prompt | Default |
|---|---|
| Number of stock polygon sides | 6 |
| Stock across flats dimension | 3/16 in |
| Collet slot width | 0.045 in |

## Output

Required collet bore diameter.

## Method

```
halfAngle = 180 / sides
offset = 0.5*acrossFlats/cos(halfAngle) - 0.5*slotWidth*sin(halfAngle)/cos(halfAngle)
boreDiameter = 2 * sqrt(0.25*slotWidth^2 + offset^2)
```

## Worked Example

No worked example was included with the original program. As an
independently verifiable check: at a slot width of zero, the
formula reduces to `acrossFlats / cos(180/sides)`, the standard
across-flats-to-across-corners conversion for a regular polygon.
For a hexagon that ratio is `2/sqrt(3)`, and for a square it is
`sqrt(2)`, both well known identities independent of this code.
