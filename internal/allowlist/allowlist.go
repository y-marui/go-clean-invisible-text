// Package allowlist implements the per-code-point/file allow-list exception
// mechanism described by the "Warn: allow-list exceptions" section of
// docs/character-policy.md (Issue #26): repeatable CLI --allow flags and a
// JSON config file each contribute Rules, which Resolve reduces to the
// per-code-point map internal/cleaner.Options.AllowRules expects for one
// file.
//
// A Rule can only ever grant an exception to a Warn-classified code point —
// nothing here can suppress a Block finding (bidi controls, tag characters,
// noncharacters, unsafe controls, NBSP, BOM, ZWSP, word joiner, soft hyphen,
// or a non-contextual ZWJ/ZWNJ/variation selector). That boundary is
// enforced by internal/cleaner, not by this package, but it is the reason
// this mechanism is safe to load from a project-committed config file: at
// worst a malicious rule mislabels one already-flagged Warn occurrence, with
// its reason kept visible in output, and it can never hide a Trojan-Source
// or steganography-class attack.
package allowlist

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/y-marui/go-clean-invisible-text/internal/cleaner"
)

// UnlimitedRun disables the run-length guard for a Rule. Mirrors
// cleaner.UnlimitedRun; re-exported here so callers parsing rules don't need
// to import internal/cleaner just for this constant.
const UnlimitedRun = cleaner.UnlimitedRun

// DefaultFileName is the config file loaded automatically from the current
// directory when no --allow-file is given, if it exists.
const DefaultFileName = ".clean-invisible-text-allow.json"

// Rule is one allow-list exception, from a CLI --allow flag or one entry of
// a config file's JSON array.
type Rule struct {
	// Codepoints lists every code point this rule covers. A Rule with no
	// Codepoints is never produced by ParseFlag or LoadFile.
	Codepoints []rune
	// Reason is required and non-empty; it is copied onto every Finding
	// this rule allows.
	Reason string
	// Paths, if non-empty, scopes the rule to files whose path matches one
	// of these patterns (see matchesPath). Empty means every file in the
	// invocation.
	Paths []string
	// MaxRun is 0 (meaning the default of 1: a single isolated occurrence),
	// UnlimitedRun, or a positive explicit cap.
	MaxRun int
}

// configRule is the JSON shape of one Rule in a config file.
type configRule struct {
	Codepoint string   `json:"codepoint"`
	Reason    string   `json:"reason"`
	Paths     []string `json:"paths,omitempty"`
	MaxRun    int      `json:"max_run,omitempty"`
}

// ParseFlag parses one --allow flag value: semicolon-separated key=value
// fields, e.g. `codepoint=U+E000;reason=Nerd Font icon;paths=*.md;max_run=3`.
// codepoint and paths accept a comma-separated list. max_run accepts a
// positive integer or the literal "unlimited".
func ParseFlag(s string) (Rule, error) {
	var r Rule
	for _, field := range strings.Split(s, ";") {
		field = strings.TrimSpace(field)
		if field == "" {
			continue
		}
		key, value, ok := strings.Cut(field, "=")
		if !ok {
			return Rule{}, fmt.Errorf("allowlist: invalid --allow field %q (want key=value)", field)
		}
		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)
		switch key {
		case "codepoint":
			cps, err := parseCodepoints(value)
			if err != nil {
				return Rule{}, err
			}
			r.Codepoints = cps
		case "reason":
			r.Reason = value
		case "paths":
			r.Paths = splitNonEmpty(value)
		case "max_run":
			n, err := parseMaxRun(value)
			if err != nil {
				return Rule{}, err
			}
			r.MaxRun = n
		default:
			return Rule{}, fmt.Errorf("allowlist: unknown --allow field %q", key)
		}
	}
	if err := validate(r); err != nil {
		return Rule{}, err
	}
	return r, nil
}

// LoadFile reads and parses a config file: a JSON array of rule objects with
// the same fields as ParseFlag (codepoint, reason, paths, max_run).
func LoadFile(path string) ([]Rule, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("allowlist: %w", err)
	}
	var raw []configRule
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&raw); err != nil {
		return nil, fmt.Errorf("allowlist: %s: %w", path, err)
	}
	rules := make([]Rule, 0, len(raw))
	for i, cr := range raw {
		cps, err := parseCodepoints(cr.Codepoint)
		if err != nil {
			return nil, fmt.Errorf("allowlist: %s: rule %d: %w", path, i, err)
		}
		if cr.MaxRun != 0 && cr.MaxRun < 1 && cr.MaxRun != UnlimitedRun {
			return nil, fmt.Errorf("allowlist: %s: rule %d: invalid max_run %d (want a positive integer or %d for unlimited)", path, i, cr.MaxRun, UnlimitedRun)
		}
		r := Rule{Codepoints: cps, Reason: cr.Reason, Paths: cr.Paths, MaxRun: cr.MaxRun}
		if err := validate(r); err != nil {
			return nil, fmt.Errorf("allowlist: %s: rule %d: %w", path, i, err)
		}
		rules = append(rules, r)
	}
	return rules, nil
}

