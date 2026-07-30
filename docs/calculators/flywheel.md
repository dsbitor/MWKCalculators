# flywheel

Tapered spoke flywheel calculations.

**Converted from:** `FLYWHEEL.C`, `MWKC/WorkshopUtilities/flywheel.zip`
**Go source:** `MWKGo/flywheel/flywheel.go`

## Purpose

Machining a tapered-spoke flywheel from solid stock starts with a
reduced-thickness web between the hub and rim, then holes drilled
through the web define the corners of cutouts; milling away the
material between those holes leaves the spokes. Given the number of
spokes, the radii and diameters of the inner and outer hole circles,
and the desired spoke taper (specified as either the outer hole's
offset from the spoke centerline or its angle from the centerline),
this program computes the rotary table angle needed to mill along a
spoke's edge (found by solving a transcendental equation via
iterative search), the table offset, spoke widths, and the full table
of rotary settings for drilling every hole.

## Inputs

| Prompt | Default |
|---|---|
| Number of spokes | 6 |
| Radius on which inner holes are located | 0.440 in |
| Radius on which outer holes are located | 1.367 in |
| Diameter of inner holes | 0.1875 in |
| Diameter of outer holes | 0.1875 in |
| Offset from spoke CL to outer hole center | 0 in |
| Angle from spoke CL to outer hole center (if offset is 0) | 7 deg |

Exactly one of the offset or angle must be given; the header
comment's own suggested demo input is offset 0 with a 7 degree angle.

## Output

Both an exact solution and a "nearest integral angle" variant (which
rounds the outer hole's angle and the table rotation angle phi to
whole degrees and re-solves the inner hole radius to match, trading
slightly different hole radii for cleaner rotary table settings):
each includes phi (rotary table angle to bring the spoke edge
parallel to the mill's y axis), the table offset, spoke widths, the
minimum web length, and the full sequence of rotary table settings
for drilling the inner and outer holes.

## Method

The key angle phi solves a transcendental equation (no closed form):
the inner and outer hole edges must align along the same line. This
conversion uses the same zoom-in grid search technique as `cseg` in
Tier 1 group 11: scan a window for the angle minimizing the
equation's residual, then repeat on a shrinking window around the
best candidate, narrowing tenfold each pass, until the found angle
stabilizes (bounded to a generous number of refinement passes per
`coding-style.md` Rule 2, since the original relies on convergence
with no iteration cap of its own).

## Worked Example

No worked numeric example was included with the original program
(`FLYWHEEL.TXT` explains the geometry and mentions that "the first
part of the solution obtained with the defaults" produces messy
angles, but includes no sample run). As an independently verifiable
check: whatever phi the search returns must make the transcendental
equation's own residual approximately zero — the defining property
of the root search, not a re-run of it — confirmed for the
documented demo input (offset 0, theta2 = 7 degrees) in this
conversion's tests, along with a check that the "nearest integral
angle" solution's own theta2 and phi are genuinely whole numbers.
