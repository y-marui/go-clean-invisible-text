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

> **Status:** specification and roadmap. The CLI commands are not implemented
> yet (the detection/cleaning engine is implemented in `internal/cleaner`).

## Requirements

Aims to run as a standalone binary on Windows, macOS, and Raspberry Pi.
Building requires only Go — no Node.js, Python, or other runtime.

## Setup

```bash
git clone https://github.com/y-marui/go-clean-invisible-text.git
cd go-clean-invisible-text
go build ./...
go test ./...
```

## Usage

The CLI commands are not implemented yet. Planned commands:

~~~console
clean-invisible-text check FILE...
clean-invisible-text fix FILE...
clean-invisible-text explain FILE...
clean-invisible-text clean
~~~

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
- [docs/decisions/](docs/decisions/) — architecture decision records (ADRs)

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
