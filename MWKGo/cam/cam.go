// cam designs a plate cam profile for a chosen follower motion law
// (straight-line/parabolic-compensated, parabolic, simple harmonic,
// or cycloidal) and follower type (flat-faced or roller), printing
// the cam profile as a table of angle/radius/x/y coordinates, the
// resulting maximum pressure angle, and a table suggesting base
// circle radii for a range of target pressure angles. It can also
// render the cam profile and the follower's displacement, velocity,
// and acceleration curves as an SVG diagram.
//
// Converted from CAM.C (M. W. Klotz), WorkshopUtilities/cam. Per
// ai/plans/c-to-go-conversion-plan.md's Tier 3 "Graphics scope"
// resolution, the original's DOS graphics calls (wline/wcircle,
// __setwindow) are replaced by internal/svgplot, emitting SVG instead
// of drawing to a screen; the interactive "press a key" pause and the
// output-file-save feature are both dropped in favor of an optional
// -svg flag and printing straight to stdout.
package main

import (
	"flag"
	"fmt"
	"math"
	"os"

	"mwkgo/internal/promptio"
	"mwkgo/internal/svgplot"
)

const rpd = math.Pi / 180

// Follower motion law types, matching CAM.C's own ctype menu.
const (
	straightLine   = 1
	parabolic      = 2
	simpleHarmonic = 3
	cycloidal      = 4
)

// Follower types, matching CAM.C's own rtype menu.
const (
	flatFaced      = 1
	rollerFollower = 2
)

// motionLaw holds every parameter needed to evaluate a chosen
// follower motion law at any point in its rise, matching CAM.C's own
// cam() local variables. Beta and the betaX/riseX fields are in
// degrees throughout, matching the original's own bare "beta"
// variable; RPD conversions happen at the point of use, exactly where
// CAM.C itself converts.
type motionLaw struct {
	Type                int
	Rise, Beta          float64
	RiseA, RiseL, RiseD float64
	BetaA, BetaL, BetaD float64
	Rob                 float64 // rise/(beta in radians), used by types 2-4
}

// followerState is the follower's displacement, velocity, and
// acceleration at one angle into the rise.
type followerState struct {
	D, Vc, Ac float64
}

// evaluateMotion returns the follower's state at theta degrees into
// the rise (theta=0 at the start of rise, theta=Beta at the end),
// matching CAM.C's own per-motion-type formulas exactly.
func evaluateMotion(m motionLaw, theta float64) followerState {
	switch m.Type {
	case straightLine:
		ba, bl, bd := m.BetaA*rpd, m.BetaL*rpd, m.BetaD*rpd
		switch {
		case theta <= m.BetaA:
			th := theta * rpd
			return followerState{
				D:  m.RiseA * (th / ba) * (th / ba),
				Vc: 2 * m.RiseA * th / (ba * ba),
				Ac: 2 * m.RiseA / (ba * ba),
			}
		case theta >= m.BetaA+m.BetaL:
			th := (theta - m.BetaA - m.BetaL) * rpd
			return followerState{
				D:  m.Rise - m.RiseD*((bd-th)/bd)*((bd-th)/bd),
				Vc: m.RiseL/bl - 2*m.RiseD*th/(bd*bd),
				Ac: -2 * m.RiseD / (bd * bd),
			}
		default:
			return followerState{
				D:  m.RiseA + m.RiseL*(theta-m.BetaA)/m.BetaL,
				Vc: m.RiseL / bl,
				Ac: 0,
			}
		}
	case parabolic:
		tob := theta / m.Beta
		if tob <= 0.5 {
			return followerState{D: 2 * m.Rise * tob * tob, Vc: 4 * m.Rob * tob, Ac: 4 * m.Rob / m.Beta}
		}
		return followerState{D: m.Rise * (1 - 2*(1-tob)*(1-tob)), Vc: 4 * m.Rob * (1 - tob), Ac: -4 * m.Rob / m.Beta}
	case simpleHarmonic:
		tob := theta / m.Beta
		return followerState{
			D:  0.5 * m.Rise * (1 - math.Cos(math.Pi*tob)),
			Vc: (math.Pi / 2) * m.Rob * math.Sin(math.Pi*tob),
			Ac: (math.Pi / 2) * m.Rob * math.Pi * math.Cos(math.Pi*tob) / m.Beta,
		}
	case cycloidal:
		tob := theta / m.Beta
		return followerState{
			D:  m.Rise * (tob - math.Sin(2*math.Pi*tob)/(2*math.Pi)),
			Vc: m.Rob * (1 - math.Cos(2*math.Pi*tob)),
			Ac: m.Rob * 2 * math.Pi * math.Sin(2*math.Pi*tob) / m.Beta,
		}
	}
	return followerState{}
}

