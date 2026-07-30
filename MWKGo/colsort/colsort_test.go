package main

import (
	"os"
	"strconv"
	"strings"
	"testing"
)

func openExample(t *testing.T) *table {
	t.Helper()
	f, err := os.Open("testdata/example.dat")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	tbl, err := loadTable(f)
	if err != nil {
		t.Fatal(err)
	}
	return tbl
}

func TestLoadTable_ExampleHeader(t *testing.T) {
	tbl := openExample(t)
	if len(tbl.columnTypes) != 2 {
		t.Fatalf("columns = %d, want 2", len(tbl.columnTypes))
	}
	if tbl.columnTypes[0] != alphabetic || tbl.columnTypes[1] != numeric {
		t.Errorf("columnTypes = %v, want [alphabetic numeric]", tbl.columnTypes)
	}
	if tbl.sortColumn != 2 {
		t.Errorf("sortColumn = %d, want 2", tbl.sortColumn)
	}
	if tbl.sortIncreasing {
		t.Error("sortIncreasing = true, want false (COLSORT.DAT specifies 'd')")
	}
	if len(tbl.rows) != 30 {
		t.Errorf("rows = %d, want 30", len(tbl.rows))
	}
}

func TestSortRows_DecreasingNumericMatchesKnownExtremes(t *testing.T) {
	tbl := openExample(t)
	sortRows(tbl)

	if got := tbl.rows[0][0]; got != "Zinc" {
		t.Errorf("first row after decreasing sort on density = %q, want Zinc (17.0, the highest density)", got)
	}
	last := tbl.rows[len(tbl.rows)-1]
	if last[0] != "Carbon" {
		t.Errorf("last row after decreasing sort on density = %q, want Carbon (1.4, the lowest density)", last[0])
	}
}

func TestSortRows_ResultIsFullyOrdered(t *testing.T) {
	tbl := openExample(t)
	original := len(tbl.rows)
	sortRows(tbl)

	if len(tbl.rows) != original {
		t.Fatalf("sort changed row count from %d to %d", original, len(tbl.rows))
	}
	for i := 1; i < len(tbl.rows); i++ {
		prev, err := strconv.ParseFloat(tbl.rows[i-1][1], 64)
		if err != nil {
			t.Fatal(err)
		}
		cur, err := strconv.ParseFloat(tbl.rows[i][1], 64)
		if err != nil {
			t.Fatal(err)
		}
		if prev < cur {
			t.Errorf("row %d (%v) < row %d (%v): not decreasing", i-1, prev, i, cur)
		}
	}
}

func TestSortRows_AlphabeticIsCaseInsensitiveAndStable(t *testing.T) {
	tbl := &table{
		columnTypes:    []columnType{alphabetic, numeric},
		sortColumn:     1,
		sortIncreasing: true,
		rows: [][]string{
			{"banana", "1"},
			{"Apple", "2"},
			{"apple", "3"}, // same key as "Apple" case-insensitively; original order must be kept
			{"cherry", "4"},
		},
	}
	sortRows(tbl)

	got := []string{tbl.rows[0][0], tbl.rows[1][0], tbl.rows[2][0], tbl.rows[3][0]}
	want := []string{"Apple", "apple", "banana", "cherry"}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("sorted order = %v, want %v", got, want)
		}
	}
	// Confirm stability: "Apple" (originally before "apple") stays first among the tie.
	if tbl.rows[0][1] != "2" || tbl.rows[1][1] != "3" {
		t.Errorf("tie not stable: got rows %v, %v", tbl.rows[0], tbl.rows[1])
	}
}

func TestLoadTable_RejectsRowWithWrongFieldCount(t *testing.T) {
	bad := `STARTOFDATA
2
a
n
1
i
onlyonefield
ENDOFDATA
`
	_, err := loadTable(strings.NewReader(bad))
	if err == nil {
		t.Fatal("expected an error for a data row with the wrong field count")
	}
}
