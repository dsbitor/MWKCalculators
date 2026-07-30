// gearatio searches an owner's own set of gears for a chain of gear
// pairs (up to a chosen chain length) whose combined ratio matches a
// required ratio within a chosen accuracy.
//
// The available gear set is specific to one owner's own gears, so —
// like diffthrd, divhead, and spaceblk in earlier conversion groups —
// this program reads it from the user's own database rather than the
// shared reference database; see
// ai/plans/c-to-go-conversion-plan.md, "Data-file strategy for
// Tier 2".
//
// Converted from GEARATIO.C (M. W. Klotz), WorkshopUtilities/gearatio.
package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"math"
	"os"
	"sort"

	"mwkgo/internal/csvtable"
	"mwkgo/internal/promptio"
	"mwkgo/internal/userdata"
)

const gearsTable = "gearatio_gears"

// maxChainLength is the largest number of gear pairs that can be
// chained together, matching GEARATIO.C's own LMAX array bound (the
// original never validates a larger request against this bound, so
// entering one would have corrupted memory; this conversion rejects
// it instead, per coding-style.md Rule 2).
const maxChainLength = 10

// maxSearchEvaluations bounds the total number of gear-pair
// combinations tried across every chain length, replacing the
// original's interactive keypress abort (there is no keyboard to
// poll here) per coding-style.md Rule 2.
const maxSearchEvaluations = 50_000_000

// gearPair is one candidate pair of gears, identified by index into
// an ascending-sorted gear list, with j < k so gears[j] <= gears[k].
type gearPair struct {
	j, k int
}

// stage is one link of a gear chain as it will be displayed: firstTeeth
// is printed before secondTeeth, following GEARATIO.C's own display
// convention of swapping the pair when the chain divides by its ratio
// rather than multiplying.
type stage struct {
	firstTeeth, secondTeeth int
}

// solution is one gear chain whose combined ratio matched the target
// within tolerance.
type solution struct {
	stages       []stage
	ratio        float64
	errorPercent float64
}

// findRatios searches, for each chain length from 1 up to maxPairs,
// every ordered chain of that many gear pairs drawn from gearsAscending
// with no gear reused across positions in the same chain (a gear, once
// mounted in a chain, cannot also be mounted elsewhere in it; two
// different gears that happen to share the same tooth count remain
// distinct positions). At each step the running ratio is multiplied by
// a pair's ratio if it is still above the target, or divided by it
// otherwise (swapping the pair's display order to match), mirroring
// GEARATIO.C's own greedy per-stage decision exactly. Every chain
// landing within tolerancePercent of targetRatio is recorded, matching
// the original's own behavior of reporting every match rather than
// stopping at the first. It returns whether the search was cut short
// by maxSearchEvaluations before exhausting every chain length up to
// maxPairs.
func findRatios(gearsAscending []int, targetRatio, tolerancePercent float64, maxPairs int) (found []solution, truncated bool) {
	var pairs []gearPair
	for j := 0; j < len(gearsAscending); j++ {
		for k := j + 1; k < len(gearsAscending); k++ {
			pairs = append(pairs, gearPair{j: j, k: k})
		}
	}

	used := make([]bool, len(gearsAscending))
	var stages []stage
	evaluations := 0
	cutShort := false

	var search func(depth, chainLength int, rtest float64)
	search = func(depth, chainLength int, rtest float64) {
		if cutShort {
			return
		}
		if depth > 0 {
			errPercent := 100 * (rtest - targetRatio) / targetRatio
			if math.Abs(errPercent) <= tolerancePercent {
				found = append(found, solution{
					stages:       append([]stage(nil), stages...),
					ratio:        rtest,
					errorPercent: errPercent,
				})
			}
		}
		if depth == chainLength {
			return
		}
		for _, p := range pairs {
			if evaluations >= maxSearchEvaluations {
				cutShort = true
				return
			}
			evaluations++

			if used[p.j] || used[p.k] {
				continue
			}
			r := float64(gearsAscending[p.j]) / float64(gearsAscending[p.k])
			if r == 1 {
				continue // identical gears yield a useless 1:1 ratio
			}

			var next float64
			var s stage
			if rtest > targetRatio {
				next = rtest * r
				s = stage{firstTeeth: gearsAscending[p.j], secondTeeth: gearsAscending[p.k]}
			} else {
				next = rtest / r
				s = stage{firstTeeth: gearsAscending[p.k], secondTeeth: gearsAscending[p.j]}
			}

			used[p.j], used[p.k] = true, true
			stages = append(stages, s)
			search(depth+1, chainLength, next)
			stages = stages[:len(stages)-1]
			used[p.j], used[p.k] = false, false

			if cutShort {
				return
			}
		}
	}

	for chainLength := 1; chainLength <= maxPairs; chainLength++ {
		search(0, chainLength, 1.0)
		if cutShort {
			break
		}
	}
	return found, cutShort
}

