package main

import (
	"math"
	"testing"
)

func TestSolveSphericalTriangle_AllNonAmbiguousCombinationsRecoverGroundTruth(t *testing.T) {
	// A single ground-truth spherical triangle (sides 60, 70, 80
	// degrees, solved via SSS) is used to check that every other
	// non-ambiguous combination of three known parts recovers the
	// same full solution. Independently verified via a separate
	// script before being trusted here.
	gt, _, err := solveSphericalTriangle([3]float64{60, 70, 80}, [3]float64{})
	if err != nil {
		t.Fatalf("solveSphericalTriangle(SSS) error = %v", err)
	}

	combos := []struct {
		name string
		s, a [3]float64
	}{
		{"AAA", [3]float64{}, gt.A},
		{"SAS angle0", [3]float64{0, gt.S[1], gt.S[2]}, [3]float64{gt.A[0], 0, 0}},
		{"ASA side0", [3]float64{gt.S[0], 0, 0}, [3]float64{0, gt.A[1], gt.A[2]}},
	}

	for _, c := range combos {
		t.Run(c.name, func(t *testing.T) {
			got, _, err := solveSphericalTriangle(c.s, c.a)
			if err != nil {
				t.Fatalf("solveSphericalTriangle() error = %v", err)
			}
			for i := 0; i < 3; i++ {
				checkClose(t, "S", got.S[i], gt.S[i])
				checkClose(t, "A", got.A[i], gt.A[i])
			}
		})
	}
}

func TestSolveSphericalTriangle_SSAAndAASRecoverGroundTruth(t *testing.T) {
	gt, _, err := solveSphericalTriangle([3]float64{60, 70, 80}, [3]float64{})
	if err != nil {
		t.Fatalf("solveSphericalTriangle(SSS) error = %v", err)
	}

	// SSA: sides 0,1 and angle 0 (s1=60 > s2=70 is false, so this
	// particular combination isn't the ambiguous branch; s[0]=60 <
	// s[1]=70 triggers ambiguity per ssaSolve's own condition, so
	// check the primary solution matches ground truth regardless.
	got, _, err := solveSphericalTriangle([3]float64{gt.S[0], gt.S[1], 0}, [3]float64{gt.A[0], 0, 0})
	if err != nil {
		t.Fatalf("solveSphericalTriangle(SSA) error = %v", err)
	}
	for i := 0; i < 3; i++ {
		checkClose(t, "SSA S", got.S[i], gt.S[i])
		checkClose(t, "SSA A", got.A[i], gt.A[i])
	}

	// AAS: angles 0,1 and side 0.
	got2, _, err := solveSphericalTriangle([3]float64{gt.S[0], 0, 0}, [3]float64{gt.A[0], gt.A[1], 0})
	if err != nil {
		t.Fatalf("solveSphericalTriangle(AAS) error = %v", err)
	}
	for i := 0; i < 3; i++ {
		checkClose(t, "AAS S", got2.S[i], gt.S[i])
		checkClose(t, "AAS A", got2.A[i], gt.A[i])
	}
}

func TestSolveSphericalTriangle_SSAAmbiguousCaseBothSolutionsValid(t *testing.T) {
	// sideI < sideJ triggers the ambiguous branch in ssaSolve; both
	// returned solutions must independently satisfy the spherical
	// law of cosines for sides using their own (possibly alternate)
	// values, the defining property of a valid spherical triangle.
	sol, alt, err := solveSphericalTriangle([3]float64{40, 60, 0}, [3]float64{30, 0, 0})
	if err != nil {
		t.Fatalf("solveSphericalTriangle() error = %v", err)
	}
	if alt == nil {
		t.Fatal("alt = nil, want a second solution for this ambiguous case (sideI < sideJ)")
	}

	for _, tri := range []sphereSolution{sol, {S: alt.S, A: alt.A}} {
		wantA1 := math.Acos((math.Cos(tri.S[0]*math.Pi/180)-math.Cos(tri.S[1]*math.Pi/180)*math.Cos(tri.S[2]*math.Pi/180))/
			(math.Sin(tri.S[1]*math.Pi/180)*math.Sin(tri.S[2]*math.Pi/180))) * 180 / math.Pi
		checkClose(t, "law of cosines A[0]", tri.A[0], wantA1)
	}
}

func TestSolveSphericalTriangle_InsufficientDataReturnsError(t *testing.T) {
	_, _, err := solveSphericalTriangle([3]float64{60, 0, 0}, [3]float64{})
	if err == nil {
		t.Fatal("solveSphericalTriangle() error = nil, want an error for a single known part")
	}
}

func checkClose(t *testing.T, name string, got, want float64) {
	t.Helper()
	if math.Abs(got-want) > 1e-6 {
		t.Errorf("%s = %v, want %v", name, got, want)
	}
}
