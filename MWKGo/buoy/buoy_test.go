package main

import (
	"errors"
	"math"
	"testing"
)

func TestSphereSubmergedVolume_FullSubmersionEqualsSphereVolume(t *testing.T) {
	// A sphere submerged to its own full diameter must displace
	// exactly its own total volume (pi*d^3/6), an identity
	// independent of the partial-submersion formula's own
	// derivation.
	d := 10.0
	got := sphereSubmergedVolume(d, d)
	want := math.Pi * d * d * d / 6
	if math.Abs(got-want) > 1e-9 {
		t.Errorf("sphereSubmergedVolume(%v, %v) = %v, want %v", d, d, got, want)
	}
}

func TestHorizontalCylinderSegmentArea_FullSubmersionEqualsCircleArea(t *testing.T) {
	// A circular cross-section submerged to its own full diameter
	// must equal the full circle's area (pi*r^2).
	d := 10.0
	got := horizontalCylinderSegmentArea(d, d)
	r := 0.5 * d
	want := math.Pi * r * r
	if math.Abs(got-want) > 1e-9 {
		t.Errorf("horizontalCylinderSegmentArea(%v, %v) = %v, want %v", d, d, got, want)
	}
}

func TestSphereImmersionDepth_ForceAtDepthMatchesWeight(t *testing.T) {
	// The whole point of the binary search: the buoyant force at the
	// returned depth (density times displaced volume) must equal the
	// target weight, within the search's own tolerance. This is
	// checked directly against Archimedes' principle rather than a
	// hand-picked expected depth.
	diameter, liquidDensity, weight := 10.0, 0.0361, 3.5
	depth, err := sphereImmersionDepth(diameter, liquidDensity, weight)
	if err != nil {
		t.Fatalf("sphereImmersionDepth() error = %v", err)
	}
	force := liquidDensity * sphereSubmergedVolume(diameter, depth)
	if math.Abs(force-weight) > 0.01 {
		t.Errorf("force at depth %v = %v, want approximately %v", depth, force, weight)
	}
}

func TestHorizontalCylinderImmersionDepth_ForceAtDepthMatchesWeight(t *testing.T) {
	diameter, length, liquidDensity, weight := 5.0, 10.0, 0.0361, 3.5
	depth, err := horizontalCylinderImmersionDepth(diameter, length, liquidDensity, weight)
	if err != nil {
		t.Fatalf("horizontalCylinderImmersionDepth() error = %v", err)
	}
	force := liquidDensity * horizontalCylinderSegmentArea(diameter, depth) * length
	if math.Abs(force-weight) > 0.1 {
		t.Errorf("force at depth %v = %v, want approximately %v", depth, force, weight)
	}
}

func TestVerticalCylinderImmersionDepth_ForceAtDepthMatchesWeightExactly(t *testing.T) {
	// A vertical cylinder's cross section doesn't change with depth,
	// so this is solved directly, not by search: the resulting force
	// should match the target weight exactly (to floating point
	// precision), not just within a search tolerance.
	diameter, length, liquidDensity, weight := 5.0, 10.0, 0.0361, 3.5
	depth, err := verticalCylinderImmersionDepth(diameter, length, liquidDensity, weight)
	if err != nil {
		t.Fatalf("verticalCylinderImmersionDepth() error = %v", err)
	}
	force := liquidDensity * (0.25 * math.Pi * diameter * diameter) * depth
	if math.Abs(force-weight) > 1e-9 {
		t.Errorf("force at depth %v = %v, want %v", depth, force, weight)
	}
}

func TestBoxImmersionDepth_ForceAtDepthMatchesWeightExactly(t *testing.T) {
	length, width, height, liquidDensity, weight := 5.0, 8.0, 6.0, 0.0361, 3.5
	depth, err := boxImmersionDepth(length, width, height, liquidDensity, weight)
	if err != nil {
		t.Fatalf("boxImmersionDepth() error = %v", err)
	}
	force := liquidDensity * length * width * depth
	if math.Abs(force-weight) > 1e-9 {
		t.Errorf("force at depth %v = %v, want %v", depth, force, weight)
	}
}

func TestSphereImmersionDepth_TooHeavySinks(t *testing.T) {
	_, err := sphereImmersionDepth(1.0, 0.0361, 1000.0)
	var sink sinkError
	if !errors.As(err, &sink) {
		t.Fatalf("sphereImmersionDepth() error = %v, want a sinkError", err)
	}
	if sink.MaxSupportableWeight <= 0 {
		t.Errorf("MaxSupportableWeight = %v, want positive", sink.MaxSupportableWeight)
	}
}
