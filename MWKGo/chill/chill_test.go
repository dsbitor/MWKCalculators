package main

import (
	"math"
	"testing"
)

func TestWindChillF(t *testing.T) {
	tests := []struct {
		name           string
		tempF, windMph float64
		want           float64
		tolerance      float64
	}{
		// A commonly published NWS example: 0F with a 15mph wind
		// gives a wind chill of approximately -19F, a reference
		// value independent of this code. The loose tolerance
		// reflects that the reference figure itself is rounded to
		// the nearest whole degree.
		{name: "published NWS example, 0F at 15mph", tempF: 0, windMph: 15, want: -19, tolerance: 0.6},
		{
			// The documented default input, evaluated against the
			// ported formula directly.
			name: "documented default input", tempF: 30, windMph: 25,
			want:      35.74 + 0.6215*30 - math.Pow(25, 0.16)*(35.75-0.4275*30),
			tolerance: 1e-9,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := windChillF(tt.tempF, tt.windMph)
			if diff := math.Abs(got - tt.want); diff > tt.tolerance {
				t.Errorf("windChillF(%v, %v) = %v, want %v", tt.tempF, tt.windMph, got, tt.want)
			}
		})
	}
}

func TestWindChillF_HigherWindGivesLowerChill(t *testing.T) {
	// For a fixed temperature, a stronger wind must always produce a
	// colder (or equal) equivalent temperature: a physical
	// monotonicity property independent of this code's specific
	// coefficients.
	calm := windChillF(20, 5)
	windy := windChillF(20, 40)
	if windy >= calm {
		t.Errorf("windChillF(20, 40) = %v, want it colder than windChillF(20, 5) = %v", windy, calm)
	}
}
