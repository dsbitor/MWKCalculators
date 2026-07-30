# protrac

Sinebar-like protractor.

**Converted from:** `PROTRAC.C` (M. W. Klotz),
`MWKC/WorkshopUtilities/protrac.zip`
**Go source:** `MWKGo/protrac/protrac.go`
**Original documentation:** `PROTRAC.TXT`, inside `MWKC/WorkshopUtilities/protrac.zip` (not included in this conversion)

## Purpose

A sinebar gives an accurate angle but a plain protractor is more
convenient. This program supports a hybrid: two bars hinged together
with two precision pins mounted at the same fixed radius from the
hinge axis. Measuring the distance between the pins with a caliper
gives the angle between the bars with sinebar-like accuracy, the same
principle as a sinebar hinged in the middle with its rolls' spacing
adjusted. After calibrating with the pin separation at the fully
closed position (nonzero, since finite-diameter pins can't fully
close the gap), the program converts between measured pin separation
and included angle in either direction.

## Inputs

| Prompt | Default |
|---|---|
| Radius on which pins are mounted | 3 in |
| Pin diameter | 0.25 in |
| Pin separation when protractor is closed | 0.22 in |

Then, repeatedly until quit:

| Choice | Prompt | Default |
|---|---|---|
| A: find angle given separation | Pin separation | 1.0 in |
| D: find separation given angle | Angle | 15.0 deg |

## Output

The closed-position angle, then, for each menu choice, both the pin
separation and the resulting angle.

## Method

```
closedAngle = 2*asind(0.5*(pinDiameter+closedGap)/radius)
angle(separation) = 2*asind(0.5*(pinDiameter+separation)/radius) - closedAngle
separation(angle) = 2*radius*sind(0.5*(angle+closedAngle)) - pinDiameter
```

## Worked Example

No worked numeric example was included with the original program
(`PROTRAC.TXT` explains the hinge and pin arrangement with an ASCII
diagram, but no sample run). As an independently verifiable check: an
included angle of exactly zero must return the calibrated closed-gap
separation, and converting an angle to a separation and back must
reproduce the original angle; both confirmed in this conversion's
tests.
