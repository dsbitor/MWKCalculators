// curfit fits a curve — polynomial, logarithmic, exponential, or
// power — to a set of (x,y) data pairs by least squares, reporting
// the fit's coefficients, a per-point error table, and a degrees-of-
// freedom-adjusted correlation coefficient. A session can try several
// fit types against the same data before quitting.
//
// The data pairs are specific to one curve-fitting job, not universal
// reference data or reusable equipment configuration, so — like
// calibrat, vrev, simul, and colsort in earlier conversion groups —
// this program reads its input fresh from a file named on the command
// line each run, in the same STARTOFDATA/ENDOFDATA text format the
// original used; see ai/plans/c-to-go-conversion-plan.md, "Data-file
// strategy for Tier 2".
//
// Converted from CURFIT.C (M. W. Klotz), Math/curfit.
package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"

	"mwkgo/internal/legacydat"
	"mwkgo/internal/promptio"
)

// maxPolynomialDegree matches CURFIT.C's own MN-1 matrix-size bound.
const maxPolynomialDegree = 49

var errSingular = errors.New("matrix is singular")

// point is one (x,y) data pair.
type point struct{ X, Y float64 }

// linearSystem and gaussJordan are simul's own full-pivoting
// Gauss-Jordan solver (the same "gaussj" method CURFIT.C's own
// gjor() implements); see MWKGo/simul/simul.go for the full
// explanation.
type linearSystem struct {
	a [][]float64
	b []float64
}

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

// polyFit returns the degree-`degree` polynomial's coefficients
// [A0, A1, ..., Adegree] (Y = A0 + A1*X + A2*X^2 + ...) by solving the
// least-squares normal equations, matching CURFIT.C's own poly()
// exactly.
func polyFit(points []point, degree int) ([]float64, error) {
	nmat := degree + 1
	a := make([][]float64, nmat)
	for i := range a {
		a[i] = make([]float64, nmat)
	}
	b := make([]float64, nmat)

	for _, p := range points {
		for i := 0; i < nmat; i++ {
			b[i] += p.Y * math.Pow(p.X, float64(i))
			for j := i; j < nmat; j++ {
				a[i][j] += math.Pow(p.X, float64(i+j))
			}
		}
	}
	for i := 0; i < nmat; i++ {
		for j := i + 1; j < nmat; j++ {
			a[j][i] = a[i][j]
		}
	}

	if err := gaussJordan(linearSystem{a: a, b: b}); err != nil {
		return nil, err
	}
	return b, nil
}

func polyEval(coeffs []float64, x float64) float64 {
	y := coeffs[0]
	for j := 1; j < len(coeffs); j++ {
		y += coeffs[j] * math.Pow(x, float64(j))
	}
	return y
}

// logFit returns A, B for Y = A + B*ln(X), by linear regression of Y
// against ln(X).
func logFit(points []point) (a, b float64) {
	n := float64(len(points))
	var sumV, sumV2, sumY, sumVY float64
	for _, p := range points {
		v := math.Log(p.X)
		sumV += v
		sumV2 += v * v
		sumY += p.Y
		sumVY += v * p.Y
	}
	x1, y1 := sumV/n, sumY/n
	b = (sumVY - n*x1*y1) / (sumV2 - n*x1*x1)
	a = y1 - b*x1
	return a, b
}

// expFit returns A, B for Y = A*exp(B*X), by linear regression of
// ln(Y) against X.
func expFit(points []point) (a, b float64) {
	n := float64(len(points))
	var sumX, sumX2, sumV, sumXV float64
	for _, p := range points {
		v := math.Log(p.Y)
		sumX += p.X
		sumX2 += p.X * p.X
		sumV += v
		sumXV += v * p.X
	}
	x1, y1 := sumX/n, sumV/n
	b = (sumXV - n*x1*y1) / (sumX2 - n*x1*x1)
	a = math.Exp(y1 - b*x1)
	return a, b
}

// powerFit returns A, B for Y = A*X^B, by linear regression of ln(Y)
// against ln(X).
func powerFit(points []point) (a, b float64) {
	n := float64(len(points))
	var sumU, sumU2, sumV, sumUV float64
	for _, p := range points {
		u, v := math.Log(p.X), math.Log(p.Y)
		sumU += u
		sumU2 += u * u
		sumV += v
		sumUV += v * u
	}
	x1, y1 := sumU/n, sumV/n
	b = (sumUV - n*x1*y1) / (sumU2 - n*x1*x1)
	a = math.Exp(y1 - b*x1)
	return a, b
}

