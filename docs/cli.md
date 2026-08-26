# CLI Contract

The installed command is `clean-invisible-text` (`cmd/clean-invisible-text`).

```console
clean-invisible-text check [--json] FILE...
clean-invisible-text fix [--json] [--keep-warnings] FILE...
clean-invisible-text explain [--json] FILE...
clean-invisible-text clean [--keep-warnings]
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

`<action>` is `remove`, `replace with "<replacement>"`, or `warn`. Line and
column are 1-indexed; column counts runes, not bytes.

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

`changed` is always present; it's only meaningful for `fix` (`check` and
`explain` never modify a file, so it's always `false` for them). `error` is
`null` on success and a string describing the failure (invalid UTF-8, a
binary-looking file, a rejected symlink, or an I/O error) when that file
could not be processed — `findings` is then empty and `changed` is `false`.
Nothing else is written to standard error in `--json` mode: the array is the
entire machine-readable contract, intended for CI and the Alfred integration.
