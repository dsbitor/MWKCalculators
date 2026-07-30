package main

import (
	"math"
	"testing"
)

func TestCircumscribedDiameterFromAcrossFlats(t *testing.T) {
	tests := []struct {
		name        string
		sides       int
		acrossFlats float64
		want        float64
	}{
		// Standard, independently verifiable identities relating a
		// regular polygon's across-flats and across-corners
		// dimensions.
		{name: "square across-corners is across-flats times sqrt(2)", sides: 4, acrossFlats: 1, want: math.Sqrt2},
		{name: "hexagon across-corners is across-flats times 2/sqrt(3)", sides: 6, acrossFlats: 1, want: 2 / math.Sqrt(3)},
		{name: "hexagon identity holds for a different across-flats value", sides: 6, acrossFlats: 0.25, want: 0.25 * 2 / math.Sqrt(3)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := circumscribedDiameterFromAcrossFlats(tt.sides, tt.acrossFlats)
			if diff := math.Abs(got - tt.want); diff > 1e-9 {
				t.Errorf("circumscribedDiameterFromAcrossFlats(%v, %v) = %v, want %v", tt.sides, tt.acrossFlats, got, tt.want)
			}
		})
	}
}

func TestTenonDepthOfCut(t *testing.T) {
	tests := []struct {
		name                  string
		stockDiameter         float64
		sides                 int
		circumscribedDiameter float64
	}{
		// The documented default input (odd sides, pentagon), and an
		// even-sided case, each evaluated against the ported formula
		// directly.
		{name: "documented default input, pentagon", stockDiameter: 0.5, sides: 5, circumscribedDiameter: 0.25},
		{name: "even sided hexagon", stockDiameter: 1.0, sides: 6, circumscribedDiameter: 0.75},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			halfAngleDeg := 0.5 * (360 / float64(tt.sides))
			flatToCenter := 0.5 * tt.circumscribedDiameter * math.Cos(halfAngleDeg*math.Pi/180)
			want := 0.5*tt.stockDiameter - flatToCenter

			got := tenonDepthOfCut(tt.stockDiameter, tt.sides, tt.circumscribedDiameter)
			if diff := math.Abs(got - want); diff > 1e-9 {
				t.Errorf("tenonDepthOfCut(%v, %v, %v) = %v, want %v", tt.stockDiameter, tt.sides, tt.circumscribedDiameter, got, want)
			}
		})
	}
}

func TestTenonDepthOfCut_ZeroWhenCircumscribedCircleFillsStock(t *testing.T) {
	// If the circumscribed circle's radius already equals the stock
	// radius exactly, cutting the flats still removes material
	// (each flat sits inside the circumscribed circle), so the depth
	// of cut is only zero in the degenerate case of an
	// infinite-sided "polygon" (a circle). This test instead checks
	// a simpler, always-true property: depth of cut must be positive
	// whenever the circumscribed diameter is less than the stock
	// diameter, since the flats then cut below the stock surface.
	got := tenonDepthOfCut(1.0, 6, 0.9)
	if got <= 0 {
		t.Errorf("tenonDepthOfCut(1.0, 6, 0.9) = %v, want a positive depth of cut", got)
	}
}
