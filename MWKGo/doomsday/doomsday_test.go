package main

import "testing"

// These worked examples are taken directly from the comments at the
// end of DOOMSDAY.C, which the original author verified by hand
// against the Julian date technique; they are the closest available
// evidence of correct behaviour for this program.

func TestJulianDay_WorkedExamples(t *testing.T) {
	tests := []struct {
		name             string
		month, day, year int
		want             int64
	}{
		{name: "epoch of the Julian day count", month: 1, day: 1, year: -4713, want: 0},
		{name: "a modern date", month: 4, day: 14, year: 1998, want: 2450918},
		{name: "a BC date", month: 9, day: 10, year: -2699, want: 735866},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := julianDay(tt.month, tt.day, tt.year); got != tt.want {
				t.Errorf("julianDay(%d, %d, %d) = %d, want %d", tt.month, tt.day, tt.year, got, tt.want)
			}
		})
	}
}

func TestFromJulianDay_WorkedExamples(t *testing.T) {
	tests := []struct {
		name string
		jd   int64
		want calendarDate
	}{
		{
			name: "12/28/-1, a Tuesday",
			jd:   1721420,
			want: calendarDate{Month: 12, Day: 28, Year: -1, Weekday: 2, YearDay: 362},
		},
		{
			name: "1/1/1, a Saturday",
			jd:   1721424,
			want: calendarDate{Month: 1, Day: 1, Year: 1, Weekday: 6, YearDay: 0},
		},
		{
			name: "12/7/1941, a Sunday: the Pearl Harbor attack",
			jd:   2430336,
			want: calendarDate{Month: 12, Day: 7, Year: 1941, Weekday: 0, YearDay: 340},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := fromJulianDay(tt.jd)
			if got != tt.want {
				t.Errorf("fromJulianDay(%d) = %+v, want %+v", tt.jd, got, tt.want)
			}
		})
	}
}

func TestJulianDayRoundTrip_WorkedExamples(t *testing.T) {
	dates := []struct {
		month, day, year int
	}{
		{12, 28, -1},
		{1, 1, 1},
		{12, 7, 1941},
		{4, 14, 1998},
	}

	for _, d := range dates {
		jd := julianDay(d.month, d.day, d.year)
		got := fromJulianDay(jd)
		if got.Month != d.month || got.Day != d.day || got.Year != d.year {
			t.Errorf("fromJulianDay(julianDay(%d, %d, %d)) = {%d, %d, %d}, want the original date back",
				d.month, d.day, d.year, got.Month, got.Day, got.Year)
		}
	}
}

func TestWeekdayViaDoomsday_AgreesWithJulianDayMethod(t *testing.T) {
	// This cross-check is the program's whole purpose, per its own
	// companion text: "the two DOWs so-computed should agree". It
	// only holds for genuinely Gregorian dates, though: the doomsday
	// algorithm has no notion of the Julian-to-Gregorian switch that
	// julianDay applies before 1752, so it is not expected to agree
	// with the julian day method for 12/28/-1 or 1/1/1, both of
	// which julianDay correctly treats as Julian calendar dates.
	dates := []struct {
		name             string
		month, day, year int
	}{
		{name: "12/7/1941", month: 12, day: 7, year: 1941},
		{name: "4/14/1998", month: 4, day: 14, year: 1998},
	}

	for _, d := range dates {
		t.Run(d.name, func(t *testing.T) {
			want := fromJulianDay(julianDay(d.month, d.day, d.year)).Weekday
			got := weekdayViaDoomsday(d.month, d.day, d.year)
			if got != want {
				t.Errorf("weekdayViaDoomsday(%d, %d, %d) = %s, want %s (from the julian day method)",
					d.month, d.day, d.year, weekdayNames[got], weekdayNames[want])
			}
		})
	}
}

func TestDoomsdayWeekday_NegativeCenturyDoesNotPanic(t *testing.T) {
	// year=-4713 has a negative century; this must return a valid
	// index rather than panicking, unlike the original C version's
	// undefined behaviour for the same input.
	got := doomsdayWeekday(-4713)
	if got < 0 || got > 6 {
		t.Errorf("doomsdayWeekday(-4713) = %d, want a value in [0,6]", got)
	}
}
