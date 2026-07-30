// unit is a general-purpose units-conversion utility: 246 named
// units, 23 metric-style prefixes, and 7 fundamental dimensions
// (length, mass, time, angle, solid angle, charge, amount), supporting
// arbitrary compound-unit expressions ((prefix)UNIT[^exponent], with
// an optional "/" denominator), dimensional-compatibility checking, a
// breakdown of any unit into its fundamental dimensions, listing
// every unit of a given physical quantity type, and defining new
// units on the fly for the rest of a session.
//
// The unit and prefix tables are universal reference data (the same
// for anyone), so they live in the shared reference.db SQLite
// database, the same mechanism used by fits, speed, gage, expand,
// findthrd, and weight earlier in this project; see
// ai/plans/c-to-go-conversion-plan.md, "Data-file strategy for
// Tier 2".
//
// Converted from UNIT.C (M. W. Klotz), Misc/unit.
package main

import (
	"context"
	"errors"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"

	"mwkgo/internal/promptio"
	"mwkgo/internal/refdata"
)

// nPrim is the number of fundamental dimensions, matching UNIT.C's
// own NPRIM.
const nPrim = 7

// primaryNames and siNames are UNIT.C's own primary[]/uprimary[]
// arrays: the display name of each fundamental dimension, and the SI
// unit that measures it.
var primaryNames = [nPrim]string{"LENGTH", "MASS", "TIME", "ANGLE", "SOLIDANGLE", "CHARGE", "AMOUNT"}
var siNames = [nPrim]string{"METER", "KILOGRAM", "SEC", "RADIAN", "STERADIAN", "COULOMB", "ND"}

// unitDef is one named unit: its conversion factor (desired units per
// primary/SI unit) and its exponents against the 7 fundamental
// dimensions.
type unitDef struct {
	Name  string
	Fact  float64
	PFact [nPrim]int
}

// prefixDef is one named metric-style prefix and its multiplier.
type prefixDef struct {
	Name string
	Val  float64
}

var errNoUnitSpecified = errors.New("no unit specified")

// ufind returns the index of the unit named name (case-insensitive
// exact match), or -1 if not found, matching UNIT.C's own ufind().
func ufind(name string, units []unitDef) int {
	for i, u := range units {
		if strings.EqualFold(name, u.Name) {
			return i
		}
	}
	return -1
}

// deplural strips a common English plural ending from n, unless n is
// already the exact name of a known unit (in which case it is
// returned unchanged) — matching UNIT.C's own deplural(), using the
// original string's own suffixes throughout rather than re-checking
// against an already-shortened copy (see UNIT.C's own comment: this
// matters only for a unit ending in "IES", none of which exist in the
// shipped data, so the two behave identically in practice).
func deplural(n string, units []unitDef) string {
	if ufind(n, units) != -1 {
		return n
	}
	l := len(n)
	switch {
	case l >= 3 && strings.EqualFold(n[l-3:], "IES"):
		return n[:l-3]
	case l >= 2 && strings.EqualFold(n[l-2:], "ES"):
		return n[:l-2]
	case l >= 1 && strings.EqualFold(n[l-1:], "S"):
		return n[:l-1]
	default:
		return n
	}
}

// unitExpr is a parsed unit expression such as "MILES/HOUR" or
// "KG/M^3", matching UNIT.C's own struct unitexp.
type unitExpr struct {
	NumPrefixIdx int // -1 if no prefix
	NumUnitIdx   int
	NumExp       int
	HasDen       bool
	DenPrefixIdx int
	DenUnitIdx   int
	DenExp       int
	PU           [nPrim]int // net dimension vector, for compatibility checks
}

