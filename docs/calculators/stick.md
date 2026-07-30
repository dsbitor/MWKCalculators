# stick

Thread-chasing carriage repositioning distance.

**Converted from:** `STICK.C` (M. W. Klotz, 5/01),
`MWKC/WorkshopUtilities/stick.zip`. Peter Lott's technique
(Machinist's Workshop, 6/01, pg. 18).
**Go source:** `MWKGo/stick/stick.go`
**Original documentation:** `STICK.TXT`, inside `MWKC/WorkshopUtilities/stick.zip` (not included in this conversion)

## Purpose

Cutting a metric thread on an Imperial leadscrew lathe (or vice
versa) usually means leaving the half-nuts engaged and running
the lathe in reverse between passes, since the thread dial can't
be used across unit systems. Peter Lott pointed out a faster
alternative: disengage the half-nuts and move the carriage a
distance that is simultaneously a whole number of leadscrew
pitches and a whole number of the thread pitch being cut, and
the half-nuts can be reengaged without losing registration. This
program finds the shortest such distance, covering all four
combinations of metric or Imperial leadscrew with metric or
Imperial thread being cut.

## Inputs

| Prompt | Default |
|---|---|
| Type of leadscrew thread | (M)etric or [I]mperial |
| Leadscrew pitch | 10 tpi (Imperial) or 4 mm (metric) |
| Type of thread | [M]etric or (I)mperial |
| Pitch of thread being cut | 15 tpi (Imperial) or 4 mm (metric) |

Blank input to either type prompt falls back to whichever branch
the original program's own default falls to: Imperial leadscrew,
metric thread.

## Output

Carriage repositioning distance, in inches and millimeters.

## Method

Both pitches are expressed as threads per inch (a metric pitch
in mm converts via `25.4/pitchMM`). Starting from one leadscrew
pitch, the distance is increased by one leadscrew pitch at a
time until the thread-pitch count at that distance is also
within a small tolerance of a whole number. The original program
bounds this search with a manual keypress bailout in case of a
runaway loop; this conversion replaces that with an explicit
maximum iteration count, returning an error rather than
searching indefinitely if it's ever exceeded.

## Worked Example

No fully worked numeric example was included with the original
program (`STICK.TXT` explains the technique and includes
correspondence about a possible future improvement, but no
sample run). As an independently verifiable check: a thread
pitch exactly double the leadscrew pitch converges on the first
candidate distance (one leadscrew pitch), confirmed in this
conversion's tests.
