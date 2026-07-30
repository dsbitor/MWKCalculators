// sun computes the sun's position (right ascension, declination,
// distance, altitude, azimuth), the equation of time, sunrise and
// sunset, and sundial shadow-angle geometry, for a chosen date, time,
// and observer location, using the long-term (1800-2100) trigonometric
// series from the 1978 "Almanac for Computers" (Nautical Almanac
// Office, US Naval Observatory).
//
// An observer's latitude, longitude, and UT offset(s) are specific to
// one location, reused across many runs from that location — not
// universal reference data, and not a piece of equipment
// configuration, but the same "machine-specific persistent" bucket by
// the same reasoning: it belongs to one owner's setup, not everyone's.
// So, like gearatio, change, and ddh in earlier conversion groups,
// this program reads it from the user's own database rather than the
// shared reference database; see
// ai/plans/c-to-go-conversion-plan.md, "Data-file strategy for
// Tier 2".
//
// Converted from SUN.C (M. W. Klotz), Misc/sun.
package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"math"
	"os"
	"strings"
	"time"

	"mwkgo/internal/csvtable"
	"mwkgo/internal/mwktrig"
	"mwkgo/internal/promptio"
	"mwkgo/internal/userdata"
)

const settingsTable = "sun_settings"

// maxMeridianTransitIterations bounds the meridian-transit secant
// search, replacing the original's unbounded convergence loop per
// coding-style.md Rule 2 (in practice this converges in a handful of
// iterations; this bound is far larger than ever needed, only
// preventing a genuinely non-convergent case from looping forever).
const maxMeridianTransitIterations = 200

// auMiles is 1 astronomical unit in miles, used only to display the
// earth-sun distance in a familiar unit alongside AU.
const auMiles = 92955807.267

func sind(deg float64) float64 { return math.Sin(deg * math.Pi / 180) }
func cosd(deg float64) float64 { return math.Cos(deg * math.Pi / 180) }
func tand(deg float64) float64 { return math.Tan(deg * math.Pi / 180) }
func atand(x float64) float64  { return math.Atan(x) * 180 / math.Pi }
func sgn(x float64) float64 {
	switch {
	case x > 0:
		return 1
	case x < 0:
		return -1
	default:
		return 0
	}
}

// normab wraps x into [a,b), assuming a period of b-a, matching
// SUN.C's own normab().
func normab(x, a, b float64) float64 {
	z := b - a
	for x < a {
		x += z
	}
	for x > b {
		x -= z
	}
	return x
}

// caz returns the azimuth (degrees, N=0, E=90, S=180, W=270) from
// point 1 to point 2 on a sphere, given as latitude/longitude pairs,
// matching SUN.C's own caz() exactly.
func caz(lat1, lon1, lat2, lon2 float64) float64 {
	slat1, clat1 := sind(lat1), cosd(lat1)
	slat2, clat2 := sind(lat2), cosd(lat2)
	cosA := slat1*slat2 + clat1*clat2*cosd(math.Abs(lon2-lon1))
	a := mwktrig.ClampedAcosDeg(cosA)
	sinA := sind(a)
	cosb := (slat2 - slat1*cosA) / (clat1 * sinA)
	az := mwktrig.ClampedAcosDeg(cosb)
	if lon1 == lon2 {
		if lat1 < lat2 {
			az = 0
		} else {
			az = 180
		}
	}
	if lon2 > lon1 {
		az = 360 - az
	}
	if math.Abs(lon2-lon1) > 180 {
		az = 360 - az
	}
	return az
}

func isLeapYear(y int) bool { return y%4 == 0 && y%100 != 0 || y%400 == 0 }

func daysInMonth(year, month int) int {
	days := [...]int{31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}
	d := days[month-1]
	if month == 2 && isLeapYear(year) {
		d = 29
	}
	return d
}

