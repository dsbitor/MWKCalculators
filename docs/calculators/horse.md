# horse

Torque, rotational speed, and horsepower solver.

**Converted from:** `HORSE.C` (M. W. Klotz, 4/01),
`MWKC/WorkshopUtilities/horse.zip`
**Go source:** `MWKGo/horse/horse.go`

## Purpose

Torque, rotational speed, and horsepower are related by a single
standard formula. Given any two of the three, this program
solves for the third. If fewer than two are given, it reports
that there isn't enough data for a solution rather than guessing.

No `.TXT` file was included with the original program; this
purpose statement is drawn from the `.C` file's own header
comment.

## Inputs

| Prompt | Default |
|---|---|
| Torque | 0 ft-lb (press return if unknown) |
| Rotational speed | 0 rpm (press return if unknown) |
| Horsepower | 0 hp, only asked if torque and speed aren't both already known |

Entering nothing (or 0) for a value means "not known"; the
program needs exactly two known values to solve for the third.

## Output

Torque (in both ft-lb and in-lb), rotational speed, and
horsepower.

## Method

```
horsepower = torque * rpm / 5252
```

solved for whichever of the three is missing, using the other
two.

## Worked Example

No worked example was included with the original program. As an
independently verifiable check: 10 ft-lb of torque at 2626 rpm
gives exactly 5 horsepower, a direct application of the standard
formula (`10 * 2626 / 5252 = 5`).
