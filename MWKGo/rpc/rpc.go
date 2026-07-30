// rpc is an RPN (Reverse Polish Notation) stack-oriented scientific
// calculator: a four-register stack (X, Y, Z, T), a memory register,
// and a wide set of named operators (trig, logs, powers, bitwise,
// unit conversion, number-base conversion) that all act on the stack.
//
// The original program is driven entirely by mouse clicks on an
// on-screen keypad (with a few keyboard shortcuts layered on top),
// building up a pending numeric entry one character at a time before
// committing it to the stack on the next operator. This conversion
// replaces both with a single command line: a line that parses as a
// number (in the current input base) pushes it onto the stack
// immediately, and any other recognized line runs the matching named
// operator. Because a whole number is always typed and committed in
// one step, there is no equivalent of the original's separate
// pending-entry buffer to convert; every operator that would have
// auto-entered a pending value now simply operates on values already
// on the stack.
//
// Converted from RPC.C (M. W. Klotz), Math/rpc.
package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"

	"mwkgo/internal/numdisplay"
	"mwkgo/internal/promptio"
)

const stackSize = 4 // X, Y, Z, T registers

// Angle modes.
const (
	angleDegrees = iota
	angleRadians
)

// Temperature scales.
const (
	tempF = iota
	tempC
	tempR
	tempK
	numTempScales
)

var tempScaleNames = [numTempScales]string{"F", "C", "R", "K"}

// Weight scales.
const (
	weightOz = iota
	weightLb
	weightG
	weightKg
	numWeightScales
)

var weightScaleNames = [numWeightScales]string{"oz", "lb", "gm", "kg"}

// Length scales.
const (
	lengthIn = iota
	lengthFt
	lengthYd
	lengthMi
	lengthMm
	lengthCm
	lengthM
	lengthKm
	numLengthScales
)

var lengthScaleNames = [numLengthScales]string{"in", "ft", "yd", "mi", "mm", "cm", "m", "km"}

// maxFactorialInput is the largest integer this calculator will take
// the factorial of; 171! overflows float64's exponent range (170! is
// approximately 7E+306, matching the original's own documented
// limit).
const maxFactorialInput = 170

// maxExpArgument bounds arguments to exp() and pow() so the result
// stays within float64's representable range (base-10 exponent up to
// about 300), matching the original's own guard.
const maxExpArgument = 300.0

// maxBitwiseMagnitude is the largest magnitude the bitwise operators
// accept, matching the original's 4-byte long integer range.
const maxBitwiseMagnitude = 0xffffffff

var (
	errDivideByZero = fmt.Errorf("division by zero")
	errDomain       = fmt.Errorf("value out of range for this operation")
)

// calcState holds the stack, memory, display settings, and unit/base
// selections, matching the original program's global variables. A
// full second copy (secondary) provides the "swap to a secondary
// calculator" feature for auxiliary calculations.
type calcState struct {
	stack [stackSize]float64
	lastX float64
	mem   float64

	undoStack [stackSize]float64
	undoLastX float64
	undoMem   float64

	secondary [stackSize + 2]float64 // stack[0..3], lastX, mem

	decimals   int
	selfAdjust bool
	notation   int
	base       int
	angleMode  int

	tempFrom, tempTo     int
	weightFrom, weightTo int
	lengthFrom, lengthTo int
}

func newCalcState() *calcState {
	return &calcState{
		decimals: 2,
		base:     10,
		weightTo: weightKg,
		lengthTo: lengthMm,
	}
}

// save snapshots the pre-operation stack, lastX, and memory into the
// undo registers, then updates lastX to the value about to be
// operated on (the x register), exactly matching the original
// program's own save() function: "last x" always reflects the x
// register's value just before the most recently performed
// operation.
func (s *calcState) save() {
	s.undoStack = s.stack
	s.undoLastX = s.lastX
	s.undoMem = s.mem
	s.lastX = s.stack[0]
}

func (s *calcState) restoreUndo() {
	s.stack = s.undoStack
	s.lastX = s.undoLastX
	s.mem = s.undoMem
}

func (s *calcState) push() {
	for i := stackSize - 1; i > 0; i-- {
		s.stack[i] = s.stack[i-1]
	}
}

