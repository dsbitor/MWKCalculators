// sinebar computes the stack height needed under a sine bar to set a
// given angle (or, given the stack height and angle already used,
// how sensitive the result is to small errors in either), then finds
// which blocks from a standard 81-piece gage block set combine to
// produce that stack height exactly.
//
// Unlike mtaper's gaugeBlocksFor, which decomposes a stack via a
// greedy digit-extraction heuristic, this program's search is a
// bounded brute-force combination search (reusing the same
// nested-loop "odometer" technique as gearfind in Tier 1 group 11):
// try every single block, then every pair, and so on up to 5 blocks,
// stopping at the first combination that sums to the target exactly.
//
// Converted from SINEBAR.C, WorkshopUtilities/sinebar.
package main

import (
	"errors"
	"fmt"
	"math"
	"os"
	"strings"

	"mwkgo/internal/promptio"
)

func sind(deg float64) float64 { return math.Sin(deg * math.Pi / 180) }

// gageBlockSet is the standard 81-piece gage block set, in
// descending order (as in the original program's own table; the
// original also bubble-sorts it at startup, redundant here since the
// table is already in that order, but noted for fidelity).
var gageBlockSet = [81]float64{
	4.000, 3.000, 2.000, 1.000, 0.950, 0.900, 0.850, 0.800, 0.750, 0.700, 0.650, 0.600, 0.550, 0.500,
	0.450, 0.400, 0.350, 0.300, 0.250, 0.200, 0.150, 0.149, 0.148, 0.147, 0.146, 0.145, 0.144, 0.143,
	0.142, 0.141, 0.140, 0.139, 0.138, 0.137, 0.136, 0.135, 0.134, 0.133, 0.132, 0.131, 0.130, 0.129,
	0.128, 0.127, 0.126, 0.125, 0.124, 0.123, 0.122, 0.121, 0.120, 0.119, 0.118, 0.117, 0.116, 0.115,
	0.114, 0.113, 0.112, 0.111, 0.110, 0.109, 0.108, 0.107, 0.106, 0.105, 0.104, 0.103, 0.102, 0.101,
	0.1009, 0.1008, 0.1007, 0.1006, 0.1005, 0.1004, 0.1003, 0.1002, 0.1001, 0.1000, 0.050,
}

// maxBlocksCombined is the largest number of blocks this search will
// try combining, matching the original program's own nlm constant.
const maxBlocksCombined = 5

// maxCombinationEvaluations bounds the total number of tooth-count
// combinations tried, replacing the original program's interactive
// keypress abort per coding-style.md Rule 2. A realistic stack
// height is found within a handful of blocks; an extremely small
// target (forcing consideration of nearly the entire 81-block set in
// combinations of up to 5) could exceed this bound and report no
// solution found even where one technically exists at greater
// computational cost, matching the practical effect of the
// original's own keypress abort under the same circumstances.
const maxCombinationEvaluations = 20_000_000

var errCombinationSearchTruncated = errors.New("search stopped early: evaluation limit reached")

// gaugeBlockCombination finds a set of up to maxBlocksCombined
// blocks (from those no larger than target) that sum to target
// exactly (within 1e-8), trying combinations of 1 block, then 2,
// then up to maxBlocksCombined, in the same order as the original
// program. It returns the blocks found, or ok=false if no
// combination within the search space sums to target.
func gaugeBlockCombination(target float64) (blocks []float64, ok bool, err error) {
	startIndex := len(gageBlockSet)
	for i, b := range gageBlockSet {
		if b <= target {
			startIndex = i
			break
		}
	}
	if startIndex == len(gageBlockSet) {
		return nil, false, nil
	}

	evaluations := 0
	for numBlocks := 1; numBlocks <= maxBlocksCombined; numBlocks++ {
		indices := make([]int, numBlocks)
		for i := range indices {
			indices[i] = startIndex
		}

		for {
			evaluations++
			if evaluations > maxCombinationEvaluations {
				return nil, false, errCombinationSearchTruncated
			}

			if indicesAreDistinct(indices) {
				sum := 0.0
				withinBudget := true
				for _, idx := range indices {
					sum += gageBlockSet[idx]
					if sum-target > 1e-8 {
						withinBudget = false
						break
					}
				}
				if withinBudget && math.Abs(sum-target) < 1e-8 {
					chosen := make([]float64, numBlocks)
					for i, idx := range indices {
						chosen[i] = gageBlockSet[idx]
					}
					return chosen, true, nil
				}
			}

			// Increment like an odometer: rightmost advances
			// fastest, carrying leftward on overflow.
			j := numBlocks - 1
			for ; j >= 0; j-- {
				indices[j]++
				if indices[j] < len(gageBlockSet) {
					break
				}
				indices[j] = startIndex
			}
			if j < 0 {
				break
			}
		}
	}
	return nil, false, nil
}

