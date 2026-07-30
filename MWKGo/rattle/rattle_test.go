package main

import (
	"math"
	"testing"
)

func TestBoreDiameter_DocumentedDefaultInput(t *testing.T) {
	diameter, diff := boreDiameter(4.0, 0.2)
	wantDiameter := 4.001252152167462
	wantDiff := wantDiameter - 4.0
	if math.Abs(diameter-wantDiameter) > 1e-9 {
		t.Errorf("diameter = %v, want %v", diameter, wantDiameter)
	}
	if math.Abs(diff-wantDiff) > 1e-9 {
		t.Errorf("diff = %v, want %v", diff, wantDiff)
	}
}

func TestBoreDiameter_ZeroRattleReturnsStickLength(t *testing.T) {
	// A stick that doesn't rattle at all is itself an exact measure
	// of the bore: the closest available caliper reading with no
	// rattle-derived correction needed.
	diameter, diff := boreDiameter(4.0, 0.0)
	if math.Abs(diameter-4.0) > 1e-9 {
		t.Errorf("diameter = %v, want 4.0", diameter)
	}
	if math.Abs(diff) > 1e-9 {
		t.Errorf("diff = %v, want 0", diff)
	}
}

func TestBoreDiameter_LargerRattleGivesLargerDiameter(t *testing.T) {
	// The whole point of the technique: the more a stick rattles in
	// the bore, the more its length underestimates the true
	// diameter, so the correction should grow with rattle distance.
	_, smallDiff := boreDiameter(4.0, 0.1)
	_, largeDiff := boreDiameter(4.0, 0.4)
	if !(largeDiff > smallDiff) {
		t.Errorf("diff for rattle=0.4 (%v) should exceed diff for rattle=0.1 (%v)", largeDiff, smallDiff)
	}
}
