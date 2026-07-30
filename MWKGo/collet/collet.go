// collet computes the bore diameter for a cylindrical collet, sawn
// with a slot for each side of the polygonal stock it needs to
// clamp, following John Way's technique (Machinist's Workshop,
// June/July 2004) for holding square or hexagonal stock without a
// dedicated broach.
//
// Converted from COLLET.C (M. W. Klotz, 5/04),
// WorkshopUtilities/collet. Valid for stock with an even number of
// sides (4, 6, or 8 in practice), per the original program's own
// notes; odd side counts are not meaningful for this construction.
package main

import (
	"fmt"
	"math"
	"os"

	"mwkgo/internal/promptio"
)

// boreDiameter returns the collet bore diameter needed to clamp
// polygonal stock with the given number of sides and across-flats
// dimension, sawn with slots of slotWidth between the sides.
func boreDiameter(sides int, acrossFlats, slotWidth float64) float64 {
	halfAngleDeg := 180 / float64(sides)
	cosHalfAngle := math.Cos(halfAngleDeg * math.Pi / 180)
	sinHalfAngle := math.Sin(halfAngleDeg * math.Pi / 180)

	offset := 0.5*acrossFlats/cosHalfAngle - 0.5*slotWidth*sinHalfAngle/cosHalfAngle
	return 2 * math.Sqrt(0.25*slotWidth*slotWidth+offset*offset)
}

func main() {
	prompter, err := promptio.New(os.Stdin, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, "collet:", err)
		os.Exit(1)
	}

	fmt.Println()
	fmt.Println("CYLINDRICAL COLLETS FOR POLYGONAL STOCK")
	fmt.Println()

	sides := prompter.Int("Number of stock polygon sides", 6)
	acrossFlats := prompter.Float("Stock across flats dimension (in)", 3.0/16)
	slotWidth := prompter.Float("Collet slot width (in)", 0.045)

	fmt.Printf("\nRequired collet bore diameter = %.4f in\n", boreDiameter(sides, acrossFlats, slotWidth))
}
