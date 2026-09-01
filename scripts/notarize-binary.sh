#!/usr/bin/env bash
# Submit a single already-codesigned binary for Apple notarization and wait
# for a result. Used by .github/workflows/release.yml for darwin release
# binaries; see docs/release-process.md for the required secrets.
#
# Required env vars (an App Store Connect API key — see
# docs/release-process.md for how to obtain one):
#   NOTARY_API_KEY_PATH - path to the .p8 private key file
#   NOTARY_KEY_ID        - the key's Key ID
#   NOTARY_ISSUER_ID     - the key's Issuer ID
set -euo pipefail

BINARY="$1"
: "${NOTARY_API_KEY_PATH:?NOTARY_API_KEY_PATH must be set}"
: "${NOTARY_KEY_ID:?NOTARY_KEY_ID must be set}"
: "${NOTARY_ISSUER_ID:?NOTARY_ISSUER_ID must be set}"

WORKDIR=$(mktemp -d)
trap 'rm -rf "$WORKDIR"' EXIT
ZIP_PATH="${WORKDIR}/notarize.zip"

# Standalone binaries can't have a notarization ticket stapled to them
# (stapling only works on .app/.pkg/.dmg); Gatekeeper looks the ticket up
# online from the binary's code signature instead when it first runs with
# the quarantine attribute set.
echo "→ Submitting $(basename "$BINARY") for notarization"
ditto -c -k --keepParent "$BINARY" "$ZIP_PATH"

OUTPUT=$(xcrun notarytool submit "$ZIP_PATH" \
  --key "$NOTARY_API_KEY_PATH" \
  --key-id "$NOTARY_KEY_ID" \
  --issuer "$NOTARY_ISSUER_ID" \
  --wait)
echo "$OUTPUT"

if ! grep -q "status: Accepted" <<<"$OUTPUT"; then
  echo "error: notarization did not succeed for ${BINARY}" >&2
  exit 1
fi
