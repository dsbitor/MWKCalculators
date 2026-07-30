# ftaper

Female taper angle, measured with two spheres.

**Converted from:** `FTAPER.C` (M. W. Klotz, 12/05, 3/07),
`MWKC/WorkshopUtilities/ftaper.zip`
**Go source:** `MWKGo/ftaper/ftaper.go`

## Purpose

Measures the half angle and included angle of a tapered hole or
socket by dropping two spheres of different diameter into it
and measuring the depth to the top of each sphere.

## Inputs

| Prompt | Default |
|---|---|
| Diameter of small sphere | 0.5 |
| Depth to top of small sphere | 2.0 |
| Diameter of large sphere | 0.75 |
| Depth to top of large sphere | 1.0 |

## Output

Taper half angle and included angle, each in degrees and as an
inch-per-inch ratio.

## Method

```
rs, rl = 0.5*smallDiameter, 0.5*largeDiameter
x      = (smallDepth+rs) - (largeDepth+rl)
z      = (rl-rs) / x
phi    = asin(z)   (clamped to +-90 degrees if z falls outside [-1,1])
```

If the two spheres' tops land at the same height (`x` is zero)
but their diameters differ, `z` goes to infinity; the arcsine is
clamped to 90 degrees rather than producing `NaN`, matching the
domain-clamped `ASND` macro used throughout the original
programs. This conversion carries that behaviour forward in the
shared `internal/mwktrig` package.

## Worked Example

No worked example was included with the original program. As an
independently verifiable check: equal sphere diameters indicate
a cylindrical (untapered) hole, so the computed angle must be
exactly zero regardless of the measured depths.
