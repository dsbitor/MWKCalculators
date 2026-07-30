// change computes lathe change-gear pairs to cut a required thread
// pitch, given the lathe's own leadscrew's effective pitch(es) and its
// own set of available change gears.
//
// A lathe's leadscrew pitch(es) and its own set of change gears are
// specific to one owner's equipment, so — like diffthrd, divhead, and
// gearatio in earlier conversion groups — this program reads them
// from the user's own database rather than the shared reference
// database; see ai/plans/c-to-go-conversion-plan.md, "Data-file
// strategy for Tier 2".
//
// Converted from CHANGE.C (M. W. Klotz), WorkshopUtilities/change.
// That archive also contains CHANGEX.C, the change gear calculator's
// earlier, redundant-output predecessor, now superseded by CHANGE
// itself; it is deferred to a future group along with combi, ratio,
// and belt.
package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"math"
	"os"
	"sort"
	"strings"

	"mwkgo/internal/csvtable"
	"mwkgo/internal/promptio"
	"mwkgo/internal/userdata"
)

const (
	settingsTable = "change_settings"
	pitchesTable  = "change_leadscrew_pitches"
	gearsTable    = "change_gears"
)

// maxChainLength is the largest number of gear pairs that can be
// chained together, matching CHANGE.C's own LMAX-1 array bound (the
// original silently clamps a data file's own maxpair value to this
// range; this conversion does the same).
const maxChainLength = 9

// maxSearchEvaluations bounds the total number of gear-pair
// combinations tried across every leadscrew pitch and chain length,
// replacing the original's interactive keypress abort (there is no
// keyboard to poll here) per coding-style.md Rule 2.
const maxSearchEvaluations = 50_000_000

type gearPair struct{ j, k int }

// stage is one link of a gear chain as it will be displayed:
// firstTeeth is printed before secondTeeth, following CHANGE.C's own
// display convention of swapping the pair when the chain divides by
// its ratio rather than multiplying.
type stage struct{ firstTeeth, secondTeeth int }

// solution is one gear chain whose combined ratio matched the target
// within tolerance, for one leadscrew pitch.
type solution struct {
	leadscrewPitch float64
	stages         []stage
	ratio          float64
	errorPercent   float64
}

