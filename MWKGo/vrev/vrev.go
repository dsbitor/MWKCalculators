// vrev lays out holes of varying diameter evenly around a bolt
// circle, given the desired spacing between adjacent hole edges: the
// tap-guide and multi-diameter punch-holder problem, where a fixed
// hole spacing (rather than a fixed angular spacing) is wanted
// despite the holes being different sizes.
//
// The hole list is specific to one job (the particular set of holes
// being laid out this time), not universal reference data or
// reusable equipment configuration, so — like calibrat in the same
// conversion group — this program reads its input fresh from a file
// named on the command line each run, in the same
// STARTOFDATA/ENDOFDATA text format the original used; see
// ai/plans/c-to-go-conversion-plan.md, "Data-file strategy for
// Tier 2".
//
// Converted from VREV.C (M. W. Klotz), WorkshopUtilities/vrev.
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
	"mwkgo/internal/promptio"
)

// maxSearchRefinements bounds decSearch's zoom-in passes (the same
// technique used by cseg and flywheel elsewhere in this project):
// each pass narrows the search window by 10x, so convergence to
// float64 precision needs nowhere near this many passes; it exists
// only so the loop has an enforced maximum.
const maxSearchRefinements = 200

// decSearch searches [low,high] for the x minimizing |f(x)-target|,
// to within tol, by repeatedly scanning ten steps across the current
// window and then zooming in around the best point found.
func decSearch(low, high, target, tol float64, f func(float64) float64) float64 {
	xs, xe := low, high
	dx := 0.1 * (high - low)
	var bestX float64
	for pass := 0; pass < maxSearchRefinements; pass++ {
		bestDiff := math.Inf(1)
		for x := xs; x <= xe; x += dx {
			if diff := math.Abs(f(x) - target); diff < bestDiff {
				bestDiff, bestX = diff, x
			}
		}
		if bestDiff <= tol {
			break
		}
		xs, xe = bestX-dx, bestX+dx
		dx *= 0.1
	}
	return bestX
}

// angleExcess returns the total angle (in degrees, minus 360) that
// holes of the given diameters, each separated by an edge-to-edge
// arc length of spacing, would subtend on a bolt circle of the given
// radius. It is zero exactly when radius is big enough to fit every
// hole with that spacing and no more.
func angleExcess(radius, spacing float64, holeDiameters []float64) float64 {
	phi := 180 / math.Pi * spacing / radius
	sum := 0.0
	for _, d := range holeDiameters {
		sum += 2*mwktrig.ClampedAsinDeg(0.5*d/radius) + phi
	}
	return sum - 360
}

// minBoltCircleDiameter returns the smallest bolt circle diameter
// that fits every hole with at least the given edge-to-edge spacing.
func minBoltCircleDiameter(spacing float64, holeDiameters []float64) float64 {
	radius := decSearch(0.01, 1000, 0, 0.0001, func(r float64) float64 {
		return angleExcess(r, spacing, holeDiameters)
	})
	return 2 * radius
}

// actualSpacing returns the edge-to-edge spacing that results from
// distributing the given holes evenly around a bolt circle of the
// given diameter (which may be larger than the minimum, in which
// case the spacing grows to fill the extra room).
func actualSpacing(boltCircleDiameter float64, holeDiameters []float64) float64 {
	radius := boltCircleDiameter / 2
	tangentSum := 0.0
	for _, d := range holeDiameters {
		tangentSum += 2 * mwktrig.ClampedAsinDeg(0.5*d/radius)
	}
	return radius * (math.Pi / 180) * (360 - tangentSum) / float64(len(holeDiameters))
}

// holePosition is one hole's location on the bolt circle.
type holePosition struct {
	Diameter float64
	X, Y     float64
	AngleDeg float64
}

