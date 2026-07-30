package main

import (
	"math"
	"os"
	"testing"
)

func loadExample(t *testing.T) []pulley {
	t.Helper()
	f, err := os.Open("testdata/example.dat")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	pulleys, err := loadPulleys(f)
	if err != nil {
		t.Fatal(err)
	}
	return pulleys
}

// TestComputeBeltLayout_SixthPulleyIsSlack reproduces BELT.DAT's own
// shipped example, whose comment explicitly says its sixth pulley
// "should be slack" — confirming the negative-wrap-length detection
// fires exactly where the original data's own author expected it to.
func TestComputeBeltLayout_SixthPulleyIsSlack(t *testing.T) {
	pulleys := loadExample(t)
	if len(pulleys) != 6 {
		t.Fatalf("got %d pulleys, want 6", len(pulleys))
	}

	results, total, slack, err := computeBeltLayout(pulleys)
	if err != nil {
		t.Fatal(err)
	}
	if !slack {
		t.Error("expected the slack flag to be set")
	}
	if results[5].WrapLength >= 0 {
		t.Errorf("pulley 6 wrap length = %v, want negative (slack)", results[5].WrapLength)
	}
	for i := 0; i < 5; i++ {
		if results[i].WrapLength < 0 {
			t.Errorf("pulley %d wrap length = %v, want non-negative", i+1, results[i].WrapLength)
		}
	}
	if total <= 0 {
		t.Errorf("total belt length = %v, want positive", total)
	}
}

// TestComputeBeltLayout_TotalLengthMatchesSumOfParts confirms the
// reported total is exactly the sum of every pulley's own wrap length
// and span, not some independently computed figure.
func TestComputeBeltLayout_TotalLengthMatchesSumOfParts(t *testing.T) {
	pulleys := loadExample(t)
	results, total, _, err := computeBeltLayout(pulleys)
	if err != nil {
		t.Fatal(err)
	}
	sum := 0.0
	for _, r := range results {
		sum += r.WrapLength + r.SpanToNext
	}
	if math.Abs(sum-total) > 1e-9 {
		t.Errorf("sum of parts = %v, want total %v", sum, total)
	}
}

func TestComputeBeltLayout_OverlappingPulleysError(t *testing.T) {
	pulleys := []pulley{
		{X: 0, Y: 0, Diam: 2, Inside: true},
		{X: 0.5, Y: 0, Diam: 2, Inside: true}, // too close together for their diameters
	}
	if _, _, _, err := computeBeltLayout(pulleys); err == nil {
		t.Error("expected an overlap error for two pulleys closer than their combined radii")
	}
}

func TestComputeBeltLayout_TwoExternalPulleysMatchesClosedForm(t *testing.T) {
	// A simple two-pulley external-tangent belt (both "inside", i.e.
	// not crossed) should match the closed-form two-pulley formula
	// QBELT/PULLEY/PCD all share (see MWKGo/qbelt/qbelt.go), cross-
	// checking BELT's general N-pulley algorithm against the simpler
	// two-pulley special case.
	d1, d2, sep := 2.0, 6.0, 6.0
	pulleys := []pulley{
		{X: 0, Y: 0, Diam: d1, Inside: true},
		{X: sep, Y: 0, Diam: d2, Inside: true},
	}
	_, total, slack, err := computeBeltLayout(pulleys)
	if err != nil {
		t.Fatal(err)
	}
	if slack {
		t.Error("did not expect a slack pulley in a simple two-pulley layout")
	}

	r1, r2 := 0.5*d1, 0.5*d2
	eps := r2 - r1
	theta := 2 * math.Asin(eps/sep)
	wrap1 := r1 * (math.Pi - theta)
	wrap2 := r2 * (math.Pi + theta)
	span := math.Sqrt(sep*sep - eps*eps)
	wantTotal := 2*span + wrap1 + wrap2

	if math.Abs(total-wantTotal) > 1e-6 {
		t.Errorf("total belt length = %v, want %v (matching the closed-form two-pulley formula)", total, wantTotal)
	}
}
