# Architecture

## Overview

A Go CLI, not yet implemented as a CLI: today the repository only contains the
single-pass character-cleaning engine. The engine classifies and rewrites UTF-8
text per `docs/character-policy.md`; a future CLI layer will wrap it for files,
stdin/stdout, and JSON output per `docs/cli.md`.

## Entry Points

- `internal/cleaner/cleaner.go` — `Clean([]byte) (Result, error)`, the only
  entry point today. There is no `cmd/` or `main` package yet (tracked by
  GitHub Issue #2).

## Directory Structure

| Directory | Role |
|---|---|
| `internal/cleaner/` | Single-pass cleaning engine: code-point tables (`tables.go`) and the scan/rewrite loop (`cleaner.go`) |
| `docs/` | Normative specification, character policy, security model, CLI contract, ADRs |
| `docs/decisions/` | Architecture decision records |
| `docs/integrations/` | Integration contracts (e.g. pre-commit) for the future CLI |
| `docs/dev-charter/` | Shared dev-charter, imported via `git subtree` — do not edit directly |

## Key Dependencies

None. The engine uses only the Go standard library (`unicode`, `unicode/utf8`).