// rawSST is the total-sum-of-squares term (y2 - n*y1^2) computed from
// the data's own Y values, used by the polynomial and logarithmic
// fits.
func rawSST(points []point) float64 {
	n := float64(len(points))
	var y1, y2 float64
	for _, p := range points {
		y1 += p.Y
		y2 += p.Y * p.Y
	}
	y1 /= n
	return y2 - n*y1*y1
}

// logTransformedSST is the same total-sum-of-squares term, but
// computed from ln(Y) rather than Y — matching CURFIT.C's own
// expx() and power(), whose "y1"/"y2" bookkeeping tracks ln(Y) sums
// throughout, not the raw dependent variable.
func logTransformedSST(points []point) float64 {
	n := float64(len(points))
	var y1, y2 float64
	for _, p := range points {
		v := math.Log(p.Y)
		y1 += v
		y2 += v * v
	}
	y1 /= n
	return y2 - n*y1*y1
}

// correlationCoefficient is CURFIT.C's own summary(): a degrees-of-
// freedom-adjusted measure of fit quality, closer to 1 for a better
// fit.
func correlationCoefficient(n, nmat int, sst, sse float64) float64 {
	return 1 - sse*float64(n-1)/(sst*float64(n-nmat-1))
}

// evaluatedPoint is one data point's fit error, plus whether it set a
// new running-maximum |error percent| as the points were scanned in
// order — matching CURFIT.C's own "**" flag, which marks every point
// that raises the running maximum, not only the single overall worst
// point.
type evaluatedPoint struct {
	X, Y, YCalc, Error, ErrorPercent float64
	NewRunningMax                    bool
}

// evaluate computes each point's fit error against ycalc and the sum
// of squared (absolute-scale) errors.
func evaluate(points []point, ycalc func(x float64) float64) (evaluated []evaluatedPoint, sse float64) {
	evaluated = make([]evaluatedPoint, len(points))
	runningMax := 0.0
	for i, p := range points {
		yy := ycalc(p.X)
		err := yy - p.Y
		var perr float64
		if p.Y != 0 {
			perr = 100 * err / p.Y
		}
		newMax := math.Abs(perr) > runningMax
		if newMax {
			runningMax = math.Abs(perr)
		}
		evaluated[i] = evaluatedPoint{X: p.X, Y: p.Y, YCalc: yy, Error: err, ErrorPercent: perr, NewRunningMax: newMax}
		sse += err * err
	}
	return evaluated, sse
}

// loadPoints parses (x,y) pairs (comma- or tab-separated, bracketed
// by STARTOFDATA/ENDOFDATA) from r, sorted ascending by x —
// CURFIT.C itself sorts before use.
func loadPoints(r io.Reader) ([]point, error) {
	rows, err := legacydat.Rows(r, legacydat.Options{StartMarker: "STARTOFDATA", EndMarker: "ENDOFDATA"})
	if err != nil {
		return nil, fmt.Errorf("scan data: %w", err)
	}

	points := make([]point, 0, len(rows))
	for i, row := range rows {
		fields := legacydat.Fields(row, ",\t;")
		if len(fields) != 2 {
			return nil, fmt.Errorf("line %d: %q does not split into x and y values", i+1, row)
		}
		x, err := strconv.ParseFloat(strings.TrimSpace(fields[0]), 64)
		if err != nil {
			return nil, fmt.Errorf("line %d: x value %q: %w", i+1, fields[0], err)
		}
		y, err := strconv.ParseFloat(strings.TrimSpace(fields[1]), 64)
		if err != nil {
			return nil, fmt.Errorf("line %d: y value %q: %w", i+1, fields[1], err)
		}
		points = append(points, point{X: x, Y: y})
	}

	sort.Slice(points, func(i, j int) bool { return points[i].X < points[j].X })
	return points, nil
}

func printMenu() {
	fmt.Println("1 - polynomial fit")
	fmt.Println("2 - logarithmic fit")
	fmt.Println("3 - exponential fit")
	fmt.Println("4 - power fit")
	fmt.Println("Q - quit")
}

