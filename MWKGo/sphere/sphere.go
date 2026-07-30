// sphere solves a spherical triangle given any three of its six
// parts (three sides and the three angles opposite them, all
// measured in degrees, with sides expressed as the angle they
// subtend at the sphere's center), then reports the triangle's area
// as a multiple of the sphere's radius squared. Unlike a plane
// triangle, three known angles alone (AAA) are enough to solve a
// spherical triangle, since spherical triangles don't have similar
// copies at different sizes. The side-side-angle (SSA) and
// angle-angle-side (AAS) cases are ambiguous, so both solutions are
// reported when they differ.
//
// Converted from SPHERE.C (M. W. Klotz, 3/99), Math/sphere.
package main

import (
	"errors"
	"fmt"
	"math"
	"os"

	"mwkgo/internal/promptio"
)

// sphereSolution is a fully solved spherical triangle: sides S[i]
// opposite angles A[i], for i = 0, 1, 2, all in degrees.
type sphereSolution struct {
	S, A [3]float64
}

// sphereAlternate is the second possible solution to an ambiguous
// SSA or AAS case, nil when the case isn't ambiguous or the two
// solutions coincide.
type sphereAlternate struct {
	S, A [3]float64
}

var errInsufficientSphereData = errors.New("insufficient data for solution")

func sind(deg float64) float64 { return math.Sin(deg * math.Pi / 180) }
func cosd(deg float64) float64 { return math.Cos(deg * math.Pi / 180) }
func tand(deg float64) float64 { return math.Tan(deg * math.Pi / 180) }
func asind(x float64) float64  { return math.Asin(clamp(x)) * 180 / math.Pi }
func acosd(x float64) float64  { return math.Acos(clamp(x)) * 180 / math.Pi }
func atand(x float64) float64  { return math.Atan(x) * 180 / math.Pi }

func clamp(x float64) float64 {
	if x > 1 {
		return 1
	}
	if x < -1 {
		return -1
	}
	return x
}

// sssSolve returns the three angles of a spherical triangle given
// all three sides, via the spherical law of cosines for sides.
func sssSolve(s1, s2, s3 float64) (a1, a2, a3 float64) {
	a1 = acosd((cosd(s1) - cosd(s2)*cosd(s3)) / (sind(s2) * sind(s3)))
	a2 = acosd((cosd(s2) - cosd(s1)*cosd(s3)) / (sind(s1) * sind(s3)))
	a3 = acosd((cosd(s3) - cosd(s1)*cosd(s2)) / (sind(s1) * sind(s2)))
	return
}

// aaaSolve returns the three sides of a spherical triangle given all
// three angles, via the spherical law of cosines for angles (the
// "polar dual" of sssSolve). Three angles alone determine a unique
// spherical triangle, unlike in plane geometry.
func aaaSolve(a1, a2, a3 float64) (s1, s2, s3 float64) {
	s1 = acosd((cosd(a1) + cosd(a2)*cosd(a3)) / (sind(a2) * sind(a3)))
	s2 = acosd((cosd(a2) + cosd(a1)*cosd(a3)) / (sind(a1) * sind(a3)))
	s3 = acosd((cosd(a3) + cosd(a1)*cosd(a2)) / (sind(a1) * sind(a2)))
	return
}

// sasSolve returns sideI (via the spherical law of cosines, given
// the angle angleI between sideJ and sideK) and the other two
// angles.
func sasSolve(angleI, sideJ, sideK float64) (sideI, angleJ, angleK float64) {
	sideI = acosd(cosd(sideJ)*cosd(sideK) + sind(sideJ)*sind(sideK)*cosd(angleI))
	angleJ = acosd((cosd(sideJ) - cosd(sideI)*cosd(sideK)) / (sind(sideI) * sind(sideK)))
	angleK = acosd((cosd(sideK) - cosd(sideI)*cosd(sideJ)) / (sind(sideI) * sind(sideJ)))
	return
}

