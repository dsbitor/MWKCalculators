// latlon computes the great circle distance and relative compass
// bearing (azimuth) between two points given by latitude and
// longitude, using the spherical law of cosines.
//
// Converted from LATLON.C, Misc/latlon. The original's constant.h
// header (not included in the archive) supplied REKM (Earth radius
// in km) and FPKM (feet per km); those constants are not otherwise
// recoverable, so this conversion substitutes the standard mean
// Earth radius (6371 km) and the exact feet-per-kilometer conversion
// factor, both documented below rather than silently assumed.
package main

import (
	"fmt"
	"math"
	"os"

	"mwkgo/internal/mwktrig"
	"mwkgo/internal/promptio"
)

// earthRadiusKm is the standard mean Earth radius, substituted for
// the original's unavailable REKM constant.
const earthRadiusKm = 6371.0

// feetPerKm is the exact number of feet in a kilometer, substituted
// for the original's unavailable FPKM constant.
const feetPerKm = 3280.839895

// feetPerMile is the standard number of feet in a mile, used to
// convert the great circle distance from kilometers to miles.
const feetPerMile = 5280.0

// normalizeDeg wraps x into the half-open interval [a, b), matching
// the original program's norm function, used to keep longitude
// within (-180, 180].
func normalizeDeg(x, a, b float64) float64 {
	span := b - a
	for x < a {
		x += span
	}
	for x > b {
		x -= span
	}
	return x
}

// centralAngleDeg returns the central angle, in degrees, between two
// points on a sphere given their latitudes and longitudes in
// degrees, via the spherical law of cosines.
func centralAngleDeg(lat1, lon1, lat2, lon2 float64) float64 {
	sLat1, cLat1 := math.Sin(lat1*math.Pi/180), math.Cos(lat1*math.Pi/180)
	sLat2, cLat2 := math.Sin(lat2*math.Pi/180), math.Cos(lat2*math.Pi/180)
	cosAngle := sLat1*sLat2 + cLat1*cLat2*math.Cos((lon2-lon1)*math.Pi/180)
	return mwktrig.ClampedAcosDeg(cosAngle)
}

// relativeAzimuthDeg returns the compass bearing, in degrees, of
// point 2 as seen from point 1, matching the original program's caz
// function (including its convention of taking longitude as
// positive westward internally, which is why main negates its own
// east-positive longitudes before calling this).
func relativeAzimuthDeg(lat1, lon1, lat2, lon2 float64) float64 {
	if lat1 == lat2 && lon1 == lon2 {
		return 0
	}

	sLat1, cLat1 := math.Sin(lat1*math.Pi/180), math.Cos(lat1*math.Pi/180)
	sLat2, cLat2 := math.Sin(lat2*math.Pi/180), math.Cos(lat2*math.Pi/180)

	cosA := sLat1*sLat2 + cLat1*cLat2*math.Cos(math.Abs(lon2-lon1)*math.Pi/180)
	a := mwktrig.ClampedAcosDeg(cosA)
	sinA := math.Sin(a * math.Pi / 180)

	cosB := (sLat2 - sLat1*cosA) / (cLat1 * sinA)
	azimuth := mwktrig.ClampedAcosDeg(cosB)

	if lon1 == lon2 {
		if lat1 <= lat2 {
			azimuth = 0
		} else {
			azimuth = 180
		}
	}
	if lon2 > lon1 {
		azimuth = 360 - azimuth
	}
	if math.Abs(lon2-lon1) > 180 {
		azimuth = 360 - azimuth
	}
	return azimuth
}

// greatCircleReport holds the computed relationship between two
// points on the Earth's surface.
type greatCircleReport struct {
	CentralAngleDeg    float64
	DistanceKm         float64
	DistanceMiles      float64
	AzimuthFrom1To2Deg float64
	AzimuthFrom2To1Deg float64
}

// computeGreatCircle returns the great circle distance and relative
// azimuths between two points given by latitude/longitude in
// degrees (east-positive longitude, north-positive latitude).
func computeGreatCircle(lat1, lon1, lat2, lon2 float64) greatCircleReport {
	angle := centralAngleDeg(lat1, lon1, lat2, lon2)
	distanceKm := earthRadiusKm * angle * math.Pi / 180

	return greatCircleReport{
		CentralAngleDeg:    angle,
		DistanceKm:         distanceKm,
		DistanceMiles:      distanceKm * feetPerKm / feetPerMile,
		AzimuthFrom1To2Deg: relativeAzimuthDeg(lat1, -lon1, lat2, -lon2),
		AzimuthFrom2To1Deg: relativeAzimuthDeg(lat2, -lon2, lat1, -lon1),
	}
}

func main() {
	prompter, err := promptio.New(os.Stdin, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, "latlon:", err)
		os.Exit(1)
	}

	fmt.Println("Latitude is measured positive northward from equator.")
	fmt.Println("Longitude is measured positive eastward from prime meridian.")
	fmt.Println()

	var lat1, lat2 float64
	for {
		lat1 = prompter.Float("Latitude of first point (deg)", 33.7775436)
		if lat1 >= -90 && lat1 <= 90 {
			break
		}
		fmt.Println("-90 <= latitude <= 90")
	}
	lon1 := normalizeDeg(prompter.Float("Longitude of first point (deg)", -118.3769558), -180, 180)

	for {
		lat2 = prompter.Float("Latitude of second point (deg)", 51.5)
		if lat2 >= -90 && lat2 <= 90 {
			break
		}
		fmt.Println("-90 <= latitude <= 90")
	}
	lon2 := normalizeDeg(prompter.Float("Longitude of second point (deg)", 0.0), -180, 180)

	report := computeGreatCircle(lat1, lon1, lat2, lon2)

	fmt.Println()
	fmt.Printf("lat1, lon1 = %.2f, %.2f\n", lat1, lon1)
	fmt.Printf("lat2, lon2 = %.2f, %.2f\n", lat2, lon2)
	fmt.Printf("Central angle between locations = %.2f deg\n", report.CentralAngleDeg)
	fmt.Printf("Great circle distance = %.2f km = %.2f miles\n", report.DistanceKm, report.DistanceMiles)
	fmt.Printf("Azimuth of point 2 from point 1 = %.2f deg\n", report.AzimuthFrom1To2Deg)
	fmt.Printf("Azimuth of point 1 from point 2 = %.2f deg\n", report.AzimuthFrom2To1Deg)
}
