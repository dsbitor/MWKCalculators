# sinenl

Sine bar made from two touching cylinders, no connecting link.

**Converted from:** `SINENL.C` (M. W. Klotz, 5/01),
`MWKC/WorkshopUtilities/sine.zip`
**Go source:** `MWKGo/sinenl/sinenl.go`
**Original documentation:** `SINE.TXT`, inside `MWKC/WorkshopUtilities/sine.zip` (not included in this conversion)

## Purpose

A sine bar sets a precise angle for machining without a
protractor. Commercial sine bars are expensive for occasional
home-shop use, but two cylinders of carefully chosen diameters
placed in contact on a flat surface, with a plate laid across
them, will present that plate at a precise, calculable angle:

```
d1 / d2 = (1 - sin(angle/2)) / (1 + sin(angle/2))
```

An earlier version of this idea (`SINE.C`, not yet converted)
used two cylinders of equal diameter connected by a pair of
drilled links spaced a fixed distance apart. This version
removes the link: with the two cylinders simply touching each
other, their own geometry provides the spacing, and only the
two diameters need to be found.

## Inputs

| Prompt | Default |
|---|---|
| Diameter of larger cylinder | 0.75 in |
| Desired angle | 1.5 deg |

## Output

Diameter of the smaller cylinder.

## Method

```
d1 = d2 * (1 - sin(angle/2)) / (1 + sin(angle/2))
```

## Worked Example

From the original author's own notes (`SINE.TXT`), describing
the day this program was written: starting with a 0.75in
cylinder on hand and wanting a 1.5 degree chamfer, the program
gave a required second cylinder diameter of 0.7306in (the
author rounded to 0.731in for turning). This matches the
documented default input exactly.
