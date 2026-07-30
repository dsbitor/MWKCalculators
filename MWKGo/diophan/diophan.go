// diophan solves a linear Diophantine equation, ax + by = c, for
// integer x and y, using the extended Euclidean algorithm, then
// reports the general integer solution family and several sample
// solutions.
//
// Converted from DIOPHAN.C, Math/diophan.
package main

import (
	"fmt"
	"os"

	"mwkgo/internal/promptio"
)

// maxContinuedFractionSteps bounds the extended Euclidean algorithm's
// internal expansion, matching the original program's own NQR
// constant (there raised to guard against a pathological input; a
// realistic a and b converge in a handful of steps).
const maxContinuedFractionSteps = 100

// gcd returns the greatest common divisor of x and y (as positive
// magnitudes); gcd(0, 0) is defined as 1, matching the original
// program's own convention.
func gcd(x, y int64) int64 {
	a, b := x, y
	if a < 0 {
		a = -a
	}
	if b < 0 {
		b = -b
	}
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

// solveDiophantine finds integers x and y satisfying a*x + b*y = c,
// via the extended Euclidean algorithm, a direct port of the
// original program's own diophan function (including its 1-based
// working-array indexing). It returns ok = false if c is not evenly
// divisible by gcd(a, b), in which case no integer solution exists,
// or if the expansion doesn't terminate within
// maxContinuedFractionSteps.
func solveDiophantine(a, b, c int64) (x, y int64, ok bool) {
	g := gcd(a, b)
	if c%g != 0 {
		return 0, 0, false
	}

	r := make([]int64, maxContinuedFractionSteps+3)
	q := make([]int64, maxContinuedFractionSteps+3)
	s := make([]int64, maxContinuedFractionSteps+3)

	i := 0
	r[1] = a / g
	r[2] = b / g
	k := int64(-1)
	for {
		i++
		if i > maxContinuedFractionSteps {
			return 0, 0, false
		}
		k *= -1
		q[i] = r[i] / r[i+1]
		r[i+2] = r[i] - q[i]*r[i+1]
		if r[i+2] == 0 {
			break
		}
	}

	s[i] = 1
	s[i+1] = 0
	i--
	for {
		s[i] = q[i]*s[i+1] + s[i+2]
		i--
		k *= -1
		if i <= 0 {
			break
		}
	}

	x = k * s[2] * c / g
	y = -k * s[1] * c / g
	return x, y, true
}

func main() {
	prompter, err := promptio.New(os.Stdin, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, "diophan:", err)
		os.Exit(1)
	}

	fmt.Println("SOLUTION OF LINEAR DIOPHANTINE EQUATIONS")
	fmt.Println()
	fmt.Println("Find integer x and y to satisfy the equation: ax + by = c")
	fmt.Println("where a,b, and c are integers.")
	fmt.Println("c must be evenly divisible by the gcd of a and b for a solution.")
	fmt.Println()

	a := int64(prompter.Float("Integer value of a", 172))
	b := int64(prompter.Float("Integer value of b", 20))
	c := int64(prompter.Float("Integer value of c", 1000))

	fmt.Println()
	x, y, ok := solveDiophantine(a, b, c)
	if !ok {
		g := gcd(a, b)
		fmt.Printf("gcd(%d,%d) = %d\n", a, b, g)
		fmt.Printf("c = %d not evenly divisible by gcd\n", c)
		fmt.Println("EQUATION NOT SOLVABLE IN INTEGERS")
		fmt.Println("NO SOLUTION WAS FOUND")
		return
	}

	g := gcd(a, b)
	fmt.Printf("x = %d + k*(%d), y = %d - k*(%d)  (k integer)\n", x, b/g, y, a/g)

	fmt.Println("\nSome sample solutions:")
	for k := int64(-4); k < 5; k++ {
		xp := x + k*b/g
		yp := y - k*a/g
		z := a*xp + b*yp
		fmt.Printf("k=%+2d: (%d)*(%d) + (%d)*(%d) = %d\n", k, a, xp, b, yp, z)
	}
}
