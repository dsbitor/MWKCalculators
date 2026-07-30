package main

import (
	"math"
	"testing"
)

func TestLinkAngles(t *testing.T) {
	tests := []struct {
		name                                              string
		smallEndRadius, bigEndRadius, holeCenterDistance  float64
		wantHalfAngle, wantSmallEndAngle, wantBigEndAngle float64
	}{
		// Equal end radii mean a parallel-sided, untapered link: an
		// identity independent of this code, true for any center
		// distance.
		{
			name:           "equal end radii mean no taper",
			smallEndRadius: 0.1, bigEndRadius: 0.1, holeCenterDistance: 1.0,
			wantHalfAngle: 0, wantSmallEndAngle: 180, wantBigEndAngle: 180,
		},
		{
			// The documented default input, checked against
			// math.Asin directly rather than a hand-typed literal.
			name:           "documented default input",
			smallEndRadius: 1.0 / 16, bigEndRadius: 3.0 / 32, holeCenterDistance: 1.0,
			wantHalfAngle:     math.Asin(1.0/32) * 180 / math.Pi,
			wantSmallEndAngle: 180 - 2*math.Asin(1.0/32)*180/math.Pi,
			wantBigEndAngle:   180 + 2*math.Asin(1.0/32)*180/math.Pi,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			halfAngle, smallEndAngle, bigEndAngle := linkAngles(tt.smallEndRadius, tt.bigEndRadius, tt.holeCenterDistance)
			if diff := math.Abs(halfAngle - tt.wantHalfAngle); diff > 1e-9 {
				t.Errorf("halfAngle = %v, want %v", halfAngle, tt.wantHalfAngle)
			}
			if diff := math.Abs(smallEndAngle - tt.wantSmallEndAngle); diff > 1e-9 {
				t.Errorf("smallEndAngle = %v, want %v", smallEndAngle, tt.wantSmallEndAngle)
			}
			if diff := math.Abs(bigEndAngle - tt.wantBigEndAngle); diff > 1e-9 {
				t.Errorf("bigEndAngle = %v, want %v", bigEndAngle, tt.wantBigEndAngle)
			}
		})
	}
}

func TestShimHeight(t *testing.T) {
	tests := []struct {
		name                                                             string
		smallEndRadius, bigEndRadius, smallHoleDiameter, bigHoleDiameter float64
		want                                                             float64
	}{
		{name: "equal ends and holes need no shim at all", smallEndRadius: 0.1, bigEndRadius: 0.1, smallHoleDiameter: 0.05, bigHoleDiameter: 0.05, want: 0},
		{
			name:           "documented default input",
			smallEndRadius: 1.0 / 16, bigEndRadius: 3.0 / 32, smallHoleDiameter: 1.0 / 16, bigHoleDiameter: 3.0 / 16,
			want: 0.09375,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shimHeight(tt.smallEndRadius, tt.bigEndRadius, tt.smallHoleDiameter, tt.bigHoleDiameter)
			if diff := math.Abs(got - tt.want); diff > 1e-9 {
				t.Errorf("shimHeight(...) = %v, want %v", got, tt.want)
			}
		})
	}
}
