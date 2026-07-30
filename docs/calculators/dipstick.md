# dipstick

Tank volume fraction from dipstick reading.

**Converted from:** `DIPSTICK.C` (M. W. Klotz, 4/01, 9/03, 4,7/04,
7,9/05), `MWKC/WorkshopUtilities/dipstick.zip`
**Go source:** `MWKGo/dipstick/dipstick.go`

## Purpose

Given a dipstick's wetted length and a tank's shape and dimensions,
computes what fraction of the tank's full capacity remains, for nine
tank shapes: horizontal cylinder, sphere, horizontal tank with
elliptical cross-section, vertical or horizontal "cartouche" (a
stadium shape: two semicircular ends joined by straight sections, as
seen on tanker trucks — not to be confused with an ellipse), bucket
(conical frustum), barrel (modeled with elliptical sides), and
horizontal cylinder with either hemispherical or "dished"
(torispherical) ends. It also calibrates a dipstick by finding the
wetted length for each of a series of percentage increments of full
capacity, for any of the nine shapes.

## Inputs

Shape choice (A-I), followed by that shape's own dimension prompts
(diameter, length, height, etc. as applicable), then the wetted
dipstick length to evaluate, then a percentage increment for
calibration.

## Output

Percentage of tank capacity remaining at the given wetted length, then
a calibration table of wetted length at each percentage increment.

## Method

Most shapes reduce to a circular-segment area (shared by several
tank cross-sections) or a direct volume formula:

```
segmentArea(r, h) = 0.5*(r^2*angle - z*chord)
  where z = r-h, angle = 2*acos(z/r), chord = 2*r*sin(angle/2)
```

The dished-ends shape's end caps have no closed form and are found
by numerical integration along the dish's torispherical profile (a
bounded 1000-step sum, matching the original program's own fixed
step size), reused for both dipstick heights below and above the
dish's own radius via a mirror-image relationship.

Calibration uses the same bounded binary search technique as `buoy`
in Tier 1 group 10.

## Worked Example

No worked numeric example was included with the original program. As
independently verifiable checks: for eight of the nine shapes, the
wetted fraction is confirmed to be exactly 0 at zero wetted height,
exactly 1 at the shape's own full height, and non-decreasing in
between — physical facts about any tank, not a re-run of each shape's
formula. The ninth shape (dished ends) does not quite reach exactly 0
or 1 at those extremes: its torispherical dish profile has interior
cross-sections *wider* than the dish's own nominal radius, so a thin
sliver of the dish registers as wetted even at zero overall dipstick
height. This is a property of the original algorithm's own dish
model, not introduced by this conversion, so it's confirmed only to
land within about 2% of 0 and 1 respectively, rather than checked
against the exact identity used for the other eight shapes.
