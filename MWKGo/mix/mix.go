// mix is a four-function calculator that tracks a single dimensioned
// quantity in six unit systems at once (meter, centimeter,
// millimeter, yard, foot, inch), plus a non-dimensional mode for
// plain arithmetic. Every entry updates all six accumulators
// together, so the running result is always available in whichever
// unit the reader wants to look at.
//
// Converted from MIX.C (M. W. Klotz), Math/mix.
package main

import (
	"errors"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"

	"mwkgo/internal/numdisplay"
	"mwkgo/internal/promptio"
)

// Unit indices. unitND (non-dimensional) has no accumulator of its
// own; the other six each have one.
const (
	unitND = iota
	unitM
	unitCM
	unitMM
	unitYD
	unitFT
	unitIN
	numUnits
)

const numAccumulators = numUnits - 1 // one accumulator per dimensioned unit

var unitNames = [numUnits]string{"nd", "m", "cm", "mm", "yd", "ft", "in"}

// unitPerInch[i] is how many of unit i equal one inch. Index unitND
// is unused (non-dimensional values have no length conversion).
var unitPerInch = [numUnits]float64{1, 0.0254, 2.54, 25.4, 1.0 / 36, 1.0 / 12, 1}

// Arithmetic operators, in the order the original program's operator
// string listed them.
const (
	opAdd = iota
	opSub
	opMul
	opDiv
)

// Display notations.
var (
	errDivideByZero = errors.New("division by zero")
	errUnknownUnit  = errors.New("unknown unit")
	errBadEntry     = errors.New("unrecognized entry")
	errMixedNonDim  = errors.New("cannot add or subtract a non-dimensional value to a dimensioned accumulator")
)

// calcState holds the six running accumulators and the current
// display and input settings, matching the original program's global
// variables.
type calcState struct {
	acc          [numAccumulators]float64
	exponent     int
	undoAcc      [numAccumulators]float64
	undoExponent int
	defaultUnit  int
	decimals     int
	selfAdjust   bool
	notation     int
	fracDenom    float64
	scale        float64
}

func newCalcState() *calcState {
	return &calcState{defaultUnit: unitIN, decimals: 3, fracDenom: 64, scale: 1}
}

func (s *calcState) clear() {
	s.acc = [numAccumulators]float64{}
	s.exponent = 0
}

func (s *calcState) undo() {
	s.acc = s.undoAcc
	s.exponent = s.undoExponent
}

// applyEntry applies one data entry, in the given operator and unit,
// to all six accumulators. value is the raw number as entered,
// before unit conversion or scaling.
func (s *calcState) applyEntry(op, unitIndex int, value float64) error {
	if value == 0 && op == opDiv {
		return errDivideByZero
	}
	if s.exponent != 0 && unitIndex == unitND && op <= opSub {
		return errMixedNonDim
	}

	x := value / unitPerInch[unitIndex]
	if unitIndex != unitND {
		x *= s.scale
	}

	s.undoAcc = s.acc
	s.undoExponent = s.exponent

	for i := 0; i < numAccumulators; i++ {
		z := 1.0
		if unitIndex != unitND {
			z = unitPerInch[i+1]
		}
		switch op {
		case opAdd:
			s.acc[i] += z * x
		case opSub:
			s.acc[i] -= z * x
		case opMul:
			s.acc[i] *= z * x
		case opDiv:
			s.acc[i] /= z * x
		}
	}

	switch {
	case (op == opAdd || op == opSub) && s.exponent == 0 && unitIndex != unitND:
		s.exponent = 1
	case op == opMul && unitIndex != unitND:
		s.exponent++
	case op == opDiv && unitIndex != unitND:
		s.exponent--
	}
	return nil
}

// gcdFloat returns the greatest common divisor of x and y, both
// assumed non-negative and not both zero, via the Euclidean
// algorithm carried out in floating point (needed here because the
// fraction-reduction caller works with float64 counts, not ints).
func gcdFloat(x, y float64) float64 {
	if x == 0 && y == 0 {
		return 1
	}
	a, b := x, y
	if a > b {
		a, b = b, a
	}
	for {
		c := b - a*math.Floor(b/a)
		b = a
		a = c
		if c == 0 {
			return b
		}
	}
}

// parseValue parses a number entered as a plain decimal ("3.75"), a
// simple fraction ("3/4"), or a mixed fraction ("1&1/2").
func parseValue(s string) (float64, error) {
	slash := strings.IndexByte(s, '/')
	if slash == -1 {
		return strconv.ParseFloat(s, 64)
	}

	whole := 0.0
	numDenom := s
	if amp := strings.IndexByte(s, '&'); amp != -1 {
		w, err := strconv.ParseFloat(s[:amp], 64)
		if err != nil {
			return 0, err
		}
		whole = w
		numDenom = s[amp+1:]
	}

	parts := strings.SplitN(numDenom, "/", 2)
	if len(parts) != 2 {
		return 0, errBadEntry
	}
	num, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return 0, err
	}
	denom, err := strconv.ParseFloat(parts[1], 64)
	if err != nil {
		return 0, err
	}
	if denom == 0 {
		return 0, errDivideByZero
	}
	return whole + num/denom, nil
}

