# ellipse

Ellipse eccentricity, area, and perimeter.

**Converted from:** `ELLIPSE.C` (M. W. Klotz), `MWKC/Math/ellipse.zip`
**Go source:** `MWKGo/ellipse/ellipse.go`

## Purpose

Unlike a rectangle or circle, an ellipse's perimeter has no simple
rational expression: exactly, it requires the complete elliptic
integral of the second kind, traditionally found by looking up and
interpolating table values. Given an ellipse's semi-major and
semi-minor axes, this program computes eccentricity, flattening, and
area exactly, and computes the perimeter both from a high-accuracy
polynomial approximation to the elliptic integral (used as the
reference value) and from three well known algebraic approximations,
reporting each approximation's error relative to the reference.

## Inputs

| Prompt | Default |
|---|---|
| Semi-major axis | 10 |
| Semi-minor axis | 3 |

## Output

Eccentricity, flattening, area, the elliptic-integral perimeter, and
three algebraic perimeter approximations (RMS, and two due to
Ramanujan) each with its percentage error relative to the
elliptic-integral value.

## Method

```
eccentricity = sqrt(a^2 - b^2) / a
flattening = (a - b) / a
area = pi*a*b

m = (b/a)^2
E(m) ~ (1 + a1*m + a2*m^2 + a3*m^3 + a4*m^4)
     + (b1*m + b2*m^2 + b3*m^3 + b4*m^4) * ln(1/m)
       (Abramowitz & Stegun 17.3.36, accuracy ~2e-8)
exactPerimeter = 4*a*E(m)

rms        = 2*pi*sqrt(0.5*(a^2+b^2))
ramanujan1 = pi*(3*(a+b) - sqrt((3a+b)*(a+3b)))
x = (a-b)/(a+b)
ramanujan2 = pi*(a+b)*(1 + 3x^2/(10+sqrt(4-3x^2)))
```

## Worked Example

No worked numeric example was included with the original program
(`ELLIPSE.TXT` explains the three approximations and their sources
but includes no sample run). The elliptic-integral approximation's
translation from the original was cross-checked independently via
numerical integration of the standard definition before being used
as an expected value in this conversion's tests. As a further,
independently verifiable check: a circle is a degenerate ellipse
(`a == b`) with an exactly known perimeter, `2*pi*r`; the elliptic
integral and all three algebraic approximations agree with it and
each other in that case, since every approximation was designed to
be exact in that limit, confirmed in this conversion's tests.
