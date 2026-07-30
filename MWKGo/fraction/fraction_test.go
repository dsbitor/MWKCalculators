package main

import "testing"

func TestParseFraction(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    fraction
		wantErr bool
	}{
		{name: "mixed number", input: "3 3/4", want: fraction{Whole: 3, Num: 3, Den: 4}},
		{name: "fraction only", input: "3/4", want: fraction{Whole: 0, Num: 3, Den: 4}},
		{name: "whole number only", input: "3", want: fraction{Whole: 3, Num: 0, Den: 1}},
		{name: "negative fraction", input: "-3/4", want: fraction{Whole: 0, Num: -3, Den: 4}},
		{name: "leading and trailing whitespace is ignored", input: "  3 3/4  ", want: fraction{Whole: 3, Num: 3, Den: 4}},
		{name: "empty input is an error", input: "", wantErr: true},
		{name: "whitespace-only input is an error", input: "   ", wantErr: true},
		{name: "non-numeric text is an error", input: "abc", wantErr: true},
		{name: "missing denominator is an error", input: "3/", wantErr: true},
		{name: "missing numerator is an error", input: "/4", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseFraction(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("parseFraction(%q) error = nil, want an error", tt.input)
				}
				return
			}
			if err != nil {
				t.Fatalf("parseFraction(%q) unexpected error: %v", tt.input, err)
			}
			if got != tt.want {
				t.Errorf("parseFraction(%q) = %+v, want %+v", tt.input, got, tt.want)
			}
		})
	}
}

func TestSplitExpression(t *testing.T) {
	tests := []struct {
		name      string
		expr      string
		wantLeft  string
		wantOp    byte
		wantRight string
		wantErr   bool
	}{
		{name: "addition", expr: "3 3/4+1 1/2", wantLeft: "3 3/4", wantOp: '+', wantRight: "1 1/2"},
		{name: "uppercase operator is lowercased", expr: "12 G 18", wantLeft: "12 ", wantOp: 'g', wantRight: " 18"},
		{name: "no operator is an error", expr: "34", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			left, op, right, err := splitExpression(tt.expr)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("splitExpression(%q) error = nil, want an error", tt.expr)
				}
				return
			}
			if err != nil {
				t.Fatalf("splitExpression(%q) unexpected error: %v", tt.expr, err)
			}
			if left != tt.wantLeft || op != tt.wantOp || right != tt.wantRight {
				t.Errorf("splitExpression(%q) = (%q, %q, %q), want (%q, %q, %q)",
					tt.expr, left, string(op), right, tt.wantLeft, string(tt.wantOp), tt.wantRight)
			}
		})
	}
}

func TestGCD(t *testing.T) {
	tests := []struct {
		name string
		x, y int64
		want int64
	}{
		{name: "normal case", x: 12, y: 18, want: 6},
		{name: "coprime values", x: 7, y: 13, want: 1},
		{name: "one operand zero returns the other", x: 0, y: 5, want: 5},
		{name: "operands reversed still returns the other", x: 5, y: 0, want: 5},
		{name: "both operands zero is defined as 1", x: 0, y: 0, want: 1},
		{name: "negative operand still returns a positive gcd", x: -4, y: 6, want: 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := gcd(tt.x, tt.y); got != tt.want {
				t.Errorf("gcd(%d, %d) = %d, want %d", tt.x, tt.y, got, tt.want)
			}
		})
	}
}

func TestEvaluate(t *testing.T) {
	tests := []struct {
		name    string
		expr    string
		want    string
		wantErr bool
	}{
		// The prompt's own example expression.
		{name: "documented example expression", expr: "3 3/4 + 1 1/2", want: "5 1/4 = 5.25"},
		{name: "addition to a whole number", expr: "1/2 + 1/2", want: "1"},
		{name: "subtraction producing a negative fraction", expr: "1/2 - 3/4", want: "-1/4 = -0.25"},
		{name: "multiplication", expr: "1/2 * 2/3", want: "1/3 = 0.3333333333333333"},
		{name: "division", expr: "1/2 \\ 1/4", want: "2"},
		{name: "gcd operates on the whole parts only", expr: "12 g 18", want: "6"},
		{name: "lcm operates on the whole parts only", expr: "4 l 6", want: "12"},
		{name: "no operator is an error", expr: "34", wantErr: true},
		{name: "a literal zero denominator is an error, not a crash", expr: "1/0 + 1/2", wantErr: true},
		{name: "dividing by a fraction that reduces to zero is an error, not a crash", expr: "1/2 \\ 0/5", wantErr: true},
		{name: "gcd of two zero whole parts is defined, not a crash", expr: "1/2 g 3/4", want: "1"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := evaluate(tt.expr)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("evaluate(%q) error = nil, want an error", tt.expr)
				}
				return
			}
			if err != nil {
				t.Fatalf("evaluate(%q) unexpected error: %v", tt.expr, err)
			}
			if got != tt.want {
				t.Errorf("evaluate(%q) = %q, want %q", tt.expr, got, tt.want)
			}
		})
	}
}
