# speed

Machining speed (recommended spindle RPM) utility.

**Converted from:** `SPEED.C` (M. W. Klotz), `MWKC/WorkshopUtilities/speed.zip`
**Go source:** `MWKGo/speed/speed.go`

## Purpose

Every material has a recommended cutting speed range, expressed as
surface feet per minute (SFPM) — the speed at which the tool's
cutting edge should pass over the material's surface, regardless of
the tool or workpiece diameter. This calculator looks up a chosen
material's recommended SFPM range and converts it to a recommended
spindle RPM range for a given tool (or workpiece) diameter.

## Inputs

| Prompt | Default |
|---|---|
| Material Number (a number from the displayed list) | 1 (Aluminum and Alloys) |
| Diameter (in) — of the workpiece if turning on a lathe, or of the cutter/drill if milling or drilling | 1.0 |

## Output

The recommended spindle RPM range for the chosen material and
diameter.

## Method

Surface speed relates to diameter and RPM by `surface speed = pi *
diameter * rpm / 12` (diameter in inches, surface speed in feet per
minute); solving for RPM gives `rpm = 12 * sfpm / (pi * diameter)`,
applied to both ends of the material's recommended SFPM range.

The thirteen materials (aluminum, brass, several grades of steel and
cast iron, copper, bronze) are universal reference data, so — like
`fits` in the same conversion group — they live in the shared
`reference.db` SQLite database rather than a flat file read at
startup; see `ai/plans/c-to-go-conversion-plan.md`, "Data-file
strategy for Tier 2", and `fits.md` for the shared mechanism
(`MWKGo/tools/build-refdb`, `internal/refdata`).

Materials are listed, and selected by number, in `SPEED.DAT`'s own
file order rather than alphabetically, the same reasoning as `fits`:
the original program's default selection (material 1) depends on
that order. `SPEED.DAT` itself has no `STARTOFDATA` marker (unlike
most of this batch's `.DAT` files) — its data region simply starts at
the first non-comment line, which `internal/legacydat.Rows` supports
directly by treating an empty start marker as "start immediately".

## Worked Example

No worked numeric example was included with the original program.
As an independently verifiable check, this conversion's tests recover
the original SFPM range from `rpmRange`'s output via the inverse of
the same formula (`sfpm = pi * diameter * rpm / 12`) rather than
comparing against hand-typed numbers, and separately check that
halving the diameter exactly doubles the recommended RPM at constant
surface speed — both properties of the formula itself, not of any one
material's numbers. An integration test against the real shipped
`reference.db` confirms "ALUMINUM AND ALLOYS" is still material
number 1.

```
$ speed
MACHINING SPEED UTILITY

number of data entries read = 13

   1  ALUMINUM AND ALLOYS
   2  BRASS AND SOFT BRONZE
   ...

Material Number (default 1): 1

Input relevant diameter.
Diameter of work if lathe, diameter of cutter/drill if mill/drillpress.
Diameter (in) (default 1): 1.0

For ALUMINUM AND ALLOYS and diameter 1.0000 in:
Recommended rpm between 764 and 1146
```
