package mwktrig

import (
	"math"
	"testing"
)

func TestClampedAsinDeg(t *testing.T) {
	tests := []struct {
		name string
		x    float64
		want float64
	}{
		{name: "zero", x: 0, want: 0},
		{name: "one half gives 30 degrees, a standard identity", x: 0.5, want: 30},
		{name: "exactly 1 gives 90 degrees", x: 1, want: 90},
		{name: "exactly -1 gives -90 degrees", x: -1, want: -90},
		{name: "slightly over 1 clamps to 90 instead of NaN", x: 1.0000000001, want: 90},
		{name: "slightly under -1 clamps to -90 instead of NaN", x: -1.0000000001, want: -90},
		{name: "far outside the domain still clamps by sign", x: 5, want: 90},
		{name: "far outside the domain still clamps by sign, negative", x: -5, want: -90},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ClampedAsinDeg(tt.x)
			if diff := math.Abs(got - tt.want); diff > 1e-9 {
				t.Errorf("ClampedAsinDeg(%v) = %v, want %v", tt.x, got, tt.want)
			}
		})
	}
}

func TestClampedAsin(t *testing.T) {
	tests := []struct {
		name string
		x    float64
		want float64
	}{
		{name: "zero", x: 0, want: 0},
		{name: "one half gives pi/6, a standard identity", x: 0.5, want: math.Pi / 6},
		{name: "exactly 1 gives pi/2", x: 1, want: math.Pi / 2},
		{name: "exactly -1 gives -pi/2", x: -1, want: -math.Pi / 2},
		{name: "slightly over 1 clamps to pi/2 instead of NaN", x: 1.0000000001, want: math.Pi / 2},
		{name: "far outside the domain still clamps by sign", x: 5, want: math.Pi / 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ClampedAsin(tt.x)
			if diff := math.Abs(got - tt.want); diff > 1e-9 {
				t.Errorf("ClampedAsin(%v) = %v, want %v", tt.x, got, tt.want)
			}
		})
	}
}

func TestClampedAcosDeg(t *testing.T) {
	tests := []struct {
		name string
		x    float64
		want float64
	}{
		{name: "one gives zero degrees", x: 1, want: 0},
		{name: "zero gives 90 degrees, a standard identity", x: 0, want: 90},
		{name: "negative one gives 180 degrees", x: -1, want: 180},
		{name: "slightly over 1 clamps to 0 instead of NaN", x: 1.0000000001, want: 0},
		{name: "slightly under -1 clamps to 180 instead of NaN", x: -1.0000000001, want: 180},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ClampedAcosDeg(tt.x)
			if diff := math.Abs(got - tt.want); diff > 1e-9 {
				t.Errorf("ClampedAcosDeg(%v) = %v, want %v", tt.x, got, tt.want)
			}
		})
	}
}
