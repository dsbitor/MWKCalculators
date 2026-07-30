package main

import (
	"math"
	"testing"
)

func TestComputeGearPair_DocumentedDefaultInput(t *testing.T) {
	// The documented default input matches John Cooper's own worked
	// example in "Spur Gears and Pinions" (Machinist's Workshop,
	// 4/99), evaluated against the ported formula directly.
	pair := computeGearPair(45, 20, 20, 20)

	checkClose(t, "GearRatio", pair.GearRatio, 45.0/20.0)
	checkClose(t, "CenterDistance", pair.CenterDistance, 0.5*(45.0/20.0+20.0/20.0))

	g1 := pair.Gears[0]
	checkClose(t, "PitchDiameter", g1.PitchDiameter, 45.0/20.0)
	checkClose(t, "OutsideDiameter", g1.OutsideDiameter, 47.0/20.0)
	checkClose(t, "Addendum", g1.Addendum, 1.0/20.0)
	checkClose(t, "Dedendum", g1.Dedendum, 1.2/20.0)
	checkClose(t, "WholeDepth", g1.WholeDepth, 2.2/20.0)
	checkClose(t, "CircularPitch", g1.CircularPitch, math.Pi/20.0)
}

func TestComputeGearPair_BaseCircleRadiusMatchesCosineIdentity(t *testing.T) {
	// The base circle radius is, by definition, the pitch radius
	// times the cosine of the pressure angle: an identity
	// independent of this code's own separate formula for pitch
	// diameter and base circle radius.
	pair := computeGearPair(45, 20, 20, 20)
	for i, g := range pair.Gears {
		pitchRadius := 0.5 * g.PitchDiameter
		want := pitchRadius * math.Cos(20*math.Pi/180)
		if math.Abs(g.BaseCircleRadius-want) > 1e-9 {
			t.Errorf("gear %d BaseCircleRadius = %v, want %v", i, g.BaseCircleRadius, want)
		}
	}
}

func TestComputeGearPair_CenterDistanceIsHalfSumOfPitchDiameters(t *testing.T) {
	// Two gears mesh with their pitch circles tangent to each other:
	// center distance is exactly half the sum of the two pitch
	// diameters, an independent geometric identity.
	pair := computeGearPair(60, 15, 24, 14.5)
	want := 0.5 * (pair.Gears[0].PitchDiameter + pair.Gears[1].PitchDiameter)
	checkClose(t, "CenterDistance", pair.CenterDistance, want)
}

func TestComputeGearPair_RatioIsOrderIndependent(t *testing.T) {
	// The gear ratio is defined as the larger tooth count over the
	// smaller, regardless of which gear is passed first.
	ab := computeGearPair(45, 20, 20, 20)
	ba := computeGearPair(20, 45, 20, 20)
	checkClose(t, "GearRatio", ab.GearRatio, ba.GearRatio)
}

func checkClose(t *testing.T, name string, got, want float64) {
	t.Helper()
	if math.Abs(got-want) > 1e-9 {
		t.Errorf("%s = %v, want %v", name, got, want)
	}
}
