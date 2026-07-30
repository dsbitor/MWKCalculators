// gearpa estimates a spur gear's pressure angle (14.5 or 20
// degrees, the two most common) by comparing a measured chordal span
// across several teeth to the span predicted for each candidate
// angle, per the chordal-span formula in Machinery's Handbook.
//
// Converted from GEARPA.C, WorkshopUtilities/gearpa.
package main

import (
	"fmt"
	"math"
	"os"

	"mwkgo/internal/promptio"
)

// standardPressureAngles are the two pressure angles this program
// checks, in the order the original program checks them.
var standardPressureAngles = [2]float64{14.5, 20.0}

// spanTeethCount returns the number of teeth to span with calipers
// for a gear with n teeth at the given pressure angle (identified by
// its index into standardPressureAngles), per Machinery's Handbook's
// table for this measurement, and whether n falls within the range
// that table covers for this pressure angle.
//
// The 14.5 degree table is ported as a sequence of overwrites, per
// the original program's own sequential (not else-if) range checks,
// rather than as a clean set of non-overlapping bins: the original's
// fifth and sixth conditions (51-62 and 53-75) overlap, and since
// each later match overwrites the previous one, the actual range
// attributed to 5 teeth is only 51-52, with 53-75 going to 6 teeth.
// This is almost certainly an unintended typo in the original (63
// for 53 would make every range contiguous, matching the clean
// pattern the 20 degree table already follows), but it's what the
// original program actually computes and prints, so it's preserved
// here rather than silently corrected.
func spanTeethCount(n int, pressureAngleIndex int) (teeth int, inRange bool) {
	if pressureAngleIndex == 0 { // 14.5 degrees
		if n >= 12 && n <= 18 {
			teeth = 2
		}
		if n >= 19 && n <= 37 {
			teeth = 3
		}
		if n >= 38 && n <= 50 {
			teeth = 4
		}
		if n >= 51 && n <= 62 {
			teeth = 5
		}
		if n >= 53 && n <= 75 {
			teeth = 6
		}
		if n >= 76 && n <= 87 {
			teeth = 7
		}
		if n >= 88 && n <= 100 {
			teeth = 8
		}
		if n >= 101 && n <= 110 {
			teeth = 9
		}
		return teeth, teeth != 0
	}

	// 20 degrees: a clean, non-overlapping table, valid only up to
	// 81 teeth (unlike the 14.5 degree table, which the outer
	// 12-110 gate in main covers completely).
	switch {
	case n >= 12 && n <= 18:
		return 2, true
	case n >= 19 && n <= 27:
		return 3, true
	case n >= 28 && n <= 36:
		return 4, true
	case n >= 37 && n <= 45:
		return 5, true
	case n >= 46 && n <= 54:
		return 6, true
	case n >= 55 && n <= 63:
		return 7, true
	case n >= 64 && n <= 72:
		return 8, true
	case n >= 73 && n <= 81:
		return 9, true
	default:
		return 0, false
	}
}

// chordalSpan returns the chordal span measured across teethSpanned
// teeth for a gear with n teeth, diametral pitch diametralPitch, and
// the given pressure angle in degrees, per Machinery's Handbook's
// chordal-span formula.
func chordalSpan(n int, diametralPitch, pressureAngleDeg float64, teethSpanned int) float64 {
	pitchRadius := 0.5 * float64(n) / diametralPitch
	toothThickness := 0.5 * math.Pi / diametralPitch
	spacesSpanned := float64(teethSpanned - 1)

	angleRad := pressureAngleDeg * math.Pi / 180
	involuteFactor := 2 * (math.Tan(angleRad) - angleRad)

	return pitchRadius * math.Cos(angleRad) * (toothThickness/pitchRadius + 2*math.Pi*spacesSpanned/float64(n) + involuteFactor)
}

func main() {
	prompter, err := promptio.New(os.Stdin, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, "gearpa:", err)
		os.Exit(1)
	}

	fmt.Println("GEAR PRESSURE ANGLE BY MEASUREMENT OF CHORDAL SPAN")
	fmt.Println()

	n := prompter.Int("Number of teeth on gear", 30)
	diametralPitch := prompter.Float("Diametral pitch of gear", 6.0)

	if n < 12 || n > 110 {
		fmt.Println("\nNot able to calculate number of teeth to measure across")
		return
	}

	fmt.Println()
	for i, pa := range standardPressureAngles {
		fmt.Printf("Pressure angle = %.2f deg...", pa)
		teeth, inRange := spanTeethCount(n, i)
		if !inRange {
			fmt.Printf("Number of teeth out of range for %.1f deg pa.\n", pa)
			continue
		}
		span := chordalSpan(n, diametralPitch, pa, teeth)
		fmt.Printf("Chordal span over %d teeth = %.4f in\n", teeth, span)
	}
}
