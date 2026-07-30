package main

import (
	"math"
	"testing"
)

// This worked example is taken directly from CONE.TXT's own sample
// run at the documented default inputs, with every printed value
// given.
func TestComputeConePattern_WorkedExample(t *testing.T) {
	got := computeConePattern(3, 5, 10, 0.25)

	tests := []struct {
		name string
		got  float64
		want float64
	}{
		{name: "IncludedAngleDeg", got: got.IncludedAngleDeg, want: 36.3915},
		{name: "SmallRadius", got: got.SmallRadius, want: 15.0748},
		{name: "SmallChord", got: got.SmallChord, want: 9.5096},
		{name: "SmallArcLength", got: got.SmallArcLength, want: 9.5748},
		{name: "LargeRadius", got: got.LargeRadius, want: 25.1247},
		{name: "LargeChord", got: got.LargeChord, want: 15.6911},
		{name: "LargeArcLength", got: got.LargeArcLength, want: 15.9580},
		{name: "EdgeLength", got: got.EdgeLength, want: 10.0499},
		{name: "ConeIncludedAngle", got: got.ConeIncludedAngle, want: 5.7106},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if diff := math.Abs(tt.got - tt.want); diff > 0.0001 {
				t.Errorf("%s = %v, want %v (per CONE.TXT)", tt.name, tt.got, tt.want)
			}
		})
	}
}

func TestComputeConePattern_EqualDiametersIsACylinderNotACone(t *testing.T) {
	// Equal small and large diameters describe a cylinder, which has
	// zero cone half-angle. Dividing by its sine (zero) sends the
	// pattern radii to +Inf; the original program does not guard
	// against this input either, so the contract here is that the
	// failure is visible as +Inf, not a panic or a silently
	// plausible-looking finite number.
	got := computeConePattern(3, 3, 10, 0.25)
	if !math.IsInf(got.LargeRadius, 1) {
		t.Errorf("LargeRadius = %v, want +Inf for the degenerate equal-diameter (cylinder) input", got.LargeRadius)
	}
}
