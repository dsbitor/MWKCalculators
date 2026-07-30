# ungula

Volume of an ungula (the shape of water in a tilted cylindrical
glass).

**Converted from:** `UNGULA.C` (M. W. Klotz, 6/04), `MWKC/Math/ungula.zip`
**Go source:** `MWKGo/ungula/ungula.go`

## Purpose

If a small amount of water sits in a tilted cylindrical glass
without fully covering the bottom, the water takes on a shape
called an ungula (Latin for "hoof", after its resemblance to a
cow's hoof). Given the cylinder's diameter, the height the water
reaches up the wall on its deep side, and the sagitta (the
distance from the wall on the wet side to the dry/wet boundary
line, measured along a diameter), this program computes the
water's volume. Valid for a sagitta greater than zero and at
most the cylinder's diameter.

## Inputs

| Prompt | Default |
|---|---|
| Cylinder diameter | 2.0 |
| Height of ungula | 10.0 |
| Sagitta of ungula base | 1.0 |

Re-prompts for the sagitta if it is outside the valid range.

## Output

Volume of the ungula, in the cube of whatever length unit was
used for the inputs.

## Method

```
a   = sqrt(2*sagitta*radius - sagitta^2)
phi = pi/2 + atan((sagitta-radius)/a)
volume = height*radius^2 * (3*sin(phi) - 3*phi*cos(phi) - sin(phi)^3) / (1-cos(phi)) / 3
```

At a sagitta equal to the full diameter, `a` is zero and the
`atan` argument becomes infinite; `math.Atan(+Inf)` correctly
returns `pi/2` in that limit under IEEE 754, so the formula
still gives the right answer at that boundary without a special
case.

## A correction to the original author's notes

The original author's own notes state: "In the case where the
sagitta equals the diameter of the cylinder, the somewhat complex
formula reduces to a simple result: volume = (2/3)*h*r*r." That
specific claim does not hold.

Verified independently, two ways: by working through the
algebra by hand, and by brute-force numerical integration of the
physical shape the formula describes (see `ungula_test.go`),
across the full range of valid sagitta values:

- At sagitta equal to the **radius** (half the diameter, a
  semicircular wet region), the volume is exactly
  `(2/3)*height*radius^2`. This is the case the original notes
  almost certainly meant.
- At sagitta equal to the full **diameter** (the entire base is
  wet), the volume is `(pi/2)*height*radius^2`, which follows
  directly from the average depth (`height/2`, since depth rises
  linearly across a full diameter) times the base area
  (`pi*radius^2`).

The formula itself is correct throughout the whole valid range,
confirmed by matching brute-force numerical integration to
floating-point precision at several sagitta values; only the
author's own stated special case, in the explanatory notes
rather than the code, was incorrect.
