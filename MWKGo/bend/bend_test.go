package main

import (
	"math"
	"testing"
)

func TestNeutralAxisOffset(t *testing.T) {
	tests := []struct {
		name      string
		thickness float64
		radius    float64
		want      float64
	}{
		{name: "tight bend, radius under 2x thickness", thickness: 1, radius: 1.9, want: 1.0 / 3},
		{name: "gentle bend, radius over 4x thickness", thickness: 1, radius: 4.1, want: 0.5},
		{name: "middle range uses the 0.4 default", thickness: 1, radius: 3, want: 0.4},
		{name: "boundary at exactly 2x thickness stays in the middle range", thickness: 1, radius: 2, want: 0.4},
		{name: "boundary at exactly 4x thickness stays in the middle range", thickness: 1, radius: 4, want: 0.4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := neutralAxisOffset(tt.thickness, tt.radius)
			if diff := math.Abs(got - tt.want); diff > 1e-9 {
				t.Errorf("neutralAxisOffset(%v, %v) = %v, want %v", tt.thickness, tt.radius, got, tt.want)
			}
		})
	}
}

func TestBendLengths(t *testing.T) {
	tests := []struct {
		name                                      string
		thickness, radius, angleDegrees           float64
		wantExterior, wantInterior, wantAllowance float64
	}{
		{
			// A 90 degree bend with round numbers, in the
			// gentle-bend range (radius 10 > 4x thickness 1), so the
			// offset is a clean 0.5. This checks the formula against
			// hand-computed values using pi/2 directly.
			name:      "90 degree gentle bend with round numbers",
			thickness: 1, radius: 10, angleDegrees: 90,
			wantExterior: math.Pi / 2 * 11, wantInterior: math.Pi / 2 * 10, wantAllowance: math.Pi / 2 * 10.5,
		},
		{
			// A zero angle bend consumes no material at all and has
			// no length on either surface: an identity independent
			// of the specific thickness, radius, or offset tier.
			name:      "zero angle bend has no length at all",
			thickness: 0.125, radius: 3, angleDegrees: 0,
			wantExterior: 0, wantInterior: 0, wantAllowance: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exterior, interior, allowance := bendLengths(tt.thickness, tt.radius, tt.angleDegrees)
			if diff := math.Abs(exterior - tt.wantExterior); diff > 1e-9 {
				t.Errorf("exterior = %v, want %v", exterior, tt.wantExterior)
			}
			if diff := math.Abs(interior - tt.wantInterior); diff > 1e-9 {
				t.Errorf("interior = %v, want %v", interior, tt.wantInterior)
			}
			if diff := math.Abs(allowance - tt.wantAllowance); diff > 1e-9 {
				t.Errorf("allowance = %v, want %v", allowance, tt.wantAllowance)
			}
		})
	}
}
