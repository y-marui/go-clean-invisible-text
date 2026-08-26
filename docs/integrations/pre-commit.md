# pre-commit Integration

The hook receives all staged text filenames in one invocation.

Two hooks are defined in [`.pre-commit-hooks.yaml`](../../.pre-commit-hooks.yaml):

- `clean-invisible-text-check`: detect only.
- `clean-invisible-text-fix`: clean files, report changes, and fail the run when any file changes.

The hook must not parse file extensions. `types: [text]` relies on pre-commit's
content-based text-file classification to filter out binary input where
possible, while the CLI retains its own validation (`internal/mutate.IsBinary`,
invalid-UTF-8 rejection) as the authoritative check.

## Installation options

### Go source (default)

Reference this repository directly; pre-commit's `language: golang` support
builds `cmd/clean-invisible-text` from source the first time the hook runs,
using whatever Go toolchain is on the machine running pre-commit. No manual
install step.

```yaml
repos:
  - repo: https://github.com/y-marui/go-clean-invisible-text
    rev: vX.Y.Z # pin to a released tag
    hooks:
      - id: clean-invisible-text-check
      - id: clean-invisible-text-fix
```

### Pre-installed binary

For environments that already carry the binary on `PATH` (e.g. `go install
github.com/y-marui/go-clean-invisible-text/cmd/clean-invisible-text@latest`,
or a downloaded release binary — see Issue #7) and want to skip the
from-source build on every run, define local hooks in the consuming
project's own `.pre-commit-config.yaml` instead of referencing this repo:

```yaml
repos:
  - repo: local
    hooks:
      - id: clean-invisible-text-check
        name: clean-invisible-text check
        language: system
        entry: clean-invisible-text check
        types: [text]
      - id: clean-invisible-text-fix
        name: clean-invisible-text fix
        language: system
        entry: clean-invisible-text fix
        types: [text]
```
