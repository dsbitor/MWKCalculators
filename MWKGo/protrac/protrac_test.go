package main

import (
	"math"
	"testing"
)

func TestClosedAngleDeg_DocumentedDefaultInput(t *testing.T) {
	got := closedAngleDeg(3.0, 0.25, 0.22)
	want := 2 * math.Asin(0.5*(0.25+0.22)/3.0) * 180 / math.Pi
	if math.Abs(got-want) > 1e-9 {
		t.Errorf("closedAngleDeg() = %v, want %v", got, want)
	}
}

func TestSeparationFromAngle_ZeroAngleReturnsClosedGap(t *testing.T) {
	// At an included angle of zero, the pin separation must be
	// exactly the closed-gap distance the protractor was calibrated
	// with, since that's what "zero angle" means by definition.
	radius, pinDiameter, closedGap := 3.0, 0.25, 0.22
	closedAngle := closedAngleDeg(radius, pinDiameter, closedGap)

	got := separationFromAngle(radius, pinDiameter, 0.0, closedAngle)
	if math.Abs(got-closedGap) > 1e-9 {
		t.Errorf("separationFromAngle(0) = %v, want closedGap %v", got, closedGap)
	}
}

func TestAngleAndSeparationAreInverses(t *testing.T) {
	// Converting an angle to a separation and back must reproduce
	// the original angle, regardless of the specific formulas used
	// for either direction.
	radius, pinDiameter, closedGap := 3.0, 0.25, 0.22
	closedAngle := closedAngleDeg(radius, pinDiameter, closedGap)

	for _, wantAngle := range []float64{0, 5, 15, 30, 45} {
		separation := separationFromAngle(radius, pinDiameter, wantAngle, closedAngle)
		gotAngle := angleFromSeparation(radius, pinDiameter, separation, closedAngle)
		if math.Abs(gotAngle-wantAngle) > 1e-6 {
			t.Errorf("round trip for angle %v: got %v", wantAngle, gotAngle)
		}
	}
}
