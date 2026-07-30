package main

import (
	"math"
	"testing"
)

func TestComputeGreatCircle_SamePointHasZeroDistance(t *testing.T) {
	got := computeGreatCircle(40, -74, 40, -74)
	if math.Abs(got.CentralAngleDeg) > 1e-9 {
		t.Errorf("CentralAngleDeg = %v, want 0", got.CentralAngleDeg)
	}
	if math.Abs(got.DistanceKm) > 1e-9 {
		t.Errorf("DistanceKm = %v, want 0", got.DistanceKm)
	}
}

func TestComputeGreatCircle_DistanceIsSymmetric(t *testing.T) {
	// The great circle distance between A and B must equal the
	// distance between B and A: an identity independent of which
	// point is labeled "first".
	ab := computeGreatCircle(33.7775436, -118.3769558, 51.5, 0.0)
	ba := computeGreatCircle(51.5, 0.0, 33.7775436, -118.3769558)
	if math.Abs(ab.DistanceKm-ba.DistanceKm) > 1e-6 {
		t.Errorf("distance A->B (%v) != distance B->A (%v)", ab.DistanceKm, ba.DistanceKm)
	}
}

func TestComputeGreatCircle_LosAngelesToLondonMatchesKnownDistance(t *testing.T) {
	// Los Angeles to London is a commonly cited great circle
	// distance, roughly 8,750-8,800 km depending on the exact mean
	// Earth radius used: an independent, real-world check on the
	// formula rather than a re-run of it.
	got := computeGreatCircle(33.7775436, -118.3769558, 51.5, 0.0)
	if got.DistanceKm < 8700 || got.DistanceKm > 8900 {
		t.Errorf("DistanceKm = %v, want roughly 8700-8900 km (LA to London)", got.DistanceKm)
	}
}

func TestRelativeAzimuthDeg_SameLongitudeGivesDueNorthOrSouth(t *testing.T) {
	// Two points on the same meridian must bear due north (0) or due
	// south (180) from each other, an explicit special case in the
	// original formula (avoiding a degenerate division when sinA is
	// tiny).
	north := relativeAzimuthDeg(10, 5, 20, 5) // point 2 is north of point 1
	if math.Abs(north) > 1e-9 {
		t.Errorf("relativeAzimuthDeg (point 2 north) = %v, want 0", north)
	}
	south := relativeAzimuthDeg(20, 5, 10, 5) // point 2 is south of point 1
	if math.Abs(south-180) > 1e-9 {
		t.Errorf("relativeAzimuthDeg (point 2 south) = %v, want 180", south)
	}
}

func TestNormalizeDeg_WrapsIntoRange(t *testing.T) {
	cases := []struct{ x, want float64 }{
		{190, -170},
		{-190, 170},
		{180, 180},
		{-180, -180},
		{45, 45},
	}
	for _, c := range cases {
		if got := normalizeDeg(c.x, -180, 180); math.Abs(got-c.want) > 1e-9 {
			t.Errorf("normalizeDeg(%v, -180, 180) = %v, want %v", c.x, got, c.want)
		}
	}
}
