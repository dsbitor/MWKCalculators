# factor

Prime factorization and complete divisor list for a whole
number.

**Converted from:** `FACTOR.C` (M. W. Klotz, 4/98),
`MWKC/Math/factor.zip`
**Go source:** `MWKGo/factor/factor.go`

## Purpose

Factors a whole number into primes and lists every divisor of
that number. The original documents a practical scope of
numbers under 2^32-1 (4,294,967,295); the Go conversion accepts
any value that fits a 64-bit unsigned integer.

## Inputs

The number to factor, given as a command-line argument
(`factor 172`) or, if none is given, at a prompt with no
default.

## Output

The prime factorization, primes separated by `x` with `^n`
shown for any prime appearing more than once, followed by the
complete list of divisors in ascending order.

## Method

Extracts factors of 2 first, then trial-divides by odd numbers
up to the square root of what remains, using a sieve of
Eratosthenes to skip non-primes. The complete divisor list is
generated separately by trial division up to the square root of
the input, pairing each divisor found with its complement.

The original computed its divisor list by repeatedly
multiplying pairs of already-found divisors until no new ones
appeared, a correct but far slower way to reach the same
result.

## Worked Example

`factor 4294967295` (2^32-1) gives `3 x 5 x 17 x 257 x 65,537`,
the product of the first five Fermat primes, a well known
factorization independent of this program, and exactly 32
divisors (2^5, since each of the five primes appears once).
