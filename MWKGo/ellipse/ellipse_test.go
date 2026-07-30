package main

import (
	"math"
	"testing"
)

func TestComputeEllipseProperties_DocumentedDefaultInput(t *testing.T) {
	// Cross-checked independently via numerical integration of the
	// standard elliptic-integral-of-the-second-kind definition
	// before trusting these as expected values.
	got := computeEllipseProperties(10, 3)

	checkClose(t, "ExactPerimeter", got.ExactPerimeter, 43.85910015176584)
	checkClose(t, "RMSPerimeter", got.RMSPerimeter, 46.38505965758242)
	checkClose(t, "Ramanujan1Perimeter", got.Ramanujan1Perimeter, 43.85673381454986)
	checkClose(t, "Ramanujan2Perimeter", got.Ramanujan2Perimeter, 43.859097438351625)
}

func TestComputeEllipseProperties_CircleIsExactForEveryFormula(t *testing.T) {
	// A circle is a degenerate ellipse (a == b) with a well known
	// exact perimeter, 2*pi*r: the elliptic integral and all three
	// algebraic approximations should agree with it and each other,
	// since every approximation was designed to be exact in this
	// limiting case.
	r := 5.0
	got := computeEllipseProperties(r, r)
	want := 2 * math.Pi * r

	checkClose(t, "Eccentricity", got.Eccentricity, 0)
	checkClose(t, "Area", got.Area, math.Pi*r*r)
	checkClose(t, "ExactPerimeter", got.ExactPerimeter, want)
	checkClose(t, "RMSPerimeter", got.RMSPerimeter, want)
	checkClose(t, "Ramanujan1Perimeter", got.Ramanujan1Perimeter, want)
	checkClose(t, "Ramanujan2Perimeter", got.Ramanujan2Perimeter, want)
}

func TestComputeEllipseProperties_HighlyElongatedHasLargeRMSError(t *testing.T) {
	// The RMS approximation is documented as having up to 5% error;
	// a highly elongated ellipse should approach that regime, unlike
	// the much more accurate Ramanujan approximations.
	got := computeEllipseProperties(10, 0.5)
	if math.Abs(got.RMSErrorPct) < 1 {
		t.Errorf("RMSErrorPct = %v, want a large error for a highly elongated ellipse", got.RMSErrorPct)
	}
	if math.Abs(got.Ramanujan2ErrorPct) > 0.1 {
		t.Errorf("Ramanujan2ErrorPct = %v, want a small error even for a highly elongated ellipse", got.Ramanujan2ErrorPct)
	}
}

func checkClose(t *testing.T, name string, got, want float64) {
	t.Helper()
	if math.Abs(got-want) > 1e-9 {
		t.Errorf("%s = %v, want %v", name, got, want)
	}
}
