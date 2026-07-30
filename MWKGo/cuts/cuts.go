// cuts solves the one-dimensional stock-cutting problem for a single
// standard bar length: given a list of piece sizes and how many of
// each are needed, work out how to cut them from as few standard-
// length bars as possible, and print the cutting schedule.
//
// For each bar, cuts searches every combination of still-needed
// pieces that fits (bounded — at most as many of a size as both fit
// in one bar and are still needed), preferring the combination with
// the least waste, and among equally-wasteful combinations, the one
// favoring larger pieces (this project's own conventional "greedier
// combinations get better tie-break scores"). It repeats, one bar at
// a time, until every piece is accounted for, printing a compressed
// schedule that groups consecutive bars cut identically.
//
// A -zero-waste flag (matching the original's own "Search for zero
// waste solutions first (Y/[N])?" prompt, default off) restricts the
// very first search per bar to zero-waste combinations only, falling
// back to the ordinary least-waste search if none exist. CUTS.TXT
// itself documents a case where this flag finds a strictly better
// overall schedule than the default, and also states plainly that
// there is no cure-all: sometimes the default finds a better answer,
// and neither approach always reaches the true theoretical optimum
// (that's what MWKGo/cutlist is for — see its own doc comment).
//
// Converted from CUTS.C (M. W. Klotz), WorkshopUtilities/cuts.zip.
// Despite this program's original grouping as a "Tier 4: stateful
// program" needing the full atomic-write pillar (per
// ai/plans/c-to-go-conversion-plan.md), it turns out to read its own
// .DAT file and print to stdout exactly like the rest of this
// project's batch — nothing here is stateful, and no output file is
// ever written. The original's "debug" mode (triggered by any
// command-line argument, printing internal search diagnostics and
// stopping after a single bar) is developer scaffolding, not a
// user-facing feature, and is dropped entirely, the same way this
// project has dropped other purely-diagnostic paths elsewhere.
package main

import (
	"flag"
	"fmt"
	"io"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"

	"mwkgo/internal/legacydat"
)

// precision is CUTS.C's own precision: two lengths within this
// tolerance of each other are treated as equal, matching floating-
// point drift in the original's own real (32-bit float) arithmetic.
const precision = 1.0e-3

// diff is CUTS.C's own diff(): -1, 0, or 1 according to whether a is
// less than, within precision of, or greater than b.
func diff(a, b float64) int {
	t := a - b
	switch {
	case math.Abs(t) <= precision:
		return 0
	case t > 0:
		return 1
	default:
		return -1
	}
}

// pieceReq is one line of the data file: how many pieces of a given
// size are still needed. Num is mutated in place as bars are cut,
// the same way CUTS.C's own num[] is.
type pieceReq struct {
	Num  int
	Size float64
}

// loadCuts parses CUTS.DAT's own format: a standard bar length,
// followed by comma-separated count,size lines, ending at ENDOFDATA
// (there is no STARTOFDATA marker — the data region starts at the
// first non-comment line, the same convention GAGE.C's own .DAT file
// uses). Pieces are sorted descending by size — CUTS.C's own
// requirement ("NOTE: order arrays by decreasing size"), enforced by
// its own bubble sort rather than left to the data file's order.
func loadCuts(r io.Reader) (bar float64, pieces []pieceReq, err error) {
	rows, err := legacydat.Rows(r, legacydat.Options{EndMarker: "ENDOFDATA"})
	if err != nil {
		return 0, nil, fmt.Errorf("scan data: %w", err)
	}
	if len(rows) == 0 {
		return 0, nil, fmt.Errorf("data file has no bar length")
	}

	barFields := legacydat.Fields(rows[0], "\t;")
	bar, err = strconv.ParseFloat(strings.TrimSpace(barFields[0]), 64)
	if err != nil {
		return 0, nil, fmt.Errorf("line 1: bar length %q: %w", barFields[0], err)
	}

	for i, row := range rows[1:] {
		fields := legacydat.Fields(row, ",;")
		if len(fields) < 2 {
			return 0, nil, fmt.Errorf("line %d: %q does not split into a count and size", i+2, row)
		}
		num, err := strconv.Atoi(strings.TrimSpace(fields[0]))
		if err != nil {
			return 0, nil, fmt.Errorf("line %d: count %q: %w", i+2, fields[0], err)
		}
		size, err := strconv.ParseFloat(strings.TrimSpace(fields[1]), 64)
		if err != nil {
			return 0, nil, fmt.Errorf("line %d: size %q: %w", i+2, fields[1], err)
		}
		pieces = append(pieces, pieceReq{Num: num, Size: size})
	}

	sort.SliceStable(pieces, func(i, j int) bool { return pieces[i].Size > pieces[j].Size })
	return bar, pieces, nil
}

