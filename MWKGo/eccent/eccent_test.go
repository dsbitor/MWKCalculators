package main

import (
	"math"
	"testing"
)

func TestPackingThickness_ThinJawBranch(t *testing.T) {
	// When the jaw is wide relative to the offset (w > sqrt3*e),
	// the packing thickness is simply 1.5 times the offset.
	got, err := packingThickness(1.0, 2.0, 0.1)
	if err != nil {
		t.Fatalf("packingThickness() error = %v", err)
	}
	want := 1.5 * 0.1
	if diff := math.Abs(got - want); diff > 1e-9 {
		t.Errorf("packingThickness(1.0, 2.0, 0.1) = %v, want %v", got, want)
	}
}

func TestPackingThickness_ThickJawBranch_DocumentedDefaultInput(t *testing.T) {
	// The documented default input falls into the thick-jaw branch
	// (w <= sqrt3*e), evaluated against the ported formula directly.
	jawWidth, workpieceDiameter, offset := 0.125, 1.5625, 0.28125
	radius := 0.5 * workpieceDiameter
	want := 1.5*offset - radius + 0.5*math.Sqrt(4*radius*radius-3*offset*offset+2*offset*jawWidth*sqrt3-jawWidth*jawWidth)

	got, err := packingThickness(jawWidth, workpieceDiameter, offset)
	if err != nil {
		t.Fatalf("packingThickness() error = %v", err)
	}
	if diff := math.Abs(got - want); diff > 1e-9 {
		t.Errorf("packingThickness(%v, %v, %v) = %v, want %v", jawWidth, workpieceDiameter, offset, got, want)
	}
}

func TestPackingThickness_OffsetTooLarge_ReturnsError(t *testing.T) {
	_, err := packingThickness(0.1, 0.2, 10.0)
	if err == nil {
		t.Fatalf("packingThickness() error = nil, want an error when the work would fall through the unpacked jaws")
	}
}

func TestPackingThickness_OffsetExactlyAtTheFallThroughBoundary(t *testing.T) {
	// At the exact boundary (offset == (radius+jawWidth)/sqrt3), the
	// original program's strict greater-than check does not treat
	// this as falling through: the calculation proceeds normally.
	jawWidth, radius := 0.1, 1.0
	workpieceDiameter := 2 * radius
	offset := (radius + jawWidth) / sqrt3

	_, err := packingThickness(jawWidth, workpieceDiameter, offset)
	if err != nil {
		t.Errorf("packingThickness() at the exact boundary returned an error, want the boundary itself to still be solvable: %v", err)
	}
}
