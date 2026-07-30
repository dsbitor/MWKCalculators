package main

import (
	"math"
	"testing"
)

func TestTwoPulleyWrap_WrapAnglesSumToFullCircle(t *testing.T) {
	// beta1 = pi - theta, beta2 = pi + theta always sum to 2*pi,
	// regardless of the pulley sizes or separation - the belt loop
	// closes exactly once around.
	beta1, beta2, _, _, _, _ := twoPulleyWrap(2, 6, 6)
	if math.Abs((beta1+beta2)-2*math.Pi) > 1e-9 {
		t.Errorf("beta1+beta2 = %v, want 2*pi", beta1+beta2)
	}
}

func TestTwoPulleyWrap_EqualDiametersGiveHalfWrapEach(t *testing.T) {
	// Two equal pulleys: eps=0, theta=0, so each wraps exactly half
	// the circle (pi radians = 180 degrees).
	beta1, beta2, _, _, _, _ := twoPulleyWrap(4, 4, 10)
	if math.Abs(beta1-math.Pi) > 1e-9 || math.Abs(beta2-math.Pi) > 1e-9 {
		t.Errorf("beta1,beta2 = %v,%v, want both pi (equal pulleys wrap half each)", beta1, beta2)
	}
}

func TestTwoPulleyWrap_TotalLengthMatchesWrapsPlusSpans(t *testing.T) {
	_, _, wrap1, wrap2, span, total := twoPulleyWrap(2, 6, 6)
	want := 2*span + wrap1 + wrap2
	if math.Abs(total-want) > 1e-9 {
		t.Errorf("total = %v, want %v (2*span + wrap1 + wrap2)", total, want)
	}
}

func TestTwoPulleyWrap_MatchesBeltDATConicalExample(t *testing.T) {
	// BELT.DAT's own "conical pulley" example rows (0,0,1.4,1) and
	// (2.5,0,0.603,1) are documented as one of a family of pulley
	// pairs sized to produce the same fixed belt length at a 2.5
	// separation. Confirm that pairing's own belt length is self-
	// consistent (a positive, physically sensible value) as a sanity
	// check independent of the PULLEY/PCD search tests, which verify
	// the actual documented target length of 8.21.
	_, _, _, _, _, total := twoPulleyWrap(1.4, 0.603, 2.5)
	if total <= 0 {
		t.Fatalf("total belt length = %v, want positive", total)
	}
	if math.Abs(total-8.21) > 0.05 {
		t.Errorf("total belt length = %v, want close to 8.21 (BELT.DAT's own documented target)", total)
	}
}
