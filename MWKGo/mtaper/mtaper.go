// mtaper measures the included angle of a tapered cylindrical part
// using a V-block of known angle set on a sine bar, per Machinery's
// Handbook's procedure (23rd ed., pg. 685): the sine bar stack is
// adjusted until an indicator swept across the part shows no
// deviation, and the part's included angle follows from the V-block
// angle and the sine bar angle. It also decomposes a sine bar stack
// height into blocks from a standard 81-piece gage block set.
//
// Converted from MTAPER.C, WorkshopUtilities/mtaper.
package main

import (
	"fmt"
	"math"
	"os"
	"strings"

	"mwkgo/internal/promptio"
)

func sind(deg float64) float64 { return math.Sin(deg * math.Pi / 180) }
func cosd(deg float64) float64 { return math.Cos(deg * math.Pi / 180) }
func asind(x float64) float64  { return math.Asin(x) * 180 / math.Pi }
func atand(x float64) float64  { return math.Atan(x) * 180 / math.Pi }

// taperResult holds a fully solved taper measurement: the sine bar
// angle B, the part's included angle A, and the angle C between the
// part centerline and the sine bar surface.
type taperResult struct {
	StackHeight, SineBarAngle, PartAngle, AngleC float64
}

// taperAngleFromStack returns the part's included angle given the
// sine bar stack height, sine bar length, and V-block included
// angle.
func taperAngleFromStack(stackHeight, barLength, vBlockAngle float64) taperResult {
	sineBarAngle := asind(stackHeight / barLength)
	d := cosd(sineBarAngle) + 1/sind(vBlockAngle/2)
	t := sind(sineBarAngle) / d
	partAngle := 2 * atand(t)
	angleC := asind(sind(partAngle/2) / sind(vBlockAngle/2))
	return taperResult{StackHeight: stackHeight, SineBarAngle: sineBarAngle, PartAngle: partAngle, AngleC: angleC}
}

// stackFromTaperAngle returns the sine bar stack height given the
// part's included angle, sine bar length, and V-block included
// angle.
func stackFromTaperAngle(partAngle, barLength, vBlockAngle float64) taperResult {
	angleC := asind(sind(partAngle/2) / sind(vBlockAngle/2))
	sineBarAngle := angleC + partAngle/2
	stackHeight := barLength * sind(sineBarAngle)
	return taperResult{StackHeight: stackHeight, SineBarAngle: sineBarAngle, PartAngle: partAngle, AngleC: angleC}
}

// decimalDigit returns the digit k places to the right of the
// decimal point in x (k=1 is tenths, k=4 is ten-thousandths), or 0
// if k is 0. Implemented via the same string-formatting approach as
// the original program's own kp function, rather than pure
// arithmetic (modulo on floating point decimal digits is easy to get
// subtly wrong at exactly the precision this function cares about).
func decimalDigit(x float64, k int) int {
	if k == 0 {
		return 0
	}
	s := fmt.Sprintf("%-+20.9f", x)
	dot := strings.IndexByte(s, '.')
	idx := dot + k
	if idx < 0 || idx >= len(s) || s[idx] < '0' || s[idx] > '9' {
		return 0
	}
	return int(s[idx] - '0')
}

// gaugeBlocksFor decomposes stack (truncated to 4 decimal places)
// into a sequence of blocks from a standard 81-piece gage block set,
// via the original program's own greedy digit-extraction algorithm:
// remove the ten-thousandths digit as a block in the 0.1000-0.1009
// range, then the thousandths and hundredths digits together as a
// block in the 0.101-0.150 range, then the tenths and remaining
// hundredths as a block in 0.05 increments up to 1.00, then whatever
// whole-inch blocks (1, 2, 3, 4) remain. It returns the truncated
// target, the blocks chosen, and whether they sum to it exactly.
func gaugeBlocksFor(stack float64) (target float64, blocks []float64, ok bool) {
	target = float64(int(1e4*stack)) * 1e-4
	if target < 0.05 {
		return target, nil, false
	}

	r := target
	takeBlock := func(t float64) bool {
		r -= t
		blocks = append(blocks, t)
		return math.Abs(r) < 1e-5
	}

	if tenThousandths := decimalDigit(r, 4); tenThousandths != 0 {
		if takeBlock(0.1000 + float64(tenThousandths)*0.0001) {
			return target, blocks, true
		}
		if r < 0.04999 {
			return target, blocks, false
		}
	}

	if thousandths := decimalDigit(r, 3); thousandths != 0 {
		hundredths := decimalDigit(r, 2)
		if hundredths >= 5 {
			hundredths -= 5
		}
		if takeBlock(0.100 + float64(hundredths)*0.01 + float64(thousandths)*0.001) {
			return target, blocks, true
		}
		if r < 0.04999 {
			return target, blocks, false
		}
	}

	if hundredths, tenths := decimalDigit(r, 2), decimalDigit(r, 1); hundredths != 0 || tenths != 0 {
		if takeBlock(0.100*float64(tenths) + float64(hundredths)*0.01) {
			return target, blocks, true
		}
		if r < 0.04999 {
			return target, blocks, false
		}
	}

	if whole := int(math.Floor(r + 1e-5)); whole != 0 {
		for i := 4; i > 0; i-- {
			if i <= whole {
				whole -= i
				if takeBlock(float64(i)) {
					return target, blocks, true
				}
				if r < 0.04999 {
					return target, blocks, false
				}
			}
		}
	}

	return target, blocks, math.Abs(r) < 1e-5
}

func main() {
	prompter, err := promptio.New(os.Stdin, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, "mtaper:", err)
		os.Exit(1)
	}

	fmt.Println("CONICAL PART TAPER MEASUREMENT")
	fmt.Println()

	vBlockAngle := prompter.Float("V-block included angle (D)", 90.0)
	barLength := prompter.Float("Sine bar length (L)", 5.0)

	for {
		fmt.Println("\n1  Compute part included angle given sine bar stack height")
		fmt.Println("2  Compute sine bar stack height given part included angle")
		fmt.Println("Q  Quit")
		choice := strings.ToLower(prompter.Line("(1,2,Q) ? "))

		var result taperResult
		switch choice {
		case "1":
			h := prompter.Float("Stack Height", 0.9947)
			result = taperAngleFromStack(h, barLength, vBlockAngle)
		case "2":
			a := prompter.Float("Part included angle", 9.5)
			result = stackFromTaperAngle(a, barLength, vBlockAngle)
		case "q":
			return
		default:
			fmt.Println("NOT A VALID OPTION")
			continue
		}

		fmt.Printf("\nSine bar length = %.4f in\n", barLength)
		fmt.Printf("Sine bar stack height = %.4f in\n", result.StackHeight)
		fmt.Printf("Sine bar angle = %.4f deg\n", result.SineBarAngle)
		fmt.Printf("Taper included angle (A) = %.4f deg\n", result.PartAngle)
		fmt.Printf("Taper half angle (A/2) = %.4f deg\n", result.PartAngle/2)
		fmt.Printf("Angle C = %.4f deg\n", result.AngleC)

		if choice == "2" {
			printGaugeBlocks(result.StackHeight)
		}
	}
}

func printGaugeBlocks(stack float64) {
	target, blocks, ok := gaugeBlocksFor(stack)
	fmt.Printf("\nBlocks from standard 81 gage block set needed to form stack = %.4f in:\n\n", target)
	r := target
	for _, b := range blocks {
		r -= b
		fmt.Printf("block = %.4f  remainder = %.4f\n", b, r)
	}
	if !ok {
		fmt.Println("\nCAN'T CONTINUE")
	}
}
