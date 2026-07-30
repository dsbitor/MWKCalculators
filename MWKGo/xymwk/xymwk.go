// xymwk transforms and measures a list of x,y points measured with a
// height gauge (per XYMWK.TXT: clamp a part to a 90-degree angle
// plate, measure feature heights, rotate 90 degrees, measure again).
// It supports the same four operations XYMWK.C offers, each
// operating on two or three user-selected points:
//
//   - reference: translate the whole point set so the selected point
//     becomes the origin.
//   - align: translate so the first selected point becomes the
//     origin, then rotate so the second selected point lies on the
//     x-axis.
//   - distance: report the distance between two selected points.
//   - pitchcircle: report the center and radius of the circle passing
//     through three selected points.
//
// XYMWK.C selects points by mouse click (optionally "snapping" to the
// nearest data point); this conversion selects points by their
// 1-based position in the data file instead, which is at least as
// precise as snapping and needs no interactive session. Per
// ai/plans/c-to-go-conversion-plan.md's Tier 3 scope discussion for
// this program, each run performs exactly one operation on the raw
// data — it does not replicate the original's interactive session,
// where each operation's result became the new working coordinate set
// for the next operation.
//
// Converted from XYMWK.C (M. W. Klotz), Math/xymwk.zip. XYMWK.C
// declares and fully implements display(), a function that prints
// exactly the table this conversion's reference/align operations
// print (original x/y, transformed x/y, radius, angle in decimal
// degrees and degrees:minutes:seconds) — matching XYMWK.TXT's own
// description of what those operations show — but never actually
// calls it anywhere in main(); in the real program, reference/align
// only silently re-plot the points on screen with no printed values
// at all. This conversion wires up that table for reference/align
// since a static CLI has no equivalent for "re-plot the points for
// the user to look at," and printing them is what the tool's own
// documentation already promises. angle() and project(), two other
// declared-but-never-called functions, are omitted as dead code, the
// same way egg omitted its own unused spline fitter copy.
package main

import (
	"flag"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
	"strings"

	"mwkgo/internal/legacydat"
	"mwkgo/internal/mwktrig"
)

// xyPoint is one raw (x,y) point as measured with the height gauge.
type xyPoint struct{ X, Y float64 }

// loadPoints parses XYMWK.DAT's own format: comma-separated x,y pairs
// bracketed by STARTOFDATA/ENDOFDATA, kept in file order (not sorted
// — unlike profile/ogive/egg, point *position* in the list is how the
// user selects a point, so reordering them would be a correctness
// bug, not a cosmetic one).
func loadPoints(r io.Reader) ([]xyPoint, error) {
	rows, err := legacydat.Rows(r, legacydat.Options{StartMarker: "STARTOFDATA", EndMarker: "ENDOFDATA"})
	if err != nil {
		return nil, fmt.Errorf("scan data: %w", err)
	}

	points := make([]xyPoint, 0, len(rows))
	for i, row := range rows {
		fields := legacydat.Fields(row, ",\t;")
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
		points = append(points, xyPoint{X: x, Y: y})
	}
	return points, nil
}

// referTo is XYMWK.C's own refer(): translate every point so that
// (xr,yr) becomes the origin.
func referTo(points []xyPoint, xr, yr float64) []xyPoint {
	out := make([]xyPoint, len(points))
	for i, p := range points {
		out[i] = xyPoint{X: p.X - xr, Y: p.Y - yr}
	}
	return out
}

// rotateDeg is XYMWK.C's own rotate(): rotate the coordinate system
// (not the points) by thetaDeg, matching the original's axis-rotation
// sign convention exactly.
func rotateDeg(points []xyPoint, thetaDeg float64) []xyPoint {
	ct := math.Cos(thetaDeg * math.Pi / 180)
	st := math.Sin(thetaDeg * math.Pi / 180)
	out := make([]xyPoint, len(points))
	for i, p := range points {
		out[i] = xyPoint{X: ct*p.X + st*p.Y, Y: -st*p.X + ct*p.Y}
	}
	return out
}

// alignTo is XYMWK.C's own align(): (x1,y1) becomes the origin, and
// the coordinate system is then rotated so (x2,y2) lies on the
// (positive) x-axis.
func alignTo(points []xyPoint, x1, y1, x2, y2 float64) []xyPoint {
	referred := referTo(points, x1, y1)
	theta := math.Atan2(y2-y1, x2-x1) * 180 / math.Pi
	return rotateDeg(referred, theta)
}

