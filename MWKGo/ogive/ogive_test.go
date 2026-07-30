package main

import (
	"math"
	"os"
	"testing"
)

const eps = 1e-6

func almostEqual(a, b, tol float64) bool { return math.Abs(a-b) < tol }

func loadExample(t *testing.T) ogiveConfig {
	t.Helper()
	f, err := os.Open("testdata/example.dat")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	cfg, err := loadConfig(f)
	if err != nil {
		t.Fatal(err)
	}
	return cfg
}

func TestLoadConfig_MatchesShippedExample(t *testing.T) {
	cfg := loadExample(t)
	want := ogiveConfig{StockDiam: 0.5, SquareTool: true, ToolSize: 0.125, AxialStep: 0.05, OgiveDiam: 0.5, OgiveLength: 1.5}
	if cfg != want {
		t.Errorf("loadConfig() = %+v, want %+v", cfg, want)
	}
}

func TestCaliber_MatchesExample(t *testing.T) {
	cfg := loadExample(t)
	if !almostEqual(cfg.caliber(), 3.0, eps) {
		t.Errorf("caliber = %v, want 3.0 (1.5/0.5)", cfg.caliber())
	}
}

// TestOgiveRadius_TipIsSharp confirms the defining geometric property
// of a tangent ogive: the radius at the very tip (x=0) is exactly
// zero, a sharp point.
func TestOgiveRadius_TipIsSharp(t *testing.T) {
	cfg := loadExample(t)
	got := ogiveRadius(0, cfg)
	if !almostEqual(got, 0, 1e-9) {
		t.Errorf("ogiveRadius(0) = %v, want 0 (a sharp tip)", got)
	}
}

// TestOgiveRadius_BlendsSmoothlyIntoStock confirms the other defining
// property: at x=oglength, the ogive profile's radius exactly equals
// the stock radius, where it blends into the cylindrical stock.
func TestOgiveRadius_BlendsSmoothlyIntoStock(t *testing.T) {
	cfg := loadExample(t)
	got := ogiveRadius(cfg.OgiveLength, cfg)
	want := cfg.stockRadius()
	if !almostEqual(got, want, 1e-6) {
		t.Errorf("ogiveRadius(oglength) = %v, want stock radius %v", got, want)
	}
}

func TestOgiveRadius_BeyondLengthIsStockRadius(t *testing.T) {
	cfg := loadExample(t)
	got := ogiveRadius(cfg.OgiveLength+1, cfg)
	want := cfg.stockRadius()
	if !almostEqual(got, want, eps) {
		t.Errorf("ogiveRadius(oglength+1) = %v, want stock radius %v", got, want)
	}
}

func TestOgiveRadius_IsMonotonicallyIncreasing(t *testing.T) {
	cfg := loadExample(t)
	prev := 0.0
	for x := 0.0; x <= cfg.OgiveLength; x += 0.05 {
		r := ogiveRadius(x, cfg)
		if r < prev-1e-9 {
			t.Errorf("ogiveRadius(%v) = %v, less than previous %v: profile should widen monotonically toward the stock", x, r, prev)
		}
		prev = r
	}
}

func TestComputeCuttingSchedule_EndsNearStockDiameter(t *testing.T) {
	cfg := loadExample(t)
	steps := computeCuttingSchedule(cfg)
	if len(steps) == 0 {
		t.Fatal("expected at least one cutting step")
	}
	last := steps[len(steps)-1]
	if !almostEqual(last.AxialPosition, cfg.OgiveLength, cfg.AxialStep+1e-9) {
		t.Errorf("last step axial position = %v, want close to ogive length %v", last.AxialPosition, cfg.OgiveLength)
	}
	if last.Diameter <= 0 {
		t.Errorf("last step diameter = %v, want positive and approaching the stock diameter", last.Diameter)
	}
}

func TestComputeCuttingSchedule_DiametersAreNonDecreasing(t *testing.T) {
	cfg := loadExample(t)
	steps := computeCuttingSchedule(cfg)
	for i := 1; i < len(steps); i++ {
		if steps[i].Diameter < steps[i-1].Diameter-1e-9 {
			t.Errorf("step %d diameter (%v) is less than step %d's (%v): the ogive widens monotonically toward the stock", i, steps[i].Diameter, i-1, steps[i-1].Diameter)
		}
	}
}
