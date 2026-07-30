# chord

Chord length for stepping off a circle into equal divisions.

**Converted from:** `CHORD.C` (M. W. Klotz, 2/00),
`MWKC/WorkshopUtilities/chord.zip`
**Go source:** `MWKGo/chord/chord.go`

## Purpose

Given the diameter of a circle and the number of equal
divisions wanted around its circumference, computes the
straight-line chord length between two adjacent division
points. Useful for stepping off equally spaced points on a
circle (bolt circles, index marks) with dividers rather than a
protractor.

## Inputs

| Prompt | Default |
|---|---|
| Number of divisions | 5 |
| Diameter of circle | 1 in |

## Output

Chord length, in the same units as the diameter.

## Method

```
angle  = 360 / divisions
chord  = diameter * sin(angle / 2)
```

## Worked Example

No worked example was included with the original program. The
default inputs (5 divisions, diameter 1) give a chord length of
0.5878, the side length of a regular pentagon inscribed in a
unit-diameter circle.