func (s *calcState) pop() {
	for i := 0; i < stackSize-1; i++ {
		s.stack[i] = s.stack[i+1]
	}
}

func (s *calcState) rollUp() {
	t := s.stack[stackSize-1]
	s.push()
	s.stack[0] = t
}

func (s *calcState) rollDown() {
	t := s.stack[0]
	s.pop()
	s.stack[stackSize-1] = t
}

// enterValue pushes x onto the stack as a new entry. Matching the
// original's entr(), this updates lastX but is not itself undoable.
func (s *calcState) enterValue(x float64) {
	s.lastX = s.stack[0]
	s.push()
	s.stack[0] = x
}

// gcdFloat returns the greatest common divisor of x and y, both
// assumed non-negative and not both zero, via the Euclidean
// algorithm carried out in floating point (needed here because gcd
// and lcm work directly on the float64 stack registers).
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

func angleToRadiansFactor(s *calcState) float64 {
	if s.angleMode == angleRadians {
		return 1
	}
	return math.Pi / 180
}

func radiansToAngleFactor(s *calcState) float64 {
	if s.angleMode == angleRadians {
		return 1
	}
	return 180 / math.Pi
}

// --- arithmetic -------------------------------------------------------

func opAdd(s *calcState) error {
	s.save()
	t := s.stack[1] + s.stack[0]
	s.pop()
	s.stack[0] = t
	return nil
}

func opSub(s *calcState) error {
	s.save()
	t := s.stack[1] - s.stack[0]
	s.pop()
	s.stack[0] = t
	return nil
}

func opMul(s *calcState) error {
	s.save()
	t := s.stack[1] * s.stack[0]
	s.pop()
	s.stack[0] = t
	return nil
}

func opDiv(s *calcState) error {
	if s.stack[0] == 0 {
		return errDivideByZero
	}
	s.save()
	t := s.stack[1] / s.stack[0]
	s.pop()
	s.stack[0] = t
	return nil
}

func opXOverY(s *calcState) error {
	if s.stack[1] == 0 {
		return errDivideByZero
	}
	s.save()
	t := s.stack[0] / s.stack[1]
	s.pop()
	s.stack[0] = t
	return nil
}

func opExchangeXY(s *calcState) error {
	s.stack[0], s.stack[1] = s.stack[1], s.stack[0]
	return nil
}

func opExchangeXMem(s *calcState) error {
	s.stack[0], s.mem = s.mem, s.stack[0]
	return nil
}

// --- stack, memory, and constants --------------------------------------

func opE(s *calcState) error {
	s.save()
	s.push()
	s.stack[0] = math.E
	return nil
}

// opLastX is not itself undoable, matching the original (no save()
// call before this operation).
func opLastX(s *calcState) error {
	s.push()
	s.stack[0] = s.lastX
	return nil
}

func opReciprocal(s *calcState) error {
	if s.stack[0] == 0 {
		return errDivideByZero
	}
	s.save()
	s.stack[0] = 1 / s.stack[0]
	return nil
}

// opChangeSign is not undoable, matching the original (chs() acts
// directly on the x register with no save() call).
func opChangeSign(s *calcState) error {
	s.stack[0] = -s.stack[0]
	return nil
}

func opPi(s *calcState) error {
	s.save()
	s.push()
	s.stack[0] = math.Pi
	return nil
}

func opSquare(s *calcState) error {
	x := s.stack[0]
	s.save()
	s.stack[0] = x * x
	return nil
}

func opSqrt(s *calcState) error {
	if s.stack[0] < 0 {
		return errDomain
	}
	x := s.stack[0]
	s.save()
	s.stack[0] = math.Sqrt(x)
	return nil
}

func opFactorial(s *calcState) error {
	intPart, frac := math.Modf(math.Abs(s.stack[0]))
	if frac != 0 || s.stack[0] < 0 || intPart > maxFactorialInput {
		return errDomain
	}
	result := 1.0
	for i := 1; i <= int(intPart); i++ {
		result *= float64(i)
	}
	s.save()
	s.stack[0] = result
	return nil
}

