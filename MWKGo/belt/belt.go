// belt computes the total belt length for an arbitrary closed-loop
// arrangement of pulleys (flat belt only, no V-belt pitch-diameter
// correction), given each pulley's center location, diameter, and
// whether it sits inside the belt loop (a normal driver/driven
// pulley) or outside it (an idler, wrapped on its far side).
//
// The pulley layout is specific to one job's own arrangement, not
// universal reference data or reusable equipment configuration — so,
// like calibrat, vrev, simul, curfit, and colsort in earlier
// conversion groups, this program reads its input fresh from a file
// named on the command line each run, in the same
// STARTOFDATA/ENDOFDATA text format the original used; see
// ai/plans/c-to-go-conversion-plan.md, "Data-file strategy for
// Tier 2".
//
// Converted from BELT.C (M. W. Klotz), WorkshopUtilities/belt. That
// archive also bundles three companion programs solving related,
// simpler two-pulley problems: QBELT (quick belt length, no data
// file needed), PULLEY (solve for the second pulley's diameter given
// belt length and separation), and PCD (solve for the pulley
// separation given both diameters and belt length) — each converted
// alongside this one; see MWKGo/qbelt, MWKGo/pulley, MWKGo/pcd.
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

// pulley is one pulley's own geometry: its center location, its
// diameter, and whether the belt wraps its inside (a normal
// driver/driven pulley) or outside (an idler).
type pulley struct {
	X, Y, Diam float64
	Inside     bool
}

// pulleyResult is one pulley's computed wrap angle, wrap length
// (radius times wrap angle), and the straight span of belt from this
// pulley to the next one in the loop.
type pulleyResult struct {
	WrapAngleDeg float64
	WrapLength   float64
	SpanToNext   float64
}

func ioSign(inside bool) float64 {
	if inside {
		return 1
	}
	return -1
}

// atan2Zero is ATAN2.C's own ATAN2 macro: math.Atan2, except it
// returns 0 for (0,0) instead of an implementation-defined angle.
func atan2Zero(s, c float64) float64 {
	if s == 0 && c == 0 {
		return 0
	}
	return math.Atan2(s, c)
}

// computeBeltLayout computes each pulley's wrap angle, wrap length,
// and span to the next pulley, plus the total belt length, matching
// BELT.C's own main() calculation exactly: pulleys sharing the same
// "inside"/"outside" side use the internal-tangent formula (the belt
// crosses between them) while pulleys on opposite sides use the
// external-tangent formula, and a wrap angle can come out negative
// (flagging the pulley as possibly slack) when the belt's span to it
// is long enough that the "short way around" is geometrically
// possible.
func computeBeltLayout(pulleys []pulley) (results []pulleyResult, totalLength float64, slack bool, err error) {
	n := len(pulleys)
	spn := make([]float64, n)
	agl := make([]float64, n)

	for i := 0; i < n; i++ {
		j := i + 1
		if j == n {
			j = 0
		}
		dx := pulleys[j].X - pulleys[i].X
		dy := pulleys[j].Y - pulleys[i].Y
		csq := dx*dx + dy*dy
		ctc := math.Sqrt(csq)
		spz := 0.5 * (pulleys[i].Diam + pulleys[j].Diam)
		if spz > ctc || ctc == 0 {
			return nil, 0, false, fmt.Errorf("pulleys %d and %d overlap", i+1, j+1)
		}
		if pulleys[i].Inside == pulleys[j].Inside {
			spz = 0.5 * (pulleys[i].Diam - pulleys[j].Diam)
		}
		spn[i] = math.Sqrt(csq - spz*spz)

		a := atan2Zero(dy, dx) + ioSign(pulleys[i].Inside)*mwktrig.ClampedAsin(spz/ctc)
		for a > 2*math.Pi {
			a -= 2 * math.Pi
		}
		for a < 0 {
			a += 2 * math.Pi
		}
		agl[i] = a
	}

	results = make([]pulleyResult, n)
	for j := 0; j < n; j++ {
		i := j - 1
		if j == 0 {
			i = n - 1
		}
		arc := ioSign(pulleys[j].Inside) * (agl[j] - agl[i])
		if arc < 0 {
			arc += 2 * math.Pi
		}
		if arc > math.Pi {
			if spn[j] > 0.5*pulleys[j].Diam*math.Tan(0.5*(2*math.Pi-arc)) {
				arc -= 2 * math.Pi
				slack = true
			}
		}
		wrapLen := arc * 0.5 * pulleys[j].Diam
		results[j] = pulleyResult{WrapAngleDeg: arc * 180 / math.Pi, WrapLength: wrapLen, SpanToNext: spn[j]}
		totalLength += wrapLen + spn[j]
	}
	return results, totalLength, slack, nil
}

