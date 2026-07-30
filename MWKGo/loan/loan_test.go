package main

import (
	"math"
	"testing"
)

func TestComputeLoanSchedule_DocumentedDefaultInput(t *testing.T) {
	got := computeLoanSchedule(10000, 10, 6.0)
	checkClose(t, "Payment", got.Payment, 1358.6795822038373)
	checkClose(t, "TotalCost", got.TotalCost, 13586.79582203837)
	if len(got.Periods) != 10 {
		t.Fatalf("len(Periods) = %d, want 10", len(got.Periods))
	}
}

func TestComputeLoanSchedule_PrincipalPaymentsSumToLoanAmount(t *testing.T) {
	// Regardless of the amortization formula's own internal
	// derivation, a fully amortizing loan must have its per-period
	// principal payments sum to exactly the original loan amount:
	// an accounting identity independent of the formula.
	loanAmount := 10000.0
	schedule := computeLoanSchedule(loanAmount, 10, 6.0)

	var totalPrincipal float64
	for _, p := range schedule.Periods {
		totalPrincipal += p.PeriodPrincipal
	}
	checkClose(t, "totalPrincipal", totalPrincipal, loanAmount)
}

func TestComputeLoanSchedule_InterestPlusPrincipalEqualsTotalCost(t *testing.T) {
	// Total cost is, by definition, principal plus all interest
	// paid: another accounting identity independent of the formula.
	loanAmount := 5000.0
	schedule := computeLoanSchedule(loanAmount, 24, 1.5)

	var totalInterest float64
	for _, p := range schedule.Periods {
		totalInterest += p.PeriodInterest
	}
	checkClose(t, "loanAmount+totalInterest", loanAmount+totalInterest, schedule.TotalCost)
}

func TestComputeLoanSchedule_LastPeriodCumulativePrincipalEqualsLoanAmount(t *testing.T) {
	// The loan must be fully repaid by the final period: cumulative
	// principal paid should equal the original loan amount.
	loanAmount := 7500.0
	schedule := computeLoanSchedule(loanAmount, 12, 4.0)
	last := schedule.Periods[len(schedule.Periods)-1]
	checkClose(t, "last.CumulativePrincipal", last.CumulativePrincipal, loanAmount)
}

func TestFormatCurrency(t *testing.T) {
	cases := []struct {
		x    float64
		want string
	}{
		{0, "0.00"},
		{1358.6795822038373, "1,358.68"},
		{13586.79582203837, "13,586.80"},
		{999.995, "1,000.00"},
	}
	for _, c := range cases {
		if got := formatCurrency(c.x); got != c.want {
			t.Errorf("formatCurrency(%v) = %q, want %q", c.x, got, c.want)
		}
	}
}

func checkClose(t *testing.T, name string, got, want float64) {
	t.Helper()
	if math.Abs(got-want) > 1e-6 {
		t.Errorf("%s = %v, want %v", name, got, want)
	}
}
