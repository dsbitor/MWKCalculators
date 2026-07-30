package main

import (
	"math"
	"testing"
)

func TestComputeSprocketDimensions_PitchDiameterMatchesRegularPolygonCircumradius(t *testing.T) {
	// The pitch diameter is, geometrically, the diameter of the
	// circle through the vertices of a regular polygon whose side
	// length is the chain pitch: an identity independent of this
	// code's own trigonometric form of the same formula.
	teeth := 15
	pitch := 0.5
	got := computeSprocketDimensions(teeth, pitch, 0.1)

	circumradius := pitch / (2 * math.Sin(math.Pi/float64(teeth)))
	want := 2 * circumradius
	if math.Abs(got.PitchDiameter-want) > 1e-9 {
		t.Errorf("PitchDiameter = %v, want %v", got.PitchDiameter, want)
	}
}

func TestComputeSprocketDimensions_DocumentedDefaultInput(t *testing.T) {
	// The documented default input (9 teeth, an odd count) evaluated
	// against the ported formula directly.
	got := computeSprocketDimensions(9, 1.0, 0.25)

	halfToothAngle := 180.0 / 9.0 * math.Pi / 180
	wantPD := 1.0 / math.Sin(halfToothAngle)
	wantOD := 1.0 * (0.6 + 1/math.Tan(halfToothAngle))
	wantCaliperFactor := wantPD * math.Cos(90.0/9.0*math.Pi/180)
	wantCD := wantCaliperFactor - 0.25
	wantMHD := 1.0*(1/math.Tan(halfToothAngle)-1) - 0.030

	if !got.OddTeeth {
		t.Errorf("OddTeeth = false, want true for 9 teeth")
	}
	if math.Abs(got.PitchDiameter-wantPD) > 1e-9 {
		t.Errorf("PitchDiameter = %v, want %v", got.PitchDiameter, wantPD)
	}
	if math.Abs(got.OutsideDiameter-wantOD) > 1e-9 {
		t.Errorf("OutsideDiameter = %v, want %v", got.OutsideDiameter, wantOD)
	}
	if math.Abs(got.CaliperFactor-wantCaliperFactor) > 1e-9 {
		t.Errorf("CaliperFactor = %v, want %v", got.CaliperFactor, wantCaliperFactor)
	}
	if math.Abs(got.CaliperDiameter-wantCD) > 1e-9 {
		t.Errorf("CaliperDiameter = %v, want %v", got.CaliperDiameter, wantCD)
	}
	if math.Abs(got.MaxHubDiameter-wantMHD) > 1e-9 {
		t.Errorf("MaxHubDiameter = %v, want %v", got.MaxHubDiameter, wantMHD)
	}
}

func TestComputeSprocketDimensions_EvenTeethSkipsCaliperFactor(t *testing.T) {
	// An even tooth count measures caliper diameter directly across
	// the sprocket, with no caliper factor involved at all.
	got := computeSprocketDimensions(10, 1.0, 0.25)
	if got.OddTeeth {
		t.Errorf("OddTeeth = true, want false for 10 teeth")
	}
	want := got.PitchDiameter - 0.25
	if math.Abs(got.CaliperDiameter-want) > 1e-9 {
		t.Errorf("CaliperDiameter = %v, want %v", got.CaliperDiameter, want)
	}
}
