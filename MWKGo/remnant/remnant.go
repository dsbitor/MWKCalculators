// remnant solves the same one-dimensional stock-cutting problem as
// cuts and cutlist, but for a heterogeneous collection of leftover
// remnant lengths rather than one uniform standard bar length, and
// accounts for a saw kerf (material lost to the cut itself) on every
// piece. Per CUTS.TXT, it was written for a user who had a pile of
// odd-length leftover stock rather than fresh standard bars.
//
// remnant reuses cutlist's best-fit heuristic (Mike Graham's: for
// each needed piece, largest first, cut it from whichever available
// remnant currently has the smallest sufficient remaining room),
// generalized so every remnant is available from the start rather
// than opened lazily one at a time — since, unlike cuts/cutlist's
// identical standard bars, which remnant a piece comes from matters.
//
// Converted from REMNANT.C (M. W. Klotz), WorkshopUtilities/cuts.zip.
// Its own best-fit refinement step accepts an exact zero-waste fit
// found later in the scan (diff(candidate.drop,piece) >= 0) — unlike
// CUTLIST.C's own refinement step, which requires a strictly tighter
// fit (drop > piece) and so can never prefer a later exact fit over
// an earlier looser one it already picked. remnant's own version does
// not carry that asymmetry, and this conversion preserves that
// difference rather than treating the two as identical algorithms.
//
// REMNANT.C's own "can't be assigned" check compares a stock array
// index against MAXPIECESPERLENGTH (a limit on pieces *per* length,
// not on the *number* of lengths) — harmless there only because that
// constant and MAXLENGTHS (the one actually meant) happen to both be
// defined as 100. This conversion has no such fixed-size bound to
// misapply in the first place (its stock list is a plain slice), so
// it simply detects "no remnant has enough room left" directly,
// which is what the original check was clearly meant to test.
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

// precision matches REMNANT.C's own precision: two lengths within
// this tolerance are treated as equal.
const precision = 1.0e-3

// diff is REMNANT.C's own diff().
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

// pieceReq is one line of the "pieces needed" section: a count and a
// (kerf-free) size.
type pieceReq struct {
	Num  int
	Size float64
}

// remnantData is everything REMNANT.DAT provides.
type remnantData struct {
	Kerf     float64
	Remnants []float64 // one entry per physical remnant, descending by size
	Pieces   []pieceReq
}

// loadRemnant parses REMNANT.DAT's own format: a saw kerf width, then
// comma-separated count,length remnant lines, a "0,0" separator line,
// then comma-separated count,size piece-needed lines, ending at
// ENDOFDATA (there is no STARTOFDATA marker, the same convention
// cuts's own .DAT file uses).
func loadRemnant(r io.Reader) (remnantData, error) {
	rows, err := legacydat.Rows(r, legacydat.Options{EndMarker: "ENDOFDATA"})
	if err != nil {
		return remnantData{}, fmt.Errorf("scan data: %w", err)
	}
	if len(rows) == 0 {
		return remnantData{}, fmt.Errorf("data file has no kerf width")
	}

	kerfFields := legacydat.Fields(rows[0], ",\t;")
	kerf, err := strconv.ParseFloat(strings.TrimSpace(kerfFields[0]), 64)
	if err != nil {
		return remnantData{}, fmt.Errorf("line 1: kerf %q: %w", kerfFields[0], err)
	}

	var remnants []float64
	var pieces []pieceReq
	inPieces := false
	for i, row := range rows[1:] {
		fields := legacydat.Fields(row, ",;")
		if len(fields) < 2 {
			return remnantData{}, fmt.Errorf("line %d: %q does not split into a count and size", i+2, row)
		}
		count, err := strconv.Atoi(strings.TrimSpace(fields[0]))
		if err != nil {
			return remnantData{}, fmt.Errorf("line %d: count %q: %w", i+2, fields[0], err)
		}
		if !inPieces && count == 0 {
			inPieces = true
			continue
		}
		size, err := strconv.ParseFloat(strings.TrimSpace(fields[1]), 64)
		if err != nil {
			return remnantData{}, fmt.Errorf("line %d: size %q: %w", i+2, fields[1], err)
		}
		if inPieces {
			pieces = append(pieces, pieceReq{Num: count, Size: size})
		} else {
			for j := 0; j < count; j++ {
				remnants = append(remnants, size)
			}
		}
	}

	sort.Sort(sort.Reverse(sort.Float64Slice(remnants)))
	sort.SliceStable(pieces, func(i, j int) bool { return pieces[i].Size > pieces[j].Size })

	return remnantData{Kerf: kerf, Remnants: remnants, Pieces: pieces}, nil
}

// stockBar is one available remnant and what's been cut from it.
type stockBar struct {
	Bar   float64
	Drop  float64
	Sizes []float64 // includes the kerf allowance, same as REMNANT.C's own stock[].size
}

