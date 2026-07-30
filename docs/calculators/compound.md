# compound

Lathe compound rest angle for a given infeed ratio.

**Converted from:** `COMPOUND.C` (M. W. Klotz, 6/01, 1/05),
`MWKC/WorkshopUtilities/compound.zip`
**Go source:** `MWKGo/compound/compound.go`

## Purpose

Angling a lathe's compound slide makes a given slide movement
produce a smaller movement of the tool toward the work. A
popular choice is an angle giving an infeed-to-slide ratio of
0.1, so the slide's graduated dial reads directly in 0.0001in of
tool movement. Given the desired ratio (which must be no more
than 1), this program computes the required compound angle and
its complement, since lathe manufacturers aren't consistent
about which one their protractor is calibrated against.

## Inputs

| Prompt | Default |
|---|---|
| Required movement ratio | 0.1 |

Re-prompts if the ratio entered is greater than 1.

## Output

The compound angle and its complement, each in decimal degrees
and in degrees/minutes/seconds.

## Method

```
angle = asin(ratio)
```

## Worked Example

No worked example was included with the original program. As an
independently verifiable check: an angle and its complement
always sum to 90 degrees, confirmed directly in this
conversion's tests.
