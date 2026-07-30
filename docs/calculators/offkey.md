# offkey

Offset key calculations.

**Converted from:** `OFFKEY.C` (M. W. Klotz, 1/03, for Ronnie
Shultz), `MWKC/WorkshopUtilities/offkey.zip`
**Go source:** `MWKGo/offkey/offkey.go`

## Purpose

An "offset key" is a shaft key rotated slightly from its normal
centered position (for example, to clear an obstruction). Given the
shaft diameter, key width, key height, and the rotation angle, this
program computes the standard (non-rotated) keyseat and keyway
depths, then the corner positions of the key's cross-section both
before and after rotation, plus the relative distances between
several of those corners, needed to lay out the offset cut.

## Inputs

| Prompt | Default |
|---|---|
| Shaft diameter | 1.0 |
| Key width | 0.5 |
| Key height | 0.5 |
| Depth of cut in shaft (Q) | (defaults to the computed standard value) |
| Key rotation angle (theta) | 5 deg |

Rotation angle must not exceed the maximum possible angle (computed
from the shaft diameter and key width); the program re-prompts
otherwise.

## Output

Standard keyseat/keyway/cut depths; the 8 corner points of the key
cross-section before and after rotation, plus 2 derived intersection
points; and the relative x/y separation between several pairs of
those points.

## Method

```
x = sqrt(D^2 - K^2); y = 0.5*(D-x)
M (keyseat depth) = D - (y + 0.5*H)
N (keyway depth)  = M + H
Q (shaft cut)      = D - M
R (hub cut)        = H - Q
maxAngle = 2*asind(K/D)
```

The 8 corners are laid out symmetrically about the key centerline,
using the round-shaft profile at ±K/2; the 5 corners nearest the
"open" side of the key (indices 2-6) are rotated by theta about the
origin using a standard 2D rotation, while the 3 corners defining the
shaft cut itself (indices 0, 1, 7) stay fixed, since the shaft cut
doesn't move with the key. Two additional points are found as the
intersections of specific edges before and after rotation, via the
standard two-line-intersection formula.

## Worked Example

No worked numeric example was included with the original program. As
independently verifiable checks: rotating a point about the origin
must preserve its exact distance from the origin (confirmed for
every point the original program actually rotates), the three
un-rotated points must be unchanged, and a zero-degree rotation must
leave every point unchanged; all confirmed in this conversion's
tests, along with a line-intersection check against a known crossing
point independent of this code's own determinant-based formula.
