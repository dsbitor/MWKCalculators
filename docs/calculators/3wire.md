# 3wire

Three-wire thread measurement.

**Converted from:** `3WIRE.C`, `MWKC/WorkshopUtilities/3wire.zip`.
Reference: Machinery's Handbook, 23rd Edition, pg. 1498.
**Go source:** `MWKGo/3wire/3wire.go`
**Original documentation:** `3WIRE.TXT`, inside `MWKC/WorkshopUtilities/3wire.zip` (not included in this conversion)

## Purpose

The three-wire method measures a 60 degree screw thread's pitch
diameter indirectly: three precision wires of a known diameter are
laid in the thread grooves, and the distance measured over the outer
two wires relates to the pitch diameter by a well known formula that
includes a small correction for the thread's lead angle. Given the
thread pitch, number of starts, wire diameter, and major diameter,
this program converts between pitch diameter and measurement over
wires in either direction.

## Inputs

| Prompt | Default |
|---|---|
| (M)etric or [I]mperial threads | [I] |
| Thread pitch | 20 tpi (Imperial) / 1 mm (metric) |
| Number of starts | 1 |
| Wire diameter used | (defaults to the computed "best wire" size) |
| Major diameter of thread | 0.25 in / 6 mm |
| Calculate pd from mow (p) or mow from pd [m] | [m] |
| Pitch diameter or Measurement over wires | (depends on the choice above) |

## Output

The suggested "best wire" diameter, then either the measurement over
wires (given pitch diameter) or the pitch diameter (given measurement
over wires), in both millimeters and inches.

## Method

```
bestWireDiameter = 0.5*pitch / cos(30 deg)
pitchDiameterFromMajor = majorDiameter - 2*(3/16)*pitch/tan(30 deg)

leadAngleTan = pitch*starts / (pi*pitchDiameter)
x = -0.5*pitch*cot(30) + wireDiameter*(1 + 1/sin(30)
                                        + 0.5*leadAngleTan^2*cos(30)*cot(30))

measurementOverWires = pitchDiameter + x
```

The lead angle correction term `x` uses whichever pitch diameter is
actually known at the time: the user's own measured value when
solving for measurement over wires, or the major-diameter-derived
estimate when solving in the other direction (since that direction's
pitch diameter isn't known until `x` is subtracted from the
measurement) — matching the original program's own approach, an
adequate approximation since this term is a small second-order
correction.

## Worked Example

No worked numeric example was included with the original program
(`3WIRE.TXT` is a reference table of standard thread dimensions, not
a sample run). As an independently verifiable check: the "best wire"
diameter is a well known standard constant, `0.577350*pitch`,
confirmed in this conversion's tests; and converting a pitch diameter
to a measurement over wires and back exactly recovers the original
value when the same pitch-diameter estimate is used for the lead
angle correction in both directions.
