# diffthrd

Differential thread calculations.

**Converted from:** `DIFFTHRD.C` (M. W. Klotz), `MWKC/WorkshopUtilities/diffthrd.zip`
**Go source:** `MWKGo/diffthrd/diffthrd.go`

## Purpose

A differential thread cuts two threads of different pitch on the
same shaft; as the shaft turns, the difference between the two
pitches produces a very fine *effective* pitch, useful for slow,
precise linear motion (fine focus mechanisms, micrometer-style
adjusters) using only ordinary screwcutting gears. This calculator
searches the pitches your lathe can actually cut for the pair whose
effective pitch comes closest to a desired value, then works out the
nut dimensions and crank turns needed to use it.

## Data setup

Unlike `fits`, `speed`, `gage`, and `expand` earlier in this
conversion group, the available thread pitches are specific to one
lathe's own change-gear train, not universal reference data — the
original `DIFFTHRD.DAT`'s own header says as much: "The default
values shown here are what my lathe can cut." So this program reads
from the user's own database (`userdata.db`) rather than the shared
`reference.db`; see `ai/plans/c-to-go-conversion-plan.md`,
"Data-file strategy for Tier 2". `userdata.db` ships with this table
empty — nothing is pre-loaded, since another user's lathe can cut a
different set of pitches entirely. Import your own list before first
use:

```
diffthrd -import my-lathe-pitches.csv
```

The CSV needs one column, `pitch_tpi` (threads per inch; convert a
metric mm/thread pitch to tpi as `25.4 / mm`). A worked example built
from the original `DIFFTHRD.DAT`'s own default pitch list (40
imperial pitches plus 24 metric pitches, metric already converted to
tpi) ships at `MWKGo/diffthrd/testdata/example-pitches.csv`.

## Inputs

| Prompt | Default |
|---|---|
| Desired effective pitch of differential thread (tpi) | 100.0 |
| Pitch of coarse thread (tpi) | the best match found |
| Pitch of fine thread (tpi) | the best match found |
| Thickness of coarse (fixed) nut (in) | 0.375 |
| Thickness of fine (movable) nut (in) | 0.25 |
| Desired motion of movable nut (in) | 0.25 |

## Output

The best available coarse/fine pitch pair and the effective pitch it
produces, then (after confirming or overriding that pair and the nut
dimensions) the crank turns, minimum thread lengths, and nut spacing
needed to achieve the desired motion.

## Method

For every ordered pair of available pitches, with the smaller (coarse)
pitch as `pc` and the larger (fine) pitch as `pf`, the differential
effective pitch is `1 / (1/pc - 1/pf)`; the pair whose effective pitch
is closest to the desired value is reported. This is an exhaustive
O(n²) search over the available pitches, matching the original's own
nested loop, and bounded by however many pitches the user has
actually imported.

The original also computes a "sort the pitches" step after loading
them, but its `do { ... } while (f)` loop never actually sets `f`
inside the swap, so it only ever performs a single bubble-sort pass —
and since the search above already tries every ordered pair
regardless of storage order, sorting (complete or not) has no effect
on the result at all; this conversion does not reproduce that step,
an internal simplification rather than a preserved quirk, since it is
provably inconsequential to any output. The original also prints a
stray `printf ("%.4lf\n",pitch[np-1]);` immediately after loading data
— an unlabeled debug leftover with no connection to anything the user
asked for — which this conversion likewise does not reproduce.

## Worked Example

No worked numeric example was included with the original program. As
an independently verifiable check, this conversion's tests recompute,
for the pair `bestPair` returns, every other pair's effective pitch
directly and confirm none is closer to the desired value — the actual
defining property of "best match" — and separately confirm
`csvtable.Import`'s round trip through the example CSV reproduces the
same pitch list `loadPitches` reads back.