// tnorm normalizes ut into [0,24), carrying day/month/year forward or
// backward as needed, matching SUN.C's own tnorm() exactly (including
// its one-off computation of ut's initial sign, which fixes the carry
// direction for the whole normalization even if a day boundary is
// crossed more than once).
func tnorm(month, day, year int, ut float64) (int, int, int, float64, error) {
	if year <= 0 || month < 1 || month > 12 {
		return 0, 0, 0, 0, fmt.Errorf("invalid date %d/%d/%d", month, day, year)
	}
	if day < 1 || day > daysInMonth(year, month) {
		return 0, 0, 0, 0, fmt.Errorf("invalid date %d/%d/%d", month, day, year)
	}

	if ut >= 0 && ut < 24 {
		return month, day, year, ut, nil
	}
	forward := sgn(ut) > 0

	for ut < 0 || ut >= 24 {
		if forward {
			ut -= 24
			day++
			if day > daysInMonth(year, month) {
				day = 1
				month++
				if month > 12 {
					month = 1
					year++
				}
			}
		} else {
			ut += 24
			day--
			if day < 1 {
				month--
				if month < 1 {
					month = 12
					year--
				}
				day = daysInMonth(year, month)
			}
		}
	}
	return month, day, year, ut, nil
}

// sunDat is one computed sun position and related quantities for a
// given date, time, and observer location, matching SUN.C's own
// struct sundat.
type sunDat struct {
	Month, Day, Year int
	UT               float64 // universal time (hr), normalized to [0,24)
	Lat, Lon         float64

	RA, Dec, Dist    float64
	Alt, Az          float64
	EQT              float64
	GMST, GAST, GHAA float64
	LHAS, GHAS       float64
	SD, HParx, AParx float64
	UTRise, UTSet    float64
}

