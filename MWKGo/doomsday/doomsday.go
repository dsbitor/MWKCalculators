// doomsday computes the day of the week for an input date two ways,
// a Julian day number calculation and Conway's doomsday algorithm,
// and shows that they agree.
//
// Converted from DOOMSDAY.C (M. W. Klotz, 5/01), Misc/doomsday.
package main

import (
	"fmt"
	"os"
	"time"

	"mwkgo/internal/promptio"
)

// dateTransitionJulianDay is the Julian day number on which this
// program switches from the Julian to the Gregorian calendar: 14
// September 1752, the UK and US adoption date, matching the
// original program's own default.
const dateTransitionJulianDay = 2361222

// weekdayNames indexes weekday 0 (Sunday) through 6 (Saturday).
var weekdayNames = [7]string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}

// julianDay computes the Julian day number (days since 1 January
// 4713 BC) for a calendar date. Years before 1 AD are given as
// negative numbers with no year zero: 2 BC is year -1. month is not
// restricted to 1 through 12; a value outside that range rolls into
// an adjacent year, matching the original program's arithmetic.
func julianDay(month, day, year int) int64 {
	y := year
	if y < 0 {
		y++
	}
	m := month - 1

	monthYears := m / 12
	y += monthYears
	m -= monthYears * 12
	if m < 0 {
		m += 12
		y--
	}
	if m < 2 {
		m += 13
		y--
	} else {
		m++
	}

	adjustedYear := int64(y) + 4716
	julian := (1461*adjustedYear)>>2 + int64(153*(m+1)/5) + int64(day) - 1524
	gregorian := julian + 2 - int64(y)/100 + int64(y)/400 - int64(y)/4000
	if gregorian >= dateTransitionJulianDay {
		return gregorian
	}
	return julian
}

// calendarDate is a calendar date recovered from a Julian day number.
type calendarDate struct {
	Month, Day, Year int
	Weekday          int // 0 = Sunday .. 6 = Saturday
	YearDay          int // 0-based day of year
}

// fromJulianDay converts a Julian day number back to a calendar
// date, the inverse of julianDay.
func fromJulianDay(jd int64) calendarDate {
	var a int64
	if jd < dateTransitionJulianDay {
		a = jd
	} else {
		aa := jd - 1721120
		ab := 31 * (aa / 1460969)
		aa %= 1460969
		ab += 3 * (aa / 146097)
		aa %= 146097
		if aa == 146096 {
			ab += 3
		} else {
			ab += aa / 36524
		}
		a = jd + (ab - 2)
	}

	b := a + 1524
	ay := (20*b - 2442) / 7305
	d := (1461 * ay) >> 2
	ee := b - d
	em := 10000 * ee / 306001

	day := int(ee - 306001*em/10000)

	month := int(em) - 1
	if month > 12 {
		month -= 12
	}

	var year int
	if month > 2 {
		year = int(ay) - 4716
	} else {
		year = int(ay) - 4715
	}
	// The original C version returns this internal astronomical year
	// unadjusted, which gives year 0 for 1 BC even though julianDay's
	// forward conversion never accepts year 0 on input (it shifts -1
	// to 0 internally, per the "no year zero" convention its own
	// comments describe). Left unadjusted, the round trip breaks:
	// this exact case, jd 1721420, is one of the worked examples in
	// the original source comments, which expect year -1, not 0.
	// Applying the same adjustment the original already makes a few
	// lines below, for its own year-length lookup, restores the round
	// trip and matches the author's documented and verified example.
	if year <= 0 {
		year--
	}

	weekday := int((jd + 1) % 7)

	// The original C version applies a second "no year zero"
	// decrement here, because it feeds this lookup its own buggy,
	// unadjusted year. Now that year above is already the corrected
	// display year, julianDay's own internal adjustment (it accepts
	// the display convention directly) gives the right start-of-year
	// reference without decrementing a second time.
	yearDay := int(jd - julianDay(1, 1, year))

	return calendarDate{Month: month, Day: day, Year: year, Weekday: weekday, YearDay: yearDay}
}

