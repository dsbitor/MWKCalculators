# unit

General-purpose units conversion.

**Converted from:** `UNIT.C` (M. W. Klotz), `MWKC/Misc/unit.zip`
**Go source:** `MWKGo/unit/unit.go`

## Purpose

Converts a value between any two compatible units, drawn from a table
of 246 named units, 23 metric-style prefixes, and 7 fundamental
dimensions (length, mass, time, angle, solid angle, charge, amount).
Supports compound expressions like `MILES/HOUR` or `KG/M^3`, checks
that a requested conversion's two sides describe the same physical
quantity before allowing it, breaks any unit down into its fundamental
dimensions, lists every unit sharing a given dimension (or a named
physical-quantity type, like "speed" or "pressure"), and lets you
define new units on the fly for the rest of a session.

This is the largest and most complex program in the whole Tier 2
batch: its own data file uses a distinct format from every other
program converted so far.

## Data setup

The unit and prefix tables are universal reference data (the same for
anyone), so they live in the shared `reference.db` SQLite database,
the same mechanism used by `fits`, `speed`, `gage`, `expand`,
`findthrd`, and `weight` earlier in this project; see
`ai/plans/c-to-go-conversion-plan.md`, "Data-file strategy for
Tier 2". No setup is required; `unit` runs immediately.

`UNIT.DAT` itself uses a section-based format unlike every other
program's `STARTOFDATA`/`ENDOFDATA` convention: `BEGINPREFIX`,
`BEGINPRIMARY`, and `BEGINMIXED` mark three different kinds of
section, each parsed differently:

- **Prefix section**: `NAME=value` per line (for example
  `KILO=1E3`).
- **Primary section**: units that map to exactly *one* fundamental
  dimension. A `NEWUNIT=d0,...,d6` directive sets which dimension
  applies to every unit line that follows, until the next `NEWUNIT` —
  individual unit lines give only a name and a factor (for example,
  after `NEWUNIT=1,0,0,0,0,0,0` (pure length), `MILE=1/1609.344` means
  a mile is `1/1609.344` meters).
- **Mixed (compound) section**: every field is explicit per line —
  `NAME=factor,d0,...,d6` (for example `NEWTON=1,1,1,-2,0,0,0,0`, a
  newton being `LENGTH¹·MASS¹·TIME⁻²`).

A factor can also be written `1/x` (matching `UNIT.C`'s own `fact()`
function, which only ever inverts the text after the slash — not a
general `a/b` fraction parser, though the shipped data never needs
one). The shipped data's own `NEWTONMETER` entry is missing its 7th
dimension field; this conversion treats a short `BEGINMIXED` line's
missing trailing dimensions as `0`, matching what the original's own
statically-initialized (and therefore zero-filled) array would have
left there.

## Inputs

A repeating menu of single-letter options:

| Option | Action |
|---|---|
| j (or space) | Convert a value between two units |
| l | Repeat the last conversion with a new "to" unit |
| a | List every unit of a chosen physical-quantity type (speed, force, pressure, ...) |
| b | Break a unit down into its fundamental dimensions |
| c | List every unit sharing a given unit's dimensions |
| d | Define a new unit for the rest of this session |
| f | Show the 7 fundamental units |
| i | List units starting with a given letter |
| p | List all prefixes |
| u | List all units |
| h | Show the menu |
| q | Quit |

A unit expression follows the syntax `(prefix)UNIT[^exponent]`,
optionally followed by `/ (prefix)UNIT[^exponent]` for a compound
unit's denominator — for example `KILOMETER`, `FT^2`, or `LB/FT^3`.
Plural forms (`MILES`, `INCHES`) are accepted and stripped down to
their singular form automatically, unless the plural form is itself
already a defined unit name (like `MILES`, which is defined directly
in `UNIT.DAT` precisely so this doesn't misfire and turn it into
`MIL`).

## Output

For a conversion: the entered value and unit, the converted value and
unit, and the value expressed in the fundamental (SI) unit system
with its dimensional breakdown — for example
`60 MILES/HOUR = 88 FEET/SEC = 26.8224 (METER) / (SEC)`.

## Method

`uanal` parses a unit expression: split on `/` into numerator and
denominator, extract an optional `^exponent` from each half, strip a
plural ending, then strip a leading prefix name (only attempted if the
bare string isn't already a recognized unit — so `FPS`, itself a
defined unit, is never mistakenly split into a nonexistent "F" prefix
plus "PS"), and resolve what remains to an exact (case-insensitive)
unit name. `uconv` then computes the actual numeric conversion by
routing the value through the fundamental unit system: from-units to
primary units, primary units to to-units, with each side's own prefix
and exponent applied along the way — matching `UNIT.C`'s own `uanal()`
and `uconv()` formula for formula, including a harmless-in-practice
quirk in the original's own denominator-exponent guard (`!= -1`
instead of `!= 1` on one comparison, which changes nothing since
raising a value to the power 1 is a no-op either way).

Two units can only be converted between if their net dimension vectors
(each unit's own 7 fundamental-dimension exponents, combined with its
expression's exponent and denominator) match exactly — attempting to
convert a mass to a force, for instance, reports `INCOMPATIBLE UNITS`
rather than silently producing a meaningless number.

`engprnt` formats a number the way `UNIT.C`'s own function does:
like C's `"%G"` (six significant digits, choosing plain decimal or
scientific notation automatically), except that when scientific
notation is used, the exponent is forced to the nearest multiple of
three below it — true "engineering notation," matching the
`kilo`/`mega`/`giga`-style prefix steps the rest of the program uses.

## Worked Example

`UNIT.TXT` ships a full worked tutorial with exact expected numeric
output for several conversions, which this conversion's tests
reproduce directly against the real shipped `reference.db`: 60
mph = 88 fps = 26.8224 m/s (checked both spelled out and using the
`MPH`/`FPS` abbreviations, and again with an implied value of 1); 10
lb/ft³ = 160.185 mg/cm³; 20 ft² = 2880 in²; 1 pound (mass) = 0.453592
kg; 1 poundforce = 4.44822 newton; and that converting pound (mass) to
newton (force) is correctly rejected as `INCOMPATIBLE UNITS`. The
tutorial's own worked "define a new unit" example
(`XX = FT^2/MILE^3`, expected to save as
`XX=4.48659e+010,-1,0,0,0,0,0,0`) is reproduced by directly checking
the accumulation logic `define` uses. `engprnt`'s forced-multiple-of-
three exponent behavior is checked both against a constructed value
and against the tutorial's own `222.886230E-012` example output.
`UNIT.DOC` (a 54KB general SI-units reference document bundled in the
same archive) contains no program-specific information and was not
otherwise used for this conversion.