// layoutHoles places each hole around a bolt circle of the given
// diameter and spacing, starting the first hole at angle zero and
// walking around by each hole's own half-width, the gap, and the
// next hole's half-width.
func layoutHoles(boltCircleDiameter, spacing float64, holeDiameters []float64) []holePosition {
	radius := boltCircleDiameter / 2
	phi := 180 / math.Pi * spacing / radius

	positions := make([]holePosition, len(holeDiameters))
	theta := 0.0
	for i, d := range holeDiameters {
		rad := theta * math.Pi / 180
		positions[i] = holePosition{Diameter: d, X: radius * math.Cos(rad), Y: radius * math.Sin(rad), AngleDeg: theta}
		if i < len(holeDiameters)-1 {
			theta += mwktrig.ClampedAsinDeg(0.5*d/radius) + mwktrig.ClampedAsinDeg(0.5*holeDiameters[i+1]/radius) + phi
		}
	}
	return positions
}

func loadHoleDiameters(r io.Reader) ([]float64, error) {
	rows, err := legacydat.Rows(r, legacydat.Options{StartMarker: "STARTOFDATA", EndMarker: "ENDOFDATA"})
	if err != nil {
		return nil, fmt.Errorf("scan data: %w", err)
	}

	holes := make([]float64, 0, len(rows))
	for i, row := range rows {
		fields := legacydat.Fields(row, ",\t;")
		if len(fields) != 1 {
			return nil, fmt.Errorf("line %d: %q does not contain exactly one diameter", i+1, row)
		}
		d, err := strconv.ParseFloat(strings.TrimSpace(fields[0]), 64)
		if err != nil {
			return nil, fmt.Errorf("line %d: diameter %q: %w", i+1, fields[0], err)
		}
		holes = append(holes, d)
	}
	return holes, nil
}

func main() {
	dataPath := flag.String("holes", "", "path to a hole-diameter data file (see MWKGo/vrev/testdata/example.dat)")
	flag.Parse()

	if *dataPath == "" {
		fmt.Println("usage: vrev -holes <file>")
		fmt.Println("see MWKGo/vrev/testdata/example.dat for the expected format")
		os.Exit(1)
	}

	f, err := os.Open(*dataPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "vrev:", err)
		os.Exit(1)
	}
	holes, err := loadHoleDiameters(f)
	f.Close()
	if err != nil {
		fmt.Fprintln(os.Stderr, "vrev:", err)
		os.Exit(1)
	}

	prompter, err := promptio.New(os.Stdin, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, "vrev:", err)
		os.Exit(1)
	}

	fmt.Printf("%d hole diameters read from data file\n\n", len(holes))

	hdMin, hdMax := holes[0], holes[0]
	for _, d := range holes {
		hdMin = math.Min(hdMin, d)
		hdMax = math.Max(hdMax, d)
	}
	fmt.Printf("Minimum hole size = %.4f\n", hdMin)
	fmt.Printf("Maximum hole size = %.4f\n", hdMax)

	spacing := prompter.Float("Desired hole-to-hole spacing", 0.1)
	minBCD := minBoltCircleDiameter(spacing, holes)
	fmt.Printf("Minimum bolt circle diameter = %.4f\n", minBCD)

	var bcd float64
	for {
		bcd = prompter.Float("Desired bolt circle diameter", minBCD)
		if bcd >= minBCD {
			break
		}
		fmt.Printf("bcd must be >= %.4f\n", minBCD)
		fmt.Println("Try again, dummy.")
	}

	spacing = actualSpacing(bcd, holes)
	fmt.Printf("For a bcd = %.4f, the hole-to-hole spacing will be %.4f\n", bcd, spacing)

	fmt.Printf("\nBolt circle diameter = %.4f\n", bcd)
	fmt.Printf("Hole-to-hole spacing = %.4f\n\n", spacing)
	fmt.Println("N, hole diameter, x, y, angle(deg))")
	fmt.Println()
	for i, p := range layoutHoles(bcd, spacing, holes) {
		fmt.Printf("%2d, %.4f, %+.4f, %+.4f, %.4f\n", i, p.Diameter, p.X, p.Y, p.AngleDeg)
	}

	fmt.Printf("\nTheoretical minimum stock diameter = %.4f\n", bcd+hdMax)
	fmt.Printf("Recommended stock diameter = %.4f\n", bcd+hdMax+2*spacing)
}
