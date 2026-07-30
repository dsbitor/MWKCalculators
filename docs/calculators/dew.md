# dew

Dew point temperature.

**Converted from:** `DEW.C` (M. W. Klotz, 3/04), `MWKC/Misc/dew.zip`
**Go source:** `MWKGo/dew/dew.go`

## Purpose

Metal tools left below the dew point accumulate condensation and
rust. Given the ambient temperature and relative humidity, this
program computes the dew point temperature, so tools can be kept
warmed above it. Valid for ambient temperatures between 0 and 60
degC (32 to 140 degF) and relative humidity between 1 and 100
percent, per the original program's own stated range; neither
the original nor this conversion validates inputs against that
range.

## Inputs

| Prompt | Default |
|---|---|
| Preferred scale | Fahrenheit (default) or Celsius |
| Measured temperature | 70 degF (or 20 degC if Celsius chosen) |
| Relative humidity | 60% |

## Output

Dew point temperature in both Celsius and Fahrenheit, and the
method's stated accuracy (+/- 0.04 degC).

## Method

The Magnus-Tetens approximation, using the standard
August-Roche-Magnus constants for Celsius (a=17.27, b=237.7):

```
alpha = a*tempC/(b+tempC) + ln(relativeHumidity)
dewPointC = b*alpha / (a-alpha)
```

## Worked Example

No worked example was included with the original program. As an
independently verifiable check: a commonly cited example for the
Magnus-Tetens formula, 20 degC at 50% relative humidity, gives a
dew point of approximately 9.3 degC, which this formula
reproduces closely. Also, 100% relative humidity makes the dew
point equal the ambient temperature exactly, since `ln(1)` is
zero.
