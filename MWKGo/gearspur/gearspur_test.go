package main

import (
	"math"
	"testing"
)

func TestSolveSpurGear_AllSixSolvableCombinations(t *testing.T) {
	// A single self-consistent gear: 45 teeth, DP 20, giving OD 2.35
	// and PD 2.25. Solving from any two of the four related
	// quantities must recover the full, consistent set.
	const wantTeeth, wantDP, wantOD, wantPD = 45.0, 20.0, 2.35, 2.25

	combos := []struct {
		name              string
		od, teeth, dp, pd float64
	}{
		{"teeth+od", wantOD, wantTeeth, 0, 0},
		{"dp+teeth", 0, wantTeeth, wantDP, 0},
		{"dp+od", wantOD, 0, wantDP, 0},
		{"pd+teeth", 0, wantTeeth, 0, wantPD},
		{"pd+od", wantOD, 0, 0, wantPD},
		{"pd+dp", 0, 0, wantDP, wantPD},
	}

	for _, c := range combos {
		t.Run(c.name, func(t *testing.T) {
			got, err := solveSpurGear(c.od, c.teeth, c.dp, c.pd)
			if err != nil {
				t.Fatalf("solveSpurGear() error = %v", err)
			}
			checkClose(t, "Teeth", got.Teeth, wantTeeth)
			checkClose(t, "DiametralPitch", got.DiametralPitch, wantDP)
			checkClose(t, "OutsideDiameter", got.OutsideDiameter, wantOD)
			checkClose(t, "PitchDiameter", got.PitchDiameter, wantPD)
		})
	}
}

func TestSolveSpurGear_WholeDepthTierBoundary(t *testing.T) {
	fine, err := solveSpurGear(0, 40, 24, 0) // DP > 20
	if err != nil {
		t.Fatalf("solveSpurGear() error = %v", err)
	}
	checkClose(t, "fine.WholeDepth", fine.WholeDepth, 2.157/24)

	coarse, err := solveSpurGear(0, 40, 10, 0) // DP <= 20
	if err != nil {
		t.Fatalf("solveSpurGear() error = %v", err)
	}
	checkClose(t, "coarse.WholeDepth", coarse.WholeDepth, 2.2/10+0.002)
}

func TestSolveSpurGear_InsufficientDataReturnsError(t *testing.T) {
	_, err := solveSpurGear(0, 45, 0, 0)
	if err == nil {
		t.Fatal("solveSpurGear() error = nil, want an error for a single known value")
	}
}

func TestBrownAndSharpeCutter(t *testing.T) {
	cases := []struct {
		teeth float64
		want  int
	}{
		{135, 1}, {200, 1},
		{55, 2}, {134, 2},
		{12, 8}, {13, 8},
		{11, 0},
	}
	for _, c := range cases {
		if got := brownAndSharpeCutter(c.teeth); got != c.want {
			t.Errorf("brownAndSharpeCutter(%v) = %d, want %d", c.teeth, got, c.want)
		}
	}
}

func checkClose(t *testing.T, name string, got, want float64) {
	t.Helper()
	if math.Abs(got-want) > 1e-6 {
		t.Errorf("%s = %v, want %v", name, got, want)
	}
}