// asaSolve returns angleR (via the spherical law of cosines for
// angles, given the side sideR between vertices P and Q) and the
// other two sides. sideP is opposite angleP, sideQ opposite angleQ,
// and sideR opposite the angleR being solved for.
func asaSolve(angleP, angleQ, sideR float64) (angleR, sideP, sideQ float64) {
	angleR = acosd(sind(angleP)*sind(angleQ)*cosd(sideR) - cosd(angleP)*cosd(angleQ))
	sideP = acosd((cosd(angleP) + cosd(angleQ)*cosd(angleR)) / (sind(angleQ) * sind(angleR)))
	sideQ = acosd((cosd(angleQ) + cosd(angleP)*cosd(angleR)) / (sind(angleP) * sind(angleR)))
	return
}

// ssaSolve returns the angle opposite sideJ, the third side, and the
// third angle, given two sides sideI, sideJ and the angle angleI
// opposite sideI. Ambiguous whenever sideI < sideJ: alt holds the
// second possible solution in that case.
func ssaSolve(sideI, sideJ, angleI float64) (angleJ, sideK, angleK float64, alt *sphereAlternate) {
	angleJ = asind(sind(sideJ) * sind(angleI) / sind(sideI))
	sideK = 2 * atand(sind(0.5*(angleI+angleJ))*tand(0.5*(sideI-sideJ))/sind(0.5*(angleI-angleJ)))
	angleK = acosd((cosd(sideK) - cosd(sideI)*cosd(sideJ)) / (sind(sideI) * sind(sideJ)))

	if sideI < sideJ {
		altAngleJ := 180 - angleJ
		altSideK := 2 * atand(sind(0.5*(angleI+altAngleJ))*tand(0.5*(sideI-sideJ))/sind(0.5*(angleI-altAngleJ)))
		altAngleK := acosd((cosd(altSideK) - cosd(sideI)*cosd(sideJ)) / (sind(sideI) * sind(sideJ)))
		alt = &sphereAlternate{
			S: [3]float64{sideI, sideJ, altSideK},
			A: [3]float64{angleI, altAngleJ, altAngleK},
		}
	}
	return angleJ, sideK, angleK, alt
}

// aasSolve returns the side opposite angleJ, the third angle, and
// the third side, given two angles angleI, angleJ and the side
// sideI opposite angleI. Ambiguous whenever angleI < angleJ: alt
// holds the second possible solution in that case.
func aasSolve(angleI, sideI, angleJ float64) (sideJ, angleK, sideK float64, alt *sphereAlternate) {
	sideJ = asind(sind(angleJ) * sind(sideI) / sind(angleI))
	z := sind(0.5*(sideI+sideJ)) * tand(0.5*(angleI-angleJ)) / sind(0.5*(sideI-sideJ))
	angleK = 2 * atand(1/z)
	sideK = acosd((cosd(angleK) + cosd(angleI)*cosd(angleJ)) / (sind(angleI) * sind(angleJ)))

	if angleI < angleJ {
		altSideJ := 180 - sideJ
		altZ := sind(0.5*(sideI+altSideJ)) * tand(0.5*(angleI-angleJ)) / sind(0.5*(sideI-altSideJ))
		altAngleK := 2 * atand(1/altZ)
		altSideK := acosd((cosd(altAngleK) + cosd(angleI)*cosd(angleJ)) / (sind(angleI) * sind(angleJ)))
		alt = &sphereAlternate{
			S: [3]float64{sideI, altSideJ, altSideK},
			A: [3]float64{angleI, angleJ, altAngleK},
		}
	}
	return sideJ, angleK, sideK, alt
}

