# diam

Machining diameter utility.

**Converted from:** `DIAM.C` (M. W. Klotz), `MWKC/WorkshopUtilities/diam.zip`
**Go source:** `MWKGo/diam/diam.go`

## Purpose

`speed` (earlier in this conversion group) starts from a diameter and
finds a recommended spindle rpm range for a material. `diam` runs the
same relationship the other way: given a material and the actual
spindle speeds your own machine offers (a lathe or mill only has so
many gear or pulley steps, not a continuously variable speed), it
reports the range of tool or workpiece diameters that keep the
cutting speed within that material's recommended range, at each
available speed step.

## Data setup

`diam` is the first program in this batch to need both databases at
once. Its material table is the same universal `machining_speeds`
table `speed` reads — the original `DIAM.DAT` literally repeats
`SPEED.DAT`'s own material list verbatim, so this conversion does not
duplicate that data into a second table; it just queries the table
`speed`'s conversion already created in the shared `reference.db`.
The list of available spindle speeds, though, is specific to one
machine, so — like `diffthrd` and `divhead` earlier in this group —
that part reads from the user's own `userdata.db` instead; see
`ai/plans/c-to-go-conversion-plan.md`, "Data-file strategy for
Tier 2". That table ships empty. Import your own machine's available
speeds before first use:

```
diam -import-speeds my-machine-speeds.csv
```

The CSV needs one column, `rpm`. A worked example built from the
original `DIAM.DAT`'s own default speed list (twelve steps, 85 to
2225 rpm) ships at `MWKGo/diam/testdata/example-speeds.csv`.

## Inputs

| Prompt | Default |
|---|---|
| Material Number (from the displayed list) | 1 (Aluminum and Alloys) |

## Output

A table of the machine's available spindle speeds, each with the
diameter range that keeps the cutting speed within the chosen
material's recommended range at that speed.

## Method

The same surface-speed relation `speed` uses, `rpm = 12 * sfpm / (pi
* diameter)`, solved for diameter instead of rpm: `diameter = 12 *
sfpm / (pi * rpm)`, applied to both ends of the material's SFPM range
at each available spindle speed.

## Worked Example

No worked numeric example was included with the original program. As
an independently verifiable check, this conversion's tests confirm
`diameterRange`'s output satisfies the surface-speed formula in
reverse (recomputing sfpm from the returned diameter and a given rpm
reproduces the material's own sfpm range) and that feeding a
`diameterRange` result back through the same relation reproduces the
original rpm — the two functions are inverses of the same formula, so
each is used to check the other rather than against hand-typed
numbers. An integration test confirms `diam` and `speed` really do
read the identical `machining_speeds` table (material 1 is "ALUMINUM
AND ALLOYS" in both).
