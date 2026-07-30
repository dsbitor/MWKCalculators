# combi

Enumerate all combinations of N things taken M at a time.

**Converted from:** `COMBI.C` (M. W. Klotz), `MWKC/Math/combi.zip`
**Go source:** `MWKGo/combi/combi.go`

## Purpose

`C(N,M) = N!/(M!(N-M)!)` gives the *count* of combinations, but often
what's actually needed is the full *list* — for exhaustively trying
every combination of some set of options, or fitting a set of
not-quite-identical parts. This prints every combination of N
letter-labeled things, taken M at a time.

This was misclassified during Tier 2's preparatory scan: it only ever
opens an *output* report file, never a data file, so it belongs with
the Tier 1 no-data-file calculators; see
`ai/plans/c-to-go-conversion-plan.md`'s clean-up phase note. Like
several other Tier 1 programs, this drops the original's file-save
feature and prints straight to stdout.

## Inputs

| Prompt | Default |
|---|---|
| Number of things | 6 |
| Taken m at a time | 4 |

N and M are both limited to 52 (the number of letters across the
upper- and lower-case alphabets, used to label each element), and M
must be less than N.

## Output

The total number of combinations, then each one numbered and printed
as a run of letters (element 1 is `A`, element 2 is `B`, and so on up
to element 52, `z`).

## Method

`next` is a direct port of `COMBI.C`'s own `combo()`: the classic
Nijenhuis & Wilf "revolving door" algorithm, which steps from one
combination to the next by complementing exactly two elements'
membership each time, visiting every combination exactly once before
returning to the start — a Hamiltonian cycle through all `C(N,M)`
combinations, not just any correct enumeration order. Its branching
depends on whether M is odd or even and where (scanning from the
right) the current combination first differs from a reference
pattern; two of its cases are reachable from either parity, so this
port expresses them once and selects them with a couple of boolean
flags instead of the original's `goto`-based jumps into shared code.

The combination count is computed exactly the way `COMBI.C` computes
it: multiply the full descending range from `M+1` to `N`, then divide
by `(N-M)!` one factor at a time — not the more commonly seen
formula that interleaves multiplication and division to stay exact at
every step. That difference matters: multiplying the full range first
can produce an intermediate value astronomically larger than the
final count (for `N=52, M=5`, roughly `8×10⁶⁴` before any division
happens), overflowing even a 64-bit integer for some otherwise
in-range inputs near the documented maximum of 52 — the same failure
`COMBI.C`'s own 32-bit `long` already had, not a regression
introduced by this conversion.

## Worked Example

`COMBI.TXT` ships an exact worked example: 6 things taken 4 at a
time, listing all 15 combinations in the algorithm's own generation
order. This conversion's tests reproduce that exact sequence,
including its very last entry: the algorithm starts at `ABCD`
internally but never reports it as combination 1 (the first `next()`
call happens before the first entry is ever recorded, matching
`COMBI.C`'s own loop) — instead, `ABCD` only reappears as combination
15, once the revolving-door cycle has visited every other combination
and returned to where it started. Further tests confirm the
combination count formula against known values, and that a larger
enumeration (8 things taken 3 at a time) produces the exact expected
number of combinations with no duplicates and no wrong-length entries.