// store operations are not undoable, matching the original (no
// save() call before this switch block): they only ever change mem,
// which is restored by Undo only when captured by a prior save().
func opStore(s *calcState) error    { s.mem = s.stack[0]; return nil }
func opStoreAdd(s *calcState) error { s.mem += s.stack[0]; return nil }
func opStoreSub(s *calcState) error { s.mem -= s.stack[0]; return nil }
func opStoreMul(s *calcState) error { s.mem *= s.stack[0]; return nil }
func opStoreDiv(s *calcState) error {
	if s.stack[0] == 0 {
		return errDivideByZero
	}
	s.mem /= s.stack[0]
	return nil
}

func opRecall(s *calcState) error {
	s.save()
	s.push()
	s.stack[0] = s.mem
	return nil
}

func opRecallAdd(s *calcState) error {
	s.save()
	s.stack[0] += s.mem
	return nil
}

func opRecallSub(s *calcState) error {
	s.save()
	s.stack[0] -= s.mem
	return nil
}

func opRecallMul(s *calcState) error {
	s.save()
	s.stack[0] *= s.mem
	return nil
}

func opRecallDiv(s *calcState) error {
	if s.mem == 0 {
		return errDivideByZero
	}
	s.save()
	s.stack[0] /= s.mem
	return nil
}

// --- trigonometry -------------------------------------------------------

func opSin(s *calcState) error {
	t := angleToRadiansFactor(s)
	s.save()
	s.stack[0] = math.Sin(s.stack[0] * t)
	return nil
}

func opCos(s *calcState) error {
	t := angleToRadiansFactor(s)
	s.save()
	s.stack[0] = math.Cos(s.stack[0] * t)
	return nil
}

func opTan(s *calcState) error {
	t := angleToRadiansFactor(s)
	s.save()
	s.stack[0] = math.Tan(s.stack[0] * t)
	return nil
}

// opRSS computes the root-sum-square of x and y via the original's
// own explicit sqrt(x*x+y*y) formula, distinct from opRectToPolar
// below, which uses the C library's hypot() as the original does.
func opRSS(s *calcState) error {
	s.save()
	t := math.Sqrt(s.stack[0]*s.stack[0] + s.stack[1]*s.stack[1])
	s.pop()
	s.stack[0] = t
	return nil
}

func opRectToPolar(s *calcState) error {
	x, y := s.stack[0], s.stack[1]
	s.save()
	r := math.Hypot(x, y)
	theta := 0.0
	if x != 0 || y != 0 {
		theta = math.Atan2(y, x) / angleToRadiansFactor(s)
	}
	s.stack[0] = r
	s.stack[1] = theta
	return nil
}

func opAsin(s *calcState) error {
	if math.Abs(s.stack[0]) > 1 {
		return errDomain
	}
	t := radiansToAngleFactor(s)
	s.save()
	s.stack[0] = t * math.Asin(s.stack[0])
	return nil
}

func opAcos(s *calcState) error {
	if math.Abs(s.stack[0]) > 1 {
		return errDomain
	}
	t := radiansToAngleFactor(s)
	s.save()
	s.stack[0] = t * math.Acos(s.stack[0])
	return nil
}

func opAtan(s *calcState) error {
	t := radiansToAngleFactor(s)
	s.save()
	s.stack[0] = t * math.Atan(s.stack[0])
	return nil
}

func opUnRSS(s *calcState) error {
	s.save()
	t := s.stack[1]*s.stack[1] - s.stack[0]*s.stack[0]
	if t < 0 {
		t = -t
	}
	s.pop()
	s.stack[0] = math.Sqrt(t)
	return nil
}

func opPolarToRect(s *calcState) error {
	r, theta := s.stack[0], s.stack[1]
	s.save()
	thetaRad := theta * angleToRadiansFactor(s)
	s.stack[0] = r * math.Cos(thetaRad)
	s.stack[1] = r * math.Sin(thetaRad)
	return nil
}

func opDegToRad(s *calcState) error {
	s.save()
	s.stack[0] *= math.Pi / 180
	return nil
}

func opRadToDeg(s *calcState) error {
	s.save()
	s.stack[0] *= 180 / math.Pi
	return nil
}