func indicesAreDistinct(indices []int) bool {
	for i := range indices {
		for j := range indices {
			if i != j && indices[i] == indices[j] {
				return false
			}
		}
	}
	return true
}

// sineBarResult holds a stack height calculation and its error
// sensitivity to small measurement errors.
type sineBarResult struct {
	Stack                 float64
	RollErrorSensitivity  float64 // deg of angle error per 0.001 error in roll distance
	StackErrorSensitivity float64 // deg of angle error per 0.001 error in stack height
}

// computeSineBarStack returns the stack height for a sine bar with
// the given roll separation and angle, plus (for angles under 90
// degrees) its error sensitivity to small measurement errors in
// either the roll distance or the stack height.
func computeSineBarStack(rollSeparation, angleDeg float64) sineBarResult {
	stack := rollSeparation * sind(angleDeg)
	result := sineBarResult{Stack: stack}
	if angleDeg < 90 {
		degPerRad := 180 / math.Pi
		cosAngle := math.Cos(angleDeg * math.Pi / 180)
		result.StackErrorSensitivity = degPerRad * 0.001 / (rollSeparation * cosAngle)
		result.RollErrorSensitivity = degPerRad * 0.001 * stack / (rollSeparation * rollSeparation * cosAngle)
	}
	return result
}

func main() {
	prompter, err := promptio.New(os.Stdin, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, "sinebar:", err)
		os.Exit(1)
	}

	fmt.Println("SINEBAR CALCULATIONS")
	fmt.Println()

	roll := prompter.Float("Distance between sine bar rolls", 5.0)
	mode := strings.ToLower(prompter.Line("Angle input mode [D]ecimal degrees, (X) deg/min/sec ? "))

	var angle float64
	if mode == "x" {
		deg := prompter.Float("degrees", 30.0)
		minutes := prompter.Float("minutes", 7.0)
		seconds := prompter.Float("seconds", 30.0)
		angle = deg + minutes/60 + seconds/3600
	} else {
		angle = prompter.Float("angle in decimal degrees", 4.93468141)
	}

	fmt.Printf("\nDistance between rolls = %.6f\n", roll)
	fmt.Printf("Angle = %.6f deg\n", angle)

	result := computeSineBarStack(roll, angle)
	fmt.Printf("Stack height = %.6f\n", result.Stack)
	fmt.Println("Stack height measured in same units as roll separation.")
	if angle < 90 {
		fmt.Printf("A .001 error in the roll distance will cause an angle error of %.6f deg\n", result.RollErrorSensitivity)
		fmt.Printf("A .001 error in the stack height will cause an angle error of %.6f deg\n", result.StackErrorSensitivity)
	}

	total := 0.0
	for _, b := range gageBlockSet {
		total += b
	}
	if result.Stack < gageBlockSet[len(gageBlockSet)-1] {
		fmt.Println("required stack size smaller than smallest available block")
		return
	}
	if result.Stack > total {
		fmt.Println("required stack size larger than total length of all blocks")
		return
	}

	target := math.Trunc(result.Stack*1e4) * 1e-4
	blocks, ok, err := gaugeBlockCombination(target)
	fmt.Printf("\nBlocks from standard 81 gage block set needed to form stack = %.4f in:\n\n", target)
	if err != nil {
		fmt.Println(err)
		return
	}
	if !ok {
		fmt.Println("CAN'T FIND A SOLUTION")
		return
	}
	for _, b := range blocks {
		fmt.Printf("%.4f\n", b)
	}
}
