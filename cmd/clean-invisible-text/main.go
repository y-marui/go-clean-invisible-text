// Command clean-invisible-text detects, explains, and safely cleans
// dangerous invisible Unicode characters in UTF-8 plain text.
package main

import (
	"os"

	"github.com/y-marui/go-clean-invisible-text/internal/cli"
)

func main() {
	os.Exit(cli.Run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
}
