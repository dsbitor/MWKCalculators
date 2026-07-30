# mradius

Radius of curvature of a part's edge, measured with two rollers.

**Converted from:** `MRADIUS.C` (M. W. Klotz, 3/02),
`MWKC/WorkshopUtilities/mradius.zip`
**Go source:** `MWKGo/mradius/mradius.go`

## Purpose

If a part has a radiused edge and the radius needs measuring,
one accurate method uses a surface plate, gage blocks, and two
rollers of well-known diameter. The original `.jpg` (not
carried forward here) shows the physical setup and the
underlying geometry; this program performs the calculation from
the measured quantities.

## Inputs

| Prompt | Default |
|---|---|
| Diameter of rollers (d) | 0.25 |
| Measurement across rollers (m) | 2.0 |
| Gap (g) | 0.1 |

## Output

Radius of curvature.

## Method

```
r = 0.5 * d
w = 0.5 * m - r
R = ((r-g)^2 + w^2 - r^2) / (2 * (d-g))
```

## Worked Example

No worked numeric example was included with the original
program. As an independently verifiable check: with a zero gage
block gap, the formula reduces algebraically to `w^2 / (2*d)`,
which the conversion's tests confirm the full formula matches
exactly.
