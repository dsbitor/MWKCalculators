// compound computes the lathe compound rest angle needed to achieve
// a given infeed-to-slide movement ratio, so the compound's
// graduated dial reads directly in a convenient fraction of diameter
// reduction.
//
// Converted from COMPOUND.C (M. W. Klotz, 6/01, 1/05),
// WorkshopUtilities/compound.
package main

import (
	"fmt"
	"math"
	"os"

	"mwkgo/internal/mwktrig"
	"mwkgo/internal/promptio"
)

// compoundAngleDeg returns the compound rest angle, in degrees, that
// gives the requested infeed-to-slide movement ratio. ratio must be
// at most 1; the original program re-prompts rather than computing
// an angle for a ratio greater than 1, since asin(ratio) has no real
// solution beyond that point without mwktrig's clamping kicking in
// and silently producing 90 degrees.
func compoundAngleDeg(ratio float64) float64 {
	return mwktrig.ClampedAsinDeg(ratio)
}

// degreesMinutesSeconds splits a decimal angle in degrees into whole
// degrees, minutes, and seconds, matching the original program's dms
// function, including its carry from 60 seconds into the next
// minute and from 60 minutes into the next degree. That carry is
// defensive: truncation (matching the original's own int-from-float
// assignment) keeps computed seconds and minutes below 60 for every
// input checked, so the carry branch is not exercised by any test
// here, but it is kept for parity with the original.
func degreesMinutesSeconds(angleDeg float64) (sign string, degrees, minutes, seconds int) {
	sign = ""
	if angleDeg < 0 {
		sign = "-"
	}

	magnitude := math.Abs(angleDeg)
	degrees = int(math.Floor(magnitude))
	fractionalSeconds := (magnitude - math.Floor(magnitude)) * 3600
	minutes = int(fractionalSeconds / 60)
	seconds = int(fractionalSeconds) - minutes*60

	if seconds == 60 {
		seconds = 0
		minutes++
	}
	if minutes == 60 {
		minutes = 0
		degrees++
	}
	return sign, degrees, minutes, seconds
}

func formatDMS(angleDeg float64) string {
	sign, degrees, minutes, seconds := degreesMinutesSeconds(angleDeg)
	return fmt.Sprintf("%s%d deg, %d min, %d sec", sign, degrees, minutes, seconds)
}

func main() {
	prompter, err := promptio.New(os.Stdin, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, "compound:", err)
		os.Exit(1)
	}

	fmt.Println("SETTING COMPOUND REST ANGLE")
	fmt.Println()

	var ratio float64
	for {
		ratio = prompter.Float("Required movement ratio", 0.1)
		if ratio <= 1 {
			break
		}
		fmt.Println("Ratio must be <= 1")
	}

	angle := compoundAngleDeg(ratio)
	fmt.Printf("Set compound angle to %.4f deg = %s\n", angle, formatDMS(angle))
	fmt.Printf("Complement of this angle = %.4f deg = %s\n", 90-angle, formatDMS(90-angle))
}
