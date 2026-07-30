# mandrel

Recommended mandrel diameter for spring winding.

**Converted from:** `MANDREL.C` (M. W. Klotz, 2/05),
`MWKC/WorkshopUtilities/mandrel.zip`
**Go source:** `MWKGo/mandrel/mandrel.go`

## Purpose

A coil spring wound on a mandrel springs back slightly once
released, ending up with a larger inside diameter than the
mandrel itself. Kozo Hiraoka's article ("Home Shop Machinist",
July/August 1987, pg. 30) gives an empirical correction factor
for this springback, separately fitted for music wire and
phosphorus bronze wire. Given the wire type, wire diameter, and
desired spring inside diameter, this program computes the
mandrel diameter to wind on.

## Inputs

| Prompt | Default |
|---|---|
| Wire type | music wire [0] or phosphorus bronze (1) |
| Wire diameter | 0.040 in |
| Spring inside diameter | 0.203 in |

Any wire type value other than 1 is treated as music wire (0),
matching the original program's own validation.

## Output

Recommended mandrel diameter.

## Method

Hiraoka's empirical coefficients, indexed by wire type:

| Wire type | Constant coefficient | Linear coefficient |
|---|---|---|
| Music wire | 0.980364 | -0.012436 |
| Phosphorus bronze | 0.012436 | -0.11018 |

```
averageSpringDiameter = springInsideDiameter + wireDiameter
factor = constantCoefficient + linearCoefficient*averageSpringDiameter/wireDiameter
mandrelDiameter = factor*averageSpringDiameter - wireDiameter
```

## Worked Example

No worked example was included with the original program. Since
Hiraoka's coefficients are empirical, there is no independent
closed-form identity to check the formula against; the
conversion's tests instead confirm the two wire types produce
different results for the same dimensions, which follows from
their distinct coefficients.
