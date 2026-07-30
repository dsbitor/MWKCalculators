# gage

Wire and sheet gage utility.

**Converted from:** `GAGE.C` (M. W. Klotz), `MWKC/WorkshopUtilities/gage.zip`
**Go source:** `MWKGo/gage/gage.go`

## Purpose

Wire and sheet metal thickness is conventionally specified by a gage
number rather than a direct measurement, and the mapping between gage
number and thickness is different for wire and for sheet. This
calculator converts in both directions: designation to size, or a
measured size to the closest matching designation.

## Inputs

| Prompt | Meaning |
|---|---|
| `[W]ire or (S)heet ?` | which gage system to use (default wire) |
| `Find (G)age or [S]ize ?` | look up a size from a designation, or the closest designation to a measured size (default size) |
| `Gage designation ?` | (size lookup) the gage number/designation to look up |
| `Thickness (in)` | (gage lookup) the measured thickness, default 0.1 |

## Output

Either the size for the given designation, or the closest
designation and its size for the given thickness.

## Method

Both the wire gage table (56 entries) and the sheet gage table (36
entries) are universal reference data, so they are stored in the
shared `reference.db` SQLite database (see
`ai/plans/c-to-go-conversion-plan.md`, "Data-file strategy for
Tier 2") rather than read from a flat file at startup, the same
mechanism `fits` and `speed` in the same conversion group use. Both
tables come from one `GAGE.DAT`, which stores the wire and sheet
tables one after another separated by a sentinel row whose
designation is literally `xx`; `MWKGo/tools/build-refdb`'s `seedGage`
splits on that sentinel into the two tables.

Exact-designation lookup is a direct match; closest-size lookup scans
every entry in the chosen table for the smallest absolute difference
from the given thickness — a plain linear scan, since the tables are
small (in the tens of rows) and the original program does the same.

`GAGE.DAT` is the one file in the whole Tier 2 batch with neither a
`STARTOFDATA` marker (its data starts at the first non-comment line)
nor a comma/tab-only delimiter scheme throughout (it also uses `;` as
an inline-comment cutoff on its sentinel row) — `internal/legacydat`
already supported both of these from earlier programs (`speed` for
the missing start marker, `fits` for the comma delimiter), so no
changes to the shared scanner were needed.

## Worked Example

No worked numeric example was included with the original program. As
independently verifiable checks, this conversion's tests confirm
`findClosest` always returns an entry whose distance is less than or
equal to every other entry's distance (the actual defining property
of "closest", not a restatement of the algorithm), and an integration
test against the real shipped `reference.db` confirms a well known
reference value: AWG wire gage 0000 is 0.4600 in.
