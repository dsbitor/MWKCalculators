package main

import (
	"math"
	"testing"
)

const eps = 1e-6

func almostEqual(a, b, tol float64) bool { return math.Abs(a-b) < tol }

// TestConsecutiveRatioPercent_MatchesMenuText checks the formula
// against RATIO.C's own menu text, which states each standard
// series' approximate step percentage directly (5 numbers ~58% apart,
// 10 numbers ~26%, 20 ~12%, 40 ~6%, 80 ~3%) — a form of worked
// example baked into the program's own UI, in place of a separate
// .TXT file (none was shipped with this program).
func TestConsecutiveRatioPercent_MatchesMenuText(t *testing.T) {
	cases := []struct {
		m    int
		want float64
	}{
		{5, 58}, {10, 26}, {20, 12}, {40, 6}, {80, 3},
	}
	for _, c := range cases {
		got := consecutiveRatioPercent(c.m)
		if math.Abs(got-c.want) > 1 {
			t.Errorf("consecutiveRatioPercent(%d) = %v, want close to %v%%", c.m, got, c.want)
		}
	}
}

func TestComputeSeries_FirstValueIsAlwaysOne(t *testing.T) {
	for _, m := range []int{5, 10, 20, 40, 80, 15} {
		series := computeSeries(m, 1)
		if !almostEqual(series[0].Value, 1, eps) {
			t.Errorf("computeSeries(%d,1)[0].Value = %v, want 1", m, series[0].Value)
		}
	}
}

func TestComputeSeries_LengthMatchesM(t *testing.T) {
	series := computeSeries(10, 1)
	if len(series) != 10 {
		t.Errorf("len(series) = %d, want 10", len(series))
	}
}

func TestComputeSeries_ScaleFactorAppliesLinearly(t *testing.T) {
	unscaled := computeSeries(10, 1)
	scaled := computeSeries(10, 100)
	for i := range unscaled {
		want := unscaled[i].Value * 100
		if !almostEqual(scaled[i].Scaled, want, 1e-9) {
			t.Errorf("scaled[%d].Scaled = %v, want %v", i, scaled[i].Scaled, want)
		}
	}
}

func TestRoundFrac_HalfRoundsDownNotUp(t *testing.T) {
	// RATIO.C's own rule: only strictly greater than 0.5 rounds up.
	if got := roundFrac(2.5); got != 2 {
		t.Errorf("roundFrac(2.5) = %d, want 2 (round-half-down, not round-half-up)", got)
	}
	if got := roundFrac(2.51); got != 3 {
		t.Errorf("roundFrac(2.51) = %d, want 3", got)
	}
	if got := roundFrac(2.49); got != 2 {
		t.Errorf("roundFrac(2.49) = %d, want 2", got)
	}
}

// TestComputeSeries_R10KnownValues checks a hand-computable case: the
// R10 series (10 steps) scaled by 1, whose 6th entry (index 5) is
// 10^0.5 = sqrt(10), a well known constant.
func TestComputeSeries_R10KnownValues(t *testing.T) {
	series := computeSeries(10, 1)
	want := math.Sqrt(10)
	if !almostEqual(series[5].Value, want, 1e-9) {
		t.Errorf("R10 series[5].Value = %v, want sqrt(10) = %v", series[5].Value, want)
	}
}
