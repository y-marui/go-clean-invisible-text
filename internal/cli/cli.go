package cli

import (
	"fmt"
	"io"
)

const usage = `Usage:
  clean-invisible-text check [--json] [--allow RULE]... [--allow-file FILE] FILE...
  clean-invisible-text fix [--json] [--keep-warnings] [--allow RULE]... [--allow-file FILE] FILE...
  clean-invisible-text explain [--json] [--allow RULE]... [--allow-file FILE] FILE...
  clean-invisible-text clean [--keep-warnings] [--allow RULE]... [--allow-file FILE]

check    report findings without modifying input
fix      modify named files and report every change
explain  show code point, Unicode name, location, category, and planned action
clean    read standard input and write cleaned text to standard output
`

// Run dispatches to the requested subcommand and returns the process exit
// code. It never calls os.Exit itself, so callers can invoke it directly in
// tests with in-memory streams.
func Run(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		fmt.Fprint(stderr, usage)
		return exitError
	}

	cmd, rest := args[0], args[1:]
	switch cmd {
	case "check":
		return runCheck(rest, stdout, stderr)
	case "fix":
		return runFix(rest, stdout, stderr)
	case "explain":
		return runExplain(rest, stdout, stderr)
	case "clean":
		return runClean(rest, stdin, stdout, stderr)
	case "-h", "--help", "help":
		fmt.Fprint(stdout, usage)
		return exitOK
	default:
		fmt.Fprintf(stderr, "clean-invisible-text: unknown command %q\n\n", cmd)
		fmt.Fprint(stderr, usage)
		return exitError
	}
}
