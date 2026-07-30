// weight computes the weight of a shape from its dimensions and a
// chosen material's density, for any of sixteen standard shapes (or a
// directly-entered volume), maintaining a running total across
// however many shapes make up one object.
//
// The material density table is universal reference data (the same
// approximate densities for anyone), so it lives in the shared
// reference.db SQLite database, the same mechanism used by fits,
// speed, gage, expand, and findthrd earlier in this project; see
// ai/plans/c-to-go-conversion-plan.md, "Data-file strategy for
// Tier 2".
//
// Converted from WEIGHT.C (M. W. Klotz), WorkshopUtilities/weight.
package main

import (
	"context"
	"database/sql"
	"fmt"
	"math"
	"os"

	"mwkgo/internal/promptio"
	"mwkgo/internal/refdata"
)

// material is one row of the density reference table.
type material struct {
	Name            string
	DensityLbPerIn3 float64
}

// loadMaterials returns every row of the weight_materials reference
// table, alphabetically (case-insensitively) — WEIGHT.DAT's own rows
// are already in this order, and WEIGHT.C's rdata() sorts them the
// same way before use.
func loadMaterials(ctx context.Context, db *sql.DB) ([]material, error) {
	rows, err := db.QueryContext(ctx, `SELECT name, density_lb_per_in3 FROM weight_materials ORDER BY name COLLATE NOCASE`)
	if err != nil {
		return nil, fmt.Errorf("query weight_materials: %w", err)
	}
	defer rows.Close()

	var materials []material
	for rows.Next() {
		var m material
		if err := rows.Scan(&m.Name, &m.DensityLbPerIn3); err != nil {
			return nil, fmt.Errorf("scan material: %w", err)
		}
		materials = append(materials, m)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate weight_materials: %w", err)
	}
	return materials, nil
}

// promptForMaterial lists every material by number and prompts for a
// choice; a number outside the list (including the offered "one past
// the end" option) instead asks for a custom density, matching
// WEIGHT.C's own mselect(): the material's name is then reported as
// "??" as the original does, since it has no name of its own.
func promptForMaterial(p *promptio.Prompter, materials []material) material {
	for i, m := range materials {
		fmt.Printf("%s-%d  ", m.Name, i+1)
	}
	customOption := len(materials) + 1
	fmt.Printf("\nUser input-%d\n", customOption)

	n := p.Int("\nMaterial number", customOption)
	if n < 1 || n > len(materials) {
		density := p.Float("Material density (lb/in^3)", 1.0)
		return material{Name: "??", DensityLbPerIn3: density}
	}
	return materials[n-1]
}

// degToRad converts degrees to radians.
func degToRad(deg float64) float64 { return deg * math.Pi / 180 }

// Shape volume functions, one per WEIGHT.C shape option, using
// whatever unit the caller's dimensions are in (the result is in that
// same unit, cubed).

func rectangularVolume(length, width, height float64) float64 {
	return length * width * height
}

func cylinderVolume(diameter, length float64) float64 {
	return math.Pi * diameter * diameter * length / 4
}

func pipeVolume(outsideDiameter, insideDiameter, length float64) float64 {
	return math.Pi * (outsideDiameter*outsideDiameter - insideDiameter*insideDiameter) * length / 4
}

func hexagonalVolume(acrossFlats, length float64) float64 {
	return 1.5 * acrossFlats * acrossFlats * length * math.Tan(degToRad(30))
}

func pyramidVolume(baseSide, altitude float64) float64 {
	return baseSide * baseSide * altitude / 3
}

func coneVolume(baseDiameter, altitude float64) float64 {
	return math.Pi * baseDiameter * baseDiameter * altitude / 12
}

func frustumVolume(largerDiameter, smallerDiameter, altitude float64) float64 {
	return math.Pi * (largerDiameter*largerDiameter + largerDiameter*smallerDiameter + smallerDiameter*smallerDiameter) * altitude / 12
}