// loadPulleys parses pulley rows (x, y, diameter, +1/-1 inside flag,
// comma-separated, bracketed by STARTOFDATA/ENDOFDATA) from r, in
// file order — pulley order is significant (counter-clockwise around
// the belt path), unlike most other Tier 2 data files, so this is not
// sorted.
func loadPulleys(r io.Reader) ([]pulley, error) {
	rows, err := legacydat.Rows(r, legacydat.Options{StartMarker: "STARTOFDATA", EndMarker: "ENDOFDATA"})
	if err != nil {
		return nil, fmt.Errorf("scan data: %w", err)
	}

	pulleys := make([]pulley, 0, len(rows))
	for i, row := range rows {
		fields := legacydat.Fields(row, ",\t;")
		if len(fields) < 4 {
			return nil, fmt.Errorf("line %d: %q does not split into x, y, diameter, and flag", i+1, row)
		}
		x, err := strconv.ParseFloat(strings.TrimSpace(fields[0]), 64)
		if err != nil {
			return nil, fmt.Errorf("line %d: x %q: %w", i+1, fields[0], err)
		}
		y, err := strconv.ParseFloat(strings.TrimSpace(fields[1]), 64)
		if err != nil {
			return nil, fmt.Errorf("line %d: y %q: %w", i+1, fields[1], err)
		}
		diam, err := strconv.ParseFloat(strings.TrimSpace(fields[2]), 64)
		if err != nil {
			return nil, fmt.Errorf("line %d: diameter %q: %w", i+1, fields[2], err)
		}
		flag, err := strconv.Atoi(strings.TrimSpace(fields[3]))
		if err != nil {
			return nil, fmt.Errorf("line %d: flag %q: %w", i+1, fields[3], err)
		}
		pulleys = append(pulleys, pulley{X: x, Y: y, Diam: diam, Inside: flag == 1})
	}
	return pulleys, nil
}

func main() {
	dataPath := flag.String("data", "", "path to a pulley layout data file (see MWKGo/belt/testdata/example.dat)")
	flag.Parse()

	if *dataPath == "" {
		fmt.Println("usage: belt -data <file>")
		fmt.Println("see MWKGo/belt/testdata/example.dat for the expected format")
		os.Exit(1)
	}

	f, err := os.Open(*dataPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "belt:", err)
		os.Exit(1)
	}
	pulleys, err := loadPulleys(f)
	f.Close()
	if err != nil {
		fmt.Fprintln(os.Stderr, "belt:", err)
		os.Exit(1)
	}

	fmt.Println("BELT LENGTH CALCULATION")
	fmt.Printf("Number of pulleys for which data was read = %d\n", len(pulleys))

	results, total, slack, err := computeBeltLayout(pulleys)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("\ndata displayed for each pulley is:")
	fmt.Println("x, y, diam, flag, wrap angle(deg), wrap length, span to next pulley")
	for i, p := range pulleys {
		flag := -1
		if p.Inside {
			flag = 1
		}
		fmt.Printf("\npulley number %d\n", i+1)
		fmt.Printf("%.3f, %.3f, %.3f, %d, %.1f, %.3f, %.3f\n",
			p.X, p.Y, p.Diam, flag, results[i].WrapAngleDeg, results[i].WrapLength, results[i].SpanToNext)
	}
	if slack {
		fmt.Println("\nNegative wrap length means belt may go slack!")
	}
	fmt.Printf("\nTotal belt length = %.3f\n", total)
}
