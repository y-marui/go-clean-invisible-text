package cli

import (
	"flag"
	"fmt"
	"io"

	"github.com/y-marui/go-clean-invisible-text/internal/allowlist"
	"github.com/y-marui/go-clean-invisible-text/internal/mutate"
)

// runFix implements the fix command: modify named files and report every
// change.
func runFix(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("fix", flag.ContinueOnError)
	fs.SetOutput(stderr)
	jsonOut := fs.Bool("json", false, "emit a machine-readable JSON report to stdout")
	keepWarnings := fs.Bool("keep-warnings", false, "preserve Warn-classified code points instead of removing them")
	var af allowFlags
	registerAllowFlags(fs, &af)
	if err := fs.Parse(args); err != nil {
		return exitError
	}
	files := fs.Args()
	if len(files) == 0 {
		fmt.Fprintln(stderr, "fix: at least one file is required")
		return exitError
	}
	rules, ok := loadRules(af, stderr)
	if !ok {
		return exitError
	}

	exit := exitOK
	reports := make([]fileReport, 0, len(files))
	for _, path := range files {
		res, err := mutate.File(path, mutate.Options{KeepWarnings: *keepWarnings, AllowRules: allowlist.Resolve(rules, path)})
		if err != nil {
			r := errorReport(path, err)
			reports = append(reports, r)
			if !*jsonOut {
				renderDetailed(stderr, r)
			}
			exit = worstExit(exit, exitError)
			continue
		}
		r := buildReport(path, res.Original, res.Findings, res.Changed)
		reports = append(reports, r)
		if !*jsonOut {
			renderDetailed(stderr, r)
		}
		if res.Changed {
			exit = worstExit(exit, exitFindings)
		}
	}

	if *jsonOut {
		if err := writeJSON(stdout, reports); err != nil {
			fmt.Fprintf(stderr, "fix: %v\n", err)
			return exitError
		}
	}
	return exit
}