func sphereVolume(diameter float64) float64 {
	return math.Pi * diameter * diameter * diameter / 6
}

func annulusVolume(outerDiameter, holeDiameter, thickness float64) float64 {
	return 0.25 * math.Pi * thickness * (outerDiameter*outerDiameter - holeDiameter*holeDiameter)
}

// polygonVolume is the volume of a regular-polygon-cross-section prism
// of n sides and length h. For an even number of sides, size is the
// distance across flats; for an odd number, size is the distance from
// a flat side to the opposite vertex — the two ways WEIGHT.C itself
// distinguishes how the polygon's size was measured.
func polygonVolume(sides int, size, length float64) (float64, error) {
	if sides < 3 {
		return 0, fmt.Errorf("%d is not a polygon (need at least 3 sides)", sides)
	}
	halfAngle := degToRad(180.0 / float64(sides))
	ca, sa := math.Cos(halfAngle), math.Sin(halfAngle)

	var r float64
	if sides%2 == 0 {
		r = size / (2 * ca)
	} else {
		r = size / (1 + ca)
	}
	f := r * ca
	c := 2 * r * sa
	return 0.5 * c * f * float64(sides) * length, nil
}

func sphericalCapVolume(sphereDiameter, capHeight float64) float64 {
	return math.Pi * capHeight * capHeight * (1.5*sphereDiameter - capHeight) / 3
}

func intersectingCylindersVolume(diameter float64) float64 {
	return 16 * diameter * diameter * diameter / 24
}

// torusVolume is the volume of a torus with the given cross-section
// diameter and overall outside diameter, by Pappus's theorem
// (2*pi*centerlineRadius*crossSectionArea), matching WEIGHT.C's own
// formula exactly.
func torusVolume(crossSectionDiameter, outsideDiameter float64) float64 {
	centerlineDiameter := outsideDiameter - crossSectionDiameter
	return 0.25 * math.Pi * math.Pi * centerlineDiameter * crossSectionDiameter * crossSectionDiameter
}

func hemisphericalShellVolume(outsideDiameter, insideDiameter float64) float64 {
	return (math.Pi / 12) * (outsideDiameter*outsideDiameter*outsideDiameter - insideDiameter*insideDiameter*insideDiameter)
}

// tangentOgiveVolume is the volume of a tangent ogive (the pointed,
// aerodynamic nose shape used on many projectiles) of base diameter d
// and length l, matching WEIGHT.C's own formula exactly.
func tangentOgiveVolume(baseDiameter, length float64) float64 {
	c := length / baseDiameter
	pp := baseDiameter * (c*c + 0.25)
	pm := baseDiameter * (c*c - 0.25)
	pp2, pm2 := pp*pp, pm*pm
	return math.Pi * ((pp2+pm2)*length - length*length*length/3 -
		pm*(length*math.Sqrt(pp2-length*length)+pp2*math.Asin(length/pp)))
}

// conicalWedgeVolume is the volume of a cone of height h and base
// radius r sliced by a plane leaving a base sagitta of s (s < r cuts
// off the smaller piece past the sagitta; s > r the larger piece;
// s == r bisects the cone exactly), matching WEIGHT.C's own wedge()
// function exactly.
func conicalWedgeVolume(height, sagitta, radius float64) float64 {
	fullCone := math.Pi * radius * radius * height / 3
	if sagitta == radius {
		return 0.5 * fullCone
	}

	var p float64
	if sagitta < radius {
		p = radius - sagitta
	} else {
		p = sagitta - radius
	}
	ct := p / radius
	t := math.Acos(ct)
	st := math.Sin(t)
	tt := math.Tan(t)
	smaller := height * radius * radius * (t - 2*st*ct + ct*ct*ct*math.Log(tt+1/ct)) / 3

	if sagitta <= radius {
		return smaller
	}
	return fullCone - smaller
}

