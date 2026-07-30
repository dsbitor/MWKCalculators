package main

import (
	"math"
	"testing"
)

func TestSolveGasLaw_StandardConditionsMatchesKnownMolarVolume(t *testing.T) {
	// One mole of an ideal gas at standard conditions (0 degC, 1
	// atmosphere) is well known to occupy 22.4140 liters,
	// independent of this program's own gas constant and formula.
	got, err := solveGasLaw(1.0, 0, 1.0, celsiusToKelvin)
	if err != nil {
		t.Fatalf("solveGasLaw() error = %v", err)
	}
	if math.Abs(got.VolumeLiters-22.4140) > 1e-3 {
		t.Errorf("VolumeLiters = %v, want approximately 22.4140", got.VolumeLiters)
	}
}

func TestSolveGasLaw_SolvingForEachQuantityAgrees(t *testing.T) {
	// A single self-consistent gas state solved for each of its four
	// quantities in turn, given the other three, must recover the
	// same full state every time. The state itself is built by
	// solving for temperature from three arbitrarily chosen
	// quantities, so it's guaranteed to actually satisfy PV=nRT
	// before being used as the round-trip's ground truth.
	full, err := solveGasLaw(2.0, 10.0, 0.8, 0)
	if err != nil {
		t.Fatalf("solveGasLaw(full) error = %v", err)
	}

	forP, err := solveGasLaw(0, full.VolumeLiters, full.Moles, full.TemperatureK)
	if err != nil {
		t.Fatalf("solveGasLaw(forP) error = %v", err)
	}
	checkClose(t, "forP.PressureAtm", forP.PressureAtm, full.PressureAtm)

	forV, err := solveGasLaw(full.PressureAtm, 0, full.Moles, full.TemperatureK)
	if err != nil {
		t.Fatalf("solveGasLaw(forV) error = %v", err)
	}
	checkClose(t, "forV.VolumeLiters", forV.VolumeLiters, full.VolumeLiters)

	forN, err := solveGasLaw(full.PressureAtm, full.VolumeLiters, 0, full.TemperatureK)
	if err != nil {
		t.Fatalf("solveGasLaw(forN) error = %v", err)
	}
	checkClose(t, "forN.Moles", forN.Moles, full.Moles)

	forT, err := solveGasLaw(full.PressureAtm, full.VolumeLiters, full.Moles, 0)
	if err != nil {
		t.Fatalf("solveGasLaw(forT) error = %v", err)
	}
	checkClose(t, "forT.TemperatureK", forT.TemperatureK, full.TemperatureK)
}

func TestSolveGasLaw_InsufficientDataReturnsError(t *testing.T) {
	_, err := solveGasLaw(1.0, 10.0, 0, 0)
	if err == nil {
		t.Fatal("solveGasLaw() error = nil, want an error for only two known quantities")
	}
}

func TestParsePressure(t *testing.T) {
	cases := []struct {
		s    string
		want float64
	}{
		{"1", 1},
		{"14.7psi", 14.7 / psiPerAtmosphere},
		{"100kpascal", 100 / kpaPerAtmosphere},
	}
	for _, c := range cases {
		if got := parsePressure(c.s); math.Abs(got-c.want) > 1e-9 {
			t.Errorf("parsePressure(%q) = %v, want %v", c.s, got, c.want)
		}
	}
}

func TestParseVolume(t *testing.T) {
	cases := []struct {
		s    string
		want float64
	}{
		{"10", 10},
		{"1cft", 1 / cubicFeetPerLiter},
		{"1cm", 1 / cubicMeterPerLiter},
	}
	for _, c := range cases {
		if got := parseVolume(c.s); math.Abs(got-c.want) > 1e-9 {
			t.Errorf("parseVolume(%q) = %v, want %v", c.s, got, c.want)
		}
	}
}

func TestParseTemperature(t *testing.T) {
	cases := []struct {
		s    string
		want float64
	}{
		{"293.16", 293.16},
		{"20c", 20 + celsiusToKelvin},
		{"68f", 5*(68.0-32)/9 + celsiusToKelvin},
	}
	for _, c := range cases {
		if got := parseTemperature(c.s); math.Abs(got-c.want) > 1e-9 {
			t.Errorf("parseTemperature(%q) = %v, want %v", c.s, got, c.want)
		}
	}
}

func checkClose(t *testing.T, name string, got, want float64) {
	t.Helper()
	if math.Abs(got-want) > 1e-6 {
		t.Errorf("%s = %v, want %v", name, got, want)
	}
}
