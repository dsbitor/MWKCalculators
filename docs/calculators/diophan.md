# diophan

Linear Diophantine equation solver.

**Converted from:** `DIOPHAN.C`, `MWKC/Math/diophan.zip`
**Go source:** `MWKGo/diophan/diophan.go`

## Purpose

Given integers a, b, and c, finds integers x and y satisfying
`ax + by = c` (a linear Diophantine equation), using the extended
Euclidean algorithm, then reports the general solution family and
several sample solutions.

## Inputs

| Prompt | Default |
|---|---|
| Integer value of a | 172 |
| Integer value of b | 20 |
| Integer value of c | 1000 |

`c` must be evenly divisible by `gcd(a, b)` for a solution to exist.

## Output

One particular solution `(x, y)`, the general solution family
`x = x0 + k*(b/g)`, `y = y0 - k*(a/g)` for integer k, and nine sample
solutions for k = -4 through 4.

## Method

Standard extended Euclidean algorithm: expand the continued fraction
of `a/g` and `b/g` (where `g = gcd(a,b)`) to find Bezout coefficients,
then scale by `c/g` to solve the general equation rather than just
`ax + by = g`. Ported directly from the original program's own
array-based implementation, including its indexing convention, since
the algorithm is intricate enough that a fresh derivation risked
introducing a translation error; the direct port was verified against
the documented default input via the equation's own defining property
before being trusted.

## Worked Example

No worked numeric example was included with the original program. A
linear Diophantine equation's own solution is self-verifying: for
any (x, y) returned, `a*x + b*y` must equal `c` exactly. This is
checked directly in this conversion's tests, both for the documented
default input and for every sample solution across the full k = -4
to 4 range, rather than comparing against a specific hand-picked
`(x, y)` pair (any of infinitely many integer solutions is equally
correct). The unsolvable case (`c` not divisible by `gcd(a,b)`) is
also covered.
