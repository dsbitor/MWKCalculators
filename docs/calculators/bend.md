# bend

Sheet metal bend allowance.

**Converted from:** `BEND.C` (M. W. Klotz, 11/98),
`MWKC/WorkshopUtilities/bend.zip`
**Go source:** `MWKGo/bend/bend.go`

## Purpose

When sheet metal is bent, the material at the bend's neutral
axis neither stretches nor compresses, but that axis sits
somewhere between the bend's inner and outer radius depending
on how tight the bend is. This program computes the length of
flat material a bend of a given angle, radius, and material
thickness will consume, along with the lengths measured at the
bend's inner and outer surfaces.

## Inputs

| Prompt | Default |
|---|---|
| Thickness of material | 0.125 in |
| Radius of bend | 3.0 in |
| Angle of bend | 180.0 deg |

## Output

Length of the bend's exterior surface, its interior surface,
and the total length of material required to form the bend.

## Method

The neutral axis offset from the inner radius is looked up by
the ratio of bend radius to material thickness, not computed
from a continuous curve: 1/3 of the thickness for a tight bend
(radius under twice the thickness), 1/2 for a gentle bend
(radius over four times the thickness), and 0.4 of the
thickness in between.

```
exterior  = angle * (radius + thickness)
interior  = angle * radius
allowance = angle * (radius + offset)
```

## Worked Example

No worked example was included with the original program. With
the documented default inputs (thickness 0.125in, radius 3in,
in the gentle-bend range, angle 180 degrees), the program
reports an exterior length of 9.8175in, an interior length of
9.4248in, and a bend allowance of 9.6211in.
