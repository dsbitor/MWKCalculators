// tenon computes the depth of cut needed to mill a regular polygonal
// tenon on the end of cylindrical stock, using a rotary indexer.
//
// Converted from TENON.C (M. W. Klotz, 6/00), WorkshopUtilities/tenon.
package main

import (
	"fmt"
	"math"
	"os"

	"mwkgo/internal/promptio"
)

// circumscribedDiameterFromAcrossFlats returns the diameter of the
// circle that just circumscribes a regular polygon of the given
// number of sides, given its across-flats dimension. Only meaningful
// for an even number of sides, where "across flats" is well defined.
func circumscribedDiameterFromAcrossFlats(sides int, acrossFlats float64) float64 {
	halfAngleDeg := 0.5 * (360 / float64(sides))
	return acrossFlats / math.Cos(halfAngleDeg*math.Pi/180)
}

// tenonDepthOfCut returns the depth of cut needed to mill a regular
// polygon of the given number of sides and circumscribed diameter on
// stock of the given diameter.
func tenonDepthOfCut(stockDiameter float64, sides int, circumscribedDiameter float64) float64 {
	halfAngleDeg := 0.5 * (360 / float64(sides))
	flatToCenterDistance := 0.5 * circumscribedDiameter * math.Cos(halfAngleDeg*math.Pi/180)
	return 0.5*stockDiameter - flatToCenterDistance
}

func main() {
	prompter, err := promptio.New(os.Stdin, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, "tenon:", err)
		os.Exit(1)
	}

	fmt.Println("DEPTH OF CUT FOR REGULAR POLYGONAL TENON")
	fmt.Println()

	stockDiameter := prompter.Float("Diameter of stock", 0.5)
	sides := prompter.Int("Number of sides on tenon", 5)
	anglePerSide := 360.0 / float64(sides)

	var circumscribedDiameter float64
	if sides%2 == 0 {
		acrossFlats := prompter.Float("Distance across flats", 0.25)
		circumscribedDiameter = circumscribedDiameterFromAcrossFlats(sides, acrossFlats)
	} else {
		circumscribedDiameter = prompter.Float("Diameter of circle circumscribing tenon", 0.25)
	}

	depthOfCut := tenonDepthOfCut(stockDiameter, sides, circumscribedDiameter)

	fmt.Println()
	fmt.Printf("Stock diameter = %.4f\n", stockDiameter)
	fmt.Printf("Number of flats on tenon = %d\n", sides)
	fmt.Printf("Angle between adjacent tenon flats = %.2f deg\n", anglePerSide)
	fmt.Printf("Diameter of circle circumscribing tenon = %.4f\n", circumscribedDiameter)
	fmt.Printf("Depth of cut = %.4f\n", depthOfCut)
}
