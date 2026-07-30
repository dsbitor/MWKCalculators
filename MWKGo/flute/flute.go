// flute computes the setup for milling a tapered flute (as on an
// engine support column) by tilting the workpiece and plunge cutting
// with a ball end mill: the depth of cut at each end of the flute
// and the workpiece inclination angle needed to produce the taper
// over a given flute length.
//
// Converted from FLUTE.C (M. W. Klotz), WorkshopUtilities/flute.
package main

import (
	"fmt"
	"math"
	"os"

	"mwkgo/internal/promptio"
)

// flutePlan holds the depth of cut at each end of the flute and the
// workpiece inclination angle needed to connect them over the
// flute's length.
type flutePlan struct {
	DepthAtSmallEnd float64
	DepthAtLargeEnd float64
	InclinationDeg  float64
}

// degreesMinutesSeconds splits a decimal angle in degrees into whole
// degrees, minutes, and seconds, matching the original program's dms
// function, including its carry from 60 seconds into the next
// minute and from 60 minutes into the next degree. As with
// compound's identically named function, truncation keeps computed
// seconds and minutes below 60 for every input checked, so the
// carry branch is not exercised by any test here, but it is kept for
// parity with the original.
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

// computeFlutePlan returns the flute cutting plan for a ball mill of
// radius millRadius, flute radii rs (small end) and rf (large end),
// over a flute of length length. rs and rf must each be no greater
// than millRadius.
func computeFlutePlan(millRadius, rs, rf, length float64) flutePlan {
	depthSmall := millRadius - math.Sqrt(millRadius*millRadius-rs*rs)
	depthLarge := millRadius - math.Sqrt(millRadius*millRadius-rf*rf)
	inclination := math.Atan((depthLarge-depthSmall)/length) * 180 / math.Pi

	return flutePlan{
		DepthAtSmallEnd: depthSmall,
		DepthAtLargeEnd: depthLarge,
		InclinationDeg:  inclination,
	}
}

func main() {
	prompter, err := promptio.New(os.Stdin, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, "flute:", err)
		os.Exit(1)
	}

	var millRadius, rs, rf float64
	for {
		millDiameter := prompter.Float("Ball mill diameter (in)", 0.5)
		millRadius = 0.5 * millDiameter
		rs = prompter.Float("Flute radius at small end (in)", 0.1)
		rf = prompter.Float("Flute radius at large end (in)", 0.2)
		if rs <= millRadius && rf <= millRadius {
			break
		}
		fmt.Printf("\nFlute radii must be less than ball mill radius = %.4f\n", millRadius)
	}

	length := prompter.Float("Length of flute (in)", 3.0)

	plan := computeFlutePlan(millRadius, rs, rf, length)

	fmt.Printf("\nDepth of cut at small end of flute = %.4f in\n", plan.DepthAtSmallEnd)
	fmt.Printf("Depth of cut at large end of flute = %.4f in\n", plan.DepthAtLargeEnd)
	fmt.Printf("Required workpiece inclination angle = %.4f deg = %s\n", plan.InclinationDeg, formatDMS(plan.InclinationDeg))
}
