# mix

Mixed dimensional units four-function calculator.

**Converted from:** `MIX.C` (M. W. Klotz), `MWKC/Math/mix.zip`
**Go source:** `MWKGo/mix/mix.go`

## Purpose

A four-function calculator that tracks a single running,
dimensioned result in six unit systems at once (meter, centimeter,
millimeter, yard, foot, inch), like having six calculator screens
side by side. Every entry updates all six accumulators together, so
the running result is always available in whichever unit is wanted,
and a non-dimensional mode allows using it as an ordinary
four-function calculator.

## Inputs

A free-form command line, entered repeatedly until `quit` or `exit`:

```
[op] entry [ud]
```

- `op` is `+`, `-`, `*`, or `/` (default `+` if omitted)
- `entry` is a number: `7`, `.2`, `1.23`, `2/3`, or `1 & 2/3`
- `ud` is a unit designator: `nd`, `m`, `cm`, `mm`, `yd`, `ft`, `in`
  (default the current default unit if omitted); `nd` marks a
  non-dimensional constant

Commands (also entered on the same prompt):

| Command | Effect |
|---|---|
| `clear` | zero all six accumulators |
| `decimals=n` | set the number of decimal places displayed (0-15) |
| `dflag=n` | `0` = fixed decimal places, `1` = self-adjusting |
| `dtype=n` | `0` = fixed, `1` = engineering, `2` = scientific display |
| `fracd=n` | set the fractional-inch display precision to 1/n (default 64) |
| `scale=n` | set an input scale factor, applied to every dimensioned entry |
| `unit=ud` | set the default unit used when an entry omits one |
| `undo` | restore the accumulators to their state before the last entry |
| `help` | redisplay the syntax and command summary |
| `quit` / `exit` | exit the program |

## Output

The six accumulators (meter, centimeter, millimeter, yard, foot,
inch), each formatted per the current decimal/notation settings and
suffixed with its unit and any dimension exponent (e.g. `in^2` after
a squaring operation); followed, whenever the result is a plain
length with a non-zero feet accumulator, by a mixed feet-and-fraction
line such as `3 ft 9 & 25/32 in`.

## Method

Every entry is converted to inches (the program's internal base
unit) and, if dimensioned, multiplied by the scale factor; the
converted value is then applied to all six accumulators together,
each scaled by its own units-per-inch factor, so they always agree
on the same physical quantity. A running dimension exponent tracks
whether the result is a plain length (exponent 1), an area (2), a
volume (3), and so on, as entries are multiplied or divided by other
dimensioned quantities; adding or subtracting a non-dimensional value
into an already-dimensioned accumulator is rejected, since that
would not be physically meaningful.

The number of decimal places shown is either fixed at the configured
`decimals` setting or, in self-adjusting mode, grown just far enough
that the displayed value reproduces the original to within double
precision (about 15 significant digits) when read back — this is
the same technique used to show `3.790551181102362 in` in the worked
example below, since 96.28 mm does not terminate in a short decimal
number of inches. Magnitudes far outside normal fixed-point range
(log10 magnitude beyond 15) automatically switch to engineering
notation regardless of the configured display type, matching the
original's own override, so that an extreme value is never shown as
an unreadably long string of digits or as all zeros.

The original program's own unit-designator parser copies at most two
raw bytes following the number for the unit code, which is exactly
right for its two-letter unit abbreviations but is fragile for
anything else reaching that code path. This conversion instead takes
the entry's full trailing run of letters as the unit and matches it
exactly against the six known unit names, an internal simplification
with the same result for every legitimate entry.

The mixed feet-and-inches display reduces the fractional remainder
of the inches accumulator to the nearest 1/`fracd` using the
program's own greatest-common-divisor reduction. Because the feet
and inches accumulators are independent running totals rather than
one derived from the other, the reduced fraction can disagree
slightly with the inches accumulator's own value; the displayed
"error" percentage is that disagreement, computed exactly as the
original does, not a warning about anything having gone wrong.

## Worked Example

Entering `+3.2in` then `+1.5cm` (default unit `in`, `decimals=6`)
produces:

```
0.096280 m
9.628000 cm
96.280000 mm
0.315879 ft
3.790551 in
 3 & 25/32 in (error = -0.2454 %)
```

Switching to self-adjusting decimals (`dflag=1`) instead shows each
accumulator to just enough places to round-trip exactly:

```
0.09628 m
9.628 cm
96.28 mm
0.315879265091864 ft
3.790551181102362 in
 3 & 25/32 in (error = -0.2454 %)
```

The accumulator values and the mixed feet-and-fraction display are
both reproduced exactly by this conversion's tests (the inch and
meter accumulators are checked against `3.2 + 1.5/2.54` computed
independently, rather than against the doc's rounded digits). The
`-0.2454%` error figure shown above is independently verifiable as
`(3.78125 − 3.2 − 1.5/2.54) / (3.2 + 1.5/2.54)`, the fractional
display's own deviation from the precise accumulator value, though it
is not separately asserted in this conversion's own test suite.
Squaring a length (multiplying a 5.25 in accumulator by itself) is
checked to double the dimension exponent to 2 and to match
`5.25 * 5.25` converted to each unit's own squared scale factor.
