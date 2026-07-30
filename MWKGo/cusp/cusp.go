// cusp computes the spacing between successive passes of a ball end
// mill needed to keep the cusp left between passes below a desired
// height.
//
// Converted from CUSP.C (M. W. Klotz, 10/02), WorkshopUtilities/cusp.
package main

import (
	"fmt"
	"math"
	"os"

	"mwkgo/internal/promptio"
)

// passSpacing returns the spacing between successive milling passes
// that leaves a cusp of at most cuspHeight between them, for a ball
// end mill of the given diameter. This is the standard circular
// segment sagitta-to-chord relationship: for a circle of radius r,
// a sagitta (cusp height) of c corresponds to a chord (pass spacing)
// of 2*sqrt(2rc-c^2). cuspHeight must not exceed the mill's radius;
// beyond that the cusp height would exceed what the ball profile can
// produce and the quantity under the square root goes negative,
// which the original program does not guard against either.
func passSpacing(millDiameter, cuspHeight float64) float64 {
	radius := 0.5 * millDiameter
	return 2 * math.Sqrt(2*radius*cuspHeight-cuspHeight*cuspHeight)
}

func main() {
	prompter, err := promptio.New(os.Stdin, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, "cusp:", err)
		os.Exit(1)
	}

	millDiameter := prompter.Float("Ball mill diameter (in)", 0.25)
	cuspHeight := prompter.Float("Desired cusp height (in)", 0.001)

	fmt.Printf("Spacing between successive cuts = %.4f in\n", passSpacing(millDiameter, cuspHeight))
}
