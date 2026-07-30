// offkey computes the geometry of an offset (rotated) key cut into a
// shaft: given the shaft diameter, key width, key height, and the
// angle the key is rotated from its normal centered position, it
// finds the resulting keyseat and keyway depths and the corner
// positions of the key cross-section before and after rotation.
//
// Converted from OFFKEY.C (M. W. Klotz, 1/03, for Ronnie Shultz),
// WorkshopUtilities/offkey.
package main

import (
	"fmt"
	"math"
	"os"

	"mwkgo/internal/promptio"
)

// point is a 2D coordinate.
type point struct{ X, Y float64 }

// keyGeometry holds the standard (non-rotated) key cutting depths.
type keyGeometry struct {
	KeyseatDepth  float64 // M
	KeywayDepth   float64 // N
	ShaftCutDepth float64 // Q
	HubCutDepth   float64 // R
	MaxAngleDeg   float64
}

// computeKeyGeometry returns the standard keyseat/keyway depths for
// a shaft of diameter shaftDiameter with a key of width keyWidth and
// height keyHeight, and the maximum angle the key could be rotated
// before its corners leave the shaft.
func computeKeyGeometry(shaftDiameter, keyWidth, keyHeight float64) keyGeometry {
	x := math.Sqrt(shaftDiameter*shaftDiameter - keyWidth*keyWidth)
	y := 0.5 * (shaftDiameter - x)
	m := shaftDiameter - (y + 0.5*keyHeight)
	n := m + keyHeight
	q := shaftDiameter - m
	r := keyHeight - q

	return keyGeometry{
		KeyseatDepth:  m,
		KeywayDepth:   n,
		ShaftCutDepth: q,
		HubCutDepth:   r,
		MaxAngleDeg:   2 * math.Asin(keyWidth/shaftDiameter) * 180 / math.Pi,
	}
}

// keyOutline returns the 8 corners of the key cross-section before
// rotation (b) and after rotating by thetaDeg about the origin (a),
// matching the original program's own point layout: index 0 and 7
// are the shaft-cut corners nearest the key's centerline, 1 and 6
// the corners where the shaft cut meets the round shaft profile, 2
// (recomputed as the round profile's own point) mirrored at 3-5 for
// the hub side.
func keyOutline(shaftDiameter, keyWidth, shaftCutDepth, hubCutDepth, thetaDeg float64) (before, after [8]point) {
	halfWidth := 0.5 * keyWidth
	roundY := 0.5 * math.Sqrt(shaftDiameter*shaftDiameter-keyWidth*keyWidth)
	topY := 0.5*shaftDiameter + hubCutDepth
	bottomY := 0.5*shaftDiameter - shaftCutDepth

	before = [8]point{
		{0, bottomY},
		{halfWidth, bottomY},
		{halfWidth, roundY},
		{halfWidth, topY},
		{0, topY},
		{-halfWidth, topY},
		{-halfWidth, roundY},
		{-halfWidth, bottomY},
	}

	cosT := math.Cos(thetaDeg * math.Pi / 180)
	sinT := math.Sin(thetaDeg * math.Pi / 180)
	for i, p := range before {
		if i > 1 && i < 7 {
			after[i] = point{X: cosT*p.X + sinT*p.Y, Y: -sinT*p.X + cosT*p.Y}
		} else {
			after[i] = p
		}
	}
	return before, after
}

// lineIntersection returns the point where the infinite lines
// through (x1,y1)-(x2,y2) and (x3,y3)-(x4,y4) cross.
func lineIntersection(x1, y1, x2, y2, x3, y3, x4, y4 float64) point {
	d := (y2-y1)*(x4-x3) - (y4-y3)*(x2-x1)
	xi := ((x1*y2-x2*y1)*(x4-x3) - (x3*y4-x4*y3)*(x2-x1)) / d
	yi := ((x1*y2-x2*y1)*(y4-y3) - (x3*y4-x4*y3)*(y2-y1)) / d
	return point{X: xi, Y: yi}
}

// separation returns the x, y, and straight-line distance between
// two points.
func separation(p1, p2 point) (dx, dy, dist float64) {
	dx = p1.X - p2.X
	dy = p1.Y - p2.Y
	return dx, dy, math.Hypot(dx, dy)
}

