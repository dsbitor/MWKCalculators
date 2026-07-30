# catenary

Droop and length of a hanging cable.

**Converted from:** `CATENARY.C` (M. W. Klotz, 6/09),
`MWKC/Math/catenary.zip`
**Go source:** `MWKGo/catenary/catenary.go`

## Purpose

A catenary is the shape a flexible, non-stretching chain or
cable takes when hung from two points under gravity. Given the
cable's tension, its weight per unit length, and the
straight-line distance between its supports, this program
computes the droop at the center of the span and the total
cable length.

## Inputs

| Prompt | Default |
|---|---|
| Tension in cable | 0.674427074 lbf |
| Cable weight per unit length | 0.00044091 lb/ft |
| Straightline distance between supports | 1.640419948 ft |

## Output

Cable droop at the center of the span, and total cable length,
both in feet.

## Method

```
param  = tension / density
droop  = param * (cosh(0.5 * distance / param) - 1)
length = 2 * param * sinh(0.5 * distance / param)
```

## Worked Example

No worked example was included with the original program. As an
independently verifiable check: with tension and density both
equal to 1, the formulas reduce to the well known constants
cosh(1) - 1 and 2 * sinh(1).
