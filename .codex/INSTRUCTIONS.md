# Codex Operating Instructions

1. Read the target task file and every listed authoritative source before editing code.
2. Preserve settled architecture. If sources conflict, newer explicit settled clarification wins, then subsystem architecture, then cross-cutting development plan, then task file, then private implementation assumption.
3. Applications own `main`; yFramework is an importable Go framework. MariaDB is the only relational database.
4. Components are passive registrants with explicit typed dependencies. No service locator, arbitrary dependency bag, independent component worker pool, direct external transport, or generic CRUD framework.
5. Stable names at dynamic boundaries are compatibility-sensitive. Internal numeric IDs stay private across subsystem/component/public boundaries.
6. Use real MariaDB tests for locking, migrations, FKs, transaction, lease, and temporal semantics. Use deterministic clocks/transports and no public internet.
7. Keep changes small. Do not scaffold future milestones with fake production implementations.
8. A task is incomplete until relevant tests and documentation pass.
9. If implementation requires a decision that would alter public API, persisted semantics, security, ownership, lifecycle, compatibility, failure recovery, or externally observable behavior, stop and record it in `STATE/open-decisions.md`.
10. Record completion details in `STATE/completed.md`.