// uanal parses a unit expression into a unitExpr, matching UNIT.C's
// own uanal() exactly: split on "/" into numerator/denominator,
// extract an optional "^exponent" from each half, depluralize, strip
// a leading prefix name (only if the bare string isn't already a
// known unit), and resolve what remains to a unit by exact
// case-insensitive name.
func uanal(expr string, units []unitDef, prefixes []prefixDef) (unitExpr, error) {
	p := unitExpr{NumPrefixIdx: -1, DenPrefixIdx: -1, NumUnitIdx: -1, DenUnitIdx: -1, NumExp: 1, DenExp: 1}

	var numStr, denStr string
	if idx := strings.Index(expr, "/"); idx >= 0 {
		numStr = strings.TrimSpace(expr[:idx])
		denStr = strings.TrimSpace(expr[idx+1:])
		p.HasDen = true
	} else {
		numStr = strings.TrimSpace(expr)
	}
	if numStr == "" {
		return unitExpr{}, errNoUnitSpecified
	}

	if idx := strings.Index(numStr, "^"); idx >= 0 {
		exp, err := strconv.Atoi(strings.TrimSpace(numStr[idx+1:]))
		if err != nil {
			return unitExpr{}, fmt.Errorf("numerator exponent in %q: %w", expr, err)
		}
		p.NumExp = exp
		numStr = numStr[:idx]
	}
	if p.HasDen {
		if idx := strings.Index(denStr, "^"); idx >= 0 {
			exp, err := strconv.Atoi(strings.TrimSpace(denStr[idx+1:]))
			if err != nil {
				return unitExpr{}, fmt.Errorf("denominator exponent in %q: %w", expr, err)
			}
			p.DenExp = exp
			denStr = denStr[:idx]
		}
	}

	numStr = deplural(numStr, units)
	if p.HasDen {
		denStr = deplural(denStr, units)
	}

	if ufind(numStr, units) == -1 {
		for i, pre := range prefixes {
			if idx := strings.Index(numStr, pre.Name); idx >= 0 {
				p.NumPrefixIdx = i
				numStr = numStr[idx+len(pre.Name):]
				break
			}
		}
	}
	if p.HasDen && ufind(denStr, units) == -1 {
		for i, pre := range prefixes {
			if idx := strings.Index(denStr, pre.Name); idx >= 0 {
				p.DenPrefixIdx = i
				denStr = denStr[idx+len(pre.Name):]
				break
			}
		}
	}

	nk := ufind(numStr, units)
	if nk == -1 {
		return unitExpr{}, fmt.Errorf("unit %q not found", numStr)
	}
	p.NumUnitIdx = nk

	if p.HasDen {
		dk := ufind(denStr, units)
		if dk == -1 {
			return unitExpr{}, fmt.Errorf("unit %q not found", denStr)
		}
		p.DenUnitIdx = dk
	}

	for i := 0; i < nPrim; i++ {
		p.PU[i] = units[p.NumUnitIdx].PFact[i] * p.NumExp
		if p.HasDen {
			p.PU[i] -= units[p.DenUnitIdx].PFact[i] * p.DenExp
		}
	}
	return p, nil
}

// uconv converts x from the "from" expression to the "to" expression,
// also returning the value expressed in primary (SI) units, matching
// UNIT.C's own uconv() exactly — including its "!= -1" guard on
// from.DenExp (apparently meant to read "!= 1" like every other
// exponent guard here, but harmless either way since pow(v,1) is a
// no-op: the guard only ever skips the computation when DenExp is
// exactly -1, which would otherwise correctly compute pow(fd,-1)
// anyway).
func uconv(x float64, from, to unitExpr, units []unitDef, prefixes []prefixDef) (converted, primary float64) {
	fn, tn, fd, td := 1.0, 1.0, 1.0, 1.0

	if from.NumPrefixIdx != -1 {
		fn = prefixes[from.NumPrefixIdx].Val
	}
	if from.NumUnitIdx != -1 {
		fn = fn / units[from.NumUnitIdx].Fact
	}
	if from.NumExp != 1 {
		fn = math.Pow(fn, float64(from.NumExp))
	}
	if to.NumPrefixIdx != -1 {
		tn = 1 / prefixes[to.NumPrefixIdx].Val
	}
	if to.NumUnitIdx != -1 {
		tn = tn * units[to.NumUnitIdx].Fact
	}
	if to.NumExp != 1 {
		tn = math.Pow(tn, float64(to.NumExp))
	}

	if from.DenPrefixIdx != -1 {
		fd = 1 / prefixes[from.DenPrefixIdx].Val
	}
	if from.DenUnitIdx != -1 {
		fd = fd * units[from.DenUnitIdx].Fact
	}
	if from.DenExp != -1 {
		fd = math.Pow(fd, float64(from.DenExp))
	}
	if to.DenPrefixIdx != -1 {
		td = prefixes[to.DenPrefixIdx].Val
	}
	if to.DenUnitIdx != -1 {
		td = td / units[to.DenUnitIdx].Fact
	}
	if to.DenExp != 1 {
		td = math.Pow(td, float64(to.DenExp))
	}

	primary = x * fn * fd
	return x * fn * fd * tn * td, primary
}

