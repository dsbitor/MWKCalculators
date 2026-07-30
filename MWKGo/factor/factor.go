// factor computes the prime factorization and the complete list of
// divisors of a whole number, taken either from the command line or
// from a prompt.
//
// Converted from FACTOR.C (M. W. Klotz, 4/98), Math/factor. The
// original computed its divisor list (pifact) by repeatedly
// multiplying pairs of already-found divisors until no new ones
// appeared; this conversion computes the same result, every divisor
// of n, directly by trial division up to sqrt(n), which is both
// simpler and asymptotically faster for the same output.
package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"

	"mwkgo/internal/mwkfmt"
	"mwkgo/internal/promptio"
)

// primeFactor is one prime raised to the power it appears in a
// factorization.
type primeFactor struct {
	Prime    uint64
	Exponent int
}

// primeSieve returns a slice where entry i is true exactly when i is
// prime, for every i from 0 up to and including limit.
func primeSieve(limit uint64) []bool {
	isPrime := make([]bool, limit+1)
	for i := range isPrime {
		isPrime[i] = i >= 2
	}
	for i := uint64(2); i*i <= limit; i++ {
		if !isPrime[i] {
			continue
		}
		for j := i * i; j <= limit; j += i {
			isPrime[j] = false
		}
	}
	return isPrime
}

// primeFactors returns the prime factorization of n as a list of
// primes with their exponents, ordered from smallest prime to
// largest. primeFactors(0) and primeFactors(1) both return an empty
// list, since neither number has a prime factorization.
func primeFactors(n uint64) []primeFactor {
	if n < 2 {
		return nil
	}

	var factors []primeFactor
	remaining := n

	for remaining%2 == 0 {
		factors = appendFactor(factors, 2)
		remaining /= 2
	}

	limit := uint64(math.Sqrt(float64(remaining))) + 1
	isPrime := primeSieve(limit)
	for candidate := uint64(3); candidate <= limit; candidate += 2 {
		if !isPrime[candidate] {
			continue
		}
		for remaining%candidate == 0 {
			factors = appendFactor(factors, candidate)
			remaining /= candidate
		}
	}

	if remaining > 1 {
		factors = appendFactor(factors, remaining)
	}
	return factors
}

func appendFactor(factors []primeFactor, prime uint64) []primeFactor {
	if len(factors) > 0 && factors[len(factors)-1].Prime == prime {
		factors[len(factors)-1].Exponent++
		return factors
	}
	return append(factors, primeFactor{Prime: prime, Exponent: 1})
}

// divisors returns every positive divisor of n in ascending order.
// divisors(0) returns an empty list, since 0 has no finite set of
// divisors under the usual definition.
func divisors(n uint64) []uint64 {
	if n == 0 {
		return nil
	}

	var low, high []uint64
	for d := uint64(1); d*d <= n; d++ {
		if n%d != 0 {
			continue
		}
		low = append(low, d)
		if partner := n / d; partner != d {
			high = append(high, partner)
		}
	}

	for i, j := 0, len(high)-1; i < j; i, j = i+1, j-1 {
		high[i], high[j] = high[j], high[i]
	}
	return append(low, high...)
}

// formatFactorization renders a prime factorization as primes joined
// by "x", with "^exponent" appended for any prime appearing more
// than once. formatFactorization(nil) returns the empty string.
func formatFactorization(factors []primeFactor) string {
	if len(factors) == 0 {
		return ""
	}
	parts := make([]string, len(factors))
	for i, f := range factors {
		parts[i] = mwkfmt.GroupedUint(f.Prime)
		if f.Exponent > 1 {
			parts[i] += fmt.Sprintf("^%d", f.Exponent)
		}
	}
	return strings.Join(parts, " x ")
}

func main() {
	n := readNumberToFactor()

	if n < 3 {
		fmt.Printf("%s = %s\n", mwkfmt.GroupedUint(n), mwkfmt.GroupedUint(n))
		return
	}

	fmt.Printf("%s = %s\n", mwkfmt.GroupedUint(n), formatFactorization(primeFactors(n)))

	fmt.Print("factors = ")
	for _, d := range divisors(n) {
		fmt.Printf(" | %s", mwkfmt.GroupedUint(d))
	}
	fmt.Println()
}

// readNumberToFactor takes the number to factor from the first
// command-line argument if one was given, matching the original
// program's argc/argv check, and otherwise prompts for it. A
// negative prompted value is treated as 0, since factoring a
// negative number is outside this calculator's scope and 0 is
// printed as-is rather than factored.
func readNumberToFactor() uint64 {
	if len(os.Args) > 1 {
		parsed, err := strconv.ParseUint(os.Args[1], 10, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "factor: %q is not a whole number\n", os.Args[1])
			os.Exit(1)
		}
		return parsed
	}

	prompter, err := promptio.New(os.Stdin, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, "factor:", err)
		os.Exit(1)
	}
	fmt.Printf("number to factor? max = %s\n", mwkfmt.GroupedUint(math.MaxUint32))
	value := prompter.Int("number to factor", 0)
	if value < 0 {
		return 0
	}
	return uint64(value)
}
