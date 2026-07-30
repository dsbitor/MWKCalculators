# profile

Turned profile roughing schedule for an arbitrary, user-specified shape.

**Converted from:** `PROFILE.C` (M. W. Klotz), `MWKC/WorkshopUtilities/profile.zip`
**Go source:** `MWKGo/profile/profile.go`

## Purpose

A generalization of [ogive](ogive.md) and [egg](egg.md): rather than a
single closed-form shape, `profile` reads an arbitrary list of
`(x, radius)` points describing a workpiece's outline and computes a
roughing schedule of axial cutting passes with a square- or
round-tipped tool, the same way [ballcut](ballcut.md) (its
spherical-shape-only predecessor) does. Between data points, the shape
is filled in either by linear interpolation (the default) or, where
the data file says so, a natural cubic spline fit over a declared
range of points.

## Inputs

A `-data <file>` flag names a `STARTOFDATA`/`ENDOFDATA` text file
containing:

- Configuration lines (`ds`=stock diameter, `toolt`=1 for a
  square-tipped tool or 0 for round, `wt`=tool tip width or diameter,
  `dxd`=axial cutting increment, `scalex`/`scaler`=scale factors
  applied to every raw x/radius value below).
- A variable-length list of comma-separated `x,r` points (radius, not
  diameter, measured from the stock's outboard end, `x=0`).
- Zero or more `sseg=start,finish` lines, each declaring that the
  *sorted* point list's index range `start` to `finish` (inclusive)
  should be spline-fitted rather than linearly interpolated. Segments
  may not overlap.

An optional `-svg <file>` flag writes a rendering of the profile and
the cutting schedule.

## Output

The fixed configuration, then a schedule: each pass's axial position
(`c`), the resulting cut diameter (`d`), and the depth of cut (`doc`).
If `-svg` is given, a diagram of the profile with each roughing pass
overlaid, plus tick marks at every original data point (a larger tick
every tenth point).

## Method

For a given axial position, `findRadius` decides how to look up the
profile's radius there: if the position falls inside a declared
spline segment, a natural cubic spline (Sedgewick's tridiagonal
method, the same fitter [spline](spline.md) uses) is fitted fresh over
just that segment's own points and evaluated; otherwise the two
bracketing raw data points are linearly interpolated. The cutting
schedule itself follows the same square/round-tool pattern as
[ogive](ogive.md) and [egg](egg.md): for a square tool, depth is set
by whichever of the tool's two edges or center cuts deepest; for a
round tool, the tool's own curved tip is scanned every 15° against the
target profile.

`PROFILE.C`'s own axial-position loop (`xl -= dxd` each pass) drifts
with repeated floating-point subtraction; this conversion computes
each pass's position as a multiple of the axial step instead
(`-(i-1)*dxd`), which avoids that drift and, on this shipped example,
happens to reproduce the exact number of passes (76) the original
produced, rather than one fewer.

Per `ai/plans/c-to-go-conversion-plan.md`'s Tier 3 "Graphics scope"
resolution, the original's interactive mouse-driven menu and
click-to-read-coordinates feature are dropped entirely (no static SVG
equivalent exists), and its `.OUT` file-save feature is replaced by
printing the schedule straight to stdout.

## Worked Example

`PROFILE.DAT`'s own shipped example (a tiny model-tool handle) is used
as this conversion's primary test oracle, checked against a genuine
`PROFILE.OUT` captured from the original DOS binary. Every row matches
to the printed digit except depth-of-cut on the very first row: the
true value there is exactly `0.0625` (half the stock diameter), an
exact decimal tie at the third digit, and Go's `fmt` (round-half-to-
even) and the original DOS `printf` (round-half-away-from-zero) round
that single tied value differently ("0.062" vs "0.063") — a display
rounding-mode difference, not a computation error. Tests also confirm
`findRadius` reproduces every one of the data file's own points
exactly, whether inside a spline segment or one of the linearly-
interpolated gaps between them.
