# cam

Plate cam profile design.

**Converted from:** `CAM.C` (M. W. Klotz), `MWKC/WorkshopUtilities/cam.zip`
**Go source:** `MWKGo/cam/cam.go`
**Original documentation:** `CAM.TXT`, inside `MWKC/WorkshopUtilities/cam.zip` (not included in this conversion)

## Purpose

Designs a plate cam profile for a chosen follower motion law and
follower type, printing the cam's own profile as a table of angle,
radius, and (x,y) coordinates, the resulting maximum pressure angle,
and a table suggesting how the base circle radius trades off against
pressure angle. Can also render the cam profile and the follower's
displacement, velocity, and acceleration curves as an SVG diagram.

This is the first of Tier 3's eight graphics-bearing calculators; see
`ai/plans/c-to-go-conversion-plan.md`'s "Graphics scope for the eight
Tier 3 programs" for the resolved design (a small internal SVG
package, `internal/svgplot`, mirroring the original's line/box/circle/
text drawing primitives). The original's interactive "press a key to
continue" pause and its `.OUT` file-save feature are both dropped: the
numeric tables print straight to stdout, and the diagram is written
via an optional `-svg <path>` flag instead.

## Inputs

| Prompt | Default |
|---|---|
| Follower motion type (1 straight-line, 2 parabolic, 3 simple harmonic, 4 cycloidal) | 4 |
| Follower type (1 flat-faced, 2 roller) | 1 |
| Radius of roller (roller follower only) | 0.25 |
| Base circle radius | 1.625 |
| Cam rise (motion types 2-4) | 1.25 |
| Cam rotation angle, deg (motion types 2-4) | 120 |
| Allowable cam rise during acceleration/linear/deceleration (motion type 1 only) | 0.1 / 1.0 / 0.15 |
| Cam rotation during linear motion, deg (motion type 1 only) | computed from the rise proportions |
| Angular step size for constructing the cam, deg | 1 |

No units are enforced anywhere except angles (always degrees) — use
any consistent system of units.

## Output

The cam profile table (angle from the highest point, radius and (x,y)
from the center of rotation, one row per angular step), the maximum
pressure angle found along that profile, and a table showing what
base circle radius would be needed to hold the pressure angle to each
of several target values from 15° to 45°.

If `-svg` is given, a diagram: the red base circle, the yellow cam
profile, and the follower's green displacement, blue velocity, and
magenta acceleration curves, all plotted against a shared synthetic
rotation-angle axis.

## Method

Four follower motion laws are supported, each giving displacement,
velocity, and acceleration as a function of how far through the rise
the cam has rotated: straight-line motion with parabolic
acceleration/deceleration blending at each end (to avoid the
infinite acceleration a literal constant-velocity cam would have),
plain parabolic (constant acceleration, then constant deceleration —
the simplest law, but with a sharp corner in the velocity curve),
simple harmonic (`d = rise/2 × (1 − cos(π×t))`), and cycloidal (`d =
rise × (t − sin(2πt)/(2π))`, the recommended law: no sharp corners in
either velocity or acceleration).

The cam's own profile radius at each angle comes from
`√((base+displacement+rollerRadius)² + velocity²) − rollerRadius`
(the roller-follower correction folds neatly into the same formula
whether or not a roller is actually in use, since `rollerRadius=0`
for a flat-faced follower is a no-op), and the pressure angle —
the angle between the follower's motion and the cam surface normal,
which produces side-thrust on the follower — is
`atan(velocity / (base+displacement+rollerRadius))`. `CAM.TXT`
recommends keeping this under 35°; the base-radius suggestion table
inverts the same relationship at the point of *maximum* pressure
angle, to show what a larger base circle would buy back.

## Worked Example

`CAM.TXT` documents its own default example directly: a cycloidal
cam, 1.625 base circle radius, 1.25 rise over 120° of rotation. As
independently verifiable checks (no exact numeric worked output was
shipped), this conversion's tests confirm: every motion law starts at
zero displacement and ends at the full rise; the piecewise parabolic
law's two halves agree exactly at the midpoint (where the original
switches formulas); the cycloidal law's midpoint displacement is
exactly half the rise (a hand-derivable value, since `sin(π) = 0`);
the cycloidal law's peak velocity (occurring at that same midpoint)
matches its own closed-form maximum; the cam profile's first point
sits exactly on the base circle and its last point at base-plus-rise;
the default example's own maximum pressure angle falls in a
physically sensible range under 35°; and that the base-radius
suggestion table is monotonic (a larger allowed pressure angle always
permits a smaller base circle).
