# colsort

Multi-column table sorting.

**Converted from:** `COLSORT.C` (M. W. Klotz), `MWKC/Misc/colsort.zip`
**Go source:** `MWKGo/colsort/colsort.go`

## Purpose

A generic utility for sorting a small comma-delimited table (up to 8
columns, up to 200 rows) on one chosen column, alphabetically or
numerically, increasing or decreasing — independent of any particular
calculator, useful for reordering any tabular data by hand.

## Data setup

Both the table itself and the sort job (how many columns, each
column's type, which column to sort on, and which direction) are
specific to one sorting job, not universal reference data or reusable
equipment configuration — so, like `calibrat`, `vrev`, and `simul` in
an earlier conversion group, this program reads its input fresh from a
file named on the command line each run, in the same
`STARTOFDATA`/`ENDOFDATA` text format the original used; see
`ai/plans/c-to-go-conversion-plan.md`, "Data-file strategy for
Tier 2".

```
colsort -data my-table.dat
```

Unlike the other ephemeral-input programs in this project, the sort
job's own parameters (column count, column types, sort column, sort
direction) are part of the data file itself, immediately before the
table rows — matching `COLSORT.DAT`'s own layout exactly. A worked
example built from the original archive's own shipped `COLSORT.DAT`
(30 materials sorted by decreasing density) ships at
`MWKGo/colsort/testdata/example.dat`.

## Inputs

None interactively; the entire job (data and parameters) comes from
the data file.

## Output

The table as read, then the same table after sorting.

## Method

The data file's header lines specify, in order: the number of columns
(1 to 8), one line per column giving its type (`a`=alphabetic,
`n`=numeric), which column to sort on (1-based), and the sort
direction (`i`=increasing, `d`=decreasing). Alphabetic columns compare
case-insensitively, with a shorter string that otherwise matches
sorting before a longer one that starts the same way (matching
`COLSORT.C`'s own `acomp` comparison exactly); numeric columns compare
by parsed value. `COLSORT.C`'s own bubble sort is already correct and
stable (unlike `diffthrd`'s, which never actually finished a second
pass), so this conversion uses `sort.SliceStable` with an equivalent
comparison rather than reproducing the bubble-sort mechanics — the
resulting order is identical either way.

## Worked Example

No separate `.TXT` worked example was shipped with the original
program; `COLSORT.DAT`'s own 30-material density table (sorted
decreasing by its second column) serves as this conversion's test
data. Tests confirm the header parses correctly, that the heaviest
entry (Zinc, 17.0) sorts first and the lightest (Carbon, 1.4) sorts
last under a decreasing sort, that the full result is monotonically
non-increasing, and — using a small synthetic table — that an
alphabetic sort is case-insensitive while still breaking a
case-insensitive tie by original row order (a stable sort).
