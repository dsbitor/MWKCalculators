// ogive computes the tangent-ogive nose-cone profile (the same
// pointed-arch shape used for both Gothic architecture and rocket
// nose cones) for a chosen ogive diameter and length, and works out
// an incremental turning schedule for roughing it out with either a
// square-tipped or round-tipped tool, printing the schedule and
// optionally rendering the profile and cutting passes as an SVG
// diagram.
//
// The ogive's own dimensions and tooling setup are specific to one
// job — like ogive's own archive-mate profile, and calibrat/vrev/
// simul/curfit/colsort/spline/belt in earlier conversion groups —
// this program reads its input fresh from a file named on the
// command line each run, in the same STARTOFDATA/ENDOFDATA text
// format the original used, though with a different per-line syntax
// (key=value rather than positional fields).
//
// Converted from OGIVE.C (M. W. Klotz), WorkshopUtilities/ogive. Per
// ai/plans/c-to-go-conversion-plan.md's Tier 3 "Graphics scope"
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
	"strconv"
	"strings"

	"mwkgo/internal/legacydat"
	"mwkgo/internal/svgplot"
)

const rpd = math.Pi / 180

// ogiveConfig is everything OGIVE.DAT provides: stock and tool
// dimensions, the axial cutting increment, and the ogive's own
// diameter and length.
type ogiveConfig struct {
	StockDiam   float64
	SquareTool  bool
	ToolSize    float64 // tool tip width (square) or diameter (round)
	AxialStep   float64
	OgiveDiam   float64
	OgiveLength float64
}

func (c ogiveConfig) caliber() float64     { return c.OgiveLength / c.OgiveDiam }
func (c ogiveConfig) stockRadius() float64 { return 0.5 * c.StockDiam }

// loadConfig parses OGIVE.DAT's own key=value format (bracketed by
// STARTOFDATA/ENDOFDATA), matching OGIVE.C's own rdata() exactly,
// including checking keys by substring containment in the same order
// (ds, toolt, wt, dxd, ogdiam, oglength) rather than an exact-match
// parser — harmless here since none of the six key names is a
// substring of another.
func loadConfig(r io.Reader) (ogiveConfig, error) {
	rows, err := legacydat.Rows(r, legacydat.Options{StartMarker: "STARTOFDATA", EndMarker: "ENDOFDATA"})
	if err != nil {
		return ogiveConfig{}, fmt.Errorf("scan data: %w", err)
	}

	var cfg ogiveConfig
	for i, row := range rows {
		if !strings.Contains(row, "=") {
			continue
		}
		fields := legacydat.Fields(row, "=\t;")
		if len(fields) < 2 {
			return ogiveConfig{}, fmt.Errorf("line %d: %q does not split into a key and value", i+1, row)
		}
		value := strings.TrimSpace(fields[1])

		switch {
		case strings.Contains(row, "ds"):
			v, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return ogiveConfig{}, fmt.Errorf("line %d: ds %q: %w", i+1, value, err)
			}
			cfg.StockDiam = v
		case strings.Contains(row, "toolt"):
			v, err := strconv.Atoi(value)
			if err != nil {
				return ogiveConfig{}, fmt.Errorf("line %d: toolt %q: %w", i+1, value, err)
			}
			cfg.SquareTool = v != 0
		case strings.Contains(row, "wt"):
			v, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return ogiveConfig{}, fmt.Errorf("line %d: wt %q: %w", i+1, value, err)
			}
			cfg.ToolSize = v
		case strings.Contains(row, "dxd"):
			v, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return ogiveConfig{}, fmt.Errorf("line %d: dxd %q: %w", i+1, value, err)
			}
			cfg.AxialStep = v
		case strings.Contains(row, "ogdiam"):
			v, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return ogiveConfig{}, fmt.Errorf("line %d: ogdiam %q: %w", i+1, value, err)
			}
			cfg.OgiveDiam = v
		case strings.Contains(row, "oglength"):
			v, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return ogiveConfig{}, fmt.Errorf("line %d: oglength %q: %w", i+1, value, err)
			}
			cfg.OgiveLength = v
		}
	}
	return cfg, nil
}

// ogiveRadius returns the ogive profile's radius at distance x from
// the tip, matching OGIVE.C's own findr() exactly: beyond the ogive's
// own length, the profile has already blended into the stock
// cylinder, so the radius is just the stock radius.
func ogiveRadius(x float64, cfg ogiveConfig) float64 {
	if x > cfg.OgiveLength {
		return cfg.stockRadius()
	}
	c := cfg.caliber()
	d := cfg.OgiveDiam
	term := d * (c*c + 0.25)
	return math.Sqrt(term*term-(cfg.OgiveLength-x)*(cfg.OgiveLength-x)) - d*(c*c-0.25)
}

// cutStep is one stage of the roughing schedule.
type cutStep struct {
	Index         int
	AxialPosition float64 // compound rest setting, measured from the tip
	Diameter      float64 // resulting cut diameter
	DepthOfCut    float64
}

