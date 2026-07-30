# temp

Temperature scale converter.

**Converted from:** `TEMP.C`, `MWKC/Misc/temp.zip`
**Go source:** `MWKGo/temp/temp.go`

## Purpose

Converts a temperature between Centigrade, Fahrenheit, Kelvin,
Rankine, and Reaumur, reading free-form input such as `100f` or
`37.5c` (a number with a trailing scale letter; no letter, or an
unrecognized one, means Fahrenheit) and reporting the equivalent in
all five scales. Runs in a loop until `q` is entered, or, given a
command-line argument, converts that one value and exits.

## Inputs

A line of the form `123.45x`, where `x` is one of `c`, `f`, `k`, `r`,
or `e` (case-insensitive; omit for Fahrenheit), or `q` to quit.

## Output

The equivalent temperature in all five scales.

## Method

```
r0 = 1.8*273.18 - 32      (chosen so 0 Kelvin = 0 Rankine)

from Centigrade: F = 1.8*C+32; R = F+r0; K = C+273.18; Re = C*0.8
from Fahrenheit: R = F+r0; C = (F-32)/1.8; K = C+273.18; Re = C*0.8
from Kelvin:     C = K-273.18; F = 1.8*C+32; R = F+r0; Re = C*0.8
from Rankine:    F = R-r0; C = (F-32)/1.8; K = C+273.18; Re = C*0.8
from Reaumur:    C = Re/0.8; F = 1.8*C+32; R = F+r0; K = C+273.18
```

Absolute-scale inputs (Kelvin, Rankine) below zero are rejected and
re-prompted, since a negative absolute temperature is physically
meaningless.

Input parsing mimics C's `atof`: a leading numeric prefix is parsed
and any trailing characters (the scale letter) are ignored, rather
than requiring the whole input to be numeric.

## Worked Example

No worked numeric example was included with the original program. As
independently verifiable checks: water's known freezing point (0°C =
32°F) and boiling point (100°C = 212°F), and the definition of
absolute zero (0 K = 0 R); confirmed in this conversion's tests,
along with an internal consistency check that converting the same
physical temperature starting from any of the five scales lands on
the same equivalent values.
