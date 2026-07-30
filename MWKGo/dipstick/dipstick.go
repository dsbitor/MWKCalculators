// dipstick computes what fraction of a tank's full capacity remains,
// given a dipstick reading (the wetted length) and the tank's shape
// and dimensions, for nine tank shapes. It also calibrates a
// dipstick by finding the wetted length corresponding to each of a
// series of percentage increments of full capacity.
//
// Converted from DIPSTICK.C (M. W. Klotz, 4/01, 9/03, 4,7/04,
// 7,9/05), WorkshopUtilities/dipstick.
package main

import (
	"errors"
	"fmt"
	"math"
	"os"
	"strings"

	"mwkgo/internal/promptio"
)

// segmentArea returns the area of the circular segment of a circle
// of radius r wetted to depth h (measured from the bottom), matching
// the shared segment-area idiom used by several of the shape
// functions below. Angles here are in radians, not degrees, matching
// the original program's own use of plain acos/sin rather than the
// ACSD/SIND macros used elsewhere in this codebase.
func segmentArea(r, h float64) float64 {
	z := r - h
	angle := 2 * math.Acos(clamp(z/r))
	chord := 2 * r * math.Sin(0.5*angle)
	return 0.5 * (r*r*angle - z*chord)
}

func clamp(x float64) float64 {
	if x > 1 {
		return 1
	}
	if x < -1 {
		return -1
	}
	return x
}

// cylFraction returns the wetted volume fraction of a horizontal
// cylindrical tank of diameter d, wetted to depth h.
func cylFraction(d, h float64) float64 {
	r := 0.5 * d
	return segmentArea(r, h) / (math.Pi * r * r)
}

// spherFraction returns the wetted volume fraction of a spherical
// tank of diameter d, wetted to depth h.
func spherFraction(d, h float64) float64 {
	r := 0.5 * d
	return h * h * (3*r - h) / (4 * r * r * r)
}

// ellipFraction returns the wetted volume fraction of a horizontal
// tank with an elliptical cross section (horizontal dimension
// horiz, vertical dimension vert), wetted to depth h.
func ellipFraction(horiz, vert, h float64) float64 {
	a, b := 0.5*horiz, 0.5*vert
	area := (a / b) * ((h-b)*math.Sqrt(h*(2*b-h)) + b*b*(math.Asin(clamp((h-b)/b))+0.5*math.Pi))
	return area / (math.Pi * a * b)
}

// vcartFraction returns the wetted volume fraction of a vertical
// cartouche tank (a "stadium" cross section: two semicircular ends
// separated by a straight section, oriented with the straight
// section vertical), given its horizontal and vertical dimensions,
// wetted to depth h.
func vcartFraction(horiz, vert, h float64) float64 {
	r := 0.5 * horiz
	l := vert - horiz
	totalArea := math.Pi*r*r + l*horiz

	var area float64
	switch {
	case h <= r:
		area = segmentArea(r, h)
	case h <= r+l:
		area = 0.5*math.Pi*r*r + horiz*(h-r)
	default:
		z := r - (h - l)
		angle := 2 * math.Acos(clamp(z/r))
		chord := 2 * r * math.Sin(0.5*angle)
		area = 0.5*(r*r*angle-z*chord) + l*horiz
	}
	return area / totalArea
}

// hcartFraction returns the wetted volume fraction of a horizontal
// cartouche tank (straight section horizontal), given its horizontal
// and vertical dimensions, wetted to depth h.
func hcartFraction(horiz, vert, h float64) float64 {
	r := 0.5 * vert
	l := horiz - vert
	totalArea := math.Pi*r*r + l*vert
	area := segmentArea(r, h) + h*l
	return area / totalArea
}

// buckFraction returns the wetted volume fraction of a bucket-shaped
// (conical frustum) tank with the big end up, given the diameters of
// its big and small ends and its height, wetted to depth h.
func buckFraction(bigDiameter, smallDiameter, height, h float64) float64 {
	rb, rs := 0.5*bigDiameter, 0.5*smallDiameter
	totalVolume := math.Pi * height * (rs*rs + rs*rb + rb*rb) / 3

	r := rs + (rb-rs)*h/height
	z := rs*rs + rs*r + r*r
	volume := math.Pi * h * z / 3
	return volume / totalVolume
}

// barrelFraction returns the wetted volume fraction of a
// barrel-shaped tank (sides modeled as sections of an ellipse),
// given its diameter at the middle and at the ends and its height,
// wetted to depth h.
func barrelFraction(midDiameter, endDiameter, height, h float64) float64 {
	rb, rs := 0.5*midDiameter, 0.5*endDiameter
	totalVolume := math.Pi * height * (2*rb*rb + rs*rs) / 3

	volume := math.Pi * (rb*rb*h + (rs*rs-rb*rb)*(height*height*h-2*height*h*h+4*h*h*h/3)/(height*height))
	return volume / totalVolume
}

