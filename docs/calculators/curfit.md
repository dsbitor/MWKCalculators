# curfit

Curve fitting to a set of data pairs.

**Converted from:** `CURFIT.C` (M. W. Klotz), `MWKC/Math/curfit.zip`
**Go source:** `MWKGo/curfit/curfit.go`

## Purpose

Fits a curve to a set of (x,y) data pairs by least squares, in one of
four forms: polynomial, logarithmic, exponential, or power. A session
can try several fit types against the same data before quitting, to
compare which shape best describes it.

## Data setup

The data pairs are specific to one curve-fitting job, not universal
reference data or reusable equipment configuration — so, like
`calibrat`, `vrev`, `simul`, and `colsort` in earlier conversion
groups, this program reads its input fresh from a file named on the
command line each run, in the same `STARTOFDATA`/`ENDOFDATA` text
format the original used; see
`ai/plans/c-to-go-conversion-plan.md`, "Data-file strategy for
Tier 2".

```
curfit -data my-points.dat
```

Each data line is `x,y`. A worked example built from the original
archive's own shipped `CURFIT.DAT` (its first, noise-free data block)
ships at `MWKGo/curfit/testdata/example.dat`.

## Inputs

A repeating menu choice (`1`, `2`, `3`, `4`, or `Q` to quit); for a
polynomial fit, also the degree to fit (1 up to `min(points-1, 49)`,
`CURFIT.C`'s own fixed matrix-size bound).

| Fit | Form |
|---|---|
| 1 Polynomial | Y = A0 + A1·X + A2·X² + A3·X³ + ... |
| 2 Logarithmic | Y = A + B·ln(X) |
| 3 Exponential | Y = A·exp(B·X) |
| 4 Power | Y = A·Xᴮ |

## Output

The fit's coefficients, then a per-point table (index, x, y, the
fit's own calculated y, the error, and the error as a percentage of
y), with `**` marking every point whose error percentage is a new
running maximum as the table is printed (not only the single worst
point overall — a point that ties or falls short of an
already-established maximum is not marked, even if it isn't the best
point either). Finally, a degrees-of-freedom-adjusted correlation
coefficient: closer to 1 means a better fit.

## Method

The polynomial fit solves the least-squares normal equations directly
(the standard sum-of-powers matrix construction), via the same
full-pivoting Gauss-Jordan solver `simul` already provides (see
`docs/calculators/simul.md`) — `CURFIT.C`'s own `gjor()` implements
the identical textbook method. The other three fits are all linear
regression in disguise: logarithmic regresses Y against ln(X);
exponential and power both regress ln(Y) against X or ln(X)
respectively, then exponentiate the resulting intercept back out.

The correlation coefficient's `sst` term (total sum of squares) is
computed from the *raw* Y values for the polynomial and logarithmic
fits, but from the *log-transformed* Y values for the exponential and
power fits — matching `CURFIT.C`'s own bookkeeping exactly, where the
same `y1`/`y2` accumulator variables track different quantities
depending on which fit is running.

## Worked Example

`CURFIT.DAT`'s own first data block (20 points, no added noise) is
this conversion's primary test: its comment claims
`y = 4 + 3*x + 2*x^2 + 1*x^2`, which is really a typo for
`1*x^3` (two terms both reading `x^2` would just combine into `3*x^2`,
not describe a fourth, distinct term) — confirmed by taking finite
differences of the data itself, which are cubic, not quadratic. A
degree-3 polynomial fit against this data recovers
`A0=4, A1=3, A2=2, A3=1` to within floating-point roundoff, and its
reported correlation coefficient is 1.0. The logarithmic, exponential,
and power fits are each checked against small, independently
constructed noise-free datasets built directly from their own
formulas, confirming exact coefficient recovery; the "new running
maximum" `**` marking rule and the zero-y "error% is 0, not a division
error" guard are each checked directly against a small synthetic
dataset built to exercise that specific rule.
