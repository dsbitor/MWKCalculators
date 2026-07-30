# chill

Wind chill equivalent temperature.

**Converted from:** `CHILL.C` (M. W. Klotz, 11/98, 12/01),
`MWKC/WorkshopUtilities/chill.zip`
**Go source:** `MWKGo/chill/chill.go`

## Purpose

Computes the equivalent wind chill temperature from air
temperature and wind speed, using the US National Weather
Service's 2001 formula. That formula replaced an older
1945-style formula and its associated lookup table; the original
source carries that older table as a comment for reference, but
it corresponds to the superseded formula, not the one this
program actually runs, and so is not carried forward here.

No `.TXT` file was included with the original program; this
purpose statement is drawn from the `.C` file's own header
comment and its embedded formula history.

## Inputs

| Prompt | Default |
|---|---|
| Temperature | 30.0 degF |
| Wind speed | 25.0 mph |

## Output

Equivalent wind chill temperature in degrees Fahrenheit.

## Method

```
windFactor = windMph ^ 0.16
windChill  = 35.74 + 0.6215*tempF - windFactor*(35.75 - 0.4275*tempF)
```

This is the NWS's published formula, valid for temperatures at
or below 50F and wind speeds at or above 3mph; neither the
original program nor this conversion validates inputs against
that range.

## Worked Example

No worked example was included with the original program. As an
independently verifiable check: a commonly published NWS
reference example, 0F with a 15mph wind, gives a wind chill of
approximately -19F, which this formula reproduces closely.
