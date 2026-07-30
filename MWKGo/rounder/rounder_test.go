package main

import (
	"math"
	"testing"
)

func TestRounderTable_EndpointsMatchKnownAngles(t *testing.T) {
	// At theta=0 the mill center hasn't moved off the starting axis
	// at all (A=0, D=0); at theta=90 it has swept a full quarter
	// turn (B=0, C=0): both independently obvious from the geometry,
	// not just a re-run of the formula.
	r := 1.25
	rows := rounderTable(r, 90)
	if len(rows) != 2 {
		t.Fatalf("len(rows) = %d, want 2", len(rows))
	}
	first, last := rows[0], rows[1]

	checkClose(t, "first.A", first.A, 0)
	checkClose(t, "first.B", first.B, r)
	checkClose(t, "first.C", first.C, r)
	checkClose(t, "first.D", first.D, 0)

	checkClose(t, "last.A", last.A, r)
	checkClose(t, "last.B", last.B, 0)
	checkClose(t, "last.C", last.C, 0)
	checkClose(t, "last.D", last.D, r)
}

func TestRounderTable_ABLieOnCircleOfCombinedRadius(t *testing.T) {
	// (A,B) is, by construction, a point on a circle of radius
	// combinedRadius centered on the origin: an identity independent
	// of the sin/cos formula used to place it there.
	r := 1.25
	for _, row := range rounderTable(r, 10) {
		got := math.Hypot(row.A, row.B)
		checkClose(t, "hypot(A,B)", got, r)
	}
}

func TestRounderTable_CAndDAreComplementsOfAAndB(t *testing.T) {
	// C and D are defined as r-A and r-B respectively; verifying
	// that relationship directly guards against the four formulas
	// drifting apart from each other.
	r := 1.25
	for _, row := range rounderTable(r, 15) {
		checkClose(t, "C", row.C, r-row.A)
		checkClose(t, "D", row.D, r-row.B)
	}
}

func checkClose(t *testing.T, name string, got, want float64) {
	t.Helper()
	if math.Abs(got-want) > 1e-9 {
		t.Errorf("%s = %v, want %v", name, got, want)
	}
}