// suncomp computes the sun's position and related quantities for the
// given date, time (UT), and observer latitude/longitude, matching
// SUN.C's own suncomp() exactly — the long-term (1800-2100)
// trigonometric series from the 1978 Almanac for Computers.
func suncomp(lat, lon float64, month, day, year int, ut float64) (sunDat, error) {
	month, day, year, ut, err := tnorm(month, day, year, ut)
	if err != nil {
		return sunDat{}, err
	}

	slat, clat := sind(lat), cosd(lat)
	et := ut + 49.3/3600.0

	ts1 := 275 * month / 9
	ts2 := (month + 9) / 12
	ts3 := year % 4
	ts4 := 1 + (ts3+2)/3
	nn := ts1 - ts2*ts4 + day - 30

	xs1 := 100.0*float64(year) + float64(month) - 190002.5
	xs2 := 0.5 * (1 - sgn(xs1))
	xjed0 := 367.0*float64(year) - float64(int(0.25*7.0*float64(year+ts2))) + float64(ts1) + float64(day) + 1721013.5 + xs2
	jd := xjed0 + ut/24.0
	jed := xjed0 + et/24.0
	tcap := (jed - 2415020.0) / 36525.0

	gmst := normab(6.67170278+.0657098232*(xjed0-2433282.5)+1.0027379093*ut, 0, 24)
	omega := normab(372.1133-.0529539*(jd-2433282.5), 0, 360)
	eqe := -0.00029 * sind(omega)
	gast := normab(gmst+eqe, 0, 24)
	ghaa := 15.0 * gast
	sunml := normab(279.697+.98564734*(jd-2415020.0), 0, 360)
	x := sunml * math.Pi / 180

	// sundial time = clock time + eqt; eqt = true sun - fictitious sun.
	eqt := -97.8*math.Sin(x) - 431.3*math.Cos(x) + 596.6*math.Sin(2*x) - 1.9*math.Cos(2*x) +
		4*math.Sin(3*x) + 19.3*math.Cos(3*x) - 12.7*math.Sin(4*x)

	t1 := tcap
	t2 := t1 * t1
	rpd := math.Pi / 180

	g := rpd * normab(358.475833+35999.04975*t1-.00015*t2, 0, 360)  // mean anomaly of sun
	l := rpd * normab(279.696678+36000.76892*t1+.000303*t2, 0, 360) // mean longitude of sun
	c := rpd * normab(270.434164+480960.0*t1+307.883142*t1-.001133*t2, 0, 360)
	n := rpd * normab(259.183275-1800.0*t1-134.142008*t1+.002078*t2, 0, 360)
	v := rpd * normab(212.603219+58320.0*t1+197.803875*t1+.001286*t2, 0, 360)
	m := rpd * normab(319.529425+19080.0*t1+59.8585*t1+.000181*t2, 0, 360)
	j := rpd * normab(225.444651+2880.0*t1+154.906654*t1, 0, 360)

	sinl := math.Sin(l)
	singml := math.Sin(g - l)
	singl := math.Sin(g + l)
	cosg := math.Cos(g)
	sing := math.Sin(g)

	tsun := .39793*sinl + 9.998999e-03*singml + .003334*singl - .000208*t1*sinl +
		.000042*math.Sin(2*g+l) - .00004*math.Cos(l) - .000039*math.Sin(n-l) -
		.00003*t1*singml - .000014*math.Sin(2*g-l) - .00001*math.Cos(g-l-j) -
		.00001*t1*singl

	rsun := 1.000421 - .033503*cosg - .00014*math.Cos(2*g) + .000084*t1*cosg -
		.000033*math.Sin(g-j) + .000027*math.Sin(2*(g-v))

	psun := -.041295*math.Sin(2*l) + .032116*sing - .001038*math.Sin(g-2*l) - .000346*math.Sin(g+2*l) -
		.000095 - .00008*t1*sing - .000079*math.Sin(n) + .000068*math.Sin(2*g) +
		.000046*t1*math.Sin(2*l) + .00003*math.Sin(c-l) - .000025*math.Cos(g-j) +
		.000024*math.Sin(4*g-8*m+3*j) - .000019*math.Sin(g-v) - .000017*math.Cos(2*(g-v))

	xr := psun / math.Sqrt(rsun-tsun*tsun)
	ra := normab((l*180/math.Pi+mwktrig.ClampedAsinDeg(xr))/15.0, 0, 24)
	dec := mwktrig.ClampedAsinDeg(tsun / math.Sqrt(rsun))
	dist := math.Sqrt(rsun)

	lha := 15.0*(gast-ra) - lon
	sinDec, cosDec := sind(dec), cosd(dec)
	// SUN.C also computes sin(lha) here, but its own alt formula never
	// actually uses it (only cos(lha) is needed) — an unused value in
	// the original too, not carried over here.
	clha := cosd(lha)
	alt := mwktrig.ClampedAsinDeg(slat*sinDec + clat*cosDec*clha)

	ghas := normab(ghaa-ra*15.0, 0, 360)
	lhas := normab(ghas-lon, 0, 360)
	az := caz(lat, lon, dec, ghas)
	semiDiam := 959.63 / (3600.0 * dist)
	hparx := 8.793999 * semiDiam / 959.63
	aparx := mwktrig.ClampedAsinDeg(sind(hparx) * cosd(alt))

	var utrs [2]float64
	for i := 0; i < 2; i++ {
		tss := float64(nn) + (6.0+12.0*float64(i)+lon/15.0)/24.0
		xm := normab(.9856*tss-3.251, 0, 360)
		xl := normab(xm+1.916*sind(xm)+.02*sind(2*xm)+282.565, 0, 360)
		alpha := atand(.91746 * tand(xl))
		for math.Abs(alpha-xl) >= 10 {
			alpha += 180
		}
		sdel := .39782 * sind(xl)
		cdel := math.Sqrt(1 - sdel*sdel)
		cosh := (cosd(90.0+5.0/6.0) - sdel*slat) / (cdel * clat)
		if math.Abs(cosh) > 1 {
			utrs[i] = 1e10
			continue
		}
		hh := mwktrig.ClampedAcosDeg(cosh)
		if i == 0 {
			hh = 360 - hh
		}
		utrs[i] = normab((hh+alpha)/15.0-.06571*tss-6.62, 0, 24) + lon/15.0
	}

	return sunDat{
		Month: month, Day: day, Year: year, UT: ut,
		Lat: lat, Lon: lon,
		RA: ra, Dec: dec, Dist: dist,
		Alt: alt, Az: az,
		EQT:    eqt,
		GMST:   gmst,
		GAST:   gast,
		GHAA:   ghaa,
		LHAS:   lhas,
		GHAS:   ghas,
		SD:     semiDiam,
		HParx:  hparx,
		AParx:  aparx,
		UTRise: utrs[0],
		UTSet:  utrs[1],
	}, nil
}

