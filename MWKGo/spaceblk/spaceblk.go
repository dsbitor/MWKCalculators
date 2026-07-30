// spaceblk selects which blocks from a set of "space blocks" (a poor
// man's precision gage block set: ground steel cylinders that screw
// together to make an accurate spacer of a required length) to
// combine to build a spacer of a desired size.
//
// The available block set is specific to one owner's own set of
// blocks, so — like diffthrd and divhead in an earlier conversion
// group — this program reads it from the user's own database rather
// than the shared reference database; see
// ai/plans/c-to-go-conversion-plan.md, "Data-file strategy for
// Tier 2".
//
// Converted from SPACEBLK.C (M. W. Klotz), WorkshopUtilities/spaceblk.
// That archive also contains GAGEBLK.C, a faster successor program
// hardcoded for one specific standard 81-block set (needing no data
// file at all); it is deferred to a future group, since SPACEBLK
// alone already covers its case (more slowly) as well as any other
// block set GAGEBLK doesn't handle.
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

const blocksTable = "spaceblk_blocks"

// maxCombinationDepth is the largest number of blocks combined in
// one spacer, matching SPACEBLK.C's own fixed limit of 5.
const maxCombinationDepth = 5

// maxCombinationEvaluations bounds the total number of block
// combinations tried, replacing the original program's interactive
// keypress abort (there is no keyboard to poll here) per
// coding-style.md Rule 2. SPACEBLK.TXT documents that a full 5-deep
// search of an 81-block set examines 3.5 billion combinations; this
// bound stops well short of that if no solution turns up first,
// rather than running effectively forever.
const maxCombinationEvaluations = 50_000_000

// matchTolerance is how close a combination's total must come to the
// target size to count as a match, matching the original's own
// 1e-8 tolerance.
const matchTolerance = 1e-8

// findSpacerBlocks searches, in order of increasing combination size
// (1 block, then 2, and so on up to maxCombinationDepth), for a
// combination of distinct blocks (by position, so two blocks of the
// same size can both be used) summing to target within
// matchTolerance, skipping any block larger than target first. It
// returns the first combination found — not necessarily the one
// using the fewest blocks — matching the original's own behavior,
// and reports whether the search was cut short by
// maxCombinationEvaluations before exhausting every combination.
func findSpacerBlocks(blocksDescending []float64, target float64) (chosen []float64, truncated bool) {
	start := 0
	for start < len(blocksDescending) && blocksDescending[start] > target {
		start++
	}
	n := len(blocksDescending)

	evaluations := 0
	for depth := 1; depth <= maxCombinationDepth && depth <= n-start; depth++ {
		indices := make([]int, depth)
		for i := range indices {
			indices[i] = start
		}

		for {
			if evaluations >= maxCombinationEvaluations {
				return nil, true
			}
			evaluations++

			if solution, ok := evaluateCombination(blocksDescending, indices, target); ok {
				return solution, false
			}

			// Increment the indices like an odometer: rightmost
			// advances fastest, carrying leftward on overflow.
			j := depth - 1
			for ; j >= 0; j-- {
				indices[j]++
				if indices[j] < n {
					break
				}
				indices[j] = start
			}
			if j < 0 {
				break // this depth's combinations are exhausted
			}
		}
	}
	return nil, false
}

// evaluateCombination checks one set of block indices: it must be
// distinct (no block position used twice) and its blocks must sum to
// target within matchTolerance, computed with the same early exit as
// the original (abandon as soon as the running sum exceeds target by
// more than the tolerance).
func evaluateCombination(blocksDescending []float64, indices []int, target float64) ([]float64, bool) {
	for i, idx := range indices {
		for j, other := range indices {
			if i != j && idx == other {
				return nil, false
			}
		}
	}

	sum := 0.0
	for _, idx := range indices {
		sum += blocksDescending[idx]
		if sum-target > matchTolerance {
			return nil, false
		}
	}
	if math.Abs(sum-target) >= matchTolerance {
		return nil, false
	}

	chosen := make([]float64, len(indices))
	for i, idx := range indices {
		chosen[i] = blocksDescending[idx]
	}
	return chosen, true
}

