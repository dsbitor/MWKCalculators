package main

import (
	"math"
	"testing"
)

func TestSecondCylinderDiameter(t *testing.T) {
	// The original author's own worked example (SINE.TXT): building
	// a 10 degree sine bar with a 3/8in first cylinder and a 3in
	// link spacing needs a second cylinder of 0.898in (the author's
	// own rounding of the exact value).
	got := secondCylinderDiameter(0.375, 3.0, 10.0)
	want := 2*3.0*math.Sin(5*math.Pi/180) + 0.375
	if diff := math.Abs(got - want); diff > 1e-9 {
		t.Errorf("secondCylinderDiameter(0.375, 3.0, 10.0) = %v, want %v", got, want)
	}
	if diff := math.Abs(got - 0.898); diff > 0.001 {
		t.Errorf("secondCylinderDiameter(0.375, 3.0, 10.0) = %v, want approximately 0.898 per SINE.TXT", got)
	}
}

func TestSecondCylinderDiameter_ZeroAngle(t *testing.T) {
	// A zero angle needs no size difference between the cylinders at
	// all.
	got := secondCylinderDiameter(0.375, 3.0, 0)
	if diff := math.Abs(got - 0.375); diff > 1e-9 {
		t.Errorf("secondCylinderDiameter(0.375, 3.0, 0) = %v, want 0.375", got)
	}
}

func TestFits(t *testing.T) {
	tests := []struct {
		name                                          string
		firstDiameter, secondDiameter, centerDistance float64
		want                                          bool
	}{
		{name: "cylinders comfortably smaller than the link spacing fit", firstDiameter: 0.375, secondDiameter: 0.898, centerDistance: 3.0, want: true},
		{name: "cylinders summing to exactly twice the spacing still fit", firstDiameter: 2, secondDiameter: 4, centerDistance: 3, want: true},
		{name: "cylinders too large for the spacing do not fit", firstDiameter: 5, secondDiameter: 5, centerDistance: 3, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fits(tt.firstDiameter, tt.secondDiameter, tt.centerDistance); got != tt.want {
				t.Errorf("fits(%v, %v, %v) = %v, want %v", tt.firstDiameter, tt.secondDiameter, tt.centerDistance, got, tt.want)
			}
		})
	}
}

func TestAngleSensitivity_DocumentedWorkedExample(t *testing.T) {
	// SINE.TXT states, for this exact example, that a 0.001in error
	// in either cylinder's diameter causes about 0.01 degrees of
	// angle error, and a 0.001in error in the link spacing causes
	// about 0.002 degrees. These are the author's own rounded
	// figures, so the comparison tolerance is loose (checking the
	// right order of magnitude and the right first significant
	// digit), not an exact match.
	d1, centerDistance, angle := 0.375, 3.0, 10.0
	d2 := secondCylinderDiameter(d1, centerDistance, angle)

	perDiameterError, perDistanceError := angleSensitivity(d1, d2, centerDistance)
	degreesPerRadian := 180 / math.Pi

	gotDiameterErrorDeg := 0.001 * degreesPerRadian * perDiameterError
	gotDistanceErrorDeg := 0.001 * degreesPerRadian * perDistanceError

	if diff := math.Abs(gotDiameterErrorDeg - 0.01); diff > 0.001 {
		t.Errorf("angle error per 0.001in diameter error = %v deg, want approximately 0.01 per SINE.TXT", gotDiameterErrorDeg)
	}
	if diff := math.Abs(gotDistanceErrorDeg - 0.002); diff > 0.001 {
		t.Errorf("angle error per 0.001in distance error = %v deg, want approximately 0.002 per SINE.TXT", gotDistanceErrorDeg)
	}
}

func TestAngleSensitivity_EqualDiameters_MinimizesDistanceSensitivity(t *testing.T) {
	// When the two cylinders are the same diameter (a zero degree
	// sine bar), the angle has no sensitivity at all to errors in
	// the link spacing, since halfDiameterDifferenceRatio is zero:
	// an identity independent of this code.
	_, perDistanceError := angleSensitivity(0.5, 0.5, 3.0)
	if perDistanceError != 0 {
		t.Errorf("perDistanceError = %v, want 0 when both cylinders are the same diameter", perDistanceError)
	}
}