// findEQTMinimum scans every day of month/year for the smallest
// |equation of time|, matching SUN.C's own 'e' option.
func findEQTMinimum(lat, lon float64, month, year int) (sunDat, error) {
	ut := 12.0 + lon/15.0
	minAbs := math.Inf(1)
	bestDay := 1
	for d := 1; d <= daysInMonth(year, month); d++ {
		s, err := suncomp(lat, lon, month, d, year, ut)
		if err != nil {
			return sunDat{}, err
		}
		if math.Abs(s.EQT) < minAbs {
			minAbs = math.Abs(s.EQT)
			bestDay = d
		}
	}
	return suncomp(lat, lon, month, bestDay, year, ut)
}

// scanFourDayWindow walks 96 hours (4 days), starting at UT=12 on the
// given fixed month/day, calling keep(candidate) whenever candidate
// should replace the running best, matching SUN.C's own 'f'/'s'/'v'/'w'
// solstice/equinox searches — each just differs in which quantity is
// tracked and how "better" is defined, so that difference is the only
// thing each caller needs to supply.
func scanFourDayWindow(lat, lon float64, month, day, year int, keep func(candidate, best sunDat) bool) (sunDat, error) {
	ut := 12.0
	s, err := suncomp(lat, lon, month, day, year, ut)
	if err != nil {
		return sunDat{}, err
	}
	best := s
	for i := 1; i < 4*24; i++ {
		s, err = suncomp(lat, lon, s.Month, s.Day, s.Year, s.UT+1)
		if err != nil {
			return sunDat{}, err
		}
		if keep(s, best) {
			best = s
		}
	}
	return best, nil
}

func findFallEquinox(lat, lon float64, year int) (sunDat, error) {
	return scanFourDayWindow(lat, lon, 9, 19, year, func(c, b sunDat) bool {
		return math.Abs(c.Dec) < math.Abs(b.Dec)
	})
}

func findSummerSolstice(lat, lon float64, year int) (sunDat, error) {
	return scanFourDayWindow(lat, lon, 6, 19, year, func(c, b sunDat) bool {
		return c.Dec > b.Dec
	})
}

func findVernalEquinox(lat, lon float64, year int) (sunDat, error) {
	return scanFourDayWindow(lat, lon, 3, 19, year, func(c, b sunDat) bool {
		return math.Abs(c.Dec) < math.Abs(b.Dec)
	})
}

func findWinterSolstice(lat, lon float64, year int) (sunDat, error) {
	return scanFourDayWindow(lat, lon, 12, 19, year, func(c, b sunDat) bool {
		return c.Dec < b.Dec
	})
}

// findMeridianTransit finds the UT at which the sun crosses the local
// meridian (azimuth 180 for northern latitudes, 0 for southern) on
// the given date, by the same secant-method iteration SUN.C's own 'm'
// option uses.
func findMeridianTransit(lat, lon float64, month, day, year int) (sunDat, error) {
	azOver := 180.0
	if lat <= 0 {
		azOver = 0.0
	}
	ut0 := 12.0 + lon/15.0
	curUT := ut0
	s, err := suncomp(lat, lon, month, day, year, curUT)
	if err != nil {
		return sunDat{}, err
	}

	azLast := math.Inf(1)
	var ulast float64
	for i := 0; math.Abs(s.Az-azOver) > 0.001; i++ {
		if i >= maxMeridianTransitIterations {
			return sunDat{}, fmt.Errorf("meridian transit search did not converge after %d iterations", i)
		}
		if azLast > 1e9 {
			azLast = s.Az
			ulast = curUT
			curUT = ut0 - s.EQT/3600.0
		} else {
			newUT := ut0 + (azOver-azLast)*(curUT-ulast)/(s.Az-azLast)
			azLast = s.Az
			ulast = newUT
			curUT = newUT
		}
		s, err = suncomp(lat, lon, month, day, year, curUT)
		if err != nil {
			return sunDat{}, err
		}
	}
	return s, nil
}

