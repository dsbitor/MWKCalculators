// loan computes the fixed periodic payment for an amortizing loan
// and a period-by-period breakdown of principal and interest, given
// the loan amount, number of payment periods, and interest rate per
// period.
//
// Converted from LOAN.C, Misc/loan. The original writes its results
// to LOAN.OUT via fopen and then pages through it with the DOS MORE
// command; this conversion prints to stdout instead and drops that
// file-save-then-page convenience, the same approach used for
// boltcirc and rotary (Tier 1 suitability review, Finding 5).
package main

import (
	"fmt"
	"math"
	"os"

	"mwkgo/internal/mwkfmt"
	"mwkgo/internal/promptio"
)

// loanPeriod is one payment period's principal and interest, both
// for that period alone and cumulatively from the start of the loan.
type loanPeriod struct {
	Period              int
	PeriodPrincipal     float64
	PeriodInterest      float64
	CumulativePrincipal float64
	CumulativeInterest  float64
}

// loanSchedule is the fixed periodic payment and total cost of an
// amortizing loan, plus its full period-by-period breakdown.
type loanSchedule struct {
	Payment   float64
	TotalCost float64
	Periods   []loanPeriod
}

// computeLoanSchedule returns the amortization schedule for a loan
// of principal loanAmount, repaid over periods payment periods at
// ratePerPeriod percent interest per period.
func computeLoanSchedule(loanAmount float64, periods int, ratePerPeriod float64) loanSchedule {
	r := 0.01 * ratePerPeriod
	growth := math.Pow(1+r, float64(periods))
	payment := loanAmount * r * growth / (growth - 1)
	totalCost := payment/r - (payment/r-loanAmount)*growth + float64(periods)*payment

	rows := make([]loanPeriod, periods)
	for i := 1; i <= periods; i++ {
		balanceAfter := payment/r - (payment/r-loanAmount)*math.Pow(1+r, float64(i))
		balanceBefore := payment/r - (payment/r-loanAmount)*math.Pow(1+r, float64(i-1))
		rows[i-1] = loanPeriod{
			Period:              i,
			PeriodPrincipal:     payment - balanceBefore*r,
			PeriodInterest:      balanceBefore * r,
			CumulativePrincipal: loanAmount - balanceAfter,
			CumulativeInterest:  float64(i)*payment - loanAmount + balanceAfter,
		}
	}

	return loanSchedule{Payment: payment, TotalCost: totalCost, Periods: rows}
}

// formatCurrency renders x with two decimal places and a comma
// grouping the whole-dollar part every three digits, matching the
// original program's own c() function. x must be non-negative, true
// of every value this program computes.
func formatCurrency(x float64) string {
	cents := int64(math.Round(x * 100))
	whole := uint64(cents / 100)
	fraction := cents % 100
	return fmt.Sprintf("%s.%02d", mwkfmt.GroupedUint(whole), fraction)
}

func main() {
	prompter, err := promptio.New(os.Stdin, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, "loan:", err)
		os.Exit(1)
	}

	fmt.Println("LOAN CALCULATIONS")
	fmt.Println()

	loanAmount := prompter.Float("Loan amount ($)", 10000.0)
	periods := prompter.Int("Number of payment periods", 10)
	ratePerPeriod := prompter.Float("Interest rate (%/period)", 6.0)

	schedule := computeLoanSchedule(loanAmount, periods, ratePerPeriod)

	fmt.Println("\nLOAN CALCULATIONS")
	fmt.Printf("Principal = $%.2f\n", loanAmount)
	fmt.Printf("Periods = %d\n", periods)
	fmt.Printf("Interest = %.2f %%\n\n", ratePerPeriod)
	fmt.Printf("Payment each period = $%s\n", formatCurrency(schedule.Payment))
	fmt.Printf("Total loan cost = $%s\n\n", formatCurrency(schedule.TotalCost))
	fmt.Println("PERIOD       PRINCIPAL       INTEREST         PTOTAL         ITOTAL")
	fmt.Println()
	for _, p := range schedule.Periods {
		fmt.Printf("%-7d%15s%15s%15s%15s\n",
			p.Period,
			formatCurrency(p.PeriodPrincipal),
			formatCurrency(p.PeriodInterest),
			formatCurrency(p.CumulativePrincipal),
			formatCurrency(p.CumulativeInterest))
	}
}
