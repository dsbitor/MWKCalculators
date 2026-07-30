# Calculator Reference

One page per converted calculator, under `docs/calculators/`.
Each page carries forward the purpose, inputs, method, and any
worked example from the original C program's `.TXT` file or,
where none was included, its `.C` header comment, so that
information isn't lost once the original `.TXT` files stop
being read directly.

This index is updated at the end of each conversion session as
new calculators are added; it is not a substitute for reading
the plan in `ai/plans/c-to-go-conversion-plan.md` for tier and
scheduling status.

A small number of the original programs depend on a capability
that existed only because DOS was single-tasking: direct access
to the screen as a block of shared memory any program could
read at any time. That capability has no equivalent on a modern
operating system, where each terminal owns its own display
privately, so these programs are excluded from conversion
rather than reimplemented as something functionally different
wearing the original name. Modern operating systems and
terminal tools already cover the same underlying needs (saving
or reviewing what has appeared on screen) in their own way. See
"Excluded" under each tier below for the specific programs this
applies to.

## Tier 0

| Calculator | Description |
|---|---|
| [chord](calculators/chord.md) | Chord length for stepping off a circle into equal divisions |
| [circ3](calculators/circ3.md) | Radius of a circle through three points, given pairwise distances |
| [factor](calculators/factor.md) | Prime factorization and complete divisor list |
| [doomsday](calculators/doomsday.md) | Day of the week for any date, by two independent methods |
| [catenary](calculators/catenary.md) | Droop and length of a hanging cable |

## Tier 1

### Group 1

| Calculator | Description |
|---|---|
| [sinenl](calculators/sinenl.md) | Sine bar made from two touching cylinders, no link |
| [eccentub](calculators/eccentub.md) | Tube size for turning offset eccentrics in a 3-jaw chuck |
| [bend](calculators/bend.md) | Sheet metal bend allowance |
| [dallow](calculators/dallow.md) | Drill tip allowance |
| [cusp](calculators/cusp.md) | Pass spacing for a ball end mill to limit cusp height |
| [knurl](calculators/knurl.md) | Workpiece diameter for a perfectly closing knurl pattern |

### Group 2

| Calculator | Description |
|---|---|
| [mradius](calculators/mradius.md) | Radius of curvature of a part's edge, measured with two rollers |
| [dplate](calculators/dplate.md) | Disk diameter for a dividing plate, no dividing head needed |
| [ftaper](calculators/ftaper.md) | Female taper angle, measured with two spheres |
| [revolver](calculators/revolver.md) | Dimensions for a revolver-cylinder style tool holder |
| [sine](calculators/sine.md) | Sine bar made from two cylinders connected by drilled links |
| [eccent](calculators/eccent.md) | Packing thickness for turning an eccentric in a 3-jaw chuck |
| [fraction](calculators/fraction.md) | Rational fraction calculator |

### Group 3

| Calculator | Description |
|---|---|
| [links](calculators/links.md) | Taper angles and shim height for a tapered, radiused-end link |
| [osborne](calculators/osborne.md) | Convergence demonstration for the "Osborne Maneuver" of centering round stock |
| [wire](calculators/wire.md) | AWG wire gage recommendation and properties |
| [tubewall](calculators/tubewall.md) | Tube wall thickness from an outside micrometer reading |
| [chill](calculators/chill.md) | Wind chill equivalent temperature |
| [dew](calculators/dew.md) | Dew point temperature |
| [ungula](calculators/ungula.md) | Volume of an ungula (water in a tilted cylindrical glass) |

### Group 4

| Calculator | Description |
|---|---|
| [chain](calculators/chain.md) | Sprocket center-to-center distance for a given chain length |
| [plate](calculators/plate.md) | Chain-drilling layout for cutting a circular plate free from flat stock |
| [slug](calculators/slug.md) | Chain-drilling layout for opening a large hole in plate stock |
| [collet](calculators/collet.md) | Bore diameter for a slotted collet holding polygonal stock |
| [mandrel](calculators/mandrel.md) | Recommended mandrel diameter for spring winding |
| [tenon](calculators/tenon.md) | Depth of cut for a regular polygonal tenon |
| [horse](calculators/horse.md) | Torque, rotational speed, and horsepower solver |

### Group 5

