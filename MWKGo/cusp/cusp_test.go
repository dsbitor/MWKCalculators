package main

import (
	"math"
	"testing"
)

func TestPassSpacing(t *testing.T) {
	tests := []struct {
		name         string
		millDiameter float64
		cuspHeight   float64
		want         float64
	}{
		// A cusp height equal to the mill's radius means the pass
		// removes exactly a semicircle, whose chord is the full
		// diameter: a geometric identity independent of this code.
		{name: "cusp height equal to the radius spans the full diameter", millDiameter: 2, cuspHeight: 1, want: 2},
		{name: "the semicircle identity holds for a different mill size too", millDiameter: 1, cuspHeight: 0.5, want: 1},
		// Zero cusp height needs zero spacing: no gap between passes
		// leaves no cusp at all.
		{name: "zero cusp height needs zero spacing", millDiameter: 0.25, cuspHeight: 0, want: 0},
		// The documented default input, evaluated against the ported
		// formula directly.
		{
			name: "documented default input", millDiameter: 0.25, cuspHeight: 0.001,
			want: 2 * math.Sqrt(2*0.125*0.001-0.001*0.001),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := passSpacing(tt.millDiameter, tt.cuspHeight)
			if diff := math.Abs(got - tt.want); diff > 1e-9 {
				t.Errorf("passSpacing(%v, %v) = %v, want %v", tt.millDiameter, tt.cuspHeight, got, tt.want)
			}
		})
	}
}

func TestPassSpacing_CuspHeightExceedsRadius_ProducesNonFiniteResult(t *testing.T) {
	// A cusp height greater than the mill's radius is not physically
	// achievable with this ball profile. The original program does
	// not guard against it either; the contract here is that the
	// failure is immediately visible as NaN, not a silently wrong
	// number.
	got := passSpacing(1, 10)
	if !math.IsNaN(got) {
		t.Errorf("passSpacing(1, 10) = %v, want NaN for a cusp height far exceeding the mill radius", got)
	}
}
