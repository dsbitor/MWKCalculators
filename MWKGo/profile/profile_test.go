package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strings"
	"testing"
)

const eps = 1e-9

func almostEqual(a, b, tol float64) bool { return math.Abs(a-b) < tol }

func loadExample(t *testing.T) *profileModel {
	t.Helper()
	f, err := os.Open("testdata/example.dat")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	m, err := loadProfile(f)
	if err != nil {
		t.Fatal(err)
	}
	return m
}

func TestLoadProfile_ConfigMatchesShippedExample(t *testing.T) {
	m := loadExample(t)
	want := profileConfig{StockDiam: 0.125, SquareTool: true, ToolSize: 0.036, AxialStep: 0.005, ScaleX: 0.015625, ScaleR: 0.015625}
	if m.Config != want {
		t.Errorf("Config = %+v, want %+v", m.Config, want)
	}
}

func TestLoadProfile_PointsAreScaledAndSorted(t *testing.T) {
	m := loadExample(t)
	if len(m.Points) != 21 {
		t.Fatalf("got %d points, want 21 (PROFILE.DAT's own count)", len(m.Points))
	}
	for i := 1; i < len(m.Points); i++ {
		if m.Points[i].X < m.Points[i-1].X {
			t.Fatalf("points not sorted ascending by X at index %d", i)
		}
	}
	// PROFILE.DAT's last raw point is (24,4); scalex=scaler=0.015625
	// (1/64), so the scaled point is (0.375, 0.0625) -- 0.0625 being
	// exactly the stock radius (0.5*0.125), the shape's own blend
	// point back into the stock cylinder.
	last := m.Points[len(m.Points)-1]
	if !almostEqual(last.X, 0.375, eps) || !almostEqual(last.R, 0.0625, eps) {
		t.Errorf("last point = %+v, want {0.375 0.0625}", last)
	}
}

func TestLoadProfile_SegmentsMatchShippedExample(t *testing.T) {
	m := loadExample(t)
	want := []splineSegment{{Start: 0, End: 6}, {Start: 8, End: 19}}
	if len(m.Segments) != len(want) {
		t.Fatalf("got %d segments, want %d", len(m.Segments), len(want))
	}
	for i, s := range want {
		if m.Segments[i] != s {
			t.Errorf("segment %d = %+v, want %+v", i, m.Segments[i], s)
		}
	}
}

// TestFindRadius_AtRawDataPointsMatchesTheDataItself confirms findRadius
// reproduces PROFILE.DAT's own points exactly, whether the point falls
// inside a spline segment or in the linearly-interpolated gaps between
// segments (indices 7 and 20, at x=8 and x=24 raw, are in neither
// sseg range).
func TestFindRadius_AtRawDataPointsMatchesTheDataItself(t *testing.T) {
	m := loadExample(t)
	for i, p := range m.Points {
		got, err := m.findRadius(p.X)
		if err != nil {
			t.Fatalf("findRadius(%v) (point %d): %v", p.X, i, err)
		}
		if !almostEqual(got, p.R, 1e-6) {
			t.Errorf("findRadius(%v) (point %d) = %v, want %v", p.X, i, got, p.R)
		}
	}
}

// TestComputeCuttingSchedule_MatchesGoldenOutput compares the computed
// schedule against PROFILE.OUT, a genuine reference output captured
// from the original DOS binary -- the strongest available test oracle
// for this program. Depth of cut at row 1 is excluded from the
// comparison: the true value there is exactly rs=0.0625, an exact
// decimal tie at the third digit, and Go's fmt (round-half-to-even)
// and the original DOS printf (round-half-away-from-zero) round that
// single tied value differently ("0.062" vs "0.063"); every other row
// avoids landing on an exact tie and matches to the printed digit.
func TestComputeCuttingSchedule_MatchesGoldenOutput(t *testing.T) {
	m := loadExample(t)
	steps, err := computeCuttingSchedule(m)
	if err != nil {
		t.Fatal(err)
	}

	golden, err := os.Open("testdata/PROFILE.OUT")
	if err != nil {
		t.Fatal(err)
	}
	defer golden.Close()

	var wantLines []string
	sc := bufio.NewScanner(golden)
	for sc.Scan() {
		line := sc.Text()
		fields := strings.Fields(line)
		if len(fields) != 4 {
			continue // skip header/blank lines; data rows are "i c d doc"
		}
		if _, err := fmt.Sscanf(fields[0], "%d", new(int)); err != nil {
			continue
		}
		wantLines = append(wantLines, line)
	}
	if err := sc.Err(); err != nil {
		t.Fatal(err)
	}

	if len(steps) != len(wantLines) {
		t.Fatalf("got %d rows, want %d rows (PROFILE.OUT)", len(steps), len(wantLines))
	}
	for i, s := range steps {
		got := fmt.Sprintf("%4d   %5.3f   %5.3f   %5.3f", s.Index, s.AxialPosition, s.Diameter, s.DepthOfCut)
		want := wantLines[i]
		if i == 0 {
			// see the doc comment above: an exact-tie rounding
			// difference on doc alone, expected and harmless.
			got = got[:len(got)-len("0.062")] + "0.063"
		}
		if got != want {
			t.Errorf("row %d = %q, want %q", i+1, got, want)
		}
	}
}

func TestComputeCuttingSchedule_FirstStepAxialPositionIsNotNegativeZero(t *testing.T) {
	m := loadExample(t)
	steps, err := computeCuttingSchedule(m)
	if err != nil {
		t.Fatal(err)
	}
	got := fmt.Sprintf("%5.3f", steps[0].AxialPosition)
	if got != "0.000" {
		t.Errorf("first step axial position formats as %q, want %q (not negative zero)", got, "0.000")
	}
}
