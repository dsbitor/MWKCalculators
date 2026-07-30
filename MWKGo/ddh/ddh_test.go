package main

import "testing"

// exampleHoleCircles and exampleGears are DDH.DAT's own shipped
// hole-circle plates (across its three plates A, B, C) and change gear
// set.
var exampleHoleCircles = []int{15, 16, 17, 18, 19, 20, 21, 23, 27, 29, 31, 33, 37, 39, 41, 43, 47, 49}
var exampleGears = []int{24, 24, 28, 32, 40, 44, 48, 56, 64, 72, 86, 100}

const exampleRatio = 40 // DDH.DAT's own worm gear ratio

func TestTurnsRequired_KnownExample(t *testing.T) {
	// 14 divisions with a 40:1 worm gear ratio is DIVHEAD's own worked
	// example (DDH's math for this stage is identical to divhead's):
	// 40/14 = 2 & 6/7.
	wholeTurns, remainderNum, remainderDenom := turnsRequired(40, 14)
	if wholeTurns != 2 || remainderNum != 6 || remainderDenom != 7 {
		t.Errorf("turnsRequired(40,14) = (%d,%d,%d), want (2,6,7)", wholeTurns, remainderNum, remainderDenom)
	}
}

func TestPlateSolutions_NoSolutionForcesArtificialDenominator(t *testing.T) {
	// 67 divisions (DDH.C's own default) with a 40:1 ratio has no plain
	// hole-plate solution among DDH.DAT's own plates: none of them is a
	// multiple of 67 (all are far smaller than 67).
	_, remainderNum, remainderDenom := turnsRequired(exampleRatio, 67)
	if remainderNum == 0 {
		t.Fatal("expected a non-exact division for this example")
	}
	solutions := plateSolutions(exampleHoleCircles, remainderNum, remainderDenom)
	if len(solutions) != 0 {
		t.Errorf("expected no plain hole-plate solution for 67 divisions, got %+v", solutions)
	}
}

func TestFindDifferentialSolutions_67DivisionsHasASolution(t *testing.T) {
	// This is the exact scenario DDH.TXT motivates: 67 divisions has no
	// plain hole-plate solution (see above), so a differential solution
	// using the change gears must be found instead.
	solutions, truncated := findDifferentialSolutions(exampleRatio, 67, exampleHoleCircles, exampleGears, 1.0, 1)
	if truncated {
		t.Fatal("search unexpectedly truncated")
	}
	if len(solutions) == 0 {
		t.Fatal("expected at least one differential solution for 67 divisions")
	}
	for _, sol := range solutions {
		gotRatio := 1.0
		for _, s := range sol.stages {
			gotRatio *= float64(s.firstTeeth) / float64(s.secondTeeth)
		}
		if diff := gotRatio - sol.achievedRatio; diff > 1e-9 || diff < -1e-9 {
			t.Errorf("stages %+v multiply to %v, want recorded achievedRatio %v", sol.stages, gotRatio, sol.achievedRatio)
		}
		if sol.holePlate <= 0 {
			t.Errorf("solution has non-positive hole plate %d", sol.holePlate)
		}
		if sol.indexIncrement <= 0 {
			t.Errorf("solution has non-positive indexing increment %d", sol.indexIncrement)
		}
	}
}

func TestFindDifferentialSolutions_NoSolutionWhenGearSetTooSparse(t *testing.T) {
	// A gear set of only two very close gears can rarely, if ever,
	// match an arbitrary required ratio within a tight tolerance.
	sparse := []int{20, 21}
	_, truncated := findDifferentialSolutions(exampleRatio, 67, exampleHoleCircles, sparse, 0.0001, 1)
	if truncated {
		t.Fatal("did not expect the small search space to be truncated")
	}
	// Not asserting zero solutions (a coincidental match is possible),
	// only that the search completes and returns cleanly.
}

func TestRapidIndexSolution_ExactMultiple(t *testing.T) {
	step, ok := rapidIndexSolution(24, 8)
	if !ok || step != 3 {
		t.Errorf("rapidIndexSolution(24,8) = (%d,%v), want (3,true)", step, ok)
	}
}

func TestRapidIndexSolution_NoPlateOrNoMultiple(t *testing.T) {
	if _, ok := rapidIndexSolution(-1, 8); ok {
		t.Error("expected no rapid index solution when no plate is configured (-1)")
	}
	if _, ok := rapidIndexSolution(24, 7); ok {
		t.Error("expected no rapid index solution when 24 is not a multiple of 7")
	}
}
