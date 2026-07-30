package main

import (
	"math"
	"testing"
)

func TestComputeSineBarStack_KnownRightTriangleValues(t *testing.T) {
	// A 30-degree sine bar angle with a 10-unit roll separation
	// gives a stack height of exactly half the roll separation
	// (sin(30) = 0.5): an independently obvious trigonometric fact.
	got := computeSineBarStack(10, 30)
	if math.Abs(got.Stack-5) > 1e-9 {
		t.Errorf("Stack = %v, want 5", got.Stack)
	}
}

func TestComputeSineBarStack_NoErrorSensitivityAt90Degrees(t *testing.T) {
	got := computeSineBarStack(5, 90)
	if got.RollErrorSensitivity != 0 || got.StackErrorSensitivity != 0 {
		t.Errorf("error sensitivities = %v, %v, want 0, 0 at 90 degrees", got.RollErrorSensitivity, got.StackErrorSensitivity)
	}
}

func TestGaugeBlockCombination_SingleBlockExactMatch(t *testing.T) {
	// 0.5 is itself one of the standard 81 block sizes, so the
	// search should find it immediately as a single-block solution.
	blocks, ok, err := gaugeBlockCombination(0.5)
	if err != nil {
		t.Fatalf("gaugeBlockCombination() error = %v", err)
	}
	if !ok {
		t.Fatal("ok = false, want true")
	}
	if len(blocks) != 1 || math.Abs(blocks[0]-0.5) > 1e-9 {
		t.Errorf("blocks = %v, want [0.5]", blocks)
	}
}

func TestGaugeBlockCombination_BlocksSumToTarget(t *testing.T) {
	// A target requiring more than one block (0.1005 + 0.9 = 1.0005,
	// neither of which alone matches): the defining property of a
	// correct combination search is that its own chosen blocks sum
	// to the target, checked directly rather than asserting a
	// specific expected combination.
	target := 1.0005
	blocks, ok, err := gaugeBlockCombination(target)
	if err != nil {
		t.Fatalf("gaugeBlockCombination() error = %v", err)
	}
	if !ok {
		t.Fatal("ok = false, want true")
	}
	var sum float64
	for _, b := range blocks {
		sum += b
	}
	if math.Abs(sum-target) > 1e-8 {
		t.Errorf("blocks %v sum to %v, want %v", blocks, sum, target)
	}
}

func TestGaugeBlockCombination_NoBlockSmallEnoughReturnsNotOK(t *testing.T) {
	_, ok, err := gaugeBlockCombination(0.01)
	if err != nil {
		t.Fatalf("gaugeBlockCombination() error = %v", err)
	}
	if ok {
		t.Error("ok = true, want false (no block is small enough to start from)")
	}
}
