package main

import (
	"math"
	"testing"
)

// descartesDiameter computes the same fourth-circle diameter using
// Descartes' Circle Theorem stated in its usual curvature form
// (curvature = 2/diameter here, not the more common 1/radius, to
// stay consistent with soddyDiameter's diameter-based inputs), as an
// independent check on soddyDiameter's diameter-based formula.
func descartesDiameter(d1, d2, d3 float64, outer bool) float64 {
	k1, k2, k3 := 2/d1, 2/d2, 2/d3
	root := 2 * math.Sqrt(k1*k2+k2*k3+k3*k1)
	k4 := k1 + k2 + k3 + root
	if outer {
		k4 = k1 + k2 + k3 - root
	}
	return 2 / math.Abs(k4)
}

func TestSoddyDiameter_MatchesDescartesCircleTheorem(t *testing.T) {
	tests := []struct {
		name       string
		d1, d2, d3 float64
	}{
		{name: "documented default input", d1: 0.245, d2: 0.249, d3: 0.253},
		{name: "three equal circles", d1: 1, d2: 1, d3: 1},
		{name: "three very different sizes", d1: 1, d2: 2, d3: 5},
	}

	for _, tt := range tests {
		for _, outer := range []bool{true, false} {
			t.Run(tt.name, func(t *testing.T) {
				got := soddyDiameter(tt.d1, tt.d2, tt.d3, outer)
				want := descartesDiameter(tt.d1, tt.d2, tt.d3, outer)
				if diff := math.Abs(got - want); diff > 1e-6 {
					t.Errorf("soddyDiameter(%v, %v, %v, %v) = %v, want %v (Descartes Circle Theorem)", tt.d1, tt.d2, tt.d3, outer, got, want)
				}
			})
		}
	}
}

func TestSoddyDiameter_ThreeEqualCircles(t *testing.T) {
	// Three equal circles of diameter 1 (radius 0.5), mutually
	// tangent, have a well known inner Soddy circle radius: r =
	// R*(2/sqrt(3)-1), an identity independent of this code. Doubled
	// to compare against soddyDiameter's diameter-based result.
	got := soddyDiameter(1, 1, 1, false)
	want := 2 * 0.5 * (2/math.Sqrt(3) - 1)
	if diff := math.Abs(got - want); diff > 1e-9 {
		t.Errorf("soddyDiameter(1, 1, 1, false) = %v, want %v", got, want)
	}
}
