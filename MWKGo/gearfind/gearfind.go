// gearfind performs an exhaustive brute-force search for sets of
// gear pairs (up to a chosen maximum number of pairs, each within a
// chosen tooth-count range) whose combined ratio matches a desired
// ratio within a chosen tolerance. It is deliberately inelegant but
// thorough: for gear trains where the available tooth counts aren't
// constrained to a fixed selection, this searches every possibility.
//
// Converted from GEARFIND.C, WorkshopUtilities/gearfind. This
// archive also contains GEARF1.C and GEARF2.C, faster successor
// programs specialized for exactly one and exactly two gear pairs
// respectively (GEARF2 using a continued-fraction technique akin to
// rfrac); both are deferred to a future group, since GEARFIND alone
// already covers their cases (more slowly) as well as the 3+ pair
// cases they don't.
package main

import (
	"errors"
	"fmt"
	"math"
	"os"

	"mwkgo/internal/promptio"
)

// maxGearPairs is the largest number of gear pairs this search
// supports, matching the original program's own LMAX/2 limit (10
// loop variables, two per pair).
const maxGearPairs = 5

// maxGearRatioEvaluations bounds the total number of tooth-count
// combinations tried, replacing the original program's interactive
// keypress abort (there is no keyboard to poll here) per
// coding-style.md Rule 2. A default search (2 pairs, a 25-tooth
// range) evaluates well under a million combinations; this bound
// exists to catch a much wider range or pair count from running
// effectively forever.
const maxGearRatioEvaluations = 20_000_000

// gearRatioSolution is one combination of gear pairs whose combined
// ratio falls within the requested tolerance of the desired ratio.
type gearRatioSolution struct {
	Pairs     [][2]int
	Ratio     float64
	ErrorPct  float64
	IsNewBest bool // true if this solution improved on every prior one found
}

var errTooManyPairs = errors.New("maximum number of gear pairs must be at most 5")

// searchGearRatios exhaustively searches every combination of
// maxPairs gear pairs, each gear's tooth count ranging from minTeeth
// to maxTeeth, for combined ratios within allowedErrorPct percent of
// desiredRatio. It returns every solution found, in the order
// encountered (matching the original program's own nested-loop
// order), and reports whether the search was cut short by
// maxGearRatioEvaluations before exhausting every combination.
func searchGearRatios(desiredRatio, allowedErrorPct float64, maxPairs, minTeeth, maxTeeth int) (solutions []gearRatioSolution, truncated bool, err error) {
	if maxPairs > maxGearPairs {
		return nil, false, errTooManyPairs
	}

	numLoops := 2 * maxPairs
	indices := make([]int, numLoops)
	for i := range indices {
		indices[i] = minTeeth
	}

	best := math.Inf(1)
	for evaluations := 0; evaluations < maxGearRatioEvaluations; evaluations++ {
		ratio := 1.0
		for i := 0; i < numLoops; i += 2 {
			ratio *= float64(indices[i]) / float64(indices[i+1])
		}
		errPct := 100 * (ratio - desiredRatio) / desiredRatio

		if math.Abs(errPct) <= allowedErrorPct {
			isNewBest := math.Abs(errPct) < best
			if isNewBest {
				best = math.Abs(errPct)
			}
			pairs := make([][2]int, maxPairs)
			for i := 0; i < maxPairs; i++ {
				pairs[i] = [2]int{indices[2*i], indices[2*i+1]}
			}
			solutions = append(solutions, gearRatioSolution{Pairs: pairs, Ratio: ratio, ErrorPct: errPct, IsNewBest: isNewBest})
		}

		// Increment the loop indices like an odometer: rightmost
		// advances fastest, carrying leftward on overflow.
		j := numLoops - 1
		for ; j >= 0; j-- {
			indices[j]++
			if indices[j] <= maxTeeth {
				break
			}
			indices[j] = minTeeth
		}
		if j < 0 {
			return solutions, false, nil
		}
	}
	return solutions, true, nil
}

func main() {
	prompter, err := promptio.New(os.Stdin, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, "gearfind:", err)
		os.Exit(1)
	}

	fmt.Println("FIND GEARS TO PRODUCE A GIVEN RATIO")
	fmt.Println()

	desiredRatio := prompter.Float("Desired ratio", 1.9945)
	allowedErrorPct := prompter.Float("Allowable ratio error (%)", 0.01)
	maxPairs := prompter.Int("Maximum number of gear pairs (< 5)", 2)
	minTeeth := prompter.Int("Minimum number of gear teeth", 16)
	maxTeeth := prompter.Int("Maximum number of gear teeth", 40)

	fmt.Println("\npatience ...")

	solutions, truncated, err := searchGearRatios(desiredRatio, allowedErrorPct, maxPairs, minTeeth, maxTeeth)
	if err != nil {
		fmt.Fprintln(os.Stderr, "gearfind:", err)
		os.Exit(1)
	}

	fmt.Printf("Desired ratio = %g\n\n", desiredRatio)
	for i, s := range solutions {
		fmt.Printf("%d,", i+1)
		for _, pair := range s.Pairs {
			fmt.Printf("%d:%d-", pair[0], pair[1])
		}
		fmt.Printf("  %g (%g %%)", s.Ratio, s.ErrorPct)
		if s.IsNewBest {
			fmt.Println(" **")
		} else {
			fmt.Println()
		}
	}

	if len(solutions) == 0 {
		fmt.Println("\nNO SOLUTION FOUND")
	}
	if truncated {
		fmt.Println("\nSEARCH STOPPED EARLY: evaluation limit reached")
	}
}
