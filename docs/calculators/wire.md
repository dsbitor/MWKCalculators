# wire

AWG wire gage recommendation and properties.

**Converted from:** `WIRE.C` (M. W. Klotz, 11/98),
`MWKC/WorkshopUtilities/wire.zip`
**Go source:** `MWKGo/wire/wire.go`

## Purpose

Given the current a wire must carry and the desired current
density (amps per circular mil), recommends the nearest standard
AWG wire gage and reports that gage's diameter, cross-sectional
area, resistance, and weight per 1000 feet.

## Inputs

| Prompt | Default |
|---|---|
| Current wire must carry | 12.0 amps |
| Desired current density | 0.0025 amps/cmil |

## Output

Recommended AWG gage, its diameter in mils, its area in circular
mils, its resistance in ohms per 1000 feet, and its weight in
pounds per 1000 feet.

## Method

AWG gage follows a standard geometric progression: each step
down in gage number multiplies the wire diameter by a fixed
ratio, `92^(1/39)`, so that 40 gage steps exactly multiply the
diameter by 92 (this program's constant, 1.12294049, is that
ratio to eight figures).

```
diameter_mils(gage) = 324.87 / ratio^gage
requiredDiameter = sqrt(current / density)
gage = round(inverse of diameter_mils at requiredDiameter)
```

## Worked Example

No worked example was included with the original program. As an
independently verifiable check: the program's diameter formula
matches the standard published AWG formula,
`diameter_mils = 5 * 92^((36-gage)/39)`, closely across a range
of gage numbers (the two constants are independent
approximations of the same standard, so they agree closely but
not to full floating-point precision).
