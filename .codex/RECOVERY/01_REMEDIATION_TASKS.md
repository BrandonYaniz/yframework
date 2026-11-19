# Recovery Tasks R00 through R12

One R task per iteration. Each task references the existing exact `.codex/TASKS/NNN_*.md` and authoritative architecture. The R-number is an overlay, never a replacement for an original task ID. Status remains **unverified** until reproduced and tested in the current worktree.

## R00: Snapshot and safety audit

**Prerequisites:** none. **Priority:** P0. **Code changes:** none except new evidence files.

- Record commit, branch, dirty state, task count, source tree, and current `git diff` without discarding changes.
- Check whether tracked state or Git history includes secrets or sensitive machine details. Do not output their contents. If leaked, record a security blocker and request credential rotation where appropriate.
- Compare `AGENTS.md`, `.gitignore`, `.codex/INSTRUCTIONS.md`, `.opencode/commands/next.md`, and task selector rules. Confirm `.codex/` is tracked and agent instructions no longer claim the opposite.
- Prepare an isolated MariaDB integration test database strategy. Protect the primary `devframework` database. Do not perform destructive tests until isolation is proven.
- Reproduce known defects on current branch without assuming past source review is still accurate.

**Deliverable:** `.codex/STATE/recovery-verification.md` created with baseline and `unverified` classifications; no repaired source claimed yet.

## R01: Reconcile task and milestone tracking

**Original:** 001 through 026, `scripts/next-task.py`, `.codex/STATE/completed.md`, `.codex/CURRENT_MILESTONE.md`, `.codex/STATE/recovery-audit.md`, `.codex/STATE/recovery-fixes.md`. **Priority:** P0.

- Historical completion log previously changed from 001 through 018 plus 023 to a five-entry file (001, 009, 012, 017, 018), several with implausible sample SHAs and inaccurate titles. Treat both historical versions as evidence needing validation, not truth.
- Preserve historical records with explicit provenance. Do not fabricate commit hashes, test counts, or successful tests. Investigate actual Git history if available.
- Create a canonical **verified-state** registry by original task ID with statuses `not_started`, `unverified`, `incomplete`, `blocked`, `verified`. The proof is tests + artifact review, not a heading in a Markdown file.
- Change `scripts/next-task.py` only after defining deterministic semantics for the verified-state registry and preserving its task-dependency parser. A verified predecessor must be required; a historical heading cannot unblock work.
- If implementation of this selector changes externally observable behavior not already settled, record a decision for review. Add tests that show stale headings **cannot** unblock Task 024 and missing predecessor evidence blocks Task 023.
- Correct progress documents to reflect verified status; do not declare Milestones 1 or 2 completed until gates pass.

**Pass:** `--status` and `--json` both distinguish attempted vs verified, report missing prerequisites, and do not claim completion based solely on historical notes.

## R02: Foundation and logging correctness

**Original:** Tasks 001–005, 008. **Sources:** `.codex/MILESTONES/01.md`, `01_CORE_RUNTIME.md`, `03_SERVICES_CONFIGURATION.md`. **Priority:** P1.

- Verify errors and validation semantics, stable public-safe error codes, wrapping, deterministic field-validation aggregation, and bootstrap stderr logging.
- Replace plaintext production logging in `logging/logging.go` with structured JSON Lines per the source, stable event/field names, level filtering, secret redaction, and correct output lifecycle. Keep bootstrap fallback independent.
- Add tests parsing actual emitted JSON, asserting levels, caller/request metadata where specified, secret redaction, and failure/shutdown behavior. Fix `logging/logging_test.go` which currently only checks constructor existence.
- Check concurrency and races for logger writes; avoid exporting third-party logger types as component interfaces.

**Pass:** explicit Task 005 requirements independently proven; Tasks 001–004 either verified or separately recorded with missing work.

## R03: Strict deployment TOML and cross-field validation

**Original:** Tasks 006–007, 009. **Sources:** `03_SERVICES_CONFIGURATION.md` top-level TOML, configuration authority, and `02_DATA_EXECUTION.md`. **Priority:** P1.

