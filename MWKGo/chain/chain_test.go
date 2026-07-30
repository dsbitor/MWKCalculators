package main

import (
	"math"
	"testing"
)

func TestCenterDistance(t *testing.T) {
	tests := []struct {
		name                   string
		pitch, chainLength     float64
		largeTeeth, smallTeeth int
		want                   float64
	}{
		// Equal sprocket sizes make the tooth-difference term vanish,
		// reducing the formula to (pitch/2)*(length-teeth): an
		// identity independent of this code, since the general
		// formula's sqrt term collapses to the tooth-sum term itself
		// when the difference is zero (for length > teeth, keeping
		// the sum positive).
		{name: "equal sprockets reduce to a simple linear formula", pitch: 1, chainLength: 48, largeTeeth: 18, smallTeeth: 18, want: 0.5 * (48 - 18)},
		{
			// The documented default input, evaluated against the
			// ported formula directly.
			name: "documented default input", pitch: 1, chainLength: 48, largeTeeth: 18, smallTeeth: 9,
			want: (1.0 / 8) * ((2*48.0 - 18 - 9) + math.Sqrt((2*48.0-18-9)*(2*48.0-18-9)-0.810*(18-9)*(18-9))),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := centerDistance(tt.pitch, tt.chainLength, tt.largeTeeth, tt.smallTeeth)
			if diff := math.Abs(got - tt.want); diff > 1e-9 {
				t.Errorf("centerDistance(%v, %v, %v, %v) = %v, want %v", tt.pitch, tt.chainLength, tt.largeTeeth, tt.smallTeeth, got, tt.want)
			}
		})
	}
}

func TestCenterDistance_ScalesLinearlyWithPitch(t *testing.T) {
	// Pitch is a simple multiplier on the whole formula: an identity
	// independent of this code's specific tooth-count arithmetic.
	single := centerDistance(1, 48, 18, 9)
	doubled := centerDistance(2, 48, 18, 9)
	if diff := math.Abs(doubled - 2*single); diff > 1e-9 {
		t.Errorf("centerDistance(2, ...) = %v, want exactly double centerDistance(1, ...) = %v", doubled, 2*single)
	}
}