// validate checks the fields every Rule must satisfy regardless of source.
func validate(r Rule) error {
	if len(r.Codepoints) == 0 {
		return fmt.Errorf("allowlist: rule requires at least one codepoint")
	}
	if strings.TrimSpace(r.Reason) == "" {
		return fmt.Errorf("allowlist: rule requires a non-empty reason")
	}
	return nil
}

func splitNonEmpty(value string) []string {
	var out []string
	for _, part := range strings.Split(value, ",") {
		part = strings.TrimSpace(part)
		if part != "" {
			out = append(out, part)
		}
	}
	return out
}

func parseCodepoints(value string) ([]rune, error) {
	var out []rune
	for _, part := range splitNonEmpty(value) {
		cp, err := parseCodepoint(part)
		if err != nil {
			return nil, err
		}
		out = append(out, cp)
	}
	return out, nil
}

func parseCodepoint(s string) (rune, error) {
	trimmed := strings.TrimPrefix(strings.TrimPrefix(s, "U+"), "u+")
	v, err := strconv.ParseInt(trimmed, 16, 64)
	if err != nil || v < 0 || v > utf8.MaxRune {
		return 0, fmt.Errorf("allowlist: invalid codepoint %q (want U+XXXX)", s)
	}
	return rune(v), nil
}

func parseMaxRun(value string) (int, error) {
	if strings.EqualFold(value, "unlimited") {
		return UnlimitedRun, nil
	}
	n, err := strconv.Atoi(value)
	if err != nil || n < 1 {
		return 0, fmt.Errorf("allowlist: invalid max_run %q (want a positive integer or \"unlimited\")", value)
	}
	return n, nil
}

// matchesPath reports whether path is covered by one of patterns: a pattern
// containing a path separator is matched against path as given (relative to
// the invocation); a pattern with no separator is also matched against just
// path's base name, so `*.md` matches a file in any directory. Neither form
// crosses directory boundaries the way a `**` glob would — patterns without
// a literal directory component only ever match by base name.
//
// path == "" means the caller (the clean command) has no file path at all;
// a nonempty patterns can never match it — without this, filepath.Base("")
// returns ".", which a wildcard pattern like paths=* would otherwise match,
// contradicting the documented "clean has no file path, so a paths-scoped
// rule never applies to it" behavior.
func matchesPath(patterns []string, path string) bool {
	if len(patterns) == 0 {
		return true
	}
	if path == "" {
		return false
	}
	base := filepath.Base(path)
	for _, p := range patterns {
		if ok, _ := filepath.Match(p, path); ok {
			return true
		}
		if !strings.ContainsAny(p, `/\`) {
			if ok, _ := filepath.Match(p, base); ok {
				return true
			}
		}
	}
	return false
}

// Resolve reduces rules to the map cleaner.Options.AllowRules expects for
// one file at path (path is "" for the clean command, which has no file
// path — only path-less rules can ever match it). Every rule matching path
// is unioned per code point: the effective MaxRun is the loosest (largest,
// or UnlimitedRun) among matching rules for that code point, and Reason
// joins every distinct matching reason so all of them stay auditable.
func Resolve(rules []Rule, path string) map[rune]cleaner.AllowRule {
	type acc struct {
		maxRun    int
		reasons   []string
		reasonSet map[string]bool
	}
	merged := map[rune]*acc{}
	for _, r := range rules {
		if !matchesPath(r.Paths, path) {
			continue
		}
		effMaxRun := r.MaxRun
		if effMaxRun == 0 {
			effMaxRun = 1
		}
		for _, cp := range r.Codepoints {
			a, ok := merged[cp]
			if !ok {
				a = &acc{maxRun: effMaxRun, reasonSet: map[string]bool{}}
				merged[cp] = a
			} else if a.maxRun != UnlimitedRun {
				if effMaxRun == UnlimitedRun || effMaxRun > a.maxRun {
					a.maxRun = effMaxRun
				}
			}
			if !a.reasonSet[r.Reason] {
				a.reasonSet[r.Reason] = true
				a.reasons = append(a.reasons, r.Reason)
			}
		}
	}
	out := make(map[rune]cleaner.AllowRule, len(merged))
	for cp, a := range merged {
		out[cp] = cleaner.AllowRule{MaxRun: a.maxRun, Reason: strings.Join(a.reasons, "; ")}
	}
	return out
}
