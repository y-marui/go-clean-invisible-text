# Character Policy

This file is the normative source for character handling. Every code point
falls into exactly one of three outcomes: **Allow** (preserved silently),
**Block** (removed or replaced), or **Warn** (flagged, and removed by
default unless explicitly kept). Nothing passes through silently just
because nobody thought to list it: a code point that is not explicitly
Allowed or Blocked is Warned on.

## Unicode version

Code-point classification (general categories such as Mn, Mc, Me, Cf, Co,
Zs, Zl, Zp) relies entirely on Go's standard library `unicode` package rather
than a vendored table, so the effective Unicode Character Database version is
whatever `unicode.Version` reports for the Go toolchain pinned in `go.mod`.
At `go 1.27.0` that is Unicode **17.0.0**. Bumping the pinned Go version may
change this value; no separate update to this document is needed unless a
category assignment used by this policy actually changes between versions.
When one does, that's a behavior change, not a documentation update — see
the versioning consequences in
[ADR 0002](decisions/0002-v1-compatibility-and-support-policy.md#unicode-table-update-policy).

## Allow

Preserved unchanged, with no finding. This is a closed, enumerated set.

- TAB, LF, and CR
- ordinary SPACE (`U+0020`) and ideographic space (`U+3000`)
- combining characters (Unicode general categories Mn, Mc, Me)
- a single contextual ZWJ or ZWNJ (see "ZWJ and ZWNJ" below)
- a single contextual variation selector (see "Variation selectors" below)

## Block: Replace

- `U+00A0 NO-BREAK SPACE` becomes `U+0020 SPACE`.

Deletion is not used because it would join neighboring words.

## Block: Remove

- Unicode bidirectional formatting and isolation controls used by Trojan Source attacks
- `U+200B ZERO WIDTH SPACE`
- `U+2060 WORD JOINER`
- `U+00AD SOFT HYPHEN`
- `U+FEFF`, both at the beginning and within text
- Unicode tag characters
- Unicode noncharacters
- unsafe control characters other than TAB, LF, and CR
- a variation selector outside the single-contextual-use case below (see "Variation selectors")

## ZWJ and ZWNJ

`U+200D ZERO WIDTH JOINER` and `U+200C ZERO WIDTH NON-JOINER` are preserved only when a single character is meaningfully surrounded by non-whitespace, non-control text.

- At the start or end of text: remove.
- Adjacent to whitespace, newline, or a control character: remove.
- Repeated identical joiners: reduce the run to one, then apply the context rule.
- A directly mixed run of ZWJ and ZWNJ: remove the entire run.
- A single contextual joiner between visible text: preserve.

This policy protects legitimate emoji and writing systems without accepting redundant hidden runs.

## Variation selectors

Unicode defines 256 variation selectors (`U+FE00`-`U+FE0F`, `U+E0100`-`U+E01EF`),
exactly one per possible byte value. That makes a sequence of variation
selectors a ready-made channel for encoding arbitrary hidden data or
instructions that render as nothing — a documented, actively exploited
technique (variously called ASCII smuggling or Unicode steganography), not a
theoretical concern:

- a malicious npm package (`os-info-checker-es6`, published May 2025) used
  variation selectors to hide command-and-control configuration data
- it is a cataloged prompt-injection and data-exfiltration technique against
  LLM agents that read untrusted text

Legitimate use is always exactly one variation selector immediately after a
visible base character (emoji presentation selection, or an Ideographic
Variation Sequence after a CJK ideograph) — never a run of several in
sequence. The rule mirrors ZWJ/ZWNJ handling above, applied to whichever
rune immediately precedes the selector (a variation selector modifies what
came before it, so no "next" context is needed):

- A run of two or more consecutive variation selectors: remove the entire run.
- A single variation selector at the start of text, or adjacent to whitespace, newline, or a control character: remove.
- A single variation selector immediately after visible text: preserve.

## Warn

Anything not covered by Allow or Block above that falls in Unicode general
category Cf (format), Co (private use), Zs (space separator — figure space,
narrow no-break space, em/en spaces, and similar), Zl, or Zp (line/paragraph
separator) is flagged as a finding rather than silently preserved.

`check` and `explain` always report every Warn finding regardless of any
option. `fix` and `clean` remove Warn-classified code points from the output
by default, matching Block; a `KeepWarnings` option preserves them in the
output instead, while the finding is still reported.

This resolves [Issue #9](https://github.com/y-marui/go-clean-invisible-text/issues/9)
("additional space-like characters"): rather than deciding each additional
space or format character individually in advance, anything beyond the
Allow list now surfaces as a warning automatically.

## No implicit normalization

NFC, NFD, NFKC, and NFKD normalization are outside scope. Homoglyph replacement and smart-punctuation conversion are outside scope.