| Calculator | Description |
|---|---|
| [soddy](calculators/soddy.md) | Outer and inner Soddy circle diameters for three mutually tangent circles |
| [compound](calculators/compound.md) | Lathe compound rest angle for a given infeed ratio |
| [cone](calculators/cone.md) | Flat sheet-metal pattern for a conical frustum |
| [cmiter](calculators/cmiter.md) | Compound table saw angles for mitred polygonal forms |
| [helixcf](calculators/helixcf.md) | Helical gear cutting dimensions |
| [bucket](calculators/bucket.md) | Slant height divisions for a bucket-shaped (frustum) container |
| [stick](calculators/stick.md) | Thread-chasing carriage repositioning distance |

### Group 6

| Calculator | Description |
|---|---|
| [rattle](calculators/rattle.md) | Diameter measurement of a large bore using Lautard's stick-and-rattle technique |
| [boltcirc](calculators/boltcirc.md) | Bolt circle hole layout |
| [protrac](calculators/protrac.md) | Sinebar-like protractor |
| [conrod](calculators/conrod.md) | Connecting rod to cylinder wall clearance |
| [dovetail](calculators/dovetail.md) | Dovetail measurement using pins |
| [rotary](calculators/rotary.md) | Rotary table division angles |
| [ellipse](calculators/ellipse.md) | Ellipse eccentricity, area, and perimeter |

### Group 7

| Calculator | Description |
|---|---|
| [wedgeco](calculators/wedgeco.md) | Volume of a conical wedge |
| [sprocket](calculators/sprocket.md) | ANSI standard roller chain sprocket dimensions |
| [loft](calculators/loft.md) | Minimum thread engagement length |
| [loan](calculators/loan.md) | Loan amortization schedule |
| [flute](calculators/flute.md) | Tapered flute milling setup |
| [rounder](calculators/rounder.md) | Ball end mill rounding-over table |
| [dot](calculators/dot.md) | Depth of thread calculations |

### Group 8

| Calculator | Description |
|---|---|
| [psych](calculators/psych.md) | Relative humidity from wet/dry bulb psychrometer |
| [polygon](calculators/polygon.md) | Properties of regular polygons |
| [feed](calculators/feed.md) | Milling feed rate solver |
| [latlon](calculators/latlon.md) | Great circle distance and bearing between two points |
| [lvern](calculators/lvern.md) | Linear vernier scale design |
| [cep](calculators/cep.md) | Circular Error Probable (CEP) |
| [gear](calculators/gear.md) | Spur gear dimensions |

### Group 9

| Calculator | Description |
|---|---|
| [atmos76](calculators/atmos76.md) | 1976 U.S. Standard Atmosphere model |
| [gearpa](calculators/gearpa.md) | Gear pressure angle estimation by chordal span |
| [ballcut](calculators/ballcut.md) | Incremental sphere turning on a lathe |
| [polycone](calculators/polycone.md) | Geometry of a "polycone" (pyramid-like) solid |
| [ugroove](calculators/ugroove.md) | Incremental cutting of a U-shaped groove |
| [rfrac](calculators/rfrac.md) | Rational fraction approximation |
| [temp](calculators/temp.md) | Temperature scale converter |

### Group 10

| Calculator | Description |
|---|---|
| [3wire](calculators/3wire.md) | Three-wire thread measurement |
| [gearspur](calculators/gearspur.md) | Spur gear dimension solver |
| [diophan](calculators/diophan.md) | Linear Diophantine equation solver |
| [mixture](calculators/mixture.md) | Mixture and dilution calculator |
| [gas](calculators/gas.md) | Perfect gas law (PV = nRT) calculator |
| [buoy](calculators/buoy.md) | Buoy immersion depth calculator |

### Group 11

| Calculator | Description |
|---|---|
| [gearfind](calculators/gearfind.md) | Exhaustive gear-ratio search |
| [offkey](calculators/offkey.md) | Offset key calculations |
| [cseg](calculators/cseg.md) | Circular segment calculations |
| [triangle](calculators/triangle.md) | Solution of plane triangles |

### Group 12

| Calculator | Description |
|---|---|
| [234](calculators/234.md) | Quadratic, cubic, and quartic equation solver |
| [dipstick](calculators/dipstick.md) | Tank volume fraction from dipstick reading |

### Group 13

| Calculator | Description |
|---|---|
| [sphere](calculators/sphere.md) | Solution of spherical triangles |
| [flywheel](calculators/flywheel.md) | Tapered spoke flywheel calculations |

### Group 14

