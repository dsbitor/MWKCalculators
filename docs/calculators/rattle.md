# rattle

Diameter measurement of a large bore using Guy Lautard's
stick-and-rattle technique.

**Converted from:** `RATTLE.C` (M. W. Klotz, 11/99),
`MWKC/WorkshopUtilities/rattle.zip`
**Go source:** `MWKGo/rattle/rattle.go`

## Purpose

Measuring a bore too large to span with available calipers: cut a
stick slightly shorter than the bore diameter, with rounded or
pointed ends, and insert it. The stick "rattles" back and forth by a
small amount. Given the stick length and the peak-to-peak rattle
distance, this program computes the actual bore diameter, sparing
the trial-and-error of directly setting internal calipers to the
right spread.

## Inputs

| Prompt | Default |
|---|---|
| Measured stick/caliper distance | 4.0 |
| Rattle distance | 0.2 |

## Output

Diameter, and the difference between the diameter and the stick
length.

## Method

```
theta = asind(0.5*rattle/stick)
beta  = asind(rattle/stick)
diameter = stick*cosd(theta) / (1 - 0.5*(1 - cosd(beta)))
```

## Worked Example

No fully worked numeric example was included with the original
program (`RATTLE.TXT` explains the technique and its origin in Guy
Lautard's "Home Machinist's Bedside Reader #1" but includes no
sample run). As an independently verifiable check: a rattle distance
of zero means the stick already fits the bore exactly, so the
diameter must equal the stick length exactly, confirmed in this
conversion's tests.
