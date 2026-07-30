package main

import (
	"math"
	"testing"
)

func TestSmallerCylinderDiameter(t *testing.T) {
	tests := []struct {
		name           string
		largerDiameter float64
		angleDegrees   float64
		want           float64
	}{
		// (1-sin(theta))/(1+sin(theta)) equals tan^2(45-theta/2), a
		// standard trigonometric identity independent of this code's
		// implementation, used here to check both the documented
		// default input (angle 1.5, so theta=0.75) and a round
		// number (angle 90, so theta=45).
		{
			name:           "documented default input, checked against the half-angle identity",
			largerDiameter: 0.75, angleDegrees: 1.5,
			want: 0.75 * math.Pow(math.Tan((45-0.75/2)*math.Pi/180), 2),
		},
		{
			name:           "90 degree angle, checked against the half-angle identity",
			largerDiameter: 1, angleDegrees: 90,
			want: math.Pow(math.Tan(22.5*math.Pi/180), 2),
		},
		// A zero angle means the two cylinders are the same size.
		{name: "zero angle keeps both cylinders equal", largerDiameter: 0.5, angleDegrees: 0, want: 0.5},
		// A 180 degree angle collapses the smaller cylinder to
		// nothing: sin(90)=1 makes the numerator zero.
		{name: "180 degree angle collapses the smaller cylinder to zero", largerDiameter: 0.5, angleDegrees: 180, want: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := smallerCylinderDiameter(tt.largerDiameter, tt.angleDegrees)
			if diff := math.Abs(got - tt.want); diff > 1e-9 {
				t.Errorf("smallerCylinderDiameter(%v, %v) = %v, want %v", tt.largerDiameter, tt.angleDegrees, got, tt.want)
			}
		})
	}
}
