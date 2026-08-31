# CLI Contract

The installed command is `clean-invisible-text` (`cmd/clean-invisible-text`).

```console
clean-invisible-text check [--json] [--allow RULE]... [--allow-file FILE] FILE...
clean-invisible-text fix [--json] [--keep-warnings] [--allow RULE]... [--allow-file FILE] FILE...
clean-invisible-text explain [--json] [--allow RULE]... [--allow-file FILE] FILE...
clean-invisible-text clean [--keep-warnings] [--allow RULE]... [--allow-file FILE]
```

Diagnostics go to standard error. Cleaned stream content goes only to
standard output (`clean` only). Exit statuses are defined in
[docs/specification.md](specification.md).

Multiple filenames are processed in one invocation of `check`/`fix`/`explain`
so pre-commit starts only one process. A file is written only when `fix`'s
cleaned content differs from what's on disk.

## Flags

- `--json`: on `check`/`fix`/`explain`, write a single JSON array to standard
  output instead of human-readable diagnostics to standard error (see "JSON
  output" below). Not available on `clean`, which only ever streams cleaned
  bytes on standard output.
- `--keep-warnings`: on `fix`/`clean`, preserve Warn-classified code points
  (see [docs/character-policy.md](character-policy.md)) in the output instead
  of removing them. The Finding is still reported either way; this only
  changes what ends up in the cleaned output. Not available on
  `check`/`explain`, which report every Warn finding regardless.
- `--allow`/`--allow-file`: on every subcommand, grant an audited exception
  to one or more Warn-classified code points. See "Allow-list flags" below.

## Allow-list flags

`--allow RULE` (repeatable) and `--allow-file FILE` implement the
[allow-list exceptions](character-policy.md#warn-allow-list-exceptions)
described by the character policy. Both are additive: rules from every
`--allow` flag and, if present, one config file are combined for the
invocation. If `--allow-file` is not given, a file named
`.clean-invisible-text-allow.json` in the current directory is loaded
automatically when it exists.

A rule has four fields: `codepoint` (required), `reason` (required),
`paths` (optional), and `max_run` (optional).

`--allow` packs one rule into a single semicolon-separated `key=value`
string:

```console
clean-invisible-text check --allow 'codepoint=U+E000;reason=Nerd Font icon glyph' notes.txt
```

- `codepoint`: one or more code points as `U+XXXX`, comma-separated for more
  than one (`U+E000,U+E001`).
- `reason`: free text; required and must be non-empty. It is copied onto
  every finding the rule allows, so the exception stays visible in output.
- `paths`: comma-separated glob patterns (`filepath.Match` syntax — no
  recursive `**`). A pattern containing `/` is matched against the file path
  as given on the command line; a pattern with no `/` is also matched
  against just the file's base name, so `*.md` matches a file in any
  directory. Omitted or empty means the rule applies to every file in the
  invocation. `clean` has no file path, so a `paths`-scoped rule never
  applies to it.
- `max_run`: a positive integer, or the literal `unlimited`. Governs the
  run-length guard described in the character policy; omitted means the
  default of 1 (a single isolated occurrence).

A config file loaded via `--allow-file` (or the default filename) is a JSON
array of the same fields:

```json
[
  {
    "codepoint": "U+E000",
    "reason": "Nerd Font icon glyph",
    "paths": ["*.md"],
    "max_run": 1
  }
]
```

`max_run` in a config file is a plain integer; use `-1` for unlimited.

When a rule's condition is met, the finding's `action` is `allow` instead of
`warn`/`remove`, and the code point is preserved in `fix`/`clean` output
regardless of `--keep-warnings`. `check`/`explain` do not fail (exit `0`,
not `1`) over a file whose only findings are `allow` — the allow-listed
finding is still printed, but it isn't why the command shows it as needing
attention.

## Human-readable diagnostics

`check` prints one line per file that has findings:

```
<path>: <N> finding(s)
```

`explain` and `fix` print one line per finding instead, since their purpose
is showing every finding in full:

```
<path>:<line>:<col>: U+<hex> <NAME> [<category>] -> <action>
```

`<action>` is `remove`, `replace with "<replacement>"`, `warn`, or
`allow (<reason>)` for a finding an [allow-list rule](#allow-list-flags)
covers. Line and column are 1-indexed; column counts runes, not bytes.

## JSON output

With `--json`, `check`/`fix`/`explain` write a single JSON array to standard
output, one object per file, in the same shape for all three commands:

```json
[
  {
    "path": "notes.txt",
    "findings": [
      {
        "line": 3,
        "column": 12,
        "offset": 34,
        "rune": "U+200B",
        "name": "ZERO WIDTH SPACE",
        "category": "zwsp",
        "action": "remove",
        "replacement": ""
      }
    ],
    "changed": false,
    "error": null
  }
]
```

`reason` is present, and non-empty, only on a finding whose `action` is
`allow` (see [Allow-list flags](#allow-list-flags)); it's omitted for every
other action.

`changed` is always present; it's only meaningful for `fix` (`check` and
`explain` never modify a file, so it's always `false` for them). `error` is
`null` on success and a string describing the failure (invalid UTF-8, a
binary-looking file, a rejected symlink, or an I/O error) when that file
could not be processed — `findings` is then empty and `changed` is `false`.
Nothing else is written to standard error in `--json` mode: the array is the
entire machine-readable contract, intended for CI and the Alfred integration.

## Compatibility

Subcommand names, flag names/semantics, exit statuses, and every JSON field
above are stable within a major version; new JSON fields may be added in a
minor version (additive only — consumers must ignore unknown fields).
Human-readable diagnostic text (non-`--json` output) is not covered and may
change in any release. Full policy, including the pre-1.0 caveat:
[ADR 0002](decisions/0002-v1-compatibility-and-support-policy.md#clijson-compatibility-policy).
