package main

import (
	"math"
	"strings"
	"testing"
)

func newTestState() *calcState {
	return newCalcState()
}

func TestArithmetic(t *testing.T) {
	s := newTestState()
	s.enterValue(3)
	s.enterValue(4)
	if err := opAdd(s); err != nil {
		t.Fatalf("opAdd error = %v", err)
	}
	if s.stack[0] != 7 {
		t.Errorf("3+4 = %v, want 7", s.stack[0])
	}
}

func TestSubtractOrderMatchesRPNConvention(t *testing.T) {
	// In RPN, "10 enter 3 -" computes y - x = 10 - 3 = 7: the second
	// entry (x, top of stack) is subtracted from the first (y).
	s := newTestState()
	s.enterValue(10)
	s.enterValue(3)
	if err := opSub(s); err != nil {
		t.Fatalf("opSub error = %v", err)
	}
	if s.stack[0] != 7 {
		t.Errorf("10 enter 3 - = %v, want 7 (y - x)", s.stack[0])
	}
}

func TestDivideByZero(t *testing.T) {
	s := newTestState()
	s.enterValue(5)
	s.enterValue(0)
	if err := opDiv(s); err != errDivideByZero {
		t.Errorf("opDiv by zero error = %v, want errDivideByZero", err)
	}
}

func TestXOverY(t *testing.T) {
	// x/y divides y (the divisor) into x (the dividend): with 2 then
	// 10 entered, x=10, y=2, so x/y computes stk[0]/stk[1] = 10/2.
	s := newTestState()
	s.enterValue(2)
	s.enterValue(10)
	if err := opXOverY(s); err != nil {
		t.Fatalf("opXOverY error = %v", err)
	}
	if s.stack[0] != 5 {
		t.Errorf("x/y = %v, want 5", s.stack[0])
	}
}

func TestStackRotation(t *testing.T) {
	s := newTestState()
	for _, v := range []float64{1, 2, 3, 4} {
		s.enterValue(v)
	}
	// stack is now X=4,Y=3,Z=2,T=1
	s.rollUp()
	want := [stackSize]float64{1, 4, 3, 2}
	if s.stack != want {
		t.Errorf("rollUp() = %v, want %v", s.stack, want)
	}
	s.rollDown()
	want = [stackSize]float64{4, 3, 2, 1}
	if s.stack != want {
		t.Errorf("rollDown() = %v, want %v", s.stack, want)
	}
}

func TestFactorial_KnownValue(t *testing.T) {
	s := newTestState()
	s.enterValue(5)
	if err := opFactorial(s); err != nil {
		t.Fatalf("opFactorial error = %v", err)
	}
	if s.stack[0] != 120 {
		t.Errorf("5! = %v, want 120", s.stack[0])
	}
}

func TestFactorial_RejectsNonIntegerAndNegative(t *testing.T) {
	for _, v := range []float64{2.5, -3} {
		s := newTestState()
		s.enterValue(v)
		if err := opFactorial(s); err != errDomain {
			t.Errorf("opFactorial(%v) error = %v, want errDomain", v, err)
		}
	}
}

func TestGCDAndLCM_KnownValues(t *testing.T) {
	s := newTestState()
	s.enterValue(12)
	s.enterValue(18)
	if err := opGCD(s); err != nil {
		t.Fatalf("opGCD error = %v", err)
	}
	if s.stack[0] != 6 {
		t.Errorf("gcd(12,18) = %v, want 6", s.stack[0])
	}

	s = newTestState()
	s.enterValue(4)
	s.enterValue(6)
	if err := opLCM(s); err != nil {
		t.Fatalf("opLCM error = %v", err)
	}
	if s.stack[0] != 12 {
		t.Errorf("lcm(4,6) = %v, want 12", s.stack[0])
	}
}

