// chain computes the center-to-center distance needed to mount two
// sprockets for a roller chain of a given length and pitch.
//
// Converted from CHAIN.C (M. W. Klotz, 2/03),
// WorkshopUtilities/sprocket.
package main

import (
	"fmt"
	"math"
	"os"

	"mwkgo/internal/promptio"
)

// centerDistance returns the center-to-center distance for two
// sprockets of largeTeeth and smallTeeth teeth, joined by a chain of
// the given pitch and length (in pitches). This is the standard
// roller chain center-distance formula.
func centerDistance(pitch, chainLengthPitches float64, largeTeeth, smallTeeth int) float64 {
	toothDifference := float64(largeTeeth - smallTeeth)
	toothSum := 2*chainLengthPitches - float64(largeTeeth) - float64(smallTeeth)

	return (pitch / 8) * (toothSum + math.Sqrt(toothSum*toothSum-0.810*toothDifference*toothDifference))
}

func main() {
	prompter, err := promptio.New(os.Stdin, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, "chain:", err)
		os.Exit(1)
	}

	fmt.Println()
	fmt.Println("SPROCKET CENTER DISTANCE GIVEN CHAIN LENGTH")
	fmt.Println()

	pitch := prompter.Float("Chain pitch (in)", 1.0)
	length := prompter.Float("Chain length (pitches)", 48.0)
	largeTeeth := prompter.Int("Number of teeth in large sprocket", 18)
	smallTeeth := prompter.Int("Number of teeth in small sprocket", 9)

	fmt.Printf("\nSprocket center-to-center distance = %.4f in\n", centerDistance(pitch, length, largeTeeth, smallTeeth))
}