// bestCombination is CUTS.C's own greedy(): an exhaustive (but
// bounded) search, across a single bar, for the combination of
// still-needed pieces with the least waste — and, among equally
// wasteful combinations, the one weighted most heavily toward larger
// pieces, since cutting the big pieces first leaves more flexibility
// for the smaller ones later.
//
// If restrictZeroWaste is true, only combinations with exactly zero
// waste are considered; if none exists, the returned waste is
// +Inf (matching the original's own large sentinel value, checked by
// its caller to decide whether to fall back to an unrestricted
// search).
func bestCombination(bar float64, pieces []pieceReq, restrictZeroWaste bool) (waste float64, best []int) {
	np := len(pieces)
	max := make([]int, np)
	x := make([]int, np)
	for i, p := range pieces {
		m := int(math.Ceil(bar / p.Size))
		if p.Num < m {
			m = p.Num
		}
		max[i] = m
	}
	min := pieces[np-1].Size

	waste = math.Inf(1)
	var score float64
	best = make([]int, np)

	// visit scores one full combination (x[0..np-1]) and, if it's an
	// improvement, saves it. The score/waste tie-break logic below is
	// a direct translation of CUTS.C's own decision table (see its
	// lengthy comment on the "waste < min / waste == min / waste >
	// min" cases): waste < min alone still requires a better score to
	// accept, since wasting a bit less material makes no practical
	// difference once the remaining waste is already too small to fit
	// another piece.
	visit := func() {
		var length, s float64
		weight := math.Ldexp(1, np)
		for i := 0; i < np; i++ {
			if x[i] != 0 {
				length += float64(x[i]) * pieces[i].Size
				if diff(length, bar) > 0 {
					return // doesn't fit; skip this combination entirely
				}
				s += float64(x[i]) * weight
			}
			weight *= 0.5
		}
		w := bar - length
		if diff(w, 0) == 0 {
			w = 0
		}
		if restrictZeroWaste && w != 0 {
			return
		}

		var accept bool
		switch {
		case waste < min:
			if w > waste && w >= min {
				accept = false
			} else {
				accept = s > score
			}
		case w < waste:
			accept = true
		case w > waste:
			accept = false
		default:
			accept = s > score
		}
		if accept {
			waste = w
			score = s
			copy(best, x)
		}
	}

	// Enumerates every combination x[0]..x[np-1] in [0,max[i]], with
	// x[0] (the largest piece size) cycling fastest. CUTS.C generates
	// the identical set of combinations via a goto-based odometer
	// (decrementing counts); which order they're visited in doesn't
	// change the result, since the score formula (a sum of powers of
	// two, one per size) gives each distinct combination a unique
	// score, so no two different combinations can ever tie on both
	// waste and score.
	var enumerate func(idx int)
	enumerate = func(idx int) {
		if idx < 0 {
			visit()
			return
		}
		for x[idx] = max[idx]; x[idx] >= 0; x[idx]-- {
			enumerate(idx - 1)
		}
	}
	enumerate(np - 1)

	return waste, best
}

func intSlicesEqual(a, b []int) bool {
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func printGroup(dup int, best []int, waste float64) {
	fmt.Printf("%d x", dup)
	for _, c := range best {
		fmt.Printf("\t%d", c)
	}
	fmt.Printf("\t%-5.2f\n", float64(dup)*waste)
}

func main() {
	dataPath := flag.String("data", "", "path to a CUTS-format data file (see MWKGo/cuts/testdata/)")
	zeroWaste := flag.Bool("zero-waste", false, "search for zero-waste combinations before the ordinary least-waste search")
	flag.Parse()

	if *dataPath == "" {
		fmt.Println("usage: cuts -data <file> [-zero-waste]")
		fmt.Println("see MWKGo/cuts/testdata/ for the expected data format")
		os.Exit(1)
	}

	f, err := os.Open(*dataPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "cuts:", err)
		os.Exit(1)
	}
	bar, pieces, err := loadCuts(f)
	f.Close()
	if err != nil {
		fmt.Fprintln(os.Stderr, "cuts:", err)
		os.Exit(1)
	}

	fmt.Printf("Number of data items read = %d\n", len(pieces))

	var reqd float64
	for _, p := range pieces {
		reqd += p.Size * float64(p.Num)
	}
	oreqd := int(math.Ceil(reqd / bar))
	owaste := float64(oreqd)*bar - reqd

	fmt.Printf("standard length = %.6g\n", bar)
	if pieces[0].Size > bar {
		fmt.Fprintf(os.Stderr, "cuts: largest piece %.6g larger than bar %.6g\n", pieces[0].Size, bar)
		os.Exit(1)
	}
	fmt.Printf("theoretical minimum waste possible = %.6g\n", owaste)
	fmt.Printf("theoretical minimum standard lengths possible = %d\n\n", oreqd)
	for _, p := range pieces {
		fmt.Printf("\t%d", p.Num)
	}
	fmt.Println()
	for _, p := range pieces {
		fmt.Printf("\t%-5.2f", p.Size)
	}
	fmt.Println()
	fmt.Println()

	var twaste float64
	var nbar, dup int
	var last []int
	var lwaste float64

	for {
		waste, best := bestCombination(bar, pieces, *zeroWaste)
		if *zeroWaste && math.IsInf(waste, 1) {
			waste, best = bestCombination(bar, pieces, false)
		}

		nbar++
		twaste += waste
		more := false
		for i := range pieces {
			pieces[i].Num -= best[i]
			if pieces[i].Num > 0 {
				more = true
			}
		}
		if waste == bar {
			more = false // prevents a runaway loop if nothing more fits
		}

		if nbar > 1 && !intSlicesEqual(best, last) {
			printGroup(dup, last, lwaste)
			dup = 1
		} else {
			dup++
		}
		last = append([]int(nil), best...)
		lwaste = waste

		if !more {
			break
		}
	}
	printGroup(dup, last, lwaste)

	fmt.Printf("\nactual waste = %.6g\n", twaste)
	fmt.Printf("actual standard lengths = %d\n", nbar)
}
