# spaceblk

Space block selection utility.

**Converted from:** `SPACEBLK.C` (M. W. Klotz), `MWKC/WorkshopUtilities/spaceblk.zip`
**Go source:** `MWKGo/spaceblk/spaceblk.go`

## Purpose

"Space blocks" are a poor man's precision gage block set: ground
steel cylinders, centrally threaded so they can be screwed together
to build an accurate spacer of a required length from a limited
assortment of individual block sizes. Given a required spacer size,
this searches the owner's own block set for a combination (of up to
five blocks) that adds up to it.

## Data setup

The available block set is specific to one owner's own blocks, so —
like `diffthrd` and `divhead` in an earlier conversion group — this
program reads it from the user's own database (`userdata.db`) rather
than the shared `reference.db`; see
`ai/plans/c-to-go-conversion-plan.md`, "Data-file strategy for
Tier 2". That table ships empty. Import your own block set before
first use:

```
spaceblk -import my-block-set.csv
```

The CSV needs one column, `size`. A worked example built from the
original archive's `SPACEBLK.36` (a typical 36-block Imperial set,
included alongside the also-shipped `SPACEBLK.81` standard 81-block
set) ships at `MWKGo/spaceblk/testdata/example-blocks.csv`.

This archive also contains `GAGEBLK.C`, a faster successor program
specialized for exactly the standard 81-block Imperial set (needing
no data file at all, since the block sizes are built into its code).
It is deferred to a future group: `spaceblk` alone already covers its
case, more slowly, as well as any other block set `GAGEBLK` doesn't
handle.

## Inputs

| Prompt | Default |
|---|---|
| Required spacer size | 0.104 |

## Output

The block set summary (count, total length, largest and smallest
block); then either the blocks to use and their total length, or a
"can't find a solution" message.

## Method

The requested size is rounded to four decimal places (matching the
original's own precision limit), then searched for as a combination
of 1, then 2, and so on up to 5 blocks (the original's own fixed
limit), each combination drawn from blocks no larger than the target,
using the same odometer-style nested index search pioneered for
`gearfind` in Tier 1: indices increment like a counter, rightmost
fastest, skipping any combination that reuses the same block position
twice (two blocks of the same *size* can still both be used, since
they occupy different positions — a real 36- or 81-block set commonly
has duplicate sizes). The search reports the first combination it
finds summing to the target within 1e-8, not necessarily the one
using the fewest blocks, matching the original exactly; a partial sum
already exceeding the target abandons that combination immediately
rather than finishing the addition, the same short-circuit the
original uses. `coding-style.md` Rule 2 replaces the original's
interactive keypress abort (there is no keyboard to poll here) with
an explicit evaluation cap: `SPACEBLK.TXT` documents that a full
5-deep search of an 81-block set examines 3.5 billion combinations,
so the cap exists to stop well short of that if no solution turns up,
rather than running effectively forever.

## Worked Example

No worked numeric example was included with the original program.
As independently verifiable checks, this conversion's tests confirm
a target exactly matching one block resolves at a single-block depth,
a target requiring two blocks sums to exactly the target within
tolerance, two same-size blocks (as in the original's own two 1.0 in
blocks) can both be used in one combination, and a target smaller
than every available block correctly reports no solution.
