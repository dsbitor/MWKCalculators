package main

import (
	"math"
	"testing"
)

func TestPerfectKnurlFit(t *testing.T) {
	tests := []struct {
		name                  string
		wheelDiameter         float64
		toothCount            int
		nominalDiameter       float64
		wantToothSpacing      float64
		wantCrestCount        int
		wantWorkpieceDiameter float64
	}{
		{
			// A single-tooth wheel gives a tooth spacing of exactly
			// pi*wheelDiameter, and a nominal diameter equal to the
			// wheel diameter fits exactly one crest: pi cancels
			// exactly in floating point when a value is divided by
			// itself, so this case is exact, not approximate.
			name:          "single tooth wheel fits exactly one crest on a matching diameter",
			wheelDiameter: 1, toothCount: 1, nominalDiameter: 1,
			wantToothSpacing: math.Pi, wantCrestCount: 1, wantWorkpieceDiameter: 1,
		},
		{
			// Doubling the nominal diameter in the same single-tooth
			// setup fits exactly two crests, another exact case.
			name:          "single tooth wheel fits exactly two crests on double the diameter",
			wheelDiameter: 1, toothCount: 1, nominalDiameter: 2,
			wantToothSpacing: math.Pi, wantCrestCount: 2, wantWorkpieceDiameter: 2,
		},
		{
			// A nominal diameter just under the next whole crest
			// boundary must floor down, not round up.
			name:          "just under a crest boundary floors down",
			wheelDiameter: 1, toothCount: 1, nominalDiameter: 1.999,
			wantToothSpacing: math.Pi, wantCrestCount: 1, wantWorkpieceDiameter: 1,
		},
		{
			// A nominal diameter smaller than one tooth spacing fits
			// no whole crests at all, and the resulting workpiece
			// diameter is zero: an edge case the original program
			// does not guard against either.
			name:          "nominal diameter smaller than one tooth spacing fits no crests",
			wheelDiameter: 1, toothCount: 1, nominalDiameter: 0.5,
			wantToothSpacing: math.Pi, wantCrestCount: 0, wantWorkpieceDiameter: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := perfectKnurlFit(tt.wheelDiameter, tt.toothCount, tt.nominalDiameter)
			if diff := math.Abs(got.ToothSpacing - tt.wantToothSpacing); diff > 1e-9 {
				t.Errorf("ToothSpacing = %v, want %v", got.ToothSpacing, tt.wantToothSpacing)
			}
			if got.CrestCount != tt.wantCrestCount {
				t.Errorf("CrestCount = %v, want %v", got.CrestCount, tt.wantCrestCount)
			}
			if diff := math.Abs(got.WorkpieceDiameter - tt.wantWorkpieceDiameter); diff > 1e-9 {
				t.Errorf("WorkpieceDiameter = %v, want %v", got.WorkpieceDiameter, tt.wantWorkpieceDiameter)
			}
		})
	}
}

func TestPerfectKnurlFit_DocumentedDefaultInput(t *testing.T) {
	// The documented default input, evaluated against the ported
	// formula directly: a regression check on the arithmetic, not
	// an independently derived value.
	wheelDiameter, toothCount, nominalDiameter := 0.625, 40, 0.87

	toothSpacing := math.Pi * wheelDiameter / float64(toothCount)
	crestCount := int(math.Floor(math.Pi * nominalDiameter / toothSpacing))
	circumference := float64(crestCount) * toothSpacing

	got := perfectKnurlFit(wheelDiameter, toothCount, nominalDiameter)
	if diff := math.Abs(got.ToothSpacing - toothSpacing); diff > 1e-9 {
		t.Errorf("ToothSpacing = %v, want %v", got.ToothSpacing, toothSpacing)
	}
	if got.CrestCount != crestCount {
		t.Errorf("CrestCount = %v, want %v", got.CrestCount, crestCount)
	}
	if diff := math.Abs(got.Circumference - circumference); diff > 1e-9 {
		t.Errorf("Circumference = %v, want %v", got.Circumference, circumference)
	}
	if diff := math.Abs(got.WorkpieceDiameter - circumference/math.Pi); diff > 1e-9 {
		t.Errorf("WorkpieceDiameter = %v, want %v", got.WorkpieceDiameter, circumference/math.Pi)
	}
}
