package cleaner

import (
	"bytes"
	"testing"
)

func mustClean(t *testing.T, input string) Result {
	t.Helper()
	res, err := Clean([]byte(input))
	if err != nil {
		t.Fatalf("Clean(%q) returned error: %v", input, err)
	}
	return res
}

func TestClean_SingleRuneCategories(t *testing.T) {
	cases := []struct {
		name         string
		input        string
		wantCleaned  string
		wantCategory Category
		wantAction   ActionKind
	}{
		{"nbsp", "a\u00A0b", "a b", CategoryNBSP, ActionReplace},
		{"bom-leading", "\uFEFFabc", "abc", CategoryBOM, ActionRemove},
		{"bom-inline", "a\uFEFFb", "ab", CategoryBOM, ActionRemove},
		{"zwsp", "a\u200Bb", "ab", CategoryZWSP, ActionRemove},
		{"word-joiner", "a\u2060b", "ab", CategoryWordJoiner, ActionRemove},
		{"soft-hyphen", "a\u00ADb", "ab", CategorySoftHyphen, ActionRemove},
		{"lre", "a\u202Ab", "ab", CategoryBidiControl, ActionRemove},
		{"rle", "a\u202Bb", "ab", CategoryBidiControl, ActionRemove},
		{"pdf", "a\u202Cb", "ab", CategoryBidiControl, ActionRemove},
		{"lro", "a\u202Db", "ab", CategoryBidiControl, ActionRemove},
		{"rlo", "a\u202Eb", "ab", CategoryBidiControl, ActionRemove},
		{"lri", "a\u2066b", "ab", CategoryBidiControl, ActionRemove},
		{"rli", "a\u2067b", "ab", CategoryBidiControl, ActionRemove},
		{"fsi", "a\u2068b", "ab", CategoryBidiControl, ActionRemove},
		{"pdi", "a\u2069b", "ab", CategoryBidiControl, ActionRemove},
		{"tag-language-tag", "a\U000E0001b", "ab", CategoryTag, ActionRemove},
		{"tag-char", "a\U000E0041b", "ab", CategoryTag, ActionRemove},
		{"noncharacter-fdd0-block", "a\uFDD0b", "ab", CategoryNoncharacter, ActionRemove},
		{"noncharacter-plane0", "a\uFFFEb", "ab", CategoryNoncharacter, ActionRemove},
		{"noncharacter-plane1", "a\U0001FFFFb", "ab", CategoryNoncharacter, ActionRemove},
		{"control-null", "a\x00b", "ab", CategoryControl, ActionRemove},
		{"control-del", "a\x7Fb", "ab", CategoryControl, ActionRemove},
		{"control-c1", "a\u0085b", "ab", CategoryControl, ActionRemove},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			res := mustClean(t, tc.input)
			if string(res.Cleaned) != tc.wantCleaned {
				t.Errorf("Cleaned = %q, want %q", res.Cleaned, tc.wantCleaned)
			}
			if len(res.Findings) != 1 {
				t.Fatalf("Findings = %v, want exactly 1", res.Findings)
			}
			f := res.Findings[0]
			if f.Category != tc.wantCategory {
				t.Errorf("Category = %q, want %q", f.Category, tc.wantCategory)
			}
			if f.Action != tc.wantAction {
				t.Errorf("Action = %q, want %q", f.Action, tc.wantAction)
			}
		})
	}
}

func TestClean_Preserved(t *testing.T) {
	cases := []struct {
		name  string
		input string
	}{
		{"tab", "a\tb"},
		{"lf", "a\nb"},
		{"cr", "a\rb"},
		{"ordinary-space", "a b"},
		{"ideographic-space", "a\u3000b"},
		{"combining-acute", "e\u0301"},
		{"variation-selector", "\u2764\uFE0F"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			res := mustClean(t, tc.input)
			if string(res.Cleaned) != tc.input {
				t.Errorf("Cleaned = %q, want unchanged %q", res.Cleaned, tc.input)
			}
			if len(res.Findings) != 0 {
				t.Errorf("Findings = %v, want none", res.Findings)
			}
		})
	}
}

func TestClean_ZWJZWNJContext(t *testing.T) {
	cases := []struct {
		name        string
		input       string
		wantCleaned string
		wantRemoved int
	}{
		{"single-between-visible-preserved", "a\u200Db", "a\u200Db", 0},
		{"single-between-visible-zwnj-preserved", "a\u200Cb", "a\u200Cb", 0},
		{"start-of-text-removed", "\u200Dab", "ab", 1},
		{"end-of-text-removed", "ab\u200D", "ab", 1},
		{"adjacent-space-removed", "a \u200Db", "a b", 1},
		{"adjacent-newline-removed", "a\n\u200Db", "a\nb", 1},
		{"adjacent-control-removed", "a\x00\u200Db", "ab", 1},
		{"repeated-identical-collapsed-and-kept", "a\u200D\u200D\u200Db", "a\u200Db", 2},
		{"repeated-identical-bad-context-all-removed", "a \u200D\u200Db", "a b", 2},
		{"mixed-run-removed", "a\u200D\u200Cb", "ab", 2},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			res := mustClean(t, tc.input)
			if string(res.Cleaned) != tc.wantCleaned {
				t.Errorf("Cleaned = %q, want %q", res.Cleaned, tc.wantCleaned)
			}
			joinerFindings := 0
			for _, f := range res.Findings {
				if f.Category == CategoryJoiner {
					joinerFindings++
				}
			}
			if joinerFindings != tc.wantRemoved {
				t.Errorf("joiner Findings = %d, want %d (all findings: %+v)", joinerFindings, tc.wantRemoved, res.Findings)
			}
		})
	}
}

func TestClean_Idempotent(t *testing.T) {
	inputs := []string{
		"plain ascii text",
		"\u65e5\u672c\u8a9e\u3000\u30c6\u30ad\u30b9\u30c8",
		"a\u00A0b\uFEFFc\u200Bd\u2060e\u00ADf",
		"a\u202Ab\u202Bc\u202Cd\u202Ee\u2066f\u2067g\u2068h\u2069i",
		"a\u200Db",
		"start\u200Dmiddle\u200Cend",
		"\u200Dleading and trailing\u200D",
		"a\u200D\u200D\u200Db",
		"a\x00\x01\x1f\x7f\u0080\u009fb",
		"e\u0301\u2764\uFE0F\u3000",
	}

	for _, in := range inputs {
		t.Run(in, func(t *testing.T) {
			first, err := Clean([]byte(in))
			if err != nil {
				t.Fatalf("first Clean error: %v", err)
			}
			second, err := Clean(first.Cleaned)
			if err != nil {
				t.Fatalf("second Clean error: %v", err)
			}
			if !bytes.Equal(first.Cleaned, second.Cleaned) {
				t.Errorf("not idempotent: first=%q second=%q", first.Cleaned, second.Cleaned)
			}
			if len(second.Findings) != 0 {
				t.Errorf("second pass produced findings: %+v", second.Findings)
			}
		})
	}
}

func TestClean_InvalidUTF8(t *testing.T) {
	_, err := Clean([]byte{0xff, 0xfe, 0x00})
	if err != ErrInvalidUTF8 {
		t.Fatalf("err = %v, want ErrInvalidUTF8", err)
	}
}
