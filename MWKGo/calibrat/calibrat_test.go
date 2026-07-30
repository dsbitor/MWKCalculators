package main

import (
	"math"
	"strings"
	"testing"
)

// TestLinearFit_WorkedExample reproduces CALIBRAT.TXT's own worked
// example verbatim: a thermometer reading 2 at a true 0, 51 at a
// true 50, and 102 at a true 100 calibrates to A=1, B=1.666667.
func TestLinearFit_WorkedExample(t *testing.T) {
	points := []calibrationPoint{
		{Truth: 0, Measured: 2},
		{Truth: 50, Measured: 51},
		{Truth: 100, Measured: 102},
	}
	a, b, err := linearFit(points)
	if err != nil {
		t.Fatalf("linearFit() error = %v", err)
	}
	if math.Abs(a-1) > 1e-9 {
		t.Errorf("A = %v, want 1", a)
	}
	if math.Abs(b-5.0/3.0) > 1e-9 {
		t.Errorf("B = %v, want 1.666667", b)
	}
}

// TestLinearFit_TableWorkedExample reproduces the table portion of
// CALIBRAT.TXT's worked example: with A=1, B=1.666667, the
// measured-to-truth and truth-to-measured tables at t=0,50,100 give
// the exact values shown in the original documentation.
func TestLinearFit_TableWorkedExample(t *testing.T) {
	points := []calibrationPoint{{0, 2}, {50, 51}, {100, 102}}
	a, b, err := linearFit(points)
	if err != nil {
		t.Fatalf("linearFit() error = %v", err)
	}

	cases := []struct {
		t                  float64
		wantMeasuredToTrue float64
		wantTrueToMeasured float64
	}{
		{0, -1.666667, 1.666667},
		{50, 48.333333, 51.666667},
		{100, 98.333333, 101.666667},
	}
	for _, c := range cases {
		gotInverse := (c.t - b) / a
		if math.Abs(gotInverse-c.wantMeasuredToTrue) > 1e-6 {
			t.Errorf("(%v-B)/A = %v, want %v", c.t, gotInverse, c.wantMeasuredToTrue)
		}
		gotForward := a*c.t + b
		if math.Abs(gotForward-c.wantTrueToMeasured) > 1e-6 {
			t.Errorf("A*%v+B = %v, want %v", c.t, gotForward, c.wantTrueToMeasured)
		}
	}
}

func TestLinearFit_TwoPointsExactSolution(t *testing.T) {
	// With exactly two points, the least-squares line must pass
	// through both exactly (the PS note in CALIBRAT.TXT makes this
	// claim explicitly).
	points := []calibrationPoint{{0, 2}, {100, 102}}
	a, b, err := linearFit(points)
	if err != nil {
		t.Fatalf("linearFit() error = %v", err)
	}
	for _, p := range points {
		got := a*p.Truth + b
		if math.Abs(got-p.Measured) > 1e-9 {
			t.Errorf("fit at truth=%v gives %v, want exactly %v", p.Truth, got, p.Measured)
		}
	}
}

func TestLinearFit_ZeroDeterminant(t *testing.T) {
	// Two coincident points (same truth value) make the fit
	// underdetermined: the determinant is exactly zero.
	points := []calibrationPoint{{5, 1}, {5, 2}}
	if _, _, err := linearFit(points); err != errZeroDeterminant {
		t.Errorf("linearFit() error = %v, want errZeroDeterminant", err)
	}
}

func TestLoadPoints_SortsByTruthValue(t *testing.T) {
	data := "STARTOFDATA\n100,102\n0,2\n50,51\nENDOFDATA\n"
	points, err := loadPoints(strings.NewReader(data))
	if err != nil {
		t.Fatalf("loadPoints() error = %v", err)
	}
	want := []float64{0, 50, 100}
	if len(points) != len(want) {
		t.Fatalf("loadPoints() = %v, want %d points", points, len(want))
	}
	for i, w := range want {
		if points[i].Truth != w {
			t.Errorf("loadPoints()[%d].Truth = %v, want %v", i, points[i].Truth, w)
		}
	}
}

func TestLoadPoints_MalformedLine(t *testing.T) {
	data := "STARTOFDATA\n100\nENDOFDATA\n"
	if _, err := loadPoints(strings.NewReader(data)); err == nil {
		t.Error("loadPoints() error = nil, want an error for a line with only one value")
	}
}
