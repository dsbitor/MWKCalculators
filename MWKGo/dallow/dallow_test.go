package main

import (
	"math"
	"testing"
)

func TestTipAllowance(t *testing.T) {
	tests := []struct {
		name                 string
		includedAngleDegrees float64
		diameter             float64
		want                 float64
	}{
		// At a 90 degree included angle the tip is a perfect right
		//-angle cone, whose height above the full diameter equals
		// exactly the drill radius: a geometric identity independent
		// of this code.
		{name: "90 degree tip allowance equals the radius", includedAngleDegrees: 90, diameter: 1, want: 0.5},
		{name: "90 degree identity holds for a different diameter too", includedAngleDegrees: 90, diameter: 2, want: 1.0},
		// The documented default angle (118 degrees, the standard
		// twist drill point angle) and diameter, checked against
		// tan computed directly rather than a hand-typed literal.
		{
			name:                 "documented default input, 118 degree standard twist drill point",
			includedAngleDegrees: 118, diameter: 0.5,
			want: 0.5 * 0.5 / math.Tan(59*math.Pi/180),
		},
		// A zero diameter needs no allowance at all, regardless of
		// angle.
		{name: "zero diameter needs no allowance", includedAngleDegrees: 118, diameter: 0, want: 0},
		// As the included angle widens toward a flat-bottomed hole
		// (180 degrees), the allowance shrinks toward zero, since
		// the tip's cone flattens out.
		{
			name:                 "a wide angle close to flat gives a small allowance",
			includedAngleDegrees: 170, diameter: 1,
			want: 0.5 / math.Tan(85*math.Pi/180),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tipAllowance(tt.includedAngleDegrees, tt.diameter)
			if diff := math.Abs(got - tt.want); diff > 1e-9 {
				t.Errorf("tipAllowance(%v, %v) = %v, want %v", tt.includedAngleDegrees, tt.diameter, got, tt.want)
			}
		})
	}
}
