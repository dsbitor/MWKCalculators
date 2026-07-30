// gas solves the ideal gas law, PV = nRT, for whichever of pressure,
// volume, moles, or temperature is unknown, given the other three,
// accepting each in a choice of common units.
//
// Converted from GAS.C, WorkshopUtilities/gas.
package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"mwkgo/internal/promptio"
)

const (
	gasConstant        = 8.20545e-2 // liter-atm/(degK-mole)
	celsiusToKelvin    = 273.16
	psiPerAtmosphere   = 14.6963
	kpaPerAtmosphere   = 101.327
	cubicFeetPerLiter  = 0.0353147
	cubicMeterPerLiter = 0.001
)

// gasState holds the four ideal-gas quantities.
type gasState struct {
	PressureAtm, VolumeLiters, Moles, TemperatureK float64
}

// errInsufficientGasData reports that fewer than three of the four
// quantities were given.
var errInsufficientGasData = errors.New("insufficient data for solution")

// solveGasLaw solves PV = nRT for whichever of p (atm), v (liters),
// n (moles), or t (Kelvin) is zero, given the other three.
func solveGasLaw(p, v, n, t float64) (gasState, error) {
	known := 0
	for _, x := range []float64{p, v, n, t} {
		if x != 0 {
			known++
		}
	}
	if known < 3 {
		return gasState{}, errInsufficientGasData
	}

	switch {
	case p == 0:
		p = n * gasConstant * t / v
	case v == 0:
		v = n * gasConstant * t / p
	case n == 0:
		n = (p * v) / (gasConstant * t)
	case t == 0:
		t = (p * v) / (n * gasConstant)
	}
	return gasState{PressureAtm: p, VolumeLiters: v, Moles: n, TemperatureK: t}, nil
}

// parseLeadingFloat parses the leading numeric prefix of s (matching
// C's atof/strtod, which reads as much of a valid number as it can
// and ignores the rest), returning 0 if no valid prefix is found.
func parseLeadingFloat(s string) float64 {
	end := 0
	for end < len(s) {
		c := s[end]
		isValid := (c >= '0' && c <= '9') || c == '.' || c == '+' || c == '-' || c == 'e' || c == 'E'
		if end > 0 && (c == '+' || c == '-') && s[end-1] != 'e' && s[end-1] != 'E' {
			isValid = false
		}
		if !isValid {
			break
		}
		end++
	}
	value, _ := strconv.ParseFloat(s[:end], 64)
	return value
}

// parsePressure parses a pressure reading such as "1", "14.7psi", or
// "100kpascal" into atmospheres. Note that the trigger substring is
// "pas", not "kpa": the original program's own prompt abbreviates
// the unit as "(kpa)scal" but its strstr check looks for "pas", so a
// natural-looking "100kpa" (without the rest of the word) falls
// through to being read as plain atmospheres instead. This is almost
// certainly an inconsistency between the prompt text and the check
// in the original, but it's what the original program actually does,
// so it's preserved here rather than silently corrected to match
// "kpa".
func parsePressure(s string) float64 {
	s = strings.ToLower(strings.TrimSpace(s))
	value := parseLeadingFloat(s)
	switch {
	case strings.Contains(s, "psi"):
		return value / psiPerAtmosphere
	case strings.Contains(s, "pas"):
		return value / kpaPerAtmosphere
	default:
		return value
	}
}

// parseVolume parses a volume reading such as "10", "5cft", or
// "2cm" (cubic meters) into liters.
func parseVolume(s string) float64 {
	s = strings.ToLower(strings.TrimSpace(s))
	value := parseLeadingFloat(s)
	switch {
	case strings.Contains(s, "cft"):
		return value / cubicFeetPerLiter
	case strings.Contains(s, "cm"):
		return value / cubicMeterPerLiter
	default:
		return value
	}
}

// parseTemperature parses a temperature reading such as "293", "20c",
// or "68f" into Kelvin.
func parseTemperature(s string) float64 {
	s = strings.ToLower(strings.TrimSpace(s))
	value := parseLeadingFloat(s)
	switch {
	case strings.Contains(s, "c"):
		return value + celsiusToKelvin
	case strings.Contains(s, "f"):
		return 5*(value-32)/9 + celsiusToKelvin
	default:
		return value
	}
}

func main() {
	prompter, err := promptio.New(os.Stdin, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, "gas:", err)
		os.Exit(1)
	}

	fmt.Println("PERFECT GAS LAW (PV = nRT) CALCULATOR")
	fmt.Println()
	fmt.Println("P = pressure; V = volume; T = temperature")
	fmt.Println("n = number of moles of (assumed perfect) gas")
	fmt.Println("R = gas constant = 8.20545E-2 liter-atm/(degK-mole)")
	fmt.Println("Avogadro's number = 6.02308E23 = number of molecules in a mole of gas")
	fmt.Println("mole of gas = weight (gm) equal to molecular weight (e.g. 32 for O2)")
	fmt.Println("A mole of ordinary air weighs ~29 gm.")
	fmt.Println("standard conditions = 0 degC (273.16 degK) and 1 atmosphere.")
	fmt.Println("One mole of gas occupies 22.4140 liters at standard conditions.")

	fmt.Println("\nInput whatever data you know; press return if not known.")
	fmt.Println("You must enter three data items to obtain a solution.")
	fmt.Println()

	p := parsePressure(prompter.Line("Pressure [atm]osphere, (psi), (kpa)scal ? "))
	v := parseVolume(prompter.Line("Volume [l]iter, (cft), (cm)eter ? "))
	n := parseLeadingFloat(prompter.Line("Moles ? "))
	known := 0
	for _, x := range []float64{p, v, n} {
		if x != 0 {
			known++
		}
	}

	var t float64
	if known < 3 {
		t = parseTemperature(prompter.Line("Temperature [K], (C), (F) ? "))
		if t != 0 {
			known++
		}
	}
	if known < 3 {
		fmt.Println("\nINSUFFICIENT DATA FOR SOLUTION")
		return
	}

	state, err := solveGasLaw(p, v, n, t)
	if err != nil {
		fmt.Println("\nINSUFFICIENT DATA FOR SOLUTION")
		return
	}

	tc := state.TemperatureK - celsiusToKelvin
	tf := 9*tc/5 + 32
	fmt.Println()
	fmt.Printf("P = %12.4f atm = %12.4f psi  = %12.4f kpas\n", state.PressureAtm, state.PressureAtm*psiPerAtmosphere, state.PressureAtm*kpaPerAtmosphere*0.001)
	fmt.Printf("V = %12.4f l   = %12.4f ft^3 = %12.4f m^3 \n", state.VolumeLiters, state.VolumeLiters*cubicFeetPerLiter, state.VolumeLiters*cubicMeterPerLiter)
	fmt.Printf("T = %12.4f K   = %12.4f F   = %12.4f C  \n", state.TemperatureK, tf, tc)
	fmt.Printf("n = %12.4f mole\n", state.Moles)
}
