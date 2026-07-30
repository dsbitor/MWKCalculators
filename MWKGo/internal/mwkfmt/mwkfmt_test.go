package mwkfmt

import "testing"

func TestGroupedUint(t *testing.T) {
	tests := []struct {
		name string
		n    uint64
		want string
	}{
		{name: "zero has no comma", n: 0, want: "0"},
		{name: "single digit has no comma", n: 7, want: "7"},
		{name: "value just below the first group boundary", n: 999, want: "999"},
		{name: "value at the first group boundary", n: 1000, want: "1,000"},
		{name: "six digit value groups twice", n: 123456, want: "123,456"},
		{name: "value with a short leading group", n: 1234567, want: "1,234,567"},
		{name: "maximum uint32 value", n: 4294967295, want: "4,294,967,295"},
		{name: "maximum uint64 value", n: 18446744073709551615, want: "18,446,744,073,709,551,615"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GroupedUint(tt.n); got != tt.want {
				t.Errorf("GroupedUint(%d) = %q, want %q", tt.n, got, tt.want)
			}
		})
	}
}
