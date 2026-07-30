package main

import (
	"math"
	"testing"
)

func TestRotaryDivisions_DocumentedDefaultInput(t *testing.T) {
	divisions := rotaryDivisions(13)
	if len(divisions) != 14 {
		t.Fatalf("len(divisions) = %d, want 14", len(divisions))
	}

	// Division 1 of 13, computed independently: 360/13 degrees.
	want := 360.0 / 13.0
	if math.Abs(divisions[1].DecimalDeg-want) > 1e-9 {
		t.Errorf("divisions[1].DecimalDeg = %v, want %v", divisions[1].DecimalDeg, want)
	}
	if divisions[1].Degrees != 27 || divisions[1].Minutes != 41 || divisions[1].Seconds != 32 {
		t.Errorf("divisions[1] DMS = %d %d %d, want 27 41 32", divisions[1].Degrees, divisions[1].Minutes, divisions[1].Seconds)
	}
}

func TestRotaryDivisions_LastDivisionWrapsToZeroDegrees(t *testing.T) {
	// The final division is a full 360 degree turn, which must wrap
	// back to 0 rather than reporting a 360 degree reading.
	divisions := rotaryDivisions(6)
	last := divisions[len(divisions)-1]
	if math.Abs(last.DecimalDeg-360) > 1e-9 {
		t.Errorf("last.DecimalDeg = %v, want 360", last.DecimalDeg)
	}
	if last.Degrees != 0 || last.Minutes != 0 || last.Seconds != 0 {
		t.Errorf("last division DMS = %d %d %d, want 0 0 0 (wrapped)", last.Degrees, last.Minutes, last.Seconds)
	}
}

func TestRotaryDivisions_SecondsCarryIntoMinute(t *testing.T) {
	// 25 equal divisions puts division 7 at exactly 100.8 degrees,
	// whose raw seconds value rounds to 60 before the carry: an
	// input that genuinely exercises the carry branch, found by a
	// brute-force search across division counts, unlike compound's
	// unreachable carry branch.
	divisions := rotaryDivisions(25)
	d := divisions[7]
	if math.Abs(d.DecimalDeg-100.8) > 1e-9 {
		t.Fatalf("divisions[7].DecimalDeg = %v, want 100.8", d.DecimalDeg)
	}
	if d.Degrees != 100 || d.Minutes != 48 || d.Seconds != 0 {
		t.Errorf("divisions[7] DMS = %d %d %d, want 100 48 0 (carried from 47 min 60 sec)", d.Degrees, d.Minutes, d.Seconds)
	}
}
