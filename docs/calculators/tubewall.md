# tubewall

Tube wall thickness from an outside micrometer reading.

**Converted from:** `TUBEWALL.C` (M. W. Klotz, 2/02),
`MWKC/WorkshopUtilities/tubewall.zip`
**Go source:** `MWKGo/tubewall/tubewall.go`
**Original documentation:** `TUBEWALL.TXT`, inside `MWKC/WorkshopUtilities/tubewall.zip` (not included in this conversion)

## Purpose

An ordinary flat-anvil micrometer reads a tube's wall as thicker
than it really is: the flat anvil bridges across the tube's
curved inside surface instead of touching it directly, an error
that gets worse the smaller the tube's bore. The correct fix is
a micrometer with a hemispherical anvil, which contacts the bore
at a point rather than bridging it, but not everyone has one.
Given the (flat) anvil's diameter, the tube's outside diameter,
and the reading obtained with an ordinary micrometer, this
program computes the actual wall thickness.

## Inputs

| Prompt | Default |
|---|---|
| Micrometer anvil diameter | 0.249 |
| Tube outside diameter | 0.879 |
| Micrometer measurement | 0.0625 |

## Output

Tube wall thickness, and the corresponding tube inside diameter.

## Method

Solves a quadratic derived from the anvil, tube, and measurement
geometry:

```
B = -tubeOutsideDiameter
C = -anvilRadius^2 + 2*tubeOutsideRadius*measurement - measurement^2
thickness = -0.5 * (B + sqrt(B^2 - 4*C))
```

## Worked Example

`TUBEWALL.TXT` gives the author's own measurement of a piece of
copper pipe at the program's default inputs, reporting a
computed wall thickness of 0.0425, which this conversion
reproduces exactly. The same note also compares this method
against a precision-ball micrometer (0.0422) and electronic
caliper knife edges (0.0415) on the same pipe, for context on
the method's accuracy.