// printReport prints ycalc's fit error table and correlation
// coefficient against points, matching CURFIT.C's own per-fit output
// layout.
func printReport(points []point, ycalc func(x float64) float64, nmat int, sst float64) {
	evaluated, sse := evaluate(points, ycalc)

	fmt.Printf("\n%6s %10s %10s %10s %10s %10s\n\n", "index", "xdata", "ydata", "ycalc", "error", "err%")
	for i, e := range evaluated {
		line := fmt.Sprintf("%6d %10.2E %10.2E %10.2E %10.2E %10.2E", i, e.X, e.Y, e.YCalc, e.Error, e.ErrorPercent)
		if e.NewRunningMax {
			line += " **"
		}
		fmt.Println(line)
	}

	r := correlationCoefficient(len(points), nmat, sst, sse)
	fmt.Printf("\ncorrelation coefficient = %E\n", r)
	fmt.Println("closer to 1, better the fit")
}

func runPolyFit(prompter *promptio.Prompter, points []point) {
	maxDegree := len(points) - 1
	if maxDegree > maxPolynomialDegree {
		maxDegree = maxPolynomialDegree
	}
	fmt.Printf("\nmaximum degree of polynomial fit possible = %d\n", maxDegree)

	var degree int
	for {
		degree = prompter.Int("Degree of polynomial fit", 3)
		if degree >= 1 && degree <= maxDegree {
			break
		}
		fmt.Println("INVALID DEGREE, TRY AGAIN")
	}

	coeffs, err := polyFit(points, degree)
	if err != nil {
		fmt.Println("SINGULAR MATRIX, CANNOT FIT")
		return
	}

	fmt.Println("polynomial fit Y = A0 + A1*X + A2*X^2 + A3*X^3 + ...")
	fmt.Printf("order of polynomial fit requested = %d\n", degree)
	for i, c := range coeffs {
		fmt.Printf("A%d = %f\n", i, c)
	}

	printReport(points, func(x float64) float64 { return polyEval(coeffs, x) }, len(coeffs), rawSST(points))
}

func runLogFit(points []point) {
	a, b := logFit(points)
	fmt.Println("logarithmic fit Y = A + B * ln(X)")
	fmt.Printf("\nA = %E\n", a)
	fmt.Printf("B = %E\n", b)
	printReport(points, func(x float64) float64 { return a + b*math.Log(x) }, 2, rawSST(points))
}

func runExpFit(points []point) {
	a, b := expFit(points)
	fmt.Println("exponential fit Y = A * exp (B * X)")
	fmt.Printf("\nA = %E\n", a)
	fmt.Printf("B = %E\n", b)
	printReport(points, func(x float64) float64 { return a * math.Exp(b*x) }, 2, logTransformedSST(points))
}

func runPowerFit(points []point) {
	a, b := powerFit(points)
	fmt.Println("power fit Y = A * X^B")
	fmt.Printf("\nA = %E\n", a)
	fmt.Printf("B = %E\n", b)
	printReport(points, func(x float64) float64 { return a * math.Pow(x, b) }, 2, logTransformedSST(points))
}

func main() {
	dataPath := flag.String("data", "", "path to an x,y curve-fit data file (see MWKGo/curfit/testdata/example.dat)")
	flag.Parse()

	if *dataPath == "" {
		fmt.Println("usage: curfit -data <file>")
		fmt.Println("see MWKGo/curfit/testdata/example.dat for the expected format")
		os.Exit(1)
	}

	f, err := os.Open(*dataPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "curfit:", err)
		os.Exit(1)
	}
	points, err := loadPoints(f)
	f.Close()
	if err != nil {
		fmt.Fprintln(os.Stderr, "curfit:", err)
		os.Exit(1)
	}

	fmt.Println("CURVE FITTING")
	fmt.Println()
	fmt.Printf("%d data pairs read from file\n", len(points))

	if len(points) < 2 {
		fmt.Println("TOO FEW DATA POINTS TO CONSTRUCT FIT")
		os.Exit(1)
	}

	prompter, err := promptio.New(os.Stdin, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, "curfit:", err)
		os.Exit(1)
	}

	printMenu()
	for {
		choice := strings.ToLower(prompter.Line("(1,2,3,4,Q) ? "))
		switch choice {
		case "1":
			runPolyFit(prompter, points)
		case "2":
			runLogFit(points)
		case "3":
			runExpFit(points)
		case "4":
			runPowerFit(points)
		case "q":
			return
		default:
			fmt.Println("NOT A VALID OPTION")
		}
	}
}
