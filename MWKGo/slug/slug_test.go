package main

import "testing"

func TestPlan_DrillsInwardFromTheFinalDiameter(t *testing.T) {
	// slug opens a hole, so the as-drilled circle must be smaller
	// than the final hole diameter (finish machining enlarges it out
	// to size). A flipped direction here would silently produce
	// plate's piece-freeing behaviour instead.
	finalDiameter := 3.0
	got := plan(finalDiameter, 0.05, 0.25, 0.05)
	if got.DrillingCircleDiameter >= finalDiameter {
		t.Errorf("DrillingCircleDiameter = %v, want it smaller than the final diameter %v", got.DrillingCircleDiameter, finalDiameter)
	}
}

func TestPlan_MatchesWorkedExamplesFromSlugTxt(t *testing.T) {
	// SLUG.TXT's own worked examples, reproduced end to end through
	// this program's specific wiring rather than just chaindrill's
	// Compute directly.
	tests := []struct {
		name          string
		drillDiameter float64
		wantHoleCount int
	}{
		{name: "1/4in drill", drillDiameter: 0.25, wantHoleCount: 27},
		{name: "3/8in drill", drillDiameter: 0.375, wantHoleCount: 18},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := plan(3, 0.05, tt.drillDiameter, 0.05)
			if got.HoleCount != tt.wantHoleCount {
				t.Errorf("HoleCount = %v, want %v", got.HoleCount, tt.wantHoleCount)
			}
		})
	}
}