// centuryDoomsday gives the doomsday weekday for a century's first
// year, indexed by century modulo 4.
var centuryDoomsday = [4]int{2, 0, 5, 3}

// doomsdayWeekday returns the weekday (0 = Sunday .. 6 = Saturday)
// shared by 4/4, 6/6, 8/8, 10/10, 12/12, and several other fixed
// dates in the given year, using Conway's doomsday algorithm.
func doomsdayWeekday(year int) int {
	century := year / 100
	yearInCentury := year % 100
	// Normalise to a non-negative index: Go's % can return a negative
	// result for a negative century (proleptic BC years), which would
	// otherwise index centuryDoomsday out of range instead of the
	// silent undefined behaviour the original C version has for the
	// same input.
	centuryAnchor := centuryDoomsday[((century%4)+4)%4]

	q := yearInCentury / 12
	r := yearInCentury - 12*q
	s := r / 4
	// The original C version's final % 7 has the same negative-result
	// problem as the array index above for a negative year, and
	// normalising it here keeps the returned value a valid weekday
	// index in every case rather than reproducing that quirk too.
	return (((centuryAnchor + q + r + s) % 7) + 7) % 7
}

// monthDoomsday returns the day of the month, in the given month and
// year, that falls on that year's doomsday weekday.
//
// This reproduces the original program's leap-year test for January
// and February exactly: year%4==0 || year%400==0. That test is not
// the true Gregorian rule (which also excludes centurial years not
// divisible by 400, such as 1900) and so gives the wrong doomsdate
// for January and February of a small number of centurial years.
// The discrepancy is confined to this convenience lookup; julianDay
// and fromJulianDay implement the correct Gregorian rule separately
// and are not affected.
func monthDoomsday(month, year int) int {
	switch {
	case month%2 == 0 && month > 2:
		return month
	case month == 3:
		return 7
	case month == 5:
		return 9
	case month == 9:
		return 5
	case month == 7:
		return 11
	case month == 11:
		return 7
	}

	isLeap := year%4 == 0 || year%400 == 0
	switch {
	case month == 1 && isLeap:
		return 25
	case month == 1:
		return 31
	case month == 2 && isLeap:
		return 29
	case month == 2:
		return 28
	}
	return 0
}

// weekdayViaDoomsday returns the weekday (0 = Sunday .. 6 = Saturday)
// of month/day/year using Conway's doomsday algorithm, as an
// independent cross-check against fromJulianDay's result for the
// same date.
func weekdayViaDoomsday(month, day, year int) int {
	doomsday := doomsdayWeekday(year)
	anchorDay := monthDoomsday(month, year)

	weekday := (doomsday + day - anchorDay) % 7
	if weekday < 0 {
		weekday += 7
	}
	return weekday
}

func main() {
	now := time.Now()

	fmt.Printf("Local time:\t%s\n", now.Format(time.RFC1123))
	fmt.Printf("UTC time:\t%s\n", now.UTC().Format(time.RFC1123))
	fmt.Println()

	prompter, err := promptio.New(os.Stdin, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, "doomsday:", err)
		os.Exit(1)
	}

	month := prompter.Int("month", int(now.Month()))
	day := prompter.Int("day", now.Day())
	year := prompter.Int("year", now.Year())

	jd := julianDay(month, day, year)
	date := fromJulianDay(jd)

	fmt.Println()
	fmt.Printf("julian day number = %d\n", jd)
	fmt.Printf("weekday (from julian day) = %s\n", weekdayNames[date.Weekday])
	fmt.Printf("day of year = %d\n", date.YearDay+1)

	fmt.Println()
	fmt.Printf("doomsday for %d = %s\n", year, weekdayNames[doomsdayWeekday(year)])
	fmt.Printf("weekday (from doomsday algorithm) = %s\n", weekdayNames[weekdayViaDoomsday(month, day, year)])
}