// compatible reports whether a and b describe the same physical
// quantity (identical dimension vectors), the check convert()
// performs before allowing a conversion.
func compatible(a, b [nPrim]int) bool {
	return a == b
}

// engprnt formats x the way UNIT.C's own engprnt() does: like C's "%G"
// (six significant digits, %E or plain decimal chosen automatically),
// except that when scientific notation is used, its exponent is
// forced to a multiple of three (true "engineering notation").
func engprnt(x float64) string {
	plain := fmt.Sprintf("%.6G", x)
	if !strings.ContainsAny(plain, "Ee") {
		return plain
	}

	absX := math.Abs(x)
	k := int(math.Floor(math.Log10(absX)))
	z := absX * math.Pow(10, float64(-k))
	for (intAbs(k)%3 != 0) || (z < 1) {
		z *= 10
		k--
	}
	sign := 1.0
	if x < 0 {
		sign = -1
	}
	return fmt.Sprintf("%10.6fE%+04d", sign*z, k)
}

func intAbs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// pprim formats a dimension vector as a "(A * B^2) / (C * D^3)"-style
// breakdown, matching UNIT.C's own pprim(): positive exponents in the
// first parenthesized group, negative exponents (shown as their
// absolute value) in a second group after a "/", omitted entirely if
// there are none.
func pprim(names [nPrim]string, pu [nPrim]int) string {
	var b strings.Builder
	b.WriteString(" (")
	hasNeg, wrote := false, false
	for i := 0; i < nPrim; i++ {
		e := pu[i]
		switch {
		case e == 0:
			continue
		case e < 0:
			hasNeg = true
			continue
		default:
			if wrote {
				b.WriteString(" * ")
			}
			if e == 1 {
				b.WriteString(names[i])
			} else {
				fmt.Fprintf(&b, "%s^%d", names[i], e)
			}
			wrote = true
		}
	}
	b.WriteString(")")
	if !hasNeg {
		return b.String()
	}

	b.WriteString(" / (")
	wrote = false
	for i := 0; i < nPrim; i++ {
		e := pu[i]
		if e >= 0 {
			continue
		}
		if wrote {
			b.WriteString(" * ")
		}
		if e == -1 {
			b.WriteString(names[i])
		} else {
			fmt.Fprintf(&b, "%s^%d", names[i], -e)
		}
		wrote = true
	}
	b.WriteString(")")
	return b.String()
}

// compoundType is one of UNIT.C's own compound() menu entries: a
// physical quantity type and the dimension vector it corresponds to.
type compoundType struct {
	Name string
	PF   [nPrim]int
}

// compoundTypes matches UNIT.C's own compound() switch statement
// exactly (dimension order: LENGTH,MASS,TIME,ANGLE,SOLIDANGLE,CHARGE,AMOUNT).
var compoundTypes = []compoundType{
	{"amount", [nPrim]int{0, 0, 0, 0, 0, 0, 1}},
	{"mass", [nPrim]int{0, 1, 0, 0, 0, 0, 0}},
	{"time", [nPrim]int{0, 0, 1, 0, 0, 0, 0}},
	{"angle", [nPrim]int{0, 0, 0, 1, 0, 0, 0}},
	{"length", [nPrim]int{1, 0, 0, 0, 0, 0, 0}},
	{"area", [nPrim]int{2, 0, 0, 0, 0, 0, 0}},
	{"volume", [nPrim]int{3, 0, 0, 0, 0, 0, 0}},
	{"speed", [nPrim]int{1, 0, -1, 0, 0, 0, 0}},
	{"acceleration", [nPrim]int{1, 0, -2, 0, 0, 0, 0}},
	{"force", [nPrim]int{1, 1, -2, 0, 0, 0, 0}},
	{"torque or energy", [nPrim]int{2, 1, -2, 0, 0, 0, 0}},
	{"power", [nPrim]int{2, 1, -3, 0, 0, 0, 0}},
	{"pressure", [nPrim]int{-1, 1, -2, 0, 0, 0, 0}},
	{"angular speed", [nPrim]int{0, 0, -1, 1, 0, 0, 0}},
	{"solid angle", [nPrim]int{0, 0, 0, 0, 1, 0, 0}},
}

