// fraction is a rational fraction calculator: it reads an
// expression of two mixed numbers and an operator, such as
// "3 3/4 + 1 1/2", and prints the reduced result.
//
// Converted from FRACTION.C (M. W. Klotz, 5/98, 6/03),
// Math/fraction.
package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// fraction is a mixed number: Whole + Num/Den. Den is always
// nonzero for a fraction produced by parseFraction or combine.
type fraction struct {
	Whole, Num, Den int64
}

// parseFraction parses a mixed-number expression such as "3 3/4",
// "3/4", or "3" into a fraction. Leading and trailing whitespace is
// ignored.
func parseFraction(s string) (fraction, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return fraction{}, errors.New("empty fraction expression")
	}

	slashIndex := strings.IndexByte(s, '/')
	if slashIndex == -1 {
		whole, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return fraction{}, fmt.Errorf("parse whole number %q: %w", s, err)
		}
		return fraction{Whole: whole, Den: 1}, nil
	}

	wholeText, fractionText := "", s
	if spaceIndex := strings.IndexByte(s, ' '); spaceIndex != -1 && spaceIndex < slashIndex {
		wholeText, fractionText = s[:spaceIndex], strings.TrimSpace(s[spaceIndex+1:])
	}

	var whole int64
	if wholeText != "" {
		var err error
		if whole, err = strconv.ParseInt(wholeText, 10, 64); err != nil {
			return fraction{}, fmt.Errorf("parse whole part %q: %w", wholeText, err)
		}
	}

	numDen := strings.SplitN(fractionText, "/", 2)
	if len(numDen) != 2 {
		return fraction{}, fmt.Errorf("parse fraction part %q: missing '/'", fractionText)
	}
	num, err := strconv.ParseInt(strings.TrimSpace(numDen[0]), 10, 64)
	if err != nil {
		return fraction{}, fmt.Errorf("parse numerator %q: %w", numDen[0], err)
	}
	den, err := strconv.ParseInt(strings.TrimSpace(numDen[1]), 10, 64)
	if err != nil {
		return fraction{}, fmt.Errorf("parse denominator %q: %w", numDen[1], err)
	}

	return fraction{Whole: whole, Num: num, Den: den}, nil
}

// splitExpression finds the first supported operator character
// (+, -, *, \, g, l, case-insensitive) in expr and splits it into
// the left operand, the operator itself (lowercased), and the right
// operand.
func splitExpression(expr string) (left string, op byte, right string, err error) {
	opIndex := strings.IndexAny(expr, "+-*\\glGL")
	if opIndex == -1 {
		return "", 0, "", errors.New("operator not found")
	}

	op = expr[opIndex]
	if op >= 'A' && op <= 'Z' {
		op += 'a' - 'A'
	}
	return expr[:opIndex], op, expr[opIndex+1:], nil
}

// gcd returns the greatest common divisor of x and y, treated as
// magnitudes: gcd(0, 0) is defined as 1 and gcd(0, n) as n, matching
// mathematical convention. The original program's Euclidean loop
// divides by zero, an undefined-behaviour crash in C, whenever one
// argument is zero and the other is not; that case is handled
// explicitly here instead of reaching the loop.
func gcd(x, y int64) int64 {
	a, b := x, y
	if a > b {
		a, b = b, a
	}
	if a == 0 {
		if b == 0 {
			return 1
		}
		return b
	}
	for {
		c := b % a
		if c == 0 {
			return a
		}
		b, a = a, c
	}
}

// combine applies operator op to fractions a and b: +, -, *, and \
// (divide) act on the full mixed numbers; g and l act only on the
// whole-number parts, returning their greatest common divisor or
// least common multiple respectively, matching the original
// program's scope for those two operators.
func combine(a, b fraction, op byte) (fraction, error) {
	improperA := a.Den*a.Whole + a.Num
	improperB := b.Den*b.Whole + b.Num

	var num, den int64
	switch op {
	case '+':
		num, den = b.Den*improperA+a.Den*improperB, a.Den*b.Den
	case '-':
		num, den = b.Den*improperA-a.Den*improperB, a.Den*b.Den
	case '*':
		num, den = improperA*improperB, a.Den*b.Den
	case '\\':
		num, den = improperA*b.Den, improperB*a.Den
	case 'g':
		return fraction{Whole: gcd(a.Whole, b.Whole), Den: 1}, nil
	case 'l':
		divisor := gcd(a.Whole, b.Whole)
		return fraction{Whole: a.Whole * b.Whole / divisor, Den: 1}, nil
	default:
		return fraction{}, fmt.Errorf("unsupported operator %q", op)
	}

	// A zero denominator here means either a literal zero
	// denominator was typed (such as "1/0") or, for \, the right
	// operand evaluated to zero. The original program has no guard
	// against this and crashes with a division-by-zero fault when
	// reduce's integer division is reached; returning an error here
	// instead keeps that a normal, reported failure rather than a
	// crash.
	if den == 0 {
		return fraction{}, errors.New("division by zero")
	}

	return reduce(num, den), nil
}

// reduce converts the improper fraction num/den into a mixed number
// in lowest terms.
func reduce(num, den int64) fraction {
	if den == 1 {
		return fraction{Whole: num, Den: 1}
	}
	if num == 0 {
		return fraction{Den: 1}
	}

	sign := int64(1)
	if num < 0 {
		sign = -1
		num = -num
	}

	whole := (num / den) * sign
	num -= (num / den) * den
	if num == 0 {
		return fraction{Whole: whole, Den: 1}
	}

	if divisor := gcd(num, den); divisor > 1 {
		num, den = num/divisor, den/divisor
	}
	return fraction{Whole: whole, Num: num * sign, Den: den}
}

// formatResult renders f the way the original program does: the
// whole part is shown unless it is zero with a nonzero fractional
// part to show instead, and the fractional part, if any, is shown
// as numerator/denominator alongside its decimal equivalent.
func formatResult(f fraction) string {
	var b strings.Builder
	if f.Whole != 0 || f.Num == 0 {
		fmt.Fprintf(&b, "%d ", f.Whole)
	}
	if f.Num != 0 {
		decimal := float64(f.Whole) + float64(f.Num)/float64(f.Den)
		fmt.Fprintf(&b, "%d/%d = %g", f.Num, f.Den, decimal)
	}
	return strings.TrimSpace(b.String())
}

// evaluate parses expr as "left OP right" and returns the formatted
// result.
func evaluate(expr string) (string, error) {
	left, op, right, err := splitExpression(expr)
	if err != nil {
		return "", err
	}
	a, err := parseFraction(left)
	if err != nil {
		return "", err
	}
	b, err := parseFraction(right)
	if err != nil {
		return "", err
	}
	result, err := combine(a, b, op)
	if err != nil {
		return "", err
	}
	return formatResult(result), nil
}

func main() {
	if args := os.Args[1:]; len(args) > 0 {
		printResult(strings.Join(args, " "))
		return
	}

	fmt.Println(`expression syntax: a b/c (+,-,*,\,g(cd),l(cm)) d e/f`)
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("expression (e.g. 3 3/4 + 1 1/2) [Enter to quit] ? ")
		if !scanner.Scan() {
			return
		}
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			return
		}
		printResult(line)
	}
}

func printResult(expr string) {
	result, err := evaluate(expr)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%s = %s\n", expr, result)
}