// hms formats t (hours) as HH:MM:SS.
func hms(t float64) string {
	h := int(t)
	f := (t - math.Floor(t)) * 3600.0
	mn := int(f / 60)
	s := int(f - float64(mn)*60)
	return fmt.Sprintf("%02d:%02d:%02d", h, mn, s)
}

// dms formats a (degrees) as [-]DDD:MM:SS.
func dms(a float64) string {
	sign := " "
	if a < 0 {
		sign = "-"
	}
	abs := math.Abs(a)
	d := int(math.Floor(abs))
	f := (abs - math.Floor(abs)) * 3600.0
	mn := int(f / 60)
	s := int(f - float64(mn)*60)
	return fmt.Sprintf("%s%3d:%02d:%02d", sign, d, mn, s)
}

// pstpdt formats ut adjusted by a standard-time or daylight-time
// offset, marking it "<ST"/"<DT" if the offset rolls it back to the
// previous day, matching SUN.C's own pstpdt().
func pstpdt(ut, ost, odt float64) (standard, daylight string) {
	format := func(offset float64, label string) string {
		pt := ut + offset
		if pt < 0 {
			pt += 24
			label = "<" + label
		}
		return fmt.Sprintf("%s = %s", label, hms(pt))
	}
	return format(ost, "ST"), format(odt, "DT")
}

// sundial is one shadow-angle geometry result for a plane (defined by
// its own normal direction) and a gnomon (a straight rod, defined by
// its own direction), matching SUN.C's own struct sundial and sd().
type sundial struct {
	Alpha, Beta, Gamma, Ratio float64
}

// computeSundial finds the shadow cast on a plane (normal direction
// thetaP/phiP) by a gnomon (direction thetaG/phiG) given the sun's
// current altitude/azimuth, matching SUN.C's own sd() exactly. If the
// gnomon lies in the plane (delta == 0, casting no meaningful shadow),
// every result field is the sentinel 1e10, matching the original.
func computeSundial(alt, az, thetaP, phiP, thetaG, phiG float64) sundial {
	d := sundial{Alpha: 1e10, Beta: 1e10, Gamma: 1e10}

	ts := cosd(alt)
	as, bs, cs := ts*sind(az), ts*cosd(az), sind(alt)
	ts = cosd(thetaP)
	ap, bp, cp := ts*sind(phiP), ts*cosd(phiP), sind(thetaP)
	ts = cosd(thetaG)
	ag, bg, cg := ts*sind(phiG), ts*cosd(phiG), sind(thetaG)

	delta := ap*as + bp*bs + cp*cs
	if delta == 0 {
		return d
	}

	ab := ag*bs - as*bg
	ac := ag*cs - as*cg
	bc := bg*cs - bs*cg
	k := 1 / delta
	x := k * (bp*ab + cp*ac)
	y := k * (cp*bc - ap*ab)
	z := -k * (ap*ac + bp*bc)
	ls := math.Sqrt(x*x + y*y + z*z)

	d.Alpha = mwktrig.ClampedAcosDeg((x*cp - z*ap) / (ls * math.Sqrt(ap*ap+cp*cp)))
	d.Beta = mwktrig.ClampedAcosDeg((y*cp - z*bp) / (ls * math.Sqrt(bp*bp+cp*cp)))
	d.Gamma = mwktrig.ClampedAcosDeg((x*ap + y*bp + z*cp) / ls)
	d.Ratio = ls
	return d
}

