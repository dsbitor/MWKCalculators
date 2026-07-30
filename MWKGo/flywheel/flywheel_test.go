package main

import (
	"math"
	"testing"
)

func TestSearchPhi_SolvesTheDefiningEquation(t *testing.T) {
	// Whatever phi the search returns, it must make the original
	// transcendental equation's residual (approximately) zero: the
	// defining property of the root search, not a re-run of it.
	r1, s1, theta1 := 0.440, 0.5*0.1875, 30.0
	r2, s2, theta2 := 1.367, 0.5*0.1875, 7.0

	phi := searchPhi(r1, s1, theta1, r2, s2, theta2)
	residual := r1*sind(theta1+phi) - s1 - (r2*sind(theta2+phi) - s2)
	if math.Abs(residual) > 1e-4 {
		t.Errorf("residual at phi=%v is %v, want approximately 0", phi, residual)
	}
}

func TestComputeFlywheelLayout_DocumentedDemoInput(t *testing.T) {
	// FLYWHEEL.C's own header comment recommends "offset = 0 and
	// theta2 = 7 deg" as a demonstration input.
	exact, nearestIntegral, err := computeFlywheelLayout(6, 0.440, 1.367, 0.1875, 0.1875, 0, 7.0)
	if err != nil {
		t.Fatalf("computeFlywheelLayout() error = %v", err)
	}

	// The exact solution's phi must satisfy the same residual check
	// as above, applied end-to-end through the full computation.
	s1, s2 := 0.5*0.1875, 0.5*0.1875
	residual := 0.440*sind(exact.Theta1+exact.Phi) - s1 - (1.367*sind(exact.Theta2+exact.Phi) - s2)
	if math.Abs(residual) > 1e-4 {
		t.Errorf("exact solution residual = %v, want approximately 0", residual)
	}

	// The nearest-integral solution rounds theta2 and phi to whole
	// degrees and then re-solves for R1 to match: its own theta2 and
	// phi must therefore actually be whole numbers.
	if nearestIntegral.Theta2 != math.Trunc(nearestIntegral.Theta2) {
		t.Errorf("nearestIntegral.Theta2 = %v, want a whole number", nearestIntegral.Theta2)
	}
	if nearestIntegral.Phi != math.Trunc(nearestIntegral.Phi) {
		t.Errorf("nearestIntegral.Phi = %v, want a whole number", nearestIntegral.Phi)
	}
}

func TestComputeFlywheelLayout_NoTaperSpecifiedReturnsError(t *testing.T) {
	_, _, err := computeFlywheelLayout(6, 0.440, 1.367, 0.1875, 0.1875, 0, 0)
	if err == nil {
		t.Fatal("computeFlywheelLayout() error = nil, want an error when neither offset nor theta2 is given")
	}
}

func TestInnerHoleSettings_EvenlySpacedByTwiceTheta1(t *testing.T) {
	theta1 := 30.0
	settings := innerHoleSettings(theta1)
	if len(settings) == 0 {
		t.Fatal("len(settings) = 0, want at least one setting")
	}
	for i, s := range settings {
		want := theta1 + float64(i)*2*theta1
		if math.Abs(s-want) > 1e-9 {
			t.Errorf("settings[%d] = %v, want %v", i, s, want)
		}
		if s >= 360 {
			t.Errorf("settings[%d] = %v, want < 360", i, s)
		}
	}
}

func TestOuterHoleSettings_AlternatesStepSize(t *testing.T) {
	theta2, sep := 7.0, 46.0
	settings := outerHoleSettings(theta2, sep)
	if len(settings) < 2 {
		t.Fatalf("len(settings) = %d, want at least 2", len(settings))
	}
	// First step from theta2 should use 2*theta2 (useSep starts false).
	checkClose(t, "settings[1]", settings[1], theta2+2*theta2)
}

func checkClose(t *testing.T, name string, got, want float64) {
	t.Helper()
	if math.Abs(got-want) > 1e-9 {
		t.Errorf("%s = %v, want %v", name, got, want)
	}
}
