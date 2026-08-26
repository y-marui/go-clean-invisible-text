package cli

import (
	"flag"
	"fmt"
	"io"
)

// runCheck implements the check command: report findings without modifying
// input.
func runCheck(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("check", flag.ContinueOnError)
	fs.SetOutput(stderr)
	jsonOut := fs.Bool("json", false, "emit a machine-readable JSON report to stdout")
	if err := fs.Parse(args); err != nil {
		return exitError
	}
	files := fs.Args()
	if len(files) == 0 {
		fmt.Fprintln(stderr, "check: at least one file is required")
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
				renderCheck(stderr, r)
			}
			exit = worstExit(exit, exitError)
			continue
		}
		r := buildReport(path, data, result.Findings, false)
		reports = append(reports, r)
		if !*jsonOut {
			renderCheck(stderr, r)
		}
		if len(result.Findings) > 0 {
			exit = worstExit(exit, exitFindings)
		}
	}

	if *jsonOut {
		if err := writeJSON(stdout, reports); err != nil {
			fmt.Fprintf(stderr, "check: %v\n", err)
			return exitError
		}
	}
	return exit
}
