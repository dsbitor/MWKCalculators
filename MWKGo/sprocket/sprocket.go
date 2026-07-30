// sprocket computes ANSI standard roller chain sprocket dimensions:
// pitch diameter, outside diameter, caliper diameter (and, for an
// odd tooth count, the caliper factor it's derived from), and the
// maximum recommended hub diameter, given the tooth count, chain
// pitch, and roller diameter.
//
// Converted from SPROCKET.C (M. W. Klotz), WorkshopUtilities/sprocket
// (the same archive as chain, converted in Tier 1 group 4).
package main

import (
	"fmt"
	"math"
	"os"

	"mwkgo/internal/promptio"
)

// sprocketDimensions holds the computed dimensions of a roller
// chain sprocket. CaliperFactor is only meaningful (and only
// printed) when the tooth count is odd.
type sprocketDimensions struct {
	PitchDiameter   float64
	OutsideDiameter float64
	CaliperDiameter float64
	CaliperFactor   float64
	OddTeeth        bool
	MaxHubDiameter  float64
}

// computeSprocketDimensions returns the dimensions of an ANSI
// standard roller chain sprocket with the given tooth count, chain
// pitch, and roller diameter.
//
// The caliper diameter for an odd tooth count is derived from a
// caliper factor of pitchDiameter*cos(90/teeth). The companion
// SPROCKET.TXT describes this as cos(180/teeth), which appears to be
// a transcription slip from Machinery's Handbook's actual formula:
// the source code's cos(90/teeth) is what the original program
// computes and prints, and is the standard formula for measuring an
// odd-tooth sprocket's caliper diameter, so it is what this
// conversion implements.
func computeSprocketDimensions(teeth int, pitch, rollerDiameter float64) sprocketDimensions {
	halfToothAngle := 180.0 / float64(teeth)
	pitchDiameter := pitch / math.Sin(halfToothAngle*math.Pi/180)
	outsideDiameter := pitch * (0.6 + 1/math.Tan(halfToothAngle*math.Pi/180))
	maxHubDiameter := pitch*(1/math.Tan(halfToothAngle*math.Pi/180)-1) - 0.030

	oddTeeth := teeth%2 != 0
	var caliperDiameter, caliperFactor float64
	if oddTeeth {
		caliperFactor = pitchDiameter * math.Cos(90.0/float64(teeth)*math.Pi/180)
		caliperDiameter = caliperFactor - rollerDiameter
	} else {
		caliperDiameter = pitchDiameter - rollerDiameter
	}

	return sprocketDimensions{
		PitchDiameter:   pitchDiameter,
		OutsideDiameter: outsideDiameter,
		CaliperDiameter: caliperDiameter,
		CaliperFactor:   caliperFactor,
		OddTeeth:        oddTeeth,
		MaxHubDiameter:  maxHubDiameter,
	}
}

func main() {
	prompter, err := promptio.New(os.Stdin, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, "sprocket:", err)
		os.Exit(1)
	}

	fmt.Println("ANSI ROLLER CHAIN SPROCKET CALCULATIONS")
	fmt.Println()

	teeth := prompter.Int("Number of teeth in sprocket", 9)
	pitch := prompter.Float("Chain pitch (in)", 1.0)
	rollerDiameter := prompter.Float("Roller diameter (in)", 0.25)

	d := computeSprocketDimensions(teeth, pitch, rollerDiameter)

	fmt.Println()
	fmt.Printf("Sprocket pitch diameter = %.4f in\n", d.PitchDiameter)
	fmt.Printf("Sprocket outside diameter = %.4f in\n", d.OutsideDiameter)
	if d.OddTeeth {
		fmt.Printf("Caliper factor = %.4f\n", d.CaliperFactor)
	}
	fmt.Printf("Sprocket caliper diameter = %.4f in\n", d.CaliperDiameter)
	fmt.Printf("Maximum hub diameter = %.4f in\n", d.MaxHubDiameter)
}
