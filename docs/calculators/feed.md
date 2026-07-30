# feed

Milling feed rate solver.

**Converted from:** `FEED.C` (M. W. Klotz, 2/01),
`MWKC/WorkshopUtilities/feed.zip`
**Go source:** `MWKGo/feed/feed.go`

## Purpose

Feed rate, chip load, cutting speed, and number of cutting edges are
all related by one formula; given any three of the five related
quantities (speed, cutting edges, chip load, and feed expressed as
either in/min or in/rev), this program solves for the rest.

No `.TXT` file was included with the original program; this purpose
statement is drawn from the `.C` file's own header comment.

## Inputs

Enter whichever of the following you know; leave the rest blank.
Exactly three data items (per one of four recognized combinations,
below) are required.

| Prompt | Unit |
|---|---|
| Speed | rpm |
| Number of cutting edges | edges/revolution |
| Chip load | in/cutting edge |
| Feed | in/min |
| Feed | in/revolution |

## Output

Cutting edges, speed, chip load, and feed in both in/min and in/rev.

## Method

The underlying relationship: `feed(in/min) = chip load * speed *
edges`, and `feed(in/rev) = feed(in/min) / speed`. Four combinations
of three known quantities are recognized:

```
speed, edges, load            -> feed(in/min), feed(in/rev)
speed, edges, +feed            -> load
speed, load, +feed             -> edges
edges, load, feed(in/min)      -> speed, feed(in/rev)
```

`edges` and `load` together with `feed(in/rev)` (but not
`feed(in/min)`) is deliberately not solvable: `feed(in/rev)` is
already fully determined by `edges` and `load` alone
(`feed(in/rev) = load*edges`), so it adds no information toward
finding `speed`. The original program never even prompts for
`feed(in/rev)` once it detects both `edges` and `load` are already
known, for exactly this reason; this conversion follows the same
prompting order so the same case is never asked for.

## Worked Example

No worked numeric example was included with the original program. As
an independently verifiable check: a single self-consistent scenario
(1000 rpm, 4 edges, 0.002 in/edge chip load, giving 8 in/min = 0.008
in/rev) is solved from each of the four recognized combinations of
three known quantities, and each recovers the full, consistent set;
confirmed in this conversion's tests, along with the two
insufficient-data cases (too few known values, and the
edges+load+feed-per-rev non-solvable combination).
