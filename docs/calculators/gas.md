# gas

Perfect gas law (PV = nRT) calculator.

**Converted from:** `GAS.C`, `MWKC/WorkshopUtilities/gas.zip`
**Go source:** `MWKGo/gas/gas.go`

## Purpose

Solves the ideal gas law for whichever of pressure, volume, moles, or
temperature is unknown, given the other three, accepting each in a
choice of common units (pressure in atmospheres, psi, or kilopascals;
volume in liters, cubic feet, or cubic meters; temperature in Kelvin,
Celsius, or Fahrenheit).

## Inputs

| Prompt | Default unit |
|---|---|
| Pressure [atm]osphere, (psi), (kpa)scal | atmospheres |
| Volume [l]iter, (cft), (cm)eter | liters |
| Moles | mole |
| Temperature [K], (C), (F) | Kelvin |

Enter whatever three you know; leave the fourth blank.

## Output

Pressure (atm/psi/kPa), volume (l/ft³/m³), temperature (K/F/C), and
moles.

## Method

```
R = 8.20545E-2 liter-atm/(degK-mole)

missing P: P = nRT/V
missing V: V = nRT/P
missing n: n = PV/(RT)
missing T: T = PV/(nR)
```

Unit-suffixed input is parsed for a leading numeric value plus a unit
hint anywhere in the string (matching C's `atof`/`strstr` combination
in the original). One quirk carried forward unchanged: the pressure
prompt abbreviates kilopascals as `(kpa)scal`, but the original
program's own unit-detection check looks for the substring `"pas"`,
not `"kpa"`. A natural-looking `"100kpa"` therefore falls through to
being read as plain atmospheres; only `"100kpascal"` (or anything
else containing `"pas"`) is actually recognized as kilopascals. This
is almost certainly an inconsistency between the prompt text and the
check in the original, but it's what the original program actually
does, so it's preserved here rather than silently corrected.

## Worked Example

No worked numeric example was included with the original program
(`GAS.TXT` explains the law and units but includes no sample run,
though it does state that one mole occupies 22.4140 liters at
standard conditions). As an independently verifiable check: one mole
of gas at standard conditions (0°C, 1 atmosphere) is confirmed to
occupy approximately 22.4140 liters, a well known physical constant
independent of this program's own formula; confirmed in this
conversion's tests, along with an internal consistency check that
solving for each of the four quantities in turn, given the other
three, recovers the same self-consistent state.