// unitConversion returns the display name and the volume-conversion
// factor (multiplying a volume measured in unit^3 by this factor gives
// in^3) for one of WEIGHT.C's four supported dimension units.
func unitConversion(choice string) (name string, volumeFactor float64) {
	switch choice {
	case "f":
		return "ft", 1728
	case "m":
		return "mm", 1 / math.Pow(25.4, 3)
	case "c":
		return "cm", 1 / math.Pow(2.54, 3)
	default:
		return "in", 1
	}
}

func main() {
	ctx := context.Background()
	db, err := refdata.Open(ctx)
	if err != nil {
		fmt.Fprintln(os.Stderr, "weight:", err)
		os.Exit(1)
	}
	defer db.Close()

	materials, err := loadMaterials(ctx, db)
	if err != nil {
		fmt.Fprintln(os.Stderr, "weight:", err)
		os.Exit(1)
	}

	prompter, err := promptio.New(os.Stdin, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, "weight:", err)
		os.Exit(1)
	}

	fmt.Println("WEIGHT CALCULATIONS")
	fmt.Println()
	fmt.Printf("Number of data items read = %d\n", len(materials))

	var totalVolume, totalWeight float64

materialLoop:
	for {
		chosen := promptForMaterial(prompter, materials)
		fmt.Printf("\nDensity of %s = %.4f lb/in^3 = %.4f gm/cc\n",
			chosen.Name, chosen.DensityLbPerIn3, chosen.DensityLbPerIn3/0.036127)

		unitChoice := prompter.Line("Dimensions of shapes in [i]n, (f)t, (m)m, (c)m ? ")
		unitName, volumeFactor := unitConversion(unitChoice)

		for {
			printShapeMenu()
			shapeID := prompter.Int("Shape number", 1)

			var volume float64
			switch shapeID {
			case -1:
				continue materialLoop
			case 0:
				return
			case 1:
				l := prompter.Float("Length ("+unitName+")", 12)
				w := prompter.Float("Width ("+unitName+")", 12)
				h := prompter.Float("Height ("+unitName+")", 12)
				volume = rectangularVolume(l, w, h)
			case 2:
				d := prompter.Float("Diameter ("+unitName+")", 1)
				h := prompter.Float("Length ("+unitName+")", 12)
				volume = cylinderVolume(d, h)
			case 3:
				d1 := prompter.Float("Outside Diameter ("+unitName+")", 1.5)
				d2 := prompter.Float("Inside Diameter ("+unitName+")", 1)
				h := prompter.Float("Length ("+unitName+")", 12)
				volume = pipeVolume(d1, d2, h)
			case 4:
				d := prompter.Float("Size across flats ("+unitName+")", 1)
				h := prompter.Float("Length ("+unitName+")", 12)
				volume = hexagonalVolume(d, h)
			case 5:
				d := prompter.Float("Length of base side ("+unitName+")", 1)
				h := prompter.Float("Altitude ("+unitName+")", 12)
				volume = pyramidVolume(d, h)
			case 6:
				d := prompter.Float("Diameter of base ("+unitName+")", 1)
				h := prompter.Float("Altitude ("+unitName+")", 12)
				volume = coneVolume(d, h)
			case 7:
				d1 := prompter.Float("Larger Diameter ("+unitName+")", 2)
				d2 := prompter.Float("Smaller Diameter ("+unitName+")", 1)
				h := prompter.Float("Altitude ("+unitName+")", 12)
				volume = frustumVolume(d1, d2, h)
			case 8:
				d := prompter.Float("Diameter ("+unitName+")", 1)
				volume = sphereVolume(d)
			case 9:
				d1 := prompter.Float("Outer Diameter ("+unitName+")", 1)
				d2 := prompter.Float("Hole Diameter ("+unitName+")", 0.5)
				h := prompter.Float("Thickness ("+unitName+")", 0.1)
				volume = annulusVolume(d1, d2, h)
			case 10:
				var n int
				for {
					n = prompter.Int("Number of polygon sides", 6)
					if n >= 3 {
						break
					}
					fmt.Println("\nNOT A POLYGON, TRY AGAIN")
				}
				var size float64
				if n%2 == 0 {
					size = prompter.Float("Size across flats", 1)
				} else {
					size = prompter.Float("Size flat-to-opposite-vertex", 1)
				}
				h := prompter.Float("Length ("+unitName+")", 12)
				v, err := polygonVolume(n, size, h)
				if err != nil {
					fmt.Println(err)
					continue
				}
				volume = v
			case 11:
				d := prompter.Float("Diameter of sphere ("+unitName+")", 1)
				h := prompter.Float("Height of cap at center ("+unitName+")", 0.5)
				volume = sphericalCapVolume(d, h)
			case 12:
				d := prompter.Float("Diameter of cylinder(s) ("+unitName+")", 1)
				volume = intersectingCylindersVolume(d)
			case 13:
				d := prompter.Float("Diameter of torus cross section ("+unitName+")", 1.5)
				d1 := prompter.Float("Outside Diameter of torus ("+unitName+")", 6.5)
				volume = torusVolume(d, d1)
			case 14:
				d2 := prompter.Float("Outside diameter of shell ("+unitName+")", 3)
				d1 := prompter.Float("Inside diameter of shell ("+unitName+")", 1)
				volume = hemisphericalShellVolume(d2, d1)
			case 15:
				d := prompter.Float("Diameter of ogive base ("+unitName+")", 2)
				l := prompter.Float("Length/height of ogive ("+unitName+")", 5)
				volume = tangentOgiveVolume(d, l)
			case 16:
				d := prompter.Float("Cone base diameter", 2)
				h := prompter.Float("Height of cone", 10)
				s := prompter.Float("Sagitta of wedge base", 0.5)
				volume = conicalWedgeVolume(h, s, 0.5*d)
			case 99:
				volume = prompter.Float("Volume of shape (in^3)", 1728)
				volumeFactor = 1 // shape 99's volume is entered directly in in^3
			default:
				fmt.Println("NOT A VALID OPTION")
				continue
			}

			weight := chosen.DensityLbPerIn3 * volume * volumeFactor
			fmt.Printf("\nVolume of shape = %.4f %s^3\n", volume, unitName)
			fmt.Printf("Weight of shape = %.4g lb = %.4g kg\n", weight, weight/2.20462)

			multiplier := prompter.Float("Number of these shapes to add to total", 1)
			totalVolume += multiplier * volume
			totalWeight += multiplier * weight
			fmt.Printf("Total volume = %.4f %s^3\n", totalVolume, unitName)
			fmt.Printf("Total weight = %.4g lb = %.4g kg\n", totalWeight, totalWeight/2.20462)
		}
	}
}

func printShapeMenu() {
	fmt.Println("\n-1 select new material")
	fmt.Println(" 0 QUIT")
	fmt.Println(" 1 Rectangular solid")
	fmt.Println(" 2 Cylindrical solid")
	fmt.Println(" 3 Cylindrical pipe")
	fmt.Println(" 4 Hexagonal solid")
	fmt.Println(" 5 Pyramid")
	fmt.Println(" 6 Cone")
	fmt.Println(" 7 Frustum of a cone")
	fmt.Println(" 8 Sphere")
	fmt.Println(" 9 Annulus (washer shape)")
	fmt.Println("10 Regular polygon cross section")
	fmt.Println("11 Spherical cap")
	fmt.Println("12 Volume formed by two intersecting cylinders")
	fmt.Println("13 Torus")
	fmt.Println("14 Hemispherical shell")
	fmt.Println("15 Tangent ogive")
	fmt.Println("16 Conical wedge")
	fmt.Println("99 User Volume Input")
}
