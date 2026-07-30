# fits

Shaft/hole fit computations.

**Converted from:** `FITS.C` (M. W. Klotz), `MWKC/WorkshopUtilities/fits.zip`
**Go source:** `MWKGo/fits/fits.go`

## Purpose

Machinists cut a shaft slightly larger or smaller than its nominal
size, depending on how tightly it needs to fit into a matching hole
(shrink and force fits are tight interference fits; running and
clearance fits leave increasing amounts of play). This calculator
looks up the standard allowance for a chosen fit class and reports
the shaft diameter to machine for a given nominal (hole) diameter.

## Inputs

| Prompt | Default |
|---|---|
| Fit desired (a number from the displayed list) | 4 (Push) |
| Nominal diameter of shaft (in) | 1.0 |

## Output

The shaft diameter to machine, to the nearest 0.00001 in.

## Method

Each fit class has a constant offset and a per-inch-of-diameter
allowance, both in thousandths of an inch; the shaft diameter is the
nominal diameter plus those two terms scaled by 0.001:

```
shaft diameter = nominal diameter + 0.001 * (allowance * nominal diameter + constant)
```

The eleven fit classes (shrink, force, drive, push, slide, precision
running, close running, normal running, easy running, small
clearance, large clearance) are universal reference data, not
specific to any one machine, so they are stored in the shared
`reference.db` SQLite database (see
`ai/plans/c-to-go-conversion-plan.md`, "Data-file strategy for
Tier 2") rather than read from a flat file at startup. `reference.db`
is built once, offline, by `MWKGo/tools/build-refdb` from the
original `FITS.DAT`, and its bytes are embedded in every calculator
that needs it via `internal/refdata`; at runtime this program just
queries the `fits` table.

The fits are listed, and selected by number, in the same order
`FITS.DAT` itself lists them — not alphabetically — since the
original program's own default selection (fit 4) depends on that
order, and the table's `list_position` column preserves it for
exactly that reason.

## Worked Example

`FITS.DAT`'s own header gives a worked example: "For a push fit on a
nominal 1 inch shaft, machine the hole to exactly 1.0000 inch, and
machine the shaft to -0.35\*(1.0)-0.15 = -0.5 thou less than the
nominal size (0.9995 inch)." This conversion's tests reproduce that
exact calculation, plus an independent identity check (a fit with
zero constant and zero allowance must leave the nominal diameter
unchanged) and an integration check against the real shipped
`reference.db` confirming "Push" is still fit number 4.

```
$ fits
SHAFT/HOLE FIT COMPUTATIONS

Number of data items read = 11
 1  Shrink
 2  Force
 3  Drive
 4  Push
 5  Slide
 6  Precision Running
 7  Close Running
 8  Normal Running
 9  Easy Running
10  Small Clearance
11  Large Clearance

Fit desired (default 4): 4
Nominal diameter of shaft (in) (default 1): 1.0

Diameter of shaft for Push Fit = 0.99950 in
```
