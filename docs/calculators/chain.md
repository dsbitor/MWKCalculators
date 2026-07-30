# chain

Sprocket center-to-center distance for a given chain length.

**Converted from:** `CHAIN.C` (M. W. Klotz, 2/03),
`MWKC/WorkshopUtilities/sprocket.zip`
**Go source:** `MWKGo/chain/chain.go`

## Purpose

Given a roller chain's pitch and length (in pitches) and the
tooth counts of the two sprockets it will run on, computes the
center-to-center mounting distance needed. The companion program
`SPROCKET.C`, from the same zip file, has not been converted
yet.

No `.TXT` file was included with the original program; this
purpose statement is drawn from the `.C` file's own header
comment.

## Inputs

| Prompt | Default |
|---|---|
| Chain pitch | 1 in |
| Chain length | 48 pitches |
| Number of teeth in large sprocket | 18 |
| Number of teeth in small sprocket | 9 |

## Output

Sprocket center-to-center distance.

## Method

This is the standard roller chain center-distance formula:

```
toothDifference = largeTeeth - smallTeeth
toothSum = 2*chainLength - largeTeeth - smallTeeth
centerDistance = (pitch/8) * (toothSum + sqrt(toothSum^2 - 0.810*toothDifference^2))
```

The constant 0.810 approximates `8/pi^2`, the standard
coefficient in this formula.

## Worked Example

No worked example was included with the original program. As an
independently verifiable check: equal sprocket sizes make the
tooth-difference term vanish, reducing the formula to
`(pitch/2)*(chainLength-teeth)`, a simpler identity independent
of this code.
