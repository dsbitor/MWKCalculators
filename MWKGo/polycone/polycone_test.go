package main

import (
	"math"
	"testing"
)

func TestComputePolycone_DocumentedDefaultInput(t *testing.T) {
	// Verified independently via direct computation from the ported
	// formula before being trusted as expected values here; the
	// identity tests below check the formula itself against
	// independent geometric facts.
	got := computePolycone(4, 6, 10)

	checkClose(t, "Base.Area", got.Base.Area, 36.0)
	checkClose(t, "Facet.Height", got.Facet.Height, 10.44030650891055)
	checkClose(t, "Facet.EdgeLength", got.Facet.EdgeLength, 10.862780491200215)
	checkClose(t, "TotalVolume", got.TotalVolume, 120.0)
}

func TestComputePolycone_PythagoreanIdentities(t *testing.T) {
	// Facet height and edge length are each the hypotenuse of the
	// polycone's height and a base dimension (apothem or
	// circumradius respectively): identities independent of this
	// code's own separate formula for each.
	got := computePolycone(6, 4, 8)
	checkClose(t, "Facet.Height", got.Facet.Height, math.Hypot(8, got.Base.Apothem))
	checkClose(t, "Facet.EdgeLength", got.Facet.EdgeLength, math.Hypot(8, got.Base.Radius))
}

func TestComputePolycone_TipAngleAndBaseAngleAreComplementary(t *testing.T) {
	// By construction, the tip angle is defined as twice the
	// complement of the base angle, so base angle + half the tip
	// angle must always equal 90 degrees, independent of the
	// specific numeric inputs.
	got := computePolycone(5, 3, 7)
	checkClose(t, "BaseAngle+TipAngle/2", got.Facet.BaseAngleDeg+got.Facet.TipAngleDeg/2, 90)
}

func TestComputePolycone_FaceAngleTangentMatchesRiseOverRun(t *testing.T) {
	// The face angle is, geometrically, the angle a facet makes with
	// the horizontal base: its tangent must equal the polycone's
	// height divided by the base apothem (rise over run along the
	// facet's own slope line), independent of this code's own
	// acos-based formula for the same angle.
	height := 9.0
	got := computePolycone(4, 6, height)
	wantTan := height / got.Base.Apothem
	gotTan := math.Tan(got.Facet.FaceAngleDeg * math.Pi / 180)
	checkClose(t, "tan(FaceAngle)", gotTan, wantTan)
}

func TestComputePolycone_TotalAreaIsBasePlusAllFacets(t *testing.T) {
	got := computePolycone(5, 4, 6)
	want := got.Base.Area + 5*got.Facet.Area
	checkClose(t, "TotalArea", got.TotalArea, want)
}

func checkClose(t *testing.T, name string, got, want float64) {
	t.Helper()
	if math.Abs(got-want) > 1e-6 {
		t.Errorf("%s = %v, want %v", name, got, want)
	}
}
