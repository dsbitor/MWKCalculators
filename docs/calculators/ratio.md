# ratio

American National Standard for preferred numbers.

**Converted from:** `RATIO.C` (M. W. Klotz), `MWKC/WorkshopUtilities/ratio.zip`
**Go source:** `MWKGo/ratio/ratio.go`

## Purpose

Generates a "preferred numbers" series (a Renard series, the American
National Standard for choosing a sensible spread of sizes): a
geometric progression from 1 up to (but not including) 10, in a
chosen number of steps, each consecutive pair separated by a fixed
percentage. Used for choosing gear ratios, resistor values, packaging
sizes, or any other range of sizes where you want evenly-spaced
options without gaps too coarse or too fine to be useful.

This was misclassified during Tier 2's preparatory scan: its own
`dfile` variable (`"DATA.DAT"`) is dead code, never actually passed
to `fopen` — only an *output* report file is opened — so it belongs
with the Tier 1 no-data-file calculators; see
`ai/plans/c-to-go-conversion-plan.md`'s clean-up phase note. (The
source tree also had a byte-for-byte duplicate extraction of this
same archive under a `ratio 2` directory name — not a second,
distinct program, just a stray duplicate copy.) Like several other
Tier 1 programs, this drops the original's file-save feature and
prints straight to stdout.

## Inputs

| Prompt | Default |
|---|---|
| Option (1-5 for a standard series, 6 for a custom length) | 2 |
| How many numbers in series (option 6 only) | 15 |
| Scale factor | 10 |

| Option | Series | Step size |
|---|---|---|
| 1 | R5 | ~58% |
| 2 | R10 | ~26% |
| 3 | R20 | ~12% |
| 4 | R40 | ~6% |
| 5 | R80 | ~3% |
| 6 | custom | depends on the chosen length |

## Output

The exact step percentage between consecutive numbers, then a table
of every number in the series: its raw value (starting at 1), the
value after applying the scale factor, and that scaled value rounded
to the nearest integer.

## Method

Step `i` of an `m`-step series is `10^(i/m)`, so the series runs from
`10^0 = 1` up to (not including) `10^((m-1)/m)`, each step a constant
multiplicative factor of `10^(1/m)` — and hence a constant
*percentage* step, `100 × (10^(1/m) − 1)`. Rounding to the nearest
integer uses `RATIO.C`'s own exact rule: a fractional part strictly
greater than 0.5 rounds up, and a fractional part of exactly 0.5 (or
less) rounds down — round-half-*down*, not the more familiar
round-half-up.

## Worked Example

The program's own menu text states each standard series' approximate
step percentage directly (5 numbers ~58% apart, 10 ~26%, 20 ~12%, 40
~6%, 80 ~3%) — a worked example baked into the program's own
interface. This conversion's
tests confirm the step-percentage formula against each of those five
values, and separately confirm the R10 series' own 6th entry equals
√10 exactly (a well known constant, since `10^(5/10) = 10^0.5`). A
manual run of the R10 series scaled by 10 reproduces the classic
industrial R10 preferred-number sequence exactly: 10, 12.5, 16, 20,
25, 31.5, 40, 50, 63, 80.
