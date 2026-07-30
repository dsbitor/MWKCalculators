package main

import (
	"math"
	"testing"
)

func TestTubeDiameter(t *testing.T) {
	tests := []struct {
		name           string
		parentDiameter float64
		offset         float64
		want           float64
	}{
		// Zero offset means no eccentricity is being cut at all, so
		// the required tube diameter must equal the parent stock
		// diameter exactly: a physical identity independent of the
		// specific formula used, true for any parentDiameter.
		{name: "zero offset needs no tube beyond the parent stock", parentDiameter: 1.0, offset: 0, want: 1.0},
		{name: "zero offset holds for a different parent diameter too", parentDiameter: 2.5, offset: 0, want: 2.5},
		// The documented default input, evaluated against the ported
		// formula directly (a regression check on the arithmetic,
		// not an independently derived value).
		{name: "documented default input", parentDiameter: 1.0, offset: 0.1, want: 2 * math.Sqrt(7*0.5*0.5-9*0.5*0.4+3*0.4*0.4)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tubeDiameter(tt.parentDiameter, tt.offset)
			if diff := math.Abs(got - tt.want); diff > 1e-9 {
				t.Errorf("tubeDiameter(%v, %v) = %v, want %v", tt.parentDiameter, tt.offset, got, tt.want)
			}
		})
	}
}

func TestTubeDiameter_OffsetEqualsRadius_InnerRadiusIsZero(t *testing.T) {
	// When the offset equals the stock radius, the inner turning
	// radius r is exactly zero: the quantity under the square root
	// reduces to 7*R*R, giving dtube = 2*sqrt(7)*R. This is always a
	// positive, finite value (3r^2-9Rr+7R^2 has a negative
	// discriminant in r for any R != 0, so it never goes negative,
	// however large the offset), so this edge case checks that the
	// boundary is handled the same as any other value rather than
	// treated as a special case that needs its own branch.
	radius := 0.5
	want := 2 * math.Sqrt(7) * radius
	got := tubeDiameter(2*radius, radius)
	if diff := math.Abs(got - want); diff > 1e-9 {
		t.Errorf("tubeDiameter(%v, %v) = %v, want %v", 2*radius, radius, got, want)
	}
}
