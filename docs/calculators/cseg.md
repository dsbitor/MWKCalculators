# cseg

Circular segment calculations.

**Converted from:** `CSEG.C`, `MWKC/Math/cseg.zip`
**Go source:** `MWKGo/cseg/cseg.go`

## Purpose

A circular segment is the region between a chord and the arc it cuts
off. Given any two of its five related dimensions (radius, included
angle, chord length, height/sagitta, or arc length), this program
solves for the rest and computes the segment's area.

## Inputs

| Prompt | Default |
|---|---|
| Radius of segment | (skip if unknown) |
| Segment included angle | (skip if unknown) |
| Chord of segment | (skip if unknown) |
| Height of segment (sagitta) | (skip if unknown) |
| Arc length | (skip if unknown) |

Exactly two of the five are required.

## Output

Radius, angle, chord, height, arc length, and area.

## Method

Eight of the ten possible pairings have closed forms, all derived
from `height = radius*(1-cos(angle/2))`, `chord =
2*radius*sin(angle/2)`, and `arc = radius*angle(radians)`. The
remaining two (chord+arc and height+arc) have no closed form, so they
use the original program's own zoom-in grid search: scan a
1-degree-resolution window for the angle that best matches the known
quantity, then repeat on a shrinking window around the best candidate
found, ten times finer each pass, for 6 passes total.

```
area = pi*radius^2*(angle/360) - 0.5*(radius-height)*chord
```

## Worked Example

No worked numeric example was included with the original program. As
an independently verifiable check: a single ground-truth segment
(radius 5, included angle 40 degrees) is used to confirm that all
eight closed-form pairings, and both search-based pairings, recover
the same consistent radius, angle, chord, height, arc, and area;
confirmed in this conversion's tests, with a looser tolerance for the
two search-based cases to reflect their finite (not floating-point
exact) convergence.
