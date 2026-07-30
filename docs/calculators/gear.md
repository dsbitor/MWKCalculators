# gear

Spur gear dimensions.

**Converted from:** `GEAR.C`, `MWKC/WorkshopUtilities/gear.zip`.
Reference: John A. Cooper, "Spur Gears and Pinions", Machinist's
Workshop, 4/99, Vol. 12, No. 2, pg. 8.
**Go source:** `MWKGo/gear/gear.go`
**Original documentation:** `GEAR.TXT`, inside `MWKC/WorkshopUtilities/gear.zip` (not included in this conversion)

## Purpose

Computes standard 20-degree full-depth involute spur gear dimensions
for a mating pair of gears, given each gear's tooth count and the
pair's shared diametral pitch and pressure angle, per John Cooper's
"Spur Gears and Pinions" article. The documented default inputs match
Cooper's own worked example in that article.

## Inputs

| Prompt | Default |
|---|---|
| Number of teeth on gear 1 | 45 |
| Number of teeth on gear 2 | 20 |
| Diametral Pitch (25.4/mod) | 20 |
| Pressure Angle | 20 deg |

## Output

Gear ratio, diametral pitch, pressure angle, and center distance for
the pair, then, for each gear individually: number of teeth, outside
diameter, addendum, dedendum, whole depth, circular pitch, tooth
thickness, pitch diameter, base circle radius, and tooth profile
radius (each in both inches and millimeters).

## Method

```
gearRatio = max(teeth1, teeth2) / min(teeth1, teeth2)
per gear:
  pitchDiameter = teeth / diametralPitch
  outsideDiameter = (teeth + 2) / diametralPitch
  addendum = 1 / diametralPitch
  dedendum = 1.200 / diametralPitch
  wholeDepth = 2.200 / diametralPitch
  circularPitch = pi / diametralPitch
  toothThickness = 0.48 * circularPitch
  baseCircleRadius = 0.5*pitchDiameter*cosd(pressureAngle)
  toothProfileRadius = 0.5*pitchDiameter*sind(pressureAngle)
centerDistance = 0.5*(pitchDiameter1 + pitchDiameter2)
```

The original program writes its results to `GEAR.OUT` via `fopen`;
this conversion prints to stdout instead and drops that DOS
file-save-then-page convenience, the same approach used for `loan`
(Tier 1 suitability review, Finding 5).

## Worked Example

No numeric worked example was included as a `.TXT` file, though
`GEAR.TXT` confirms the documented defaults are Cooper's own example
values. As an independently verifiable check: the base circle radius
is, by definition, the pitch radius times the cosine of the pressure
angle, and the center distance is exactly half the sum of the two
gears' pitch diameters (their pitch circles are tangent where they
mesh); both confirmed in this conversion's tests as identities
independent of this code's own separate formulas for those
quantities.
