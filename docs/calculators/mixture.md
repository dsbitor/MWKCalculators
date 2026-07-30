# mixture

Mixture and dilution calculator.

**Converted from:** `MIXTURE.C`, `MWKC/WorkshopUtilities/mixture.zip`
**Go source:** `MWKGo/mixture/mixture.go`

## Purpose

Given the concentrations of two solutions (a dilutant such as pure
water is treated as a solution with 0% concentration) and any two of:
the amount of solution a, the amount of solution b, the total amount
of mixture, or the mixture's resulting concentration, this program
solves for the rest. Units for "amount" are unspecified and must be
used consistently.

## Inputs

| Prompt | Default |
|---|---|
| Ca = concentration of solution a | 10 % |
| Cb = concentration of solution b | 0 % |
| amount of solution a | (skip if unknown) |
| amount of solution b | (skip if unknown) |
| amount of mixture | (skip if unknown) |
| concentration of mixture | (skip if unknown) |

Exactly two of the four "amount"/"concentration of mixture" prompts
are required.

## Output

Amount and concentration of solution a, amount and concentration of
solution b, amount and concentration of the mixture.

## Method

```
pa, pb, pm = Ca/100, Cb/100, Cm/100  (fractions)

a+b known:        vm = va+vb;  pm = (va*pa+vb*pb)/vm
a+cm known:        vb = va*(pa-pm)/(pm-pb);  vm = va+vb
a+vm known:        vb = vm-va;  pm = (va*pa+vb*pb)/vm
b+cm known:        va = vb*(pm-pb)/(pa-pm);  vm = va+vb
b+vm known:        va = vm-vb;  pm = (va*pa+vb*pb)/vm
vm+cm known:       va = vm*(pm-pb)/(pa-pb);  vb = vm-va
```

## Worked Example

No worked numeric example was included with the original program. As
independently verifiable checks: a single self-consistent mixture (4
units of 10% solution plus 6 units of 0% solution giving 10 units of
4% mixture) is solved from each of the six recognized two-quantity
combinations, and each recovers the same full, consistent set; and,
regardless of which quantities were given, the mixture must conserve
total active-ingredient mass (an accounting identity independent of
the formula used to reach the answer); both confirmed in this
conversion's tests.
