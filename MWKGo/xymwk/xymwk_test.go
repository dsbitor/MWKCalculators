package main

import (
	"math"
	"os"
	"testing"
)

const eps = 1e-6

func almostEqual(a, b, tol float64) bool { return math.Abs(a-b) < tol }

func loadExample(t *testing.T) []xyPoint {
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

// TestLoadPoints_MatchesShippedExample confirms the file's five
// points load in file order (not sorted -- point position is how the
// user selects a point, so order is load-bearing here, unlike
// profile/ogive/egg).
func TestLoadPoints_MatchesShippedExample(t *testing.T) {
	points := loadExample(t)
	want := []xyPoint{
		{150, 100},
		{193.301, 75},
		{150, 150},
		{106.699, 75},
		{180, 140},
	}
	if len(points) != len(want) {
		t.Fatalf("got %d points, want %d", len(points), len(want))
	}
	for i, w := range want {
		if !almostEqual(points[i].X, w.X, eps) || !almostEqual(points[i].Y, w.Y, eps) {
			t.Errorf("point %d = %+v, want %+v", i+1, points[i], w)
		}
	}
}

// TestReferTo_OriginPointLandsAtZero confirms XYMWK.DAT's own shipped
// example: points 2-5 are each exactly 50 units from point 1, so
// referencing to point 1 puts point 1 at the origin and every other
// point at radius exactly 50.
func TestReferTo_OriginPointLandsAtZero(t *testing.T) {
	points := loadExample(t)
	origin := points[0]
	transformed := referTo(points, origin.X, origin.Y)

	if !almostEqual(transformed[0].X, 0, eps) || !almostEqual(transformed[0].Y, 0, eps) {
		t.Fatalf("referenced origin point = %+v, want (0,0)", transformed[0])
	}
	// The shipped example's own coordinates are rounded to 3 decimal
	// places (e.g. 193.301 for 150+50*cos(30 deg)), so "radius
	// exactly 50" only holds to about 1e-4, not eps's tighter bound.
	for i := 1; i < len(transformed); i++ {
		r := math.Hypot(transformed[i].X, transformed[i].Y)
		if !almostEqual(r, 50, 1e-3) {
			t.Errorf("point %d radius from origin = %v, want 50", i+1, r)
		}
	}
}

// TestAlignTo_SecondPointLandsOnXAxis confirms the defining property
// of align: after aligning to points 1 and 3, point 3's y coordinate
// (and therefore its angle) is exactly zero.
func TestAlignTo_SecondPointLandsOnXAxis(t *testing.T) {
	points := loadExample(t)
	a, b := points[0], points[2]
	transformed := alignTo(points, a.X, a.Y, b.X, b.Y)

	if !almostEqual(transformed[0].X, 0, eps) || !almostEqual(transformed[0].Y, 0, eps) {
		t.Fatalf("aligned origin point = %+v, want (0,0)", transformed[0])
	}
	if !almostEqual(transformed[2].Y, 0, eps) {
		t.Errorf("aligned second point y = %v, want 0 (on the x-axis)", transformed[2].Y)
	}
	if transformed[2].X <= 0 {
		t.Errorf("aligned second point x = %v, want positive (on the +x axis)", transformed[2].X)
	}
}

func TestDist_MatchesShippedExampleRadius(t *testing.T) {
	points := loadExample(t)
	got := dist(points[0].X, points[0].Y, points[2].X, points[2].Y)
	if !almostEqual(got, 50, eps) {
		t.Errorf("dist(point1,point3) = %v, want 50", got)
	}
}

// TestCircumcircle_MatchesShippedExampleCenter confirms the strongest
// available hand-derivable check: points 2, 4, and 5 all lie exactly
// on the radius-50 circle centered at point 1 (150,100), so their
// circumcircle must reproduce that center and radius exactly.
func TestCircumcircle_MatchesShippedExampleCenter(t *testing.T) {
	points := loadExample(t)
	p2, p4, p5 := points[1], points[3], points[4]
	xc, yc, r, err := circumcircle(p2.X, p2.Y, p4.X, p4.Y, p5.X, p5.Y)
	if err != nil {
		t.Fatal(err)
	}
	// Same 3-decimal-place shipped-data precision limit as above.
	if !almostEqual(xc, 150, 1e-3) || !almostEqual(yc, 100, 1e-3) {
		t.Errorf("center = (%v,%v), want (150,100)", xc, yc)
	}
	if !almostEqual(r, 50, 1e-3) {
		t.Errorf("radius = %v, want 50", r)
	}
}

// TestCircumcircle_AllThreePointsEquidistantFromCenter is a second,
// independent check on an unrelated triple from the same file: rather
// than a known center, it confirms the defining property of any
// circumcircle directly.
func TestCircumcircle_AllThreePointsEquidistantFromCenter(t *testing.T) {
	points := loadExample(t)
	p1, p2, p3 := points[0], points[1], points[2]
	xc, yc, r, err := circumcircle(p1.X, p1.Y, p2.X, p2.Y, p3.X, p3.Y)
	if err != nil {
		t.Fatal(err)
	}
	for i, p := range []xyPoint{p1, p2, p3} {
		got := dist(p.X, p.Y, xc, yc)
		if !almostEqual(got, r, eps) {
			t.Errorf("point %d distance from center = %v, want %v (the circle's own radius)", i+1, got, r)
		}
	}
}

// TestCircumcircle_CollinearPointsIsAnError confirms the fix for the
// gap XYMWK.C's own comment on circ3 flags but never implements: three
// collinear points have no unique circumscribing circle, and this
// conversion reports that as an error instead of silently returning a
// circle centered on the origin.
func TestCircumcircle_CollinearPointsIsAnError(t *testing.T) {
	_, _, _, err := circumcircle(0, 0, 1, 1, 2, 2)
	if err == nil {
		t.Fatal("expected an error for three collinear points, got nil")
	}
}

func TestFormatDMS_WholeDegreesExample(t *testing.T) {
	got := formatDMS(90)
	want := "(  90:00:00)"
	if got != want {
		t.Errorf("formatDMS(90) = %q, want %q", got, want)
	}
}

func TestFormatDMS_NegativeAngle(t *testing.T) {
	got := formatDMS(-36.8699)
	want := "(- 36:52:12)"
	if got != want {
		t.Errorf("formatDMS(-36.8699) = %q, want %q", got, want)
	}
}
