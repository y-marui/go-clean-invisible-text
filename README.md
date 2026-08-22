# Clean Invisible Text

A planned cross-platform Go CLI for detecting, explaining, and safely cleaning dangerous invisible Unicode characters in UTF-8 plain text.

> Status: specification and roadmap. The CLI is not implemented yet.

## Goals

- Preserve text semantics and bytes unrelated to an intentional edit.
- Run as a standalone binary on Windows, macOS, and Raspberry Pi.
- Support files, standard input/output, CI, and pre-commit.
- Explain every detected, removed, or replaced code point.
- Keep legitimate combining marks, variation selectors, and contextual joiners.

## Planned commands

~~~console
clean-invisible-text check FILE...
clean-invisible-text fix FILE...
clean-invisible-text explain FILE...
clean-invisible-text clean
~~~

The normative behavior is defined in [docs/specification.md](docs/specification.md) and [docs/character-policy.md](docs/character-policy.md). Work in progress belongs in GitHub Issues and Projects, not in specification files.

## Alfred integration

The separate [alfred-clean-invisible-text](https://github.com/y-marui/alfred-clean-invisible-text) repository will package this CLI for Alfred.

## License

MIT
