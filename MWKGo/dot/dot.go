// dot computes depth-of-thread figures for single point threading
// on a lathe: the four standard combinations of sharp/flat crest and
// sharp/flat root depths, both measured perpendicular to the thread
// axis and along the compound feed at a given compound rest angle,
// plus advice on which lines of the threading dial can be used for
// the given threads-per-inch.
//
// Converted from DOT.C (M. W. Klotz, 2/99, 12/02, 3/04, 6/05),
// WorkshopUtilities/dot.
package main

import (
	"fmt"
	"math"
	"os"

	"mwkgo/internal/promptio"
)

// threadDepths holds the four standard depth-of-thread figures
// (labeled A-D to match the original program's own output) and the
// doubled sharp-crest-to-sharp-root figure some lathe operators use
// directly when infeeding on the cross slide rather than the
// compound.
type threadDepths struct {
	SharpCrestSharpRoot float64 // A
	FlatCrestFlatRoot   float64 // B
	SharpCrestFlatRoot  float64 // C
	FlatCrestSharpRoot  float64 // D
	DoubleSharpToSharp  float64 // E
}

// computeThreadDepths returns the standard depth-of-thread figures
// for a thread of the given included angle and pitch. For a 60
// degree thread, SharpCrestSharpRoot reduces to the well known
// H = 0.866025*pitch constant from ANSI B1.1.
func computeThreadDepths(threadAngleDeg, pitch float64) threadDepths {
	h := 0.5 * pitch / math.Tan(0.5*threadAngleDeg*math.Pi/180)
	return threadDepths{
		SharpCrestSharpRoot: h,
		FlatCrestFlatRoot:   0.625 * h,
		SharpCrestFlatRoot:  0.75 * h,
		FlatCrestSharpRoot:  0.875 * h,
		DoubleSharpToSharp:  2 * h,
	}
}

// threadingDialHint reports which lines of a lathe's threading dial
// can be used to pick up the thread, based on the threads-per-inch:
// any line works for an even, whole tpi; only numbered lines for an
// odd, whole tpi; and only odd-numbered lines whenever tpi has a
// fractional part.
func threadingDialHint(threadsPerInch float64) string {
	whole, fraction := math.Modf(threadsPerInch)
	if fraction != 0 {
		return "use any odd-numbered line on threading dial"
	}
	if int64(whole)%2 != 0 {
		return "use any numbered line on threading dial"
	}
	return "use any line on threading dial"
}

func main() {
	prompter, err := promptio.New(os.Stdin, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, "dot:", err)
		os.Exit(1)
	}

	fmt.Println("DEPTH OF THREAD CALCULATIONS")
	fmt.Println()

	threadAngle := prompter.Float("Thread angle (deg)", 60.0)
	threadsPerInch := prompter.Float("Threads per inch", 20.0)
	compoundAngle := prompter.Float("Compound rest angle (deg)", 29.0)
	cosCompoundAngle := math.Cos(compoundAngle * math.Pi / 180)
	pitch := 1 / threadsPerInch

	d := computeThreadDepths(threadAngle, pitch)

	fmt.Println()
	fmt.Printf("thread angle = %.2f deg\n", threadAngle)
	fmt.Printf("threads per inch = %.1f", threadsPerInch)
	fmt.Printf("   (pitch = %.5f in/thread)\n", pitch)
	fmt.Printf("{compound feed at compound angle = %.1f deg}\n\n", compoundAngle)
	fmt.Printf("(A) dot sharp crest - sharp root = %.5f in {%.5f in}\n", d.SharpCrestSharpRoot, d.SharpCrestSharpRoot/cosCompoundAngle)

	if threadAngle != 60.0 {
		fmt.Println("\n***********************")
		fmt.Println("NB... For thread forms other than American National, the following")
		fmt.Println("information may not be valid unless the root and crest flats are defined")
		fmt.Println("as in the American National form.")
		fmt.Println("***********************")
		fmt.Println()
	}

	fmt.Printf("(B) dot flat  crest - flat  root = %.5f in {%.5f in}\n", d.FlatCrestFlatRoot, d.FlatCrestFlatRoot/cosCompoundAngle)
	fmt.Printf("(C) dot sharp crest - flat  root = %.5f in {%.5f in}\n", d.SharpCrestFlatRoot, d.SharpCrestFlatRoot/cosCompoundAngle)
	fmt.Printf("(D) dot flat  crest - sharp root = %.5f in {%.5f in}\n", d.FlatCrestSharpRoot, d.FlatCrestSharpRoot/cosCompoundAngle)
	fmt.Printf("(E) double dot sharp crest - sharp root = %.5f in\n", d.DoubleSharpToSharp)

	z := (3.0 / 16.0) * pitch / math.Tan(30*math.Pi/180)
	fmt.Printf("\nFor American National (60 deg) thread form, subtract %.4f in from\n", 2*z)
	fmt.Println("major diameter (assumes p/8 flat on crest) to obtain pitch diameter")
	fmt.Println()

	fmt.Println(threadingDialHint(threadsPerInch))
}
