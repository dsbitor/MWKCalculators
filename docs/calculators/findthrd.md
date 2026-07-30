# findthrd

Find a thread from its diameter or pitch.

**Converted from:** `FINDTHRD.C` (M. W. Klotz), `MWKC/WorkshopUtilities/findthrd.zip`
**Go source:** `MWKGo/findthrd/findthrd.go`

## Purpose

Andy Pugh (Sheffield, UK) compiled a table of roughly 470 standard
threads — British, metric, and a number of historic and specialist
standards (Whitworth, BSF, BSP, BA, ISO metric, Loewenherz, Thury,
watch and instrument threads, and others) — arranged by size and
identified by name. This calculator searches that table by measured
major diameter or pitch (in either inches/tpi or millimeters) to
identify a "mystery" thread, within a chosen tolerance.

## Inputs

A menu (A-D search by diameter-in, diameter-mm, pitch-tpi, or
pitch-mm; M redisplays the menu; Q quits), then per search:

| Prompt | Default |
|---|---|
| Thread diameter (in) | 0.138 |
| Thread diameter (mm) | 6.0 |
| Thread pitch (tpi) | 32.0 |
| Thread pitch (mm) | 1.0 |
| Allowable error (%) | 1.0 |

## Output

Every thread in the reference table whose chosen field is within
tolerance of the entered value: name, diameter (in and mm), and
pitch (tpi and mm).

## Method

The ~470 threads are universal reference data (a published,
third-party-compiled table), so they live in the shared
`reference.db` SQLite database, the same mechanism used by `fits`,
`speed`, `gage`, and `expand` earlier in this project; see
`ai/plans/c-to-go-conversion-plan.md`, "Data-file strategy for
Tier 2".

`FINDTHRD.DAT` is the one file in the whole Tier 2 batch using
fixed-width columns instead of a delimiter (name in the first 12
characters, then diameter-in, diameter-mm, pitch-tpi, and pitch-mm
each in their own fixed-width field), reflecting a tabular layout
meant to be read by eye as much as by the original program. Two more
columns (core diameter and thread depth) follow in the full-width
rows, which the original's own `struct threads` reads into memory but
never actually uses anywhere in the program; this conversion does not
carry them over either. A handful of entries are shorter lines that
simply stop after the pitch-mm column, which `MWKGo/tools/build-refdb`
treats as those columns being present but the trailing ones absent —
not as an error, and deliberately not as whatever leftover bytes the
original's own reused, non-cleared line buffer would have copied from
the previous line's data for a `strncpy` reading past a short C
string's end (undefined behavior no Go conversion should try to
reproduce).

A match brackets the *search value* between a table entry scaled by
`1 - tolerance%` and `1 + tolerance%` — scaling the table's own value,
not the search value — matching `FINDTHRD.C`'s exact comparison
(`d >= a*th[i].diam && d <= b*th[i].diam`); the two are only
approximately the same thing for small tolerances, so this
conversion's tests check the exact boundary the original's formula
produces rather than a symmetric approximation of it.

## Worked Example

No worked numeric example was included with the original program. As
independently verifiable checks, this conversion's tests confirm a
search exactly at a table-scaled tolerance boundary matches and one
just outside it does not (the defining property of the comparison
formula), and an integration test against the real shipped
`reference.db` confirms a well known standard thread (1/4 UNEF: 0.25
in diameter, 32 tpi) is found by a diameter search.
