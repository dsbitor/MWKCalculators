// Package promptio implements the prompt-with-default input pattern
// used throughout the original MWK calculators: the vin function from
// the unavailable MWK.LIB. A prompt shows its default value. A blank
// line accepts the default. A line that cannot be parsed as the
// expected type also falls back to the default, after a message
// explaining why, matching vin's own forgiving behaviour rather than
// treating bad input as a fatal error.
package promptio

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// Prompter reads prompted values from in and writes prompts and
// fallback messages to out. It has no dependency on os.Stdin or
// os.Stdout, so a program under test can drive it with a fixed script.
type Prompter struct {
	in  *bufio.Reader
	out io.Writer
}

// New returns a Prompter that reads from in and writes to out.
// Neither argument may be nil.
func New(in io.Reader, out io.Writer) (*Prompter, error) {
	if in == nil {
		return nil, errors.New("promptio.New: in must not be nil")
	}
	if out == nil {
		return nil, errors.New("promptio.New: out must not be nil")
	}
	return &Prompter{in: bufio.NewReader(in), out: out}, nil
}

// Float prompts with label and the given default, shown in the
// prompt text, and returns the value entered. A blank line, an
// input stream that cannot be read further, or text that does not
// parse as a floating-point number all return def.
func (p *Prompter) Float(label string, def float64) float64 {
	fmt.Fprintf(p.out, "%s (default %g): ", label, def)

	line := p.readLine()
	if line == "" {
		return def
	}

	value, err := strconv.ParseFloat(line, 64)
	if err != nil {
		fmt.Fprintf(p.out, "%q is not a number; using default %g\n", line, def)
		return def
	}
	return value
}

// Int prompts with label and the given default, shown in the prompt
// text, and returns the value entered. A blank line, an input stream
// that cannot be read further, or text that does not parse as a
// whole number all return def.
func (p *Prompter) Int(label string, def int) int {
	fmt.Fprintf(p.out, "%s (default %d): ", label, def)

	line := p.readLine()
	if line == "" {
		return def
	}

	value, err := strconv.Atoi(line)
	if err != nil {
		fmt.Fprintf(p.out, "%q is not a whole number; using default %d\n", line, def)
		return def
	}
	return value
}

// Line prompts with label, shown as given with no default appended,
// and returns the trimmed line entered. It returns an empty string
// both for a genuinely blank line and for an input stream that ends
// or fails before a newline. Line is for prompts with no single
// natural default value, such as a choice among a few letters; Float
// and Int cover the common case of a numeric value with a default.
func (p *Prompter) Line(label string) string {
	fmt.Fprint(p.out, label)
	return p.readLine()
}

// readLine returns the next line with leading and trailing whitespace
// removed. It returns an empty string both for a genuinely blank line
// and for an input stream that ends or fails before a newline, since
// both cases mean the same thing to the caller: no usable input was
// supplied, so the default value applies.
func (p *Prompter) readLine() string {
	line, _ := p.in.ReadString('\n')
	return strings.TrimSpace(line)
}
