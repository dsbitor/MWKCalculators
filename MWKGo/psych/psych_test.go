package main

import (
	"math"
	"testing"
)

func TestComputePsychrometer_EqualBulbsGivesFullSaturation(t *testing.T) {
	// When the wet and dry bulb read the same temperature, there is
	// no evaporative cooling at all: the air must already be at 100%
	// relative humidity, and the dew point must equal the measured
	// temperature. This follows from the vapor pressure and dew
	// point formulas being exact inverses of each other, verified
	// algebraically before trusting it as a test here.
	got := computePsychrometer(300, 20, 20)
	if math.Abs(got.RelativeHumidityPct-100) > 1e-6 {
		t.Errorf("RelativeHumidityPct = %v, want 100", got.RelativeHumidityPct)
	}
	if math.Abs(got.DewPointC-20) > 1e-6 {
		t.Errorf("DewPointC = %v, want 20", got.DewPointC)
	}
}

func TestComputePsychrometer_DriedAirHasLowerHumidityThanSaturated(t *testing.T) {
	// A larger wet/dry bulb spread means drier air: relative
	// humidity should decrease as the spread widens.
	small := computePsychrometer(300, 20, 19)
	large := computePsychrometer(300, 20, 10)
	if !(large.RelativeHumidityPct < small.RelativeHumidityPct) {
		t.Errorf("larger wet/dry spread (%v%%) should give lower humidity than smaller spread (%v%%)", large.RelativeHumidityPct, small.RelativeHumidityPct)
	}
}

func TestFahrenheitCelsiusRoundTrip(t *testing.T) {
	for _, f := range []float64{-40, 0, 32, 70, 98.6, 212} {
		c := fahrenheitToCelsius(f)
		back := celsiusToFahrenheit(c)
		if math.Abs(back-f) > 1e-9 {
			t.Errorf("round trip for %v: got %v", f, back)
		}
	}
}
