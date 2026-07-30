// temp converts a temperature between Centigrade, Fahrenheit,
// Kelvin, Rankine, and Reaumur scales, reading a value with a
// trailing scale letter (e.g. "100f") and reporting the equivalent
// in all five scales.
//
// Converted from TEMP.C, Misc/temp.
package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"mwkgo/internal/promptio"
)

// celsiusToKelvinOffset is the original program's own c0 constant.
const celsiusToKelvinOffset = 273.18

// fahrenheitToRankineOffset is chosen, per the original program's own
// comment, so that 0 Kelvin corresponds to 0 Rankine.
var fahrenheitToRankineOffset = 1.8*celsiusToKelvinOffset - 32

// temperatures holds a single temperature expressed in all five
// scales this program supports.
type temperatures struct {
	Centigrade, Fahrenheit, Kelvin, Rankine, Reaumur float64
}

// convertTemperature returns value (given on the scale identified by
// scale: 'c', 'f', 'k', 'r', or 'e') expressed in every supported
// scale. Any scale byte other than the five recognized ones is
// treated as Fahrenheit, matching the original program's own
// default.
func convertTemperature(value float64, scale byte) temperatures {
	var t temperatures
	switch scale {
	case 'k':
		t.Kelvin = value
		t.Centigrade = t.Kelvin - celsiusToKelvinOffset
		t.Fahrenheit = 1.8*t.Centigrade + 32
		t.Rankine = t.Fahrenheit + fahrenheitToRankineOffset
		t.Reaumur = t.Centigrade * 0.8
	case 'c':
		t.Centigrade = value
		t.Fahrenheit = 1.8*t.Centigrade + 32
		t.Rankine = t.Fahrenheit + fahrenheitToRankineOffset
		t.Kelvin = t.Centigrade + celsiusToKelvinOffset
		t.Reaumur = t.Centigrade * 0.8
	case 'r':
		t.Rankine = value
		t.Fahrenheit = t.Rankine - fahrenheitToRankineOffset
		t.Centigrade = (t.Fahrenheit - 32) / 1.8
		t.Kelvin = t.Centigrade + celsiusToKelvinOffset
		t.Reaumur = t.Centigrade * 0.8
	case 'e':
		t.Reaumur = value
		t.Centigrade = t.Reaumur / 0.8
		t.Fahrenheit = 1.8*t.Centigrade + 32
		t.Rankine = t.Fahrenheit + fahrenheitToRankineOffset
		t.Kelvin = t.Centigrade + celsiusToKelvinOffset
	default:
		t.Fahrenheit = value
		t.Rankine = t.Fahrenheit + fahrenheitToRankineOffset
		t.Centigrade = (t.Fahrenheit - 32) / 1.8
		t.Kelvin = t.Centigrade + celsiusToKelvinOffset
		t.Reaumur = t.Centigrade * 0.8
	}
	return t
}

// parseLeadingFloat parses the leading numeric prefix of s (matching
// C's atof, which reads as much of a valid number as it can and
// ignores the rest, rather than requiring the entire string to be
// numeric), returning 0 if no valid prefix is found.
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

// parseTemperatureInput parses an input line such as "100f", "37.5c",
// or "q" (to quit), matching the original program's own free-form
// input format: a leading number followed by a scale letter anywhere
// in the string (c, k, r, or e; anything else, including no letter
// at all, means Fahrenheit). If more than one scale letter appears,
// the last one found wins, matching the original's sequence of
// independent (not else-if) checks.
func parseTemperatureInput(line string) (value float64, scale byte, quit bool) {
	s := strings.ToLower(strings.TrimSpace(line))
	if strings.Contains(s, "q") {
		return 0, 0, true
	}

	value = parseLeadingFloat(s)
	scale = 'f'
	for _, candidate := range []byte{'c', 'k', 'r', 'e'} {
		if strings.IndexByte(s, candidate) != -1 {
			scale = candidate
		}
	}
	return value, scale, false
}

func main() {
	prompter, err := promptio.New(os.Stdin, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, "temp:", err)
		os.Exit(1)
	}

	fmt.Println("TEMPERATURE CONVERTOR")
	fmt.Println()

	if len(os.Args) > 1 {
		printConversion(os.Args[1])
		return
	}

	fmt.Println("Temperature scales are")
	fmt.Println("  Centigrade")
	fmt.Println("  Fahrenheit")
	fmt.Println("  Kelvin")
	fmt.Println("  Rankine")
	fmt.Println("  rEaumur")

	for {
		line := prompter.Line("Enter temperature as 123.45x, x=c,[f],k,r,e or q(uit)) ")
		if !printConversion(line) {
			return
		}
	}
}

// printConversion parses and prints the conversion for one input
// line, returning false if the line requested quitting.
func printConversion(line string) bool {
	value, scale, quit := parseTemperatureInput(line)
	if quit {
		return false
	}

	if value < 0 && (scale == 'k' || scale == 'r') {
		fmt.Println("Absolute temperatures are always positive!")
		return true
	}

	t := convertTemperature(value, scale)
	fmt.Println()
	fmt.Printf("Centigrade = %10.3f\n", t.Centigrade)
	fmt.Printf("Fahrenheit = %10.3f\n", t.Fahrenheit)
	fmt.Printf("Kelvin     = %10.3f\n", t.Kelvin)
	fmt.Printf("Rankine    = %10.3f\n", t.Rankine)
	fmt.Printf("Reaumur    = %10.3f\n", t.Reaumur)
	fmt.Println()
	return true
}