| Calculator | Description |
|---|---|
| [mtaper](calculators/mtaper.md) | Conical part taper measurement |
| [sinebar](calculators/sinebar.md) | Sine bar calculations |

### Group 15

| Calculator | Description |
|---|---|
| [vernier](calculators/vernier.md) | Two-plate angular vernier design |
| [scalc](calculators/scalc.md) | A collection of specialized calculators |

### Group 16

| Calculator | Description |
|---|---|
| [mix](calculators/mix.md) | Mixed dimensional units four-function calculator |
| [rpc](calculators/rpc.md) | RPN (Reverse Polish Notation) stack-oriented scientific calculator |

This group closes out Tier 1 review of the last four programs
originally listed as Tier 1 candidates (`mix`, `rpc`, `sun`,
`unit`). Both `sun` and `unit` turned out to be misclassified: a
scan for `fopen` during the original Tier 1 suitability review
silently skipped both `SUN.C` and `UNIT.C`, because both files
contain extended-ASCII degree-symbol bytes that make `grep` treat
them as binary and suppress its output unless `-a` is passed
explicitly. `UNIT.C` `fopen`s `UNIT.DAT` to load its entire
unit-definition database at startup (the same kind of miss already
corrected once for `colsort`); `SUN.C` `fopen`s `SUN.DAT` to load
default latitude/longitude/timezone-offset data. Both move to Tier 2
rather than being converted here.

### Clean-up

| Calculator | Description |
|---|---|
| [combi](calculators/combi.md) | Combination (n choose m) enumeration |
| [ratio](calculators/ratio.md) | Preferred-numbers (Renard) series generator |

`combi` and `ratio` were misclassified as Tier 2 by the original prep
scan (neither reads a data file; both only `fopen` an *output* file)
and are converted here as ordinary Tier 1 programs instead, alongside
the exclusion of `data.zip` (see "Excluded" below) — together, the
Tier 2 clean-up phase this project's plan called for before starting
Tier 3.

## Tier 2

Programs that read a shipped `.DAT` reference or configuration table
at startup. See `ai/plans/c-to-go-conversion-plan.md`, "Data-file
strategy for Tier 2", for the reference-data/user-data split and the
shared tooling (`internal/legacydat`, `internal/refdata`,
`MWKGo/tools/build-refdb`) both groups below rely on.

### Group 1 (trial)

| Calculator | Description |
|---|---|
| [fits](calculators/fits.md) | Shaft/hole fit computations |
| [speed](calculators/speed.md) | Machining speed (recommended spindle RPM) utility |

