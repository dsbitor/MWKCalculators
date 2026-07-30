// cone computes the flat sheet-metal pattern for a conical part (a
// frustum): the two arc radii, chord lengths, and included angle
// needed to lay out and roll a flat pattern that forms the required
// cone shape, with an overlap allowance for joining the seam.
//
// Converted from CONE.C (M. W. Klotz, 10-11/01),
// WorkshopUtilities/cone.
package main

import (
	"fmt"
	"math"
	"os"

	"mwkgo/internal/promptio"
)

// conePattern is the flat pattern needed to roll a conical frustum.
type conePattern struct {
	IncludedAngleDeg  float64
	SmallRadius       float64
	SmallChord        float64
	SmallArcLength    float64
	LargeRadius       float64
	LargeChord        float64
	LargeArcLength    float64
	EdgeLength        float64
	ConeIncludedAngle float64
}

// computeConePattern returns the flat pattern for a conical frustum
// of the given small and large diameters and height, with overlap
// added to the large circumference for joining the seam.
func computeConePattern(smallDiameter, largeDiameter, height, overlap float64) conePattern {
	halfDiameterDifference := 0.5 * (largeDiameter - smallDiameter)
	slopeRatio := halfDiameterDifference / height
	coneHalfAngleDeg := math.Atan(slopeRatio) * 180 / math.Pi
	sinConeHalfAngle := math.Sin(coneHalfAngleDeg * math.Pi / 180)

	var smallRadius, smallChord float64
	if sinConeHalfAngle != 0 {
		smallRadius = 0.5 * smallDiameter / sinConeHalfAngle
		smallArcAngle := (math.Pi*smallDiameter + overlap) / smallRadius
		smallChord = 2 * smallRadius * math.Sin(0.5*smallArcAngle)
	}

	largeRadius := 0.5 * largeDiameter / sinConeHalfAngle
	largeArcAngle := (math.Pi*largeDiameter + overlap) / largeRadius
	largeChord := 2 * largeRadius * math.Sin(0.5*largeArcAngle)
	includedAngleDeg := largeArcAngle * 180 / math.Pi

	return conePattern{
		IncludedAngleDeg:  includedAngleDeg,
		SmallRadius:       smallRadius,
		SmallChord:        smallChord,
		SmallArcLength:    includedAngleDeg * math.Pi / 180 * smallRadius,
		LargeRadius:       largeRadius,
		LargeChord:        largeChord,
		LargeArcLength:    includedAngleDeg * math.Pi / 180 * largeRadius,
		EdgeLength:        largeRadius - smallRadius,
		ConeIncludedAngle: coneHalfAngleDeg,
	}
}

func main() {
	prompter, err := promptio.New(os.Stdin, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, "cone:", err)
		os.Exit(1)
	}

	fmt.Println("FLAT PATTERNS FOR CONICAL PARTS")
	fmt.Println()

	smallDiameter := prompter.Float("small diameter of cone", 3.0)
	largeDiameter := prompter.Float("large diameter of cone", 5.0)
	height := prompter.Float("height of cone", 10.0)
	overlap := prompter.Float("overlap allowance for joining", 0.25)

	pattern := computeConePattern(smallDiameter, largeDiameter, height, overlap)

	fmt.Printf("\nsmall diameter = %.4f\n", smallDiameter)
	fmt.Printf("small circumference = %.4f\n", math.Pi*smallDiameter)
	fmt.Printf("large diameter = %.4f\n", largeDiameter)
	fmt.Printf("large circumference = %.4f\n", math.Pi*largeDiameter)
	fmt.Printf("cone height = %.4f\n", height)
	fmt.Printf("overlap allowance = %.4f\n", overlap)

	fmt.Printf("\nincluded angle of pattern = %.4f deg\n", pattern.IncludedAngleDeg)
	fmt.Printf("smaller radius of pattern = %.4f\n", pattern.SmallRadius)
	fmt.Printf("  chord of smaller radius = %.4f\n", pattern.SmallChord)
	fmt.Printf("  arc length for smaller radius = %.4f\n", pattern.SmallArcLength)
	fmt.Printf("larger radius of pattern = %.4f\n", pattern.LargeRadius)
	fmt.Printf("  chord of larger radius = %.4f\n", pattern.LargeChord)
	fmt.Printf("  arc length for larger radius = %.4f\n", pattern.LargeArcLength)
	fmt.Printf("length of edge = %.4f\n", pattern.EdgeLength)
	fmt.Printf("cone included angle = %.4f deg\n", pattern.ConeIncludedAngle)
}
