# buoy

Buoy immersion depth calculator.

**Converted from:** `BUOY.C`, `MWKC/Misc/buoy.zip`
**Go source:** `MWKGo/buoy/buoy.go`

## Purpose

Per Archimedes' principle, a floating buoy sinks into a liquid until
the weight of the liquid it displaces equals its own weight plus any
load it carries. Given the liquid's density, the buoy's weight, and
its shape and dimensions (sphere, horizontal cylinder, vertical
cylinder, or rectangular box), this program computes how deep the
buoy sinks.

## Inputs

| Prompt | Default |
|---|---|
| Density of liquid | 0.0361 lb/in^3 (fresh water) |
| Weight of buoy and applied load | 3.5 lb |
| Buoy shape (a-d) | (required) |
| Diameter/Length/Width/Height (depends on shape) | varies |

## Output

Immersion depth (in), or a message that the buoy will sink along
with the maximum weight that shape and size could support.

## Method

Vertical cylinders and boxes have a submerged volume that's a simple
linear function of depth, so their immersion depth is solved for
directly:

```
verticalCylinder: depth = weight / (0.25*pi*diameter^2*liquidDensity)
box:               depth = weight / (length*width*liquidDensity)
```

Spheres and horizontal cylinders have a submerged volume that's a
nonlinear function of depth (a spherical cap, or a circular segment
respectively), so their depth is found by a bounded binary search
over the original program's own `binsearch` function, converging on
the depth at which the displaced liquid's weight equals the target:

```
sphere:     submergedVolume(h) = pi*h^2*(1.5*d - h)/3
horizontal
cylinder:   submergedArea(h) = 0.5*(r^2*angle - z*chord),
            where z = r-h, angle = 2*acos(z/r), chord = 2*r*sin(angle/2)
```

## Worked Example

No worked numeric example was included with the original program. As
independently verifiable checks: a sphere or cylinder submerged to
its own full diameter must displace exactly its own full volume (or
full cross-sectional area), confirmed algebraically; and, for every
shape, the buoyant force computed at the returned immersion depth
(density times displaced volume) must equal the target weight, either
exactly (for the two directly-solved shapes) or within the binary
search's own tolerance (for the two searched shapes) — Archimedes'
principle itself used as the test oracle, not a re-run of the
formula; both confirmed in this conversion's tests.
