# gearatio

Find a chain of gears matching a required ratio.

**Converted from:** `GEARATIO.C` (M. W. Klotz), `MWKC/WorkshopUtilities/gearatio.zip`
**Go source:** `MWKGo/gearatio/gearatio.go`

## Purpose

Anyone with a scrounged or purchased set of gears (commonly a lathe's
own set of change gears) can use this to find which of them, chained
together, produce a required step-up or step-down ratio for some other
project — a gearbox, a jig, a geared attachment — independent of any
particular threading calculation.

## Data setup

The available gear set is specific to one owner's own gears, so — like
`diffthrd`, `divhead`, and `spaceblk` in earlier conversion groups —
this program reads it from the user's own database (`userdata.db`)
rather than the shared `reference.db`; see
`ai/plans/c-to-go-conversion-plan.md`, "Data-file strategy for
Tier 2". That table ships empty. Import your own gear set before first
use:

```
gearatio -import my-gears.csv
```

The CSV needs one column, `teeth`. A worked example built from the
original archive's own shipped `GEARATIO.DAT` (11 gears from 24 to 100
teeth) ships at `MWKGo/gearatio/testdata/example-gears.csv`.

## Inputs

| Prompt | Default |
|---|---|
| Required ratio | 0.5 |
| Allowable ratio accuracy (%) | 0.01 |
| Maximum gear pairs to examine | 1 |

## Output

Every gear chain found within tolerance: the gears used, printed as
`teeth:teeth` per stage joined by `-` for chains longer than one
stage, then the chain's actual combined ratio and its error from the
target. If none is found within the examined chain lengths, a
"NO SOLUTION FOUND" message.

## Method

Gears are considered in pairs (two gears meshed together, or in a
compound train, two gears on the same shaft each meshed with another
gear). A chain of 1 to the requested maximum number of pairs is
searched: for chains longer than one pair, no single gear (by
position — two gears sharing the same tooth count, as several of
`GEARATIO.DAT`'s own alternate data sets have, remain distinct
positions and can both be used) may appear twice within one chain,
since a single physical gear cannot be mounted in two places on the
same train at once.

At each stage, the running ratio is multiplied by that stage's own
gear ratio if the running ratio is still above the target, or divided
by it otherwise — swapping which gear of the pair is displayed first
to match — exactly matching `GEARATIO.C`'s own greedy per-stage
decision. Every complete chain landing within the requested tolerance
is reported, not just the first found, matching the original's own
behavior of examining chain lengths from 1 up to the requested maximum
and reporting every match at every length. `coding-style.md` Rule 2
replaces the original's interactive keypress abort (there is no
keyboard to poll here) with an explicit search-evaluation cap, and
also caps the requested chain length itself at 10 — `GEARATIO.C`'s own
fixed array bound, which the original never actually validated a
user's input against.

## Worked Example

No worked numeric example was included with the original program
(`GEARATIO.TXT` is a general description only). As independently
verifiable checks, this conversion's tests confirm an exact single-pair
ratio present in the shipped gear set (24:48 = 0.5) is found, that
every reported solution's stages actually multiply out to its recorded
ratio and fall within the requested tolerance, that identical-size gear
pairs (a useless 1:1 ratio) are skipped, that an unreachable target
reports no solution, and that a chain requiring every available gear
position uses each position exactly once.
