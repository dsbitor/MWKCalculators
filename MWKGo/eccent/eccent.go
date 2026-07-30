// eccent computes the packing thickness needed to turn an eccentric
// in a 3-jaw chuck, by shimming one jaw with packing material to
// shift the workpiece's turning axis away from the spindle axis.
//
// Converted from ECCENT.C (M. W. Klotz, 11/98, 9/01),
// WorkshopUtilities/eccent. The later ECCENTUB.C, converted
// separately as eccentub, replaces this packing method with a slotted
// tube, which avoids the awkward jaw-width measurement this method
// requires.
package main

import (
	"errors"
	"fmt"
	"math"
	"os"

	"mwkgo/internal/promptio"
)

var sqrt3 = math.Sqrt(3)

// packingThickness returns the thickness of packing material needed
// behind one chuck jaw to turn an eccentric of the given offset in a
// workpiece of the given diameter, held by jaws of the given width.
// It returns an error if the workpiece would fall through the
// unpacked jaws at the requested offset, in which case no packing
// thickness can achieve it.
func packingThickness(jawWidth, workpieceDiameter, offset float64) (float64, error) {
	radius := 0.5 * workpieceDiameter

	if offset > (radius+jawWidth)/sqrt3 {
		return 0, errors.New("work will fall through unpacked jaws")
	}

	if jawWidth > sqrt3*offset {
		return 1.5 * offset, nil
	}

	return 1.5*offset - radius + 0.5*math.Sqrt(4*radius*radius-3*offset*offset+2*offset*jawWidth*sqrt3-jawWidth*jawWidth), nil
}

func main() {
	prompter, err := promptio.New(os.Stdin, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, "eccent:", err)
		os.Exit(1)
	}

	fmt.Println("THREE JAW ECCENTRIC PACKING CALCULATION")
	fmt.Println()

	jawWidth := prompter.Float("Width of chuck jaws", 0.125)
	workpieceDiameter := prompter.Float("Diameter of workpiece", 1.5625)
	offset := prompter.Float("Required eccentric offset", 0.28125)

	packing, err := packingThickness(jawWidth, workpieceDiameter, offset)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("\nRequired Packing Size = %.4f in\n", packing)
}