// assignRemnant is REMNANT.C's own cutlist(): every remnant is
// available from the start (unlike cutlist's lazily-opened uniform
// bars), so for each piece (largest first, kerf already added), pick
// whichever remnant has the smallest sufficient remaining room,
// scanning every remnant rather than only previously-used ones.
// Pieces that fit no remnant at all are returned in unassigned
// (with the kerf allowance removed again, matching what the original
// prints for this case).
func assignRemnant(bars []float64, piecesWithKerf []float64, kerf float64) (stock []stockBar, unassigned []float64) {
	stock = make([]stockBar, len(bars))
	for i, b := range bars {
		stock[i] = stockBar{Bar: b, Drop: b}
	}

	for _, piece := range piecesWithKerf {
		j := 0
		for j < len(stock) && diff(stock[j].Drop, piece) < 0 {
			j++
		}
		if j == len(stock) {
			unassigned = append(unassigned, piece-kerf)
			continue
		}
		best := j
		for k := j + 1; k < len(stock); k++ {
			if diff(stock[best].Drop, stock[k].Drop) > 0 && diff(stock[k].Drop, piece) >= 0 {
				best = k
			}
		}
		stock[best].Drop -= piece
		stock[best].Sizes = append(stock[best].Sizes, piece)
	}
	return stock, unassigned
}

func main() {
	dataPath := flag.String("data", "", "path to a REMNANT-format data file (see MWKGo/remnant/testdata/)")
	flag.Parse()

	fmt.Println("REMNANT CUTTING LIST")
	fmt.Println()

	if *dataPath == "" {
		fmt.Println("usage: remnant -data <file>")
		fmt.Println("see MWKGo/remnant/testdata/ for the expected data format")
		os.Exit(1)
	}

	f, err := os.Open(*dataPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "remnant:", err)
		os.Exit(1)
	}
	data, err := loadRemnant(f)
	f.Close()
	if err != nil {
		fmt.Fprintln(os.Stderr, "remnant:", err)
		os.Exit(1)
	}

	var tbars, reqd, kreqd float64
	for _, b := range data.Remnants {
		tbars += b
	}
	for _, p := range data.Pieces {
		reqd += p.Size * float64(p.Num)
		kreqd += (p.Size + data.Kerf) * float64(p.Num)
	}

	fmt.Printf("saw kerf = %.6g\n", data.Kerf)
	fmt.Print("remnants = ")
	for _, b := range data.Remnants {
		fmt.Printf("%.6g, ", b)
	}
	fmt.Println()
	fmt.Printf("total remnant pieces, length = %d, %.6g\n", len(data.Remnants), tbars)
	fmt.Printf("total required pieces length = %.6g\n\n", reqd)

	if data.Pieces[0].Size > data.Remnants[0] {
		fmt.Fprintf(os.Stderr, "remnant: largest piece %.6g larger than largest remnant %.6g\n", data.Pieces[0].Size, data.Remnants[0])
		os.Exit(1)
	}
	if reqd > tbars {
		fmt.Fprintf(os.Stderr, "remnant: total piece length %.6g greater than total remnant length %.6g\n", reqd, tbars)
		os.Exit(1)
	}
	if kreqd > tbars {
		fmt.Fprintf(os.Stderr, "remnant: total piece length including kerf allowance %.6g greater than total remnant length %.6g\n", kreqd, tbars)
		os.Exit(1)
	}
	if data.Pieces[len(data.Pieces)-1].Size > data.Remnants[len(data.Remnants)-1] {
		fmt.Printf("smallest piece %.6g longer than shortest remnant %.6g\n\n", data.Pieces[len(data.Pieces)-1].Size, data.Remnants[len(data.Remnants)-1])
	}

	for _, p := range data.Pieces {
		fmt.Printf("\t%d", p.Num)
	}
	fmt.Println()
	for _, p := range data.Pieces {
		fmt.Printf("\t%-5.2f", p.Size)
	}
	fmt.Println()
	fmt.Println()

	var piecesWithKerf []float64
	var npiece int
	for _, p := range data.Pieces {
		npiece += p.Num
		for j := 0; j < p.Num; j++ {
			piecesWithKerf = append(piecesWithKerf, p.Size+data.Kerf)
		}
	}

	stock, unassigned := assignRemnant(data.Remnants, piecesWithKerf, data.Kerf)
	for _, u := range unassigned {
		fmt.Printf("*** PIECE %.6g CAN'T BE ASSIGNED ***\n", u)
	}

	var tcuts int
	var kwaste, dwaste float64
	for _, s := range stock {
		tcuts += len(s.Sizes)
		kw := float64(len(s.Sizes)) * data.Kerf
		kwaste += kw
		dwaste += s.Drop
		fmt.Printf("Length %.6g - ", s.Bar)
		for _, size := range s.Sizes {
			fmt.Printf("%.6g, ", size-data.Kerf)
		}
		fmt.Printf("(%.6g kerf + %.6g drop = %.6g waste)\n", kw, s.Drop, s.Drop+kw)
	}

	if tcuts < npiece {
		fmt.Printf("\n*** WARNING: %d PIECE(S) NOT ASSIGNED ***\n", npiece-tcuts)
	}

	fmt.Printf("\ntheoretical possible waste = %.6g - %.6g = %.6g\n", tbars, reqd, tbars-reqd)
	fmt.Printf("kerf waste = %.6g\n", kwaste)
	fmt.Printf("drop waste = %.6g\n", dwaste)
	fmt.Printf("kerf + drop waste = %.6g\n", kwaste+dwaste)
}