- `config/config.go` currently only reads three environment variables. Implement strictly decoded **typed TOML** and compiled defaults. Reject unknown sections/keys, type mismatches, invalid enums, malformed durations, and unregistered component sections.
- Preserve secret classification/redaction. No generic `map[string]any` primary configuration or blanket environment overrides. CLI overrides remain limited by documented policy.
- Implement static cross-field validation: heartbeat vs lease, foreground/background pool bounds, DB mode, explicit public origin/WebAuthn RP ID relationships, socket and timeout constraints, and adapter settings as settled.
- `mariadb/config.go` must consume typed values, validate DSN/TLS/escaping rules through the driver, and avoid printing password-bearing strings.
- Test default cases, rejected unknown keys/types, secret redaction, mutually invalid settings, and failure before DB mutation.

**Pass:** Milestone 1 TOML contract is fully exercised and no environment-only fallback silently replaces it.

## R04: MariaDB compatibility and transactions

**Original:** Tasks 009–011. **Sources:** `02_DATA_EXECUTION.md` MariaDB section and `04_IMPLEMENTATION_CODEX.md`. **Priority:** P1.

- Implement explicit connection construction from typed config and a separate dedicated startup/migration connection.
- Validate MariaDB version and required capabilities **before** mutation using the FreeBSD-reference per-release support policy; do not hardcode a guessed version or add other SQL dialects.
- Verify transaction commit, rollback, and context cancellation on real data, not only a `SELECT 1`.
- Replace tests that log unexpected DB errors as acceptable and goroutines that end without joining/validating outcomes.

**Pass:** unsupported server/capability detection is nonmutating, supported server passes, and real transaction effects are verified.

## R05: Migration identity, registration, and durable history

**Original:** Tasks 012–013. **Sources:** `02_DATA_EXECUTION.md`, `04_IMPLEMENTATION_CODEX.md`. **Priority:** P0.

- `migrations/registry.go` currently keys records by owner + name and has mutable `AppliedAt` state in registration objects. Align declared stable owner type/name, migration identity, deterministic source checksum, and DB-free registration with architecture.
- Use `framework_migrations` with ownership-qualified unique identity, checksum, applied timestamp, and documented useful timing metadata. The present table has `name` + `version` uniqueness without owner, which risks collisions.
- Encode SHA-256 as **hex**, never raw bytes masquerading as a string. Reject checksum drift of applied identities.
- Define a safe compatibility path if any local DB already has the old `migration_history` or the incomplete `framework_migrations` table. **Do not silently rewrite real history or alter a populated DB.** If not settled, raise a decision.
- Test duplicate registration, checksum identity/drift, owner collision, immutable history, and exact schema constraints.

**Pass:** a component and application can use the same human-readable migration name without history collision, and changed applied content is detected.

## R06: Advisory lock ownership and failure handling

**Original:** Task 014. **Sources:** `02_DATA_EXECUTION.md` dedicated migration connection and locking. **Priority:** P0.

- In `migrations/lock.go`, `Acquire()` currently obtains a `*sql.Conn` then defers `Close()` before returning. `Release()` obtains another connection. Replace this invalid lifetime model with explicit lock-guard ownership of **one held connection for the entire migration critical section**.
- Ensure timeout/cancellation returns no false success, lock release occurs on the owner connection, connection close releases ownership, and unexpected loss is detected; abort migrations on lost ownership.
- Test with two simultaneous connections/runner instances, bounded lock wait, connection loss, and subsequent successful acquisition. Use real MariaDB.
- Do not add a persistent lock table or mix advisory lock ownership with framework logical concurrency locks.

**Pass:** an independent runner cannot enter while lock is held; release/crash correctly allows the next runner.

## R07: Real forward-only migration execution

**Original:** Task 015. **Sources:** `02_DATA_EXECUTION.md` migration implementation contract. **Priority:** P0.

- `migrations/runner.go` currently calls `executeMigration()` which **returns nil without running any SQL**. Implement actual registered migration execution and required post-execution verification.
- Deterministic order: framework, application, components in resolved dependency order, and declared sequential order within each owner.
- Validate compatibility, registration consistency, previously-applied checksums, and expected initial states before mutation.
- Record migration success **only after** execution and verification. A failed migration leaves no success record. Previously successful migrations remain recorded.
- No automatic down migrations, no destructive fake rollback, and no deletion of applied-history records to pretend schema changes were reversed. Remove or explicitly deprecate `Runner.Rollback`, `HistoryStore.MarkMigrationRolledBack`, and CLI rollback in accordance with compatibility/decision authority.
- Handle nontransactional MariaDB DDL and recognize/recover known partial schema states without indiscriminate `IF NOT EXISTS` masking drift.
- Structured logs include stable identity, checksum, duration, result; no sensitive row data.

