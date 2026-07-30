# sprocket

ANSI standard roller chain sprocket dimensions.

**Converted from:** `SPROCKET.C` (M. W. Klotz),
`MWKC/WorkshopUtilities/sprocket.zip` (the same archive as `chain`,
converted in Tier 1 group 4)
**Go source:** `MWKGo/sprocket/sprocket.go`

## Purpose

Written for Martin Vinnicombe, to compute dimensions for ANSI
standard roller chain sprockets, per the reference data in
Machinery's Handbook. Given the tooth count, chain pitch, and roller
diameter, this program computes the pitch diameter, outside diameter,
caliper diameter (via a caliper factor, for an odd tooth count), and
the maximum recommended hub diameter.

## Inputs

| Prompt | Default |
|---|---|
| Number of teeth in sprocket | 9 |
| Chain pitch | 1 in |
| Roller diameter | 0.25 in |

## Output

Pitch diameter, outside diameter, caliper factor (odd tooth count
only), caliper diameter, and maximum hub diameter.

## Method

```
halfToothAngle = 180 / teeth
pitchDiameter = pitch / sind(halfToothAngle)
outsideDiameter = pitch * (0.6 + 1/tand(halfToothAngle))
maxHubDiameter = pitch*(1/tand(halfToothAngle) - 1) - 0.030

odd teeth:  caliperFactor = pitchDiameter * cosd(90/teeth)
            caliperDiameter = caliperFactor - rollerDiameter
even teeth: caliperDiameter = pitchDiameter - rollerDiameter
```

`SPROCKET.TXT` describes the odd-tooth caliper factor as
`pitchDiameter * cos(180/teeth)`, but the source code computes and
prints `cos(90/teeth)`, which is also the standard Machinery's
Handbook formula for this measurement; the `.TXT` appears to have a
transcription slip (180 for 90). This conversion implements what the
code (and the original `.EXE`) actually computes.

## Worked Example

No worked numeric example was included with the original program. As
an independently verifiable check: the pitch diameter is,
geometrically, the diameter of the circle passing through the
vertices of a regular polygon whose side length is the chain pitch,
`pitch / (2*sin(pi/teeth))` doubled; confirmed against this
program's own trigonometric form of the same formula in this
conversion's tests.
