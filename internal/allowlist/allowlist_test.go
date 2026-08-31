package allowlist

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/y-marui/go-clean-invisible-text/internal/cleaner"
)

func TestParseFlag(t *testing.T) {
	cases := []struct {
		name    string
		input   string
		want    Rule
		wantErr bool
	}{
		{
			name:  "minimal",
			input: "codepoint=U+E000;reason=Nerd Font icon glyph",
			want:  Rule{Codepoints: []rune{0xE000}, Reason: "Nerd Font icon glyph"},
		},
		{
			name:  "multiple codepoints and paths and max_run",
			input: "codepoint=U+E000,U+E001;reason=icons;paths=*.md,internal/foo.txt;max_run=3",
			want: Rule{
				Codepoints: []rune{0xE000, 0xE001},
				Reason:     "icons",
				Paths:      []string{"*.md", "internal/foo.txt"},
				MaxRun:     3,
			},
		},
		{
			name:  "unlimited max_run",
			input: "codepoint=U+E000;reason=icons;max_run=unlimited",
			want:  Rule{Codepoints: []rune{0xE000}, Reason: "icons", MaxRun: UnlimitedRun},
		},
		{
			name:  "lowercase u prefix",
			input: "codepoint=u+e000;reason=icons",
			want:  Rule{Codepoints: []rune{0xE000}, Reason: "icons"},
		},
		{name: "missing reason", input: "codepoint=U+E000", wantErr: true},
		{name: "missing codepoint", input: "reason=icons", wantErr: true},
		{name: "empty reason", input: "codepoint=U+E000;reason=", wantErr: true},
		{name: "invalid codepoint", input: "codepoint=not-hex;reason=icons", wantErr: true},
		{name: "invalid max_run", input: "codepoint=U+E000;reason=icons;max_run=0", wantErr: true},
		{name: "invalid max_run word", input: "codepoint=U+E000;reason=icons;max_run=lots", wantErr: true},
		{name: "unknown field", input: "codepoint=U+E000;reason=icons;bogus=1", wantErr: true},
		{name: "field without equals", input: "codepoint", wantErr: true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ParseFlag(tc.input)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("ParseFlag(%q) = %+v, want error", tc.input, got)
				}
				return
			}
			if err != nil {
				t.Fatalf("ParseFlag(%q) returned error: %v", tc.input, err)
			}
			if !rulesEqual(got, tc.want) {
				t.Errorf("ParseFlag(%q) = %+v, want %+v", tc.input, got, tc.want)
			}
		})
	}
}

func rulesEqual(a, b Rule) bool {
	if a.Reason != b.Reason || a.MaxRun != b.MaxRun || len(a.Codepoints) != len(b.Codepoints) || len(a.Paths) != len(b.Paths) {
		return false
	}
	for i := range a.Codepoints {
		if a.Codepoints[i] != b.Codepoints[i] {
			return false
		}
	}
	for i := range a.Paths {
		if a.Paths[i] != b.Paths[i] {
			return false
		}
	}
	return true
}

func TestLoadFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "allow.json")
	content := `[
		{"codepoint": "U+E000", "reason": "Nerd Font icon glyph"},
		{"codepoint": "U+E001", "reason": "icons", "paths": ["*.md"], "max_run": 2}
	]`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	rules, err := LoadFile(path)
	if err != nil {
		t.Fatalf("LoadFile: %v", err)
	}
	if len(rules) != 2 {
		t.Fatalf("len(rules) = %d, want 2", len(rules))
	}
	if rules[0].Codepoints[0] != 0xE000 || rules[0].Reason != "Nerd Font icon glyph" {
		t.Errorf("rules[0] = %+v", rules[0])
	}
	if rules[1].MaxRun != 2 || len(rules[1].Paths) != 1 || rules[1].Paths[0] != "*.md" {
		t.Errorf("rules[1] = %+v", rules[1])
	}
}

func TestLoadFile_MissingReason(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "allow.json")
	if err := os.WriteFile(path, []byte(`[{"codepoint": "U+E000"}]`), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	if _, err := LoadFile(path); err == nil {
		t.Fatal("LoadFile: want error for missing reason")
	}
}

func TestLoadFile_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "allow.json")
	if err := os.WriteFile(path, []byte(`not json`), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	if _, err := LoadFile(path); err == nil {
		t.Fatal("LoadFile: want error for invalid JSON")
	}
}

