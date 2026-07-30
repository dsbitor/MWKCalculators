package main

import (
	"math"
	"testing"
)

func TestComputeKeyGeometry_DocumentedDefaultInput(t *testing.T) {
	got := computeKeyGeometry(1.0, 0.5, 0.5)
	x := math.Sqrt(1.0*1.0 - 0.5*0.5)
	y := 0.5 * (1.0 - x)
	wantM := 1.0 - (y + 0.25)
	wantN := wantM + 0.5
	wantQ := 1.0 - wantM
	wantR := 0.5 - wantQ

	checkClose(t, "KeyseatDepth", got.KeyseatDepth, wantM)
	checkClose(t, "KeywayDepth", got.KeywayDepth, wantN)
	checkClose(t, "ShaftCutDepth", got.ShaftCutDepth, wantQ)
	checkClose(t, "HubCutDepth", got.HubCutDepth, wantR)
}

func TestKeyOutline_RotationPreservesDistanceFromOrigin(t *testing.T) {
	// Rotating a point about the origin must preserve its distance
	// from the origin exactly: an identity independent of this
	// program's own rotation formula, checked for every point the
	// original program actually rotates (indices 2 through 6).
	shaftDiameter, keyWidth := 1.0, 0.5
	g := computeKeyGeometry(shaftDiameter, keyWidth, 0.5)
	before, after := keyOutline(shaftDiameter, keyWidth, g.ShaftCutDepth, g.HubCutDepth, 15.0)

	for i := 2; i <= 6; i++ {
		beforeDist := math.Hypot(before[i].X, before[i].Y)
		afterDist := math.Hypot(after[i].X, after[i].Y)
		if math.Abs(beforeDist-afterDist) > 1e-9 {
			t.Errorf("point %d: distance from origin changed from %v to %v after rotation", i, beforeDist, afterDist)
		}
	}
}

func TestKeyOutline_UnrotatedIndicesAreUnchanged(t *testing.T) {
	// Points 0, 1, and 7 are explicitly excluded from rotation in
	// the original program (they represent the shaft cut, which
	// doesn't move with the key).
	shaftDiameter, keyWidth := 1.0, 0.5
	g := computeKeyGeometry(shaftDiameter, keyWidth, 0.5)
	before, after := keyOutline(shaftDiameter, keyWidth, g.ShaftCutDepth, g.HubCutDepth, 15.0)

	for _, i := range []int{0, 1, 7} {
		if before[i] != after[i] {
			t.Errorf("point %d: before %v != after %v, want unchanged", i, before[i], after[i])
		}
	}
}

func TestKeyOutline_ZeroRotationLeavesEveryPointUnchanged(t *testing.T) {
	shaftDiameter, keyWidth := 1.0, 0.5
	g := computeKeyGeometry(shaftDiameter, keyWidth, 0.5)
	before, after := keyOutline(shaftDiameter, keyWidth, g.ShaftCutDepth, g.HubCutDepth, 0.0)

	for i := 0; i < 8; i++ {
		if math.Abs(before[i].X-after[i].X) > 1e-9 || math.Abs(before[i].Y-after[i].Y) > 1e-9 {
			t.Errorf("point %d: before %v, after %v, want equal at zero rotation", i, before[i], after[i])
		}
	}
}

func TestLineIntersection_KnownCrossingPoint(t *testing.T) {
	// Two lines crossing at a known point: y=x and y=-x+4 cross at
	// (2,2), independent of this code's own determinant-based
	// formula.
	got := lineIntersection(0, 0, 1, 1, 0, 4, 4, 0)
	checkClose(t, "X", got.X, 2)
	checkClose(t, "Y", got.Y, 2)
}

func checkClose(t *testing.T, name string, got, want float64) {
	t.Helper()
	if math.Abs(got-want) > 1e-9 {
		t.Errorf("%s = %v, want %v", name, got, want)
	}
}
