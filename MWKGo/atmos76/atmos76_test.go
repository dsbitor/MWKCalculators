package main

import (
	"math"
	"testing"
)

func TestComputeAtmosphereRatios_SeaLevelIsExactlyOne(t *testing.T) {
	got := computeAtmosphereRatios(0)
	checkClose(t, "TemperatureRatio", got.TemperatureRatio, 1.0)
	checkClose(t, "PressureRatio", got.PressureRatio, 1.0)
	checkClose(t, "DensityRatio", got.DensityRatio, 1.0)
}

func TestComputeAtmosphereRatios_MatchesTableBreakpointsExactly(t *testing.T) {
	// At each layer's own tabulated base geopotential height, the
	// computed ratios must exactly reproduce that layer's tabulated
	// values (deltaH = 0 there), independent of the interpolation
	// formula used between breakpoints. Geopotential height is
	// converted back to geometric altitude for the call, since
	// computeAtmosphereRatios takes geometric altitude as input.
	for i, hGeopotential := range layerBaseHeightKm {
		if i == len(layerBaseHeightKm)-1 {
			continue // last entry is only a table bound, not itself a modeled layer
		}
		altitudeKm := hGeopotential * earthRadiusKm / (earthRadiusKm - hGeopotential)
		got := computeAtmosphereRatios(altitudeKm)

		wantTempRatio := layerBaseTempK[i] / layerBaseTempK[0]
		if diff := math.Abs(got.TemperatureRatio - wantTempRatio); diff > 1e-6 {
			t.Errorf("layer %d: TemperatureRatio = %v, want %v", i, got.TemperatureRatio, wantTempRatio)
		}
		if diff := math.Abs(got.PressureRatio - layerBasePressRatio[i]); diff > 1e-6 {
			t.Errorf("layer %d: PressureRatio = %v, want %v", i, got.PressureRatio, layerBasePressRatio[i])
		}
	}
}

func TestComputeAtmosphereRatios_RatiosDecreaseWithAltitude(t *testing.T) {
	low := computeAtmosphereRatios(5)
	high := computeAtmosphereRatios(20)
	if !(high.PressureRatio < low.PressureRatio) {
		t.Errorf("pressure ratio at 20km (%v) should be less than at 5km (%v)", high.PressureRatio, low.PressureRatio)
	}
	if !(high.DensityRatio < low.DensityRatio) {
		t.Errorf("density ratio at 20km (%v) should be less than at 5km (%v)", high.DensityRatio, low.DensityRatio)
	}
}

func checkClose(t *testing.T, name string, got, want float64) {
	t.Helper()
	if math.Abs(got-want) > 1e-9 {
		t.Errorf("%s = %v, want %v", name, got, want)
	}
}
