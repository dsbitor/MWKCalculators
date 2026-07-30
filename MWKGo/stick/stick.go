// stick implements Peter Lott's technique (Machinist's Workshop,
// 6/01, pg. 18) for chasing threads on a lathe when the half-nuts
// must be disengaged and reengaged partway through: it finds the
// shortest carriage repositioning distance that is simultaneously a
// whole number of leadscrew pitches and a whole number of the thread
// pitch being cut, so reengaging the half-nuts does not lose
// registration on the thread. Covers all four combinations of
// metric or Imperial leadscrew with metric or Imperial thread.
//
// Converted from STICK.C (M. W. Klotz, 5/01), WorkshopUtilities/stick.
package main

import (
	"errors"
	"fmt"
	"math"
	"os"

	"mwkgo/internal/promptio"
)

// maxStickSearchIterations bounds the search for a carriage
// repositioning distance. Realistic leadscrew and thread pitch
// combinations converge within at most a few hundred iterations;
// this generous bound catches a pathological combination that would
// otherwise search indefinitely, the role the original program's own
// keypress-based manual bailout played.
const maxStickSearchIterations = 1_000_000

// nearInteger reports whether x is within tolerance of a whole
// number.
func nearInteger(x, tolerance float64) bool {
	_, fraction := math.Modf(x)
	fraction = math.Abs(fraction)
	return fraction < tolerance || fraction > 1-tolerance
}

// stickLength finds the shortest carriage repositioning distance, in
// inches, that is simultaneously a whole number of leadscrew pitches
// and a whole number of thread pitches. leadscrewFactor and
// threadFactor are both threads per inch: given directly for an
// Imperial leadscrew or thread, or as 25.4/pitchMM for a metric one.
// It returns an error if no such distance is found within
// maxIterations steps.
func stickLength(leadscrewFactor, threadFactor float64, maxIterations int) (float64, error) {
	const tolerance = 0.000001
	step := 1 / leadscrewFactor

	stick := step
	for i := 0; i < maxIterations; i++ {
		if nearInteger(threadFactor*stick, tolerance) {
			return stick, nil
		}
		stick += step
	}
	return 0, errors.New("no common repositioning distance found within the search limit")
}

func main() {
	prompter, err := promptio.New(os.Stdin, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, "stick:", err)
		os.Exit(1)
	}

	fmt.Println("THREAD CHASING OFFSETS")
	fmt.Println()

	leadscrewChoice := prompter.Line("Type of leadscrew thread (M)etric or [I]mperial? ")
	var leadscrewFactor float64
	if leadscrewChoice == "m" || leadscrewChoice == "M" {
		pitch := prompter.Float("Leadscrew pitch (mm)", 4.0)
		leadscrewFactor = 25.4 / pitch
	} else {
		pitch := prompter.Float("Leadscrew pitch (tpi)", 10.0)
		leadscrewFactor = pitch
	}

	threadChoice := prompter.Line("Type of thread [M]etric or (I)mperial? ")
	var threadFactor float64
	if threadChoice == "i" || threadChoice == "I" {
		pitch := prompter.Float("Pitch of thread being cut (tpi)", 15.0)
		threadFactor = pitch
	} else {
		pitch := prompter.Float("Pitch of thread being cut (mm)", 4.0)
		threadFactor = 25.4 / pitch
	}

	stick, err := stickLength(leadscrewFactor, threadFactor, maxStickSearchIterations)
	if err != nil {
		fmt.Fprintln(os.Stderr, "stick:", err)
		os.Exit(1)
	}

	fmt.Printf("\nMove carriage %.3f in = %.3f mm to right to recapture thread\n", stick, 25.4*stick)
}
