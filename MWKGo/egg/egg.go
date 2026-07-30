// egg computes the profile of an asymmetrical-ellipse "egg" shape (a
// tribute, per EGG.TXT, to a scene in Nevil Shute's novel "Trustee
// from the Toolroom" in which a model engineer turns metal eggs on a
// lathe) and an incremental turning schedule for roughing it out,
// printing the schedule and optionally rendering the profile and
// cutting passes as an SVG diagram.
//
// Unlike every other Tier 3 program, egg takes no input at all —
// EGG.C hardwires its own stock diameter, tool width, axial step,
// and egg shape parameters (semi-axes and asymmetry factor), and
// EGG.TXT says so explicitly: "I hardwired these values into the
// program. They can't be changed without recompiling." This
// conversion keeps that behavior rather than inventing configurable
// input the original never had.
//
// Converted from EGG.C (M. W. Klotz), WorkshopUtilities/egg. Per
// ai/plans/c-to-go-conversion-plan.md's Tier 3 "Graphics scope"
// resolution, DOS graphics calls are replaced by internal/svgplot;
// the original's interactive mouse-driven menu and click-to-read-
// coordinates feature are dropped entirely, and the .OUT file-save
// feature is replaced by printing straight to stdout. EGG.C's own
// dead cubic-spline machinery (a spline()/seval() pair identical to
// spline's own, declared and fully implemented but never actually
// called anywhere in the file) is omitted from this port entirely
// rather than translated.
package main

import (
	"flag"
	"fmt"
	"math"
	"os"

	"mwkgo/internal/svgplot"
)

// Hardwired parameters, matching EGG.C's own init() exactly.
const (
	stockDiam = 1.0 // ds
	toolWidth = 1.0 / 16.0
	axialStep = 0.01 // dxd
)

// eggSemiAxes returns the egg shape's own semi-major axis, semi-minor
// axis, and asymmetry factor, all derived from the stock radius,
// matching EGG.C's own findr() exactly: b=rs, a=1.5*b, k=0.6 (the
// values EGG.TXT documents as chosen "by fiddling" to look right).
func eggSemiAxes(stockRadius float64) (a, b, k float64) {
	b = stockRadius
	a = 1.5 * b
	k = 0.6
	return a, b, k
}

// eggShape returns the y coordinate of the asymmetrical-ellipse egg
// curve x²/a² + (1−k·x)·y²/b² = 1, solved for y, matching EGG.C's own
// egg() exactly.
func eggShape(x, a, b, k float64) float64 {
	return b * math.Sqrt((1-x*x/(a*a))/(1-k*x))
}

// findRadius returns the egg profile's radius at distance x from the
// pointed end of the stock, matching EGG.C's own findr() exactly:
// beyond twice the semi-major axis, the profile has flattened out to
// a small nub (rend) rather than the egg formula itself, which would
// otherwise become undefined there.
func findRadius(x, stockRadius float64) float64 {
	a, b, k := eggSemiAxes(stockRadius)
	if x < 2*a {
		return eggShape(a-x, a, b, k)
	}
	return 0.1 * stockRadius
}

// eggDimensions bundles the hardwired parameters init() computes from
// the three constants above.
type eggDimensions struct {
	StockRadius   float64
	OverallLength float64 // oal
	PlotEnd       float64 // xend
}

func computeDimensions() eggDimensions {
	rs := 0.5 * stockDiam
	oal := 1.75 * stockDiam
	return eggDimensions{StockRadius: rs, OverallLength: oal, PlotEnd: oal + 0.125}
}

// cutStep is one roughing pass, printed only when material actually
// needs removing there (the profile hasn't already reached full
// stock radius).
type cutStep struct {
	Index         int
	AxialPosition float64
	DepthOfCut    float64
	Diameter      float64
	GapBefore     bool // one or more steps were skipped before this one
}

