// colsort sorts a table of comma-delimited rows (each column marked
// alphabetic or numeric) by one chosen column, increasing or
// decreasing.
//
// The table, its column types, and which column to sort by are all
// specific to one sorting job — not universal reference data, and not
// reusable equipment configuration — so, like calibrat, vrev, and
// simul in an earlier conversion group, this program reads its input
// fresh from a file named on the command line each run, in the same
// STARTOFDATA/ENDOFDATA text format the original used; see
// ai/plans/c-to-go-conversion-plan.md, "Data-file strategy for
// Tier 2".
//
// Converted from COLSORT.C (M. W. Klotz), Misc/colsort.
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

// columnType matches COLSORT.C's own coltype convention.
type columnType int

const (
	numeric columnType = iota
	alphabetic
)

// table is a parsed COLSORT.DAT-format job: its column types, which
// column to sort on, the sort direction, and the row data itself
// (rows[i][j] is row i's column j value, as text).
type table struct {
	columnTypes    []columnType
	sortColumn     int // 1-based, as in the original
	sortIncreasing bool
	rows           [][]string
}

// maxColumns matches COLSORT.C's own MAXCOL fixed array bound.
const maxColumns = 8

// loadTable parses a COLSORT.DAT-format job from r: a column count, one
// a/n line per column, the 1-based column to sort on, an i/d sort
// direction, then the data rows themselves.
func loadTable(r io.Reader) (*table, error) {
	rows, err := legacydat.Rows(r, legacydat.Options{StartMarker: "STARTOFDATA", EndMarker: "ENDOFDATA"})
	if err != nil {
		return nil, fmt.Errorf("scan data: %w", err)
	}
	if len(rows) == 0 {
		return nil, fmt.Errorf("no data found")
	}

	pos := 0
	next := func() (string, error) {
		if pos >= len(rows) {
			return "", fmt.Errorf("data file ends before job header is complete")
		}
		line := rows[pos]
		pos++
		return line, nil
	}

	ncolLine, err := next()
	if err != nil {
		return nil, err
	}
	ncolFields := legacydat.Fields(ncolLine, ",\t;")
	if len(ncolFields) == 0 {
		return nil, fmt.Errorf("column count line %q is empty", ncolLine)
	}
	ncol, err := strconv.Atoi(strings.TrimSpace(ncolFields[0]))
	if err != nil {
		return nil, fmt.Errorf("column count %q: %w", ncolFields[0], err)
	}
	if ncol < 1 || ncol > maxColumns {
		return nil, fmt.Errorf("column count %d out of range 1..%d", ncol, maxColumns)
	}

	columnTypes := make([]columnType, ncol)
	for c := 0; c < ncol; c++ {
		line, err := next()
		if err != nil {
			return nil, err
		}
		if strings.HasPrefix(line, "a") {
			columnTypes[c] = alphabetic
		} else {
			columnTypes[c] = numeric
		}
	}

	sortColLine, err := next()
	if err != nil {
		return nil, err
	}
	sortColFields := legacydat.Fields(sortColLine, ",\t;")
	if len(sortColFields) == 0 {
		return nil, fmt.Errorf("sort column line %q is empty", sortColLine)
	}
	sortColumn, err := strconv.Atoi(strings.TrimSpace(sortColFields[0]))
	if err != nil {
		return nil, fmt.Errorf("sort column %q: %w", sortColFields[0], err)
	}
	if sortColumn < 1 || sortColumn > ncol {
		return nil, fmt.Errorf("sort column %d out of range 1..%d", sortColumn, ncol)
	}

	sortTypeLine, err := next()
	if err != nil {
		return nil, err
	}
	sortIncreasing := strings.HasPrefix(sortTypeLine, "i")

	var dataRows [][]string
	for pos < len(rows) {
		line, _ := next()
		fields := legacydat.Fields(line, ",\t;")
		if len(fields) != ncol {
			return nil, fmt.Errorf("data row %q has %d field(s), want %d", line, len(fields), ncol)
		}
		for i, f := range fields {
			fields[i] = strings.TrimSpace(f)
		}
		dataRows = append(dataRows, fields)
	}

	return &table{
		columnTypes:    columnTypes,
		sortColumn:     sortColumn,
		sortIncreasing: sortIncreasing,
		rows:           dataRows,
	}, nil
}

// sortRows sorts t.rows in place on t.sortColumn, matching COLSORT.C's
// own comparison exactly: alphabetic columns compare case-insensitively
// (with a shorter, otherwise-matching prefix sorting first), numeric
// columns compare by parsed value. The original's own bubble sort is
// already a correct, stable sort (unlike diffthrd's, which never
// finished); sort.SliceStable reproduces the same total order without
// reproducing the bubble-sort mechanics.
func sortRows(t *table) {
	col := t.sortColumn - 1
	var less func(a, b string) bool
	if t.columnTypes[col] == alphabetic {
		less = func(a, b string) bool { return strings.ToLower(a) < strings.ToLower(b) }
	} else {
		less = func(a, b string) bool {
			av, _ := strconv.ParseFloat(a, 64)
			bv, _ := strconv.ParseFloat(b, 64)
			return av < bv
		}
	}

	sort.SliceStable(t.rows, func(i, j int) bool {
		a, b := t.rows[i][col], t.rows[j][col]
		if t.sortIncreasing {
			return less(a, b)
		}
		return less(b, a)
	})
}

func printRows(w io.Writer, rows [][]string) {
	for _, row := range rows {
		fmt.Fprintln(w, strings.Join(row, ",")+",")
	}
}

func main() {
	dataPath := flag.String("data", "", "path to a COLSORT-format data file (see MWKGo/colsort/testdata/example.dat)")
	flag.Parse()

	if *dataPath == "" {
		fmt.Println("usage: colsort -data <file>")
		fmt.Println("see MWKGo/colsort/testdata/example.dat for the expected format")
		os.Exit(1)
	}

	f, err := os.Open(*dataPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "colsort:", err)
		os.Exit(1)
	}
	t, err := loadTable(f)
	f.Close()
	if err != nil {
		fmt.Fprintln(os.Stderr, "colsort:", err)
		os.Exit(1)
	}

	fmt.Println("MULTI-COLUMN SORTING")
	fmt.Println()
	fmt.Printf("Number of columns in a data line = %d\n", len(t.columnTypes))
	fmt.Printf("Number of data lines read = %d\n", len(t.rows))

	fmt.Println("\nData before sorting:")
	fmt.Println()
	printRows(os.Stdout, t.rows)

	sortRows(t)

	direction := "decreasing"
	if t.sortIncreasing {
		direction = "increasing"
	}
	kind := "numerical"
	if t.columnTypes[t.sortColumn-1] == alphabetic {
		kind = "alphabetic"
	}
	fmt.Printf("\nData after %s %s sort on column %d:\n\n", direction, kind, t.sortColumn)
	printRows(os.Stdout, t.rows)
}
