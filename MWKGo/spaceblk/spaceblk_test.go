package main

import (
	"math"
	"testing"
)

func TestFindSpacerBlocks_SingleBlockMatch(t *testing.T) {
	blocks := []float64{1.0, 0.5, 0.25, 0.125}
	chosen, truncated := findSpacerBlocks(blocks, 0.5)
	if truncated {
		t.Fatal("search was truncated, want it to complete")
	}
	if len(chosen) != 1 || chosen[0] != 0.5 {
		t.Errorf("findSpacerBlocks(0.5) = %v, want [0.5]", chosen)
	}
}

func TestFindSpacerBlocks_TwoBlockCombination(t *testing.T) {
	blocks := []float64{1.0, 0.5, 0.3, 0.2}
	chosen, truncated := findSpacerBlocks(blocks, 0.8)
	if truncated {
		t.Fatal("search was truncated, want it to complete")
	}
	sum := 0.0
	for _, b := range chosen {
		sum += b
	}
	if math.Abs(sum-0.8) > matchTolerance {
		t.Errorf("findSpacerBlocks(0.8) = %v, sums to %v, want 0.8", chosen, sum)
	}
}

// TestFindSpacerBlocks_UsesDuplicateValueBlocksAsDistinctPositions
// confirms two blocks of the same physical size (a real space block
// set commonly has duplicates, like the two 1.0 blocks in
// SPACEBLK.DAT) can both be used in the same combination, since they
// occupy different array positions.
func TestFindSpacerBlocks_UsesDuplicateValueBlocksAsDistinctPositions(t *testing.T) {
	blocks := []float64{1.0, 1.0, 0.5}
	chosen, truncated := findSpacerBlocks(blocks, 2.0)
	if truncated {
		t.Fatal("search was truncated, want it to complete")
	}
	if len(chosen) != 2 || chosen[0] != 1.0 || chosen[1] != 1.0 {
		t.Errorf("findSpacerBlocks(2.0) = %v, want [1.0 1.0]", chosen)
	}
}

func TestFindSpacerBlocks_NoSolution(t *testing.T) {
	blocks := []float64{1.0, 0.5}
	chosen, truncated := findSpacerBlocks(blocks, 0.9)
	if truncated {
		t.Fatal("search was truncated, want it to complete (small search space)")
	}
	if chosen != nil {
		t.Errorf("findSpacerBlocks(0.9) = %v, want no solution", chosen)
	}
}

func TestFindSpacerBlocks_SkipsBlocksLargerThanTarget(t *testing.T) {
	// The 5.0 block is larger than the target and must never appear
	// in a returned combination.
	blocks := []float64{5.0, 0.6, 0.4}
	chosen, _ := findSpacerBlocks(blocks, 1.0)
	for _, b := range chosen {
		if b > 1.0 {
			t.Errorf("findSpacerBlocks(1.0) chose block %v, larger than the target", b)
		}
	}
}

// TestEvaluateCombination_RejectsRepeatedIndex confirms the same
// block position can't be used twice within one combination, even
// though the search itself only ever generates distinct-length
// index tuples from a single depth loop (this tests the guard
// directly, independent of whether the odometer ever produces such a
// tuple in practice).
func TestEvaluateCombination_RejectsRepeatedIndex(t *testing.T) {
	blocks := []float64{0.5, 0.5}
	if _, ok := evaluateCombination(blocks, []int{0, 0}, 1.0); ok {
		t.Error("evaluateCombination with a repeated index ok = true, want false")
	}
}

func TestEvaluateCombination_AcceptsWithinTolerance(t *testing.T) {
	blocks := []float64{0.30000001, 0.5}
	if _, ok := evaluateCombination(blocks, []int{0, 1}, 0.8); !ok {
		t.Error("evaluateCombination just within tolerance ok = false, want true")
	}
}
