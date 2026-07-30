# loft

Minimum thread engagement length.

**Converted from:** `LOFT.C` (M. W. Klotz, 1/05),
`MWKC/WorkshopUtilities/loft.zip`
**Go source:** `MWKGo/loft/loft.go`

## Purpose

When a screw is threaded into a tapped hole, too little engagement
lets the internal threads strip before the screw itself fails; too
much wastes material and effort. This program computes the minimum
length of thread engagement at which the screw's own tensile strength
and the internal thread's shear strength balance, so a longer
engagement guarantees the screw breaks before the threads strip
(assuming screw and hole are the same material; multiply by the
strength ratio J otherwise).

## Inputs

| Prompt | Default |
|---|---|
| (m)etric or [i]mperial | [i]mperial |
| Basic diameter of screw | 0.25 in / 4 mm |
| Pitch of screw | 20 tpi (Imperial) / 0.7 mm (metric) |

## Output

Pitch circle diameter, screw thread tensile area, thread shear area,
length of thread engagement (in both length units and number of
threads).

## Method

```
pitchCircleDiameter = diameter - 0.64952*pitch
tensileArea = 0.25*pi*(diameter - 0.938194*pitch)^2
shearArea = 0.5*pi*pitchCircleDiameter
engagementLength = 2*tensileArea / shearArea
```

The original program's "Pitch circle diameter of thread" line prints
this length-valued figure with the area unit label (`in^2`/`mm^2`)
rather than the length unit; this conversion preserves that label
choice for fidelity with the original output, though it is almost
certainly a copy-paste label slip in the original rather than an
intentional unit.

## Worked Example

No worked numeric example was included with the original program. As
an independently verifiable check: as the thread pitch shrinks toward
zero, the tensile area approaches the full bolt cross section
(`0.25*pi*D^2`) and the shear area approaches `0.5*pi*D`, so the
engagement length approaches the diameter itself; confirmed in this
conversion's tests.