// hcylHemisphereFraction returns the wetted volume fraction of a
// horizontal cylindrical tank with hemispherical ends, given its
// diameter and overall length, wetted to depth h.
func hcylHemisphereFraction(diameter, length, h float64) float64 {
	r := 0.5 * diameter
	totalVolume := math.Pi * r * r * (length - diameter + 4*r/3)

	area := segmentArea(r, h)
	cylindricalVolume := area * (length - diameter)
	hemisphereVolume := math.Pi * h * h * (r - h/3)
	return (cylindricalVolume + hemisphereVolume) / totalVolume
}

// maxDishIntegrationSteps bounds the numerical integration in
// dishedEndVolume, matching the original program's own fixed
// 0.1%-of-dish-height step size (1000 steps total).
const maxDishIntegrationSteps = 1000

// dishedEndVolume returns the wetted volume of one torispherical
// ("dished") tank end, of radius dishRadius and depth dishHeight,
// wetted to depth h (measured from the bottom of the whole tank),
// via numerical integration along the dish's profile, matching the
// original program's own vdish function. tankRadius is the overall
// tank radius, used (along with h) for the "more than half full"
// mirror-image case.
func dishedEndVolume(h, dishRadius, dishHeight, tankRadius float64) float64 {
	var wettedHeight float64
	mirrored := h > dishRadius
	if mirrored {
		wettedHeight = 2*tankRadius - h
	} else {
		wettedHeight = h
	}

	totalDishVolume := math.Pi * dishHeight * (3*dishRadius*dishRadius + dishHeight*dishHeight) / 6
	profileRadius := 0.5 * (dishRadius*dishRadius + dishHeight*dishHeight) / dishHeight
	z := dishRadius - wettedHeight
	step := 0.001 * dishHeight

	var wettedVolume float64
	for i := 0; i < maxDishIntegrationSteps; i++ {
		rho := float64(i) * step
		if rho > dishHeight {
			break
		}
		localRadius := math.Sqrt(math.Max(0, profileRadius*profileRadius-(profileRadius-dishHeight+rho)*(profileRadius-dishHeight+rho)))
		if math.Abs(z) > localRadius {
			break
		}
		angle := 2 * math.Acos(clamp(z/localRadius))
		chord := 2 * localRadius * math.Sin(0.5*angle)
		area := 0.5 * (localRadius*localRadius*angle - z*chord)
		wettedVolume += area * step
	}

	if mirrored {
		return totalDishVolume - wettedVolume
	}
	return wettedVolume
}

// dcylDishedFraction returns the wetted volume fraction of a
// horizontal cylindrical tank with dished ends, given its diameter,
// maximum overall length (over the dished ends), and the length of
// its purely cylindrical portion, wetted to depth h.
func dcylDishedFraction(diameter, overallLength, cylindricalLength, h float64) float64 {
	r := 0.5 * diameter
	dishHeight := 0.5 * (overallLength - cylindricalLength)
	dishVolume := math.Pi * dishHeight * (3*r*r + dishHeight*dishHeight) / 6
	totalVolume := math.Pi*r*r*cylindricalLength + 2*dishVolume

	area := segmentArea(r, h)
	volume := area*cylindricalLength + 2*dishedEndVolume(h, r, dishHeight, r)
	return volume / totalVolume
}

// maxBinarySearchSteps bounds the dipstick-calibration search,
// matching the original program's own runaway guard.
const maxBinarySearchSteps = 500

var errSearchDidNotConverge = errors.New("search did not converge within the step limit")

// findWettedHeight finds the wetted height in [0, maxHeight] at
// which fractionAt(height) equals targetFraction, to within
// tolerance, by binary search (fractionAt must be non-decreasing
// over the search range).
func findWettedHeight(maxHeight, targetFraction, tolerance float64, fractionAt func(float64) float64) (float64, error) {
	low, high := 0.0, maxHeight
	for i := 0; i < maxBinarySearchSteps; i++ {
		mid := low + 0.5*(high-low)
		y := fractionAt(mid)
		if math.Abs(y-targetFraction) <= tolerance {
			return mid, nil
		}
		if targetFraction < y {
			high = mid
		} else {
			low = mid
		}
	}
	return 0, errSearchDidNotConverge
}

