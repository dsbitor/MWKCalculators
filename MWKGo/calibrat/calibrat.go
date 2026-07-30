// calibrat finds the best-fit linear calibration equation (y = A*x +
// B, where x is a known "truth" value and y is what an instrument
// measured for it) from a set of truth/measured checkpoint pairs, by
// least squares, and can then print out a calibration table over a
// chosen range.
//
// The checkpoint data is specific to one calibration run of one
// instrument — not universal reference data, and not a fixed piece
// of equipment configuration reused indefinitely like a dividing
// head's plates — so, unlike the reference- and machine-config-
// bucket programs elsewhere in this project, this program reads its
// input fresh from a file named on the command line each run, in the
// same STARTOFDATA/ENDOFDATA text format the original used, rather
// than from either SQLite database; see
// ai/plans/c-to-go-conversion-plan.md, "Data-file strategy for
// Tier 2".
//
// Converted from CALIBRAT.C (M. W. Klotz), Math/calibrat.
package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"

	"mwkgo/internal/legacydat"
	"mwkgo/internal/promptio"
)

// calibrationPoint is one truth/measured checkpoint.
type calibrationPoint struct {
	Truth, Measured float64
}

var errZeroDeterminant = errors.New("determinant is zero, impossible to proceed")

// linearFit returns the least-squares slope and intercept (y = A*x +
// B) for points.
func linearFit(points []calibrationPoint) (a, b float64, err error) {
	n := float64(len(points))
	var sumX, sumX2, sumY, sumXY float64
	for _, p := range points {
		sumX += p.Truth
		sumX2 += p.Truth * p.Truth
		sumY += p.Measured
		sumXY += p.Truth * p.Measured
	}
	det := sumX2*n - sumX*sumX
	if det == 0 {
		return 0, 0, errZeroDeterminant
	}
	a = (n*sumXY - sumX*sumY) / det
	b = (sumX2*sumY - sumX*sumXY) / det
	return a, b, nil
}

// loadPoints parses truth/measured pairs (comma- or tab-separated,
// bracketed by STARTOFDATA/ENDOFDATA) from r, sorted ascending by
// truth value — CALIBRAT.C itself sorts before use, and its default
// table start/end values depend on the sort.
func loadPoints(r io.Reader) ([]calibrationPoint, error) {
	rows, err := legacydat.Rows(r, legacydat.Options{StartMarker: "STARTOFDATA", EndMarker: "ENDOFDATA"})
	if err != nil {
		return nil, fmt.Errorf("scan data: %w", err)
	}

	points := make([]calibrationPoint, 0, len(rows))
	for i, row := range rows {
		fields := legacydat.Fields(row, ",\t;")
		if len(fields) != 2 {
			return nil, fmt.Errorf("line %d: %q does not split into truth and measured values", i+1, row)
		}
		truth, err := strconv.ParseFloat(strings.TrimSpace(fields[0]), 64)
		if err != nil {
			return nil, fmt.Errorf("line %d: truth value %q: %w", i+1, fields[0], err)
		}
		measured, err := strconv.ParseFloat(strings.TrimSpace(fields[1]), 64)
		if err != nil {
			return nil, fmt.Errorf("line %d: measured value %q: %w", i+1, fields[1], err)
		}
		points = append(points, calibrationPoint{Truth: truth, Measured: measured})
	}

	sort.Slice(points, func(i, j int) bool { return points[i].Truth < points[j].Truth })
	return points, nil
}

func main() {
	dataPath := flag.String("data", "", "path to a truth/measured calibration data file (see MWKGo/calibrat/testdata/example.dat)")
	flag.Parse()

	if *dataPath == "" {
		fmt.Println("usage: calibrat -data <file>")
		fmt.Println("see MWKGo/calibrat/testdata/example.dat for the expected format")
		os.Exit(1)
	}

	f, err := os.Open(*dataPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "calibrat:", err)
		os.Exit(1)
	}
	points, err := loadPoints(f)
	f.Close()
	if err != nil {
		fmt.Fprintln(os.Stderr, "calibrat:", err)
		os.Exit(1)
	}

	prompter, err := promptio.New(os.Stdin, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, "calibrat:", err)
		os.Exit(1)
	}

	fmt.Println("CALIBRATE A LINEAR SCALE")
	fmt.Println()
	fmt.Printf("%d data pairs read\n", len(points))

	if len(points) < 2 {
		fmt.Println("Too few data pairs for solution.  Try again.")
		os.Exit(1)
	}

	a, b, err := linearFit(points)
	if err != nil {
		fmt.Println("Determinat=0, impossible to proceed.")
		os.Exit(1)
	}

	fmt.Println()
	fmt.Println("Calibration Equations:")
	fmt.Println()
	fmt.Println("x=truth value, y=measured value")
	fmt.Println("y = A*x + B  or x = (y-B)/A")
	fmt.Printf("A = %.6f\n", a)
	fmt.Printf("B = %.6f\n", b)

	if strings.ToLower(prompter.Line("\nDo you want to construct a calibration table [Y]/N ? ")) == "n" {
		return
	}

	tstart := prompter.Float("Table starting truth value", points[0].Truth)
	tstop := prompter.Float("Table ending truth value", points[len(points)-1].Truth)

	var tinc float64
	for {
		tinc = prompter.Float("Table increment", (tstop-tstart)/float64(len(points)))
		if tinc > 0 {
			break
		}
		fmt.Println("Bad increment.  Try again.")
	}

	printTable(a, b, tstart, tstop, tinc)
}

// printTable reproduces CALIBRAT.C's own wdata(): two tables over the
// same tstart..tstop range by tinc, the first treating each value as
// a measured reading and inverting the fit to recover truth, the
// second treating it as a truth value and applying the fit forward
// to predict what would be measured.
func printTable(a, b, tstart, tstop, tinc float64) {
	fmt.Println("\nCalibration Table")
	fmt.Println("\nFirst entry is measured, second is truth value")
	fmt.Println()
	for t := tstart; t <= tstop; t += tinc {
		fmt.Printf("%f <=> %f\n", t, (t-b)/a)
	}

	fmt.Println("\nFirst entry is truth, second is measured value")
	fmt.Println()
	for t := tstart; t <= tstop; t += tinc {
		fmt.Printf("%f <=> %f\n", t, a*t+b)
	}
}
