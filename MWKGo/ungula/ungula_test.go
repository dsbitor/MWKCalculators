package main

import (
	"math"
	"testing"
)

// TestUngulaVolume_HalfDiameterSagitta checks the ported formula
// against a closed-form result independent of this code. When the
// sagitta equals the cylinder's radius (the wet region is exactly a
// semicircular half of the base), the volume is exactly
// (2/3)*height*radius^2.
//
// UNGULA.TXT states this reduction happens when the sagitta equals
// the full diameter instead. That claim does not hold: verified by
// numerical integration of the physical shape (see
// TestUngulaVolume_MatchesNumericalIntegration below and
// docs/calculators/ungula.md), the diameter case actually gives
// pi/2 * height * radius^2, not (2/3)*height*radius^2. The formula
// itself is correct throughout; only the original author's stated
// special case was wrong.
func TestUngulaVolume_HalfDiameterSagittaMatchesTwoThirdsIdentity(t *testing.T) {
	height, radius := 10.0, 1.0
	got := ungulaVolume(height, radius, radius) // sagitta == radius
	want := (2.0 / 3.0) * height * radius * radius
	if diff := math.Abs(got - want); diff > 1e-9 {
		t.Errorf("ungulaVolume(%v, %v, %v) = %v, want %v", height, radius, radius, got, want)
	}
}

func TestUngulaVolume_FullDiameterSagittaMatchesHalfCylinderIdentity(t *testing.T) {
	// At sagitta == diameter, the entire base is wet and depth rises
	// linearly from 0 to height across a full diameter; by symmetry
	// the average depth over the circle is height/2, so the volume
	// is average depth times base area: (height/2)*(pi*radius^2).
	// This is algebraically equal to pi/2 * height * radius^2.
	height, radius := 10.0, 1.0
	got := ungulaVolume(height, 2*radius, radius) // sagitta == diameter
	want := math.Pi / 2 * height * radius * radius
	if diff := math.Abs(got - want); diff > 1e-9 {
		t.Errorf("ungulaVolume(%v, %v, %v) = %v, want %v", height, 2*radius, radius, got, want)
	}
}

// TestUngulaVolume_MatchesNumericalIntegration independently verifies
// the formula at several sagitta values by brute-force numerical
// integration of the same physical shape it claims to compute: a
// circular base of the given radius, wet over the region from the
// wall at x=-radius to the chord at x=-radius+sagitta, with depth
// rising linearly from 0 at the chord to height at the wall.
func TestUngulaVolume_MatchesNumericalIntegration(t *testing.T) {
	height, radius := 10.0, 1.0

	numericalVolume := func(sagitta float64) float64 {
		chordX := -radius + sagitta
		const steps = 200000
		stepWidth := (chordX - (-radius)) / steps
		total := 0.0
		for i := 0; i < steps; i++ {
			x := -radius + (float64(i)+0.5)*stepWidth
			depth := height * (chordX - x) / sagitta
			halfChordLength := math.Sqrt(math.Max(radius*radius-x*x, 0))
			total += depth * 2 * halfChordLength * stepWidth
		}
		return total
	}

	for _, sagitta := range []float64{0.2, 0.5, 1.0, 1.5, 1.9} {
		got := ungulaVolume(height, sagitta, radius)
		want := numericalVolume(sagitta)
		if diff := math.Abs(got - want); diff > 1e-3 {
			t.Errorf("ungulaVolume(%v, %v, %v) = %v, want approximately %v (numerical integration)", height, sagitta, radius, got, want)
		}
	}
}

func TestUngulaVolume_ScalesLinearlyWithHeight(t *testing.T) {
	// The height parameter enters the formula as a simple linear
	// multiplier, independent of the geometry captured in phi: an
	// identity independent of this code's specific implementation.
	radius, sagitta := 1.0, 0.7
	single := ungulaVolume(1, sagitta, radius)
	doubled := ungulaVolume(2, sagitta, radius)
	if diff := math.Abs(doubled - 2*single); diff > 1e-9 {
		t.Errorf("ungulaVolume(2, ...) = %v, want exactly double ungulaVolume(1, ...) = %v", doubled, 2*single)
	}
}
