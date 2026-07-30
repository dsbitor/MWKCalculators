package main

import (
	"math"
	"testing"
)

const eps = 1e-9

func almostEqual(a, b, tol float64) bool { return math.Abs(a-b) < tol }

// TestGudgeonPinPosition_TopDeadCenter confirms the well known
// physical fact that at TDC (theta=0), the piston is at its maximum
// extension: pin position relative to the crank center equals
// throw+rodLen exactly (the crank and rod are fully "stretched out"
// in a straight line).
func TestGudgeonPinPosition_TopDeadCenter(t *testing.T) {
	throw, rodLen := 0.6, 2.4
	got := gudgeonPinPosition(throw, rodLen, 0)
	if !almostEqual(got, throw+rodLen, eps) {
		t.Errorf("gudgeonPinPosition at TDC = %v, want %v", got, throw+rodLen)
	}
}

// TestGudgeonPinPosition_BottomDeadCenter confirms the complementary
// fact at BDC (theta=180): minimum extension, rodLen-throw.
func TestGudgeonPinPosition_BottomDeadCenter(t *testing.T) {
	throw, rodLen := 0.6, 2.4
	got := gudgeonPinPosition(throw, rodLen, 180)
	if !almostEqual(got, rodLen-throw, eps) {
		t.Errorf("gudgeonPinPosition at BDC = %v, want %v", got, rodLen-throw)
	}
}

func TestGudgeonPinPosition_StaysWithinThrowBounds(t *testing.T) {
	throw, rodLen := 0.6, 2.4
	for theta := 0.0; theta <= 180; theta += 5 {
		x := gudgeonPinPosition(throw, rodLen, theta)
		if x < rodLen-throw-eps || x > rodLen+throw+eps {
			t.Errorf("gudgeonPinPosition(%v) = %v, want within [%v,%v]", theta, x, rodLen-throw, rodLen+throw)
		}
	}
}

func TestComputeCrodTable_TDCAndBDCValues(t *testing.T) {
	throw, rodLen, dtheta := 0.6, 2.4, 5.0
	points := computeCrodTable(throw, rodLen, dtheta)

	first := points[0]
	if first.ThetaDeg != 0 || !almostEqual(first.XCenterRelative, 3.0, eps) || !almostEqual(first.ZTDCRelative, 0, eps) {
		t.Errorf("first point = %+v, want theta=0, X=3.0, Z=0", first)
	}

	last := points[len(points)-1]
	if !almostEqual(last.ThetaDeg, 180, eps) {
		t.Fatalf("last point theta = %v, want 180", last.ThetaDeg)
	}
	if !almostEqual(last.XCenterRelative, 1.8, eps) {
		t.Errorf("last point X = %v, want 1.8 (rodLen-throw)", last.XCenterRelative)
	}
	if !almostEqual(last.ZTDCRelative, 1.2, eps) {
		t.Errorf("last point Z = %v, want 1.2 (2*throw, full travel from TDC)", last.ZTDCRelative)
	}
}

func TestComputeCrodTable_ZIsAlwaysXReflectedFromTDC(t *testing.T) {
	throw, rodLen, dtheta := 0.6, 2.4, 10.0
	for _, p := range computeCrodTable(throw, rodLen, dtheta) {
		want := throw + rodLen - p.XCenterRelative
		if !almostEqual(p.ZTDCRelative, want, eps) {
			t.Errorf("theta=%v: Z = %v, want %v (throw+rodLen-X)", p.ThetaDeg, p.ZTDCRelative, want)
		}
	}
}
