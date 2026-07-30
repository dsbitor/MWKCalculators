# soddy

Outer and inner Soddy circle diameters for three mutually tangent
circles.

**Converted from:** `SODDY.C` (M. W. Klotz, 3/02, 12/04),
`MWKC/WorkshopUtilities/plug.zip`
**Go source:** `MWKGo/soddy/soddy.go`

## Purpose

Given the diameters of three mutually tangent circles, computes
the diameter of a fourth circle tangent to all three: either the
outer Soddy circle, which encloses all three, or the inner Soddy
circle, nestled in the gap between them. This is Descartes'
Circle Theorem.

No `.TXT` file was included with the original program; this
purpose statement is drawn from the `.C` file's own header
comment. The same zip file also contains `PLUG.C`, an unrelated
program not yet converted.

## Inputs

| Prompt | Default |
|---|---|
| Diameter of first circle | 0.245 |
| Diameter of second circle | 0.249 |
| Diameter of third circle | 0.253 |

## Output

Diameter of the outer Soddy circle and the inner Soddy circle.

## Method

```
radical = +-2 * sqrt(d1*d2*d3*(d1+d2+d3))   (+ for outer, - for inner)
diameter = |d1*d2*d3 / (d2*d3 + d1*(d2+d3) - radical)|
```

## Worked Example

No worked example was included with the original program. As an
independently verifiable check: this conversion's tests confirm
the formula matches Descartes' Circle Theorem stated in its usual
curvature form (`curvature = 2/diameter`) across several inputs,
and that three equal circles of diameter 1 give the well known
inner tangent circle radius `R*(2/sqrt(3)-1)`.
