// triangle solves a plane triangle given any three of its six parts
// (three sides and the three angles opposite them), provided at
// least one side is given, then reports the area and the radii of
// the circumscribed and inscribed circles. The side-side-angle case
// is ambiguous (two different triangles can share the same two
// sides and non-included angle), so both solutions are reported when
// they differ.
//
// Converted from TRIANGLE.C (M. W. Klotz, 11/98, 3/99, 4/11),
// Math/triangle.
package main

import (
	"errors"
	"fmt"
	"math"
	"os"

	"mwkgo/internal/promptio"
)

// triangleSolution is a fully solved triangle: sides S[i] opposite
// angles A[i], for i = 0, 1, 2.
type triangleSolution struct {
	S, A [3]float64
}

// ssaAlternate is the second possible solution to an ambiguous
// side-side-angle case, nil when the case isn't ambiguous or the two
// solutions coincide.
type ssaAlternate struct {
	S, A [3]float64
}

var errInsufficientTriangleData = errors.New("insufficient data for solution")
var errNoSideSpecified = errors.New("at least one side must be specified")

// sssSolve returns the three angles of a triangle given all three
// sides, via the half-angle formula.
func sssSolve(s1, s2, s3 float64) (a1, a2, a3 float64) {
	p := 0.5 * (s1 + s2 + s3)
	a1 = 2 * math.Acos(math.Sqrt(p*(p-s1)/(s2*s3))) * 180 / math.Pi
	a2 = 2 * math.Acos(math.Sqrt(p*(p-s2)/(s1*s3))) * 180 / math.Pi
	a3 = 180 - a1 - a2
	return
}

// asaSolve returns the angle opposite sideI and the two other sides,
// given sideI and the two angles opposite the other two sides
// (angleJ, angleK).
func asaSolve(sideI, angleJ, angleK float64) (angleI, sideJ, sideK float64) {
	angleI = 180 - angleJ - angleK
	sideJ = sideI * sind(angleJ) / sind(angleI)
	sideK = sideI*cosd(angleJ) + sideJ*cosd(angleI)
	return
}

// saaSolve returns the third angle and the other two sides, given
// sideI, its own opposite angle angleI, and another known angle
// angleK.
func saaSolve(sideI, angleI, angleK float64) (angleJ, sideJ, sideK float64) {
	angleJ = 180 - angleI - angleK
	sideJ = sideI * sind(angleJ) / sind(angleI)
	sideK = sideI*cosd(angleJ) + sideJ*cosd(angleI)
	return
}

// sasSolve returns the third side and the other two angles, given
// two sides sideI, sideJ and the angle angleK between them (the
// angle opposite the third, unknown side).
func sasSolve(sideI, sideJ, angleK float64) (sideK, angleI, angleJ float64) {
	sideK = math.Sqrt(sideI*sideI + sideJ*sideJ - 2*sideI*sideJ*cosd(angleK))
	p := 0.5 * (sideI + sideJ + sideK)
	angleI = 2 * math.Acos(math.Sqrt(p*(p-sideI)/(sideJ*sideK))) * 180 / math.Pi
	angleJ = 2 * math.Acos(math.Sqrt(p*(p-sideJ)/(sideI*sideK))) * 180 / math.Pi
	return
}

// ssaSolve returns the other two angles and the third side, given
// two sides sideI, sideJ and the angle angleI opposite sideI. This
// case is ambiguous whenever sideJ > sideI and the computed angleJ
// isn't a right angle: alt holds the second possible solution in
// that case, nil otherwise.
func ssaSolve(sideI, sideJ, angleI float64) (angleJ, angleK, sideK float64, alt *ssaAlternate) {
	angleJ = asind(sideJ * sind(angleI) / sideI)
	angleK = 180 - angleI - angleJ
	sideK = sideI*cosd(angleJ) + sideJ*cosd(angleI)

	if sideJ > sideI && angleJ != 90 {
		altAngleJ := 180 - angleJ
		altAngleK := 180 - angleI - altAngleJ
		altSideK := sideI*cosd(altAngleJ) + sideJ*cosd(angleI)
		alt = &ssaAlternate{A: [3]float64{angleI, altAngleJ, altAngleK}, S: [3]float64{sideI, sideJ, altSideK}}
	}
	return angleJ, angleK, sideK, alt
}

func sind(deg float64) float64 { return math.Sin(deg * math.Pi / 180) }
func cosd(deg float64) float64 { return math.Cos(deg * math.Pi / 180) }
func asind(x float64) float64  { return math.Asin(x) * 180 / math.Pi }

// otherTwoAscending returns the two indices, in ascending order,
// other than i, from {0, 1, 2}.
func otherTwoAscending(i int) (j, k int) {
	rest := make([]int, 0, 2)
	for x := 0; x < 3; x++ {
		if x != i {
			rest = append(rest, x)
		}
	}
	return rest[0], rest[1]
}

