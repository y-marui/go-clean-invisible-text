# Clean Invisible Text

> **This is the reference (English) version.**
> The canonical (Japanese) version is [README-jp.md](README-jp.md).

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![CI](https://github.com/y-marui/go-clean-invisible-text/actions/workflows/ci.yml/badge.svg)](https://github.com/y-marui/go-clean-invisible-text/actions/workflows/ci.yml)
[![Charter Check](https://github.com/y-marui/go-clean-invisible-text/actions/workflows/dev-charter-check.yml/badge.svg)](https://github.com/y-marui/go-clean-invisible-text/actions/workflows/dev-charter-check.yml)
[![GitHub Sponsors](https://img.shields.io/github/sponsors/y-marui?style=social)](https://github.com/sponsors/y-marui)
[![Buy Me a Coffee](https://img.shields.io/badge/Buy%20Me%20a%20Coffee-donate-yellow.svg)](https://www.buymeacoffee.com/y.marui)

A cross-platform Go CLI for detecting, explaining, and safely cleaning
dangerous invisible Unicode characters in UTF-8 plain text.

> **Status:** v0.1 in progress. `check`/`fix`/`explain`/`clean` and pre-commit
> packaging are implemented; cross-platform release automation is not yet
> (see the
> [roadmap](https://github.com/y-marui/go-clean-invisible-text/issues/1)).

## Requirements

Aims to run as a standalone binary on Windows, macOS, and Raspberry Pi.
Building requires only Go — no Node.js, Python, or other runtime. Minimum OS
versions for official release binaries and supported Raspberry Pi
architectures are defined in
[ADR 0002](docs/decisions/0002-v1-compatibility-and-support-policy.md).
Source-built binaries follow the requirements of the Go toolchain that builds
them.

## Setup

```bash
go install github.com/y-marui/go-clean-invisible-text/cmd/clean-invisible-text@latest
```

Or build from a clone:

```bash
git clone https://github.com/y-marui/go-clean-invisible-text.git
cd go-clean-invisible-text
go build -o clean-invisible-text ./cmd/clean-invisible-text
go test ./...
```

## Usage

| Command | Description |
|---|---|
| `clean-invisible-text check FILE...` | Report findings without modifying input |
| `clean-invisible-text fix FILE...` | Modify named files and report every change |
| `clean-invisible-text explain FILE...` | Show code point, Unicode name, location, category, and planned action |
| `clean-invisible-text clean` | Read standard input and write cleaned text to standard output |

```console
$ clean-invisible-text check notes.txt
notes.txt: 2 finding(s)

$ echo "hello world" | clean-invisible-text clean
hello world
```

Add `--json` to `check`/`fix`/`explain` for machine-readable output, and
`--keep-warnings` to `fix`/`clean` to preserve Warn-classified code points
instead of removing them. Full contract: [docs/cli.md](docs/cli.md).

The normative behavior is defined in
[docs/specification.md](docs/specification.md) and
[docs/character-policy.md](docs/character-policy.md). Work in progress belongs
in GitHub Issues, not in specification files.

## Documentation

- [docs/specification.md](docs/specification.md) — functional specification
- [docs/character-policy.md](docs/character-policy.md) — character policy (normative)
- [docs/security-model.md](docs/security-model.md) — security model
- [docs/cli.md](docs/cli.md) — CLI contract
- [docs/architecture.md](docs/architecture.md) — architecture
- [docs/integrations/pre-commit.md](docs/integrations/pre-commit.md) — pre-commit hook contract
- [docs/decisions/](docs/decisions/) — architecture decision records (ADRs)

## pre-commit Integration

```yaml
repos:
  - repo: https://github.com/y-marui/go-clean-invisible-text
    rev: vX.Y.Z # pin to a released tag
    hooks:
      - id: clean-invisible-text-check
      - id: clean-invisible-text-fix
```

Full contract, including the pre-installed-binary option: [docs/integrations/pre-commit.md](docs/integrations/pre-commit.md).

## Alfred Integration

The separate
[alfred-clean-invisible-text](https://github.com/y-marui/alfred-clean-invisible-text)
repository will package this CLI for Alfred.

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md).

## Changelog

See [CHANGELOG.md](CHANGELOG.md).

## License

[MIT](LICENSE)

---
*This document has a Japanese canonical version [README-jp.md](README-jp.md). Update both in the same commit when editing.*
