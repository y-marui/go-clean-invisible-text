// Package cleaner implements the single-pass cleaning engine defined by
// docs/character-policy.md. This file holds the versioned code-point tables;
// see character-policy.md for the normative source.
//
// Code points are written as Go \u escape sequences rather than literal
// characters so this source file does not itself contain the invisible and
// bidi-control characters the engine is built to detect.
package cleaner

import "unicode"

// Category classifies why a code point was changed.
type Category string

const (
	CategoryBidiControl       Category = "bidi-control"
	CategoryTag               Category = "tag"
	CategoryNoncharacter      Category = "noncharacter"
	CategoryControl           Category = "control"
	CategoryNBSP              Category = "nbsp"
	CategoryBOM               Category = "bom"
	CategorySoftHyphen        Category = "soft-hyphen"
	CategoryZWSP              Category = "zwsp"
	CategoryWordJoiner        Category = "word-joiner"
	CategoryJoiner            Category = "zwj-zwnj"
	CategoryVariationSelector Category = "variation-selector"
	CategoryUnclassified      Category = "unclassified"
)

// ActionKind is what the engine did to a code point.
type ActionKind string

const (
	ActionRemove  ActionKind = "remove"
	ActionReplace ActionKind = "replace"
	ActionWarn    ActionKind = "warn"
	// ActionAllow marks a Warn-classified code point permitted by an
	// AllowRule (see Options.AllowRules and internal/allowlist). It is never
	// produced for a Block-classified code point — the allow-list mechanism
	// cannot override Block, only grant an audited exception to Warn.
	ActionAllow ActionKind = "allow"
)

const (
	runeNBSP       = '\u00A0' // NO-BREAK SPACE
	runeBOM        = '\uFEFF' // ZERO WIDTH NO-BREAK SPACE (BOM)
	runeZWSP       = '\u200B' // ZERO WIDTH SPACE
	runeWordJoiner = '\u2060' // WORD JOINER
	runeSoftHyphen = '\u00AD' // SOFT HYPHEN
	runeZWJ        = '\u200D' // ZERO WIDTH JOINER
	runeZWNJ       = '\u200C' // ZERO WIDTH NON-JOINER
)

// bidiControls is the Trojan Source bidirectional formatting/isolation
// control set (character-policy.md "Remove", first bullet).
var bidiControls = map[rune]string{
	'\u202A': "left-to-right embedding",
	'\u202B': "right-to-left embedding",
	'\u202C': "pop directional formatting",
	'\u202D': "left-to-right override",
	'\u202E': "right-to-left override",
	'\u2066': "left-to-right isolate",
	'\u2067': "right-to-left isolate",
	'\u2068': "first strong isolate",
	'\u2069': "pop directional isolate",
}

// isTagCharacter reports whether r is one of the Unicode tag characters
// (U+E0000-U+E007F, including U+E0001 LANGUAGE TAG).
func isTagCharacter(r rune) bool {
	return r >= 0xE0000 && r <= 0xE007F
}

// isNoncharacter reports whether r is a Unicode noncharacter: U+FDD0-U+FDEF,
// or the last two code points of any plane (U+xFFFE, U+xFFFF).
func isNoncharacter(r rune) bool {
	if r >= 0xFDD0 && r <= 0xFDEF {
		return true
	}
	low := r & 0xFFFF
	return low == 0xFFFE || low == 0xFFFF
}

// isUnsafeControl reports whether r is a control character other than TAB,
// LF, or CR: C0 controls, DEL, and C1 controls.
func isUnsafeControl(r rune) bool {
	switch r {
	case '\t', '\n', '\r':
		return false
	}
	switch {
	case r >= 0x00 && r <= 0x1F:
		return true
	case r == 0x7F:
		return true
	case r >= 0x80 && r <= 0x9F:
		return true
	}
	return false
}

