// helixcf computes cutting and setup dimensions for a helical gear,
// following Chuck Fellows' method for machining helical gears on a
// manual mill using a mandrel and angled template.
//
// Converted from HELIXCF.C (M. W. Klotz, 7/10),
// WorkshopUtilities/helixcf. Reference:
// http://www.homemodelenginemachinist.com/index.php?topic=9916.0
package main

import (
	"fmt"
	"math"
	"os"

	"mwkgo/internal/promptio"
)

// helicalGear is the computed dimensions for cutting a helical gear.
type helicalGear struct {
	BlankDiameter    float64
	WholeDepth       float64
	PitchDiameter    float64
	HelixLead        float64
	TemplateAngleDeg float64
}

// computeHelicalGear returns the cutting dimensions for a helical
// gear of the given tooth count, diametral pitch, and helix angle,
// to be cut using a mandrel of the given hub diameter and a template
// of the given thickness.
func computeHelicalGear(teeth, diametralPitch, helixAngleDeg, mandrelHubDiameter, templateThickness float64) helicalGear {
	cosHelixAngle := math.Cos(helixAngleDeg * math.Pi / 180)
	tanHelixAngle := math.Tan(helixAngleDeg * math.Pi / 180)

	pitchDiameter := teeth / (diametralPitch * cosHelixAngle)
	helixLead := math.Pi * pitchDiameter / tanHelixAngle
	blankDiameter := pitchDiameter + 2/diametralPitch

	var wholeDepth float64
	if diametralPitch <= 20 {
		wholeDepth = 2.2/diametralPitch + 0.002
	} else {
		wholeDepth = 2.157 / diametralPitch
	}

	templateAngleDeg := math.Atan(helixLead/(math.Pi*mandrelHubDiameter+templateThickness)) * 180 / math.Pi

	return helicalGear{
		BlankDiameter:    blankDiameter,
		WholeDepth:       wholeDepth,
		PitchDiameter:    pitchDiameter,
		HelixLead:        helixLead,
		TemplateAngleDeg: templateAngleDeg,
	}
}

func main() {
	prompter, err := promptio.New(os.Stdin, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, "helixcf:", err)
		os.Exit(1)
	}

	fmt.Println("Chuck Fellows' Helical Gear Calculations")

	teeth := prompter.Float("\nNumber of teeth", 6.0)
	diametralPitch := prompter.Float("Diametral Pitch", 40.0)
	helixAngle := prompter.Float("Helix angle (deg)", 80.0)
	mandrelHubDiameter := prompter.Float("Mandrel hub diameter (in)", 1.0)
	templateThickness := prompter.Float("Template thickness (in)", 0.125)

	gear := computeHelicalGear(teeth, diametralPitch, helixAngle, mandrelHubDiameter, templateThickness)

	fmt.Printf("\nDiametral Pitch = %.4f\n", diametralPitch)
	fmt.Printf("Number of teeth = %.0f\n", teeth)
	fmt.Printf("Helix Angle = %.4f deg\n", helixAngle)
	fmt.Printf("Gear Blank Diameter = %.4f in\n", gear.BlankDiameter)
	fmt.Printf("Whole Depth = %.4f in\n", gear.WholeDepth)
	fmt.Printf("Pitch Diameter = %.4f in\n", gear.PitchDiameter)
	fmt.Printf("Helix Lead = %.4f in\n", gear.HelixLead)
	fmt.Printf("Template angle = %.4f deg\n", gear.TemplateAngleDeg)
}
