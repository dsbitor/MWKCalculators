// divhead computes how to set up a dividing head (worm-gear ratio
// fixed by the head itself) to divide a workpiece into an arbitrary
// number of equal divisions: how many full crank turns are needed,
// and, when that isn't a whole number, either a rapid indexing plate
// shortcut or which hole circle and hole count on a standard index
// plate makes up the remainder.
//
// A dividing head's worm-gear ratio, rapid indexing plate, and set of
// available hole-circle plates are all specific to one owner's
// equipment, not universal, so this program reads from the user's
// own database rather than the shared reference database; see
// ai/plans/c-to-go-conversion-plan.md, "Data-file strategy for
// Tier 2".
//
// Converted from DIVHEAD.C (M. W. Klotz), WorkshopUtilities/divhead.
package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"os"

	"mwkgo/internal/csvtable"
	"mwkgo/internal/promptio"
	"mwkgo/internal/userdata"
)

const (
	settingsTable    = "divhead_settings"
	holeCirclesTable = "divhead_hole_circles"
)

// settings is the dividing head's own fixed configuration: its
// worm-gear ratio, and the hole count of its rapid indexing plate (0
// if it has none).
type settings struct {
	WormGearRatio   int
	RapidIndexHoles int
}

// gcd returns the greatest common divisor of x and y, both assumed
// positive.
func gcd(x, y int) int {
	for y != 0 {
		x, y = y, x%y
	}
	return x
}

// turnsRequired returns the number of full crank turns and, when the
// division is not exact, the remaining fraction of a turn reduced to
// lowest terms (remainderHoles/plateHoles). remainderNum is 0 when
// the division is exact.
func turnsRequired(ratio, divisions int) (wholeTurns, remainderNum, remainderDenom int) {
	wholeTurns = ratio / divisions
	remainderNum = ratio - wholeTurns*divisions
	remainderDenom = divisions
	if remainderNum != 0 {
		g := gcd(remainderNum, remainderDenom)
		remainderNum /= g
		remainderDenom /= g
	}
	return wholeTurns, remainderNum, remainderDenom
}

// rapidIndexSolution returns how many holes to step on the rapid
// indexing plate, if the plate's hole count is an exact multiple of
// divisions.
func rapidIndexSolution(rapidIndexHoles, divisions int) (stepHoles int, ok bool) {
	if rapidIndexHoles <= 0 || rapidIndexHoles%divisions != 0 {
		return 0, false
	}
	return rapidIndexHoles / divisions, true
}

// plateSolution is one usable hole-circle plate for the remainder
// fraction of a turn.
type plateSolution struct {
	PlateHoles int // total holes on this plate
	StepHoles  int // holes to step on this plate
}

// plateSolutions returns every hole-circle plate whose hole count is
// an exact multiple of remainderDenom, each paired with how many
// holes to step on it.
func plateSolutions(holeCircles []int, remainderNum, remainderDenom int) []plateSolution {
	var solutions []plateSolution
	for _, holes := range holeCircles {
		if holes%remainderDenom == 0 {
			solutions = append(solutions, plateSolution{
				PlateHoles: holes,
				StepHoles:  holes * remainderNum / remainderDenom,
			})
		}
	}
	return solutions
}

func loadSettings(ctx context.Context, db *sql.DB) (settings, bool, error) {
	var s settings
	err := db.QueryRowContext(ctx, `SELECT worm_gear_ratio, rapid_index_holes FROM `+settingsTable+` WHERE id = 1`).
		Scan(&s.WormGearRatio, &s.RapidIndexHoles)
	if err == sql.ErrNoRows {
		return settings{}, false, nil
	}
	if err != nil {
		return settings{}, false, fmt.Errorf("query %s: %w", settingsTable, err)
	}
	return s, true, nil
}

func loadHoleCircles(ctx context.Context, db *sql.DB) ([]int, error) {
	rows, err := db.QueryContext(ctx, `SELECT holes FROM `+holeCirclesTable)
	if err != nil {
		return nil, fmt.Errorf("query %s: %w", holeCirclesTable, err)
	}
	defer rows.Close()

	var holes []int
	for rows.Next() {
		var h int
		if err := rows.Scan(&h); err != nil {
			return nil, fmt.Errorf("scan hole circle: %w", err)
		}
		holes = append(holes, h)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate %s: %w", holeCirclesTable, err)
	}
	return holes, nil
}

