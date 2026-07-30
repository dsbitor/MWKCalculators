// drill looks up drill sizes and designations, computes tap drill
// diameters (for both cutting and thread-forming taps, imperial and
// metric), and works out a step-drilling schedule for enlarging a
// pilot hole to a final size in stages.
//
// The 371-entry drill size table is universal reference data (the
// same for anyone), so it lives in the shared reference.db SQLite
// database, the same mechanism used by fits, speed, gage, expand,
// findthrd, weight, and unit earlier in this project; see
// ai/plans/c-to-go-conversion-plan.md, "Data-file strategy for
// Tier 2". drill itself has no drawing of any kind (unlike the other
// seven Tier 3 programs it's grouped with by the plan's own
// data-dependency note) — its own DRILL.C never calls a single
// graphics primitive, only colored text output.
//
// Converted from DRILL.C (M. W. Klotz), WorkshopUtilities/drill.
package main

import (
	"context"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"

	"mwkgo/internal/promptio"
	"mwkgo/internal/refdata"
)

// mmPerInch is DRILL.C's own mpi constant.
const mmPerInch = 25.4

// sixEighthsTan60 and fiveEighthsTan60 are the two published tap
// thread-form constants DRILL.TXT discusses at length: (6/8)*tan(60)
// for the American National (imperial) / Standard (metric) thread
// form, and (5/8)*tan(60) for the Unified (imperial) / ISO (metric)
// thread form.
var (
	sixEighthsTan60  = (6.0 / 8.0) * math.Tan(60*math.Pi/180)
	fiveEighthsTan60 = (5.0 / 8.0) * math.Tan(60*math.Pi/180)
)

// drill is one entry in the drill size table.
type drill struct {
	Name   string
	SizeIn float64
}

// loadDrills returns every drill, ordered ascending by size — DRILL.C
// itself sorts before use, and finddrill()'s own boundary checks and
// nearest-size search both depend on ascending order.
func loadDrills(ctx context.Context) ([]drill, error) {
	db, err := refdata.Open(ctx)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.QueryContext(ctx, `SELECT name, diameter_in FROM drills ORDER BY diameter_in`)
	if err != nil {
		return nil, fmt.Errorf("query drills: %w", err)
	}
	defer rows.Close()

	var drills []drill
	for rows.Next() {
		var d drill
		if err := rows.Scan(&d.Name, &d.SizeIn); err != nil {
			return nil, fmt.Errorf("scan drill: %w", err)
		}
		drills = append(drills, d)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate drills: %w", err)
	}
	return drills, nil
}

// findDrill returns the index of the drill closest to size, matching
// DRILL.C's own finddrill(): sizes at or beyond either end of the
// table clamp to that end; otherwise the closest match by absolute
// difference wins. If excludeMetric is set, a metric-only match (its
// name contains "mm" but not "=" — a dual-labeled entry like
// "13=4.70 mm" still counts as acceptable) is walked back to the
// nearest smaller non-metric-only entry instead.
func findDrill(drills []drill, size float64, excludeMetric bool) int {
	if size <= drills[0].SizeIn {
		return 0
	}
	if size >= drills[len(drills)-1].SizeIn {
		return len(drills) - 1
	}

	best := math.Inf(1)
	k := 0
	for i, d := range drills {
		if diff := math.Abs(size - d.SizeIn); diff < best {
			best, k = diff, i
		}
	}

	if excludeMetric {
		for k > 0 && strings.Contains(drills[k].Name, "mm") && !strings.Contains(drills[k].Name, "=") {
			k--
		}
	}
	return k
}

// findByDesignation looks up a drill by its exact name (case-
// insensitive) first, then falls back to a substring match — matching
// DRILL.C's own fhole() exactly, including its own inconsistency: the
// exact match is case-insensitive but the substring fallback is
// case-sensitive, since only the exact-match branch uses strcmpi.
func findByDesignation(drills []drill, designation string) (int, bool) {
	for i, d := range drills {
		if strings.EqualFold(d.Name, designation) {
			return i, true
		}
	}
	for i, d := range drills {
		if strings.Contains(d.Name, designation) {
			return i, true
		}
	}
	return 0, false
}

