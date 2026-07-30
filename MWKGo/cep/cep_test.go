package main

import (
	"math"
	"testing"
)

func TestComputeCEP_CircularCaseMatchesKnownConstant(t *testing.T) {
	// For two equal, uncorrelated Gaussian error sources (a circular
	// bivariate normal distribution), the exact 50% CEP radius is a
	// well known closed-form constant: sigma*sqrt(2*ln(2)) =~
	// 1.1774*sigma, independent of both this program's numerical
	// integration and its Wilcox cross-check formula.
	got, err := computeCEP(0, 1, 0, 1, 0, 0.5)
	if err != nil {
		t.Fatalf("computeCEP() error = %v", err)
	}
	want := math.Sqrt(2 * math.Log(2))
	if math.Abs(got.Radius-want) > 0.01 {
		t.Errorf("Radius = %v, want approximately %v", got.Radius, want)
	}
}

func TestComputeCEP_NumericalIntegrationAgreesWithWilcoxCheck(t *testing.T) {
	// The program's own Wilcox analytic approximation, computed only
	// for zero-mean/50%-probability inputs, should closely agree
	// with the independently derived numerical integration result:
	// two different methods converging on the same answer.
	got, err := computeCEP(0, 2, 0, 3, 0.3, 0.5)
	if err != nil {
		t.Fatalf("computeCEP() error = %v", err)
	}
	if got.WilcoxChk == 0 {
		t.Fatal("WilcoxChk = 0, want the analytic cross-check to be computed for this input")
	}
	if diff := math.Abs(got.Radius - got.WilcoxChk); diff/got.Radius > 0.01 {
		t.Errorf("Radius = %v, WilcoxChk = %v, differ by more than 1%%", got.Radius, got.WilcoxChk)
	}
}

func TestComputeCEP_WilcoxCheckOnlyForZeroMeanHalfProbability(t *testing.T) {
	got, err := computeCEP(1, 1, 0, 1, 0, 0.5)
	if err != nil {
		t.Fatalf("computeCEP() error = %v", err)
	}
	if got.WilcoxChk != 0 {
		t.Errorf("WilcoxChk = %v, want 0 (Wilcox check does not apply to a nonzero mean)", got.WilcoxChk)
	}
}

func TestComputeCEP_LargerProbabilityGivesLargerRadius(t *testing.T) {
	small, err := computeCEP(0, 1, 0, 1, 0, 0.5)
	if err != nil {
		t.Fatalf("computeCEP(0.5) error = %v", err)
	}
	large, err := computeCEP(0, 1, 0, 1, 0, 0.9)
	if err != nil {
		t.Fatalf("computeCEP(0.9) error = %v", err)
	}
	if !(large.Radius > small.Radius) {
		t.Errorf("90%% CEP radius (%v) should exceed 50%% CEP radius (%v)", large.Radius, small.Radius)
	}
}

func TestComputeCEP_ZeroPrincipalScaleReturnsNegativeOne(t *testing.T) {
	// A degenerate distribution (both principal variances collapsing
	// to zero, i.e. sigma1 == sigma2 == 0) has no meaningful CEP;
	// the original program signals this with -1 rather than dividing
	// by zero.
	got, err := computeCEP(0, 0, 0, 0, 0, 0.5)
	if err != nil {
		t.Fatalf("computeCEP() error = %v", err)
	}
	if got.Radius != -1 {
		t.Errorf("Radius = %v, want -1 for a degenerate distribution", got.Radius)
	}
}
