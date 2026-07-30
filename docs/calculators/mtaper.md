# mtaper

Conical part taper measurement.

**Converted from:** `MTAPER.C`, `MWKC/WorkshopUtilities/mtaper.zip`.
Reference: Machinery's Handbook, 23rd Edition, pg. 685.
**Go source:** `MWKGo/mtaper/mtaper.go`

## Purpose

Measures the included angle of a tapered cylindrical part by setting
it on a V-block of known angle, itself set on a sine bar, and
adjusting the sine bar stack until an indicator swept across the part
shows no deviation. Given the V-block angle and sine bar length, this
program converts between the sine bar stack height and the part's
included angle in either direction, and, given a stack height,
decomposes it into blocks from a standard 81-piece gage block set.

## Inputs

| Prompt | Default |
|---|---|
| V-block included angle (D) | 90 deg |
| Sine bar length (L) | 5 in |
| 1 (angle from stack) or 2 (stack from angle) or Q | (menu, repeats) |
| Stack Height (mode 1) | 0.9947 in |
| Part included angle (mode 2) | 9.5 deg |

## Output

Sine bar angle, part included angle and half angle, and angle C
(between the part centerline and the sine bar surface); mode 2 also
prints the gage block decomposition of the resulting stack height.

## Method

```
mode 1 (angle from stack):
  B = asind(stackHeight / L)
  A = 2*atand(sind(B) / (cosd(B) + 1/sind(D/2)))
  C = asind(sind(A/2) / sind(D/2))

mode 2 (stack from angle):
  C = asind(sind(A/2) / sind(D/2))
  B = C + A/2
  stackHeight = L*sind(B)
```

The gage block decomposition uses a greedy digit-extraction
algorithm: the stack height's ten-thousandths digit picks a block in
the 0.1000-0.1009 range, then its thousandths and hundredths digits
together pick a block in the 0.101-0.150 range, then its remaining
tenths and hundredths pick a block in 0.05 increments up to 1.00,
and finally whatever whole-inch blocks (1, 2, 3, 4) remain. Digit
extraction is done via the same string-formatting approach as the
original program's own `kp` function (formatting the number to 9
decimal places and reading off a character), rather than pure
arithmetic, since floating-point decimal-digit extraction by modulo
is easy to get subtly wrong at exactly the precision this matters.

`WorkshopUtilities/sinebar.zip` solves the same "which gage blocks
sum to a target stack height" problem with a completely different
algorithm (a bounded combinatorial search, converted separately as
`sinebar`), rather than this greedy heuristic; both are kept as
distinct, faithful conversions of their respective originals rather
than merged into one "better" implementation.

## Worked Example

No worked numeric example was included with the original program. As
independently verifiable checks: converting a part angle to a stack
height and back exactly reproduces the original angle regardless of
which formula is used for either direction, and every gage block
decomposition's own chosen blocks sum back to the target stack height
within the algorithm's own tolerance; both confirmed in this
conversion's tests.
