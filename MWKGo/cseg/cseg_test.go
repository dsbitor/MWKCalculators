package main

import (
	"math"
	"testing"
)

// A single ground-truth circular segment (radius 5, included angle
// 40 deg) used to check that every combination of two known
// quantities recovers the same full, consistent set.
var groundTruth = computeFromRadiusAngle(5.0, 40.0)

func TestSolveCircularSegment_AllClosedFormCombinations(t *testing.T) {
	gt := groundTruth
	combos := []struct {
		name                              string
		radius, angle, chord, height, arc float64
	}{
		{"radius+angle", gt.Radius, gt.AngleDeg, 0, 0, 0},
		{"radius+chord", gt.Radius, 0, gt.Chord, 0, 0},
		{"radius+height", gt.Radius, 0, 0, gt.Height, 0},
		{"radius+arc", gt.Radius, 0, 0, 0, gt.Arc},
		{"angle+chord", 0, gt.AngleDeg, gt.Chord, 0, 0},
		{"angle+height", 0, gt.AngleDeg, 0, gt.Height, 0},
		{"angle+arc", 0, gt.AngleDeg, 0, 0, gt.Arc},
		{"chord+height", 0, 0, gt.Chord, gt.Height, 0},
	}

	for _, c := range combos {
		t.Run(c.name, func(t *testing.T) {
			got, err := solveCircularSegment(c.radius, c.angle, c.chord, c.height, c.arc)
			if err != nil {
				t.Fatalf("solveCircularSegment() error = %v", err)
			}
			checkClose(t, "Radius", got.Radius, gt.Radius)
			checkClose(t, "AngleDeg", got.AngleDeg, gt.AngleDeg)
			checkClose(t, "Chord", got.Chord, gt.Chord)
			checkClose(t, "Height", got.Height, gt.Height)
			checkClose(t, "Arc", got.Arc, gt.Arc)
			checkClose(t, "Area", got.Area, gt.Area)
		})
	}
}

func TestSolveCircularSegment_SearchBasedCombinations(t *testing.T) {
	gt := groundTruth
	combos := []struct {
		name               string
		chord, height, arc float64
	}{
		{"chord+arc", gt.Chord, 0, gt.Arc},
		{"height+arc", 0, gt.Height, gt.Arc},
	}

	for _, c := range combos {
		t.Run(c.name, func(t *testing.T) {
			got, err := solveCircularSegment(0, 0, c.chord, c.height, c.arc)
			if err != nil {
				t.Fatalf("solveCircularSegment() error = %v", err)
			}
			// The search-based combinations only converge to a
			// finite (not floating-point-exact) precision, so a
			// looser tolerance is used here than for the
			// closed-form combinations above.
			if math.Abs(got.Radius-gt.Radius) > 1e-3 {
				t.Errorf("Radius = %v, want %v", got.Radius, gt.Radius)
			}
			if math.Abs(got.AngleDeg-gt.AngleDeg) > 1e-3 {
				t.Errorf("AngleDeg = %v, want %v", got.AngleDeg, gt.AngleDeg)
			}
		})
	}
}

func TestSolveCircularSegment_InsufficientDataReturnsError(t *testing.T) {
	_, err := solveCircularSegment(5.0, 0, 0, 0, 0)
	if err == nil {
		t.Fatal("solveCircularSegment() error = nil, want an error for a single known value")
	}
}

func checkClose(t *testing.T, name string, got, want float64) {
	t.Helper()
	if math.Abs(got-want) > 1e-6 {
		t.Errorf("%s = %v, want %v", name, got, want)
	}
}
