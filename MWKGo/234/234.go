// 234 solves quadratic, cubic, and quartic equations with real
// coefficients, returning all roots (real or complex).
//
// Converted from 234.C (M. W. Klotz, 2/01), Math/234.
package main

import (
	"errors"
	"fmt"
	"math"
	"os"

	"mwkgo/internal/promptio"
)

func sgn(x float64) float64 {
	switch {
	case x > 0:
		return 1
	case x < 0:
		return -1
	default:
		return 0
	}
}

// solveQuadratic returns the two roots of c2*x^2 + c1*x + c0 = 0.
func solveQuadratic(c0, c1, c2 float64) [2]complex128 {
	a, b, c := c2, c1, c0
	d := b*b - 4*a*c
	e := math.Sqrt(math.Abs(d))

	if d >= 0 {
		return [2]complex128{
			complex((-b+e)/(2*a), 0),
			complex((-b-e)/(2*a), 0),
		}
	}
	return [2]complex128{
		complex(-b/(2*a), e/(2*a)),
		complex(-b/(2*a), -e/(2*a)),
	}
}

// solveCubic returns the three roots of c3*x^3 + c2*x^2 + c1*x + c0 = 0,
// via Cardano's method (trigonometric form for three real roots).
func solveCubic(c0, c1, c2, c3 float64) [3]complex128 {
	p := c2 / c3
	q := c1 / c3
	r := c0 / c3

	a := (3*q - p*p) / 3
	b := (2*p*p*p - 9*p*q + 27*r) / 27
	d := b*b/4 + a*a*a/27
	u := p / 3

	if d < 0 {
		// three real, unequal roots
		phi := math.Acos(clamp(-0.5*b/math.Sqrt(-a*a*a/27))) * 180 / math.Pi
		t := 2 * math.Sqrt(-a/3)
		x := phi / 3
		return [3]complex128{
			complex(t*cosd(x)-u, 0),
			complex(t*cosd(x+120)-u, 0),
			complex(t*cosd(x+240)-u, 0),
		}
	}

	// one real root and a conjugate pair, or a repeated real root
	x := math.Sqrt(d)
	t := -0.5*b + x
	aa := sgn(t) * math.Pow(math.Abs(t), 1.0/3.0)
	t = -0.5*b - x
	bb := sgn(t) * math.Pow(math.Abs(t), 1.0/3.0)

	return [3]complex128{
		complex(aa+bb-u, 0),
		complex(-0.5*(aa+bb)-u, 0.5*(aa-bb)*math.Sqrt(3)),
		complex(-0.5*(aa+bb)-u, -0.5*(aa-bb)*math.Sqrt(3)),
	}
}

func clamp(x float64) float64 {
	if x > 1 {
		return 1
	}
	if x < -1 {
		return -1
	}
	return x
}
func sind(deg float64) float64 { return math.Sin(deg * math.Pi / 180) }
func cosd(deg float64) float64 { return math.Cos(deg * math.Pi / 180) }

// errNoResultantCubicRoot reports that none of the resultant cubic's
// three roots produced a usable (real, non-negative) value for R^2 in
// Ferrari's method. The original program has the same limitation:
// per its own comment, it hasn't explored the case where R might be
// imaginary.
var errNoResultantCubicRoot = errors.New("no usable real solution to resultant cubic")

// solveQuartic returns the four roots of c4*x^4 + c3*x^3 + c2*x^2 +
// c1*x + c0 = 0, via Ferrari's method (solving a resultant cubic,
// then two quadratics).
func solveQuartic(c0, c1, c2, c3, c4 float64) ([4]complex128, error) {
	a := c3 / c4
	b := c2 / c4
	c := c1 / c4
	d := c0 / c4

	resolventCubic := solveCubic(-a*a*d+4*b*d-c*c, a*c-4*d, -b, 1)

	var y, rr float64
	found := false
	for _, root := range resolventCubic {
		if imag(root) != 0 {
			continue
		}
		y = real(root)
		rr = a*a/4 - b + y
		if rr < 0 {
			continue
		}
		found = true
		break
	}
	if !found {
		return [4]complex128{}, errNoResultantCubicRoot
	}

	r := math.Sqrt(math.Abs(rr))

	var dd, ee float64
	if r == 0 {
		dd = 0.75*a*a - 2*b + 2*math.Sqrt(y*y-4*d)
		ee = 0.75*a*a - 2*b - 2*math.Sqrt(y*y-4*d)
	} else {
		dd = 0.75*a*a - rr - 2*b + (4*a*b-8*c-a*a*a)/(4*r)
		ee = 0.75*a*a - rr - 2*b - (4*a*b-8*c-a*a*a)/(4*r)
	}
	bigD := math.Sqrt(math.Abs(dd))
	bigE := math.Sqrt(math.Abs(ee))

	var re, im [4]float64
	for i := range re {
		re[i] = -0.25 * a
	}

	if rr >= 0 {
		re[0] += 0.5 * r
		re[1] += 0.5 * r
		re[2] -= 0.5 * r
		re[3] -= 0.5 * r
	} else {
		im[0] += 0.5 * r
		im[1] += 0.5 * r
		im[2] -= 0.5 * r
		im[3] -= 0.5 * r
	}

	if dd >= 0 {
		re[0] += 0.5 * bigD
		re[1] -= 0.5 * bigD
	} else {
		im[0] += 0.5 * bigD
		im[1] -= 0.5 * bigD
	}

	if ee >= 0 {
		re[2] += 0.5 * bigE
		re[3] -= 0.5 * bigE
	} else {
		im[2] += 0.5 * bigE
		im[3] -= 0.5 * bigE
	}

	return [4]complex128{
		complex(re[0], im[0]),
		complex(re[1], im[1]),
		complex(re[2], im[2]),
		complex(re[3], im[3]),
	}, nil
}

func main() {
	prompter, err := promptio.New(os.Stdin, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, "234:", err)
		os.Exit(1)
	}

	fmt.Println("SOLVE QUADRATIC, CUBIC, AND QUARTIC EQUATIONS")
	fmt.Println()

	var order int
	for {
		order = prompter.Int("Order of equation to solve (2,3,4)", 2)
		if order >= 2 && order <= 4 {
			break
		}
		fmt.Println("NOT A VALID OPTION")
	}
	fmt.Println()

	coeffs := make([]float64, order+1)
	for i := order; i >= 0; i-- {
		coeffs[i] = prompter.Float(fmt.Sprintf("For order %d term, input coefficient", i), 0)
	}

	var roots []complex128
	switch order {
	case 2:
		r := solveQuadratic(coeffs[0], coeffs[1], coeffs[2])
		roots = r[:]
	case 3:
		r := solveCubic(coeffs[0], coeffs[1], coeffs[2], coeffs[3])
		roots = r[:]
	case 4:
		r, err := solveQuartic(coeffs[0], coeffs[1], coeffs[2], coeffs[3], coeffs[4])
		if err != nil {
			fmt.Fprintln(os.Stderr, "234:", err)
			os.Exit(1)
		}
		roots = r[:]
	}

	fmt.Println("\nSolutions to:")
	for i := order; i >= 0; i-- {
		if i != 0 {
			fmt.Printf(" %+.4f * x^%d", coeffs[i], i)
		} else {
			fmt.Printf(" %+.4f", coeffs[i])
		}
	}
	fmt.Println(" = 0")
	fmt.Println("are:")
	for _, root := range roots {
		fmt.Printf("real:  %f  imaginary:  %f\n", real(root), imag(root))
	}
}
