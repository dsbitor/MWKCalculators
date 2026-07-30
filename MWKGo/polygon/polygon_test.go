package main

import (
	"math"
	"testing"
)

func TestComputeRegularPolygon_HexagonSideEqualsCircumradius(t *testing.T) {
	// A well known property of a regular hexagon: its side length
	// equals its circumradius exactly, independent of this code's
	// own trigonometric formula for side length.
	got := computeRegularPolygon(6, 1.0)
	if math.Abs(got.SideLength-1.0) > 1e-9 {
		t.Errorf("SideLength = %v, want 1.0 (equal to circumradius for a hexagon)", got.SideLength)
	}
}

func TestComputeRegularPolygon_SquareAcrossFlatsEqualsSideLength(t *testing.T) {
	// A well known property of a square: the distance across flats
	// (its own width) equals its side length exactly.
	circumradius := circumradiusFromSideLength(4, 1.0)
	got := computeRegularPolygon(4, circumradius)
	if math.Abs(got.SideLength-1.0) > 1e-9 {
		t.Errorf("SideLength = %v, want 1.0", got.SideLength)
	}
	if math.Abs(got.AcrossFlats-1.0) > 1e-9 {
		t.Errorf("AcrossFlats = %v, want 1.0 (equal to side length for a square)", got.AcrossFlats)
	}
}

func TestCircumradiusInversesRoundTrip(t *testing.T) {
	// Each circumradiusFromX helper must exactly invert
	// computeRegularPolygon's own corresponding output field,
	// regardless of which one the user happened to supply.
	sideCount := 7
	wantCircumradius := 3.3

	p := computeRegularPolygon(sideCount, wantCircumradius)

	if got := circumradiusFromSideLength(sideCount, p.SideLength); math.Abs(got-wantCircumradius) > 1e-9 {
		t.Errorf("circumradiusFromSideLength = %v, want %v", got, wantCircumradius)
	}
	if got := circumradiusFromFlatToOppositeVertex(sideCount, p.FlatToOppositeVertex); math.Abs(got-wantCircumradius) > 1e-9 {
		t.Errorf("circumradiusFromFlatToOppositeVertex = %v, want %v", got, wantCircumradius)
	}
	if got := circumradiusFromInscribedDiameter(sideCount, p.InscribedDiam); math.Abs(got-wantCircumradius) > 1e-9 {
		t.Errorf("circumradiusFromInscribedDiameter = %v, want %v", got, wantCircumradius)
	}

	evenSides := 8
	pEven := computeRegularPolygon(evenSides, wantCircumradius)
	if got := circumradiusFromAcrossFlats(evenSides, pEven.AcrossFlats); math.Abs(got-wantCircumradius) > 1e-9 {
		t.Errorf("circumradiusFromAcrossFlats = %v, want %v", got, wantCircumradius)
	}
}
