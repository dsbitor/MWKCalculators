// gearspur solves for standard spur gear dimensions given any two of
// four related quantities (outside diameter, tooth count, diametral
// pitch, or pitch diameter; a module value is accepted as an
// alternate way to specify diametral pitch), then reports the full
// set of standard dimensions plus the Brown & Sharpe cutter number
// needed to cut the gear.
//
// Converted from GEARSPUR.C, WorkshopUtilities/gearspur.
package main

import (
	"errors"
	"fmt"
	"math"
	"os"
	"strings"

	"mwkgo/internal/promptio"
)

// spurGear holds the standard dimensions of a single spur gear.
type spurGear struct {
	Teeth                                                         float64
	DiametralPitch                                                float64
	Module                                                        float64
	OutsideDiameter, PitchDiameter                                float64
	Addendum, Dedendum, WholeDepth, CircularPitch, ToothThickness float64
}

// errInsufficientGearData reports that fewer than two of the four
// solvable quantities (outside diameter, teeth, diametral pitch,
// pitch diameter) were given.
var errInsufficientGearData = errors.New("insufficient data for solution")

// solveSpurGear solves for a spur gear's teeth, diametral pitch,
// outside diameter, and pitch diameter given any two of those four
// (zero meaning unknown), then returns the full set of standard
// dimensions. It recognizes the same six combinations as the
// original program.
func solveSpurGear(outsideDiameter, teeth, diametralPitch, pitchDiameter float64) (spurGear, error) {
	switch {
	case teeth != 0 && outsideDiameter != 0:
		diametralPitch = (teeth + 2) / outsideDiameter
		pitchDiameter = teeth / diametralPitch
	case diametralPitch != 0 && teeth != 0:
		outsideDiameter = (teeth + 2) / diametralPitch
		pitchDiameter = teeth / diametralPitch
	case diametralPitch != 0 && outsideDiameter != 0:
		teeth = diametralPitch*outsideDiameter - 2
		pitchDiameter = teeth / diametralPitch
	case pitchDiameter != 0 && teeth != 0:
		diametralPitch = teeth / pitchDiameter
		outsideDiameter = (teeth + 2) / diametralPitch
	case pitchDiameter != 0 && outsideDiameter != 0:
		teeth = 2 * pitchDiameter / (outsideDiameter - pitchDiameter)
		diametralPitch = teeth / pitchDiameter
	case pitchDiameter != 0 && diametralPitch != 0:
		teeth = diametralPitch * pitchDiameter
		outsideDiameter = (teeth + 2) / diametralPitch
	default:
		return spurGear{}, errInsufficientGearData
	}

	addendum := 1.0 / diametralPitch
	wholeDepth := 2.2/diametralPitch + 0.002
	if diametralPitch > 20 {
		wholeDepth = 2.157 / diametralPitch
	}
	circularPitch := math.Pi / diametralPitch

	return spurGear{
		Teeth:           teeth,
		DiametralPitch:  diametralPitch,
		Module:          mmPerInch / diametralPitch,
		OutsideDiameter: outsideDiameter,
		PitchDiameter:   pitchDiameter,
		Addendum:        addendum,
		Dedendum:        wholeDepth - addendum,
		WholeDepth:      wholeDepth,
		CircularPitch:   circularPitch,
		ToothThickness:  0.48 * circularPitch,
	}, nil
}

const mmPerInch = 25.4

// brownAndSharpeCutter returns the Brown & Sharpe involute gear
// cutter number needed to cut a gear with teeth teeth.
func brownAndSharpeCutter(teeth float64) int {
	switch {
	case teeth >= 135:
		return 1
	case teeth >= 55:
		return 2
	case teeth >= 35:
		return 3
	case teeth >= 26:
		return 4
	case teeth >= 21:
		return 5
	case teeth >= 17:
		return 6
	case teeth >= 14:
		return 7
	case teeth >= 12:
		return 8
	default:
		return 0
	}
}

func main() {
	prompter, err := promptio.New(os.Stdin, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, "gearspur:", err)
		os.Exit(1)
	}

	metric := strings.ToLower(prompter.Line("[I]mperial or (M)etric units? ")) == "m"
	unit := "in"
	if metric {
		unit = "mm"
	}

	for {
		fmt.Println("\nEnter whatever data you know.  Enter zero (0) for unknowns.")
		fmt.Println("You must enter two data items to obtain an answer.")
		fmt.Println()

		outsideDiameter := prompter.Float("OD of gear ("+unit+")", 2.35)
		if metric {
			outsideDiameter /= mmPerInch
		}
		teeth := prompter.Float("Number of teeth", 45.0)

		known := 0
		for _, v := range []float64{outsideDiameter, teeth} {
			if v != 0 {
				known++
			}
		}

		var diametralPitch, pitchDiameter, module float64
		if known < 2 {
			diametralPitch = prompter.Float("Diametral Pitch", 20.0)
			if diametralPitch != 0 {
				known++
			}
		}
		if known < 2 {
			module = prompter.Float("Module", 0.7874)
			if module != 0 {
				diametralPitch = mmPerInch / module
				known++
			}
		}
		if known < 2 {
			pitchDiameter = prompter.Float("Pitch Diameter ("+unit+")", 2.25)
			if metric {
				pitchDiameter /= mmPerInch
			}
			if pitchDiameter != 0 {
				known++
			}
		}
		if known < 2 {
			fmt.Println("\nINSUFFICENT DATA FOR SOLUTION. TRY AGAIN.")
			continue
		}

		gear, err := solveSpurGear(outsideDiameter, teeth, diametralPitch, pitchDiameter)
		if err != nil {
			fmt.Println("\nINSUFFICENT DATA FOR SOLUTION. TRY AGAIN.")
			continue
		}

		fmt.Printf("\nDiametral Pitch = %.4f\n", gear.DiametralPitch)
		fmt.Printf("Module = %.4f\n", gear.Module)
		fmt.Printf("Number of teeth = %.0f\n", gear.Teeth)
		fmt.Printf("Outside Diameter = %.4f in = %.4f mm\n", gear.OutsideDiameter, gear.OutsideDiameter*mmPerInch)
		fmt.Printf("Pitch Diameter = %.4f in = %.4f mm\n", gear.PitchDiameter, gear.PitchDiameter*mmPerInch)
		fmt.Printf("Addendum = %.4f in = %.4f mm\n", gear.Addendum, gear.Addendum*mmPerInch)
		fmt.Printf("Dedendum = %.4f in = %.4f mm\n", gear.Dedendum, gear.Dedendum*mmPerInch)
		fmt.Printf("Whole Depth = %.4f in = %.4f mm\n", gear.WholeDepth, gear.WholeDepth*mmPerInch)
		fmt.Printf("Circular Pitch = %.4f in = %.4f mm\n", gear.CircularPitch, gear.CircularPitch*mmPerInch)
		fmt.Printf("Tooth Thickness = %.4f in = %.4f mm\n", gear.ToothThickness, gear.ToothThickness*mmPerInch)
		fmt.Printf("\nB & S cutter number used to cut this gear = %d\n", brownAndSharpeCutter(gear.Teeth))
		return
	}
}
