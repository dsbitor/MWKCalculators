# sphere

Solution of spherical triangles.

**Converted from:** `SPHERE.C` (M. W. Klotz, 3/99), `MWKC/Math/sphere.zip`
**Go source:** `MWKGo/sphere/sphere.go`

## Purpose

Given any three of a spherical triangle's six parts (three sides and
the three angles opposite them, all measured in degrees, with sides
expressed as the angle they subtend at the sphere's center), solves
for the rest and reports the triangle's area as a multiple of the
sphere's radius squared. Unlike a plane triangle, three known angles
alone (AAA) are enough to solve a spherical triangle, since spherical
triangles of different sizes aren't similar to each other. As with
`triangle`'s SSA case, the side-side-angle and angle-angle-side cases
here are inherently ambiguous, so both solutions are reported when
they genuinely differ.

## Inputs

| Prompt | Default |
|---|---|
| side 1, side 2, side 3 | (skip if unknown) |
| angle opposite side 1, 2, 3 | (skip if unknown) |

Exactly three data items are required (unlike the plane `triangle`,
no side is strictly required, since AAA alone determines a spherical
triangle).

## Output

All three sides, all three angles, and the area (as a multiple of the
sphere's radius squared); plus a second full solution for an
ambiguous SSA or AAS input, if it differs from the first.

## Method

Six classical spherical-triangle-solving techniques:

```
SSS: spherical law of cosines (for sides) for each angle
AAA: spherical law of cosines (for angles) for each side — the
     "polar dual" of SSS, with no plane-triangle counterpart
SAS: law of cosines for the third side, then two more law-of-cosines calls
ASA: the dual law of cosines for the third angle, then two sides
SSA: law of sines for the second angle (two solutions when ambiguous),
     then a formula involving the half-angle/half-side tangent for the
     third side, then the law of cosines for the third angle
AAS: the dual of SSA (sides and angles swapped throughout)
```

Following the same approach used for `triangle` in Tier 1 group 11,
the original program's roughly twenty near-identical branches (each
calling one of six shared helper functions with a different
permutation of its own six values) are generalized here into six
index-independent helper functions, with the correct index
combinations for each technique derived and verified against a
single ground-truth spherical triangle before being trusted, rather
than transcribed branch-by-branch.

## Worked Example

No worked numeric example was included with the original program. As
independently verifiable checks: a ground-truth spherical triangle
(sides 60, 70, 80 degrees, solved via SSS) is used to confirm that
AAA, SAS, ASA, SSA, and AAS all recover the same full solution when
given the corresponding three parts of that same triangle; and the
classic ambiguous SSA case (a smaller side opposite a known angle,
paired with a larger second side) produces two solutions, both of
which independently satisfy the spherical law of cosines. All
confirmed in this conversion's tests.
