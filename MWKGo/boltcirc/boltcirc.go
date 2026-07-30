// boltcirc computes the polar and Cartesian coordinates of the
// holes in a bolt circle: an equally-spaced ring of holes at a fixed
// radius, given the hole count, circle radius, hole diameter, an
// angular offset for the first hole, and an offset for the circle's
// center.
//
// Converted from BOLTCIRC.C (M. W. Klotz, 11/98, 1/00),
// WorkshopUtilities/boltcirc. The original writes its results to
// BOLTCIRC.DAT via fopen; this conversion prints to stdout instead
// and drops that DOS file-save-then-page convenience, the same
// approach used for loan (Tier 1 suitability review, Finding 5).
package main

import (
	"fmt"
	"math"
	"os"

	"mwkgo/internal/promptio"
)

// boltHole is one hole's angular position and Cartesian coordinates.
type boltHole struct {
	AngleDeg float64
	X, Y     float64
}

// boltCircleLayout returns the edge-to-edge spacing between adjacent
// holes and the position of every hole, given the hole count,
// circle radius, hole diameter, angular offset of the first hole,
// and offset of the circle's center.
func boltCircleLayout(numHoles int, radius, holeDiameter, angularOffsetDeg, xOffset, yOffset float64) (spacing float64, holes []boltHole) {
	step := 360.0 / float64(numHoles)
	spacing = 2*radius*math.Sin(0.5*step*math.Pi/180) - holeDiameter

	holes = make([]boltHole, numHoles)
	for i := 0; i < numHoles; i++ {
		angle := angularOffsetDeg + float64(i)*step
		holes[i] = boltHole{
			AngleDeg: angle,
			X:        radius*math.Cos(angle*math.Pi/180) + xOffset,
			Y:        radius*math.Sin(angle*math.Pi/180) + yOffset,
		}
	}
	return spacing, holes
}

func main() {
	prompter, err := promptio.New(os.Stdin, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, "boltcirc:", err)
		os.Exit(1)
	}

	fmt.Println("BOLTCIRCLE COMPUTATIONS")
	fmt.Println()

	numHoles := prompter.Int("Number of holes", 5)
	radius := prompter.Float("Radius of bolt circle", 1.0)
	holeDiameter := prompter.Float("Diameter of bolt holes", 0.5)
	angularOffset := prompter.Float("Angular offset of first hole (deg)", 0.0)
	xOffset := prompter.Float("X offset of bolt circle center", 0.0)
	yOffset := prompter.Float("Y offset of bolt circle center", 0.0)

	spacing, holes := boltCircleLayout(numHoles, radius, holeDiameter, angularOffset, xOffset, yOffset)
	if spacing < 0 {
		fmt.Println("\nWARNING: HOLES WILL OVERLAP!")
	}

	fmt.Println("\nBoltcircle specification:")
	fmt.Printf("Radius of bolt circle = %.4f\n", radius)
	fmt.Printf("Bolt hole diameter = %.4f\n", holeDiameter)
	fmt.Printf("Spacing between hole edges = %.4f\n", spacing)
	fmt.Printf("Angular offset of first hole = %.4f deg\n", angularOffset)
	fmt.Printf("X offset of bolt circle center = %.4f\n", xOffset)
	fmt.Printf("Y offset of bolt circle center = %.4f\n\n", yOffset)
	fmt.Println("HOLE       ANGLE     X-COORD     Y-COORD")
	fmt.Println()
	for i, hole := range holes {
		fmt.Printf("%4d  %10.4f  %10.4f  %10.4f\n", i+1, hole.AngleDeg, hole.X, hole.Y)
	}
}
