# Task 029: Passkey repository

**Milestone:** 4, Identity and Authorization  
**Primary predecessor:** 028  
**Expected package area:** auth, authorization, mariadb, transaction, logging

## Purpose

Deliver the smallest reviewable implementation increment for **Passkey repository** while preserving all settled yFramework architecture boundaries.

## Source documents

- `docs/architecture/01_CORE_RUNTIME.md`
- `docs/architecture/02_DATA_EXECUTION.md`
- `docs/architecture/04_IMPLEMENTATION_CODEX.md`
- `.codex/MILESTONES/04.md`
- `.codex/DECISION_AUTHORITY.md`

## Scope

Implement **Passkey repository** and no unrelated feature expansion. Treat `.codex/MILESTONES/04.md` as the compact milestone contract, then verify the authoritative architecture sources listed below before changing persisted or externally observable semantics.

## Behavioral requirements

- Expose narrow domain/subsystem operations rather than generic CRUD.
- Use explicit context and transaction dependencies.
- Do not expose internal numeric database IDs across public boundaries.
- Preserve passkey-only human auth, hash-only opaque secrets, deny-by-default permission semantics, and stable UUID identities.

## Required tests

- Add focused unit/package tests for deterministic behavior introduced by this task.
- Add real MariaDB integration coverage when this task depends on MariaDB locking, transactions, temporal behavior, FK behavior, or query plans.
- Add race coverage for concurrent in-memory/runtime behavior when applicable.
- Add failure/restart coverage when this task changes durable coordination or lifecycle behavior.
- Do not require public internet access.

## Acceptance criteria

- The task-specific behavior is implemented and documented.
- Relevant tests pass, including race/integration tests when applicable.
- No public API, persisted semantic, security rule, lifecycle rule, compatibility promise, or failure-recovery behavior was invented beyond the governing sources.
- The repository remains compiling and suitable as the base for the next task.

## Forbidden architectural changes

- Do not add multi-database support. MariaDB remains the only relational database.
- Do not add a framework-owned application `main` or executable template.
- Do not introduce a service locator or generic dependency bag.
- Do not move application/component domain ownership into framework packages.
- Do not store executable scripts or arbitrary expressions in MariaDB.
- Do not change public, persisted, security, lifecycle, compatibility, or recovery semantics unless an authoritative source already settles the change.

## Completion notes

When completing this task, record:
- files changed,
- tests run,
- notable private implementation choices,
- follow-up tasks discovered,
- confirmation that no unresolved architectural decision was silently made.
