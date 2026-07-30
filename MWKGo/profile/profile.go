// profile computes an arbitrary turned profile from a user-supplied
// list of (x,radius) points — some ranges fitted with a natural cubic
// spline, the rest linearly interpolated, whichever the data file
// itself specifies — and works out an incremental roughing schedule
// for a square- or round-tipped tool, printing the schedule and
// optionally rendering the profile and cutting passes as an SVG
// diagram.
//
// The profile's own points are specific to one job — like ogive,
// egg, and calibrat/vrev/simul/curfit/colsort/spline/belt in earlier
// conversion groups — this program reads its input fresh from a file
// named on the command line each run, in the same
// STARTOFDATA/ENDOFDATA text format the original used (with ogive's
// same key=value syntax for its configuration lines, plus a variable-
// length list of (x,r) points and zero or more spline-segment
// declarations).
//
// Converted from PROFILE.C (M. W. Klotz), WorkshopUtilities/profile.
// Per ai/plans/c-to-go-conversion-plan.md's Tier 3 "Graphics scope"
// resolution, DOS graphics calls are replaced by internal/svgplot;
// the original's interactive mouse-driven menu and click-to-read-
// coordinates feature are dropped entirely, and the .OUT file-save
// feature is replaced by printing straight to stdout.
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

const rpd = math.Pi / 180

// maxSplineSegments and maxPointsPerSegment match PROFILE.C's own NS
// and NP array bounds (documented in a 5/2001 update as hard limits:
// 400 data points total, 20 spline segments, 200 points per segment).
const (
	maxDataPoints       = 400
	maxSplineSegments   = 20
	maxPointsPerSegment = 200
)

// profileConfig is everything PROFILE.DAT's key=value lines provide.
type profileConfig struct {
	StockDiam  float64
	SquareTool bool
	ToolSize   float64
	AxialStep  float64
	ScaleX     float64
	ScaleR     float64
}

func (c profileConfig) stockRadius() float64 { return 0.5 * c.StockDiam }

// profilePoint is one (x, radius) point of the turned profile.
type profilePoint struct{ X, R float64 }

// splineSegment is one index range (inclusive, into profileModel's
// own Points slice) that should be spline-fitted rather than linearly
// interpolated.
type splineSegment struct{ Start, End int }

// profileModel is a fully-loaded, ready-to-query profile: its
// configuration, its points (sorted ascending by X and scaled by
// ScaleX/ScaleR), and its spline segments. Segment indices are
// assigned during parsing, before the final sort — matching
// PROFILE.C's own rdata() exactly, including its latent assumption
// that the data file's own points are already close enough to sorted
// order that this doesn't matter (true of every shipped example).
type profileModel struct {
	Config   profileConfig
	Points   []profilePoint
	Segments []splineSegment
}

