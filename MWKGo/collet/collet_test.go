package main

import (
	"math"
	"testing"
)

func TestBoreDiameter(t *testing.T) {
	tests := []struct {
		name                   string
		sides                  int
		acrossFlats, slotWidth float64
		want                   float64
	}{
		// A zero slot width reduces the formula to
		// acrossFlats/cos(180/sides): the standard across-flats to
		// across-corners conversion for a regular polygon. For a
		// hexagon that ratio is 2/sqrt(3), a well known identity
		// independent of this code.
		{name: "zero slot width matches the hexagon across-corners identity", sides: 6, acrossFlats: 1, slotWidth: 0, want: 2 / math.Sqrt(3)},
		{name: "hexagon identity holds for a different across-flats dimension", sides: 6, acrossFlats: 0.5, slotWidth: 0, want: 0.5 * 2 / math.Sqrt(3)},
		{
			// A zero slot width square: across-corners = across-flats
			// * sqrt(2), another standard, independently verifiable
			// identity.
			name: "zero slot width matches the square across-corners identity", sides: 4, acrossFlats: 1, slotWidth: 0, want: math.Sqrt2,
		},
		{
			// The documented default input, evaluated against the
			// ported formula directly.
			name: "documented default input", sides: 6, acrossFlats: 3.0 / 16, slotWidth: 0.045,
			want: func() float64 {
				cosHalf := math.Cos(30 * math.Pi / 180)
				sinHalf := math.Sin(30 * math.Pi / 180)
				offset := 0.5*(3.0/16)/cosHalf - 0.5*0.045*sinHalf/cosHalf
				return 2 * math.Sqrt(0.25*0.045*0.045+offset*offset)
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := boreDiameter(tt.sides, tt.acrossFlats, tt.slotWidth)
			if diff := math.Abs(got - tt.want); diff > 1e-9 {
				t.Errorf("boreDiameter(%v, %v, %v) = %v, want %v", tt.sides, tt.acrossFlats, tt.slotWidth, got, tt.want)
			}
		})
	}
}
