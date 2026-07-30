# osborne

Convergence demonstration for the "Osborne Maneuver" of centering
round stock.

**Converted from:** `OSBORNE.C` (M. W. Klotz, 11/99),
`MWKC/WorkshopUtilities/osborne.zip`
**Go source:** `MWKGo/osborne/osborne.go`
**Original documentation:** `OSBORNE.TXT`, inside `MWKC/WorkshopUtilities/osborne.zip` (not included in this conversion)

## Purpose

In Guy Lautard's "Home Machinist's Bedside Reader #2", the
"Osborne Maneuver" centers round stock in a milling machine
using nothing more than an edge finder: find the edge on one
axis, move in by half the diameter, find the edge on the other
axis, move in by half the diameter, and repeat. Each repetition
uses the leftover offset from the previous axis as the starting
offset for the next. This program shows how quickly that process
converges on the true center, for a given workpiece diameter and
initial alignment error.

## Inputs

| Prompt | Default |
|---|---|
| Workpiece diameter | 2.0 |
| Initial offset | 0.1 |

## Output

Six iterations, each printing the offset going in, the resulting
offset on the other axis, and the combined radial error (the
root-sum-square of the two offsets).

## Method

```
theta = asin(offset / radius)
nextOffset = radius * (1 - cos(theta))
error = sqrt(offset^2 + nextOffset^2)
```

Each iteration's `nextOffset` becomes the next iteration's
`offset`.

## Worked Example

`OSBORNE.TXT` gives the author's own example run at the
program's default inputs (diameter 2, initial offset 0.1):

| Iteration | del1 | del2 | error |
|---|---|---|---|
| 1 | 0.10000000 | 0.00501256 | 0.10012555 |
| 2 | 0.00501256 | 0.00001256 | 0.00501258 |
| 3 | 0.00001256 | 0.00000000 | 0.00001256 |

This conversion's output matches these figures exactly.