// loadProfile parses a PROFILE-format data file: key=value
// configuration lines, a variable-length list of comma-separated
// (x,r) points, and zero or more `sseg=start,finish` spline-segment
// declarations, matching PROFILE.C's own rdata() exactly.
func loadProfile(r io.Reader) (*profileModel, error) {
	rows, err := legacydat.Rows(r, legacydat.Options{StartMarker: "STARTOFDATA", EndMarker: "ENDOFDATA"})
	if err != nil {
		return nil, fmt.Errorf("scan data: %w", err)
	}

	var cfg profileConfig
	var points []profilePoint
	var segments []splineSegment

	for i, row := range rows {
		if strings.Contains(row, "=") {
			fields := legacydat.Fields(row, "=,\t;")
			if len(fields) < 2 {
				return nil, fmt.Errorf("line %d: %q does not split into a key and value", i+1, row)
			}
			value := strings.TrimSpace(fields[1])
			switch {
			case strings.Contains(row, "ds"):
				v, err := strconv.ParseFloat(value, 64)
				if err != nil {
					return nil, fmt.Errorf("line %d: ds %q: %w", i+1, value, err)
				}
				cfg.StockDiam = v
			case strings.Contains(row, "toolt"):
				v, err := strconv.Atoi(value)
				if err != nil {
					return nil, fmt.Errorf("line %d: toolt %q: %w", i+1, value, err)
				}
				cfg.SquareTool = v != 0
			case strings.Contains(row, "wt"):
				v, err := strconv.ParseFloat(value, 64)
				if err != nil {
					return nil, fmt.Errorf("line %d: wt %q: %w", i+1, value, err)
				}
				cfg.ToolSize = v
			case strings.Contains(row, "dxd"):
				v, err := strconv.ParseFloat(value, 64)
				if err != nil {
					return nil, fmt.Errorf("line %d: dxd %q: %w", i+1, value, err)
				}
				cfg.AxialStep = v
			case strings.Contains(row, "scalex"):
				v, err := strconv.ParseFloat(value, 64)
				if err != nil {
					return nil, fmt.Errorf("line %d: scalex %q: %w", i+1, value, err)
				}
				cfg.ScaleX = v
			case strings.Contains(row, "scaler"):
				v, err := strconv.ParseFloat(value, 64)
				if err != nil {
					return nil, fmt.Errorf("line %d: scaler %q: %w", i+1, value, err)
				}
				cfg.ScaleR = v
			case strings.Contains(row, "sseg"):
				if len(fields) < 3 {
					return nil, fmt.Errorf("line %d: %q does not split into a start and finish index", i+1, row)
				}
				start, err := strconv.Atoi(strings.TrimSpace(fields[1]))
				if err != nil {
					return nil, fmt.Errorf("line %d: sseg start %q: %w", i+1, fields[1], err)
				}
				finish, err := strconv.Atoi(strings.TrimSpace(fields[2]))
				if err != nil {
					return nil, fmt.Errorf("line %d: sseg finish %q: %w", i+1, fields[2], err)
				}
				if finish-start+1 > maxPointsPerSegment {
					return nil, fmt.Errorf("line %d: spline segment has more than %d points", i+1, maxPointsPerSegment)
				}
				if len(segments) >= maxSplineSegments {
					return nil, fmt.Errorf("line %d: more than %d spline segments", i+1, maxSplineSegments)
				}
				segments = append(segments, splineSegment{Start: start, End: finish})
			}
			continue
		}
		if strings.Contains(row, ",") {
			fields := legacydat.Fields(row, ",\t;")
			if len(fields) < 2 {
				return nil, fmt.Errorf("line %d: %q does not split into x and r values", i+1, row)
			}
			x, err := strconv.ParseFloat(strings.TrimSpace(fields[0]), 64)
			if err != nil {
				return nil, fmt.Errorf("line %d: x %q: %w", i+1, fields[0], err)
			}
			rad, err := strconv.ParseFloat(strings.TrimSpace(fields[1]), 64)
			if err != nil {
				return nil, fmt.Errorf("line %d: r %q: %w", i+1, fields[1], err)
			}
			if len(points) >= maxDataPoints {
				return nil, fmt.Errorf("line %d: more than %d data points", i+1, maxDataPoints)
			}
			points = append(points, profilePoint{X: x, R: rad})
		}
	}

	sort.SliceStable(points, func(i, j int) bool { return points[i].X < points[j].X })

	for i := range points {
		points[i].X *= cfg.ScaleX
		points[i].R *= cfg.ScaleR
	}

	return &profileModel{Config: cfg, Points: points, Segments: segments}, nil
}

// splinePoint, cubicSpline, fitCubicSpline, and cubicHermiteTerm are
// PROFILE.C's own spline()/seval(), duplicated verbatim from spline's
// own fitter (see spline.go) rather than shared, matching this
// project's established per-program convention — profile fits a fresh
// spline over just one segment's own points at a time, rather than
// spline's whole-dataset fit, so the two are still separate call
// sites even though the underlying algorithm is identical.
type splinePoint struct{ X, Y float64 }

type cubicSpline struct {
	points []splinePoint
	u, p   []float64
}

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

func cubicHermiteTerm(a float64) float64 { return a*a*a - a }

// Eval evaluates the fitted spline at z, matching SPLINE.C's own
// seval() exactly; z outside the segment's own x range is an error.
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

// linInterp is PROFILE.C's own lin(): the line through (x1,y1) and
// (x2,y2), evaluated at x.
func linInterp(x1, y1, x2, y2, x float64) float64 {
	a := (y2 - y1) / (x2 - x1)
	b := 0.5 * (y1 + y2 - a*(x1+x2))
	return a*x + b
}

