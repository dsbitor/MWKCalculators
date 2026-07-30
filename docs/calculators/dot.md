# dot

Depth of thread calculations.

**Converted from:** `DOT.C` (M. W. Klotz, 2/99, 12/02, 3/04, 6/05),
`MWKC/WorkshopUtilities/dot.zip`
**Go source:** `MWKGo/dot/dot.go`

## Purpose

Single point threading on a lathe needs the depth of cut, which
depends on how the crest and root of the thread are actually formed:
sharp-to-sharp is the theoretical maximum, but crest and root are
usually truncated with a flat, giving four standard combinations.
Given the thread angle, threads per inch, and compound rest angle,
this program reports all four depths (both perpendicular to the
thread axis and along the compound feed), the doubled
sharp-crest-to-sharp-root figure some operators infeed on the cross
slide directly, an American National (60 degree) pitch-diameter
correction figure, and advice on which threading dial lines apply.

## Inputs

| Prompt | Default |
|---|---|
| Threads angle | 60 deg |
| Threads per inch | 20 |
| Compound rest angle | 29 deg |

## Output

Four depth-of-thread figures (A: sharp crest/sharp root, B: flat
crest/flat root, C: sharp crest/flat root, D: flat crest/sharp
root), each with its compound-feed equivalent; the doubled A figure
(E); a major-diameter correction for American National threads; and
a threading dial hint.

## Method

```
H = 0.5*pitch / tand(0.5*threadAngle)     (A: sharp crest - sharp root)
B = 0.625*H
C = 0.75*H
D = 0.875*H
E = 2*H

compound-feed equivalent = figure / cosd(compoundAngle)
```

The threading dial hint depends on threads per inch: any dial line
works for an even whole tpi, only numbered lines for an odd whole
tpi, and only odd-numbered lines whenever tpi has a fractional part.

## Worked Example

No worked numeric example was included with the original program. As
an independently verifiable check: for the standard 60 degree thread
form, H reduces to the well known ANSI B1.1 constant `0.866025 *
pitch`, independent of this program's own tangent-based derivation
of the same value; confirmed in this conversion's tests, along with
the documented 5/8, 3/4, and 7/8 fractions for B, C, and D.