// polarPoint adds radius and four-quadrant angle (degrees, positive
// ccw from the x-axis) to an xyPoint, matching XYMWK.C's own polar().
type polarPoint struct {
	X, Y, R, Ang float64
}

func toPolar(points []xyPoint) []polarPoint {
	out := make([]polarPoint, len(points))
	for i, p := range points {
		out[i] = polarPoint{
			X: p.X, Y: p.Y,
			R:   math.Hypot(p.X, p.Y),
			Ang: math.Atan2(p.Y, p.X) * 180 / math.Pi,
		}
	}
	return out
}

// formatDMS renders an angle in degree:minute:second notation,
// matching XYMWK.C's own dms() exactly.
func formatDMS(a float64) string {
	sign := " "
	z := math.Abs(a)
	if a < 0 {
		sign = "-"
	}
	d := math.Floor(z)
	x := int(math.Round((z - d) * 3600))
	m := x / 60
	s := x - m*60
	if s >= 60 {
		s -= 60
		m++
	}
	if m >= 60 {
		m -= 60
		d++
	}
	return fmt.Sprintf("(%s%3d:%02d:%02d)", sign, int(d), m, s)
}

// dist is XYMWK.C's own dist().
func dist(x1, y1, x2, y2 float64) float64 {
	return math.Hypot(x2-x1, y2-y1)
}

// loc is XYMWK.C's own loc(): the two points a distance d1 from
// (x1,y1) and d2 from (x2,y2) — the two intersections of the circles
// of those radii centered on the two given points.
func loc(x1, y1, d1, x2, y2, d2 float64) (xs, ys, xt, yt float64) {
	d3 := dist(x1, y1, x2, y2)
	p := 0.5 * (d1 + d2 + d3)
	a1 := 2 * mwktrig.ClampedAcosDeg(math.Sqrt(p*(p-d1)/(d2*d3)))
	a3 := 2 * mwktrig.ClampedAcosDeg(math.Sqrt(p*(p-d3)/(d1*d2)))
	a2 := 180 - a1 - a3
	beta := math.Atan2(y2-y1, x2-x1) * 180 / math.Pi
	phi := beta - a2
	xs = x1 + d1*math.Cos(phi*math.Pi/180)
	ys = y1 + d1*math.Sin(phi*math.Pi/180)
	xt = x1 + d1*math.Cos((phi+2*a2)*math.Pi/180)
	yt = y1 + d1*math.Sin((phi+2*a2)*math.Pi/180)
	return xs, ys, xt, yt
}

// bisect is XYMWK.C's own bisect(): a line, given by two points on
// it, perpendicular to and bisecting the segment (x1,y1)-(x2,y2).
func bisect(x1, y1, x2, y2 float64) (x3, y3, x4, y4 float64) {
	d := dist(x1, y1, x2, y2)
	return loc(x1, y1, d, x2, y2, d)
}

// intersect is XYMWK.C's own intersect(): the intersection of lines
// (x1,y1)-(x2,y2) and (x3,y3)-(x4,y4), or ok=false if they're
// parallel.
func intersect(x1, y1, x2, y2, x3, y3, x4, y4 float64) (xi, yi float64, ok bool) {
	d := (y2-y1)*(x4-x3) - (y4-y3)*(x2-x1)
	if d == 0 {
		return 0, 0, false
	}
	xi = ((x1*y2-x2*y1)*(x4-x3) - (x3*y4-x4*y3)*(x2-x1)) / d
	yi = ((x1*y2-x2*y1)*(y4-y3) - (x3*y4-x4*y3)*(y2-y1)) / d
	return xi, yi, true
}

