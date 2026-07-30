// feed solves for milling feed rate given any three of five related
// quantities: spindle speed, number of cutting edges, chip load per
// edge, and feed rate expressed as either in/min or in/rev. The
// underlying relationship is feed(in/min) = chip load * speed *
// edges, so knowing any three (with a couple of degenerate
// exceptions) determines the rest.
//
// Converted from FEED.C (M. W. Klotz, 2/01),
// WorkshopUtilities/feed.
package main

import (
	"errors"
	"fmt"
	"os"

	"mwkgo/internal/promptio"
)

// feedRateResult is the fully solved set of feed rate quantities. A
// zero value in the input to solveFeedRate means "unknown", matching
// the original program's own use of zero as "not entered".
type feedRateResult struct {
	Speed, Edges, Load, FeedPerMin, FeedPerRev float64
}

// errInsufficientFeedData reports that the given combination of
// known quantities does not determine the rest: either too few were
// given, or the specific combination given, of the four combinations
// this solver recognizes, does not include enough independent
// information.
var errInsufficientFeedData = errors.New("insufficient data for solution")

// solveFeedRate solves for the unknown feed rate quantities given
// whichever of speed, edges, load, feedPerMin, and feedPerRev are
// already known (zero meaning unknown). It recognizes the same four
// solvable combinations as the original program: all of
// speed/edges/load; speed and edges plus either feed figure; speed
// and load plus either feed figure; or edges, load, and feedPerMin
// specifically (not feedPerRev, which the original program never
// even prompts for in that case, since edges and load alone already
// fix feedPerRev without determining speed).
func solveFeedRate(speed, edges, load, feedPerMin, feedPerRev float64) (feedRateResult, error) {
	switch {
	case speed != 0 && edges != 0 && load != 0:
		feedPerMin = load * speed * edges
		feedPerRev = feedPerMin / speed

	case speed != 0 && edges != 0 && (feedPerRev != 0 || feedPerMin != 0):
		if feedPerMin != 0 {
			feedPerRev = feedPerMin / speed
		} else {
			feedPerMin = feedPerRev * speed
		}
		load = feedPerMin / (speed * edges)

	case speed != 0 && load != 0 && (feedPerRev != 0 || feedPerMin != 0):
		if feedPerMin != 0 {
			feedPerRev = feedPerMin / speed
		} else {
			feedPerMin = feedPerRev * speed
		}
		edges = feedPerMin / (load * speed)

	case edges != 0 && load != 0 && feedPerMin != 0:
		speed = feedPerMin / (load * edges)
		feedPerRev = feedPerMin / speed

	default:
		return feedRateResult{}, errInsufficientFeedData
	}

	return feedRateResult{Speed: speed, Edges: edges, Load: load, FeedPerMin: feedPerMin, FeedPerRev: feedPerRev}, nil
}

func main() {
	prompter, err := promptio.New(os.Stdin, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, "feed:", err)
		os.Exit(1)
	}

	fmt.Println("FEED RATE CALCULATIONS")
	fmt.Println()
	fmt.Println("Input whatever data you know; press return if not known.")
	fmt.Println("You must enter three data items to obtain a solution.")
	fmt.Println("If finding speed from edges, chip load and feed,")
	fmt.Println("feed must be entered in in/min.")
	fmt.Println()

	speed := prompter.Float("Speed (rpm)", 0)
	edges := prompter.Float("Number of cutting edges (edges/revolution)", 0)
	load := prompter.Float("Chip load (in/cutting edge)", 0)

	known := 0
	for _, v := range []float64{speed, edges, load} {
		if v != 0 {
			known++
		}
	}

	var feedPerMin, feedPerRev float64
	if known < 3 {
		feedPerMin = prompter.Float("Feed (in/min)", 0)
		if feedPerMin != 0 {
			known++
		}
	}
	if known < 3 && !(edges != 0 && load != 0) {
		feedPerRev = prompter.Float("Feed (in/revolution)", 0)
		if feedPerRev != 0 {
			known++
		}
	}

	result, err := solveFeedRate(speed, edges, load, feedPerMin, feedPerRev)
	if err != nil {
		fmt.Println()
		fmt.Fprintln(os.Stderr, "feed:", err)
		os.Exit(1)
	}

	fmt.Println()
	fmt.Printf("Cutting edges = %.0f\n", result.Edges)
	fmt.Printf("Speed = %.0f rpm\n", result.Speed)
	fmt.Printf("Chip load = %.4f in/edge\n", result.Load)
	fmt.Printf("Feed = %.4f in/min\n", result.FeedPerMin)
	fmt.Printf("Feed = %.4f in/rev\n", result.FeedPerRev)
}