// unitsOfType returns every unit (in the given display order) whose
// dimension vector exactly matches pf.
func unitsOfType(order []int, units []unitDef, pf [nPrim]int) []string {
	var names []string
	for _, idx := range order {
		if units[idx].PFact == pf {
			names = append(names, units[idx].Name)
		}
	}
	return names
}

// applyPrimaryUnit accumulates one (unit, exponent) step of an
// on-the-fly unit definition, matching UNIT.C's own define() inner
// loop: the new unit's factor is repeatedly multiplied (positive
// exponent) or divided (negative exponent) by the given primary
// unit's own factor, and its dimension vector accumulates exp times
// that unit's own dimension vector.
func applyPrimaryUnit(fact float64, pfact [nPrim]int, primaryFact float64, primaryPFact [nPrim]int, exp int) (float64, [nPrim]int) {
	for i := 0; i < intAbs(exp); i++ {
		if exp > 0 {
			fact *= primaryFact
		} else {
			fact /= primaryFact
		}
	}
	for i := 0; i < nPrim; i++ {
		pfact[i] += exp * primaryPFact[i]
	}
	return fact, pfact
}

func loadUnitsAndPrefixes(ctx context.Context) ([]unitDef, []prefixDef, error) {
	db, err := refdata.Open(ctx)
	if err != nil {
		return nil, nil, err
	}
	defer db.Close()

	prows, err := db.QueryContext(ctx, `SELECT name, value FROM unit_prefixes`)
	if err != nil {
		return nil, nil, fmt.Errorf("query unit_prefixes: %w", err)
	}
	var prefixes []prefixDef
	for prows.Next() {
		var p prefixDef
		if err := prows.Scan(&p.Name, &p.Val); err != nil {
			prows.Close()
			return nil, nil, fmt.Errorf("scan prefix: %w", err)
		}
		prefixes = append(prefixes, p)
	}
	if err := prows.Err(); err != nil {
		prows.Close()
		return nil, nil, fmt.Errorf("iterate unit_prefixes: %w", err)
	}
	prows.Close()

	urows, err := db.QueryContext(ctx,
		`SELECT name, factor, dim_length, dim_mass, dim_time, dim_angle, dim_solidangle, dim_charge, dim_amount FROM unit_definitions`)
	if err != nil {
		return nil, nil, fmt.Errorf("query unit_definitions: %w", err)
	}
	defer urows.Close()
	var units []unitDef
	for urows.Next() {
		var u unitDef
		if err := urows.Scan(&u.Name, &u.Fact, &u.PFact[0], &u.PFact[1], &u.PFact[2], &u.PFact[3], &u.PFact[4], &u.PFact[5], &u.PFact[6]); err != nil {
			return nil, nil, fmt.Errorf("scan unit: %w", err)
		}
		units = append(units, u)
	}
	if err := urows.Err(); err != nil {
		return nil, nil, fmt.Errorf("iterate unit_definitions: %w", err)
	}
	return units, prefixes, nil
}

// alphabeticalOrder returns unit indices in case-insensitive
// alphabetical order by name, matching UNIT.C's own sort()/alpha[].
func alphabeticalOrder(units []unitDef) []int {
	order := make([]int, len(units))
	for i := range order {
		order[i] = i
	}
	sort.SliceStable(order, func(i, j int) bool {
		return strings.ToLower(units[order[i]].Name) < strings.ToLower(units[order[j]].Name)
	})
	return order
}

