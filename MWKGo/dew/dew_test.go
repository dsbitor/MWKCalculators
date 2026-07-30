package main

import (
	"math"
	"testing"
)

func TestDewPointCelsius(t *testing.T) {
	tests := []struct {
		name               string
		tempC, relHumidity float64
		want               float64
		tolerance          float64
	}{
		// A commonly cited example for the Magnus-Tetens formula: 20
		// degC at 50% relative humidity gives a dew point of
		// approximately 9.3 degC, a reference value independent of
		// this code.
		{name: "commonly cited Magnus formula example", tempC: 20, relHumidity: 0.5, want: 9.3, tolerance: 0.1},
		// 100% relative humidity means the air is already saturated:
		// the dew point equals the ambient temperature exactly,
		// since ln(1)=0 makes alpha reduce to the temperature term
		// alone. An identity independent of this code.
		{name: "100 percent humidity means dew point equals ambient temperature", tempC: 15, relHumidity: 1.0, want: 15, tolerance: 1e-9},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := dewPointCelsius(tt.tempC, tt.relHumidity)
			if diff := math.Abs(got - tt.want); diff > tt.tolerance {
				t.Errorf("dewPointCelsius(%v, %v) = %v, want %v", tt.tempC, tt.relHumidity, got, tt.want)
			}
		})
	}
}

func TestDewPointCelsius_LowerHumidityGivesLowerDewPoint(t *testing.T) {
	// For a fixed temperature, drier air must have a lower dew
	// point: a physical monotonicity property independent of this
	// code's specific constants.
	dry := dewPointCelsius(20, 0.2)
	humid := dewPointCelsius(20, 0.9)
	if dry >= humid {
		t.Errorf("dewPointCelsius(20, 0.2) = %v, want it lower than dewPointCelsius(20, 0.9) = %v", dry, humid)
	}
}
