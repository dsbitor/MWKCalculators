// Package legacydat scans a legacy MWK.LIB-era ".DAT" reference file
// for its data region, the part shared by nearly every original
// program's own rdata() function: a run of content lines, usually
// bracketed by a "STARTOFDATA"/"ENDOFDATA" marker pair, with blank
// lines and ';'-prefixed comment lines already stripped out and any
// trailing inline comment left for the caller's own field split to
// discard.
//
// This package is deliberately narrow. Column layout, delimiter
// choice, and how many fields make up one logical row are different
// for every program's .DAT file, so that mapping stays in each
// program's own one-off conversion tool under MWKGo/tools; only the
// marker-scanning and comment/blank-line stripping is common enough
// across the batch to share.
package legacydat

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// Options controls how Rows finds the data region of a file.
type Options struct {
	// StartMarker is matched as a substring against each trimmed
	// line (the same way the original programs used strstr), and
	// every line up to and including the match is discarded. An
	// empty StartMarker means the data region starts at the first
	// line, matching the handful of original programs (GAGE.C, for
	// one) whose .DAT file has no start marker at all.
	StartMarker string

	// EndMarker is matched the same way as StartMarker; the line it
	// matches, and everything after, is discarded. An empty
	// EndMarker means the data region runs to the end of the file.
	EndMarker string
}

// Rows returns the content lines of r's data region, in file order,
// with leading and trailing whitespace trimmed, blank lines removed,
// and lines starting with ';' (the original convention for a
// comment line) removed. Each returned line still carries any
// trailing inline comment (for example "40 ; DH worm gear ratio");
// splitting that off, along with splitting a line into fields at
// all, is the caller's job, since the delimiter set is different for
// almost every program's data file.
//
// It is not an error for the data region to end up empty (matching
// the original: if the marker text has a typo, rdata() silently
// returns zero rows rather than failing), but it is an error for a
// non-empty StartMarker to never be found, since a throwaway import
// tool that silently produced zero rows would ship a much harder
// problem to notice than one that fails immediately.
func Rows(r io.Reader, opts Options) ([]string, error) {
	scanner := bufio.NewScanner(r)

	var rows []string
	active := opts.StartMarker == ""
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if !active {
			if strings.Contains(line, opts.StartMarker) {
				active = true
			}
			continue
		}
		if opts.EndMarker != "" && strings.Contains(line, opts.EndMarker) {
			break
		}
		if line == "" || strings.HasPrefix(line, ";") {
			continue
		}
		rows = append(rows, line)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scan legacy data: %w", err)
	}
	if !active {
		return nil, fmt.Errorf("start marker %q not found", opts.StartMarker)
	}
	return rows, nil
}

// Fields splits line on any character in cutset, discarding empty
// fields, matching the behaviour of the C standard library's
// strtok(line, cutset): a run of consecutive delimiters counts as one
// separator, and a delimiter at the very start or end of the line
// produces no leading or trailing empty field. Most of the original
// rdata() functions call strtok once or twice per line (typically
// splitting on some combination of ',', '\t', and ';', the last of
// which doubles as an inline-comment cutoff) to pull the row's fields
// off one at a time.
func Fields(line, cutset string) []string {
	return strings.FieldsFunc(line, func(r rune) bool {
		return strings.ContainsRune(cutset, r)
	})
}
