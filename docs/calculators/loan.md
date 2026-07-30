# loan

Loan amortization schedule.

**Converted from:** `LOAN.C`, `MWKC/Misc/loan.zip`
**Go source:** `MWKGo/loan/loan.go`

## Purpose

Given a loan amount, number of payment periods, and interest rate per
period, computes the fixed periodic payment and total cost of the
loan, plus a period-by-period breakdown of principal and interest,
both for that period alone and cumulatively.

## Inputs

| Prompt | Default |
|---|---|
| Loan amount | $10,000 |
| Number of payment periods | 10 |
| Interest rate | 6 %/period |

## Output

Payment amount per period, total loan cost, then a table with each
period's principal payment, interest payment, cumulative principal
paid, and cumulative interest paid.

## Method

Standard amortization formulas:

```
r = rate/100
growth = (1+r)^periods
payment = loanAmount*r*growth / (growth - 1)
totalCost = payment/r - (payment/r - loanAmount)*growth + periods*payment

for period i = 1..periods:
  balanceBefore = payment/r - (payment/r - loanAmount)*(1+r)^(i-1)
  balanceAfter  = payment/r - (payment/r - loanAmount)*(1+r)^i
  periodInterest = balanceBefore * r
  periodPrincipal = payment - periodInterest
  cumulativePrincipal = loanAmount - balanceAfter
  cumulativeInterest = i*payment - loanAmount + balanceAfter
```

The original program writes its results to `LOAN.OUT` via `fopen`
and pages through it with the DOS `MORE` command; this conversion
prints to stdout instead and drops that file-save-then-page
convenience, the same approach used for `boltcirc` and `rotary`
(Tier 1 suitability review, Finding 5). The original's `c()`
function, which comma-groups a formatted dollar amount, is
reimplemented as `formatCurrency`, built on the existing
`mwkfmt.GroupedUint` helper rather than duplicating comma-insertion
logic.

## Worked Example

No worked numeric example was included with the original program. As
independently verifiable checks: the per-period principal payments
must sum to exactly the original loan amount, and cumulative
principal after the final period must equal the loan amount, since a
fully amortizing loan is by definition paid off exactly on schedule;
both confirmed in this conversion's tests.
