package main

import (
	"math"
	"testing"
)

// exampleGears is CHANGE.DAT's own shipped change gear set.
var exampleGears = []int{20, 25, 30, 35, 40, 45, 50, 55, 60, 65, 70, 75, 80, 85, 90, 95, 100}

func TestFindChains_SingleStageExactRatio(t *testing.T) {
	// 20:40 is exactly 0.5, one of the gear set's own pairs.
	evaluations := 0
	solutions, truncated := findChains(exampleGears, 0.5, 0.01, 1, &evaluations)
	if truncated {
		t.Fatal("search unexpectedly truncated")
	}
	found := false
	for _, sol := range solutions {
		if len(sol.stages) == 1 && sol.stages[0].firstTeeth == 20 && sol.stages[0].secondTeeth == 40 {
			found = true
		}
	}
	if !found {
		t.Errorf("expected 20:40 among solutions, got %+v", solutions)
	}
}

func TestFindAllChains_MatchesChangeDATWorkedScenario(t *testing.T) {
	// CHANGE.DAT ships an effective leadscrew pitch of 8 tpi and
	// recommends starting with a maximum of 2 gear pairs. Requiring a
	// common 20 tpi thread from an 8 tpi leadscrew (ratio 0.4) should
	// be solvable with this gear set within two pairs.
	solutions, truncated := findAllChains([]float64{8}, exampleGears, 20.0, 0.5, 2)
	if truncated {
		t.Fatal("search unexpectedly truncated")
	}
	if len(solutions) == 0 {
		t.Fatal("expected at least one solution for a common 20 tpi thread")
	}
	for _, sol := range solutions {
		if sol.leadscrewPitch != 8 {
			t.Errorf("solution tagged with leadscrew pitch %v, want 8", sol.leadscrewPitch)
		}
		if math.Abs(sol.errorPercent) > 0.5+1e-9 {
			t.Errorf("solution %+v has error %v%%, outside requested 0.5%% tolerance", sol, sol.errorPercent)
		}
	}
}

func TestFindAllChains_TagsEachSolutionWithItsOwnPitch(t *testing.T) {
	solutions, _ := findAllChains([]float64{8, 10}, exampleGears, 20.0, 1.0, 1)
	sawEight, sawTen := false, false
	for _, sol := range solutions {
		gotRatio := 1.0
		for _, s := range sol.stages {
			gotRatio *= float64(s.firstTeeth) / float64(s.secondTeeth)
		}
		if math.Abs(gotRatio-sol.ratio) > 1e-9 {
			t.Errorf("stages %+v multiply to %v, want recorded ratio %v", sol.stages, gotRatio, sol.ratio)
		}
		wantRatio := sol.leadscrewPitch / 20.0
		errPercent := 100 * (sol.ratio - wantRatio) / wantRatio
		if math.Abs(errPercent-sol.errorPercent) > 1e-6 {
			t.Errorf("solution for pitch %v has recorded error %v%%, recomputed %v%%", sol.leadscrewPitch, sol.errorPercent, errPercent)
		}
		if sol.leadscrewPitch == 8 {
			sawEight = true
		}
		if sol.leadscrewPitch == 10 {
			sawTen = true
		}
	}
	if !sawEight || !sawTen {
		t.Errorf("expected solutions for both leadscrew pitches, sawEight=%v sawTen=%v", sawEight, sawTen)
	}
}

func TestFindChains_NoSolutionWhenUnreachable(t *testing.T) {
	evaluations := 0
	solutions, truncated := findChains(exampleGears, 1000.0, 0.001, 1, &evaluations)
	if len(solutions) != 0 {
		t.Errorf("expected no solutions, got %+v", solutions)
	}
	if truncated {
		t.Error("did not expect the small search space to be truncated")
	}
}
