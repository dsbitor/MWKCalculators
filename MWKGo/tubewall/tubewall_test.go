package main

import (
	"math"
	"testing"
)

func TestWallThickness_DocumentedWorkedExample(t *testing.T) {
	// From TUBEWALL.TXT: the author's own measurement of a piece of
	// copper pipe, all three inputs and the resulting thickness
	// given explicitly.
	got := wallThickness(0.249, 0.879, 0.0625)
	want := 0.0425
	if diff := math.Abs(got - want); diff > 0.0001 {
		t.Errorf("wallThickness(0.249, 0.879, 0.0625) = %v, want %v (per TUBEWALL.TXT)", got, want)
	}
}

func TestWallThickness_ZeroDiameterAnvilNeedsNoCorrection(t *testing.T) {
	// The whole reason this program exists is that a flat anvil
	// "bridges" the tube's curved inside wall instead of touching it
	// directly. A point-contact anvil (diameter zero) has no such
	// bridging error to correct for, so the formula must reduce to
	// the reading itself: t equals the micrometer measurement
	// exactly, independent of the tube's outside diameter. Verified
	// independently by working through the quadratic by hand, not
	// just by re-running this code with different inputs.
	tests := []struct {
		tubeOutsideDiameter, reading float64
	}{
		{tubeOutsideDiameter: 1.0, reading: 0.4},
		{tubeOutsideDiameter: 2.0, reading: 0.3},
	}

	for _, tt := range tests {
		got := wallThickness(0, tt.tubeOutsideDiameter, tt.reading)
		if diff := math.Abs(got - tt.reading); diff > 1e-9 {
			t.Errorf("wallThickness(0, %v, %v) = %v, want %v", tt.tubeOutsideDiameter, tt.reading, got, tt.reading)
		}
	}
}

func TestWallThickness_InsideDiameterIsConsistentWithThickness(t *testing.T) {
	thickness := wallThickness(0.249, 0.879, 0.0625)
	insideDiameter := 0.879 - 2*thickness
	if insideDiameter <= 0 || insideDiameter >= 0.879 {
		t.Errorf("insideDiameter = %v, want a value between 0 and the outside diameter 0.879", insideDiameter)
	}
}
