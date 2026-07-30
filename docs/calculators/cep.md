# cep

Circular Error Probable (CEP).

**Converted from:** `CEP.C`, `MWKC/Math/cep.zip`
**Go source:** `MWKGo/cep/cep.go`

## Purpose

When two random error sources combine (a missile's crossrange and
downrange error, or a milling machine's X and Y axis error), the
combined error is characterized by CEP: the radius of a circle,
centered on the combined mean, within which a given fraction of
outcomes are expected to fall (conventionally 50%, but any
probability can be requested). Given the mean and standard deviation
of each error source and their correlation coefficient, this program
numerically finds the CEP radius for a requested probability.

## Inputs

| Prompt | Default |
|---|---|
| Mean1 | 0 |
| Sigma1 | 1 |
| Mean2 | 0 |
| Sigma2 | 1 |
| Rho (correlation coefficient) | 0 |
| CEP probability | 0.5 |

## Output

The CEP radius for the requested probability, plus, for the special
case of zero means and 50% probability, an independent analytic
cross-check value (Wilcox's rational polynomial approximation).

## Method

The two correlated Gaussian sources are diagonalized into principal
axes, then the bivariate density is integrated outward in growing
concentric rings (72 angular samples per ring, radial step 0.1% of
the smaller principal standard deviation) until the accumulated
probability reaches the target, with a final linear interpolation
between the last two radial steps. The original program's own
diagonalizing rotation angle formula divides by `(c11-c12)`; the
standard formula for this rotation instead divides by `(c11-c22)`.
This conversion reproduces the original's own formula unchanged,
following the same policy as `cone`'s formula quirk in Tier 1 group
5: preserve rather than silently "fix" an inherited formula. The
discrepancy has no effect at all when `sigma1 == sigma2` (the case
the program's own Wilcox cross-check applies to), since `c11-c22` is
zero either way regardless of which formula is used.

The original's unbounded `do...while` radial search is bounded to a
maximum step count per `coding-style.md` Rule 2, returning an error
if a requested probability doesn't converge within that bound (a
realistic probability converges in well under a thousand steps at
the program's own step size).

## Worked Example

No worked numeric example was included with the original program
(`CEP.TXT` explains the concept but includes no sample run). As an
independently verifiable check: for two equal, uncorrelated Gaussian
sources (a circular bivariate normal distribution), the exact 50%
CEP radius is a well known closed-form constant,
`sigma*sqrt(2*ln(2))` (~1.1774*sigma), confirmed against this
conversion's numerical integration in its tests; the program's own
Wilcox analytic cross-check was also confirmed to agree closely with
the numerical integration for an asymmetric, correlated case.
