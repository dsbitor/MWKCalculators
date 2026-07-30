# ballcut

Incremental sphere turning on a lathe.

**Converted from:** `BALLCUT.C`, `MWKC/WorkshopUtilities/ballcut.zip`
**Go source:** `MWKGo/ballcut/ballcut.go`
**Original documentation:** `BALLCUT.TXT`, inside `MWKC/WorkshopUtilities/ballcut.zip` (not included in this conversion)

## Purpose

Turning a spherical shape on a lathe can be approximated with a
series of plunge cuts using a squared-off cutoff tool, producing a
"staircase" shape that's then filed smooth. Given the desired sphere
diameter, the stock diameter it's being cut from, and a step size
(either a fixed angular increment around the quarter-circle profile,
or a fixed linear increment along the lathe bed), this program
computes each cut's axial position and depth.

## Inputs

| Prompt | Default |
|---|---|
| Sphere diameter | 1.0 in |
| Stock diameter | (defaults to sphere diameter) |
| Constant [A]ngle step or constant (X) step | [A] |
| Angular increment (if angle step) | 5.0 deg |
| X increment (if X step) | 0.02 in |

## Output

A table of cut number, axial tool position, its increment from the
last cut, depth of cut, its increment from the last cut, and the
resulting work diameter.

## Method

Angle-stepped mode walks theta from 0 to 90 degrees:

```
axial = r - r*cos(theta)
radial = r*sin(theta)
depth = stockRadius - radial      (stop once depth < 0)
workDiameter = 2*radial
```

Linear-stepped mode walks the axial position `x` directly from the
sphere radius down to 0:

```
radial = sqrt(r^2 - (x-r)^2)
depth = stockRadius - radial      (stop once depth < 0)
```

The original program writes its results to `BALLCUT.OUT` via
`fopen`; this conversion prints to stdout instead and drops that DOS
file-save-then-page convenience, the same approach used for `loan`
(Tier 1 suitability review, Finding 5).

## Worked Example

`BALLCUT.TXT` includes two complete worked example tables (angular
increment mode), both reproduced in this conversion's tests: a 1"
sphere from 1" stock (19 rows, N=0 to N=18) and a 2" radius cut on 1"
stock (terminating once the required depth would exceed the
available stock). The second example's exact termination point sits
at a genuine floating-point boundary: at theta=30 degrees the
required depth computes to a value that is mathematically exactly
zero but numerically an infinitesimal epsilon either side of zero
depending on the runtime's trigonometric rounding, so the original
DOS executable's own cutoff (after N=5) isn't guaranteed to reproduce
bit-for-bit in Go; this conversion's test checks the unambiguous
early rows only, rather than asserting a specific total row count at
that boundary.
