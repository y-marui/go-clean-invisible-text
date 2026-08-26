# Release Process

`.github/workflows/release.yml` builds and publishes standalone binaries for
every architecture in [Issue #7](https://github.com/y-marui/go-clean-invisible-text/issues/7):

| OS | Architecture | Runner | `go test` |
|---|---|---|---|
| Linux | amd64 | `ubuntu-latest` | native |
| Linux | arm64 | `ubuntu-24.04-arm` | native |
| Linux | arm (`GOARM=7`) | `ubuntu-latest` (cross-built) | QEMU smoke test (`--help`) only — no hosted `linux/arm/v7` runner exists |
| macOS | amd64 | `macos-15-intel` | native |
| macOS | arm64 | `macos-latest` | native |
| Windows | amd64 | `windows-latest` | native |
| Windows | arm64 | `windows-11-arm` | native |

All builds run with `CGO_ENABLED=0`, so cross-compilation needs nothing
beyond `GOOS`/`GOARCH`/`GOARM` — no cross-compiler toolchain.

## Cutting a release

1. Update `CHANGELOG.md`: move `[Unreleased]` entries under a new
   `## [vX.Y.Z] - YYYY-MM-DD` heading.
2. Tag and push: `git tag vX.Y.Z && git push origin vX.Y.Z`.
3. The workflow builds all seven targets, generates `checksums.txt`
   (SHA-256), attests build provenance for every binary via
   [`actions/attest-build-provenance`](https://github.com/actions/attest-build-provenance),
   and publishes a GitHub Release with all of it attached plus
   auto-generated release notes.

`workflow_dispatch` runs the same build/test matrix without publishing
anything (the provenance and release-creation steps only run on an actual
`refs/tags/*` push) — use it to validate the matrix after changing the
workflow, before ever pushing a version tag.

## Verifying a downloaded binary

```bash
sha256sum -c checksums.txt --ignore-missing
gh attestation verify clean-invisible-text-<os>-<arch> --repo y-marui/go-clean-invisible-text
```

## Supported architectures

See [ADR 0002](decisions/0002-v1-compatibility-and-support-policy.md) for
which OS versions and Raspberry Pi architectures these binaries are expected
to run on.
