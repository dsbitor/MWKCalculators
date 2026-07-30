# helixcf

Helical gear cutting dimensions.

**Converted from:** `HELIXCF.C` (M. W. Klotz, 7/10),
`MWKC/WorkshopUtilities/helixcf.zip`. Chuck Fellows' method.
Reference: http://www.homemodelenginemachinist.com/index.php?topic=9916.0
**Go source:** `MWKGo/helixcf/helixcf.go`

## Purpose

Chuck Fellows' method machines a helical gear on a manual mill
using a mandrel and an angled template to guide the workpiece
through the cutter at the correct helix angle. Given the tooth
count, diametral pitch, helix angle, mandrel hub diameter, and
template thickness, this program computes the gear blank
diameter, whole tooth depth, pitch diameter, helix lead, and the
angle to cut the guide template at.

No `.TXT` file was included with the original program; this
purpose statement is drawn from the `.C` file's own header
comment and the referenced forum thread's subject.

## Inputs

| Prompt | Default |
|---|---|
| Number of teeth | 6 |
| Diametral Pitch | 40 |
| Helix angle | 80 deg |
| Mandrel hub diameter | 1 in |
| Template thickness | 0.125 in |

## Output

Gear blank diameter, whole depth, pitch diameter, helix lead,
and template angle.

## Method

```
pitchDiameter = teeth / (diametralPitch * cos(helixAngle))
helixLead = pi * pitchDiameter / tan(helixAngle)
blankDiameter = pitchDiameter + 2/diametralPitch
wholeDepth = 2.2/diametralPitch + 0.002   if diametralPitch <= 20
           = 2.157/diametralPitch          otherwise
templateAngle = atan(helixLead / (pi*mandrelHubDiameter + templateThickness))
```

## Worked Example

No worked example was included with the original program. As an
independently verifiable check: a zero degree helix angle is a
plain spur gear, whose pitch diameter reduces to
`teeth/diametralPitch` exactly, the standard spur gear formula
independent of this code's helix-specific trigonometry.
