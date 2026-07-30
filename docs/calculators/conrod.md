# conrod

Connecting rod to cylinder wall clearance.

**Converted from:** `CONROD.C` (M. W. Klotz, with Tom Roach),
`MWKC/WorkshopUtilities/conrod.zip`
**Go source:** `MWKGo/conrod/conrod.go`

## Purpose

When designing a piston engine, the connecting rod must not foul the
bottom of the cylinder at its maximum lateral extension, which occurs
when the crank radius and the connecting rod are at right angles.
Given the rod length, crank radius, height of the cylinder bottom
above the crank center, and cylinder diameter, this program computes
the worst-case clearance between the rod centerline and the cylinder
wall at that position, plus several related distances useful for
engine layout. Any consistent unit system may be used.

## Inputs

| Prompt | Default |
|---|---|
| Connecting rod length (center-to-center) | 2.4 |
| Crank radius | 0.6 |
| Height of cylinder bottom above crank center | 1.5 |
| Cylinder diameter | 1.0 |

## Output

Phi (the angle between the con rod and the line to the crank
center), the gudgeon-pin-to-crank-center distance, the worst-case
clearance, and several intermediate distances (`d34`, `d45`, `d35`,
`d23`, `d13`, `d12`, `d14`, `d25`) matching the labeled points in the
original program's accompanying diagram. If the worst-case clearance
(`d23`) is greater than half the connecting rod's own thickness,
the rod will not foul the cylinder wall; otherwise the designer needs
to increase the crank-to-cylinder-bottom height, lengthen the rod, or
relieve the cylinder wall for clearance. For a complete time history
of the gudgeon pin's position through a full rotation rather than
just this worst-case snapshot, see [crod](crod.md).

## Method

```
gudgeonToCrank = sqrt(crankRadius^2 + rodLength^2)
phi = acosd(rodLength / gudgeonToCrank)
d34 = 0.5*cylinderDiameter
d45 = (gudgeonToCrank - cylinderBottomHeight) * tan(phi)
d35 = d34 - d45
d23 = d35 / cos(phi)      <- the worst-case clearance
d13 = d34 / cos(phi)
d12 = d13 - d23
d14 = d34 * tan(phi)
d25 = d35 * sin(phi)
```

## Worked Example

No fully worked numeric example is available. As an
independently verifiable check: the gudgeon-to-crank distance is the
hypotenuse of the crank radius and rod length by the Pythagorean
theorem regardless of the other inputs, and a crank radius of zero
degenerates every lateral-offset quantity to zero; both confirmed in
this conversion's tests.
