// crod computes gudgeon (wrist) pin position for a slider-crank
// mechanism (an engine's connecting rod and crank) across a half
// revolution, both relative to the crank center and relative to its
// top-dead-center position, printing a table and optionally rendering
// both curves as an SVG diagram.
//
// Converted from CROD.C (M. W. Klotz), WorkshopUtilities/crod. Per
// ai/plans/c-to-go-conversion-plan.md's Tier 3 "Graphics scope"
// resolution, DOS graphics calls are replaced by internal/svgplot;
// the original's interactive mouse-driven menu and click-to-read-
// coordinates feature are dropped entirely (no equivalent in a
// static SVG image), and the .OUT file-save feature is replaced by
// printing straight to stdout. Note: CROD.TXT itself refers to the
// output file as "CROD.DAT", but the source actually writes
// "CROD.OUT" — a documentation/code mismatch in the original, not
// reproduced here since neither name is meaningful once stdout
// replaces the file entirely.
package main

import (
	"flag"
	"fmt"
	"math"
	"os"

	"mwkgo/internal/promptio"
	"mwkgo/internal/svgplot"
)

const rpd = math.Pi / 180

// crankPoint is the gudgeon pin position at one crank angle.
type crankPoint struct {
	ThetaDeg        float64
	XCenterRelative float64 // measured relative to the crank center
	ZTDCRelative    float64 // measured relative to the top-dead-center position
}

// gudgeonPinPosition is the standard slider-crank displacement
// formula, matching CROD.C's own computation exactly: throw is the
// crank radius, rodLen the connecting rod length.
func gudgeonPinPosition(throw, rodLen, thetaDeg float64) float64 {
	theta := thetaDeg * rpd
	st, ct := math.Sin(theta), math.Cos(theta)
	return throw*ct + math.Sqrt(rodLen*rodLen-throw*throw*st*st)
}

// computeCrodTable walks the crank from 0 (top dead center) to 180
// degrees (bottom dead center) in dtheta increments, matching CROD.C's
// own loop exactly.
func computeCrodTable(throw, rodLen, dtheta float64) []crankPoint {
	var points []crankPoint
	for theta := 0.0; theta <= 180.0; theta += dtheta {
		x := gudgeonPinPosition(throw, rodLen, theta)
		points = append(points, crankPoint{ThetaDeg: theta, XCenterRelative: x, ZTDCRelative: throw + rodLen - x})
	}
	return points
}

// renderCrod draws both gudgeon pin position curves (relative to
// crank center, color 14; relative to TDC, color 12) against a
// dotted-gridline crank-angle axis, matching CROD.C's own graph()
// exactly.
func renderCrod(throw, rodLen float64, points []crankPoint) *svgplot.Canvas {
	xmin, xmax := -20.0, 190.0
	ymin, ymax := -0.25*throw, 1.2*(throw+rodLen)

	c := svgplot.New(xmin, ymin, xmax, ymax, 640, 480)
	c.Box(xmin, ymin, xmax, ymax, 15, false)
	c.DashedLine(xmin, 0, xmax, 0, 15)
	c.DashedLine(0, ymin, 0, ymax, 15)
	for theta := 0.0; theta <= 180.0; theta += 15 {
		c.DashedLine(theta, 0, theta, ymax, 15)
	}
	c.DashedLine(0, rodLen+throw, 180, rodLen+throw, 15)
	c.DashedLine(0, rodLen-throw, 180, rodLen-throw, 15)
	c.DashedLine(0, rodLen-2*throw, 180, rodLen-2*throw, 15)

	z := -0.002 * (ymax - ymin)
	c.Text(0, z, "0", 7)
	c.Text(45-3, z, "45", 7)
	c.Text(90-3, z, "90", 7)
	c.Text(135-4, z, "135", 7)
	c.Text(180-4, z, "180", 7)
	c.Text(-15, rodLen+throw, fmt.Sprintf("%.1f", rodLen+throw), 7)
	c.Text(-15, rodLen-throw, fmt.Sprintf("%.1f", rodLen-throw), 7)
	c.Text(-15, rodLen-2*throw, fmt.Sprintf("%.1f", rodLen-2*throw), 7)

	tl, zl := 0.0, throw+rodLen
	for _, p := range points {
		c.Line(tl, zl, p.ThetaDeg, p.XCenterRelative, 14)
		tl, zl = p.ThetaDeg, p.XCenterRelative
	}
	tl, zl = 0.0, 0.0
	for _, p := range points {
		c.Line(tl, zl, p.ThetaDeg, p.ZTDCRelative, 12)
		tl, zl = p.ThetaDeg, p.ZTDCRelative
	}
	return c
}

func main() {
	svgPath := flag.String("svg", "", "path to write an SVG rendering of the crank diagram (optional)")
	flag.Parse()

	prompter, err := promptio.New(os.Stdin, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, "crod:", err)
		os.Exit(1)
	}

	fmt.Println("Units of measurement don't matter but must be consistent")
	rodLen := prompter.Float("Connecting rod length", 2.4)
	fmt.Println("Radius measured from crank center to connecting rod driver center")
	throw := prompter.Float("Crank radius", 0.6)
	dtheta := prompter.Float("Angular increment (deg)", 5.0)

	points := computeCrodTable(throw, rodLen, dtheta)

	fmt.Println("\nGudgeon Pin Position")
	fmt.Printf("Connecting rod length = %.4f\n", rodLen)
	fmt.Printf("Crank throw = %.4f\n", throw)
	fmt.Printf("Angular increment = %.4f deg\n\n", dtheta)
	fmt.Println("A = crank angle (deg) (0=tdc,180=bdc)")
	fmt.Println("X = gudgeon pin position (measured relative to crank center)")
	fmt.Println("Z = gudgeon pin position (measured relative to tdc)")
	fmt.Println()
	fmt.Printf("%9s%9s%9s\n\n", "A", "X", "Z")
	for _, p := range points {
		fmt.Printf("%9.3f%9.3f%9.3f\n", p.ThetaDeg, p.XCenterRelative, p.ZTDCRelative)
	}

	if *svgPath != "" {
		if err := os.WriteFile(*svgPath, []byte(renderCrod(throw, rodLen, points).String()), 0o644); err != nil {
			fmt.Fprintln(os.Stderr, "crod:", err)
			os.Exit(1)
		}
		fmt.Printf("\nWrote crank diagram to %s\n", *svgPath)
	}
}
