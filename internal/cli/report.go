// Package cli implements the clean-invisible-text command-line interface:
// check, fix, explain, and clean, per docs/cli.md and docs/specification.md.
package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"unicode/utf8"

	"github.com/y-marui/go-clean-invisible-text/internal/cleaner"
	"github.com/y-marui/go-clean-invisible-text/internal/mutate"
)

const (
	exitOK       = 0
	exitFindings = 1
	exitError    = 2
)

// worstExit returns whichever of a, b represents the more severe outcome
// (2 > 1 > 0).
func worstExit(a, b int) int {
	if b > a {
		return b
	}
	return a
}

// findingJSON is the JSON representation of one cleaner.Finding, per the
// --json contract documented in docs/cli.md.
type findingJSON struct {
	Line        int                `json:"line"`
	Column      int                `json:"column"`
	Offset      int                `json:"offset"`
	Rune        string             `json:"rune"`
	Name        string             `json:"name"`
	Category    cleaner.Category   `json:"category"`
	Action      cleaner.ActionKind `json:"action"`
	Replacement string             `json:"replacement"`
}

// fileReport is the JSON representation of one file's outcome, per the
// --json contract documented in docs/cli.md. The shape is the same for
// check, fix, and explain so consumers don't need to special-case the
// command that produced it.
type fileReport struct {
	Path     string        `json:"path"`
	Findings []findingJSON `json:"findings"`
	Changed  bool          `json:"changed"`
	Error    *string       `json:"error"`
}

// locate returns the 1-indexed line and column (in runes, not bytes) of the
// byte offset into data. Column counts runes since the start of the line.
func locate(data []byte, offset int) (line, col int) {
	line = 1
	lineStart := 0
	for i := 0; i < offset && i < len(data); i++ {
		if data[i] == '\n' {
			line++
			lineStart = i + 1
		}
	}
	col = utf8.RuneCount(data[lineStart:offset]) + 1
	return line, col
}

// buildReport converts cleaner.Findings against the original file content
// into a fileReport.
func buildReport(path string, data []byte, findings []cleaner.Finding, changed bool) fileReport {
	fjs := make([]findingJSON, 0, len(findings))
	for _, f := range findings {
		line, col := locate(data, f.Offset)
		fjs = append(fjs, findingJSON{
			Line:        line,
			Column:      col,
			Offset:      f.Offset,
			Rune:        fmt.Sprintf("U+%04X", f.Rune),
			Name:        f.Name,
			Category:    f.Category,
			Action:      f.Action,
			Replacement: f.Replacement,
		})
	}
	return fileReport{Path: path, Findings: fjs, Changed: changed}
}

// errorReport builds a fileReport for a file that could not be processed.
func errorReport(path string, err error) fileReport {
	msg := err.Error()
	return fileReport{Path: path, Findings: []findingJSON{}, Error: &msg}
}

// writeJSON writes reports as a single JSON array to w.
func writeJSON(w io.Writer, reports []fileReport) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(reports)
}

// renderCheck writes a terse per-file summary to w: one line naming the
// finding count when there are any, nothing when the file is clean.
func renderCheck(w io.Writer, r fileReport) {
	if r.Error != nil {
		fmt.Fprintf(w, "%s: %s\n", r.Path, *r.Error)
		return
	}
	if len(r.Findings) == 0 {
		return
	}
	fmt.Fprintf(w, "%s: %d finding(s)\n", r.Path, len(r.Findings))
}

// renderDetailed writes one line per Finding to w, in the form
// "<path>:<line>:<col>: U+<hex> <NAME> [<category>] -> <action>". Used by
// explain and fix, whose whole point is showing every finding in full.
func renderDetailed(w io.Writer, r fileReport) {
	if r.Error != nil {
		fmt.Fprintf(w, "%s: %s\n", r.Path, *r.Error)
		return
	}
	for _, f := range r.Findings {
		action := string(f.Action)
		if f.Action == cleaner.ActionReplace {
			action = fmt.Sprintf("replace with %q", f.Replacement)
		}
		fmt.Fprintf(w, "%s:%d:%d: %s %s [%s] -> %s\n",
			r.Path, f.Line, f.Column, f.Rune, f.Name, f.Category, action)
	}
}

// analyzeFile reads path and runs it through cleaner.Clean without writing
// anything back, for the read-only check and explain commands. It applies
// the same binary-content rejection as mutate.File, but skips the
// symlink-mutation guard entirely: a read-only command has nothing to
// protect by refusing to follow a symlink.
func analyzeFile(path string) ([]byte, cleaner.Result, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, cleaner.Result{}, err
	}
	if mutate.IsBinary(data) {
		return nil, cleaner.Result{}, fmt.Errorf("%w: %s", mutate.ErrBinary, path)
	}
	result, err := cleaner.Clean(data, cleaner.Options{})
	if err != nil {
		return nil, cleaner.Result{}, fmt.Errorf("%s: %w", path, err)
	}
	return data, result, nil
}
