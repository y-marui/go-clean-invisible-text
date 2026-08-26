# File Map

_Last updated: 2026-08-26_

## cleaner

| File | Role | Key dependencies |
|---|---|---|
| `internal/cleaner/cleaner.go` | `Clean()`: single-pass scan, ZWJ/ZWNJ context resolution, output assembly | `internal/cleaner/tables.go` |
| `internal/cleaner/tables.go` | Versioned code-point classification tables and `classify()` | — |
| `internal/cleaner/cleaner_test.go` | Table-driven, context-matrix, and idempotency tests | `internal/cleaner` |
