// eccentub computes the tube diameter needed to turn an offset
// eccentric in a 3-jaw chuck, given the parent stock diameter and
// the required eccentric offset.
//
// Converted from ECCENTUB.C (M. W. Klotz, 02/01), WorkshopUtilities/eccent.
package main

import (
	"fmt"
	"math"
	"os"

	"mwkgo/internal/promptio"
)

// tubeDiameter returns the diameter of tube stock needed to turn an
// eccentric of the given offset from parent stock of parentDiameter.
// A zero offset needs no tube at all: the parent stock itself is
// used directly.
func tubeDiameter(parentDiameter, offset float64) float64 {
	radius := 0.5 * parentDiameter
	innerRadius := radius - offset
	return 2 * math.Sqrt(7*radius*radius-9*radius*innerRadius+3*innerRadius*innerRadius)
}

func main() {
	prompter, err := promptio.New(os.Stdin, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, "eccentub:", err)
		os.Exit(1)
	}

	fmt.Println("TUBE SIZE FOR TURNING ECCENTRICS")
	fmt.Println()

	parentDiameter := prompter.Float("Diameter of parent stock (in)", 1.0)
	offset := prompter.Float("Required eccentric offset (in)", 0.1)

	fmt.Printf("\nDiameter of required tube = %.4f in\n", tubeDiameter(parentDiameter, offset))
}