// parseDataEntry splits a data-entry line into a leading operator
// (default +), a numeric portion, and a trailing unit designator
// (default the calculator's current default unit). The original
// program's own parser copies at most two raw bytes past the number
// for the unit, which is exactly right for its two-letter unit codes
// but is fragile for anything else; this conversion instead takes the
// full trailing letters as the unit and matches it exactly against
// the known unit names, an internal simplification with the same
// result for every legitimate entry.
func parseDataEntry(input string, defaultUnit int) (op int, value float64, unitIndex int, err error) {
	b := strings.TrimLeft(input, " ")
	op = opAdd
	if len(b) > 0 {
		switch b[0] {
		case '+':
			op, b = opAdd, b[1:]
		case '-':
			op, b = opSub, b[1:]
		case '*':
			op, b = opMul, b[1:]
		case '/':
			op, b = opDiv, b[1:]
		}
	}
	b = strings.TrimLeft(b, " ")

	i := 0
	for i < len(b) && !isLetter(b[i]) {
		i++
	}
	valStr := strings.ReplaceAll(b[:i], " ", "")
	unitStr := strings.TrimSpace(b[i:])
	if unitStr == "" {
		unitStr = unitNames[defaultUnit]
	}

	unitIndex = -1
	for idx, name := range unitNames {
		if name == unitStr {
			unitIndex = idx
			break
		}
	}
	if unitIndex == -1 {
		return 0, 0, 0, errUnknownUnit
	}

	if valStr == "" {
		return 0, 0, 0, errBadEntry
	}
	value, err = parseValue(valStr)
	if err != nil {
		return 0, 0, 0, err
	}
	return op, value, unitIndex, nil
}

func isLetter(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}

// fractionalFeetInches renders the feet and inches accumulators
// together as a mixed feet-and-fraction-of-an-inch string (e.g.
// "3 ft 9 & 25/32 in"), the display the original program shows
// whenever the running result is a plain length (exponent 1) with a
// non-zero feet accumulator. Because the feet and inches
// accumulators are independent running totals rather than one
// derived from the other, reducing the inches remainder to the
// nearest fracd'th can disagree slightly with the inches
// accumulator's own value; the reported error percentage is that
// disagreement, exactly as the original computes it.
func (s *calcState) fractionalFeetInches() (string, bool) {
	ftAcc := s.acc[unitFT-1]
	inAcc := s.acc[unitIN-1]
	if ftAcc == 0 || s.exponent != 1 {
		return "", false
	}

	feet, _ := math.Modf(math.Abs(ftAcc))
	totalIn := math.Abs(inAcc)
	inches, frac := math.Modf(totalIn - 12*feet)
	numerator, _ := math.Modf(frac * s.fracDenom)
	denominator := 0.0
	if numerator != 0 {
		g := gcdFloat(numerator, s.fracDenom)
		numerator /= g
		denominator = s.fracDenom / g
	}

	reconstructed := 12*feet + inches
	if numerator != 0 {
		reconstructed += numerator / denominator
	}
	errPct := 0.0
	if totalIn != 0 {
		errPct = 100 * (reconstructed - totalIn) / totalIn
	}

	var b strings.Builder
	if ftAcc < 0 {
		b.WriteString("-")
	}
	if feet != 0 {
		fmt.Fprintf(&b, "%d ft", int(feet))
	}
	switch {
	case inches != 0 && numerator == 0:
		fmt.Fprintf(&b, " %d in", int(inches))
	case inches != 0:
		fmt.Fprintf(&b, " %d & %d/%d in", int(inches), int(numerator), int(denominator))
	case numerator != 0:
		fmt.Fprintf(&b, " %d/%d in", int(numerator), int(denominator))
	default:
		b.WriteString(" 0 in")
	}
	fmt.Fprintf(&b, " (error = %.4f %%)", errPct)
	return b.String(), true
}

// displayLines renders the six accumulators, each with its unit
// suffix and any dimension exponent, followed by the mixed
// feet-and-inches line when applicable.
func (s *calcState) displayLines() []string {
	lines := make([]string, 0, numAccumulators+1)
	for i := 0; i < numAccumulators; i++ {
		line := numdisplay.Format(s.acc[i], s.decimals, s.selfAdjust, s.notation)
		if s.exponent != 0 {
			line += " " + unitNames[i+1]
			if s.exponent != 1 {
				line += fmt.Sprintf("^%d", s.exponent)
			}
		}
		lines = append(lines, line)
	}
	if frac, ok := s.fractionalFeetInches(); ok {
		lines = append(lines, frac)
	}
	return lines
}

