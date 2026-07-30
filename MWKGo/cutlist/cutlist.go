// cutlist solves the same one-dimensional stock-cutting problem
// cuts does — cutting a list of needed piece sizes from as few
// standard-length bars as possible — using a different heuristic:
// Mike Graham's "best fit decreasing" approach, contributed to
// address cases where cuts's own greedy search falls short. Per
// CUTS.TXT, sort every needed piece largest first, then for each
// piece in turn, cut it from whichever already-opened bar currently
// has the smallest remaining room that still fits it, opening a new
// bar only when none of the already-opened ones do.
//
// CUTS.TXT quotes the original author calling this approach
// "definitely superior... runs faster, and is just generally a
// better way to do the problem" than his own cuts — but also that
// cuts sometimes still wins, so both are worth trying on a given
// problem (they share the same data file format).
//
// Converted from CUTLIST.C (Mike Graham, Feb 2002),
// WorkshopUtilities/cuts.zip. The original works entirely in
// thousandths of a unit as integers (each parsed length is multiplied
// by 1000 and truncated), specifically to avoid a floating-point drift
// problem its own header comment mentions correcting — this
// conversion keeps that fixed-point representation for the same
// reason, rather than reintroducing the float comparisons cuts itself
// still needs its own diff()/precision helper to work around.
package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"

	"mwkgo/internal/legacydat"
)

// millis is a length in thousandths of the data file's own units —
// CUTLIST.C's own fixed-point representation (multiply by 1000,
// truncate), which sidesteps floating-point equality comparisons
// entirely rather than needing a tolerance like cuts's diff() does.
type millis int64

// toMillis matches CUTLIST.C's own "f*1000" cast to an unsigned
// integer: truncation toward zero, not rounding.
func toMillis(f float64) millis {
	return millis(f * 1000)
}

func (m millis) units() float64 { return float64(m) / 1000 }

// loadCutlist parses CUTLIST.DAT's own format: a standard bar length,
// then comma-separated count,size lines, expanded eagerly into one
// entry per individual piece (matching Pieces[]). Unlike cuts and
// remnant, CUTLIST.C doesn't require an ENDOFDATA marker — it simply
// reads until end of file, acting on ENDOFDATA if present — so
// EndMarker here is satisfied whether or not the file has one.
func loadCutlist(r io.Reader) (stockLength millis, pieces []millis, err error) {
	rows, err := legacydat.Rows(r, legacydat.Options{EndMarker: "ENDOFDATA"})
	if err != nil {
		return 0, nil, fmt.Errorf("scan data: %w", err)
	}
	if len(rows) == 0 {
		return 0, nil, fmt.Errorf("data file has no stock length")
	}

	lengthFields := legacydat.Fields(rows[0], "\t;")
	length, err := strconv.ParseFloat(strings.TrimSpace(lengthFields[0]), 64)
	if err != nil {
		return 0, nil, fmt.Errorf("line 1: stock length %q: %w", lengthFields[0], err)
	}
	stockLength = toMillis(length)

	for i, row := range rows[1:] {
		fields := legacydat.Fields(row, ",;")
		if len(fields) < 2 {
			return 0, nil, fmt.Errorf("line %d: %q does not split into a count and size", i+2, row)
		}
		count, err := strconv.Atoi(strings.TrimSpace(fields[0]))
		if err != nil {
			return 0, nil, fmt.Errorf("line %d: count %q: %w", i+2, fields[0], err)
		}
		size, err := strconv.ParseFloat(strings.TrimSpace(fields[1]), 64)
		if err != nil {
			return 0, nil, fmt.Errorf("line %d: size %q: %w", i+2, fields[1], err)
		}
		m := toMillis(size)
		for j := 0; j < count; j++ {
			pieces = append(pieces, m)
		}
	}

	sort.SliceStable(pieces, func(i, j int) bool { return pieces[i] > pieces[j] })
	return stockLength, pieces, nil
}

// stockBar is one opened standard-length bar and what's been cut
// from it so far.
type stockBar struct {
	Drop  millis
	Sizes []millis
}

// assignCutlist is CUTLIST.C's own cutlist(): for each piece (largest
// first), cut it from whichever already-opened bar has the smallest
// remaining room that still fits it, opening a new bar only if none
// of the already-opened ones do.
//
// The original's two-phase scan — an initial pass across every bar
// slot, including still-unopened ones, accepting the first with
// drop >= piece; then a refinement pass across only already-opened
// bars, replacing the pick only on a strictly tighter fit (drop >
// piece, not >=) — has a subtle asymmetry: a later already-opened bar
// whose remaining room is an exact zero-waste fit is never preferred
// over an earlier, looser-fitting one found first, since the
// refinement pass demands a strict inequality. This conversion
// preserves that quirk exactly (see the two loops below) rather than
// silently tightening it into an always-optimal best-fit, since nine
// times out of ten it makes no difference and it's not this
// conversion's place to second-guess Mike Graham's algorithm.
func assignCutlist(stockLength millis, pieces []millis) []stockBar {
	var stock []stockBar
	for _, piece := range pieces {
		j := 0
		for j < len(stock) && stock[j].Drop < piece {
			j++
		}
		best := j
		if best == len(stock) {
			stock = append(stock, stockBar{Drop: stockLength})
		} else {
			for k := best + 1; k < len(stock); k++ {
				if stock[k].Drop < stock[best].Drop && stock[k].Drop > piece {
					best = k
				}
			}
		}
		stock[best].Drop -= piece
		stock[best].Sizes = append(stock[best].Sizes, piece)
	}
	return stock
}

func main() {
	dataPath := flag.String("data", "", "path to a CUTLIST-format data file (see MWKGo/cutlist/testdata/)")
	flag.Parse()

	fmt.Println("Cutlist 1.0")
	fmt.Println("Public Domain by Mike Graham")

	if *dataPath == "" {
		fmt.Println("usage: cutlist -data <file>")
		fmt.Println("see MWKGo/cutlist/testdata/ for the expected data format")
		os.Exit(1)
	}

	f, err := os.Open(*dataPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "cutlist:", err)
		os.Exit(1)
	}
	stockLength, pieces, err := loadCutlist(f)
	f.Close()
	if err != nil {
		fmt.Fprintln(os.Stderr, "cutlist:", err)
		os.Exit(1)
	}

	stock := assignCutlist(stockLength, pieces)

	var total millis
	for _, p := range pieces {
		total += p
	}
	var waste millis
	for _, s := range stock {
		waste += s.Drop
	}

	fmt.Printf("Number of pieces to cut = %d\n", len(pieces))
	fmt.Printf("Total length of material being cut = %.3f units\n", total.units())
	fmt.Printf("Standard length = %.3f units \n", stockLength.units())
	remainder := total % stockLength
	fmt.Printf("Theoretical minimum waste possible = %.3f units \n", (stockLength - remainder).units())
	fmt.Printf("Theoretical minimum standard lengths possible = %d\n\n", int(total/stockLength)+1)

	i := 0
	for i < len(pieces) {
		j := i
		for j < len(pieces) && pieces[j] == pieces[i] {
			j++
		}
		fmt.Printf("%3d piece(s) %.3f units long.\n", j-i, pieces[i].units())
		i = j
	}

	fmt.Println("\n\nCutlist:")
	for i, s := range stock {
		fmt.Printf("Length #%d - \n", i+1)
		for _, size := range s.Sizes {
			fmt.Printf("   %.3f units\n", size.units())
		}
	}

	fmt.Printf("\n\nActual waste = %.3f\n", waste.units())
	fmt.Printf("Actual standard lengths = %d\n\n", len(stock))
}
