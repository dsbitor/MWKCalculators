# fraction

Rational fraction calculator.

**Converted from:** `FRACTION.C` (M. W. Klotz, 5/98, 6/03),
`MWKC/Math/fraction.zip`
**Go source:** `MWKGo/fraction/fraction.go`

## Purpose

Evaluates an expression combining two mixed numbers with an
operator, such as `3 3/4 + 1 1/2`, and prints the reduced
result as a mixed number and its decimal equivalent. Also
supports `g` and `l`, which compute the greatest common divisor
or least common multiple of the whole-number parts of the two
operands, for expressions like `12 g 18`.

## Inputs

An expression, given as a command-line argument or, if none is
given, at a repeated prompt (blank input exits). Supported
operators: `+`, `-`, `*`, `\` (divide), `g` (gcd), `l` (lcm),
case-insensitive. Each operand is a mixed number: a whole part
and a `numerator/denominator` part separated by a space, either
part alone, or just a whole number.

## Output

The expression followed by its reduced result: a whole part
(when nonzero, or when there is no fractional part to show
instead), and a `numerator/denominator = decimal` part (when the
result has a fractional part).

## Method

Both operands are converted to a common denominator (or, for
divide, cross-multiplied) and the result is reduced to lowest
terms by extracting the whole part and dividing the remaining
numerator and denominator by their greatest common divisor.

Converting this program surfaced two crash bugs in the original,
both fixed in the conversion and documented at the point of the
fix: its `gcd` function divides by zero whenever one operand is
zero and the other isn't, which happens whenever `g` or `l` is
used with an operand that has no whole part (such as `1/2 g
3/4`); and dividing by an operand that evaluates to zero (such as
`1/2 \ 0/5`) produces a zero denominator that the original never
checks for before dividing by it.

## Worked Example

The program's own usage line offers `3 3/4 + 1 1/2` as an
example, which gives `5 1/4 = 5.25`.
