# triangle

Solution of plane triangles.

**Converted from:** `TRIANGLE.C` (M. W. Klotz, 11/98, 3/99, 4/11),
`MWKC/Math/triangle.zip`
**Go source:** `MWKGo/triangle/triangle.go`

## Purpose

Given any three of a triangle's six parts (three sides and the three
angles opposite them), with at least one side among them, solves for
the rest and reports the area, circumscribed circle radius, and
inscribed circle radius. The side-side-angle (SSA) case is
inherently ambiguous — two different triangles can share the same two
sides and the angle opposite one of them — so both solutions are
reported whenever they genuinely differ.

## Inputs

| Prompt | Default |
|---|---|
| side 1, side 2, side 3 | (skip if unknown) |
| angle opposite side 1, 2, 3 | (skip if unknown) |

Exactly three data items are required, at least one of which must be
a side.

## Output

All three sides, all three angles, area, circumscribed circle radius,
inscribed circle radius; and, for an ambiguous SSA input, a second
full solution if it differs from the first.

## Method

Five classical triangle-solving techniques, dispatched according to
which three parts are known:

```
SSS: half-angle formula from all three sides
ASA: third angle from the angle sum, then law of sines for the other two sides
SAA: same as ASA, generalized to any one side plus its own angle and one other
SAS: law of cosines for the third side, then the SSS half-angle formula
SSA: law of sines for the second angle (two solutions when ambiguous),
     then the angle sum and a projection formula for the third side
```

The original program implements this dispatch as roughly twenty
near-identical branches, each calling one of five shared helper
functions with a different permutation of its own six values bound
to the helper's generic parameter names (playing "side 1" for one
call, "side 2" for another, and so on). This conversion instead
generalizes each helper to take explicit index-independent arguments
(the "own" side/angle, an ascending pair of indices, and so on) and
iterates over the relevant index combinations directly, avoiding
manual re-derivation of all twenty branches; each generalized case
was independently verified, alongside every closed-form combination,
against a single ground-truth triangle before being trusted.

## Worked Example

No worked numeric example was included with the original program. As
independently verifiable checks: the 3-4-5 right triangle (solved via
SSS) has a well known exact 90 degree angle; every other non-ambiguous
combination of three known parts recovers that same triangle exactly;
the classic ambiguous case (sides 8 and 10, with a 40 degree angle
opposite the side of length 8) produces two solutions, both of which
independently satisfy the law of sines; and the circumscribed
circle's radius matches the standard `abc/(4*Area)` formula,
independent of this code's own law-of-sines-based computation of the
same value. All confirmed in this conversion's tests.
