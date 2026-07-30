package main

import (
	"reflect"
	"testing"
)

func TestPrimeFactors(t *testing.T) {
	tests := []struct {
		name string
		n    uint64
		want []primeFactor
	}{
		{name: "zero has no prime factorization", n: 0, want: nil},
		{name: "one has no prime factorization", n: 1, want: nil},
		{name: "smallest prime", n: 2, want: []primeFactor{{Prime: 2, Exponent: 1}}},
		{name: "a prime with no small factors", n: 17, want: []primeFactor{{Prime: 17, Exponent: 1}}},
		{name: "repeated single prime factor", n: 12, want: []primeFactor{{Prime: 2, Exponent: 2}, {Prime: 3, Exponent: 1}}},
		{name: "two squared prime factors", n: 100, want: []primeFactor{{Prime: 2, Exponent: 2}, {Prime: 5, Exponent: 2}}},
		{
			// 2^32-1 is the product of the first five Fermat primes,
			// a well known factorization independent of this code.
			name: "2^32-1, the product of the first five Fermat primes",
			n:    4294967295,
			want: []primeFactor{
				{Prime: 3, Exponent: 1},
				{Prime: 5, Exponent: 1},
				{Prime: 17, Exponent: 1},
				{Prime: 257, Exponent: 1},
				{Prime: 65537, Exponent: 1},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := primeFactors(tt.n)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("primeFactors(%d) = %v, want %v", tt.n, got, tt.want)
			}
		})
	}
}

func TestDivisors(t *testing.T) {
	tests := []struct {
		name string
		n    uint64
		want []uint64
	}{
		{name: "zero has no divisors", n: 0, want: nil},
		{name: "one divides only itself", n: 1, want: []uint64{1}},
		{name: "a prime has exactly two divisors", n: 17, want: []uint64{1, 17}},
		{name: "twelve has six divisors", n: 12, want: []uint64{1, 2, 3, 4, 6, 12}},
		{name: "a perfect square includes its square root once", n: 16, want: []uint64{1, 2, 4, 8, 16}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := divisors(tt.n)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("divisors(%d) = %v, want %v", tt.n, got, tt.want)
			}
		})
	}
}

func TestFormatFactorization(t *testing.T) {
	tests := []struct {
		name    string
		factors []primeFactor
		want    string
	}{
		{name: "no factors formats as empty", factors: nil, want: ""},
		{name: "single prime with no exponent shown", factors: []primeFactor{{Prime: 17, Exponent: 1}}, want: "17"},
		{
			name: "360 as 2^3 x 3^2 x 5",
			factors: []primeFactor{
				{Prime: 2, Exponent: 3},
				{Prime: 3, Exponent: 2},
				{Prime: 5, Exponent: 1},
			},
			want: "2^3 x 3^2 x 5",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := formatFactorization(tt.factors); got != tt.want {
				t.Errorf("formatFactorization(%v) = %q, want %q", tt.factors, got, tt.want)
			}
		})
	}
}
