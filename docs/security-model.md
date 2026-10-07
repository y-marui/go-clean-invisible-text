# Security Model

## Threats addressed

- Bidirectional control characters that make reviewed text differ from logical order
- Invisible code points used to hide data or alter comparisons
- Accidental NBSP and zero-width characters introduced by copying text
- Unicode tag characters used as hidden payload carriers
- Dangerous control and noncharacter code points

## Safety goals

- deterministic, local-only processing;
- explicit per-code-point findings;
- no silent normalization;
- idempotent cleaning;
- no unrelated file changes;
- review required after automatic file modification;
- an allow-list exception (`--allow`/`--allow-file`,
  [docs/cli.md](cli.md#allow-list-flags)) can only ever affect a
  Warn-classified code point, never a Block-classified one — so even a
  project-committed allow-list file can at most relabel one already-flagged
  Warn occurrence, with its reason kept visible in output, and can never
  suppress a bidi-control, tag-character, or other Block-classified finding.

## Invariants

The tool processes text locally and must not add telemetry or network transmission to scanning or cleaning paths. A change that adds network I/O to the scanning/cleaning path is a security regression, reported through the vulnerability reporting channel described in [SECURITY.md](../SECURITY.md).

The support and response policy (best-effort acknowledgment, patch releases on the latest minor version, supported versions) is defined in [ADR 0002](decisions/0002-v1-compatibility-and-support-policy.md#security-response-and-release-policy).

## Out of scope

- homoglyph and mixed-script identifier analysis;
- syntax-aware source-code analysis;
- decoding URLs, HTML entities, base64, or archives;
- malware detection;
- proving that arbitrary text is trustworthy;
- rich-text clipboard formats.
