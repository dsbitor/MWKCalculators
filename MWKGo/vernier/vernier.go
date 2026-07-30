// vernier designs a two-plate angular vernier: two plates, with n1
// and n2 holes respectively, mounted concentrically so that some
// hole in each plate can always be brought into alignment (measured
// through a fixed reference) to mark off any of a chosen number of
// equal angular subdivisions, using far fewer total holes than a
// single plate with that many subdivisions would need directly.
//
// Converted from VERNIER.C, WorkshopUtilities/vernier. The original
// writes its results to VERNIER.OUT via fopen; this conversion
// prints to stdout instead and drops that DOS file-save-then-page
// convenience, the same approach used for loan (Tier 1 suitability
// review, Finding 5).
package main

import (
	"errors"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"

	"mwkgo/internal/mwkfmt"
	"mwkgo/internal/promptio"
)

// primeFactorization is a number's prime factors and their
// exponents, in ascending order of factor.
type primeFactorization struct {
	Factors []uint64
	Exps    []int
}

// factorize returns the prime factorization of n. The original
// program builds the same factorization using a sieve of small
// primes (a performance optimization for the DOS-era hardware it
// targeted); this conversion uses plain trial division instead,
// which produces an identical factorization without needing to
// reproduce the sieve, since the sieve is an internal performance
// detail with no effect on the program's output.
func factorize(n uint64) primeFactorization {
	var f primeFactorization
	add := func(factor uint64) {
		if len(f.Factors) > 0 && f.Factors[len(f.Factors)-1] == factor {
			f.Exps[len(f.Exps)-1]++
		} else {
			f.Factors = append(f.Factors, factor)
			f.Exps = append(f.Exps, 1)
		}
	}
	for n%2 == 0 {
		add(2)
		n /= 2
	}
	for factor := uint64(3); factor*factor <= n; factor += 2 {
		for n%factor == 0 {
			add(factor)
			n /= factor
		}
	}
	if n > 1 {
		add(n)
	}
	return f
}

func (f primeFactorization) String() string {
	var b strings.Builder
	for i, factor := range f.Factors {
		if i > 0 {
			b.WriteString(" x ")
		}
		b.WriteString(mwkfmt.GroupedUint(factor))
		if f.Exps[i] > 1 {
			fmt.Fprintf(&b, "^%d", f.Exps[i])
		}
	}
	return b.String()
}

// holeCountsAreCompatible reports whether plates of n1 and n2 holes
// can produce every one of numDivisions equal angular subdivisions:
// for each subdivision, some pair of holes (one per plate) must
// align, measured through a fixed reference, to within 0.0001
// degrees.
func holeCountsAreCompatible(n1, n2, numDivisions uint64) bool {
	divs := 360.0 / float64(numDivisions)
	div1 := 360.0 / float64(n1)
	div2 := 360.0 / float64(n2)

	for nd := uint64(1); nd <= n1; nd++ {
		angle := float64(nd) * divs
		if !alignmentExists(n1, n2, div1, div2, angle) {
			return false
		}
	}
	return true
}

// alignmentExists reports whether some pair of holes (i1, i2) brings
// the two plates into alignment at the given reference angle.
func alignmentExists(n1, n2 uint64, div1, div2, angle float64) bool {
	for i1 := uint64(0); i1 < n1; i1++ {
		p1 := float64(i1) * div1
		for i2 := uint64(0); i2 < n2; i2++ {
			p2 := float64(i2)*div2 + angle
			if p2 >= 360 {
				p2 -= 360
			}
			if math.Abs(p2-360) < 0.0001 {
				p2 = 0
			}
			if math.Abs(p2-p1) < 0.0001 {
				return true
			}
		}
	}
	return false
}

// maxHoleCountSearchCalls bounds the number of times
// holeCountsAreCompatible is called while searching for the best
// n1/n2 pair, replacing the original program's interactive keypress
// abort per coding-style.md Rule 2. A realistic numDivisions (a few
// hundred) needs well under a hundred such calls, since each
// successful pair immediately tightens the pruning bound for every
// later candidate.
const maxHoleCountSearchCalls = 200_000

var errVernierSearchTruncated = errors.New("search stopped early: evaluation limit reached")

// vernierPlates is the best pair of plate hole counts found for a
// given number of divisions.
type vernierPlates struct {
	Holes1, Holes2, TotalHoles uint64
}