func clampInt(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

func parseAssignedFloat(input string) (float64, bool) {
	idx := strings.IndexByte(input, '=')
	if idx == -1 {
		return 0, false
	}
	v, err := strconv.ParseFloat(strings.TrimSpace(input[idx+1:]), 64)
	if err != nil {
		return 0, false
	}
	return v, true
}

func setDefaultUnit(state *calcState, input string) bool {
	idx := strings.IndexByte(input, '=')
	if idx == -1 {
		return false
	}
	name := strings.TrimSpace(input[idx+1:])
	for i, u := range unitNames {
		if u == name {
			state.defaultUnit = i
			return true
		}
	}
	return false
}

func processEntry(state *calcState, input string) error {
	op, value, unitIndex, err := parseDataEntry(input, state.defaultUnit)
	if err != nil {
		return err
	}
	return state.applyEntry(op, unitIndex, value)
}

func printHelp() {
	fmt.Println()
	fmt.Println("SYNTAX:")
	fmt.Println("entry => number expressed in any of the forms: 7, .2, 1.23, 3/4, 4&3/8")
	fmt.Println("op => operator (+,-,*,/) if omitted '+' is assumed")
	fmt.Println("ud => unit designator (nd,m,cm,mm,yd,ft,in) if omitted default units assumed")
	fmt.Println("  nd is used to enter non-dimensional constants")
	fmt.Println("COMMANDS:")
	fmt.Println("clear => clear accumulators")
	fmt.Println("decimals=n => set number of decimal places to display (0<=n<=15)")
	fmt.Println("dflag=n => n=0 (decimal places); n=1 (self-adjusting)")
	fmt.Println("dtype=n => n=0 (fixed); n=1 (engineering); n=2 (scientific) display")
	fmt.Println("exit/quit => exit program")
	fmt.Println("fracd=n => set precision of fractional inches to 1/n in (default = 64)")
	fmt.Println("help => display this screen")
	fmt.Println("scale=n => set input scale factor")
	fmt.Println("undo => undo the last calculation performed")
	fmt.Println("unit=ud => set default units to ud")
	fmt.Println("DATA ENTRY:")
	fmt.Println("[op] entry [ud] => enter a value into the calculation")
	fmt.Println()
}

func printDisplay(state *calcState) {
	for _, line := range state.displayLines() {
		fmt.Println(line)
	}
	fmt.Println()
	fmt.Println("SYNTAX: op=+,-,*,/ ; entry = 7, .2, 1.23, 2/3, 1 & 2/3")
	fmt.Println("unit designators (ud) = nd,m,cm,mm,yd,ft,in")
	fmt.Println("COMMANDS: clear; decimals=(0-15); dflag=(0-1); dtype=(0-2); exit")
	fmt.Println("fracd=n; help; quit; scale=n; undo; unit=(ud)")
	fmt.Println("DATA ENTRY: [op] entry [ud]")
	fmt.Printf("Fractional approximation accuracy = 1/%.0f in\n", state.fracDenom)
	fmt.Printf("Scale factor = %.5f\n", state.scale)
	fmt.Printf("Default unit = %s\n", unitNames[state.defaultUnit])
}

func main() {
	prompter, err := promptio.New(os.Stdin, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, "mix:", err)
		os.Exit(1)
	}

	state := newCalcState()
	fmt.Println("MIXED DIMENSIONAL UNITS FOUR FUNCTION CALCULATOR")
	printHelp()

	for {
		printDisplay(state)
		input := strings.ToLower(prompter.Line("? "))

		switch {
		case strings.Contains(input, "quit"), strings.Contains(input, "exit"):
			return
		case strings.Contains(input, "help"):
			printHelp()
		case strings.Contains(input, "clear"):
			state.clear()
		case strings.Contains(input, "unit"):
			if !setDefaultUnit(state, input) {
				fmt.Println("INPUT ERROR")
			}
		case strings.Contains(input, "decimals"):
			if v, ok := parseAssignedFloat(input); ok {
				state.decimals = clampInt(int(v), 0, 15)
			} else {
				fmt.Println("INPUT ERROR")
			}
		case strings.Contains(input, "dflag"):
			if v, ok := parseAssignedFloat(input); ok {
				state.selfAdjust = v != 0
			} else {
				fmt.Println("INPUT ERROR")
			}
		case strings.Contains(input, "dtype"):
			if v, ok := parseAssignedFloat(input); ok {
				state.notation = int(v)
			} else {
				fmt.Println("INPUT ERROR")
			}
		case strings.Contains(input, "fracd"):
			if v, ok := parseAssignedFloat(input); ok {
				state.fracDenom = v
			} else {
				fmt.Println("INPUT ERROR")
			}
		case strings.Contains(input, "undo"):
			state.undo()
		case strings.Contains(input, "scale"):
			if v, ok := parseAssignedFloat(input); ok {
				state.scale = v
			} else {
				fmt.Println("INPUT ERROR")
			}
		default:
			if err := processEntry(state, input); err != nil {
				fmt.Println("INPUT ERROR")
			}
		}
	}
}
