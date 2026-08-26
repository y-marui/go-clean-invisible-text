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

## pre-commit Hooks

`.pre-commit-hooks.yaml` defines the `clean-invisible-text-check` and
`clean-invisible-text-fix` hooks this repository ships to consumers — see
[docs/integrations/pre-commit.md](docs/integrations/pre-commit.md) for the
contract and both installation options. `test/precommit/precommit_test.go` is
an end-to-end fixture for that contract; it skips itself when `pre-commit` is
not on `PATH`, so `make test` stays runnable without it, but run it directly
after touching either hook's behavior:

```bash
go test ./test/precommit/...
```

## Release Process

`.github/workflows/release.yml` builds every architecture in
[Issue #7](https://github.com/y-marui/go-clean-invisible-text/issues/7) on a
`v*` tag push. See [docs/release-process.md](docs/release-process.md) for the
target matrix, how to cut a release, and how to verify a downloaded binary's
checksum and provenance attestation. `workflow_dispatch` runs the same build
matrix without publishing anything — use it to validate the workflow after
changing it.
