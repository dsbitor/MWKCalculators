# ogive

Tangent ogive nose-cone profile and turning schedule.

**Converted from:** `OGIVE.C` (M. W. Klotz), `MWKC/WorkshopUtilities/ogive.zip`
**Go source:** `MWKGo/ogive/ogive.go`

## Purpose

An ogive is the pointed-arch shape formed by two intersecting
circular arcs — the same shape used in both Gothic architecture and
rocket nose cones. Given the ogive's own diameter and length, this
computes its exact profile and works out an incremental turning
schedule for roughing it out on a lathe with either a square-tipped
or round-tipped tool (the same "rough it out, then finish by filing
and sanding" workflow `ogive`'s archive-mate `profile` uses for
arbitrary tool-handle shapes).

## Data setup

The ogive's own dimensions and tooling setup are specific to one
job — like `calibrat`, `vrev`, `simul`, `curfit`, `colsort`, `spline`,
and `belt` in earlier conversion groups — this program reads its
input fresh from a file named on the command line each run, in the
same `STARTOFDATA`/`ENDOFDATA` text format the original used, though
with a `key=value` syntax rather than positional fields:

```
ogive -data my-ogive.dat [-svg my-ogive.svg]
```

The file needs `ds` (stock diameter), `toolt` (1 for a square-tipped
tool, 0 for round-tipped), `wt` (tool tip width for a square tool, or
diameter for a round one), `dxd` (axial cutting increment), `ogdiam`
(ogive diameter), and `oglength` (ogive length). Unlike most other
Tier 2/3 data files, `OGIVE.C` never prompts interactively at
all — every parameter comes from this one file. A worked example
built from the original archive's own shipped `OGIVE.DAT` ships at
`MWKGo/ogive/testdata/example.dat`.

## Inputs

None interactively; the entire job comes from the data file.

## Output

The configuration as read, then a table of the roughing schedule:
each pass's axial position (measured from the tip), the resulting cut
diameter, and its depth of cut. If `-svg` is given, a diagram of the
ogive profile with each roughing pass overlaid (a rectangle for a
square tool, a circle for a round one).

The original's interactive mouse-driven menu and click-to-read-
coordinates feature are dropped entirely (no equivalent in a static
SVG image), and its `.OUT` file-save feature is replaced by printing
straight to stdout; see
`ai/plans/c-to-go-conversion-plan.md`'s Tier 3 "Graphics scope"
resolution for the general policy this follows.

## Method

The ogive profile's radius at distance `x` from the tip is
`√((d×(c²+0.25))² − (L−x)²) − d×(c²−0.25)`, where `L` is the ogive
length, `d` its diameter, and `c = L/d` its caliber — the standard
tangent-ogive formula. Beyond the ogive's own length, the profile has
already blended into the stock cylinder, so the radius is simply the
stock radius there.

At each axial cutting position, a square-tipped tool's required depth
is set by whichever of its two edges (or its center) would cut
deepest into the ogive shape; a round-tipped tool's depth is instead
found by scanning its own curved tip profile (sampled every 15°
around its leading half) against the ogive shape, since a round
tool's own curvature means no single point on it is necessarily the
deepest.

## Worked Example

No worked numeric example with expected output was included with the
original program (`OGIVE.TXT` only explains the shape and points to
an external mathematical reference). As independently verifiable
checks, this conversion's tests confirm two defining geometric
properties of a tangent ogive using the shipped example's own
dimensions (0.5 diameter, 1.5 length, caliber 3): the radius at the
very tip is exactly zero (a sharp point), and the radius at the
ogive's own length exactly equals the stock radius (a smooth blend
into the cylindrical stock, not a step). Further tests confirm the
profile widens monotonically from tip to base, and that the roughing
schedule's own cut diameters are correspondingly non-decreasing,
ending near the full stock diameter.
