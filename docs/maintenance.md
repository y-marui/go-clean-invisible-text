# Maintenance

Prompts for a local LLM (read-only) to draft updates to `docs/` files. The
primary AI must verify every claim against the real files before saving.

## Update docs/architecture.md

```
Survey this Go repository's top-level and internal/ package structure.

Steps:
1. List directories under the repository root and internal/, one level deep.
2. Read go.mod and any existing docs/architecture.md.
3. Produce an updated docs/architecture.md following this format, with every
   claim backed by a file path:

   # Architecture

   ## Overview
   <!-- 3 lines max -->

   ## Entry Points
   - `path/file` — description

   ## Directory Structure
   | Directory | Role |
   |---|---|

   ## Key Dependencies
   | Library / module | Purpose |
   |---|---|

Note: Overview must stay within 3 lines. Do not list every file — that belongs
in file-map.md. List only genuinely key dependencies, not the full go.mod.
```

## Update docs/file-map.md

```
List the Go files under internal/ (and cmd/ if it exists) that were read or
edited in the current task.

Steps:
1. For each such file, note its role in one short phrase and its key
   same-repository dependencies (imported internal packages, not stdlib).
2. Read the existing docs/file-map.md.
3. Append or update entries following this format, grouped by package:

   # File Map

   _Last updated: YYYY-MM-DD_

   ## [package or feature name]
   | File | Role | Key dependencies |
   |---|---|---|
   | `internal/foo/bar.go` | description | `internal/foo/baz.go` |

Note: Do not attempt to cover every file in the repository — only files
actually read or edited. Do not list files that were not explored (avoid
inaccurate entries). Update the "Last updated" date.
```
