# expand

Material thermal expansion calculations.

**Converted from:** `EXPAND.C` (M. W. Klotz), `MWKC/WorkshopUtilities/expand.zip`
**Go source:** `MWKGo/expand/expand.go`

## Purpose

Given a material's coefficient of thermal expansion, a nominal
length, and either a temperature change or a length change, this
calculator computes the other one — how much a part grows (or
shrinks) over a given temperature swing, or how much temperature
change would produce a given length change.

## Inputs

| Prompt | Default |
|---|---|
| Material number (from the displayed list), or a number past the end for a custom value | one past the end of the list |
| Material coefficient of thermal expansion (ppm/degF) — only if a custom value was chosen | 10.0 |
| Nominal length of object (in) | 1.0 |
| `([A],B) ?` — find length change given temperature change, or the reverse | A |
| Temperature change (degF) — option A | 100.0 |
| Length change (in) — option B | 0.0001 |

## Output

The material's CTE, the nominal length, the length change, and the
temperature change (whichever of the last two was computed rather
than entered).

## Method

```
length change = CTE * nominal length * temperature change
```

The thirty materials are universal reference data, so — like `fits`,
`speed`, and `gage` in the same conversion group — they live in the
shared `reference.db` SQLite database rather than a flat file read at
startup; see `ai/plans/c-to-go-conversion-plan.md`, "Data-file
strategy for Tier 2", and `fits.md` for the shared mechanism.
Materials are listed, and selected by number, alphabetically
(case-insensitively), matching `EXPAND.C`'s own sort before display;
one past the end of the list is offered as a "User input" option that
prompts for a coefficient directly, exactly as the original does,
reported back under the placeholder name `??` the original also uses
for it. Published values for a material's coefficient of thermal
expansion vary considerably between sources; the shipped values are
representative averages across several references rather than a
single authoritative source, and a user with a precise value for
their own material should enter it via the custom-value option
instead.

## Worked Example

No worked numeric example is available. As an
independently verifiable check, this conversion's tests confirm that
computing a length change and then recovering the temperature change
from it reproduces the original temperature change, for several
materials and lengths — the defining property of the formula's own
invertibility, not a re-derivation of it. An integration test against
the real shipped `reference.db` confirms aluminum's coefficient
(12.44 ppm/degF, a widely cited reference value) is present and that
materials come back sorted case-insensitively.
