# Task 253: Coalesce policy

**Milestone:** 8, Workflows and Scheduling  
**Primary predecessor:** 252  
**Expected package area:** workflow, schedule, job, mariadb, transaction

## Purpose

Deliver the smallest reviewable implementation increment for **Coalesce policy** while preserving all settled yFramework architecture boundaries.

## Source documents

- `docs/architecture/02_DATA_EXECUTION.md`
- `docs/architecture/04_IMPLEMENTATION_CODEX.md`
- `.codex/MILESTONES/08.md`
- `.codex/DECISION_AUTHORITY.md`

## Scope

Implement **Coalesce policy** and no unrelated feature expansion. Treat `.codex/MILESTONES/08.md` as the compact milestone contract, then verify the authoritative architecture sources listed below before changing persisted or externally observable semantics.

## Behavioral requirements

- Implement only the behavior implied by this task and its governing milestone contract.
- Read the listed authoritative source files before changing public or persisted behavior.
- Keep the change small, reviewable, compiling, and covered by relevant tests.

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