func main() {
	prompter, err := promptio.New(os.Stdin, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, "offkey:", err)
		os.Exit(1)
	}

	fmt.Println("OFFSET KEY CALCULATIONS")
	fmt.Println()
	fmt.Println("Shaft Diameter   Recommended Key Size")
	fmt.Println()
	fmt.Println("0.3125-0.4375    0.0625-0.1250")
	fmt.Println("0.4375-0.5000    0.0625-0.1562")
	fmt.Println("0.5000-0.8750    0.0938-0.3125")
	fmt.Println("0.8750-1.1250    0.1250-0.4375")
	fmt.Println("1.1250-1.1500    0.1562-0.5000")
	fmt.Println("1.5000-2.0000    0.1875-0.5000")
	fmt.Println()

	shaftDiameter := prompter.Float("Shaft diameter", 1.0)
	keyWidth := prompter.Float("Key width", 0.5)
	keyHeight := prompter.Float("Key height", 0.5)

	g := computeKeyGeometry(shaftDiameter, keyWidth, keyHeight)

	fmt.Println()
	fmt.Printf("Shaft diameter (D)         = %.4f\n", shaftDiameter)
	fmt.Printf("Key width (K)              = %.4f\n", keyWidth)
	fmt.Printf("Key height (H)             = %.4f\n", keyHeight)
	fmt.Printf("Keyseat depth (M)          = %.4f\n", g.KeyseatDepth)
	fmt.Printf("Keyway depth (N)           = %.4f\n", g.KeywayDepth)
	fmt.Printf("Depth of cut in shaft (Q)  = %.4f\n", g.ShaftCutDepth)
	fmt.Printf("Depth of cut in hub (R)    = %.4f\n", g.HubCutDepth)
	fmt.Println()
	fmt.Println("Keyseat depth measured bottom of keyseat to opposite side of shaft on cl.")
	fmt.Println("Keyway depth measured top of hub keyway to opposide side of hole on cl.")
	fmt.Println()

	shaftCutDepth := prompter.Float("Depth of cut in shaft (Q)", g.ShaftCutDepth)
	hubCutDepth := keyHeight - shaftCutDepth

	var theta float64
	for {
		fmt.Printf("Maximum possible rotation angle = %.4f deg\n", g.MaxAngleDeg)
		theta = math.Abs(prompter.Float("Key rotation angle (theta)", 5.0))
		if theta <= g.MaxAngleDeg {
			break
		}
		fmt.Println("ANGLE TOO BIG.  TRY AGAIN.")
	}

	before, after := keyOutline(shaftDiameter, keyWidth, shaftCutDepth, hubCutDepth, theta)

	fmt.Println()
	fmt.Println("x,y locations of points (b)efore => (a)fter rotation")
	fmt.Println()
	for i := 0; i < 8; i++ {
		fmt.Printf("b[%d] = %.4f, %.4f => a[%d] = %.4f, %.4f\n", i, before[i].X, before[i].Y, i, after[i].X, after[i].Y)
	}

	point8 := lineIntersection(before[6].X, before[6].Y, before[2].X, before[2].Y, after[6].X, after[6].Y, after[5].X, after[5].Y)
	point9 := lineIntersection(before[6].X, before[6].Y, before[2].X, before[2].Y, after[2].X, after[2].Y, after[3].X, after[3].Y)

	fmt.Printf("a[8] = %.4f, %.4f\n", point8.X, point8.Y)
	fmt.Printf("a[9] = %.4f, %.4f\n", point9.X, point9.Y)

	fmt.Println("\nrelative locations")
	fmt.Println()
	printRel := func(aIdx string, ap point, bLabel string, bp point) {
		dx, dy, dist := separation(ap, bp)
		fmt.Printf("a[%s] wrt b[%s]:  x sep = %.4f, y sep = %.4f, sep = %.4f\n", aIdx, bLabel, dx, dy, dist)
	}
	printRel("6", after[6], "6", before[6])
	printRel("2", after[2], "2", before[2])
	printRel("3", after[3], "7", before[7])
	printRel("5", after[5], "7", before[7])
	printRel("8", point8, "6", before[6])
	printRel("9", point9, "2", before[2])

	dx, dy, dist := separation(after[5], point8)
	fmt.Printf("a[5] wrt a[8]:  x sep = %.4f, y sep = %.4f, sep = %.4f\n", dx, dy, dist)
}
