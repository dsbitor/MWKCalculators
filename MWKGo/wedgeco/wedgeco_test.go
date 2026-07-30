package main

import (
	"math"
	"testing"
)

func TestWedgeVolume_HalfSagittaIsHalfCone(t *testing.T) {
	// A sagitta equal to the radius means the cutting plane passes
	// through the axis, splitting the cone exactly in half.
	got := wedgeVolume(10, 1, 1)
	want := 0.5 * math.Pi * 1 * 1 * 10 / 3
	if math.Abs(got-want) > 1e-9 {
		t.Errorf("wedgeVolume(10, 1, 1) = %v, want %v", got, want)
	}
}

func TestWedgeVolume_ComplementaryChordsSumToConeVolume(t *testing.T) {
	// The wedge cut at sagitta s and the wedge cut at sagitta 2r-s
	// (the complementary chord on the other side of the base circle)
	// together must account for the entire cone, an identity
	// independent of the formula's own derivation. Cross-checked by
	// numerical integration of the underlying cross-sectional
	// segment areas before trusting it here.
	h, r := 10.0, 1.0
	coneVolume := math.Pi * r * r * h / 3

	for _, s := range []float64{0.2, 0.5, 0.9, 1.3, 1.8} {
		sum := wedgeVolume(h, s, r) + wedgeVolume(h, 2*r-s, r)
		if math.Abs(sum-coneVolume) > 1e-6 {
			t.Errorf("wedgeVolume(%v)+wedgeVolume(%v) = %v, want cone volume %v", s, 2*r-s, sum, coneVolume)
		}
	}
}

func TestWedgeVolume_DocumentedDefaultInput(t *testing.T) {
	// Verified independently by numerical integration of the cone's
	// circular cross-sections against the cutting plane before being
	// trusted as an expected value here.
	got := wedgeVolume(10, 0.5, 1)
	want := 1.152640
	if math.Abs(got-want) > 1e-5 {
		t.Errorf("wedgeVolume(10, 0.5, 1) = %v, want %v", got, want)
	}
}
