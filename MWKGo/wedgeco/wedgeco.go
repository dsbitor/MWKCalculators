// wedgeco computes the volume of a conical wedge: the smaller of
// the two pieces formed when a cone is sliced by a plane parallel to
// its axis, offset from the axis by a given sagitta measured along a
// diameter of the cone's base. The exact volume requires a triple
// integral with no simple closed form; this program implements the
// resulting formula directly.
//
// Converted from WEDGECO.C (M. W. Klotz), Math/wedgeco.
package main

import (
	"fmt"
	"math"
	"os"

	"mwkgo/internal/promptio"
)

// wedgeVolume returns the volume of the conical wedge cut from a
// cone of height h and base radius r by a plane parallel to the
// axis, offset so that the sagitta (distance from the base
// circumference to the plane, measured along a diameter) is s.
// s must be strictly between 0 and 2r; s == r (the plane through the
// axis) is handled as the exact half-cone case, since the general
// formula is singular there.
func wedgeVolume(h, s, r float64) float64 {
	coneVolume := math.Pi * r * r * h / 3

	if s == r {
		return 0.5 * coneVolume
	}

	offset := r - s
	if offset < 0 {
		offset = -offset
	}
	cosAngle := offset / r
	angle := math.Acos(cosAngle)
	sinAngle := math.Sin(angle)
	tanAngle := math.Tan(angle)

	smallerPiece := h * r * r * (angle - 2*sinAngle*cosAngle + cosAngle*cosAngle*cosAngle*math.Log(tanAngle+1/cosAngle)) / 3

	if s <= r {
		return smallerPiece
	}
	return coneVolume - smallerPiece
}

func main() {
	prompter, err := promptio.New(os.Stdin, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, "wedgeco:", err)
		os.Exit(1)
	}

	fmt.Println("COMPUTE VOLUME OF A CONICAL WEDGE")
	fmt.Println()

	coneDiameter := prompter.Float("Cone base diameter", 2.0)
	coneRadius := 0.5 * coneDiameter
	height := prompter.Float("Height of cone", 10.0)
	coneVolume := math.Pi * coneRadius * coneRadius * height / 3

	var sagitta float64
	for {
		sagitta = prompter.Float("Sagitta of wedge base", 0.5)
		if sagitta > 0 && sagitta <= coneDiameter {
			break
		}
		fmt.Println("\n0 < sagitta <= diameter of cone base is required")
	}

	wedge := wedgeVolume(height, sagitta, coneRadius)

	fmt.Printf("\nVolume of complete cone = %.4f\n", coneVolume)
	fmt.Printf("Volume of conical wedge = %.4f\n", wedge)
	fmt.Printf("Ratio wedge volume/cone volume = %.4f\n", wedge/coneVolume)
}