func opAtan2(s *calcState) error {
	t := 0.0
	if s.stack[0] != 0 || s.stack[1] != 0 {
		t = math.Atan2(s.stack[1], s.stack[0])
	}
	if s.angleMode == angleDegrees {
		t *= 180 / math.Pi
	}
	s.save()
	s.pop()
	s.stack[0] = t
	return nil
}

// --- fractional and rounding ---------------------------------------------

func opFrac(s *calcState) error {
	_, frac := math.Modf(s.stack[0])
	s.save()
	s.stack[0] = frac
	return nil
}

func opInt(s *calcState) error {
	intPart, _ := math.Modf(s.stack[0])
	s.save()
	s.stack[0] = intPart
	return nil
}

func opSplit(s *calcState) error {
	intPart, frac := math.Modf(s.stack[0])
	s.save()
	s.push()
	s.stack[0] = frac
	s.stack[1] = intPart
	return nil
}

func opYModX(s *calcState) error {
	s.save()
	s.stack[0] = math.Mod(s.stack[1], s.stack[0])
	return nil
}

func opFloor(s *calcState) error {
	s.save()
	s.stack[0] = math.Floor(s.stack[0])
	return nil
}

func opCeil(s *calcState) error {
	s.save()
	s.stack[0] = math.Ceil(s.stack[0])
	return nil
}

func opRound(s *calcState) error {
	s.save()
	rounded, err := strconv.ParseFloat(strconv.FormatFloat(s.stack[0], 'f', s.decimals, 64), 64)
	if err != nil {
		return errDomain
	}
	s.stack[0] = rounded
	return nil
}

func opGCD(s *calcState) error {
	x, fracX := math.Modf(math.Abs(s.stack[0]))
	y, fracY := math.Modf(math.Abs(s.stack[1]))
	if fracX != 0 || x == 0 || fracY != 0 || y == 0 {
		return errDomain
	}
	s.save()
	s.pop()
	s.stack[0] = gcdFloat(x, y)
	return nil
}

func opLCM(s *calcState) error {
	x, fracX := math.Modf(math.Abs(s.stack[0]))
	y, fracY := math.Modf(math.Abs(s.stack[1]))
	if fracX != 0 || x == 0 || fracY != 0 || y == 0 {
		return errDomain
	}
	s.save()
	s.pop()
	s.stack[0] = x * y / gcdFloat(x, y)
	return nil
}

// --- exponentials and logarithms -----------------------------------------

func opExpE(s *calcState) error {
	if math.Abs(s.stack[0]) > maxExpArgument*math.Log(10) {
		return errDomain
	}
	s.save()
	s.stack[0] = math.Exp(s.stack[0])
	return nil
}

func opExp10(s *calcState) error {
	if math.Abs(s.stack[0]) > maxExpArgument {
		return errDomain
	}
	s.save()
	s.stack[0] = math.Pow(10, s.stack[0])
	return nil
}

func opExp2(s *calcState) error {
	if math.Abs(s.stack[0]) > maxExpArgument*math.Log(10)/math.Log(2) {
		return errDomain
	}
	s.save()
	s.stack[0] = math.Pow(2, s.stack[0])
	return nil
}

func opPowYX(s *calcState) error {
	if s.stack[1] <= 0 || math.Abs(s.stack[0]) > maxExpArgument*math.Log(10)/math.Log(s.stack[1]) {
		return errDomain
	}
	s.save()
	t := math.Pow(s.stack[1], s.stack[0])
	s.pop()
	s.stack[0] = t
	return nil
}

func opRootYX(s *calcState) error {
	if s.stack[0] == 0 || s.stack[1] <= 0 {
		return errDomain
	}
	if math.Abs(1/s.stack[0]) > maxExpArgument*math.Log(10)/math.Log(s.stack[1]) {
		return errDomain
	}
	s.save()
	t := math.Pow(s.stack[1], 1/s.stack[0])
	s.pop()
	s.stack[0] = t
	return nil
}

func opLn(s *calcState) error {
	if s.stack[0] <= 0 {
		return errDomain
	}
	s.save()
	s.stack[0] = math.Log(s.stack[0])
	return nil
}

func opLog10(s *calcState) error {
	if s.stack[0] <= 0 {
		return errDomain
	}
	s.save()
	s.stack[0] = math.Log10(s.stack[0])
	return nil
}

