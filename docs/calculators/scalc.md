# scalc

A collection of specialized calculators.

**Converted from:** `SCALC.C` (M. W. Klotz), `MWKC/Math/scalc.zip`
**Go source:** `MWKGo/scalc/scalc.go`

## Purpose

A menu of ten small, independent calculators bundled into one
program: a proportion solver, running parallel-resistance/RSS/RMS
accumulators, running descriptive statistics, degree/DMS conversions
in both directions, machine screw number to major diameter, linear
interpolation, and drawing-scale conversion.

## Inputs

A menu (A-J, Q to quit) selects which calculator to run; each has its
own prompts, described below.

## Output

Varies by calculator; see each below.

## Method

```
A: solve a/b = c/d for whichever of a,b,c,d is unknown, given the other three
B: 1/z = sum(1/xi), a running parallel-combination accumulator
C: rss = sqrt(sum(xi*|xi|)), a running root-sum-square accumulator
D: rms = sqrt(sum(xi*|xi|)/n), a running root-mean-square accumulator
E: running max, min, median, mean, and sample standard deviation
F: degrees/minutes/seconds -> fractional degrees (carrying seconds
   into minutes and minutes into degrees whenever they exceed 60)
G: fractional degrees -> degrees/minutes/seconds (rounding seconds)
H: screw major diameter = 0.06 + screwNumber*0.013
I: linear interpolation through two given points
J: drawing scale = drawingLength / realLength
```

The accumulators (B, C, D) use the original program's own
sign-aware squaring (`x*|x|` rather than `x*x`), which lets typing a
negative value cancel a previously entered positive one of the same
magnitude, matching the original's documented "remove a previously
entered value" feature.

The running-statistics calculator (E) finds the median via the
original program's own in-place quickselect partition; this
conversion sorts a copy of the values instead, an internal
simplification with no effect on the result, since both methods
select the same middle element.

## Worked Example

No worked numeric example was included with the original program. As
independently verifiable checks: two equal resistors in parallel
combine to exactly half of one; a 3-4-5 right triangle's legs give an
RSS of exactly 5; the RMS of a set of identical values is that value;
a well known example statistics dataset's mean and sample standard
deviation match textbook values; converting an angle to
degrees/minutes/seconds and back reproduces the original value; and
a #10 screw's major diameter (0.19 in) matches the standard machine
screw gauge formula. All confirmed in this conversion's tests.
