// dallow computes the drill tip allowance: the extra depth a drilled
// hole must go to give a flat-bottomed hole of the intended depth,
// once the drill's conical tip is accounted for.
//
// Converted from DALLOW.C (M. W. Klotz, 5/99), WorkshopUtilities/dallow.
package main

import (
	"fmt"
	"math"
	"os"

	"mwkgo/internal/promptio"
)

// tipAllowance returns the extra drilling depth needed for a drill
// of the given diameter and included tip angle, measured from the
// point where the drill's conical tip first touches the material to
// the point where the full diameter is reached.
func tipAllowance(includedAngleDegrees, diameter float64) float64 {
	halfAngleTan := math.Tan(0.5 * includedAngleDegrees * math.Pi / 180)
	return 0.5 * diameter / halfAngleTan
}

func main() {
	prompter, err := promptio.New(os.Stdin, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, "dallow:", err)
		os.Exit(1)
	}

	fmt.Println("DRILL TIP ALLOWANCE COMPUTATION")
	fmt.Println()

	angle := prompter.Float("Included angle of drill tip (deg)", 118)
	diameter := prompter.Float("Drill diameter (in)", 0.5)

	fmt.Printf("\nAllowance for drill tip = %.4f in\n", tipAllowance(angle, diameter))
}
