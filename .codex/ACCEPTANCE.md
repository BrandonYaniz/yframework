# Acceptance

Every task must leave the repository compiling and relevant tests passing. Every milestone gate must satisfy the corresponding acceptance criteria in `docs/architecture/04_IMPLEMENTATION_CODEX.md`.

Before a v1 release tag, all behavioral acceptance tests, race tests, relevant fuzz tests, zero-to-current and supported-upgrade migrations, critical MariaDB EXPLAIN reviews, performance baselines, security review, documentation, and final `.codex` synchronization must pass.

A green build alone is not sufficient when a documented behavioral acceptance criterion is unmet.
