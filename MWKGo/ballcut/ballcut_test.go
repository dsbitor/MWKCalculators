package main

import (
	"math"
	"testing"
)

func TestBallCutScheduleByAngle_FirstWorkedExample(t *testing.T) {
	// BALLCUT.TXT's first worked example: a 1" sphere from 1" stock,
	// 5 degree angular steps, producing 19 rows (N=0 to N=18).
	steps := ballCutScheduleByAngle(1.0, 1.0, 5.0)
	if len(steps) != 19 {
		t.Fatalf("len(steps) = %d, want 19", len(steps))
	}

	want := []struct {
		xf, yf, wd float64
	}{
		{0.000, 0.500, 0.000},
		{0.002, 0.456, 0.087},
		{0.067, 0.250, 0.500},
		{0.500, 0.000, 1.000},
	}
	checkStep(t, steps[0], want[0].xf, want[0].yf, want[0].wd)
	checkStep(t, steps[1], want[1].xf, want[1].yf, want[1].wd)
	checkStep(t, steps[6], want[2].xf, want[2].yf, want[2].wd)
	checkStep(t, steps[18], want[3].xf, want[3].yf, want[3].wd)
}

func TestBallCutScheduleByAngle_SecondWorkedExample(t *testing.T) {
	// BALLCUT.TXT's second worked example: a 2" radius cut on 1"
	// stock. The worked example's last row (N=5) has a small
	// positive YF (0.077); the exact step at which YF crosses zero
	// sits at a floating-point boundary (theta=30 degrees, where the
	// original DOS runtime's trig rounding evidently differed from
	// Go's), so only the unambiguous early rows are checked here
	// rather than asserting a specific total row count.
	steps := ballCutScheduleByAngle(2.0, 1.0, 5.0)
	if len(steps) < 6 {
		t.Fatalf("len(steps) = %d, want at least 6", len(steps))
	}

	want := []struct {
		xf, yf, wd float64
	}{
		{0.000, 0.500, 0.000},
		{0.004, 0.413, 0.174},
		{0.015, 0.326, 0.347},
		{0.034, 0.241, 0.518},
		{0.060, 0.158, 0.684},
		{0.094, 0.077, 0.845},
	}
	for i, w := range want {
		checkStep(t, steps[i], w.xf, w.yf, w.wd)
	}
}

func TestBallCutScheduleByAxialStep_FirstRowMatchesAngleMode(t *testing.T) {
	// Both stepping modes start from the same physical position
	// (the pole of the sphere): axial position 0, depth of cut equal
	// to the stock radius, and zero work diameter.
	steps := ballCutScheduleByAxialStep(1.0, 1.0, 0.02)
	if len(steps) == 0 {
		t.Fatal("len(steps) = 0, want at least 1")
	}
	checkStep(t, steps[0], 0.0, 0.5, 0.0)
}

func checkStep(t *testing.T, got ballCutStep, wantXF, wantYF, wantWD float64) {
	t.Helper()
	if math.Abs(got.XF-wantXF) > 5e-4 {
		t.Errorf("step %d: XF = %v, want %v", got.N, got.XF, wantXF)
	}
	if math.Abs(got.YF-wantYF) > 5e-4 {
		t.Errorf("step %d: YF = %v, want %v", got.N, got.YF, wantYF)
	}
	if math.Abs(got.WD-wantWD) > 5e-4 {
		t.Errorf("step %d: WD = %v, want %v", got.N, got.WD, wantWD)
	}
}
