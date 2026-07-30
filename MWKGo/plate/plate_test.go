package main

import "testing"

func TestPlan_DrillsOutwardFromTheFinalDiameter(t *testing.T) {
	// plate cuts a piece free, so the as-drilled circle must be
	// larger than the final plate diameter (finish machining removes
	// the scalloped edge down to size). A flipped direction here
	// would silently produce slug's hole-opening behaviour instead.
	finalDiameter := 3.0
	got := plan(finalDiameter, 0.05, 0.25, 0.05)
	if got.DrillingCircleDiameter <= finalDiameter {
		t.Errorf("DrillingCircleDiameter = %v, want it larger than the final diameter %v", got.DrillingCircleDiameter, finalDiameter)
	}
}
