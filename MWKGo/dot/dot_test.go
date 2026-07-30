package main

import (
	"math"
	"testing"
)

func TestComputeThreadDepths_SixtyDegreeMatchesANSIConstant(t *testing.T) {
	// For a standard 60 degree thread form, H (sharp crest to sharp
	// root depth) is a well known ANSI B1.1 constant times the
	// pitch: 0.866025*pitch, independent of this program's own
	// tan-based derivation of the same value.
	pitch := 1.0 / 20.0
	got := computeThreadDepths(60.0, pitch)
	want := 0.866025 * pitch
	if math.Abs(got.SharpCrestSharpRoot-want) > 1e-6 {
		t.Errorf("SharpCrestSharpRoot = %v, want %v", got.SharpCrestSharpRoot, want)
	}
}

func TestComputeThreadDepths_FractionsOfH(t *testing.T) {
	// B, C, and D are documented as fixed fractions of H (5/8, 3/4,
	// 7/8); verifying those fractions directly guards against them
	// drifting apart from H during refactoring.
	got := computeThreadDepths(60.0, 0.05)
	checkClose(t, "FlatCrestFlatRoot", got.FlatCrestFlatRoot, 0.625*got.SharpCrestSharpRoot)
	checkClose(t, "SharpCrestFlatRoot", got.SharpCrestFlatRoot, 0.75*got.SharpCrestSharpRoot)
	checkClose(t, "FlatCrestSharpRoot", got.FlatCrestSharpRoot, 0.875*got.SharpCrestSharpRoot)
	checkClose(t, "DoubleSharpToSharp", got.DoubleSharpToSharp, 2*got.SharpCrestSharpRoot)
}

func TestThreadingDialHint(t *testing.T) {
	cases := []struct {
		tpi  float64
		want string
	}{
		{20, "use any line on threading dial"},
		{9, "use any numbered line on threading dial"},
		{11.5, "use any odd-numbered line on threading dial"},
		{20.5, "use any odd-numbered line on threading dial"},
	}
	for _, c := range cases {
		if got := threadingDialHint(c.tpi); got != c.want {
			t.Errorf("threadingDialHint(%v) = %q, want %q", c.tpi, got, c.want)
		}
	}
}

func checkClose(t *testing.T, name string, got, want float64) {
	t.Helper()
	if math.Abs(got-want) > 1e-9 {
		t.Errorf("%s = %v, want %v", name, got, want)
	}
}
