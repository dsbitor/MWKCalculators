# xymwk

Height-gauge coordinate transforms and measurements: reference,
align, distance, and pitch circle.

**Converted from:** `XYMWK.C` (M. W. Klotz), `MWKC/Math/xymwk.zip`
**Go source:** `MWKGo/xymwk/xymwk.go`

## Purpose

Clamp a part to a precision 90-degree angle plate, measure feature
heights with a height gauge to get one coordinate, rotate 90 degrees
and measure again to get the other, and you have a
raw x,y coordinate for each feature — referenced to the plate, not to
whatever is actually convenient for checking the part. This program
takes that raw list of points and performs one of four operations on
it:

- **reference** — translate the whole point set so a chosen point
  becomes the new origin.
- **align** — translate so a first chosen point becomes the origin,
  then rotate so a second chosen point lies on the (positive) x-axis.
- **distance** — report the distance between two chosen points.
- **pitchcircle** — report the center and radius of the circle passing
  through three chosen points.

Unlike the other seven Tier 3 programs, `xymwk` is not a curve/profile
calculator with a roughing schedule — it is a coordinate-transform and
measurement tool, and its scope was discussed and decided separately
(see `ai/plans/c-to-go-conversion-plan.md`'s Tier 3 Group 3 section)
before converting it.

`XYMWK.C` selects its two or three working points by mouse click,
optionally "snapping" to the nearest data point. This conversion
selects points by their 1-based position in the data file instead —
at least as precise as snapping, and needs no interactive session.
Per that same scope discussion, each run performs exactly **one**
operation on the raw data; it does not replicate the original's
interactive session, where each operation's result became the new
working coordinate set for the next operation (so, unlike the
original, running `align` and then `distance` in two separate
invocations measures distance in the *original* frame, not the
aligned one).

## Inputs

A `-data <file>` flag names a `STARTOFDATA`/`ENDOFDATA` text file of
comma-separated `x,y` points, one per line, in the same order they'll
be selected by index — unlike every other program in this tier, point
order is not cosmetic here, since it's how a point is chosen.

An `-op` flag selects the operation (`reference`, `align`, `distance`,
or `pitchcircle`), and `-p1`/`-p2`/`-p3` name the 1-based point
indices it needs: one for `reference`, two for `align` and `distance`,
three for `pitchcircle`.

## Output

For `reference` and `align`, a table of every point: its original
coordinates, its transformed coordinates, its radius from the new
origin, and its angle (both decimal degrees and degree:minute:second
notation). For `distance`, a single distance value. For `pitchcircle`,
the circle's center and radius.

## Method

`reference` and `align` are plain translation and translation-plus-
rotation. `pitchcircle` finds the circle through three points via the
perpendicular bisectors of two of the triangle's chords (the same
construction as its original `circ3`/`bisect`/`loc`/`intersect`
functions).

`XYMWK.C` declares and fully implements `display()` — a function that
prints exactly the table this conversion's `reference`/`align` print
(original and transformed x/y, radius, angle in both notations) — and
the original program's own accompanying documentation describes
`reference` and `align` as showing exactly this table, matching what
`display()` would show. But `display()` is never actually called
anywhere in `main()`: in the real program, `reference` and `align`
only silently re-plot the
points on screen with no printed values at all, and only `distance`
and `pitchcircle` draw short text results onto the graphics screen (no
output file exists for this program at all, unlike most others in this
project). Since a static CLI has no equivalent for "re-plot the points
for the user to look at," this conversion wires up `display()`'s table
for `reference`/`align` — not inventing new behavior, but connecting
code that already matches the tool's own documented intent. Two other
declared-but-never-called functions, `angle()` and `project()`, are
omitted entirely as dead code, the same way [egg](egg.md) omitted its
own unused copy of `spline`'s cubic-spline fitter.

`XYMWK.C`'s own comment on `circ3` flags a known gap it never fixes:
"no error checking is done for the anomalous case in which the three
points are colinear" — and even suggests how: "can be done by aligning
to points 1 and 2; if y coordinate of point 3 is zero after said
alignment, points are colinear." This conversion implements exactly
that check (reusing its own `align` code), rather than relying on
`intersect`'s own exact `d==0.0` equality test the way the original
did: two bisector directions computed through trigonometric functions
are essentially never bit-for-bit parallel even when the three points
are mathematically colinear, so the original's unchecked call could
return a wild, wrong circle instead of visibly failing.

## Worked Example

`XYMWK.DAT`'s own shipped example turns out to be a handy,
hand-checkable dataset: its five points are a center point and four
points each exactly 50 units from it (`193.301 = 150 + 50*cos(30°)`,
and `(180,140)` relative to `(150,100)` is a 3-4-5 right triangle
scaled by ten). This conversion's tests confirm `reference`d-to-center
points land at radius exactly 50; that `align`ing to two points puts
the second exactly on the x-axis; that the circumcircle through three
of the four outer points reproduces the known center and radius
exactly; and, as an independent check not tied to the known center,
that any circumcircle's own three defining points are equidistant from
its computed center. A separate test confirms three collinear points
now report an error rather than a wrong circle.
