package main

import (
	"math"
	"testing"
)

func TestTaperAngle(t *testing.T) {
	tests := []struct {
		name                                                 string
		smallDiameter, smallDepth, largeDiameter, largeDepth float64
		wantHalfAngle, wantRatio                             float64
	}{
		// Equal sphere diameters indicate a cylindrical hole, not a
		// tapered one: a physical identity independent of this
		// code, true for any depths that give a nonzero vertical
		// separation.
		{
			name:          "equal sphere diameters mean zero taper",
			smallDiameter: 0.5, smallDepth: 2.0, largeDiameter: 0.5, largeDepth: 1.0,
			wantHalfAngle: 0, wantRatio: 0,
		},
		{
			// The documented default input, checked against
			// math.Asin directly rather than a hand-typed literal.
			name:          "documented default input",
			smallDiameter: 0.5, smallDepth: 2.0, largeDiameter: 0.75, largeDepth: 1.0,
			wantRatio:     0.125 / 0.875,
			wantHalfAngle: math.Asin(0.125/0.875) * 180 / math.Pi,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotHalfAngle, gotRatio := taperAngle(tt.smallDiameter, tt.smallDepth, tt.largeDiameter, tt.largeDepth)
			if diff := math.Abs(gotRatio - tt.wantRatio); diff > 1e-9 {
				t.Errorf("ratio = %v, want %v", gotRatio, tt.wantRatio)
			}
			if diff := math.Abs(gotHalfAngle - tt.wantHalfAngle); diff > 1e-9 {
				t.Errorf("halfAngle = %v, want %v", gotHalfAngle, tt.wantHalfAngle)
			}
		})
	}
}

func TestTaperAngle_ZeroVerticalSeparation_ClampsToNinetyDegrees(t *testing.T) {
	// When the two spheres' tops land at the same height but their
	// diameters differ, the ratio's denominator is zero, sending
	// the raw ratio to +-Inf. mwktrig.ClampedAsinDeg clamps this to
	// +-90 degrees rather than propagating NaN, matching the
	// original ASND macro's domain-clamping behaviour.
	halfAngle, ratio := taperAngle(0.5, 1.0, 0.75, 0.875)
	if !math.IsInf(ratio, 1) {
		t.Fatalf("ratio = %v, want +Inf", ratio)
	}
	if halfAngle != 90 {
		t.Errorf("halfAngle = %v, want 90 (clamped)", halfAngle)
	}
}
