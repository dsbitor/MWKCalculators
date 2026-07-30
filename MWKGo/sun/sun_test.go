package main

import (
	"math"
	"testing"
)

// exampleLat, exampleLon, exampleOST, exampleODT are SUN.DAT's own
// shipped location (Los Angeles/Long Beach area).
const (
	exampleLat = 33.7775436
	exampleLon = 118.3769558
	exampleOST = -8.0
	exampleODT = -7.0
)

func TestNormab_WrapsIntoRange(t *testing.T) {
	cases := []struct{ x, a, b, want float64 }{
		{370, 0, 360, 10},
		{-10, 0, 360, 350},
		{25, 0, 24, 1},
		{-1, 0, 24, 23},
		{180, 0, 360, 180},
	}
	for _, c := range cases {
		if got := normab(c.x, c.a, c.b); math.Abs(got-c.want) > 1e-9 {
			t.Errorf("normab(%v,%v,%v) = %v, want %v", c.x, c.a, c.b, got, c.want)
		}
	}
}

func TestTnorm_NoOpWithinRange(t *testing.T) {
	m, d, y, ut, err := tnorm(6, 15, 2024, 12.5)
	if err != nil {
		t.Fatal(err)
	}
	if m != 6 || d != 15 || y != 2024 || ut != 12.5 {
		t.Errorf("tnorm(6,15,2024,12.5) = (%d,%d,%d,%v), want unchanged", m, d, y, ut)
	}
}

func TestTnorm_CarriesForwardAcrossMonthAndYear(t *testing.T) {
	// Dec 31, 2023 at 25 hours UT should roll to Jan 1, 2024 at 1:00.
	m, d, y, ut, err := tnorm(12, 31, 2023, 25.0)
	if err != nil {
		t.Fatal(err)
	}
	if m != 1 || d != 1 || y != 2024 || math.Abs(ut-1) > 1e-9 {
		t.Errorf("tnorm(12,31,2023,25) = (%d,%d,%d,%v), want (1,1,2024,1)", m, d, y, ut)
	}
}

func TestTnorm_CarriesBackwardAcrossMonthAndYear(t *testing.T) {
	// Jan 1, 2024 at -1 hour UT should roll back to Dec 31, 2023 at 23:00.
	m, d, y, ut, err := tnorm(1, 1, 2024, -1.0)
	if err != nil {
		t.Fatal(err)
	}
	if m != 12 || d != 31 || y != 2023 || math.Abs(ut-23) > 1e-9 {
		t.Errorf("tnorm(1,1,2024,-1) = (%d,%d,%d,%v), want (12,31,2023,23)", m, d, y, ut)
	}
}

func TestTnorm_RespectsLeapYearFebruary(t *testing.T) {
	// Feb 28, 2024 (a leap year) at 25 hours should roll to Feb 29, not Mar 1.
	m, d, _, _, err := tnorm(2, 28, 2024, 25.0)
	if err != nil {
		t.Fatal(err)
	}
	if m != 2 || d != 29 {
		t.Errorf("tnorm(2,28,2024,25) = (%d,%d,...), want (2,29,...) since 2024 is a leap year", m, d)
	}

	// The same date in a non-leap year should roll to Mar 1.
	m, d, _, _, err = tnorm(2, 28, 2023, 25.0)
	if err != nil {
		t.Fatal(err)
	}
	if m != 3 || d != 1 {
		t.Errorf("tnorm(2,28,2023,25) = (%d,%d,...), want (3,1,...) since 2023 is not a leap year", m, d)
	}
}

func TestTnorm_RejectsInvalidDate(t *testing.T) {
	if _, _, _, _, err := tnorm(13, 1, 2024, 12); err == nil {
		t.Error("expected an error for month 13")
	}
	if _, _, _, _, err := tnorm(2, 30, 2024, 12); err == nil {
		t.Error("expected an error for February 30")
	}
}

func TestCaz_DueNorthAndSouth(t *testing.T) {
	// Same longitude, target due north (higher latitude) => azimuth 0.
	if az := caz(30, -120, 40, -120); math.Abs(az) > 1e-9 {
		t.Errorf("caz due north = %v, want 0", az)
	}
	// Same longitude, target due south (lower latitude) => azimuth 180.
	if az := caz(40, -120, 30, -120); math.Abs(az-180) > 1e-9 {
		t.Errorf("caz due south = %v, want 180", az)
	}
}

func TestSuncomp_RightAscensionAndDeclinationInRange(t *testing.T) {
	s, err := suncomp(exampleLat, exampleLon, 6, 21, 2024, 12)
	if err != nil {
		t.Fatal(err)
	}
	if s.RA < 0 || s.RA >= 24 {
		t.Errorf("RA = %v, want in [0,24)", s.RA)
	}
	if math.Abs(s.Dec) > 23.5 {
		t.Errorf("Dec = %v, want within +/-23.5 degrees (the sun's own declination range)", s.Dec)
	}
}

func TestSuncomp_DistanceStaysNearOneAU(t *testing.T) {
	// Earth's orbital eccentricity is small; distance should always be
	// within a couple percent of 1 AU.
	s, err := suncomp(exampleLat, exampleLon, 1, 3, 2024, 12) // near perihelion
	if err != nil {
		t.Fatal(err)
	}
	if s.Dist < 0.98 || s.Dist > 1.02 {
		t.Errorf("Dist = %v AU, want within 0.98-1.02", s.Dist)
	}
}

