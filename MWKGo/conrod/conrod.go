// conrod finds the worst-case clearance between an engine's
// connecting rod centerline and the bottom of the cylinder wall,
// which occurs when the crank radius and the connecting rod are at
// right angles. Given the rod length, crank radius, the height of
// the cylinder bottom above the crank center, and the cylinder
// diameter, it reports that clearance along with several related
// distances useful when laying out the design.
//
// Converted from CONROD.C (M. W. Klotz, with Tom Roach),
// WorkshopUtilities/conrod.
package main

import (
	"fmt"
	"math"
	"os"

	"mwkgo/internal/mwktrig"
	"mwkgo/internal/promptio"
)

// conRodGeometry holds the worst-case distances computed for a
// connecting rod and crank at right angles, using the point labels
// from the original program's accompanying diagram.
type conRodGeometry struct {
	Phi                    float64 // angle between the con rod and the line to the crank center
	GudgeonToCrankDistance float64 // xa: distance, gudgeon pin to crank center
	D34, D45, D35          float64
	D23                    float64 // clearance, con rod centerline to cylinder bottom
	D13, D12, D14, D25     float64
}

// computeConRodGeometry returns the worst-case connecting rod
// clearance geometry for a rod of length rodLength, a crank of
// radius crankRadius, a cylinder bottom at height cylinderBottomHeight
// above the crank center, and a cylinder of diameter cylinderDiameter.
func computeConRodGeometry(rodLength, crankRadius, cylinderBottomHeight, cylinderDiameter float64) conRodGeometry {
	gudgeonToCrankDistance := math.Hypot(crankRadius, rodLength)
	cosPhi := rodLength / gudgeonToCrankDistance
	sinPhi := crankRadius / gudgeonToCrankDistance
	tanPhi := sinPhi / cosPhi

	d34 := 0.5 * cylinderDiameter
	d45 := (gudgeonToCrankDistance - cylinderBottomHeight) * tanPhi
	d35 := d34 - d45
	d23 := d35 / cosPhi
	d13 := d34 / cosPhi
	d12 := d13 - d23
	d14 := d34 * tanPhi
	d25 := d35 * sinPhi

	return conRodGeometry{
		Phi:                    mwktrig.ClampedAcosDeg(cosPhi),
		GudgeonToCrankDistance: gudgeonToCrankDistance,
		D34:                    d34,
		D45:                    d45,
		D35:                    d35,
		D23:                    d23,
		D13:                    d13,
		D12:                    d12,
		D14:                    d14,
		D25:                    d25,
	}
}

func main() {
	prompter, err := promptio.New(os.Stdin, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, "conrod:", err)
		os.Exit(1)
	}

	fmt.Println("CONNECTING ROD CLEARANCE")
	fmt.Println()
	fmt.Println("Units of measurement don't matter but must be consistent")

	rodLength := prompter.Float("Connecting rod length (center-to-center)", 2.4)
	fmt.Println("Radius measured from crank center to connecting rod driver center")
	crankRadius := prompter.Float("Crank radius", 0.6)
	cylinderBottomHeight := prompter.Float("Height of cylinder bottom above crank center", 1.5)
	cylinderDiameter := prompter.Float("Cylinder diameter", 1.0)

	g := computeConRodGeometry(rodLength, crankRadius, cylinderBottomHeight, cylinderDiameter)

	fmt.Println("\nAt worst case, maximum con rod lateral offset:")
	fmt.Printf("Phi = %.4f deg\n", g.Phi)
	fmt.Printf("Distance, gudgeon pin to crank center = %.4f\n", g.GudgeonToCrankDistance)
	fmt.Printf("Clearance, con rod centerline to cylinder bottom = %.4f\n", g.D23)
	fmt.Printf("d34 = %.4f\n", g.D34)
	fmt.Printf("d45 = %.4f\n", g.D45)
	fmt.Printf("d35 = %.4f\n", g.D35)
	fmt.Printf("d23 = %.4f\n", g.D23)
	fmt.Printf("d13 = %.4f\n", g.D13)
	fmt.Printf("d12 = %.4f\n", g.D12)
	fmt.Printf("d14 = %.4f\n", g.D14)
	fmt.Printf("d25 = %.4f\n", g.D25)
}
