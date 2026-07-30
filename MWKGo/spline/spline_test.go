package main

import (
	"math"
	"os"
	"testing"
)

const eps = 1e-9

func almostEqual(a, b, tol float64) bool { return math.Abs(a-b) < tol }

func loadExample(t *testing.T) []splinePoint {
	t.Helper()
	f, err := os.Open("testdata/example.dat")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	points, err := loadPoints(f)
	if err != nil {
		t.Fatal(err)
	}
	return points
}

// TestLoadPoints_SortsOutOfOrderData confirms SPLINE.DAT's own
// deliberately-out-of-order first entry ("-50,-30 ;this entry is
// purposely out of order") ends up in the correct ascending-x
// position after loading.
func TestLoadPoints_SortsOutOfOrderData(t *testing.T) {
	points := loadExample(t)
	want := []splinePoint{{-180, 50}, {-150, -50}, {-100, -80}, {-50, -30}}
	if len(points) != len(want) {
		t.Fatalf("got %d points, want %d", len(points), len(want))
	}
	for i, w := range want {
		if points[i] != w {
			t.Errorf("point %d = %+v, want %+v", i, points[i], w)
		}
	}
}

func TestFitCubicSpline_RejectsFewerThanThreePoints(t *testing.T) {
	if _, err := fitCubicSpline([]splinePoint{{0, 0}, {1, 1}}); err == nil {
		t.Error("expected an error for fewer than 3 points")
	}
}

// TestEval_PassesThroughOriginalPoints confirms the defining property
// of spline interpolation: the fitted curve must reproduce each
// original data point's own y value exactly at its own x.
func TestEval_PassesThroughOriginalPoints(t *testing.T) {
	points := loadExample(t)
	s, err := fitCubicSpline(points)
	if err != nil {
		t.Fatal(err)
	}
	for _, p := range points {
		got, err := s.Eval(p.X)
		if err != nil {
			t.Fatal(err)
		}
		if !almostEqual(got, p.Y, 1e-9) {
			t.Errorf("Eval(%v) = %v, want %v (must pass through the original point)", p.X, got, p.Y)
		}
	}
}

func TestEval_OutOfRangeIsAnError(t *testing.T) {
	points := loadExample(t)
	s, err := fitCubicSpline(points)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := s.Eval(points[0].X - 1); err == nil {
		t.Error("expected an error for z below the data range")
	}
	if _, err := s.Eval(points[len(points)-1].X + 1); err == nil {
		t.Error("expected an error for z above the data range")
	}
}

// TestFitCubicSpline_CollinearPointsFitExactlyToTheLine confirms a
// well known property of natural cubic splines: fitting collinear
// points produces a curve that is itself exactly linear (the
// second-derivative coefficients should all vanish), not just a
// curve that happens to pass through them.
func TestFitCubicSpline_CollinearPointsFitExactlyToTheLine(t *testing.T) {
	points := []splinePoint{{0, 0}, {1, 2}, {2, 4}, {3, 6}, {4, 8}}
	s, err := fitCubicSpline(points)
	if err != nil {
		t.Fatal(err)
	}
	for x := 0.0; x <= 4; x += 0.5 {
		got, err := s.Eval(x)
		if err != nil {
			t.Fatal(err)
		}
		want := 2 * x
		if !almostEqual(got, want, 1e-9) {
			t.Errorf("Eval(%v) = %v, want %v (collinear points should fit exactly to the line y=2x)", x, got, want)
		}
	}
}

func TestSampleCurve_EndpointsMatchOriginalData(t *testing.T) {
	points := loadExample(t)
	s, err := fitCubicSpline(points)
	if err != nil {
		t.Fatal(err)
	}
	samples := sampleCurve(s, 40)
	if len(samples) != 41 {
		t.Fatalf("got %d samples, want 41 (40 segments)", len(samples))
	}
	if !almostEqual(samples[0].Y, points[0].Y, eps) {
		t.Errorf("first sample y = %v, want %v", samples[0].Y, points[0].Y)
	}
	last := samples[len(samples)-1]
	lastPoint := points[len(points)-1]
	if !almostEqual(last.Y, lastPoint.Y, eps) {
		t.Errorf("last sample y = %v, want %v", last.Y, lastPoint.Y)
	}
}
