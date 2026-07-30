package main

import (
	"math"
	"testing"
)

func TestSearchGearRatios_DocumentedDefaultInputFindsSolutions(t *testing.T) {
	// The documented default input (ratio 1.9945, 0.01% tolerance, 2
	// pairs, teeth 16-40) is independently confirmed, via a
	// straightforward brute-force check in a separate script, to
	// have 48 solutions with a best error of about 0.000275%.
	solutions, truncated, err := searchGearRatios(1.9945, 0.01, 2, 16, 40)
	if err != nil {
		t.Fatalf("searchGearRatios() error = %v", err)
	}
	if truncated {
		t.Fatal("truncated = true, want false (this search space is well within the evaluation limit)")
	}
	if len(solutions) != 48 {
		t.Errorf("len(solutions) = %d, want 48", len(solutions))
	}
}

func TestSearchGearRatios_EverySolutionMeetsTolerance(t *testing.T) {
	// Whatever solutions are returned, each one's own reported error
	// must actually be within the requested tolerance: a check on
	// the search's own contract, not just a specific count.
	desiredRatio, allowedErrorPct := 1.9945, 0.01
	solutions, _, err := searchGearRatios(desiredRatio, allowedErrorPct, 2, 16, 40)
	if err != nil {
		t.Fatalf("searchGearRatios() error = %v", err)
	}
	for _, s := range solutions {
		ratio := 1.0
		for _, pair := range s.Pairs {
			ratio *= float64(pair[0]) / float64(pair[1])
		}
		if math.Abs(ratio-s.Ratio) > 1e-12 {
			t.Errorf("solution pairs %v recompute to ratio %v, want %v", s.Pairs, ratio, s.Ratio)
		}
		errPct := 100 * (ratio - desiredRatio) / desiredRatio
		if math.Abs(errPct) > allowedErrorPct {
			t.Errorf("solution %v has error %v%%, exceeds tolerance %v%%", s.Pairs, errPct, allowedErrorPct)
		}
	}
}

func TestSearchGearRatios_NoSolutionForImpossibleRatio(t *testing.T) {
	// A ratio far outside what any pair of teeth in [16,18] can
	// produce, with an extremely tight tolerance, should find
	// nothing.
	solutions, _, err := searchGearRatios(100.0, 0.0001, 1, 16, 18)
	if err != nil {
		t.Fatalf("searchGearRatios() error = %v", err)
	}
	if len(solutions) != 0 {
		t.Errorf("len(solutions) = %d, want 0", len(solutions))
	}
}

func TestSearchGearRatios_TooManyPairsReturnsError(t *testing.T) {
	_, _, err := searchGearRatios(2.0, 1.0, 6, 16, 40)
	if err == nil {
		t.Fatal("searchGearRatios() error = nil, want an error for more than 5 pairs")
	}
}

func TestSearchGearRatios_SinglePairExactMatch(t *testing.T) {
	// A single pair of 20:10 teeth gives an exact ratio of 2.0: the
	// simplest possible case, verified directly.
	solutions, _, err := searchGearRatios(2.0, 0.001, 1, 10, 20)
	if err != nil {
		t.Fatalf("searchGearRatios() error = %v", err)
	}
	found := false
	for _, s := range solutions {
		if s.Pairs[0] == [2]int{20, 10} {
			found = true
			checkClose(t, "Ratio", s.Ratio, 2.0)
		}
	}
	if !found {
		t.Error("expected to find the pair (20, 10) among the solutions")
	}
}

func checkClose(t *testing.T, name string, got, want float64) {
	t.Helper()
	if math.Abs(got-want) > 1e-9 {
		t.Errorf("%s = %v, want %v", name, got, want)
	}
}
