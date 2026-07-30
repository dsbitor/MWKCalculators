package main

import (
	"math"
	"testing"
)

func TestComputeThreadEngagement_DocumentedDefaultInput(t *testing.T) {
	// The documented Imperial default input (1/4-20), evaluated
	// against the ported formula directly.
	diameter, pitch := 0.25, 1.0/20.0
	got := computeThreadEngagement(diameter, pitch)

	wantPCD := diameter - 0.64952*pitch
	wantTensile := 0.25 * math.Pi * (diameter - 0.938194*pitch) * (diameter - 0.938194*pitch)
	wantShear := 0.5 * math.Pi * wantPCD
	wantLength := 2 * wantTensile / wantShear

	checkClose(t, "PitchCircleDiameter", got.PitchCircleDiameter, wantPCD)
	checkClose(t, "TensileArea", got.TensileArea, wantTensile)
	checkClose(t, "ShearArea", got.ShearArea, wantShear)
	checkClose(t, "EngagementLength", got.EngagementLength, wantLength)
}

func TestComputeThreadEngagement_NegligiblePitchApproachesDiameter(t *testing.T) {
	// As the thread pitch shrinks toward zero, the tensile area
	// approaches the full bolt cross section (0.25*pi*D^2) and the
	// shear area approaches 0.5*pi*D, so their ratio, doubled,
	// approaches D itself: an independent limiting-case identity
	// rather than a re-run of the formula.
	diameter := 1.0
	got := computeThreadEngagement(diameter, 1e-9)
	if math.Abs(got.EngagementLength-diameter) > 1e-6 {
		t.Errorf("EngagementLength = %v, want approximately %v (the diameter)", got.EngagementLength, diameter)
	}
}

func TestComputeThreadEngagement_LengthConsistentWithAreas(t *testing.T) {
	// le is defined as 2*tensileArea/shearArea; this checks the
	// three returned fields actually satisfy that relationship
	// rather than having drifted apart during refactoring.
	got := computeThreadEngagement(0.5, 0.05)
	want := 2 * got.TensileArea / got.ShearArea
	checkClose(t, "EngagementLength", got.EngagementLength, want)
}

func checkClose(t *testing.T, name string, got, want float64) {
	t.Helper()
	if math.Abs(got-want) > 1e-9 {
		t.Errorf("%s = %v, want %v", name, got, want)
	}
}
