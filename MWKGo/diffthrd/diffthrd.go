// diffthrd finds the best pair of available thread pitches to cut a
// differential thread (two threads, of different pitch, cut on the
// same shaft, whose combined effect is a very fine effective pitch)
// approximating a desired effective pitch, then works out the nut
// dimensions and travel needed to use it.
//
// The available thread pitches are specific to one lathe's own
// screwcutting gear train, not universal, so unlike fits/speed/gage/
// expand in the same conversion group, this program reads from the
// user's own database rather than the shared reference database; see
// ai/plans/c-to-go-conversion-plan.md, "Data-file strategy for
// Tier 2".
//
// Converted from DIFFTHRD.C (M. W. Klotz), WorkshopUtilities/diffthrd.
package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"math"
	"os"

	"mwkgo/internal/csvtable"
	"mwkgo/internal/promptio"
	"mwkgo/internal/userdata"
)

const pitchesTable = "diffthrd_pitches"

// bestPair is the pitch pair (coarse and fine, coarse <= fine) whose
// differential effective pitch is closest to a desired effective
// pitch, found by an exhaustive search of every ordered pair drawn
// from the available pitches. The search is O(n^2) in the number of
// available pitches, which is bounded by however many entries the
// user has actually imported (typically a few dozen at most, one
// lathe's own gear train), matching the original program's own
// exhaustive nested loop.
func bestPair(pitches []float64, desired float64) (coarse, fine, effective float64, found bool) {
	bestError := math.Inf(1)
	for i, pi := range pitches {
		for j, pj := range pitches {
			if i == j {
				continue
			}
			pc, pf := pi, pj
			if pc > pf {
				pc, pf = pf, pc
			}
			diff := 1/pc - 1/pf
			if diff == 0 {
				continue
			}
			eff := 1 / diff
			if err := math.Abs(eff - desired); err < bestError {
				bestError = err
				coarse, fine, effective = pc, pf, eff
				found = true
			}
		}
	}
	return coarse, fine, effective, found
}

// exactMatch reports whether desired is already one of the available
// pitches, meaning it can be cut directly without a differential
// setup.
func exactMatch(pitches []float64, desired float64) bool {
	for _, p := range pitches {
		if p == desired {
			return true
		}
	}
	return false
}

func loadPitches(ctx context.Context, db *sql.DB) ([]float64, error) {
	rows, err := db.QueryContext(ctx, `SELECT pitch_tpi FROM `+pitchesTable)
	if err != nil {
		return nil, fmt.Errorf("query %s: %w", pitchesTable, err)
	}
	defer rows.Close()

	var pitches []float64
	for rows.Next() {
		var p float64
		if err := rows.Scan(&p); err != nil {
			return nil, fmt.Errorf("scan pitch: %w", err)
		}
		pitches = append(pitches, p)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate %s: %w", pitchesTable, err)
	}
	return pitches, nil
}

func main() {
	importPath := flag.String("import", "", "import available thread pitches from a CSV file (one pitch_tpi column) and exit")
	flag.Parse()

	ctx := context.Background()
	db, err := userdata.Open(ctx)
	if err != nil {
		fmt.Fprintln(os.Stderr, "diffthrd:", err)
		os.Exit(1)
	}
	defer db.Close()

	if *importPath != "" {
		if err := importPitches(ctx, db, *importPath); err != nil {
			fmt.Fprintln(os.Stderr, "diffthrd:", err)
			os.Exit(1)
		}
		return
	}

	pitches, err := loadPitches(ctx, db)
	if err != nil {
		fmt.Fprintln(os.Stderr, "diffthrd:", err)
		os.Exit(1)
	}
	if len(pitches) == 0 {
		fmt.Println("No available thread pitches configured yet.")
		fmt.Println("This is machine-specific data (the pitches your own lathe's change gears can cut),")
		fmt.Println("so nothing is pre-loaded. Import a CSV with a pitch_tpi column:")
		fmt.Println()
		fmt.Println("    diffthrd -import my-lathe-pitches.csv")
		fmt.Println()
		fmt.Println("See docs/calculators/diffthrd.md for the expected format.")
		os.Exit(1)
	}

	prompter, err := promptio.New(os.Stdin, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, "diffthrd:", err)
		os.Exit(1)
	}

	fmt.Println("DIFFERENTIAL THREAD CALCULATIONS")
	fmt.Println()
	fmt.Printf("Number of data items read = %d\n\n", len(pitches))

	desired := prompter.Float("Desired effective pitch of differential thread (tpi)", 100.0)
	if exactMatch(pitches, desired) {
		fmt.Println("You can cut this thread directly with available screwcutting gear.")
	}

	coarse, fine, effective, found := bestPair(pitches, desired)
	if !found {
		fmt.Println("No usable pair of available pitches was found.")
		return
	}
	fmt.Printf("\nOf available threads, best match to %.3f tpi is:\n", desired)
	fmt.Printf("Coarse thread = %.3f tpi = %.3f mm/thrd\n", coarse, 25.4/coarse)
	fmt.Printf("Fine thread = %.3f tpi = %.3f mm/thrd\n", fine, 25.4/fine)
	fmt.Printf("with an effective pitch of %.3f tpi\n\n", effective)

	coarse = prompter.Float("Pitch of coarse thread (tpi)", coarse)
	fine = prompter.Float("Pitch of fine thread (tpi)", fine)
	coarseNut := prompter.Float("Thickness of coarse (fixed) nut (in)", 0.375)
	fineNut := prompter.Float("Thickness of fine (movable) nut (in)", 0.25)
	motion := prompter.Float("Desired motion of movable nut (in)", 0.25)

	if coarse > fine {
		coarse, fine = fine, coarse
	}
	pe := 1 / (1/coarse - 1/fine)
	turns := motion * pe

	fmt.Printf("\nEffective pitch = %.3f tpi\n", pe)
	fmt.Printf("Motion for one revolution = %.5f in\n", 1/pe)
	fmt.Printf("Total turns to obtain desired motion = %.3f\n", turns)
	fmt.Printf("Minimum length of coarse thread needed = %.3f in\n", coarseNut+turns/coarse)
	fmt.Printf("Minimum length of fine thread needed = %.3f in\n", fineNut+turns/fine)
	fmt.Printf("Maximum distance between nuts = %.3f in\n", turns/coarse)
	fmt.Printf("Minimum distance between nuts = %.3f in\n", turns/fine)
}

func importPitches(ctx context.Context, db *sql.DB, path string) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("open %s: %w", path, err)
	}
	defer f.Close()

	if err := csvtable.Import(ctx, db, pitchesTable, f); err != nil {
		return fmt.Errorf("import %s: %w", path, err)
	}
	fmt.Printf("imported thread pitches from %s\n", path)
	return nil
}