// parseSize parses a decimal or fractional dimension string like
// "1.5", "3/8", or "1-1/2", matching DRILL.C's own dfinput(): an
// unparseable string returns 0 rather than an error, the same
// forgiving fallback atof() itself has.
func parseSize(s string) float64 {
	s = strings.TrimSpace(s)
	if idx := strings.IndexByte(s, '/'); idx >= 0 {
		whole := 0.0
		rest := s
		if dash := strings.IndexByte(s, '-'); dash >= 0 {
			whole, _ = strconv.ParseFloat(strings.TrimSpace(s[:dash]), 64)
			rest = s[dash+1:]
		}
		slash := strings.IndexByte(rest, '/')
		if slash < 0 {
			return whole
		}
		num, _ := strconv.ParseFloat(strings.TrimSpace(rest[:slash]), 64)
		den, _ := strconv.ParseFloat(strings.TrimSpace(rest[slash+1:]), 64)
		if den == 0 {
			return whole
		}
		return whole + num/den
	}
	v, _ := strconv.ParseFloat(s, 64)
	return v
}

// reportLine is one line of a drill lookup report: a candidate drill,
// its size difference from the target, whether it's the actual match,
// and (for tap-drill reports only) the percentage depth of thread it
// would cut.
type reportLine struct {
	Drill      drill
	Diff       float64
	IsMatch    bool
	DotPercent float64
	ShowDot    bool
}

// buildReport returns up to 7 lines (3 on either side of the match,
// clipped at the table's ends) describing drill k and its neighbors
// relative to the target size, matching DRILL.C's own dprint().
func buildReport(drills []drill, k int, size float64, td, pitchPerIn, threadFormConstant float64, showDot bool) []reportLine {
	var lines []reportLine
	for l := k - 3; l <= k+3; l++ {
		if l < 0 || l >= len(drills) {
			continue
		}
		line := reportLine{Drill: drills[l], Diff: drills[l].SizeIn - size, IsMatch: l == k}
		if showDot {
			line.ShowDot = true
			line.DotPercent = 100 * (td - drills[l].SizeIn) * pitchPerIn / threadFormConstant
		}
		lines = append(lines, line)
	}
	return lines
}

func printReport(lines []reportLine) {
	fmt.Println()
	for _, l := range lines {
		marker := ""
		if l.IsMatch {
			marker = "***** "
		}
		fmt.Printf("%s(%s) with size %.4f (%+.4f)", marker, l.Drill.Name, l.Drill.SizeIn, l.Diff)
		if l.ShowDot {
			fmt.Printf(" %.0f%% dot", l.DotPercent)
		}
		if l.IsMatch {
			fmt.Print(" *****")
		}
		fmt.Println()
	}
	fmt.Println()
}

// circleArea and diameterFromArea match DRILL.C's own carea()/cdiam().
func circleArea(diameter float64) float64   { return 0.25 * math.Pi * diameter * diameter }
func diameterFromArea(area float64) float64 { return 2 * math.Sqrt(area/math.Pi) }

// stepDrillEntry is one stage of a step-drilling schedule.
type stepDrillEntry struct {
	Drill          drill
	RemovedPercent float64
}

// pilotRemovedPercent reports what fraction of the total material to
// remove the pilot hole alone already accounts for, or false if the
// final hole size has zero area (never true in practice, but matches
// DRILL.C's own explicit zero-area guard).
func pilotRemovedPercent(drills []drill, finalIdx, pilotIdx int) (percent float64, ok bool) {
	area := circleArea(drills[pilotIdx].SizeIn)
	if area == 0 {
		return 0, false
	}
	return 100 * area / circleArea(drills[finalIdx].SizeIn), true
}

// stepDrillingSchedule computes the sequence of drills to step up
// through from a pilot hole to a final hole size, each step removing
// approximately stepPercent of the total remaining material, snapped
// to the nearest actually-available drill size after each step —
// matching DRILL.C's own step() exactly, including compounding each
// step from the previous step's *snapped* area rather than an ideal
// geometric progression.
func stepDrillingSchedule(drills []drill, finalIdx, pilotIdx int, stepPercent float64, excludeMetric bool) []stepDrillEntry {
	targetArea := circleArea(drills[finalIdx].SizeIn)
	area := circleArea(drills[pilotIdx].SizeIn)
	step := 0.01 * stepPercent

	if area < step*targetArea {
		area = step * targetArea
	} else {
		area += step * targetArea
	}
	if area > targetArea {
		area = targetArea
	}

	var entries []stepDrillEntry
	for {
		k := findDrill(drills, diameterFromArea(area), excludeMetric)
		area = circleArea(drills[k].SizeIn)
		entries = append(entries, stepDrillEntry{Drill: drills[k], RemovedPercent: 100 * area / targetArea})
		if k == finalIdx {
			break
		}
		area += step * targetArea
		if area > targetArea {
			area = targetArea
		}
	}
	return entries
}

