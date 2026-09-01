# Changelog

All notable changes to this project are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/).

## [1.1.1] - 2026-09-01

### Added

- macOS release binaries (`darwin-amd64`, `darwin-arm64`) are now signed with
  a Developer ID Application certificate and notarised by Apple as part of
  `.github/workflows/release.yml`, before checksumming/attestation/publishing
  ([Issue #31](https://github.com/y-marui/go-clean-invisible-text/issues/31)).
  Requires five repository secrets to be configured; see
  [docs/release-process.md](docs/release-process.md#macos-code-signing-and-notarization).

## [1.1.0] - 2026-08-31

### Added

- Per-code-point/file allow-list ([Issue #26](https://github.com/y-marui/go-clean-invisible-text/issues/26),
  [ADR 0003](docs/decisions/0003-per-codepoint-allow-list.md)): repeatable
  `--allow` flags and a `--allow-file`/`.clean-invisible-text-allow.json`
  config file grant an audited, required-`reason` exception for a specific
  Warn-classified code point, optionally scoped to file path globs, with a
  run-length guard (default: a single isolated occurrence) that falls back
  to ordinary Warn treatment for a longer run. The mechanism can never
  affect a Block-classified code point. A new `internal/allowlist` package
  implements rule parsing and resolution; `internal/cleaner` gains
  `ActionAllow` and `Finding.Reason`; the `--json` contract gains the
  `allow` action value and an optional `reason` field (additive, per
  [ADR 0002](docs/decisions/0002-v1-compatibility-and-support-policy.md)).
  `check`/`explain` no longer fail over a file whose only findings are
  allow-listed.

## [1.0.0] - 2026-08-27

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
- v1.0 compatibility and support policy
  ([ADR 0002](docs/decisions/0002-v1-compatibility-and-support-policy.md)):
  minimum OS versions tied to the exact official-release toolchain, supported
  Raspberry Pi architectures, the Unicode table update policy's versioning
  consequences, a CLI/JSON compatibility contract, and the security
  response/release policy.
