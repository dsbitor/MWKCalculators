// chill computes the equivalent wind chill temperature using the US
// National Weather Service's 2001 formula.
//
// Converted from CHILL.C (M. W. Klotz, 11/98, 12/01),
// WorkshopUtilities/chill. The original source also carries, as a
// comment, the older 1945-style wind chill lookup table this formula
// replaced; that table corresponds to a different, superseded
// formula and is not reproduced here since it was never the code
// path this program executed.
package main

import (
	"fmt"
	"math"
	"os"

	"mwkgo/internal/promptio"
)

// windChillF returns the equivalent wind chill temperature in
// degrees Fahrenheit, using the US National Weather Service's 2001
// formula. The formula (and this conversion) is only meaningful for
// temperatures at or below 50F and wind speeds at or above 3mph, the
// NWS's stated range of validity; it is not guarded against inputs
// outside that range, matching the original program.
func windChillF(tempF, windMph float64) float64 {
	windFactor := math.Pow(windMph, 0.16)
	return 35.74 + 0.6215*tempF - windFactor*(35.75-0.4275*tempF)
}

func main() {
	prompter, err := promptio.New(os.Stdin, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, "chill:", err)
		os.Exit(1)
	}

	tempF := prompter.Float("Temperature (degF)", 30.0)
	windMph := prompter.Float("Wind speed (mph)", 25.0)

	fmt.Printf("\nEquivalent wind chill temperature = %.1f degF\n", windChillF(tempF, windMph))
}