// TestFracIntSplit is self-verifying: frac + int must reconstruct the
// original value for any input.
func TestFracIntSplit(t *testing.T) {
	x := 3.75
	s := newTestState()
	s.enterValue(x)
	if err := opFrac(s); err != nil {
		t.Fatalf("opFrac error = %v", err)
	}
	frac := s.stack[0]

	s = newTestState()
	s.enterValue(x)
	if err := opInt(s); err != nil {
		t.Fatalf("opInt error = %v", err)
	}
	intPart := s.stack[0]

	if math.Abs((intPart+frac)-x) > 1e-12 {
		t.Errorf("int(%v)+frac(%v) = %v, want %v", intPart, frac, intPart+frac, x)
	}
	if intPart != 3 || math.Abs(frac-0.75) > 1e-12 {
		t.Errorf("int=%v frac=%v, want int=3 frac=0.75", intPart, frac)
	}

	s = newTestState()
	s.enterValue(x)
	if err := opSplit(s); err != nil {
		t.Fatalf("opSplit error = %v", err)
	}
	if s.stack[0] != frac || s.stack[1] != intPart {
		t.Errorf("split gives x=%v y=%v, want x=%v y=%v (frac in x, int in y)", s.stack[0], s.stack[1], frac, intPart)
	}
}

// TestTrig_DegreesRoundTrip is self-verifying: asin(sin(x)) should
// reproduce x for x in [-90,90], in the default degrees angle mode.
func TestTrig_DegreesRoundTrip(t *testing.T) {
	for _, deg := range []float64{0, 30, 45, 60, -30} {
		s := newTestState()
		s.enterValue(deg)
		if err := opSin(s); err != nil {
			t.Fatalf("opSin(%v) error = %v", deg, err)
		}
		if err := opAsin(s); err != nil {
			t.Fatalf("opAsin error = %v", err)
		}
		if math.Abs(s.stack[0]-deg) > 1e-9 {
			t.Errorf("asin(sin(%v deg)) = %v", deg, s.stack[0])
		}
	}
}

func TestTrig_RadiansMode(t *testing.T) {
	s := newTestState()
	s.angleMode = angleRadians
	s.enterValue(math.Pi / 2)
	if err := opSin(s); err != nil {
		t.Fatalf("opSin error = %v", err)
	}
	if math.Abs(s.stack[0]-1) > 1e-12 {
		t.Errorf("sin(pi/2 rad) = %v, want 1", s.stack[0])
	}
}

// TestRectPolarRoundTrip is self-verifying: converting rect -> polar
// -> rect should reproduce the original x, y.
func TestRectPolarRoundTrip(t *testing.T) {
	s := newTestState()
	s.enterValue(3) // y register after next entry
	s.enterValue(4) // x register
	wantX, wantY := s.stack[0], s.stack[1]

	if err := opRectToPolar(s); err != nil {
		t.Fatalf("opRectToPolar error = %v", err)
	}
	if err := opPolarToRect(s); err != nil {
		t.Fatalf("opPolarToRect error = %v", err)
	}
	if math.Abs(s.stack[0]-wantX) > 1e-9 || math.Abs(s.stack[1]-wantY) > 1e-9 {
		t.Errorf("rect->polar->rect = (%v,%v), want (%v,%v)", s.stack[0], s.stack[1], wantX, wantY)
	}
}

// TestExpLogInverses is self-verifying: ln(e^x) and log10(10^x)
// should reproduce x.
func TestExpLogInverses(t *testing.T) {
	x := 2.5
	s := newTestState()
	s.enterValue(x)
	if err := opExpE(s); err != nil {
		t.Fatalf("opExpE error = %v", err)
	}
	if err := opLn(s); err != nil {
		t.Fatalf("opLn error = %v", err)
	}
	if math.Abs(s.stack[0]-x) > 1e-9 {
		t.Errorf("ln(e^%v) = %v", x, s.stack[0])
	}

	s = newTestState()
	s.enterValue(x)
	if err := opExp10(s); err != nil {
		t.Fatalf("opExp10 error = %v", err)
	}
	if err := opLog10(s); err != nil {
		t.Fatalf("opLog10 error = %v", err)
	}
	if math.Abs(s.stack[0]-x) > 1e-9 {
		t.Errorf("log10(10^%v) = %v", x, s.stack[0])
	}
}

func TestPowYX_KnownValue(t *testing.T) {
	s := newTestState()
	s.enterValue(2) // y
	s.enterValue(3) // x
	if err := opPowYX(s); err != nil {
		t.Fatalf("opPowYX error = %v", err)
	}
	if s.stack[0] != 8 {
		t.Errorf("2^3 = %v, want 8", s.stack[0])
	}
}

