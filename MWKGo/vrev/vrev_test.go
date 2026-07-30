package main

import (
	"math"
	"strings"
	"testing"
)

// vrevTXTHoles is the exact hole list from VREV.TXT's worked example
// (a pin-punch holder): two each of three sizes plus one larger one.
var vrevTXTHoles = []float64{0.265625, 0.265625, 0.328125, 0.328125, 0.40625, 0.40625, 0.40625, 0.484375}

// TestMinBoltCircleDiameter_WorkedExample reproduces VREV.TXT's own
// worked example: with the default 0.1 spacing, the minimum bolt
// circle diameter for this hole set is documented as 1.1915.
func TestMinBoltCircleDiameter_WorkedExample(t *testing.T) {
	got := minBoltCircleDiameter(0.1, vrevTXTHoles)
	if math.Abs(got-1.1915) > 1e-3 {
		t.Errorf("minBoltCircleDiameter() = %v, want approximately 1.1915", got)
	}
}

// TestActualSpacing_WorkedExample reproduces VREV.TXT's own worked
// example: choosing a bolt circle diameter of 1.25 (larger than the
// 1.1915 minimum) for this hole set increases the resulting spacing
// to 0.1236.
func TestActualSpacing_WorkedExample(t *testing.T) {
	got := actualSpacing(1.25, vrevTXTHoles)
	if math.Abs(got-0.1236) > 1e-3 {
		t.Errorf("actualSpacing(1.25, ...) = %v, want approximately 0.1236", got)
	}
}

// TestLayoutHoles_WorkedExample reproduces VREV.TXT's own worked
// example's full hole table for a 1.25 bolt circle diameter.
func TestLayoutHoles_WorkedExample(t *testing.T) {
	spacing := actualSpacing(1.25, vrevTXTHoles)
	positions := layoutHoles(1.25, spacing, vrevTXTHoles)

	want := []holePosition{
		{Diameter: 0.2656, X: 0.6250, Y: 0.0000, AngleDeg: 0.0000},
		{Diameter: 0.2656, X: 0.5065, Y: 0.3662, AngleDeg: 35.8702},
		{Diameter: 0.3281, X: 0.1650, Y: 0.6028, AngleDeg: 74.6899},
		{Diameter: 0.3281, X: -0.2785, Y: 0.5595, AngleDeg: 116.4592},
		{Diameter: 0.4063, X: -0.5943, Y: 0.1934, AngleDeg: 161.9756},
		{Diameter: 0.4063, X: -0.5344, Y: -0.3241, AngleDeg: 211.2392},
		{Diameter: 0.4063, X: -0.1031, Y: -0.6164, AngleDeg: 260.5027},
		{Diameter: 0.4844, X: 0.4310, Y: -0.4526, AngleDeg: 313.5997},
	}
	if len(positions) != len(want) {
		t.Fatalf("layoutHoles() returned %d positions, want %d", len(positions), len(want))
	}
	for i, w := range want {
		p := positions[i]
		if math.Abs(p.X-w.X) > 1e-3 || math.Abs(p.Y-w.Y) > 1e-3 || math.Abs(p.AngleDeg-w.AngleDeg) > 1e-3 {
			t.Errorf("layoutHoles()[%d] = %+v, want %+v", i, p, w)
		}
	}
}

// TestLayoutHoles_StockDiameters reproduces VREV.TXT's own worked
// example: theoretical minimum stock diameter 1.7344, recommended
// 1.9816, for bcd=1.25.
func TestLayoutHoles_StockDiameters(t *testing.T) {
	bcd := 1.25
	spacing := actualSpacing(bcd, vrevTXTHoles)
	hdMax := 0.484375

	theoretical := bcd + hdMax
	recommended := bcd + hdMax + 2*spacing

	if math.Abs(theoretical-1.7344) > 1e-3 {
		t.Errorf("theoretical stock diameter = %v, want 1.7344", theoretical)
	}
	if math.Abs(recommended-1.9816) > 1e-3 {
		t.Errorf("recommended stock diameter = %v, want 1.9816", recommended)
	}
}

// TestAngleExcess_IsZeroAtMinBoltCircleDiameter is self-verifying:
// the whole point of minBoltCircleDiameter is that the holes' total
// subtended angle at that diameter is exactly 360 degrees.
func TestAngleExcess_IsZeroAtMinBoltCircleDiameter(t *testing.T) {
	spacing := 0.1
	minBCD := minBoltCircleDiameter(spacing, vrevTXTHoles)
	excess := angleExcess(minBCD/2, spacing, vrevTXTHoles)
	if math.Abs(excess) > 1e-3 {
		t.Errorf("angleExcess at minBoltCircleDiameter = %v, want approximately 0", excess)
	}
}

func TestLoadHoleDiameters_Parses(t *testing.T) {
	data := "STARTOFDATA\n0.25\n0.5\n0.75\nENDOFDATA\n"
	holes, err := loadHoleDiameters(strings.NewReader(data))
	if err != nil {
		t.Fatalf("loadHoleDiameters() error = %v", err)
	}
	want := []float64{0.25, 0.5, 0.75}
	if len(holes) != len(want) {
		t.Fatalf("loadHoleDiameters() = %v, want %v", holes, want)
	}
	for i := range want {
		if holes[i] != want[i] {
			t.Errorf("loadHoleDiameters()[%d] = %v, want %v", i, holes[i], want[i])
		}
	}
}
