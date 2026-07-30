// simul solves a system of N simultaneous linear equations in N
// unknowns by Gauss-Jordan elimination with full pivoting.
//
// The equation system is specific to one problem the user wants
// solved right now, not universal reference data or reusable
// equipment configuration, so — like calibrat and vrev in the same
// conversion group — this program reads its input fresh from a file
// named on the command line each run, in the same
// STARTOFDATA/ENDOFDATA text format the original used; see
// ai/plans/c-to-go-conversion-plan.md, "Data-file strategy for
// Tier 2".
//
// Converted from SIMUL.C (M. W. Klotz), Math/simul.
package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
	"strings"

	"mwkgo/internal/legacydat"
)

var errSingular = errors.New("matrix is singular")

// linearSystem is a * x = b for an n-unknown system.
type linearSystem struct {
	a [][]float64 // n x n coefficient matrix
	b []float64   // n constants
}

// gaussJordan solves sys.a * x = sys.b in place by full-pivoting
// Gauss-Jordan elimination (the classic Numerical Recipes gaussj
// method, which SIMUL.C itself ports): sys.a ends up holding its own
// inverse, and sys.b ends up holding the solution vector. Only the
// solution vector is used by this program; the inverse is a
// byproduct of the method, not something SIMUL.C's own output uses
// either.
func gaussJordan(sys linearSystem) error {
	n := len(sys.a)
	pivoted := make([]bool, n)
	pivotRow := make([]int, n)
	pivotCol := make([]int, n)

	for i := 0; i < n; i++ {
		irow, icol, big := -1, -1, 0.0
		for j := 0; j < n; j++ {
			if pivoted[j] {
				continue
			}
			for k := 0; k < n; k++ {
				if pivoted[k] {
					continue
				}
				if abs := math.Abs(sys.a[j][k]); abs > big {
					big, irow, icol = abs, j, k
				}
			}
		}
		if irow == -1 {
			return errSingular
		}

		if irow != icol {
			sys.a[irow], sys.a[icol] = sys.a[icol], sys.a[irow]
			sys.b[irow], sys.b[icol] = sys.b[icol], sys.b[irow]
		}
		pivotRow[i], pivotCol[i] = irow, icol

		pivot := sys.a[icol][icol]
		pivoted[icol] = true
		sys.a[icol][icol] = 1
		for l := 0; l < n; l++ {
			sys.a[icol][l] /= pivot
		}
		sys.b[icol] /= pivot

		for row := 0; row < n; row++ {
			if row == icol {
				continue
			}
			factor := sys.a[row][icol]
			sys.a[row][icol] = 0
			for l := 0; l < n; l++ {
				sys.a[row][l] -= sys.a[icol][l] * factor
			}
			sys.b[row] -= sys.b[icol] * factor
		}
	}

	for l := n - 1; l >= 0; l-- {
		if pivotRow[l] != pivotCol[l] {
			for k := 0; k < n; k++ {
				sys.a[k][pivotRow[l]], sys.a[k][pivotCol[l]] = sys.a[k][pivotCol[l]], sys.a[k][pivotRow[l]]
			}
		}
	}
	return nil
}

// loadSystem parses a linear system from r: a line with just the
// number of unknowns n, then n lines each with n coefficients and a
// constant term, comma-separated, bracketed by
// STARTOFDATA/ENDOFDATA.
func loadSystem(r io.Reader) (linearSystem, error) {
	rows, err := legacydat.Rows(r, legacydat.Options{StartMarker: "STARTOFDATA", EndMarker: "ENDOFDATA"})
	if err != nil {
		return linearSystem{}, fmt.Errorf("scan data: %w", err)
	}
	if len(rows) == 0 {
		return linearSystem{}, errors.New("no data found")
	}

	nFields := legacydat.Fields(rows[0], ",\t;")
	n, err := strconv.Atoi(strings.TrimSpace(nFields[0]))
	if err != nil {
		return linearSystem{}, fmt.Errorf("line 1: number of unknowns %q: %w", nFields[0], err)
	}
	if len(rows) != n+1 {
		return linearSystem{}, fmt.Errorf("expected %d equations after the unknown count, found %d", n, len(rows)-1)
	}

	sys := linearSystem{a: make([][]float64, n), b: make([]float64, n)}
	for i := 0; i < n; i++ {
		fields := legacydat.Fields(rows[i+1], ",\t;")
		if len(fields) != n+1 {
			return linearSystem{}, fmt.Errorf("line %d: %q has %d values, want %d coefficients plus a constant", i+2, rows[i+1], len(fields), n+1)
		}
		sys.a[i] = make([]float64, n)
		for j := 0; j < n; j++ {
			v, err := strconv.ParseFloat(strings.TrimSpace(fields[j]), 64)
			if err != nil {
				return linearSystem{}, fmt.Errorf("line %d: coefficient %q: %w", i+2, fields[j], err)
			}
			sys.a[i][j] = v
		}
		v, err := strconv.ParseFloat(strings.TrimSpace(fields[n]), 64)
		if err != nil {
			return linearSystem{}, fmt.Errorf("line %d: constant %q: %w", i+2, fields[n], err)
		}
		sys.b[i] = v
	}
	return sys, nil
}

func main() {
	dataPath := flag.String("equations", "", "path to a linear-system data file (see MWKGo/simul/testdata/example.dat)")
	flag.Parse()

	if *dataPath == "" {
		fmt.Println("usage: simul -equations <file>")
		fmt.Println("see MWKGo/simul/testdata/example.dat for the expected format")
		os.Exit(1)
	}

	f, err := os.Open(*dataPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "simul:", err)
		os.Exit(1)
	}
	sys, err := loadSystem(f)
	f.Close()
	if err != nil {
		fmt.Fprintln(os.Stderr, "simul:", err)
		os.Exit(1)
	}
	n := len(sys.a)

	fmt.Println()
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			fmt.Printf(" %+.3f *x%d", sys.a[i][j], j+1)
		}
		fmt.Printf(" = %.3f\n", sys.b[i])
	}

	if err := gaussJordan(sys); err != nil {
		fmt.Println("ERROR")
		os.Exit(1)
	}

	fmt.Println()
	for i := 0; i < n; i++ {
		fmt.Printf("x%d = %f\n", i+1, sys.b[i])
	}
}
