# ADR 0002: v1.0 compatibility and support policy

## Status

Accepted on 2026-08-27 (see [Issue #8](https://github.com/y-marui/go-clean-invisible-text/issues/8)).

## Context

`docs/character-policy.md`, `docs/cli.md`, `docs/security-model.md`, and
`README.md` describe today's behavior, but none of them commit to what stays
stable across releases. Issue #8 asks for five provisional guarantees to be
turned into stable contracts before v1.0: minimum OS versions, supported
Raspberry Pi architectures, the Unicode table update policy, CLI/JSON
compatibility, and the security response/release policy.

This project already publishes a `CHANGELOG.md` "based on Keep a Changelog"
and uses Conventional Commits; Keep a Changelog itself recommends pairing
with [Semantic Versioning](https://semver.org/), so this ADR adopts SemVer
explicitly rather than leaving it implicit.

## Decision

### Minimum OS versions

Official release binaries track the OS floor required by the exact Go version
used by the release workflow, rather than a hardcoded number that silently
goes stale. The workflow reads the patch-qualified `go` directive from
`go.mod`; `actions/setup-go` resolves that value exactly. Changing the release
toolchain or its build configuration must be called out in `CHANGELOG.md` when
it changes the floor. See the
[Go Wiki's MinimumRequirements page](https://go.dev/wiki/MinimumRequirements)
for current values; for release binaries built with `go 1.27.0`:

- Windows 10 / Windows Server 2016 or later
- macOS 13 (Ventura) or later
- Linux kernel 3.2 or later

The `go` directive is a minimum version for source builds, not a general
toolchain pin. A developer or `go install ...@latest` may use a newer
compatible toolchain, so a source-built binary follows that toolchain's OS
requirements and is not covered by the official-binary floor above.

### Supported Raspberry Pi architectures

Matches the `linux/arm64` and `linux/armv7` targets in
[Issue #7](https://github.com/y-marui/go-clean-invisible-text/issues/7):

- `linux/arm64`: Raspberry Pi 3 and later running a 64-bit OS (Raspberry Pi
  OS 64-bit, Ubuntu, etc.)
- `linux/armv7` (`GOARM=7`): Raspberry Pi 2 and later running 32-bit
  Raspberry Pi OS

Raspberry Pi 1/Zero (ARMv6) is explicitly out of scope — it was never in
Issue #7's target list. Support may be added later as a new decision, not by
silent inference from this ADR.

### Unicode table update policy

Formalizes what `docs/character-policy.md`'s "Unicode version" section
already implies: classification tracks `unicode.Version` for the toolchain
that built the binary, not a vendored table. Official release binaries use
the exact patch-qualified Go version selected by the release workflow;
source builds, including `go install ...@latest`, use the toolchain actually
selected by the invoking `go` command and may therefore use a newer Unicode
table.

A release-toolchain bump that changes the Cf/Co/Zs/Zl/Zp category (or any
other category this policy relies on — Mn/Mc/Me for the Allow list) of a code
point this policy already classifies is a **behavior change** — it can change
`fix`/`clean` output — and requires at least a **minor** version bump plus a
`CHANGELOG.md` entry. It is never silently absorbed into a patch release. A
Go bump that changes no relevant category assignment needs no
character-policy update at all, per the existing text.

### CLI/JSON compatibility policy

Within a major version:

- Subcommand names (`check`/`fix`/`explain`/`clean`), flag names and
  semantics (`--json`, `--keep-warnings`), and exit statuses (`0`/`1`/`2`,
  `docs/specification.md`) are stable.
- Every JSON field documented in `docs/cli.md` is stable in name, meaning,
  and type. New fields may be added in a **minor** version (additive only);
  consumers must ignore unknown fields.
- Removing or renaming a field, or changing its type or exit-status meaning,
  requires a **major** version bump.

Human-readable diagnostics (stderr text from `check`/`fix`/`explain` without
`--json`) are explicitly **not** covered by this guarantee and may change in
any release, including a patch — `--json` is the stable contract for
scripting, CI, pre-commit, and the Alfred integration.

Pre-1.0, SemVer §4 already permits breaking changes in any release; this
project follows that, but still calls out breaking changes under a
"Changed"/"Removed" heading in `CHANGELOG.md` rather than folding them
silently into "Added".

### Security response and release policy

- Report suspected vulnerabilities via GitHub Security Advisories, per
  `SECURITY.md`.
- Best-effort acknowledgment within 5 business days; this is a solo-
  maintainer project, not a guarantee backed by an SLA.
- Fixes ship as a patch release on the latest minor version. Only the latest
  released minor version is supported — no backport commitment to older
  lines pre-v1.0.
- No commitment to request a CVE; the maintainer may do so via GitHub's CVE
  numbering when warranted.
- `docs/security-model.md`'s "local-only processing" and
  `SECURITY.md`'s "must not add telemetry or network transmission" are
  themselves security invariants: a change that adds network I/O to the
  scanning/cleaning path is treated as a security regression through this
  same disclosure path, not an ordinary bug.

## Consequences

- `README.md`/`README-jp.md`, `docs/cli.md`, `docs/character-policy.md`, and
  `SECURITY.md` each gain a short pointer to the relevant part of this
  decision rather than restating it, so there is one place to update when
  the policy itself changes.
- The release workflow must continue selecting an exact patch-qualified Go
  version. Future release-toolchain bumps must check the Unicode category and
  OS-floor consequences described above before merging, not just run the test
  suite.
- A CLI/JSON breaking change now has a concrete bar (major version bump,
  explicit CHANGELOG callout) instead of being a case-by-case judgment call.
