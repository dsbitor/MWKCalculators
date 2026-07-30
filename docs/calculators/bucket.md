# bucket

Slant height divisions for a bucket-shaped (frustum) container.

**Converted from:** `BUCKET.C` (M. W. Klotz, 11/99),
`MWKC/WorkshopUtilities/bucket.zip`
**Go source:** `MWKGo/bucket/bucket.go`

## Purpose

A cylindrical container can be divided into equal-volume parts
just by marking equal height divisions, since its cross-section
never changes. A bucket-shaped container (a conical frustum)
can't: its cross-section varies with height, and every frustum
volume formula uses the height measured perpendicular to the
base, not the slant height measured along the sloping side a
person would actually mark and measure. Given the bucket's big
and small diameters and its slant height, this program finds the
slant heights to mark to divide its volume into equal parts.
Measurements should be taken on the inside of the bucket if it
has substantial wall thickness, since wall thickness isn't
accounted for.

## Inputs

| Prompt | Default |
|---|---|
| Diameter of bucket big end | 4.0 |
| Diameter of bucket small end | 3.0 |
| Slant height of bucket | 6.0 |
| Divide volume into how many parts | 4 |

## Output

Total volume, then the slant height for each division from the
small end.

## Method

The overall shape uses the standard frustum-of-a-cone volume
formula:

```
halfAngle = asin((largeRadius-smallRadius) / slantHeight)
height = slantHeight * cos(halfAngle)
volume = pi*height*(largeRadius^2 + largeRadius*smallRadius + smallRadius^2) / 3
```

Each division's slant height is found by an incremental search
(matching the original program's own method, bounded to a
maximum number of steps rather than relying on a keypress to
bail out of a runaway loop): starting from an initial guess,
nudge the height by a small step, in whichever direction reduces
the volume error, until the error falls within tolerance.

## Worked Example

No worked numeric example was included with the original
program. As an independently verifiable check: a frustum with
equal top and bottom radii is just a cylinder, whose volume is
`pi*height*radius^2` exactly, confirmed in this conversion's
tests; and searching for the height containing the bucket's
entire volume converges back to the bucket's own full slant
height.
