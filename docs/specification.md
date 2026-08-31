# Specification

## Scope

Clean Invisible Text processes UTF-8 plain text. It does not parse programming languages, Markdown, rich text, URLs, HTML entities, or encoded payload layers.

## Inputs

- One or more files
- Standard input

Invalid UTF-8 is an error and must not be modified. Binary-looking input is rejected. Symbolic-link mutation is rejected by default.

## Preservation

Unless a code point is explicitly covered by the character policy:

- input bytes are preserved;
- LF and CRLF are preserved;
- a final newline is preserved;
- Unicode normalization is not performed;
- visible punctuation and homoglyphs are not rewritten;
- files without changes are not rewritten.

## Operations

- `check`: report findings without modifying input.
- `fix`: modify named files and report every change.
- `explain`: show code point, Unicode name, location, category, and planned action.
- `clean`: read standard input and write cleaned text to standard output.

## Exit status

- `0`: no finding, or successful standard-stream cleaning. For `check`/
  `explain`, a file whose only findings are allow-listed exceptions (see
  [docs/cli.md](cli.md#allow-list-flags)) also exits `0`.
- `1`: findings were detected; for `fix`, changes were applied and must be reviewed.
- `2`: invalid arguments, input, encoding, or I/O failure.

A pre-commit fix must exit with status 1 after changing a file so the resulting diff is reviewed before commit.

## File mutation

File replacement must avoid partial output on failure and preserve applicable permissions. Platform-specific metadata guarantees will be documented before v1.0.
