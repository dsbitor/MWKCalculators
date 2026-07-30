package legacydat

import (
	"strings"
	"testing"
)

func TestRows_StartAndEndMarkers(t *testing.T) {
	data := `Data for DRILL

Entries are:   drill name		 drill diameter (in)

anything above the line below is ignored
STARTOFDATA

1/64		.0156
1/32		.0313

;a comment line, ignored
3/64		.0469
ENDOFDATA

You can add notes here.
`
	rows, err := Rows(strings.NewReader(data), Options{StartMarker: "STARTOFDATA", EndMarker: "ENDOFDATA"})
	if err != nil {
		t.Fatalf("Rows() error = %v", err)
	}
	want := []string{"1/64\t\t.0156", "1/32\t\t.0313", "3/64\t\t.0469"}
	if len(rows) != len(want) {
		t.Fatalf("Rows() = %q, want %d rows", rows, len(want))
	}
	for i := range want {
		if rows[i] != want[i] {
			t.Errorf("Rows()[%d] = %q, want %q", i, rows[i], want[i])
		}
	}
}

func TestRows_NoStartMarker_StartsAtFirstLine(t *testing.T) {
	data := ";Data for GAGE (wire and sheet gage utility)\n\n000000,0.5800\n00000,0.5165\nENDOFDATA\n\ntrailer text\n"
	rows, err := Rows(strings.NewReader(data), Options{EndMarker: "ENDOFDATA"})
	if err != nil {
		t.Fatalf("Rows() error = %v", err)
	}
	want := []string{"000000,0.5800", "00000,0.5165"}
	if len(rows) != len(want) || rows[0] != want[0] || rows[1] != want[1] {
		t.Errorf("Rows() = %q, want %q", rows, want)
	}
}

func TestRows_NoEndMarker_ReadsToEOF(t *testing.T) {
	data := "STARTOFDATA\n1\n2\n3\n"
	rows, err := Rows(strings.NewReader(data), Options{StartMarker: "STARTOFDATA"})
	if err != nil {
		t.Fatalf("Rows() error = %v", err)
	}
	if len(rows) != 3 {
		t.Errorf("Rows() = %q, want 3 rows", rows)
	}
}

func TestRows_MissingStartMarker_ReturnsError(t *testing.T) {
	data := "no marker anywhere in this file\njust some lines\n"
	_, err := Rows(strings.NewReader(data), Options{StartMarker: "STARTOFDATA", EndMarker: "ENDOFDATA"})
	if err == nil {
		t.Fatal("Rows() error = nil, want an error when the start marker is never found")
	}
}

func TestRows_EmptyDataRegion_IsNotAnError(t *testing.T) {
	data := "STARTOFDATA\nENDOFDATA\n"
	rows, err := Rows(strings.NewReader(data), Options{StartMarker: "STARTOFDATA", EndMarker: "ENDOFDATA"})
	if err != nil {
		t.Fatalf("Rows() error = %v, want no error for a marker pair with nothing between them", err)
	}
	if len(rows) != 0 {
		t.Errorf("Rows() = %q, want zero rows", rows)
	}
}

func TestRows_InlineCommentIsNotStripped(t *testing.T) {
	// Rows only strips whole comment lines (leading ';'); an inline
	// trailing comment is left for the caller's own Fields split to
	// discard, since some callers include ';' in their cutset and
	// some don't.
	data := "STARTOFDATA\n40   ; DH worm gear ratio\nENDOFDATA\n"
	rows, err := Rows(strings.NewReader(data), Options{StartMarker: "STARTOFDATA", EndMarker: "ENDOFDATA"})
	if err != nil {
		t.Fatalf("Rows() error = %v", err)
	}
	if len(rows) != 1 || !strings.Contains(rows[0], "DH worm gear ratio") {
		t.Errorf("Rows() = %q, want the inline comment preserved on the one row", rows)
	}
}

func TestFields_SplitsOnCutsetAndDropsEmpties(t *testing.T) {
	got := Fields("40   ; DH worm gear ratio", ",\t;")
	want := []string{"40", " DH worm gear ratio"}
	if len(got) != 2 || strings.TrimSpace(got[0]) != "40" || got[1] != want[1] {
		t.Errorf("Fields() = %q, want first field \"40\"", got)
	}
}

func TestFields_TabDelimitedTwoColumnRow(t *testing.T) {
	got := Fields("1/64\t\t.0156", "\t")
	want := []string{"1/64", ".0156"}
	if len(got) != 2 || got[0] != want[0] || got[1] != want[1] {
		t.Errorf("Fields() = %q, want %q", got, want)
	}
}

func TestFields_CommaDelimited(t *testing.T) {
	got := Fields("000000,0.5800", ",")
	want := []string{"000000", "0.5800"}
	if len(got) != 2 || got[0] != want[0] || got[1] != want[1] {
		t.Errorf("Fields() = %q, want %q", got, want)
	}
}

func TestFields_NoDelimiterPresent_ReturnsWholeLineAsOneField(t *testing.T) {
	got := Fields("100", ",\t;")
	if len(got) != 1 || got[0] != "100" {
		t.Errorf("Fields() = %q, want [\"100\"]", got)
	}
}
