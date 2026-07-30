package main

import (
	"math"
	"testing"
)

func TestDiskDiameter(t *testing.T) {
	tests := []struct {
		name             string
		divisions        float64
		mountingDiameter float64
		want             float64
	}{
		// Six equal circles arranged around a central circle of the
		// same radius, all mutually touching, is the classic
		// hexagonal circle-packing arrangement: an identity
		// independent of this code. The disk diameter must equal
		// the mounting circle's diameter exactly.
		{name: "six divisions matches the hexagonal packing identity", divisions: 6, mountingDiameter: 112, want: 112},
		{name: "hexagonal packing identity holds for a different mounting diameter", divisions: 6, mountingDiameter: 10, want: 10},
		// The documented default input, evaluated against the
		// ported formula directly using math.Sin rather than a
		// hand-typed literal.
		{
			name: "documented default input", divisions: 14, mountingDiameter: 112,
			want: 112 * math.Sin(0.5*(360.0/14)*math.Pi/180) / (1 - math.Sin(0.5*(360.0/14)*math.Pi/180)),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := diskDiameter(tt.divisions, tt.mountingDiameter)
			if diff := math.Abs(got - tt.want); diff > 1e-9 {
				t.Errorf("diskDiameter(%v, %v) = %v, want %v", tt.divisions, tt.mountingDiameter, got, tt.want)
			}
		})
	}
}

func TestDiskDiameter_TwoDivisions_ProducesNonFiniteResult(t *testing.T) {
	// At exactly two divisions the half-angle is 90 degrees, whose
	// sine is 1, making the denominator (1 - sin) zero. The
	// original program does not guard against this either; two
	// disks cannot mutually touch a central disk and each other in
	// the way this construction assumes, so the failure is
	// immediately visible as an infinite result rather than a
	// silently wrong number.
	got := diskDiameter(2, 112)
	if !math.IsInf(got, 0) {
		t.Errorf("diskDiameter(2, 112) = %v, want +Inf", got)
	}
}
