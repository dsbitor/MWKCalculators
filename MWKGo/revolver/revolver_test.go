package main

import (
	"math"
	"testing"
)

func TestCylinderDimensions(t *testing.T) {
	tests := []struct {
		name                                                string
		holeCount, holeDiameter, edgeSpacing, wallThickness float64
		wantHoleRadius                                      float64
	}{
		// For six holes, adjacent hole centers 60 degrees apart lie
		// at a chord distance equal to the placement radius itself
		// (the hexagon side-equals-radius identity, independent of
		// this code). That chord is exactly edgeSpacing+holeDiameter
		// (one hole diameter plus the gap between edges), so the
		// placement radius must equal edgeSpacing+holeDiameter
		// exactly.
		{
			name:      "six holes matches the hexagon side-equals-radius identity",
			holeCount: 6, holeDiameter: 0.25, edgeSpacing: 0.5, wallThickness: 0.25,
			wantHoleRadius: 0.75,
		},
		{
			name:      "hexagon identity holds for different hole and spacing values",
			holeCount: 6, holeDiameter: 0.1, edgeSpacing: 0.2, wallThickness: 0.1,
			wantHoleRadius: 0.3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotHoleRadius, _ := cylinderDimensions(tt.holeCount, tt.holeDiameter, tt.edgeSpacing, tt.wallThickness)
			if diff := math.Abs(gotHoleRadius - tt.wantHoleRadius); diff > 1e-9 {
				t.Errorf("holeRadius = %v, want %v", gotHoleRadius, tt.wantHoleRadius)
			}
		})
	}
}

func TestCylinderDimensions_DocumentedDefaultInput(t *testing.T) {
	// The documented default input, evaluated against the ported
	// formula directly using math.Sin rather than a hand-typed
	// literal.
	holeCount, holeDiameter, edgeSpacing, wallThickness := 6.0, 0.25, 0.5, 0.25

	anglePerHole := 360 / holeCount
	wantHoleRadius := (edgeSpacing + holeDiameter) / (2 * math.Sin(0.5*anglePerHole*math.Pi/180))
	wantCylinderDiameter := 2 * (wantHoleRadius + 0.5*holeDiameter + wallThickness)

	gotHoleRadius, gotCylinderDiameter := cylinderDimensions(holeCount, holeDiameter, edgeSpacing, wallThickness)
	if diff := math.Abs(gotHoleRadius - wantHoleRadius); diff > 1e-9 {
		t.Errorf("holeRadius = %v, want %v", gotHoleRadius, wantHoleRadius)
	}
	if diff := math.Abs(gotCylinderDiameter - wantCylinderDiameter); diff > 1e-9 {
		t.Errorf("cylinderDiameter = %v, want %v", gotCylinderDiameter, wantCylinderDiameter)
	}
}

func TestCylinderDimensions_OneHole_ProducesNonFiniteResult(t *testing.T) {
	// With a single hole, the angle per hole is 360 degrees, whose
	// half-angle sine is sin(180)=0 (up to floating-point error),
	// making the denominator effectively zero. The original program
	// does not guard against this either: a "cylinder" with one
	// hole is outside what this construction means, so the failure
	// is a very large or infinite result rather than a silently
	// plausible one.
	holeRadius, _ := cylinderDimensions(1, 0.25, 0.5, 0.25)
	if holeRadius < 1e6 {
		t.Errorf("holeRadius = %v, want an extremely large value for a single hole", holeRadius)
	}
}