func printSunDat(s sunDat, ost, odt float64) {
	fmt.Printf("%d/%d/%d  UT = %s", s.Month, s.Day, s.Year, hms(s.UT))
	st, dt := pstpdt(s.UT, ost, odt)
	fmt.Printf("  %s  %s\n", st, dt)

	fmt.Printf("RA   = %.4f hr = %.4f deg = %s\n", s.RA, normab(s.RA*15, 0, 360), dms(normab(s.RA*15, 0, 360)))
	fmt.Printf("DEC  = %.4f deg = %s\n", s.Dec, dms(s.Dec))
	fmt.Printf("DIST = %.4f au = %.0f miles\n", s.Dist, s.Dist*auMiles)
	fmt.Printf("ALT  = %.4f deg = %s\n", s.Alt, dms(s.Alt))
	fmt.Printf("ZEN  = %.4f deg = %s\n", 90-s.Alt, dms(90-s.Alt))
	fmt.Printf("AZ   = %.4f deg = %s\n", s.Az, dms(s.Az))
	fmt.Printf("SD   = %.4f deg = %s\n", s.SD, dms(s.SD))
	fmt.Printf("GMST = %.4f hr\n", s.GMST)
	fmt.Printf("GAST = %.4f hr\n", s.GAST)
	fmt.Printf("GHAA = %.4f deg\n", s.GHAA)
	fmt.Printf("LHAS = %.4f deg\n", s.LHAS)
	fmt.Printf("GHAS = %.4f deg\n", s.GHAS)
	dialNote := "dial fast"
	if s.EQT < 0 {
		dialNote = "dial slow"
	}
	fmt.Printf("EQT  = %.4f sec = %s  %s\n", s.EQT, dms(s.EQT/3600), dialNote)

	riseST, riseDT := pstpdt(s.UTRise, ost, odt)
	fmt.Printf("SUNRISE  UT = %s  %s  %s\n", hms(s.UTRise), riseST, riseDT)
	setST, setDT := pstpdt(s.UTSet, ost, odt)
	fmt.Printf("SUNSET   UT = %s  %s  %s\n", hms(s.UTSet), setST, setDT)
	daylight := s.UTSet - s.UTRise
	fmt.Printf("DAYLIGHT %s  =  %.2f %%\n", hms(daylight), 100*daylight/24)
}

func printSundial(label string, d sundial) {
	fmt.Println(label)
	fmt.Printf("alpha (wrt N x P) = %.2f | beta (wrt P x E) = %.2f\n", d.Alpha, d.Beta)
	fmt.Printf("gamma (wrt P) = %.2f | shadow/(gnomon length) ratio = %.3f\n", d.Gamma, d.Ratio)
}

func printMenu() {
	fmt.Println("a  current time")
	fmt.Println("b  input time")
	fmt.Println("e  equation of time minimum for a month")
	fmt.Println("f  fall equinox")
	fmt.Println("m  meridian transit")
	fmt.Println("s  summer solstice")
	fmt.Println("v  vernal equinox")
	fmt.Println("w  winter solstice")
	fmt.Println("h  show this menu")
	fmt.Println("q  quit")
}

type settings struct {
	Lat, Lon, OST, ODT float64
}

func loadSettings(ctx context.Context, db *sql.DB) (settings, bool, error) {
	var s settings
	err := db.QueryRowContext(ctx,
		`SELECT latitude_deg, longitude_deg, standard_time_offset_hr, daylight_time_offset_hr FROM `+settingsTable+` WHERE id = 1`).
		Scan(&s.Lat, &s.Lon, &s.OST, &s.ODT)
	if err == sql.ErrNoRows {
		return settings{}, false, nil
	}
	if err != nil {
		return settings{}, false, fmt.Errorf("query %s: %w", settingsTable, err)
	}
	return s, true, nil
}