// TestLogYX_UsesLastXAsBase is a direct test of the original's own
// ylogx quirk: the base of the logarithm comes from lastX (whatever
// was in the x register before the previous operation), not from the
// current x register, even though the current x register is checked
// for positivity. See opLogYX's doc comment.
func TestLogYX_UsesLastXAsBase(t *testing.T) {
	s := newTestState()
	s.enterValue(8) // this becomes lastX after the next entry
	s.enterValue(2) // x register (only checked for sign, not used as the base)
	// After these two enters, lastX = 8 (the x register's value
	// before the second enterValue call).
	s.stack[1] = 100 // y register: the actual log argument
	if err := opLogYX(s); err != nil {
		t.Fatalf("opLogYX error = %v", err)
	}
	want := math.Log(100) / math.Log(8)
	if math.Abs(s.stack[0]-want) > 1e-9 {
		t.Errorf("ylogx = %v, want log_8(100) = %v (base from lastX, not x register)", s.stack[0], want)
	}
}

func TestBitwise_KnownValues(t *testing.T) {
	s := newTestState()
	s.enterValue(12) // 1100
	s.enterValue(10) // 1010
	if err := opAnd(s); err != nil {
		t.Fatalf("opAnd error = %v", err)
	}
	if s.stack[0] != 8 { // 1000
		t.Errorf("12 and 10 = %v, want 8", s.stack[0])
	}

	s = newTestState()
	s.enterValue(12)
	s.enterValue(10)
	if err := opOr(s); err != nil {
		t.Fatalf("opOr error = %v", err)
	}
	if s.stack[0] != 14 { // 1110
		t.Errorf("12 or 10 = %v, want 14", s.stack[0])
	}

	s = newTestState()
	s.enterValue(12)
	s.enterValue(10)
	if err := opXor(s); err != nil {
		t.Fatalf("opXor error = %v", err)
	}
	if s.stack[0] != 6 { // 0110
		t.Errorf("12 xor 10 = %v, want 6", s.stack[0])
	}
}

func TestOnesComplement_KnownValue(t *testing.T) {
	s := newTestState()
	s.enterValue(0)
	if err := opOnesComplement(s); err != nil {
		t.Fatalf("opOnesComplement error = %v", err)
	}
	if s.stack[0] != -1 {
		t.Errorf("~0 = %v, want -1 (all bits set, as a signed 4-byte int)", s.stack[0])
	}
}

// TestFormatIntegerBases_NegativeTwosComplement checks the documented
// two's-complement display: -1, as a 4-byte long, is
// 0xFFFFFFFF/037777777777/all-ones.
func TestFormatIntegerBases_NegativeTwosComplement(t *testing.T) {
	hex, oct, bin, ok := formatIntegerBases(-1)
	if !ok {
		t.Fatal("formatIntegerBases(-1) ok = false, want true")
	}
	if hex != "ffffffff" {
		t.Errorf("formatIntegerBases(-1) hex = %q, want ffffffff", hex)
	}
	if oct != "37777777777" {
		t.Errorf("formatIntegerBases(-1) oct = %q, want 37777777777", oct)
	}
	if bin != strings.Repeat("1", 32) {
		t.Errorf("formatIntegerBases(-1) bin = %q, want 32 ones", bin)
	}
}

// TestWeightConvert_RoundTrip is self-verifying: converting a value
// to another weight scale and back should reproduce the original.
func TestWeightConvert_RoundTrip(t *testing.T) {
	x := 5.0
	converted := weightConvert(x, weightLb, weightOz)
	back := weightConvert(converted, weightOz, weightLb)
	if math.Abs(back-x) > 1e-9 {
		t.Errorf("weight round trip: %v -> %v -> %v", x, converted, back)
	}
	// 1 lb = 16 oz exactly, independent of the internal kg conversion factor.
	oneLbInOz := weightConvert(1, weightLb, weightOz)
	if math.Abs(oneLbInOz-16) > 1e-6 {
		t.Errorf("1 lb in oz = %v, want 16", oneLbInOz)
	}
}

func TestLengthConvert_KnownValue(t *testing.T) {
	got := lengthConvert(12, lengthIn, lengthFt)
	if math.Abs(got-1) > 1e-9 {
		t.Errorf("lengthConvert(12 in, ft) = %v, want 1", got)
	}
	got = lengthConvert(1, lengthM, lengthCm)
	if math.Abs(got-100) > 1e-9 {
		t.Errorf("lengthConvert(1 m, cm) = %v, want 100", got)
	}
}

