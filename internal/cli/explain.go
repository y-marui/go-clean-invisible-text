package cli

import (
	"flag"
	"fmt"
	"io"
)

// runExplain implements the explain command: show code point, Unicode name,
// location, category, and planned action for every finding.
func runExplain(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("explain", flag.ContinueOnError)
	fs.SetOutput(stderr)
	jsonOut := fs.Bool("json", false, "emit a machine-readable JSON report to stdout")
	if err := fs.Parse(args); err != nil {
		return exitError
	}
	files := fs.Args()
	if len(files) == 0 {
		fmt.Fprintln(stderr, "explain: at least one file is required")
		return exitError
	}

	exit := exitOK
	reports := make([]fileReport, 0, len(files))
	for _, path := range files {
		data, result, err := analyzeFile(path)
		if err != nil {
			r := errorReport(path, err)
			reports = append(reports, r)
			if !*jsonOut {
				renderDetailed(stderr, r)
			}
			exit = worstExit(exit, exitError)
			continue
		}
		r := buildReport(path, data, result.Findings, false)
		reports = append(reports, r)
		if !*jsonOut {
			renderDetailed(stderr, r)
		}
		if len(result.Findings) > 0 {
			exit = worstExit(exit, exitFindings)
		}
	}

	if *jsonOut {
		if err := writeJSON(stdout, reports); err != nil {
			fmt.Fprintf(stderr, "explain: %v\n", err)
			return exitError
		}
	}
	return exit
}
