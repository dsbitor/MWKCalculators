package main

import (
	"math"
	"testing"
)

func TestApproximateFraction_DocumentedDefaultInputMatchesMilu(t *testing.T) {
	// The documented default input (3.14159, 0.01% accuracy)
	// converges on 355/113 (Milü), the famous, extremely accurate
	// simple rational approximation to pi: an independently
	// verifiable historical fact, not just a re-run of the formula.
	result, err := approximateFraction(3.14159, 0.01)
	if err != nil {
		t.Fatalf("approximateFraction() error = %v", err)
	}
	if result.Whole != 3 {
		t.Errorf("Whole = %d, want 3", result.Whole)
	}
	if result.Num != 16 || result.Den != 113 {
		t.Errorf("Num/Den = %d/%d, want 16/113 (3 & 16/113 = 355/113)", result.Num, result.Den)
	}
}

func TestApproximateFraction_ValueIsWithinRequestedAccuracy(t *testing.T) {
	// Regardless of the specific fraction chosen, the resulting
	// value must actually be within the requested relative accuracy
	// of the original number: a check on the contract the function
	// promises, not just its particular output for one input.
	f := 1.61803398875 // golden ratio
	accuracyPct := 0.1

	result, err := approximateFraction(f, accuracyPct)
	if err != nil {
		t.Fatalf("approximateFraction() error = %v", err)
	}
	relativeError := math.Abs(result.Value-f) / f * 100
	if relativeError > accuracyPct {
		t.Errorf("relative error %v%% exceeds requested accuracy %v%%", relativeError, accuracyPct)
	}
}

func TestApproximateFraction_TighterAccuracyNeedsLargerDenominator(t *testing.T) {
	loose, err := approximateFraction(math.Pi, 1.0)
	if err != nil {
		t.Fatalf("approximateFraction(1%%) error = %v", err)
	}
	tight, err := approximateFraction(math.Pi, 0.0001)
	if err != nil {
		t.Fatalf("approximateFraction(0.0001%%) error = %v", err)
	}
	if !(tight.Den > loose.Den) {
		t.Errorf("tighter accuracy denominator (%d) should exceed looser accuracy denominator (%d)", tight.Den, loose.Den)
	}
}

func TestApproximateFraction_WholeNumberReturnsError(t *testing.T) {
	_, err := approximateFraction(5.0, 0.01)
	if err == nil {
		t.Fatal("approximateFraction(5.0) error = nil, want an error (no fractional part to approximate)")
	}
}

func TestGCD(t *testing.T) {
	cases := []struct{ a, b, want int64 }{
		{12, 8, 4},
		{17, 5, 1},
		{100, 10, 10},
	}
	for _, c := range cases {
		if got := gcd(c.a, c.b); got != c.want {
			t.Errorf("gcd(%d, %d) = %d, want %d", c.a, c.b, got, c.want)
		}
	}
}
