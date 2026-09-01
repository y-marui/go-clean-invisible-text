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

The QEMU action and both container images used by the `linux/armv7` smoke test
are pinned to immutable commits or image digests. When updating them, update
the adjacent human-readable version tag and immutable revision together, then
run `workflow_dispatch` before cutting a release. The smoke-test container has
no network access and receives the built binary through a read-only mount so it
cannot alter the artifact that is subsequently checksummed and attested.

## Verifying a downloaded binary

```bash
sha256sum -c checksums.txt --ignore-missing
gh attestation verify clean-invisible-text-<os>-<arch> --repo y-marui/go-clean-invisible-text
```

macOS binaries are also signed and notarised (see below); verify with:

```bash
codesign -dv --verbose=2 clean-invisible-text-darwin-<arch>
spctl -a -t exec -vv clean-invisible-text-darwin-<arch>
```

## macOS code signing and notarization

For a tagged release, the two `darwin` build jobs additionally sign the
binary with a Developer ID Application certificate
(`codesign --options runtime --timestamp`) and submit it to Apple's
notary service (`xcrun notarytool submit --wait`,
[`scripts/notarize-binary.sh`](../scripts/notarize-binary.sh)), before
`checksums.txt`/attestation/release publishing — so the checksummed,
attested, and published bytes are the signed-and-notarised ones. Standalone
binaries can't have a notarization ticket stapled to them, so Gatekeeper
looks the ticket up online from the code signature the first time the binary
runs with the quarantine attribute set (e.g. after download via a browser).
`workflow_dispatch` runs never touch these secrets (see the
`if: startsWith(github.ref, 'refs/tags/')` guards in `release.yml`), so
iterating on the workflow doesn't burn notarization requests.

Consuming this project's own binaries as a third-party dependency (as
[alfred-clean-invisible-text](https://github.com/y-marui/alfred-clean-invisible-text)
does — see [`go-clean-invisible-text#31`](https://github.com/y-marui/go-clean-invisible-text/issues/31))
should verify the signature after checksum/attestation verification, not
instead of it: the checksum/attestation chain proves the bytes came from
this repository's CI; `codesign`/`spctl` prove those bytes carry a valid
Apple signature.

### One-time secret setup (manual, @y-marui)

The workflow expects five repository secrets. None of them should ever be
committed or pasted into an issue/PR/chat — set them directly via
`gh secret set` or the GitHub UI.

**1. Export the Developer ID Application certificate as `.p12`:**

Keychain Access → login keychain → My Certificates → find
"Developer ID Application: Yukihiro Marui (7TEQWKRRX7)" → expand it, select
both the certificate and its private key → right-click → Export 2 items… →
save as `certificate.p12`, choosing an export password (this becomes
`MACOS_CERTIFICATE_PASSWORD` below).

```bash
base64 -i certificate.p12 | gh secret set MACOS_CERTIFICATE_P12 \
  --repo y-marui/go-clean-invisible-text
gh secret set MACOS_CERTIFICATE_PASSWORD \
  --repo y-marui/go-clean-invisible-text  # paste the export password when prompted
rm certificate.p12  # don't leave the exported key sitting on disk
```

Two GUI pitfalls hit during the first real setup (Issue #31): the account has four
similarly-named certificates (3rd Party Mac Developer Installer, **Apple
Distribution**, **Developer ID Application**, Developer ID Installer) — picking
the wrong one imports fine but makes the "Import signing certificate" step fail
with `no Developer ID Application identity found`, since the CI script greps the
imported keychain by that exact name. Double-check the Team ID `(7TEQWKRRX7)`
suffix, not just "Developer ID" in the name. Separately, Keychain Access's "My
Certificates" view sometimes fails to show a certificate as paired with its
private key (export format options stay greyed out even though the key exists)
— `security find-identity -v -p codesigning` is the reliable way to confirm a
given identity actually has an exportable private key, and
`security export -k login.keychain-db -t identities -f pkcs12 -P <password> -o certificate.p12`
is a working CLI fallback when the GUI won't cooperate (it exports every
identity in the keychain at once, which is fine — the import step above filters
by name).

**2. Generate an App Store Connect API key for notarization:**

[appstoreconnect.apple.com](https://appstoreconnect.apple.com/) → Users and
Access → Integrations → App Store Connect API → generate a key with the
"Developer" role. Apple lets you download the `.p8` private key file only
once — save it somewhere safe until it's registered below, then delete it.
Note the Key ID and Issuer ID shown on the same page.

```bash
base64 -i AuthKey_XXXXXXXXXX.p8 | gh secret set NOTARY_API_KEY \
  --repo y-marui/go-clean-invisible-text
gh secret set NOTARY_KEY_ID --repo y-marui/go-clean-invisible-text     # the Key ID
gh secret set NOTARY_ISSUER_ID --repo y-marui/go-clean-invisible-text  # the Issuer ID
rm AuthKey_XXXXXXXXXX.p8
```

Once all five secrets (`MACOS_CERTIFICATE_P12`, `MACOS_CERTIFICATE_PASSWORD`,
`NOTARY_API_KEY`, `NOTARY_KEY_ID`, `NOTARY_ISSUER_ID`) exist, the next tagged
release automatically signs and notarises both `darwin` binaries — no
further workflow changes needed. The same Developer ID certificate and App
Store Connect key already registered on
[alfred-clean-invisible-text](https://github.com/y-marui/alfred-clean-invisible-text)
can be reused here, but GitHub Actions secrets are scoped per repository, so
they must be registered again against this repository.

## Supported architectures

See [ADR 0002](decisions/0002-v1-compatibility-and-support-policy.md) for
which OS versions and Raspberry Pi architectures these binaries are expected
to run on.
