package main

import (
	"math"
	"testing"
)

func TestMandrelDiameter(t *testing.T) {
	tests := []struct {
		name                 string
		wt                   wireType
		wireDiameter         float64
		springInsideDiameter float64
	}{
		{name: "documented default input, music wire", wt: musicWire, wireDiameter: 0.040, springInsideDiameter: 0.203},
		{name: "same dimensions, phosphorus bronze", wt: phosphorusBronze, wireDiameter: 0.040, springInsideDiameter: 0.203},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Evaluated against the ported formula directly: a
			// regression check on the arithmetic, since Hiraoka's
			// coefficients are empirical and have no independent
			// closed-form identity to check against.
			averageSpringDiameter := tt.springInsideDiameter + tt.wireDiameter
			factor := constantCoefficient[tt.wt] + linearCoefficient[tt.wt]*averageSpringDiameter/tt.wireDiameter
			want := factor*averageSpringDiameter - tt.wireDiameter

			got := mandrelDiameter(tt.wt, tt.wireDiameter, tt.springInsideDiameter)
			if diff := math.Abs(got - want); diff > 1e-9 {
				t.Errorf("mandrelDiameter(%v, %v, %v) = %v, want %v", tt.wt, tt.wireDiameter, tt.springInsideDiameter, got, want)
			}
		})
	}
}

func TestNormalizeWireType(t *testing.T) {
	tests := []struct {
		name  string
		input int
		want  wireType
	}{
		{name: "zero is music wire", input: 0, want: musicWire},
		{name: "one is phosphorus bronze", input: 1, want: phosphorusBronze},
		{name: "any other value falls back to music wire", input: 2, want: musicWire},
		{name: "negative value falls back to music wire", input: -1, want: musicWire},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := normalizeWireType(tt.input); got != tt.want {
				t.Errorf("normalizeWireType(%v) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestMandrelDiameter_MusicAndBronzeGiveDifferentResults(t *testing.T) {
	// The two wire types use distinct empirical coefficients, so the
	// same dimensions must produce different recommendations;
	// otherwise the wire type selection would have no effect.
	music := mandrelDiameter(musicWire, 0.040, 0.203)
	bronze := mandrelDiameter(phosphorusBronze, 0.040, 0.203)
	if music == bronze {
		t.Errorf("mandrelDiameter gave the same result (%v) for both wire types, want them to differ", music)
	}
}
