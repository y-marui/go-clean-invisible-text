# Character Policy

This file is the normative source for character handling.

## Replace

- `U+00A0 NO-BREAK SPACE` becomes `U+0020 SPACE`.

Deletion is not used because it would join neighboring words.

## Remove

- Unicode bidirectional formatting and isolation controls used by Trojan Source attacks
- `U+200B ZERO WIDTH SPACE`
- `U+2060 WORD JOINER`
- `U+00AD SOFT HYPHEN`
- `U+FEFF`, both at the beginning and within text
- Unicode tag characters
- Unicode noncharacters
- unsafe control characters other than TAB, LF, and CR

The exact code-point tables must be versioned and covered by table-driven tests before implementation is released.

## Preserve

- TAB, LF, and CR
- combining characters
- variation selectors
- ordinary and ideographic visible spaces unless later specified
- a single contextual ZWJ or ZWNJ

## ZWJ and ZWNJ

`U+200D ZERO WIDTH JOINER` and `U+200C ZERO WIDTH NON-JOINER` are preserved only when a single character is meaningfully surrounded by non-whitespace, non-control text.

- At the start or end of text: remove.
- Adjacent to whitespace, newline, or a control character: remove.
- Repeated identical joiners: reduce the run to one, then apply the context rule.
- A directly mixed run of ZWJ and ZWNJ: remove the entire run.
- A single contextual joiner between visible text: preserve.

This policy protects legitimate emoji and writing systems without accepting redundant hidden runs.

## No implicit normalization

NFC, NFD, NFKC, and NFKD normalization are outside scope. Homoglyph replacement and smart-punctuation conversion are outside scope.
