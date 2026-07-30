package main

import (
	"math"
	"testing"
)

// These worked values are taken directly from OSBORNE.TXT's own
// example run (workpiece diameter 2, initial offset 0.1), which the
// original author printed as evidence of how fast the maneuver
// converges. The tolerance accounts for the .TXT's 8-decimal-place
// rounding.
func TestOsborneStep_WorkedExample(t *testing.T) {
	tests := []struct {
		name            string
		offset          float64
		wantNextOffset  float64
		wantRadialError float64
	}{
		{name: "iteration 1", offset: 0.1, wantNextOffset: 0.00501256, wantRadialError: 0.10012555},
		{name: "iteration 2", offset: 0.00501256, wantNextOffset: 0.00001256, wantRadialError: 0.00501258},
	}

	const radius = 1.0 // workpiece diameter 2

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nextOffset, radialError := osborneStep(tt.offset, radius)
			if diff := math.Abs(nextOffset - tt.wantNextOffset); diff > 5e-8 {
				t.Errorf("nextOffset = %v, want %v", nextOffset, tt.wantNextOffset)
			}
			if diff := math.Abs(radialError - tt.wantRadialError); diff > 5e-8 {
				t.Errorf("radialError = %v, want %v", radialError, tt.wantRadialError)
			}
		})
	}
}

func TestOsborneStep_ZeroOffset_StaysAtCenter(t *testing.T) {
	// Zero initial offset means the stock is already centered: an
	// identity independent of this code, true for any radius.
	nextOffset, radialError := osborneStep(0, 1.0)
	if nextOffset != 0 || radialError != 0 {
		t.Errorf("osborneStep(0, 1.0) = (%v, %v), want (0, 0)", nextOffset, radialError)
	}
}

func TestOsborneStep_RepeatedIterations_ConvergeTowardZero(t *testing.T) {
	// Regardless of how large the initial offset is (short of being
	// larger than the radius itself), repeated iterations converge
	// toward the true center, matching OSBORNE.TXT's central claim.
	offset, radius := 0.5, 1.0
	for i := 0; i < 6; i++ {
		offset, _ = osborneStep(offset, radius)
	}
	if offset > 1e-6 {
		t.Errorf("offset after 6 iterations = %v, want it to have converged close to zero", offset)
	}
}