func opLog2(s *calcState) error {
	if s.stack[0] <= 0 {
		return errDomain
	}
	s.save()
	s.stack[0] = math.Log(s.stack[0]) / math.Log(2)
	return nil
}

// opLogYX computes a logarithm of the y register to a base taken
// from lastX, not from the x register, exactly reproducing the
// original's own ylogx operator: the x register is checked for
// positivity (shared with ln/log/log2 above) but the actual base
// used in the formula is lastX, the value that was in the x register
// before the previous operation. This looks like a mismatch between
// the guard and the formula, but it is exactly what the original
// computes, so it is preserved rather than "corrected".
func opLogYX(s *calcState) error {
	if s.stack[0] <= 0 {
		return errDomain
	}
	base := s.lastX
	arg := s.stack[1]
	if arg <= 0 || math.Log(base) == 0 {
		return errDomain
	}
	s.save()
	s.pop()
	s.stack[0] = math.Log(arg) / math.Log(base)
	return nil
}

// --- bitwise (operating on the original's 4-byte long range) -------------

func opAnd(s *calcState) error {
	if math.Abs(s.stack[0]) > maxBitwiseMagnitude || math.Abs(s.stack[1]) > maxBitwiseMagnitude {
		return errDomain
	}
	s.save()
	x := int32(s.stack[0]) & int32(s.stack[1])
	s.pop()
	s.stack[0] = float64(x)
	return nil
}

func opOr(s *calcState) error {
	if math.Abs(s.stack[0]) > maxBitwiseMagnitude || math.Abs(s.stack[1]) > maxBitwiseMagnitude {
		return errDomain
	}
	s.save()
	x := int32(s.stack[0]) | int32(s.stack[1])
	s.pop()
	s.stack[0] = float64(x)
	return nil
}

func opXor(s *calcState) error {
	if math.Abs(s.stack[0]) > maxBitwiseMagnitude || math.Abs(s.stack[1]) > maxBitwiseMagnitude {
		return errDomain
	}
	s.save()
	x := int32(s.stack[0]) ^ int32(s.stack[1])
	s.pop()
	s.stack[0] = float64(x)
	return nil
}

func opOnesComplement(s *calcState) error {
	if math.Abs(s.stack[0]) > maxBitwiseMagnitude {
		return errDomain
	}
	s.save()
	x := ^int32(s.stack[0])
	s.stack[0] = float64(x)
	return nil
}

// --- unit conversion -------------------------------------------------

// tempConvert converts x from the from temperature scale to the to
// scale. Like the original tconv(), an out-of-physical-range input
// (below absolute zero on the source scale) leaves x unchanged
// (reported to the caller via ok=false) rather than erroring loudly.
func tempConvert(x float64, from, to int) (result float64, ok bool) {
	const c0 = 273.16
	r0 := 9*c0/5 - 32
	var c float64
	switch from {
	case tempF:
		if x < -r0 {
			return x, false
		}
		c = 5 * (x - 32) / 9
	case tempC:
		if x < -c0 {
			return x, false
		}
		c = x
	case tempR:
		if x < 0 {
			return x, false
		}
		c = 5 * ((x - r0) - 32) / 9
	case tempK:
		if x < 0 {
			return x, false
		}
		c = x - c0
	}
	switch to {
	case tempF:
		return 9*c/5 + 32, true
	case tempC:
		return c, true
	case tempR:
		return 9*c/5 + 32 + r0, true
	case tempK:
		return c + c0, true
	}
	return x, false
}

func weightConvert(x float64, from, to int) float64 {
	const poundsPerKg = 2.20462
	var kg float64
	switch from {
	case weightOz:
		kg = x / (16 * poundsPerKg)
	case weightLb:
		kg = x / poundsPerKg
	case weightG:
		kg = x * 1e-3
	case weightKg:
		kg = x
	}
	switch to {
	case weightOz:
		return 16 * poundsPerKg * kg
	case weightLb:
		return kg * poundsPerKg
	case weightG:
		return kg * 1000
	case weightKg:
		return kg
	}
	return kg
}