A two-program trial run before committing to a full group size for
the remaining ~19 Tier 2 programs, picked as the two shortest,
structurally representative candidates (both a straightforward
"read a small reference table, look something up by number, apply a
formula" shape). Both are universal reference data — the same
allowance table or cutting-speed chart is correct for every user — so
both are the first real consumers of the shared `reference.db`
SQLite database that `MWKGo/tools/build-refdb` builds offline from
the original `FITS.DAT`/`SPEED.DAT` and every calculator embeds via
`internal/refdata`. `fits` and `speed` both number their reference
rows by the original `.DAT` file's own order rather than
alphabetically, since each program's documented default selection
(a "push" fit; aluminum) depends on that order; a `list_position`
column on each table preserves it. `speed.DAT` also exercises the one
real structural variant found across the batch: it has no
`STARTOFDATA` marker, which `internal/legacydat.Rows` already
supported (an empty start marker means the data region starts
immediately).

### Group 2

| Calculator | Description |
|---|---|
| [diam](calculators/diam.md) | Machining diameter utility |
| [diffthrd](calculators/diffthrd.md) | Differential thread calculations |
| [gage](calculators/gage.md) | Wire and sheet gage utility |
| [divhead](calculators/divhead.md) | Dividing head calculations |
| [expand](calculators/expand.md) | Material thermal expansion calculations |

The first full-size group (5 programs, following Tier 1's
ascending-line-count convention), and the first to actually reach the
machine-specific-data half of the Tier 2 split: `gage` and `expand`
are universal reference data, added to the same shared `reference.db`
`fits`/`speed` created; `diffthrd` and `divhead` turned out to be
entirely machine-specific once their `.DAT` files were read closely
(`diffthrd` was tentatively guessed "universal" back when the Tier 2
list was first drafted, a guess this group's closer reading
corrects), reading instead from a new shared, empty-by-default
`userdata.db` (`internal/userdata`) populated only by CSV import,
never a fabricated default; and `diam` needs both at once — it
reuses `speed`'s own `machining_speeds` table for its material data
(the original `DIAM.DAT` literally repeats `SPEED.DAT`'s list
verbatim) while reading its machine's own available spindle speeds
from `userdata.db`, the first calculator in the batch to draw on both
databases together. Every machine-specific-data program in this group
ships a worked example CSV, built from its own original `.DAT` file's
default configuration, under its own `testdata/` directory, and
supports a `-import`-style flag that loads a user's CSV via
`internal/csvtable.Import` and exits.

### Group 3

| Calculator | Description |
|---|---|
| [calibrat](calculators/calibrat.md) | Calibrate a linear scale |
| [vrev](calculators/vrev.md) | Calculations for various sized holes arranged on a bolt circle |
| [findthrd](calculators/findthrd.md) | Find a thread from its diameter or pitch |
| [simul](calculators/simul.md) | Simultaneous linear equation solver |
| [spaceblk](calculators/spaceblk.md) | Space block selection utility |

This group introduces a third data category alongside universal
reference data and reusable machine configuration: per-invocation
job data. `calibrat` (one instrument's calibration checkpoints),
`vrev` (one job's list of hole diameters), and `simul` (one system of
equations) are all specific to a single run, not reusable across
many — so instead of either SQLite database, all three read a file
named on the command line fresh each run, in the same
`STARTOFDATA`/`ENDOFDATA` text format the originals used, via
`internal/legacydat` at runtime rather than only as an offline
build-time tool. `findthrd` (a third-party-compiled table of ~470
standard threads) is universal reference data, added to
`reference.db`; its `.DAT` file is also the one file in the batch
using fixed-width columns instead of a delimiter, needing its own
column-offset parsing in `MWKGo/tools/build-refdb` rather than
`internal/legacydat.Fields`. `spaceblk` (an owner's own space block
set) is machine-specific, added to `userdata.db`, and reuses the
odometer-style nested search `gearfind` pioneered in Tier 1, bounded
per `coding-style.md` Rule 2 in place of the original's interactive
abort. `calibrat`, `vrev`, and `simul` each ship a worked example data
file built from their own original `.DAT`'s default content or
`.TXT`'s documented example, letting a user try the program
immediately without writing their own input file first.

### Group 4

| Calculator | Description |
|---|---|
| [gearatio](calculators/gearatio.md) | Find a chain of gears matching a required ratio |
| [colsort](calculators/colsort.md) | Multi-column table sorting |
| [change](calculators/change.md) | Lathe change gear calculations |
| [weight](calculators/weight.md) | Weight of standard shapes from volume and material density |
| [ddh](calculators/ddh.md) | Differential dividing head calculations |

`gearatio` is a stand-alone gear-ratio finder directly adapted from
`change`'s own gear-chain search; both, along with `ddh`'s
differential-indexing fallback, share the same search (no gear
position reused within one chain, each stage's ratio applied by
multiplying or dividing toward the target, every match within
tolerance reported rather than only the first) reimplemented in each
program rather than factored into a shared package, consistent with
this project's one-file-per-calculator convention. All three, plus
`change`'s leadscrew pitch(es) and `ddh`'s hole plates and change
gears, are machine-specific data in `userdata.db`. `weight` (a
material density table) is universal reference data in `reference.db`,
alongside a from-scratch port of sixteen closed-form shape-volume
formulas, verified in this conversion's tests against known geometric
identities (Pappus's theorem for a torus, a many-sided polygon
converging on a cylinder, and several degenerate-case checks) rather
than hand-copied numbers, since neither original program shipped a
worked numeric example. `colsort` — previously deferred out of Tier 1
into this project's Tier 2 backlog, see "Excluded" below — turned out
to be a third example of this project's per-invocation "ephemeral job
data" category introduced in Group 3: both the table to sort and the
sort job's own parameters (column count, column types, sort column,
direction) live in one data file read fresh each run.

### Group 5

| Calculator | Description |
|---|---|
| [curfit](calculators/curfit.md) | Curve fitting to a set of data pairs |
| [sun](calculators/sun.md) | Solar position, equation of time, sunrise/sunset, and sundial geometry |
| [unit](calculators/unit.md) | General-purpose units conversion |
| [belt](calculators/belt.md) | Flat belt length for an arbitrary arrangement of pulleys |
| [qbelt](calculators/qbelt.md) | Quick two-pulley belt length calculation |
| [pulley](calculators/pulley.md) | Solve for an unknown pulley diameter given belt length and separation |
| [pcd](calculators/pcd.md) | Solve for pulley center distance given both diameters and belt length |

The last group before the Tier 2 clean-up phase. `curfit` (per-
invocation data pairs, like `calibrat`/`vrev`/`simul`/`colsort`)
reuses the same full-pivoting Gauss-Jordan solver `simul` already
built for its polynomial fit. `sun` (an observer's own latitude,
longitude, and UT offsets, reused across many runs from one location)
is machine-specific data in `userdata.db`, alongside a from-scratch
port of the 1978 "Almanac for Computers" long-term solar-position
series; its tests check the summer and winter solstice declinations
against the sun's well known ±23.44° obliquity rather than a
hand-copied worked example, since none was shipped. `unit` (246
units, 23 prefixes, universal reference data in `reference.db`) is
the largest and most complex program in the whole batch, using a
section-based data-file format (`BEGINPREFIX`/`BEGINPRIMARY`/
`BEGINMIXED`) unlike every other Tier 2 program's
`STARTOFDATA`/`ENDOFDATA` convention; its own `.TXT` file shipped a
full worked tutorial with exact numeric output, which this
conversion's tests reproduce directly against the real shipped
`reference.db`. `belt` (a job's own pulley layout, per-invocation
ephemeral data like `curfit`) closes out the regular Tier 2 sequence,
bringing its three bundled companion programs along in the same
group as the plan's own clean-up note anticipated: `qbelt` (a quick,
data-file-free two-pulley calculation), `pulley` (solve for an
unknown pulley diameter given belt length and separation), and `pcd`
(solve for pulley separation given both diameters and belt length,
added five years after the other three). All three companions share
one closed-form two-pulley formula, duplicated across each rather
than factored into a shared package per this project's convention;
`belt`'s own general N-pulley algorithm is cross-checked against that
same closed-form formula for the plain two-pulley case, and its
slack-pulley detection is verified against `BELT.DAT`'s own shipped
example, whose comment names exactly which pulley should come out
slack. `pulley` and `pcd` are verified against `BELT.DAT`'s own
"conical pulley" example rows, a family of pulley-diameter pairs the
original data file documents as all sized for the same fixed belt
length and separation.

