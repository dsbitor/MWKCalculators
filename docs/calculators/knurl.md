# knurl

Workpiece diameter for a perfectly closing knurl pattern.

**Converted from:** `KNURL.C` (M. W. Klotz, 11/98),
`MWKC/WorkshopUtilities/knurl.zip`
**Go source:** `MWKGo/knurl/knurl.go`

## Purpose

Knurling a workpiece whose circumference isn't an exact whole
multiple of the knurl wheel's tooth spacing leaves a visible
seam where the pattern fails to close up. Given the knurl
wheel's diameter and tooth count, and the workpiece's nominal
diameter, this program finds the closest diameter that fits a
whole number of teeth and closes the pattern perfectly.

No `.TXT` file was included with the original program; this
purpose statement is drawn from the `.C` file's own header
comment.

## Inputs

| Prompt | Default |
|---|---|
| Diameter of knurl wheel | 0.625 in |
| Number of teeth on knurl wheel | 40 |
| Nominal diameter of workpiece | 0.87 in |

## Output

Spacing between knurl teeth, the whole number of crests that
fit, the required workpiece circumference, and the workpiece
diameter that produces it.

## Method

```
tooth spacing = pi * wheel diameter / tooth count
crest count   = floor(pi * nominal diameter / tooth spacing)
circumference = crest count * tooth spacing
diameter      = circumference / pi
```

## Worked Example

No worked example was included with the original program. With
the documented default inputs, the program reports a tooth
spacing of 0.049in, 55 crests, a required circumference of
2.700in, and a workpiece diameter of 0.859in.