// circumcircle is XYMWK.C's own circ3(): the center and radius of the
// circle passing through three points, found via the perpendicular
// bisectors of two of the chords. XYMWK.C's own comment on circ3
// flags a known gap it never fixes: "no error checking is done for
// the anomalous case in which the three points are colinear," and
// even suggests how: "can be done by aligning to points 1 and 2; if y
// coordinate of point 3 is zero after said alignment, points are
// colinear." This conversion implements exactly that check (reusing
// alignTo itself) rather than relying on intersect's own d==0 test,
// since two bisector directions computed through trig functions are
// essentially never bit-for-bit parallel even when the three points
// are mathematically colinear — the same reason the original's
// unchecked call could return a wild, wrong circle instead of
// visibly failing.
func circumcircle(x1, y1, x2, y2, x3, y3 float64) (xc, yc, r float64, err error) {
	aligned := alignTo([]xyPoint{{X: x3, Y: y3}}, x1, y1, x2, y2)
	tol := 1e-9 * math.Max(dist(x1, y1, x2, y2), 1)
	if math.Abs(aligned[0].Y) < tol {
		return 0, 0, 0, fmt.Errorf("points are collinear: no unique circle passes through them")
	}

	xa, ya, xb, yb := bisect(x1, y1, x2, y2)
	xd, yd, xe, ye := bisect(x2, y2, x3, y3)
	xc, yc, ok := intersect(xa, ya, xb, yb, xd, yd, xe, ye)
	if !ok {
		return 0, 0, 0, fmt.Errorf("points are collinear: no unique circle passes through them")
	}
	r = dist(x1, y1, xc, yc)
	return xc, yc, r, nil
}

// printTable prints XYMWK.C's own display() table: for each point,
// its original (hd) and transformed (hp) coordinates, radius, and
// angle in both decimal degrees and degree:minute:second notation.
func printTable(raw []xyPoint, transformed []polarPoint) {
	fmt.Printf("%2s: %8s %8s  %8s %8s  %8s %9s  %10s\n",
		"PT", "hdx", "hdy", "hpx", "hpy", "hpr", "hpang(d)", "hpang(d:m:s)")
	for i := range raw {
		fmt.Printf("%2d: %8.3f %8.3f  %8.3f %8.3f  %8.3f %9.4f  %s\n",
			i+1, raw[i].X, raw[i].Y,
			transformed[i].X, transformed[i].Y, transformed[i].R, transformed[i].Ang,
			formatDMS(transformed[i].Ang))
	}
}

func main() {
	dataPath := flag.String("data", "", "path to an XYMWK-format data file (see MWKGo/xymwk/testdata/)")
	op := flag.String("op", "", "operation: reference, align, distance, or pitchcircle")
	p1 := flag.Int("p1", 0, "1-based index of the first selected point")
	p2 := flag.Int("p2", 0, "1-based index of the second selected point (align, distance, pitchcircle)")
	p3 := flag.Int("p3", 0, "1-based index of the third selected point (pitchcircle)")
	flag.Parse()

	if *dataPath == "" || *op == "" {
		fmt.Println("usage: xymwk -data <file> -op reference|align|distance|pitchcircle -p1 <n> [-p2 <n>] [-p3 <n>]")
		fmt.Println("see MWKGo/xymwk/testdata/ for the expected data format")
		os.Exit(1)
	}

	f, err := os.Open(*dataPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "xymwk:", err)
		os.Exit(1)
	}
	points, err := loadPoints(f)
	f.Close()
	if err != nil {
		fmt.Fprintln(os.Stderr, "xymwk:", err)
		os.Exit(1)
	}

	point := func(n int, label string) xyPoint {
		if n < 1 || n > len(points) {
			fmt.Fprintf(os.Stderr, "xymwk: %s point %d out of range (have %d points)\n", label, n, len(points))
			os.Exit(1)
		}
		return points[n-1]
	}

	switch *op {
	case "reference":
		p := point(*p1, "-p1")
		transformed := referTo(points, p.X, p.Y)
		printTable(points, toPolar(transformed))

	case "align":
		a := point(*p1, "-p1")
		b := point(*p2, "-p2")
		transformed := alignTo(points, a.X, a.Y, b.X, b.Y)
		printTable(points, toPolar(transformed))

	case "distance":
		a := point(*p1, "-p1")
		b := point(*p2, "-p2")
		fmt.Printf("Distance = %8.3f\n", dist(a.X, a.Y, b.X, b.Y))

	case "pitchcircle":
		a := point(*p1, "-p1")
		b := point(*p2, "-p2")
		c := point(*p3, "-p3")
		xc, yc, r, err := circumcircle(a.X, a.Y, b.X, b.Y, c.X, c.Y)
		if err != nil {
			fmt.Fprintln(os.Stderr, "xymwk:", err)
			os.Exit(1)
		}
		fmt.Printf("Circle center at x,y = %8.3f,%8.3f\n", xc, yc)
		fmt.Printf("Circle radius = %8.3f\n", r)

	default:
		fmt.Fprintf(os.Stderr, "xymwk: unknown -op %q (want reference, align, distance, or pitchcircle)\n", *op)
		os.Exit(1)
	}
}
