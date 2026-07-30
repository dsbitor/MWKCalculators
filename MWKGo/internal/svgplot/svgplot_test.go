package svgplot

import "testing"

func TestColorHex_ValidAndOutOfRangeIndices(t *testing.T) {
	if got := colorHex(0); got != "#000000" {
		t.Errorf("colorHex(0) = %q, want black", got)
	}
	if got := colorHex(15); got != "#FFFFFF" {
		t.Errorf("colorHex(15) = %q, want white", got)
	}
	if got := colorHex(99); got != "#FFFFFF" {
		t.Errorf("colorHex(99) = %q, want white fallback for an out-of-range index", got)
	}
	if got := colorHex(-1); got != "#FFFFFF" {
		t.Errorf("colorHex(-1) = %q, want white fallback for a negative index", got)
	}
}

func TestCanvas_CenterOfWindowMapsToCenterOfImage(t *testing.T) {
	c := New(-10, -10, 10, 10, 200, 200)
	if got := c.px(0); got != 100 {
		t.Errorf("px(0) = %v, want 100 (center of a 200px-wide image)", got)
	}
	if got := c.py(0); got != 100 {
		t.Errorf("py(0) = %v, want 100 (center of a 200px-tall image)", got)
	}
}

// TestCanvas_YAxisIsFlipped confirms the Cartesian-y-up convention:
// a larger window-coordinate y must map to a *smaller* pixel y,
// matching the original _setwindow(TRUE,...) behavior.
func TestCanvas_YAxisIsFlipped(t *testing.T) {
	c := New(0, 0, 10, 10, 100, 100)
	topPx := c.py(10)
	bottomPx := c.py(0)
	if topPx >= bottomPx {
		t.Errorf("py(10)=%v should be less than py(0)=%v (y increases upward in window space, downward in pixel space)", topPx, bottomPx)
	}
}

func TestCanvas_StringProducesWellFormedSVG(t *testing.T) {
	c := New(0, 0, 10, 10, 100, 100)
	c.Line(0, 0, 10, 10, 15)
	c.DashedLine(0, 5, 10, 5, 7)
	c.Box(2, 2, 8, 8, 14, false)
	c.Circle(5, 5, 2, 12, true)
	c.Text(1, 1, "hello", 15)

	out := c.String()
	if out[:4] != "<svg" {
		t.Errorf("output does not start with <svg: %q", out[:20])
	}
	wantSubstrings := []string{
		"<line", "stroke-dasharray", "<rect", "<circle", "<text", "hello", "</svg>",
	}
	for _, s := range wantSubstrings {
		if !contains(out, s) {
			t.Errorf("output missing expected substring %q", s)
		}
	}
}

func TestCanvas_TextEscapesSpecialCharacters(t *testing.T) {
	c := New(0, 0, 10, 10, 100, 100)
	c.Text(0, 0, "a < b & c", 15)
	out := c.String()
	if contains(out, "a < b & c") {
		t.Error("text was not escaped: raw '<' or '&' found in output")
	}
	if !contains(out, "&lt;") || !contains(out, "&amp;") {
		t.Error("expected escaped '<' and '&' in output")
	}
}

func contains(s, substr string) bool {
	for i := 0; i+len(substr) <= len(s); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
