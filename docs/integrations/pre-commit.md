# pre-commit Integration

The hook receives all staged text filenames in one invocation.

Two hooks are planned:

- `clean-invisible-text-check`: detect only.
- `clean-invisible-text-fix`: clean files, report changes, and fail the run when any file changes.

The hook must not parse file extensions. pre-commit's text-file classification should prevent binary input where possible, while the CLI retains its own validation.
