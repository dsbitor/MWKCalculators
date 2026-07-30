package main

import (
	"math"
	"os"
	"testing"
)

const eps = 1e-6

func almostEqual(a, b, tol float64) bool { return math.Abs(a-b) < tol }

// TestPolyFit_RecoversCurfitDATCubic reproduces CURFIT.DAT's own
// shipped example data exactly. Despite its comment claiming
// "y = 4 + 3*x + 2*x^2 + 1*x^2" (a transcription slip: the last term's
// exponent should read ^3, not ^2, since two terms both reading x^2
// would just be 3*x^2, not a distinct term), the data itself is a
// noise-free cubic. Solving for its finite differences confirms
// y = 4 + 3*x + 2*x^2 + 1*x^3 exactly (checked against x=18 and x=19's
// values below), so a degree-3 polynomial fit should recover
// A0=4, A1=3, A2=2, A3=1 almost exactly.
func TestPolyFit_RecoversCurfitDATCubic(t *testing.T) {
	f, err := os.Open("testdata/example.dat")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	points, err := loadPoints(f)
	if err != nil {
		t.Fatal(err)
	}
	if len(points) != 20 {
		t.Fatalf("got %d points, want 20", len(points))
	}

	coeffs, err := polyFit(points, 3)
	if err != nil {
		t.Fatal(err)
	}
	want := []float64{4, 3, 2, 1}
	for i, w := range want {
		if !almostEqual(coeffs[i], w, 1e-3) {
			t.Errorf("A%d = %v, want %v", i, coeffs[i], w)
		}
	}
}

// TestPolyFit_QuadraticExactRecovery checks polyFit against a small,
// hand-picked noise-free quadratic (y = 1 + 2x + 3x^2), independent of
// the CURFIT.DAT data.
func TestPolyFit_QuadraticExactRecovery(t *testing.T) {
	var points []point
	for x := -2.0; x <= 2.0; x++ {
		points = append(points, point{X: x, Y: 1 + 2*x + 3*x*x})
	}
	coeffs, err := polyFit(points, 2)
	if err != nil {
		t.Fatal(err)
	}
	want := []float64{1, 2, 3}
	for i, w := range want {
		if !almostEqual(coeffs[i], w, 1e-9) {
			t.Errorf("A%d = %v, want %v", i, coeffs[i], w)
		}
	}
}

// TestLogFit_ExactRecovery checks logFit against noise-free
// Y = 2 + 3*ln(X) data.
func TestLogFit_ExactRecovery(t *testing.T) {
	var points []point
	for _, x := range []float64{1, 2, 4, 8, 16} {
		points = append(points, point{X: x, Y: 2 + 3*math.Log(x)})
	}
	a, b := logFit(points)
	if !almostEqual(a, 2, 1e-9) || !almostEqual(b, 3, 1e-9) {
		t.Errorf("logFit() = (%v,%v), want (2,3)", a, b)
	}
}

// TestExpFit_ExactRecovery checks expFit against noise-free
// Y = 2*exp(0.5*X) data.
func TestExpFit_ExactRecovery(t *testing.T) {
	var points []point
	for _, x := range []float64{0, 1, 2, 3, 4} {
		points = append(points, point{X: x, Y: 2 * math.Exp(0.5*x)})
	}
	a, b := expFit(points)
	if !almostEqual(a, 2, 1e-9) || !almostEqual(b, 0.5, 1e-9) {
		t.Errorf("expFit() = (%v,%v), want (2,0.5)", a, b)
	}
}

// TestPowerFit_ExactRecovery checks powerFit against noise-free
// Y = 2*X^0.5 data (CURFIT.DAT's own commented power-fit example
// series, without its added noise variant).
func TestPowerFit_ExactRecovery(t *testing.T) {
	var points []point
	for _, x := range []float64{1, 2, 3, 4, 5, 6, 7, 8, 9} {
		points = append(points, point{X: x, Y: 2 * math.Pow(x, 0.5)})
	}
	a, b := powerFit(points)
	if !almostEqual(a, 2, 1e-9) || !almostEqual(b, 0.5, 1e-9) {
		t.Errorf("powerFit() = (%v,%v), want (2,0.5)", a, b)
	}
}

// TestCorrelationCoefficient_PerfectFitIsOne confirms a fit with zero
// error (sse=0) reports a correlation coefficient of exactly 1,
// regardless of sst, matching the formula f = 1 - sse*(n-1)/(sst*(n-nmat-1)).
func TestCorrelationCoefficient_PerfectFitIsOne(t *testing.T) {
	r := correlationCoefficient(10, 2, 100.0, 0.0)
	if !almostEqual(r, 1.0, eps) {
		t.Errorf("correlationCoefficient() = %v, want 1", r)
	}
}

// TestEvaluate_NewRunningMaxFlagsOnlyIncreasingPoints reproduces
// CURFIT.C's own "**" marking behavior exactly: the flag is set on
// every point whose |error%| exceeds the running maximum seen so far
// (which can be more than one point, since the running max only ever
// increases), not merely the single overall worst point.
func TestEvaluate_NewRunningMaxFlagsOnlyIncreasingPoints(t *testing.T) {
	// ycalc always returns 0, so error = -y, and |err%| = 100 for every
	// nonzero y. Using y values whose |err%| sequence is 100,100,100
	// (all equal) means only the FIRST point sets a new running max;
	// later ties do not exceed it.
	points := []point{{X: 0, Y: 1}, {X: 1, Y: 2}, {X: 2, Y: 3}}
	evaluated, _ := evaluate(points, func(float64) float64 { return 0 })
	if !evaluated[0].NewRunningMax {
		t.Error("first point should set the initial running max")
	}
	if evaluated[1].NewRunningMax || evaluated[2].NewRunningMax {
		t.Error("later points tying (not exceeding) the running max should not be flagged")
	}
}

func TestEvaluate_ZeroYGivesZeroErrorPercent(t *testing.T) {
	// CURFIT.C guards against division by zero: err% is 0 when the
	// data's own y value is 0, regardless of the computed error.
	points := []point{{X: 5, Y: 0}}
	evaluated, _ := evaluate(points, func(float64) float64 { return 100 })
	if evaluated[0].ErrorPercent != 0 {
		t.Errorf("ErrorPercent = %v, want 0 when Y is 0", evaluated[0].ErrorPercent)
	}
}
