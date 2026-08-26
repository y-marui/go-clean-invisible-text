# Developing

## Build and Test

```bash
make build   # go build ./...
make test    # go test ./...
make lint    # gofmt -l . && go vet ./...
make fmt     # gofmt -w .
```

See the [Makefile](Makefile) for the exact commands each target runs.

## Requirements

- Go (see [go.mod](go.mod) for the minimum version)
- [pre-commit](https://pre-commit.com/) for the security and documentation
  hooks in [.pre-commit-config.yaml](.pre-commit-config.yaml)

## Project Layout

See [docs/architecture.md](docs/architecture.md) for directory structure and
[docs/file-map.md](docs/file-map.md) for file-level dependencies.

## Conventions

- The normative source for character-handling behavior is
  [docs/character-policy.md](docs/character-policy.md). Code in
  `internal/cleaner/` implements exactly that policy; a behavior change starts
  with a specification update, not the other way around.
- Table-driven tests for every classification category, plus an idempotency
  test (`Clean(Clean(x).Cleaned)` produces no findings) for any new
  code-point handling — see `internal/cleaner/cleaner_test.go`.
- Package names and file layout follow standard Go conventions
  (`internal/<package>/<file>.go`, `_test.go` alongside the code it tests).
- No comments that restate what the code does; comments explain non-obvious
  *why* only (see [docs/dev-charter/CODE_STYLE.md](docs/dev-charter/CODE_STYLE.md)).

## Commit Messages

[Conventional Commits](https://www.conventionalcommits.org/) format
(`feat:`, `fix:`, `docs:`, `refactor:`, `chore:`, `test:`).

## Debugging

`internal/cleaner` is a pure function over `[]byte` with no I/O — reproduce any
reported issue as a table-driven test case first (see
[.github/ISSUE_TEMPLATE/bug.yml](.github/ISSUE_TEMPLATE/bug.yml) for the report
format), then fix.
