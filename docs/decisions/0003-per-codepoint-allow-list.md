# ADR 0003: per-code-point/file allow-list

## Status

Accepted on 2026-08-31 (see [Issue #26](https://github.com/y-marui/go-clean-invisible-text/issues/26)).

## Context

`--keep-warnings` is all-or-nothing per invocation: there was no way to
permit one specific Warn finding (e.g. a Nerd Font private-use icon glyph)
in one file while still catching unexpected Warn findings elsewhere. Issue
#26 proposed CLI flags plus a config file, a required `reason` per rule, and
a run-length guard, but explicitly left several implementation questions
undecided. This ADR records those decisions.

## Decision

### Config file format: JSON, not YAML

The issue's example used a `.yml` filename, but this project has no runtime
dependency beyond the Go standard library (`docs/architecture.md`), and
`encoding/json` covers the same "structured, hand-editable" need without
adding one. The config file is `.clean-invisible-text-allow.json` (a JSON
array of rule objects), loaded automatically from the current directory when
present, or from an explicit `--allow-file` path.

### Scope: only Warn-classified code points are ever eligible

An allow-list rule can only grant an exception to a Warn-classified code
point (`internal/cleaner.ActionWarn`). It has no effect on any
Block-classified code point — bidi controls, tag characters, noncharacters,
unsafe controls, NBSP, BOM, ZWSP, word joiner, soft hyphen, or a
non-contextual ZWJ/ZWNJ/variation selector. This is enforced in
`internal/cleaner.Clean` itself, not just by convention in the CLI layer.

This is a specific application of the character policy's ordinary Warn
outcome, not a general text-preservation exception: it holds even when a
project-committed `.clean-invisible-text-allow.json` is influenced by
someone other than the person reviewing a given change, since at most it
relabels one already-flagged Warn occurrence (with its `reason` kept visible
in output, per the "auditable" requirement in Issue #26); it can never
suppress a Trojan-Source-class or steganography-class Block finding.

### Overlapping-rule resolution: union

When multiple rules (from `--allow` flags and/or a config file) match the
same code point for the same file, every matching rule's `reason` is
recorded (joined for the finding, so every applicable exception stays
visible) and the effective run-length cap is the loosest among them — the
largest finite `max_run`, or unlimited if any matching rule is unlimited.
Rejected alternative: "most specific `paths` wins," which would silently
discard a broader rule's reason and read as a smaller allowance than at
least one rule the operator actually configured.

### Run-length guard: whole run reverts, not just the excess

A run of consecutive identical occurrences of an allow-listed code point
longer than the rule's `max_run` (default 1) falls back to ordinary Warn
treatment for every occurrence in the run, not just the ones past the cap.
This mirrors the existing ZWJ/ZWNJ and variation-selector policies in
`docs/character-policy.md`, where a run past the single-contextual-use case
is rejected as a whole rather than partially kept.

### check/explain exit status: an allow-listed finding is not "actionable"

`check` and `explain` continue to report every finding, including
`ActionAllow` ones (still visible for audit), but exit `0` (not `1`) when a
file's only findings are `ActionAllow`. Otherwise the allow-list would have
no operational effect on pre-commit/CI, defeating its purpose. `fix`'s exit
status is already based on whether the file's content changed, which an
`ActionAllow` finding never does, so it needed no equivalent change.

### Path matching: `filepath.Match`, no `**`

A `paths` pattern is matched with the standard library's `filepath.Match`
against the path as given to the command; a pattern with no path separator
is also matched against just the file's base name (so `*.md` matches a file
in any directory without needing a recursive-glob syntax). Recursive `**`
globbing is not supported — it isn't in the Go standard library, and adding
it would mean either a new dependency or a hand-rolled matcher for a need
that hasn't come up yet.

## Consequences

- `internal/allowlist` is a new package: parses `--allow` flag values and
  the JSON config file into `Rule`s, and reduces a rule set plus a file path
  to the `map[rune]cleaner.AllowRule` `internal/cleaner.Options.AllowRules`
  expects.
- `internal/cleaner.Finding` gains a `Reason` field (populated only for
  `ActionAllow`); `ActionAllow` is a new `ActionKind`. The `--json` `action`
  field gains a new possible value and each finding gains an optional
  `reason` field — additive per the CLI/JSON compatibility policy
  ([ADR 0002](0002-v1-compatibility-and-support-policy.md)), shipped as a
  minor version.
- `docs/character-policy.md` and `docs/cli.md` document the allow-list as
  normative behavior.