func main() {
	prompter, err := promptio.New(os.Stdin, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, "dipstick:", err)
		os.Exit(1)
	}

	fmt.Println("TANK VOLUME FRACTION FROM DIPSTICK READING")
	fmt.Println()
	fmt.Println("A - horizontal cylindrical tank")
	fmt.Println("B - spherical tank")
	fmt.Println("C - elliptical tank")
	fmt.Println("D - vertical cartouche")
	fmt.Println("E - horizontal cartouche")
	fmt.Println("F - bucket (conic frustum)")
	fmt.Println("G - barrel (elliptical sides)")
	fmt.Println("H - horizontal cylindrical tank with hemispherical ends")
	fmt.Println("I - horizontal cylindrical tank with dished ends")
	fmt.Println()

	var shape string
	var maxHeight float64
	var fractionAt func(float64) float64

	for {
		shape = strings.ToLower(prompter.Line("(A-I) ? "))
		ok := true
		switch shape {
		case "a":
			d := prompter.Float("Tank diameter", 10.0)
			maxHeight = d
			fractionAt = func(h float64) float64 { return cylFraction(d, h) }
		case "b":
			d := prompter.Float("Tank diameter", 10.0)
			maxHeight = d
			fractionAt = func(h float64) float64 { return spherFraction(d, h) }
			fmt.Printf("Total volume of tank = %.3f\n", math.Pi*d*d*d/6)
		case "c":
			horiz := prompter.Float("Tank horizontal dimension", 10.0)
			vert := prompter.Float("Tank vertical dimension", 7.0)
			maxHeight = vert
			fractionAt = func(h float64) float64 { return ellipFraction(horiz, vert, h) }
		case "d":
			horiz := prompter.Float("Tank horizontal dimension", 4.0)
			vert := prompter.Float("Tank vertical dimension", 10.0)
			maxHeight = vert
			fractionAt = func(h float64) float64 { return vcartFraction(horiz, vert, h) }
		case "e":
			horiz := prompter.Float("Tank horizontal dimension", 10.0)
			vert := prompter.Float("Tank vertical dimension", 4.0)
			maxHeight = vert
			fractionAt = func(h float64) float64 { return hcartFraction(horiz, vert, h) }
		case "f":
			big := prompter.Float("Diameter of bucket big end", 4.0)
			small := prompter.Float("Diameter of bucket small end", 3.0)
			height := prompter.Float("Vertical height of bucket", 6.0)
			maxHeight = height
			fractionAt = func(h float64) float64 { return buckFraction(big, small, height, h) }
			rb, rs := 0.5*big, 0.5*small
			fmt.Printf("Total volume of bucket = %.3f\n", math.Pi*height*(rs*rs+rs*rb+rb*rb)/3)
		case "g":
			mid := prompter.Float("Diameter of barrel at middle", 8.0)
			end := prompter.Float("Diameter of barrel at end(s)", 6.0)
			height := prompter.Float("Vertical height of barrel", 5.0)
			maxHeight = height
			fractionAt = func(h float64) float64 { return barrelFraction(mid, end, height, h) }
			rb, rs := 0.5*mid, 0.5*end
			fmt.Printf("Total volume of barrel = %.3f\n", math.Pi*height*(2*rb*rb+rs*rs)/3)
		case "h":
			d := prompter.Float("Tank diameter", 10.0)
			length := prompter.Float("Tank length", 20.0)
			maxHeight = d
			fractionAt = func(h float64) float64 { return hcylHemisphereFraction(d, length, h) }
			r := 0.5 * d
			fmt.Printf("Total volume of tank = %.3f\n", math.Pi*r*r*(length-d+4*r/3))
		case "i":
			d := prompter.Float("Tank diameter", 100.0)
			overall := prompter.Float("(Maximum) Tank length over ends", 1000.0)
			cylindrical := prompter.Float("Length of cylindrical portion of tank", 800.0)
			maxHeight = d
			fractionAt = func(h float64) float64 { return dcylDishedFraction(d, overall, cylindrical, h) }
			r := 0.5 * d
			dishHeight := 0.5 * (overall - cylindrical)
			dishVolume := math.Pi * dishHeight * (3*r*r + dishHeight*dishHeight) / 6
			fmt.Printf("Total volume of tank = %.3f\n", math.Pi*r*r*cylindrical+2*dishVolume)
		default:
			fmt.Println("NOT A VALID OPTION")
			ok = false
		}
		if ok {
			break
		}
	}

	var h float64
	for {
		h = prompter.Float("Wetted dipstick length", maxHeight/3)
		if h >= 0 && h <= maxHeight {
			break
		}
		fmt.Println("INPUT MAKES NO SENSE. Try again.")
	}

	f := fractionAt(h)
	fmt.Printf("\nTank has %.2f%% of its full capacity remaining\n", 100*f)

	fmt.Println("\nCalibration of dip stick as function of tank volume.")
	fmt.Println()

	var increment float64
	for {
		increment = prompter.Float("Percentage increment of tank volume", 10.0)
		if increment > 0 && increment <= 100 {
			break
		}
		fmt.Println("INPUT MAKES NO SENSE. Try again.")
	}

	for target := increment; target <= 100-increment; target += increment {
		height, err := findWettedHeight(maxHeight, 0.01*target, 0.00001, fractionAt)
		if err != nil {
			fmt.Fprintln(os.Stderr, "dipstick:", err)
			os.Exit(1)
		}
		fmt.Printf("Tank volume = %.2f%%, stick wetted length = %.2f\n", target, height)
	}
	fmt.Printf("Tank volume = %.2f%%, stick wetted length = %.2f\n", 100.0, maxHeight)
}
