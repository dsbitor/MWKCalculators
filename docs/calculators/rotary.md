# rotary

Rotary table division angles.

**Converted from:** `ROTARY.C` (M. W. Klotz),
`MWKC/WorkshopUtilities/rotary.zip`
**Go source:** `MWKGo/rotary/rotary.go`

## Purpose

Given the number of equal divisions wanted around a full circle, this
program lists the angular position of each division, both as a
decimal degree value and split into degrees, minutes, and seconds,
for setting a rotary table when the division count doesn't correspond
to a convenient round number of degrees.

No `.TXT` file was included with the original program; this purpose
statement is drawn from the `.C` file's own header comment.

## Inputs

| Prompt | Default |
|---|---|
| Number of divisions | 13 |

## Output

Every division from 0 up to and including the full circle (which
wraps back to 0 degrees), each as a decimal degree value and as
degrees/minutes/seconds.

## Method

```
step = 360 / numDivisions
theta_k = k * step             for k = 0 .. numDivisions
degrees = floor(theta_k)
minutes = floor((theta_k - degrees) * 60)
seconds = round((theta_k - degrees - minutes/60) * 3600)
```

with seconds carrying into minutes, and minutes into degrees, if
rounding pushes either to 60, and degrees wrapping from 360 back to
0 for the final (full-circle) division.

Unlike `compound`'s degrees/minutes/seconds conversion, which
truncates seconds and so never exercises its own carry branch, this
program rounds seconds to the nearest whole second (matching the
original's own `RND`), which does make the carry reachable: it fires
for `numDivisions = 25`, division 7 (exactly 100.8 degrees).

The original program writes its results to `ROTARY.OUT` via `fopen`;
this conversion prints to stdout instead and drops that DOS
file-save-then-page convenience, the same approach used for `loan`
(Tier 1 suitability review, Finding 5).

## Worked Example

No worked numeric example was included with the original program. As
an independently verifiable check: division 1 of 13 is exactly
360/13 decimal degrees, and the final division of any division count
must wrap to exactly 0 degrees rather than reporting 360; both
confirmed in this conversion's tests, along with the genuine
seconds-carry case found by brute-force search across division
counts.
