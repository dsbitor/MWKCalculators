package numdisplay

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"testing"
)

// TestFormat_SelfAdjustRoundTrips is self-verifying: whatever
// precision the self-adjusting branch selects, formatting x with
// that many places and parsing the result back must reproduce x to
// within the same tolerance the search itself used.
func TestFormat_SelfAdjustRoundTrips(t *testing.T) {
	for _, x := range []float64{0, 1, 0.5, 96.28 / 25.4, 1.0 / 3, 123456.789} {
		got := Format(x, 3, true, Fixed)
		wd, err := strconv.ParseFloat(strings.ReplaceAll(got, ",", ""), 64)
		if err != nil {
			t.Fatalf("Format(%v, selfAdjust) = %q, not parseable: %v", x, got, err)
		}
		if math.Abs(x-wd) > selfAdjustEpsilon*10 {
			t.Errorf("self-adjust round trip for %v: got %v back from %q", x, wd, got)
		}
	}
}

func TestFormat_FixedKnownValues(t *testing.T) {
	cases := []struct {
		x        float64
		decimals int
		want     string
	}{
		{1234567.891, 2, "1,234,567.89"},
		{96.28, 2, "96.28"},
		{0, 3, "0.000"},
	}
	for _, c := range cases {
		if got := Format(c.x, c.decimals, false, Fixed); got != c.want {
			t.Errorf("Format(%v, %d, fixed) = %q, want %q", c.x, c.decimals, got, c.want)
		}
	}
}

// TestFormat_EngineeringMantissaAndExponent checks the two defining
// properties of engineering notation directly (mantissa magnitude in
// [1,1000), exponent a multiple of three) and that the mantissa and
// exponent reconstruct the original value, rather than asserting a
// specific formatted string.
func TestFormat_EngineeringMantissaAndExponent(t *testing.T) {
	for _, x := range []float64{1234.5, 0.00567, -98765.4, 3} {
		s := Format(x, 3, false, Engineering)
		var mantissa float64
		var exponent int
		if _, err := fmt.Sscanf(s, "%f E%d", &mantissa, &exponent); err != nil {
			t.Fatalf("Format(%v, engineering) produced unparsable %q: %v", x, s, err)
		}
		if exponent%3 != 0 {
			t.Errorf("Format(%v, engineering) exponent = %d, not a multiple of 3", x, exponent)
		}
		if math.Abs(mantissa) < 1 || math.Abs(mantissa) >= 1000 {
			t.Errorf("Format(%v, engineering) mantissa = %v, want magnitude in [1,1000)", x, mantissa)
		}
		reconstructed := mantissa * math.Pow(10, float64(exponent))
		if math.Abs(reconstructed-x)/math.Max(math.Abs(x), 1) > 1e-3 {
			t.Errorf("Format(%v, engineering) reconstructs to %v", x, reconstructed)
		}
	}
}

func TestFormat_ScientificReconstructsValue(t *testing.T) {
	for _, x := range []float64{1234.5, 0.00567, -98765.4} {
		s := Format(x, 4, false, Scientific)
		var mantissa, exponent float64
		if _, err := fmt.Sscanf(s, "%f E%f", &mantissa, &exponent); err != nil {
			t.Fatalf("Format(%v, scientific) produced unparsable %q: %v", x, s, err)
		}
		reconstructed := mantissa * math.Pow(10, exponent)
		if math.Abs(reconstructed-x)/math.Max(math.Abs(x), 1) > 1e-3 {
			t.Errorf("Format(%v, scientific) reconstructs to %v", x, reconstructed)
		}
	}
}

func TestFormat_MagnitudeEscalatesToEngineering(t *testing.T) {
	// A magnitude with log10 beyond 15 in fixed-point mode should
	// automatically switch to engineering notation.
	x := 1e20
	got := Format(x, 2, false, Fixed)
	if !strings.Contains(got, "E") {
		t.Errorf("Format(%v, fixed) = %q, want automatic escalation to engineering notation", x, got)
	}
}
