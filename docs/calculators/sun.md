# sun

Solar position, equation of time, sunrise/sunset, and sundial geometry.

**Converted from:** `SUN.C` (M. W. Klotz), `MWKC/Misc/sun.zip`
**Go source:** `MWKGo/sun/sun.go`
**Original documentation:** `SUN.TXT`, inside `MWKC/Misc/sun.zip` (not included in this conversion)

## Purpose

Computes the sun's position (right ascension, declination, distance,
altitude, azimuth), the equation of time, sunrise and sunset, and
shadow-angle geometry for two standard sundial orientations, for a
chosen date, time, and observer location. It also finds specific
solar events by search: the equinoxes and solstices, the day of the
year with the smallest equation of time in a given month, and the UT
at which the sun crosses the local meridian on a given date.

The underlying algorithm is a long-term (1800-2100) trigonometric
series from the 1978 "Almanac for Computers" (Nautical Almanac
Office, US Naval Observatory) — the same one described in the
original archive's own `SUN.TXT`, which quotes the Almanac's stated
accuracy: sun position to 1E-4 AU, right ascension to 0.1 minute,
declination to 1 arc minute.

## Data setup

An observer's latitude, longitude, and UT offset(s) belong to one
location, reused across many runs from that place — not universal
reference data, and not equipment configuration, but the same
"specific to one owner's own setup" reasoning that puts `gearatio`,
`change`, and `ddh`'s data in `userdata.db` rather than
`reference.db`; see `ai/plans/c-to-go-conversion-plan.md`, "Data-file
strategy for Tier 2". Import your location before first use:

```
sun -import my-location.csv
```

The CSV needs `id` (always `1`), `latitude_deg` (positive north),
`longitude_deg` (positive **west** of Greenwich, matching the
original's own convention), `standard_time_offset_hr`, and
`daylight_time_offset_hr` (both measured negative westward from UT;
use the same value for both if your area doesn't observe daylight
saving). A worked example built from the original archive's own
shipped `SUN.DAT` (Los Angeles/Long Beach area) ships at
`MWKGo/sun/testdata/example-settings.csv`.

## Inputs

A repeating menu of single-letter options:

| Option | Action |
|---|---|
| a | Use the current UTC date/time |
| b | Enter a date/time directly |
| e | Find the day in a given month with the smallest equation of time |
| f | Find the fall equinox in a given year |
| m | Find the meridian transit time on a given date |
| s | Find the summer solstice in a given year |
| v | Find the vernal equinox in a given year |
| w | Find the winter solstice in a given year |
| h | Show the menu again |
| q | Quit |

## Output

Date/time (UT, plus standard and daylight-saving local time), right
ascension, declination, distance (AU and miles), altitude, zenith
angle, azimuth, semi-diameter, sidereal times, the equation of time
(with a "dial slow"/"dial fast" note), sunrise, sunset, and the
daylight fraction of the day. If the sun is above the horizon, two
sundial geometries follow: a horizontal sundial with a vertical
gnomon, and a south-facing vertical sundial with a horizontal gnomon
— each reporting the shadow's angle relative to the plane's own
reference directions and the shadow-length-to-gnomon-length ratio.

## Method

`suncomp` is a direct, formula-for-formula port of `SUN.C`'s own
function of the same name — the mean anomalies and longitudes of the
sun, moon, and several planets, combined into the sun's true anomaly
components, right ascension, and declination. `tnorm` normalizes a
time-of-day value outside `[0,24)` by carrying whole days into the
month/year fields (correctly crossing month and year boundaries, and
respecting leap years), matching the original's own date-rollover
logic exactly, including its quirk of fixing the carry direction once
from the input's initial sign rather than re-checking it every step
(harmless, since the direction never needs to change mid-normalization
for a single call).

The four solstice/equinox searches share one 96-hour (4-day) hourly
scan starting at UT=12 on a fixed date (September 19 for the fall
equinox, June 19 for the summer solstice, and so on) — each differing
only in what "better" means (smallest `|declination|` for the
equinoxes, largest declination for the summer solstice, smallest
signed declination for the winter solstice). The equation-of-time
minimum search instead scans every day of one chosen month
independently. The meridian transit search is a secant-method
root-find on azimuth versus target (180° for northern latitudes, 0°
for southern), bounded by an explicit iteration cap per
`coding-style.md` Rule 2 in place of the original's unbounded
convergence loop (in practice it converges in a handful of
iterations).

This conversion uses Go's own `time.Now().UTC()` for the "current
time" option rather than porting `SUN.C`'s own `dtinit`/`jdnum`
machinery (a DOS-era system-clock read plus a from-scratch Julian date
calculation whose result, `dt.jd`, the original never actually uses
anywhere besides computing itself) — a direct, simpler equivalent for
a value nothing downstream depends on.

Two known original quirks are preserved rather than "fixed": `caz`'s
supporting `sin(local hour angle)` value is computed but never
actually used in the altitude formula (only its cosine is), matching
the original's own unused calculation; and a location's raw computed
sunrise/sunset UT can come out above 24 or below 0 (the longitude
correction is added *after* the value is already wrapped into
`[0,24)`), which the original's own display does nothing to correct
either — the demonstration run below shows a sunset UT of `27:07:18`,
which is correct and simply means "past midnight in UT terms," not a
bug.

## Worked Example

No worked numeric example with expected output was included with the
original program; `SUN.TXT` only documents the underlying algorithm's
source and accuracy claims. As independently verifiable checks, this
conversion's tests confirm: the summer solstice's declination lands
within 0.1° of the sun's well known +23.44° obliquity maximum (and
the winter solstice within 0.1° of −23.44°); both equinoxes land
within 0.5° of 0° declination; distance stays within 0.98-1.02 AU
(Earth's orbital eccentricity is small); a northern-latitude location
has more daylight at the summer solstice than the winter solstice;
the meridian transit search converges to the requested target azimuth
within its own 0.001° tolerance; a sundial with the sun at zenith
above a horizontal plane and a vertical gnomon reports a
shadow-length ratio of essentially zero (a gnomon under a sun
directly overhead casts no shadow); and `tnorm` is checked directly
against known calendar rollovers, including a leap-year February 29.
A manual run using `SUN.DAT`'s own Los Angeles-area coordinates for
June 20, 2024 shows a declination of 23.44°, an altitude of 72.47° at
21:00 UT, and 60.01% daylight — all consistent with expectations for
a summer date near that latitude.
