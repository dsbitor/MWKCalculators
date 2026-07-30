package main

import (
	"math"
	"testing"
)

func TestBucketVolume_EqualRadiiMatchesCylinderVolume(t *testing.T) {
	// A frustum with equal top and bottom radii is just a cylinder,
	// whose volume is pi*height*radius^2: a standard identity
	// independent of this code.
	got := bucketVolume(10, 2, 2)
	want := math.Pi * 10 * 2 * 2
	if diff := math.Abs(got - want); diff > 1e-9 {
		t.Errorf("bucketVolume(10, 2, 2) = %v, want %v", got, want)
	}
}

func TestComputeBucketGeometry_DocumentedDefaultInput(t *testing.T) {
	// The documented default input, evaluated against the ported
	// formula directly.
	bigDiameter, smallDiameter, slantHeight := 4.0, 3.0, 6.0
	smallRadius, largeRadius := 1.5, 2.0
	wantHalfAngle := math.Asin((largeRadius-smallRadius)/slantHeight) * 180 / math.Pi
	wantHeight := slantHeight * math.Cos(wantHalfAngle*math.Pi/180)
	wantVolume := math.Pi * wantHeight * (largeRadius*largeRadius + largeRadius*smallRadius + smallRadius*smallRadius) / 3

	got := computeBucketGeometry(bigDiameter, smallDiameter, slantHeight)
	if diff := math.Abs(got.Height - wantHeight); diff > 1e-9 {
		t.Errorf("Height = %v, want %v", got.Height, wantHeight)
	}
	if diff := math.Abs(got.Volume - wantVolume); diff > 1e-9 {
		t.Errorf("Volume = %v, want %v", got.Volume, wantVolume)
	}
}

func TestSlantHeightForVolume_FullVolumeReturnsFullSlantHeight(t *testing.T) {
	// Searching for the height containing the entire volume must
	// converge back to (approximately) the bucket's own full slant
	// height: a self-consistency identity independent of the search
	// algorithm's internal step size.
	geometry := computeBucketGeometry(4, 3, 6)

	got, err := slantHeightForVolume(geometry, geometry.Volume, geometry.Height)
	if err != nil {
		t.Fatalf("slantHeightForVolume() error = %v", err)
	}
	if diff := math.Abs(got - 6.0); diff > 0.01 {
		t.Errorf("slantHeightForVolume() = %v, want approximately 6.0 (the full slant height)", got)
	}
}

func TestSlantHeightForVolume_HalfCountUpAndCountDownAgree(t *testing.T) {
	// Searching from below (increasing height) and from above
	// (decreasing height) toward the same target volume must land on
	// the same slant height, regardless of which direction the
	// search approaches from.
	geometry := computeBucketGeometry(4, 3, 6)
	target := 0.5 * geometry.Volume

	fromBelow, err := slantHeightForVolume(geometry, target, 0.1*geometry.Height)
	if err != nil {
		t.Fatalf("slantHeightForVolume() from below error = %v", err)
	}
	fromAbove, err := slantHeightForVolume(geometry, target, 0.9*geometry.Height)
	if err != nil {
		t.Fatalf("slantHeightForVolume() from above error = %v", err)
	}
	if diff := math.Abs(fromBelow - fromAbove); diff > 0.01 {
		t.Errorf("search from below (%v) and from above (%v) disagree", fromBelow, fromAbove)
	}
}

func TestSlantHeightForVolume_UnreachableTargetReturnsErrorNotPanic(t *testing.T) {
	// A target volume the search can never satisfy within the
	// bounded iteration count must return an error, matching the
	// original program's own "iteration error" safety valve, rather
	// than looping forever or panicking.
	geometry := computeBucketGeometry(4, 3, 6)

	_, err := slantHeightForVolume(geometry, -1000*geometry.Volume, geometry.Height)
	if err == nil {
		t.Fatalf("slantHeightForVolume() error = nil, want an error for a target volume the search cannot reach")
	}
}