// motionMax returns the motion law's own peak velocity and
// acceleration, used only to scale the follower-motion plots to fit
// the drawing window, matching CAM.C's own vmax/amax computation.
func motionMax(m motionLaw) (vmax, amax float64) {
	switch m.Type {
	case straightLine:
		ba, bd := m.BetaA*rpd, m.BetaD*rpd
		return m.RiseL / (m.BetaL * rpd), 2 * math.Max(m.RiseA/(ba*ba), m.RiseD/(bd*bd))
	case parabolic:
		return 2 * m.Rob, 4 * m.Rob / m.Beta
	case simpleHarmonic:
		return (math.Pi / 2) * m.Rob, (math.Pi / 2) * m.Rob * math.Pi / m.Beta
	case cycloidal:
		return 2 * m.Rob, m.Rob * 2 * math.Pi / m.Beta
	}
	return 0, 0
}

// camPoint is one point on the cam's own profile.
type camPoint struct {
	AngleDeg         float64 // measured from the highest point of the cam
	Radius           float64 // measured from the center of cam rotation
	X, Y             float64
	PressureAngleRad float64
}

// computeCamProfile walks the rise from Beta down to 0 in astep
// increments, matching CAM.C's own main loop exactly.
func computeCamProfile(m motionLaw, rbase, rroll, astep float64) []camPoint {
	var points []camPoint
	for ang := m.Beta; ang >= 0; ang -= astep {
		theta := m.Beta - ang
		s := evaluateMotion(m, theta)
		rcam := math.Hypot(rbase+s.D+rroll, s.Vc) - rroll
		pcam := math.Atan(s.Vc / (rbase + s.D + rroll))
		tcam := ang*rpd - pcam
		points = append(points, camPoint{
			AngleDeg:         ang,
			Radius:           rcam,
			X:                rcam * math.Cos(tcam),
			Y:                rcam * math.Sin(tcam),
			PressureAngleRad: pcam,
		})
	}
	return points
}

// maxPressurePoint records the follower state at the point of maximum
// pressure angle along the cam profile.
type maxPressurePoint struct {
	AngleRad float64
	Vc, D    float64
}

// findMaxPressureAngle returns the point of maximum pressure angle
// among the same profile points computeCamProfile produces, matching
// CAM.C's own running-maximum tracking.
func findMaxPressureAngle(m motionLaw, rbase, rroll, astep float64) maxPressurePoint {
	var best maxPressurePoint
	for ang := m.Beta; ang >= 0; ang -= astep {
		theta := m.Beta - ang
		s := evaluateMotion(m, theta)
		pcam := math.Atan(s.Vc / (rbase + s.D + rroll))
		if pcam > best.AngleRad {
			best = maxPressurePoint{AngleRad: pcam, Vc: s.Vc, D: s.D}
		}
	}
	return best
}

// baseRadiusSuggestion is one entry of the base-radius-vs-pressure-
// angle table CAM.C prints after the main profile calculation.
type baseRadiusSuggestion struct {
	PressureAngleDeg float64
	BaseRadius       float64
}

// suggestBaseRadii returns the base radius that would produce each
// pressure angle from 15 to 45 degrees in 2.5 degree steps, given the
// cam's own maximum-pressure-angle point, matching CAM.C's own
// trailing table exactly.
func suggestBaseRadii(maxP maxPressurePoint, rroll float64) []baseRadiusSuggestion {
	var out []baseRadiusSuggestion
	for pa := 15.0; pa <= 45.0; pa += 2.5 {
		rbase := maxP.Vc/math.Tan(pa*rpd) - maxP.D - rroll
		out = append(out, baseRadiusSuggestion{PressureAngleDeg: pa, BaseRadius: rbase})
	}
	return out
}

