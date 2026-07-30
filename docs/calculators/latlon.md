# latlon

Great circle distance and bearing between two points.

**Converted from:** `LATLON.C`, `MWKC/Misc/latlon.zip`
**Go source:** `MWKGo/latlon/latlon.go`

## Purpose

Given the latitude and longitude of two points on the Earth's
surface, computes the central angle between them, the great circle
distance (the shortest path along the surface of a sphere), and the
compass bearing (azimuth) from each point toward the other.

## Inputs

| Prompt | Default |
|---|---|
| Latitude of first point | 33.7775436 deg (Los Angeles) |
| Longitude of first point | -118.3769558 deg |
| Latitude of second point | 51.5 deg (London) |
| Longitude of second point | 0 deg |

Latitude is positive northward from the equator; longitude is
positive eastward from the prime meridian. Latitude must be between
-90 and 90; longitude is normalized into (-180, 180].

## Output

Central angle between the two points, great circle distance in
kilometers and miles, and the azimuth of each point as seen from the
other.

## Method

Spherical law of cosines:

```
cosAngle = sin(lat1)*sin(lat2) + cos(lat1)*cos(lat2)*cos(lon2-lon1)
angle = acos(cosAngle)
distanceKm = earthRadiusKm * angle(radians)
```

with a similar spherical-trigonometry formula for relative azimuth
(see the original `caz` function, ported directly).

The original's `constant.h` header (not included in the source
archive) supplied `REKM` (Earth radius in km) and `FPKM` (feet per
km); since that header isn't available and no exact output values
needed to be reproduced, this conversion substitutes the standard
mean Earth radius (6371 km) and the exact feet-per-kilometer
conversion factor, documented as such in the source rather than
silently assumed to match the original's own unknown constants.

## Worked Example

No worked numeric example was included with the original program. As
independently verifiable checks: the distance between a point and
itself is zero, the distance between two points doesn't depend on
which is labeled "first" (an identity independent of any specific
formula), and the Los Angeles/London default distance falls within
the commonly cited real-world range for that route (roughly
8,700-8,900 km, depending on the exact Earth radius assumed); all
confirmed in this conversion's tests.
