# weight

Weight of standard shapes from volume and material density.

**Converted from:** `WEIGHT.C` (M. W. Klotz), `MWKC/WorkshopUtilities/weight.zip`
**Go source:** `MWKGo/weight/weight.go`

## Purpose

Computes how much a part weighs: `weight = density × volume`. Pick a
material (or enter your own density), pick one of sixteen standard
shapes and give its dimensions (or just type a volume directly if
you've already worked it out), and the calculator reports the volume
and weight. A running total of volume and weight is kept across a
session, so a complex object built from several simpler shapes (a tank
body plus its two end caps, for example) can be totaled up.

## Data setup

No setup is required; `weight` runs immediately using the 46 materials
shipped with the original program's own `WEIGHT.DAT`, listed below.
This is universal reference data (the same approximate densities for
anyone), so it lives in the shared `reference.db` SQLite database, the
same mechanism used by `fits`, `speed`, `gage`, `expand`, and
`findthrd` earlier in this project; see
`ai/plans/c-to-go-conversion-plan.md`, "Data-file strategy for
Tier 2".

## Materials

Densities are in lb/in³ (the program also shows the gm/cc equivalent
at runtime). If your material isn't listed, choose the "User input"
option at the bottom of the on-screen list to enter its density
directly.

| Material | Density (lb/in³) | Material | Density (lb/in³) |
|---|---|---|---|
| Aluminum | 0.0924 | Molybdenum | 0.370 |
| Asbestos | 0.021 | Monel | 0.319 |
| Beryllium | 0.067 | Nickel | 0.310 |
| Brass | 0.3032 | Oak | 0.029 |
| Brick | 0.065 | Palladium | 0.434 |
| Bronze | 0.3195 | Pine | 0.018 |
| Cast Iron | 0.26 | Platinum | 0.775 |
| Concrete | 0.076 | Rhenium | 0.760 |
| Copper | 0.3184 | Rhodium | 0.447 |
| Cork | 0.021 | Sand | 0.060 |
| Earth | 0.058 | Sea Water | 0.037 |
| Gasoline | 0.024 | Silver | 0.379 |
| Glass | 0.097 | Stainless | 0.290 |
| Gold | 0.697 | Steel | 0.2836 |
| Granite | 0.098 | Tantalum | 0.600 |
| Gunmetal | 0.3194 | Thorium | 0.420 |
| Ice | 0.033 | Titanium | 0.163 |
| Inconel | 0.304 | Tungsten | 0.700 |
| Indium | 0.264 | Uranium | 0.683 |
| Iridium | 0.813 | Water | 0.0361 |
| Lead | 0.4082 | White Metal | 0.2641 |
| Magnesium | 0.065 | Zinc | 0.2551 |
| Mercury | 0.4894 | Zirconium | 0.235 |

## Inputs

A material (by number from the list, or a custom density), a
dimension unit ([i]n, (f)t, (m)m, (c)m), then repeatedly: a shape
number, that shape's own dimensions, and how many of that shape to add
to the running total. Entering `-1` for the shape number returns to
material selection; `0` quits.

| Shape | Dimensions asked |
|---|---|
| 1 Rectangular solid | length, width, height |
| 2 Cylindrical solid | diameter, length |
| 3 Cylindrical pipe | outside diameter, inside diameter, length |
| 4 Hexagonal solid | size across flats, length |
| 5 Pyramid | length of base side, altitude |
| 6 Cone | diameter of base, altitude |
| 7 Frustum of a cone | larger diameter, smaller diameter, altitude |
| 8 Sphere | diameter |
| 9 Annulus (washer shape) | outer diameter, hole diameter, thickness |
| 10 Regular polygon cross section | number of sides, size (across flats if even, flat-to-opposite-vertex if odd), length |
| 11 Spherical cap | diameter of sphere, height of cap at center |
| 12 Two intersecting cylinders | diameter of cylinder(s) |
| 13 Torus | cross-section diameter, outside diameter |
| 14 Hemispherical shell | outside diameter, inside diameter |
| 15 Tangent ogive | diameter of base, length/height |
| 16 Conical wedge | cone base diameter, height of cone, sagitta of wedge base |
| 99 User Volume Input | volume directly (in³) |

## Output

For each shape: its volume, its weight (lb and kg), then the running
total volume and weight after adding the requested multiple of it.

## Method

Weight is `density × volume × unit conversion factor` (in³ per unit³:
1 for inches, 1728 for feet, 1/25.4³ for millimeters, 1/2.54³ for
centimeters), except shape 99 (volume entered directly in in³), which
needs no conversion. The rest of this section is about how that
volume figure itself is obtained for each of the sixteen shapes; the
weight step afterward is always the same one multiplication.

Each shape's volume is computed by its own closed-form formula, ported
directly from `WEIGHT.C`.

The regular polygon cross section (shape 10) measures its size two
different ways depending on whether it has an even or odd number of
sides — across flats for even, flat-to-opposite-vertex for odd —
matching `WEIGHT.C`'s own distinction (an even-sided polygon has flats
directly opposite each other to measure across; an odd-sided one
doesn't). The conical wedge (shape 16) and tangent ogive (shape 15)
are ported as closed-form functions directly from `WEIGHT.C`'s own
`wedge()` function and inline formula, respectively, rather than
redesigned, since both involve the same non-obvious trigonometric
derivations as the original.

## Worked Example

No worked numeric example with expected output is available. As
independently verifiable checks, this conversion's tests confirm
several shape
formulas against known geometric identities rather than hand-copied
numbers: a cylinder and sphere against their standard textbook
formulas; an annulus against an algebraically identical pipe; a
frustum with equal end diameters against a cylinder; a spherical cap
of height equal to the sphere's radius against half a sphere; a torus
against Pappus's theorem; a conical wedge at its degenerate cases
(zero sagitta cuts nothing, a sagitta equal to the base diameter
leaves the whole cone, equal sagitta and radius exactly halves it);
and a many-sided regular polygon prism converging on a cylinder of the
same diameter.

Separately, the actual weight computation (density × volume × unit
conversion factor) is tested end to end against the real shipped
`reference.db`, not just the volume formulas in isolation: Aluminum's
real density (0.0924 lb/in³) is loaded from the reference table and
multiplied by a 1 in diameter sphere's volume, and again by the same
physical sphere expressed in millimeters (25.4 mm diameter), checking
that both unit paths agree on the same weight.
