# flute

Tapered flute milling setup.

**Converted from:** `FLUTE.C` (M. W. Klotz),
`MWKC/WorkshopUtilities/flute.zip`
**Go source:** `MWKGo/flute/flute.go`

## Purpose

One way to cut a tapered flute (such as on an engine support column)
is to tilt the workpiece and plunge cut along it with a ball end
mill; as the cut progresses along the inclined workpiece it widens,
forming the taper. Given the ball mill diameter and the flute's
radius at each end, this program computes the depth of cut at each
end and the workpiece inclination angle needed to connect them over a
given flute length.

## Inputs

| Prompt | Default |
|---|---|
| Ball mill diameter | 0.5 in |
| Flute radius at small end | 0.1 in |
| Flute radius at large end | 0.2 in |
| Length of flute | 3.0 in |

Both flute radii must be no greater than the ball mill radius; the
program re-prompts otherwise.

## Output

Depth of cut at each end of the flute, and the required workpiece
inclination angle (in decimal degrees and degrees/minutes/seconds).

## Method

```
depth = millRadius - sqrt(millRadius^2 - fluteRadius^2)   (per end)
inclination = atand((depthLarge - depthSmall) / length)
```

## Worked Example

No worked numeric example was included with the original program. As
independently verifiable checks: a flute radius of zero has zero
depth of cut (the sagitta of a zero-width chord is zero), and a
flute radius equal to the mill's own radius cuts to the full depth of
the mill radius; both confirmed in this conversion's tests.