// findVernierPlates searches for the pair of plate hole counts
// (n1, n2) satisfying 1/n1 - 1/n2 = 1/numDivisions with the smallest
// total hole count n1+n2, among those that can actually produce
// every division (per holeCountsAreCompatible). It returns ok=false
// if no such pair exists.
func findVernierPlates(numDivisions uint64) (plates vernierPlates, ok bool, err error) {
	bestTotal := uint64(math.MaxUint64)
	calls := 0

	for n1 := uint64(2); n1 < numDivisions; n1++ {
		n2 := numDivisions * n1 / (numDivisions - n1)
		if n2 == 0 || n1 == n2 || n2 == 1 {
			continue
		}
		if n1+n2 >= bestTotal {
			continue
		}

		calls++
		if calls > maxHoleCountSearchCalls {
			return plates, ok, errVernierSearchTruncated
		}
		if !holeCountsAreCompatible(n1, n2, numDivisions) {
			continue
		}

		bestTotal = n1 + n2
		plates = vernierPlates{Holes1: n1, Holes2: n2, TotalHoles: bestTotal}
		ok = true
	}
	return plates, ok, nil
}

// divisionAlignment is which pair of holes aligns a given division,
// and at what angle.
type divisionAlignment struct {
	Division uint64
	AngleDeg float64
	Label    string
}

// holeLabel returns the letter-based label for hole index i1 on the
// "lettered" plate, matching the original program's own scheme: the
// 52 letters A-Z, a-z, then a trailing apostrophe added for every
// additional 52 holes beyond that.
func holeLabel(i1 uint64) string {
	const alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	return string(alphabet[i1%52]) + strings.Repeat("'", int(i1/52))
}

// computeAlignments returns, for every division from 0 to
// numDivisions inclusive, which pair of holes aligns it, along with
// which holes on each plate are used by at least one division.
func computeAlignments(n1, n2, numDivisions uint64) (alignments []divisionAlignment, used1, used2 []bool) {
	divs := 360.0 / float64(numDivisions)
	div1 := 360.0 / float64(n1)
	div2 := 360.0 / float64(n2)
	used1 = make([]bool, n1)
	used2 = make([]bool, n2)

	for nd := uint64(0); nd <= numDivisions; nd++ {
		angle := float64(nd) * divs
	search:
		for i1 := uint64(0); i1 < n1; i1++ {
			p1 := float64(i1) * div1
			for i2 := uint64(0); i2 < n2; i2++ {
				p2 := float64(i2)*div2 + angle
				if p2 >= 360 {
					p2 -= 360
				}
				if math.Abs(p2-360) < 0.0001 {
					p2 = 0
				}
				if math.Abs(p2-p1) < 0.0001 {
					label := holeLabel(i1) + strconv.FormatUint(i2+1, 10)
					alignments = append(alignments, divisionAlignment{Division: nd, AngleDeg: angle, Label: label})
					used1[i1] = true
					used2[i2] = true
					break search
				}
			}
		}
	}
	return alignments, used1, used2
}

func countUsed(used []bool) int {
	n := 0
	for _, u := range used {
		if u {
			n++
		}
	}
	return n
}

func main() {
	prompter, err := promptio.New(os.Stdin, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, "vernier:", err)
		os.Exit(1)
	}

	fmt.Println("TWO PLATE ANGULAR VERNIER COMPUTATIONS")
	fmt.Println()

	num := uint64(prompter.Int("Number of angular subdivisions required", 360))
	if num <= 1 {
		fmt.Println("MUST BE GREATER THAN 1")
		return
	}

	f := factorize(num)
	fmt.Printf("%s = %s\n", mwkfmt.GroupedUint(num), f.String())

	plates, ok, err := findVernierPlates(num)
	if err != nil {
		fmt.Fprintln(os.Stderr, "vernier:", err)
		os.Exit(1)
	}
	if !ok {
		fmt.Println("NO SOLUTION WAS FOUND")
		return
	}

	fmt.Printf("selected %d,%d\n", plates.Holes1, plates.Holes2)
	if plates.TotalHoles > num {
		fmt.Printf("\nTwo plate vernier requires drilling more holes (%d) \n", plates.TotalHoles)
		fmt.Println("than a single plate. Use a single plate.")
		return
	}

	alignments, used1, used2 := computeAlignments(plates.Holes1, plates.Holes2, num)

	fmt.Printf("DATA FOR %d DIVISION TWO PLATE VERNIER\n\n", num)
	fmt.Printf("Total number of holes = %d\n\n", plates.Holes1+plates.Holes2)
	fmt.Printf("Assume holes in plate with %d holes are labeled with LETTERS.\n", plates.Holes1)
	fmt.Printf("Assume holes in plate with %d holes are labeled with NUMBERS.\n", plates.Holes2)
	fmt.Println("Either plate can be the movable plate.")
	fmt.Println("Assume 'A' and '1' holes are aligned initially.")
	fmt.Println("To reverse the direction of rotation, read the list backwards.")
	fmt.Println()
	fmt.Println("Division  Angle(deg)  Holes to Align")
	fmt.Println()
	for _, a := range alignments {
		fmt.Printf("%8d  %10.4f  %14s\n", a.Division, a.AngleDeg, a.Label)
	}
	fmt.Printf("\nTotal holes to drill on letter plate = %d\n", countUsed(used1))
	fmt.Printf("Total holes to drill on number plate = %d\n", countUsed(used2))
}
