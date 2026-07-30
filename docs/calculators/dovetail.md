# dovetail

Dovetail measurement using pins.

**Converted from:** `DOVETAIL.C` (M. W. Klotz),
`MWKC/WorkshopUtilities/dovetail.zip`
**Go source:** `MWKGo/dovetail/dovetail.go`

## Purpose

One way to measure a dovetail is to put cylindrical pins into the
angles of the dovetail and measure across (male, external dovetail)
or between (female, internal dovetail) the pins with calipers.
Knowing the dovetail angle and pin diameter, this program derives the
measurements across the top and bottom of the dovetail from that pin
measurement. Any consistent unit system may be used; output units
match the input units.

## Inputs

| Prompt | Default |
|---|---|
| (M)ale or (F)emale dovetail | (required, re-prompts until valid) |
| Dovetail angle | 60 deg |
| Pin diameter | 0.375 |
| Height (male) or depth (female) of dovetail | 0.5 |
| Measurement across pins (male) or between pins (female) | 2.5 (male) / 1.0283 (female) |

## Output

Measurement across the top of the dovetail, and across the bottom.

## Method

```
verticalOffset = height / tand(angle)
radiusFactor = 1 + 1/tand(0.5*angle)

male:   bottom = acrossPins - pinDiameter*radiusFactor
        top    = bottom + 2*verticalOffset

female: bottom = betweenPins + pinDiameter*radiusFactor
        top    = bottom - 2*verticalOffset
```

## Worked Example

No worked numeric example was included with the original program.
The original program's own default inputs, however, form a matched
pair: a male dovetail measured at the default 2.5 across-pins
separation and a female dovetail measured at the default 1.0283
between-pins separation describe the same physical dovetail, so the
male program's top and bottom measurements must equal the female
program's bottom and top measurements (swapped) respectively. This
conversion verifies that cross-check directly in its tests, a
stronger check than re-deriving either formula alone since it also
confirms the male and female formulas are consistent with each
other.
