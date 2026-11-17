# Task 422: LLM schema migration

**Milestone:** 11, LLM and yllmd  
**Primary predecessor:** 421  
**Expected package area:** llm, external, health, logging, mariadb

## Purpose

Deliver the smallest reviewable implementation increment for **LLM schema migration** while preserving all settled yFramework architecture boundaries.

## Source documents

- `docs/architecture/03_SERVICES_CONFIGURATION.md`
- `docs/architecture/02_DATA_EXECUTION.md`
- `docs/architecture/04_IMPLEMENTATION_CODEX.md`
- `.codex/MILESTONES/11.md`
- `.codex/DECISION_AUTHORITY.md`

## Scope

Implement **LLM schema migration** and no unrelated feature expansion. Treat `.codex/MILESTONES/11.md` as the compact milestone contract, then verify the authoritative architecture sources listed below before changing persisted or externally observable semantics.

## Behavioral requirements

- Implement only the schema owned by this subsystem.
- Use framework_ tables only for framework state.
- Keep high-write current rows narrow and separate large content/history where specified.
- Add only indexes justified by documented queries and tests.
- Use the registered migration engine and checksum rules.
- Do not auto-diff or perform destructive schema sync.
- Test zero/current and failure behavior against real MariaDB where database semantics matter.
- Components must use the framework LLM API, never the yllmd protocol directly; prompt/output content stays out of ordinary logs.

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
