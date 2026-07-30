# boltcirc

Bolt circle hole layout.

**Converted from:** `BOLTCIRC.C` (M. W. Klotz, 11/98, 1/00),
`MWKC/WorkshopUtilities/boltcirc.zip`
**Go source:** `MWKGo/boltcirc/boltcirc.go`

## Purpose

Given a bolt circle's hole count, radius, hole diameter, angular
offset for the first hole, and center offset, computes the
edge-to-edge spacing between adjacent holes (warning if they would
overlap) and the angular and Cartesian coordinates of every hole, for
laying out or programming the pattern.

No `.TXT` file was included with the original program; this purpose
statement is drawn from the `.C` file's own header comment.

## Inputs

| Prompt | Default |
|---|---|
| Number of holes | 5 |
| Radius of bolt circle | 1.0 |
| Diameter of bolt holes | 0.5 |
| Angular offset of first hole | 0.0 deg |
| X offset of bolt circle center | 0.0 |
| Y offset of bolt circle center | 0.0 |

## Output

Spacing between hole edges (with an overlap warning if negative),
then each hole's angle and X/Y coordinates.

## Method

```
step = 360 / numHoles
spacing = 2*radius*sind(0.5*step) - holeDiameter
angle_i = angularOffset + i*step
x_i = radius*cosd(angle_i) + xOffset
y_i = radius*sind(angle_i) + yOffset
```

The original program writes its results to `BOLTCIRC.DAT` via
`fopen`; this conversion prints to stdout instead and drops that DOS
file-save-then-page convenience, the same approach used for `loan`
(Tier 1 suitability review, Finding 5).

## Worked Example

No worked numeric example was included with the original program. As
an independently verifiable check: four holes on a unit-radius circle
with no offset fall exactly on the cardinal points of the unit
circle, confirmed in this conversion's tests.
