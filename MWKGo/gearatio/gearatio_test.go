package main

import (
	"math"
	"testing"
)

// exampleGears is GEARATIO.DAT's own shipped gear set.
var exampleGears = []int{24, 28, 32, 40, 44, 48, 56, 64, 72, 86, 100}

func TestFindRatios_SingleStageExactRatio(t *testing.T) {
	// 24:48 is exactly 0.5, one of the gear set's own pairs.
	solutions, truncated := findRatios(exampleGears, 0.5, 0.01, 1)
	if truncated {
		t.Fatal("search unexpectedly truncated")
	}
	found := false
	for _, sol := range solutions {
		if len(sol.stages) == 1 && sol.stages[0].firstTeeth == 24 && sol.stages[0].secondTeeth == 48 {
			found = true
			if math.Abs(sol.ratio-0.5) > 1e-9 {
				t.Errorf("ratio = %v, want 0.5", sol.ratio)
			}
		}
	}
	if !found {
		t.Errorf("expected 24:48 among solutions, got %+v", solutions)
	}
}

func TestFindRatios_EveryResultWithinTolerance(t *testing.T) {
	solutions, _ := findRatios(exampleGears, 0.35, 2.0, 2)
	if len(solutions) == 0 {
		t.Fatal("expected at least one solution")
	}
	for _, sol := range solutions {
		if math.Abs(sol.errorPercent) > 2.0+1e-9 {
			t.Errorf("solution %+v has error %v%%, outside requested 2%% tolerance", sol, sol.errorPercent)
		}
		gotRatio := 1.0
		// Recompute the ratio directly from the displayed stages to
		// confirm errorPercent matches what the stages actually produce
		// (independent of the search's own bookkeeping).
		for _, s := range sol.stages {
			gotRatio *= float64(s.firstTeeth) / float64(s.secondTeeth)
		}
		if math.Abs(gotRatio-sol.ratio) > 1e-9 {
			t.Errorf("stages %+v multiply to %v, want recorded ratio %v", sol.stages, gotRatio, sol.ratio)
		}
	}
}

func TestFindRatios_NoGearReusedWithinAChain(t *testing.T) {
	// With exactly 4 distinct gears and a 2-stage chain (4 gear
	// positions), the "no gear reused within a chain" rule forces every
	// solution to use all four gears exactly once each.
	gears := []int{10, 20, 30, 60}
	solutions, _ := findRatios(gears, 0.5, 50.0, 2)
	if len(solutions) == 0 {
		t.Fatal("expected at least one solution")
	}
	checked := 0
	for _, sol := range solutions {
		if len(sol.stages) != 2 {
			continue // findRatios also reports shorter chains that meet tolerance
		}
		checked++
		used := map[int]int{}
		for _, s := range sol.stages {
			used[s.firstTeeth]++
			used[s.secondTeeth]++
		}
		for _, g := range gears {
			if used[g] != 1 {
				t.Errorf("chain %+v uses gear %d %d time(s), want exactly once", sol.stages, g, used[g])
			}
		}
	}
	if checked == 0 {
		t.Fatal("expected at least one 2-stage solution")
	}
}

func TestFindRatios_NoSolutionWhenUnreachable(t *testing.T) {
	// A ratio far outside anything the gear set can produce, with a
	// tight tolerance and a short chain, should find nothing.
	solutions, truncated := findRatios(exampleGears, 1000.0, 0.001, 1)
	if len(solutions) != 0 {
		t.Errorf("expected no solutions, got %+v", solutions)
	}
	if truncated {
		t.Error("did not expect the small search space to be truncated")
	}
}

func TestFindRatios_IdenticalGearsSkipped(t *testing.T) {
	// Two gears of the same size (as in GEARATIO.DAT's own alternate
	// data sets, e.g. two 20-tooth gears) would produce a useless 1:1
	// ratio and must be skipped rather than reported as a solution.
	gears := []int{20, 20, 25, 30}
	solutions, _ := findRatios(gears, 1.0, 0.001, 1)
	for _, sol := range solutions {
		if sol.stages[0].firstTeeth == sol.stages[0].secondTeeth {
			t.Errorf("expected 1:1 pairs to be skipped, got %+v", sol)
		}
	}
}
