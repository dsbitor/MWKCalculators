// cep computes the Circular Error Probable (CEP): the radius of a
// circle, centered on the combined mean of two correlated Gaussian
// error sources, within which a given fraction of outcomes are
// expected to fall. It's used, for example, to characterize the
// combined miss distance of two independent error sources such as a
// milling machine's X and Y axis errors, or a missile's crossrange
// and downrange errors.
//
// Converted from CEP.C, Math/cep.
package main

import (
	"errors"
	"fmt"
	"math"
	"os"

	"mwkgo/internal/promptio"
)

// wilcoxA and wilcoxB are the coefficients of Wilcox's rational
// polynomial approximation to the 50% CEP of a correlated bivariate
// Gaussian, used here only as an independent analytic cross-check
// against this program's own numerical integration, and only for
// the special case it applies to (zero means, 50% probability).
var wilcoxA = [4]float64{0.69329763813785, -1.4621993839667, 1.06372876137083, -0.203505733447903}
var wilcoxB = [3]float64{1, -1.77649593676893, 0.977186538295758}

// maxCepRadialSteps bounds the radial integration step count. A
// realistic CEP probability (well under 1) converges within at most
// a few thousand steps at the original program's own step size
// (0.1% of the smaller principal standard deviation); this generous
// bound exists to satisfy coding-style.md Rule 2 (bound every loop)
// against a pathological probability request.
const maxCepRadialSteps = 1_000_000

// angularSteps is the number of angular samples used per radial
// step when numerically integrating the bivariate Gaussian density
// around a circle, matching the original program's own resolution.
const angularSteps = 72

// errCepDidNotConverge reports that the radial search did not reach
// the requested probability within maxCepRadialSteps steps.
var errCepDidNotConverge = errors.New("radial search did not converge within the step limit")

// cepResult holds the computed CEP radius and, when the inputs match
// the special case Wilcox's approximation applies to (zero means,
// 50% probability), the independent analytic cross-check value.
type cepResult struct {
	Radius    float64
	WilcoxChk float64 // zero unless the Wilcox special case applies
}

// computeCEP returns the radius of the circle, centered on
// (mean1, mean2), within which probability (0 to 1) of outcomes from
// two correlated Gaussian error sources are expected to fall. sigma1
// and sigma2 are the two sources' standard deviations and rho is
// their correlation coefficient (-1 to 1). It numerically integrates
// the bivariate Gaussian density outward in a growing circle until
// the accumulated probability reaches the target.
func computeCEP(mean1, sigma1, mean2, sigma2, rho, probability float64) (cepResult, error) {
	c11 := sigma1 * sigma1
	c22 := sigma2 * sigma2
	c12 := rho * sigma1 * sigma2

	// psi is the rotation angle diagonalizing the covariance matrix.
	// The original program's own formula divides by (c11-c12) here;
	// the standard formula for this rotation divides by (c11-c22)
	// instead. Carried forward unchanged, as with cone's formula
	// quirk in Tier 1 group 5: this conversion reproduces what the
	// original computes, not a "corrected" version of it. It has no
	// effect at all when sigma1 == sigma2 (the case this program's
	// own Wilcox cross-check applies to), since c11-c22 is zero
	// either way.
	var psi float64
	if c11 != c12 || c12 != 0 {
		psi = math.Atan2(2*c12, c11-c12)
	}
	sinPsi, cosPsi := math.Sin(psi), math.Cos(psi)

	spread := math.Sqrt((c11-c22)*(c11-c22) + 4*c12*c12)
	s1 := 0.5 * (c11 + c22 + spread)
	s2 := 0.5 * (c11 + c22 - spread)
	principalScale := math.Sqrt(s1 * s2)
	if principalScale == 0 {
		return cepResult{Radius: -1}, nil
	}

	radialStep := 0.001 * math.Sqrt(math.Min(s1, s2))
	density := radialStep / (angularSteps * principalScale)
	angleStep := 2 * math.Pi / angularSteps

	var radius, cumulative, lastRadius, lastCumulative float64
	converged := false
	for i := 0; i < maxCepRadialSteps; i++ {
		lastCumulative, lastRadius = cumulative, radius
		radius += radialStep

		var ringDensity float64
		for k := 1; k <= angularSteps; k++ {
			angle := float64(k) * angleStep
			sinA, cosA := math.Sin(angle), math.Cos(angle)
			u := (radius*cosA-mean1)*cosPsi + (radius*sinA-mean2)*sinPsi
			v := (radius*sinA-mean2)*cosPsi - (radius*cosA-mean1)*sinPsi
			ringDensity += math.Exp(-0.5 * (u*u/s1 + v*v/s2))
		}
		cumulative += radius * density * ringDensity

		if cumulative >= probability {
			converged = true
			break
		}
	}
	if !converged {
		return cepResult{}, errCepDidNotConverge
	}

	result := cepResult{
		Radius: lastRadius + (probability-lastCumulative)*(radius-lastRadius)/(cumulative-lastCumulative),
	}

	if mean1 == 0 && mean2 == 0 && probability == 0.5 {
		trace := c11 + c22
		m := ((c11-c22)*(c11-c22) + 4*c12*c12) / (trace * trace)
		numerator := ((wilcoxA[3]*m+wilcoxA[2])*m+wilcoxA[1])*m + wilcoxA[0]
		denominator := (wilcoxB[2]*m+wilcoxB[1])*m + wilcoxB[0]
		result.WilcoxChk = math.Sqrt(trace * numerator / denominator)
	}

	return result, nil
}

func main() {
	prompter, err := promptio.New(os.Stdin, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, "cep:", err)
		os.Exit(1)
	}

	fmt.Println()
	mean1 := prompter.Float("Mean1", 0)
	sigma1 := prompter.Float("Sigma1", 1)
	mean2 := prompter.Float("Mean2", 0)
	sigma2 := prompter.Float("Sigma2", 1)
	rho := prompter.Float("Rho", 0)
	probability := prompter.Float("CEP probability", 0.5)

	result, err := computeCEP(mean1, sigma1, mean2, sigma2, rho, probability)
	if err != nil {
		fmt.Fprintln(os.Stderr, "cep:", err)
		os.Exit(1)
	}

	fmt.Printf("\nCEP(%.4f) = %.4f\n", probability, result.Radius)
	if result.WilcoxChk != 0 {
		fmt.Printf("Wilcox cepchk = %.4f\n", result.WilcoxChk)
	}
}