// computeCuttingSchedule works from the tip toward the stock in
// AxialStep increments, matching OGIVE.C's own draw() cutting-depth
// loop exactly: a square-tipped tool's depth is set by the deepest of
// its two edges and its center; a round-tipped tool's depth is found
// by scanning its own curved tip profile against the ogive shape.
func computeCuttingSchedule(cfg ogiveConfig) []cutStep {
	var steps []cutStep
	i := 1
	for xl := 0.0; xl >= -cfg.OgiveLength; xl -= cfg.AxialStep {
		xr := xl + cfg.ToolSize
		rt := 0.5 * cfg.ToolSize
		xm := xl + rt
		rl := ogiveRadius(-xl, cfg)

		var doc float64
		if cfg.SquareTool {
			var rr, rm float64
			if xr < 0 {
				rr = ogiveRadius(-xr, cfg)
			}
			dmax := math.Max(rl, rr)
			if xm < 0 {
				rm = ogiveRadius(-xm, cfg)
			}
			dmax = math.Max(dmax, rm)
			doc = cfg.stockRadius() - dmax
		} else {
			dmax := 0.0
			for phi := 0.0; phi <= 180; phi += 15 {
				xt := xm + rt*math.Cos(phi*rpd)
				yt := rt * math.Sin(phi*rpd)
				var rr float64
				if xt < 0 {
					rr = ogiveRadius(-xt, cfg)
				}
				if rr+yt > dmax {
					dmax = rr + yt
				}
			}
			doc = rt + (cfg.stockRadius() - dmax)
		}

		steps = append(steps, cutStep{
			Index:         i,
			AxialPosition: -xl,
			Diameter:      2 * (cfg.stockRadius() - doc),
			DepthOfCut:    doc,
		})
		i++
	}
	return steps
}

// renderOgive draws the ogive profile and each roughing pass,
// matching OGIVE.C's own draw() exactly.
func renderOgive(cfg ogiveConfig, steps []cutStep) *svgplot.Canvas {
	rs := cfg.stockRadius()
	xmax := 1.1 * cfg.ToolSize
	xmin := -1.1 * cfg.OgiveLength
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

	xl, rl := 0.0, 0.0
	for x := -cfg.AxialStep; x >= xmin; x -= cfg.AxialStep {
		r := ogiveRadius(-x, cfg)
		c.Line(xl, rl, x, r, 15)
		c.Line(xl, -rl, x, -r, 15)
		xl, rl = x, r
	}

	for _, step := range steps {
		xlPos := -step.AxialPosition
		xr := xlPos + cfg.ToolSize
		rt := 0.5 * cfg.ToolSize
		xm := xlPos + rt
		if cfg.SquareTool {
			dmax := rs - step.DepthOfCut
			if dmax < rs {
				c.Box(xlPos, -dmax, xr, -(dmax + rs), 14, false)
			}
		} else {
			// doc = rt + (rs - dmax), so dmax = rs - doc + rt.
			dmax := rs - step.DepthOfCut + rt
			c.Circle(xm, -dmax, rt, 14, false)
		}
	}
	return c
}

func main() {
	dataPath := flag.String("data", "", "path to an OGIVE-format data file (see MWKGo/ogive/testdata/example.dat)")
	svgPath := flag.String("svg", "", "path to write an SVG rendering of the ogive profile (optional)")
	flag.Parse()

	if *dataPath == "" {
		fmt.Println("usage: ogive -data <file> [-svg <file>]")
		fmt.Println("see MWKGo/ogive/testdata/example.dat for the expected format")
		os.Exit(1)
	}

	f, err := os.Open(*dataPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "ogive:", err)
		os.Exit(1)
	}
	cfg, err := loadConfig(f)
	f.Close()
	if err != nil {
		fmt.Fprintln(os.Stderr, "ogive:", err)
		os.Exit(1)
	}

	fmt.Printf("stock diameter = %.3f\n", cfg.StockDiam)
	if cfg.SquareTool {
		fmt.Printf("tool width = %.3f\n", cfg.ToolSize)
	} else {
		fmt.Printf("tool diameter = %.3f\n", cfg.ToolSize)
	}
	fmt.Printf("axial cutting increment = %.3f\n", cfg.AxialStep)
	fmt.Printf("ogive diameter = %.3f\n", cfg.OgiveDiam)
	fmt.Printf("ogive length = %.3f\n", cfg.OgiveLength)
	fmt.Printf("ogive caliber = %.3f\n", cfg.caliber())
	fmt.Println("\nstart with left side of tool against end of stock")
	fmt.Println("\nc=compound rest setting (axial feed)")
	fmt.Println("d=diameter of resulting cut")
	fmt.Println("doc=depth of cut")
	fmt.Println("\n   i       c       d     doc")
	fmt.Println()

	steps := computeCuttingSchedule(cfg)
	for _, s := range steps {
		fmt.Printf("%4d   %5.3f   %5.3f   %5.3f\n", s.Index, s.AxialPosition, s.Diameter, s.DepthOfCut)
	}

	if *svgPath != "" {
		if err := os.WriteFile(*svgPath, []byte(renderOgive(cfg, steps).String()), 0o644); err != nil {
			fmt.Fprintln(os.Stderr, "ogive:", err)
			os.Exit(1)
		}
		fmt.Printf("\nWrote ogive profile diagram to %s\n", *svgPath)
	}
}
