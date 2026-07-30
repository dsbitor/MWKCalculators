# drill

Drill size lookup, tap drill, and step-drilling calculations.

**Converted from:** `DRILL.C` (M. W. Klotz), `MWKC/WorkshopUtilities/drill.zip`
**Go source:** `MWKGo/drill/drill.go`
**Original documentation:** `DRILL.TXT`, inside `MWKC/WorkshopUtilities/drill.zip` (not included in this conversion)

## Purpose

Looks up drills by hole size or by designation (fractional, number,
letter, or metric), computes the correct tap drill diameter for a
given tap (cutting or thread-forming, imperial or metric), and works
out a step-drilling schedule for opening up a pilot hole to a final
size in a chosen number of roughly-equal material-removal stages.

Although grouped with Tier 3 by the plan's own data-dependency note
(its shipped `.DAT` table also appears in the Tier 2 preparatory
scan), `drill` itself draws nothing at all — `DRILL.C` never calls a
single graphics primitive, only colored text output — so unlike its
seven group-mates, this conversion needed no SVG rendering work.

## Data setup

The 371-entry drill size table (every fractional, number, letter, and
common metric drill) is universal reference data (the same for
anyone), so it lives in the shared `reference.db` SQLite database, the
same mechanism used by `fits`, `speed`, `gage`, `expand`, `findthrd`,
`weight`, and `unit` earlier in this project; see
`ai/plans/c-to-go-conversion-plan.md`, "Data-file strategy for
Tier 2". No setup is required.

## Inputs

A repeating menu of single-letter options:

| Option | Action |
|---|---|
| D | Find the nearest drill given a hole size |
| S | Find a drill's size given its designation |
| T | Compute a tap drill diameter (cutting tap) |
| F | Compute a tap drill diameter (thread-forming tap) |
| X | Compute a step-drilling schedule |
| H, M | Show the menu |
| Q | Quit |

Sizes can be entered as a decimal (`1.5`) or a fraction, with an
optional whole-number part (`3/8`, `1-1/2`).

## Output

**D/S**: the matched drill plus three drills on either side (clipped
at the table's ends), each with its size difference from the target,
the match marked with `*****`.

**T/F**: the tap's diameter and pitch (both units), the resulting tap
drill diameter, and the same neighbor report as D/S — with each
candidate drill's own percentage depth of thread shown alongside it
for a cutting tap (T), since a driller might reasonably pick a
neighboring size to trade a little thread depth for a size actually
in their drawer.

**X**: what fraction of the total material the pilot hole itself
already accounts for, then each stepped-up drill size in turn with
its own cumulative percentage removed, ending at the requested final
size.

## Method

`findDrill` clamps to either end of the table for a size beyond its
range, otherwise picks the closest match by absolute difference; an
"exclude metric" mode (used by step drilling when the user declines
metric drills) walks back past metric-only entries to the nearest one
that isn't — a dual-labeled entry like `13=4.70 mm` still counts as
acceptable, since it also has a non-metric designation.

The tap drill formula for a cutting tap is
`diameter − 0.01 × depth% × K / pitch`, where `K` is one of two
published thread-form constants: `(6/8)×tan(60°) ≈ 1.299` for the
American National (imperial, the default) or Standard (metric) thread
form, or `(5/8)×tan(60°) ≈ 1.082` for the American Unified (imperial)
or ISO (metric) thread form — see `DRILL.TXT`'s own lengthy discussion
of why both exist and which one matches published tap-drill charts
(the default, sharp-crested-thread assumption does). A thread-forming
(roll) tap instead uses `diameter − 0.0068 × depth% / pitch`
(imperial) or `(tapDiameterMM − depth% × pitchMM / 147.06) / 25.4`
(metric) — a different, non-thread-form-constant-based formula, since
a forming tap doesn't cut material away the same way.

Step drilling repeatedly increases the removed circular area by a
fixed percentage of the total area to remove, converts back to a
diameter, and snaps to the nearest actually-available drill after
*each* step — so later steps compound from the previous step's
snapped-to-a-real-drill size, not an idealized geometric progression,
matching `DRILL.C`'s own `step()` exactly.

## Worked Example

`DRILL.TXT` is entirely about the two thread-form constants and gives
their values as "1.299..." and "1.082...", both reproduced exactly by
this conversion's `sixEighthsTan60`/`fiveEighthsTan60`. As a further,
independently verifiable check, this conversion's tests confirm a
1/4-20 UNC tap (0.25 in diameter, 20 threads/in, National thread form,
75% depth of thread) computes a tap drill diameter that resolves to
the `7` drill (0.201 in) — the standard tap-drill recommendation
printed on nearly every tap drill chart for this extremely common
thread size. Further tests confirm drill lookup boundary clamping,
the exact-designation lookup's case-insensitivity, the fraction
parser, and that a step-drilling schedule always ends at the
requested final drill with exactly 100% of the material removed.