// computeCuttingSchedule matches EGG.C's own draw() cutting-depth loop
// exactly: at each axial position, the tool's required depth is set
// by whichever of its two edges or its center would cut deepest into
// the egg profile; a step is only recorded if material actually needs
// removing there (r < stock radius), and a gap marker notes any run
// of skipped steps in between two recorded ones.
func computeCuttingSchedule(dim eggDimensions) []cutStep {
	var steps []cutStep
	lastPrinted := 0
	i := 1
	for xl := 0.0; xl >= -dim.PlotEnd; xl -= axialStep {
		xr := xl + toolWidth
		xm := xl + 0.5*toolWidth
		rl := findRadius(-xl, dim.StockRadius)
		rr := 0.0
		if xr < 0 {
			rr = findRadius(-xr, dim.StockRadius)
		}
		r := math.Max(rl, rr)
		rm := 0.0
		if xm < 0 {
			rm = findRadius(-xm, dim.StockRadius)
		}
		r = math.Max(r, rm)

		if r < dim.StockRadius {
			steps = append(steps, cutStep{
				Index:         i,
				AxialPosition: -xl,
				DepthOfCut:    dim.StockRadius - r,
				Diameter:      2 * r,
				GapBefore:     i-lastPrinted > 1,
			})
			lastPrinted = i
		}
		i++
	}
	return steps
}

// renderEgg draws the egg profile and each roughing pass, matching
// EGG.C's own draw() exactly.
func renderEgg(dim eggDimensions, steps []cutStep) *svgplot.Canvas {
	rs := dim.StockRadius
	xmin := -dim.PlotEnd
	xmax := 1.1 * toolWidth
	ymax := 2 * stockDiam
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
	c.Line(-dim.OverallLength, -rs, -dim.OverallLength, rs, 15)
	c.Text(0, 1.5*rs, "x=0", 2)
	c.Text(xmin+0.05*wx, 3.0*rs, fmt.Sprintf("stock diam = %.3f in", stockDiam), 2)
	c.Text(xmin+0.05*wx, 2.8*rs, fmt.Sprintf("tool width = %.3f in", toolWidth), 2)
	c.Text(xmin+0.05*wx, 2.6*rs, fmt.Sprintf("axial step = %.3f in", axialStep), 2)

	xl, rl := 0.0, 0.0
	for x := -axialStep; x >= -dim.PlotEnd; x -= axialStep {
		r := findRadius(-x, rs)
		c.Line(xl, rl, x, r, 15)
		c.Line(xl, -rl, x, -r, 15)
		xl, rl = x, r
	}

	for _, step := range steps {
		xlPos := -step.AxialPosition
		xr := xlPos + toolWidth
		r := 0.5 * step.Diameter
		c.Box(xlPos, -r, xr, -(r + rs), 14, false)
	}
	return c
}

func main() {
	svgPath := flag.String("svg", "", "path to write an SVG rendering of the egg profile (optional)")
	flag.Parse()

	dim := computeDimensions()

	fmt.Printf("stock diameter = %.3f in\n", stockDiam)
	fmt.Printf("tool width = %.3f in\n", toolWidth)
	fmt.Printf("axial cutting increment = %.3f in\n", axialStep)
	fmt.Println("start with left side of tool against end of stock")
	fmt.Println("\nc=compound rest setting (axial feed)")
	fmt.Println("d=diameter of resulting cut")
	fmt.Println("doc=depth of cut")
	fmt.Println("\n   i       c     doc       d")
	fmt.Println()

	steps := computeCuttingSchedule(dim)
	for _, s := range steps {
		if s.GapBefore {
			fmt.Println("  ...")
		}
		fmt.Printf("%4d   %5.3f   %5.3f   %5.3f\n", s.Index, s.AxialPosition, s.DepthOfCut, s.Diameter)
	}

	if *svgPath != "" {
		if err := os.WriteFile(*svgPath, []byte(renderEgg(dim, steps).String()), 0o644); err != nil {
			fmt.Fprintln(os.Stderr, "egg:", err)
			os.Exit(1)
		}
		fmt.Printf("\nWrote egg profile diagram to %s\n", *svgPath)
	}
}
