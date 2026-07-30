# ugroove

Incremental cutting of a U-shaped groove.

**Converted from:** `UGROOVE.C`, `MWKC/WorkshopUtilities/ugroove.zip`
**Go source:** `MWKGo/ugroove/ugroove.go`

## Purpose

Written for cutting a U-shaped groove in a workpiece (a giant O-ring
groove, for making pipe dies), this program computes a series of
incremental plunge cuts with a round-nosed tool, spaced either by a
fixed angular increment around the groove's profile or by a fixed
linear increment, leaving small cusps that are then filed smooth.

## Inputs

| Prompt | Default |
|---|---|
| Groove radius | 1.0 in |
| Lathe tool diameter | 0.25 in |
| [A]ngular or (L)inear increment | [A] |
| Angular increment (if angular) | 5.0 deg |
| Linear increment (if linear) | 0.0625 in |

## Output

A table of offset from the groove's center (X) and the depth of cut
to make at that offset (DOC).

## Method

```
r = grooveRadius - 0.5*toolDiameter

angular mode (theta from 0 to 90 deg):
  x = r*cos(theta)
  DOC = sqrt(r^2 - x^2) + 0.5*toolDiameter

linear mode (x from r down to 0):
  DOC = sqrt(r^2 - x^2) + 0.5*toolDiameter
```

The original program writes its results to `UGROOVE.OUT` via
`fopen`; this conversion prints to stdout instead and drops that DOS
file-save-then-page convenience, the same approach used for `loan`
(Tier 1 suitability review, Finding 5).

## Worked Example

No worked numeric example was included with the original program. As
independently verifiable checks: at the groove's centerline the cut
depth equals exactly the tool's own radius (the shallowest cut), and
at the groove's edge the cut depth equals exactly the full groove
radius (the deepest cut); both hold identically whether stepping by
angle or by linear distance, confirmed in this conversion's tests.
