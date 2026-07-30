// ellipse computes an ellipse's eccentricity, flattening, area, and
// perimeter, given its semi-major and semi-minor axes. The
// perimeter has no closed-form rational expression, so it is found
// from a high-accuracy polynomial approximation to the complete
// elliptic integral of the second kind (Abramowitz & Stegun 17.3.36,
// accurate to 2e-8) and used as a reference against which three
// well known algebraic approximations (an RMS approximation and two
// due to Ramanujan) are compared.
//
// Converted from ELLIPSE.C (M. W. Klotz), Math/ellipse.
package main

import (
	"fmt"
	"math"
	"os"

	"mwkgo/internal/promptio"
)

// ellipticE2 approximates the complete elliptic integral of the
// second kind as a function of the complementary parameter
// m = (b/a)^2, per Abramowitz & Stegun 17.3.36. Accuracy is
// approximately 2e-8 over the full domain 0 < m <= 1.
func ellipticE2(m float64) float64 {
	const a1, a2, a3, a4 = 0.44325141463, 0.06260601220, 0.04757383546, 0.01736506451
	const b1, b2, b3, b4 = 0.24998368310, 0.09200180037, 0.04069697526, 0.00526449639

	polyA := 1 + a1*m + a2*m*m + a3*m*m*m + a4*m*m*m*m
	polyB := b1*m + b2*m*m + b3*m*m*m + b4*m*m*m*m
	return polyA + polyB*math.Log(1/m)
}

// ellipseProperties holds the computed geometric properties of an
// ellipse with semi-major axis a and semi-minor axis b, including
// its exact perimeter (via the elliptic integral) and three
// algebraic approximations to it, each with its error relative to
// the exact value.
type ellipseProperties struct {
	Eccentricity, Flattening, Area float64
	ExactPerimeter                 float64
	RMSPerimeter, RMSErrorPct      float64
	Ramanujan1Perimeter            float64
	Ramanujan1ErrorPct             float64
	Ramanujan2Perimeter            float64
	Ramanujan2ErrorPct             float64
}

// computeEllipseProperties returns the geometric properties of an
// ellipse with semi-major axis a and semi-minor axis b.
func computeEllipseProperties(a, b float64) ellipseProperties {
	eccentricity := math.Sqrt(a*a-b*b) / a
	flattening := (a - b) / a
	area := math.Pi * a * b

	exact := 4 * a * ellipticE2((b*b)/(a*a))

	rms := 2 * math.Pi * math.Sqrt(0.5*(a*a+b*b))
	ramanujan1 := math.Pi * (3*(a+b) - math.Sqrt((3*a+b)*(a+3*b)))

	x := (a - b) / (a + b)
	ramanujan2 := math.Pi * (a + b) * (1 + 3*x*x/(10+math.Sqrt(4-3*x*x)))

	return ellipseProperties{
		Eccentricity:        eccentricity,
		Flattening:          flattening,
		Area:                area,
		ExactPerimeter:      exact,
		RMSPerimeter:        rms,
		RMSErrorPct:         100 * (exact - rms) / exact,
		Ramanujan1Perimeter: ramanujan1,
		Ramanujan1ErrorPct:  100 * (exact - ramanujan1) / exact,
		Ramanujan2Perimeter: ramanujan2,
		Ramanujan2ErrorPct:  100 * (exact - ramanujan2) / exact,
	}
}

func main() {
	prompter, err := promptio.New(os.Stdin, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, "ellipse:", err)
		os.Exit(1)
	}

	a := prompter.Float("Semi-major axis", 10.0)
	b := prompter.Float("Semi-minor axis", 3.0)

	p := computeEllipseProperties(a, b)

	fmt.Println()
	fmt.Printf("eccentricity = %.4f\n", p.Eccentricity)
	fmt.Printf("flattening = %.4f\n", p.Flattening)
	fmt.Printf("area = %.4f\n", p.Area)

	fmt.Println("\nPerimeter = ")
	fmt.Println()
	fmt.Printf("        elliptic integral: %.4f\n", p.ExactPerimeter)
	fmt.Printf("        rms approximation: %.4f (error = %.4f %%)\n", p.RMSPerimeter, p.RMSErrorPct)
	fmt.Printf("Ramanujuan1 approximation: %.4f (error = %.4f %%)\n", p.Ramanujan1Perimeter, p.Ramanujan1ErrorPct)
	fmt.Printf("Ramanujuan2 approximation: %.4f (error = %.4f %%)\n", p.Ramanujan2Perimeter, p.Ramanujan2ErrorPct)
}
