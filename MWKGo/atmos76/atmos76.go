// atmos76 computes atmospheric temperature, pressure, and density
// ratios (relative to sea level) at a given altitude, per the 1976
// U.S. Standard Atmosphere model, valid from 0 to 86 km.
//
// Converted from ATMOS76.C, Misc/atmos76.
package main

import (
	"fmt"
	"math"
	"os"

	"mwkgo/internal/promptio"
)

const (
	feetPerKm       = 3280.839895
	earthRadiusKm   = 6369.0
	gasConstant     = 34.163195
	seaLevelTempK   = 288.15
	seaLevelPressPa = 101325.0
	seaLevelDensity = 1.225
)

// Standard Atmosphere layer base tables: geopotential height (km),
// base temperature (K), base pressure ratio (relative to sea
// level), and temperature lapse rate (K/km), one entry per layer.
var (
	layerBaseHeightKm   = [...]float64{0.0, 11.0, 20.0, 32.0, 47.0, 51.0, 71.0, 84.852}
	layerBaseTempK      = [...]float64{288.15, 216.65, 216.65, 228.65, 270.65, 270.65, 214.65, 186.946}
	layerBasePressRatio = [...]float64{1.0, 2.233611e-1, 5.403295e-2, 8.5666784e-3, 1.0945601e-3, 6.6063531e-4, 3.9046834e-5, 3.68501e-6}
	layerLapseRate      = [...]float64{-6.5, 0.0, 1.0, 2.8, 0.0, -2.8, -2.0, 0.0}
)

// atmosphereRatios holds the temperature, pressure, and density
// ratios (each relative to sea level) at a given altitude.
type atmosphereRatios struct {
	TemperatureRatio float64
	PressureRatio    float64
	DensityRatio     float64
}

// computeAtmosphereRatios returns the standard-atmosphere ratios at
// geometric altitude altitudeKm (kilometers above sea level).
func computeAtmosphereRatios(altitudeKm float64) atmosphereRatios {
	// Convert geometric to geopotential altitude.
	h := altitudeKm * earthRadiusKm / (altitudeKm + earthRadiusKm)

	layer := len(layerBaseHeightKm) - 2
	for i := 0; i < len(layerBaseHeightKm)-1; i++ {
		if h >= layerBaseHeightKm[i] && h < layerBaseHeightKm[i+1] {
			layer = i
			break
		}
	}

	lapseRate := layerLapseRate[layer]
	baseTemp := layerBaseTempK[layer]
	deltaH := h - layerBaseHeightKm[layer]
	localTemp := baseTemp + lapseRate*deltaH
	temperatureRatio := localTemp / layerBaseTempK[0]

	var pressureRatio float64
	if lapseRate == 0 {
		pressureRatio = layerBasePressRatio[layer] * math.Exp(-gasConstant*deltaH/baseTemp)
	} else {
		pressureRatio = layerBasePressRatio[layer] * math.Pow(baseTemp/localTemp, gasConstant/lapseRate)
	}

	return atmosphereRatios{
		TemperatureRatio: temperatureRatio,
		PressureRatio:    pressureRatio,
		DensityRatio:     pressureRatio / temperatureRatio,
	}
}

func main() {
	prompter, err := promptio.New(os.Stdin, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, "atmos76:", err)
		os.Exit(1)
	}

	fmt.Println("ATMOSPHERIC COMPUTATIONS (valid 0 to 86 km)")
	fmt.Println()
	fmt.Println("1 km = 3280.84 ft")

	altitudeFt := prompter.Float("Altitude (ft)", 1000.0)
	altitudeKm := altitudeFt / feetPerKm

	ratios := computeAtmosphereRatios(altitudeKm)

	fmt.Printf("\nAltitude = %.4E km = %.4E mi = %.4E ft\n", altitudeKm, altitudeFt/5280, altitudeFt)
	fmt.Printf("Temperature ratio (wrt sl) = %.4E\n", ratios.TemperatureRatio)
	fmt.Printf("Pressure ratio (wrt sl) = %.4E\n", ratios.PressureRatio)
	fmt.Printf("Density ratio (wrt sl) = %.4E\n", ratios.DensityRatio)

	fmt.Println("\nFor sea level values of temperature, pressure, density of:")
	fmt.Printf("%.4E K, %.4E N/m^2, %.4E kg/m^3\n\n", seaLevelTempK, seaLevelPressPa, seaLevelDensity)

	ta := ratios.TemperatureRatio * seaLevelTempK
	pa := ratios.PressureRatio * seaLevelPressPa
	da := ratios.DensityRatio * seaLevelDensity

	tc := ta - 273.18
	tf := 1.8*tc + 32
	fmt.Printf("Temperature = %.4f K = %.4f C = %.4f F \n", ta, tc, tf)
	fmt.Printf("Pressure = %.4E N/m^2 = %.4E psi\n", pa, pa/6894.76)
	fmt.Printf("Density = %.4E kg/m^3 = %.4E slug/ft^3\n", da, da/515.379)

	speedOfSoundMph := 1116.45 * math.Sqrt((tf+459.67)/518.67) * 3600 / 5280
	fmt.Printf("Speed of sound = %.4f kph = %.4f mph\n", speedOfSoundMph*1.60934, speedOfSoundMph)
}
