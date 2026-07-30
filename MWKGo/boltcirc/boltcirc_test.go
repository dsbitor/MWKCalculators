package main

import (
	"math"
	"testing"
)

func TestBoltCircleLayout_FourHolesNoOffset(t *testing.T) {
	spacing, holes := boltCircleLayout(4, 1.0, 0.5, 0.0, 0.0, 0.0)

	wantSpacing := 2*1.0*math.Sin(45*math.Pi/180) - 0.5
	if math.Abs(spacing-wantSpacing) > 1e-9 {
		t.Errorf("spacing = %v, want %v", spacing, wantSpacing)
	}

	if len(holes) != 4 {
		t.Fatalf("len(holes) = %d, want 4", len(holes))
	}
	// The four holes of a square bolt circle at radius 1 are the
	// cardinal points of the unit circle: an independently
	// verifiable identity, not just a re-run of the formula.
	wantAngles := []float64{0, 90, 180, 270}
	wantX := []float64{1, 0, -1, 0}
	wantY := []float64{0, 1, 0, -1}
	for i, hole := range holes {
		if math.Abs(hole.AngleDeg-wantAngles[i]) > 1e-9 {
			t.Errorf("hole %d angle = %v, want %v", i, hole.AngleDeg, wantAngles[i])
		}
		if math.Abs(hole.X-wantX[i]) > 1e-9 {
			t.Errorf("hole %d X = %v, want %v", i, hole.X, wantX[i])
		}
		if math.Abs(hole.Y-wantY[i]) > 1e-9 {
			t.Errorf("hole %d Y = %v, want %v", i, hole.Y, wantY[i])
		}
	}
}

func TestBoltCircleLayout_CenterOffsetShiftsEveryHole(t *testing.T) {
	_, unshifted := boltCircleLayout(6, 2.0, 0.25, 0.0, 0.0, 0.0)
	_, shifted := boltCircleLayout(6, 2.0, 0.25, 0.0, 5.0, -3.0)

	for i := range unshifted {
		if math.Abs(shifted[i].X-(unshifted[i].X+5.0)) > 1e-9 {
			t.Errorf("hole %d X = %v, want %v", i, shifted[i].X, unshifted[i].X+5.0)
		}
		if math.Abs(shifted[i].Y-(unshifted[i].Y-3.0)) > 1e-9 {
			t.Errorf("hole %d Y = %v, want %v", i, shifted[i].Y, unshifted[i].Y-3.0)
		}
	}
}

func TestBoltCircleLayout_LargeHolesOverlap(t *testing.T) {
	// Holes much larger than the circle can accommodate must report
	// negative spacing, the signal main uses to print the overlap
	// warning.
	spacing, _ := boltCircleLayout(8, 1.0, 5.0, 0.0, 0.0, 0.0)
	if spacing >= 0 {
		t.Errorf("spacing = %v, want negative (holes should overlap)", spacing)
	}
}
