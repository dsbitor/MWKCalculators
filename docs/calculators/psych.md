# psych

Relative humidity from wet/dry bulb psychrometer.

**Converted from:** `PSYCH.C`, `MWKC/Misc/psych.zip`. Reference:
http://www.uswcl.ars.ag.gov/exper/relhumeq.htm
**Go source:** `MWKGo/psych/psych.go`

## Purpose

A conventional psychrometer uses two identical thermometers, one
with its bulb kept wet by an evaporating wick; the wet bulb reads
lower than the dry bulb due to evaporative cooling, and the size of
that difference indicates relative humidity. Given the site elevation
and the two temperature readings, this program computes relative
humidity and dew point.

## Inputs

| Prompt | Default |
|---|---|
| Elevation above sea level | 970 ft |
| Preferred temperature scale | [F], C |
| Dry bulb temperature | 70 degF / 20 degC |
| Wet bulb temperature | 60 degF / 14 degC |

## Output

Relative humidity (%), and dew point temperature in both Celsius and
Fahrenheit.

## Method

Magnus-Tetens saturation vapor pressure approximation:

```
airPressure = 101.325 * exp(-0.0001184*elevationM)
a = 0.00066*(1+0.00115*wetBulbC)
esdb = exp((16.78*dryBulbC - 116.9) / (dryBulbC + 237.3))
eswb = exp((16.78*wetBulbC - 116.9) / (wetBulbC + 237.3))
ed = eswb - a*airPressure*(dryBulbC - wetBulbC)
relativeHumidity = 100*ed/esdb
dewPointC = (116.9 + 237.3*ln(ed)) / (16.78 - ln(ed))
```

## Worked Example

No worked numeric example was included with the original program. As
an independently verifiable check: when the wet and dry bulb read
the same temperature, there is no evaporative cooling at all, so
relative humidity must be exactly 100% and the dew point must equal
the measured temperature exactly; confirmed algebraically (the vapor
pressure and dew point formulas are exact inverses of each other) and
in this conversion's tests.