// solveTriangle solves for every side and angle of a triangle given
// any three of the six (zero meaning unknown), matching the original
// program's own dispatch to whichever of SSS, ASA, SAA, SAS, or SSA
// applies. At least one side must be given.
func solveTriangle(s, a [3]float64) (triangleSolution, *ssaAlternate, error) {
	if s[0] == 0 && s[1] == 0 && s[2] == 0 {
		return triangleSolution{}, nil, errNoSideSpecified
	}

	// SSS
	if s[0] != 0 && s[1] != 0 && s[2] != 0 {
		a[0], a[1], a[2] = sssSolve(s[0], s[1], s[2])
		return triangleSolution{S: s, A: a}, nil, nil
	}

	// ASA: one side and the other two angles.
	for i := 0; i < 3; i++ {
		j, k := otherTwoAscending(i)
		if s[i] != 0 && a[j] != 0 && a[k] != 0 {
			a[i], s[j], s[k] = asaSolve(s[i], a[j], a[k])
			return triangleSolution{S: s, A: a}, nil, nil
		}
	}

	// SAA: one side, its own opposite angle, and one other angle.
	for i := 0; i < 3; i++ {
		for k := 0; k < 3; k++ {
			if k == i {
				continue
			}
			j := 3 - i - k
			if s[i] != 0 && a[i] != 0 && a[k] != 0 {
				a[j], s[j], s[k] = saaSolve(s[i], a[i], a[k])
				return triangleSolution{S: s, A: a}, nil, nil
			}
		}
	}

	// SAS: two sides and the angle between them.
	for _, pair := range [3][2]int{{0, 1}, {1, 2}, {0, 2}} {
		i, j := pair[0], pair[1]
		k := 3 - i - j
		if s[i] != 0 && s[j] != 0 && a[k] != 0 {
			s[k], a[i], a[j] = sasSolve(s[i], s[j], a[k])
			return triangleSolution{S: s, A: a}, nil, nil
		}
	}

	// SSA: two sides and the angle opposite one of them (ambiguous).
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if j == i {
				continue
			}
			k := 3 - i - j
			if s[i] != 0 && s[j] != 0 && a[i] != 0 {
				var alt *ssaAlternate
				a[j], a[k], s[k], alt = ssaSolve(s[i], s[j], a[i])
				return triangleSolution{S: s, A: a}, alt, nil
			}
		}
	}

	return triangleSolution{}, nil, errInsufficientTriangleData
}

// triangleProperties holds the area and circle radii derived from a
// solved triangle.
type triangleProperties struct {
	Area, CircumscribedRadius, InscribedRadius float64
}

func computeTriangleProperties(sol triangleSolution) triangleProperties {
	s := 0.5 * (sol.S[0] + sol.S[1] + sol.S[2])
	return triangleProperties{
		Area:                0.5 * sol.S[0] * sol.S[1] * sind(sol.A[2]),
		CircumscribedRadius: 0.5 * sol.S[0] / sind(sol.A[0]),
		InscribedRadius:     math.Sqrt((s - sol.S[0]) * (s - sol.S[1]) * (s - sol.S[2]) / s),
	}
}

func main() {
	prompter, err := promptio.New(os.Stdin, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, "triangle:", err)
		os.Exit(1)
	}

	fmt.Println("SOLUTION OF PLANE TRIANGLES")
	fmt.Println()
	fmt.Println("Number the sides of your triangle 1-3 and answer the questions to follow.")
	fmt.Println("Input whatever data you know; press return if not known.")
	fmt.Println("You need three data items (at least one side) to solve a triangle.")
	fmt.Println()

	var s, a [3]float64
	s[0] = prompter.Float("side 1", 0)
	s[1] = prompter.Float("side 2", 0)
	s[2] = prompter.Float("side 3", 0)

	known := 0
	for _, v := range s {
		if v != 0 {
			known++
		}
	}
	if known == 0 {
		fmt.Println("\nAT LEAST ONE SIDE MUST BE SPECIFIED")
		return
	}

	if known < 3 {
		a[0] = prompter.Float("angle opposite side 1 (deg)", 0)
		if a[0] != 0 {
			known++
		}
	}
	if known < 3 {
		a[1] = prompter.Float("angle opposite side 2 (deg)", 0)
		if a[1] != 0 {
			known++
		}
	}
	if known < 3 {
		a[2] = prompter.Float("angle opposite side 3 (deg)", 0)
	}

	sol, alt, err := solveTriangle(s, a)
	if err != nil {
		fmt.Println("\nINSUFFICIENT DATA FOR SOLUTION")
		return
	}

	printSolution(sol)
	if alt != nil {
		altSol := triangleSolution{S: alt.S, A: alt.A}
		if !solutionsMatch(sol, altSol) {
			fmt.Println("AN ALTERNATE SOLUTION EXISTS:")
			printSolution(altSol)
		}
	}
}

func solutionsMatch(a, b triangleSolution) bool {
	const tol = 1e-6
	for i := 0; i < 3; i++ {
		if math.Abs(a.S[i]-b.S[i]) > tol || math.Abs(a.A[i]-b.A[i]) > tol {
			return false
		}
	}
	return true
}

func printSolution(sol triangleSolution) {
	p := computeTriangleProperties(sol)
	fmt.Printf("\nside  1 = %g\n", sol.S[0])
	fmt.Printf("side  2 = %g\n", sol.S[1])
	fmt.Printf("side  3 = %g\n", sol.S[2])
	fmt.Printf("angle opposite side 1 = %g deg\n", sol.A[0])
	fmt.Printf("angle opposite side 2 = %g deg\n", sol.A[1])
	fmt.Printf("angle opposite side 3 = %g deg\n", sol.A[2])
	fmt.Printf("area of triangle = %g\n", p.Area)
	fmt.Printf("circumscribed circle radius = %g\n", p.CircumscribedRadius)
	fmt.Printf("inscribed circle radius = %g\n", p.InscribedRadius)
}
