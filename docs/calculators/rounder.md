# rounder

Ball end mill rounding-over table.

**Converted from:** `ROUNDER.C` (M. W. Klotz),
`MWKC/WorkshopUtilities/rounder.zip`
**Go source:** `MWKGo/rounder/rounder.go`

## Purpose

Rounding over an edge with a fixed radius, with no lathe access and
no through-hole to pivot against, can be roughed out on the milling
machine with a ball end mill: successive scallop cuts, each tangent
to the desired radius profile, are made by stepping through an angle
theta and positioning the mill's center at the corresponding
coordinates. Given the desired workpiece radius, the ball mill
diameter, and an angular step size, this program tabulates those
coordinates (and two related derived values) from 0 to 90 degrees.

## Inputs

| Prompt | Default |
|---|---|
| Workpiece radius | 1.0 in |
| Ball-end mill diameter | 0.25 in |
| Angular increment | 5.0 deg |

## Output

A table, for each angle from 0 to 90 degrees in steps of the
increment, of:

```
A = (R+r)*sin(theta)
B = (R+r)*cos(theta)
C = (R+r)*(1 - sin(theta))
D = (R+r)*(1 - cos(theta))
```

where R is the workpiece radius and r is the ball mill radius.

## Worked Example

No worked numeric example is available. As independently
verifiable checks: at theta=0 and theta=90 the values reduce to the
obvious endpoints of a
quarter turn, `(A,B)` always lies on a circle of radius `R+r`
centered on the origin, and `C`/`D` are always exactly `(R+r)-A` and
`(R+r)-B` respectively; all confirmed in this conversion's tests.
