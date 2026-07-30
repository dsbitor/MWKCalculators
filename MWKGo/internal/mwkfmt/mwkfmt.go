// Package mwkfmt implements small output-formatting helpers shared
// across the converted calculators that the Go standard library does
// not provide directly, starting with the comma-grouped integer
// formatting used by FACTOR.C's cprnt function.
package mwkfmt

import "strconv"

// GroupedUint formats n with a comma inserted every three digits from
// the right, matching the original cprnt function's output (for
// example 1234567 becomes "1,234,567"). Values below 1000 are
// returned unchanged, with no comma.
func GroupedUint(n uint64) string {
	digits := strconv.FormatUint(n, 10)

	groupCount := (len(digits) - 1) / 3
	if groupCount == 0 {
		return digits
	}

	grouped := make([]byte, 0, len(digits)+groupCount)
	firstGroupLen := len(digits) % 3
	if firstGroupLen == 0 {
		firstGroupLen = 3
	}

	grouped = append(grouped, digits[:firstGroupLen]...)
	for i := firstGroupLen; i < len(digits); i += 3 {
		grouped = append(grouped, ',')
		grouped = append(grouped, digits[i:i+3]...)
	}
	return string(grouped)
}
