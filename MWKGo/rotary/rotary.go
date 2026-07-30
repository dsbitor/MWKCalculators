// rotary computes the angular position of each division of a
// rotary table, both as a decimal degree value and as
// degrees/minutes/seconds, given the number of equal divisions
// wanted around a full circle.
//
// Converted from ROTARY.C (M. W. Klotz), WorkshopUtilities/rotary.
// The original writes its results to ROTARY.OUT via fopen; this
// conversion prints to stdout instead and drops that DOS
// file-save-then-page convenience, the same approach used for loan
// (Tier 1 suitability review, Finding 5).
package main

import (
	"fmt"
	"math"
	"os"

	"mwkgo/internal/promptio"
)

// division is one rotary table division's angular position.
type division struct {
	Index            int
	DecimalDeg       float64
	Degrees, Minutes int
	Seconds          int
}

// rotaryDivisions returns the angular position of each of
// numDivisions equal divisions around a full circle, from 0 up to
// and including the full circle (which wraps back to 0 degrees).
// Unlike compound's degreesMinutesSeconds, seconds here are rounded
// to the nearest whole second (matching the original program's own
// RND, rather than truncated), so the minute-carry branch is
// genuinely reachable: it fires for numDivisions=25, division 7.
func rotaryDivisions(numDivisions int) []division {
	step := 360.0 / float64(numDivisions)
	divisions := make([]division, numDivisions+1)

	for k := 0; k <= numDivisions; k++ {
		theta := float64(k) * step
		degrees := int(math.Floor(theta))
		minutes := int(math.Floor((theta - float64(degrees)) * 60))
		seconds := int(math.Floor((theta-float64(degrees)-float64(minutes)/60)*3600 + 0.5))

		if seconds == 60 {
			seconds = 0
			minutes++
		}
		if minutes == 60 {
			minutes = 0
			degrees++
		}
		if degrees == 360 {
			degrees = 0
		}

		divisions[k] = division{Index: k, DecimalDeg: theta, Degrees: degrees, Minutes: minutes, Seconds: seconds}
	}
	return divisions
}

func main() {
	prompter, err := promptio.New(os.Stdin, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, "rotary:", err)
		os.Exit(1)
	}

	fmt.Println("ROTARY TABLE CALCULATIONS")
	fmt.Println()

	numDivisions := prompter.Int("Number of divisions", 13)

	fmt.Printf("Number of divisions = %d\n\n", numDivisions)
	fmt.Printf("%8s    %8s  %6s  %6s  %6s\n", "DIVISION", "degdec", "deg", "min", "sec")
	for _, d := range rotaryDivisions(numDivisions) {
		fmt.Printf("%8d    %8.4f  %6d  %6d  %6d\n", d.Index, d.DecimalDeg, d.Degrees, d.Minutes, d.Seconds)
	}
}