func loadBlocks(ctx context.Context, db *sql.DB) ([]float64, error) {
	rows, err := db.QueryContext(ctx, `SELECT size FROM `+blocksTable)
	if err != nil {
		return nil, fmt.Errorf("query %s: %w", blocksTable, err)
	}
	defer rows.Close()

	var blocks []float64
	for rows.Next() {
		var size float64
		if err := rows.Scan(&size); err != nil {
			return nil, fmt.Errorf("scan block: %w", err)
		}
		blocks = append(blocks, size)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate %s: %w", blocksTable, err)
	}
	sort.Sort(sort.Reverse(sort.Float64Slice(blocks)))
	return blocks, nil
}

func main() {
	importPath := flag.String("import", "", "import a block set from a CSV file (one size column) and exit")
	flag.Parse()

	ctx := context.Background()
	db, err := userdata.Open(ctx)
	if err != nil {
		fmt.Fprintln(os.Stderr, "spaceblk:", err)
		os.Exit(1)
	}
	defer db.Close()

	if *importPath != "" {
		if err := importBlocks(ctx, db, *importPath); err != nil {
			fmt.Fprintln(os.Stderr, "spaceblk:", err)
			os.Exit(1)
		}
		return
	}

	blocks, err := loadBlocks(ctx, db)
	if err != nil {
		fmt.Fprintln(os.Stderr, "spaceblk:", err)
		os.Exit(1)
	}
	if len(blocks) == 0 {
		fmt.Println("No space block set configured yet.")
		fmt.Println("This is machine-specific data (the blocks in your own set), so nothing")
		fmt.Println("is pre-loaded. Import a CSV with a size column:")
		fmt.Println()
		fmt.Println("    spaceblk -import my-block-set.csv")
		fmt.Println()
		fmt.Println("See docs/calculators/spaceblk.md for the expected format.")
		os.Exit(1)
	}

	prompter, err := promptio.New(os.Stdin, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, "spaceblk:", err)
		os.Exit(1)
	}

	fmt.Println("SPACEBLOCK SELECTION UTILITY")
	fmt.Println()
	fmt.Printf("number of blocks read from data file = %d\n", len(blocks))
	total := 0.0
	for _, b := range blocks {
		total += b
	}
	fmt.Printf("total length of all blocks = %.4f\n", total)
	fmt.Printf("largest  block = %.4f\n", blocks[0])
	fmt.Printf("smallest block = %.4f\n\n", blocks[len(blocks)-1])

	target := prompter.Float("Required spacer size", 0.104)
	target = 1e-4 * math.Floor(target*1e4+0.5) // round to 4 decimal places

	if target < blocks[len(blocks)-1] {
		fmt.Println("required size smaller than smallest available block")
		os.Exit(1)
	}
	if target > total {
		fmt.Println("required size larger than total length of all blocks")
		os.Exit(1)
	}

	chosen, truncated := findSpacerBlocks(blocks, target)
	if chosen == nil {
		fmt.Println("\nCAN'T FIND A SOLUTION")
		if truncated {
			fmt.Println("(search stopped early after reaching the evaluation limit)")
		}
		os.Exit(1)
	}

	fmt.Println("\nUSE BLOCKS:")
	sum := 0.0
	for _, b := range chosen {
		fmt.Printf("%.4f\n", b)
		sum += b
	}
	fmt.Printf("TOTAL LENGTH = %.4f\n", sum)
}

func importBlocks(ctx context.Context, db *sql.DB, path string) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("open %s: %w", path, err)
	}
	defer f.Close()

	if err := csvtable.Import(ctx, db, blocksTable, f); err != nil {
		return fmt.Errorf("import %s: %w", path, err)
	}
	fmt.Printf("imported block set from %s\n", path)
	return nil
}
