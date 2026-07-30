package main

import (
	"math"
	"testing"
)

func TestComputeFlutePlan_DocumentedDefaultInput(t *testing.T) {
	got := computeFlutePlan(0.25, 0.1, 0.2, 3.0)

	wantSmall := 0.25 - math.Sqrt(0.25*0.25-0.1*0.1)
	wantLarge := 0.25 - math.Sqrt(0.25*0.25-0.2*0.2)
	wantIncl := math.Atan((wantLarge-wantSmall)/3.0) * 180 / math.Pi

	if math.Abs(got.DepthAtSmallEnd-wantSmall) > 1e-9 {
		t.Errorf("DepthAtSmallEnd = %v, want %v", got.DepthAtSmallEnd, wantSmall)
	}
	if math.Abs(got.DepthAtLargeEnd-wantLarge) > 1e-9 {
		t.Errorf("DepthAtLargeEnd = %v, want %v", got.DepthAtLargeEnd, wantLarge)
	}
	if math.Abs(got.InclinationDeg-wantIncl) > 1e-9 {
		t.Errorf("InclinationDeg = %v, want %v", got.InclinationDeg, wantIncl)
	}
}

func TestComputeFlutePlan_ZeroFluteRadiusHasZeroDepth(t *testing.T) {
	// A flute radius of zero means the ball's tip barely grazes the
	// surface: an independently obvious identity, since the
	// sagitta of a zero-width chord is zero.
	got := computeFlutePlan(0.25, 0, 0.2, 3.0)
	if math.Abs(got.DepthAtSmallEnd) > 1e-9 {
		t.Errorf("DepthAtSmallEnd = %v, want 0", got.DepthAtSmallEnd)
	}
}

func TestComputeFlutePlan_FluteRadiusEqualsMillRadiusCutsToCenter(t *testing.T) {
	// A flute radius equal to the ball mill's own radius means the
	// full hemisphere is buried: depth of cut equals the mill
	// radius itself.
	millRadius := 0.25
	got := computeFlutePlan(millRadius, millRadius, millRadius, 3.0)
	if math.Abs(got.DepthAtSmallEnd-millRadius) > 1e-9 {
		t.Errorf("DepthAtSmallEnd = %v, want %v", got.DepthAtSmallEnd, millRadius)
	}
	if math.Abs(got.DepthAtLargeEnd-millRadius) > 1e-9 {
		t.Errorf("DepthAtLargeEnd = %v, want %v", got.DepthAtLargeEnd, millRadius)
	}
	if got.InclinationDeg != 0 {
		t.Errorf("InclinationDeg = %v, want 0 (equal depths at both ends)", got.InclinationDeg)
	}
}
