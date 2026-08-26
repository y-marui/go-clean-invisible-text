## Reference Order

AI reads the following in order at the start of a task:

1. `README.md` (overview, setup)
2. `DEVELOPING.md` (build, implementation conventions, naming)

Read as needed (any order):
- `CONTRIBUTING.md` (PR/Issue rules)
- `docs/architecture.md` (module/component structure)
- `docs/file-map.md` (file-level dependencies; explore and append if stale or missing)
- `docs/specification.md`, `docs/character-policy.md`, `docs/security-model.md`, `docs/cli.md` (behavior specification, the normative source)
- `docs/ui-design.md` (not applicable — this is a CLI tool)

## Project Overview

A cross-platform Go CLI that detects, explains, and safely cleans dangerous
invisible Unicode characters in UTF-8 plain text. Currently early v0.1: the
character-cleaning engine exists; the CLI commands do not yet.

### Technology Stack

- Go (see `go.mod` for the toolchain version)
- No external runtime dependencies beyond the Go standard library at this stage

### Main Directories

| Path | Role |
|---|---|
| `internal/cleaner/` | The single-pass cleaning engine (code-point tables, `Clean()`) |
| `docs/` | Normative specification, character policy, security model, ADRs |
| `docs/dev-charter/` | Shared dev-charter (`git subtree`, see below) |

## Applied Charter Principles

- Charter reference: use `docs/dev-charter/CHARTER_INDEX.md` to find the relevant topic, then read only that file
- YAGNI, minimal diff scope, reuse existing patterns before adding new ones — `docs/dev-charter/PRINCIPLES.md`
- Secrets and pre-commit security gates — `docs/dev-charter/SECURITY_POLICY.md`
- Public-facing text (README, CLI/error output, commit/PR text) is English; internal Japanese is fine — `docs/dev-charter/LANGUAGE_POLICY.md`
- Do not directly edit files under `docs/dev-charter/`; changes go through an Issue in the dev-charter repository and `git subtree pull`

## Document Sync Rule

When a spec, rule, or structural change happens, update the related documentation
in the same piece of work. This includes files under `docs/` as well as root files
such as `AI_CONTEXT.md` and `README.md`.

## Project-Specific Rules

- The normative source for character handling is `docs/character-policy.md`. A
  change to `internal/cleaner/tables.go` that isn't already covered there requires
  updating that spec (or opening a decision Issue) in the same change — see
  `CONTRIBUTING.md`.
- A change that alters observable CLI/engine behavior must update the relevant
  specification file or add an ADR under `docs/decisions/`. Closed Issues are not
  the source of truth.
- Roadmap and task tracking live in GitHub Issues/Milestones under this repository
  (Issue #1 is the roadmap), not in Markdown files in this repo.

## AI Tool Assignments

- **Tools in use**: Claude Code, Codex, GitHub Copilot, Gemini CLI, local LLM (Ollama)
- **Canonical responsibilities**: `docs/dev-charter/AI_COLLABORATION_RULES.md`, "AI Tool Responsibilities" and "Rules for Multi-AI Usage"
- **Project-specific overrides**: none

## Prohibited Actions

- Adding telemetry or network transmission to the scanning/cleaning code paths (`docs/security-model.md`, `SECURITY.md`) — this tool is local-only
- Committing secrets or credentials
- Direct edits under `docs/dev-charter/`
