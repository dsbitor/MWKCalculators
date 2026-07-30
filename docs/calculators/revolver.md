# revolver

Dimensions for a revolver-cylinder style tool holder.

**Converted from:** `REVOLVER.C` (M. W. Klotz, 8/00),
`MWKC/WorkshopUtilities/revolver.zip`
**Go source:** `MWKGo/revolver/revolver.go`

## Purpose

A rack of small tools (jeweler's screwdrivers, for example) can
be held in a rotating cylinder drilled with off-axis holes,
mounted on a thin shaft above a base, much like a revolver's
cylinder. Given the number of holes, the hole diameter, the
spacing wanted between adjacent hole edges, and the wall
thickness required beyond the holes, this program computes the
radius at which to place the holes and the stock diameter
needed for the cylinder.

## Inputs

| Prompt | Default |
|---|---|
| Number of holes | 6 |
| Diameter of holes | 0.25 in |
| Spacing between hole edges | 0.5 in |
| Thickness required at outer edge of holes | 0.25 in |

## Output

Radius for hole placement, and required stock diameter for the
cylinder.

## Method

```
anglePerHole = 360 / holeCount
holeRadius   = (edgeSpacing+holeDiameter) / (2*sin(anglePerHole/2))
cylinderDiameter = 2 * (holeRadius + 0.5*holeDiameter + wallThickness)
```

## Worked Example

No worked example was included with the original program. As an
independently verifiable check: at exactly six holes, adjacent
hole centers are 60 degrees apart, which is the hexagon
side-equals-radius identity, so the hole placement radius must
equal exactly `edgeSpacing + holeDiameter`, which the
conversion's tests confirm.