// renderCam draws the cam profile and the follower's displacement,
// velocity, and acceleration curves, matching CAM.C's own drawing
// loop: the cam profile itself (color 14) as a sequence of short
// chords between consecutive computed points, and the three follower
// curves (colors 10/11/13) plotted against a shared synthetic
// "rotation angle" axis, each independently scaled to fill the
// drawing window's own vertical range.
func renderCam(m motionLaw, rbase, rroll float64, points []camPoint) *svgplot.Canvas {
	scale := 200.0 / (rbase + m.Rise)
	rbs := rbase * scale
	ocam := scale * (2*rbase + m.Rise)
	xmin := -(320 - 0.5*ocam + rbs)
	xmax := 640 + xmin
	ymin, ymax := -240.0, 240.0

	c := svgplot.New(xmin, ymin, xmax, ymax, 640, 480)
	c.Circle(0, 0, rbs, 12, false)
	c.Line(-0.05*640, 0, 0.05*640, 0, 12)
	c.Line(0, -0.05*480, 0, 0.05*480, 12)

	vmax, amax := motionMax(m)

	xl := rbs * math.Cos(m.Beta*rpd)
	yl := rbs * math.Sin(m.Beta*rpd)
	c.Line(0, 0, xl, yl, 15)
	c.Line(0, 0, xl, -yl, 15)

	xt, yt, vt, at := xmin, ymin, ymin, ymin
	for i, p := range points {
		xc, yc := scale*p.X, scale*p.Y
		c.Line(xl, yl, xc, yc, 14)
		c.Line(xl, -yl, xc, -yc, 14)
		xl, yl = xc, yc

		theta := m.Beta - p.AngleDeg
		s := evaluateMotion(m, theta)
		xp := theta*640/m.Beta + xmin
		yp := s.D*480/m.Rise + ymin
		if i == 0 {
			xt, yt = xp, yp
		}
		c.Line(xt, yt, xp, yp, 10)

		vp := s.Vc*480/vmax + ymin
		if i == 0 {
			vt = vp
		}
		c.Line(xt, vt, xp, vp, 11)

		ap := s.Ac * 240 / amax
		if i == 0 {
			at = ap
		}
		c.Line(xt, at, xp, ap, 13)

		xt, yt, vt, at = xp, yp, vp, ap
	}
	return c
}

