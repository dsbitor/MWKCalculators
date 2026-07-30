# dplate

Disk diameter for building a dividing plate without a dividing
head.

**Converted from:** `DPLATE.C` (M. W. Klotz, 1/02),
`MWKC/WorkshopUtilities/dplate.zip`
**Go source:** `MWKGo/dplate/dplate.go`

## Purpose

A dividing head plate with a given number of divisions can be
made without already having a dividing head, or a plate with
that hole count, by a different construction: turn a top-hat
shaped central disk (the "crown"), then turn a set of equal
disks that, when placed around the crown's edge, simultaneously
touch the crown and their own neighbors. A detent between
adjacent disks then locates each division. Given the crown's
diameter and the number of divisions wanted, this program
computes the disk diameter needed.

## Inputs

| Prompt | Default |
|---|---|
| Number of divisions | 14 |
| Diameter of mounting circle | 112 |

## Output

Required disk diameter.

## Method

```
theta = 360 / divisions
rd    = 0.5*mountingDiameter * sin(theta/2) / (1 - sin(theta/2))
```

## Worked Example

No worked example was included with the original program. As an
independently verifiable check: at exactly six divisions, this
is the classic hexagonal circle-packing arrangement (six equal
circles around a central circle of the same size, all mutually
touching), so the disk diameter must equal the mounting circle
diameter exactly, which the conversion's tests confirm.
