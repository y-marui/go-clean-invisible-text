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