func main() {
	svgPath := flag.String("svg", "", "path to write an SVG rendering of the cam profile (optional)")
	flag.Parse()

	prompter, err := promptio.New(os.Stdin, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, "cam:", err)
		os.Exit(1)
	}

	fmt.Println("CAM DESIGN")
	fmt.Println()
	fmt.Println("TYPE OF FOLLOWER MOTION:")
	fmt.Println("  1 - STRAIGHT LINE (Infinite acceleration parabolically compensated)")
	fmt.Println("  2 - PARABOLIC (Smallest acceleration but produces impact shocks)")
	fmt.Println("  3 - SIMPLE HARMONIC (Smooth motion - acceleration causes vibration)")
	fmt.Println("  4 - CYCLOIDAL (RECOMMENDED - No shocks, low vibration and wear)")
	ctype := prompter.Int("Type desired", cycloidal)

	fmt.Println("\nTYPE OF CAM FOLLOWER:")
	fmt.Println("  1 - FLAT-FACED")
	fmt.Println("  2 - ROLLER")
	rtype := prompter.Int("Type desired", flatFaced)
	rroll := 0.0
	if rtype == rollerFollower {
		rroll = prompter.Float("Radius of roller", 0.25)
	}

	fmt.Println("\nBase circle radius is the radius of the part of the cam that provides")
	fmt.Println("the minimum follower extension (the 'lowest' part of the cam)")
	rbase := prompter.Float("Base circle radius", 1.625)

	m := motionLaw{Type: ctype}
	if ctype == straightLine {
		fmt.Println("\nConstant velocity cams have infinite start/end-of-stroke")
		fmt.Println("accelerations that cause unacceptable impact loads and wear.")
		fmt.Println("This is avoided by using constant acceleration and deceleration")
		fmt.Println("(PARABOLIC motion) at the beginning and end of rise.")
		m.RiseA = prompter.Float("Allowable cam rise during acceleration", 0.1)
		m.RiseL = prompter.Float("Cam rise during linear motion portion of stroke", 1.0)
		m.RiseD = prompter.Float("Allowable cam rise during deceleration", 0.15)
		defaultBetaL := 120.0 / (1.0 + 2.0*(m.RiseA+m.RiseD)/m.RiseL)
		m.BetaL = prompter.Float("Cam rotation during which linear motion occurs (deg)", defaultBetaL)
		m.BetaA = 2 * m.BetaL * m.RiseA / m.RiseL
		m.BetaD = 2 * m.BetaL * m.RiseD / m.RiseL
		m.Beta = m.BetaA + m.BetaL + m.BetaD
		m.Rise = m.RiseA + m.RiseL + m.RiseD
		fmt.Printf("Cam will rise %.4f during %.4f deg of rotation\n", m.Rise, m.Beta)
	} else {
		fmt.Println("\nCam rise is the distance the follower travels from the point where")
		fmt.Println("it touches the base circle to the point where it touches the highest")
		fmt.Println("point of the cam.")
		m.Rise = prompter.Float("Cam rise", 1.25)
		fmt.Println("\nCam rotation angle is the angular arc of the cam during which")
		fmt.Println("the follower rise occurs.")
		m.Beta = prompter.Float("Cam rotation angle (deg)", 120.0)
		m.Rob = m.Rise / (m.Beta * rpd)
	}
	fmt.Println()
	astep := prompter.Float("Angular step size for constructing cam (deg)", 1.0)

	motionNames := map[int]string{straightLine: "Straight Line", parabolic: "Parabolic", simpleHarmonic: "Simple Harmonic", cycloidal: "Cycloidal"}
	fmt.Printf("Cam Motion = %s\n", motionNames[ctype])
	if rtype == flatFaced {
		fmt.Println("Follower Type = Flat-Faced")
	} else {
		fmt.Println("Follower Type = Roller")
		fmt.Printf("Follower Roller Radius = %.4f\n", rroll)
	}
	fmt.Printf("Base Circle Radius = %.4f\n", rbase)
	fmt.Printf("Cam Rise (Stroke) = %.4f\n", m.Rise)
	fmt.Printf("Cam Rotation Angle = %.4f deg\n", m.Beta)
	fmt.Println("\nAngle is measured from highest point of cam.")
	fmt.Println("Radius & x,y coordinates measured from center of cam rotation.")
	fmt.Println("\n     Angle      Radius     X-coord     Y-coord")

	points := computeCamProfile(m, rbase, rroll, astep)
	for _, p := range points {
		fmt.Printf("%10.4f  %10.4f  %10.4f  %10.4f\n", p.AngleDeg, p.Radius, p.X, p.Y)
	}

	maxP := findMaxPressureAngle(m, rbase, rroll, astep)
	fmt.Printf("\nThe maximum pressure angle was %.4f deg\n", maxP.AngleRad*180/math.Pi)
	fmt.Println("\nMaximum pressure angle should be less than 35 deg. to limit side")
	fmt.Println("thrust on the follower.  Pressure angle can be adjusted by modifying")
	fmt.Println("the base radius.  For the cam just specified, the following holds:")
	fmt.Println("\nPressure Angle (deg)   Base Radius")
	for _, s := range suggestBaseRadii(maxP, rroll) {
		fmt.Printf("%20.4f   %11.4f\n", s.PressureAngleDeg, s.BaseRadius)
	}

	if *svgPath != "" {
		if err := os.WriteFile(*svgPath, []byte(renderCam(m, rbase, rroll, points).String()), 0o644); err != nil {
			fmt.Fprintln(os.Stderr, "cam:", err)
			os.Exit(1)
		}
		fmt.Printf("\nWrote cam profile diagram to %s\n", *svgPath)
	}
}