// findRadius returns the profile's radius at distance x from the
// outboard end of stock, matching PROFILE.C's own findr() exactly: if
// x falls within a declared spline segment, a fresh cubic spline fit
// over just that segment's own points is evaluated; otherwise the two
// bracketing points are linearly interpolated.
func (m *profileModel) findRadius(x float64) (float64, error) {
	for _, seg := range m.Segments {
		if x < m.Points[seg.Start].X || x > m.Points[seg.End].X {
			continue
		}
		pts := make([]splinePoint, 0, seg.End-seg.Start+1)
		for i := seg.Start; i <= seg.End; i++ {
			pts = append(pts, splinePoint{X: m.Points[i].X, Y: m.Points[i].R})
		}
		s, err := fitCubicSpline(pts)
		if err != nil {
			return 0, fmt.Errorf("spline segment [%d,%d]: %w", seg.Start, seg.End, err)
		}
		return s.Eval(x)
	}

	k := -1
	for i, p := range m.Points {
		if x <= p.X {
			k = i - 1
			break
		}
	}
	if k < 0 {
		k = 0
	}
	return linInterp(m.Points[k].X, m.Points[k].R, m.Points[k+1].X, m.Points[k+1].R, x), nil
}

// cleanNegativeZero replaces a negative-zero float (produced here by
// negating an exact 0.0 axial position) with plain positive zero, so
// it prints as "0.000" rather than "-0.000" — matching PROFILE.OUT's
// own first row.
func cleanNegativeZero(v float64) float64 {
	if v == 0 {
		return 0
	}
	return v
}

// cutStep is one roughing pass.
type cutStep struct {
	Index         int
	AxialPosition float64
	Diameter      float64
	DepthOfCut    float64
}

// computeCuttingSchedule matches PROFILE.C's own draw() cutting-depth
// loop exactly — identical in structure to ogive's and egg's own
// schedules, but querying the profile's own findRadius instead of a
// closed-form shape formula.
func computeCuttingSchedule(m *profileModel) ([]cutStep, error) {
	cfg := m.Config
	rs := cfg.stockRadius()
	lastX := m.Points[len(m.Points)-1].X

	var steps []cutStep
	for i := 1; ; i++ {
		// Computed as a multiple of AxialStep rather than by
		// repeated subtraction: PROFILE.C's own loop (xl-=dxd each
		// pass) hits the same boundary-floating-point question at
		// the very last row, and repeated subtraction's particular
		// rounding drift comes out one row short of PROFILE.OUT's
		// own 76-row output; this form matches it exactly.
		xl := -float64(i-1) * cfg.AxialStep
		if xl < -lastX {
			break
		}
		xr := xl + cfg.ToolSize
		rt := 0.5 * cfg.ToolSize
		xm := xl + rt

		rl, err := m.findRadius(-xl)
		if err != nil {
			return nil, err
		}

		var doc float64
		if cfg.SquareTool {
			rr, rm := 0.0, 0.0
			if xr < 0 {
				if rr, err = m.findRadius(-xr); err != nil {
					return nil, err
				}
			}
			dmax := math.Max(rl, rr)
			if xm < 0 {
				if rm, err = m.findRadius(-xm); err != nil {
					return nil, err
				}
			}
			dmax = math.Max(dmax, rm)
			doc = rs - dmax
		} else {
			dmax := 0.0
			for phi := 0.0; phi <= 180; phi += 15 {
				xt := xm + rt*math.Cos(phi*rpd)
				yt := rt * math.Sin(phi*rpd)
				rr := 0.0
				if xt < 0 {
					if rr, err = m.findRadius(-xt); err != nil {
						return nil, err
					}
				}
				if rr+yt > dmax {
					dmax = rr + yt
				}
			}
			doc = rt + (rs - dmax)
		}

		steps = append(steps, cutStep{
			Index:         i,
			AxialPosition: cleanNegativeZero(-xl),
			Diameter:      2 * (rs - doc),
			DepthOfCut:    doc,
		})
	}
	return steps, nil
}

