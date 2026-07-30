// spline fits a natural cubic spline through a set of (x,y) points
// (Sedgewick's tridiagonal method, "Algorithms" p.545) and prints a
// sampled table of the fitted curve; it can also render the curve and
// the original data points as an SVG diagram.
//
// The data points are specific to one curve-fitting job — SPLINE.DAT's
// own shipped example is explicitly a demonstration dataset for the
// algorithm itself, not reference data or reusable configuration — so,
// like calibrat, vrev, simul, curfit, and colsort in earlier
// conversion groups, this program reads its input fresh from a file
// named on the command line each run, in the same
// STARTOFDATA/ENDOFDATA text format the original used.
//
// Converted from SPLINE.C (M. W. Klotz), Math/spline. Per
// ai/plans/c-to-go-conversion-plan.md's Tier 3 "Graphics scope"
// resolution, the original's DOS graphics calls are replaced by
// internal/svgplot; the interactive mouse-driven menu (exit/toggle
// menu/redraw, and a live x,y coordinate readout under the cursor) is
// dropped entirely, since none of it has any equivalent in a static,
// non-interactive SVG output — this program has no companion .TXT
// file to separately record that decision in, so it is recorded here
// and in docs/calculators/spline.md instead.
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
	"mwkgo/internal/svgplot"
)

// splinePoint is one (x,y) data point.
type splinePoint struct{ X, Y float64 }

// cubicSpline is a fitted natural cubic spline: the sorted input
// points, plus the working arrays SPLINE.C's own spline()/seval()
// compute and reuse (u: consecutive x differences; p: the spline's
// second-derivative-like coefficients from the tridiagonal solve).
type cubicSpline struct {
	points []splinePoint
	u, p   []float64
}

// fitCubicSpline fits a natural cubic spline through points (already
// sorted ascending by X), matching SPLINE.C's own spline() exactly:
// a tridiagonal system for the curve's second-derivative-like
// coefficients, with natural (zero second derivative) boundary
// conditions at both ends.
func fitCubicSpline(points []splinePoint) (*cubicSpline, error) {
	n := len(points)
	if n < 3 {
		return nil, fmt.Errorf("need at least 3 points for a cubic spline fit, got %d", n)
	}

	d := make([]float64, n)
	u := make([]float64, n)
	w := make([]float64, n)
	p := make([]float64, n)

	for i := 1; i < n-1; i++ {
		d[i] = 2 * (points[i+1].X - points[i-1].X)
	}
	for i := 0; i < n-1; i++ {
		u[i] = points[i+1].X - points[i].X
	}
	for i := 1; i < n-1; i++ {
		w[i] = 6 * ((points[i+1].Y-points[i].Y)/u[i] - (points[i].Y-points[i-1].Y)/u[i-1])
	}
	p[0] = 0
	p[n-1] = 0
	for i := 1; i < n-2; i++ {
		w[i+1] -= w[i] * u[i] / d[i]
		d[i+1] -= u[i] * u[i] / d[i]
	}
	for i := n - 2; i > 0; i-- {
		p[i] = (w[i] - u[i]*p[i+1]) / d[i]
	}

	return &cubicSpline{points: points, u: u, p: p}, nil
}

// cubicHermiteTerm is SPLINE.C's own f(a) = a^3 - a, used by seval().
func cubicHermiteTerm(a float64) float64 { return a*a*a - a }

// Eval evaluates the fitted spline at z, matching SPLINE.C's own
// seval() exactly; z outside the original points' own x range is an
// error (the original prints a message and returns 0.0, which this
// conversion treats as a reportable error instead of a silently wrong
// value).
func (s *cubicSpline) Eval(z float64) (float64, error) {
	n := len(s.points)
	if z < s.points[0].X || z > s.points[n-1].X {
		return 0, fmt.Errorf("z=%v out of range [%v,%v]", z, s.points[0].X, s.points[n-1].X)
	}
	i := 0
	for ; i < n-1; i++ {
		if z <= s.points[i+1].X {
			break
		}
	}
	t := (z - s.points[i].X) / s.u[i]
	return t*s.points[i+1].Y + (1-t)*s.points[i].Y +
		s.u[i]*s.u[i]*(cubicHermiteTerm(t)*s.p[i+1]+cubicHermiteTerm(1-t)*s.p[i])/6.0, nil
}

// sampleCurve samples the fitted spline at n+1 evenly-spaced points
// from the first to the last original x value, matching SPLINE.C's
// own draw() sampling (nt=40 there).
func sampleCurve(s *cubicSpline, n int) []splinePoint {
	t0 := s.points[0].X
	tLast := s.points[len(s.points)-1].X
	dt := (tLast - t0) / float64(n)

	samples := make([]splinePoint, n+1)
	for i := 0; i <= n; i++ {
		t := t0 + float64(i)*dt
		if t > tLast {
			t = tLast
		}
		y, _ := s.Eval(t) // always in range by construction
		samples[i] = splinePoint{X: t, Y: y}
	}
	return samples
}

