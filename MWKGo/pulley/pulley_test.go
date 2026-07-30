package main

import (
	"math"
	"testing"
)

// TestFindUnknownDiameter_MatchesBeltDATConicalExample uses
// PULLEY.C's own default prompt values, which are BELT.DAT's own
// "conical pulley" worked example (d1=1.4, sep=2.5, belt length
// 8.21): the second pulley's diameter should converge to
// approximately 0.603, matching BELT.DAT's own paired row
// "2.5,0,0.603,1".
func TestFindUnknownDiameter_MatchesBeltDATConicalExample(t *testing.T) {
	r, err := findUnknownDiameter(1.4, 2.5, 8.21, 0.001)
	if err != nil {
		t.Fatal(err)
	}
	if math.Abs(r.D2-0.603) > 0.01 {
		t.Errorf("D2 = %v, want close to 0.603", r.D2)
	}
}

func TestFindUnknownDiameter_ResultReproducesTargetLength(t *testing.T) {
	r, err := findUnknownDiameter(1.4, 2.5, 8.21, 0.0001)
	if err != nil {
		t.Fatal(err)
	}
	gotLength := 2*r.Span + r.Wrap1 + r.Wrap2
	if math.Abs(gotLength-8.21) > 0.0001 {
		t.Errorf("reconstructed belt length = %v, want within 0.0001 of 8.21", gotLength)
	}
}
