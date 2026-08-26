# File Map

_Last updated: 2026-08-27_

## cleaner

| File | Role | Key dependencies |
|---|---|---|
| `internal/cleaner/cleaner.go` | `Clean()`: single-pass scan, ZWJ/ZWNJ context resolution, output assembly | `internal/cleaner/tables.go` |
| `internal/cleaner/tables.go` | Versioned code-point classification tables and `classify()` | — |
| `internal/cleaner/cleaner_test.go` | Table-driven, context-matrix, and idempotency tests | `internal/cleaner` |

## pre-commit integration

| File | Role | Key dependencies |
|---|---|---|
| `.pre-commit-hooks.yaml` | Shipped hook manifest: `clean-invisible-text-check`, `clean-invisible-text-fix` | `cmd/clean-invisible-text` |
| `test/precommit/precommit_test.go` | End-to-end fixture: builds the binary, runs both hooks via `pre-commit` in a scratch git repo, asserts exit codes and file content | `pre-commit` CLI (self-skips if absent), `cmd/clean-invisible-text` |

Note: this file does not yet cover `internal/cli/`, `internal/mutate/`, or
`cmd/clean-invisible-text/` — out of scope for this change; worth a follow-up
pass.