// loadPoints parses (x,y) pairs (comma-separated, bracketed by
// STARTOFDATA/ENDOFDATA) from r, sorted ascending by x — SPLINE.C
// itself sorts before use, and its own shipped example is
// deliberately entered out of order to demonstrate this.
func loadPoints(r io.Reader) ([]splinePoint, error) {
	rows, err := legacydat.Rows(r, legacydat.Options{StartMarker: "STARTOFDATA", EndMarker: "ENDOFDATA"})
	if err != nil {
		return nil, fmt.Errorf("scan data: %w", err)
	}

	points := make([]splinePoint, 0, len(rows))
	for i, row := range rows {
		fields := legacydat.Fields(row, ",;")
		if len(fields) < 2 {
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
		points = append(points, splinePoint{X: x, Y: y})
	}

	sort.Slice(points, func(i, j int) bool { return points[i].X < points[j].X })
	return points, nil
}

// renderSpline draws the fitted curve (color 14, solid) and the
// original data points (as small crosses, color 15), framed by a
// dotted x/y axis pair, matching SPLINE.C's own draw() exactly. The
// window bounds pad the data's own bounding box by 20% on each side,
// matching SPLINE.C's own setup().
func renderSpline(points []splinePoint, s *cubicSpline) *svgplot.Canvas {
	xmin, xmax := points[0].X, points[0].X
	ymin, ymax := points[0].Y, points[0].Y
	for _, p := range points {
		xmin, xmax = math.Min(xmin, p.X), math.Max(xmax, p.X)
		ymin, ymax = math.Min(ymin, p.Y), math.Max(ymax, p.Y)
	}
	const pad = 0.2
	dx, dy := math.Abs(xmax-xmin), math.Abs(ymax-ymin)
	xmin -= pad * dx
	xmax += pad * dx
	ymin -= pad * dy
	ymax += pad * dy
	wx, wy := math.Abs(xmax-xmin), math.Abs(ymax-ymin)

	c := svgplot.New(xmin, ymin, xmax, ymax, 640, 480)
	c.Box(xmin, ymin, xmax, ymax, 15, false)
	c.DashedLine(xmin, 0, xmax, 0, 15)
	c.DashedLine(0, ymin, 0, ymax, 15)

	samples := sampleCurve(s, 40)
	for i := 1; i < len(samples); i++ {
		c.Line(samples[i-1].X, samples[i-1].Y, samples[i].X, samples[i].Y, 14)
	}

	for _, p := range points {
		drawCross(c, p.X, p.Y, 0.005, wx, wy, 15)
	}
	return c
}

// drawCross draws a small "+" at (x,y), matching SPLINE.C's own
// cross()/tic(): both arms have the same length in window-coordinate
// units (del*wx), a deliberate original quirk that keeps the cross
// visually square-ish on screen even when the window's x and y spans
// differ.
func drawCross(c *svgplot.Canvas, x, y, del, wx, wy float64, color int) {
	t := del * wx
	c.Line(x-t, y, x+t, y, color)
	c.Line(x, y-t, x, y+t, color)
}

func main() {
	dataPath := flag.String("data", "", "path to an x,y spline data file (see MWKGo/spline/testdata/example.dat)")
	svgPath := flag.String("svg", "", "path to write an SVG rendering of the fitted curve (optional)")
	flag.Parse()

	if *dataPath == "" {
		fmt.Println("usage: spline -data <file> [-svg <file>]")
		fmt.Println("see MWKGo/spline/testdata/example.dat for the expected format")
		os.Exit(1)
	}

	f, err := os.Open(*dataPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "spline:", err)
		os.Exit(1)
	}
	points, err := loadPoints(f)
	f.Close()
	if err != nil {
		fmt.Fprintln(os.Stderr, "spline:", err)
		os.Exit(1)
	}

	s, err := fitCubicSpline(points)
	if err != nil {
		fmt.Fprintln(os.Stderr, "spline:", err)
		os.Exit(1)
	}

	fmt.Printf("%d data points read\n\n", len(points))
	fmt.Println("     X            Y (fitted)")
	for _, p := range sampleCurve(s, 40) {
		fmt.Printf("%10.4f   %10.4f\n", p.X, p.Y)
	}

	if *svgPath != "" {
		if err := os.WriteFile(*svgPath, []byte(renderSpline(points, s).String()), 0o644); err != nil {
			fmt.Fprintln(os.Stderr, "spline:", err)
			os.Exit(1)
		}
		fmt.Printf("\nWrote spline curve diagram to %s\n", *svgPath)
	}
}