func printMenu() {
	fmt.Println("D    find Drill designation given hole size")
	fmt.Println("S    find hole Size given drill designation")
	fmt.Println("T    find Tapdrill for any tap and dot")
	fmt.Println("F    find tapdrill for thread-Forming tap")
	fmt.Println("X    step drilling calculations")
	fmt.Println("H    display this Help/Menu")
	fmt.Println("M    display this Help/Menu")
	fmt.Println("Q    Quit (Esc also)")
}

func promptSize(p *promptio.Prompter, label string) float64 {
	line := p.Line(fmt.Sprintf("{[d].dd or [d-]d/d (e.g. 1.5 or 1-1/2)}  %s ? ", label))
	return parseSize(line)
}

func main() {
	ctx := context.Background()
	drills, err := loadDrills(ctx)
	if err != nil {
		fmt.Fprintln(os.Stderr, "drill:", err)
		os.Exit(1)
	}

	prompter, err := promptio.New(os.Stdin, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, "drill:", err)
		os.Exit(1)
	}

	fmt.Printf("number of drills read from data file = %d\n", len(drills))
	printMenu()

	for {
		op := strings.ToLower(strings.TrimSpace(prompter.Line("(D,S,T,F,X,H,M,Q) ? ")))
		if op == "" {
			continue
		}
		switch op[0] {
		case 'd':
			size := promptSize(prompter, "Hole Size")
			k := findDrill(drills, size, false)
			printReport(buildReport(drills, k, size, 0, 0, 1, false))
		case 's':
			runFindHoleSize(prompter, drills)
		case 't':
			runTapDrill(prompter, drills, false)
		case 'f':
			runTapDrill(prompter, drills, true)
		case 'x':
			runStepDrilling(prompter, drills)
		case 'h', 'm':
			printMenu()
		case 'q':
			return
		default:
			fmt.Println("NOT A VALID OPTION")
		}
	}
}

func runFindHoleSize(prompter *promptio.Prompter, drills []drill) {
	designation := prompter.Line("Drill Designation (e.g. 6, F, 3.10 mm, 3/8, 1-1/4) ? ")
	k, ok := findByDesignation(drills, designation)
	if !ok {
		fmt.Println("NO MATCH FOUND")
		return
	}
	printReport(buildReport(drills, k, drills[k].SizeIn, 0, 0, 1, false))
}