// findChains searches, for each chain length from 1 up to maxPairs,
// every ordered chain of that many gear pairs drawn from
// gearsAscending with no gear reused across positions in the same
// chain, exactly as gearatio's own search does; see
// MWKGo/gearatio/gearatio.go for the full explanation of the
// algorithm this mirrors. It returns every chain within
// tolerancePercent of targetRatio and whether the search was cut short
// by maxSearchEvaluations.
func findChains(gearsAscending []int, targetRatio, tolerancePercent float64, maxPairs int, evaluations *int) (found []solution, truncated bool) {
	var pairs []gearPair
	for j := 0; j < len(gearsAscending); j++ {
		for k := j + 1; k < len(gearsAscending); k++ {
			pairs = append(pairs, gearPair{j: j, k: k})
		}
	}

	used := make([]bool, len(gearsAscending))
	var stages []stage
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
			if *evaluations >= maxSearchEvaluations {
				cutShort = true
				return
			}
			*evaluations++

			if used[p.j] || used[p.k] {
				continue
			}
			r := float64(gearsAscending[p.j]) / float64(gearsAscending[p.k])
			if r == 1 {
				continue
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

// findAllChains runs findChains once per leadscrew pitch, tagging each
// resulting solution with the pitch it was found for.
func findAllChains(leadscrewPitches []float64, gearsAscending []int, desiredPitch, tolerancePercent float64, maxPairs int) (found []solution, truncated bool) {
	evaluations := 0
	for _, lsp := range leadscrewPitches {
		ratio := lsp / desiredPitch
		solutions, cutShort := findChains(gearsAscending, ratio, tolerancePercent, maxPairs, &evaluations)
		for i := range solutions {
			solutions[i].leadscrewPitch = lsp
		}
		found = append(found, solutions...)
		if cutShort {
			return found, true
		}
	}
	return found, false
}

func loadSettings(ctx context.Context, db *sql.DB) (maxPairs int, ok bool, err error) {
	err = db.QueryRowContext(ctx, `SELECT max_pairs FROM `+settingsTable+` WHERE id = 1`).Scan(&maxPairs)
	if err == sql.ErrNoRows {
		return 0, false, nil
	}
	if err != nil {
		return 0, false, fmt.Errorf("query %s: %w", settingsTable, err)
	}
	if maxPairs < 1 {
		maxPairs = 1
	}
	if maxPairs > maxChainLength {
		maxPairs = maxChainLength
	}
	return maxPairs, true, nil
}

func loadFloatColumn(ctx context.Context, db *sql.DB, table, column string) ([]float64, error) {
	rows, err := db.QueryContext(ctx, `SELECT `+column+` FROM `+table)
	if err != nil {
		return nil, fmt.Errorf("query %s: %w", table, err)
	}
	defer rows.Close()

	var values []float64
	for rows.Next() {
		var v float64
		if err := rows.Scan(&v); err != nil {
			return nil, fmt.Errorf("scan %s: %w", table, err)
		}
		values = append(values, v)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate %s: %w", table, err)
	}
	sort.Float64s(values)
	return values, nil
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
	importSettingsPath := flag.String("import-settings", "", "import the maximum chain length to search from a CSV file and exit")
	importPitchesPath := flag.String("import-pitches", "", "import effective leadscrew pitch(es) from a CSV file and exit")
	importGearsPath := flag.String("import-gears", "", "import the available change gear set from a CSV file and exit")
	flag.Parse()

	ctx := context.Background()
	db, err := userdata.Open(ctx)
	if err != nil {
		fmt.Fprintln(os.Stderr, "change:", err)
		os.Exit(1)
	}
	defer db.Close()

	if *importSettingsPath != "" || *importPitchesPath != "" || *importGearsPath != "" {
		if err := runImports(ctx, db, *importSettingsPath, *importPitchesPath, *importGearsPath); err != nil {
			fmt.Fprintln(os.Stderr, "change:", err)
			os.Exit(1)
		}
		return
	}

	maxPairs, ok, err := loadSettings(ctx, db)
	if err != nil {
		fmt.Fprintln(os.Stderr, "change:", err)
		os.Exit(1)
	}
	leadscrewPitches, err := loadFloatColumn(ctx, db, pitchesTable, "pitch_tpi")
	if err != nil {
		fmt.Fprintln(os.Stderr, "change:", err)
		os.Exit(1)
	}
	gears, err := loadGears(ctx, db)
	if err != nil {
		fmt.Fprintln(os.Stderr, "change:", err)
		os.Exit(1)
	}
	if !ok || len(leadscrewPitches) == 0 || len(gears) < 2 {
		printSetupHelp()
		os.Exit(1)
	}

	prompter, err := promptio.New(os.Stdin, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, "change:", err)
		os.Exit(1)
	}

	fmt.Println("CHANGE GEAR CALCULATIONS")
	fmt.Println()
	fmt.Printf("Maximum number of change gear pairs to examine = %d\n", maxPairs)
	fmt.Printf("Number of leadscrew pitches read = %d\n", len(leadscrewPitches))
	fmt.Printf("Number of change gears read = %d\n", len(gears))

	metric := strings.EqualFold(prompter.Line("Type of thread to cut - [I]mperial or (M)etric? "), "m")

	var desiredPitchTPI float64
	if metric {
		pmm := prompter.Float("Thread pitch to cut (mm)", 1.0)
		desiredPitchTPI = 25.4 / pmm
	} else {
		desiredPitchTPI = prompter.Float("Thread pitch to cut (tpi)", 21.0)
	}
	tolerancePercent := prompter.Float("Allowable thread pitch accuracy (%)", 0.01)

	fmt.Printf("\nThread pitch to cut = %.2f tpi = %.2f mm/thread\n", desiredPitchTPI, 25.4/desiredPitchTPI)
	fmt.Println("\npatience...")

	solutions, truncated := findAllChains(leadscrewPitches, gears, desiredPitchTPI, tolerancePercent, maxPairs)
	if len(solutions) == 0 {
		fmt.Println("\nNO SOLUTION FOUND")
		if truncated {
			fmt.Println("(search stopped early after reaching the evaluation limit)")
		}
		os.Exit(1)
	}

	printSolutions(leadscrewPitches, solutions, len(leadscrewPitches) > 1)
	if truncated {
		fmt.Println("\n(search stopped early after reaching the evaluation limit;")
		fmt.Println(" further solutions may exist)")
	}
}

// printSolutions groups solutions by leadscrew pitch (in the order the
// pitches were tried) and, within each group, marks the lowest-error
// solution with "**", matching CHANGE.C's own per-elsp "best" marker.
func printSolutions(leadscrewPitches []float64, solutions []solution, showPitch bool) {
	for _, lsp := range leadscrewPitches {
		var group []solution
		for _, s := range solutions {
			if s.leadscrewPitch == lsp {
				group = append(group, s)
			}
		}
		if len(group) == 0 {
			continue
		}

		fmt.Printf("\nEffective lead screw pitch = %.2f\n", lsp)

		best := math.Inf(1)
		for _, s := range group {
			if math.Abs(s.errorPercent) < best {
				best = math.Abs(s.errorPercent)
			}
		}

		for _, s := range group {
			if showPitch {
				fmt.Printf("[%.2f]", lsp)
			}
			for i, st := range s.stages {
				sep := " "
				if i > 0 {
					sep = "-"
				}
				fmt.Printf("%s%d:%d", sep, st.firstTeeth, st.secondTeeth)
			}
			fmt.Printf("  %.6f  %.3E %%", s.ratio, s.errorPercent)
			if math.Abs(s.errorPercent) == best {
				fmt.Print("  **")
			}
			fmt.Println()
		}
	}
}

func printSetupHelp() {
	fmt.Println("No change gear setup configured yet.")
	fmt.Println("This is machine-specific data (your lathe's own leadscrew pitch(es) and")
	fmt.Println("change gear set), so nothing is pre-loaded. Import all three parts as CSV:")
	fmt.Println()
	fmt.Println("    change -import-settings my-change-settings.csv")
	fmt.Println("    change -import-pitches my-change-pitches.csv")
	fmt.Println("    change -import-gears my-change-gears.csv")
	fmt.Println()
	fmt.Println("See docs/calculators/change.md for the expected format.")
}

func runImports(ctx context.Context, db *sql.DB, settingsPath, pitchesPath, gearsPath string) error {
	if settingsPath != "" {
		if err := importCSV(ctx, db, settingsTable, settingsPath); err != nil {
			return err
		}
		fmt.Printf("imported change gear settings from %s\n", settingsPath)
	}
	if pitchesPath != "" {
		if err := importCSV(ctx, db, pitchesTable, pitchesPath); err != nil {
			return err
		}
		fmt.Printf("imported leadscrew pitches from %s\n", pitchesPath)
	}
	if gearsPath != "" {
		if err := importCSV(ctx, db, gearsTable, gearsPath); err != nil {
			return err
		}
		fmt.Printf("imported change gears from %s\n", gearsPath)
	}
	return nil
}

func importCSV(ctx context.Context, db *sql.DB, table, path string) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("open %s: %w", path, err)
	}
	defer f.Close()

	if err := csvtable.Import(ctx, db, table, f); err != nil {
		return fmt.Errorf("import %s into %s: %w", path, table, err)
	}
	return nil
}
