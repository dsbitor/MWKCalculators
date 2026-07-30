# spline

Natural cubic spline curve fitting.

**Converted from:** `SPLINE.C` (M. W. Klotz), `MWKC/Math/spline.zip`
**Go source:** `MWKGo/spline/spline.go`

## Purpose

Fits a natural cubic spline through a set of (x,y) points (Sedgewick's
tridiagonal method, "Algorithms" p.545) and prints a sampled table of
the fitted curve; can also render the curve and the original data
points as an SVG diagram.

This is the second of Tier 3's eight graphics-bearing calculators; see
`ai/plans/c-to-go-conversion-plan.md`'s "Graphics scope for the eight
Tier 3 programs" for the resolved design (`internal/svgplot`, a small
package mirroring the original's line/box/circle/text drawing
primitives, emitting SVG instead of drawing to a screen). Unusually
among the eight, `SPLINE.C` calls the underlying DOS primitives
(`_moveto_w`/`_lineto_w`) directly rather than through the shared
`wline`-style wrapper every other Tier 3 program uses — the strongest
direct evidence available anywhere in the source tree for what that
wrapper actually does, which is what let this project's survey work
confirm the `wline` → line-segment mapping used throughout
`internal/svgplot`.

## Data setup

The data points are specific to one curve-fitting job — `SPLINE.DAT`'s
own shipped example is explicitly a demonstration dataset for the
algorithm itself, not reference data or reusable configuration — so,
like `calibrat`, `vrev`, `simul`, `curfit`, and `colsort` in earlier
conversion groups, this program reads its input fresh from a file
named on the command line each run, in the same
`STARTOFDATA`/`ENDOFDATA` text format the original used.

```
spline -data my-points.dat [-svg my-curve.svg]
```

Each data line is `x,y`. Points need not be pre-sorted — the shipped
example's own first data row is deliberately out of order, to
demonstrate this. A worked example built from the original archive's
own shipped `SPLINE.DAT` (4 points) ships at
`MWKGo/spline/testdata/example.dat`.

## Inputs

None interactively; the entire dataset comes from the data file. The
original's interactive mouse-driven menu (exit / toggle menu / redraw,
plus a live x,y coordinate readout following the mouse cursor) has no
equivalent in a static SVG output and is dropped entirely, not merely
adapted — a materially different feature loss than the color/pixel
rendering change every Tier 3 program shares, so it's called out
explicitly here per the plan's own documentation requirement.

## Output

A table of 41 evenly-spaced sample points along the fitted curve
(matching the original's own 40-segment sampling). If `-svg` is
given, a diagram: a dotted x/y axis pair, the fitted curve (solid,
color 14), and the original data points marked with small crosses.

## Method

The spline's shape is controlled by a set of per-point coefficients
(informally, how much the curve bends at each knot), found by solving
a tridiagonal system built from the points' own x-spacing and
y-differences, with "natural" boundary conditions (zero curvature) at
both ends — the same textbook method `profile` also uses internally
for its own spline-fit segments (see `docs/calculators/profile.md`).
Evaluating the fitted curve at any x within the data's own range
blends the two bracketing points' y values with a cubic correction
term weighted by each point's own coefficient.

## Worked Example

No numeric worked example is available, since `SPLINE.C` never
prints anything itself — it only draws. As
independently verifiable checks, this conversion's tests confirm the
shipped example's deliberately-out-of-order first point sorts into
its correct position; that the fitted curve passes through every
original data point exactly (the defining property of interpolation,
not a re-derivation of one); that evaluating outside the data's own
range is correctly rejected; and, as a strong sanity check
independent of any specific dataset, that fitting a natural cubic
spline through perfectly collinear points reproduces the underlying
straight line exactly at several intermediate x values — not merely
close to it, since a natural cubic spline of collinear data should
have zero curvature everywhere, matching a straight line without
approximation error.
