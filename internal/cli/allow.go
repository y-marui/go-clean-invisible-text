package cli

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/y-marui/go-clean-invisible-text/internal/allowlist"
)

const allowFlagUsage = `add an allow-list rule for one Warn finding (repeatable): ` +
	`codepoint=U+XXXX[,U+YYYY];reason=...;paths=glob[,glob];max_run=N|unlimited ` +
	`(reason is required; paths and max_run are optional, see docs/cli.md)`

const allowFileFlagUsage = `path to an allow-list config file (default: load ` +
	allowlist.DefaultFileName + ` from the current directory if present)`

// allowFlags holds the --allow/--allow-file values registered on one
// command's FlagSet, before they're parsed into allowlist.Rules.
type allowFlags struct {
	rules []string
	file  string
}

// stringSliceFlag implements flag.Value for a repeatable string flag.
type stringSliceFlag struct{ values *[]string }

func (f stringSliceFlag) String() string { return "" }

func (f stringSliceFlag) Set(v string) error {
	*f.values = append(*f.values, v)
	return nil
}

// registerAllowFlags adds --allow and --allow-file to fs.
func registerAllowFlags(fs *flag.FlagSet, af *allowFlags) {
	fs.Var(stringSliceFlag{&af.rules}, "allow", allowFlagUsage)
	fs.StringVar(&af.file, "allow-file", "", allowFileFlagUsage)
}

// loadRules parses every --allow flag and the config file (the explicit
// --allow-file, or allowlist.DefaultFileName if it exists and no
// --allow-file was given) into one combined rule set. On a parse error it
// writes a message to stderr and returns ok=false.
func loadRules(af allowFlags, stderr io.Writer) (rules []allowlist.Rule, ok bool) {
	filePath := af.file
	if filePath == "" {
		if _, err := os.Stat(allowlist.DefaultFileName); err == nil {
			filePath = allowlist.DefaultFileName
		}
	}
	if filePath != "" {
		fileRules, err := allowlist.LoadFile(filePath)
		if err != nil {
			fmt.Fprintln(stderr, err)
			return nil, false
		}
		rules = append(rules, fileRules...)
	}
	for _, s := range af.rules {
		r, err := allowlist.ParseFlag(s)
		if err != nil {
			fmt.Fprintln(stderr, err)
			return nil, false
		}
		rules = append(rules, r)
	}
	return rules, true
}
