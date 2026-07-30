# eccentub

Tube size for turning offset eccentrics in a 3-jaw chuck.

**Converted from:** `ECCENTUB.C` (M. W. Klotz, 02/01),
`MWKC/WorkshopUtilities/eccent.zip`
**Go source:** `MWKGo/eccentub/eccentub.go`

## Purpose

Turning an eccentric (a cylindrical part whose turned axis is
offset from its stock's original centerline) in a lathe's 3-jaw
chuck normally means packing one jaw with shim stock, which is
fiddly to measure and not very stable.

`ECCENTUB.C` implements a better method described in
`ECCENT.TXT`: bore a tube to a sliding fit on the parent stock,
mill a slot in the tube wide enough to pass one chuck jaw, and
clamp the assembly so one jaw seats on the parent stock through
the slot while the other two seat on the tube's outer surface.
This offsets the stock's centerline from the spindle axis by a
calculable amount. Given the parent stock diameter and the
required offset, this program computes the outer diameter the
tube itself needs.

The related program `ECCENT.C` (not yet converted) solves the
same problem using the older shim-packing method instead.

## Inputs

| Prompt | Default |
|---|---|
| Diameter of parent stock | 1.0 in |
| Required eccentric offset | 0.1 in |

## Output

Diameter of the tube required.

## Method

```
R     = 0.5 * parent stock diameter
r     = R - offset
dtube = 2 * sqrt(7*R^2 - 9*R*r + 3*r^2)
```

## Worked Example

No worked numeric example was included in `ECCENT.TXT`. As an
independently verifiable check: an offset of zero means no
eccentricity is being cut, so the formula must and does return
the parent stock diameter unchanged, for any stock diameter.
