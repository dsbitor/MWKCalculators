# crod

Gudgeon (wrist) pin position for a slider-crank mechanism.

**Converted from:** `CROD.C` (M. W. Klotz), `MWKC/WorkshopUtilities/crod.zip`
**Go source:** `MWKGo/crod/crod.go`
**Original documentation:** `CROD.TXT`, inside `MWKC/WorkshopUtilities/crod.zip` (not included in this conversion)

## Purpose

Computes the position of a piston's gudgeon (wrist) pin across a
half-revolution of the crankshaft, both relative to the crank center
and relative to its position at top dead center (TDC) — the latter is
what's needed to convert an engine timing diagram (valve or port
opening/closing given in crank degrees) into an actual piston
position for blueprinting or building a two-stroke engine's ports.

## Inputs

| Prompt | Default |
|---|---|
| Connecting rod length | 2.4 |
| Crank radius (throw) | 0.6 |
| Angular increment, deg | 5.0 |

No units are enforced — use any consistent system. `CROD.TXT`
recommends a 5° interval for small engines, 1° for something engine
enough to power a bike or car.

## Output

A table of crank angle (0° at TDC, 180° at bottom dead center), pin
position relative to the crank center, and pin position relative to
TDC. If `-svg` is given, a diagram: both curves plotted against a
gridded crank-angle axis, with reference lines at the pin's TDC,
BDC, and (BDC minus one more throw) positions.

## Method

The classic slider-crank displacement formula:
`x(θ) = throw × cos(θ) + √(rodLen² − throw² × sin²(θ))`, the pin's
distance from the crank center. Its position relative to TDC is just
`throw + rodLen − x(θ)`, since `throw + rodLen` is the pin's own
maximum (TDC) extension.

The original's interactive mouse-driven menu and its click-anywhere-
on-the-curve coordinate readout are dropped entirely — there is no
equivalent for either in a static SVG image — rather than adapted;
see `ai/plans/c-to-go-conversion-plan.md`'s Tier 3 "Graphics scope"
resolution for the general policy this follows.

`CROD.TXT` itself refers to the original's output file as
`CROD.DAT`, but `CROD.C` actually writes `CROD.OUT` — a
documentation/code mismatch in the original; noted here rather than
silently perpetuated, though neither name is meaningful anymore now
that the table prints straight to stdout instead of a file.

## Worked Example

No numeric worked example was shipped with the original program. As
independently verifiable checks, this conversion's tests confirm two
well known physical facts about a slider-crank mechanism: at TDC
(θ=0°) the pin sits at its maximum extension, exactly
`throw + rodLen`; at BDC (θ=180°) it sits at its minimum extension,
exactly `rodLen − throw`. Further tests confirm the pin position
never leaves the `[rodLen−throw, rodLen+throw]` range at any angle,
and that the TDC-relative value is always exactly
`throw + rodLen − x(θ)`, the algebraic relationship the original
itself depends on rather than a redundant recomputation.
