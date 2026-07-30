// psych computes relative humidity and dew point from a wet/dry
// bulb psychrometer reading: two identical thermometers, one with
// its bulb kept wet by an evaporating wick, read a lower temperature
// on the wet bulb due to evaporative cooling, and the size of that
// difference indicates how humid the air is.
//
// Converted from PSYCH.C, Misc/psych. Reference:
// http://www.uswcl.ars.ag.gov/exper/relhumeq.htm
package main

import (
	"fmt"
	"math"
	"os"
	"strings"

	"mwkgo/internal/promptio"
)

// zeroCelsiusOffset is the constant in the saturation vapor pressure
// approximation (Magnus-Tetens form) used throughout this program.
const zeroCelsiusOffset = 237.3

// fahrenheitToCelsius and celsiusToFahrenheit convert between the
// two temperature scales, matching the original program's dc and df
// functions.
func fahrenheitToCelsius(f float64) float64 { return (f - 32) / 1.8 }
func celsiusToFahrenheit(c float64) float64 { return 1.8*c + 32 }

// saturationVaporPressure returns the saturation vapor pressure, in
// kPa, at temperature tempC (Celsius), per the Magnus-Tetens
// approximation used by the original program.
func saturationVaporPressure(tempC float64) float64 {
	return math.Exp((16.78*tempC - 116.9) / (tempC + zeroCelsiusOffset))
}

// psychrometerResult holds the humidity and dew point derived from a
// wet/dry bulb reading.
type psychrometerResult struct {
	RelativeHumidityPct float64
	DewPointC           float64
	DewPointF           float64
}

// computePsychrometer returns the relative humidity and dew point
// for a dry bulb temperature dryBulbC and wet bulb temperature
// wetBulbC (both Celsius), at elevation elevationM meters above sea
// level.
func computePsychrometer(elevationM, dryBulbC, wetBulbC float64) psychrometerResult {
	airPressure := 101.325 * math.Exp(-0.0001184*elevationM)
	conversionFactor := 0.00066 * (1 + 0.00115*wetBulbC)

	satDryBulb := saturationVaporPressure(dryBulbC)
	satWetBulb := saturationVaporPressure(wetBulbC)

	vaporPressure := satWetBulb - conversionFactor*airPressure*(dryBulbC-wetBulbC)
	relativeHumidity := 100 * vaporPressure / satDryBulb

	logVaporPressure := math.Log(vaporPressure)
	dewPointC := (116.9 + zeroCelsiusOffset*logVaporPressure) / (16.78 - logVaporPressure)

	return psychrometerResult{
		RelativeHumidityPct: relativeHumidity,
		DewPointC:           dewPointC,
		DewPointF:           celsiusToFahrenheit(dewPointC),
	}
}

func main() {
	prompter, err := promptio.New(os.Stdin, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, "psych:", err)
		os.Exit(1)
	}

	fmt.Println("RELATIVE HUMIDITY FROM WET/DRY BULB PSYCHROMETER")
	fmt.Println()

	elevationFt := prompter.Float("Elevation above sea level (ft)", 970.0)
	elevationM := elevationFt * 0.3048

	fmt.Println("Preferred temperature scale - Fahrenheit (default) or Celsius")
	fahrenheit := strings.ToLower(prompter.Line("[F], C ? ")) != "c"

	var dryBulbDefault, wetBulbDefault float64
	unit := "degC"
	if fahrenheit {
		dryBulbDefault, wetBulbDefault, unit = 70.0, 60.0, "degF"
	} else {
		dryBulbDefault, wetBulbDefault = 20.0, 14.0
	}

	dryBulb := prompter.Float(fmt.Sprintf("Dry bulb temperature (%s)", unit), dryBulbDefault)
	if fahrenheit {
		dryBulb = fahrenheitToCelsius(dryBulb)
	}
	wetBulb := prompter.Float(fmt.Sprintf("Wet bulb temperature (%s)", unit), wetBulbDefault)
	if fahrenheit {
		wetBulb = fahrenheitToCelsius(wetBulb)
	}

	result := computePsychrometer(elevationM, dryBulb, wetBulb)

	fmt.Printf("Relative humidity = %.1f %%\n", result.RelativeHumidityPct)
	fmt.Printf("Dewpoint temperature = %.2f degC = %.2f degF\n", result.DewPointC, result.DewPointF)
}