func lengthConvert(x float64, from, to int) float64 {
	const metersPerInch = 0.0254
	var m float64
	switch from {
	case lengthIn:
		m = x * metersPerInch
	case lengthFt:
		m = x * 12 * metersPerInch
	case lengthYd:
		m = x * 36 * metersPerInch
	case lengthMi:
		m = x * 12 * 5280 * metersPerInch
	case lengthMm:
		m = x * 1e-3
	case lengthCm:
		m = x * 1e-2
	case lengthM:
		m = x
	case lengthKm:
		m = x * 1000
	}
	switch to {
	case lengthIn:
		return m / metersPerInch
	case lengthFt:
		return m / (12 * metersPerInch)
	case lengthYd:
		return m / (36 * metersPerInch)
	case lengthMi:
		return m / (12 * 5280 * metersPerInch)
	case lengthMm:
		return m * 1000
	case lengthCm:
		return m * 100
	case lengthM:
		return m
	case lengthKm:
		return m * 1e-3
	}
	return m
}

func opTempConvert(s *calcState) error {
	s.save()
	result, ok := tempConvert(s.stack[0], s.tempFrom, s.tempTo)
	s.stack[0] = result
	if !ok {
		return errDomain
	}
	return nil
}

func opWeightConvert(s *calcState) error {
	s.save()
	s.stack[0] = weightConvert(s.stack[0], s.weightFrom, s.weightTo)
	return nil
}

func opLengthConvert(s *calcState) error {
	s.save()
	s.stack[0] = lengthConvert(s.stack[0], s.lengthFrom, s.lengthTo)
	return nil
}

// --- secondary calculator -------------------------------------------------

func opSwap(s *calcState) error {
	var saved [stackSize + 2]float64
	copy(saved[:stackSize], s.stack[:])
	saved[stackSize] = s.lastX
	saved[stackSize+1] = s.mem

	copy(s.stack[:], s.secondary[:stackSize])
	s.lastX = s.secondary[stackSize]
	s.mem = s.secondary[stackSize+1]
	s.secondary = saved
	return nil
}

func opSwapX(s *calcState) error {
	s.stack[0], s.secondary[0] = s.secondary[0], s.stack[0]
	return nil
}

// operatorTable maps a typed command name to the operator it runs.
var operatorTable = map[string]func(*calcState) error{
	"+": opAdd, "-": opSub, "*": opMul, "/": opDiv, "x/y": opXOverY,
	"xy": opExchangeXY, "xm": opExchangeXMem,
	"e": opE, "lastx": opLastX, "1/x": opReciprocal, "chs": opChangeSign,
	"pi": opPi, "sqr": opSquare, "sqrt": opSqrt, "fact": opFactorial,
	"store": opStore, "store+": opStoreAdd, "store-": opStoreSub, "store*": opStoreMul, "store/": opStoreDiv,
	"rcall": opRecall, "rcall+": opRecallAdd, "rcall-": opRecallSub, "rcall*": opRecallMul, "rcall/": opRecallDiv,
	"sin": opSin, "cos": opCos, "tan": opTan, "rss": opRSS, "topolar": opRectToPolar,
	"asin": opAsin, "acos": opAcos, "atan": opAtan, "unrss": opUnRSS, "torect": opPolarToRect,
	"deg2rad": opDegToRad, "rad2deg": opRadToDeg, "atan2": opAtan2,
	"frac": opFrac, "int": opInt, "split": opSplit, "ymodx": opYModX,
	"floor": opFloor, "ceil": opCeil, "round": opRound, "gcd": opGCD, "lcm": opLCM,
	"e^x": opExpE, "10^x": opExp10, "2^x": opExp2, "y^x": opPowYX, "y^(1/x)": opRootYX,
	"ln": opLn, "log": opLog10, "log2": opLog2, "ylogx": opLogYX,
	"and": opAnd, "or": opOr, "1comp": opOnesComplement, "xor": opXor,
	"temp": opTempConvert, "weight": opWeightConvert, "length": opLengthConvert,
	"swap": opSwap, "swapx": opSwapX,
	"roll":   func(s *calcState) error { s.rollUp(); return nil },
	"rolldn": func(s *calcState) error { s.rollDown(); return nil },
}