// renderProfile draws the turned profile and each roughing pass,
// matching PROFILE.C's own draw() exactly, including tick marks at
// each original data point (larger every 10th point).
func renderProfile(m *profileModel, steps []cutStep) (*svgplot.Canvas, error) {
	cfg := m.Config
	rs := cfg.stockRadius()
	lastX := m.Points[len(m.Points)-1].X

	xmin := -lastX
	xmax := 1.1 * cfg.ToolSize
	ymax := 2 * cfg.StockDiam
	ymin := -ymax
	wx := math.Abs(xmax - xmin)
	wy := math.Abs(ymax - ymin)
	if wx > 4*wy/3 {
		ymax = 3 * wx / 8
		ymin = -ymax
	} else {
		wx = 4 * wy / 3
		xmin = xmax - wx
	}

	c := svgplot.New(xmin, ymin, xmax, ymax, 640, 480)
	c.Box(xmin, ymin, xmax, ymax, 15, false)
	c.DashedLine(xmin, 0, xmax, 0, 15)
	c.DashedLine(xmin, rs, xmax, rs, 12)
	c.DashedLine(xmin, -rs, xmax, -rs, 12)
	c.DashedLine(0, -rs, 0, 1.5*rs, 12)
	c.Text(0, 1.5*rs, "x=0", 2)
	c.Text(xmin+0.05*wx, 3.0*rs, fmt.Sprintf("stock diam = %.3f", cfg.StockDiam), 2)
	if cfg.SquareTool {
		c.Text(xmin+0.05*wx, 2.8*rs, fmt.Sprintf("tool width = %.3f", cfg.ToolSize), 2)
	} else {
		c.Text(xmin+0.05*wx, 2.8*rs, fmt.Sprintf("tool diameter = %.3f", cfg.ToolSize), 2)
	}
	c.Text(xmin+0.05*wx, 2.6*rs, fmt.Sprintf("axial step = %.3f", cfg.AxialStep), 2)

	xl, rl := -m.Points[0].X, m.Points[0].R
	for x := -cfg.AxialStep; x >= -lastX; x -= cfg.AxialStep {
		r, err := m.findRadius(-x)
		if err != nil {
			return nil, err
		}
		c.Line(xl, rl, x, r, 15)
		c.Line(xl, -rl, x, -r, 15)
		xl, rl = x, r
	}

	for i, p := range m.Points {
		size := 0.005
		if i%10 == 0 {
			size = 0.005 * 2.5
		}
		c.Line(-p.X, rs-size*wy, -p.X, rs+size*wy, 2)
		c.Line(-p.X, -size*wy, -p.X, size*wy, 2)
	}

	for _, step := range steps {
		xlPos := -step.AxialPosition
		xr := xlPos + cfg.ToolSize
		if cfg.SquareTool {
			dmax := rs - step.DepthOfCut
			if dmax < rs {
				c.Box(xlPos, -dmax, xr, -(dmax + rs), 14, false)
			}
		} else {
			rt := 0.5 * cfg.ToolSize
			dmax := rs - step.DepthOfCut + rt
			xm := xlPos + rt
			c.Circle(xm, -dmax, rt, 14, false)
		}
	}
	return c, nil
}

func main() {
	dataPath := flag.String("data", "", "path to a PROFILE-format data file (see MWKGo/profile/testdata/)")
	svgPath := flag.String("svg", "", "path to write an SVG rendering of the profile (optional)")
	flag.Parse()

	if *dataPath == "" {
		fmt.Println("usage: profile -data <file> [-svg <file>]")
		fmt.Println("see MWKGo/profile/testdata/ for expected formats")
		os.Exit(1)
	}

	f, err := os.Open(*dataPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "profile:", err)
		os.Exit(1)
	}
	m, err := loadProfile(f)
	f.Close()
	if err != nil {
		fmt.Fprintln(os.Stderr, "profile:", err)
		os.Exit(1)
	}

	fmt.Printf("stock diameter = %.3f\n", m.Config.StockDiam)
	if m.Config.SquareTool {
		fmt.Printf("tool width = %.3f\n", m.Config.ToolSize)
	} else {
		fmt.Printf("tool diameter = %.3f\n", m.Config.ToolSize)
	}
	fmt.Printf("axial cutting increment = %.3f\n", m.Config.AxialStep)
	fmt.Println("start with left side of tool against end of stock")
	fmt.Println("\nc=compound rest setting (axial feed)")
	fmt.Println("d=diameter of resulting cut")
	fmt.Println("doc=depth of cut")
	fmt.Println("\n   i       c       d     doc")
	fmt.Println()

	steps, err := computeCuttingSchedule(m)
	if err != nil {
		fmt.Fprintln(os.Stderr, "profile:", err)
		os.Exit(1)
	}
	for _, s := range steps {
		fmt.Printf("%4d   %5.3f   %5.3f   %5.3f\n", s.Index, s.AxialPosition, s.Diameter, s.DepthOfCut)
	}

	if *svgPath != "" {
		canvas, err := renderProfile(m, steps)
		if err != nil {
			fmt.Fprintln(os.Stderr, "profile:", err)
			os.Exit(1)
		}
		if err := os.WriteFile(*svgPath, []byte(canvas.String()), 0o644); err != nil {
			fmt.Fprintln(os.Stderr, "profile:", err)
			os.Exit(1)
		}
		fmt.Printf("\nWrote profile diagram to %s\n", *svgPath)
	}
}