This closes out every regularly-numbered Tier 2 group and the
clean-up phase along with it: `combi` and `ratio` joined Tier 1
proper (see above), and `data.zip` — a folder of 16 loosely-formatted
reference `.TXT` tables plus one image, with no source code and no
consistent structure across files — was excluded rather than
converted; see "Excluded" below.

### Excluded

`photo` (`MWKC/Misc/photo.zip`) is not a C program: it is x86
assembly (`PHOTO.ASM`, with a compiled `PHOTO.COM`). It reads
DOS text-mode video memory directly to capture whatever is
currently displayed on screen and appends it, as plain text, to
a growing log file. Besides being outside this project's
C-to-Go scope on language grounds alone, the capability it
depends on, reading the screen as shared memory from outside
whatever program put it there, is specific to DOS's
single-tasking model and has no counterpart on macOS or Linux;
a Go program calling itself `photo` would necessarily do
something else entirely (shell out to a terminal multiplexer's
own pane-capture feature, for instance) rather than the same
thing in a new language. It will not be converted.

`hints` (`MWKC/WorkshopUtilities/hints.zip`) contains no C
program at all: it is a plain text file of shop tips and mini
projects (`HINTS.TXT`), with no accompanying `.C` source. There
is nothing to convert.

`data` (`MWKC/WorkshopUtilities/data.zip`) likewise contains no C
program: it is a folder of 16 reference `.TXT` tables (material
density, wood screw sizes, pipe/thread/gauge charts, tapers, a
metric-conversion sheet, and similar shop reference sheets) plus one
JPEG, with no source code and no shared format across files — one
file (`FILES.TXT`) is not even a manifest, despite its name, but just
another unrelated table. There is nothing to convert into a
calculator.

`data.txt` and `hints.txt`, from the original Marv Klotz archive
these came from, are planned to be made available as formatted PDF
files in a future release, rather than ported as software.

