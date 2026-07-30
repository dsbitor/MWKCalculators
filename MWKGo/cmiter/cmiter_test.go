package main

import (
	"math"
	"testing"
)

func TestComputeCompoundMiterAngles_NinetyDegreeSlope(t *testing.T) {
	// A 90 degree slope is a flat form cut in the table's plane: the
	// miter gauge angle is always 90, and the mitred blade tilt is
	// just the polygon's own half angle (180/sides), independent of
	// this code's general-case trigonometry.
	got := computeCompoundMiterAngles(4, 90)
	if got.MiterGaugeAngleDeg != 90 {
		t.Errorf("MiterGaugeAngleDeg = %v, want 90", got.MiterGaugeAngleDeg)
	}
	if diff := math.Abs(got.BladeTiltMitredDeg - 45); diff > 1e-9 {
		t.Errorf("BladeTiltMitredDeg = %v, want 45 (180/4)", got.BladeTiltMitredDeg)
	}
	if got.BladeTiltButtedDeg != 0 {
		t.Errorf("BladeTiltButtedDeg = %v, want 0", got.BladeTiltButtedDeg)
	}
}

func TestComputeCompoundMiterAngles_ZeroSlope(t *testing.T) {
	// A zero degree slope (vertical, flat-sided form) makes the
	// miter gauge angle the complement of the polygon's half angle:
	// atan(cot(x)) = 90-x, a standard trigonometric identity
	// independent of this code.
	got := computeCompoundMiterAngles(4, 0)
	want := 90 - 180.0/4
	if diff := math.Abs(got.MiterGaugeAngleDeg - want); diff > 1e-9 {
		t.Errorf("MiterGaugeAngleDeg = %v, want %v", got.MiterGaugeAngleDeg, want)
	}
	if got.BladeTiltMitredDeg != 0 {
		t.Errorf("BladeTiltMitredDeg = %v, want 0", got.BladeTiltMitredDeg)
	}
	if got.BladeTiltButtedDeg != 90 {
		t.Errorf("BladeTiltButtedDeg = %v, want 90", got.BladeTiltButtedDeg)
	}
}

func TestComputeCompoundMiterAngles_DocumentedDefaultInput(t *testing.T) {
	// The documented default input (a 4 sided form at 30 degree
	// slope), evaluated against the ported formula directly.
	sides, slope := 4, 30.0
	halfAngleDeg := 180 / float64(sides)
	wantMiterGauge := math.Atan(1/(math.Cos(slope*math.Pi/180)*math.Tan(halfAngleDeg*math.Pi/180))) * 180 / math.Pi
	cosMiterGauge := math.Cos(wantMiterGauge * math.Pi / 180)
	tanSlope := math.Tan(slope * math.Pi / 180)
	wantBTM := math.Atan(cosMiterGauge*tanSlope) * 180 / math.Pi
	wantBTB := math.Atan(cosMiterGauge/tanSlope) * 180 / math.Pi

	got := computeCompoundMiterAngles(sides, slope)
	if diff := math.Abs(got.MiterGaugeAngleDeg - wantMiterGauge); diff > 1e-9 {
		t.Errorf("MiterGaugeAngleDeg = %v, want %v", got.MiterGaugeAngleDeg, wantMiterGauge)
	}
	if diff := math.Abs(got.BladeTiltMitredDeg - wantBTM); diff > 1e-9 {
		t.Errorf("BladeTiltMitredDeg = %v, want %v", got.BladeTiltMitredDeg, wantBTM)
	}
	if diff := math.Abs(got.BladeTiltButtedDeg - wantBTB); diff > 1e-9 {
		t.Errorf("BladeTiltButtedDeg = %v, want %v", got.BladeTiltButtedDeg, wantBTB)
	}
}
