# doomsday

Day of the week for any date, computed two independent ways.

**Converted from:** `DOOMSDAY.C` (M. W. Klotz, 5/01),
`MWKC/Misc/doomsday.zip`
**Go source:** `MWKGo/doomsday/doomsday.go`

## Purpose

The "doomsday algorithm" is a mental technique for working out
the day of the week for any date, based on the observation that
several fixed dates each year (4/4, 6/6, 8/8, 10/10, 12/12,
and, via the mnemonic "I work from 9 to 5 in the 7-11", 5/9,
9/5, 7/11, 11/7) all fall on the same weekday, the year's
"doomsday". This program computes the weekday for an input date
two ways: once using a Julian day number calculation, and once
using the doomsday algorithm, and shows that they agree.

## Inputs

| Prompt | Default |
|---|---|
| month | today's month |
| day | today's day |
| year | today's year |

## Output

The Julian day number and day of the year for the input date,
the weekday computed from the Julian day, the year's doomsday
weekday, and the weekday computed from the doomsday algorithm
for the input date.

## Method

`julianDay` and `fromJulianDay` convert between a calendar date
and a Julian day number, switching from the Julian to the
Gregorian calendar at 14 September 1752 (the UK and US adoption
date). `doomsdayWeekday` and `monthDoomsday` implement Conway's
doomsday algorithm directly. The doomsday algorithm assumes the
Gregorian calendar throughout, so it is only expected to agree
with the Julian day method for dates from 1752 onward.

Converting this program surfaced two bugs in the original: the
returned year from a Julian day number was never adjusted back
to the "no year zero" display convention the forward conversion
uses (1 BC round-tripped as year 0, not year -1), and the
doomsday weekday formula could produce an out-of-range index or
a negative result for years before year 1. Both are fixed and
documented in `doomsday.go` at the point of the fix.

## Worked Example

From the original source comments, verified by the author by
hand against the Julian date technique:

| Date | Julian day | Weekday | Day of year |
|---|---|---|---|
| 28 December 1 BC | 1721420 | Tuesday | 362 |
| 1 January 1 AD | 1721424 | Saturday | 0 |
| 7 December 1941 | 2430336 | Sunday | 340 |