`colsort` (`MWKC/Misc/colsort.zip`) reads an external data file
(`COLSORT.DAT`) via `fopen` at runtime; the Tier 1 suitability
review's `fopen` scan missed it, the same kind of correction made
for `ratio`/`ratio 2` in Tier 1 group 1. It moved to Tier 2
(data-file-backed) rather than staying in Tier 1, and has since been
converted; see Group 4 above.

## Tier 3

The eight graphics-bearing calculators. Per
`ai/plans/c-to-go-conversion-plan.md`'s "Graphics scope for the eight
Tier 3 programs," the original DOS graphics calls (line, box, and
circle primitives, called through each program's own `wline`/`wbox`/
`wcircle`-style wrapper functions — the bodies of which exist nowhere
in the source tree; their behavior is inferred from naming and the
one program, `spline`, that calls the underlying `_moveto_w`/
`_lineto_w` primitives directly) are replaced by a small internal
package, `internal/svgplot`, emitting SVG instead of drawing to a
screen. Every original interactive feature building on those DOS
graphics calls — a "press a key to continue" pause, or (for the
programs that had one) a full mouse-driven click-to-read-coordinates
menu — is dropped entirely, not adapted, since none of it has any
equivalent in a static SVG image; each program's own doc page states
this explicitly per the plan's own checklist requirement. Where a
program supports it, an optional `-svg <path>` flag writes the
diagram; the numeric tables every program's own calculation produces
print to stdout either way, replacing the original's `.OUT`/`.DAT`
file-save feature the same way earlier tiers did.

### Group 1

| Calculator | Description |
|---|---|
| [cam](calculators/cam.md) | Plate cam profile design |
| [spline](calculators/spline.md) | Natural cubic spline curve fitting |
| [drill](calculators/drill.md) | Drill size lookup, tap drill, and step-drilling calculations |

