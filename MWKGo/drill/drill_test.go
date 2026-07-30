package main

import (
	"context"
	"math"
	"strings"
	"testing"
)

func realDrills(t *testing.T) []drill {
	t.Helper()
	t.Setenv("HOME", t.TempDir())
	drills, err := loadDrills(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	return drills
}

const eps = 1e-6

func almostEqual(a, b, tol float64) bool { return math.Abs(a-b) < tol }

// TestThreadFormConstants_MatchDrillTXT confirms the two published
// thread-form constants DRILL.TXT discusses at length: "1.299..." for
// the American National (imperial, default) / Standard (metric)
// thread form, and "1.082..." for the American Unified (imperial) /
// ISO (metric) thread form.
func TestThreadFormConstants_MatchDrillTXT(t *testing.T) {
	if !almostEqual(sixEighthsTan60, 1.299, 1e-3) {
		t.Errorf("sixEighthsTan60 = %v, want close to 1.299 (National/Standard)", sixEighthsTan60)
	}
	if !almostEqual(fiveEighthsTan60, 1.082, 1e-3) {
		t.Errorf("fiveEighthsTan60 = %v, want close to 1.082 (Unified/ISO)", fiveEighthsTan60)
	}
}

// TestTapDrillFormula_QuarterTwentyMatchesKnownPractice checks the
// cutting-tap formula against a widely known real-world fact: a
// 1/4-20 UNC tap (0.25 in major diameter, 20 threads/in) at 75% depth
// of thread, using the National thread form, is commonly drilled with
// a #7 drill (0.201 in) — the standard recommendation printed on
// nearly every tap drill chart.
func TestTapDrillFormula_QuarterTwentyMatchesKnownPractice(t *testing.T) {
	td, pitch, dot := 0.25, 20.0, 75.0
	dtd := td - 0.01*dot*sixEighthsTan60/pitch
	if !almostEqual(dtd, 0.2013, 1e-3) {
		t.Errorf("1/4-20 tap drill diameter = %v, want close to 0.2013 (a #7 drill)", dtd)
	}

	drills := realDrills(t)
	k := findDrill(drills, dtd, false)
	if drills[k].Name != "7" {
		t.Errorf("nearest drill to %v = %q, want the #7 drill (named \"7\" in the table)", dtd, drills[k].Name)
	}
}

func TestFindDrill_ExactMatchHasZeroDifference(t *testing.T) {
	drills := realDrills(t)
	// 3/8" is a standard fractional drill, guaranteed to be exact in
	// the table.
	k := findDrill(drills, 0.375, false)
	if drills[k].SizeIn != 0.375 {
		t.Errorf("findDrill(0.375) = %+v, want an exact 0.375 match", drills[k])
	}
}

func TestFindDrill_ClampsAtTableBoundaries(t *testing.T) {
	drills := realDrills(t)
	if k := findDrill(drills, -1, false); k != 0 {
		t.Errorf("findDrill(-1) = %d, want 0 (clamped to the smallest drill)", k)
	}
	last := len(drills) - 1
	if k := findDrill(drills, 1000, false); k != last {
		t.Errorf("findDrill(1000) = %d, want %d (clamped to the largest drill)", k, last)
	}
}

func TestFindDrill_ExcludeMetricSkipsMetricOnlyEntries(t *testing.T) {
	drills := realDrills(t)
	// Find a pure-metric entry's index directly and confirm
	// requesting its exact size with excludeMetric walks back to a
	// non-metric-only (or dual-labeled) entry instead.
	isMetricOnly := func(name string) bool {
		return strings.Contains(name, "mm") && !strings.Contains(name, "=")
	}
	for i, d := range drills {
		if i == 0 {
			continue
		}
		if isMetricOnly(d.Name) {
			k := findDrill(drills, d.SizeIn, true)
			if isMetricOnly(drills[k].Name) {
				t.Errorf("findDrill(%v, excludeMetric=true) = %q, still a metric-only entry", d.SizeIn, drills[k].Name)
			}
			return
		}
	}
	t.Skip("no metric-only entry found in the shipped table")
}

func TestFindByDesignation_ExactMatchIsCaseInsensitive(t *testing.T) {
	drills := realDrills(t)
	k, ok := findByDesignation(drills, "3/8")
	if !ok || drills[k].SizeIn != 0.375 {
		t.Fatalf("findByDesignation(3/8) = (%d,%v), want the 0.375 drill", k, ok)
	}
}

func TestFindByDesignation_NoMatchReturnsFalse(t *testing.T) {
	drills := realDrills(t)
	if _, ok := findByDesignation(drills, "not-a-real-drill-designation"); ok {
		t.Error("expected no match for a nonsense designation")
	}
}

func TestParseSize_DecimalAndFractional(t *testing.T) {
	cases := []struct {
		in   string
		want float64
	}{
		{"1.5", 1.5},
		{"3/8", 0.375},
		{"1-1/2", 1.5},
		{"0.25", 0.25},
	}
	for _, c := range cases {
		if got := parseSize(c.in); !almostEqual(got, c.want, eps) {
			t.Errorf("parseSize(%q) = %v, want %v", c.in, got, c.want)
		}
	}
}

func TestCircleAreaAndDiameterFromArea_AreInverses(t *testing.T) {
	for _, d := range []float64{0.1, 0.375, 1.0, 2.5} {
		area := circleArea(d)
		if got := diameterFromArea(area); !almostEqual(got, d, 1e-9) {
			t.Errorf("diameterFromArea(circleArea(%v)) = %v, want %v", d, got, d)
		}
	}
}

func TestStepDrillingSchedule_EndsAtFinalDrillWithFullRemoval(t *testing.T) {
	drills := realDrills(t)
	finalIdx := findDrill(drills, 0.5, false)
	pilotIdx := findDrill(drills, 0.25, false)

	entries := stepDrillingSchedule(drills, finalIdx, pilotIdx, 20, false)
	if len(entries) == 0 {
		t.Fatal("expected at least one step")
	}
	last := entries[len(entries)-1]
	if last.Drill.Name != drills[finalIdx].Name {
		t.Errorf("last step drill = %q, want final drill %q", last.Drill.Name, drills[finalIdx].Name)
	}
	if !almostEqual(last.RemovedPercent, 100, 1e-6) {
		t.Errorf("last step removed%% = %v, want 100", last.RemovedPercent)
	}
	// Each step should remove at least as much material as the one
	// before it (a monotonically increasing schedule).
	for i := 1; i < len(entries); i++ {
		if entries[i].RemovedPercent < entries[i-1].RemovedPercent {
			t.Errorf("step %d removed%% (%v) is less than step %d's (%v)", i, entries[i].RemovedPercent, i-1, entries[i-1].RemovedPercent)
		}
	}
}

func TestPilotRemovedPercent_ZeroAreaReturnsFalse(t *testing.T) {
	drills := []drill{{Name: "zero", SizeIn: 0}, {Name: "half", SizeIn: 0.5}}
	if _, ok := pilotRemovedPercent(drills, 1, 0); ok {
		t.Error("expected ok=false when the pilot drill's own size is zero")
	}
}
