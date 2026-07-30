# rfrac

Rational fraction approximation.

**Converted from:** `RFRAC.C`, `MWKC/Math/rfrac.zip`
**Go source:** `MWKGo/rfrac/rfrac.go`

## Purpose

Approximates any decimal number as a rational fraction, expanding a
continued fraction one term at a time only as far as needed to reach
a requested approximation accuracy (rather than lvern's fixed, very
tight tolerance), so a coarser accuracy request yields a simpler
fraction with a smaller denominator.

## Inputs

| Prompt | Default |
|---|---|
| Number to approximate | 3.14159 |
| Desired approximation accuracy | 0.01 % |

## Output

The number as a mixed number (whole part plus reduced
numerator/denominator), the equivalent improper fraction, its decimal
value, and the approximation's percentage error.

## Method

Standard continued-fraction convergent expansion of the number's
fractional part, evaluated after each new term until the relative
error against the original fractional part falls within the
requested accuracy (or a bound of 100 terms is reached, matching the
original program's own `NCF` constant), then reduced to lowest terms
via the greatest common divisor.

A whole-number input (no fractional part at all) has nothing to
approximate; the original program has no guard for this and divides
by zero in its own relative-error calculation. This conversion
reports it as an explicit error instead, the same policy as
`fraction`'s divide-by-zero fix in Tier 1 group 2.

## Worked Example

No worked numeric example is available. As an independently
verifiable check: the documented default input (3.14159 to 0.01%
accuracy) converges on 355/113 (Milü), the famous, extremely accurate
historical rational approximation to pi, independent of both this
program's own algorithm and any hand-picked expected value; confirmed
in this conversion's tests.
