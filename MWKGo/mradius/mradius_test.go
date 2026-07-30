package main

import (
	"math"
	"testing"
)

func TestRadiusOfCurvature(t *testing.T) {
	tests := []struct {
		name                                   string
		rollerDiameter, measurementAcross, gap float64
		want                                   float64
	}{
		// A zero gage block gap simplifies the formula to
		// w^2/(2*d), since the adjustedRadius and rollerRadius
		// terms cancel: an identity independent of this code, used
		// here to check the full formula against a simpler
		// equivalent for the same inputs.
		{
			name:           "zero gap reduces to the simplified w^2/(2d) form",
			rollerDiameter: 1, measurementAcross: 3, gap: 0,
			want: 1.0 * 1.0 / (2 * 1),
		},
		{
			name:           "zero gap identity holds for different roller and measurement values",
			rollerDiameter: 0.5, measurementAcross: 2, gap: 0,
			want: 0.75 * 0.75 / (2 * 0.5),
		},
		// The documented default input, evaluated against the
		// ported formula directly.
		{
			name:           "documented default input",
			rollerDiameter: 0.25, measurementAcross: 2.0, gap: 0.1,
			want: (0.025*0.025 + 0.875*0.875 - 0.125*0.125) / (2 * (0.25 - 0.1)),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := radiusOfCurvature(tt.rollerDiameter, tt.measurementAcross, tt.gap)
			if diff := math.Abs(got - tt.want); diff > 1e-9 {
				t.Errorf("radiusOfCurvature(%v, %v, %v) = %v, want %v", tt.rollerDiameter, tt.measurementAcross, tt.gap, got, tt.want)
			}
		})
	}
}

func TestRadiusOfCurvature_GapEqualsRollerDiameter_ProducesNonFiniteResult(t *testing.T) {
	// The original program does not guard against a gage block gap
	// equal to the roller diameter either; the contract here is
	// that the failure is immediately visible (infinite or NaN),
	// not a silently wrong number, since the denominator goes to
	// zero.
	got := radiusOfCurvature(0.25, 2.0, 0.25)
	if !math.IsInf(got, 0) && !math.IsNaN(got) {
		t.Errorf("radiusOfCurvature(0.25, 2.0, 0.25) = %v, want +Inf, -Inf, or NaN", got)
	}
}