// isVariationSelector reports whether r is one of the 256 Unicode variation
// selectors (U+FE00-U+FE0F, U+E0100-U+E01EF). These are handled contextually
// (see resolveVariationSelectorRun in cleaner.go), not by classify, because a
// run of several in sequence is a documented steganography/ASCII-smuggling
// technique (see docs/character-policy.md).
func isVariationSelector(r rune) bool {
	return (r >= 0xFE00 && r <= 0xFE0F) || (r >= 0xE0100 && r <= 0xE01EF)
}

// isWhitelisted reports whether r is explicitly allowed to pass through
// unchanged, per the "Allow" section of docs/character-policy.md: ordinary
// and ideographic space, and combining marks (Unicode categories Mn, Mc,
// Me). TAB/LF/CR and a contextual ZWJ/ZWNJ are preserved elsewhere and are
// not covered by this function.
func isWhitelisted(r rune) bool {
	if r == ' ' || r == '\u3000' {
		return true
	}
	return unicode.In(r, unicode.Mn, unicode.Mc, unicode.Me)
}

// classifyWarn returns a Warn classification for r when it falls in a
// Unicode category that is invisible-or-format-like but not explicitly
// allowed or blocked: Cf (format), Co (private use), Zs/Zl/Zp (space, line,
// and paragraph separators). This is the "Warn" section of
// docs/character-policy.md: anything not on the Allow or Block list is
// flagged rather than silently passed through.
func classifyWarn(r rune) (classification, bool) {
	name := ""
	switch {
	case unicode.In(r, unicode.Cf):
		name = "Unicode format character"
	case unicode.In(r, unicode.Co):
		name = "Unicode private-use character"
	case unicode.In(r, unicode.Zs):
		name = "Unicode space separator"
	case unicode.In(r, unicode.Zl, unicode.Zp):
		name = "Unicode line/paragraph separator"
	default:
		return classification{}, false
	}
	return classification{CategoryUnclassified, ActionWarn, name, ""}, true
}

// classification is the outcome of classifying a single non-joiner rune.
type classification struct {
	category    Category
	action      ActionKind
	name        string
	replacement string
}

// classify returns the classification for r, or ok=false when r is preserved
// unchanged. Callers must handle ZWJ/ZWNJ (runeZWJ, runeZWNJ) separately
// before calling classify, since their handling is contextual rather than
// per-rune.
func classify(r rune) (classification, bool) {
	switch r {
	case runeNBSP:
		return classification{CategoryNBSP, ActionReplace, "NO-BREAK SPACE", " "}, true
	case runeBOM:
		return classification{CategoryBOM, ActionRemove, "ZERO WIDTH NO-BREAK SPACE (BOM)", ""}, true
	case runeZWSP:
		return classification{CategoryZWSP, ActionRemove, "ZERO WIDTH SPACE", ""}, true
	case runeWordJoiner:
		return classification{CategoryWordJoiner, ActionRemove, "WORD JOINER", ""}, true
	case runeSoftHyphen:
		return classification{CategorySoftHyphen, ActionRemove, "SOFT HYPHEN", ""}, true
	}
	if name, ok := bidiControls[r]; ok {
		return classification{CategoryBidiControl, ActionRemove, name, ""}, true
	}
	if isTagCharacter(r) {
		return classification{CategoryTag, ActionRemove, "Unicode tag character", ""}, true
	}
	if isNoncharacter(r) {
		return classification{CategoryNoncharacter, ActionRemove, "noncharacter", ""}, true
	}
	if isUnsafeControl(r) {
		return classification{CategoryControl, ActionRemove, "unsafe control character", ""}, true
	}
	if isWhitelisted(r) {
		return classification{}, false
	}
	if c, ok := classifyWarn(r); ok {
		return c, true
	}
	return classification{}, false
}

// joinerName returns the descriptive name for a ZWJ/ZWNJ rune.
func joinerName(r rune) string {
	if r == runeZWJ {
		return "ZERO WIDTH JOINER"
	}
	return "ZERO WIDTH NON-JOINER"
}
