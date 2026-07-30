# sinebar

Sine bar calculations.

**Converted from:** `SINEBAR.C`, `MWKC/WorkshopUtilities/sinebar.zip`
**Go source:** `MWKGo/sinebar/sinebar.go`

## Purpose

Given the distance between a sine bar's rolls and the angle to set
(entered as decimal degrees or degrees/minutes/seconds), computes the
gage block stack height needed under one roll, how sensitive that
angle is to small measurement errors in either the roll distance or
the stack height, and which blocks from a standard 81-piece gage
block set combine to produce that stack height exactly.

## Inputs

| Prompt | Default |
|---|---|
| Distance between sine bar rolls | 5 |
| Angle input mode [D]ecimal degrees, (X) deg/min/sec | [D] |
| angle in decimal degrees (if D) | 4.93468141 deg |
| degrees, minutes, seconds (if X) | 30, 7, 30 |

## Output

Stack height; for angles under 90 degrees, the angle error caused by
a 0.001 error in the roll distance or in the stack height; and the
gage blocks needed to form the stack height.

## Method

```
stack = rollSeparation * sind(angle)
```

Unlike `mtaper`'s `gaugeBlocksFor`, which decomposes a stack via a
greedy digit-extraction heuristic, this program's block search is a
bounded brute-force combination search: try every single block from
the 81-piece set no larger than the target, then every valid pair,
and so on up to 5 blocks, stopping at the first combination that sums
to the target exactly. This reuses the same nested-loop "odometer"
technique as `gearfind` in Tier 1 group 11. The original program's
interactive keypress abort (there being no keyboard to poll here) is
replaced with an explicit evaluation-count bound per
`coding-style.md` Rule 2; an extremely small target stack (forcing
the search to consider nearly the full 81-block set in combinations
of up to 5) could exceed that bound and report no solution found even
where one exists at much greater computational cost, the same
practical effect the original's own keypress abort would have on an
impatient user.

## Worked Example

No worked numeric example was included with the original program. As
independently verifiable checks: a 30 degree sine bar angle with a
10-unit roll separation gives a stack height of exactly half the roll
separation (`sin(30) = 0.5`, an independently obvious trigonometric
fact), and every combination the block search returns sums back to
its target exactly, rather than being asserted against one specific
expected combination; both confirmed in this conversion's tests.
