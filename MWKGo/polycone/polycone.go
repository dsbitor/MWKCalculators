// polycone computes the geometry of a "polycone": a solid whose base
// is a regular polygon and whose sides are triangular facets all
// meeting at a single apex point above the base's center (a
// tetrahedron is a three-sided polycone; a pyramid is a four-sided
// one). Given the number of base sides, the length of one base side,
// and the polycone's height, this program computes the base's
// dimensions, each triangular facet's dimensions and angles, and the
// solid's total surface area and volume.
//
// Converted from POLYCONE.C, WorkshopUtilities/polycone.
package main

import (
	"fmt"
	"math"
	"os"

	"mwkgo/internal/promptio"
)

// polyconeBase holds the dimensions of a polycone's regular-polygon
// base.
type polyconeBase struct {
	Diameter float64 // of the circumscribed circle
	Radius   float64 // of the circumscribed circle
	Apothem  float64 // perpendicular distance from center to a side
	Area     float64
}

// polyconeFacet holds the dimensions and angles of one of a
// polycone's triangular side facets, all of which are congruent.
type polyconeFacet struct {
	Height           float64 // from base edge midpoint to apex
	EdgeLength       float64 // from a base vertex to apex
	BaseAngleDeg     float64 // angle of a facet's base corner
	TipAngleDeg      float64 // angle at the apex
	FaceAngleDeg     float64 // angle of the facet plane above horizontal
	EdgeAngleDeg     float64 // angle of a facet edge above horizontal
	ExteriorAngleDeg float64 // between the normals of two adjacent facets
	InteriorAngleDeg float64 // dihedral angle between two adjacent facets
	Area             float64
}

// polycone holds the full computed geometry of a polycone.
type polycone struct {
	Base        polyconeBase
	Facet       polyconeFacet
	TotalArea   float64
	TotalVolume float64
}

// vectorAngleDeg returns the angle, in degrees, between two 3D
// vectors.
func vectorAngleDeg(v1, v2 [3]float64) float64 {
	dot := v1[0]*v2[0] + v1[1]*v2[1] + v1[2]*v2[2]
	mag1 := math.Sqrt(v1[0]*v1[0] + v1[1]*v1[1] + v1[2]*v1[2])
	mag2 := math.Sqrt(v2[0]*v2[0] + v2[1]*v2[1] + v2[2]*v2[2])
	return math.Acos(dot/(mag1*mag2)) * 180 / math.Pi
}

// computePolycone returns the geometry of a polycone with sides
// base sides, each of length sideLength, and a perpendicular height
// of height from the base's center to the apex.
func computePolycone(sides int, sideLength, height float64) polycone {
	halfChord := 0.5 * sideLength
	halfAngleDeg := 180.0 / float64(sides)
	halfAngleRad := halfAngleDeg * math.Pi / 180

	radius := halfChord / math.Sin(halfAngleRad)
	apothem := radius * math.Cos(halfAngleRad)
	baseArea := float64(sides) * apothem * halfChord

	facetHeight := math.Hypot(height, apothem)
	edgeLength := math.Hypot(height, radius)
	baseAngle := math.Acos(halfChord/edgeLength) * 180 / math.Pi
	tipAngle := 2 * (90 - baseAngle)
	faceAngle := math.Acos(apothem/facetHeight) * 180 / math.Pi
	edgeAngle := math.Acos(radius/edgeLength) * 180 / math.Pi
	facetArea := 0.5 * facetHeight * halfChord

	// Exterior angle between the normals of two adjacent facets: v1
	// is a vector in the plane of one facet, perpendicular to its
	// base edge; v2 is the same vector rotated by the base's full
	// central angle (2*halfAngleDeg) to lie in the adjacent facet.
	faceAngleRad := faceAngle * math.Pi / 180
	sinFace, cosFace := math.Sin(faceAngleRad), math.Cos(faceAngleRad)
	v1 := [3]float64{apothem * sinFace * sinFace, 0, apothem * sinFace * cosFace}
	twoPhiRad := 2 * halfAngleRad
	sp, cp := math.Sin(twoPhiRad), math.Cos(twoPhiRad)
	v2 := [3]float64{cp*v1[0] - sp*v1[1], sp*v1[0] + cp*v1[1], v1[2]}
	exteriorAngle := vectorAngleDeg(v1, v2)

	return polycone{
		Base: polyconeBase{
			Diameter: 2 * radius,
			Radius:   radius,
			Apothem:  apothem,
			Area:     baseArea,
		},
		Facet: polyconeFacet{
			Height:           facetHeight,
			EdgeLength:       edgeLength,
			BaseAngleDeg:     baseAngle,
			TipAngleDeg:      tipAngle,
			FaceAngleDeg:     faceAngle,
			EdgeAngleDeg:     edgeAngle,
			ExteriorAngleDeg: exteriorAngle,
			InteriorAngleDeg: 180 - exteriorAngle,
			Area:             facetArea,
		},
		TotalArea:   baseArea + float64(sides)*facetArea,
		TotalVolume: baseArea * height / 3,
	}
}

func main() {
	prompter, err := promptio.New(os.Stdin, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, "polycone:", err)
		os.Exit(1)
	}

	sides := prompter.Int("Number of polygon sides", 4)
	sideLength := prompter.Float("Length of polygon side", 6.0)
	height := prompter.Float("Cone height", 10.0)

	p := computePolycone(sides, sideLength, height)

	fmt.Println("\nPOLYCONE BASE:")
	fmt.Printf("Diameter / radius of circumscribed circle = %.4f / %.4f\n", p.Base.Diameter, p.Base.Radius)
	fmt.Printf("Length of perpendicular to side from center = %.4f\n", p.Base.Apothem)
	fmt.Printf("Area of base = %.4f\n", p.Base.Area)

	fmt.Println("\nPOLYCONE TRIANGULAR FACE FACETS:")
	fmt.Printf("Facet height = %.4f\n", p.Facet.Height)
	fmt.Printf("Edge length = %.4f\n", p.Facet.EdgeLength)
	fmt.Printf("Base angle = %.4f deg\n", p.Facet.BaseAngleDeg)
	fmt.Printf("Tip angle = %.4f deg\n", p.Facet.TipAngleDeg)
	fmt.Printf("Angle of face wrt horizontal = %.4f deg\n", p.Facet.FaceAngleDeg)
	fmt.Printf("Angle of edge wrt horizontal = %.4f deg\n", p.Facet.EdgeAngleDeg)
	fmt.Printf("Exterior angle between perpendiculars to two adjacent faces = %.4f deg\n", p.Facet.ExteriorAngleDeg)
	fmt.Printf("Interior angle between two adjacent faces = %.4f deg\n", p.Facet.InteriorAngleDeg)
	fmt.Printf("Area of single facet = %.4f\n", p.Facet.Area)

	fmt.Println("\nCOMPLETE POLYCONE:")
	fmt.Printf("Total area = %.4f\n", p.TotalArea)
	fmt.Printf("Total volume = %.4f\n", p.TotalVolume)
}
