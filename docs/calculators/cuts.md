# cuts

One-dimensional stock-cutting: cut a list of needed piece lengths
from as few standard-length bars as possible.

**Converted from:** `CUTS.C` (M. W. Klotz), `MWKC/WorkshopUtilities/cuts.zip`
**Go source:** `MWKGo/cuts/cuts.go`
**Original documentation:** `CUTS.TXT`, inside `MWKC/WorkshopUtilities/cuts.zip` (not included in this conversion)

## Purpose

The classic 1D cutting-stock problem: given a standard bar length and
a list of piece sizes and counts needed, work out how to cut them from
as few bars as possible, minimizing waste. Per `CUTS.TXT`, this applies
to any one-dimensional medium — pipe, lumber, wire, rebar — not just
literal bars.

`cuts` is one of two different heuristics this project converts for
the same problem (see also [cutlist](cutlist.md)): its own author's
greedy, per-bar exhaustive combinatorial search, optionally biased
toward zero-waste combinations first.

## Inputs

A `-data <file>` flag names a data file: a standard bar length,
followed by comma-separated `count,size` lines for each piece needed
(order doesn't matter — the program sorts by size), ending at
`ENDOFDATA`. There is no `STARTOFDATA` marker; the data region starts
at the first non-comment line.

A `-zero-waste` flag (default off) matches the original's own
interactive "Search for zero waste solutions first (Y/[N])?" prompt:
each bar's search first tries to find a zero-waste combination,
falling back to the ordinary least-waste search if none exists.

## Output

For each bar (or group of identically-cut consecutive bars), the
combination of pieces cut and the waste. A final summary reports total
waste and total bars used, alongside the theoretical minimum of each
computed up front from the total length needed.

## Method

For each bar, an exhaustive but bounded search considers every
combination of still-needed pieces that fits (bounded by how many of
a size both fit in one bar and are still needed), preferring least
waste, and — among equally-wasteful combinations — the one weighted
more heavily toward larger pieces, since cutting big pieces first
leaves more flexibility for the smaller ones later. This repeats one
bar at a time until every piece is accounted for.

Per `CUTS.TXT`'s own "Update 2/02" note, neither the default search
nor `-zero-waste` always reaches the true theoretical optimum, and
which one wins varies by problem — the original author explicitly
recommends trying both, which is why this project converts both
`cuts` and its sibling `cutlist` rather than picking one.

Despite this program's original planning note describing it as
"stateful" and needing an atomic-write persistence pattern, it turns
out to only read its own data file and print to stdout — nothing here
is stateful, and no output file is ever written. The original's
"debug" mode (triggered by any command-line argument, printing
internal search diagnostics and stopping after a single bar) is
developer scaffolding, not a user-facing feature, and is dropped
entirely.

## Worked Example

`CUTS.TXT` includes a full worked example (its own shipped `CUTS.DAT`:
a 6-unit bar, 7 piece sizes) with a complete expected cutting schedule
and totals (waste 5.25, 12 bars), reproduced exactly by this
conversion's tests. It also documents a second, smaller case (cutting
`1,10 / 2,7 / 1,6 / 2,4` from 20-unit stock) specifically to
demonstrate `-zero-waste`'s effect: the default search wastes 22 units
across 3 bars, while `-zero-waste` finds a 2-bar, 2-unit-waste
solution — both reproduced exactly by this conversion's tests.
