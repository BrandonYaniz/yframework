# Task 007: Configuration cross-field validation

**Milestone:** 1, Foundation and Repository Bootstrap  
**Primary predecessor:** 006  
**Expected package area:** errors, validation, logging, config, repository root, .codex

## Purpose

Deliver the smallest reviewable implementation increment for **Configuration cross-field validation** while preserving all settled yFramework architecture boundaries.

## Source documents

- `docs/architecture/04_IMPLEMENTATION_CODEX.md`
- `docs/architecture/00_FOUNDATION_ARCHITECTURE.md`
- `docs/architecture/03_SERVICES_CONFIGURATION.md`
- `docs/architecture/01_CORE_RUNTIME.md`
- `.codex/MILESTONES/01.md`
- `.codex/DECISION_AUTHORITY.md`

## Original settled task detail

Purpose:
- Add static validation for database pools, background timings, server/auth origin relationships, and other settled invariants.

Tests:
- heartbeat < lease
- max idle <= max open
- invalid/conflicting DB connection modes
- required public origin/RP ID validation

Depends on: 006.

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