func TestTempConvert_KnownValue(t *testing.T) {
	got, ok := tempConvert(32, tempF, tempC)
	if !ok || math.Abs(got) > 1e-9 {
		t.Errorf("tempConvert(32F, C) = %v, ok=%v, want 0, true", got, ok)
	}
	got, ok = tempConvert(100, tempC, tempF)
	if !ok || math.Abs(got-212) > 1e-9 {
		t.Errorf("tempConvert(100C, F) = %v, ok=%v, want 212, true", got, ok)
	}
}

// TestStoreIsNotUndoable directly tests the preserved quirk: store
// operations do not call save(), so Undo reverts to whatever the
// last save()'d checkpoint was, discarding the store entirely (mem
// goes back to what it was at that earlier checkpoint, not to what
// it was just before the store).
func TestStoreIsNotUndoable(t *testing.T) {
	s := newTestState()
	s.enterValue(5)
	s.enterValue(3)
	if err := opAdd(s); err != nil { // establishes an undo checkpoint: mem=0 at this point
		t.Fatalf("opAdd error = %v", err)
	}
	if err := opStore(s); err != nil { // mem = 8, no save()
		t.Fatalf("opStore error = %v", err)
	}
	if s.mem != 8 {
		t.Fatalf("mem = %v, want 8", s.mem)
	}
	s.restoreUndo()
	if s.mem != 0 {
		t.Errorf("mem after undo = %v, want 0 (store bypassed the undo checkpoint entirely)", s.mem)
	}
}

// TestRecallIsUndoable is the counterpart: recall operations do call
// save() first, so Undo reverts them.
func TestRecallIsUndoable(t *testing.T) {
	s := newTestState()
	s.mem = 42
	s.enterValue(1)
	before := s.stack
	if err := opRecall(s); err != nil {
		t.Fatalf("opRecall error = %v", err)
	}
	if s.stack[0] != 42 {
		t.Fatalf("recall stack[0] = %v, want 42", s.stack[0])
	}
	s.restoreUndo()
	if s.stack != before {
		t.Errorf("stack after undo = %v, want %v (recall is undoable)", s.stack, before)
	}
}

func TestLastX_TracksPreOperationX(t *testing.T) {
	s := newTestState()
	s.enterValue(3)
	s.enterValue(4)
	if err := opAdd(s); err != nil {
		t.Fatalf("opAdd error = %v", err)
	}
	if s.lastX != 4 {
		t.Errorf("lastX after 3+4 = %v, want 4 (x register value just before the operation)", s.lastX)
	}
}

func TestSwapAndSwapX(t *testing.T) {
	s := newTestState()
	s.enterValue(1)
	s.enterValue(2)
	primaryBefore := s.stack

	if err := opSwap(s); err != nil {
		t.Fatalf("opSwap error = %v", err)
	}
	if s.stack != ([stackSize]float64{}) {
		t.Errorf("stack after first swap = %v, want zero (secondary calculator starts empty)", s.stack)
	}

	if err := opSwap(s); err != nil {
		t.Fatalf("opSwap error = %v", err)
	}
	if s.stack != primaryBefore {
		t.Errorf("stack after second swap = %v, want %v (swap back)", s.stack, primaryBefore)
	}
}

func TestParseNumericEntry(t *testing.T) {
	if v, ok := parseNumericEntry("3.5", 10); !ok || v != 3.5 {
		t.Errorf("parseNumericEntry(3.5, base 10) = %v, %v, want 3.5, true", v, ok)
	}
	if v, ok := parseNumericEntry("ff", 16); !ok || v != 255 {
		t.Errorf("parseNumericEntry(ff, base 16) = %v, %v, want 255, true", v, ok)
	}
	if v, ok := parseNumericEntry("101", 2); !ok || v != 5 {
		t.Errorf("parseNumericEntry(101, base 2) = %v, %v, want 5, true", v, ok)
	}
	if _, ok := parseNumericEntry("xyz", 10); ok {
		t.Error("parseNumericEntry(xyz, base 10) ok = true, want false")
	}
}
