# pulley

Solve for an unknown pulley diameter given belt length and separation.

**Converted from:** `PULLEY.C` (M. W. Klotz), `MWKC/WorkshopUtilities/belt.zip`
**Go source:** `MWKGo/pulley/pulley.go`

## Purpose

Given one pulley's known diameter, the separation between two pulley
centers, and a fixed belt length (for example, a belt you already
own), finds the diameter the *second* pulley needs to be for that
belt to fit exactly — the reverse of the usual belt-length
calculation `qbelt` performs.

## Inputs

| Prompt | Default |
|---|---|
| Diameter of known pulley | 1.4 |
| Separation between pulley centers | 2.5 |
| Belt length | 8.21 |
| Calculation accuracy | 0.001 |

(The defaults are `BELT.DAT`'s own "conical pulley" worked example —
see below.)

## Output

For the known pulley and the calculated pulley: diameter, wrap angle
(degrees), and wrap length. Then the belt span between them and the
total belt length (which matches the requested belt length within the
chosen accuracy).

## Method

The second pulley's diameter is searched for between a tenth and ten
times the known pulley's own diameter, using [qbelt](qbelt.md)'s own
two-pulley belt-length formula evaluated at each trial diameter: a
coarse 10-step scan finds where the belt length crosses the target
length (going from too short to too long), then each pass narrows the
search to that bracket and repeats, converging geometrically until
the result is within the requested accuracy. `coding-style.md` Rule 2
replaces the original's interactive keypress abort with an explicit
pass-count cap (in practice, convergence takes only a handful of
passes).

## Worked Example

`PULLEY.C`'s own default prompt values are themselves `BELT.DAT`'s
own "conical pulley" example: a 1.4-diameter pulley at a 2.5
separation with a target belt length of 8.21. This conversion's tests
confirm that scenario converges to a second-pulley diameter close to
0.603 — matching `BELT.DAT`'s own paired data row
(`2.5,0,0.603,1`) exactly — and that the found diameter's own
reconstructed belt length matches the requested 8.21 within the
chosen accuracy.
