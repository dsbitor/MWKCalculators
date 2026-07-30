// circ3 finds the radius of the circle that passes through three
// points, given only the measured distances between each pair of
// points, not their coordinates.
//
// Converted from CIRC3.C (M. W. Klotz, 7/99), Math/circ3.
package main

import (
	"errors"
	"fmt"
	"math"
	"os"

	"mwkgo/internal/promptio"
)

// point is a 2D Cartesian coordinate.
type point struct {
	X, Y float64
}

func distance(a, b point) float64 {
	return math.Hypot(b.X-a.X, b.Y-a.Y)
}

// circleIntersection returns the two points that lie at distance r1
// from center1 and distance r2 from center2: the intersection points
// of two circles. ok is false when the centers coincide or when no
// real intersection exists (the circles are too far apart, one is
// nested inside the other with no overlap, or the computed geometry
// is otherwise degenerate).
func circleIntersection(center1 point, r1 float64, center2 point, r2 float64) (p, q point, ok bool) {
	centerDistance := distance(center1, center2)
	if centerDistance == 0 {
		return point{}, point{}, false
	}

	a := (r1*r1 - r2*r2 + centerDistance*centerDistance) / (2 * centerDistance)
	h2 := r1*r1 - a*a
	if h2 < 0 {
		return point{}, point{}, false
	}
	h := math.Sqrt(h2)

	midX := center1.X + a*(center2.X-center1.X)/centerDistance
	midY := center1.Y + a*(center2.Y-center1.Y)/centerDistance
	offsetX := h * (center2.Y - center1.Y) / centerDistance
	offsetY := h * (center2.X - center1.X) / centerDistance

	p = point{X: midX + offsetX, Y: midY - offsetY}
	q = point{X: midX - offsetX, Y: midY + offsetY}
	return p, q, true
}

// bisect returns two points on the perpendicular bisector of the
// segment p1-p2, constructed the same way a compass-and-straightedge
// bisector is: as the intersection of two equal circles, one
// centered on each endpoint, with radius equal to the segment
// length. ok is false when p1 and p2 coincide.
func bisect(p1, p2 point) (a, b point, ok bool) {
	segmentLength := distance(p1, p2)
	return circleIntersection(p1, segmentLength, p2, segmentLength)
}

// lineIntersection returns the point where line p1-p2 crosses line
// p3-p4. ok is false when the two lines are parallel.
func lineIntersection(p1, p2, p3, p4 point) (point, bool) {
	d := (p2.Y-p1.Y)*(p4.X-p3.X) - (p4.Y-p3.Y)*(p2.X-p1.X)
	if d == 0 {
		return point{}, false
	}

	cross12 := p1.X*p2.Y - p2.X*p1.Y
	cross34 := p3.X*p4.Y - p4.X*p3.Y
	xi := (cross12*(p4.X-p3.X) - cross34*(p2.X-p1.X)) / d
	yi := (cross12*(p4.Y-p3.Y) - cross34*(p2.Y-p1.Y)) / d
	return point{X: xi, Y: yi}, true
}

// circumcircle returns the center and radius of the circle that
// passes through p1, p2, and p3, found as the intersection of the
// perpendicular bisectors of two of the triangle's sides. ok is
// false when the three points are collinear, in which case no
// circle passes through all three.
func circumcircle(p1, p2, p3 point) (center point, radius float64, ok bool) {
	a1, b1, ok1 := bisect(p1, p2)
	a2, b2, ok2 := bisect(p2, p3)
	if !ok1 || !ok2 {
		return point{}, 0, false
	}

	center, ok = lineIntersection(a1, b1, a2, b2)
	if !ok {
		return point{}, 0, false
	}
	return center, distance(p1, center), true
}

// radiusFromDistances returns the radius of the circle passing
// through three points, given only the three pairwise distances
// between them. It returns an error when the three distances do not
// form a triangle (the triangle inequality fails) or describe three
// collinear points.
func radiusFromDistances(d12, d13, d23 float64) (float64, error) {
	// Reorder so d12 is the longest side. This mirrors the original
	// program's placement strategy: p1=(0,0) and p2=(d12,0) puts the
	// two most-separated points down first, which keeps the
	// remaining construction numerically well conditioned.
	if d13 > d12 {
		d12, d13 = d13, d12
	}
	if d23 > d12 {
		d12, d23 = d23, d12
	}

	if d13+d23 <= d12 {
		return 0, fmt.Errorf("distances %.4g, %.4g, %.4g do not satisfy the triangle inequality", d12, d13, d23)
	}

	p1 := point{X: 0, Y: 0}
	p2 := point{X: d12, Y: 0}

	p3, _, ok := circleIntersection(p1, d13, p2, d23)
	if !ok {
		return 0, fmt.Errorf("distances %.4g, %.4g, %.4g do not define a third point", d12, d13, d23)
	}

	_, radius, ok := circumcircle(p1, p2, p3)
	if !ok {
		return 0, errors.New("the three points are collinear; no circle passes through all three")
	}
	return radius, nil
}

func main() {
	prompter, err := promptio.New(os.Stdin, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, "circ3:", err)
		os.Exit(1)
	}

	fmt.Println("CIRCLE THROUGH THREE POINTS")
	fmt.Println()
	fmt.Println("Number the three points 1-3 (in any order).")
	fmt.Println("Enter lengths between points as specified below:")
	fmt.Println()

	for {
		d12 := prompter.Float("length 1-2", 2)
		d13 := prompter.Float("length 1-3", 1.1)
		d23 := prompter.Float("length 2-3", 1.1)

		radius, err := radiusFromDistances(d12, d13, d23)
		if err != nil {
			fmt.Println("NO SOLUTION FOR THOSE DIMENSIONS - CHECK YOUR MEASUREMENTS")
			continue
		}

		fmt.Printf("RADIUS = %f\n", radius)
		fmt.Printf("DIAMETER = %f\n", 2*radius)
		return
	}
}
