package main

import (
	"math"
	"testing"
)

// A single self-consistent mixture: 4 units of 10% solution mixed
// with 6 units of 0% solution (dilution) gives 10 units of 4%
// mixture. Used to verify each of the six solvable combinations
// recovers the same full result.
const (
	testConcA, testConcB = 10.0, 0.0
	testAmountA          = 4.0
	testAmountB          = 6.0
	testAmountMixture    = 10.0
	testConcMixture      = 4.0
)

func TestSolveMixture_AllSixSolvableCombinations(t *testing.T) {
	cases := []struct {
		name                             string
		amountA, amountB, amountM, concM float64
	}{
		{"a+b", testAmountA, testAmountB, 0, 0},
		{"a+concM", testAmountA, 0, 0, testConcMixture},
		{"a+amountM", testAmountA, 0, testAmountMixture, 0},
		{"b+concM", 0, testAmountB, 0, testConcMixture},
		{"b+amountM", 0, testAmountB, testAmountMixture, 0},
		{"amountM+concM", 0, 0, testAmountMixture, testConcMixture},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := solveMixture(testConcA, testConcB, c.amountA, c.amountB, c.amountM, c.concM)
			if err != nil {
				t.Fatalf("solveMixture() error = %v", err)
			}
			checkClose(t, "AmountA", got.AmountA, testAmountA)
			checkClose(t, "AmountB", got.AmountB, testAmountB)
			checkClose(t, "AmountMixture", got.AmountMixture, testAmountMixture)
			checkClose(t, "ConcentrationMixturePct", got.ConcentrationMixturePct, testConcMixture)
		})
	}
}

func TestSolveMixture_MassBalanceIdentity(t *testing.T) {
	// Regardless of which two quantities were given, the resulting
	// mixture must conserve total active-ingredient mass: an
	// accounting identity independent of the specific formula used
	// to reach the answer.
	got, err := solveMixture(25, 5, testAmountA, testAmountB, 0, 0)
	if err != nil {
		t.Fatalf("solveMixture() error = %v", err)
	}
	activeIngredientIn := testAmountA*0.25 + testAmountB*0.05
	activeIngredientOut := got.AmountMixture * (got.ConcentrationMixturePct / 100)
	checkClose(t, "activeIngredientOut", activeIngredientOut, activeIngredientIn)
}

func TestSolveMixture_InsufficientDataReturnsError(t *testing.T) {
	_, err := solveMixture(testConcA, testConcB, testAmountA, 0, 0, 0)
	if err == nil {
		t.Fatal("solveMixture() error = nil, want an error for a single known amount")
	}
}

func checkClose(t *testing.T, name string, got, want float64) {
	t.Helper()
	if math.Abs(got-want) > 1e-9 {
		t.Errorf("%s = %v, want %v", name, got, want)
	}
}
