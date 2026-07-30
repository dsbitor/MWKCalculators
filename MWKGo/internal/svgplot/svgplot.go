// Package svgplot renders the same line/box/circle/text primitives
// the original Tier 3 calculators drew with DOS graphics calls
// (_moveto_w, _lineto_w, _rectangle, _ellipse, _outtext, by way of
// each program's own wline/wbox/wcircle/pwin wrapper functions,
// documented in ai/plans/c-to-go-conversion-plan.md's "Graphics
// scope" resolution), but emits SVG instead of drawing to a screen.
//
// Every coordinate given to a Canvas method is in the program's own
// "window" (user, floating-point, Cartesian-with-y-up) coordinate
// system, matching the original _setwindow(TRUE, xmin, ymin, xmax,
// ymax) convention — the TRUE argument is what makes y increase
// upward, the natural convention for an engineering plot, rather
// than SVG's own default of y increasing downward. Canvas handles
// that flip internally so callers never need to think about pixel
// space at all.
//
// Only line, box (rectangle), circle, and text primitives are
// implemented, because a survey of all eight Tier 3 programs' source
// found no call site anywhere for the wider set their shared,
// unrecovered graphics library declares (wpaint, wset, warc): those
// three are declared but never actually used by any of the eight.
package svgplot

import (
	"fmt"
	"html"
	"strings"
)

// ega is the standard 16-color EGA/CGA palette, index 0-15, which is
// what every original program's own color arguments (plain integers)
// index into. The original's own custom() palette-remapping call has
// no recoverable implementation anywhere in the source tree (see the
// Tier 3 survey), so this is the standard table rather than a
// per-program remapping.
var ega = [16]string{
	"#000000", "#0000AA", "#00AA00", "#00AAAA",
	"#AA0000", "#AA00AA", "#AA5500", "#AAAAAA",
	"#555555", "#5555FF", "#55FF55", "#55FFFF",
	"#FF5555", "#FF55FF", "#FFFF55", "#FFFFFF",
}

// colorHex returns the hex color for an EGA palette index, clamping
// out-of-range indices to white rather than panicking or silently
// indexing garbage.
func colorHex(color int) string {
	if color < 0 || color > 15 {
		return "#FFFFFF"
	}
	return ega[color]
}

// Canvas accumulates drawing commands in window coordinates and
// renders them to SVG on demand.
type Canvas struct {
	xmin, ymin, xmax, ymax float64
	widthPx, heightPx      int
	elements               []string
}

// New creates a Canvas whose window spans [xmin,xmax] horizontally and
// [ymin,ymax] vertically (y increasing upward), rendered into an SVG
// image widthPx by heightPx pixels.
func New(xmin, ymin, xmax, ymax float64, widthPx, heightPx int) *Canvas {
	return &Canvas{xmin: xmin, ymin: ymin, xmax: xmax, ymax: ymax, widthPx: widthPx, heightPx: heightPx}
}

func (c *Canvas) px(x float64) float64 {
	return (x - c.xmin) / (c.xmax - c.xmin) * float64(c.widthPx)
}

func (c *Canvas) py(y float64) float64 {
	return float64(c.heightPx) - (y-c.ymin)/(c.ymax-c.ymin)*float64(c.heightPx)
}

// pxLen converts a window-coordinate length (not a point) to pixels,
// for radii and similar magnitudes where no coordinate flip applies.
func (c *Canvas) pxLen(length float64) float64 {
	return length / (c.xmax - c.xmin) * float64(c.widthPx)
}

// Line draws a solid line segment from (x1,y1) to (x2,y2) in the
// given EGA color, matching wline().
func (c *Canvas) Line(x1, y1, x2, y2 float64, color int) {
	c.elements = append(c.elements, fmt.Sprintf(
		`<line x1="%.3f" y1="%.3f" x2="%.3f" y2="%.3f" stroke="%s" stroke-width="1"/>`,
		c.px(x1), c.py(y1), c.px(x2), c.py(y2), colorHex(color)))
}

// DashedLine draws a dotted/dashed line segment, matching the
// original's own _setlinestyle(0xAAAA) convention used for axes and
// gridlines (as opposed to 0xFFFF, solid, used for curves and frames).
func (c *Canvas) DashedLine(x1, y1, x2, y2 float64, color int) {
	c.elements = append(c.elements, fmt.Sprintf(
		`<line x1="%.3f" y1="%.3f" x2="%.3f" y2="%.3f" stroke="%s" stroke-width="1" stroke-dasharray="3,3"/>`,
		c.px(x1), c.py(y1), c.px(x2), c.py(y2), colorHex(color)))
}

// Box draws an axis-aligned rectangle between (x1,y1) and (x2,y2),
// matching wbox(): outline only if filled is false, solid fill of the
// same color if filled is true.
func (c *Canvas) Box(x1, y1, x2, y2 float64, color int, filled bool) {
	left, right := c.px(x1), c.px(x2)
	if left > right {
		left, right = right, left
	}
	top, bottom := c.py(y1), c.py(y2)
	if top > bottom {
		top, bottom = bottom, top
	}
	fill := "none"
	if filled {
		fill = colorHex(color)
	}
	c.elements = append(c.elements, fmt.Sprintf(
		`<rect x="%.3f" y="%.3f" width="%.3f" height="%.3f" stroke="%s" fill="%s"/>`,
		left, top, right-left, bottom-top, colorHex(color), fill))
}

// Circle draws a circle of radius r centered at (x,y), matching
// wcircle(): outline only if filled is false, solid fill of the same
// color if filled is true.
func (c *Canvas) Circle(x, y, r float64, color int, filled bool) {
	fill := "none"
	if filled {
		fill = colorHex(color)
	}
	c.elements = append(c.elements, fmt.Sprintf(
		`<circle cx="%.3f" cy="%.3f" r="%.3f" stroke="%s" fill="%s"/>`,
		c.px(x), c.py(y), c.pxLen(r), colorHex(color), fill))
}

// Text places a text label with its baseline-left corner at (x,y),
// matching pwin()/gtext(): the original's own off-screen-buffer
// double-buffering trick (needed only because DOS text output
// couldn't be positioned in floating-point window coordinates
// directly) collapses to a single <text> element here.
func (c *Canvas) Text(x, y float64, s string, color int) {
	c.elements = append(c.elements, fmt.Sprintf(
		`<text x="%.3f" y="%.3f" fill="%s" font-family="monospace" font-size="12">%s</text>`,
		c.px(x), c.py(y), colorHex(color), html.EscapeString(s)))
}

// String renders the accumulated drawing commands as a complete,
// standalone SVG document.
func (c *Canvas) String() string {
	var b strings.Builder
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 %d %d" width="%d" height="%d">`+"\n",
		c.widthPx, c.heightPx, c.widthPx, c.heightPx)
	b.WriteString(fmt.Sprintf(`<rect x="0" y="0" width="%d" height="%d" fill="#000000"/>`+"\n", c.widthPx, c.heightPx))
	for _, el := range c.elements {
		b.WriteString(el)
		b.WriteByte('\n')
	}
	b.WriteString("</svg>\n")
	return b.String()
}
