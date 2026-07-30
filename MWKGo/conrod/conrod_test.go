package main

import (
	"math"
	"testing"
)

func TestComputeConRodGeometry_DocumentedDefaultInput(t *testing.T) {
	got := computeConRodGeometry(2.4, 0.6, 1.5, 1.0)

	want := conRodGeometry{
		Phi:                    14.03624346792651,
		GudgeonToCrankDistance: 2.4738633753705965,
		D34:                    0.5,
		D45:                    0.24346584384264913,
		D35:                    0.25653415615735087,
		D23:                    0.2644293556038632,
		D13:                    0.5153882032022077,
		D12:                    0.25095884759834447,
		D14:                    0.125,
		D25:                    0.06221867190679133,
	}

	checkClose(t, "Phi", got.Phi, want.Phi)
	checkClose(t, "GudgeonToCrankDistance", got.GudgeonToCrankDistance, want.GudgeonToCrankDistance)
	checkClose(t, "D34", got.D34, want.D34)
	checkClose(t, "D45", got.D45, want.D45)
	checkClose(t, "D35", got.D35, want.D35)
	checkClose(t, "D23", got.D23, want.D23)
	checkClose(t, "D13", got.D13, want.D13)
	checkClose(t, "D12", got.D12, want.D12)
	checkClose(t, "D14", got.D14, want.D14)
	checkClose(t, "D25", got.D25, want.D25)
}

func TestComputeConRodGeometry_ZeroCrankRadiusIsDegenerate(t *testing.T) {
	// With no crank throw at all there is no "worst case" lateral
	// offset: the gudgeon pin never leaves the cylinder centerline,
	// phi is zero, and the gudgeon-to-crank distance is just the rod
	// length.
	got := computeConRodGeometry(2.4, 0.0, 1.5, 1.0)

	checkClose(t, "Phi", got.Phi, 0)
	checkClose(t, "GudgeonToCrankDistance", got.GudgeonToCrankDistance, 2.4)
	checkClose(t, "D45", got.D45, 0)
	checkClose(t, "D14", got.D14, 0)
	checkClose(t, "D25", got.D25, 0)
}

func TestComputeConRodGeometry_PythagoreanIdentity(t *testing.T) {
	// The gudgeon-to-crank distance is the hypotenuse of the crank
	// radius and rod length regardless of the other inputs: an
	// identity independent of this code's own formula for it.
	rodLength, crankRadius := 3.1, 0.9
	got := computeConRodGeometry(rodLength, crankRadius, 1.0, 1.0)
	want := math.Sqrt(rodLength*rodLength + crankRadius*crankRadius)
	checkClose(t, "GudgeonToCrankDistance", got.GudgeonToCrankDistance, want)
}

func checkClose(t *testing.T, name string, got, want float64) {
	t.Helper()
	if math.Abs(got-want) > 1e-9 {
		t.Errorf("%s = %v, want %v", name, got, want)
	}
}