func loadGears(ctx context.Context, db *sql.DB) ([]int, error) {
	rows, err := db.QueryContext(ctx, `SELECT teeth FROM `+gearsTable)
	if err != nil {
		return nil, fmt.Errorf("query %s: %w", gearsTable, err)
	}
	defer rows.Close()

	var gears []int
	for rows.Next() {
		var teeth int
		if err := rows.Scan(&teeth); err != nil {
			return nil, fmt.Errorf("scan gear: %w", err)
		}
		gears = append(gears, teeth)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate %s: %w", gearsTable, err)
	}
	sort.Ints(gears)
	return gears, nil
}

func main() {
	importPath := flag.String("import", "", "import a gear set from a CSV file (one teeth column) and exit")
	flag.Parse()

	ctx := context.Background()
	db, err := userdata.Open(ctx)
	if err != nil {
		fmt.Fprintln(os.Stderr, "gearatio:", err)
		os.Exit(1)
	}
	defer db.Close()

	if *importPath != "" {
		if err := importGears(ctx, db, *importPath); err != nil {
			fmt.Fprintln(os.Stderr, "gearatio:", err)
			os.Exit(1)
		}
		return
	}

	gears, err := loadGears(ctx, db)
	if err != nil {
		fmt.Fprintln(os.Stderr, "gearatio:", err)
		os.Exit(1)
	}
	if len(gears) < 2 {
		fmt.Println("No gear set configured yet.")
		fmt.Println("This is machine-specific data (the gears you own), so nothing is")
		fmt.Println("pre-loaded. Import a CSV with a teeth column:")
		fmt.Println()
		fmt.Println("    gearatio -import my-gears.csv")
		fmt.Println()
		fmt.Println("See docs/calculators/gearatio.md for the expected format.")
		os.Exit(1)
	}

	prompter, err := promptio.New(os.Stdin, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, "gearatio:", err)
		os.Exit(1)
	}

	fmt.Println("GEAR RATIO CALCULATIONS")
	fmt.Println()
	fmt.Printf("Number of gears read = %d\n", len(gears))

	ratio := prompter.Float("Required ratio", 0.5)
	tolerancePercent := prompter.Float("Allowable ratio accuracy (%)", 0.01)
	maxPairs := prompter.Int("Maximum gear pairs to examine", 1)
	if maxPairs > maxChainLength {
		fmt.Printf("Maximum gear pairs to examine is capped at %d\n", maxChainLength)
		maxPairs = maxChainLength
	}

	fmt.Println("\npatience...")
	solutions, truncated := findRatios(gears, ratio, tolerancePercent, maxPairs)

	if len(solutions) == 0 {
		fmt.Println("\nNO SOLUTION FOUND")
		if truncated {
			fmt.Println("(search stopped early after reaching the evaluation limit)")
		}
		os.Exit(1)
	}

	for _, sol := range solutions {
		printSolution(sol)
	}
	if truncated {
		fmt.Println("\n(search stopped early after reaching the evaluation limit;")
		fmt.Println(" further solutions may exist)")
	}
}

func printSolution(sol solution) {
	for i, s := range sol.stages {
		sep := " "
		if i > 0 {
			sep = "-"
		}
		fmt.Printf("%s%d:%d", sep, s.firstTeeth, s.secondTeeth)
	}
	fmt.Printf("  %.6f  %.3E %%\n", sol.ratio, sol.errorPercent)
}

func importGears(ctx context.Context, db *sql.DB, path string) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("open %s: %w", path, err)
	}
	defer f.Close()

	if err := csvtable.Import(ctx, db, gearsTable, f); err != nil {
		return fmt.Errorf("import %s: %w", path, err)
	}
	fmt.Printf("imported gear set from %s\n", path)
	return nil
}
