package main

import (
	"math"
	"testing"
)

func TestStackAndAngleAreInverses(t *testing.T) {
	// Converting a part angle to a stack height and back must
	// reproduce the original angle, regardless of which formula is
	// used for either direction.
	barLength, vBlockAngle := 5.0, 90.0
	for _, wantAngle := range []float64{5, 9.5, 15, 20} {
		toStack := stackFromTaperAngle(wantAngle, barLength, vBlockAngle)
		back := taperAngleFromStack(toStack.StackHeight, barLength, vBlockAngle)
		if math.Abs(back.PartAngle-wantAngle) > 1e-6 {
			t.Errorf("round trip for angle %v: got %v", wantAngle, back.PartAngle)
		}
	}
}

func TestStackFromTaperAngle_DocumentedDefaultInput(t *testing.T) {
	got := stackFromTaperAngle(9.5, 5.0, 90.0)
	// Verified independently before being trusted as an expected
	// value here.
	want := 0.9947234806125259
	if math.Abs(got.StackHeight-want) > 1e-9 {
		t.Errorf("StackHeight = %v, want %v", got.StackHeight, want)
	}
}

func TestDecimalDigit(t *testing.T) {
	cases := []struct {
		x    float64
		k    int
		want int
	}{
		{0.9947, 1, 9},
		{0.9947, 2, 9},
		{0.9947, 3, 4},
		{0.9947, 4, 7},
		{0.9947, 0, 0},
		{1.1345, 1, 1},
		{1.1345, 4, 5},
	}
	for _, c := range cases {
		if got := decimalDigit(c.x, c.k); got != c.want {
			t.Errorf("decimalDigit(%v, %d) = %d, want %d", c.x, c.k, got, c.want)
		}
	}
}

func TestGaugeBlocksFor_BlocksSumToTarget(t *testing.T) {
	// The chosen blocks' own sum must reconstruct the target exactly
	// (to within the algorithm's own 1e-5 tolerance): the defining
	// property of a correct decomposition, not a re-run of the
	// algorithm.
	for _, stack := range []float64{0.9947, 1.1345, 3.7654, 2.0001, 1.0} {
		target, blocks, ok := gaugeBlocksFor(stack)
		if !ok {
			t.Fatalf("gaugeBlocksFor(%v) ok = false, want true", stack)
		}
		var sum float64
		for _, b := range blocks {
			sum += b
		}
		if math.Abs(sum-target) > 1e-5 {
			t.Errorf("stack %v: blocks %v sum to %v, want %v", stack, blocks, sum, target)
		}
	}
}

func TestGaugeBlocksFor_TooSmallReturnsNotOK(t *testing.T) {
	_, _, ok := gaugeBlocksFor(0.03)
	if ok {
		t.Error("gaugeBlocksFor(0.03) ok = true, want false (smaller than the smallest available block)")
	}
}
