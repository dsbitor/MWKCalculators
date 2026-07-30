// ddh computes how to set up a differential dividing head (DDH) — a
// dividing head augmented with gearing between the spindle and hole
// plate — to divide a workpiece into an arbitrary number of equal
// divisions. It first tries the same non-differential solution
// divhead itself would (a rapid indexing plate shortcut, or a plain
// hole-circle-and-turns solution); if neither works, it searches for
// a differential solution: a hole plate, an indexing increment near
// it, and a change gear train between spindle and hole plate that
// makes the combination work.
//
// A DDH's worm-gear ratio, rapid indexing plate, available hole-circle
// plates, and available change gears are all specific to one owner's
// equipment, so — like diffthrd, divhead, gearatio, and change in
// earlier conversion groups — this program reads them from the user's
// own database rather than the shared reference database; see
// ai/plans/c-to-go-conversion-plan.md, "Data-file strategy for
// Tier 2".
//
// Converted from DDH.C (M. W. Klotz), WorkshopUtilities/ddh.
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

const (
	settingsTable    = "ddh_settings"
	holeCirclesTable = "ddh_hole_circles"
	gearsTable       = "ddh_gears"
)

// maxChainLength is the largest number of gear pairs that can be
// chained together in a differential solution, matching DDH.C's own
// LMAX array bound.
const maxChainLength = 10

// maxSearchEvaluations bounds the total number of gear-pair
// combinations tried across every hole plate and indexing increment
// candidate, replacing the original's interactive keypress abort
// (there is no keyboard to poll here) per coding-style.md Rule 2.
const maxSearchEvaluations = 50_000_000

// settings is the DDH's own fixed configuration.
type settings struct {
	WormGearRatio   int
	RapidIndexHoles int // <= 0 means no rapid indexing plate
}

func gcd(x, y int) int {
	for y != 0 {
		x, y = y, x%y
	}
	return x
}

// turnsRequired mirrors divhead's own function of the same name; see
// MWKGo/divhead/divhead.go.
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

func rapidIndexSolution(rapidIndexHoles, divisions int) (stepHoles int, ok bool) {
	if rapidIndexHoles <= 0 || rapidIndexHoles%divisions != 0 {
		return 0, false
	}
	return rapidIndexHoles / divisions, true
}

type plateSolution struct {
	PlateHoles int
	StepHoles  int
}

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

// gearPair and stage mirror gearatio's own types; see
// MWKGo/gearatio/gearatio.go for the full explanation of the search
// they support.
type gearPair struct{ j, k int }
type stage struct{ firstTeeth, secondTeeth int }

// chainSolution is one gear chain whose combined ratio matched a
// target within tolerance.
type chainSolution struct {
	stages       []stage
	ratio        float64
	errorPercent float64
}

