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
- review required after automatic file modification.

## Out of scope

- homoglyph and mixed-script identifier analysis;
- syntax-aware source-code analysis;
- decoding URLs, HTML entities, base64, or archives;
- malware detection;
- proving that arbitrary text is trustworthy;
- rich-text clipboard formats.