`drill`, despite being grouped into Tier 3 by the plan's own
data-dependency note (its `.DAT` table was caught by the same
preparatory scan that found Tier 2's data-file-backed programs), turns
out to have no graphics at all: `DRILL.C` never calls a single drawing
primitive, only colored text output. It needed no `internal/svgplot`
work — only its own 371-entry drill table, added to `reference.db`
here (the tool that was meant to do this during Tier 2,
`tools/import-drill`, was only ever a throwaway demonstration of the
pattern and had never actually been run against the real, committed
database). `cam` and `spline` are the two shortest genuinely-graphical
programs, and establish the group's shared drawing conventions
(`internal/svgplot`'s line/box/circle/text primitives, a Cartesian-
y-up window matching the original's own `_setwindow(TRUE,...)`
convention) that the remaining five programs reuse. `spline`'s own
natural-cubic-spline fit (Sedgewick's tridiagonal method) is also
reused, in a more elaborate form, by `profile` later in this tier.

### Group 2

| Calculator | Description |
|---|---|
| [crod](calculators/crod.md) | Gudgeon (wrist) pin position for a slider-crank mechanism |
| [ogive](calculators/ogive.md) | Tangent ogive nose-cone profile and turning schedule |
| [egg](calculators/egg.md) | Turned "egg" shape profile and roughing schedule |

All three share the same `wline`-chain drawing pattern the trial
group established. `crod` is the first program with the interactive
mouse-driven menu (present in five of the eight Tier 3 programs) this
project actually needed to make an explicit drop decision about —
dropped entirely rather than adapted, since a static SVG image has no
equivalent for a live coordinate readout following the mouse cursor.
Its own `.TXT` file has a small but real documentation/code mismatch,
worth noting rather than silently perpetuating: it calls the output
file `CROD.DAT`, but `CROD.C` actually writes `CROD.OUT`. `ogive` and
`egg` are the two Tier 3 programs with no interactive prompts of any
kind — `ogive`'s parameters come entirely from its own `.DAT` file
(a `key=value` syntax, distinct from every other program's positional-
field convention), and `egg`'s are hardwired into the source itself,
confirmed by `EGG.TXT`'s own admission that they "can't be changed
without recompiling." Both compute a nose-cone/egg-shaped turning
profile and an incremental roughing schedule for either a square- or
round-tipped tool; `egg`'s own asymmetric-ellipse formula has a subtle
property worth documenting rather than treating as a bug: it bulges
slightly past its own nominal semi-minor axis before reaching the
formula's algebraically clean midpoint, which is why its roughing
schedule has exactly one gap where no material needs removing.
`EGG.C` also carries a complete, unused copy of `spline`'s own cubic-
spline fitter (declared, fully implemented, never called anywhere in
the file) — dead code, omitted from this port rather than translated.

### Group 3

| Calculator | Description |
|---|---|
| [profile](calculators/profile.md) | Turned profile roughing schedule for an arbitrary, user-specified shape |
| [xymwk](calculators/xymwk.md) | Height-gauge coordinate transforms and measurements: reference, align, distance, pitch circle |

`profile` generalizes `ogive` and `egg`: instead of one closed-form
shape, it takes an arbitrary user-supplied list of `(x, radius)`
points and, between them, interpolates either linearly (the default)
or, over data-file-declared ranges, with a fresh natural cubic spline
fit reusing `spline`'s own fitter. Its shipped example (`PROFILE.DAT`,
a tiny model-tool handle) came with a genuine reference output
(`PROFILE.OUT`) captured from the original DOS binary — the strongest
test oracle available anywhere in Tier 3 — and this conversion matches
it to the printed digit except for one single exact decimal tie
(`0.0625`, at the very first cutting pass), where Go's round-half-to-
even and the original's round-half-away-from-zero `printf` disagree on
the third digit; a display rounding-mode difference, not a computation
error. `PROFILE.C`'s own axial-position loop drifts under repeated
floating-point subtraction; this conversion computes each pass's
position as a multiple of the axial step instead, which avoids that
drift and happens to reproduce the original's exact pass count.

`xymwk` is the tier's odd one out: not a curve/profile calculator with
a roughing schedule, but a coordinate-transform and measurement tool
for points measured with a height gauge (reference, align, distance,
and pitch-circle-through-three-points), and its scope was discussed
and decided separately before conversion — each run performs exactly
one operation on the raw data, rather than replicating the original's
interactive session where each operation's result became the working
coordinate set for the next. It selects its working points by 1-based
index rather than mouse click. Its own `display()` function — fully
implemented, matching `XYMWK.TXT`'s own description of what
`reference`/`align` show, but never actually called anywhere in the
original's `main()` — is wired up here, since a static CLI has no
equivalent for the original's silent on-screen re-plot. Its `circ3`
pitch-circle function carries a known, author-flagged, never-fixed gap
("no error checking... for the anomalous case in which the three
points are colinear," with the author's own suggested fix in the same
comment) — implemented here as an actual error return instead of the
wrong-circle result the original's unchecked call could silently
produce.

## Tier 4

The one-dimensional stock-cutting batch: three programs solving the
same problem (cut a list of needed piece lengths from as few
standard-length bars, or from a pile of odd-length remnants, as
possible) with different heuristics. Originally planned as "the only
stateful program" needing a full atomic-write persistence pattern —
reading and re-checked against the actual source before conversion,
none of the three turns out to write any file at all: all three just
read their own data file and print to stdout, exactly like the rest of
this project's batch.

| Calculator | Description |
|---|---|
| [cuts](calculators/cuts.md) | Per-bar exhaustive combinatorial search, optionally biased toward zero-waste combinations first |
| [cutlist](calculators/cutlist.md) | "Best fit decreasing" heuristic for the same uniform-bar problem |
| [remnant](calculators/remnant.md) | `cutlist`'s heuristic generalized to a heterogeneous pile of remnants, plus saw kerf |

`cuts` is the original author's own greedy, per-bar exhaustive search;
`cutlist` (contributed by Mike Graham) is a different, generally
better heuristic — largest piece first, cut from whichever open bar
has the smallest sufficient remaining room — that the original author
calls "definitely superior... runs faster" in `CUTS.TXT`, while also
noting `cuts` sometimes still wins, hence converting both rather than
picking a winner. `remnant` generalizes `cutlist`'s own heuristic to a
heterogeneous set of available remnant lengths (not one uniform bar)
plus a saw-kerf allowance per cut, and, in doing so, fixes (or simply
never inherited) a subtle asymmetry in `cutlist`'s own best-fit
refinement step: `cutlist` can never prefer a later exact zero-waste
fit over an earlier looser one it already found, while `remnant`'s
version has no such restriction — this conversion preserves both
behaviors exactly rather than treating the two algorithms as
identical. `cuts`'s own "debug" mode (triggered by any command-line
argument in the original, printing internal search diagnostics and
stopping after a single bar) is developer scaffolding, not a
user-facing feature, and is dropped entirely.
