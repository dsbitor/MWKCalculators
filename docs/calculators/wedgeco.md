# wedgeco

Volume of a conical wedge.

**Converted from:** `WEDGECO.C` (M. W. Klotz), `MWKC/Math/wedgeco.zip`
**Go source:** `MWKGo/wedgeco/wedgeco.go`

## Purpose

If a cone is sliced by a plane parallel to its axis, offset from the
axis, the cone divides into two pieces; the smaller is a "conical
wedge." Its volume requires a triple integral with no simple closed
form. Given the cone's base diameter and height, and the wedge's
sagitta (the distance from the edge of the cone's base to the
slicing plane, measured along a diameter), this program computes the
wedge's volume, the complete cone's volume, and their ratio.

## Inputs

| Prompt | Default |
|---|---|
| Cone base diameter | 2.0 |
| Height of cone | 10.0 |
| Sagitta of wedge base | 0.5 |

Sagitta must be strictly between 0 and the cone's base diameter; the
program re-prompts otherwise.

## Output

Volume of the complete cone, volume of the conical wedge, and the
ratio of the two.

## Method

```
coneVolume = pi*r^2*h/3
if sagitta == r: wedge = 0.5*coneVolume      (plane through the axis)
else:
  offset = |r - sagitta|
  angle = acos(offset/r)
  wedge' = h*r^2*(angle - 2*sin(angle)*cos(angle)
                  + cos(angle)^3*ln(tan(angle)+1/cos(angle))) / 3
  wedge = wedge'                 if sagitta <= r
        = coneVolume - wedge'    if sagitta > r
```

## Worked Example

No worked numeric example was included with the original program
(`WEDGECO.TXT` explains the geometry but includes no sample run). The
formula's translation from the original was cross-checked
independently by numerical integration of the cone's circular
cross-sections against the cutting plane before being trusted for
this conversion's tests. As a further, independently verifiable
check: the wedge volumes cut at sagitta `s` and at the complementary
sagitta `2r - s` (the chord on the opposite side of the base circle)
must sum to exactly the complete cone's volume, confirmed in this
conversion's tests.
