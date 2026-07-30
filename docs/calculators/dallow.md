# dallow

Drill tip allowance.

**Converted from:** `DALLOW.C` (M. W. Klotz, 5/99),
`MWKC/WorkshopUtilities/dallow.zip`
**Go source:** `MWKGo/dallow/dallow.go`

## Purpose

A twist drill's conical tip means a hole must be drilled deeper
than its intended flat-bottomed depth to reach full diameter at
that depth. This program computes that extra depth, the "tip
allowance", from the drill's diameter and its included tip
angle.

## Inputs

| Prompt | Default |
|---|---|
| Included angle of drill tip | 118 deg |
| Drill diameter | 0.5 in |

## Output

Tip allowance, in the same units as the diameter.

## Method

```
allowance = 0.5 * diameter / tan(angle / 2)
```

## Worked Example

No worked example was included with the original program. As an
independently verifiable check: at a 90 degree included angle,
the tip forms a perfect right-angle cone, and the allowance
equals exactly the drill's radius, for any diameter. The
documented default (118 degrees, the standard twist drill point
angle, 0.5in diameter) gives an allowance of 0.1502in.
