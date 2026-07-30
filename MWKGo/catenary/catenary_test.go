package main

import (
	"math"
	"testing"
)

func TestCatenaryShape(t *testing.T) {
	tests := []struct {
		name                       string
		tension, density, distance float64
		wantDroop, wantLength      float64
	}{
		{
			name:    "zero span has no droop and no length",
			tension: 1, density: 1, distance: 0,
			wantDroop: 0, wantLength: 0,
		},
		{
			// tension=density=1 makes param=1, so the result reduces
			// to cosh(1)-1 and 2*sinh(1): well known constants,
			// independent of this code.
			name:    "unit parameter uses cosh(1) and sinh(1) directly",
			tension: 1, density: 1, distance: 2,
			wantDroop: math.Cosh(1) - 1, wantLength: 2 * math.Sinh(1),
		},
		{
			name:    "negative distance mirrors the droop and negates the length",
			tension: 1, density: 1, distance: -2,
			wantDroop: math.Cosh(1) - 1, wantLength: -2 * math.Sinh(1),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			droop, length := catenaryShape(tt.tension, tt.density, tt.distance)
			if diff := math.Abs(droop - tt.wantDroop); diff > 1e-9 {
				t.Errorf("droop = %v, want %v", droop, tt.wantDroop)
			}
			if diff := math.Abs(length - tt.wantLength); diff > 1e-9 {
				t.Errorf("length = %v, want %v", length, tt.wantLength)
			}
		})
	}
}

func TestCatenaryShape_ZeroDensityProducesNonFiniteResult(t *testing.T) {
	// The original program does not guard against a weightless
	// cable either; the contract here is that the failure is
	// immediately visible (NaN or Inf), not a silently wrong number.
	droop, length := catenaryShape(1, 0, 1)
	if !math.IsInf(droop, 0) && !math.IsNaN(droop) {
		t.Errorf("droop = %v, want +Inf, -Inf, or NaN when density is zero", droop)
	}
	if !math.IsInf(length, 0) && !math.IsNaN(length) {
		t.Errorf("length = %v, want +Inf, -Inf, or NaN when density is zero", length)
	}
}