// solveSphericalTriangle solves for every side and angle of a
// spherical triangle given any three of the six (zero meaning
// unknown), matching the original program's own dispatch to
// whichever of SSS, AAA, SAS, ASA, SSA, or AAS applies.
func solveSphericalTriangle(s, a [3]float64) (sphereSolution, *sphereAlternate, error) {
	if s[0] != 0 && s[1] != 0 && s[2] != 0 {
		a[0], a[1], a[2] = sssSolve(s[0], s[1], s[2])
		return sphereSolution{S: s, A: a}, nil, nil
	}
	if a[0] != 0 && a[1] != 0 && a[2] != 0 {
		s[0], s[1], s[2] = aaaSolve(a[0], a[1], a[2])
		return sphereSolution{S: s, A: a}, nil, nil
	}

	// SAS: angle i, and the two sides adjacent to it.
	for i := 0; i < 3; i++ {
		j, k := (i+1)%3, (i+2)%3
		if a[i] != 0 && s[j] != 0 && s[k] != 0 {
			s[i], a[j], a[k] = sasSolve(a[i], s[j], s[k])
			return sphereSolution{S: s, A: a}, nil, nil
		}
	}

	// ASA: side r, and the two angles adjacent to it.
	for r := 0; r < 3; r++ {
		p, q := (r+1)%3, (r+2)%3
		if s[r] != 0 && a[p] != 0 && a[q] != 0 {
			a[r], s[p], s[q] = asaSolve(a[p], a[q], s[r])
			return sphereSolution{S: s, A: a}, nil, nil
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
				var alt *sphereAlternate
				a[j], s[k], a[k], alt = ssaSolve(s[i], s[j], a[i])
				return sphereSolution{S: s, A: a}, alt, nil
			}
		}
	}

	// AAS: two angles and the side opposite one of them (ambiguous).
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if j == i {
				continue
			}
			k := 3 - i - j
			if a[i] != 0 && s[i] != 0 && a[j] != 0 {
				var alt *sphereAlternate
				s[j], a[k], s[k], alt = aasSolve(a[i], s[i], a[j])
				return sphereSolution{S: s, A: a}, alt, nil
			}
		}
	}

	return sphereSolution{}, nil, errInsufficientSphereData
}

// sphericalExcessArea returns the triangle's area as a multiple of
// the sphere's radius squared, via the spherical excess formula.
func sphericalExcessArea(sol sphereSolution) float64 {
	return (sol.A[0] + sol.A[1] + sol.A[2] - 180) * math.Pi / 180
}

func main() {
	prompter, err := promptio.New(os.Stdin, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, "sphere:", err)
		os.Exit(1)
	}

	fmt.Println("SOLUTION OF SPHERICAL TRIANGLES")
	fmt.Println()
	fmt.Println("Number the sides of your triangle 1-3 and answer the questions to follow.")
	fmt.Println("Input whatever data you know; press return if not known.")
	fmt.Println("You need three data items to solve a triangle.")
	fmt.Println()

	var s, a [3]float64
	s[0] = prompter.Float("side 1 (deg)", 0)
	s[1] = prompter.Float("side 2 (deg)", 0)
	s[2] = prompter.Float("side 3 (deg)", 0)

	known := 0
	for _, v := range s {
		if v != 0 {
			known++
		}
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

	sol, alt, err := solveSphericalTriangle(s, a)
	if err != nil {
		fmt.Println("\nINSUFFICIENT DATA FOR SOLUTION")
		return
	}

	printSphereSolution(sol)
	if alt != nil {
		altSol := sphereSolution{S: alt.S, A: alt.A}
		if !sphereSolutionsMatch(sol, altSol) {
			fmt.Println("AN ALTERNATE SOLUTION EXISTS:")
			printSphereSolution(altSol)
		}
	}
}

func sphereSolutionsMatch(a, b sphereSolution) bool {
	const tol = 1e-6
	for i := 0; i < 3; i++ {
		if math.Abs(a.S[i]-b.S[i]) > tol || math.Abs(a.A[i]-b.A[i]) > tol {
			return false
		}
	}
	return true
}

func printSphereSolution(sol sphereSolution) {
	fmt.Printf("\nside  1 = %g deg\n", sol.S[0])
	fmt.Printf("side  2 = %g deg\n", sol.S[1])
	fmt.Printf("side  3 = %g deg\n", sol.S[2])
	fmt.Printf("angle opposite side 1 = %g deg\n", sol.A[0])
	fmt.Printf("angle opposite side 2 = %g deg\n", sol.A[1])
	fmt.Printf("angle opposite side 3 = %g deg\n", sol.A[2])
	fmt.Printf("area of triangle = %g * (radius of sphere)^2\n", sphericalExcessArea(sol))
}
