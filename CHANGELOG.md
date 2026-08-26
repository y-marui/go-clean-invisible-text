# Changelog

All notable changes to this project are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/).

## [Unreleased]

### Added

- Single-pass cleaning engine (`internal/cleaner`) implementing
  `docs/character-policy.md`: NBSP replacement, removal of BOM/ZWSP/WORD
  JOINER/SOFT HYPHEN, the Trojan Source bidirectional control set, Unicode tag
  characters, noncharacters, and unsafe control characters, plus contextual
  ZWJ/ZWNJ preservation.
- Safe UTF-8 file mutation (`internal/mutate`): rewrites a file in place via a
  same-directory temp file, fsync, and atomic rename; rejects symlinks and
  binary-looking content by default; skips the write when content is
  unchanged.
- Allow/Block/Warn model for invisible characters: an explicit whitelist
  (combining marks, ordinary/ideographic space, a single contextual ZWJ/ZWNJ
  or variation selector), the existing blacklist unchanged, and a new Warn
  outcome for anything else in Unicode category Cf/Co/Zs/Zl/Zp. Variation
  selectors get a dedicated contextual rule (a run of 2+, or one with no
  valid preceding base character, is blocked) since they're a documented
  ASCII-smuggling/steganography vector.
- CLI (`cmd/clean-invisible-text`): `check`, `fix`, `explain`, and `clean`
  commands, `--json` machine-readable output, `--keep-warnings`, and stable
  exit statuses per `docs/specification.md`.
- pre-commit hooks (`.pre-commit-hooks.yaml`): `clean-invisible-text-check`
  (detect only) and `clean-invisible-text-fix` (clean, then fail so the
  change is reviewed), plus an end-to-end fixture
  (`test/precommit/precommit_test.go`) and installation docs covering both
  the Go-source and pre-installed-binary paths
  (`docs/integrations/pre-commit.md`).
- Cross-platform release automation (`.github/workflows/release.yml`):
  builds darwin/windows/linux amd64+arm64 plus linux/armv7 with
  `CGO_ENABLED=0`, publishes SHA-256 checksums and build provenance
  attestations, and creates the GitHub Release on a `v*` tag push. See
  [docs/release-process.md](docs/release-process.md).
