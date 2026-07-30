package main

import (
	"math"
	"testing"
)

func TestGrooveCutsByAngle_CenterAndEdgeDepths(t *testing.T) {
	// At the groove's centerline (theta=0, x=r), the tool is only
	// as deep as its own radius: the shallowest cut. At the groove's
	// edge (theta=90, x=0), the cut reaches the full groove radius:
	// the deepest cut. Both are independently obvious from the
	// geometry, not just a re-run of the formula.
	grooveRadius, toolDiameter := 1.0, 0.25
	cuts := grooveCutsByAngle(grooveRadius, toolDiameter, 90)
	if len(cuts) != 2 {
		t.Fatalf("len(cuts) = %d, want 2", len(cuts))
	}
	checkClose(t, "first.DOC", cuts[0].DOC, 0.5*toolDiameter)
	checkClose(t, "last.DOC", cuts[1].DOC, grooveRadius)
}

func TestGrooveCutsByAngle_DepthNeverExceedsGrooveRadius(t *testing.T) {
	cuts := grooveCutsByAngle(1.0, 0.25, 5)
	for _, c := range cuts {
		if c.DOC > 1.0+1e-9 {
			t.Errorf("DOC = %v at X = %v, want <= groove radius 1.0", c.DOC, c.X)
		}
	}
}

func TestGrooveCutsByLinearStep_EndpointsMatchAngularMode(t *testing.T) {
	// Both stepping modes describe the same physical groove profile,
	// so their endpoints (edge and centerline) must agree exactly,
	// regardless of which direction each mode steps in.
	grooveRadius, toolDiameter := 1.0, 0.25
	step := effectiveRadius(grooveRadius, toolDiameter)
	cuts := grooveCutsByLinearStep(grooveRadius, toolDiameter, step)
	if len(cuts) != 2 {
		t.Fatalf("len(cuts) = %d, want 2", len(cuts))
	}
	checkClose(t, "first.DOC", cuts[0].DOC, 0.5*toolDiameter)
	checkClose(t, "last.DOC", cuts[1].DOC, grooveRadius)
}

func checkClose(t *testing.T, name string, got, want float64) {
	t.Helper()
	if math.Abs(got-want) > 1e-9 {
		t.Errorf("%s = %v, want %v", name, got, want)
	}
}
