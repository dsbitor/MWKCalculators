package main

import (
	"math"
	"testing"
)

// TestFindCenterDistance_MatchesBeltDATConicalExample uses PCD.C's
// own default prompt values, which are BELT.DAT's own "conical
// pulley" worked example (d1=1.4, d2=0.603, belt length 8.21): the
// separation should converge to approximately 2.5, matching
// BELT.DAT's own paired rows "0,0,1.4,1" and "2.5,0,0.603,1".
func TestFindCenterDistance_MatchesBeltDATConicalExample(t *testing.T) {
	r, err := findCenterDistance(1.4, 0.603, 8.21, 0.0001)
	if err != nil {
		t.Fatal(err)
	}
	if math.Abs(r.Sep-2.5) > 0.01 {
		t.Errorf("Sep = %v, want close to 2.5", r.Sep)
	}
}

func TestFindCenterDistance_ResultReproducesTargetLength(t *testing.T) {
	r, err := findCenterDistance(1.4, 0.603, 8.21, 0.0001)
	if err != nil {
		t.Fatal(err)
	}
	gotLength := 2*r.Span + r.Wrap1 + r.Wrap2
	if math.Abs(gotLength-8.21) > 0.0001 {
		t.Errorf("reconstructed belt length = %v, want within 0.0001 of 8.21", gotLength)
	}
}

// TestFindCenterDistance_AgreesWithOtherBeltDATRows checks every one
// of BELT.DAT's own "conical pulley" example rows (documented as a
// family of pulley-diameter pairs all sized to the same fixed 2.5
// separation at a belt length of 8.21), not just the first.
func TestFindCenterDistance_AgreesWithOtherBeltDATRows(t *testing.T) {
	cases := []struct{ d1, d2 float64 }{
		{1.5, 0.477},
		{1.6, 0.343},
		{1.7, 0.200},
	}
	for _, c := range cases {
		r, err := findCenterDistance(c.d1, c.d2, 8.21, 0.001)
		if err != nil {
			t.Fatalf("d1=%v d2=%v: %v", c.d1, c.d2, err)
		}
		if math.Abs(r.Sep-2.5) > 0.05 {
			t.Errorf("d1=%v d2=%v: Sep = %v, want close to 2.5 (BELT.DAT's own fixed separation)", c.d1, c.d2, r.Sep)
		}
	}
}
