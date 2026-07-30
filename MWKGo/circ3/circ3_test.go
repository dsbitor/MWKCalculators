package main

import (
	"math"
	"testing"
)

const epsilon = 1e-9

func nearlyEqual(a, b float64) bool {
	return math.Abs(a-b) < epsilon
}

func TestCircleIntersection(t *testing.T) {
	tests := []struct {
		name    string
		center1 point
		r1      float64
		center2 point
		r2      float64
		wantOK  bool
		wantP   point
		wantQ   point
	}{
		{
			// Two radius-5 circles 6 apart intersect at (3,4) and
			// (3,-4): the classic 3-4-5 right triangle.
			name:    "classic 3-4-5 configuration",
			center1: point{X: 0, Y: 0}, r1: 5,
			center2: point{X: 6, Y: 0}, r2: 5,
			wantOK: true,
			wantP:  point{X: 3, Y: -4},
			wantQ:  point{X: 3, Y: 4},
		},
		{
			name:    "coincident centers have no defined intersection",
			center1: point{X: 1, Y: 1}, r1: 2,
			center2: point{X: 1, Y: 1}, r2: 3,
			wantOK: false,
		},
		{
			name:    "circles too far apart do not intersect",
			center1: point{X: 0, Y: 0}, r1: 1,
			center2: point{X: 10, Y: 0}, r2: 1,
			wantOK: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, q, ok := circleIntersection(tt.center1, tt.r1, tt.center2, tt.r2)
			if ok != tt.wantOK {
				t.Fatalf("ok = %v, want %v", ok, tt.wantOK)
			}
			if !ok {
				return
			}
			match := (nearlyEqual(p.X, tt.wantP.X) && nearlyEqual(p.Y, tt.wantP.Y) &&
				nearlyEqual(q.X, tt.wantQ.X) && nearlyEqual(q.Y, tt.wantQ.Y)) ||
				(nearlyEqual(p.X, tt.wantQ.X) && nearlyEqual(p.Y, tt.wantQ.Y) &&
					nearlyEqual(q.X, tt.wantP.X) && nearlyEqual(q.Y, tt.wantP.Y))
			if !match {
				t.Errorf("intersection = (%v, %v), want %v and %v in either order", p, q, tt.wantP, tt.wantQ)
			}
		})
	}
}

func TestLineIntersection(t *testing.T) {
	tests := []struct {
		name   string
		p1, p2 point
		p3, p4 point
		wantOK bool
		want   point
	}{
		{
			name: "crossing diagonals meet at the center",
			p1:   point{X: 0, Y: 0}, p2: point{X: 2, Y: 2},
			p3: point{X: 0, Y: 2}, p4: point{X: 2, Y: 0},
			wantOK: true,
			want:   point{X: 1, Y: 1},
		},
		{
			name: "parallel lines never meet",
			p1:   point{X: 0, Y: 0}, p2: point{X: 1, Y: 0},
			p3: point{X: 0, Y: 1}, p4: point{X: 1, Y: 1},
			wantOK: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := lineIntersection(tt.p1, tt.p2, tt.p3, tt.p4)
			if ok != tt.wantOK {
				t.Fatalf("ok = %v, want %v", ok, tt.wantOK)
			}
			if ok && (!nearlyEqual(got.X, tt.want.X) || !nearlyEqual(got.Y, tt.want.Y)) {
				t.Errorf("intersection = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRadiusFromDistances(t *testing.T) {
	tests := []struct {
		name          string
		d12, d13, d23 float64
		want          float64
		wantErr       bool
	}{
		// A 3-4-5 right triangle's circumradius is half the
		// hypotenuse, a well known identity independent of this code.
		{name: "3-4-5 triangle, distances in order", d12: 5, d13: 3, d23: 4, want: 2.5},
		{name: "3-4-5 triangle, distances reordered", d12: 3, d13: 4, d23: 5, want: 2.5},
		{name: "3-4-5 triangle, distances reordered again", d12: 4, d13: 5, d23: 3, want: 2.5},
		// An equilateral triangle's circumradius is side/sqrt(3).
		{name: "equilateral triangle side 2", d12: 2, d13: 2, d23: 2, want: 2 / math.Sqrt(3)},
		{name: "violates the triangle inequality", d12: 10, d13: 1, d23: 1, wantErr: true},
		{name: "exactly collinear points", d12: 3, d13: 1, d23: 2, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := radiusFromDistances(tt.d12, tt.d13, tt.d23)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("radiusFromDistances() error = nil, want an error")
				}
				return
			}
			if err != nil {
				t.Fatalf("radiusFromDistances() unexpected error: %v", err)
			}
			if !nearlyEqual(got, tt.want) {
				t.Errorf("radiusFromDistances() = %v, want %v", got, tt.want)
			}
		})
	}
}
