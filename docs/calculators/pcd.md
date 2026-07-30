# pcd

Solve for pulley center distance given both diameters and belt length.

**Converted from:** `PCD.C` (M. W. Klotz), `MWKC/WorkshopUtilities/belt.zip`
**Go source:** `MWKGo/pcd/pcd.go`

## Purpose

Given both pulley diameters and a fixed belt length, finds the
separation between the pulley centers that makes the belt fit exactly
— the complement of [pulley](pulley.md), which instead solves for an
unknown diameter when the separation is already fixed. Added later
than `BELT`, `QBELT`, and `PULLEY` (2005 versus the others' 2000), but
it shares the same underlying two-pulley formula and search structure.

## Inputs

| Prompt | Default |
|---|---|
| Diameter of driver pulley | 1.4 |
| Diameter of driven pulley | 0.603 |
| Belt length | 8.21 |
| Calculation accuracy | 0.0001 |

(The defaults are `BELT.DAT`'s own "conical pulley" worked example —
see below.)

## Output

For the driver and driven pulleys: diameter, wrap angle (degrees),
and wrap length. Then the belt span, the total belt length (matching
the requested belt length within the chosen accuracy), and the
pulley center distance found.

## Method

Identical in structure to [pulley](pulley.md)'s own search, but
searching the pulley separation instead of a diameter: a coarse
10-step scan between a tenth and ten times the target belt length
finds where the computed belt length crosses the target, then each
pass narrows to that bracket and repeats until within the requested
accuracy. `coding-style.md` Rule 2 replaces the original's
interactive keypress abort with an explicit pass-count cap.

## Worked Example

`PCD.C`'s own default prompt values are `BELT.DAT`'s own "conical
pulley" example: diameters 1.4 and 0.603 with a target belt length of
8.21. This conversion's tests confirm that scenario converges to a
separation close to 2.5 — matching `BELT.DAT`'s own data exactly —
and, going further, that all three of `BELT.DAT`'s *other* documented
conical-pulley diameter pairs (1.5/0.477, 1.6/0.343, 1.7/0.200, each
commented as sized for the same fixed belt) also converge to that
same 2.5 separation.