func parseNumericEntry(token string, base int) (float64, bool) {
	if base == 10 {
		v, err := strconv.ParseFloat(token, 64)
		if err != nil {
			return 0, false
		}
		return v, true
	}
	v, err := strconv.ParseUint(token, base, 64)
	if err != nil {
		return 0, false
	}
	return float64(v), true
}

// formatIntegerBases returns the hex, octal, and binary
// representations of x, matching the original display's automatic
// base breakdown, shown whenever x is a whole number that fits a
// 4-byte long. A negative value is shown as the two's-complement bit
// pattern of that 4-byte long, matching how the original C compiler
// truncated a negative double into an unsigned long for display.
func formatIntegerBases(x float64) (hex, oct, bin string, ok bool) {
	intPart, frac := math.Modf(x)
	if frac != 0 || math.Abs(intPart) > maxBitwiseMagnitude {
		return "", "", "", false
	}
	u := uint32(int32(intPart))
	return strconv.FormatUint(uint64(u), 16), strconv.FormatUint(uint64(u), 8), strconv.FormatUint(uint64(u), 2), true
}

func displayLines(s *calcState) []string {
	fmtVal := func(x float64) string { return numdisplay.Format(x, s.decimals, s.selfAdjust, s.notation) }
	lines := []string{
		"M (mem)   = " + fmtVal(s.mem),
		"L (lastx) = " + fmtVal(s.lastX),
		"T         = " + fmtVal(s.stack[3]),
		"Z         = " + fmtVal(s.stack[2]),
		"Y         = " + fmtVal(s.stack[1]),
		"X         = " + fmtVal(s.stack[0]),
	}
	if hex, oct, bin, ok := formatIntegerBases(s.stack[0]); ok {
		lines = append(lines, fmt.Sprintf("HEX: %s  OCT: %s  BIN: %s", hex, oct, bin))
	}
	lines = append(lines, fmt.Sprintf(
		"temp %s->%s  weight %s->%s  length %s->%s  base %d  angle %s  %s",
		tempScaleNames[s.tempFrom], tempScaleNames[s.tempTo],
		weightScaleNames[s.weightFrom], weightScaleNames[s.weightTo],
		lengthScaleNames[s.lengthFrom], lengthScaleNames[s.lengthTo],
		s.base, angleModeName(s.angleMode), notationName(s.notation)))
	return lines
}

func angleModeName(mode int) string {
	if mode == angleRadians {
		return "RAD"
	}
	return "DEG"
}

func notationName(notation int) string {
	switch notation {
	case numdisplay.Engineering:
		return "ENG"
	case numdisplay.Scientific:
		return "SCI"
	default:
		return "FIX"
	}
}

func printMenu() {
	fmt.Println("RPN SCIENTIFIC CALCULATOR")
	fmt.Println("enter a number (in the current base) to push it onto the stack")
	fmt.Println()
	fmt.Println("+ - * / x/y          basic arithmetic (x/y divides y by x)")
	fmt.Println("roll rolldn xy xm    stack rearrangement (roll=up, xy=x<>y, xm=x<>mem)")
	fmt.Println("e lastx 1/x chs      constants and single-value ops")
	fmt.Println("pi sqr sqrt fact     pi, x^2, sqrt(x), x!")
	fmt.Println("store store+-*/      store into memory (default/add/sub/mul/div)")
	fmt.Println("rcall rcall+-*/      recall from memory (default/add/sub/mul/div)")
	fmt.Println("sin cos tan rss      trig functions; rss = sqrt(x^2+y^2)")
	fmt.Println("asin acos atan unrss inverse trig; unrss = sqrt(|y^2-x^2|)")
	fmt.Println("topolar torect       rect(x,y) <-> polar(r,theta)")
	fmt.Println("deg2rad rad2deg atan2")
	fmt.Println("frac int split ymodx fractional/integer parts, y mod x")
	fmt.Println("floor ceil round gcd lcm")
	fmt.Println("e^x 10^x 2^x y^x y^(1/x)")
	fmt.Println("ln log log2 ylogx    natural, base 10, base 2, log base lastx of y")
	fmt.Println("and or 1comp xor     bitwise (4-byte integer range)")
	fmt.Println("temp weight length   unit conversion using the current from/to scales")
	fmt.Println("tempfrom tempto weightfrom weightto lengthfrom lengthto")
	fmt.Println("                     cycle the from/to scale used by temp/weight/length")
	fmt.Println("fix eng sci          display notation")
	fmt.Println("dp=n adj             set decimal places / self-adjusting decimals")
	fmt.Println("dec hex bin oct      number input/display base")
	fmt.Println("deg rad              angle mode")
	fmt.Println("clearx clrstk clrmem clrall")
	fmt.Println("swap swapx           swap with / swap x with the secondary calculator")
	fmt.Println("undo                 undo the last operation")
	fmt.Println("notes                show usage notes")
	fmt.Println("quit                 exit")
	fmt.Println()
}

