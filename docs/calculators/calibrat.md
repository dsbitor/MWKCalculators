# calibrat

Calibrate a linear scale.

**Converted from:** `CALIBRAT.C` (M. W. Klotz), `MWKC/Math/calibrat.zip`
**Go source:** `MWKGo/calibrat/calibrat.go`

## Purpose

Given a set of checkpoints for an instrument with a linear response
(a truth value and what the instrument actually read at that value),
this finds the best-fit calibration equation `y = A*x + B` (x = truth,
y = measured) by least squares, and can print a calibration table
over a chosen range.

## Data setup

The checkpoint data is specific to one calibration run of one
instrument — not universal reference data, and not a fixed piece of
reusable equipment configuration — so, unlike the reference- and
machine-config-bucket programs elsewhere in this project, `calibrat`
reads its input fresh from a file named on the command line each run,
in the same `STARTOFDATA`/`ENDOFDATA` text format the original
program used, rather than from either SQLite database; see
`ai/plans/c-to-go-conversion-plan.md`, "Data-file strategy for
Tier 2".

```
calibrat -data my-calibration-run.dat
```

The file needs one `truth,measured` pair per line (comma or tab
separated). A worked example built from the original `CALIBRAT.DAT`'s
own default data (the thermometer example below) ships at
`MWKGo/calibrat/testdata/example.dat`.

## Inputs

| Prompt | Default |
|---|---|
| Do you want to construct a calibration table [Y]/N ? | Y |
| Table starting truth value | the smallest truth value in the data |
| Table ending truth value | the largest truth value in the data |
| Table increment | (end − start) / number of data pairs |

## Output

The least-squares slope (A) and intercept (B); optionally, two
tables over the chosen range — one treating each value as a measured
reading and inverting the fit to recover truth, the other treating
it as a truth value and predicting what would be measured.

## Method

Ordinary least squares on `y = A*x + B`. At least two data pairs are
required; the check pairs are sorted by truth value before use (this
also fixes the table's default start/end values, which are the
smallest and largest truth values). A determinant of zero (all data
pairs sharing the same truth value) is reported as an error, matching
the original's own check, reproducing its message verbatim (including
its "Determinat" typo, kept rather than silently corrected, the same
treatment given to other characterful original text such as `vrev`'s
"Try again, dummy.").

## Worked Example

The original program's own worked example is a thermometer reading
2°C in ice water (truth 0) and 102°C in boiling water (truth 100),
checked against a calibrated oven at 50°C reading 51°C (truth 50).
This conversion's tests reproduce it exactly: `A = 1.000000`,
`B = 1.666667`, and the calibration table at truth values 0, 50, 100
(`0.000000 <=> -1.666667`, `50.000000 <=> 48.333333`,
`100.000000 <=> 98.333333` for the measured-to-truth table; `0.000000
<=> 1.666667`, `50.000000 <=> 51.666667`, `100.000000 <=> 101.666667`
for truth-to-measured).
