# remnant

One-dimensional stock-cutting from a pile of leftover remnant
lengths, accounting for saw kerf.

**Converted from:** `REMNANT.C` (M. W. Klotz), `MWKC/WorkshopUtilities/cuts.zip`
**Go source:** `MWKGo/remnant/remnant.go`
**Original documentation:** `CUTS.TXT`, inside `MWKC/WorkshopUtilities/cuts.zip` (not included in this conversion)

## Purpose

The same cutting-stock problem [cuts](cuts.md) and
[cutlist](cutlist.md) solve, but for a heterogeneous collection of
leftover remnant lengths rather than one uniform standard bar length,
and accounting for a saw kerf (material lost to the cut itself) on
every piece. Per `CUTS.TXT`, it was written for a user who had a pile
of odd-length leftover stock rather than fresh standard bars, reusing
Mike Graham's `cutlist` heuristic rather than the original author's
own (inferior, by his own account) algorithm.

## Inputs

A `-data <file>` flag names a data file: a saw kerf width, then
comma-separated `count,length` lines for the available remnants, a
`0,0` separator line, then comma-separated `count,size` lines for the
pieces needed, ending at `ENDOFDATA`. There is no `STARTOFDATA`
marker, the same convention `cuts`'s own data file uses.

## Output

The problem restated (kerf, remnants available, pieces needed), then
one line per remnant showing what was cut from it and its waste
(kerf plus leftover drop), followed by waste totals. A piece that fits
no remnant at all is reported rather than silently dropped.

## Method

Reuses `cutlist`'s best-fit heuristic (largest piece first, cut from
whichever available remnant currently has the smallest sufficient
remaining room), generalized so every remnant is available from the
start rather than opened lazily one at a time — since, unlike
`cuts`/`cutlist`'s identical standard bars, *which* remnant a piece
comes from matters here.

Unlike `cutlist`'s own refinement step (see its doc page), `REMNANT.C`'s
version accepts an exact zero-waste fit found later in the scan rather
than requiring a strictly tighter one — so, unlike `cutlist`, `remnant`
does not carry that particular asymmetry forward from the same
underlying algorithm. This conversion preserves that difference rather
than treating the two as identical.

`REMNANT.C`'s own "can't be assigned" check compares a stock array
index against `MAXPIECESPERLENGTH` — a limit on pieces *per* length,
not on the *number* of lengths available — which only works there
because that constant and `MAXLENGTHS` (the one actually meant) happen
to both be defined as 100. This conversion has no fixed-size array to
misapply a wrong bound against in the first place (its stock list is a
plain slice), so it simply detects "no remnant has enough room left"
directly, which is what the original check was clearly meant to test.

## Worked Example

`REMNANT.DAT`'s own shipped example (11 remnants across four lengths,
7 piece sizes needed, 0.03 kerf) is fully solvable with no pieces left
unassigned. This conversion's tests confirm that: every piece is cut
from some remnant, no remnant's cuts plus drop ever exceed its own
original length, and — the strongest available check, since no golden
reference output exists for this program — that kerf waste plus drop
waste reconciles exactly against the theoretical possible waste
(total remnant length minus total required length), which can only
hold if every single piece was actually assigned. A further test
confirms the exact-fit-preference difference from `cutlist` directly:
given remnants of length 10 and 6, a single piece of length 6 is cut
from the exact-fit remnant, not the looser 10-length one.
