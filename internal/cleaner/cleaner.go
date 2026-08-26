package cleaner

import (
	"errors"
	"unicode"
	"unicode/utf8"
)

// ErrInvalidUTF8 is returned by Clean when the input is not valid UTF-8.
var ErrInvalidUTF8 = errors.New("cleaner: invalid UTF-8 input")

// Options configures Clean.
type Options struct {
	// KeepWarnings, if true, preserves Warn-classified code points in
	// Cleaned instead of removing them. A Finding is still reported either
	// way — this only controls what ends up in the output bytes.
	KeepWarnings bool
}

// Finding describes one code point the engine changed or flagged. There is
// no Finding for a code point that was preserved unchanged.
type Finding struct {
	Offset      int // byte offset in the original input
	Rune        rune
	Name        string
	Category    Category
	Action      ActionKind
	Replacement string // "" for ActionRemove/ActionWarn, the replacement text for ActionReplace
}

// Result is the outcome of Clean.
type Result struct {
	Cleaned  []byte
	Findings []Finding
}

// pendingRune is one buffered rune (a ZWJ/ZWNJ or a variation selector)
// awaiting context resolution.
type pendingRune struct {
	offset int
	r      rune
}

// Clean processes UTF-8 input in a single deterministic pass and returns the
// cleaned bytes plus a Finding for every code point that was changed or
// flagged, per docs/character-policy.md.
func Clean(input []byte, opts Options) (Result, error) {
	if !utf8.Valid(input) {
		return Result{}, ErrInvalidUTF8
	}

	out := make([]byte, 0, len(input))
	var findings []Finding

	var joinerRun []pendingRune
	var vsRun []pendingRune
	var prev rune
	havePrev := false

	isBreak := func(r rune, have bool) bool {
		if !have {
			return true // start or end of text
		}
		return unicode.IsSpace(r) || unicode.IsControl(r)
	}

	// resolveJoinerRun decides the fate of any buffered ZWJ/ZWNJ run using
	// the rune immediately before it (prev/havePrev) and the rune
	// immediately after it (next/haveNext), both from the original input.
	resolveJoinerRun := func(next rune, haveNext bool) {
		if len(joinerRun) == 0 {
			return
		}
		mixed := false
		for _, pr := range joinerRun {
			if pr.r != joinerRun[0].r {
				mixed = true
				break
			}
		}
		contextOK := !mixed && !isBreak(prev, havePrev) && !isBreak(next, haveNext)
		if contextOK {
			out = utf8.AppendRune(out, joinerRun[0].r)
			joinerRun = joinerRun[1:]
		}
		for _, pr := range joinerRun {
			findings = append(findings, Finding{
				Offset:   pr.offset,
				Rune:     pr.r,
				Name:     joinerName(pr.r),
				Category: CategoryJoiner,
				Action:   ActionRemove,
			})
		}
		joinerRun = nil
	}

	// resolveVSRun decides the fate of any buffered variation-selector run
	// using only the rune immediately before it (prev/havePrev) — a
	// variation selector modifies what precedes it, so unlike a joiner it
	// doesn't need "next" context. A run of more than one, or one with no
	// valid preceding base character, is the documented steganography shape
	// (see docs/character-policy.md) and is removed.
	resolveVSRun := func() {
		if len(vsRun) == 0 {
			return
		}
		if len(vsRun) == 1 && !isBreak(prev, havePrev) {
			out = utf8.AppendRune(out, vsRun[0].r)
			vsRun = nil
			return
		}
		for _, pr := range vsRun {
			findings = append(findings, Finding{
				Offset:   pr.offset,
				Rune:     pr.r,
				Name:     "VARIATION SELECTOR",
				Category: CategoryVariationSelector,
				Action:   ActionRemove,
			})
		}
		vsRun = nil
	}

	offset := 0
	for offset < len(input) {
		r, size := utf8.DecodeRune(input[offset:])
		curOffset := offset
		offset += size

		if r == runeZWJ || r == runeZWNJ {
			resolveVSRun()
			joinerRun = append(joinerRun, pendingRune{curOffset, r})
			continue
		}

		if isVariationSelector(r) {
			resolveJoinerRun(r, true)
			vsRun = append(vsRun, pendingRune{curOffset, r})
			continue
		}

		resolveJoinerRun(r, true)
		resolveVSRun()

		if c, ok := classify(r); ok {
			findings = append(findings, Finding{
				Offset:      curOffset,
				Rune:        r,
				Name:        c.name,
				Category:    c.category,
				Action:      c.action,
				Replacement: c.replacement,
			})
			switch c.action {
			case ActionReplace:
				out = append(out, c.replacement...)
			case ActionWarn:
				if opts.KeepWarnings {
					out = utf8.AppendRune(out, r)
				}
			}
		} else {
			out = utf8.AppendRune(out, r)
		}

		prev = r
		havePrev = true
	}
	resolveJoinerRun(0, false)
	resolveVSRun()

	return Result{Cleaned: out, Findings: findings}, nil
}
