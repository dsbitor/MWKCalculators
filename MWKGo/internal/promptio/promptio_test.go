package promptio

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

// errReader always fails on Read, simulating a closed or broken
// input stream (the forced-failure case for this package).
type errReader struct{}

func (errReader) Read([]byte) (int, error) {
	return 0, errors.New("simulated input stream failure")
}

func TestNew_NilReaderOrWriter_ReturnsError(t *testing.T) {
	tests := []struct {
		name string
		in   *strings.Reader
		out  *bytes.Buffer
	}{
		{name: "nil reader", in: nil, out: &bytes.Buffer{}},
		{name: "nil writer", in: strings.NewReader(""), out: nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var (
				p   *Prompter
				err error
			)
			switch {
			case tt.in == nil:
				p, err = New(nil, tt.out)
			case tt.out == nil:
				p, err = New(tt.in, nil)
			}
			if err == nil {
				t.Fatalf("New() error = nil, want an error")
			}
			if p != nil {
				t.Fatalf("New() Prompter = %v, want nil on error", p)
			}
		})
	}
}

func TestPrompterFloat(t *testing.T) {
	tests := []struct {
		name  string
		input string
		def   float64
		want  float64
	}{
		{name: "valid input overrides default", input: "3.5\n", def: 1.0, want: 3.5},
		{name: "blank line keeps default", input: "\n", def: 5.0, want: 5.0},
		{name: "negative value parses correctly", input: "-2.25\n", def: 0, want: -2.25},
		{name: "zero value parses correctly", input: "0\n", def: 9.0, want: 0},
		{name: "non-numeric text falls back to default", input: "abc\n", def: 2.0, want: 2.0},
		{name: "end of input with no newline keeps default", input: "", def: 4.0, want: 4.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := &bytes.Buffer{}
			p, err := New(strings.NewReader(tt.input), out)
			if err != nil {
				t.Fatalf("New() error = %v", err)
			}

			got := p.Float("value", tt.def)
			if got != tt.want {
				t.Errorf("Float() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPrompterFloat_BrokenInputStream_ReturnsDefault(t *testing.T) {
	out := &bytes.Buffer{}
	p, err := New(errReader{}, out)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	got := p.Float("value", 7.5)
	if got != 7.5 {
		t.Errorf("Float() = %v, want default 7.5 when the input stream fails", got)
	}
}

func TestPrompterInt(t *testing.T) {
	tests := []struct {
		name  string
		input string
		def   int
		want  int
	}{
		{name: "valid input overrides default", input: "12\n", def: 1, want: 12},
		{name: "blank line keeps default", input: "\n", def: 5, want: 5},
		{name: "negative value parses correctly", input: "-3\n", def: 0, want: -3},
		{name: "zero value parses correctly", input: "0\n", def: 9, want: 0},
		{name: "non-numeric text falls back to default", input: "abc\n", def: 2, want: 2},
		{name: "decimal text falls back to default", input: "3.5\n", def: 2, want: 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := &bytes.Buffer{}
			p, err := New(strings.NewReader(tt.input), out)
			if err != nil {
				t.Fatalf("New() error = %v", err)
			}

			got := p.Int("value", tt.def)
			if got != tt.want {
				t.Errorf("Int() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPrompter_SequentialPrompts_ReadInOrder(t *testing.T) {
	out := &bytes.Buffer{}
	p, err := New(strings.NewReader("10\n\n"), out)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	first := p.Float("first", 1.0)
	second := p.Float("second", 2.0)

	if first != 10 {
		t.Errorf("first Float() = %v, want 10", first)
	}
	if second != 2.0 {
		t.Errorf("second Float() = %v, want default 2.0 for the blank second line", second)
	}
}

func TestPrompterFloat_InvalidInput_WritesFallbackMessage(t *testing.T) {
	out := &bytes.Buffer{}
	p, err := New(strings.NewReader("not-a-number\n"), out)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	p.Float("value", 2.0)

	if !strings.Contains(out.String(), "not-a-number") {
		t.Errorf("output %q does not explain why the default was used", out.String())
	}
}

func TestPrompterLine(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{name: "normal input is trimmed", input: "  c  \n", want: "c"},
		{name: "blank line returns empty string", input: "\n", want: ""},
		{name: "end of input with no newline returns empty string", input: "", want: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := &bytes.Buffer{}
			p, err := New(strings.NewReader(tt.input), out)
			if err != nil {
				t.Fatalf("New() error = %v", err)
			}

			if got := p.Line("choice? "); got != tt.want {
				t.Errorf("Line() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestPrompterLine_BrokenInputStream_ReturnsEmptyString(t *testing.T) {
	out := &bytes.Buffer{}
	p, err := New(errReader{}, out)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	if got := p.Line("choice? "); got != "" {
		t.Errorf("Line() = %q, want empty string when the input stream fails", got)
	}
}

func TestPrompter_LineThenFloat_ShareOneUnderlyingReader(t *testing.T) {
	// Line and Float must read from the same buffered reader so that
	// a program combining both prompt styles (as dew does) never
	// strands input the buffer already read ahead past the first
	// prompt's newline.
	out := &bytes.Buffer{}
	p, err := New(strings.NewReader("c\n70\n"), out)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	choice := p.Line("scale? ")
	value := p.Float("value", 0)

	if choice != "c" {
		t.Errorf("Line() = %q, want %q", choice, "c")
	}
	if value != 70 {
		t.Errorf("Float() = %v, want 70", value)
	}
}
