// findthrd identifies a "mystery" thread from its measured major
// diameter or pitch, searched against a compiled reference table of
// roughly 470 standard threads (British, metric, and other historic
// and current standards), within a chosen tolerance.
//
// Converted from FINDTHRD.C (M. W. Klotz), WorkshopUtilities/findthrd.
package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"

	"mwkgo/internal/promptio"
	"mwkgo/internal/refdata"
)

// thread is one row of the reference thread table.
type thread struct {
	Name              string
	DiamIn, DiamMM    float64
	PitchTPI, PitchMM float64
}

func loadThreads(ctx context.Context, db *sql.DB) ([]thread, error) {
	rows, err := db.QueryContext(ctx, `SELECT name, diam_in, diam_mm, pitch_tpi, pitch_mm FROM threads`)
	if err != nil {
		return nil, fmt.Errorf("query threads: %w", err)
	}
	defer rows.Close()

	var threads []thread
	for rows.Next() {
		var t thread
		if err := rows.Scan(&t.Name, &t.DiamIn, &t.DiamMM, &t.PitchTPI, &t.PitchMM); err != nil {
			return nil, fmt.Errorf("scan thread: %w", err)
		}
		threads = append(threads, t)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate threads: %w", err)
	}
	return threads, nil
}

// searchField selects which field of a thread a search matches
// against.
type searchField int

const (
	searchDiamIn searchField = iota
	searchDiamMM
	searchPitchTPI
	searchPitchMM
)

func (f searchField) value(t thread) float64 {
	switch f {
	case searchDiamIn:
		return t.DiamIn
	case searchDiamMM:
		return t.DiamMM
	case searchPitchTPI:
		return t.PitchTPI
	default:
		return t.PitchMM
	}
}

// findMatches returns every thread whose chosen field, scaled by
// +/-tolerancePercent, brackets target — matching FINDTHRD.C's own
// comparison of the search value against each table entry scaled by
// tolerance, not the other way around (the two are only
// approximately the same for small tolerances).
func findMatches(threads []thread, field searchField, target, tolerancePercent float64) []thread {
	tol := tolerancePercent * 0.01

	var matches []thread
	for _, t := range threads {
		v := field.value(t)
		if target >= v*(1-tol) && target <= v*(1+tol) {
			matches = append(matches, t)
		}
	}
	return matches
}

func main() {
	ctx := context.Background()
	db, err := refdata.Open(ctx)
	if err != nil {
		fmt.Fprintln(os.Stderr, "findthrd:", err)
		os.Exit(1)
	}
	defer db.Close()

	threads, err := loadThreads(ctx, db)
	if err != nil {
		fmt.Fprintln(os.Stderr, "findthrd:", err)
		os.Exit(1)
	}

	prompter, err := promptio.New(os.Stdin, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, "findthrd:", err)
		os.Exit(1)
	}

	printMenu()
	for {
		switch strings.ToLower(prompter.Line("(A,B,C,D,M,Q) ? ")) {
		case "a":
			runSearch(prompter, threads, searchDiamIn, "Thread diameter (in)", 0.138, "Allowable diameter error (%)")
		case "b":
			runSearch(prompter, threads, searchDiamMM, "Thread diameter (mm)", 6.0, "Allowable diameter error (%)")
		case "c":
			runSearch(prompter, threads, searchPitchTPI, "Thread pitch (tpi)", 32.0, "Allowable pitch error (%)")
		case "d":
			runSearch(prompter, threads, searchPitchMM, "Thread pitch (mm)", 1.0, "Allowable pitch error (%)")
		case "m":
			printMenu()
		case "q":
			return
		default:
			fmt.Println("NOT A VALID OPTION")
		}
	}
}

func printMenu() {
	fmt.Println("A - Find thread from diameter in inches")
	fmt.Println("B - Find thread from diameter in millimeters")
	fmt.Println("C - Find thread from pitch in tpi (threads per inch)")
	fmt.Println("D - Find thread from pitch in millimeters")
	fmt.Println("M - Display this menu")
	fmt.Println("Q - Quit")
	fmt.Println()
}

func runSearch(p *promptio.Prompter, threads []thread, field searchField, valueLabel string, valueDefault float64, tolLabel string) {
	value := p.Float(valueLabel, valueDefault)
	tolerance := p.Float(tolLabel, 1.0)

	fmt.Println("\nName, diam (in), pitch (tpi), diam (mm), pitch (mm)")
	fmt.Println()
	for _, t := range findMatches(threads, field, value, tolerance) {
		fmt.Printf("%s, %.4f, %.4f, %.4f, %.4f\n", t.Name, t.DiamIn, t.PitchTPI, t.DiamMM, t.PitchMM)
	}
}
