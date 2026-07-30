// 3wire implements the three-wire method for measuring the pitch
// diameter of a 60 degree screw thread: three precision wires of a
// known diameter are laid in the thread grooves, and the distance
// measured over the outer two wires (or, in reverse, the pitch
// diameter implied by a measured distance) is related to the thread's
// pitch diameter by a well known formula that includes a small
// correction for the thread's lead angle.
//
// Converted from 3WIRE.C, WorkshopUtilities/3wire.
package main

import (
	"fmt"
	"math"
	"os"
	"strings"

	"mwkgo/internal/promptio"
)

const mmPerInch = 25.4

// halfThreadAngleDeg is half of the standard 60 degree thread angle
// this program is restricted to, matching the original program's own
// fixed choice.
const halfThreadAngleDeg = 30.0

// bestWireDiameter returns the "best size" wire diameter for
// measuring a thread of the given pitch: the size at which the wire
// touches the thread flanks exactly at the pitch diameter.
func bestWireDiameter(pitch float64) float64 {
	return 0.5 * pitch / math.Cos(halfThreadAngleDeg*math.Pi/180)
}

// pitchDiameterFromMajor returns the approximate pitch diameter for
// a 60 degree thread given its major diameter and pitch, assuming a
// p/8 flat on the crest (the American National thread form
// assumption the original program uses).
func pitchDiameterFromMajor(majorDiameter, pitch float64) float64 {
	flatCorrection := (3.0 / 16.0) * pitch / math.Tan(halfThreadAngleDeg*math.Pi/180)
	return majorDiameter - 2*flatCorrection
}

// wireOffset returns the offset (the original program's own x) added
// to the pitch diameter to get the measurement over wires, for a
// thread of the given pitch, number of starts, wire diameter, and a
// pitch diameter used only for the lead angle correction (this is
// the actual pitch diameter when it's known, or an estimate from
// pitchDiameterFromMajor otherwise; the correction is small enough
// that an estimate is adequate).
func wireOffset(pitch, starts, pitchDiameterForLeadAngle, wireDiameter float64) float64 {
	halfAngleRad := halfThreadAngleDeg * math.Pi / 180
	cotHalfAngle := 1 / math.Tan(halfAngleRad)
	axialThreadWidth := 0.5 * pitch
	leadAngleTan := pitch * starts / (math.Pi * pitchDiameterForLeadAngle)

	return -axialThreadWidth*cotHalfAngle + wireDiameter*(1+1/math.Sin(halfAngleRad)+0.5*leadAngleTan*leadAngleTan*math.Cos(halfAngleRad)*cotHalfAngle)
}

func main() {
	prompter, err := promptio.New(os.Stdin, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, "3wire:", err)
		os.Exit(1)
	}

	fmt.Println("3 WIRE THREAD MEASUREMENTS (60 degree threads)")
	fmt.Println()

	metric := strings.ToLower(prompter.Line("(M)etric or [I]mperial threads? ")) == "m"

	var pitch float64
	if metric {
		pitchMM := prompter.Float("Thread pitch (mm)", 1.0)
		pitch = pitchMM / mmPerInch
	} else {
		threadsPerInch := prompter.Float("Thread pitch (tpi)", 20.0)
		pitch = 1 / threadsPerInch
	}

	fmt.Println("\nMost common threads are single start.")
	starts := prompter.Float("Number of starts", 1.0)

	suggestedWire := bestWireDiameter(pitch)
	fmt.Printf("\nBest wire diameter to use = %.4f mm = %.4f in\n", suggestedWire*mmPerInch, suggestedWire)

	var wireDiameter float64
	if metric {
		wireDiameterMM := prompter.Float("Wire diameter used (mm)", suggestedWire*mmPerInch)
		wireDiameter = wireDiameterMM / mmPerInch
	} else {
		wireDiameter = prompter.Float("Wire diameter used (in)", suggestedWire)
	}

	var majorDiameter float64
	if metric {
		majorDiameterMM := prompter.Float("Major diameter of thread (mm)", 6.0)
		majorDiameter = majorDiameterMM / mmPerInch
	} else {
		majorDiameter = prompter.Float("Major diameter of thread (in)", 0.25)
	}

	if !metric {
		flatCorrection := (3.0 / 16.0) * pitch / math.Tan(halfThreadAngleDeg*math.Pi/180)
		fmt.Printf("\nFor American National Thread form, subtract %.4f in from\n", 2*flatCorrection)
		fmt.Println("major diameter (assumes p/8 flat on crest) to obtain pitch diameter.")
	}

	estimatedPitchDiameter := pitchDiameterFromMajor(majorDiameter, pitch)

	fmt.Println("\npd = pitch diameter, mow = measurement over wires")
	solveForPD := strings.ToLower(prompter.Line("Calculate pd from mow (p) or mow from pd [m]? ")) == "p"

	if !solveForPD {
		var pitchDiameter float64
		if metric {
			pdMM := prompter.Float("Pitch diameter of thread (mm)", estimatedPitchDiameter*mmPerInch)
			pitchDiameter = pdMM / mmPerInch
		} else {
			pitchDiameter = prompter.Float("Pitch diameter of thread (in)", estimatedPitchDiameter)
		}
		x := wireOffset(pitch, starts, pitchDiameter, wireDiameter)
		mow := pitchDiameter + x
		fmt.Printf("\nMeasurement over wires = %.4f mm = %.4f in\n", mow*mmPerInch, mow)
	} else {
		x := wireOffset(pitch, starts, estimatedPitchDiameter, wireDiameter)
		var mow float64
		if metric {
			mowMM := prompter.Float("Measurement over wires (mm)", (estimatedPitchDiameter+x)*mmPerInch)
			mow = mowMM / mmPerInch
		} else {
			mow = prompter.Float("Measurement over wires (in)", estimatedPitchDiameter+x)
		}
		pitchDiameter := mow - x
		fmt.Printf("\nPitch diameter = %.4f mm = %.4f in\n", pitchDiameter*mmPerInch, pitchDiameter)
	}
}
