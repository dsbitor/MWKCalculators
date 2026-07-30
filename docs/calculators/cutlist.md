# cutlist

One-dimensional stock-cutting via a "best fit decreasing" heuristic.

**Converted from:** `CUTLIST.C` (Mike Graham), `MWKC/WorkshopUtilities/cuts.zip`
**Go source:** `MWKGo/cutlist/cutlist.go`
**Original documentation:** `CUTS.TXT`, inside `MWKC/WorkshopUtilities/cuts.zip` (not included in this conversion)

## Purpose

The same 1D cutting-stock problem [cuts](cuts.md) solves — cutting a
list of needed piece lengths from as few standard-length bars as
possible — using a different heuristic, contributed by Mike Graham to
address cases where `cuts`'s own greedy search falls short. Per
`CUTS.TXT`, sort every needed piece largest first, then for each piece
in turn, cut it from whichever already-opened bar currently has the
smallest remaining room that still fits it, opening a new bar only
when none of the already-opened ones do ("best fit decreasing").

`CUTS.TXT` quotes the original author calling this "definitely
superior... runs faster, and is just generally a better way to do the
problem" than his own `cuts` — but also that `cuts` sometimes still
wins, hence converting both.

## Inputs

A `-data <file>` flag names a data file in the same format `cuts`
uses: a standard bar length, then comma-separated `count,size` lines,
ending at `ENDOFDATA`. Unlike `cuts`, the original doesn't actually
require the `ENDOFDATA` marker — it reads until end of file, acting on
the marker only if present — and this conversion keeps that leniency.

## Output

The problem restated (piece counts and sizes, theoretical minimums),
then the cutting list: one entry per bar used, listing what was cut
from it, followed by actual waste and bar count totals.

## Method

`CUTLIST.C` works entirely in thousandths of a unit as integers (every
parsed length is multiplied by 1000 and truncated) specifically to
avoid a floating-point drift problem its own header comment mentions
correcting. This conversion keeps that fixed-point representation for
the same reason, rather than reintroducing the float comparisons
[cuts](cuts.md) itself still needs a tolerance helper to work around.

The original's best-fit search has a subtle asymmetry, preserved here
rather than silently tightened: an initial scan across every bar slot,
including still-unopened ones, accepts the first with enough room
(inclusive); a refinement scan across only already-opened bars then
replaces that pick only on a *strictly* tighter fit. A later
already-opened bar that turns out to be an exact zero-waste fit is
therefore never preferred over an earlier, looser-fitting one already
found — see [remnant](remnant.md), which generalizes this same
algorithm and does not carry this particular asymmetry forward.

## Worked Example

`cutlist`'s tests run it against the same shipped `CUTS.DAT` dataset
`cuts` uses, confirming it reaches the same theoretical-minimum totals
(waste 5.25 across 12 standard lengths) via its own, different
assignment — along with the basic correctness invariant that every
bar's cuts plus its remaining drop always add back up to exactly the
standard length.
