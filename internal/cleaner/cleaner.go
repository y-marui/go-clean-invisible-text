package cleaner

import (
	"errors"
	"unicode"
	"unicode/utf8"
)

// ErrInvalidUTF8 is returned by Clean when the input is not valid UTF-8.
var ErrInvalidUTF8 = errors.New("cleaner: invalid UTF-8 input")

// Finding describes one code point the engine changed. There is no Finding
// for a code point that was preserved unchanged.
type Finding struct {
	Offset      int // byte offset in the original input
	Rune        rune
	Name        string
	Category    Category
	Action      ActionKind
	Replacement string // "" for ActionRemove, the replacement text for ActionReplace
}

// Result is the outcome of Clean.
type Result struct {
	Cleaned  []byte
	Findings []Finding
}

// pendingJoiner is one buffered ZWJ/ZWNJ rune awaiting context resolution.
type pendingJoiner struct {
	offset int
	r      rune
}

// Clean processes UTF-8 input in a single deterministic pass and returns the
// cleaned bytes plus a Finding for every code point that was changed, per
// docs/character-policy.md.
func Clean(input []byte) (Result, error) {
	if !utf8.Valid(input) {
		return Result{}, ErrInvalidUTF8
	}

	out := make([]byte, 0, len(input))
	var findings []Finding

	var run []pendingJoiner
	var prev rune
	havePrev := false

	isBreak := func(r rune, have bool) bool {
		if !have {
			return true // start or end of text
		}
		return unicode.IsSpace(r) || unicode.IsControl(r)
	}

	// resolveRun decides the fate of any buffered ZWJ/ZWNJ run using the rune
	// immediately before it (prev/havePrev) and the rune immediately after it
	// (next/haveNext), both from the original input.
	resolveRun := func(next rune, haveNext bool) {
		if len(run) == 0 {
			return
		}
		mixed := false
		for _, pr := range run {
			if pr.r != run[0].r {
				mixed = true
				break
			}
		}
		contextOK := !mixed && !isBreak(prev, havePrev) && !isBreak(next, haveNext)
		if contextOK {
			out = utf8.AppendRune(out, run[0].r)
			run = run[1:]
		}
		for _, pr := range run {
			findings = append(findings, Finding{
				Offset:   pr.offset,
				Rune:     pr.r,
				Name:     joinerName(pr.r),
				Category: CategoryJoiner,
				Action:   ActionRemove,
			})
		}
		run = nil
	}

	offset := 0
	for offset < len(input) {
		r, size := utf8.DecodeRune(input[offset:])
		curOffset := offset
		offset += size

		if r == runeZWJ || r == runeZWNJ {
			run = append(run, pendingJoiner{curOffset, r})
			continue
		}

		resolveRun(r, true)

		if c, ok := classify(r); ok {
			findings = append(findings, Finding{
				Offset:      curOffset,
				Rune:        r,
				Name:        c.name,
				Category:    c.category,
				Action:      c.action,
				Replacement: c.replacement,
			})
			if c.action == ActionReplace {
				out = append(out, c.replacement...)
			}
		} else {
			out = utf8.AppendRune(out, r)
		}

		prev = r
		havePrev = true
	}
	resolveRun(0, false)

	return Result{Cleaned: out, Findings: findings}, nil
}