func TestLoadFile_NotFound(t *testing.T) {
	if _, err := LoadFile(filepath.Join(t.TempDir(), "missing.json")); err == nil {
		t.Fatal("LoadFile: want error for missing file")
	}
}

func TestResolve_PathScoping(t *testing.T) {
	rules := []Rule{
		{Codepoints: []rune{0xE000}, Reason: "repo-wide"},
		{Codepoints: []rune{0xE001}, Reason: "md-only", Paths: []string{"*.md"}},
		{Codepoints: []rune{0xE002}, Reason: "dir-scoped", Paths: []string{"internal/notes.txt"}},
	}

	cases := []struct {
		name string
		path string
		want []rune
	}{
		{"repo-wide file", "anything.go", []rune{0xE000}},
		{"markdown by basename", "docs/readme.md", []rune{0xE000, 0xE001}},
		{"exact relative path", "internal/notes.txt", []rune{0xE000, 0xE002}},
		{"exact path elsewhere does not match", "other/notes.txt", []rune{0xE000}},
		{"empty path (clean command) only matches path-less rules", "", []rune{0xE000}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := Resolve(rules, tc.path)
			for _, cp := range tc.want {
				if _, ok := got[cp]; !ok {
					t.Errorf("Resolve(%q): missing rule for U+%04X, got %v", tc.path, cp, got)
				}
			}
			if len(got) != len(tc.want) {
				t.Errorf("Resolve(%q) = %v, want exactly %v", tc.path, got, tc.want)
			}
		})
	}
}

func TestResolve_UnionsOverlappingRules(t *testing.T) {
	rules := []Rule{
		{Codepoints: []rune{0xE000}, Reason: "first reason", MaxRun: 1, Paths: []string{"a.txt"}},
		{Codepoints: []rune{0xE000}, Reason: "second reason", MaxRun: 5},
	}

	got := Resolve(rules, "a.txt")
	rule, ok := got[0xE000]
	if !ok {
		t.Fatalf("Resolve: missing rule for U+E000: %v", got)
	}
	if rule.MaxRun != 5 {
		t.Errorf("MaxRun = %d, want 5 (loosest of matching rules)", rule.MaxRun)
	}
	if rule.Reason != "first reason; second reason" {
		t.Errorf("Reason = %q, want both reasons joined", rule.Reason)
	}

	// The path-scoped rule must not apply to a file it doesn't cover, but the
	// repo-wide rule still does, with only its own reason and MaxRun.
	got = Resolve(rules, "b.txt")
	rule, ok = got[0xE000]
	if !ok {
		t.Fatalf("Resolve: missing rule for U+E000 on b.txt: %v", got)
	}
	if rule.MaxRun != 5 || rule.Reason != "second reason" {
		t.Errorf("rule for b.txt = %+v, want {MaxRun:5 Reason:\"second reason\"}", rule)
	}
}

func TestResolve_UnlimitedWins(t *testing.T) {
	rules := []Rule{
		{Codepoints: []rune{0xE000}, Reason: "capped", MaxRun: 2},
		{Codepoints: []rune{0xE000}, Reason: "uncapped", MaxRun: UnlimitedRun},
	}
	got := Resolve(rules, "any.txt")
	if got[0xE000].MaxRun != UnlimitedRun {
		t.Errorf("MaxRun = %d, want UnlimitedRun once any matching rule is unlimited", got[0xE000].MaxRun)
	}
}

func TestResolve_DefaultsMaxRunToOne(t *testing.T) {
	rules := []Rule{{Codepoints: []rune{0xE000}, Reason: "icons"}}
	got := Resolve(rules, "any.txt")
	if got[0xE000].MaxRun != 1 {
		t.Errorf("MaxRun = %d, want 1 (default)", got[0xE000].MaxRun)
	}
}

// TestResolveFeedsIntoCleanerAllowRule is a light integration check that the
// map Resolve produces is exactly the shape internal/cleaner.Options expects.
func TestResolveFeedsIntoCleanerAllowRule(t *testing.T) {
	rules := []Rule{{Codepoints: []rune{0xE000}, Reason: "icons"}}
	var _ map[rune]cleaner.AllowRule = Resolve(rules, "any.txt")
}
