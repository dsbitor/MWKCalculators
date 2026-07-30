# links

Taper angles and shim height for a tapered, radiused-end link.

**Converted from:** `LINKS.C` (M. W. Klotz, 10/00),
`MWKC/WorkshopUtilities/links.zip`
**Go source:** `MWKGo/links/links.go`

## Purpose

Many engine and model builds need links: flat strips holding two
holes a fixed distance apart, each end radiused, with the sides
tapered to blend smoothly into the end radii. The usual shop
method mills each tapered side with the strip pinned between two
vise-mounted pins, shimming the small-end pin to set the
required taper angle. Given both end radii, both hole diameters,
and the distance between hole centers, this program computes the
taper angle, the included angle at each end, and the shim height
needed.

## Inputs

| Prompt | Default |
|---|---|
| Small end radius | 1/16 in |
| Small end hole diameter | 1/16 in |
| Big end radius | 3/32 in |
| Big end hole diameter | 3/16 in |
| Distance between hole centers | 1 in |

## Output

Angle of the tapered side, included angle at each end, and the
shim height needed under the small end pin.

## Method

```
halfAngle = asin((bigRadius - smallRadius) / centerDistance)
smallEndAngle = 180 - 2*halfAngle
bigEndAngle   = 180 + 2*halfAngle
shimHeight = (bigRadius - smallRadius) + 0.5*(bigHoleDiameter - smallHoleDiameter)
```

## Worked Example

No worked numeric example was included with the original
program. As an independently verifiable check: equal end radii
describe a parallel-sided, untapered link, so the taper angle
must be exactly zero and both end angles exactly 180 degrees,
regardless of hole center distance.