**Pass:** from empty DB actually creates tables and reruns idempotently; checksum drift fails; failed DDL is unrecorded; partial-DDL recovery matches source contract.

## R08: Runtime foreground and background pools

**Original:** Task 016. **Priority:** P1.

- `mariadb/pools.go` still has not-implemented `GetForegroundPool` and `GetBackgroundPool`. Implement explicit construction through typed dependencies, after migration phase, with separate limits and cleanup if either pool construction fails.
- Avoid untyped globals/service locators. Confirm lifetime and ownership at application shutdown.
- Test independent pool limits, context-aware connectivity/health, partial failure cleanup, and shutdown.

**Pass:** app can create two correctly configured pools and safely dispose of them.

## R09: CLI operations and executable MariaDB gate

**Original:** Tasks 017–018. **Priority:** P0.

- `migrations/cli.go` currently returns `"placeholder_checksum"` from Checksum and success from Validate. Implement registered, non-executable CLI service primitives that match the shared migration engine, including **nonmutating** validate/status, and removing unsafe rollback behavior.
- `migrations/integration_test.go` currently logs placeholder success for empty migration, concurrency, checksum drift, unsupported DB, and migration-runner tests. `migrations/runner_test.go`, `lock_test.go`, and `history_test.go` contain `t.Skip` in essential cases. Replace with meaningful assertions on real MariaDB.
- Required gate tests must **fail** if test DB unavailable. No `t.Skip` and no `t.Log` as proof. Use explicit isolated fixtures and protect non-test DBs.
- Add `scripts/verify-milestone.sh` (or equivalent) with an explicit named-test manifest. The optional included `scripts/check_required_go_tests.py` detects missing/skipped/failed test identities from `go test -json` but cannot establish assertion quality.

**Pass:** Milestone 2 complete acceptance cases physically execute and verify expected DB state, including unsupported-version no mutation and lock crash/reacquire.

## R10: Typed component graph and registry

**Original:** Tasks 019–022. **Priority:** P1.

- `internal/registry/dependency_graph.go` removed the reversed result but still iterates maps without deterministic sorting; sort roots and dependency lists. Test repeated order, chain, diamond, missing dependency, and cycles.
- Verify frozen registry rejects duplicates and mutation after freeze, supports concurrent reads, and protects internal maps from external mutation.
- Replace generic `constructor func() interface{}` and string-based service lookup patterns where they violate the settled typed registrar boundary. Do not invent a broad service container.
- Reconcile existing code for Tasks 019–022 despite missing historical completion entries.

**Pass:** stable deterministic dependency-first order, validated dependencies, typed registrar, frozen state, no untyped global services.

## R11: Application builder and minimal consuming fixture

**Original:** Tasks 023–026. **Priority:** P1.

- `application/builder.go` currently has no-op `ValidateComposition`, `GetComponent` always false, and `Run` stub. Keep lifecycle behavior owned by later milestones; remove misleading fake functionality or implement only the documented composition surface.
- `Build` must validate all static dependencies, duplicates, and registrations *before* any DB mutation, and freeze after successful validation in deterministic order.
- Implement internal consumer fixture using public APIs only. Do not create a framework-owned `main` or a public `examples/` app.
- Fix mismatch where Task 023 was recorded complete while its required predecessors 021 and 022 were not verified.

**Pass:** public consumer composes a real graph; invalid graph fails without DB mutation; registrar freezes predictably.

## R12: Recovery gate, resume, and forward roadmap

**Prerequisites:** R00–R11 verified. **Priority:** P0 gate.

- Run and review Milestone 1, 2, and 3 acceptance gates. Verify required tests ran, had meaningful assertions, used real MariaDB as needed, and were not skipped.
- Correct `.codex/STATE/` to show factual progress and unresolved risks. Link commits and evidence. Update `.codex/CURRENT_MILESTONE.md` only when a milestone gate actually passes.
- Repair default task flow so the next original task is chosen only when verified predecessors and milestones permit it.
- Resume original tasks 027 onward, then Milestones 4–15 in original dependency order. Do not implement beyond a blocked foundation.
- Later, support separate optional `yframework-admin` module with left-side navigation and component contributions; this is **not** a foundational repair task.

**Pass:** the normal task flow is safe to resume and Milestones 1–3 have auditable evidence.