func printNotes() {
	fmt.Println("Computational accuracy: about 15 significant digits.")
	fmt.Println("Usable range: about 1E+/-300.")
	fmt.Println("Bitwise operators are limited to a 4-byte integer range.")
	fmt.Println("ADJ sets decimal places to display each number to full (about 1E-15) precision.")
	fmt.Println("A value that would display as 0.0... in FIX mode instead grows enough")
	fmt.Println(" decimal places to show something non-zero.")
	fmt.Println("Magnitudes with log10 beyond 15 automatically switch to engineering notation")
	fmt.Println(" even if FIX is active.")
	fmt.Println()
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

func parseAssignedInt(input string) (int, bool) {
	idx := strings.IndexByte(input, '=')
	if idx == -1 {
		return 0, false
	}
	v, err := strconv.Atoi(strings.TrimSpace(input[idx+1:]))
	if err != nil {
		return 0, false
	}
	return v, true
}

func main() {
	prompter, err := promptio.New(os.Stdin, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, "rpc:", err)
		os.Exit(1)
	}

	state := newCalcState()
	printMenu()

	for {
		for _, line := range displayLines(state) {
			fmt.Println(line)
		}
		input := strings.ToLower(strings.TrimSpace(prompter.Line("> ")))

		switch {
		case input == "quit" || input == "q":
			return
		case input == "notes" || input == "help":
			printNotes()
		case input == "undo":
			state.restoreUndo()
		case input == "fix":
			state.notation = numdisplay.Fixed
		case input == "eng":
			state.notation = numdisplay.Engineering
		case input == "sci":
			state.notation = numdisplay.Scientific
		case input == "adj":
			state.selfAdjust = true
		case strings.HasPrefix(input, "dp="):
			if v, ok := parseAssignedInt(input); ok {
				state.decimals = clampInt(v, 0, 15)
				state.selfAdjust = false
			} else {
				fmt.Println("INPUT ERROR")
			}
		case input == "dec":
			state.base = 10
		case input == "hex":
			state.base = 16
		case input == "bin":
			state.base = 2
		case input == "oct":
			state.base = 8
		case input == "deg":
			state.angleMode = angleDegrees
		case input == "rad":
			state.angleMode = angleRadians
		case input == "clearx":
			state.stack[0] = 0
		case input == "clrstk":
			state.stack = [stackSize]float64{}
		case input == "clrmem":
			state.mem = 0
		case input == "clrall":
			state.stack = [stackSize]float64{}
			state.mem = 0
			state.lastX = 0
		case input == "tempfrom":
			state.tempFrom = (state.tempFrom + 1) % numTempScales
		case input == "tempto":
			state.tempTo = (state.tempTo + 1) % numTempScales
		case input == "weightfrom":
			state.weightFrom = (state.weightFrom + 1) % numWeightScales
		case input == "weightto":
			state.weightTo = (state.weightTo + 1) % numWeightScales
		case input == "lengthfrom":
			state.lengthFrom = (state.lengthFrom + 1) % numLengthScales
		case input == "lengthto":
			state.lengthTo = (state.lengthTo + 1) % numLengthScales
		default:
			if fn, ok := operatorTable[input]; ok {
				if err := fn(state); err != nil {
					fmt.Println(err)
				}
			} else if v, ok := parseNumericEntry(input, state.base); ok {
				state.enterValue(v)
			} else {
				fmt.Println("INPUT ERROR")
			}
		}
	}
}
