# sine

Sine bar made from two cylinders connected by drilled links.

**Converted from:** `SINE.C` (M. W. Klotz, 11/98),
`MWKC/WorkshopUtilities/sine.zip`
**Go source:** `MWKGo/sine/sine.go`

## Purpose

The original method for building a custom sine bar from two
cylinders: hold two cylinders a fixed center-to-center distance
apart using a pair of identically drilled links, and a plate
laid across both cylinders sits at a calculable angle. Given
the first cylinder's diameter, the link spacing, and the
desired angle, this program computes the required diameter of
the second cylinder, checks that the two cylinders are small
enough to fit within the link spacing, and reports how
sensitive the resulting angle is to small manufacturing errors
in either the first cylinder's diameter or the link spacing.

The later `SINENL.C` (converted separately as `sinenl`) removes
the links entirely: two cylinders simply placed in contact
provide their own spacing.

## Inputs

| Prompt | Default |
|---|---|
| Diameter of first cylinder (d1) | 0.375 in |
| Distance between cylinder centers (r) | 3.0 in |
| Desired angle | 10.0 deg |

## Output

Diameter of the second cylinder, or a message that the second
cylinder is too large to fit at the given link spacing, plus the
angle error caused by a 0.001in error in the first cylinder's
diameter and a 0.001in error in the link spacing.

## Method

```
d2 = 2*r*sin(angle/2) + d1
```

fits if `0.5*(d1+d2) <= r`. The sensitivity figures come from
differentiating the angle with respect to `d1` and `r`.

## Worked Example

From the original author's own notes: building a 10 degree sine bar
with a 0.375in first cylinder and a 3in link spacing needs a second
cylinder of 0.898in, and the same notes
state that a 0.001in error in either cylinder's diameter causes
about 0.01 degrees of angle error, while the same error in the
link spacing causes about 0.002 degrees. Checking the ported
formula's sensitivity output against these figures surfaced a
transcription error in this conversion's own arithmetic (a
squared term where the original formula does not have one); the
error was in the new code, not the original, and is fixed.
