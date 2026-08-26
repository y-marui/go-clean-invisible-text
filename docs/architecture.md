# Architecture

## Overview

A Go CLI (`clean-invisible-text`) built around a single-pass character-cleaning
engine. The engine classifies and rewrites UTF-8 text per
`docs/character-policy.md`; the CLI layer wraps it for files, stdin/stdout,
and JSON output per `docs/cli.md`, and `internal/mutate` applies it to files
on disk with failure-safe replacement per `docs/security-model.md`.

## Entry Points

- `cmd/clean-invisible-text/main.go` — the `clean-invisible-text` binary;
  dispatches to `internal/cli.Run`.
- `internal/cleaner/cleaner.go` — `Clean([]byte) (Result, error)`, the pure
  engine entry point used by both the CLI and `internal/mutate`.

## Directory Structure

| Directory | Role |
|---|---|
| `cmd/clean-invisible-text/` | `main` package; thin wrapper over `internal/cli.Run` |
| `internal/cleaner/` | Single-pass cleaning engine: code-point tables (`tables.go`) and the scan/rewrite loop (`cleaner.go`) |
| `internal/mutate/` | Safe in-place file rewriting: binary/symlink rejection, atomic same-directory replace |
| `internal/cli/` | Subcommand dispatch (`check`/`fix`/`explain`/`clean`) and human-readable/JSON reporting |
| `test/precommit/` | End-to-end fixture exercising the pre-commit hook contract from `.pre-commit-hooks.yaml` |
| `docs/` | Normative specification, character policy, security model, CLI contract, ADRs |
| `docs/decisions/` | Architecture decision records |
| `docs/integrations/` | Integration contracts (e.g. pre-commit) |
| `docs/dev-charter/` | Shared dev-charter, imported via `git subtree` — do not edit directly |

## Key Dependencies

None. The engine and CLI use only the Go standard library.
