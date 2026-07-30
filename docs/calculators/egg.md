# egg

Turned "egg" shape profile and roughing schedule.

**Converted from:** `EGG.C` (M. W. Klotz), `MWKC/WorkshopUtilities/egg.zip`
**Go source:** `MWKGo/egg/egg.go`

## Purpose

A tribute to a scene in Nevil Shute's novel "Trustee from the
Toolroom," in which a model engineer turns metal eggs for a toy duck
to sit on. Computes the profile of an asymmetrical-ellipse "egg" shape
and an incremental turning schedule for roughing it out on a lathe.

Unlike every other Tier 3 program, `egg` takes **no input at all**.
`EGG.C` hardwires its own stock diameter, tool width, axial cutting
step, and the egg shape's own semi-axes and asymmetry factor; the
original author explained why directly: *"Since it's questionable
just how many people want to turn metal eggs, I hardwired these
values into the program. They can't be changed without recompiling."*
This conversion keeps that behavior rather than inventing configurable
input the original never had.

## Inputs

None. Every parameter is fixed: stock diameter 1 in, tool width
1/16 in, axial cutting increment 0.01 in, egg semi-axes `a=0.75`,
`b=0.5`, asymmetry `k=0.6`.

## Output

The fixed configuration, then a roughing schedule: each pass's axial
position, depth of cut, and resulting diameter. A step is only listed
if material actually needs removing there; a `...` line marks any gap
where one or more positions needed no cut (this happens once, in the
egg's own bulge region — see Method below). If `-svg` is given, a
diagram of the egg profile with each roughing pass overlaid.

## Method

The egg curve is an asymmetrical ellipse,
`x²/a² + (1 − k·x)·y²/b² = 1`, solved for `y`:
`y = b × √((1 − x²/a²) / (1 − k·x))`. The diameter of the egg at its
thickest part is nominally `2×b`, and `x=a` is an exact,
algebraically clean point on the curve where `y=b` regardless of the
asymmetry factor `k` — but it is *not* actually the profile's single
widest point: the `(1 − k·x)` term shrinks the denominator for
`x` between 0 and `a`, which lets the curve bulge slightly *past* `b`
a bit before reaching `x=a`. That's the entire purpose of the
asymmetry factor (an egg is visibly lopsided, wider toward one end
than a plain ellipse would be), and it's why the roughing schedule can
briefly need no cut at all in that bulge region — the theoretical
profile there already exceeds the stock's own radius, so there's
nothing left to remove.

`EGG.C` also declares and fully implements a natural cubic spline
fitter (`spline()`/`seval()`), identical in every particular to
[spline](spline.md)'s own — but never actually calls it anywhere in
the file. This conversion omits that dead code entirely rather than
translating unused machinery.

## Worked Example

The original program's own hardwired parameters are stock diameter
1", `b = 0.5"`, `a = 0.75"`, `k = 0.6`. This conversion's
tests confirm those exact values are what `eggSemiAxes` computes; that
the profile's radius at the pointed tip (`x=0`) is exactly zero; that
at `x=a` the radius is exactly `b` (the hand-derivable point described
above); that beyond `x=2a` the profile correctly falls back to a small
flat nub rather than evaluating the egg formula where it would become
mathematically undefined; and, as an independent sanity check on the
formula itself, that setting the asymmetry factor to zero reduces it
to a plain, symmetric ellipse. A manual run confirms the roughing
schedule's own single expected gap, in the bulge region described
above.