func main() {
	ctx := context.Background()
	units, prefixes, err := loadUnitsAndPrefixes(ctx)
	if err != nil {
		fmt.Fprintln(os.Stderr, "unit:", err)
		os.Exit(1)
	}
	order := alphabeticalOrder(units)

	prompter, err := promptio.New(os.Stdin, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, "unit:", err)
		os.Exit(1)
	}

	fmt.Printf("NUMBER OF FUNDAMENTAL UNITS = %d\n", nPrim)
	fmt.Printf("NUMBER OF PREFIXES = %d\n", len(prefixes))
	fmt.Printf("NUMBER OF UNITS = %d\n", len(units))
	printMenu()

	var lastFrom string

	for {
		op := strings.ToLower(prompter.Line("? "))
		if op == "" {
			op = "j"
		}
		switch op[0] {
		case 'a':
			runCompound(prompter, units, order)
		case 'b':
			runBreakdown(prompter, units)
		case 'c':
			runCompatible(prompter, units, order)
		case 'd':
			units, order = runDefine(prompter, units, order)
		case 'f':
			for i := 0; i < nPrim; i++ {
				fmt.Printf("%-15s %-15s\n", primaryNames[i], siNames[i])
			}
		case 'h':
			printMenu()
		case 'i':
			runListByLetter(prompter, units, order)
		case ' ', 'j':
			lastFrom = runConvert(prompter, units, prefixes, "")
		case 'l':
			lastFrom = runConvert(prompter, units, prefixes, lastFrom)
		case 'p':
			for _, p := range prefixes {
				fmt.Printf("%-12s%-5.0E\n", p.Name, p.Val)
			}
		case 'u':
			for _, idx := range order {
				fmt.Printf("%-15s ", units[idx].Name)
			}
			fmt.Println()
		case 'q':
			return
		default:
			fmt.Println("NOT A VALID OPTION")
		}
	}
}

func printMenu() {
	fmt.Println("a  show all units of a given physical quantity type")
	fmt.Println("b  breakdown of a unit into fundamental dimensions")
	fmt.Println("c  show units compatible with (same dimensions as) a given unit")
	fmt.Println("d  define a new unit")
	fmt.Println("f  show the 7 fundamental units")
	fmt.Println("h  show this menu")
	fmt.Println("i  list units starting with a given letter")
	fmt.Println("j  convert units")
	fmt.Println("l  repeat last conversion with a new \"to\" unit")
	fmt.Println("p  list prefixes")
	fmt.Println("u  list all units")
	fmt.Println("q  quit")
}

func runCompound(p *promptio.Prompter, units []unitDef, order []int) {
	for i, c := range compoundTypes {
		fmt.Printf("%d  %s\n", i, c.Name)
	}
	t := p.Int("Desired unit type", 0)
	if t < 0 || t >= len(compoundTypes) {
		fmt.Println("NOT A VALID OPTION")
		return
	}
	names := unitsOfType(order, units, compoundTypes[t].PF)
	fmt.Println(strings.Join(names, " "))
}

func runBreakdown(p *promptio.Prompter, units []unitDef) {
	name := p.Line("UNIT NAME ? ")
	k := ufind(strings.TrimSpace(name), units)
	if k == -1 {
		fmt.Println("NOT A PRIMARY UNIT")
		return
	}
	u := units[k]
	fmt.Printf("%s = %s = %G %s\n", u.Name, pprim(primaryNames, u.PFact), 1/u.Fact, pprim(siNames, u.PFact))
}

func runCompatible(p *promptio.Prompter, units []unitDef, order []int) {
	name := p.Line("PRIMARY UNIT NAME ? ")
	k := ufind(strings.TrimSpace(name), units)
	if k == -1 {
		fmt.Println("UNIT NOT FOUND")
		return
	}
	var names []string
	for _, idx := range order {
		if units[idx].PFact == units[k].PFact {
			names = append(names, units[idx].Name)
		}
	}
	fmt.Println(strings.Join(names, " "))
}

