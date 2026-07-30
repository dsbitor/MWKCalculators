package main

import (
	"math"
	"testing"
)

const eps = 1e-9

func almostEqual(a, b, tol float64) bool { return math.Abs(a-b) < tol }

// TestEggSemiAxes_MatchEggTXTDocumentedValues confirms the hardwired
// values EGG.TXT itself documents: "stock diameter = 1", b = 0.5",
// a = 0.75", k = 0.6".
func TestEggSemiAxes_MatchEggTXTDocumentedValues(t *testing.T) {
	rs := 0.5 * stockDiam
	a, b, k := eggSemiAxes(rs)
	if !almostEqual(rs, 0.5, eps) {
		t.Errorf("stock radius = %v, want 0.5", rs)
	}
	if !almostEqual(a, 0.75, eps) {
		t.Errorf("a = %v, want 0.75", a)
	}
	if !almostEqual(b, 0.5, eps) {
		t.Errorf("b = %v, want 0.5", b)
	}
	if !almostEqual(k, 0.6, eps) {
		t.Errorf("k = %v, want 0.6", k)
	}
}

// TestFindRadius_TipIsZero confirms the pointed end of the egg (x=0,
// the outboard end of the stock) has zero radius.
func TestFindRadius_TipIsZero(t *testing.T) {
	rs := 0.5 * stockDiam
	got := findRadius(0, rs)
	if !almostEqual(got, 0, 1e-9) {
		t.Errorf("findRadius(0) = %v, want 0 (the pointed end)", got)
	}
}

// TestFindRadius_BeyondEggFormulaIsSmallNub confirms the profile
// flattens to a small nub (0.1*stockRadius) once x reaches twice the
// semi-major axis, rather than evaluating the egg formula where it
// would become mathematically undefined.
func TestFindRadius_BeyondEggFormulaIsSmallNub(t *testing.T) {
	rs := 0.5 * stockDiam
	a, _, _ := eggSemiAxes(rs)
	got := findRadius(2*a, rs)
	want := 0.1 * rs
	if !almostEqual(got, want, eps) {
		t.Errorf("findRadius(2a) = %v, want %v (the small end nub)", got, want)
	}
}

// TestFindRadius_AtSemiMajorAxisEqualsSemiMinorAxis confirms a
// hand-derivable point on the curve: at x=a (the egg formula's own
// "a-x=0" case), the radius is exactly b, regardless of asymmetry k
// (since the (1-k*x) denominator term vanishes to 1 there). Note this
// is not actually the profile's single widest point — the asymmetric
// formula bulges slightly past b a bit before x=a, which is the
// entire point of the "asymmetry factor" EGG.TXT describes; it is
// only a well known, algebraically exact point on the curve.
func TestFindRadius_AtSemiMajorAxisEqualsSemiMinorAxis(t *testing.T) {
	rs := 0.5 * stockDiam
	a, b, _ := eggSemiAxes(rs)
	got := findRadius(a, rs)
	if !almostEqual(got, b, 1e-9) {
		t.Errorf("findRadius(a) = %v, want b = %v", got, b)
	}
}

func TestFindRadius_StaysWithinReasonableBounds(t *testing.T) {
	rs := 0.5 * stockDiam
	dim := computeDimensions()
	for x := 0.0; x <= dim.PlotEnd; x += 0.01 {
		r := findRadius(x, rs)
		if r < 0 || r > 1.5*rs {
			t.Errorf("findRadius(%v) = %v, want within [0, 1.5*stockRadius]", x, r)
		}
	}
}

// TestEggShape_SymmetricCaseMatchesPlainEllipse confirms eggShape
// reduces to a standard (symmetric) ellipse formula when k=0 — a
// sanity check on the formula independent of the egg-specific
// asymmetry feature.
func TestEggShape_SymmetricCaseMatchesPlainEllipse(t *testing.T) {
	a, b := 2.0, 1.0
	x := 1.0
	got := eggShape(x, a, b, 0)
	want := b * math.Sqrt(1-x*x/(a*a))
	if !almostEqual(got, want, eps) {
		t.Errorf("eggShape(x,a,b,0) = %v, want %v (plain ellipse)", got, want)
	}
}

func TestComputeCuttingSchedule_AllDepthsPositive(t *testing.T) {
	dim := computeDimensions()
	steps := computeCuttingSchedule(dim)
	if len(steps) == 0 {
		t.Fatal("expected at least one cutting step")
	}
	for _, s := range steps {
		if s.DepthOfCut <= 0 {
			t.Errorf("step %d depth of cut = %v, want positive", s.Index, s.DepthOfCut)
		}
		if s.Diameter < 0 || s.Diameter >= 2*dim.StockRadius {
			t.Errorf("step %d diameter = %v, want within [0, stock diameter)", s.Index, s.Diameter)
		}
	}
}

func TestComputeCuttingSchedule_FirstStepIsNearTheTip(t *testing.T) {
	dim := computeDimensions()
	steps := computeCuttingSchedule(dim)
	if steps[0].Index != 1 {
		t.Errorf("first step index = %d, want 1", steps[0].Index)
	}
	if steps[0].GapBefore {
		t.Error("first step should never report a gap before it")
	}
}
