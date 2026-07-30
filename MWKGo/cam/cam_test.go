package main

import (
	"math"
	"testing"
)

const eps = 1e-6

func almostEqual(a, b, tol float64) bool { return math.Abs(a-b) < tol }

// cycloidalDefault matches CAM.TXT's own documented default example:
// a cycloidal cam, 1.625 base circle, 1.25 rise over 120 degrees.
func cycloidalDefault() motionLaw {
	rise, beta := 1.25, 120.0
	return motionLaw{Type: cycloidal, Rise: rise, Beta: beta, Rob: rise / (beta * rpd)}
}

func TestEvaluateMotion_AllTypesStartAtZeroDisplacement(t *testing.T) {
	laws := []motionLaw{
		cycloidalDefault(),
		{Type: parabolic, Rise: 1.25, Beta: 120, Rob: 1.25 / (120 * rpd)},
		{Type: simpleHarmonic, Rise: 1.25, Beta: 120, Rob: 1.25 / (120 * rpd)},
		{Type: straightLine, Rise: 1.35, Beta: 128, RiseA: 0.1, RiseL: 1.0, RiseD: 0.25, BetaA: 20, BetaL: 80, BetaD: 28},
	}
	for _, m := range laws {
		s := evaluateMotion(m, 0)
		if !almostEqual(s.D, 0, eps) {
			t.Errorf("type %d: displacement at theta=0 = %v, want 0", m.Type, s.D)
		}
	}
}

func TestEvaluateMotion_AllTypesEndAtFullRise(t *testing.T) {
	laws := []motionLaw{
		cycloidalDefault(),
		{Type: parabolic, Rise: 1.25, Beta: 120, Rob: 1.25 / (120 * rpd)},
		{Type: simpleHarmonic, Rise: 1.25, Beta: 120, Rob: 1.25 / (120 * rpd)},
	}
	for _, m := range laws {
		s := evaluateMotion(m, m.Beta)
		if !almostEqual(s.D, m.Rise, eps) {
			t.Errorf("type %d: displacement at theta=beta = %v, want rise %v", m.Type, s.D, m.Rise)
		}
	}
}

func TestEvaluateMotion_StraightLineEndsAtFullRise(t *testing.T) {
	m := motionLaw{Type: straightLine, RiseA: 0.1, RiseL: 1.0, RiseD: 0.25, BetaA: 20, BetaL: 80, BetaD: 28}
	m.Beta = m.BetaA + m.BetaL + m.BetaD
	m.Rise = m.RiseA + m.RiseL + m.RiseD
	s := evaluateMotion(m, m.Beta)
	if !almostEqual(s.D, m.Rise, eps) {
		t.Errorf("straight-line displacement at theta=beta = %v, want rise %v", s.D, m.Rise)
	}
}

// TestEvaluateMotion_ParabolicIsContinuousAtMidpoint confirms the
// piecewise parabolic formula's two halves agree exactly at tob=0.5,
// where CAM.C itself switches from one branch to the other.
func TestEvaluateMotion_ParabolicIsContinuousAtMidpoint(t *testing.T) {
	m := motionLaw{Type: parabolic, Rise: 1.25, Beta: 120, Rob: 1.25 / (120 * rpd)}
	theta := m.Beta * 0.5
	s := evaluateMotion(m, theta)
	wantD := 0.5 * m.Rise // at the midpoint, half the rise is complete
	if !almostEqual(s.D, wantD, eps) {
		t.Errorf("parabolic displacement at midpoint = %v, want %v (half the rise)", s.D, wantD)
	}
}

// TestEvaluateMotion_CycloidalMidpointIsHalfRise checks a
// hand-derivable value: cycloidal displacement at tob=0.5 is exactly
// half the rise, since sin(2*pi*0.5) = sin(pi) = 0.
func TestEvaluateMotion_CycloidalMidpointIsHalfRise(t *testing.T) {
	m := cycloidalDefault()
	s := evaluateMotion(m, m.Beta*0.5)
	if !almostEqual(s.D, 0.5*m.Rise, 1e-9) {
		t.Errorf("cycloidal displacement at midpoint = %v, want %v", s.D, 0.5*m.Rise)
	}
}

// TestMotionMax_CycloidalPeakVelocityMatchesFormula confirms the
// cycloidal motion's own peak velocity (which occurs exactly at the
// midpoint) matches motionMax's independently-stated formula.
func TestMotionMax_CycloidalPeakVelocityMatchesFormula(t *testing.T) {
	m := cycloidalDefault()
	vmax, _ := motionMax(m)
	s := evaluateMotion(m, m.Beta*0.5)
	if !almostEqual(s.Vc, vmax, 1e-9) {
		t.Errorf("cycloidal velocity at midpoint = %v, want vmax = %v", s.Vc, vmax)
	}
}

func TestComputeCamProfile_AtStartMatchesBaseCirclePlusRise(t *testing.T) {
	m := cycloidalDefault()
	rbase := 1.625
	points := computeCamProfile(m, rbase, 0, 1.0)
	// The first point (ang=beta, theta=0) should sit exactly on the
	// base circle radius, since displacement is 0 there.
	first := points[0]
	if !almostEqual(first.Radius, rbase, 1e-6) {
		t.Errorf("cam radius at start of rise = %v, want base circle radius %v", first.Radius, rbase)
	}
	// The last point (ang=0, theta=beta) should sit at base+rise.
	last := points[len(points)-1]
	if !almostEqual(last.Radius, rbase+m.Rise, 1e-6) {
		t.Errorf("cam radius at end of rise = %v, want %v", last.Radius, rbase+m.Rise)
	}
}

func TestFindMaxPressureAngle_IsWithinProfileRange(t *testing.T) {
	m := cycloidalDefault()
	rbase := 1.625
	maxP := findMaxPressureAngle(m, rbase, 0, 1.0)
	if maxP.AngleRad <= 0 || maxP.AngleRad >= math.Pi/2 {
		t.Errorf("max pressure angle = %v rad, want a value strictly between 0 and pi/2", maxP.AngleRad)
	}
	// CAM.TXT recommends staying under 35 degrees; the default
	// example should be a reasonable, realistic cam design.
	if deg := maxP.AngleRad * 180 / math.Pi; deg > 60 {
		t.Errorf("max pressure angle = %v deg, suspiciously large for the documented default example", deg)
	}
}

func TestSuggestBaseRadii_LargerPressureAngleAllowsSmallerBase(t *testing.T) {
	m := cycloidalDefault()
	maxP := findMaxPressureAngle(m, 1.625, 0, 1.0)
	suggestions := suggestBaseRadii(maxP, 0)
	if len(suggestions) == 0 {
		t.Fatal("expected at least one suggestion")
	}
	for i := 1; i < len(suggestions); i++ {
		if suggestions[i].BaseRadius >= suggestions[i-1].BaseRadius {
			t.Errorf("suggestion %d base radius (%v) should be smaller than suggestion %d's (%v): a larger allowed pressure angle permits a smaller base circle",
				i, suggestions[i].BaseRadius, i-1, suggestions[i-1].BaseRadius)
		}
	}
}