func runTapDrill(prompter *promptio.Prompter, drills []drill, forming bool) {
	metric := strings.EqualFold(strings.TrimSpace(prompter.Line("[I]nch [Default] or (M)etric Tap ? ")), "m")

	var td, pitchPerIn, tapDiamMM, pitchMM, aa float64

	if metric {
		if forming {
			fmt.Println("For metric taps the approximate tap drill size can be found from:")
			fmt.Println("tapdrill diameter (mm) = tap diameter (mm) - pitch of tap (mm)")
		}
		tapDiamMM = prompter.Float("Tap diameter (mm)", 6)
		td = tapDiamMM / mmPerInch
		pitchMM = prompter.Float("Pitch of tap (mm/thread)", 1)
		pitchPerIn = mmPerInch / pitchMM
		iso := strings.HasPrefix(strings.ToLower(strings.TrimSpace(prompter.Line("[S]tandard or (I)SO Thread Form ? "))), "i")
		if iso {
			aa = fiveEighthsTan60
		} else {
			aa = sixEighthsTan60
		}
	} else {
		fmt.Println("Diameters of numbered taps:")
		fmt.Println("#0000 (160) = .021  #000 (120) = .034  #00 (90) = .047")
		fmt.Println("#0 = .060  #1 = .073  #2 = .085  #3  = .099  #4  = .112")
		fmt.Println("#5 = .125  #6 = .138  #8 = .164  #10 = .190  #12 = .216")
		fmt.Println("(for number sizes >= 0, you can enter number with a leading minus)")
		td = promptSize(prompter, "Tap Diameter (in)")
		if td <= 0 {
			td = 0.060 - td*0.013
		}
		pitchPerIn = promptSize(prompter, "Pitch of Tap (threads/in)")
	}

	if !forming && !metric {
		unified := strings.HasPrefix(strings.ToLower(strings.TrimSpace(prompter.Line("American [N]ational or (U)nified Thread Form ? "))), "u")
		if unified {
			aa = fiveEighthsTan60
		} else {
			aa = sixEighthsTan60
		}
	}
	if !forming {
		fmt.Println("Recommended Depth of Thread:")
		fmt.Println("MATERIAL                                          % DOT")
		fmt.Println()
		fmt.Println("MILD AND UNTREATED STEELS                         60-65")
		fmt.Println("HIGH CARBON STEEL                                 50")
		fmt.Println("HIGH SPEED STEEL                                  55")
		fmt.Println("STAINLESS STEEL                                   50")
		fmt.Println("FREE CUTTING STAINLESS STEEL                      60")
		fmt.Println("CAST IRON                                         70-75")
		fmt.Println("WROUGHT ALUMINUM                                  65")
		fmt.Println("CAST ALUMINUM                                     75")
		fmt.Println("WROUGHT COPPER                                    60")
		fmt.Println("FREE CUTTING YELLOW BRASS                         70")
		fmt.Println("DRAWN BRASS                                       65")
		fmt.Println("MANGANESE BRONZE                                  55")
		fmt.Println("MONEL METAL                                       55-60")
		fmt.Println("NICKEL SILVER (GERMAN SILVER)                     50-60")
	}
	dot := prompter.Float("Depth of thread desired (%)", 75)

	fmt.Printf("\nTap Diameter = %.4f in = %.4f mm\n", td, td*mmPerInch)
	fmt.Printf("Tap Pitch = %.4f threads/in = %.4f in/thread\n", pitchPerIn, 1/pitchPerIn)
	fmt.Printf("          = %.4f threads/mm = %.4f mm/thread\n", pitchPerIn/mmPerInch, mmPerInch/pitchPerIn)

	var dtd float64
	if !forming {
		fmt.Printf("Thread form constant used = %.4f\n", aa)
		dtd = td - 0.01*dot*aa/pitchPerIn
	} else {
		dtd = td - 0.0068*dot/pitchPerIn
		if metric {
			dtd = (tapDiamMM - dot*pitchMM/147.06) / mmPerInch
		}
	}
	fmt.Printf("TAP DRILL DIAMETER FOR %.0f%% DOT = %.4f in = %.4f mm\n", dot, dtd, dtd*mmPerInch)

	k := findDrill(drills, dtd, false)
	printReport(buildReport(drills, k, dtd, td, pitchPerIn, aa, !forming))
}

func runStepDrilling(prompter *promptio.Prompter, drills []drill) {
	excludeMetric := strings.EqualFold(strings.TrimSpace(prompter.Line("Include metric drills in search [Y]/N? ")), "n")

	finalSize := promptSize(prompter, "Final Hole Size")
	kf := findDrill(drills, finalSize, excludeMetric)
	if drills[kf].SizeIn != finalSize {
		fmt.Println("Final hole size drill NOT in data file.  Will use closest value:")
		fmt.Printf("CLOSEST DRILL IS (%s) WITH SIZE %.4f (%+.4f)\n", drills[kf].Name, drills[kf].SizeIn, drills[kf].SizeIn-finalSize)
	}

	pilotSize := promptSize(prompter, "Pilot Hole Size")
	kp := findDrill(drills, pilotSize, excludeMetric)

	stepPercent := prompter.Float("Percentage of material to remove with each step", 20)

	fmt.Println()
	if percent, ok := pilotRemovedPercent(drills, kf, kp); ok {
		fmt.Printf("Pilot drill removes %.0f%% of material to be removed\n", percent)
	}

	for _, entry := range stepDrillingSchedule(drills, kf, kp, stepPercent, excludeMetric) {
		fmt.Printf("Drill (%s) with size %.4f removes %.0f%%\n", entry.Drill.Name, entry.Drill.SizeIn, entry.RemovedPercent)
	}
}