func main() {
	importSettingsPath := flag.String("import-settings", "", "import worm-gear ratio and rapid index plate from a CSV file and exit")
	importHolesPath := flag.String("import-holes", "", "import available hole-circle plates from a CSV file and exit")
	flag.Parse()

	ctx := context.Background()
	db, err := userdata.Open(ctx)
	if err != nil {
		fmt.Fprintln(os.Stderr, "divhead:", err)
		os.Exit(1)
	}
	defer db.Close()

	if *importSettingsPath != "" || *importHolesPath != "" {
		if err := runImports(ctx, db, *importSettingsPath, *importHolesPath); err != nil {
			fmt.Fprintln(os.Stderr, "divhead:", err)
			os.Exit(1)
		}
		return
	}

	s, ok, err := loadSettings(ctx, db)
	if err != nil {
		fmt.Fprintln(os.Stderr, "divhead:", err)
		os.Exit(1)
	}
	holeCircles, err := loadHoleCircles(ctx, db)
	if err != nil {
		fmt.Fprintln(os.Stderr, "divhead:", err)
		os.Exit(1)
	}
	if !ok || len(holeCircles) == 0 {
		printSetupHelp()
		os.Exit(1)
	}

	prompter, err := promptio.New(os.Stdin, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, "divhead:", err)
		os.Exit(1)
	}

	fmt.Println("DIVIDING HEAD CALCULATIONS")
	fmt.Println()

	divisions := prompter.Int("Number of workpiece divisions", 14)
	turns := float64(s.WormGearRatio) / float64(divisions)
	wholeTurns, remainderNum, remainderDenom := turnsRequired(s.WormGearRatio, divisions)

	fmt.Printf("\nDH Worm Gear Ratio = %d:1\n", s.WormGearRatio)
	if s.RapidIndexHoles > 0 {
		fmt.Printf("Holes in rapid indexing plate on spindle = %d\n", s.RapidIndexHoles)
	} else {
		fmt.Println("No rapid indexing plate available")
	}
	fmt.Printf("Divisions of Workpiece = %d\n", divisions)
	fmt.Printf("Ratio/Divisions = %d/%d = %.4f\n", s.WormGearRatio, divisions, turns)
	if remainderNum != 0 {
		fmt.Printf("Turns required = %d/%d = %d & %d/%d\n",
			remainderDenom*wholeTurns+remainderNum, remainderDenom, wholeTurns, remainderNum, remainderDenom)
	} else {
		fmt.Printf("Turns required = %d\n", wholeTurns)
	}

	if step, ok := rapidIndexSolution(s.RapidIndexHoles, divisions); ok {
		fmt.Printf("\nStep %d holes on rapid indexing plate\n", step)
		return
	}

	if remainderNum == 0 {
		fmt.Printf("\n%d full turns of crank\n", wholeTurns)
		return
	}

	solutions := plateSolutions(holeCircles, remainderNum, remainderDenom)
	if len(solutions) == 0 {
		fmt.Println("\nNO SOLUTION USING AVAILABLE HOLE PLATES WAS FOUND")
		fmt.Printf("a plate with an integer multiple of %d holes is required\n", remainderDenom)
		return
	}
	fmt.Printf("\n%d full turns of crank\n", wholeTurns)
	for i, sol := range solutions {
		prefix := "and "
		if i > 0 {
			prefix = "or  "
		}
		fmt.Printf("%s%d holes on %d hole plate\n", prefix, sol.StepHoles, sol.PlateHoles)
	}
}

func printSetupHelp() {
	fmt.Println("No dividing head configured yet.")
	fmt.Println("This is machine-specific data (your dividing head's own worm-gear ratio,")
	fmt.Println("rapid indexing plate, and available hole-circle plates), so nothing is")
	fmt.Println("pre-loaded. Import both parts as CSV:")
	fmt.Println()
	fmt.Println("    divhead -import-settings my-divhead-settings.csv")
	fmt.Println("    divhead -import-holes my-divhead-holes.csv")
	fmt.Println()
	fmt.Println("See docs/calculators/divhead.md for the expected format.")
}

func runImports(ctx context.Context, db *sql.DB, settingsPath, holesPath string) error {
	if settingsPath != "" {
		if err := importCSV(ctx, db, settingsTable, settingsPath); err != nil {
			return err
		}
		fmt.Printf("imported dividing head settings from %s\n", settingsPath)
	}
	if holesPath != "" {
		if err := importCSV(ctx, db, holeCirclesTable, holesPath); err != nil {
			return err
		}
		fmt.Printf("imported hole circles from %s\n", holesPath)
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
