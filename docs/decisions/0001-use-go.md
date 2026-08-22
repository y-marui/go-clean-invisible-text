# ADR 0001: Implement a new Go CLI

## Status

Accepted on 2026-08-22.

## Context

The tool must run on Windows, macOS, and Raspberry Pi, integrate with pre-commit, and provide a backend for an Alfred Workflow. Users should not need Node.js, Python, or another runtime.

Existing tools were reviewed:

- `unspook` is an MIT-licensed TypeScript library, web app, and Node.js CLI. It overlaps substantially, including NBSP replacement, invisible-character cleaning, scanning, and reports.
- `anti-trojan-source` is a Node.js detector with pre-commit support focused on Trojan Source and confusable Unicode.
- `unicleaner` is a Rust security scanner focused on source repositories.
- `out-of-character` is a JavaScript detector and remover.

`unspook` requires Node.js 18, normalizes CRLF and CR to LF by default, handles zero-width characters as categories rather than using this project's contextual ZWJ/ZWNJ policy, and does not ship a standard `.pre-commit-hooks.yaml` manifest at the reviewed revision.

## Decision

Create an independent Go implementation rather than fork `unspook`.

Go provides standalone binaries for the required operating systems and CPU architectures. The implementation will use a narrowly specified character policy and preserve unrelated text.

## Consequences

- Unicode tables and behavior tests must be maintained independently.
- Relevant prior art must be credited in documentation.
- Changes may be proposed upstream when generally useful, but compatibility with `unspook` is not a goal.
- The Alfred Workflow remains in a separate repository.
