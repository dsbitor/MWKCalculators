// dew computes dew point temperature from ambient temperature and
// relative humidity, using the Magnus-Tetens approximation. Metal
// tools left below the dew point will accumulate condensation and
// rust; knowing the dew point tells a shop when tools need warming
// to stay above it.
//
// Converted from DEW.C (M. W. Klotz, 3/04), Misc/dew. Valid for
// ambient temperatures between 0 and 60 degC (32 to 140 degF) and
// relative humidity between 1 and 100 percent, per the original
// program's own stated range.
package main

import (
	"fmt"
	"math"
	"os"

	"mwkgo/internal/promptio"
)

// magnusA and magnusB are the August-Roche-Magnus approximation's
// standard constants for temperatures in degrees Celsius.
const (
	magnusA = 17.27
	magnusB = 237.7 // degC
)

// dewPointCelsius returns the dew point temperature in degrees
// Celsius, given the ambient temperature in degrees Celsius and the
// relative humidity as a fraction (0.5 for 50%, not 50).
func dewPointCelsius(tempC, relativeHumidity float64) float64 {
	alpha := magnusA*tempC/(magnusB+tempC) + math.Log(relativeHumidity)
	return magnusB * alpha / (magnusA - alpha)
}

const statedAccuracyC = 0.04

func main() {
	prompter, err := promptio.New(os.Stdin, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, "dew:", err)
		os.Exit(1)
	}

	fmt.Println()
	fmt.Println("DEWPOINT CALCULATION")
	fmt.Println()

	fmt.Println("Preferred temperature scale - Fahrenheit (default) or Celsius")
	scale := prompter.Line("[F], C ? ")
	useCelsius := scale == "c" || scale == "C"

	var tempC float64
	if useCelsius {
		tempC = prompter.Float("Measured temperature (degC)", 20.0)
	} else {
		tempF := prompter.Float("Measured temperature (degF)", 70.0)
		tempC = (tempF - 32) / 1.8
	}

	humidityPercent := prompter.Float("Relative humidity (%)", 60.0)
	dewPointC := dewPointCelsius(tempC, humidityPercent*0.01)

	fmt.Printf("Dew point temperature = %.2f degC = %.2f degF\n", dewPointC, 1.8*dewPointC+32)
	fmt.Printf("Accuracy = +/- %.2f degC = +/- %.2f degF\n", statedAccuracyC, 1.8*statedAccuracyC)
}
