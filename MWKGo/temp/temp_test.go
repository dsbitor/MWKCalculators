package main

import (
	"math"
	"testing"
)

func TestConvertTemperature_KnownReferencePoints(t *testing.T) {
	// Water's freezing and boiling points, and absolute zero, are
	// independently known fixed points across all five scales.
	freezing := convertTemperature(0, 'c')
	checkClose(t, "freezing.Fahrenheit", freezing.Fahrenheit, 32)
	checkClose(t, "freezing.Kelvin", freezing.Kelvin, 273.18)

	boiling := convertTemperature(100, 'c')
	checkClose(t, "boiling.Fahrenheit", boiling.Fahrenheit, 212)

	absoluteZero := convertTemperature(0, 'k')
	checkClose(t, "absoluteZero.Rankine", absoluteZero.Rankine, 0)
}

func TestConvertTemperature_AllScalesAgreeRegardlessOfStartingScale(t *testing.T) {
	// Converting the same physical temperature starting from each of
	// the five scales must land on the same set of equivalent
	// values: an internal consistency check independent of which
	// scale's formula is used as the entry point.
	fromC := convertTemperature(37, 'c')
	fromF := convertTemperature(fromC.Fahrenheit, 'f')
	fromK := convertTemperature(fromC.Kelvin, 'k')
	fromR := convertTemperature(fromC.Rankine, 'r')
	fromE := convertTemperature(fromC.Reaumur, 'e')

	for name, got := range map[string]temperatures{"F": fromF, "K": fromK, "R": fromR, "E": fromE} {
		checkClose(t, name+".Centigrade", got.Centigrade, fromC.Centigrade)
		checkClose(t, name+".Fahrenheit", got.Fahrenheit, fromC.Fahrenheit)
	}
}

func TestConvertTemperature_UnrecognizedScaleDefaultsToFahrenheit(t *testing.T) {
	got := convertTemperature(98.6, 'x')
	want := convertTemperature(98.6, 'f')
	checkClose(t, "Centigrade", got.Centigrade, want.Centigrade)
}

func TestParseTemperatureInput(t *testing.T) {
	cases := []struct {
		line      string
		wantValue float64
		wantScale byte
		wantQuit  bool
	}{
		{"100f", 100, 'f', false},
		{"37c", 37, 'c', false},
		{"-40c", -40, 'c', false},
		{"300k", 300, 'k', false},
		{"70", 70, 'f', false},
		{"q", 0, 0, true},
		{"quit", 0, 0, true},
	}
	for _, c := range cases {
		value, scale, quit := parseTemperatureInput(c.line)
		if quit != c.wantQuit {
			t.Errorf("parseTemperatureInput(%q) quit = %v, want %v", c.line, quit, c.wantQuit)
			continue
		}
		if quit {
			continue
		}
		if math.Abs(value-c.wantValue) > 1e-9 || scale != c.wantScale {
			t.Errorf("parseTemperatureInput(%q) = %v, %q, want %v, %q", c.line, value, scale, c.wantValue, c.wantScale)
		}
	}
}

func TestPrintConversion_NegativeAbsoluteTemperatureIsRejected(t *testing.T) {
	if !printConversion("-10k") {
		t.Error("printConversion(-10k) = false, want true (rejected input should re-prompt, not quit)")
	}
}

func checkClose(t *testing.T, name string, got, want float64) {
	t.Helper()
	if math.Abs(got-want) > 1e-9 {
		t.Errorf("%s = %v, want %v", name, got, want)
	}
}
