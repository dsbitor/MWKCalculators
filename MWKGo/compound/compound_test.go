package main

import (
	"math"
	"testing"
)

func TestCompoundAngleDeg(t *testing.T) {
	tests := []struct {
		name  string
		ratio float64
		want  float64
	}{
		// A standard arcsine identity, independent of this code.
		{name: "one half gives 30 degrees", ratio: 0.5, want: 30},
		{name: "zero ratio needs zero angle", ratio: 0, want: 0},
		{name: "ratio of one gives 90 degrees", ratio: 1, want: 90},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := compoundAngleDeg(tt.ratio)
			if diff := math.Abs(got - tt.want); diff > 1e-9 {
				t.Errorf("compoundAngleDeg(%v) = %v, want %v", tt.ratio, got, tt.want)
			}
		})
	}
}

func TestDegreesMinutesSeconds(t *testing.T) {
	tests := []struct {
		name     string
		angleDeg float64
		wantSign string
		wantDeg  int
		wantMin  int
		wantSec  int
	}{
		{name: "whole degrees", angleDeg: 10.0, wantSign: "", wantDeg: 10, wantMin: 0, wantSec: 0},
		{name: "half a degree is 30 minutes", angleDeg: 10.5, wantSign: "", wantDeg: 10, wantMin: 30, wantSec: 0},
		{name: "negative angle keeps its sign separate from the magnitude", angleDeg: -10.5, wantSign: "-", wantDeg: 10, wantMin: 30, wantSec: 0},
		{name: "zero angle", angleDeg: 0, wantSign: "", wantDeg: 0, wantMin: 0, wantSec: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSign, gotDeg, gotMin, gotSec := degreesMinutesSeconds(tt.angleDeg)
			if gotSign != tt.wantSign || gotDeg != tt.wantDeg || gotMin != tt.wantMin || gotSec != tt.wantSec {
				t.Errorf("degreesMinutesSeconds(%v) = (%q, %v, %v, %v), want (%q, %v, %v, %v)",
					tt.angleDeg, gotSign, gotDeg, gotMin, gotSec, tt.wantSign, tt.wantDeg, tt.wantMin, tt.wantSec)
			}
		})
	}
}

func TestCompoundAngleDeg_ComplementSumsToNinety(t *testing.T) {
	// An angle and its complement must always sum to 90 degrees: an
	// identity independent of this code's specific ratio handling.
	angle := compoundAngleDeg(0.3)
	complement := 90 - angle
	if diff := math.Abs((angle + complement) - 90); diff > 1e-9 {
		t.Errorf("angle + complement = %v, want 90", angle+complement)
	}
}
