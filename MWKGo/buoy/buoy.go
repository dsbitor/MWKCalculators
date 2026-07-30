// buoy computes how deep a buoy of a given shape and weight sinks
// into a liquid of a given density, per Archimedes' principle: the
// buoy sinks until the weight of displaced liquid equals the buoy's
// own weight plus any load it carries.
//
// Converted from BUOY.C, Misc/buoy.
package main

import (
	"errors"
	"fmt"
	"math"
	"os"

	"mwkgo/internal/promptio"
)

// maxBinarySearchSteps bounds the binary search for immersion depth,
// matching the original program's own runaway guard.
const maxBinarySearchSteps = 500

// sinkError reports that the buoy's total displaced-liquid weight,
// even fully submerged, is less than its own weight: it cannot float
// at all. MaxSupportableWeight is the most this buoy shape and size
// could support before sinking.
type sinkError struct {
	MaxSupportableWeight float64
}

func (e sinkError) Error() string {
	return fmt.Sprintf("buoy will sink: a buoy of this size can support <= %.4f lb", e.MaxSupportableWeight)
}

// sphereSubmergedVolume returns the volume of a sphere of diameter
// diameter submerged to depth depth (measured from the bottom of the
// sphere).
func sphereSubmergedVolume(diameter, depth float64) float64 {
	return math.Pi * depth * depth * (1.5*diameter - depth) / 3
}

// horizontalCylinderSegmentArea returns the cross-sectional area of
// the circular segment of a cylinder of diameter diameter submerged
// to depth depth, matching the original program's own fcyl function.
func horizontalCylinderSegmentArea(diameter, depth float64) float64 {
	r := 0.5 * diameter
	z := r - depth
	angle := 2 * math.Acos(z/r)
	chord := 2 * r * math.Sin(0.5*angle)
	return 0.5 * (r*r*angle - z*chord)
}

// binarySearchDepth finds the depth in [low, high] at which f(depth)
// equals target, to within tolerance, matching the original
// program's own binsearch function. It returns an error if no such
// depth is found within maxBinarySearchSteps steps.
func binarySearchDepth(low, high, target, tolerance float64, f func(float64) float64) (float64, error) {
	for i := 0; i < maxBinarySearchSteps; i++ {
		mid := low + 0.5*(high-low)
		y := f(mid)
		if math.Abs(y-target) <= tolerance {
			return mid, nil
		}
		if target < y {
			high = mid
		} else {
			low = mid
		}
	}
	return 0, errors.New("search did not converge within the step limit")
}

// sphereImmersionDepth returns how deep a sphere-shaped buoy of
// diameter diameter, in a liquid of density liquidDensity, sinks
// under weight weight.
func sphereImmersionDepth(diameter, liquidDensity, weight float64) (float64, error) {
	volume := math.Pi * diameter * diameter * diameter / 6
	if maxSupportable := volume * liquidDensity; maxSupportable < weight {
		return 0, sinkError{MaxSupportableWeight: maxSupportable}
	}
	return binarySearchDepth(0, diameter, weight, 0.01, func(depth float64) float64 {
		return liquidDensity * sphereSubmergedVolume(diameter, depth)
	})
}

// horizontalCylinderImmersionDepth returns how deep a horizontal
// cylindrical buoy of diameter diameter and length length, in a
// liquid of density liquidDensity, sinks under weight weight.
func horizontalCylinderImmersionDepth(diameter, length, liquidDensity, weight float64) (float64, error) {
	volume := 0.25 * math.Pi * diameter * diameter * length
	if maxSupportable := volume * liquidDensity; maxSupportable < weight {
		return 0, sinkError{MaxSupportableWeight: maxSupportable}
	}
	return binarySearchDepth(0, diameter, weight, 0.1, func(depth float64) float64 {
		return liquidDensity * horizontalCylinderSegmentArea(diameter, depth) * length
	})
}

// verticalCylinderImmersionDepth returns how deep a vertical
// cylindrical buoy of diameter diameter and length length, in a
// liquid of density liquidDensity, sinks under weight weight. Unlike
// the sphere and horizontal cylinder, a vertical cylinder's
// submerged volume is a simple linear function of depth, so this
// solves directly rather than by search.
func verticalCylinderImmersionDepth(diameter, length, liquidDensity, weight float64) (float64, error) {
	volume := 0.25 * math.Pi * diameter * diameter * length
	if maxSupportable := volume * liquidDensity; maxSupportable < weight {
		return 0, sinkError{MaxSupportableWeight: maxSupportable}
	}
	return weight / (0.25 * math.Pi * diameter * diameter * liquidDensity), nil
}

// boxImmersionDepth returns how deep a rectangular-parallelepiped
// buoy of the given length, width, and height, in a liquid of
// density liquidDensity, sinks under weight weight.
func boxImmersionDepth(length, width, height, liquidDensity, weight float64) (float64, error) {
	volume := length * width * height
	if maxSupportable := volume * liquidDensity; maxSupportable < weight {
		return 0, sinkError{MaxSupportableWeight: maxSupportable}
	}
	return weight / (length * width * liquidDensity), nil
}

func main() {
	prompter, err := promptio.New(os.Stdin, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, "buoy:", err)
		os.Exit(1)
	}

	fmt.Println("Densities (lb/in^3):")
	fmt.Println("Sea Water = 0.037")
	fmt.Println("Water = 0.0361")
	fmt.Println("Mercury = 0.4894")
	fmt.Println()

	liquidDensity := prompter.Float("Density of liquid (lb/in^3)", 0.0361)
	weight := prompter.Float("Weight of buoy and applied load (lb)", 3.5)

	fmt.Println("\nBuoy shape:")
	fmt.Println("a - Sphere")
	fmt.Println("b - Horizontal Cylinder")
	fmt.Println("c - Vertical Cylinder")
	fmt.Println("d - Rectangular Parallelepiped")

	shape := prompter.Line("(a-d) ? ")
	fmt.Println()

	var depth float64
	switch shape {
	case "a":
		diameter := prompter.Float("Diameter of sphere (in)", 10.0)
		depth, err = sphereImmersionDepth(diameter, liquidDensity, weight)
	case "b":
		diameter := prompter.Float("Diameter of cylinder (in)", 5.0)
		length := prompter.Float("Length of cylinder (in)", 10.0)
		depth, err = horizontalCylinderImmersionDepth(diameter, length, liquidDensity, weight)
	case "c":
		diameter := prompter.Float("Diameter of cylinder (in)", 5.0)
		length := prompter.Float("Length of cylinder (in)", 10.0)
		depth, err = verticalCylinderImmersionDepth(diameter, length, liquidDensity, weight)
	case "d":
		length := prompter.Float("Length of buoy (in)", 5.0)
		width := prompter.Float("Width of buoy (in)", 8.0)
		height := prompter.Float("Height of buoy (in)", 6.0)
		depth, err = boxImmersionDepth(length, width, height, liquidDensity, weight)
	default:
		fmt.Println("NOT A VALID OPTION")
		return
	}

	var sink sinkError
	if errors.As(err, &sink) {
		fmt.Println("Buoy will sink!")
		fmt.Printf("A buoy of this size can support <= %.4f lb\n", sink.MaxSupportableWeight)
		return
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "buoy:", err)
		os.Exit(1)
	}

	fmt.Printf("\nBuoy immersion depth = %.2f in\n", depth)
}
