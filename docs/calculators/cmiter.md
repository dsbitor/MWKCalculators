# cmiter

Compound table saw angles for mitred polygonal forms.

**Converted from:** `CMITER.C` (M. W. Klotz, 3/02),
`MWKC/WorkshopUtilities/cmiter.zip`. Written for Jeff Jarnberg.
Reference: http://www.betterwoodworking.com/compound_miter.htm
**Go source:** `MWKGo/cmiter/cmiter.go`

## Purpose

Cutting the sides of a sloped polygonal form (a tapered box or
lampshade, for example) on a table saw needs both a miter gauge
angle and a blade tilt, set together. Given the number of sides
and the slope of the sides, this program computes the miter
gauge angle and the blade tilt for both a mitred and a butted
joint style.

## Inputs

| Prompt | Default |
|---|---|
| Number of sides | 4 |
| Slope | 30 deg |

## Output

Miter gauge angle, and blade tilt for both mitred and butted
joints.

## Method

Three cases, by slope:

```
slope == 90: miterGauge = 90, bladeTiltMitred = 180/sides, bladeTiltButted = 0
slope == 0:  miterGauge = atan(1 / tan(180/sides)), bladeTiltMitred = 0, bladeTiltButted = 90
otherwise:
  miterGauge = atan(1 / (cos(slope) * tan(180/sides)))
  bladeTiltMitred = atan(cos(miterGauge) * tan(slope))
  bladeTiltButted = atan(cos(miterGauge) / tan(slope))
```

## Worked Example

No worked example was included with the original program. As an
independently verifiable check: at a zero degree slope, the
miter gauge angle reduces to `atan(cot(180/sides))`, which
equals `90 - 180/sides` exactly, a standard trigonometric
identity independent of this code.
