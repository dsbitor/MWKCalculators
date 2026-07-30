# lvern

Linear vernier scale design.

**Converted from:** `LVERN.C` (M. W. Klotz, 11/98),
`MWKC/WorkshopUtilities/lvern.zip`
**Go source:** `MWKGo/lvern/lvern.go`

## Purpose

A vernier scale lets a sliding secondary scale read a measurement to
a finer precision than the primary scale's own major divisions allow.
Given the spacing of the main scale's major divisions, how many
subdivisions the main scale already has, and the desired vernier
resolution, this program computes how many divisions the vernier
scale needs, and, once a total vernier scale length is chosen, how
long each of those divisions is.

Inputs may be entered as a plain decimal (`1.5`) or as a fraction
(`1/2` or a mixed number, `1-1/2`); the program also reports the
nearest simple rational fraction for several of its own decimal
results, since machinists' scales are usually marked in fractions.

## Inputs

| Prompt | Default |
|---|---|
| Distance spanned by one major division on main scale | 1.0 |
| Number of subdivisions of one major division on main scale | 8 |
| Distance measured by one division on vernier scale | 1/32 |
| Distance spanned by vernier scale | (defaults to a computed "typical" span) |

## Output

The nearest rational fraction for the vernier resolution; the
distance spanned by one main-scale division (and its nearest
fraction); the number of vernier subdivisions required; a typical
vernier scale span; and, once a vernier span is chosen, the length of
one vernier division (and its nearest fraction).

## Method

```
mainDivisionLength = mainDivisionSpacing / mainSubdivisions
vernierSubdivisions = floor(mainDivisionLength / vernierResolution)
  (must come out to an exact whole number, or there is no solution)
typicalVernierSpan = mainDivisionSpacing - mainDivisionLength
vernierDivisionLength = vernierSpan / vernierSubdivisions
```

The nearest rational fraction for a decimal value uses a bounded
continued-fraction convergent search (`An = Xn*An-1 + An-2`, per the
original program's own algorithm), the same algebraic technique used
in `fraction`'s reduction logic but applied to approximate an
arbitrary decimal rather than reduce an already-exact ratio.

## Worked Example

`LVERN.TXT` includes two complete worked examples, both reproduced
exactly in this conversion's tests:

1. A one-inch major division split into eighths, with a vernier
   resolving thirty-seconds: main division length 0.1250 (1/8), 4
   vernier subdivisions, and, for the typical 0.8750 vernier span, a
   vernier division length of 0.2188 (nearest fraction 7/32).
2. A one-inch major division split into tenths, with a vernier
   resolving to 0.01: main division length 0.1000 (1/10), 10 vernier
   subdivisions, and a vernier division length of 0.0900 (9/100).
