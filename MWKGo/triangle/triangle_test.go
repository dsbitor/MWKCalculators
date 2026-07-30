package main

import (
	"math"
	"testing"
)

func TestSolveTriangle_SSS_345RightTriangle(t *testing.T) {
	// The 3-4-5 right triangle is a well known exact case: its
	// largest angle is exactly 90 degrees.
	sol, alt, err := solveTriangle([3]float64{3, 4, 5}, [3]float64{})
	if err != nil {
		t.Fatalf("solveTriangle() error = %v", err)
	}
	if alt != nil {
		t.Error("alt != nil, want nil for an SSS solution (never ambiguous)")
	}
	checkClose(t, "A[2]", sol.A[2], 90)
	checkClose(t, "A[0]+A[1]+A[2]", sol.A[0]+sol.A[1]+sol.A[2], 180)
}

func TestSolveTriangle_AllCombinationsRecoverGroundTruth(t *testing.T) {
	// A single ground-truth triangle (the 3-4-5 right triangle,
	// solved via SSS) is used to check that every other
	// non-ambiguous combination of three known parts recovers the
	// same full solution.
	gt, _, err := solveTriangle([3]float64{3, 4, 5}, [3]float64{})
	if err != nil {
		t.Fatalf("solveTriangle(SSS) error = %v", err)
	}

	combos := []struct {
		name string
		s, a [3]float64
	}{
		{"ASA side0", [3]float64{gt.S[0], 0, 0}, [3]float64{0, gt.A[1], gt.A[2]}},
		{"ASA side1", [3]float64{0, gt.S[1], 0}, [3]float64{gt.A[0], 0, gt.A[2]}},
		{"SAA side0+angle2", [3]float64{gt.S[0], 0, 0}, [3]float64{gt.A[0], 0, gt.A[2]}},
		{"SAS sides0,1", [3]float64{gt.S[0], gt.S[1], 0}, [3]float64{0, 0, gt.A[2]}},
	}

	for _, c := range combos {
		t.Run(c.name, func(t *testing.T) {
			got, _, err := solveTriangle(c.s, c.a)
			if err != nil {
				t.Fatalf("solveTriangle() error = %v", err)
			}
			for i := 0; i < 3; i++ {
				checkClose(t, "S", got.S[i], gt.S[i])
				checkClose(t, "A", got.A[i], gt.A[i])
			}
		})
	}
}

func TestSolveTriangle_SSAAmbiguousCase(t *testing.T) {
	// s1=8, s2=10, a1=40deg is a classic ambiguous SSA case: since
	// s2 > s1, two distinct triangles satisfy the same given parts.
	sol, alt, err := solveTriangle([3]float64{8, 10, 0}, [3]float64{40, 0, 0})
	if err != nil {
		t.Fatalf("solveTriangle() error = %v", err)
	}
	if alt == nil {
		t.Fatal("alt = nil, want a second solution for this ambiguous case")
	}
	checkClose(t, "primary A[1]", sol.A[1], 53.46414901438847)
	checkClose(t, "alt A[1]", alt.A[1], 126.53585098561153)

	// Both solutions must independently satisfy the law of sines for
	// the given data (s1/sin(a1) = s2/sin(a2)): the defining property
	// of the ambiguous case, not just a re-run of the formula.
	for _, s := range []struct {
		name           string
		s1, a1, s2, a2 float64
	}{
		{"primary", sol.S[0], sol.A[0], sol.S[1], sol.A[1]},
		{"alt", alt.S[0], alt.A[0], alt.S[1], alt.A[1]},
	} {
		ratio1 := s.s1 / math.Sin(s.a1*math.Pi/180)
		ratio2 := s.s2 / math.Sin(s.a2*math.Pi/180)
		if math.Abs(ratio1-ratio2) > 1e-6 {
			t.Errorf("%s: law of sines violated: %v != %v", s.name, ratio1, ratio2)
		}
	}
}

func TestSolveTriangle_NoSideSpecifiedReturnsError(t *testing.T) {
	_, _, err := solveTriangle([3]float64{}, [3]float64{60, 60, 60})
	if err == nil {
		t.Fatal("solveTriangle() error = nil, want an error when no side is given")
	}
}

func TestComputeTriangleProperties_CircumradiusMatchesStandardFormula(t *testing.T) {
	// R = (a*b*c) / (4*Area) is the standard circumradius formula,
	// independent of this code's own law-of-sines-based computation
	// of the same value.
	sol, _, err := solveTriangle([3]float64{3, 4, 5}, [3]float64{})
	if err != nil {
		t.Fatalf("solveTriangle() error = %v", err)
	}
	p := computeTriangleProperties(sol)
	want := (sol.S[0] * sol.S[1] * sol.S[2]) / (4 * p.Area)
	checkClose(t, "CircumscribedRadius", p.CircumscribedRadius, want)
}

func checkClose(t *testing.T, name string, got, want float64) {
	t.Helper()
	if math.Abs(got-want) > 1e-6 {
		t.Errorf("%s = %v, want %v", name, got, want)
	}
}
