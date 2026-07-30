# gearpa

Gear pressure angle estimation by chordal span.

**Converted from:** `GEARPA.C`, `MWKC/WorkshopUtilities/gearpa.zip`.
Reference: Machinery's Handbook, "Checking Gear Sizes".
**Go source:** `MWKGo/gearpa/gearpa.go`

## Purpose

The two most common gear pressure angles, 14.5 and 20 degrees,
produce a measurably different chordal span (the distance measured
across several teeth with calipers) for gears of the same tooth
count and diametral pitch. Given the tooth count and diametral pitch,
this program predicts the span for each candidate pressure angle, so
a hobbyist without dedicated gear-checking equipment can measure the
actual span and compare it against the two predictions to guess which
pressure angle a gear was cut to.

## Inputs

| Prompt | Default |
|---|---|
| Number of teeth on gear | 30 |
| Diametral pitch of gear | 6 |

Number of teeth must be between 12 and 110.

## Output

For each of the two pressure angles: the number of teeth to span with
calipers and the predicted chordal span, or a message that the tooth
count is out of the table's range for that pressure angle.

## Method

```
halfToothAngle = 180/teeth
pitchRadius = 0.5*teeth/diametralPitch
toothThickness = 0.5*pi/diametralPitch
involuteFactor = 2*(tan(pa) - pa in radians)
span = pitchRadius*cos(pa) * (toothThickness/pitchRadius
                               + 2*pi*(teethSpanned-1)/teeth
                               + involuteFactor)
```

`teethSpanned` comes from a lookup table (per Machinery's Handbook)
keyed by tooth count, separate for each pressure angle; the 20 degree
table only covers up to 81 teeth.

The original program's 14.5 degree table is a sequence of
overlapping range checks (`51-62` then `53-75`), each overwriting the
previous rather than a clean partition; since C evaluates them as
independent `if` statements, not `else if`, the actual effect is that
teeth counts 51-52 get a span of 5 teeth while 53-75 get 6, not the
51-62/63-75 split the pattern otherwise suggests. This is almost
certainly an unintended typo (`63` for `53` would match the 20 degree
table's own clean, non-overlapping pattern), but it's what the
original program actually computes and prints, so this conversion
preserves it rather than silently correcting it; see the code comment
on `spanTeethCount` for the exact boundary.

## Worked Example

`GEARPA.TXT` includes two complete worked examples, both reproduced
exactly in this conversion's tests:

- 30 teeth, DP 6: 14.5° spans 3 teeth at 1.2941 in; 20° spans 4 teeth
  at 1.7921 in.
- 18 teeth, DP 6: both pressure angles span 2 teeth, at 0.7765 in
  (14.5°) and 0.7800 in (20°) respectively, illustrating how close the
  two measurements can be for a small gear.
