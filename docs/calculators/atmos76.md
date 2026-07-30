# atmos76

1976 U.S. Standard Atmosphere model.

**Converted from:** `ATMOS76.C`, `MWKC/Misc/atmos76.zip`
**Go source:** `MWKGo/atmos76/atmos76.go`

## Purpose

Given an altitude, computes temperature, pressure, and density
ratios (relative to sea level), plus their absolute values and the
local speed of sound, per the 1976 U.S. Standard Atmosphere model
(valid from 0 to 86 km).

## Inputs

| Prompt | Default |
|---|---|
| Altitude | 1000 ft |

## Output

Altitude in km/miles/ft; temperature, pressure, and density ratios
relative to sea level; absolute temperature (K/C/F), pressure
(N/m^2/psi), and density (kg/m^3/slug per ft^3); and speed of sound
(kph/mph).

## Method

The atmosphere is modeled as 7 layers, each with a known base
geopotential height, base temperature, base pressure ratio, and
temperature lapse rate. Geometric altitude is converted to
geopotential altitude (`h = alt*earthRadius/(alt+earthRadius)`), the
containing layer is found, and:

```
localTemp = layerBaseTemp + lapseRate*(h - layerBaseHeight)
temperatureRatio = localTemp / seaLevelTemp

pressureRatio = layerBasePressureRatio * exp(-gasConstant*deltaH/layerBaseTemp)
              (isothermal layers, lapseRate == 0)
            = layerBasePressureRatio * (layerBaseTemp/localTemp)^(gasConstant/lapseRate)
              (otherwise)

densityRatio = pressureRatio / temperatureRatio
```

This is the standard, widely published 1976 U.S. Standard Atmosphere
model (also known as ICAO/ISA below 32 km), not specific to this
program.

## Worked Example

No worked numeric example was included with the original program. As
an independently verifiable check: at each layer's own tabulated
base height (where `deltaH = 0`), the computed ratios must exactly
reproduce that layer's own tabulated base values, independent of the
interpolation/extrapolation formula used between breakpoints;
confirmed in this conversion's tests for sea level and all 7 layer
boundaries.
