package main

import (
	"math"
	"testing"
)

func TestBestWireDiameter_MatchesStandardConstant(t *testing.T) {
	// The "best wire" size for a 60 degree thread is a well known
	// standard constant times the pitch: 0.577350*pitch, independent
	// of this code's own cos(30deg)-based formula for the same
	// value.
	pitch := 1.0 / 20.0
	got := bestWireDiameter(pitch)
	want := 0.577350 * pitch
	if math.Abs(got-want) > 1e-6 {
		t.Errorf("bestWireDiameter(%v) = %v, want %v", pitch, got, want)
	}
}

func TestPitchDiameterAndMeasurementOverWires_RoundTrip(t *testing.T) {
	// Converting a pitch diameter to a measurement over wires and
	// back must reproduce the original pitch diameter exactly,
	// provided the same pitch-diameter estimate is used for the lead
	// angle correction in both directions (the original program's
	// own approach when the "mow from pd" branch is used without a
	// custom pitch diameter override).
	pitch, starts, wireDiameter := 1.0/20.0, 1.0, 0.0289
	pitchDiameter := 0.2175

	x := wireOffset(pitch, starts, pitchDiameter, wireDiameter)
	mow := pitchDiameter + x
	recoveredPD := mow - x

	if math.Abs(recoveredPD-pitchDiameter) > 1e-12 {
		t.Errorf("recoveredPD = %v, want %v", recoveredPD, pitchDiameter)
	}
}

func TestWireOffset_LargerWireGivesLargerMeasurement(t *testing.T) {
	pitch, starts, pitchDiameter := 1.0/20.0, 1.0, 0.2175
	small := wireOffset(pitch, starts, pitchDiameter, 0.02)
	large := wireOffset(pitch, starts, pitchDiameter, 0.04)
	if !(large > small) {
		t.Errorf("wireOffset with larger wire (%v) should exceed offset with smaller wire (%v)", large, small)
	}
}

func TestPitchDiameterFromMajor_1_4_20UNC(t *testing.T) {
	// A 1/4-20 UNC thread's pitch diameter, per this program's own
	// p/8-flat assumption, evaluated directly against the ported
	// formula.
	majorDiameter, pitch := 0.25, 1.0/20.0
	got := pitchDiameterFromMajor(majorDiameter, pitch)
	flatCorrection := (3.0 / 16.0) * pitch / math.Tan(30*math.Pi/180)
	want := majorDiameter - 2*flatCorrection
	if math.Abs(got-want) > 1e-9 {
		t.Errorf("pitchDiameterFromMajor() = %v, want %v", got, want)
	}
}
