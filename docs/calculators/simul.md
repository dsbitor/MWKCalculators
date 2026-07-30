# simul

Simultaneous linear equation solver.

**Converted from:** `SIMUL.C` (M. W. Klotz), `MWKC/Math/simul.zip`
**Go source:** `MWKGo/simul/simul.go`
**Original documentation:** `SIMUL.TXT`, inside `MWKC/Math/simul.zip` (not included in this conversion)

## Purpose

Solves a system of N simultaneous linear equations in N unknowns,
given in standard canonical form (N coefficients and a constant term
per equation), by Gauss-Jordan elimination with full pivoting.

## Data setup

The equation system is specific to one problem the user wants solved
right now, not universal reference data or reusable equipment
configuration, so — like `calibrat` and `vrev` in the same conversion
group — this program reads its input fresh from a file named on the
command line each run, in the same `STARTOFDATA`/`ENDOFDATA` text
format the original program used; see
`ai/plans/c-to-go-conversion-plan.md`, "Data-file strategy for
Tier 2".

```
simul -equations my-system.dat
```

The file needs a line with the number of unknowns N, then N lines
each with N comma-separated coefficients followed by the equation's
constant term. A worked example built from `SIMUL.TXT`'s own 4x4
example ships at `MWKGo/simul/testdata/example.dat`.

## Inputs

None interactively; the entire problem comes from the data file.

## Output

Each equation as entered (for confirmation), then each unknown's
solved value, or `ERROR` if the system is singular (has no unique
solution, such as one equation being a multiple of another).

## Method

Gauss-Jordan elimination with full (row and column) pivoting — at
each step, the largest-magnitude remaining candidate entry in the
whole unpivoted submatrix is chosen as the pivot, which both improves
numerical stability and detects a singular system directly (no valid
pivot remains before every row has been used). This is the same
"gaussj" method described in *Numerical Recipes*, which the original
program itself is a direct port of.

## Worked Example

`SIMUL.TXT`'s own worked example (a 4-unknown system) is this
conversion's test: rather than asserting specific solution values,
the test substitutes the solved values back into each original
equation and confirms the residual is negligible — the actual
defining property of a correct solution, not a re-derivation of one.
A second test confirms `SIMUL.TXT`'s own example of a *non*-
independent system (`x + y = 3` and `2x + 2y = 6`, the second being
just the first doubled) is correctly detected as singular.
