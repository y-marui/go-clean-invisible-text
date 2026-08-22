# CLI Contract

The installed command is planned as `clean-invisible-text`.

Diagnostics go to standard error. Cleaned stream content goes only to standard output. Machine-readable JSON output is planned for CI and Alfred integration.

Multiple filenames must be processed in one invocation so pre-commit starts only one process. A file is written only when its cleaned content differs.
