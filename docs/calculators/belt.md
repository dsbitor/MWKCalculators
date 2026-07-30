# belt

Flat belt length for an arbitrary arrangement of pulleys.

**Converted from:** `BELT.C` (M. W. Klotz), `MWKC/WorkshopUtilities/belt.zip`
**Go source:** `MWKGo/belt/belt.go`

## Purpose

Given the locations, diameters, and belt-side (inside the loop, like a
normal driver/driven pulley, or outside it, like an idler) of any
number of pulleys arranged in a closed loop, computes the total length
of flat belt needed to connect them all, and warns if any pulley
would end up with a slack (unsupported) belt wrap.

This archive also bundles three companion programs for simpler,
two-pulley-only problems: [qbelt](qbelt.md) (quick belt length, no
data file needed), [pulley](pulley.md) (solve for the second pulley's
diameter given belt length and separation), and [pcd](pcd.md) (solve
for the pulley separation given both diameters and belt length).

## Data setup

A pulley layout is specific to one job's own arrangement — not
universal reference data, and not reusable equipment configuration —
so, like `calibrat`, `vrev`, `simul`, `curfit`, and `colsort` in
earlier conversion groups, this program reads its input fresh from a
file named on the command line each run, in the same
`STARTOFDATA`/`ENDOFDATA` text format the original used; see
`ai/plans/c-to-go-conversion-plan.md`, "Data-file strategy for
Tier 2".

```
belt -data my-pulleys.dat
```

Each data line is `x,y,diam,flag`, one per pulley, entered in
**counter-clockwise order around the belt path** (unlike most other
Tier 2 data files, pulley order is significant and is not sorted).
`flag` is `1` if the pulley sits inside the belt loop (a normal
driver/driven pulley), or `-1` if it sits outside (an idler, wrapped
on its far side). A worked example built from the original archive's
own shipped `BELT.DAT` (a 6-pulley layout the original's own comment
notes should have one slack pulley) ships at
`MWKGo/belt/testdata/example.dat`.

## Inputs

None interactively; the entire pulley layout comes from the data
file.

## Output

For each pulley: its own data as entered, its wrap angle (degrees),
wrap length, and the straight belt span to the next pulley in the
loop. If any pulley's wrap length comes out negative, a note that the
belt may go slack there. Finally, the total belt length.

## Method

Between each adjacent pair of pulleys, the straight tangent-line span
is found from the center-to-center distance and each pulley's radius:
pulleys on the *same* side of the belt loop (both normal, or both
idlers) use the internal-tangent formula (the belt crosses between
them), while pulleys on *opposite* sides use the external-tangent
formula. If the required tangent offset exceeds the center-to-center
distance, or the two centers coincide, the pulleys physically overlap
and no answer is possible. Each pulley's wrap angle is then the
angular difference between its incoming and outgoing tangent points;
if that difference indicates a wrap of more than half a turn, an
alternate ("the short way around") solution is checked, and adopted
(reported as a negative wrap length) if it is the physically valid
one for the given span — the same condition `BELT.C` itself uses to
flag a pulley as likely slack.

## Worked Example

No worked total-belt-length figure was included with the original
program's own text files, but `BELT.DAT`'s own 6-pulley example
carries an inline comment flagging its last pulley as one that
"should be slack" — an author-supplied expected result this
conversion's tests confirm directly: pulley 6 comes back with a
negative wrap length and the overall slack flag set, while every
other pulley's wrap length is non-negative. A second test cross-checks
`computeBeltLayout`'s general N-pulley algorithm against the simpler,
independently-formulated closed-form two-pulley external-tangent
calculation `qbelt`/`pulley`/`pcd` all share, confirming both
approaches agree on the same total length for a plain two-pulley
layout. A third test confirms two overlapping pulleys are correctly
rejected.