func TestFindSummerSolstice_DeclinationNearMaximum(t *testing.T) {
	// The sun's declination peaks at approximately +23.44 degrees at
	// the summer solstice (around June 20-21).
	s, err := findSummerSolstice(exampleLat, exampleLon, 2024)
	if err != nil {
		t.Fatal(err)
	}
	if math.Abs(s.Dec-23.44) > 0.1 {
		t.Errorf("summer solstice declination = %v, want close to +23.44", s.Dec)
	}
	if s.Month != 6 {
		t.Errorf("summer solstice month = %d, want 6", s.Month)
	}
}

func TestFindWinterSolstice_DeclinationNearMinimum(t *testing.T) {
	s, err := findWinterSolstice(exampleLat, exampleLon, 2024)
	if err != nil {
		t.Fatal(err)
	}
	if math.Abs(s.Dec+23.44) > 0.1 {
		t.Errorf("winter solstice declination = %v, want close to -23.44", s.Dec)
	}
	if s.Month != 12 {
		t.Errorf("winter solstice month = %d, want 12", s.Month)
	}
}

func TestFindVernalAndFallEquinox_DeclinationNearZero(t *testing.T) {
	spring, err := findVernalEquinox(exampleLat, exampleLon, 2024)
	if err != nil {
		t.Fatal(err)
	}
	if math.Abs(spring.Dec) > 0.5 {
		t.Errorf("vernal equinox declination = %v, want close to 0", spring.Dec)
	}

	fall, err := findFallEquinox(exampleLat, exampleLon, 2024)
	if err != nil {
		t.Fatal(err)
	}
	if math.Abs(fall.Dec) > 0.5 {
		t.Errorf("fall equinox declination = %v, want close to 0", fall.Dec)
	}
}

func TestSuncomp_SummerDaylightLongerThanWinterAtNorthernLatitude(t *testing.T) {
	summer, err := suncomp(exampleLat, exampleLon, 6, 21, 2024, 12+exampleLon/15)
	if err != nil {
		t.Fatal(err)
	}
	winter, err := suncomp(exampleLat, exampleLon, 12, 21, 2024, 12+exampleLon/15)
	if err != nil {
		t.Fatal(err)
	}
	summerDaylight := summer.UTSet - summer.UTRise
	winterDaylight := winter.UTSet - winter.UTRise
	if summerDaylight <= winterDaylight {
		t.Errorf("summer daylight (%v hr) should exceed winter daylight (%v hr) at a northern latitude", summerDaylight, winterDaylight)
	}
}

func TestFindMeridianTransit_AzimuthMatchesTarget(t *testing.T) {
	s, err := findMeridianTransit(exampleLat, exampleLon, 6, 21, 2024)
	if err != nil {
		t.Fatal(err)
	}
	azOver := 180.0
	if exampleLat <= 0 {
		azOver = 0.0
	}
	if math.Abs(s.Az-azOver) > 0.001 {
		t.Errorf("meridian transit azimuth = %v, want within 0.001 of %v", s.Az, azOver)
	}
}

func TestComputeSundial_SunAtZenithCastsNoShadow(t *testing.T) {
	// A sun directly overhead (alt=90) shining on a vertical gnomon
	// standing in a horizontal plane casts an essentially zero-length
	// shadow: the shadow/gnomon-length ratio should be vanishingly
	// small. (delta's own zero-shadow sentinel branch, for a plane
	// exactly edge-on to the sun, is not separately tested here: since
	// this port uses float64 rather than SUN.C's own 4-byte float,
	// residues like cos(90 deg) never land on an exact 0 in either
	// precision, so that branch is unreachable through realistic
	// degree inputs in practice — it is kept only as the same
	// defensive guard the original has.)
	d := computeSundial(90, 0, 90, 0, 90, 0)
	if math.Abs(d.Ratio) > 1e-6 {
		t.Errorf("shadow/gnomon ratio with the sun at zenith = %v, want ~0", d.Ratio)
	}
}

func TestHmsAndDms_Formatting(t *testing.T) {
	if got := hms(13.5); got != "13:30:00" {
		t.Errorf("hms(13.5) = %q, want 13:30:00", got)
	}
	if got := dms(45.5); got != "  45:30:00" {
		t.Errorf("dms(45.5) = %q, want '  45:30:00' (a sign space, then a %%3d-padded degree field)", got)
	}
	if got := dms(-45.5); got != "- 45:30:00" {
		t.Errorf("dms(-45.5) = %q, want '- 45:30:00'", got)
	}
}

func TestPstpdt_MarksPreviousDayWhenNegative(t *testing.T) {
	st, dt := pstpdt(2, -8, -7)
	if st != "<ST = 18:00:00" {
		t.Errorf("pstpdt standard = %q, want previous-day marker", st)
	}
	if dt != "<DT = 19:00:00" {
		t.Errorf("pstpdt daylight = %q, want previous-day marker", dt)
	}

	st, dt = pstpdt(12, -8, -7)
	if st != "ST = 04:00:00" || dt != "DT = 05:00:00" {
		t.Errorf("pstpdt(12,-8,-7) = (%q,%q), want no previous-day marker", st, dt)
	}
}
