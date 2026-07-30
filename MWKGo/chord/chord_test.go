package main

import (
	"math"
	"testing"
)

func TestChordLength(t *testing.T) {
	tests := []struct {
		name      string
		divisions float64
		diameter  float64
		want      float64
		wantIsNaN bool
	}{
		// A hexagon's side length equals the circumscribing circle's
		// radius (a well known identity, independent of this code).
		{name: "hexagon side equals the radius", divisions: 6, diameter: 1, want: 0.5},
		// A square inscribed in a unit-diameter circle has side
		// diameter*sin(45 degrees) = diameter*sqrt(2)/2.
		{name: "square side matches sqrt(2)/2", divisions: 4, diameter: 1, want: math.Sqrt2 / 2},
		// An equilateral triangle's side is diameter*sin(60 degrees).
		{name: "triangle side matches sqrt(3)/2", divisions: 3, diameter: 1, want: math.Sqrt(3) / 2},
		{name: "negative diameter scales the result", divisions: 6, diameter: -1, want: -0.5},
		{name: "negative divisions still produce a defined value", divisions: -4, diameter: 1, want: -math.Sqrt2 / 2},
		{name: "zero divisions produces NaN through ordinary float division", divisions: 0, diameter: 1, wantIsNaN: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := chordLength(tt.divisions, tt.diameter)
			if tt.wantIsNaN {
				if !math.IsNaN(got) {
					t.Errorf("chordLength(%v, %v) = %v, want NaN", tt.divisions, tt.diameter, got)
				}
				return
			}
			if diff := math.Abs(got - tt.want); diff > 1e-9 {
				t.Errorf("chordLength(%v, %v) = %v, want %v", tt.divisions, tt.diameter, got, tt.want)
			}
		})
	}
}
