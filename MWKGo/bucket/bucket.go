// bucket computes the slant heights at which to mark a bucket-shaped
// container (a conical frustum) to divide its volume into equal
// parts, since a frustum's varying cross-section means equal height
// divisions do not give equal volume divisions the way they would
// for a cylinder.
//
// Converted from BUCKET.C (M. W. Klotz, 11/99),
// WorkshopUtilities/bucket.
package main

import (
	"fmt"
	"math"
	"os"

	"mwkgo/internal/mwktrig"
	"mwkgo/internal/promptio"
)

// maxSearchIterations bounds the incremental search
// slantHeightForVolume uses to locate each division height. The
// original program treats reaching this many iterations without
// converging as an error and stops rather than looping forever.
const maxSearchIterations = 500

// bucketVolume returns the volume of a conical frustum of the given
// perpendicular height and base and top radii, using the standard
// frustum-of-a-cone volume formula.
func bucketVolume(height, smallRadius, largeRadius float64) float64 {
	return math.Pi * height * (largeRadius*largeRadius + largeRadius*smallRadius + smallRadius*smallRadius) / 3
}

// bucketGeometry is a conical frustum's shape, derived from its
// measured dimensions.
type bucketGeometry struct {
	SmallRadius  float64
	LargeRadius  float64
	Height       float64 // perpendicular height, not slant height
	Volume       float64
	HalfAngleDeg float64
}

// computeBucketGeometry derives a bucket's perpendicular height,
// total volume, and half angle from its measured big and small
// diameters and its slant height (measured along the sloping side).
func computeBucketGeometry(bigDiameter, smallDiameter, slantHeight float64) bucketGeometry {
	smallRadius := 0.5 * smallDiameter
	largeRadius := 0.5 * bigDiameter
	halfAngleDeg := mwktrig.ClampedAsinDeg((largeRadius - smallRadius) / slantHeight)
	height := slantHeight * math.Cos(halfAngleDeg*math.Pi/180)

	return bucketGeometry{
		SmallRadius:  smallRadius,
		LargeRadius:  largeRadius,
		Height:       height,
		Volume:       bucketVolume(height, smallRadius, largeRadius),
		HalfAngleDeg: halfAngleDeg,
	}
}

// slantHeightForVolume finds the slant height, measured from the
// small end, at which the bucket's partial volume equals
// targetVolume, starting the search from initialHeightGuess. It
// searches by the same incremental stepping the original program
// uses: nudge the height by a small fixed step, in whichever
// direction reduces the volume error, until the error falls within
// tolerance or maxSearchIterations is reached.
func slantHeightForVolume(geometry bucketGeometry, targetVolume, initialHeightGuess float64) (float64, error) {
	tolerance := 0.001 * geometry.Volume
	step := 0.001 * geometry.Height
	cosHalfAngle := math.Cos(geometry.HalfAngleDeg * math.Pi / 180)

	height := initialHeightGuess
	for i := 0; i < maxSearchIterations; i++ {
		radius := geometry.SmallRadius + height*(geometry.LargeRadius-geometry.SmallRadius)/geometry.Height
		volume := bucketVolume(height, geometry.SmallRadius, radius)

		if math.Abs(targetVolume-volume) < tolerance {
			return height / cosHalfAngle, nil
		}
		if targetVolume > volume {
			height += step
		} else {
			height -= step
		}
	}
	return 0, fmt.Errorf("iteration error: no convergence within %d steps", maxSearchIterations)
}

func main() {
	prompter, err := promptio.New(os.Stdin, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, "bucket:", err)
		os.Exit(1)
	}

	fmt.Println("BUCKET CALCULATIONS")
	fmt.Println()

	bigDiameter := prompter.Float("Diameter of bucket big end", 4.0)
	smallDiameter := prompter.Float("Diameter of bucket small end", 3.0)
	slantHeight := prompter.Float("Slant height of bucket", 6.0)
	divisions := prompter.Int("Divide volume into how many parts", 4)
	fmt.Println()

	geometry := computeBucketGeometry(bigDiameter, smallDiameter, slantHeight)
	fmt.Printf("Total volume of bucket = %.3f\n", geometry.Volume)

	for i := 1; i <= divisions; i++ {
		fraction := float64(i) / float64(divisions)
		targetVolume := fraction * geometry.Volume
		initialGuess := fraction * geometry.Height

		slant, err := slantHeightForVolume(geometry, targetVolume, initialGuess)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Printf("Slant height for %d/%d of total volume = %.3f\n", i, divisions, slant)
	}
}
