# gearfind

Exhaustive gear-ratio search.

**Converted from:** `GEARFIND.C` (M. W. Klotz, 8/03),
`MWKC/WorkshopUtilities/gearfind.zip`
**Go source:** `MWKGo/gearfind/gearfind.go`

## Purpose

For a gear train where the available tooth counts aren't limited to a
fixed selection (any gear within some fabrication-limited range can
be cut), finds every combination of gear pairs (up to a chosen
maximum count) whose combined ratio matches a desired ratio within a
chosen tolerance, by brute-force search over every tooth count in a
chosen range. Deliberately inelegant but thorough — modern computers
make an exhaustive search practical where it once wasn't.

This archive also contains `GEARF1.C` and `GEARF2.C`, faster
successor programs specialized for exactly one and exactly two gear
pairs respectively (`GEARF2` using a continued-fraction technique
akin to `rfrac`). Both are deferred to a future group, since
`GEARFIND` alone already covers their cases (more slowly) as well as
the 3+ pair cases they don't handle at all.

## Inputs

| Prompt | Default |
|---|---|
| Desired ratio | 1.9945 |
| Allowable ratio error | 0.01 % |
| Maximum number of gear pairs (< 5) | 2 |
| Minimum number of gear teeth | 16 |
| Maximum number of gear teeth | 40 |

## Output

Every combination of gear pairs found within tolerance, each shown as
`teeth1:teeth2-teeth3:teeth4-...`, its resulting ratio, and its
percentage error; `**` marks each new best (closest) result found so
far.

## Method

An odometer-style nested search: for `maxPairs` gear pairs, `2 *
maxPairs` tooth-count values are enumerated across every combination
in `[minTeeth, maxTeeth]`, computing the combined ratio
`(teeth1/teeth2) * (teeth3/teeth4) * ...` for each and keeping every
combination within tolerance. The original program's interactive
keypress abort (there being no keyboard to poll in this conversion)
is replaced with an explicit evaluation-count bound per
`coding-style.md` Rule 2.

## Worked Example

No worked numeric example was included with the original program.
The documented default input was independently confirmed, via a
separate brute-force script, to have exactly 48 solutions with a best
error of about 0.000275%; this exact count is checked in this
conversion's tests, along with a check that every returned solution's
own reported ratio and error are self-consistent with its tooth
counts.
