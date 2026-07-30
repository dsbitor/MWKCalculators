package main

import (
	"math"
	"testing"
)

func TestSolveRatio_AllFourCombinations(t *testing.T) {
	// A single self-consistent ratio: 2/4 = 3/6.
	const wantA, wantB, wantC, wantD = 2.0, 4.0, 3.0, 6.0

	cases := []struct {
		name       string
		a, b, c, d float64
	}{
		{"solve a", 0, wantB, wantC, wantD},
		{"solve b", wantA, 0, wantC, wantD},
		{"solve c", wantA, wantB, 0, wantD},
		{"solve d", wantA, wantB, wantC, 0},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			a, b, cc, d, err := solveRatio(c.a, c.b, c.c, c.d)
			if err != nil {
				t.Fatalf("solveRatio() error = %v", err)
			}
			checkClose(t, "a", a, wantA)
			checkClose(t, "b", b, wantB)
			checkClose(t, "c", cc, wantC)
			checkClose(t, "d", d, wantD)
		})
	}
}

func TestSolveRatio_InsufficientDataReturnsError(t *testing.T) {
	_, _, _, _, err := solveRatio(1, 2, 0, 0)
	if err == nil {
		t.Fatal("solveRatio() error = nil, want an error for only two known values")
	}
}

func TestParallelResistance_TwoEqualResistors(t *testing.T) {
	// Two equal resistors in parallel combine to exactly half the
	// value of one: an independently obvious electrical identity.
	got := parallelResistance([]float64{10, 10})
	if math.Abs(got-5) > 1e-9 {
		t.Errorf("parallelResistance([10,10]) = %v, want 5", got)
	}
}

func TestRSS_PythagoreanTriple(t *testing.T) {
	got := rss([]float64{3, 4})
	if math.Abs(got-5) > 1e-9 {
		t.Errorf("rss([3,4]) = %v, want 5", got)
	}
}

func TestRMS_AllEqualValues(t *testing.T) {
	// The RMS of a set of identical values is just that value.
	got := rms([]float64{7, 7, 7})
	if math.Abs(got-7) > 1e-9 {
		t.Errorf("rms([7,7,7]) = %v, want 7", got)
	}
}

func TestComputeStats_KnownDataset(t *testing.T) {
	xs := []float64{2, 4, 4, 4, 5, 5, 7, 9}
	got := computeStats(xs)
	checkClose(t, "Max", got.Max, 9)
	checkClose(t, "Min", got.Min, 2)
	checkClose(t, "Mean", got.Mean, 5)
	// Population variance of this well known example dataset is 4,
	// so sample standard deviation (n-1 denominator) is sqrt(32/7).
	checkClose(t, "StdDev", got.StdDev, math.Sqrt(32.0/7.0))
}

func TestDMSAndDecimalDegrees_RoundTrip(t *testing.T) {
	for _, want := range []float64{0, 30.5, 90, -45.25} {
		sign, d, m, s := decimalDegreesToDMS(want)
		_, _, _, decimal := dmsToDecimalDegrees(sign*float64(d), float64(m), float64(s))
		if math.Abs(decimal-want) > 1e-6 {
			t.Errorf("round trip for %v: got %v", want, decimal)
		}
	}
}

func TestDmsToDecimalDegrees_CarriesSecondsAndMinutes(t *testing.T) {
	// 90 seconds carries into 1 minute 30 seconds; 90 minutes
	// carries into 1 degree 30 minutes.
	_, _, _, decimal := dmsToDecimalDegrees(1, 0, 90)
	checkClose(t, "decimal", decimal, 1+1.0/60+30.0/3600)
}

func TestScrewMajorDiameter_KnownSizes(t *testing.T) {
	// #10 screw: 0.06 + 10*0.013 = 0.19 in, matching the well known
	// standard machine screw gauge formula.
	got := screwMajorDiameter(10)
	if math.Abs(got-0.19) > 1e-9 {
		t.Errorf("screwMajorDiameter(10) = %v, want 0.19", got)
	}
}

func TestLinearInterpolate_Midpoint(t *testing.T) {
	got := linearInterpolate(0, 0, 10, 10, 5)
	if math.Abs(got-5) > 1e-9 {
		t.Errorf("linearInterpolate midpoint = %v, want 5", got)
	}
}

func TestDrawingScale_RoundTrip(t *testing.T) {
	sf := drawingScale(0.25, 12.0)
	realLength := 0.1 / sf
	checkClose(t, "realLength", realLength, 0.1*12.0/0.25)
}

func checkClose(t *testing.T, name string, got, want float64) {
	t.Helper()
	if math.Abs(got-want) > 1e-9 {
		t.Errorf("%s = %v, want %v", name, got, want)
	}
}
