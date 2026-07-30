package main

import (
	"math"
	"testing"
)

func TestGearpa_WorkedExample_30Teeth(t *testing.T) {
	// GEARPA.TXT's first worked example: 30 teeth, DP 6.
	teeth145, ok := spanTeethCount(30, 0)
	if !ok || teeth145 != 3 {
		t.Fatalf("spanTeethCount(30, 14.5deg) = %d, %v, want 3, true", teeth145, ok)
	}
	span145 := chordalSpan(30, 6, 14.5, teeth145)
	checkClose(t, "span145", span145, 1.2941)

	teeth20, ok := spanTeethCount(30, 1)
	if !ok || teeth20 != 4 {
		t.Fatalf("spanTeethCount(30, 20deg) = %d, %v, want 4, true", teeth20, ok)
	}
	span20 := chordalSpan(30, 6, 20, teeth20)
	checkClose(t, "span20", span20, 1.7921)
}

func TestGearpa_WorkedExample_18Teeth(t *testing.T) {
	// GEARPA.TXT's second worked example: 18 teeth, DP 6. Both
	// pressure angles happen to span the same tooth count here,
	// which is exactly why the TXT uses it to illustrate how close
	// the two measurements can be.
	teeth145, ok := spanTeethCount(18, 0)
	if !ok || teeth145 != 2 {
		t.Fatalf("spanTeethCount(18, 14.5deg) = %d, %v, want 2, true", teeth145, ok)
	}
	span145 := chordalSpan(18, 6, 14.5, teeth145)
	checkClose(t, "span145", span145, 0.7765)

	teeth20, ok := spanTeethCount(18, 1)
	if !ok || teeth20 != 2 {
		t.Fatalf("spanTeethCount(18, 20deg) = %d, %v, want 2, true", teeth20, ok)
	}
	span20 := chordalSpan(18, 6, 20, teeth20)
	checkClose(t, "span20", span20, 0.7800)
}

func TestSpanTeethCount_145DegreeOverlapQuirk(t *testing.T) {
	// The original program's overlapping range checks give teeth=5
	// only for n=51-52 at 14.5 degrees; n=53 onward, though also
	// matching the "51-62" condition, is overwritten by the
	// following "53-75" condition to teeth=6.
	if teeth, ok := spanTeethCount(52, 0); !ok || teeth != 5 {
		t.Errorf("spanTeethCount(52, 14.5deg) = %d, %v, want 5, true", teeth, ok)
	}
	if teeth, ok := spanTeethCount(53, 0); !ok || teeth != 6 {
		t.Errorf("spanTeethCount(53, 14.5deg) = %d, %v, want 6, true", teeth, ok)
	}
	if teeth, ok := spanTeethCount(62, 0); !ok || teeth != 6 {
		t.Errorf("spanTeethCount(62, 14.5deg) = %d, %v, want 6, true", teeth, ok)
	}
}

func TestSpanTeethCount_20DegreeOutOfRangeAbove81(t *testing.T) {
	if _, ok := spanTeethCount(90, 1); ok {
		t.Error("spanTeethCount(90, 20deg) ok = true, want false (90 exceeds the 20 degree table)")
	}
	if _, ok := spanTeethCount(90, 0); !ok {
		t.Error("spanTeethCount(90, 14.5deg) ok = false, want true (the 14.5 degree table covers the full 12-110 range)")
	}
}

func checkClose(t *testing.T, name string, got, want float64) {
	t.Helper()
	if math.Abs(got-want) > 1e-4 {
		t.Errorf("%s = %v, want %v", name, got, want)
	}
}