// findChains is gearatio's own chain search, single-target and with a
// caller-supplied shared evaluation counter so a whole differential
// search (many hole plates and indexing increments, each needing its
// own gear chain search) shares one evaluation cap.
func findChains(gearsAscending []int, targetRatio, tolerancePercent float64, maxPairs int, evaluations *int) (found []chainSolution, truncated bool) {
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
				found = append(found, chainSolution{
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

// diffSolution is one differential-indexing solution: a hole plate, an
// indexing increment near the ratio*plate/divisions estimate, whether
// the hole plate and crank turn the same way, and a gear chain
// implementing the required change-gear ratio.
type diffSolution struct {
	holePlate      int
	indexIncrement int
	sameDirection  bool
	gearRatio      float64
	stages         []stage
	achievedRatio  float64
	errorPercent   float64
}

// findDifferentialSolutions searches, for every hole plate and every
// indexing increment within 2 of the estimate ratio*plate/divisions
// (matching DDH.C's own fgear() search window exactly), for a gear
// chain whose ratio matches the resulting required change-gear ratio
// within tolerancePercent.
func findDifferentialSolutions(ratio, divisions int, holeCircles, gearsAscending []int, tolerancePercent float64, maxPairs int) (solutions []diffSolution, truncated bool) {
	evaluations := 0
	for _, p := range holeCircles {
		estimate := float64(ratio) * float64(p) / float64(divisions)
		center := int(math.Floor(estimate))
		for j := center - 2; j <= center+2; j++ {
			if j <= 0 {
				continue // a non-positive indexing increment is not physically meaningful
			}
			approx := float64(ratio) * float64(p) / float64(j)
			rnum := math.Abs((approx - float64(divisions)) * float64(ratio))
			rdenom := approx
			if rnum == 0 {
				continue
			}
			target := rnum / rdenom

			chains, cutShort := findChains(gearsAscending, target, tolerancePercent, maxPairs, &evaluations)
			for _, c := range chains {
				solutions = append(solutions, diffSolution{
					holePlate:      p,
					indexIncrement: j,
					sameDirection:  approx > float64(divisions),
					gearRatio:      target,
					stages:         c.stages,
					achievedRatio:  c.ratio,
					errorPercent:   c.errorPercent,
				})
			}
			if cutShort {
				return solutions, true
			}
		}
	}
	return solutions, false
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

func loadIntColumn(ctx context.Context, db *sql.DB, table, column string) ([]int, error) {
	rows, err := db.QueryContext(ctx, `SELECT `+column+` FROM `+table)
	if err != nil {
		return nil, fmt.Errorf("query %s: %w", table, err)
	}
	defer rows.Close()

	var values []int
	for rows.Next() {
		var v int
		if err := rows.Scan(&v); err != nil {
			return nil, fmt.Errorf("scan %s: %w", table, err)
		}
		values = append(values, v)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate %s: %w", table, err)
	}
	sort.Ints(values)
	return values, nil
}

func main() {
	importSettingsPath := flag.String("import-settings", "", "import worm-gear ratio and rapid index plate from a CSV file and exit")
	importHolesPath := flag.String("import-holes", "", "import available hole-circle plates from a CSV file and exit")
	importGearsPath := flag.String("import-gears", "", "import the available change gear set from a CSV file and exit")
	flag.Parse()

	ctx := context.Background()
	db, err := userdata.Open(ctx)
	if err != nil {
		fmt.Fprintln(os.Stderr, "ddh:", err)
		os.Exit(1)
	}
	defer db.Close()

	if *importSettingsPath != "" || *importHolesPath != "" || *importGearsPath != "" {
		if err := runImports(ctx, db, *importSettingsPath, *importHolesPath, *importGearsPath); err != nil {
			fmt.Fprintln(os.Stderr, "ddh:", err)
			os.Exit(1)
		}
		return
	}

	s, ok, err := loadSettings(ctx, db)
	if err != nil {
		fmt.Fprintln(os.Stderr, "ddh:", err)
		os.Exit(1)
	}
	holeCircles, err := loadIntColumn(ctx, db, holeCirclesTable, "holes")
	if err != nil {
		fmt.Fprintln(os.Stderr, "ddh:", err)
		os.Exit(1)
	}
	gears, err := loadIntColumn(ctx, db, gearsTable, "teeth")
	if err != nil {
		fmt.Fprintln(os.Stderr, "ddh:", err)
		os.Exit(1)
	}
	if !ok || len(holeCircles) == 0 || len(gears) < 2 {
		printSetupHelp()
		os.Exit(1)
	}

	prompter, err := promptio.New(os.Stdin, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, "ddh:", err)
		os.Exit(1)
	}

	fmt.Println("(DIFFERENTIAL) DIVIDING HEAD CALCULATIONS")
	fmt.Println()
	fmt.Printf("%d hole plates and %d gears read from data file\n\n", len(holeCircles), len(gears))

	divisions := prompter.Int("Number of workpiece divisions", 67)
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

	if solutions := plateSolutions(holeCircles, remainderNum, remainderDenom); len(solutions) > 0 {
		fmt.Printf("\n%d full turns of crank\n", wholeTurns)
		for i, sol := range solutions {
			prefix := "and "
			if i > 0 {
				prefix = "or  "
			}
			fmt.Printf("%s%d holes on %d hole plate\n", prefix, sol.StepHoles, sol.PlateHoles)
		}
		return
	}

	fmt.Println("\nNO SOLUTION USING AVAILABLE HOLE PLATES WAS FOUND")
	fmt.Printf("a plate with an integer multiple of %d holes is required\n", remainderDenom)
	fmt.Println("\nTRYING TO FIND A DIFFERENTIAL SOLUTION")
	fmt.Println()

	maxPairs := prompter.Int("Maximum gear pairs to examine", 1)
	if maxPairs > maxChainLength {
		fmt.Printf("Maximum gear pairs to examine is capped at %d\n", maxChainLength)
		maxPairs = maxChainLength
	}
	tolerancePercent := prompter.Float("Allowable gear ratio matching accuracy (%)", 0.01)

	solutions, truncated := findDifferentialSolutions(s.WormGearRatio, divisions, holeCircles, gears, tolerancePercent, maxPairs)
	if len(solutions) == 0 {
		fmt.Println("\nNO SOLUTION USING DIFFERENTIAL INDEXING WAS FOUND")
		os.Exit(1)
	}

	printDifferentialSolutions(divisions, solutions)
	if truncated {
		fmt.Println("\n(search stopped early after reaching the evaluation limit;")
		fmt.Println(" further solutions may exist)")
	}
}

func printDifferentialSolutions(divisions int, solutions []diffSolution) {
	lastPlate, lastIncrement := -1, -1
	for _, sol := range solutions {
		if sol.holePlate != lastPlate || sol.indexIncrement != lastIncrement {
			lastPlate, lastIncrement = sol.holePlate, sol.indexIncrement
			fmt.Println("\n------------")
			fmt.Printf("Hole plate = %d\n", sol.holePlate)
			fmt.Printf("Indexing increment = %d\n", sol.indexIncrement)
			if sol.sameDirection {
				fmt.Println("Hole plate and crank rotate in same direction")
			} else {
				fmt.Println("Hole plate and crank rotate in opposite directions")
			}
			fmt.Printf("Gear ratio = %.6f\n", sol.gearRatio)
			fmt.Println("\nAcceptable gear trains include:")
			fmt.Println()
		}
		for i, st := range sol.stages {
			sep := " "
			if i > 0 {
				sep = "-"
			}
			fmt.Printf("%s%d:%d", sep, st.firstTeeth, st.secondTeeth)
		}
		fmt.Printf("  %.6f  %.3E %%\n", sol.achievedRatio, sol.errorPercent)
	}
}

func printSetupHelp() {
	fmt.Println("No DDH configured yet.")
	fmt.Println("This is machine-specific data (your DDH's own worm-gear ratio, rapid")
	fmt.Println("indexing plate, hole-circle plates, and change gear set), so nothing is")
	fmt.Println("pre-loaded. Import all three parts as CSV:")
	fmt.Println()
	fmt.Println("    ddh -import-settings my-ddh-settings.csv")
	fmt.Println("    ddh -import-holes my-ddh-holes.csv")
	fmt.Println("    ddh -import-gears my-ddh-gears.csv")
	fmt.Println()
	fmt.Println("See docs/calculators/ddh.md for the expected format.")
}

func runImports(ctx context.Context, db *sql.DB, settingsPath, holesPath, gearsPath string) error {
	if settingsPath != "" {
		if err := importCSV(ctx, db, settingsTable, settingsPath); err != nil {
			return err
		}
		fmt.Printf("imported DDH settings from %s\n", settingsPath)
	}
	if holesPath != "" {
		if err := importCSV(ctx, db, holeCirclesTable, holesPath); err != nil {
			return err
		}
		fmt.Printf("imported hole circles from %s\n", holesPath)
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
