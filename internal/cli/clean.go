package cli

import (
	"flag"
	"fmt"
	"io"

	"github.com/y-marui/go-clean-invisible-text/internal/cleaner"
)

// runClean implements the clean command: read standard input and write
// cleaned text to standard output. Unlike check/fix/explain, it never
// reports findings and never returns exitFindings — per the spec, a
// successful stream cleaning always exits 0, since the cleaned output is
// already the result there is nothing further to review.
func runClean(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("clean", flag.ContinueOnError)
	fs.SetOutput(stderr)
	keepWarnings := fs.Bool("keep-warnings", false, "preserve Warn-classified code points instead of removing them")
	if err := fs.Parse(args); err != nil {
		return exitError
	}
	if fs.NArg() > 0 {
		fmt.Fprintln(stderr, "clean: reads standard input and takes no file arguments")
		return exitError
	}

	data, err := io.ReadAll(stdin)
	if err != nil {
		fmt.Fprintf(stderr, "clean: %v\n", err)
		return exitError
	}

	result, err := cleaner.Clean(data, cleaner.Options{KeepWarnings: *keepWarnings})
	if err != nil {
		fmt.Fprintf(stderr, "clean: %v\n", err)
		return exitError
	}

	if _, err := stdout.Write(result.Cleaned); err != nil {
		fmt.Fprintf(stderr, "clean: %v\n", err)
		return exitError
	}
	return exitOK
}
