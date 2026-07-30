# gearspur

Spur gear dimension solver.

**Converted from:** `GEARSPUR.C`, `MWKC/WorkshopUtilities/gearspur.zip`
**Go source:** `MWKGo/gearspur/gearspur.go`

## Purpose

Given any two of a spur gear's outside diameter, tooth count,
diametral pitch, or pitch diameter (a module value is also accepted
as an alternate way to specify diametral pitch), this program solves
for the rest and reports the full set of standard gear dimensions,
plus the Brown & Sharpe involute cutter number needed to cut the
gear.

## Inputs

| Prompt | Default |
|---|---|
| [I]mperial or (M)etric units | [I] |
| OD of gear | 2.35 |
| Number of teeth | 45 |
| Diametral Pitch (if still needed) | 20 |
| Module (if still needed) | 0.7874 |
| Pitch Diameter (if still needed) | 2.25 |

Exactly two data items (from OD, teeth, diametral pitch/module, or
pitch diameter) are required.

## Output

Diametral pitch, module, tooth count, outside diameter, pitch
diameter, addendum, dedendum, whole depth, circular pitch, tooth
thickness (each length in both inches and millimeters), and the
Brown & Sharpe cutter number.

## Method

```
teeth+OD:        DP = (teeth+2)/OD;         PD = teeth/DP
DP+teeth:        OD = (teeth+2)/DP;         PD = teeth/DP
DP+OD:           teeth = DP*OD-2;           PD = teeth/DP
PD+teeth:        DP = teeth/PD;             OD = (teeth+2)/DP
PD+OD:           teeth = 2*PD/(OD-PD);      DP = teeth/PD
PD+DP:           teeth = DP*PD;             OD = (teeth+2)/DP

addendum = 1/DP
wholeDepth = 2.2/DP + 0.002   if DP <= 20
           = 2.157/DP          otherwise
dedendum = wholeDepth - addendum
circularPitch = pi/DP
toothThickness = 0.48*circularPitch
```

Brown & Sharpe cutter number is looked up from tooth count via a
standard table (cutter 8 for 12-13 teeth, down to cutter 1 for 135+
teeth).

## Worked Example

No worked numeric example was included with the original program. As
an independently verifiable check: a single self-consistent gear (45
teeth, DP 20, giving OD 2.35 and PD 2.25) is solved from each of the
six recognized two-quantity combinations, and each recovers the same
full, consistent set; confirmed in this conversion's tests, along
with the documented `DP <= 20` / `DP > 20` whole-depth tier boundary.
