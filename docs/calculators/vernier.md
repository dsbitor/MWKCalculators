# vernier

Two-plate angular vernier design.

**Converted from:** `VERNIER.C`, `MWKC/WorkshopUtilities/vernier.zip`
**Go source:** `MWKGo/vernier/vernier.go`

## Purpose

A single index plate needing many equal angular divisions needs a
hole for every one of them. A two-plate vernier instead uses two
plates, with `n1` and `n2` holes respectively, mounted concentrically
so that some pair of holes (one per plate) can always be brought into
alignment to mark off any of the required divisions — using far
fewer total holes than one plate would need directly. Given the
number of divisions required, this program factors that number,
searches for the pair of plate hole counts needing the fewest total
holes, and, if a two-plate solution actually saves holes over a
single plate, produces the full table of which holes to align for
each division.

## Inputs

| Prompt | Default |
|---|---|
| Number of angular subdivisions required | 360 |

## Output

The prime factorization of the division count; the selected plate
hole counts; and, if a two-plate solution is worthwhile, a table of
each division's angle and the letter/number hole pair that aligns it,
plus the total holes actually needed to drill on each plate.

## Method

Candidate plate pairs satisfy `1/n1 - 1/n2 = 1/numDivisions`
(computed via integer division, so only near-exact candidates are
tried); each candidate is verified by checking that every one of the
`numDivisions` divisions actually has some pair of holes that aligns
it to within 0.0001 degrees, and the pair needing the fewest total
holes is kept.

The original program computes the same prime factorization using a
sieve of small primes, a performance optimization suited to the
DOS-era hardware it targeted; this conversion uses plain trial
division instead, which produces an identical factorization without
reproducing the sieve, since the sieve is an internal performance
detail with no effect on the program's output. The original's
interactive keypress abort during the plate-pair search (there being
no keyboard to poll here) is replaced with an explicit bound on the
number of alignment checks performed, per `coding-style.md` Rule 2.

This is the third distinct algorithm in this project tackling a
gage-block/index-plate-adjacent combinatorial problem, after
`mtaper`'s greedy digit extraction and `sinebar`'s bounded
combination search in Tier 1 group 14; all three are kept as
separate, faithful conversions of their respective originals.

## Worked Example

No worked numeric example was included with the original program. As
independently verifiable checks: for the documented default (360
divisions), the search is independently confirmed (via a separate
brute-force script) to select 36 and 40 holes for 76 total; and,
given that pair, every one of the 360 requested divisions must
actually have a valid hole alignment, checked directly against the
found plates rather than assuming the search's own internal
verification was applied correctly. Both confirmed in this
conversion's tests.
