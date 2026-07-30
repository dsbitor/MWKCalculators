# eccent

Packing thickness for turning an eccentric in a 3-jaw chuck.

**Converted from:** `ECCENT.C` (M. W. Klotz, 11/98, 9/01),
`MWKC/WorkshopUtilities/eccent.zip`
**Go source:** `MWKGo/eccent/eccent.go`

## Purpose

Rather than setting up a four-jaw chuck, an eccentric can be
turned in a three-jaw chuck by packing one jaw with shim stock,
shifting the workpiece's turning axis away from the spindle
axis. Given the width of the chuck jaws (specifically the part
of the jaw that actually contacts the work, which can be
narrower than the jaw body), the workpiece diameter, and the
required offset, this program computes the packing thickness
needed.

The later `ECCENTUB.C` (converted separately as `eccentub`)
replaces this method with a slotted tube, which avoids the
awkward jaw-width measurement this method requires.

## Inputs

| Prompt | Default |
|---|---|
| Width of chuck jaws | 0.125 in |
| Diameter of workpiece | 1.5625 in |
| Required eccentric offset | 0.28125 in |

## Output

Required packing size, or a message that the work would fall
through the unpacked jaws at the requested offset.

## Method

If the offset exceeds `(radius+jawWidth)/sqrt(3)`, the workpiece
would fall through the unpacked jaws and no packing size is
computed. Otherwise:

```
if jawWidth > sqrt(3)*offset:
    packing = 1.5 * offset
else:
    packing = 1.5*offset - radius
            + 0.5*sqrt(4*radius^2 - 3*offset^2 + 2*offset*jawWidth*sqrt(3) - jawWidth^2)
```

The original source comments note that David Hoskins of
Australia pointed out an error in this equation in September
2001, corrected in the version converted here; the comments
don't identify which specific term the correction touched, so
that detail isn't claimed here either.

## Worked Example

No worked numeric example was included with the original
program.
