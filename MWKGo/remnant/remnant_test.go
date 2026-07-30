package main

import (
	"math"
	"os"
	"testing"
)

const eps = 1e-6

func almostEqual(a, b, tol float64) bool { return math.Abs(a-b) < tol }

func loadExample(t *testing.T) remnantData {
	t.Helper()
	f, err := os.Open("testdata/example.dat")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	data, err := loadRemnant(f)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func TestLoadRemnant_MatchesShippedExample(t *testing.T) {
	data := loadExample(t)
	if data.Kerf != 0.03 {
		t.Errorf("kerf = %v, want 0.03", data.Kerf)
	}
	// 6 remnants of 4.0, 3 of 5.0, 1 of 7.0, 1 of 10.0 = 11 total,
	// sorted descending.
	if len(data.Remnants) != 11 {
		t.Fatalf("got %d remnants, want 11", len(data.Remnants))
	}
	want := []float64{10, 7, 5, 5, 5, 4, 4, 4, 4, 4, 4}
	for i, w := range want {
		if !almostEqual(data.Remnants[i], w, eps) {
			t.Errorf("remnant %d = %v, want %v", i, data.Remnants[i], w)
		}
	}
	if len(data.Pieces) != 7 {
		t.Fatalf("got %d piece requirements, want 7", len(data.Pieces))
	}
	if data.Pieces[0].Size != 2.55 || data.Pieces[0].Num != 5 {
		t.Errorf("first piece = %+v, want {Num:5 Size:2.55}", data.Pieces[0])
	}
}

// TestAssignRemnant_EveryPieceAssignedOnShippedExample confirms the
// shipped example is fully solvable (no unassigned pieces) and that
// the accounting reconciles exactly: kerf waste + drop waste must
// equal the theoretical possible waste (total remnant length minus
// total required length) whenever every piece finds a home.
func TestAssignRemnant_EveryPieceAssignedOnShippedExample(t *testing.T) {
	data := loadExample(t)

	var tbars, reqd, kreqd float64
	for _, b := range data.Remnants {
		tbars += b
	}
	var piecesWithKerf []float64
	for _, p := range data.Pieces {
		reqd += p.Size * float64(p.Num)
		kreqd += (p.Size + data.Kerf) * float64(p.Num)
		for i := 0; i < p.Num; i++ {
			piecesWithKerf = append(piecesWithKerf, p.Size+data.Kerf)
		}
	}

	stock, unassigned := assignRemnant(data.Remnants, piecesWithKerf, data.Kerf)
	if len(unassigned) != 0 {
		t.Fatalf("unassigned = %v, want none", unassigned)
	}

	var kwaste, dwaste float64
	var cutCount int
	for _, s := range stock {
		kwaste += float64(len(s.Sizes)) * data.Kerf
		dwaste += s.Drop
		cutCount += len(s.Sizes)
	}
	if cutCount != len(piecesWithKerf) {
		t.Errorf("pieces cut = %d, want %d", cutCount, len(piecesWithKerf))
	}
	if !almostEqual(kwaste+dwaste, tbars-reqd, eps) {
		t.Errorf("kerf+drop waste = %v, want %v (tbars-reqd)", kwaste+dwaste, tbars-reqd)
	}
}

// TestAssignRemnant_NeverExceedsRemnantLength confirms the basic
// correctness invariant: no remnant's assigned cuts (including kerf)
// ever exceed its own original length.
func TestAssignRemnant_NeverExceedsRemnantLength(t *testing.T) {
	data := loadExample(t)
	var piecesWithKerf []float64
	for _, p := range data.Pieces {
		for i := 0; i < p.Num; i++ {
			piecesWithKerf = append(piecesWithKerf, p.Size+data.Kerf)
		}
	}
	stock, _ := assignRemnant(data.Remnants, piecesWithKerf, data.Kerf)
	for i, s := range stock {
		var used float64
		for _, size := range s.Sizes {
			used += size
		}
		if !almostEqual(used+s.Drop, s.Bar, eps) {
			t.Errorf("remnant %d: cuts (%v) + drop (%v) = %v, want %v", i, used, s.Drop, used+s.Drop, s.Bar)
		}
		if s.Drop < -eps {
			t.Errorf("remnant %d: drop = %v, want >= 0", i, s.Drop)
		}
	}
}

// TestAssignRemnant_UnfittablePieceIsReported confirms a piece larger
// than every available remnant is reported as unassigned (with the
// kerf allowance removed again) rather than silently dropped or
// causing an out-of-range panic.
func TestAssignRemnant_UnfittablePieceIsReported(t *testing.T) {
	bars := []float64{5, 3}
	kerf := 0.1
	piecesWithKerf := []float64{10 + kerf, 2 + kerf}
	stock, unassigned := assignRemnant(bars, piecesWithKerf, kerf)

	if len(unassigned) != 1 || !almostEqual(unassigned[0], 10, eps) {
		t.Fatalf("unassigned = %v, want [10]", unassigned)
	}
	var cutCount int
	for _, s := range stock {
		cutCount += len(s.Sizes)
	}
	if cutCount != 1 {
		t.Errorf("pieces cut = %d, want 1 (the piece that does fit)", cutCount)
	}
}

// TestAssignRemnant_PrefersExactFitFoundLater is the key documented
// difference from cutlist's own best-fit refinement: a later remnant
// offering an exact zero-waste fit must be preferred over an earlier,
// looser-fitting one already picked.
func TestAssignRemnant_PrefersExactFitFoundLater(t *testing.T) {
	bars := []float64{10, 6} // descending, as loadRemnant would sort them
	kerf := 0.0
	piecesWithKerf := []float64{6}
	stock, unassigned := assignRemnant(bars, piecesWithKerf, kerf)
	if len(unassigned) != 0 {
		t.Fatalf("unassigned = %v, want none", unassigned)
	}
	if len(stock[1].Sizes) != 1 {
		t.Errorf("piece should have been cut from the exact-fit 6-length remnant (index 1), got sizes %v/%v", stock[0].Sizes, stock[1].Sizes)
	}
}
