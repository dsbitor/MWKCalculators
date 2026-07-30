// polygon computes the properties of a regular polygon (equal side
// lengths and angles) from any one of four possible size
// specifications: side length, distance across flats (even side
// count), distance from a flat to the opposite vertex (odd side
// count), or the diameter of the circumscribed or inscribed circle.
//
// Converted from POLYGON.C, Math/polygon.
package main

import (
	"errors"
	"fmt"
	"math"
	"os"

	"mwkgo/internal/promptio"
)

// regularPolygon holds the computed properties of a regular polygon
// with sideCount sides and circumradius circumradius.
type regularPolygon struct {
	SideCount            int
	CentralAngleDeg      float64
	VertexAngleDeg       float64
	CircumscribedDiam    float64
	InscribedDiam        float64
	AcrossFlats          float64 // only meaningful when SideCount is even
	FlatToOppositeVertex float64 // only meaningful when SideCount is odd
	SideLength           float64
	Perimeter            float64
	Area                 float64
}

// computeRegularPolygon returns the properties of a regular polygon
// with the given side count and circumradius (the radius of the
// circle passing through every vertex).
func computeRegularPolygon(sideCount int, circumradius float64) regularPolygon {
	centralAngle := 360.0 / float64(sideCount)
	halfAngleRad := 0.5 * centralAngle * math.Pi / 180
	apothem := circumradius * math.Cos(halfAngleRad) // f: center to midpoint of a side
	sideLength := 2 * circumradius * math.Sin(halfAngleRad)

	return regularPolygon{
		SideCount:            sideCount,
		CentralAngleDeg:      centralAngle,
		VertexAngleDeg:       180 - centralAngle,
		CircumscribedDiam:    2 * circumradius,
		InscribedDiam:        2 * apothem,
		AcrossFlats:          2 * apothem,
		FlatToOppositeVertex: apothem + circumradius,
		SideLength:           sideLength,
		Perimeter:            float64(sideCount) * sideLength,
		Area:                 0.5 * sideLength * apothem * float64(sideCount),
	}
}

// circumradiusFromSideLength, circumradiusFromAcrossFlats,
// circumradiusFromFlatToOppositeVertex, circumradiusFromDiameter,
// and circumradiusFromInscribedDiameter each invert
// computeRegularPolygon's own formulas to recover the circumradius
// from one known size specification.
func circumradiusFromSideLength(sideCount int, sideLength float64) float64 {
	halfAngleRad := math.Pi / float64(sideCount)
	return 0.5 * sideLength / math.Sin(halfAngleRad)
}

func circumradiusFromAcrossFlats(sideCount int, acrossFlats float64) float64 {
	halfAngleRad := math.Pi / float64(sideCount)
	return acrossFlats / (2 * math.Cos(halfAngleRad))
}

func circumradiusFromFlatToOppositeVertex(sideCount int, flatToOppositeVertex float64) float64 {
	halfAngleRad := math.Pi / float64(sideCount)
	return flatToOppositeVertex / (1 + math.Cos(halfAngleRad))
}

func circumradiusFromInscribedDiameter(sideCount int, inscribedDiam float64) float64 {
	halfAngleRad := math.Pi / float64(sideCount)
	return 0.5 * inscribedDiam / math.Cos(halfAngleRad)
}

// errInsufficientData reports that none of the size specifications
// needed to solve for the polygon's circumradius was given.
var errInsufficientData = errors.New("insufficient information for solution")

func main() {
	prompter, err := promptio.New(os.Stdin, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, "polygon:", err)
		os.Exit(1)
	}

	fmt.Println("PROPERTIES OF REGULAR POLYGONS")
	fmt.Println()

	var sideCount int
	for {
		sideCount = prompter.Int("Number of polygon sides", 6)
		if sideCount >= 3 {
			break
		}
		fmt.Println("\nNOT A POLYGON, TRY AGAIN")
	}
	even := sideCount%2 == 0

	fmt.Println("\nYou must answer at least one of the following four questions.")
	fmt.Println("Press enter if you don't know quantity requested.")
	fmt.Println()

	var circumradius float64
	sideLength := prompter.Float("Length of side", 0)
	switch {
	case sideLength != 0:
		circumradius = circumradiusFromSideLength(sideCount, sideLength)
	case even:
		if acrossFlats := prompter.Float("Size across flats", 0); acrossFlats != 0 {
			circumradius = circumradiusFromAcrossFlats(sideCount, acrossFlats)
		}
	default:
		if flatToVertex := prompter.Float("Size flat-to-opposite-vertex", 0); flatToVertex != 0 {
			circumradius = circumradiusFromFlatToOppositeVertex(sideCount, flatToVertex)
		}
	}

	if circumradius == 0 {
		if circumDiam := prompter.Float("Diameter of circumscribed circle", 0); circumDiam != 0 {
			circumradius = 0.5 * circumDiam
		}
	}
	if circumradius == 0 {
		if inscribedDiam := prompter.Float("Diameter of inscribed circle", 0); inscribedDiam != 0 {
			circumradius = circumradiusFromInscribedDiameter(sideCount, inscribedDiam)
		}
	}
	if circumradius == 0 {
		fmt.Println()
		fmt.Fprintln(os.Stderr, "polygon:", errInsufficientData)
		os.Exit(1)
	}

	p := computeRegularPolygon(sideCount, circumradius)

	fmt.Printf("\nNumber of sides = %d\n", p.SideCount)
	fmt.Printf("Central angle between vertices = %.4f deg\n", p.CentralAngleDeg)
	fmt.Printf("Vertex angle = %.4f deg\n", p.VertexAngleDeg)
	fmt.Printf("Diameter of circumscribed circle = %.4f\n", p.CircumscribedDiam)
	fmt.Printf("Diameter of inscribed circle = %.4f\n", p.InscribedDiam)
	if even {
		fmt.Printf("Distance across flats = %.4f\n", p.AcrossFlats)
	} else {
		fmt.Printf("Distance flat to opposite vertex = %.4f\n", p.FlatToOppositeVertex)
	}
	fmt.Printf("Length of side = %.4f\n", p.SideLength)
	fmt.Printf("Perimeter = %.4f\n", p.Perimeter)
	fmt.Printf("Area = %.4f\n", p.Area)
}