func main() {
	importPath := flag.String("import", "", "import latitude/longitude/UT offsets from a CSV file and exit")
	flag.Parse()

	ctx := context.Background()
	db, err := userdata.Open(ctx)
	if err != nil {
		fmt.Fprintln(os.Stderr, "sun:", err)
		os.Exit(1)
	}
	defer db.Close()

	if *importPath != "" {
		f, err := os.Open(*importPath)
		if err != nil {
			fmt.Fprintln(os.Stderr, "sun:", err)
			os.Exit(1)
		}
		defer f.Close()
		if err := csvtable.Import(ctx, db, settingsTable, f); err != nil {
			fmt.Fprintln(os.Stderr, "sun:", err)
			os.Exit(1)
		}
		fmt.Printf("imported location settings from %s\n", *importPath)
		return
	}

	loc, ok, err := loadSettings(ctx, db)
	if err != nil {
		fmt.Fprintln(os.Stderr, "sun:", err)
		os.Exit(1)
	}
	if !ok {
		fmt.Println("No location configured yet.")
		fmt.Println("This is specific to your own observing location, so nothing is pre-loaded.")
		fmt.Println("Import your latitude, longitude, and UT offsets as CSV:")
		fmt.Println()
		fmt.Println("    sun -import my-location.csv")
		fmt.Println()
		fmt.Println("See docs/calculators/sun.md for the expected format.")
		os.Exit(1)
	}

	prompter, err := promptio.New(os.Stdin, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, "sun:", err)
		os.Exit(1)
	}

	now := time.Now().UTC()
	month, day, year := int(now.Month()), now.Day(), now.Year()
	ut := float64(now.Hour()) + float64(now.Minute())/60 + float64(now.Second())/3600

	printMenu()
	for {
		op := strings.ToLower(prompter.Line("Option ? "))
		if op == "" {
			op = " "
		}
		var s sunDat
		var computeErr error

		switch op[0] {
		case 'a':
			now = time.Now().UTC()
			month, day, year = int(now.Month()), now.Day(), now.Year()
			ut = float64(now.Hour()) + float64(now.Minute())/60 + float64(now.Second())/3600
		case 'b':
			month = prompter.Int("month", month)
			day = prompter.Int("day", day)
			year = prompter.Int("year", year)
			ut = prompter.Float("ut", ut)
		case 'e':
			month = prompter.Int("month", month)
			year = prompter.Int("year", year)
			s, computeErr = findEQTMinimum(loc.Lat, loc.Lon, month, year)
		case 'f':
			year = prompter.Int("year", year)
			s, computeErr = findFallEquinox(loc.Lat, loc.Lon, year)
		case 'm':
			month = prompter.Int("month", month)
			day = prompter.Int("day", day)
			year = prompter.Int("year", year)
			s, computeErr = findMeridianTransit(loc.Lat, loc.Lon, month, day, year)
		case 's':
			year = prompter.Int("year", year)
			s, computeErr = findSummerSolstice(loc.Lat, loc.Lon, year)
		case 'v':
			year = prompter.Int("year", year)
			s, computeErr = findVernalEquinox(loc.Lat, loc.Lon, year)
		case 'w':
			year = prompter.Int("year", year)
			s, computeErr = findWinterSolstice(loc.Lat, loc.Lon, year)
		case 'h', ' ':
			printMenu()
			continue
		case 'q':
			return
		default:
			fmt.Println("NOT A VALID OPTION")
			continue
		}

		if computeErr != nil {
			fmt.Fprintln(os.Stderr, "sun:", computeErr)
			continue
		}
		if op[0] == 'a' || op[0] == 'b' {
			s, computeErr = suncomp(loc.Lat, loc.Lon, month, day, year, ut)
			if computeErr != nil {
				fmt.Fprintln(os.Stderr, "sun:", computeErr)
				continue
			}
		}
		month, day, year, ut = s.Month, s.Day, s.Year, s.UT

		printSunDat(s, loc.OST, loc.ODT)
		if s.Alt > 0 {
			fmt.Println("\nfor a horizontal sundial with a vertical gnomon:")
			fmt.Println("P = normal to sundial plane | dial time = clock time + EQT")
			printSundial("horizontal sundial:", computeSundial(s.Alt, s.Az, 90, 0, 90, 0))
			fmt.Println("\nfor a vertical (south facing) sundial with a horizontal gnomon:")
			printSundial("vertical sundial:", computeSundial(s.Alt, s.Az, 0, 180, 0, 180))
		}
		fmt.Println()
	}
}