func runListByLetter(p *promptio.Prompter, units []unitDef, order []int) {
	letter := strings.ToUpper(strings.TrimSpace(p.Line("INITIAL LETTER ? ")))
	if letter == "" {
		return
	}
	var names []string
	for _, idx := range order {
		if strings.HasPrefix(units[idx].Name, letter) {
			names = append(names, units[idx].Name)
		}
	}
	fmt.Println(strings.Join(names, " "))
}

func runConvert(p *promptio.Prompter, units []unitDef, prefixes []prefixDef, repeatFrom string) string {
	fmt.Println("UNIT SYNTAX:    (prefix)UNIT[^m [[/] (prefix)UNIT[^n]]]")

	var valueStr, fromStr string
	if repeatFrom == "" {
		for {
			line := strings.ToUpper(strings.TrimSpace(p.Line("[VALUE] [UNITS] TO CONVERT (or Quit) ? ")))
			if line == "" {
				continue
			}
			if line == "Q" {
				return repeatFrom
			}
			if idx := strings.IndexByte(line, ' '); idx >= 0 {
				valueStr, fromStr = line[:idx], strings.TrimSpace(line[idx+1:])
			} else if len(line) > 0 && isAlpha(line[0]) {
				valueStr, fromStr = "1", line
			} else {
				valueStr, fromStr = line, ""
			}
			if fromStr == "" {
				fromStr = strings.ToUpper(strings.TrimSpace(p.Line("FROM UNIT(S) ? ")))
			}
			break
		}
	} else {
		valueStr, fromStr = "1", repeatFrom
	}

	x, err := strconv.ParseFloat(valueStr, 64)
	if err != nil {
		fmt.Printf("%q is not a number\n", valueStr)
		return repeatFrom
	}

	from, err := uanal(fromStr, units, prefixes)
	if err != nil {
		fmt.Println(err)
		return repeatFrom
	}

	toStr := strings.ToUpper(strings.TrimSpace(p.Line("TO UNIT(S) ? ")))
	to, err := uanal(toStr, units, prefixes)
	if err != nil {
		fmt.Println(err)
		return fromStr
	}

	if !compatible(from.PU, to.PU) {
		fmt.Println("INCOMPATIBLE UNITS")
		return fromStr
	}

	y, primary := uconv(x, from, to, units, prefixes)
	fmt.Printf("%s %s = %s %s = %s%s\n",
		engprnt(x), fromStr, engprnt(y), toStr, engprnt(primary), pprim(siNames, from.PU))
	return fromStr
}

func isAlpha(b byte) bool {
	return (b >= 'A' && b <= 'Z') || (b >= 'a' && b <= 'z')
}

func runDefine(p *promptio.Prompter, units []unitDef, order []int) ([]unitDef, []int) {
	name := strings.ToUpper(strings.TrimSpace(p.Line("NEW UNIT NAME ? ")))
	if name == "" {
		return units, order
	}
	if ufind(name, units) != -1 {
		fmt.Println("name already in use")
		return units, order
	}
	if len(name) == 0 || !isAlpha(name[0]) {
		fmt.Println("must start with a-z")
		return units, order
	}

	fact := 1.0
	var pfact [nPrim]int
	for {
		unitName := strings.ToUpper(strings.TrimSpace(p.Line("PRIMARY UNIT [default=done] ? ")))
		if unitName == "" {
			break
		}
		k := ufind(unitName, units)
		if k == -1 {
			fmt.Println("unit not found")
			continue
		}
		exp := p.Int("INTEGER EXPONENT", 1)
		if exp == 0 {
			fmt.Println("zero not allowed")
			continue
		}
		fact, pfact = applyPrimaryUnit(fact, pfact, units[k].Fact, units[k].PFact, exp)
	}

	fmt.Printf("to save this unit, add the line below to UNIT.DAT\n")
	fmt.Printf("%s=%g", name, fact)
	for i := 0; i < nPrim; i++ {
		fmt.Printf(",%d", pfact[i])
	}
	fmt.Println()

	units = append(units, unitDef{Name: name, Fact: fact, PFact: pfact})
	return units, alphabeticalOrder(units)
}
