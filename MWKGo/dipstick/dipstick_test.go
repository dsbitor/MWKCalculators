package main

import (
	"math"
	"testing"
)

// Each tank shape's wetted-volume fraction must start at exactly 0
// (empty) and reach exactly 1 (full) at its own maximum wetted
// height, and increase monotonically in between: independently
// obvious physical facts about any tank, not just a re-run of each
// shape's own formula.
func TestShapeFractions_EmptyAndFullBoundaries(t *testing.T) {
	cases := []struct {
		name      string
		maxHeight float64
		f         func(float64) float64
	}{
		{"cyl", 10, func(h float64) float64 { return cylFraction(10, h) }},
		{"spher", 10, func(h float64) float64 { return spherFraction(10, h) }},
		{"ellip", 7, func(h float64) float64 { return ellipFraction(10, 7, h) }},
		{"vcart", 10, func(h float64) float64 { return vcartFraction(4, 10, h) }},
		{"hcart", 4, func(h float64) float64 { return hcartFraction(10, 4, h) }},
		{"buck", 6, func(h float64) float64 { return buckFraction(4, 3, 6, h) }},
		{"barrel", 5, func(h float64) float64 { return barrelFraction(8, 6, 5, h) }},
		{"hcylHemisphere", 10, func(h float64) float64 { return hcylHemisphereFraction(10, 20, h) }},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := c.f(0); math.Abs(got) > 1e-9 {
				t.Errorf("f(0) = %v, want 0", got)
			}
			if got := c.f(c.maxHeight); math.Abs(got-1) > 1e-9 {
				t.Errorf("f(maxHeight) = %v, want 1", got)
			}
			prev := -1.0
			for h := 0.0; h <= c.maxHeight; h += c.maxHeight / 20 {
				v := c.f(h)
				if v < prev-1e-9 {
					t.Errorf("f(%v) = %v is less than f at previous height (%v): not monotonic", h, v, prev)
				}
				prev = v
			}
		})
	}
}

func TestDcylDishedFraction_ApproximatelyEmptyAndFull(t *testing.T) {
	// The dished-ends shape's numerical dish integration doesn't
	// reach exactly 0 or exactly 1 at the extremes (see dipstick.md
	// for why: the torispherical dish profile has some cross
	// sections wider than the nominal dish radius, so a small sliver
	// registers as "wetted" even at zero overall dipstick height).
	// This is inherited from the original algorithm, not introduced
	// by this conversion, so a looser tolerance is used here instead
	// of the exact 0/1 identity checked for the other eight shapes.
	d, overall, cylindrical := 100.0, 1000.0, 800.0
	f0 := dcylDishedFraction(d, overall, cylindrical, 0)
	fFull := dcylDishedFraction(d, overall, cylindrical, d)
	if f0 < 0 || f0 > 0.02 {
		t.Errorf("f(0) = %v, want a small value near 0 (within 2%%)", f0)
	}
	if fFull < 0.98 || fFull > 1 {
		t.Errorf("f(d) = %v, want a value near 1 (within 2%%)", fFull)
	}
}

func TestFindWettedHeight_RecoversConsistentFraction(t *testing.T) {
	// Whatever height the search returns, evaluating the shape's own
	// fraction function at that height must reproduce the requested
	// target fraction within the search's own tolerance: the
	// defining property of the search, not a re-run of it.
	target := 0.35
	height, err := findWettedHeight(10, target, 0.00001, func(h float64) float64 { return cylFraction(10, h) })
	if err != nil {
		t.Fatalf("findWettedHeight() error = %v", err)
	}
	if got := cylFraction(10, height); math.Abs(got-target) > 0.00001 {
		t.Errorf("cylFraction(10, %v) = %v, want %v", height, got, target)
	}
}

func TestSegmentArea_FullCircleMatchesKnownArea(t *testing.T) {
	r := 5.0
	got := segmentArea(r, 2*r)
	want := math.Pi * r * r
	if math.Abs(got-want) > 1e-9 {
		t.Errorf("segmentArea(%v, %v) = %v, want %v", r, 2*r, got, want)
	}
}
